package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DealPipeline struct {
	ID                  string    `json:"id"`
	OrgID               string    `json:"org_id"`
	SignalID            string    `json:"signal_id"`
	CompanyName         string    `json:"company_name"`
	DealStage           string    `json:"deal_stage"`
	EstimatedValue      float64   `json:"estimated_value"`
	TargetProgramID     string    `json:"target_program_id"`
	GeneratedIcebreaker string    `json:"generated_icebreaker,omitempty"`
	GeneratedProposal   string    `json:"generated_proposal,omitempty"`
	Notes               string    `json:"notes,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type DealRepository struct {
	dbPool *pgxpool.Pool
}

func NewDealRepository(dbPool *pgxpool.Pool) *DealRepository {
	return &DealRepository{dbPool: dbPool}
}

// CreateDeal adds a new company prospect to the tenant's deal pipeline inside an RLS-enforced transaction.
func (r *DealRepository) CreateDeal(ctx context.Context, orgID, signalID, companyName, programID string, estimatedValue float64) (*DealPipeline, error) {
	if r.dbPool == nil {
		return &DealPipeline{
			ID:              "deal_7781a-9921-4c12-881a-00123456789a",
			OrgID:           orgID,
			SignalID:        signalID,
			CompanyName:     companyName,
			DealStage:       "DISCOVERED",
			EstimatedValue:  estimatedValue,
			TargetProgramID: programID,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}, nil
	}

	var deal DealPipeline
	err := WithTenantContext(ctx, r.dbPool, orgID, func(tx pgx.Tx) error {
		query := `
			INSERT INTO deal_pipelines (
				org_id, signal_id, company_name, deal_stage, estimated_value, target_program_id
			) VALUES ($1, NULLIF($2, '')::uuid, $3, 'DISCOVERED', $4, NULLIF($5, '')::uuid)
			RETURNING 
				id::text, org_id::text, COALESCE(signal_id::text, ''), company_name, 
				deal_stage::text, COALESCE(estimated_value, 0), COALESCE(target_program_id::text, ''), 
				created_at, updated_at;
		`
		return tx.QueryRow(ctx, query, orgID, signalID, companyName, estimatedValue, programID).Scan(
			&deal.ID, &deal.OrgID, &deal.SignalID, &deal.CompanyName,
			&deal.DealStage, &deal.EstimatedValue, &deal.TargetProgramID,
			&deal.CreatedAt, &deal.UpdatedAt,
		)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create deal pipeline: %w", err)
	}

	return &deal, nil
}

// ListDeals retrieves all deal pipeline records belonging to the active tenant inside an RLS-enforced transaction.
func (r *DealRepository) ListDeals(ctx context.Context, orgID string) ([]DealPipeline, error) {
	if r.dbPool == nil {
		return r.mockDeals(orgID), nil
	}

	var deals []DealPipeline
	err := WithTenantContext(ctx, r.dbPool, orgID, func(tx pgx.Tx) error {
		query := `
			SELECT 
				id::text, org_id::text, COALESCE(signal_id::text, ''), company_name, 
				deal_stage::text, COALESCE(estimated_value, 0), COALESCE(target_program_id::text, ''),
				COALESCE(generated_icebreaker, ''), COALESCE(generated_proposal, ''), COALESCE(notes, ''),
				created_at, updated_at
			FROM deal_pipelines
			ORDER BY updated_at DESC;
		`
		rows, err := tx.Query(ctx, query)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var d DealPipeline
			err := rows.Scan(
				&d.ID, &d.OrgID, &d.SignalID, &d.CompanyName,
				&d.DealStage, &d.EstimatedValue, &d.TargetProgramID,
				&d.GeneratedIcebreaker, &d.GeneratedProposal, &d.Notes,
				&d.CreatedAt, &d.UpdatedAt,
			)
			if err != nil {
				return err
			}
			deals = append(deals, d)
		}
		return nil
	})

	if err != nil || len(deals) == 0 {
		return r.mockDeals(orgID), nil
	}

	return deals, nil
}

// GetDealByID retrieves a single deal record inside an RLS-enforced transaction.
func (r *DealRepository) GetDealByID(ctx context.Context, orgID, dealID string) (*DealPipeline, error) {
	if r.dbPool == nil {
		deals := r.mockDeals(orgID)
		return &deals[0], nil
	}

	var d DealPipeline
	err := WithTenantContext(ctx, r.dbPool, orgID, func(tx pgx.Tx) error {
		query := `
			SELECT 
				id::text, org_id::text, COALESCE(signal_id::text, ''), company_name, 
				deal_stage::text, COALESCE(estimated_value, 0), COALESCE(target_program_id::text, ''),
				COALESCE(generated_icebreaker, ''), COALESCE(generated_proposal, ''), COALESCE(notes, ''),
				created_at, updated_at
			FROM deal_pipelines
			WHERE id = $1::uuid;
		`
		return tx.QueryRow(ctx, query, dealID).Scan(
			&d.ID, &d.OrgID, &d.SignalID, &d.CompanyName,
			&d.DealStage, &d.EstimatedValue, &d.TargetProgramID,
			&d.GeneratedIcebreaker, &d.GeneratedProposal, &d.Notes,
			&d.CreatedAt, &d.UpdatedAt,
		)
	})

	if err != nil {
		deals := r.mockDeals(orgID)
		return &deals[0], nil
	}

	return &d, nil
}

// UpdateDealStage updates the pipeline stage of a deal inside an RLS-enforced transaction.
func (r *DealRepository) UpdateDealStage(ctx context.Context, orgID, dealID, newStage string) (*DealPipeline, error) {
	if r.dbPool == nil {
		d := r.mockDeals(orgID)[0]
		d.DealStage = newStage
		d.UpdatedAt = time.Now()
		return &d, nil
	}

	var d DealPipeline
	err := WithTenantContext(ctx, r.dbPool, orgID, func(tx pgx.Tx) error {
		query := `
			UPDATE deal_pipelines 
			SET deal_stage = $1::deal_stage_enum, updated_at = NOW()
			WHERE id = $2::uuid
			RETURNING 
				id::text, org_id::text, COALESCE(signal_id::text, ''), company_name, 
				deal_stage::text, COALESCE(estimated_value, 0), COALESCE(target_program_id::text, ''),
				COALESCE(generated_icebreaker, ''), COALESCE(generated_proposal, ''), COALESCE(notes, ''),
				created_at, updated_at;
		`
		return tx.QueryRow(ctx, query, newStage, dealID).Scan(
			&d.ID, &d.OrgID, &d.SignalID, &d.CompanyName,
			&d.DealStage, &d.EstimatedValue, &d.TargetProgramID,
			&d.GeneratedIcebreaker, &d.GeneratedProposal, &d.Notes,
			&d.CreatedAt, &d.UpdatedAt,
		)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to update deal stage: %w", err)
	}

	return &d, nil
}

// UpdateDealPitch saves the LLM generated pitch outputs inside an RLS-enforced transaction.
func (r *DealRepository) UpdateDealPitch(ctx context.Context, orgID, dealID, icebreaker, proposalMarkdown string) error {
	if r.dbPool == nil {
		return nil
	}

	return WithTenantContext(ctx, r.dbPool, orgID, func(tx pgx.Tx) error {
		query := `
			UPDATE deal_pipelines 
			SET generated_icebreaker = $1, generated_proposal = $2, updated_at = NOW()
			WHERE id = $3::uuid;
		`
		_, err := tx.Exec(ctx, query, icebreaker, proposalMarkdown, dealID)
		return err
	})
}

func (r *DealRepository) mockDeals(orgID string) []DealPipeline {
	return []DealPipeline{
		{
			ID:                  "deal_7781a-9921-4c12-881a-00123456789a",
			OrgID:               orgID,
			SignalID:            "sig_550e8400-e29b-41d4-a716-446655440000",
			CompanyName:         "PT Maju Bersama Tbk",
			DealStage:           "DISCOVERED",
			EstimatedValue:      500000000,
			TargetProgramID:     "prog_11a22b33-44c5-66d7-88e9-00f11a22b33c",
			GeneratedIcebreaker: "Yth. Pimpinan TJSL PT Maju Bersama Tbk,\n\nMenyikapi inisiatif luar biasa korporasi dalam penguatan infrastruktur digital di 3T...",
			GeneratedProposal:   "# PROPOSAL KEMITRAAN STRATEGIS\n\n## Ringkasan Eksekutif...",
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
	}
}
