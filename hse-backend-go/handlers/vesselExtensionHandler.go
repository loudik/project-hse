package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hse-backend-go/db"
	"hse-backend-go/utils"
)

type extensionRequest struct {
	ID                   string  `json:"id"`
	ApplicationID        string  `json:"applicationId"`
	ApplicationNumber    string  `json:"applicationNumber"`
	VesselName           string  `json:"vesselName"`
	CurrentEntryDateTo   string  `json:"currentEntryDateTo"`
	RequestedEntryDateTo string  `json:"requestedEntryDateTo"`
	Reason               string  `json:"reason"`
	Status               string  `json:"status"`
	RejectionReason      *string `json:"rejectionReason"`
	RequestedByName      string  `json:"requestedByName"`
	DecidedByName        *string `json:"decidedByName"`
	DecidedAt            *string `json:"decidedAt"`
	CreatedAt            string  `json:"createdAt"`
}

// POST /api/vessel-applications/:id/extension - owner only, app must be Approved
// Body: { "requestedEntryDateTo": "YYYY-MM-DD", "reason": "..." }
func RequestExtension(c *gin.Context) {
	appID := c.Param("id")
	userID := c.GetString("userId")

	app, err := fetchVesselApplication(appID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Application not found"})
		return
	}
	if app.CreatedBy != userID {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to this application"})
		return
	}
	if app.Status != "Approved" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "You can only request an extension for an Approved application"})
		return
	}

	var input struct {
		RequestedEntryDateTo string `json:"requestedEntryDateTo" binding:"required"`
		Reason               string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "requestedEntryDateTo and reason are required"})
		return
	}

	newDate, err := time.Parse("2006-01-02", input.RequestedEntryDateTo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid date format"})
		return
	}
	if app.EntryDateTo != nil {
		currentDate, _ := time.Parse("2006-01-02", *app.EntryDateTo)
		if !newDate.After(currentDate) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "The requested date must be after the current entry date"})
			return
		}
	}

	id := uuid.NewString()
	_, err = db.DB.Exec(`
		INSERT INTO vessel_extension_requests
			(id, application_id, current_entry_date_to, requested_entry_date_to, reason, status, requested_by)
		VALUES (?, ?, ?, ?, ?, 'Pending', ?)
	`, id, appID, app.EntryDateTo, input.RequestedEntryDateTo, strings.TrimSpace(input.Reason), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to submit extension request"})
		return
	}

	vesselName := "a vessel"
	if app.VesselName != nil {
		vesselName = *app.VesselName
	}
	utils.NotifyReviewers(
		"Entry authorisation extension requested",
		app.ApplicationNumber+" - "+vesselName+" - new date: "+input.RequestedEntryDateTo,
		"/vessel/extension-requests",
	)
	if app.AssignedHSEOfficerID != nil && *app.AssignedHSEOfficerID != "" {
		utils.CreateNotification(
			*app.AssignedHSEOfficerID,
			"Extension requested for an application you're assigned to",
			app.ApplicationNumber+" - new date: "+input.RequestedEntryDateTo,
			"/vessel/review/"+appID,
		)
	}
	utils.LogAudit("vessel_application", appID, "extension_requested", "", "", userID, "Requested new date: "+input.RequestedEntryDateTo)

	c.JSON(http.StatusOK, gin.H{"message": "Extension request submitted"})
}

// GET /api/vessel-applications/extension-requests?status=Pending - Admin/ANP HSE only
func ListExtensionRequests(c *gin.Context) {
	status := c.DefaultQuery("status", "Pending")

	rows, err := db.DB.Query(`
		SELECT r.id, r.application_id, a.application_number, a.vessel_name,
		       r.current_entry_date_to, r.requested_entry_date_to, r.reason,
		       r.status, r.rejection_reason, u.name, r.decided_at, r.created_at,
		       d.name
		FROM vessel_extension_requests r
		JOIN vessel_applications a ON a.id = r.application_id
		JOIN users u ON u.id = r.requested_by
		LEFT JOIN users d ON d.id = r.decided_by
		WHERE r.status = ?
		ORDER BY r.created_at DESC
	`, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch extension requests"})
		return
	}
	defer rows.Close()

	list := []extensionRequest{}
	for rows.Next() {
		var e extensionRequest
		var vesselName sql.NullString
		var rejectionReason sql.NullString
		var decidedByName sql.NullString
		var decidedAt sql.NullTime
		var createdAt time.Time
		var currentDate, requestedDate time.Time

		if err := rows.Scan(&e.ID, &e.ApplicationID, &e.ApplicationNumber, &vesselName,
			&currentDate, &requestedDate, &e.Reason, &e.Status, &rejectionReason,
			&e.RequestedByName, &decidedAt, &createdAt, &decidedByName); err != nil {
			continue
		}
		if vesselName.Valid {
			e.VesselName = vesselName.String
		}
		e.CurrentEntryDateTo = currentDate.Format("2006-01-02")
		e.RequestedEntryDateTo = requestedDate.Format("2006-01-02")
		if rejectionReason.Valid {
			e.RejectionReason = &rejectionReason.String
		}
		if decidedByName.Valid {
			e.DecidedByName = &decidedByName.String
		}
		if decidedAt.Valid {
			s := decidedAt.Time.Format(time.RFC3339)
			e.DecidedAt = &s
		}
		e.CreatedAt = createdAt.Format(time.RFC3339)
		list = append(list, e)
	}

	c.JSON(http.StatusOK, list)
}

// GET /api/vessel-applications/:id/extension-requests - owner or reviewer
func ListExtensionRequestsForApplication(c *gin.Context) {
	appID := c.Param("id")
	userID := c.GetString("userId")
	roleName := c.GetString("roleName")

	app, err := fetchVesselApplication(appID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Application not found"})
		return
	}
	isOwner := app.CreatedBy == userID
	isReviewer := roleName == "Admin" || roleName == "ANP HSE"
	if !isOwner && !isReviewer {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to this application"})
		return
	}

	rows, err := db.DB.Query(`
		SELECT r.id, r.application_id, a.application_number, a.vessel_name,
		       r.current_entry_date_to, r.requested_entry_date_to, r.reason,
		       r.status, r.rejection_reason, u.name, r.decided_at, r.created_at,
		       COALESCE(d.name, '')
		FROM vessel_extension_requests r
		JOIN vessel_applications a ON a.id = r.application_id
		JOIN users u ON u.id = r.requested_by
		LEFT JOIN users d ON d.id = r.decided_by
		WHERE r.application_id = ?
		ORDER BY r.created_at DESC
	`, appID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch extension requests"})
		return
	}
	defer rows.Close()

	list := []extensionRequest{}
	for rows.Next() {
		var e extensionRequest
		var vesselName sql.NullString
		var rejectionReason sql.NullString
		var decidedByName string
		var decidedAt sql.NullTime
		var createdAt time.Time
		var currentDate, requestedDate time.Time

		if err := rows.Scan(&e.ID, &e.ApplicationID, &e.ApplicationNumber, &vesselName,
			&currentDate, &requestedDate, &e.Reason, &e.Status, &rejectionReason,
			&e.RequestedByName, &decidedAt, &createdAt, &decidedByName); err != nil {
			continue
		}
		if vesselName.Valid {
			e.VesselName = vesselName.String
		}
		e.CurrentEntryDateTo = currentDate.Format("2006-01-02")
		e.RequestedEntryDateTo = requestedDate.Format("2006-01-02")
		if rejectionReason.Valid {
			e.RejectionReason = &rejectionReason.String
		}
		if decidedByName != "" {
			e.DecidedByName = &decidedByName
		}
		if decidedAt.Valid {
			s := decidedAt.Time.Format(time.RFC3339)
			e.DecidedAt = &s
		}
		e.CreatedAt = createdAt.Format(time.RFC3339)
		list = append(list, e)
	}

	c.JSON(http.StatusOK, list)
}

// PATCH /api/vessel-applications/extension-requests/:id/decision - Admin/ANP HSE only
// Body: { "status": "Approved" | "Rejected", "rejectionReason": "..." (required if Rejected) }
func DecideExtension(c *gin.Context) {
	reqID := c.Param("id")
	deciderID := c.GetString("userId")

	var input struct {
		Status          string `json:"status" binding:"required"`
		RejectionReason string `json:"rejectionReason"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || (input.Status != "Approved" && input.Status != "Rejected") {
		c.JSON(http.StatusBadRequest, gin.H{"message": "status must be Approved or Rejected"})
		return
	}
	if input.Status == "Rejected" && strings.TrimSpace(input.RejectionReason) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "rejectionReason is required when rejecting"})
		return
	}

	var applicationID, requestedEntryDateTo, currentStatus string
	err := db.DB.QueryRow(`
		SELECT application_id, requested_entry_date_to, status FROM vessel_extension_requests WHERE id = ?
	`, reqID).Scan(&applicationID, &requestedEntryDateTo, &currentStatus)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Extension request not found"})
		return
	}
	if currentStatus != "Pending" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "This request has already been decided"})
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to start transaction"})
		return
	}

	if _, err := tx.Exec(`
		UPDATE vessel_extension_requests
		SET status = ?, rejection_reason = ?, decided_by = ?, decided_at = NOW()
		WHERE id = ?
		`, input.Status, nullableStringPtr(input.RejectionReason), deciderID, reqID); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update request"})
		return
	}

	if input.Status == "Approved" {
		if _, err := tx.Exec(`
			UPDATE vessel_applications
			SET entry_date_to = ?, authorisation_reminder_sent_at = NULL
			WHERE id = ?
		`, requestedEntryDateTo, applicationID); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update application"})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save decision"})
		return
	}

	app, _ := fetchVesselApplication(applicationID)
	if app != nil {
		if input.Status == "Approved" {
			utils.CreateNotification(
				app.CreatedBy,
				"Your extension request was approved",
				app.ApplicationNumber+" - new entry date: "+requestedEntryDateTo,
				"/vessel/"+applicationID,
			)
		} else {
			utils.CreateNotification(
				app.CreatedBy,
				"Your extension request was rejected",
				input.RejectionReason,
				"/vessel/"+applicationID,
			)
		}
		utils.LogAudit("vessel_application", applicationID, "extension_"+strings.ToLower(input.Status), "", "", deciderID, "")
	}

	c.JSON(http.StatusOK, gin.H{"message": "Decision saved"})
}
