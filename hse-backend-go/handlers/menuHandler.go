package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"hse-backend-go/db"
	"hse-backend-go/models"
)

type menuRow struct {
	ID                  int
	ParentID            *int
	Name                string
	Path                *string
	Icon                *string
	RequiresOrgApproval bool
}

func GetMyMenu(c *gin.Context) {
	roleID := c.GetInt("roleId")
	userID := c.GetString("userId")

	orgApproved := userHasApprovedOrganization(userID)

	rows, err := db.DB.Query(`
		SELECT m.id, m.parent_id, m.name, m.path, m.icon, m.requires_org_approval
		FROM menus m
		JOIN role_menu_access rma ON rma.menu_id = m.id
		WHERE rma.role_id = ? AND rma.can_view = TRUE
		ORDER BY m.order_index ASC
	`, roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch menu"})
		return
	}
	defer rows.Close()

	var records []menuRow
	for rows.Next() {
		var r menuRow
		if err := rows.Scan(&r.ID, &r.ParentID, &r.Name, &r.Path, &r.Icon, &r.RequiresOrgApproval); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read menu"})
			return
		}
		if r.RequiresOrgApproval && !orgApproved {
			continue // hide this menu item - organization not approved yet
		}
		records = append(records, r)
	}

	// Pass 1: collect top-level menus (parent_id NULL), store their position
	roots := []models.MenuNode{}
	posByID := map[int]int{}
	for _, r := range records {
		if r.ParentID == nil {
			roots = append(roots, models.MenuNode{ID: r.ID, Name: r.Name, Path: r.Path, Icon: r.Icon})
			posByID[r.ID] = len(roots) - 1
		}
	}

	// Pass 2: attach submenus to their parent
	for _, r := range records {
		if r.ParentID == nil {
			continue
		}
		if pos, ok := posByID[*r.ParentID]; ok {
			roots[pos].Children = append(roots[pos].Children, models.MenuNode{
				ID: r.ID, Name: r.Name, Path: r.Path, Icon: r.Icon,
			})
		}
	}

	c.JSON(http.StatusOK, roots)
}

func userHasApprovedOrganization(userID string) bool {
	var status sql.NullString
	err := db.DB.QueryRow(`
		SELECT o.status
		FROM users u JOIN organizations o ON o.id = u.organization_id
		WHERE u.id = ?
	`, userID).Scan(&status)
	if err != nil {
		return false // no organization linked yet, or query failed -> not approved
	}
	return status.Valid && status.String == "Approved"
}
