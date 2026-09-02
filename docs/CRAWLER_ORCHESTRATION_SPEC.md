# Crawler Orchestration & Scheduler Specification (CRAWLER_ORCHESTRATION_SPEC.md)

**Product:** Sovera (FundIQ) Core API & Ingestion Engine  
**Component:** Backend Orchestrator & External Scraper Dispatcher  
**Target Runtime:** Go 1.22+ (`sovera-core-api`), Redis (Asynq Scheduler), WebScraper Service  
**Document Version:** 1.0  

---

## 1. Overview & Architectural Philosophy

Dokumen ini mendefinisikan arsitektur dan spesifikasi teknis untuk **Orkestrasi & Penjadwalan Pemindaian Data** (*Crawler Scheduling & Dispatching Pattern*).

### Prinsip Utama
1. **Decoupled Stateless Scraper:** *WebScraper Service* berperan sebagai *Execution Engine* murni yang bersifat *stateless* (mengeksekusi tugas crawling/parsing PDF/Web berbasis HTTP request). Scraper **tidak** mengelola state database target atau jadwal cron internal.
2. **Backend Engine as Orchestrator:** *Sovera Core Backend API* bertindak sebagai *Orchestrator & Scheduler*. Backend menyimpan registri situs target, mengelola jadwal eksekusi (Asynq Cron Scheduler), mentrigger scraper, dan menerima *callback webhook* hasil ekstraksi.
3. **Idempotent Ingestion & Deduplication:** Setiap hasil crawling diverifikasi menggunakan tanda tangan HMAC SHA-256 dan di-deduplikasi berbasis hash konten (`content_hash`) untuk menghemat penggunaan kuota AI LLM.

---

## 2. Architecture & Data Flow Diagram

```text
┌─────────────────────────────────────────────────────────────────────────┐
│                    Sovera Core API (Cron Orchestrator)                  │
│                                                                         │
│   ┌──────────────────────────┐         ┌─────────────────────────────┐  │
│   │ PostgreSQL Database      │         │ Asynq Cron Scheduler        │  │
│   │ - `crawling_targets`     │────────►│ - Periodic Tick (e.g. Daily)│  │
│   │ - `crawling_logs`        │         │ - Fetch Due Targets         │  │
│   └──────────────────────────┘         └──────────────┬──────────────┘  │
└───────────────────────────────────────────────────────┼─────────────────┘
                                                        │ HTTP POST Request
                                                        │ (Target URL + Callback URL)
                                                        ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                     External Scraper Service Engine                     │
│                                                                         │
│   ┌─────────────────────────────────────────────────────────────────┐   │
│   │ WebScraper Service (`SCRAPER_SERVICE_URL`)                      │   │
│   │ - Render Headless Browser / Anti-Bot Bypass / PDF Extractor     │   │
│   └────────────────────────────────┬────────────────────────────────┘   │
└────────────────────────────────────┼────────────────────────────────────┘
                                     │ HTTP Webhook Delivery
                                     │ POST `WEBHOOK_URL` + HMAC SHA256
                                     ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                    Sovera Core API (Webhook Receiver)                   │
│                                                                         │
│   ┌─────────────────────────────────────────────────────────────────┐   │
│   │ Webhook Ingestion Controller (`POST /api/v1/webhooks/crawler`)  │   │
│   │ 1. Verify HMAC SHA-256 Signature (`X-Hub-Signature-256`)        │   │
│   │ 2. Compute Content Hash (`content_hash = SHA256(raw_text)`)     │   │
│   │ 3. Check Duplicate -> Enqueue `task:llm_extraction` to Redis    │   │
│   └────────────────────────────────┬────────────────────────────────┘   │
│                                    │ Async Background Processing        │
│                                    ▼                                    │
│   ┌─────────────────────────────────────────────────────────────────┐   │
│   │ Gemini LLM Extraction & Vector Embedding Engine (1536-dim)      │   │
│   │ Save to `public_corporate_signals` Database Table               │   │
│   └─────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Database Schema Specification

### 3.1 Table `crawling_targets` (Registri Target Pemindaian)
Tabel ini menyimpan daftar URL situs target, laporan PDF, atau RSS feed yang perlu dipantau secara periodik.

```sql
CREATE TABLE crawling_targets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_name VARCHAR(255) NOT NULL,            -- Misal: "IDX Sustainability Report Telkom"
    source_type VARCHAR(50) NOT NULL,             -- 'PDF_DOCUMENT', 'NEWS_ARTICLE', 'RAW_WEB', 'SOCIAL_POST'
    target_url TEXT NOT NULL,                     -- URL target dokumen/halaman
    check_interval_hours INT DEFAULT 24,          -- Interval pemeriksaan dalam jam
    last_scraped_at TIMESTAMPTZ,                  -- Waktu pemindaian terakhir
    next_run_at TIMESTAMPTZ DEFAULT NOW(),        -- Jadwal eksekusi berikutnya
    is_active BOOLEAN DEFAULT TRUE,               -- Status keaktifan target
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Index untuk efisiensi query scheduler
CREATE INDEX idx_crawling_targets_active_next_run 
ON crawling_targets (is_active, next_run_at) 
WHERE is_active = TRUE;
```

### 3.2 Table `crawling_logs` (Histori Eksekusi Dispatched Tasks)
Tabel ini mencatat histori pengiriman tugas dari Backend Orchestrator ke Scraper Service beserta hasilnya.

```sql
CREATE TABLE crawling_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    target_id UUID REFERENCES crawling_targets(id) ON DELETE SET NULL,
    task_id VARCHAR(100) NOT NULL UNIQUE,         -- Unique ID tugas yang dikirim ke scraper
    status VARCHAR(50) NOT NULL,                  -- 'DISPATCHED', 'COMPLETED', 'FAILED'
    http_status_code INT,
    error_message TEXT,
    execution_time_ms INT,
    content_hash VARCHAR(64),                     -- SHA-256 hash hasil ekstraksi
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 4. Orchestration & Dispatching Lifecycle

### Step 1: Periodic Cron Scheduler Tick
1. Asynq Scheduler (atau Go Worker) menjalankan job `task:dispatch_crawling_tasks` setiap interval (misal: setiap jam 00:00 WIB atau setiap 1 jam).
2. Worker mengambil daftar target aktif dari database:
   ```sql
   SELECT id, source_name, source_type, target_url 
   FROM crawling_targets 
   WHERE is_active = TRUE AND next_run_at <= NOW()
   LIMIT 50;
   ```

### Step 2: HTTP Task Dispatching to WebScraper Service
Untuk setiap target yang siap dieksekusi, Backend menyusun payload HTTP POST ke `SCRAPER_SERVICE_URL` (`POST /api/v1/scrape-tasks`):

```json
{
  "task_id": "job_csr_idx_991823",
  "client_origin": "sovera_b2b_engine",
  "source_type": "PDF_DOCUMENT",
  "target_url": "https://example-corpo.com/reports/sustainability-report-2025.pdf",
  "callback_url": "http://localhost:4000/api/v1/webhooks/crawler",
  "config": {
    "render_js": false,
    "bypass_anti_bot": true,
    "max_pages": 50
  }
}
```

Backend memperbarui `last_scraped_at` dan `next_run_at` pada tabel `crawling_targets`:
```sql
UPDATE crawling_targets 
SET last_scraped_at = NOW(), 
    next_run_at = NOW() + (check_interval_hours || ' hours')::INTERVAL,
    updated_at = NOW()
WHERE id = :target_id;
```

### Step 3: Scraper Execution & Webhook Delivery
1. Scraper Service menerima request, mengeksekusi ekstraksi PDF / HTML, dan segera mengembalikan `202 Accepted`.
2. Setelah ekstraksi selesai, Scraper mengirimkan *callback webhook* ke `WEBHOOK_URL` (`POST /api/v1/webhooks/crawler`) dengan header `X-Hub-Signature-256`.

### Step 4: Webhook Ingestion & AI Processing
1. Endpoint `/api/v1/webhooks/crawler` di Backend memverifikasi tanda tangan HMAC.
2. Backend memeriksa `content_hash`. Jika belum ada di `public_corporate_signals`, backend memasukkan `task:llm_extraction` ke antrean Asynq Redis.
3. Gemini LLM mengekstraksi entitas JSON dan `text-embedding-004` menghasilkan vektor 1536-dimensi untuk disimpan ke `public_corporate_signals`.

---

## 5. Environment Variables Configuration

| Variable Name | Example Value | Description |
| :--- | :--- | :--- |
| `SCRAPER_SERVICE_URL` | `http://localhost:8000/api/scrape-tasks` | Endpoint WebScraper Service untuk penerimaan tugas scraping |
| `WEBHOOK_URL` | `https://api.sovera.id/api/v1/webhooks/crawler` | URL Webhook Callback Backend penerima payload hasil scraping |
| `WEBHOOK_SECRET_KEY` | `super_secret_crawler_key_123` | Secret key pre-shared untuk enkripsi HMAC SHA-256 signature |

---

## 6. Error Handling & Resilience Strategy

1. **Scraper Unreachable / Timeout:**
   Jika pemanggilan HTTP POST ke `SCRAPER_SERVICE_URL` mengalami timeout atau error 5xx, Backend mencatat error ke `crawling_logs` dengan status `FAILED` dan menjadwalkan ulang *retry* 1 jam kemudian.
2. **Duplicate Content Protection:**
   Jika berkas PDF tidak mengalami perubahan konten (nilai `content_hash` identik dengan histori sebelumnya), Backend menghentikan pemrosesan AI di tingkat Ingestion Controller untuk menghemat token LLM.
3. **Graceful Rate Limiting:**
   Dispatching tugas dilakukan secara bertahap (*batched dispatching* 5–10 request per detik) untuk mencegah penumpukan antrean pada WebScraper Service.
