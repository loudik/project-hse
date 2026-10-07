package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hse-backend-go/db"
	"hse-backend-go/models"
	"hse-backend-go/utils"
)

func genApplicationNumber() string {
	ymd := time.Now().Format("20060102")
	rand.Seed(time.Now().UnixNano())
	return "VEA-" + ymd + "-" + randDigits(4)
}
func randDigits(n int) string {
	digits := "0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = digits[rand.Intn(len(digits))]
	}
	return string(b)
}

func getUserOrganizationID(userID string) (string, error) {
	var orgID sql.NullString
	err := db.DB.QueryRow(`SELECT organization_id FROM users WHERE id = ?`, userID).Scan(&orgID)
	if err != nil {
		return "", err
	}
	if !orgID.Valid {
		return "", sql.ErrNoRows
	}
	return orgID.String, nil
}

// POST /api/vessel-applications - create a new Draft (US 3.1)
// POST /api/vessel-applications - create a new Draft (US 3.1)
func CreateVesselApplication(c *gin.Context) {
	userID := c.GetString("userId")

	orgID, err := getUserOrganizationID(userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "You need an approved organization profile first"})
		return
	}

	var orgStatus string
	if err := db.DB.QueryRow(`SELECT status FROM organizations WHERE id = ?`, orgID).Scan(&orgStatus); err != nil || orgStatus != "Approved" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Your organization must be approved before creating applications"})
		return
	}

	id := uuid.NewString()
	appNumber := genApplicationNumber()

	_, err = db.DB.Exec(`
		INSERT INTO vessel_applications (id, application_number, organization_id, status, created_by)
		VALUES (?, ?, ?, 'Draft', ?)
	`, id, appNumber, orgID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create draft: " + err.Error()})
		return
	}

	app, err := fetchVesselApplication(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Draft saved but could not be read back"})
		return
	}
	c.JSON(http.StatusCreated, app)
}

// PATCH /api/vessel-applications/:id - update a Draft (US 3.1 "save and continue later")
// PATCH /api/vessel-applications/:id - partial update, dipakai di semua step wizard
func UpdateVesselApplication(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userId")

	app, err := fetchVesselApplication(id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Application not found"})
		return
	}
	if app.CreatedBy != userID {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to this application"})
		return
	}
	if app.Status != "Draft" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Only Draft applications can be edited"})
		return
	}

	var input models.VesselApplicationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid data"})
		return
	}

	setClauses := []string{}
	args := []any{}

	addStr := func(col string, val *string) {
		if val != nil {
			setClauses = append(setClauses, col+" = ?")
			args = append(args, nullableString(*val))
		}
	}
	addBool := func(col string, val *bool) {
		if val != nil {
			setClauses = append(setClauses, col+" = ?")
			args = append(args, *val)
		}
	}
	addInt := func(col string, val *int) {
		if val != nil {
			setClauses = append(setClauses, col+" = ?")
			args = append(args, *val)
		}
	}

	addStr("entry_condition", input.EntryCondition)
	addStr("purpose", input.Purpose)
	if input.NotificationEmails != nil {
		b, _ := json.Marshal(input.NotificationEmails)
		setClauses = append(setClauses, "notification_emails = ?")
		args = append(args, string(b))
	}

	addStr("contract_status", input.ContractStatus)
	addStr("contract_type", input.ContractType)
	addStr("contract_type_other", input.ContractTypeOther)
	addStr("contract_number", input.ContractNumber)
	addStr("contract_contact_email", input.ContractContactEmail)

	if input.EntryApplicationTypes != nil {
		b, _ := json.Marshal(input.EntryApplicationTypes)
		setClauses = append(setClauses, "entry_application_types = ?")
		args = append(args, string(b))
	}
	addStr("entry_date_from", input.EntryDateFrom)
	addStr("entry_date_to", input.EntryDateTo)
	addStr("entry_type", input.EntryType)

	if input.ScopeOfWork != nil {
		b, _ := json.Marshal(input.ScopeOfWork)
		setClauses = append(setClauses, "scope_of_work = ?")
		args = append(args, string(b))
	}
	addStr("scope_of_work_other", input.ScopeOfWorkOther)
	addStr("operation_mode", input.OperationMode)
	addStr("vessel_name", input.VesselName)
	addStr("vessel_imo_number", input.VesselIMONumber)
	addStr("vessel_owner", input.VesselOwner)
	addStr("vessel_type", input.VesselType)
	addStr("flag_state", input.FlagState)
	addStr("port_of_registry", input.PortOfRegistry)
	addStr("classification_society", input.ClassificationSociety)
	addStr("class_id_number", input.ClassIDNumber)
	addStr("length_overall", input.LengthOverall)
	addStr("draft_value", input.DraftValue)
	addStr("gross_tonnage", input.GrossTonnage)
	addStr("call_sign", input.CallSign)

	addBool("no_major_deficiencies", input.NoMajorDeficiencies)
	addBool("no_detention_12_months", input.NoDetention12Months)
	addBool("safety_equipment_operational", input.SafetyEquipmentOperational)
	addBool("firefighting_operational", input.FirefightingOperational)
	addBool("lifesaving_operational", input.LifesavingOperational)
	addInt("crew_count", input.CrewCount)
	addInt("survey_crew_count", input.SurveyCrewCount)
	addBool("cargo_onboard", input.CargoOnboard)
	addStr("cargo_description", input.CargoDescription)
	addBool("hazardous_cargo", input.HazardousCargo)
	addBool("waste_discharge", input.WasteDischarge)
	addBool("oily_waste_onboard", input.OilyWasteOnboard)
	addBool("sewage_disposal_required", input.SewageDisposalRequired)

	addStr("proposed_activity", input.ProposedActivity)
	addStr("planned_arrival_date", input.PlannedArrivalDate)
	addStr("planned_departure_date", input.PlannedDepartureDate)

	if len(setClauses) == 0 {
		c.JSON(http.StatusOK, app) // nothing to update
		return
	}

	query := "UPDATE vessel_applications SET " + strings.Join(setClauses, ", ") + " WHERE id = ?"
	args = append(args, id)

	if _, err := db.DB.Exec(query, args...); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update draft: " + err.Error()})
		return
	}

	updated, _ := fetchVesselApplication(id)
	c.JSON(http.StatusOK, updated)
}

// POST /api/vessel-applications/:id/submit - Draft -> Submitted (US 3.2)
func SubmitVesselApplication(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userId")

	app, err := fetchVesselApplication(id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Application not found"})
		return
	}
	if app.CreatedBy != userID {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to this application"})
		return
	}
	if app.Status != "Draft" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Only Draft applications can be submitted"})
		return
	}

	missing := []string{}
	if app.VesselName == nil || *app.VesselName == "" {
		missing = append(missing, "vessel name")
	}
	if app.VesselType == nil || *app.VesselType == "" {
		missing = append(missing, "vessel type")
	}
	if app.EntryDateFrom == nil {
		missing = append(missing, "entry date from")
	}
	if app.EntryDateTo == nil {
		missing = append(missing, "entry date to")
	}
	if len(missing) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing required fields before submission",
			"missing": missing,
		})
		return
	}

	_, err = db.DB.Exec(`
		UPDATE vessel_applications SET status = 'Submitted', submitted_at = NOW() WHERE id = ?
	`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to submit application"})
		return
	}
	utils.LogAudit("vessel_application", id, "submitted", "Draft", "Submitted", userID, "")

	vesselNameForNotif := "a vessel"
	if app.VesselName != nil {
		vesselNameForNotif = *app.VesselName
	}
	utils.NotifyReviewers(
		"New vessel entry application submitted",
		app.ApplicationNumber+" - "+vesselNameForNotif,
		"/vessel/review/"+id,
	)
	notifyAcknowledgmentDepartments(id, app.ApplicationNumber, vesselNameForNotif)

	if app.AssignedHSEOfficerID != nil && *app.AssignedHSEOfficerID != "" {
		utils.CreateNotification(
			*app.AssignedHSEOfficerID,
			"A vessel application you're assigned to was submitted",
			app.ApplicationNumber+" - "+vesselNameForNotif,
			"/vessel/review/"+id,
		)
	}
	// TODO: notify applicant + responsible HSE officer once email is fully wired for this flow

	// Generate the outcome document now too (not just at decision time) so
	// the operator has a record of exactly what they submitted, and it can
	// be attached to their confirmation email below.
	app.Status = "Submitted" // reflect the current status in the generated document
	var outcomeDocxBytes []byte
	if paths, genErr := generateOutcomeDocument(app); genErr != nil {
		log.Printf("[vessel] outcome document issue at submit for %s: %v", app.ApplicationNumber, genErr)
	} else if paths != nil && paths.DocxPath != "" {
		if _, err := db.DB.Exec(`UPDATE vessel_applications SET outcome_document_path = ? WHERE id = ?`, paths.DocxPath, id); err != nil {
			log.Printf("[vessel] failed to save outcome_document_path at submit for %s: %v", app.ApplicationNumber, err)
		}
		if paths.PDFPath != "" {
			if _, err := db.DB.Exec(`UPDATE vessel_applications SET outcome_document_pdf_path = ? WHERE id = ?`, paths.PDFPath, id); err != nil {
				log.Printf("[vessel] failed to save outcome_document_pdf_path at submit for %s: %v", app.ApplicationNumber, err)
			}
		}
		if data, err := utils.DownloadFileBytes(paths.DocxPath); err == nil {
			outcomeDocxBytes = data
		} else {
			log.Printf("[vessel] failed to re-read outcome docx for email attachment (%s): %v", app.ApplicationNumber, err)
		}
	}

	// Confirm submission to the applicant, with the outcome document attached.
	// Confirm submission to the applicant - both in-app and email, with the
	// outcome document attached to the email.
	utils.CreateNotification(
		app.CreatedBy,
		"Your vessel entry application has been submitted",
		app.ApplicationNumber+" - now awaiting review by ANP HSE",
		"/vessel/"+id,
	)

	var applicantEmail, applicantName string
	if err := db.DB.QueryRow(`SELECT email, name FROM users WHERE id = ?`, app.CreatedBy).Scan(&applicantEmail, &applicantName); err == nil {
		vesselName := "your vessel"
		if app.VesselName != nil {
			vesselName = *app.VesselName
		}
		body := fmt.Sprintf(`
			<p>Hi %s,</p>
			<p>Your vessel entry application <strong>%s</strong> for <strong>%s</strong> has been
			submitted and is now awaiting review by ANP HSE.</p>
			<p>A copy of your submission is attached for your records.</p>
		`, applicantName, app.ApplicationNumber, vesselName)

		var attachments []utils.EmailAttachment
		if outcomeDocxBytes != nil {
			attachments = append(attachments, utils.EmailAttachment{
				Filename:    app.ApplicationNumber + "-outcome.docx",
				ContentType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
				Data:        outcomeDocxBytes,
			})
		}
		if emailErr := utils.SendEmailWithAttachments(applicantEmail, "Vessel entry application submitted", body, attachments); emailErr != nil {
			log.Printf("[vessel] failed to send submission confirmation email to %s: %v", applicantEmail, emailErr)
		}
	} else {
		log.Printf("[vessel] failed to look up applicant email for %s: %v", app.ApplicationNumber, err)
	}

	updated, err := fetchVesselApplication(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Submitted but could not read back the application"})
		return
	}
	c.JSON(http.StatusOK, updated)

}

// GET /api/vessel-applications - list current user's organization's applications
// GET /api/vessel-applications - list current user's organization's applications
func ListMyVesselApplications(c *gin.Context) {
	userID := c.GetString("userId")

	orgID, err := getUserOrganizationID(userID)
	if err != nil {
		c.JSON(http.StatusOK, []models.VesselApplication{})
		return
	}

	rows, err := db.DB.Query(`
		SELECT id, application_number, organization_id, status, created_by, submitted_at, created_at, updated_at,
					 outcome_document_path, outcome_document_pdf_path,assigned_hse_officer_id,substituted_from_application_id,
		       entry_condition, purpose, entry_date_from, entry_date_to, entry_type,notification_emails, contract_status, contract_type, contract_type_other, contract_number, contract_contact_email,
					 entry_application_types,
		       scope_of_work, scope_of_work_other, operation_mode,
		       vessel_name, vessel_imo_number, vessel_owner, vessel_type, flag_state, port_of_registry,
		       classification_society, class_id_number, length_overall, draft_value, gross_tonnage, call_sign,
		       no_major_deficiencies, no_detention_12_months, safety_equipment_operational,
		       firefighting_operational, lifesaving_operational, crew_count, survey_crew_count,
		       cargo_onboard, cargo_description, hazardous_cargo, waste_discharge, oily_waste_onboard, sewage_disposal_required,
		       proposed_activity, planned_arrival_date, planned_departure_date
		FROM vessel_applications WHERE organization_id = ?
		ORDER BY created_at DESC
	`, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch applications"})
		return
	}
	defer rows.Close()

	list := []models.VesselApplication{}
	for rows.Next() {
		app, err := scanVesselApplication(rows)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read application data"})
			return
		}
		list = append(list, *app)
	}
	c.JSON(http.StatusOK, list)
}

// GET /api/vessel-applications/:id
func GetVesselApplication(c *gin.Context) {
	userID := c.GetString("userId")
	app, err := fetchVesselApplication(c.Param("id"))
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Application not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch application"})
		return
	}
	if app.CreatedBy != userID {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to this application"})
		return
	}
	c.JSON(http.StatusOK, app)
}

// --- helpers ---

func nullableDate(s string) any {
	if s == "" {
		return nil
	}
	return s
}

type vesselRowScanner interface {
	Scan(dest ...any) error
}

func fetchVesselApplication(id string) (*models.VesselApplication, error) {
	row := db.DB.QueryRow(`
		SELECT id, application_number, organization_id, status, created_by, submitted_at, created_at, updated_at,outcome_document_path,outcome_document_pdf_path,assigned_hse_officer_id,substituted_from_application_id,
		       entry_condition, purpose, entry_date_from, entry_date_to, entry_type,
					 notification_emails, contract_status, contract_type, contract_type_other, contract_number, contract_contact_email,
					entry_application_types,
		       scope_of_work, scope_of_work_other, operation_mode,
		       vessel_name, vessel_imo_number, vessel_owner, vessel_type, flag_state, port_of_registry,
		       classification_society, class_id_number, length_overall, draft_value, gross_tonnage, call_sign,
		       no_major_deficiencies, no_detention_12_months, safety_equipment_operational,
		       firefighting_operational, lifesaving_operational, crew_count, survey_crew_count,
		       cargo_onboard, cargo_description, hazardous_cargo, waste_discharge, oily_waste_onboard, sewage_disposal_required,
		       proposed_activity, planned_arrival_date, planned_departure_date
		FROM vessel_applications WHERE id = ?
	`, id)
	return scanVesselApplication(row)
}

func scanVesselApplication(row vesselRowScanner) (*models.VesselApplication, error) {
	var a models.VesselApplication
	var submittedAt sql.NullTime
	var outcomeDocumentPath, outcomeDocumentPDFPath sql.NullString
	var assignedHSEOfficerID sql.NullString
	var substitutedFromApplicationID sql.NullString
	var entryCondition, purpose, entryType sql.NullString
	var notificationEmailsRaw sql.NullString
	var contractStatus, contractType, contractTypeOther, contractNumber, contractContactEmail sql.NullString
	var entryAppTypesRaw sql.NullString
	var entryDateFrom, entryDateTo sql.NullTime
	var scopeOfWorkRaw sql.NullString
	var scopeOfWorkOther, operationMode sql.NullString
	var vesselName, vesselIMO, vesselOwner, vesselType, flagState, portOfRegistry sql.NullString
	var classSociety, classID, lengthOverall, draftValue, grossTonnage, callSign sql.NullString
	var noMajorDef, noDetention, safetyOp, firefightingOp, lifesavingOp sql.NullBool
	var crewCount, surveyCrewCount sql.NullInt64
	var cargoOnboard, hazardousCargo, wasteDischarge, oilyWaste, sewageDisposal sql.NullBool
	var cargoDescription sql.NullString
	var proposedActivity sql.NullString
	var plannedArrival, plannedDeparture sql.NullTime

	err := row.Scan(
		&a.ID, &a.ApplicationNumber, &a.OrganizationID, &a.Status, &a.CreatedBy, &submittedAt, &a.CreatedAt, &a.UpdatedAt, &outcomeDocumentPath, &outcomeDocumentPDFPath, &assignedHSEOfficerID, &substitutedFromApplicationID,
		&entryCondition, &purpose, &entryDateFrom, &entryDateTo, &entryType, &notificationEmailsRaw, &contractStatus, &contractType, &contractTypeOther, &contractNumber, &contractContactEmail,
		&entryAppTypesRaw,
		&scopeOfWorkRaw, &scopeOfWorkOther, &operationMode,
		&vesselName, &vesselIMO, &vesselOwner, &vesselType, &flagState, &portOfRegistry,
		&classSociety, &classID, &lengthOverall, &draftValue, &grossTonnage, &callSign,
		&noMajorDef, &noDetention, &safetyOp, &firefightingOp, &lifesavingOp, &crewCount, &surveyCrewCount,
		&cargoOnboard, &cargoDescription, &hazardousCargo, &wasteDischarge, &oilyWaste, &sewageDisposal,
		&proposedActivity, &plannedArrival, &plannedDeparture,
	)
	if err != nil {
		return nil, err
	}

	strPtr := func(ns sql.NullString) *string {
		if ns.Valid {
			return &ns.String
		}
		return nil
	}
	boolPtr := func(nb sql.NullBool) *bool {
		if nb.Valid {
			return &nb.Bool
		}
		return nil
	}
	intPtr := func(ni sql.NullInt64) *int {
		if ni.Valid {
			v := int(ni.Int64)
			return &v
		}
		return nil
	}
	datePtr := func(nt sql.NullTime) *string {
		if nt.Valid {
			s := nt.Time.Format("2006-01-02")
			return &s
		}
		return nil
	}

	if submittedAt.Valid {
		a.SubmittedAt = &submittedAt.Time
	}

	if outcomeDocumentPath.Valid { // <-- baris baru
		a.OutcomeDocumentPath = &outcomeDocumentPath.String // <-- baris baru
	}

	if outcomeDocumentPDFPath.Valid { // <-- baris baru
		a.OutcomeDocumentPDFPath = &outcomeDocumentPDFPath.String // <-- baris baru
	}

	if assignedHSEOfficerID.Valid {
		a.AssignedHSEOfficerID = &assignedHSEOfficerID.String
	}
	if substitutedFromApplicationID.Valid {
		a.SubstitutedFromApplicationID = &substitutedFromApplicationID.String
	}

	a.EntryCondition = strPtr(entryCondition)
	a.Purpose = strPtr(purpose)
	if notificationEmailsRaw.Valid && notificationEmailsRaw.String != "" {
		_ = json.Unmarshal([]byte(notificationEmailsRaw.String), &a.NotificationEmails)
	}
	a.ContractStatus = strPtr(contractStatus)
	a.ContractType = strPtr(contractType)
	a.ContractTypeOther = strPtr(contractTypeOther)
	a.ContractNumber = strPtr(contractNumber)
	a.ContractContactEmail = strPtr(contractContactEmail)
	if entryAppTypesRaw.Valid && entryAppTypesRaw.String != "" {
		_ = json.Unmarshal([]byte(entryAppTypesRaw.String), &a.EntryApplicationTypes)
	}
	a.EntryDateFrom = datePtr(entryDateFrom)
	a.EntryDateTo = datePtr(entryDateTo)
	a.EntryType = strPtr(entryType)

	if scopeOfWorkRaw.Valid && scopeOfWorkRaw.String != "" {
		_ = json.Unmarshal([]byte(scopeOfWorkRaw.String), &a.ScopeOfWork)
	}
	a.ScopeOfWorkOther = strPtr(scopeOfWorkOther)
	a.OperationMode = strPtr(operationMode)
	a.VesselName = strPtr(vesselName)
	a.VesselIMONumber = strPtr(vesselIMO)
	a.VesselOwner = strPtr(vesselOwner)
	a.VesselType = strPtr(vesselType)
	a.FlagState = strPtr(flagState)
	a.PortOfRegistry = strPtr(portOfRegistry)
	a.ClassificationSociety = strPtr(classSociety)
	a.ClassIDNumber = strPtr(classID)
	a.LengthOverall = strPtr(lengthOverall)
	a.DraftValue = strPtr(draftValue)
	a.GrossTonnage = strPtr(grossTonnage)
	a.CallSign = strPtr(callSign)

	a.NoMajorDeficiencies = boolPtr(noMajorDef)
	a.NoDetention12Months = boolPtr(noDetention)
	a.SafetyEquipmentOperational = boolPtr(safetyOp)
	a.FirefightingOperational = boolPtr(firefightingOp)
	a.LifesavingOperational = boolPtr(lifesavingOp)
	a.CrewCount = intPtr(crewCount)
	a.SurveyCrewCount = intPtr(surveyCrewCount)
	a.CargoOnboard = boolPtr(cargoOnboard)
	a.CargoDescription = strPtr(cargoDescription)
	a.HazardousCargo = boolPtr(hazardousCargo)
	a.WasteDischarge = boolPtr(wasteDischarge)
	a.OilyWasteOnboard = boolPtr(oilyWaste)
	a.SewageDisposalRequired = boolPtr(sewageDisposal)

	a.ProposedActivity = strPtr(proposedActivity)
	a.PlannedArrivalDate = datePtr(plannedArrival)
	a.PlannedDepartureDate = datePtr(plannedDeparture)

	return &a, nil
}

// PATCH /api/vessel-applications/:id/assign - Admin/ANP HSE only
// Body: { "hseOfficerId": "..." } or { "hseOfficerId": null } to unassign.
// The assigned HSE Officer is notified in parallel at each stage
// (submit/decision/expiry reminders) as the coordinator for this
// application - they don't approve/reject themselves.
func AssignHSEOfficer(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		HSEOfficerID *string `json:"hseOfficerId"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	if _, err := fetchVesselApplication(id); err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Application not found"})
		return
	}

	if input.HSEOfficerID != nil && *input.HSEOfficerID != "" {
		var roleName string
		err := db.DB.QueryRow(`
			SELECT r.name FROM users u JOIN roles r ON r.id = u.role_id WHERE u.id = ?
		`, *input.HSEOfficerID).Scan(&roleName)
		if err != nil || (roleName != "HSE Officer" && roleName != "Staff") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "That user must be an HSE Officer or Staff"})
			return
		}
	}

	if _, err := db.DB.Exec(`UPDATE vessel_applications SET assigned_hse_officer_id = ? WHERE id = ?`, input.HSEOfficerID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to assign HSE Officer"})
		return
	}

	if input.HSEOfficerID != nil && *input.HSEOfficerID != "" {
		app, _ := fetchVesselApplication(id)
		vesselName := "a vessel"
		if app != nil && app.VesselName != nil {
			vesselName = *app.VesselName
		}
		utils.CreateNotification(
			*input.HSEOfficerID,
			"You've been assigned to a vessel application",
			app.ApplicationNumber+" - "+vesselName,
			"/vessel/review/"+id,
		)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assigned"})
}

// GET /api/vessel-applications/hse-officers - Admin/ANP HSE only
// For the assignment dropdown.
func ListHSEOfficers(c *gin.Context) {
	rows, err := db.DB.Query(`
		SELECT u.id, u.name FROM users u JOIN roles r ON r.id = u.role_id
		WHERE r.name IN ('HSE Officer', 'Staff') AND u.status = 'Active'
		ORDER BY u.name
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch HSE Officers"})
		return
	}
	defer rows.Close()

	type officer struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	list := []officer{}
	for rows.Next() {
		var o officer
		if err := rows.Scan(&o.ID, &o.Name); err == nil {
			list = append(list, o)
		}
	}
	c.JSON(http.StatusOK, list)
}

// GET /api/vessel-applications/:id/history
func GetVesselApplicationHistory(c *gin.Context) {
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

	history, err := utils.FetchAuditLog("vessel_application", appID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch history"})
		return
	}
	c.JSON(http.StatusOK, history)
}

// GET /api/vessel-applications/lookup?number=XXX
func LookupVesselApplication(c *gin.Context) {
	userID := c.GetString("userId")
	number := c.Query("number")
	if number == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Application number is required"})
		return
	}

	orgID, err := getUserOrganizationID(userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "You need an approved organization profile first"})
		return
	}

	row := db.DB.QueryRow(`
		SELECT id, application_number, organization_id, vessel_name, vessel_type,
		       proposed_activity, planned_arrival_date, planned_departure_date,
		       status, created_by, submitted_at, created_at, updated_at
		FROM vessel_applications WHERE application_number = ? AND organization_id = ?
	`, number, orgID)
	app, err := scanVesselApplication(row)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "No application found with that number for your organization"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to look up application"})
		return
	}
	c.JSON(http.StatusOK, app)
}

// POST /api/vessel-applications/:id/withdraw
func WithdrawVesselApplication(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userId")

	app, err := fetchVesselApplication(id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Application not found"})
		return
	}
	if app.CreatedBy != userID {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to this application"})
		return
	}
	if app.Status != "Submitted" && app.Status != "Approved" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Only Submitted or Approved applications can be withdrawn"})
		return
	}

	_, err = db.DB.Exec(`UPDATE vessel_applications SET status = 'Withdrawn' WHERE id = ?`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to withdraw application"})
		return
	}
	utils.LogAudit("vessel_application", id, "withdrawn", app.Status, "Withdrawn", userID, "")
	updated, _ := fetchVesselApplication(id)
	c.JSON(http.StatusOK, updated)
}

// POST /api/vessel-applications/:id/reopen
func ReopenVesselApplication(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userId")

	app, err := fetchVesselApplication(id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Application not found"})
		return
	}
	if app.CreatedBy != userID {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to this application"})
		return
	}
	if app.Status != "Approved" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Only Approved applications can be updated"})
		return
	}

	_, err = db.DB.Exec(`UPDATE vessel_applications SET status = 'Draft' WHERE id = ?`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to reopen application"})
		return
	}
	updated, _ := fetchVesselApplication(id)
	c.JSON(http.StatusOK, updated)
}

// GET /api/vessel-applications/review - ANP HSE/Admin only, lihat SEMUA aplikasi (bukan cuma milik sendiri)
func ListVesselApplicationsForReview(c *gin.Context) {
	statusFilter := c.Query("status")

	query := `
		SELECT id, application_number, organization_id, status, created_by, submitted_at, created_at, updated_at,
		       outcome_document_path, outcome_document_pdf_path,assigned_hse_officer_id,substituted_from_application_id,
		       entry_condition, purpose, entry_date_from, entry_date_to, entry_type,notification_emails, contract_status, contract_type, contract_type_other, contract_number, contract_contact_email,
					 entry_application_types,
		       scope_of_work, scope_of_work_other, operation_mode,
		       vessel_name, vessel_imo_number, vessel_owner, vessel_type, flag_state, port_of_registry,
		       classification_society, class_id_number, length_overall, draft_value, gross_tonnage, call_sign,
		       no_major_deficiencies, no_detention_12_months, safety_equipment_operational,
		       firefighting_operational, lifesaving_operational, crew_count, survey_crew_count,
		       cargo_onboard, cargo_description, hazardous_cargo, waste_discharge, oily_waste_onboard, sewage_disposal_required,
		       proposed_activity, planned_arrival_date, planned_departure_date
		FROM vessel_applications
	`
	args := []any{}
	if statusFilter != "" {
		query += ` WHERE status = ?`
		args = append(args, statusFilter)
	} else {
		query += ` WHERE status != 'Draft'` // reviewer doesn't need to see drafts
	}
	query += ` ORDER BY submitted_at DESC`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch applications"})
		return
	}
	defer rows.Close()

	list := []models.VesselApplication{}
	for rows.Next() {
		app, err := scanVesselApplication(rows)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read application data"})
			return
		}
		list = append(list, *app)
	}
	c.JSON(http.StatusOK, list)
}

// GET /api/vessel-applications/review/:id - ANP HSE/Admin, plus the
// assigned HSE Officer and any Staff specifically asked to review this one.
func GetVesselApplicationForReview(c *gin.Context) {
	userID := c.GetString("userId")
	roleName := c.GetString("roleName")

	app, err := fetchVesselApplication(c.Param("id"))
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Application not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch application"})
		return
	}
	if !canCollaborateOnApplication(app, userID, roleName) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to this application"})
		return
	}
	c.JSON(http.StatusOK, app)
}

// PATCH /api/vessel-applications/:id/decision - ANP HSE/Admin only
func DecideVesselApplication(c *gin.Context) {
	id := c.Param("id")
	reviewerID := c.GetString("userId")

	var input struct {
		Status string `json:"status" binding:"required,oneof=Approved Rejected"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "status must be Approved or Rejected"})
		return
	}

	app, err := fetchVesselApplication(id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Application not found"})
		return
	}
	if app.Status != "Submitted" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Only Submitted applications can be decided"})
		return
	}
	if app.AssignedHSEOfficerID == nil || *app.AssignedHSEOfficerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please assign an HSE Officer before making a decision"})
		return
	}

	if _, err := db.DB.Exec(`UPDATE vessel_applications SET status = ?, decided_at = NOW() WHERE id = ?`, input.Status, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save decision"})
		return
	}
	utils.LogAudit("vessel_application", id, strings.ToLower(input.Status), "Submitted", input.Status, reviewerID, "")
	utils.CreateNotification(
		app.CreatedBy,
		"Your vessel entry application has been "+strings.ToLower(input.Status),
		app.ApplicationNumber,
		"/vessel/"+id,
	)
	if app.AssignedHSEOfficerID != nil && *app.AssignedHSEOfficerID != "" {
		utils.CreateNotification(
			*app.AssignedHSEOfficerID,
			"A vessel application you're assigned to was "+strings.ToLower(input.Status),
			app.ApplicationNumber,
			"/vessel/review/"+id,
		)
	}
	documentGenerated := false
	if input.Status == "Approved" {
		app.Status = "Approved" // reflect the new status in the generated document
		paths, genErr := generateOutcomeDocument(app)
		if genErr != nil {
			log.Printf("[vessel] outcome document issue for %s: %v", app.ApplicationNumber, genErr)
		}
		if paths != nil && paths.DocxPath != "" {
			if _, err := db.DB.Exec(`UPDATE vessel_applications SET outcome_document_path = ? WHERE id = ?`, paths.DocxPath, id); err != nil {
				log.Printf("[vessel] failed to save outcome_document_path for %s: %v", app.ApplicationNumber, err)
			} else {
				documentGenerated = true
			}
		}
		if paths != nil && paths.PDFPath != "" {
			if _, err := db.DB.Exec(`UPDATE vessel_applications SET outcome_document_pdf_path = ? WHERE id = ?`, paths.PDFPath, id); err != nil {
				log.Printf("[vessel] failed to save outcome_document_pdf_path for %s: %v", app.ApplicationNumber, err)
			}
		}
	}

	// Notify the applicant by email
	var applicantEmail, applicantName string
	if err := db.DB.QueryRow(`SELECT email, name FROM users WHERE id = ?`, app.CreatedBy).Scan(&applicantEmail, &applicantName); err == nil {
		vesselName := "your vessel"
		if app.VesselName != nil {
			vesselName = *app.VesselName
		}
		subject := "Vessel entry application " + strings.ToLower(input.Status)
		outcomeNote := ""
		if documentGenerated {
			outcomeNote = "<p>The outcome document is available in the system - log in to view or download it.</p>"
		}
		body := fmt.Sprintf(`
			<p>Hi %s,</p>
			<p>Your vessel entry application <strong>%s</strong> for <strong>%s</strong> has been
			<strong>%s</strong> by ANP HSE.</p>
			%s
			<p>You can log in to view the full details.</p>
		`, applicantName, app.ApplicationNumber, vesselName, input.Status, outcomeNote)
		_ = utils.SendEmail(applicantEmail, subject, body)
	}
	updated, err := fetchVesselApplication(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Decision saved but could not read back the application"})
		return
	}
	c.JSON(http.StatusOK, updated)
}
