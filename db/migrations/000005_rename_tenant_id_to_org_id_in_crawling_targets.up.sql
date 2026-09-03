-- Migration 000005: Rename tenant_id to org_id in crawling_targets for naming convention consistency

-- 1. Drop old RLS policy referencing tenant_id
DROP POLICY IF EXISTS crawling_targets_tenant_isolation ON crawling_targets;

-- 2. Rename column tenant_id → org_id (safely handling pre-existing org_id)
DO $$ 
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'crawling_targets' AND column_name = 'tenant_id'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'crawling_targets' AND column_name = 'org_id'
    ) THEN
        ALTER TABLE crawling_targets RENAME COLUMN tenant_id TO org_id;
    ELSIF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'crawling_targets' AND column_name = 'tenant_id'
    ) AND EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'crawling_targets' AND column_name = 'org_id'
    ) THEN
        ALTER TABLE crawling_targets DROP COLUMN tenant_id;
    END IF;
END $$;

-- 3. Drop old index and recreate with correct name
DROP INDEX IF EXISTS idx_crawling_targets_tenant;
CREATE INDEX IF NOT EXISTS idx_crawling_targets_org ON crawling_targets (org_id);

-- 4. Recreate RLS policy using org_id
DROP POLICY IF EXISTS crawling_targets_org_isolation ON crawling_targets;
CREATE POLICY crawling_targets_org_isolation ON crawling_targets
    FOR ALL
    USING (
        org_id IS NULL
        OR org_id = NULLIF(current_setting('app.current_tenant_id', true), '')::uuid
    );
