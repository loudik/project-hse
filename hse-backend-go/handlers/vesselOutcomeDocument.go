package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"hse-backend-go/db"
	"hse-backend-go/models"
	"hse-backend-go/utils"
)

const outcomeTemplatePath = "templates/submit_vessel.docx"

// ---- small formatting helpers (dash for "no data", consistent with how
// the rest of the app already shows optional fields) ----

func strOrDash(s *string) string {
	if s == nil || *s == "" {
		return "-"
	}
	return *s
}

func boolLabel(b *bool) string {
	if b == nil {
		return "-"
	}
	if *b {
		return "Yes"
	}
	return "No"
}

func intOrDash(n *int) string {
	if n == nil {
		return "-"
	}
	return strconv.Itoa(*n)
}

func formatDisplayDate(iso *string) string {
	if iso == nil || *iso == "" {
		return "-"
	}
	t, err := time.Parse("2006-01-02", *iso)
	if err != nil {
		return *iso
	}
	return t.Format("02 Jan 2006")
}

// computeDuration renders both the date range and the day count, e.g.
// "01 Jan 2026 - 30 Jan 2026 (29 days)".
func computeDuration(from, to *string) string {
	if from == nil || to == nil || *from == "" || *to == "" {
		return "-"
	}
	tFrom, err1 := time.Parse("2006-01-02", *from)
	tTo, err2 := time.Parse("2006-01-02", *to)
	if err1 != nil || err2 != nil {
		return "-"
	}
	days := int(tTo.Sub(tFrom).Hours() / 24)
	return fmt.Sprintf("%s - %s (%d days)", tFrom.Format("02 Jan 2006"), tTo.Format("02 Jan 2006"), days)
}

func joinOrDash(items []string, sep string) string {
	if len(items) == 0 {
		return "-"
	}
	return strings.Join(items, sep)
}

// bulletList renders each item on its own line prefixed with "- ",
// joined with real line breaks (see escapeXMLText in utils/docxgen.go).
func bulletList(items []string) string {
	if len(items) == 0 {
		return "-"
	}
	lines := make([]string, len(items))
	for i, it := range items {
		lines[i] = "- " + it
	}
	return strings.Join(lines, "\n")
}

func buildScopeOfWork(app *models.VesselApplication) string {
	items := app.ScopeOfWork
	if app.ScopeOfWorkOther != nil && *app.ScopeOfWorkOther != "" {
		items = append(items, *app.ScopeOfWorkOther)
	}
	return joinOrDash(items, ", ")
}

func buildContractSummary(app *models.VesselApplication) string {
	contractType := strOrDash(app.ContractType)
	if app.ContractType != nil && *app.ContractType == "Other" && app.ContractTypeOther != nil && *app.ContractTypeOther != "" {
		contractType = *app.ContractTypeOther
	}
	lines := []string{
		"Status: " + strOrDash(app.ContractStatus),
		"Type: " + contractType,
		"Contract Number: " + strOrDash(app.ContractNumber),
		"Contact Email: " + strOrDash(app.ContractContactEmail),
	}
	return strings.Join(lines, "\n")
}

// ---- documents (vessel_application_documents) ----

type outcomeDocRow struct {
	Category      string
	DocumentKey   string
	Label         *string
	NotApplicable bool
}

func fetchOutcomeDocuments(appID string) ([]outcomeDocRow, error) {
	rows, err := db.DB.Query(`
		SELECT category, document_key, label, not_applicable
		FROM vessel_application_documents WHERE application_id = ?
	`, appID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []outcomeDocRow
	for rows.Next() {
		var d outcomeDocRow
		var label sql.NullString
		if err := rows.Scan(&d.Category, &d.DocumentKey, &label, &d.NotApplicable); err != nil {
			return nil, err
		}
		if label.Valid {
			d.Label = &label.String
		}
		list = append(list, d)
	}
	return list, rows.Err()
}

// formatDocumentList renders one line per document in `category`:
// "- <name>: Uploaded" or "- <name>: Not Applicable".
func formatDocumentList(docs []outcomeDocRow, category string) string {
	var lines []string
	for _, d := range docs {
		if d.Category != category {
			continue
		}
		name := d.DocumentKey
		if d.Label != nil && *d.Label != "" {
			name = *d.Label
		}
		status := "Uploaded"
		if d.NotApplicable {
			status = "Not Applicable"
		}
		lines = append(lines, fmt.Sprintf("- %s: %s", name, status))
	}
	if len(lines) == 0 {
		return "-"
	}
	return strings.Join(lines, "\n")
}

// findDocumentStatus looks up a single document by its document_key -
// used for the one-off clearance/declaration placeholders that aren't
// full lists.
//
// NOTE: the document_key values below ("dangerous_goods_declaration",
// "cargo_manifest", "customs_clearance", "immigration_clearance",
// "quarantine_clearance") are a best guess based on naming convention -
// please confirm these match what the frontend wizard actually saves
// (Vessel step 4/5, category = cargo_document / clearance_document) and
// adjust here if the real keys differ.
func findDocumentStatus(docs []outcomeDocRow, documentKey string) string {
	for _, d := range docs {
		if d.DocumentKey == documentKey {
			if d.NotApplicable {
				return "Not Applicable"
			}
			return "Submitted"
		}
	}
	return "Not provided"
}

// buildOutcomeDocumentValues maps a VesselApplication + its documents to
// every {{PLACEHOLDER}} used in templates/submit_vessel.docx.
func buildOutcomeDocumentValues(app *models.VesselApplication) (map[string]string, error) {
	docs, err := fetchOutcomeDocuments(app.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load documents: %w", err)
	}

	return map[string]string{
		"APPLICATION_NUMBER": app.ApplicationNumber,
		"STATUS":             app.Status,

		"ENTRY_CONDITION": strOrDash(app.EntryCondition),
		"PURPOSE":         strOrDash(app.Purpose),
		"ENTRY_DATE_FROM": formatDisplayDate(app.EntryDateFrom),
		"ENTRY_DATE_TO":   formatDisplayDate(app.EntryDateTo),
		"ENTRY_TYPE":      strOrDash(app.EntryType),
		"DURATION":        computeDuration(app.EntryDateFrom, app.EntryDateTo),

		"NOTIFICATION_EMAILS":     joinOrDash(app.NotificationEmails, ", "),
		"ENTRY_APPLICATION_TYPES": bulletList(app.EntryApplicationTypes),
		"CONTRACT_SUMMARY":        buildContractSummary(app),

		"VESSEL_NAME":            strOrDash(app.VesselName),
		"VESSEL_IMO_NUMBER":      strOrDash(app.VesselIMONumber),
		"VESSEL_OWNER":           strOrDash(app.VesselOwner),
		"VESSEL_TYPE":            strOrDash(app.VesselType),
		"FLAG_STATE":             strOrDash(app.FlagState),
		"PORT_OF_REGISTRY":       strOrDash(app.PortOfRegistry),
		"CLASSIFICATION_SOCIETY": strOrDash(app.ClassificationSociety),
		"CLASS_ID_NUMBER":        strOrDash(app.ClassIDNumber),
		"LENGTH_OVERALL":         strOrDash(app.LengthOverall),
		"DRAFT_VALUE":            strOrDash(app.DraftValue),
		"GROSS_TONNAGE":          strOrDash(app.GrossTonnage),
		"CALL_SIGN":              strOrDash(app.CallSign),
		"SCOPE_OF_WORK":          buildScopeOfWork(app),
		"OPERATION_MODE":         strOrDash(app.OperationMode),

		"NO_MAJOR_DEFICIENCIES":        boolLabel(app.NoMajorDeficiencies),
		"NO_DETENTION_12_MONTHS":       boolLabel(app.NoDetention12Months),
		"SAFETY_EQUIPMENT_OPERATIONAL": boolLabel(app.SafetyEquipmentOperational),
		"FIREFIGHTING_OPERATIONAL":     boolLabel(app.FirefightingOperational),
		"LIFESAVING_OPERATIONAL":       boolLabel(app.LifesavingOperational),
		"CREW_COUNT":                   intOrDash(app.CrewCount),
		"SURVEY_CREW_COUNT":            intOrDash(app.SurveyCrewCount),
		"CARGO_ONBOARD":                boolLabel(app.CargoOnboard),
		"CARGO_DESCRIPTION":            strOrDash(app.CargoDescription),
		"HAZARDOUS_CARGO":              boolLabel(app.HazardousCargo),
		"WASTE_DISCHARGE":              boolLabel(app.WasteDischarge),
		"OILY_WASTE_ONBOARD":           boolLabel(app.OilyWasteOnboard),
		"SEWAGE_DISPOSAL_REQUIRED":     boolLabel(app.SewageDisposalRequired),

		"REGULATORY_DOCUMENTS_LIST":   formatDocumentList(docs, "regulatory_document"),
		"STATUTORY_CERTIFICATES_LIST": formatDocumentList(docs, "statutory_certificate"),
		"SUPPORTING_DOCUMENTS_LIST":   formatDocumentList(docs, "supporting_document"),

		"DANGEROUS_GOOD_DECLARATION_STATUS": findDocumentStatus(docs, "dangerous_good_declaration"),
		"CARGO_MANIFEST_STATUS":             findDocumentStatus(docs, "cargo_manifest"),
		"CLEARANCE_CUSTOMS":                 findDocumentStatus(docs, "customs"),
		"CLEARANCE_IMMIGRATION":             findDocumentStatus(docs, "immigration"),
		"CLEARANCE_QUARANTINE":              findDocumentStatus(docs, "quarantine"),
	}, nil
}

// outcomeDocumentPaths holds where each generated format landed in MinIO.
// PDFPath may be empty if PDF conversion failed - the docx is still saved
// and usable even without a preview-friendly PDF.
type outcomeDocumentPaths struct {
	DocxPath string
	PDFPath  string
}

// generateOutcomeDocument renders the template for `app`, converts it to
// PDF, and uploads both to MinIO at a stable, predictable path (re-approving
// overwrites the same objects rather than piling up duplicates).
func generateOutcomeDocument(app *models.VesselApplication) (*outcomeDocumentPaths, error) {
	values, err := buildOutcomeDocumentValues(app)
	if err != nil {
		return nil, err
	}

	docBytes, err := utils.RenderDocxTemplate(outcomeTemplatePath, values)
	if err != nil {
		return nil, fmt.Errorf("failed to render outcome document: %w", err)
	}

	base := fmt.Sprintf("vessel-applications/%s/outcome/%s-outcome", app.ID, app.ApplicationNumber)
	docxPath := base + ".docx"
	docxContentType := "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	if err := utils.UploadBytes(docBytes, docxPath, docxContentType); err != nil {
		return nil, fmt.Errorf("failed to upload outcome docx: %w", err)
	}
	paths := &outcomeDocumentPaths{DocxPath: docxPath}

	pdfBytes, err := utils.ConvertDocxToPDF(docBytes)
	if err != nil {
		// docx is already saved and usable - PDF is only needed for inline
		// preview, so surface the error but don't treat it as fatal
		return paths, fmt.Errorf("docx saved but PDF conversion failed: %w", err)
	}

	pdfPath := base + ".pdf"
	if err := utils.UploadBytes(pdfBytes, pdfPath, "application/pdf"); err != nil {
		return paths, fmt.Errorf("docx saved but failed to upload PDF: %w", err)
	}
	paths.PDFPath = pdfPath

	return paths, nil
}

// GET /api/vessel-applications/:id/outcome-document/view
// Returns a fresh presigned URL to the PDF version, meant for inline
// preview in the frontend (e.g. <iframe src={url}> or <embed>) - browsers
// render PDF natively, unlike .docx.
func ViewVesselOutcomeDocument(c *gin.Context) {
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

	var path sql.NullString
	err = db.DB.QueryRow(`SELECT outcome_document_pdf_path FROM vessel_applications WHERE id = ?`, appID).Scan(&path)
	if err != nil || !path.Valid {
		c.JSON(http.StatusNotFound, gin.H{"message": "No PDF preview available for this application yet"})
		return
	}

	url, err := utils.GetFileDownloadURL(path.String)
	if err != nil {
		log.Printf("[outcome-document] failed to generate URL for %s: %v", path.String, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate download link"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
}

// GET /api/vessel-applications/:id/outcome-document/download
// Accessible by the applicant who owns the application, AND by ANP
// HSE/Admin reviewers (they need to pull this up any time from the
// review side, not just right after approving).
// Returns a fresh 1-hour presigned URL each time it's called (rather than
// embedding a link in the approval email, which would go stale after an
// hour) - same pattern as DownloadVesselDocument for regular documents.
func DownloadVesselOutcomeDocument(c *gin.Context) {
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

	var path sql.NullString
	err = db.DB.QueryRow(`SELECT outcome_document_path FROM vessel_applications WHERE id = ?`, appID).Scan(&path)
	if err != nil || !path.Valid {
		c.JSON(http.StatusNotFound, gin.H{"message": "No outcome document available for this application yet"})
		return
	}

	url, err := utils.GetFileDownloadURL(path.String)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate download link"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
}
