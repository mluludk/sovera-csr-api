package model

import "time"

type CompanyCSRFocus struct {
	CompanyID  string     `json:"company_id"`
	FocusID    string     `json:"focus_id"`
	Priority   *int16     `json:"priority,omitempty"`
	Confidence *float64   `json:"confidence,omitempty"`
	SourceID   *string    `json:"source_id,omitempty"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`

	Focus *CSRFocus `json:"focus,omitempty"`
}
