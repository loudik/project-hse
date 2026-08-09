package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hse-backend-go/db"
	"hse-backend-go/models"
)

func genCertNumber() string {
	ymd := time.Now().Format("20060102")
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("HSE-CERT-%s-%04d", ymd, rand.Intn(9000)+1000)
}

// GET /api/hse/submissions
func ListSubmissions(c *gin.Context) {
	rows, err := db.DB.Query(`
		SELECT id, cert_number, applicant_name, work_type, status, submitted_at
		FROM hse_submissions
		ORDER BY submitted_at ASC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch data"})
		return
	}
	defer rows.Close()

	list := []models.SubmissionSummary{}
	for rows.Next() {
		var s models.SubmissionSummary
		if err := rows.Scan(&s.ID, &s.CertNumber, &s.ApplicantName, &s.WorkType, &s.Status, &s.SubmittedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read data"})
			return
		}
		list = append(list, s)
	}

	c.JSON(http.StatusOK, list)
}

// GET /api/hse/submissions/:id
func GetSubmission(c *gin.Context) {
	id := c.Param("id")
	record, err := fetchSubmission(id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Submission not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch data"})
		return
	}
	c.JSON(http.StatusOK, record)
}

// POST /api/hse/submissions
func CreateSubmission(c *gin.Context) {
	var input models.CreateSubmissionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incomplete or invalid data: " + err.Error()})
		return
	}

	id := uuid.NewString()
	certNumber := genCertNumber()

	photosJSON, _ := json.Marshal(input.Photos)

	_, err := db.DB.Exec(`
		INSERT INTO hse_submissions
			(id, cert_number, applicant_name, applicant_position, department, work_type,
			 location, start_date, end_date, description, hazards, controls,
			 photos, applicant_signature, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'Submitted')
	`,
		id, certNumber, input.ApplicantName, nullableString(input.ApplicantPosition), input.Department, input.WorkType,
		input.Location, input.StartDate, input.EndDate, input.Description, input.Hazards, input.Controls,
		string(photosJSON), input.ApplicantSignature,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save submission: " + err.Error()})
		return
	}

	record, err := fetchSubmission(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Submission saved but could not be read back"})
		return
	}
	c.JSON(http.StatusCreated, record)
}

// PATCH /api/hse/submissions/:id/decision
func DecideSubmission(c *gin.Context) {
	id := c.Param("id")

	var input models.DecisionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid data: " + err.Error()})
		return
	}

	_, err := fetchSubmission(id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Submission not found"})
		return
	}

	_, err = db.DB.Exec(`
		UPDATE hse_submissions
		SET status = ?, approver_name = ?, approver_note = ?, approver_signature = ?, decided_at = NOW()
		WHERE id = ?
	`, input.Status, input.ApproverName, nullableString(input.ApproverNote), input.ApproverSignature, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save decision"})
		return
	}

	record, err := fetchSubmission(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read latest data"})
		return
	}
	c.JSON(http.StatusOK, record)
}

// --- helpers ---

func fetchSubmission(id string) (*models.HseSubmission, error) {
	row := db.DB.QueryRow(`
		SELECT id, cert_number, applicant_name, applicant_position, department, work_type,
		       location, start_date, end_date, description, hazards, controls,
		       photos, applicant_signature, status,
		       approver_name, approver_note, approver_signature,
		       submitted_at, decided_at
		FROM hse_submissions WHERE id = ?
	`, id)

	var s models.HseSubmission
	var applicantPosition, approverName, approverNote, approverSignature sql.NullString
	var photosRaw sql.NullString
	var decidedAt sql.NullTime
	var startDate, endDate time.Time

	err := row.Scan(
		&s.ID, &s.CertNumber, &s.ApplicantName, &applicantPosition, &s.Department, &s.WorkType,
		&s.Location, &startDate, &endDate, &s.Description, &s.Hazards, &s.Controls,
		&photosRaw, &s.ApplicantSignature, &s.Status,
		&approverName, &approverNote, &approverSignature,
		&s.SubmittedAt, &decidedAt,
	)
	if err != nil {
		return nil, err
	}

	s.StartDate = startDate.Format("2006-01-02")
	s.EndDate = endDate.Format("2006-01-02")
	if applicantPosition.Valid {
		s.ApplicantPosition = &applicantPosition.String
	}
	if approverName.Valid {
		s.ApproverName = &approverName.String
	}
	if approverNote.Valid {
		s.ApproverNote = &approverNote.String
	}
	if approverSignature.Valid {
		s.ApproverSignature = &approverSignature.String
	}
	if decidedAt.Valid {
		s.DecidedAt = &decidedAt.Time
	}
	if photosRaw.Valid && photosRaw.String != "" {
		_ = json.Unmarshal([]byte(photosRaw.String), &s.Photos)
	} else {
		s.Photos = []string{}
	}

	return &s, nil
}

func nullableString(s string) any {
	if s == "" {
		return nil
	}
	return s
}
