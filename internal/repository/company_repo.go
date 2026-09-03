package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"sovera-core-api/internal/model"
)

type CompanyRepository struct {
	pool *pgxpool.Pool
}

func NewCompanyRepository(pool *pgxpool.Pool) *CompanyRepository {
	return &CompanyRepository{pool: pool}
}

// ListCompanies retrieves a paginated list of companies with target and signal counts.
func (r *CompanyRepository) ListCompanies(ctx context.Context, limit, offset int, search, sector string) ([]model.CompanyDetail, int, error) {
	if r.pool == nil {
		return nil, 0, nil
	}

	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIdx := 1

	if search != "" {
		whereClause += fmt.Sprintf(" AND (c.name ILIKE $%d OR array_to_string(c.alias_keywords, ' ') ILIKE $%d OR c.stock_code ILIKE $%d OR c.ticker ILIKE $%d)", argIdx, argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	if sector != "" && sector != "ALL" {
		whereClause += fmt.Sprintf(" AND c.industry_sector = $%d", argIdx)
		args = append(args, sector)
		argIdx++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM companies c %s", whereClause)
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count companies: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT 
			c.id::text, c.name, c.legal_name, c.slug, c.industry_id, c.industry_sector,
			c.company_type, c.website, c.website_url, c.linkedin_url, c.headquarters,
			c.employee_range, c.revenue_range, c.is_public, c.ticker, c.stock_code,
			c.parent_company_id::text, COALESCE(c.alias_keywords, '{}'), c.created_at, c.updated_at,
			(SELECT COUNT(*) FROM crawling_targets t WHERE t.company_id = c.id) AS target_count,
			(SELECT COUNT(*) FROM public_corporate_signals s WHERE s.company_id = c.id OR s.company_name ILIKE c.name) AS signal_count,
			COALESCE((SELECT SUM(s.estimated_budget_signal) FROM public_corporate_signals s WHERE s.company_id = c.id OR s.company_name ILIKE c.name), 0) AS total_budget
		FROM companies c
		%s
		ORDER BY signal_count DESC, c.name ASC
		LIMIT $%d OFFSET $%d;
	`, whereClause, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query companies: %w", err)
	}
	defer rows.Close()

	var result []model.CompanyDetail
	for rows.Next() {
		var cd model.CompanyDetail
		err := rows.Scan(
			&cd.ID, &cd.Name, &cd.LegalName, &cd.Slug, &cd.IndustryID, &cd.IndustrySector,
			&cd.CompanyType, &cd.Website, &cd.WebsiteURL, &cd.LinkedinURL, &cd.Headquarters,
			&cd.EmployeeRange, &cd.RevenueRange, &cd.IsPublic, &cd.Ticker, &cd.StockCode,
			&cd.ParentCompanyID, &cd.AliasKeywords, &cd.CreatedAt, &cd.UpdatedAt,
			&cd.TargetCount, &cd.SignalCount, &cd.TotalBudget,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan company row: %w", err)
		}
		result = append(result, cd)
	}

	return result, total, nil
}

// GetCompanyByID retrieves a single company detail by ID or Slug.
func (r *CompanyRepository) GetCompanyByID(ctx context.Context, idOrSlug string) (*model.CompanyDetail, error) {
	if r.pool == nil {
		return nil, fmt.Errorf("database pool is nil")
	}

	query := `
		SELECT 
			c.id::text, c.name, c.legal_name, c.slug, c.industry_id, c.industry_sector,
			c.company_type, c.website, c.website_url, c.linkedin_url, c.headquarters,
			c.employee_range, c.revenue_range, c.is_public, c.ticker, c.stock_code,
			c.parent_company_id::text, COALESCE(c.alias_keywords, '{}'), c.created_at, c.updated_at,
			(SELECT COUNT(*) FROM crawling_targets t WHERE t.company_id = c.id) AS target_count,
			(SELECT COUNT(*) FROM public_corporate_signals s WHERE s.company_id = c.id OR s.company_name ILIKE c.name) AS signal_count,
			COALESCE((SELECT SUM(s.estimated_budget_signal) FROM public_corporate_signals s WHERE s.company_id = c.id OR s.company_name ILIKE c.name), 0) AS total_budget
		FROM companies c
		WHERE c.id::text = $1 OR c.slug = $1
		LIMIT 1;
	`

	var cd model.CompanyDetail
	err := r.pool.QueryRow(ctx, query, idOrSlug).Scan(
		&cd.ID, &cd.Name, &cd.LegalName, &cd.Slug, &cd.IndustryID, &cd.IndustrySector,
		&cd.CompanyType, &cd.Website, &cd.WebsiteURL, &cd.LinkedinURL, &cd.Headquarters,
		&cd.EmployeeRange, &cd.RevenueRange, &cd.IsPublic, &cd.Ticker, &cd.StockCode,
		&cd.ParentCompanyID, &cd.AliasKeywords, &cd.CreatedAt, &cd.UpdatedAt,
		&cd.TargetCount, &cd.SignalCount, &cd.TotalBudget,
	)
	if err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}

	return &cd, nil
}
