# Intent Score & AI Scoring Pipeline Specification

Document Version: `1.0.0`  
Last Updated: `2026-09-02`  
Target Module: `Sovera Core API & Intelligence Engine`

---

## 1. Overview

**Intent Score** ($0 - 100$) adalah metrik kuantitatif berbasis AI yang merepresentasikan tingkat kesiapan (*fundraising readiness*) dan besarnya peluang kemitraan CSR/TJSL/Zakat Korporasi dari suatu emiten. 

Setiap sinyal intelijen korporasi yang diekstrak dari Laporan Keberlanjutan (POJK 51), Keterbukaan Informasi BEI, atau berita resmi korporasi dihitung skornya secara otomatis oleh **Sovera AI Extractor Engine** (menggunakan Google Gemini LLM & Heuristic Scoring Pipeline).

---

## 2. Intent Score Formula & Weighting Matrix

Perhitungan `intent_score` didasarkan pada **4 Indikator Kunci** dengan alokasi bobot sebagai berikut:

$$\text{Intent Score} = (w_1 \times S_{\text{budget}}) + (w_2 \times S_{\text{urgency}}) + (w_3 \times S_{\text{esg}}) + (w_4 \times S_{\text{location}})$$

| Indikator Penilaian | Simbol | Bobot ($w$) | Kriteria Evaluasi Gemini AI |
| :--- | :---: | :---: | :--- |
| **Ketersediaan Anggaran** | $S_{\text{budget}}$ | **35%** | Apakah emiten mencantumkan alokasi nilai nominal rupiah secara eksplisit? (Misal: *"Menganggarkan TJSL Rp 25 M"* memberi skor maksimal dibanding naskah komitmen tanpa nominal). |
| **Urgensi Peristiwa Pemicu** | $S_{\text{urgency}}$ | **25%** | Apakah sinyal didorong oleh peristiwa resmi terkini? (Publikasi SR POJK 51, RUPS, atau Peluncuran Program Baru). |
| **Keselarasan Pilar ESG & Filantropi** | $S_{\text{esg}}$ | **25%** | Seberapa spesifik program menyasar 8 Asnaf / SDG (Pendidikan 3T, UMKM Syariah, Kesehatan, Lingkungan)? |
| **Kejelasan Wilayah Target** | $S_{\text{location}}$ | **15%** | Apakah emiten menyebutkan lokasi intervensi spesifik (misal: daerah 3T, Jawa Barat, NTT) dibanding hanya skala nasional umum? |

---

## 3. Detailed Scoring Rules & Rubric

### 3.1. Budget Signal Scoring ($S_{\text{budget}}$ - 35%)
- **Skor 100**: Nominal anggaran dicantumkan secara eksplisit $> \text{Rp } 10 \text{ Miliar}$.
- **Skor 75**: Nominal anggaran dicantumkan $\text{Rp } 1 \text{ Miliar} - \text{Rp } 10 \text{ Miliar}$.
- **Skor 50**: Nominal anggaran dicantumkan $< \text{Rp } 1 \text{ Miliar}$.
- **Skor 25**: Menyebutkan alokasi anggaran CSR tetapi tanpa mencantumkan angka nominal.
- **Skor 0**: Tidak ada informasi anggaran sama sekali.

### 3.2. Urgency & Trigger Event Scoring ($S_{\text{urgency}}$ - 25%)
- **Skor 100**: Terbit dari dokumen resmi POJK 51 / Keterbukaan Informasi BEI $\le 30 \text{ hari}$.
- **Skor 75**: Pengumuman RUPS / Press Release resmi korporasi $\le 60 \text{ hari}$.
- **Skor 50**: Berita media nasional terverifikasi $\le 90 \text{ hari}$.
- **Skor 25**: Laporan arsip tahun lalu ($> 90 \text{ hari}$).

### 3.3. ESG & Philanthropy Alignment ($S_{\text{esg}}$ - 25%)
- **Skor 100**: Sangat selaras dengan pilar prioritas Lembaga Zakat/Filantropi (Beasiswa 3T, Pesantren Digital, UMKM Dhuafa, Konservasi Air/Mangrove).
- **Skor 60**: Program CSR umum (Penghijauan umum, bantuan insidental).
- **Skor 20**: Program internal/komersial korporasi.

### 3.4. Geographic Specificity ($S_{\text{location}}$ - 15%)
- **Skor 100**: Menyebutkan kabupaten/kota/provinsi spesifik (misal: NTT, Pasuruan, Halmahera).
- **Skor 50**: Menyebutkan wilayah regional luas (misal: Indonesia Timur, Pulau Jawa).
- **Skor 20**: Tidak menyebutkan lokasi (hanya skala nasional).

---

## 4. LLM Prompt Template for Gemini AI

Berikut adalah prompt sistem yang digunakan oleh Sovera AI Extractor:

```text
You are an expert Corporate CSR & ESG Intelligence Analyst for Islamic Philanthropy.
Analyze the following corporate report text and extract structured signal JSON.

Calculate intent_score (integer 0-100) using this exact weighted formula:
- Budget Signal (35% weight): Is explicit IDR budget mentioned?
- Trigger Event Urgency (25% weight): Is this a formal POJK 51 / IDX report or recent RUPS?
- ESG & Philanthropy Alignment (25% weight): Does it align with Education, Poverty Alleviation (8 Asnaf), or SDG 4/8/9/17?
- Geographic Specificity (15% weight): Are specific target regions mentioned?

Return JSON output with schema:
{
  "company_name": "string",
  "industry_sector": "string",
  "source_type": "BEI_REPORT | NEWS | CSR_PDF",
  "summary": "string",
  "extracted_pillar": "string",
  "target_regions": ["string"],
  "estimated_budget_signal": number,
  "intent_score": number
}
```

---

## 5. UI Classification Thresholds in Frontend

Di dasbor frontend (`sovera-web-dashboard`), `intent_score` dikelompokkan ke dalam 3 kualifikasi visual:

| Rentang Skor | Kategori | Warna Badge UI | Rekomendasi Tindakan Fundraiser |
| :---: | :---: | :---: | :--- |
| **$\ge 80$** | **High Intent** | 🟢 Emerald (`bg-emerald-100`) | Prospek sangat matang. Langsung jalankan *AI Match Program* & buat proposal. |
| **$50 - 79$** | **Medium Intent** | 🟡 Amber (`bg-amber-100`) | Prospek potensial. Masukkan ke tahap *Research* untuk pendalaman PIC. |
| **$< 50$** | **Low Intent** | ⚪ Slate (`bg-slate-100`) | Sinyal eksploratif awal. Pantau pembaruan laporan berikutnya. |
