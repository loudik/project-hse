package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hse-backend-go/db"
	"hse-backend-go/utils"
)

type department struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	NotifyOnSubmission bool   `json:"notifyOnSubmission"`
	MemberCount        int    `json:"memberCount"`
}

// GET /api/departments - Admin/ANP HSE
func ListDepartments(c *gin.Context) {
	rows, err := db.DB.Query(`
		SELECT d.id, d.name, d.notify_on_submission, COUNT(u.id) as member_count
		FROM departments d
		LEFT JOIN users u ON u.department_id = d.id
		GROUP BY d.id, d.name, d.notify_on_submission
		ORDER BY d.name
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch departments"})
		return
	}
	defer rows.Close()

	list := []department{}
	for rows.Next() {
		var d department
		if err := rows.Scan(&d.ID, &d.Name, &d.NotifyOnSubmission, &d.MemberCount); err == nil {
			list = append(list, d)
		}
	}
	c.JSON(http.StatusOK, list)
}

// POST /api/departments - Admin only
// Body: { "name": "..." }
func CreateDepartment(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || strings.TrimSpace(input.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "name is required"})
		return
	}

	id := uuid.NewString()
	if _, err := db.DB.Exec(`INSERT INTO departments (id, name) VALUES (?, ?)`, id, strings.TrimSpace(input.Name)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create department (name may already exist)"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "name": input.Name})
}

// PATCH /api/departments/:id - Admin only
// Body: { "name": "..." (optional), "notifyOnSubmission": true/false (optional) }
func UpdateDepartment(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Name               *string `json:"name"`
		NotifyOnSubmission *bool   `json:"notifyOnSubmission"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	if input.Name != nil && strings.TrimSpace(*input.Name) != "" {
		if _, err := db.DB.Exec(`UPDATE departments SET name = ? WHERE id = ?`, strings.TrimSpace(*input.Name), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update name"})
			return
		}
	}
	if input.NotifyOnSubmission != nil {
		if _, err := db.DB.Exec(`UPDATE departments SET notify_on_submission = ? WHERE id = ?`, *input.NotifyOnSubmission, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update notification setting"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Department updated"})
}

// DELETE /api/departments/:id - Admin only
func DeleteDepartment(c *gin.Context) {
	id := c.Param("id")
	if _, err := db.DB.Exec(`DELETE FROM departments WHERE id = ?`, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete department"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Department deleted"})
}

// PATCH /api/users/:id/department - Admin only
// Body: { "departmentId": "..." } or { "departmentId": null } to unassign
func UpdateUserDepartment(c *gin.Context) {
	userID := c.Param("id")

	var input struct {
		DepartmentID *string `json:"departmentId"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	if _, err := db.DB.Exec(`UPDATE users SET department_id = ? WHERE id = ?`, input.DepartmentID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user's department"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User's department updated"})
}

// GET /api/departments/:id/members - Admin/ANP HSE
func ListDepartmentMembers(c *gin.Context) {
	deptID := c.Param("id")

	rows, err := db.DB.Query(`SELECT id, name, email FROM users WHERE department_id = ? ORDER BY name`, deptID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch members"})
		return
	}
	defer rows.Close()

	type member struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	list := []member{}
	for rows.Next() {
		var m member
		if err := rows.Scan(&m.ID, &m.Name, &m.Email); err == nil {
			list = append(list, m)
		}
	}
	c.JSON(http.StatusOK, list)
}

// notifyAcknowledgmentDepartments is called from SubmitVesselApplication.
// Notifies every member of every department flagged notify_on_submission -
// this is FYI-only, read access is granted separately via
// canCollaborateOnApplication checking department membership.
func notifyAcknowledgmentDepartments(appID, applicationNumber, vesselName string) {
	rows, err := db.DB.Query(`
		SELECT u.id FROM users u
		JOIN departments d ON d.id = u.department_id
		WHERE d.notify_on_submission = TRUE AND u.status = 'Active'
	`)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			continue
		}
		utils.CreateNotification(userID,
			"New vessel entry application submitted (FYI)",
			applicationNumber+" - "+vesselName,
			"/vessel/review/"+appID,
		)
	}
}
