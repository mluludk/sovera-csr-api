# API Contract Specification

Dokumen ini berisi spesifikasi lengkap API Contract untuk **WebScraper Service (v1 & v2)**. Semua request yang dilindungi wajib menyertakan header otentikasi.

---

## 🔑 Authentication & Standards

### Header Otentikasi
```http
Authorization: Bearer <API_AUTH_KEY>
Content-Type: application/json
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

### 1. Universal Ingestion Task (v2 — CSR & B2B Document Intelligence)

Membuat pekerjaan ekstraksi dokumen PDF (*Sustainability Report*), portal berita, rilis bursa (IDX/RSS), atau postingan media sosial korporasi.

* **Method & Path:** `POST /api/v1/scrape-tasks`
* **Content-Type:** `application/json`

#### **Request Body**
```json
{
  "task_id": "job_csr_idx_991823",
  "client_origin": "sovera_b2b_engine",
  "source_type": "PDF_DOCUMENT",
  "target_url": "https://example-corpo.com/reports/sustainability-report-2025.pdf",
  "callback_url": "https://engine.domain.com/webhooks/crawler-response",
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
| `client_origin` | `string` | Ya | Identifikasi asal sistem pemanggil (misal: `sovera_b2b_engine`) |
| `source_type` | `string` | Ya | Pilihan: `PDF_DOCUMENT`, `NEWS_ARTICLE`, `RAW_WEB`, `SOCIAL_POST` |
| `target_url` | `string` | Ya | URL target dokumen PDF / Artikel / Web / Postingan |
| `callback_url` | `string` | Ya | URL Webhook penerima hasil callback |
| `config.render_js` | `boolean` | Tidak | Batasi render JS (Default: `false`) |
| `config.max_pages` | `integer` | Tidak | Batas maks halaman PDF (Default: `100`) |

#### **Success Response (`HTTP 202 Accepted`)**
```json
{
  "task_id": "job_csr_idx_991823",
  "status": "ACCEPTED"
}
```

#### **Webhook Callback Delivery Payload (`POST` to `callback_url`)**
```json
{
  "event": "csr.task_completed",
  "job_id": "job_csr_idx_991823",
  "target_type": "PDF_DOCUMENT",
  "platform": "CORPORATE_SITE",
  "status": "COMPLETED",
  "timestamp": "2026-08-31T16:20:00Z",
  "data": {
    "execution_time_ms": 1565,
    "content_hash": "ff67a9d764d6a2367a187734e697f6a53217db9a21c101d410a113ca871a299d",
    "source_metadata": {
      "source_type": "PDF_DOCUMENT",
      "target_url": "https://example-corpo.com/reports/sustainability-report-2025.pdf",
      "domain": "example-corpo.com",
      "platform": "CORPORATE_SITE",
      "published_date": "2026-04-12T00:00:00Z",
      "author_or_account": "PT Maju Sejahtera Tbk"
    },
    "raw_data": {
      "title": "Laporan Keberlanjutan 2025",
      "raw_text": "Teks mentah hasil ekstraksi...",
      "markdown_content": "# Laporan Keberlanjutan 2025\n\n...",
      "media_urls": [],
      "attachments": [
        {
          "file_name": "sustainability-report-2025.pdf",
          "file_size_bytes": 14205820
        }
      ],
      "engagement_metrics": {
        "likes_count": 0,
        "comments_count": 0,
        "shares_count": 0
      }
    }
  }
}
```

#### **Get Task Status & Details (v2 Polling API)**
Mengambil status dan informasi task ingestion CSR berbasis `task_id`.

* **Method & Path:** `GET /api/v1/tasks/{task_id}`

##### **Success Response (`HTTP 200 OK`)**
```json
{
  "task_id": "job_csr_idx_991823",
  "client_origin": "sovera_b2b_engine",
  "source_type": "PDF_DOCUMENT",
  "target_url": "https://example-corpo.com/reports/sustainability-report-2025.pdf",
  "callback_url": "https://engine.domain.com/webhooks/crawler-response",
  "status": "COMPLETED",
  "execution_time_ms": 1565,
  "content_hash": "ff67a9d764d6a2367a187734e697f6a53217db9a21c101d410a113ca871a299d",
  "created_at": "2026-08-31T16:20:00Z",
  "updated_at": "2026-08-31T16:20:01Z"
}
```

---

### 2. Marketplace & Social Search (v1)

Membuat pekerjaan scraping kata kunci produk e-commerce atau profil/post medsos.

* **Method & Path:** `POST /api/v1/search`
* **Content-Type:** `application/json`

#### **Request Body**
```json
{
  "platform": "tokopedia",
  "target_type": "marketplace_search",
  "keyword": "sepatu running",
  "limit": 50,
  "webhook_url": "https://client.domain.com/webhook"
}
```

| Parameter | Tipe | Wajib | Keterangan |
| :--- | :--- | :---: | :--- |
| `platform` | `string` | Ya | `tokopedia`, `shopee`, `lazada`, `tiktokshop`, `google`, `tiktok` |
| `target_type` | `string` | Tidak | `marketplace_search` (default), `social_profile`, `social_post`, `social_search` |
| `keyword` | `string` | Kondisional | Wajib untuk pencarian kata kunci (`marketplace_search`, `social_search`) |
| `username` | `string` | Kondisional | Wajib untuk `social_profile` dan `social_post` |
| `limit` | `integer` | Ya | Jumlah maksimum item (1 - 1000) |
| `webhook_url` | `string` | Tidak | URL Webhook penerima callback saat job selesai |

#### **Success Response (`HTTP 202 Accepted`)**
```json
{
  "job_id": "01JGH8X...",
  "status": "queued"
}
```

---

### 3. Get Job Status (v1)

Mengambil status terkini pekerjaan scraping.

* **Method & Path:** `GET /api/v1/jobs/{job_id}`

#### **Success Response (`HTTP 200 OK`)**
```json
{
  "job_id": "01JGH8X...",
  "platform": "tokopedia",
  "target_type": "marketplace_search",
  "keyword": "sepatu running",
  "requested_limit": 50,
  "status": "completed",
  "created_at": "2026-08-31T10:00:00Z",
  "started_at": "2026-08-31T10:00:01Z",
  "completed_at": "2026-08-31T10:00:05Z"
}
```

---

### 4. Get Job Results (v1 — dengan Pagination & Export)

Mengambil hasil scraping produk atau data sosial.

* **Method & Path:** `GET /api/v1/jobs/{job_id}/results`
* **Query Parameters:**
  * `format`: `json` (default), `csv`, `xlsx`
  * `limit`: Jumlah data per halaman (opsional)
  * `offset`: Indeks awal baris data (opsional)

#### **Success Response JSON (`HTTP 200 OK`)**
```json
{
  "job_id": "01JGH8X...",
  "platform": "tokopedia",
  "items": [
    {
      "platform": "tokopedia",
      "productId": "9812301",
      "name": "Sepatu Running Air",
      "url": "https://tokopedia.com/...",
      "price": 450000,
      "seller": {
        "id": "3131144",
        "slug": "makassar_shop",
        "name": "Makassar Shop",
        "location": "Kota Makassar",
        "badge": "Power Merchant Pro"
      },
      "position": 1
    }
  ],
  "total": 50,
  "limit": 10,
  "offset": 0
}
```

---

### 5. Health & Readiness Check

Mengecek kesehatan server API, koneksi PostgreSQL, dan Redis.

* **Method & Path:** `GET /health` atau `GET /ready`

#### **Success Response (`HTTP 200 OK`)**
```json
{
  "status": "ok"
}
```
