# Product Requirement Document (PRD) - Sovera Multi-Sector Expansion

**Product Name:** Sovera / FundIQ  
**Target Sector:** Seluruh Lembaga Kemanusiaan (NGO, Yayasan Kebencanaan, Lingkungan Hidup, Pendidikan, Kesehatan, Pemberdayaan Masyarakat, serta Lembaga Zakat/Wakaf)  
**Document Version:** 2.0 (Universal Humanitarian Expansion)  
**Document Status:** Approved for Implementation  

---

## 1. Background & Business Objective

### 1.1 Latar Belakang
Pada versi awal, Sovera dirancang dengan taksonomi spesifik filantropi Islam (8 Asnaf). Namun, kebutuhan riset intelijen B2B untuk menyerap anggaran CSR, TJSL BUMN, hibah korporasi, dan Zakat Perusahaan dirasakan oleh seluruh ekosistem lembaga nirlaba. 

Korporasi dan emiten bursa mengalokasikan anggaran sosial mereka berbasis pilar **ESG (Environmental, Social, Governance)** dan **UN Sustainable Development Goals (SDGs)**, bukan semata-mata nomenklatur keagamaan.

### 1.2 Tujuan Utama (Objectives)
1. **Perluasan Total Addressable Market (TAM):** Membuka akses platform bagi ribuan NGO nasional dan internasional yang beroperasi di Indonesia.
2. **Standardisasi Taksonomi Universal:** Mengadopsi pilar SDG dan Klaster Kemanusiaan Standar Nasional/Internasional sebagai jembatan semantik utama antara program lembaga dan fokus korporasi.
3. **Konfigurasi Dinamis Berbasis Tenant (`org_type`):** Antarmuka, taksonomi form, dan prompt AI secara otomatis beradaptasi dengan identitas masing-masing lembaga (misal: tetap menampilkan opsi 8 Asnaf jika profilnya lembaga zakat, dan menyembunyikannya jika NGO kemanusiaan umum).

---

## 2. Target User Personas & Organization Types

### 2.1 Klasifikasi Organisasi (`organization_type`)
* `HUMANITARIAN_NGO`: Yayasan kemanusiaan umum, perlindungan anak, kelompok rentan, hak asasi, dan kesejahteraan sosial.
* `DISASTER_RELIEF`: Lembaga tanggap darurat, evakuasi, logistik kebencanaan, dan rekonstruksi pasca-bencana.
* `ENVIRONMENT_CONSERVATION`: Lembaga konservasi alam, reboisasi, energi terbarukan, dan adaptasi perubahan iklim.
* `HEALTH_EDUCATION`: Yayasan beasiswa, peningkatan literasi sekolah 3T, penanganan stunting, dan bantuan medis darurat.
* `ZAKAT_WAQF_INSTITUTION`: Lembaga Pengelola Zakat (BAZNAS/LAZNAS) dan Nazhir Wakaf (BMT/Yayasan Wakaf).
* `UNIVERSITY_ENDOWMENT`: Dana abadi perguruan tinggi dan lembaga riset independen.

### 2.2 User Personas
* **Director of Institutional Partnership / Kemitraan Strategis:** Membutuhkan ringkasan analitik portofolio corporate funding dan konversi pipeline.
* **Corporate Fundraiser / Account Executive NGO:** Memerlukan rekomendasi emiten yang memiliki pemicu (*trigger event*) selaras dengan isu intervensi lembaga mereka.
* **Proposal & Program Development Specialist:** Membutuhkan alat bantu generatif untuk menyusun Term of Reference (TOR), proposal kemitraan, dan Rencana Anggaran Biaya (RAB) dampak.

---

## 3. Taxonomy & Data Schema Expansion

### 3.1 Universal Sector & Intervention Framework
Program lembaga tidak lagi dikunci oleh kolom fikih. Taksonomi program mengadopsi 3 layer:
1. **Primary Cluster (Klaster Intervensi):**
   * Bencana & Tanggap Darurat (*Disaster & Emergency*)
   * Pendidikan & Literasi Anak (*Education & Literacy*)
   * Kesehatan, Gizi & Sanitasi (*Health & WASH*)
   * Pemberdayaan Ekonomi & UMKM (*Economic Empowerment*)
   * Kelestarian Lingkungan & Iklim (*Climate & Environment*)
   * Perlindungan Sosial & Difabel (*Social Protection & Vulnerable Groups*)
   * Fikih Zakat / Wakaf (*Spesifik jika `org_type` = ZAKAT_WAQF_INSTITUTION*)
2. **SDG Mapping:** SDG 1 hingga SDG 17.
3. **ESG Alignment:** Environmental (E), Social (S), Governance (G).

### 3.2 Database Schema Migration
Penyesuaian pada skema PostgreSQL:

```sql
-- 1. Penambahan org_type pada tabel organizations
ALTER TABLE organizations 
ADD COLUMN IF NOT EXISTS org_type VARCHAR(50) DEFAULT 'HUMANITARIAN_NGO';

-- 2. Generalisasi tabel institution_programs
ALTER TABLE institution_programs
ADD COLUMN IF NOT EXISTS primary_cluster VARCHAR(100) NOT NULL DEFAULT 'COMMUNITY_DEVELOPMENT',
ADD COLUMN IF NOT EXISTS target_sdgs TEXT[] DEFAULT '{}',
ADD COLUMN IF NOT EXISTS esg_pillar VARCHAR(20) DEFAULT 'SOCIAL',
ALTER COLUMN asnaf_category DROP NOT NULL;
```

---

## 4. Key Functional Requirements

### 4.1 Multi-Source Corporate Intelligence Feed (`/signals`)
* **Pencarian & Filter Universal:** Pengguna dapat memfilter sinyal berdasarkan Klaster Kemanusiaan, Target SDG, dan Sektor Industri Emiten.
* **Intent Scoring Engine:** AI mendeteksi alokasi dana korporasi untuk:
  - Program Bantuan Bencana Spontan (Respons Cepat).
  - Program TJSL/CSR Berkelanjutan (Multi-years).
  - Dana Hibah Riset/Pendidikan.
  - Zakat Perusahaan / Infaq Korporasi.

### 4.2 Universal Semantic Program Matcher
* **Dual-Vector / Hybrid Matching:**
  - Program lembaga diubah menjadi embedding 1536 dimensi yang merepresentasikan deskripsi masalah, wilayah intervensi, target penerima, dan dampak terukur.
  - Pencocokan kosinus mengukur kedekatan semantik antara laporan keberlanjutan emiten dan program lembaga, tanpa memedulikan perbedaan terminologi agama/sekuler.

### 4.3 Adaptive Proposal Studio (`/pipeline/:dealId`)
Proposal Studio mengadaptasi struktur naskah dan *tone of voice* berdasarkan tipe lembaga:
* **Mode Kemanusiaan & Kebencanaan:** Menitikberatkan pada rasio kecepatan respons, akuntabilitas logistik, transparansi data penerima manfaat langsung, dan dokumentasi lapangan.
* **Mode Pemberdayaan & ESG:** Menitikberatkan pada Social Return on Investment (SROI), indikator SDG, dan kesesuaian dengan matriks POJK 51 / Global Reporting Initiative (GRI).
* **Mode Zakat Korporasi:** Menitikberatkan pada akad syariah, penyaluran 8 Asnaf, dan legalitas izin amil resmi.

---

## 5. Non-Functional Requirements & Security

* **Multi-Tenant Data Isolation:** RLS di kernel PostgreSQL tetap membatasi data proposal, catatan internal, dan draf program agar hanya dapat dibaca oleh staf organisasi terkait (`SET LOCAL app.current_org_id = :orgId`).
* **Zero AI Contamination:** Data proposal sensitif dari seluruh NGO tidak akan digunakan untuk pelatihan ulang model AI publik.
* **Interoperabilitas Ekspor:** Dukungan ekspor berkas naskah proposal ke `.docx` dan `.pdf` dengan kop surat dan logo kustom masing-masing organisasi.

---

## 6. Success Metrics (KPIs)

1. **User Acquisition:** Minimal 30% dari tenant aktif berasal dari NGO non-keagamaan dalam 6 bulan pertama rilis.
2. **Match Relevance:** Rasio kecocokan program (*similarity match*) di atas 75% dinilai akurat dan relevan oleh Corporate Fundraiser lintas sektor.
3. **Time to Pitch:** Mempersingkat waktu riset prospek dan pembuatan draf proposal dari 3 hari menjadi kurang dari 20 menit.