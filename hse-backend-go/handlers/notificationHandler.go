package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"hse-backend-go/db"
)

type notificationItem struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Message   *string `json:"message"`
	Link      *string `json:"link"`
	IsRead    bool    `json:"isRead"`
	CreatedAt string  `json:"createdAt"`
}

// GET /api/notifications - current user's notifications, newest first
func ListNotifications(c *gin.Context) {
	userID := c.GetString("userId")

	rows, err := db.DB.Query(`
		SELECT id, title, message, link, is_read, created_at
		FROM notifications WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT 50
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch notifications"})
		return
	}
	defer rows.Close()

	list := []notificationItem{}
	for rows.Next() {
		var n notificationItem
		var message, link sql.NullString
		var createdAt time.Time
		if err := rows.Scan(&n.ID, &n.Title, &message, &link, &n.IsRead, &createdAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read notification data"})
			return
		}
		if message.Valid {
			n.Message = &message.String
		}
		if link.Valid {
			n.Link = &link.String
		}
		n.CreatedAt = createdAt.Format("2006-01-02T15:04:05Z07:00")
		list = append(list, n)
	}
	c.JSON(http.StatusOK, list)
}

// GET /api/notifications/unread-count - for the bell icon badge
func GetUnreadNotificationCount(c *gin.Context) {
	userID := c.GetString("userId")
	var count int
	if err := db.DB.QueryRow(`SELECT COUNT(*) FROM notifications WHERE user_id = ? AND is_read = FALSE`, userID).Scan(&count); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to count notifications"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}

// PATCH /api/notifications/:id/read
func MarkNotificationRead(c *gin.Context) {
	userID := c.GetString("userId")
	id := c.Param("id")

	res, err := db.DB.Exec(`UPDATE notifications SET is_read = TRUE WHERE id = ? AND user_id = ?`, id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update notification"})
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Notification not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// PATCH /api/notifications/read-all
func MarkAllNotificationsRead(c *gin.Context) {
	userID := c.GetString("userId")
	if _, err := db.DB.Exec(`UPDATE notifications SET is_read = TRUE WHERE user_id = ? AND is_read = FALSE`, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update notifications"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
