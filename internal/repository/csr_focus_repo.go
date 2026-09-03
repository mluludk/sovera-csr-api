package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"sovera-core-api/internal/model"
)

type CSRFocusRepository struct {
	pool *pgxpool.Pool
}

func NewCSRFocusRepository(pool *pgxpool.Pool) *CSRFocusRepository {
	return &CSRFocusRepository{pool: pool}
}

// List retrieves all CSR focus entries ordered by category and name.
func (r *CSRFocusRepository) List(ctx context.Context) ([]model.CSRFocus, error) {
	if r.pool == nil {
		return nil, fmt.Errorf("database pool is nil")
	}

	query := `
		SELECT id::text, code, name, category, description, created_at, updated_at
		FROM csr_focuses
		ORDER BY category ASC, name ASC;
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query csr_focuses: %w", err)
	}
	defer rows.Close()

	var results []model.CSRFocus
	for rows.Next() {
		var f model.CSRFocus
		if err := rows.Scan(&f.ID, &f.Code, &f.Name, &f.Category, &f.Description, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan csr_focus row: %w", err)
		}
		results = append(results, f)
	}

	return results, nil
}

// GetByCode retrieves a single CSR focus by its unique code.
func (r *CSRFocusRepository) GetByCode(ctx context.Context, code string) (*model.CSRFocus, error) {
	if r.pool == nil {
		return nil, fmt.Errorf("database pool is nil")
	}

	query := `
		SELECT id::text, code, name, category, description, created_at, updated_at
		FROM csr_focuses
		WHERE code = $1
		LIMIT 1;
	`

	var f model.CSRFocus
	err := r.pool.QueryRow(ctx, query, code).Scan(&f.ID, &f.Code, &f.Name, &f.Category, &f.Description, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("csr_focus not found for code %s: %w", code, err)
	}

	return &f, nil
}
