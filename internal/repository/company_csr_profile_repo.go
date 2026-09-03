package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"sovera-core-api/internal/model"
)

type CompanyCSRProfileRepository struct {
	pool *pgxpool.Pool
}

func NewCompanyCSRProfileRepository(pool *pgxpool.Pool) *CompanyCSRProfileRepository {
	return &CompanyCSRProfileRepository{pool: pool}
}

// GetByCompanyID retrieves the CSR profile for a given company ID.
func (r *CompanyCSRProfileRepository) GetByCompanyID(ctx context.Context, companyID string) (*model.CompanyCSRProfile, error) {
	if r.pool == nil {
		return nil, fmt.Errorf("database pool is nil")
	}

	query := `
		SELECT 
			id::text, company_id::text, has_csr, csr_department_name, csr_email_public,
			COALESCE(csr_focus, '{}'), budget_range, proposal_acceptance, website_source,
			last_verified_at, created_at, updated_at
		FROM company_csr_profiles
		WHERE company_id::text = $1
		LIMIT 1;
	`

	var p model.CompanyCSRProfile
	err := r.pool.QueryRow(ctx, query, companyID).Scan(
		&p.ID, &p.CompanyID, &p.HasCSR, &p.CSRDepartmentName, &p.CSREmailPublic,
		&p.CSRFocus, &p.BudgetRange, &p.ProposalAcceptance, &p.WebsiteSource,
		&p.LastVerifiedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("company CSR profile not found: %w", err)
	}

	return &p, nil
}

// Upsert inserts or updates a company CSR profile.
func (r *CompanyCSRProfileRepository) Upsert(ctx context.Context, profile *model.CompanyCSRProfile) (*model.CompanyCSRProfile, error) {
	if r.pool == nil {
		return nil, fmt.Errorf("database pool is nil")
	}

	query := `
		INSERT INTO company_csr_profiles (
			company_id, has_csr, csr_department_name, csr_email_public,
			csr_focus, budget_range, proposal_acceptance, website_source,
			last_verified_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, NOW()
		)
		ON CONFLICT (company_id) DO UPDATE SET
			has_csr = EXCLUDED.has_csr,
			csr_department_name = EXCLUDED.csr_department_name,
			csr_email_public = EXCLUDED.csr_email_public,
			csr_focus = EXCLUDED.csr_focus,
			budget_range = EXCLUDED.budget_range,
			proposal_acceptance = EXCLUDED.proposal_acceptance,
			website_source = EXCLUDED.website_source,
			last_verified_at = EXCLUDED.last_verified_at,
			updated_at = NOW()
		RETURNING 
			id::text, company_id::text, has_csr, csr_department_name, csr_email_public,
			COALESCE(csr_focus, '{}'), budget_range, proposal_acceptance, website_source,
			last_verified_at, created_at, updated_at;
	`

	var p model.CompanyCSRProfile
	err := r.pool.QueryRow(ctx, query,
		profile.CompanyID, profile.HasCSR, profile.CSRDepartmentName, profile.CSREmailPublic,
		profile.CSRFocus, profile.BudgetRange, profile.ProposalAcceptance, profile.WebsiteSource,
		profile.LastVerifiedAt,
	).Scan(
		&p.ID, &p.CompanyID, &p.HasCSR, &p.CSRDepartmentName, &p.CSREmailPublic,
		&p.CSRFocus, &p.BudgetRange, &p.ProposalAcceptance, &p.WebsiteSource,
		&p.LastVerifiedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert company CSR profile: %w", err)
	}

	return &p, nil
}

// Delete removes the CSR profile for a given company ID.
func (r *CompanyCSRProfileRepository) Delete(ctx context.Context, companyID string) error {
	if r.pool == nil {
		return fmt.Errorf("database pool is nil")
	}

	query := `DELETE FROM company_csr_profiles WHERE company_id::text = $1;`
	_, err := r.pool.Exec(ctx, query, companyID)
	if err != nil {
		return fmt.Errorf("failed to delete company CSR profile: %w", err)
	}

	return nil
}
