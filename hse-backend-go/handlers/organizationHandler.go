package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hse-backend-go/db"
	"hse-backend-go/models"
)

// POST /api/organizations - applicant creates their organization profile
func CreateOrganization(c *gin.Context) {
	userID := c.GetString("userId")

	var input models.CreateOrganizationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incomplete or invalid data: " + err.Error()})
		return
	}

	// One organization per applicant for now - if they already have one, block
	var existingOrgID sql.NullString
	_ = db.DB.QueryRow(`SELECT organization_id FROM users WHERE id = ?`, userID).Scan(&existingOrgID)
	if existingOrgID.Valid {
		c.JSON(http.StatusConflict, gin.H{"message": "You already have an organization profile"})
		return
	}

	var dupeID string
	err := db.DB.QueryRow(`SELECT id FROM organizations WHERE registration_number = ?`, input.RegistrationNumber).Scan(&dupeID)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "This registration number is already registered"})
		return
	}
	if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to check for duplicate organization"})
		return
	}

	id := uuid.NewString()
	_, err = db.DB.Exec(`
		INSERT INTO organizations
			(id, name, registration_number, type, address, country, phone_number, email, website, status, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 'Pending', ?)
	`, id, input.Name, input.RegistrationNumber, input.Type,
		nullableString(input.Address), nullableString(input.Country), nullableString(input.PhoneNumber),
		nullableString(input.Email), nullableString(input.Website), userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create organization: " + err.Error()})
		return
	}

	if _, err := db.DB.Exec(`UPDATE users SET organization_id = ? WHERE id = ?`, id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Organization created but failed to link to user"})
		return
	}

	org, err := fetchOrganization(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Organization saved but could not be read back"})
		return
	}
	c.JSON(http.StatusCreated, org)
}

// GET /api/organizations/mine - current user's own organization (if any)
func GetMyOrganization(c *gin.Context) {
	userID := c.GetString("userId")

	var orgID sql.NullString
	err := db.DB.QueryRow(`SELECT organization_id FROM users WHERE id = ?`, userID).Scan(&orgID)
	if err != nil || !orgID.Valid {
		c.JSON(http.StatusNotFound, gin.H{"message": "No organization profile yet"})
		return
	}

	org, err := fetchOrganization(orgID.String)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch organization"})
		return
	}
	c.JSON(http.StatusOK, org)
}

// GET /api/organizations - list all (for ANP HSE review queue)
func ListOrganizations(c *gin.Context) {
	statusFilter := c.Query("status")

	query := `
		SELECT id, name, registration_number, type, address, country, phone_number, email, website,
		       status, created_by, approved_by, approved_at, rejection_reason, created_at
		FROM organizations
	`
	args := []any{}
	if statusFilter != "" {
		query += ` WHERE status = ?`
		args = append(args, statusFilter)
	}
	query += ` ORDER BY created_at DESC`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch organizations"})
		return
	}
	defer rows.Close()

	list := []models.Organization{}
	for rows.Next() {
		org, err := scanOrganization(rows)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read organization data"})
			return
		}
		list = append(list, *org)
	}
	c.JSON(http.StatusOK, list)
}

// GET /api/organizations/:id
func GetOrganization(c *gin.Context) {
	org, err := fetchOrganization(c.Param("id"))
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Organization not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch organization"})
		return
	}
	c.JSON(http.StatusOK, org)
}

// PATCH /api/organizations/:id/decision - ANP HSE approve/reject
func DecideOrganization(c *gin.Context) {
	approverID := c.GetString("userId")
	orgID := c.Param("id")

	var input models.OrganizationDecisionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid data: " + err.Error()})
		return
	}
	if input.Status == "Rejected" && input.RejectionReason == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "A rejection reason is required"})
		return
	}

	if _, err := fetchOrganization(orgID); err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Organization not found"})
		return
	}

	_, err := db.DB.Exec(`
		UPDATE organizations
		SET status = ?, approved_by = ?, approved_at = NOW(), rejection_reason = ?
		WHERE id = ?
	`, input.Status, approverID, nullableString(input.RejectionReason), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save decision"})
		return
	}

	// TODO: send email notification to the applicant once an email/SMTP
	// service is configured (US 2.x mentions "system sends a confirmation").

	org, err := fetchOrganization(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read updated organization"})
		return
	}
	c.JSON(http.StatusOK, org)
}

// --- helpers ---

func fetchOrganization(id string) (*models.Organization, error) {
	row := db.DB.QueryRow(`
		SELECT id, name, registration_number, type, address, country, phone_number, email, website,
		       status, created_by, approved_by, approved_at, rejection_reason, created_at
		FROM organizations WHERE id = ?
	`, id)
	return scanOrganization(row)
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanOrganization(row rowScanner) (*models.Organization, error) {
	var o models.Organization
	var address, country, phone, email, website, approvedBy, rejectionReason sql.NullString
	var approvedAt sql.NullTime

	err := row.Scan(
		&o.ID, &o.Name, &o.RegistrationNumber, &o.Type,
		&address, &country, &phone, &email, &website,
		&o.Status, &o.CreatedBy, &approvedBy, &approvedAt, &rejectionReason, &o.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if address.Valid {
		o.Address = &address.String
	}
	if country.Valid {
		o.Country = &country.String
	}
	if phone.Valid {
		o.PhoneNumber = &phone.String
	}
	if email.Valid {
		o.Email = &email.String
	}
	if website.Valid {
		o.Website = &website.String
	}
	if approvedBy.Valid {
		o.ApprovedBy = &approvedBy.String
	}
	if approvedAt.Valid {
		o.ApprovedAt = &approvedAt.Time
	}
	if rejectionReason.Valid {
		o.RejectionReason = &rejectionReason.String
	}

	return &o, nil
}
