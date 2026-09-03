-- Migration 000012 Up: Create csr_focuses table and seed 11 focus areas

CREATE TABLE IF NOT EXISTS csr_focuses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_csr_focuses_code ON csr_focuses(code);
CREATE INDEX IF NOT EXISTS idx_csr_focuses_category ON csr_focuses(category);

-- Seed initial 11 CSR focus taxonomy codes
INSERT INTO csr_focuses (code, name, category, description) VALUES
('EDUCATION', 'Pendidikan & Literasi', 'Pendidikan', 'Program beasiswa, pembenahan fasilitas sekolah, pelatihan tenaga pengajar, dan peningkatan literasi digital.'),
('HEALTH', 'Kesehatan & Gizi Masyarakat', 'Kesehatan', 'Pencegahan stunting, layanan kesehatan gratis, fasilitas medis dasar, dan edukasi pola hidup sehat.'),
('HUMANITARIAN', 'Bantuan Kemanusiaan & Karitatif', 'Kemanusiaan', 'Penyaluran bantuan sosial bagi masyarakat marjinal, penanganan krisis kemanusiaan, dan dukungan pengungsi.'),
('ECONOMIC_EMPOWERMENT', 'Pemberdayaan Ekonomi & UMKM', 'Ekonomi', 'Pelatihan kewirausahaan, pendampingan UMKM, permodalan usaha mikro, dan akses pasar lokal.'),
('ENVIRONMENT', 'Lingkungan & Keanekaragaman Hayati', 'Lingkungan', 'Penanaman pohon dan mangrove, pengolahan sampah/limbah, energi terbarukan, dan reboisasi hutan.'),
('DISASTER', 'Penanggulangan & Tanggap Bencana', 'Kemanusiaan', 'Bantuan logistik darurat bencana, pemulihan sarana fisik pascabencana, serta simulasi mitigasi risiko.'),
('WATER_SANITATION', 'Air Bersih & Sanitasi (WASH)', 'Infrastruktur & Kesehatan', 'Pembangunan sumur air bersih, penyediaan sarana MCK layak, serta edukasi higiene masyarakat.'),
('FOOD_SECURITY', 'Ketahanan Pangan & Pertanian', 'Pertanian & Gizi', 'Dukungan bagi kelompok tani lokal, inovasi urban farming, dan penyediaan bank makanan.'),
('CHILDREN', 'Perlindungan & Kesejahteraan Anak', 'Sosial', 'Pemenuhan hak anak, dukungan panti asuhan, dan penyediaan sarana ramah anak.'),
('WOMEN', 'Pemberdayaan Perempuan & Kesetaraan', 'Sosial & Ekonomi', 'Pelatihan keterampilan wanita karir/ibu rumah tangga, pendampingan usaha perempuan, dan kesetaraan gender.'),
('DISABILITY', 'Inklusi & Pemberdayaan Disabilitas', 'Sosial', 'Penyediaan sarana ramah disabilitas, pelatihan keterampilan kerja inklusif, dan pembagian alat bantu mobilitas')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    category = EXCLUDED.category,
    description = EXCLUDED.description,
    updated_at = NOW();
