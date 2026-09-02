package repository

import (
	"context"
	"fmt"
	"time"

	"sovera-core-api/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CrawlerRepository struct {
	pool *pgxpool.Pool
}

func NewCrawlerRepository(pool *pgxpool.Pool) *CrawlerRepository {
	return &CrawlerRepository{pool: pool}
}

// GetDueTargets retrieves active crawling targets that are due for scraping
func (r *CrawlerRepository) GetDueTargets(ctx context.Context, limit int) ([]model.CrawlingTarget, error) {
	query := `
		SELECT id, source_name, source_type, target_url, check_interval_hours, 
		       last_scraped_at, next_run_at, is_active, created_at, updated_at
		FROM crawling_targets
		WHERE is_active = TRUE AND next_run_at <= $1
		ORDER BY next_run_at ASC
		LIMIT $2;
	`
	rows, err := r.pool.Query(ctx, query, time.Now(), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query due targets: %w", err)
	}
	defer rows.Close()

	var targets []model.CrawlingTarget
	for rows.Next() {
		var t model.CrawlingTarget
		err := rows.Scan(
			&t.ID, &t.SourceName, &t.SourceType, &t.TargetURL, &t.CheckIntervalHours,
			&t.LastScrapedAt, &t.NextRunAt, &t.IsActive, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan crawling target: %w", err)
		}
		targets = append(targets, t)
	}
	return targets, nil
}

// UpdateTargetNextRun updates last_scraped_at and schedules next_run_at
func (r *CrawlerRepository) UpdateTargetNextRun(ctx context.Context, targetID string, intervalHours int) error {
	query := `
		UPDATE crawling_targets
		SET last_scraped_at = NOW(),
		    next_run_at = NOW() + ($2 || ' hours')::INTERVAL,
		    updated_at = NOW()
		WHERE id = $1;
	`
	_, err := r.pool.Exec(ctx, query, targetID, intervalHours)
	if err != nil {
		return fmt.Errorf("failed to update target next run: %w", err)
	}
	return nil
}

// CreateLog records a task dispatch attempt
func (r *CrawlerRepository) CreateLog(ctx context.Context, log model.CrawlingLog) error {
	query := `
		INSERT INTO crawling_logs (target_id, task_id, status, http_status_code, error_message, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		ON CONFLICT (task_id) DO UPDATE 
		SET status = EXCLUDED.status, updated_at = NOW();
	`
	_, err := r.pool.Exec(ctx, query, log.TargetID, log.TaskID, log.Status, log.HTTPStatusCode, log.ErrorMessage)
	if err != nil {
		return fmt.Errorf("failed to insert crawling log: %w", err)
	}
	return nil
}

// UpdateLogStatus updates log status when callback arrives
func (r *CrawlerRepository) UpdateLogStatus(ctx context.Context, taskID, status string, execTimeMs *int, contentHash *string, errMsg *string) error {
	query := `
		UPDATE crawling_logs
		SET status = $2,
		    execution_time_ms = COALESCE($3, execution_time_ms),
		    content_hash = COALESCE($4, content_hash),
		    error_message = COALESCE($5, error_message),
		    updated_at = NOW()
		WHERE task_id = $1;
	`
	_, err := r.pool.Exec(ctx, query, taskID, status, execTimeMs, contentHash, errMsg)
	if err != nil {
		return fmt.Errorf("failed to update crawling log status: %w", err)
	}
	return nil
}
