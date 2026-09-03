-- Migration 000015 Down: Drop ESG tables

DROP TABLE IF EXISTS company_esg_material_topics;
DROP TABLE IF EXISTS esg_material_topics;
DROP TABLE IF EXISTS company_esg_profiles;
