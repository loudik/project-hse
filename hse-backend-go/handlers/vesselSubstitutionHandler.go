package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hse-backend-go/db"
)

// POST /api/vessel-applications/substitute - Operator only
// Body: { "originalApplicationNumber": "VEA-..." }
// Creates a new Draft application linked to the original (which must be
// Approved and belong to the same operator), for the substitute vessel's
// details to be filled in via the normal wizard. It goes through the full
// review flow when submitted - a different vessel needs its own document
// check from ANP HSE.
func SubstituteVessel(c *gin.Context) {
	userID := c.GetString("userId")

	var input struct {
		OriginalApplicationNumber string `json:"originalApplicationNumber" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "originalApplicationNumber is required"})
		return
	}

	var originalID, organizationID, status, createdBy string
	err := db.DB.QueryRow(`
		SELECT id, organization_id, status, created_by
		FROM vessel_applications
		WHERE application_number = ?
	`, strings.TrimSpace(input.OriginalApplicationNumber)).Scan(&originalID, &organizationID, &status, &createdBy)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "No application found with that number"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to look up application"})
		return
	}
	if createdBy != userID {
		c.JSON(http.StatusForbidden, gin.H{"message": "That application does not belong to you"})
		return
	}
	if status != "Approved" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Substitute Vessel is only available for an Approved application"})
		return
	}

	newID := uuid.NewString()
	newApplicationNumber := genApplicationNumber()

	if _, err := db.DB.Exec(`
		INSERT INTO vessel_applications (id, application_number, organization_id, status, created_by, substituted_from_application_id)
		VALUES (?, ?, ?, 'Draft', ?, ?)
	`, newID, newApplicationNumber, organizationID, userID, originalID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create substitute application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":                newID,
		"applicationNumber": newApplicationNumber,
	})
}
