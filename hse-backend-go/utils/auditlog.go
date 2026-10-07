package utils

import (
	"database/sql"
	"log"

	"github.com/google/uuid"

	"hse-backend-go/db"
)

func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// LogAudit records a status-change event in the audit_logs table:
// who did what, to which entity, and what the status was before/after.
// Best-effort by design - a logging failure must never block the actual
// business action (approval/submission/etc.) that triggered it, so this
// only logs errors, never returns one to the caller.
//
// entityType: "vessel_application" | "organization"
// action: "submitted" | "approved" | "rejected" | "withdrawn" | "reopened"
// performedBy: users.id of whoever triggered it - pass "" for system-generated events
func LogAudit(entityType, entityID, action, oldStatus, newStatus, performedBy, notes string) {
	id := uuid.NewString()
	_, err := db.DB.Exec(`
		INSERT INTO audit_logs (id, entity_type, entity_id, action, old_status, new_status, performed_by, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, id, entityType, entityID, action, nullIfEmpty(oldStatus), nullIfEmpty(newStatus), nullIfEmpty(performedBy), nullIfEmpty(notes))
	if err != nil {
		log.Printf("[audit] failed to record %s on %s/%s: %v", action, entityType, entityID, err)
	}
}

// AuditLogEntry is what GET .../history endpoints return.
type AuditLogEntry struct {
	ID              string  `json:"id"`
	Action          string  `json:"action"`
	OldStatus       *string `json:"oldStatus"`
	NewStatus       *string `json:"newStatus"`
	PerformedBy     *string `json:"performedBy"`     // users.id
	PerformedByName *string `json:"performedByName"` // users.name, resolved via JOIN
	Notes           *string `json:"notes"`
	CreatedAt       string  `json:"createdAt"`
}

// FetchAuditLog returns the history for one entity, oldest first (a
// natural reading order for a timeline).
func FetchAuditLog(entityType, entityID string) ([]AuditLogEntry, error) {
	rows, err := db.DB.Query(`
		SELECT al.id, al.action, al.old_status, al.new_status, al.performed_by, u.name, al.notes, al.created_at
		FROM audit_logs al
		LEFT JOIN users u ON u.id = al.performed_by
		WHERE al.entity_type = ? AND al.entity_id = ?
		ORDER BY al.created_at ASC
	`, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []AuditLogEntry{}
	for rows.Next() {
		var e AuditLogEntry
		var oldStatus, newStatus, performedBy, performedByName, notes sql.NullString
		var createdAt sql.NullTime
		if err := rows.Scan(&e.ID, &e.Action, &oldStatus, &newStatus, &performedBy, &performedByName, &notes, &createdAt); err != nil {
			return nil, err
		}
		if oldStatus.Valid {
			e.OldStatus = &oldStatus.String
		}
		if newStatus.Valid {
			e.NewStatus = &newStatus.String
		}
		if performedBy.Valid {
			e.PerformedBy = &performedBy.String
		}
		if performedByName.Valid {
			e.PerformedByName = &performedByName.String
		}
		if notes.Valid {
			e.Notes = &notes.String
		}
		if createdAt.Valid {
			e.CreatedAt = createdAt.Time.Format("2006-01-02T15:04:05Z07:00")
		}
		list = append(list, e)
	}
	return list, rows.Err()
}
