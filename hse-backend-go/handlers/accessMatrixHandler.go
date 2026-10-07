package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"hse-backend-go/db"
)

type menuItem struct {
	ID   int     `json:"id"`
	Name string  `json:"name"`
	Path *string `json:"path"`
	Icon *string `json:"icon"`
}

type roleItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type roleAccessPair struct {
	RoleID  int  `json:"roleId"`
	MenuID  int  `json:"menuId"`
	CanView bool `json:"canView"`
}

// GET /api/admin/access-matrix - Admin only
// Returns every menu, every role, and the current role->menu visibility -
// the frontend renders this as a checkbox grid.
func GetAccessMatrix(c *gin.Context) {
	menuRows, err := db.DB.Query(`SELECT id, name, path, icon FROM menus ORDER BY order_index, id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch menus"})
		return
	}
	menus := []menuItem{}
	for menuRows.Next() {
		var m menuItem
		var path, icon sql.NullString
		if err := menuRows.Scan(&m.ID, &m.Name, &path, &icon); err == nil {
			if path.Valid {
				m.Path = &path.String
			}
			if icon.Valid {
				m.Icon = &icon.String
			}
			menus = append(menus, m)
		}
	}
	menuRows.Close()

	roleRows, err := db.DB.Query(`SELECT id, name FROM roles ORDER BY id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch roles"})
		return
	}
	roles := []roleItem{}
	for roleRows.Next() {
		var r roleItem
		if err := roleRows.Scan(&r.ID, &r.Name); err == nil {
			roles = append(roles, r)
		}
	}
	roleRows.Close()

	accessRows, err := db.DB.Query(`SELECT role_id, menu_id, can_view FROM role_menu_access`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch access"})
		return
	}
	access := []roleAccessPair{}
	for accessRows.Next() {
		var a roleAccessPair
		if err := accessRows.Scan(&a.RoleID, &a.MenuID, &a.CanView); err == nil {
			access = append(access, a)
		}
	}
	accessRows.Close()

	c.JSON(http.StatusOK, gin.H{
		"menus":  menus,
		"roles":  roles,
		"access": access,
	})
}

// PATCH /api/admin/access-matrix - Admin only
// Body: { "access": [{ "roleId": 1, "menuId": 2, "canView": true }, ...] }
// The frontend sends the complete desired state for every role/menu pair
// it rendered (checked = grant, unchecked = revoke) - applied in one
// transaction so a partial failure doesn't leave the matrix half-updated.
func UpdateAccessMatrix(c *gin.Context) {
	var input struct {
		Access []roleAccessPair `json:"access" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to start transaction"})
		return
	}

	for _, a := range input.Access {
		if a.CanView {
			if _, err := tx.Exec(`
				INSERT INTO role_menu_access (role_id, menu_id, can_view) VALUES (?, ?, TRUE)
				ON DUPLICATE KEY UPDATE can_view = TRUE
			`, a.RoleID, a.MenuID); err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update access: " + err.Error()})
				return
			}
		} else {
			if _, err := tx.Exec(`DELETE FROM role_menu_access WHERE role_id = ? AND menu_id = ?`, a.RoleID, a.MenuID); err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update access: " + err.Error()})
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save changes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Access updated"})
}
