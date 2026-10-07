package models

import "time"

type VesselApplication struct {
	ID                string     `json:"id"`
	ApplicationNumber string     `json:"applicationNumber"`
	OrganizationID    string     `json:"organizationId"`
	Status            string     `json:"status"`
	CreatedBy         string     `json:"createdBy"`
	SubmittedAt       *time.Time `json:"submittedAt"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`

	OutcomeDocumentPath          *string `json:"outcomeDocumentPath"`
	OutcomeDocumentPDFPath       *string `json:"outcomeDocumentPdfPath"`
	AssignedHSEOfficerID         *string `json:"assignedHseOfficerId"`
	SubstitutedFromApplicationID *string `json:"substitutedFromApplicationId"`

	// Vessel-2: condition & entry info
	EntryCondition     *string  `json:"entryCondition"`
	Purpose            *string  `json:"purpose"`
	NotificationEmails []string `json:"notificationEmails"`

	ContractStatus       *string `json:"contractStatus"`
	ContractType         *string `json:"contractType"`
	ContractTypeOther    *string `json:"contractTypeOther"`
	ContractNumber       *string `json:"contractNumber"`
	ContractContactEmail *string `json:"contractContactEmail"`

	EntryApplicationTypes []string `json:"entryApplicationTypes"`

	EntryDateFrom *string `json:"entryDateFrom"`
	EntryDateTo   *string `json:"entryDateTo"`
	EntryType     *string `json:"entryType"`

	// Vessel 3: scope & vessel description
	ScopeOfWork           []string `json:"scopeOfWork"`
	ScopeOfWorkOther      *string  `json:"scopeOfWorkOther"`
	OperationMode         *string  `json:"operationMode"`
	VesselName            *string  `json:"vesselName"`
	VesselIMONumber       *string  `json:"vesselImoNumber"`
	VesselOwner           *string  `json:"vesselOwner"`
	VesselType            *string  `json:"vesselType"`
	FlagState             *string  `json:"flagState"`
	PortOfRegistry        *string  `json:"portOfRegistry"`
	ClassificationSociety *string  `json:"classificationSociety"`
	ClassIDNumber         *string  `json:"classIdNumber"`
	LengthOverall         *string  `json:"lengthOverall"`
	DraftValue            *string  `json:"draftValue"`
	GrossTonnage          *string  `json:"grossTonnage"`
	CallSign              *string  `json:"callSign"`

	// Vessel-5: declarations
	NoMajorDeficiencies        *bool   `json:"noMajorDeficiencies"`
	NoDetention12Months        *bool   `json:"noDetention12Months"`
	SafetyEquipmentOperational *bool   `json:"safetyEquipmentOperational"`
	FirefightingOperational    *bool   `json:"firefightingOperational"`
	LifesavingOperational      *bool   `json:"lifesavingOperational"`
	CrewCount                  *int    `json:"crewCount"`
	SurveyCrewCount            *int    `json:"surveyCrewCount"`
	CargoOnboard               *bool   `json:"cargoOnboard"`
	CargoDescription           *string `json:"cargoDescription"`
	HazardousCargo             *bool   `json:"hazardousCargo"`
	WasteDischarge             *bool   `json:"wasteDischarge"`
	OilyWasteOnboard           *bool   `json:"oilyWasteOnboard"`
	SewageDisposalRequired     *bool   `json:"sewageDisposalRequired"`

	// legacy fields kept for backward compat with the original 3-step wizard
	ProposedActivity     *string `json:"proposedActivity"`
	PlannedArrivalDate   *string `json:"plannedArrivalDate"`
	PlannedDepartureDate *string `json:"plannedDepartureDate"`
}

type VesselApplicationInput struct {
	EntryCondition     *string  `json:"entryCondition"`
	Purpose            *string  `json:"purpose"`
	NotificationEmails []string `json:"notificationEmails"`

	ContractStatus       *string `json:"contractStatus"`
	ContractType         *string `json:"contractType"`
	ContractTypeOther    *string `json:"contractTypeOther"`
	ContractNumber       *string `json:"contractNumber"`
	ContractContactEmail *string `json:"contractContactEmail"`

	EntryApplicationTypes []string `json:"entryApplicationTypes"`
	EntryDateFrom         *string  `json:"entryDateFrom"`
	EntryDateTo           *string  `json:"entryDateTo"`
	EntryType             *string  `json:"entryType"`

	ScopeOfWork           []string `json:"scopeOfWork"`
	ScopeOfWorkOther      *string  `json:"scopeOfWorkOther"`
	OperationMode         *string  `json:"operationMode"`
	VesselName            *string  `json:"vesselName"`
	VesselIMONumber       *string  `json:"vesselImoNumber"`
	VesselOwner           *string  `json:"vesselOwner"`
	VesselType            *string  `json:"vesselType"`
	FlagState             *string  `json:"flagState"`
	PortOfRegistry        *string  `json:"portOfRegistry"`
	ClassificationSociety *string  `json:"classificationSociety"`
	ClassIDNumber         *string  `json:"classIdNumber"`
	LengthOverall         *string  `json:"lengthOverall"`
	DraftValue            *string  `json:"draftValue"`
	GrossTonnage          *string  `json:"grossTonnage"`
	CallSign              *string  `json:"callSign"`

	NoMajorDeficiencies        *bool   `json:"noMajorDeficiencies"`
	NoDetention12Months        *bool   `json:"noDetention12Months"`
	SafetyEquipmentOperational *bool   `json:"safetyEquipmentOperational"`
	FirefightingOperational    *bool   `json:"firefightingOperational"`
	LifesavingOperational      *bool   `json:"lifesavingOperational"`
	CrewCount                  *int    `json:"crewCount"`
	SurveyCrewCount            *int    `json:"surveyCrewCount"`
	CargoOnboard               *bool   `json:"cargoOnboard"`
	CargoDescription           *string `json:"cargoDescription"`
	HazardousCargo             *bool   `json:"hazardousCargo"`
	WasteDischarge             *bool   `json:"wasteDischarge"`
	OilyWasteOnboard           *bool   `json:"oilyWasteOnboard"`
	SewageDisposalRequired     *bool   `json:"sewageDisposalRequired"`

	ProposedActivity     *string `json:"proposedActivity"`
	PlannedArrivalDate   *string `json:"plannedArrivalDate"`
	PlannedDepartureDate *string `json:"plannedDepartureDate"`
}
