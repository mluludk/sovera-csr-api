package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"sovera-core-api/internal/model"
)

type CompanyCSRFocusRepository struct {
	pool *pgxpool.Pool
}

func NewCompanyCSRFocusRepository(pool *pgxpool.Pool) *CompanyCSRFocusRepository {
	return &CompanyCSRFocusRepository{pool: pool}
}

// ListByCompanyID retrieves all CSR focus associations for a company, joining details from csr_focuses.
func (r *CompanyCSRFocusRepository) ListByCompanyID(ctx context.Context, companyID string) ([]model.CompanyCSRFocus, error) {
	if r.pool == nil {
		return nil, fmt.Errorf("database pool is nil")
	}

	query := `
		SELECT 
			cf.company_id::text, cf.focus_id::text, cf.priority, cf.confidence,
			cf.source_id::text, cf.verified_at, cf.created_at, cf.updated_at,
			f.id::text, f.code, f.name, f.category, f.description
		FROM company_csr_focuses cf
		JOIN csr_focuses f ON f.id = cf.focus_id
		WHERE cf.company_id::text = $1
		ORDER BY cf.priority ASC NULLS LAST, f.name ASC;
	`

	rows, err := r.pool.Query(ctx, query, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to query company_csr_focuses: %w", err)
	}
	defer rows.Close()

	var results []model.CompanyCSRFocus
	for rows.Next() {
		var cf model.CompanyCSRFocus
		var f model.CSRFocus
		err := rows.Scan(
			&cf.CompanyID, &cf.FocusID, &cf.Priority, &cf.Confidence,
			&cf.SourceID, &cf.VerifiedAt, &cf.CreatedAt, &cf.UpdatedAt,
			&f.ID, &f.Code, &f.Name, &f.Category, &f.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan company_csr_focus row: %w", err)
		}
		cf.Focus = &f
		results = append(results, cf)
	}

	return results, nil
}

// Upsert inserts or updates a company CSR focus junction entry.
func (r *CompanyCSRFocusRepository) Upsert(ctx context.Context, item *model.CompanyCSRFocus) error {
	if r.pool == nil {
		return fmt.Errorf("database pool is nil")
	}

	query := `
		INSERT INTO company_csr_focuses (
			company_id, focus_id, priority, confidence, source_id, verified_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, NOW()
		)
		ON CONFLICT (company_id, focus_id) DO UPDATE SET
			priority = EXCLUDED.priority,
			confidence = EXCLUDED.confidence,
			source_id = EXCLUDED.source_id,
			verified_at = EXCLUDED.verified_at,
			updated_at = NOW();
	`

	_, err := r.pool.Exec(ctx, query,
		item.CompanyID, item.FocusID, item.Priority, item.Confidence, item.SourceID, item.VerifiedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert company_csr_focus: %w", err)
	}

	return nil
}

// Delete removes a specific focus assignment for a company.
func (r *CompanyCSRFocusRepository) Delete(ctx context.Context, companyID, focusID string) error {
	if r.pool == nil {
		return fmt.Errorf("database pool is nil")
	}

	query := `DELETE FROM company_csr_focuses WHERE company_id::text = $1 AND focus_id::text = $2;`
	_, err := r.pool.Exec(ctx, query, companyID, focusID)
	if err != nil {
		return fmt.Errorf("failed to delete company_csr_focus: %w", err)
	}

	return nil
}
