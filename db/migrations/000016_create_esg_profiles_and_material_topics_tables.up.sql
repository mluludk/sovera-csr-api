-- Migration 000015 Up: Create company_esg_profiles, esg_material_topics, and company_esg_material_topics tables with seed data

CREATE TABLE IF NOT EXISTS company_esg_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    reporting_year SMALLINT NOT NULL,
    report_date DATE,
    overall_score NUMERIC(5,2),
    environmental_score NUMERIC(5,2),
    social_score NUMERIC(5,2),
    governance_score NUMERIC(5,2),
    esg_rating VARCHAR(50),
    sustainability_strategy TEXT,
    sdg_alignment JSONB DEFAULT '{}'::jsonb,
    source_id UUID,
    confidence NUMERIC(5,2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, reporting_year)
);

CREATE TABLE IF NOT EXISTS esg_material_topics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS company_esg_material_topics (
    esg_profile_id UUID NOT NULL REFERENCES company_esg_profiles(id) ON DELETE CASCADE,
    topic_id UUID NOT NULL REFERENCES esg_material_topics(id) ON DELETE CASCADE,
    materiality_score NUMERIC(5,2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (esg_profile_id, topic_id)
);

CREATE INDEX IF NOT EXISTS idx_company_esg_profiles_company_id ON company_esg_profiles(company_id);
CREATE INDEX IF NOT EXISTS idx_company_esg_profiles_year ON company_esg_profiles(reporting_year);
CREATE INDEX IF NOT EXISTS idx_esg_material_topics_code ON esg_material_topics(code);
CREATE INDEX IF NOT EXISTS idx_esg_material_topics_category ON esg_material_topics(category);
CREATE INDEX IF NOT EXISTS idx_company_esg_material_topics_profile ON company_esg_material_topics(esg_profile_id);

-- Seed ESG Material Topics (GRI & SASB standard taxonomy)
INSERT INTO esg_material_topics (code, name, category, description) VALUES
-- Environmental (E)
('CLIMATE_CHANGE', 'Perubahan Iklim & Emisi GRK', 'ENVIRONMENTAL', 'Manajemen krisis iklim, pengurangan emisi Gas Rumah Kaca (Scope 1, 2, 3), dan target Net-Zero.'),
('CARBON_EMISSIONS', 'Jejak Karbon & Efisiensi Energi', 'ENVIRONMENTAL', 'Pengisian daya hijau, efisiensi energi operasional, dan transisi ke Energi Baru Terbarukan (EBT).'),
('BIODIVERSITY', 'Keanekaragaman Hayati & Lahan', 'ENVIRONMENTAL', 'Perlindungan ekosistem lokal, restorasi hutan/mangrove, dan pencegahan degradasi keanekaragaman hayati.'),
('WATER_MANAGEMENT', 'Manajemen Air & Limbah Cair', 'ENVIRONMENTAL', 'Penghematan penggunaan air baku, daur ulang air operasional, dan pengelolaan limbah cair cair.'),
('WASTE_CIRCULARITY', 'Ekonomi Sirkular & Pengelolaan Sampah', 'ENVIRONMENTAL', 'Reduksi sampah plastik, daur ulang limbah padat, serta efisiensi sirkularitas rantai pasok.'),

-- Social (S)
('HUMAN_RIGHTS', 'Hak Asasi Manusia & Perlindungan Masyarakat', 'SOCIAL', 'Penegakan HAM di lingkungan kerja dan masyarakat terdampak operasional bisnis.'),
('LABOR_STANDARDS', 'Ketenagakerjaan & Kesejahteraan Pekerja', 'SOCIAL', 'Upah layak, jam kerja adil, kebebasan berserikat, dan pencegahan kerja paksa/anak.'),
('HEALTH_SAFETY', 'Kesehatan & Keselamatan Kerja (K3)', 'SOCIAL', 'Manajemen risiko kecelakaan kerja, sertifikasi K3, dan program kesehatan mental karyawan.'),
('COMMUNITY_RELATIONS', 'Hubungan Masyarakat & Dampak Lokal', 'SOCIAL', 'Program CSR berlanjut, konsultasi publik masyarakat adat/lokal, dan penyelesaian konflik sosial.'),
('DIVERSITY_INCLUSION', 'Keberagaman, Kesetaraan & Inklusi (DEI)', 'SOCIAL', 'Kesetaraan gender, representasi kelompok minoritas, serta akses inklusif pekerja disabilitas.'),

-- Governance (G)
('BUSINESS_ETHICS', 'Etika Bisnis & Kepatuhan Hukum', 'SOCIAL', 'Penegakan kode etik perusahaan, aturan persaingan sehat, dan kepatuhan regulasi publik.'),
('TRANSPARENCY_ANTI_CORRUPTION', 'Transparansi & Anti-Korupsi', 'GOVERNANCE', 'Kebijakan anti-suap/gratifikasi, sistem perlindungan whistleblower, dan keterbukaan publik.'),
('DATA_PRIVACY_SECURITY', 'Keamanan Data & Privasi Konsumen', 'GOVERNANCE', 'Perlindungan data pribadi (PDP), keamanan siber, dan transparansi tata kelola data pelanggan.'),
('BOARD_DIVERSITY_GOVERNANCE', 'Tata Kelola & Keberagaman Dewan', 'GOVERNANCE', 'Independensi komisaris, keberagaman dewan direksi, serta transparansi remunerasi eksekutif.')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    category = EXCLUDED.category,
    description = EXCLUDED.description,
    updated_at = NOW();
