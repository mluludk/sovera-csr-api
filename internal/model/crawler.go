package model

import (
	"time"
)

type CrawlingTarget struct {
	ID                 string     `json:"id" db:"id"`
	SourceName         string     `json:"source_name" db:"source_name"`
	SourceType         string     `json:"source_type" db:"source_type"`
	TargetURL          string     `json:"target_url" db:"target_url"`
	CheckIntervalHours int        `json:"check_interval_hours" db:"check_interval_hours"`
	LastScrapedAt      *time.Time `json:"last_scraped_at,omitempty" db:"last_scraped_at"`
	NextRunAt          time.Time  `json:"next_run_at" db:"next_run_at"`
	IsActive           bool       `json:"is_active" db:"is_active"`
	ConsecutiveFailures int        `json:"consecutive_failures" db:"consecutive_failures"`
	LastHTTPStatus      *int       `json:"last_http_status,omitempty" db:"last_http_status"`
	LastErrorMsg        *string    `json:"last_error_message,omitempty" db:"last_error_message"`
	HealthStatus        string     `json:"health_status" db:"health_status"` // HEALTHY, DEGRADED, DISABLED_DEAD_LINK
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
}

type CrawlingLog struct {
	ID              string    `json:"id" db:"id"`
	TargetID        *string   `json:"target_id,omitempty" db:"target_id"`
	TaskID          string    `json:"task_id" db:"task_id"`
	Status          string    `json:"status" db:"status"` // DISPATCHED, COMPLETED, FAILED
	HTTPStatusCode  *int      `json:"http_status_code,omitempty" db:"http_status_code"`
	ErrorMessage    *string   `json:"error_message,omitempty" db:"error_message"`
	ExecutionTimeMs *int      `json:"execution_time_ms,omitempty" db:"execution_time_ms"`
	ContentHash     *string   `json:"content_hash,omitempty" db:"content_hash"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type ScrapeTaskConfig struct {
	RenderJS      bool `json:"render_js"`
	BypassAntiBot bool `json:"bypass_anti_bot"`
	MaxPages      int  `json:"max_pages,omitempty"`
}

type ScrapeTaskPayload struct {
	TaskID       string            `json:"task_id"`
	TargetID     string            `json:"target_id,omitempty"`
	ClientOrigin string            `json:"client_origin"`
	SourceType   string            `json:"source_type"`
	TargetURL    string            `json:"target_url"`
	CallbackURL  string            `json:"callback_url"`
	Config       *ScrapeTaskConfig `json:"config,omitempty"`
}
