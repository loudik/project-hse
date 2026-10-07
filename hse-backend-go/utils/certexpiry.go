package utils

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"hse-backend-go/db"
)

type expiringDoc struct {
	DocID             string
	AppID             string
	ApplicationNumber string
	DocumentKey       string
	Label             sql.NullString
	DateExpired       time.Time
	CreatedBy         string
	ApplicantEmail    string
	ApplicantName     string
	AssignedOfficerID sql.NullString
}

// CheckExpiringCertificates finds certificates on Approved vessel applications
// that expire within the next 7 days and haven't been reminded about yet,
// notifies the applicant and all HSE reviewers (in-app + email), and marks
// them as reminded so the same expiry doesn't trigger repeat notifications
// every day. Re-uploading a renewed certificate resets the reminder state
// (see the ON DUPLICATE KEY UPDATE in UploadVesselDocument), so a fresh
// reminder cycle starts for the new expiry date.
func CheckExpiringCertificates() {
	rows, err := db.DB.Query(`
		SELECT d.id, d.application_id, a.application_number, d.document_key, d.label, d.date_expired,
		       a.created_by, u.email, u.name, a.assigned_hse_officer_id
		FROM vessel_application_documents d
		JOIN vessel_applications a ON a.id = d.application_id
		JOIN users u ON u.id = a.created_by
		WHERE a.status = 'Approved'
		  AND d.not_applicable = FALSE
		  AND d.date_expired IS NOT NULL
		  AND d.date_expired BETWEEN CURDATE() AND DATE_ADD(CURDATE(), INTERVAL 7 DAY)
		  AND d.reminder_sent_at IS NULL
	`)
	if err != nil {
		log.Printf("[cert-expiry] failed to query expiring certificates: %v", err)
		return
	}

	var docs []expiringDoc
	for rows.Next() {
		var d expiringDoc
		var label sql.NullString
		if err := rows.Scan(&d.DocID, &d.AppID, &d.ApplicationNumber, &d.DocumentKey, &label, &d.DateExpired, &d.CreatedBy, &d.ApplicantEmail, &d.ApplicantName, &d.AssignedOfficerID); err != nil {
			log.Printf("[cert-expiry] failed to scan row: %v", err)
			continue
		}
		d.Label = label
		docs = append(docs, d)
	}
	rows.Close()

	for _, d := range docs {
		notifyExpiringCertificate(d)
	}
}

// CheckExpiringAuthorisations notifies Operator + ANP HSE + assigned HSE
// Officer when a vessel's overall entry authorisation (entry_date_to) is
// within 10 days of expiring - distinct from individual certificate
// expiry (CheckExpiringCertificates, which checks each document).
func CheckExpiringAuthorisations() {
	rows, err := db.DB.Query(`
		SELECT a.id, a.application_number, a.vessel_name, a.entry_date_to,
		       a.created_by, u.email, u.name, a.assigned_hse_officer_id
		FROM vessel_applications a
		JOIN users u ON u.id = a.created_by
		WHERE a.status = 'Approved'
		  AND a.entry_date_to IS NOT NULL
		  AND a.entry_date_to BETWEEN CURDATE() AND DATE_ADD(CURDATE(), INTERVAL 10 DAY)
		  AND a.authorisation_reminder_sent_at IS NULL
	`)
	if err != nil {
		log.Printf("[authorisation-expiry] failed to query: %v", err)
		return
	}

	type expiringApp struct {
		ID                string
		ApplicationNumber string
		VesselName        sql.NullString
		EntryDateTo       time.Time
		CreatedBy         string
		ApplicantEmail    string
		ApplicantName     string
		AssignedOfficerID sql.NullString
	}

	var apps []expiringApp
	for rows.Next() {
		var a expiringApp
		if err := rows.Scan(&a.ID, &a.ApplicationNumber, &a.VesselName, &a.EntryDateTo,
			&a.CreatedBy, &a.ApplicantEmail, &a.ApplicantName, &a.AssignedOfficerID); err != nil {
			continue
		}
		apps = append(apps, a)
	}
	rows.Close()

	for _, a := range apps {
		vesselName := "your vessel"
		if a.VesselName.Valid {
			vesselName = a.VesselName.String
		}
		expiryDate := a.EntryDateTo.Format("02 Jan 2006")

		operatorMsg := fmt.Sprintf("%s (%s) expires on %s. Please request an extension if needed.", a.ApplicationNumber, vesselName, expiryDate)
		CreateNotification(a.CreatedBy, "Entry authorisation expiring soon", operatorMsg, "/vessel/"+a.ID)
		if a.ApplicantEmail != "" {
			_ = SendEmail(a.ApplicantEmail, "Entry authorisation expiring soon - "+a.ApplicationNumber, "<p>"+operatorMsg+"</p>")
		}

		reviewerMsg := fmt.Sprintf("%s (%s) - operator's entry authorisation expires on %s.", a.ApplicationNumber, vesselName, expiryDate)
		NotifyReviewers("Entry authorisation expiring soon", reviewerMsg, "/vessel/review/"+a.ID)
		notifyReviewersByEmail("Entry authorisation expiring soon - "+a.ApplicationNumber, "<p>"+reviewerMsg+"</p>")

		if a.AssignedOfficerID.Valid && a.AssignedOfficerID.String != "" {
			CreateNotification(a.AssignedOfficerID.String, "Entry authorisation expiring soon", reviewerMsg, "/vessel/review/"+a.ID)
		}

		if _, err := db.DB.Exec(`UPDATE vessel_applications SET authorisation_reminder_sent_at = NOW() WHERE id = ?`, a.ID); err != nil {
			log.Printf("[authorisation-expiry] failed to mark reminder sent for %s: %v", a.ApplicationNumber, err)
		}
	}
}

func notifyExpiringCertificate(d expiringDoc) {
	docName := d.DocumentKey
	if d.Label.Valid && d.Label.String != "" {
		docName = d.Label.String
	}
	expiryStr := d.DateExpired.Format("02 Jan 2006")

	// Operator
	operatorMsg := fmt.Sprintf(
		"%s for application %s expires on %s. Please upload a renewed certificate before then.",
		docName, d.ApplicationNumber, expiryStr,
	)
	CreateNotification(d.CreatedBy, "Certificate expiring soon", operatorMsg, "/vessel/summary/"+d.AppID)
	if d.ApplicantEmail != "" {
		_ = SendEmail(d.ApplicantEmail, "Certificate expiring soon - "+d.ApplicationNumber,
			fmt.Sprintf("<p>Hi %s,</p><p>%s</p>", d.ApplicantName, operatorMsg))
	}

	// HSE reviewers
	reviewerMsg := fmt.Sprintf(
		"%s for application %s (submitted by %s) expires on %s.",
		docName, d.ApplicationNumber, d.ApplicantName, expiryStr,
	)
	NotifyReviewers("Certificate expiring soon", reviewerMsg, "/vessel/review/"+d.AppID)
	notifyReviewersByEmail("Certificate expiring soon - "+d.ApplicationNumber, "<p>"+reviewerMsg+"</p>")

	if d.AssignedOfficerID.Valid && d.AssignedOfficerID.String != "" {
		CreateNotification(d.AssignedOfficerID.String, "Certificate expiring soon", reviewerMsg, "/vessel/review/"+d.AppID)
	}

	if _, err := db.DB.Exec(`UPDATE vessel_application_documents SET reminder_sent_at = NOW() WHERE id = ?`, d.DocID); err != nil {
		log.Printf("[cert-expiry] failed to mark reminder sent for document %s: %v", d.DocID, err)
	}
}

func notifyReviewersByEmail(subject, htmlBody string) {
	rows, err := db.DB.Query(`
		SELECT u.email FROM users u JOIN roles r ON r.id = u.role_id
		WHERE r.name IN ('Admin', 'ANP HSE') AND u.status = 'Active'
	`)
	if err != nil {
		log.Printf("[cert-expiry] failed to list reviewer emails: %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err == nil && email != "" {
			_ = SendEmail(email, subject, htmlBody)
		}
	}
}
