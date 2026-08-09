package models

import "time"

// HseSubmission - full detail of a single submission
type HseSubmission struct {
	ID                 string     `json:"id"`
	CertNumber         string     `json:"certNumber"`
	ApplicantName      string     `json:"applicantName"`
	ApplicantPosition  *string    `json:"applicantPosition"`
	Department         string     `json:"department"`
	WorkType           string     `json:"workType"`
	Location           string     `json:"location"`
	StartDate          string     `json:"startDate"` // format YYYY-MM-DD
	EndDate            string     `json:"endDate"`
	Description        string     `json:"description"`
	Hazards            string     `json:"hazards"`
	Controls           string     `json:"controls"`
	Photos             []string   `json:"photos"`
	ApplicantSignature string     `json:"applicantSignature"`
	Status             string     `json:"status"`
	ApproverName       *string    `json:"approverName"`
	ApproverNote       *string    `json:"approverNote"`
	ApproverSignature  *string    `json:"approverSignature"`
	SubmittedAt        time.Time  `json:"submittedAt"`
	DecidedAt          *time.Time `json:"decidedAt"`
}

// SubmissionSummary - for list/dashboard views (lighter payload)
type SubmissionSummary struct {
	ID            string    `json:"id"`
	CertNumber    string    `json:"certNumber"`
	ApplicantName string    `json:"applicantName"`
	WorkType      string    `json:"workType"`
	Status        string    `json:"status"`
	SubmittedAt   time.Time `json:"submittedAt"`
}

// CreateSubmissionInput - payload from the new-submission form
type CreateSubmissionInput struct {
	ApplicantName      string   `json:"applicantName" binding:"required"`
	ApplicantPosition  string   `json:"applicantPosition"`
	Department         string   `json:"department" binding:"required"`
	WorkType           string   `json:"workType" binding:"required"`
	Location           string   `json:"location" binding:"required"`
	StartDate          string   `json:"startDate" binding:"required"`
	EndDate            string   `json:"endDate" binding:"required"`
	Description        string   `json:"description" binding:"required"`
	Hazards            string   `json:"hazards" binding:"required"`
	Controls           string   `json:"controls" binding:"required"`
	Photos             []string `json:"photos"`
	ApplicantSignature string   `json:"applicantSignature" binding:"required"`
}

// DecisionInput - approve/reject payload
type DecisionInput struct {
	Status            string `json:"status" binding:"required,oneof=Approved Rejected"`
	ApproverName      string `json:"approverName" binding:"required"`
	ApproverNote      string `json:"approverNote"`
	ApproverSignature string `json:"approverSignature" binding:"required"`
}
