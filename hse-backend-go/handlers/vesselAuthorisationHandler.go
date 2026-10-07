package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hse-backend-go/db"
	"hse-backend-go/utils"
)

const authorisationTemplatePath = "templates/vessel_entry_authorisation.docx"

type authorisationLetter struct {
	ID              string  `json:"id"`
	ReferenceNo     string  `json:"referenceNo"`
	FacilityName    string  `json:"facilityName"`
	DecreeLawNumber string  `json:"decreeLawNumber"`
	DecreeLawClause string  `json:"decreeLawClause"`
	InspectionDate  *string `json:"inspectionDate"`
	PSCName         string  `json:"pscName"`
	DigitalHash     string  `json:"digitalHash"`
	IssuedByName    string  `json:"issuedByName"`
	IssuedAt        string  `json:"issuedAt"`
}

func nullableStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// POST /api/vessel-applications/:id/authorisation - Admin/ANP HSE only
// multipart form-data: facilityName, decreeLawNumber, decreeLawClause,
// inspectionDate (optional, YYYY-MM-DD), signature (file - either a
// drawn-then-exported PNG or an uploaded image, same field either way)
func IssueAuthorisationLetter(c *gin.Context) {
	appID := c.Param("id")
	issuerID := c.GetString("userId")

	app, err := fetchVesselApplication(appID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Application not found"})
		return
	}
	if app.Status != "Approved" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "The authorisation letter can only be issued for an Approved application"})
		return
	}

	org, err := fetchOrganization(app.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to load organization"})
		return
	}

	facilityName := strings.TrimSpace(c.PostForm("facilityName"))
	decreeLawNumber := strings.TrimSpace(c.PostForm("decreeLawNumber"))
	decreeLawClause := strings.TrimSpace(c.PostForm("decreeLawClause"))
	inspectionDate := strings.TrimSpace(c.PostForm("inspectionDate"))
	if facilityName == "" || decreeLawNumber == "" || decreeLawClause == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "facilityName, decreeLawNumber and decreeLawClause are required"})
		return
	}

	fileHeader, err := c.FormFile("signature")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "A signature (drawn or uploaded) is required"})
		return
	}
	sigFile, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read signature file"})
		return
	}
	defer sigFile.Close()
	sigBytes, err := io.ReadAll(sigFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read signature file"})
		return
	}
	sigExt := strings.TrimPrefix(strings.ToLower(filepath.Ext(fileHeader.Filename)), ".")
	if sigExt == "" {
		sigExt = "png"
	}

	// Reserve a reference number: insert with a temporary unique
	// placeholder, capture the auto-increment seq, then fill in the real
	// "ANP/HSE/S/NNNN" reference in the same transaction.
	id := uuid.NewString()
	tempRef := "PENDING-" + id
	issuedAt := time.Now()

	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to start transaction"})
		return
	}

	res, err := tx.Exec(`
		INSERT INTO vessel_authorisation_letters
			(id, application_id, reference_no, facility_name, decree_law_number, decree_law_clause,
			 inspection_date, psc_name, signature_path, digital_hash, docx_path, pdf_path, issued_by, issued_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, '', '', '', NULL, ?, ?)
	`, id, appID, tempRef, facilityName, decreeLawNumber, decreeLawClause,
		nullableDate(inspectionDate), org.Name, issuerID, issuedAt)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to reserve reference number: " + err.Error()})
		return
	}
	seq, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to determine reference number"})
		return
	}
	referenceNo := fmt.Sprintf("ANP/HSE/S/%04d", seq)

	// Digital hash: a SHA-256 fingerprint tying this specific letter to the
	// application, reference number, issuer, and timestamp - lets anyone
	// with the letter verify it wasn't altered or reissued under a
	// different reference.
	hashInput := fmt.Sprintf("%s|%s|%s|%s|%s", app.ID, referenceNo, issuerID, issuedAt.Format(time.RFC3339), app.ApplicationNumber)
	hashBytes := sha256.Sum256([]byte(hashInput))
	digitalHash := hex.EncodeToString(hashBytes[:])

	signaturePath := fmt.Sprintf("vessel-applications/%s/authorisation/%s-signature.%s", appID, referenceNo, sigExt)
	if err := utils.UploadBytes(sigBytes, signaturePath, "image/"+sigExt); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to upload signature"})
		return
	}

	values := map[string]string{
		"ISSUANCE_DATE":          issuedAt.Format("02 Jan 2006"),
		"REFERENCE_NO":           referenceNo,
		"OPERATOR_NAME":          strOrDash(nil),
		"OPERATOR_POSITION":      strOrDash(nil),
		"OPERATOR_ADDRESS":       strOrDash(org.Address),
		"FACILITY_NAME":          facilityName,
		"DECREE_LAW_NUMBER":      decreeLawNumber,
		"DECREE_LAW_CLAUSE":      decreeLawClause,
		"VESSEL_NAME":            strOrDash(app.VesselName),
		"VESSEL_IMO_NUMBER":      strOrDash(app.VesselIMONumber),
		"VESSEL_TYPE":            strOrDash(app.VesselType),
		"PORT_OF_REGISTRY":       strOrDash(app.PortOfRegistry),
		"CLASSIFICATION_SOCIETY": strOrDash(app.ClassificationSociety),
		"CLASS_ID_NUMBER":        strOrDash(app.ClassIDNumber),
		"PSC_NAME":               org.Name,
		"INSPECTION_DATE":        formatDisplayDate(nullableStringPtr(inspectionDate)),
		"VALIDITY_END":           formatDisplayDate(app.EntryDateTo),
		"DIGITAL_HASH":           digitalHash,
	}

	docBytes, err := utils.RenderDocxTemplate(authorisationTemplatePath, values)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to render authorisation letter: " + err.Error()})
		return
	}

	docBytes, err = utils.InsertImageIntoDocx(docBytes, sigBytes, sigExt, 1828800, 731520)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to embed signature: " + err.Error()})
		return
	}

	docxPath := fmt.Sprintf("vessel-applications/%s/authorisation/%s.docx", appID, referenceNo)
	docxContentType := "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	if err := utils.UploadBytes(docBytes, docxPath, docxContentType); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to upload authorisation letter"})
		return
	}

	var pdfPath *string
	if pdfBytes, err := utils.ConvertDocxToPDF(docBytes); err == nil {
		path := fmt.Sprintf("vessel-applications/%s/authorisation/%s.pdf", appID, referenceNo)
		if err := utils.UploadBytes(pdfBytes, path, "application/pdf"); err == nil {
			pdfPath = &path
		} else {
			log.Printf("[authorisation] PDF generated but failed to upload for %s: %v", referenceNo, err)
		}
	} else {
		log.Printf("[authorisation] PDF conversion failed for %s: %v", referenceNo, err)
	}

	if _, err := tx.Exec(`
		UPDATE vessel_authorisation_letters
		SET reference_no = ?, signature_path = ?, digital_hash = ?, docx_path = ?, pdf_path = ?
		WHERE id = ?
	`, referenceNo, signaturePath, digitalHash, docxPath, pdfPath, id); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to finalize authorisation letter"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save authorisation letter"})
		return
	}

	utils.LogAudit("vessel_application", appID, "authorisation_issued", "", "", issuerID, "Reference: "+referenceNo)
	utils.CreateNotification(app.CreatedBy, "Your vessel entry authorisation letter is ready", referenceNo, "/vessel/"+appID)

	c.JSON(http.StatusOK, gin.H{"message": "Authorisation letter issued", "referenceNo": referenceNo})
}

// GET /api/vessel-applications/:id/authorisation - latest issued letter, if any
func GetAuthorisationLetter(c *gin.Context) {
	appID := c.Param("id")
	userID := c.GetString("userId")
	roleName := c.GetString("roleName")

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

	var l authorisationLetter
	var inspectionDate sql.NullTime
	var issuedByName string
	var issuedAt time.Time
	err = db.DB.QueryRow(`
		SELECT l.id, l.reference_no, l.facility_name, l.decree_law_number, l.decree_law_clause,
		       l.inspection_date, l.psc_name, l.digital_hash, u.name, l.issued_at
		FROM vessel_authorisation_letters l
		JOIN users u ON u.id = l.issued_by
		WHERE l.application_id = ?
		ORDER BY l.issued_at DESC
		LIMIT 1
	`, appID).Scan(&l.ID, &l.ReferenceNo, &l.FacilityName, &l.DecreeLawNumber, &l.DecreeLawClause,
		&inspectionDate, &l.PSCName, &l.DigitalHash, &issuedByName, &issuedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "No authorisation letter issued yet"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch authorisation letter"})
		return
	}
	if inspectionDate.Valid {
		d := inspectionDate.Time.Format("2006-01-02")
		l.InspectionDate = &d
	}
	l.IssuedByName = issuedByName
	l.IssuedAt = issuedAt.Format(time.RFC3339)

	c.JSON(http.StatusOK, l)
}

// GET /api/vessel-applications/:id/authorisation/download?format=docx|pdf
func DownloadAuthorisationLetter(c *gin.Context) {
	appID := c.Param("id")
	userID := c.GetString("userId")
	roleName := c.GetString("roleName")
	format := c.DefaultQuery("format", "docx")

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

	column := "docx_path"
	if format == "pdf" {
		column = "pdf_path"
	}
	var path sql.NullString
	err = db.DB.QueryRow(fmt.Sprintf(`
		SELECT %s FROM vessel_authorisation_letters WHERE application_id = ? ORDER BY issued_at DESC LIMIT 1
	`, column), appID).Scan(&path)
	if err != nil || !path.Valid {
		c.JSON(http.StatusNotFound, gin.H{"message": "No authorisation letter available"})
		return
	}

	url, err := utils.GetFileDownloadURL(path.String)
	if err != nil {
		log.Printf("[authorisation] failed to generate download URL for path %s: %v", path.String, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate download link"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
}
