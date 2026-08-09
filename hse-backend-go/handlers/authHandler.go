package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"hse-backend-go/db"
	"hse-backend-go/models"
	"hse-backend-go/utils"
)

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
		// account created via OAuth (Microsoft/Google) - no password set
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

// POST /api/sign-out - called by AuthProvider.signOut(). No server-side
// token invalidation needed (JWT is stateless) - just respond 200 so the
// frontend proceeds to clear its local token/session.
func SignOut(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// GET /api/auth/me - optional, verify token & fetch current user data
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

// POST /api/sign-up - public registration for external vessel
// operators/applicants (US 1.1). Account starts as "Pending" and only
// becomes usable after email verification - see VerifyEmail below.
// Always assigned the "Operator" role (internal ANP staff never
// self-register - see oauthHandler.go for their Microsoft sign-in flow).
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

	// Public Sign Up is for external vessel operators/applicants -
	// they always get the "Operator" role (internal ANP staff never
	// self-register; they sign in via Microsoft - see oauthHandler.go).
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

	// Create a verification token, valid 24h
	verificationID := uuid.NewString()
	token := uuid.NewString()
	_, err = db.DB.Exec(`
		INSERT INTO email_verifications (id, user_id, token, expires_at)
		VALUES (?, ?, ?, ?)
	`, verificationID, newID, token, time.Now().Add(24*time.Hour))
	if err != nil {
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

	response := gin.H{
		"message": "Registration successful. Please check your email to verify your account.",
	}
	// SMTP not configured (local dev) or sending failed - surface the link
	// directly in the response so the flow can still be tested end-to-end.
	if os.Getenv("SMTP_HOST") == "" || emailErr != nil {
		response["verificationUrl"] = verificationURL
	}

	c.JSON(http.StatusCreated, response)
}

// GET /api/verify-email?token=... - activates a Pending account (US 1.1)
func VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Verification token is required"})
		return
	}

	var (
		id, userID string
		expiresAt  time.Time
		usedAt     sql.NullTime
	)
	err := db.DB.QueryRow(`
		SELECT id, user_id, expires_at, used_at FROM email_verifications WHERE token = ?
	`, token).Scan(&id, &userID, &expiresAt, &usedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid verification link"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to verify email"})
		return
	}
	if usedAt.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"message": "This verification link has already been used"})
		return
	}
	if time.Now().After(expiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "This verification link has expired"})
		return
	}

	if _, err := db.DB.Exec(`UPDATE users SET status = 'Active' WHERE id = ?`, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to activate account"})
		return
	}
	if _, err := db.DB.Exec(`UPDATE email_verifications SET used_at = NOW() WHERE id = ?`, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Account activated, but failed to record verification usage"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Your email has been verified. You can now sign in."})
}
