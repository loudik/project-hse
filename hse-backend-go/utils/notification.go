package utils

import (
	"log"

	"github.com/google/uuid"

	"hse-backend-go/db"
)

// CreateNotification inserts an in-app notification for one user.
// Best-effort by design, same reasoning as LogAudit - a notification
// failure should never block the action that triggered it.
func CreateNotification(userID, title, message, link string) {
	if userID == "" {
		return
	}
	id := uuid.NewString()
	_, err := db.DB.Exec(`
		INSERT INTO notifications (id, user_id, title, message, link)
		VALUES (?, ?, ?, ?, ?)
	`, id, userID, title, nullIfEmpty(message), nullIfEmpty(link))
	if err != nil {
		log.Printf("[notification] failed to create for user %s: %v", userID, err)
	}
}

// NotifyReviewers sends the same notification to every active Admin/ANP HSE
// user - used when something needs their attention (e.g. a new submission).
func NotifyReviewers(title, message, link string) {
	rows, err := db.DB.Query(`
		SELECT u.id FROM users u JOIN roles r ON r.id = u.role_id
		WHERE r.name IN ('Admin', 'ANP HSE') AND u.status = 'Active'
	`)
	if err != nil {
		log.Printf("[notification] failed to list reviewers: %v", err)
		return
	}
	defer rows.Close()

	var userIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			userIDs = append(userIDs, id)
		}
	}
	for _, uid := range userIDs {
		CreateNotification(uid, title, message, link)
	}
}
