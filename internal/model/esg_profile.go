package model

import "time"

type CompanyESGProfile struct {
	ID                     string                 `json:"id"`
	CompanyID              string                 `json:"company_id"`
	ReportingYear          int16                  `json:"reporting_year"`
	ReportDate             *time.Time             `json:"report_date,omitempty"`
	OverallScore           *float64               `json:"overall_score,omitempty"`
	EnvironmentalScore     *float64               `json:"environmental_score,omitempty"`
	SocialScore            *float64               `json:"social_score,omitempty"`
	GovernanceScore        *float64               `json:"governance_score,omitempty"`
	ESGRating              *string                `json:"esg_rating,omitempty"`
	SustainabilityStrategy *string                `json:"sustainability_strategy,omitempty"`
	SDGAlignment           map[string]interface{} `json:"sdg_alignment,omitempty"`
	SourceID               *string                `json:"source_id,omitempty"`
	Confidence             *float64               `json:"confidence,omitempty"`
	CreatedAt              time.Time              `json:"created_at"`
	UpdatedAt              time.Time              `json:"updated_at"`

	MaterialTopics []CompanyESGMaterialTopic `json:"material_topics,omitempty"`
}

type ESGMaterialTopic struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CompanyESGMaterialTopic struct {
	ESGProfileID    string            `json:"esg_profile_id"`
	TopicID         string            `json:"topic_id"`
	MaterialityScore *float64         `json:"materiality_score,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	Topic           *ESGMaterialTopic `json:"topic,omitempty"`
}
