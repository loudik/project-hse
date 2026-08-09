package models

import "time"

type Organization struct {
	ID                 string     `json:"id"`
	Name               string     `json:"name"`
	RegistrationNumber string     `json:"registrationNumber"`
	Type               string     `json:"type"`
	Address            *string    `json:"address"`
	Country            *string    `json:"country"`
	PhoneNumber        *string    `json:"phoneNumber"`
	Email              *string    `json:"email"`
	Website            *string    `json:"website"`
	Status             string     `json:"status"` // Pending | Approved | Rejected
	CreatedBy          string     `json:"createdBy"`
	ApprovedBy         *string    `json:"approvedBy"`
	ApprovedAt         *time.Time `json:"approvedAt"`
	RejectionReason    *string    `json:"rejectionReason"`
	CreatedAt          time.Time  `json:"createdAt"`
}

// CreateOrganizationInput - payload from "Create an organisation profile" (US 2.1)
type CreateOrganizationInput struct {
	Name               string `json:"name" binding:"required"`
	RegistrationNumber string `json:"registrationNumber" binding:"required"`
	Type               string `json:"type" binding:"required"`
	Address            string `json:"address"`
	Country            string `json:"country"`
	PhoneNumber        string `json:"phoneNumber"`
	Email              string `json:"email"`
	Website            string `json:"website"`
}

// OrganizationDecisionInput - ANP HSE approve/reject
type OrganizationDecisionInput struct {
	Status           string `json:"status" binding:"required,oneof=Approved Rejected"`
	RejectionReason  string `json:"rejectionReason"`
}
