# API Contract Specification (WebScraper Service v2)

Dokumen ini berisi spesifikasi lengkap API Contract untuk **WebScraper Service (v1 & v2)**, termasuk penambahan dukungan pengayaan profil perusahaan (*Company Enrichment*) dan pencarian peluang CSR aktif (*CSR Opportunity Feed*). Semua request yang dilindungi wajib menyertakan header otentikasi.

---

## 🔑 Authentication & Standards

### Header Otentikasi API
```http
Authorization: Bearer <API_AUTH_KEY>
Content-Type: application/json
```

### Header Otentikasi Webhook Callback (SHA-256 HMAC)
Crawler service **wajib** menyertakan header tanda tangan digital pada setiap pengiriman callback ke backend Sovera:
```http
Content-Type: application/json
X-Hub-Signature-256: sha256=<hmac_hex_digest>
```

### Standard Error Response (HTTP 4xx / 5xx)
```json
{
  "error": "invalid_request",
  "message": "penjelasan detail kesalahan"
}
```

---

## 🛠️ Endpoints API

### 1. Universal Ingestion Task (`POST /api/v1/scrape-tasks`)

Membuat pekerjaan pemindaian dokumen PDF (*Sustainability Report*), portal berita, rilis bursa (IDX/RSS), pencarian peluang CSR aktif, atau pengayaan profil entitas korporasi.

* **Method & Path:** `POST /api/v1/scrape-tasks`
* **Content-Type:** `application/json`

#### **Request Body**
```json
{
  "task_id": "job_csr_idx_991823",
  "target_id": "a0ac2304-d0d4-4d14-a5db-ea22966c71c6",
  "client_origin": "sovera_b2b_engine",
  "source_type": "COMPANY_ENRICHMENT",
  "target_url": "https://telkom.co.id/csr",
  "callback_url": "https://api.sovera.id/api/v1/webhooks/crawler",
  "config": {
    "render_js": false,
    "bypass_anti_bot": true,
    "max_pages": 50
  }
}
```

| Parameter | Tipe | Wajib | Keterangan |
| :--- | :--- | :---: | :--- |
| `task_id` | `string` | Ya | ID Unik pekerjaan dari sistem pemanggil |
| `target_id` | `string` | Tidak | UUID Target dari `crawling_targets` (untuk pelacakan Circuit Breaker) |
| `client_origin` | `string` | Ya | Identifikasi asal sistem pemanggil (`sovera_b2b_engine`) |
| `source_type` | `string` | Ya | Pilihan: `PDF_DOCUMENT`, `NEWS_ARTICLE`, `NEWS_RSS`, `BUMN_PORTAL`, `GRANTS_PORTAL`, **`CSR_OPPORTUNITY_SEARCH`**, **`COMPANY_ENRICHMENT`** |
| `target_url` | `string` | Ya | URL target dokumen PDF / Artikel / Portal Web / Halaman Kontak CSR Perusahaan |
| `callback_url` | `string` | Ya | URL Webhook tunggal (*single callback*) penerima hasil callback di backend |
| `config.render_js` | `boolean` | Tidak | Batasi render Headless JS Browser (Default: `false`) |
| `config.max_pages` | `integer` | Tidak | Batas maks halaman PDF (Default: `100`) |

#### **Success Response (`HTTP 202 Accepted`)**
```json
{
  "task_id": "job_csr_idx_991823",
  "status": "ACCEPTED"
}
```

---

### 2. Single Webhook Callback Delivery Specification (`POST` to `callback_url`)

Crawler service wajib mengirimkan hasil pemindaian ke **`callback_url` yang tunggal** dengan header otentikasi `X-Hub-Signature-256`. Atribut `source_type` menentukan kategori data yang diproses backend.

#### **A. Payload Callback: Laporan PDF & Berita CSR (`source_type: "PDF_DOCUMENT"` / `"NEWS_ARTICLE"`)**
```json
{
  "task_id": "job_csr_idx_991823",
  "target_id": "a0ac2304-d0d4-4d14-a5db-ea22966c71c6",
  "status": "COMPLETED",
  "http_status_code": 200,
  "error_message": "",
  "source_type": "PDF_DOCUMENT",
  "source_url": "https://example-corpo.com/reports/sustainability-report-2025.pdf",
  "author_or_account": "PT Maju Sejahtera Tbk",
  "published_date": "2026-04-12T00:00:00Z",
  "raw_text": "PT Maju Sejahtera Tbk mengalokasikan dana TJSL sebesar Rp 25 Miliar...",
  "markdown_content": "# Laporan Keberlanjutan 2025\n\nRealisasi pilar pendidikan...",
  "execution_time_ms": 1565
}
```

#### **B. Payload Callback: Peluang Hibah CSR Aktif (`source_type: "CSR_OPPORTUNITY_SEARCH"`)**
```json
{
  "task_id": "job_opp_idx_881204",
  "target_id": "b1bd3405-e1e5-5e25-b6fc-fb33077d82d7",
  "status": "COMPLETED",
  "http_status_code": 200,
  "error_message": "",
  "source_type": "CSR_OPPORTUNITY_SEARCH",
  "source_url": "https://yayasan.djarumfoundation.org/call-for-proposals-2026",
  "author_or_account": "Djarum Foundation",
  "published_date": "2026-09-01T00:00:00Z",
  "raw_text": "Djarum Foundation membuka pendaftaran proposal hibah program konservasi air...",
  "markdown_content": "# Panggilan Proposal Hibah Konservasi Air 2026\n\nSyarat dan ketentuan...",
  "execution_time_ms": 1820
}
```

#### **C. Payload Callback: Pengayaan Profil Perusahaan (`source_type: "COMPANY_ENRICHMENT"`)**
```json
{
  "task_id": "job_enrich_telkom_3391",
  "target_id": "c2ce4516-f2f6-6f36-c7ad-gc44188e93e8",
  "status": "COMPLETED",
  "http_status_code": 200,
  "error_message": "",
  "source_type": "COMPANY_ENRICHMENT",
  "source_url": "https://telkom.co.id/tjsl-contact",
  "author_or_account": "PT Telkom Indonesia (Persero) Tbk",
  "published_date": "2026-09-03T00:00:00Z",
  "raw_text": "Kontak Departemen TJSL Telkom: tjsl@telkom.co.id. Kantor Pusat: Jl. Japati No. 1 Bandung...",
  "markdown_content": "## Profil TJSL & ESG Telkom Indonesia\n\n- Email Public CSR: tjsl@telkom.co.id\n- Alamat HQ: Bandung\n- Fokus: Digitalisasi Pendidikan & EBT",
  "execution_time_ms": 2100
}
```

#### **D. Payload Callback Gagal / Dead Link (`status: "FAILED"`)**
```json
{
  "task_id": "job_csr_idx_991823",
  "target_id": "a0ac2304-d0d4-4d14-a5db-ea22966c71c6",
  "status": "FAILED",
  "http_status_code": 404,
  "error_message": "HTTP 404 Not Found - Target URL does not exist or has been removed",
  "source_type": "PDF_DOCUMENT",
  "source_url": "https://example-corpo.com/reports/sustainability-report-2025.pdf",
  "execution_time_ms": 450
}
```

---

### 3. Task Status Polling API (`GET /api/v1/tasks/{task_id}`)

Mengambil status dan detail pekerjaan pemindaian berbasis `task_id`.

* **Method & Path:** `GET /api/v1/tasks/{task_id}`

#### **Success Response (`HTTP 200 OK`)**
```json
{
  "task_id": "job_csr_idx_991823",
  "target_id": "a0ac2304-d0d4-4d14-a5db-ea22966c71c6",
  "client_origin": "sovera_b2b_engine",
  "source_type": "COMPANY_ENRICHMENT",
  "target_url": "https://telkom.co.id/csr",
  "callback_url": "https://api.sovera.id/api/v1/webhooks/crawler",
  "status": "COMPLETED",
  "http_status_code": 200,
  "execution_time_ms": 1565,
  "content_hash": "ff67a9d764d6a2367a187734e697f6a53217db9a21c101d410a113ca871a299d",
  "created_at": "2026-09-03T16:20:00Z",
  "updated_at": "2026-09-03T16:20:01Z"
}
```

---

### 4. Health & Readiness Check

Mengecek kesehatan server WebScraper API.

* **Method & Path:** `GET /health` atau `GET /ready`

#### **Success Response (`HTTP 200 OK`)**
```json
{
  "status": "ok"
}
```
