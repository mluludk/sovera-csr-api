-- Migration 000004: Universal Humanitarian Expansion Schema Update

-- 1. Add org_type column to organizations
ALTER TABLE organizations 
ADD COLUMN IF NOT EXISTS org_type VARCHAR(50) DEFAULT 'ZAKAT_WAQF_INSTITUTION';

-- 2. Generalize institution_programs for multi-sector NGOs
ALTER TABLE institution_programs
ADD COLUMN IF NOT EXISTS primary_cluster VARCHAR(100) NOT NULL DEFAULT 'COMMUNITY_DEVELOPMENT',
ADD COLUMN IF NOT EXISTS target_sdgs TEXT[] DEFAULT '{}',
ADD COLUMN IF NOT EXISTS esg_pillar VARCHAR(50) DEFAULT 'SOCIAL',
ALTER COLUMN asnaf_category DROP NOT NULL;

-- 3. Create Indexes for multi-sector queries
CREATE INDEX IF NOT EXISTS idx_organizations_org_type ON organizations (org_type);
CREATE INDEX IF NOT EXISTS idx_institution_programs_cluster ON institution_programs (primary_cluster);
