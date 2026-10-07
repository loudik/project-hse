package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hse-backend-go/db"
	"hse-backend-go/models"
	"hse-backend-go/utils"
)

func canCollaborateOnApplication(app *models.VesselApplication, userID, roleName string) bool {
	if roleName == "Admin" || roleName == "ANP HSE" {
		return true
	}
	if app.AssignedHSEOfficerID != nil && *app.AssignedHSEOfficerID == userID {
		return true
	}
	var exists int
	if err := db.DB.QueryRow(`
		SELECT 1 FROM vessel_review_requests WHERE application_id = ? AND requested_staff_id = ? LIMIT 1
	`, app.ID, userID).Scan(&exists); err == nil {
		return true
	}
	// Members of a department flagged "notify on submission" get read-only
	// access to every application, for acknowledgment/FYI purposes.
	if err := db.DB.QueryRow(`
		SELECT 1 FROM users u JOIN departments d ON d.id = u.department_id
		WHERE u.id = ? AND d.notify_on_submission = TRUE LIMIT 1
	`, userID).Scan(&exists); err == nil {
		return true
	}
	return false
}

// GET /api/vessel-applications/staff-users - Admin/ANP HSE/HSE Officer
// For the "request review from" picker.
func ListStaffUsers(c *gin.Context) {
	roleName := c.GetString("roleName")
	if roleName != "Admin" && roleName != "ANP HSE" && roleName != "HSE Officer" {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to this"})
		return
	}

	rows, err := db.DB.Query(`
		SELECT u.id, u.name FROM users u JOIN roles r ON r.id = u.role_id
		WHERE r.name = 'Staff' AND u.status = 'Active'
		ORDER BY u.name
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch staff"})
		return
	}
	defer rows.Close()

	type staffUser struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	list := []staffUser{}
	for rows.Next() {
		var s staffUser
		if err := rows.Scan(&s.ID, &s.Name); err == nil {
			list = append(list, s)
		}
	}
	c.JSON(http.StatusOK, list)
}

// POST /api/vessel-applications/:id/review-requests
// Assigned HSE Officer (or Admin) only. Body: { "staffIds": ["...", "..."] }
func RequestStaffReview(c *gin.Context) {
	appID := c.Param("id")
	userID := c.GetString("userId")
	roleName := c.GetString("roleName")

	app, err := fetchVesselApplication(appID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Application not found"})
		return
	}
	isAssignedOfficer := app.AssignedHSEOfficerID != nil && *app.AssignedHSEOfficerID == userID
	if roleName != "Admin" && !isAssignedOfficer {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only the assigned HSE Officer can request a staff review"})
		return
	}

	var input struct {
		StaffIDs []string `json:"staffIds" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || len(input.StaffIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "staffIds is required"})
		return
	}

	vesselName := "a vessel"
	if app.VesselName != nil {
		vesselName = *app.VesselName
	}

	for _, staffID := range input.StaffIDs {
		var roleCheck string
		if err := db.DB.QueryRow(`
			SELECT r.name FROM users u JOIN roles r ON r.id = u.role_id WHERE u.id = ?
		`, staffID).Scan(&roleCheck); err != nil || roleCheck != "Staff" {
			continue // silently skip invalid ids rather than failing the whole batch
		}

		if _, err := db.DB.Exec(`
			INSERT IGNORE INTO vessel_review_requests (id, application_id, requested_staff_id, requested_by)
			VALUES (?, ?, ?, ?)
		`, uuid.NewString(), appID, staffID, userID); err != nil {
			continue
		}

		utils.CreateNotification(
			staffID,
			"You've been asked to review a vessel application",
			app.ApplicationNumber+" - "+vesselName,
			"/vessel/review/"+appID,
		)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review requests sent"})
}

// GET /api/vessel-applications/:id/review-requests
func ListReviewRequests(c *gin.Context) {
	appID := c.Param("id")
	userID := c.GetString("userId")
	roleName := c.GetString("roleName")

	app, err := fetchVesselApplication(appID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Application not found"})
		return
	}
	if !canCollaborateOnApplication(app, userID, roleName) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to this application"})
		return
	}

	rows, err := db.DB.Query(`
		SELECT r.requested_staff_id, u.name, r.created_at
		FROM vessel_review_requests r
		JOIN users u ON u.id = r.requested_staff_id
		WHERE r.application_id = ?
		ORDER BY r.created_at
	`, appID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch review requests"})
		return
	}
	defer rows.Close()

	type reviewRequest struct {
		StaffID   string `json:"staffId"`
		StaffName string `json:"staffName"`
		CreatedAt string `json:"createdAt"`
	}
	list := []reviewRequest{}
	for rows.Next() {
		var r reviewRequest
		var createdAt time.Time
		if err := rows.Scan(&r.StaffID, &r.StaffName, &createdAt); err == nil {
			r.CreatedAt = createdAt.Format(time.RFC3339)
			list = append(list, r)
		}
	}
	c.JSON(http.StatusOK, list)
}

// POST /api/vessel-applications/:id/comments - any collaborator
// Body: { "comment": "..." }
func AddReviewComment(c *gin.Context) {
	appID := c.Param("id")
	userID := c.GetString("userId")
	roleName := c.GetString("roleName")

	app, err := fetchVesselApplication(appID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Application not found"})
		return
	}
	if !canCollaborateOnApplication(app, userID, roleName) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to this application"})
		return
	}

	var input struct {
		Comment string `json:"comment" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || strings.TrimSpace(input.Comment) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "comment is required"})
		return
	}

	if _, err := db.DB.Exec(`
		INSERT INTO vessel_review_comments (id, application_id, author_id, comment)
		VALUES (?, ?, ?, ?)
	`, uuid.NewString(), appID, userID, strings.TrimSpace(input.Comment)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save comment"})
		return
	}

	// Let the assigned HSE Officer know a new comment came in, so they
	// don't have to keep checking back manually.
	if app.AssignedHSEOfficerID != nil && *app.AssignedHSEOfficerID != "" && *app.AssignedHSEOfficerID != userID {
		utils.CreateNotification(
			*app.AssignedHSEOfficerID,
			"New review comment on "+app.ApplicationNumber,
			input.Comment,
			"/vessel/review/"+appID,
		)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment added"})
}

// GET /api/vessel-applications/:id/comments - any collaborator
func ListReviewComments(c *gin.Context) {
	appID := c.Param("id")
	userID := c.GetString("userId")
	roleName := c.GetString("roleName")

	app, err := fetchVesselApplication(appID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Application not found"})
		return
	}
	if !canCollaborateOnApplication(app, userID, roleName) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to this application"})
		return
	}

	rows, err := db.DB.Query(`
		SELECT c.id, c.author_id, u.name, r.name, c.comment, c.created_at
		FROM vessel_review_comments c
		JOIN users u ON u.id = c.author_id
		JOIN roles r ON r.id = u.role_id
		WHERE c.application_id = ?
		ORDER BY c.created_at
	`, appID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch comments"})
		return
	}
	defer rows.Close()

	type comment struct {
		ID         string `json:"id"`
		AuthorID   string `json:"authorId"`
		AuthorName string `json:"authorName"`
		AuthorRole string `json:"authorRole"`
		Comment    string `json:"comment"`
		CreatedAt  string `json:"createdAt"`
	}
	list := []comment{}
	for rows.Next() {
		var cm comment
		var createdAt time.Time
		if err := rows.Scan(&cm.ID, &cm.AuthorID, &cm.AuthorName, &cm.AuthorRole, &cm.Comment, &createdAt); err == nil {
			cm.CreatedAt = createdAt.Format(time.RFC3339)
			list = append(list, cm)
		}
	}
	c.JSON(http.StatusOK, list)
}

// POST /api/vessel-applications/:id/notify-anp-ready
// Assigned HSE Officer (or Admin) only - signals ANP HSE that the internal
// review/comment round is done and this one's ready for their decision.
// Doesn't change status (it's already Submitted and visible to ANP HSE) -
// this is purely a courtesy heads-up notification.
func NotifyANPReady(c *gin.Context) {
	appID := c.Param("id")
	userID := c.GetString("userId")
	roleName := c.GetString("roleName")

	app, err := fetchVesselApplication(appID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Application not found"})
		return
	}
	isAssignedOfficer := app.AssignedHSEOfficerID != nil && *app.AssignedHSEOfficerID == userID
	if roleName != "Admin" && !isAssignedOfficer {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only the assigned HSE Officer can send this"})
		return
	}

	vesselName := "a vessel"
	if app.VesselName != nil {
		vesselName = *app.VesselName
	}
	utils.NotifyReviewers(
		"Ready for decision",
		app.ApplicationNumber+" - "+vesselName+" - internal review complete",
		"/vessel/review/"+appID,
	)

	c.JSON(http.StatusOK, gin.H{"message": "ANP HSE notified"})
}
