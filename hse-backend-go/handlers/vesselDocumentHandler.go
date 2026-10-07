package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hse-backend-go/db"
	"hse-backend-go/utils"
)

func UploadVesselDocument(c *gin.Context) {
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

	category := c.PostForm("category")
	documentKey := c.PostForm("documentKey")
	label := c.PostForm("label")
	notApplicable := c.PostForm("notApplicable") == "true"
	naReason := c.PostForm("naReason")
	dateIssued := c.PostForm("dateIssued")
	dateExpired := c.PostForm("dateExpired")

	if category == "" || documentKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "category and documentKey are required"})
		return
	}

	if !notApplicable {
		if err := validateCertificateDates(dateIssued, dateExpired); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
	}

	var filePath, fileName string
	var fileSize int

	if !notApplicable {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "A file is required unless marked Not Applicable"})
			return
		}
		if fileHeader.Size > 20*1024*1024 { // 20MB limit
			c.JSON(http.StatusBadRequest, gin.H{"message": "File is too large (max 20MB)"})
			return
		}

		objectPath, err := utils.UploadFile(fileHeader, "vessel-applications/"+appID+"/"+category)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to upload file: " + err.Error()})
			return
		}
		filePath = objectPath
		fileName = fileHeader.Filename
		fileSize = int(fileHeader.Size)
	}

	id := uuid.NewString()
	_, err = db.DB.Exec(`
		INSERT INTO vessel_application_documents
			(id, application_id, category, document_key, label, not_applicable, na_reason,
			 date_issued, date_expired, file_path, file_name, file_size, uploaded_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
		ON DUPLICATE KEY UPDATE
			label = VALUES(label), not_applicable = VALUES(not_applicable), na_reason = VALUES(na_reason),
			date_issued = VALUES(date_issued), date_expired = VALUES(date_expired),
			file_path = VALUES(file_path), file_name = VALUES(file_name), file_size = VALUES(file_size),
			uploaded_at = NOW(), reminder_sent_at = NULL
	`, id, appID, category, documentKey, nullableString(label), notApplicable, nullableString(naReason),
		nullableDate(dateIssued), nullableDate(dateExpired), nullableString(filePath), nullableString(fileName), fileSize,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save document record: " + err.Error()})
		return
	}

	// Warn (don't block) if the certificate expires before the vessel is
	// due to leave the contract area - the operator will likely need a
	// renewed certificate before then.
	if !notApplicable && dateExpired != "" && app.EntryDateTo != nil && *app.EntryDateTo != "" {
		expTime, err1 := time.Parse("2006-01-02", dateExpired)
		exitTime, err2 := time.Parse("2006-01-02", *app.EntryDateTo)
		if err1 == nil && err2 == nil && expTime.Before(exitTime) {
			docLabel := documentKey
			if label != "" {
				docLabel = label
			}
			msg := fmt.Sprintf(
				"%s expires on %s, before the vessel's planned exit on %s. You'll need to upload a renewed certificate before then.",
				docLabel, dateExpired, *app.EntryDateTo,
			)
			utils.CreateNotification(userID, "Certificate expires before vessel exit date", msg, "/vessel/summary/"+appID)
			if applicantEmail := lookupUserEmail(userID); applicantEmail != "" {
				_ = utils.SendEmail(applicantEmail, "Certificate expires before vessel exit date", "<p>"+msg+"</p>")
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document saved"})
}

func lookupUserEmail(userID string) string {
	var email string
	if err := db.DB.QueryRow(`SELECT email FROM users WHERE id = ?`, userID).Scan(&email); err != nil {
		return ""
	}
	return email
}

// GET /api/vessel-applications/:id/documents - list all documents for an application
func ListVesselDocuments(c *gin.Context) {
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

	rows, err := db.DB.Query(`
		SELECT id, category, document_key, label, not_applicable, na_reason,
		       date_issued, date_expired, file_path, file_name, file_size, uploaded_at
		FROM vessel_application_documents WHERE application_id = ?
	`, appID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch documents"})
		return
	}
	defer rows.Close()

	type doc struct {
		ID            string  `json:"id"`
		Category      string  `json:"category"`
		DocumentKey   string  `json:"documentKey"`
		Label         *string `json:"label"`
		NotApplicable bool    `json:"notApplicable"`
		NaReason      *string `json:"naReason"`
		DateIssued    *string `json:"dateIssued"`
		DateExpired   *string `json:"dateExpired"`
		FileName      *string `json:"fileName"`
		FileSize      *int    `json:"fileSize"`
		UploadedAt    *string `json:"uploadedAt"`
	}

	list := []doc{}
	for rows.Next() {
		var d doc
		var label, naReason, fileName sql.NullString
		var dateIssued, dateExpired sql.NullTime
		var filePath sql.NullString
		var fileSize sql.NullInt64
		var uploadedAt sql.NullTime

		if err := rows.Scan(&d.ID, &d.Category, &d.DocumentKey, &label, &d.NotApplicable, &naReason,
			&dateIssued, &dateExpired, &filePath, &fileName, &fileSize, &uploadedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read document data"})
			return
		}
		if label.Valid {
			d.Label = &label.String
		}
		if naReason.Valid {
			d.NaReason = &naReason.String
		}
		if dateIssued.Valid {
			s := dateIssued.Time.Format("2006-01-02")
			d.DateIssued = &s
		}
		if dateExpired.Valid {
			s := dateExpired.Time.Format("2006-01-02")
			d.DateExpired = &s
		}
		if fileName.Valid {
			d.FileName = &fileName.String
		}
		if fileSize.Valid {
			size := int(fileSize.Int64)
			d.FileSize = &size
		}
		if uploadedAt.Valid {
			s := uploadedAt.Time.Format("2006-01-02T15:04:05Z07:00")
			d.UploadedAt = &s
		}
		list = append(list, d)
	}
	c.JSON(http.StatusOK, list)
}

// GET /api/vessel-applications/:id/documents/:docId/download
func DownloadVesselDocument(c *gin.Context) {
	docID := c.Param("docId")
	userID := c.GetString("userId")
	roleName := c.GetString("roleName")
	appID := c.Param("id")

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

	var filePath sql.NullString
	err = db.DB.QueryRow(`SELECT file_path FROM vessel_application_documents WHERE id = ? AND application_id = ?`, docID, appID).Scan(&filePath)
	if err == sql.ErrNoRows || !filePath.Valid {
		c.JSON(http.StatusNotFound, gin.H{"message": "File not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch document"})
		return
	}

	url, err := utils.GetFileDownloadURL(filePath.String)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate download link"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
}

// GET /api/vessel-applications/review/:id/documents - ANP HSE/Admin only
func ListVesselDocumentsForReview(c *gin.Context) {
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
		SELECT id, category, document_key, label, not_applicable, na_reason,
		       date_issued, date_expired, file_path, file_name, file_size, uploaded_at
		FROM vessel_application_documents WHERE application_id = ?
	`, appID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch documents"})
		return
	}
	defer rows.Close()

	type doc struct {
		ID            string  `json:"id"`
		Category      string  `json:"category"`
		DocumentKey   string  `json:"documentKey"`
		Label         *string `json:"label"`
		NotApplicable bool    `json:"notApplicable"`
		NaReason      *string `json:"naReason"`
		DateIssued    *string `json:"dateIssued"`
		DateExpired   *string `json:"dateExpired"`
		FileName      *string `json:"fileName"`
		FileSize      *int    `json:"fileSize"`
		UploadedAt    *string `json:"uploadedAt"`
	}

	list := []doc{}
	for rows.Next() {
		var d doc
		var label, naReason, fileName sql.NullString
		var dateIssued, dateExpired sql.NullTime
		var filePath sql.NullString
		var fileSize sql.NullInt64
		var uploadedAt sql.NullTime

		if err := rows.Scan(&d.ID, &d.Category, &d.DocumentKey, &label, &d.NotApplicable, &naReason,
			&dateIssued, &dateExpired, &filePath, &fileName, &fileSize, &uploadedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read document data"})
			return
		}
		if label.Valid {
			d.Label = &label.String
		}
		if naReason.Valid {
			d.NaReason = &naReason.String
		}
		if dateIssued.Valid {
			s := dateIssued.Time.Format("2006-01-02")
			d.DateIssued = &s
		}
		if dateExpired.Valid {
			s := dateExpired.Time.Format("2006-01-02")
			d.DateExpired = &s
		}
		if fileName.Valid {
			d.FileName = &fileName.String
		}
		if fileSize.Valid {
			size := int(fileSize.Int64)
			d.FileSize = &size
		}
		if uploadedAt.Valid {
			s := uploadedAt.Time.Format("2006-01-02T15:04:05Z07:00")
			d.UploadedAt = &s
		}
		list = append(list, d)
	}
	c.JSON(http.StatusOK, list)
}

func validateCertificateDates(dateIssued, dateExpired string) error {
	today := time.Now().Truncate(24 * time.Hour)

	var issued, expired time.Time
	var hasIssued, hasExpired bool

	if dateIssued != "" {
		t, err := time.Parse("2006-01-02", dateIssued)
		if err != nil {
			return fmt.Errorf("issue date is not a valid date")
		}
		if t.After(today) {
			return fmt.Errorf("issue date cannot be in the future")
		}
		issued = t
		hasIssued = true
	}

	if dateExpired != "" {
		t, err := time.Parse("2006-01-02", dateExpired)
		if err != nil {
			return fmt.Errorf("expiry date is not a valid date")
		}
		if t.Before(today) {
			return fmt.Errorf("this certificate has expired - please upload a valid, non-expired certificate")
		}
		expired = t
		hasExpired = true
	}

	if hasIssued && hasExpired && !expired.After(issued) {
		return fmt.Errorf("expiry date must be after the issue date")
	}

	return nil
}
