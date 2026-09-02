-- Rollback migration 000007
DELETE FROM crawling_targets WHERE source_name IN (
    'IDX Keterbukaan Informasi Publik',
    'IDX Sustainability Reports Hub',
    'Bisnis.com RSS Feed - CSR & Finansial',
    'Kontan.co.id RSS Feed - Keuangan & Industri',
    'Antara News RSS - Warta Ekonomi & Humaniora',
    'Republika Online RSS - Filantropi & Zakat Korporasi',
    'Kementerian BUMN - Siaran Pers TJSL',
    'Pertamina Newsroom - CSR & TJSL Updates',
    'PLN Newsroom - Program CSR & Desa Berdaya',
    'Pelindo TJSL & Social Impact Newsroom',
    'UNDP Indonesia Procurement Notices & Grants',
    'Kedutaan Besar Australia - Direct Aid Program (DAP)',
    'Kedutaan Besar Jepang - Hibah Grassroots (GGP)'
);
