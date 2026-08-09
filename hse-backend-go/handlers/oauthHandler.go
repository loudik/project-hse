package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hse-backend-go/db"
	"hse-backend-go/models"
	"hse-backend-go/utils"
)

type msTokenResponse struct {
	AccessToken string `json:"access_token"`
	Error       string `json:"error"`
	ErrorDesc   string `json:"error_description"`
}

type msProfile struct {
	DisplayName       string `json:"displayName"`
	Mail              string `json:"mail"`
	UserPrincipalName string `json:"userPrincipalName"`
}

// POST /api/oauth/microsoft  { "code": "..." }
func MicrosoftOAuthSignIn(c *gin.Context) {
	var body struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Authorization code is required"})
		return
	}

	accessToken, err := exchangeMicrosoftCode(body.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to exchange Microsoft code: " + err.Error()})
		return
	}

	profile, err := fetchMicrosoftProfile(accessToken)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "Failed to fetch profile from Microsoft: " + err.Error()})
		return
	}

	email := profile.Mail
	if email == "" {
		email = profile.UserPrincipalName
	}
	if email == "" {
		c.JSON(http.StatusBadGateway, gin.H{"message": "This Microsoft account has no usable email"})
		return
	}

	user, roleName, err := findOrCreateOAuthUser(email, profile.DisplayName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to process user account: " + err.Error()})
		return
	}
	if user.Status == "Suspended" {
		c.JSON(http.StatusForbidden, gin.H{"message": "This account has been suspended. Please contact an administrator."})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.RoleID, roleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		User: models.EcmeUser{
			UserName:  user.Name,
			Email:     user.Email,
			Avatar:    "",
			Authority: []string{roleName},
		},
	})
}

func exchangeMicrosoftCode(code string) (string, error) {
	tenantID := os.Getenv("MS_TENANT_ID")
	clientID := os.Getenv("MS_CLIENT_ID")
	clientSecret := os.Getenv("MS_CLIENT_SECRET")
	redirectURI := os.Getenv("MS_REDIRECT_URI")

	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("code", code)
	form.Set("redirect_uri", redirectURI)
	form.Set("grant_type", "authorization_code")
	form.Set("scope", "openid profile email User.Read")

	req, err := http.NewRequest(
		http.MethodPost,
		"https://login.microsoftonline.com/"+tenantID+"/oauth2/v2.0/token",
		bytes.NewBufferString(form.Encode()),
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[oauth] request to Microsoft token endpoint failed: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	raw, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		log.Printf("[oauth] failed reading Microsoft response body: %v", readErr)
		return "", readErr
	}
	log.Printf("[oauth] Microsoft token endpoint responded with status %d, body length %d", resp.StatusCode, len(raw))
	if len(raw) == 0 {
		return "", errString("Microsoft returned an empty response (HTTP " +
			strconv.Itoa(resp.StatusCode) + ") - check outbound network access from the backend container")
	}

	var parsed msTokenResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		log.Printf("[oauth] failed to parse Microsoft response as JSON. Status: %d, Body: %s", resp.StatusCode, string(raw))
		return "", err
	}
	if parsed.Error != "" {
		log.Printf("[oauth] Microsoft rejected the request: %s - %s", parsed.Error, parsed.ErrorDesc)
		return "", errString(parsed.ErrorDesc)
	}
	return parsed.AccessToken, nil
}

func fetchMicrosoftProfile(accessToken string) (*msProfile, error) {
	req, err := http.NewRequest(http.MethodGet, "https://graph.microsoft.com/v1.0/me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var profile msProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

func findOrCreateOAuthUser(email, displayName string) (*models.User, string, error) {
	row := db.DB.QueryRow(`
		SELECT u.id, u.name, u.email, u.role_id, u.status, r.name
		FROM users u JOIN roles r ON r.id = u.role_id
		WHERE u.email = ?
	`, email)

	var u models.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.RoleID, &u.Status, &u.RoleName)
	if err == nil {
		return &u, u.RoleName, nil
	}
	if err != sql.ErrNoRows {
		return nil, "", err
	}

	// New internal (Microsoft) sign-in - gets ZERO menu access until an
	// Admin manually assigns their real role (Staff / HSE Officer /
	// ANP HSE / Admin). They're immediately "Active" though (Microsoft
	// already verified their identity, unlike Sign Up's email verification).
	var unassignedRoleID int
	if err := db.DB.QueryRow(`SELECT id FROM roles WHERE name = 'Unassigned'`).Scan(&unassignedRoleID); err != nil {
		return nil, "", err
	}

	newID := uuid.NewString()
	name := displayName
	if name == "" {
		name = email
	}

	_, err = db.DB.Exec(`
		INSERT INTO users (id, name, email, password_hash, role_id, status)
		VALUES (?, ?, ?, NULL, ?, 'Active')
	`, newID, name, email, unassignedRoleID)
	if err != nil {
		return nil, "", err
	}

	return &models.User{ID: newID, Name: name, Email: email, RoleID: unassignedRoleID, RoleName: "Unassigned", Status: "Active"}, "Unassigned", nil
}

func errString(s string) error {
	if s == "" {
		s = "unknown error from Microsoft"
	}
	return &simpleError{s}
}

type simpleError struct{ msg string }

func (e *simpleError) Error() string { return e.msg }
