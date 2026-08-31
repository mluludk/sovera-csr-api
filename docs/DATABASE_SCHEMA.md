# Database Schema & RLS Specification (DATABASE_SCHEMA.md)

**Product:** Sovera (FundIQ) Core API  
**Database Engine:** PostgreSQL 16+ with `pgvector` and `uuid-ossp`  
**Security Model:** Multi-Tenant Row-Level Security (RLS) Enforced at Kernel Level  

---

## 1. Architectural Overview & Security Model

PostgreSQL dirancang dengan pemisahan tegas antara **Data Intelijen Publik (Shared Data)** dan **Data Lembaga Privat (Tenant-Isolated Data)**:

* **Tabel Publik / Bersama (`public_corporate_signals`):** Berisi data hasil crawling bursa, rilis berita, dan laporan keberlanjutan yang dapat dibaca oleh seluruh tenant terdaftar.
* **Tabel Privat (`institution_programs`, `deal_pipelines`, `users`):** Diberikan proteksi `ENABLE ROW LEVEL SECURITY` dan `FORCE ROW LEVEL SECURITY`. Akses data diisolasi secara otomatis berdasarkan parameter transaksi `app.current_org_id`.

---

## 2. Complete SQL DDL & Migration Script

```sql
-- 1. SETUP EXTENSIONS
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "vector";

-- 2. ENUM TYPES
DO $$ BEGIN
    CREATE TYPE deal_stage_enum AS ENUM (
        'DISCOVERED', 
        'RESEARCH', 
        'PITCHED', 
        'NEGOTIATION', 
        'CLOSED_WON', 
        'CLOSED_LOST'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE signal_source_enum AS ENUM (
        'BEI_REPORT', 
        'NEWS', 
        'CSR_PDF', 
        'SOCIAL'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- 3. TENANT & AUTH TABLES (PUBLIC READ / SYSTEM MANAGED)
CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    subscription_tier VARCHAR(50) DEFAULT 'PRO',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'FUNDRAISER',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 4. SHARED MARKET DATA (NO RLS - SHARED ACROSS TENANTS)
CREATE TABLE IF NOT EXISTS public_corporate_signals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_name VARCHAR(255) NOT NULL,
    industry_sector VARCHAR(100),
    source_type signal_source_enum NOT NULL,
    source_url TEXT,
    summary TEXT,
    extracted_pillar VARCHAR(100),
    target_regions TEXT[],
    estimated_budget_signal NUMERIC(15, 2),
    trigger_event VARCHAR(255),
    intent_score INT CHECK (intent_score BETWEEN 0 AND 100),
    content_hash VARCHAR(64) UNIQUE,
    signal_embedding vector(1536),
    published_date DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 5. TENANT ISOLATED TABLES (RLS ENFORCED)
CREATE TABLE IF NOT EXISTS institution_programs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    asnaf_category VARCHAR(100),
    esg_pillar VARCHAR(100),
    target_beneficiaries VARCHAR(255),
    program_embedding vector(1536),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS deal_pipelines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    signal_id UUID REFERENCES public_corporate_signals(id) ON DELETE SET NULL,
    company_name VARCHAR(255) NOT NULL,
    deal_stage deal_stage_enum DEFAULT 'DISCOVERED',
    estimated_value NUMERIC(15, 2),
    target_program_id UUID REFERENCES institution_programs(id) ON DELETE SET NULL,
    generated_icebreaker TEXT,
    generated_proposal TEXT,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================================
-- 6. ROW-LEVEL SECURITY (RLS) POLICIES
-- ============================================================

-- Activate RLS
ALTER TABLE institution_programs ENABLE ROW LEVEL SECURITY;
ALTER TABLE institution_programs FORCE ROW LEVEL SECURITY;

ALTER TABLE deal_pipelines ENABLE ROW LEVEL SECURITY;
ALTER TABLE deal_pipelines FORCE ROW LEVEL SECURITY;

-- Policies for institution_programs
CREATE POLICY tenant_isolation_programs ON institution_programs
    FOR ALL
    USING (org_id = NULLIF(current_setting('app.current_org_id', true), '')::UUID)
    WITH CHECK (org_id = NULLIF(current_setting('app.current_org_id', true), '')::UUID);

-- Policies for deal_pipelines
CREATE POLICY tenant_isolation_deals ON deal_pipelines
    FOR ALL
    USING (org_id = NULLIF(current_setting('app.current_org_id', true), '')::UUID)
    WITH CHECK (org_id = NULLIF(current_setting('app.current_org_id', true), '')::UUID);

-- ============================================================
-- 7. PERFORMANCE INDEXES (HNSW & BTREE)
-- ============================================================

CREATE INDEX IF NOT EXISTS idx_signals_embedding 
    ON public_corporate_signals 
    USING hnsw (signal_embedding vector_cosine_ops);

CREATE INDEX IF NOT EXISTS idx_programs_embedding 
    ON institution_programs 
    USING hnsw (program_embedding vector_cosine_ops);

CREATE INDEX IF NOT EXISTS idx_programs_org_id ON institution_programs(org_id);
CREATE INDEX IF NOT EXISTS idx_deals_org_id ON deal_pipelines(org_id);
CREATE INDEX IF NOT EXISTS idx_deals_stage ON deal_pipelines(deal_stage);
CREATE INDEX IF NOT EXISTS idx_signals_intent ON public_corporate_signals(intent_score DESC);
CREATE INDEX IF NOT EXISTS idx_signals_hash ON public_corporate_signals(content_hash);

```

---

## 3. Data Dictionary

| Table Name | Description | Isolation Level | Primary Keys / Foreign Keys |
| --- | --- | --- | --- |
| `organizations` | Data tenant/lembaga zakat/NGO terdaftar | Global System | PK: `id` |
| `users` | Akun pengguna/fundraiser di bawah tenant | Global System | PK: `id`, FK: `org_id` -> `organizations.id` |
| `public_corporate_signals` | Data intelijen CSR/ESG hasil ekstraksi crawler | Shared (Public Read) | PK: `id` |
| `institution_programs` | Portofolio program unggulan & embedding program | **RLS Enforced** | PK: `id`, FK: `org_id` -> `organizations.id` |
| `deal_pipelines` | Prospek, naskah proposal, dan tracking negosiasi | **RLS Enforced** | PK: `id`, FK: `org_id`, FK: `signal_id`, FK: `target_program_id` |

---

## 4. Query Execution Pattern (Go pgx/v5 Backend Wrapper)

Semua query pada tabel privat wajib dieksekusi melalui transaksi berpagar RLS:

```go
package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WithTenantContext executes database operations within an RLS-enforced PostgreSQL transaction.
func WithTenantContext(ctx context.Context, pool *pgxpool.Pool, orgID string, fn func(tx pgx.Tx) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	// Set tenant context for PostgreSQL RLS policies
	_, err = tx.Exec(ctx, "SET LOCAL app.current_org_id = $1", orgID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("failed to set tenant RLS context: %w", err)
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
```
