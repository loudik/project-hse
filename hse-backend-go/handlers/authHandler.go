package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"hse-backend-go/db"
	"hse-backend-go/models"
	"hse-backend-go/utils"
)

// emailVerifyKey builds the Redis key holding a pending email-verification
// token. Value = userID, TTL = how long the link stays valid.
func emailVerifyKey(token string) string {
	return "email_verify:" + token
}

// POST /api/sign-in (endpointConfig.signIn on the Ecme frontend)
func SignIn(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email and password are required"})
		return
	}

	row := db.DB.QueryRow(`
		SELECT u.id, u.name, u.email, u.password_hash, u.role_id, u.status, r.name
		FROM users u JOIN roles r ON r.id = u.role_id
		WHERE u.email = ?
	`, input.Email)

	var (
		id, name, email, roleName, status string
		passwordHash                      sql.NullString
		roleID                            int
	)
	err := row.Scan(&id, &name, &email, &passwordHash, &roleID, &status, &roleName)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to process login"})
		return
	}

	if !passwordHash.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "This account uses social sign-in. Please use the Microsoft/Google button instead."})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(passwordHash.String), []byte(input.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		return
	}

	switch status {
	case "Pending":
		c.JSON(http.StatusForbidden, gin.H{"message": "Please verify your email before signing in. Check your inbox for the verification link."})
		return
	case "Suspended":
		c.JSON(http.StatusForbidden, gin.H{"message": "This account has been suspended. Please contact an administrator."})
		return
	}

	token, err := utils.GenerateToken(id, roleID, roleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		User: models.EcmeUser{
			UserName:  name,
			Email:     email,
			Avatar:    "",
			Authority: []string{roleName},
		},
	})
}

func SignOut(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if strings.HasPrefix(header, "Bearer ") {
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		if claims, err := utils.ParseToken(tokenStr); err == nil && claims.ExpiresAt != nil {
			_ = utils.BlacklistToken(claims.ID, claims.ExpiresAt.Time)
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// GET /api/auth/me
func Me(c *gin.Context) {
	userID := c.GetString("userId")

	row := db.DB.QueryRow(`
		SELECT u.name, u.email, r.name
		FROM users u JOIN roles r ON r.id = u.role_id
		WHERE u.id = ?
	`, userID)

	var name, email, roleName string
	if err := row.Scan(&name, &email, &roleName); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, models.EcmeUser{UserName: name, Email: email, Authority: []string{roleName}})
}

// POST /api/sign-up - public registration for external vessel operators
func SignUp(c *gin.Context) {
	var input struct {
		Name             string `json:"name" binding:"required"`
		Email            string `json:"email" binding:"required,email"`
		Password         string `json:"password" binding:"required,min=6"`
		PhoneNumber      string `json:"phoneNumber"`
		Position         string `json:"position"`
		OrganizationName string `json:"organizationName"`
		Location         string `json:"location"`
		Website          string `json:"website"`
		AcceptTerms      bool   `json:"acceptTerms" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incomplete or invalid data: " + err.Error()})
		return
	}
	if !input.AcceptTerms {
		c.JSON(http.StatusBadRequest, gin.H{"message": "You must accept the Terms and Conditions"})
		return
	}

	var existingID string
	err := db.DB.QueryRow(`SELECT id FROM users WHERE email = ?`, input.Email).Scan(&existingID)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "This email is already registered"})
		return
	}
	if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to check email"})
		return
	}

	var operatorRoleID int
	if err := db.DB.QueryRow(`SELECT id FROM roles WHERE name = 'Operator'`).Scan(&operatorRoleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Default role not found"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to process password"})
		return
	}

	newID := uuid.NewString()
	_, err = db.DB.Exec(`
		INSERT INTO users
			(id, name, email, password_hash, role_id, status,
			 phone_number, position, organization_name, location, website, accepted_terms_at)
		VALUES (?, ?, ?, ?, ?, 'Pending', ?, ?, ?, ?, ?, NOW())
	`, newID, input.Name, input.Email, string(hash), operatorRoleID,
		nullableString(input.PhoneNumber), nullableString(input.Position), nullableString(input.OrganizationName),
		nullableString(input.Location), nullableString(input.Website),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create account: " + err.Error()})
		return
	}

	// Create a verification token, valid 24h - stored in Redis so it
	// expires itself, no manual expiry/used-flag bookkeeping needed
	token := uuid.NewString()
	if err := db.RDB.Set(context.Background(), emailVerifyKey(token), newID, 24*time.Hour).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Account created but failed to prepare email verification"})
		return
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	verificationURL := frontendURL + "/verify-email?token=" + token

	emailBody := fmt.Sprintf(`
		<p>Hi %s,</p>
		<p>Thanks for registering. Please verify your email to activate your account:</p>
		<p><a href="%s">%s</a></p>
		<p>This link expires in 24 hours.</p>
	`, input.Name, verificationURL, verificationURL)

	emailErr := utils.SendEmail(input.Email, "Verify your HSE Approval account", emailBody)
	if emailErr != nil {
		log.Printf("[sign-up] failed to send verification email to %s: %v", input.Email, emailErr)
	}
	response := gin.H{
		"message": "Registration successful. Please check your email to verify your account.",
	}
	if os.Getenv("SMTP_HOST") == "" || emailErr != nil {
		response["verificationUrl"] = verificationURL
	}

	c.JSON(http.StatusCreated, response)
}

// GET /api/verify-email?token=... - activates a Pending account
func VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Verification token is required"})
		return
	}

	ctx := context.Background()
	key := emailVerifyKey(token)

	userID, err := db.RDB.Get(ctx, key).Result()
	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "This verification link is invalid or has expired"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to verify email"})
		return
	}

	if _, err := db.DB.Exec(`UPDATE users SET status = 'Active' WHERE id = ?`, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to activate account"})
		return
	}

	// one-time use - delete right after successful verification
	db.RDB.Del(ctx, key)

	c.JSON(http.StatusOK, gin.H{"message": "Your email has been verified. You can now sign in."})
}
