package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"hse-backend-go/db"
)

type userListItem struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	RoleID    int     `json:"roleId"`
	RoleName  string  `json:"roleName"`
	Status    string  `json:"status"`
	OrgName   *string `json:"organizationName"`
	CreatedAt string  `json:"createdAt"`
}

// GET /api/users - Admin only
func ListUsers(c *gin.Context) {
	rows, err := db.DB.Query(`
		SELECT u.id, u.name, u.email, u.role_id, r.name, u.status, u.organization_name, u.created_at
		FROM users u JOIN roles r ON r.id = u.role_id
		ORDER BY u.created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	list := []userListItem{}
	for rows.Next() {
		var u userListItem
		var orgName sql.NullString
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.RoleID, &u.RoleName, &u.Status, &orgName, &u.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read user data"})
			return
		}
		if orgName.Valid {
			u.OrgName = &orgName.String
		}
		list = append(list, u)
	}
	c.JSON(http.StatusOK, list)
}

// GET /api/roles - list of assignable roles (for the role dropdown)
func ListRoles(c *gin.Context) {
	rows, err := db.DB.Query(`SELECT id, name FROM roles ORDER BY name ASC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch roles"})
		return
	}
	defer rows.Close()

	type role struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	list := []role{}
	for rows.Next() {
		var r role
		if err := rows.Scan(&r.ID, &r.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read role data"})
			return
		}
		list = append(list, r)
	}
	c.JSON(http.StatusOK, list)
}

// PATCH /api/users/:id/role - Admin only. Body: { "roleId": 2 }
func UpdateUserRole(c *gin.Context) {
	userID := c.Param("id")
	actingUserID := c.GetString("userId")

	var body struct {
		RoleID int `json:"roleId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "roleId is required"})
		return
	}

	if userID == actingUserID {
		c.JSON(http.StatusBadRequest, gin.H{"message": "You cannot change your own role"})
		return
	}

	var exists int
	if err := db.DB.QueryRow(`SELECT 1 FROM roles WHERE id = ?`, body.RoleID).Scan(&exists); err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Role not found"})
		return
	}

	res, err := db.DB.Exec(`UPDATE users SET role_id = ? WHERE id = ?`, body.RoleID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update role"})
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role updated"})
}

// PATCH /api/users/:id/status - Admin only. Body: { "status": "Active" | "Suspended" }
func UpdateUserStatus(c *gin.Context) {
	userID := c.Param("id")
	actingUserID := c.GetString("userId")

	var body struct {
		Status string `json:"status" binding:"required,oneof=Active Suspended"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "status must be Active or Suspended"})
		return
	}

	if userID == actingUserID {
		c.JSON(http.StatusBadRequest, gin.H{"message": "You cannot change your own status"})
		return
	}

	res, err := db.DB.Exec(`UPDATE users SET status = ? WHERE id = ?`, body.Status, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update status"})
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}
