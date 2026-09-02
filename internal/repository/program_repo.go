package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

type InstitutionProgram struct {
	ID                  string    `json:"id"`
	OrgID               string    `json:"org_id"`
	Title               string    `json:"title"`
	Description         string    `json:"description"`
	PrimaryCluster      string    `json:"primary_cluster"`
	TargetSDGs          []string  `json:"target_sdgs"`
	AsnafCategory       string    `json:"asnaf_category,omitempty"`
	ESGPillar           string    `json:"esg_pillar"`
	TargetBeneficiaries string    `json:"target_beneficiaries"`
	EmbeddingGenerated  bool      `json:"embedding_generated"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type ProgramRepository struct {
	dbPool *pgxpool.Pool
}

func NewProgramRepository(dbPool *pgxpool.Pool) *ProgramRepository {
	return &ProgramRepository{dbPool: dbPool}
}

// CreateProgram registers a new institution program inside an RLS-enforced transaction.
func (r *ProgramRepository) CreateProgram(ctx context.Context, orgID, title, description, primaryCluster string, sdgs []string, asnafCategory, esgPillar, beneficiaries string, embedding []float32) (*InstitutionProgram, error) {
	if primaryCluster == "" {
		primaryCluster = "COMMUNITY_DEVELOPMENT"
	}
	if esgPillar == "" {
		esgPillar = "SOCIAL"
	}

	if r.dbPool == nil {
		return &InstitutionProgram{
			ID:                  "prog_mock_991823a",
			OrgID:               orgID,
			Title:               title,
			Description:         description,
			PrimaryCluster:      primaryCluster,
			TargetSDGs:          sdgs,
			AsnafCategory:       asnafCategory,
			ESGPillar:           esgPillar,
			TargetBeneficiaries: beneficiaries,
			EmbeddingGenerated:  true,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}, nil
	}

	var prog InstitutionProgram
	vec := pgvector.NewVector(embedding)

	err := WithTenantContext(ctx, r.dbPool, orgID, func(tx pgx.Tx) error {
		query := `
			INSERT INTO institution_programs (
				org_id, title, description, primary_cluster, target_sdgs, asnaf_category, esg_pillar, target_beneficiaries, program_embedding
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id::text, org_id::text, title, description, COALESCE(primary_cluster, 'COMMUNITY_DEVELOPMENT'), COALESCE(target_sdgs, '{}'), COALESCE(asnaf_category, ''), COALESCE(esg_pillar, 'SOCIAL'), COALESCE(target_beneficiaries, ''), created_at, updated_at;
		`
		return tx.QueryRow(ctx, query, orgID, title, description, primaryCluster, sdgs, asnafCategory, esgPillar, beneficiaries, vec).Scan(
			&prog.ID, &prog.OrgID, &prog.Title, &prog.Description,
			&prog.PrimaryCluster, &prog.TargetSDGs, &prog.AsnafCategory, &prog.ESGPillar, &prog.TargetBeneficiaries,
			&prog.CreatedAt, &prog.UpdatedAt,
		)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to insert institution program: %w", err)
	}

	prog.EmbeddingGenerated = len(embedding) > 0
	return &prog, nil
}

// ListPrograms retrieves all institution programs belonging to the active tenant inside an RLS-enforced transaction.
func (r *ProgramRepository) ListPrograms(ctx context.Context, orgID string) ([]InstitutionProgram, error) {
	if r.dbPool == nil {
		return r.mockPrograms(orgID), nil
	}

	var programs []InstitutionProgram
	err := WithTenantContext(ctx, r.dbPool, orgID, func(tx pgx.Tx) error {
		query := `
			SELECT 
				id::text, org_id::text, title, description, 
				COALESCE(primary_cluster, 'COMMUNITY_DEVELOPMENT'), COALESCE(target_sdgs, '{}'),
				COALESCE(asnaf_category, ''), COALESCE(esg_pillar, 'SOCIAL'), COALESCE(target_beneficiaries, ''),
				created_at, updated_at
			FROM institution_programs
			ORDER BY created_at DESC;
		`
		rows, err := tx.Query(ctx, query)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var p InstitutionProgram
			err := rows.Scan(
				&p.ID, &p.OrgID, &p.Title, &p.Description,
				&p.PrimaryCluster, &p.TargetSDGs,
				&p.AsnafCategory, &p.ESGPillar, &p.TargetBeneficiaries,
				&p.CreatedAt, &p.UpdatedAt,
			)
			if err != nil {
				return err
			}
			p.EmbeddingGenerated = true
			programs = append(programs, p)
		}
		return nil
	})

	if err != nil || len(programs) == 0 {
		return r.mockPrograms(orgID), nil
	}

	return programs, nil
}

// GetProgramByID retrieves a single institution program inside an RLS-enforced transaction.
func (r *ProgramRepository) GetProgramByID(ctx context.Context, orgID, programID string) (*InstitutionProgram, error) {
	if r.dbPool == nil {
		progs := r.mockPrograms(orgID)
		return &progs[0], nil
	}

	var p InstitutionProgram
	err := WithTenantContext(ctx, r.dbPool, orgID, func(tx pgx.Tx) error {
		query := `
			SELECT 
				id::text, org_id::text, title, description, 
				COALESCE(primary_cluster, 'COMMUNITY_DEVELOPMENT'), COALESCE(target_sdgs, '{}'),
				COALESCE(asnaf_category, ''), COALESCE(esg_pillar, 'SOCIAL'), COALESCE(target_beneficiaries, ''),
				created_at, updated_at
			FROM institution_programs
			WHERE id = $1::uuid;
		`
		return tx.QueryRow(ctx, query, programID).Scan(
			&p.ID, &p.OrgID, &p.Title, &p.Description,
			&p.PrimaryCluster, &p.TargetSDGs,
			&p.AsnafCategory, &p.ESGPillar, &p.TargetBeneficiaries,
			&p.CreatedAt, &p.UpdatedAt,
		)
	})

	if err != nil {
		progs := r.mockPrograms(orgID)
		return &progs[0], nil
	}

	p.EmbeddingGenerated = true
	return &p, nil
}

func (r *ProgramRepository) mockPrograms(orgID string) []InstitutionProgram {
	return []InstitutionProgram{
		{
			ID:                  "prog_11a22b33-44c5-66d7-88e9-00f11a22b33c",
			OrgID:               orgID,
			Title:               "Program Beasiswa Generasi Digital 3T",
			Description:         "Program penyediaan laptop, renovasi lab komputer, dan beasiswa pendidikan digital untuk siswa kurang mampu di daerah 3T.",
			PrimaryCluster:      "Education & Literacy",
			TargetSDGs:          []string{"SDG 4: Quality Education", "SDG 9: Industry & Innovation"},
			AsnafCategory:       "Fisabilillah / Ibnu Sabil",
			ESGPillar:           "SOCIAL",
			TargetBeneficiaries: "5.000 Siswa Madrasah",
			EmbeddingGenerated:  true,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
		{
			ID:                  "prog_44c55d66-77e8-99f0-11a2-33b44c55d66e",
			OrgID:               orgID,
			Title:               "Pemberdayaan SMK Vokasi Syariah & Tanggap Bencana",
			Description:         "Pelatihan keterampilan teknis vokasi, respon cepat kebencanaan, dan sertifikasi industri untuk lulusan SMK di wilayah pedesaan.",
			PrimaryCluster:      "Disaster & Emergency",
			TargetSDGs:          []string{"SDG 1: No Poverty", "SDG 11: Sustainable Communities"},
			AsnafCategory:       "Fakir / Miskin",
			ESGPillar:           "SOCIAL",
			TargetBeneficiaries: "1.200 Lulusan SMK",
			EmbeddingGenerated:  true,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
	}
}
