-- Migration 000007: Seed public crawling targets from CRAWL_SOURCE.md
-- org_id = NULL denotes Global Public Targets accessible by platform worker

INSERT INTO crawling_targets (source_name, source_type, target_url, check_interval_hours, is_active, org_id)
VALUES
    -- 1. Keterbukaan Informasi & Laporan Emiten Bursa (BEI / IDX)
    ('IDX Keterbukaan Informasi Publik', 'IDX_ANNOUNCEMENT', 'https://www.idx.co.id/id/berita/pengumuman', 6, true, NULL),
    ('IDX Sustainability Reports Hub', 'PDF_REPORTS', 'https://www.idx.co.id/id/perusahaan-tercatat/laporan-keberlanjutan', 24, true, NULL),

    -- 2. Portal Berita Finansial, Bisnis, dan CSR (RSS Feed)
    ('Bisnis.com RSS Feed - CSR & Finansial', 'NEWS_RSS', 'https://www.bisnis.com/rss', 4, true, NULL),
    ('Kontan.co.id RSS Feed - Keuangan & Industri', 'NEWS_RSS', 'https://www.kontan.co.id/rss', 4, true, NULL),
    ('Antara News RSS - Warta Ekonomi & Humaniora', 'NEWS_RSS', 'https://www.antaranews.com/rss/ekonomi', 4, true, NULL),
    ('Republika Online RSS - Filantropi & Zakat Korporasi', 'NEWS_RSS', 'https://www.republika.co.id/rss/ekonomi/syariah-filantropi', 4, true, NULL),

    -- 3. Portal Resmi BUMN & Press Release Korporasi
    ('Kementerian BUMN - Siaran Pers TJSL', 'BUMN_PORTAL', 'https://bumn.go.id/media/press-release', 12, true, NULL),
    ('Pertamina Newsroom - CSR & TJSL Updates', 'CORPORATE_NEWSROOM', 'https://www.pertamina.com/id/news-room', 12, true, NULL),
    ('PLN Newsroom - Program CSR & Desa Berdaya', 'CORPORATE_NEWSROOM', 'https://web.pln.co.id/media/siaran-pers', 12, true, NULL),
    ('Pelindo TJSL & Social Impact Newsroom', 'CORPORATE_NEWSROOM', 'https://pelindo.co.id/media/berita', 12, true, NULL),

    -- 4. Portal Pengadaan & Hibah Lembaga Donor / Kedutaan (Grants)
    ('UNDP Indonesia Procurement Notices & Grants', 'GRANTS_PORTAL', 'https://procurement-notices.undp.org/', 24, true, NULL),
    ('Kedutaan Besar Australia - Direct Aid Program (DAP)', 'GRANTS_PORTAL', 'https://indonesia.embassy.gov.au/jkti/dap.html', 24, true, NULL),
    ('Kedutaan Besar Jepang - Hibah Grassroots (GGP)', 'GRANTS_PORTAL', 'https://www.id.emb-japan.go.jp/ggp.html', 24, true, NULL)
ON CONFLICT DO NOTHING;
