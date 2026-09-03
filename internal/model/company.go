package model

import "time"

type Company struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Slug           string    `json:"slug"`
	StockCode      *string   `json:"stock_code,omitempty"`
	IndustrySector string    `json:"industry_sector"`
	AliasKeywords  []string  `json:"alias_keywords"`
	WebsiteURL     *string   `json:"website_url,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CompanyDetail struct {
	Company
	TargetCount  int     `json:"target_count"`
	SignalCount  int     `json:"signal_count"`
	TotalBudget  float64 `json:"total_budget_signal"`
}
