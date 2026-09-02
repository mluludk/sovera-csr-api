-- Migration: Add tenant_id column and RLS policy to crawling_targets
-- NOTE: Column renamed to org_id for naming consistency in migration 000005

-- 1. Add tenant_id column (NULL = Public/Global target, UUID = Tenant-specific private target)
ALTER TABLE crawling_targets 
ADD COLUMN IF NOT EXISTS tenant_id UUID REFERENCES organizations(id) ON DELETE CASCADE;

-- 2. Index for tenant query optimization
CREATE INDEX IF NOT EXISTS idx_crawling_targets_tenant 
ON crawling_targets (tenant_id);

-- 3. Enable Row Level Security (RLS)
ALTER TABLE crawling_targets ENABLE ROW LEVEL SECURITY;

-- 4. Drop policy if exists to ensure idempotency
DROP POLICY IF EXISTS crawling_targets_tenant_isolation ON crawling_targets;

-- 5. Create Hybrid RLS Policy (Public Global Targets + Tenant Private Targets)
CREATE POLICY crawling_targets_tenant_isolation ON crawling_targets
    FOR ALL
    USING (
        tenant_id IS NULL 
        OR tenant_id = NULLIF(current_setting('app.current_tenant_id', true), '')::uuid
    );
