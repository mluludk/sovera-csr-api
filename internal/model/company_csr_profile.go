package model

import "time"

type CompanyCSRProfile struct {
	ID                 string     `json:"id"`
	CompanyID          string     `json:"company_id"`
	HasCSR             bool       `json:"has_csr"`
	CSRDepartmentName  *string    `json:"csr_department_name,omitempty"`
	CSREmailPublic     *string    `json:"csr_email_public,omitempty"`
	CSRFocus           []string   `json:"csr_focus"`
	BudgetRange        *string    `json:"budget_range,omitempty"`
	ProposalAcceptance *string    `json:"proposal_acceptance,omitempty"`
	WebsiteSource      *string    `json:"website_source,omitempty"`
	LastVerifiedAt     *time.Time `json:"last_verified_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}
