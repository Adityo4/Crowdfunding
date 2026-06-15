# 🗄️ Database Harness

> **Tujuan:** Memastikan desain database yang matang — schema yang solid, indexing yang tepat, caching yang efisien, dan query yang optimal sejak awal.
> **Gunakan harness ini saat:** mendesain schema baru, melakukan database review, optimasi query lambat, atau merancang strategi caching.

---

## 🧭 Konteks untuk AI

```
Kamu adalah database engineer yang berpengalaman.
Database yang buruk tidak bisa diperbaiki hanya dengan kode yang bagus.

Prinsipmu:
1. Desain schema berdasarkan ACCESS PATTERN — bagaimana data dibaca, bukan hanya disimpan
2. Indexing adalah seni — terlalu sedikit = lambat baca, terlalu banyak = lambat tulis
3. Normalisasi adalah titik awal — denormalisasi hanya saat ada alasan performa yang terbukti
4. Migration harus reversible — selalu siapkan rollback
5. Cache adalah layer, bukan solusi — pahami data di baliknya sebelum men-cache

Sebelum menulis DDL apapun, SELALU tanya:
"Bagaimana data ini akan dibaca? Berapa sering? Oleh siapa?"
```

---

## 1. 📐 Schema Design

### Prinsip Utama

```
Pilih normalisasi level yang tepat:

3NF (Third Normal Form) — titik awal yang baik:
  ✅ Tidak ada data duplikat
  ✅ Setiap kolom bergantung penuh pada primary key
  ✅ Cocok untuk data transaksional (OLTP)

Denormalisasi — hanya saat terbukti perlu:
  ✅ Performa read sangat kritis dan tidak bisa dioptimasi dengan index
  ✅ Data yang di-denormalisasi jarang berubah
  ⚠️ Harus ada mekanisme sinkronisasi data yang konsisten
```

### Pemilihan Primary Key

```
UUID vs Auto-increment Integer — pilih dengan sadar:

UUID (Universally Unique Identifier):
  ✅ Aman untuk expose ke public (tidak bisa ditebak jumlah record)
  ✅ Bisa di-generate di aplikasi tanpa round-trip ke DB
  ✅ Aman untuk distributed system / multi-database
  ⚠️ Index lebih besar, insert lebih lambat (random ordering)
  → Gunakan UUID v7 atau ULID untuk performa lebih baik (time-sortable)

Auto-increment Integer (BIGINT):
  ✅ Index lebih kecil dan efisien
  ✅ Insert lebih cepat (sequential)
  ✅ Lebih mudah untuk debugging
  ⚠️ Bocorkan jumlah record ke public
  ⚠️ Masalah di distributed system

Rekomendasi:
  Public-facing ID  → UUID v7 atau ULID
  Internal join key → BIGINT auto-increment
  Expose ke API     → Selalu UUID, bukan integer
```

### Tipe Data yang Tepat

```sql
-- Teks
VARCHAR(n)    -- Gunakan saat ada batas maksimum yang jelas
TEXT          -- Gunakan untuk konten panjang tanpa batas
CHAR(n)       -- Hanya untuk fixed-length string (kode negara: CHAR(2))

-- Angka
INT / INTEGER    -- Bilangan bulat biasa (max ~2.1 miliar)
BIGINT           -- ID, counter besar (max ~9.2 kuadriliun)
DECIMAL(p,s)     -- Uang, harga — JANGAN pakai FLOAT (rounding error!)
NUMERIC(p,s)     -- Sama dengan DECIMAL, alias
-- ❌ Jangan pakai FLOAT/DOUBLE untuk nilai finansial

-- Waktu
TIMESTAMP WITH TIME ZONE  -- Default untuk semua timestamp (simpan dalam UTC)
DATE                      -- Jika hanya perlu tanggal (ulang tahun, tanggal lahir)
-- ❌ Jangan simpan timestamp sebagai VARCHAR atau Unix epoch integer

-- Boolean
BOOLEAN       -- Gunakan true/false native DB
-- ❌ Jangan pakai TINYINT(1) atau VARCHAR 'Y'/'N'

-- JSON
JSONB         -- PostgreSQL: untuk data semi-structured, bisa di-index
JSON          -- Hanya jika butuh preserve whitespace/order (jarang)
-- ⚠️ Jangan semua data dijadikan JSON — relasional tetap lebih efisien

-- Enum
-- Pilih antara:
-- 1. DB-level ENUM (rigid, sulit diubah)
-- 2. SMALLINT + aplikasi-level mapping (lebih fleksibel)
-- 3. Lookup table (paling fleksibel, bisa di-JOIN)
```

### Kolom Audit Wajib

```sql
-- Tambahkan di setiap tabel yang perlu audit trail:
created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
created_by  UUID REFERENCES users(id),      -- opsional
updated_by  UUID REFERENCES users(id),      -- opsional

-- Soft delete (jika diperlukan):
deleted_at  TIMESTAMP WITH TIME ZONE,       -- NULL = aktif, NOT NULL = dihapus
deleted_by  UUID REFERENCES users(id)
```

### Naming Convention

```sql
-- Tabel: snake_case, noun jamak
users, orders, order_items, payment_methods

-- Kolom: snake_case, deskriptif
user_id, email_verified_at, total_amount_cents

-- Primary Key: selalu bernama "id"
id UUID PRIMARY KEY DEFAULT gen_random_uuid()

-- Foreign Key: <table_singular>_id
user_id, order_id, product_id

-- Boolean: awali dengan is_, has_, can_
is_active, has_verified_email, can_withdraw

-- Timestamp: akhiri dengan _at
created_at, updated_at, deleted_at, paid_at

-- Index: idx_<table>_<columns>
idx_orders_user_id
idx_orders_status_created_at
idx_users_email
```

---

## 2. 🔍 Indexing Strategy

### Prinsip Indexing

```
Golden Rules:
1. Index kolom yang ada di WHERE, JOIN ON, ORDER BY, GROUP BY
2. Urutan kolom di composite index SANGAT penting (leftmost prefix rule)
3. Index yang tidak pernah dipakai lebih berbahaya dari tidak ada index
4. Setiap index = overhead di INSERT, UPDATE, DELETE
5. Selalu EXPLAIN ANALYZE sebelum dan sesudah menambah index

Kapan TIDAK perlu index:
  ❌ Tabel kecil (< 1.000 rows) — full scan lebih cepat
  ❌ Kolom dengan cardinality sangat rendah (boolean, gender, status sedikit)
  ❌ Kolom yang sangat jarang di-query
  ❌ Tabel yang sangat sering di-INSERT (overhead terlalu besar)
```

### Jenis Index dan Kapan Dipakai

```sql
-- 1. SINGLE COLUMN INDEX — paling umum
CREATE INDEX idx_users_email ON users(email);
-- Kapan: kolom sering di-WHERE secara independen

-- 2. COMPOSITE INDEX — untuk query multi-kolom
CREATE INDEX idx_orders_user_status ON orders(user_id, status);
-- Kapan: query sering filter user_id DAN status bersamaan
-- PENTING: Urutan kolom = urutan selectivity (yang paling selektif dulu)
-- Query yang memanfaatkan: WHERE user_id = ? AND status = ?
-- Query yang TIDAK memanfaatkan: WHERE status = ? (tanpa user_id)

-- 3. UNIQUE INDEX — constraint + performa
CREATE UNIQUE INDEX idx_users_email_unique ON users(email);
-- Kapan: kolom harus unik (email, username, SKU produk)

-- 4. PARTIAL INDEX — index subset data (PostgreSQL)
CREATE INDEX idx_orders_pending ON orders(created_at)
WHERE status = 'pending';
-- Kapan: query sering filter ke subset tertentu (active records, unread, dll)
-- Lebih kecil = lebih cepat dari full index

-- 5. COVERING INDEX — include kolom tambahan
CREATE INDEX idx_orders_user ON orders(user_id)
INCLUDE (status, total_amount, created_at);  -- PostgreSQL
-- Kapan: query butuh kolom tertentu tapi tidak perlu full table scan
-- Index-only scan: tidak perlu baca tabel sama sekali

-- 6. EXPRESSION INDEX — index hasil ekspresi
CREATE INDEX idx_users_email_lower ON users(LOWER(email));
-- Kapan: query pakai fungsi: WHERE LOWER(email) = ?

-- 7. GIN INDEX — untuk array, JSONB, full-text search (PostgreSQL)
CREATE INDEX idx_products_tags ON products USING GIN(tags);
CREATE INDEX idx_posts_body_fts ON posts USING GIN(to_tsvector('english', body));
-- Kapan: query ke dalam array atau JSONB field
```

### Composite Index — Urutan Kolom

```sql
-- Skenario: tabel orders dengan query berikut:
-- 1. WHERE user_id = ? AND status = ? ORDER BY created_at DESC
-- 2. WHERE user_id = ? ORDER BY created_at DESC
-- 3. WHERE status = ?

-- ✅ Index yang efisien untuk query 1 dan 2:
CREATE INDEX idx_orders_user_status_date ON orders(user_id, status, created_at DESC);

-- Leftmost prefix rule — index di atas bisa dipakai untuk:
-- ✅ WHERE user_id = ?
-- ✅ WHERE user_id = ? AND status = ?
-- ✅ WHERE user_id = ? AND status = ? ORDER BY created_at
-- ❌ WHERE status = ? (tidak mulai dari kolom pertama)
-- ❌ WHERE created_at > ? (skip user_id dan status)

-- Untuk query 3 (WHERE status = ?), butuh index terpisah:
CREATE INDEX idx_orders_status ON orders(status);
```

### Menganalisis Index dengan EXPLAIN

```sql
-- Selalu EXPLAIN ANALYZE — bukan hanya EXPLAIN
EXPLAIN ANALYZE
SELECT o.id, o.total_amount, u.email
FROM orders o
JOIN users u ON o.user_id = u.id
WHERE o.status = 'pending'
  AND o.created_at > NOW() - INTERVAL '7 days'
ORDER BY o.created_at DESC
LIMIT 20;

-- Yang perlu diperhatikan di output:
-- ✅ "Index Scan"     → menggunakan index
-- ✅ "Index Only Scan"→ covering index (paling efisien)
-- ⚠️ "Bitmap Heap Scan" → index dipakai tapi butuh heap access
-- ❌ "Seq Scan"       → full table scan, tidak ada index yang cocok
-- ❌ "Hash Join" pada tabel besar → pertimbangkan index pada join key

-- Perhatikan juga:
-- actual rows vs estimated rows: jika jauh berbeda → statistik perlu di-update
-- Execution Time: total waktu eksekusi aktual
```

### Index Maintenance

```sql
-- Cek index yang tidak pernah dipakai (PostgreSQL):
SELECT schemaname, tablename, indexname, idx_scan
FROM pg_stat_user_indexes
WHERE idx_scan = 0
ORDER BY schemaname, tablename;

-- Cek index yang bloat (perlu di-REINDEX):
SELECT tablename, indexname,
       pg_size_pretty(pg_relation_size(indexrelid)) AS index_size
FROM pg_stat_user_indexes
ORDER BY pg_relation_size(indexrelid) DESC;

-- Update statistik agar query planner akurat:
ANALYZE users;
ANALYZE orders;
-- Atau semua sekaligus:
ANALYZE;
```

---

## 3. 🔄 Query Optimization

### Anti-Pattern Query yang Paling Umum

```sql
-- ❌ SELECT * — ambil semua kolom padahal tidak perlu
SELECT * FROM orders WHERE user_id = $1;
-- ✅ Pilih kolom yang dibutuhkan saja
SELECT id, status, total_amount, created_at FROM orders WHERE user_id = $1;

-- ❌ Fungsi di WHERE — membuat index tidak bisa dipakai
WHERE YEAR(created_at) = 2026
-- ✅ Gunakan range comparison
WHERE created_at >= '2026-01-01' AND created_at < '2027-01-01'

-- ❌ LIKE dengan leading wildcard — tidak bisa pakai B-tree index
WHERE name LIKE '%keyword%'
-- ✅ Full-text search atau trigram index (pg_trgm)
WHERE to_tsvector('indonesian', name) @@ plainto_tsquery('keyword')

-- ❌ N+1 Query — loop query di aplikasi
for user in users:
    orders = db.query("SELECT * FROM orders WHERE user_id = ?", user.id)
-- ✅ JOIN atau subquery
SELECT u.*, o.* FROM users u
LEFT JOIN orders o ON u.id = o.user_id
WHERE u.id IN (...)

-- ❌ COUNT(*) untuk cek keberadaan record
SELECT COUNT(*) FROM orders WHERE user_id = $1 -- harus scan semua
-- ✅ EXISTS lebih cepat — berhenti di record pertama
SELECT EXISTS(SELECT 1 FROM orders WHERE user_id = $1)

-- ❌ OFFSET pagination pada data besar
SELECT * FROM orders LIMIT 20 OFFSET 10000 -- skip 10000 rows tetap dibaca!
-- ✅ Cursor pagination menggunakan keyed column
SELECT * FROM orders WHERE id > $lastId ORDER BY id LIMIT 20
```

### Query Pattern yang Efisien

```sql
-- Batch insert daripada individual insert:
INSERT INTO order_items (order_id, product_id, quantity, unit_price)
VALUES
  ($1, $2, $3, $4),
  ($5, $6, $7, $8),
  ($9, $10, $11, $12);  -- insert banyak row sekaligus

-- Upsert (INSERT OR UPDATE):
INSERT INTO user_preferences (user_id, key, value)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, key)
DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();

-- Window function daripada subquery berulang:
SELECT
  id, user_id, total_amount,
  ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY created_at DESC) AS rn,
  SUM(total_amount) OVER (PARTITION BY user_id) AS user_total
FROM orders;

-- CTE untuk query yang kompleks dan mudah dibaca:
WITH recent_active_users AS (
  SELECT DISTINCT user_id
  FROM orders
  WHERE created_at > NOW() - INTERVAL '30 days'
),
user_stats AS (
  SELECT user_id, COUNT(*) as order_count, SUM(total_amount) as lifetime_value
  FROM orders
  GROUP BY user_id
)
SELECT u.email, s.order_count, s.lifetime_value
FROM recent_active_users r
JOIN users u ON r.user_id = u.id
JOIN user_stats s ON r.user_id = s.user_id;
```

---

## 4. ⚡ Caching Strategy

### Cache Layer untuk Database

```
Hierarki cache (dari tercepat):

L1: Application-level cache (in-process)
    → Untuk data yang sama diakses berkali-kali dalam satu request
    → Map / dictionary di dalam proses aplikasi
    → Hidup hanya selama request

L2: Distributed cache (Redis / Memcached)
    → Untuk data yang dishare antar instance aplikasi
    → Hidup antar request, bisa dikonfigurasi TTL
    → Sweet spot untuk database caching

L3: Query cache / Read replica
    → Database-level caching
    → Read replica untuk offload read query dari primary
```

### Apa yang Layak Di-cache

```
✅ Bagus untuk di-cache:
  - Data yang sering dibaca, jarang berubah (produk, kategori, config)
  - Hasil agregasi yang mahal (laporan, dashboard stats)
  - Session data & token
  - Rate limiting counter
  - Data referensi (list negara, list provinsi)

❌ Jangan di-cache:
  - Data real-time (stok yang berubah tiap detik)
  - Data personal/sensitif (kecuali dengan enkripsi dan TTL ketat)
  - Data yang sering berubah dan harus konsisten
  - Query yang jarang dieksekusi
```

### Redis Patterns

```
[KUSTOMISASI] Sesuaikan key naming dengan project:

Key naming convention:
  <app>:<entity>:<identifier>:<field?>
  myapp:user:usr_123
  myapp:user:usr_123:profile
  myapp:order:ord_456:items
  myapp:rate_limit:ip:192.168.1.1
```

```python
# Pattern 1: Cache-Aside (paling umum)
async def get_user(user_id: str) -> User:
    cache_key = f"user:{user_id}"

    # 1. Cek cache
    cached = await redis.get(cache_key)
    if cached:
        return User.parse_raw(cached)  # Cache hit

    # 2. Cache miss — ambil dari DB
    user = await db.query("SELECT * FROM users WHERE id = $1", user_id)
    if not user:
        raise NotFoundError("User", user_id)

    # 3. Simpan ke cache
    await redis.setex(cache_key, 3600, user.json())  # TTL 1 jam
    return user

async def update_user(user_id: str, data: dict) -> User:
    user = await db.update("UPDATE users SET ... WHERE id = $1", user_id, data)

    # Invalidate cache setelah update
    await redis.delete(f"user:{user_id}")
    return user
```

```python
# Pattern 2: Cache dengan Lock (mencegah Cache Stampede)
async def get_expensive_report(report_id: str):
    cache_key = f"report:{report_id}"
    lock_key = f"lock:report:{report_id}"

    # Cek cache dulu
    cached = await redis.get(cache_key)
    if cached:
        return json.loads(cached)

    # Ambil lock untuk mencegah thundering herd
    async with redis.lock(lock_key, timeout=30):
        # Double-check setelah dapat lock
        cached = await redis.get(cache_key)
        if cached:
            return json.loads(cached)

        # Generate report (operasi mahal)
        report = await generate_report(report_id)
        await redis.setex(cache_key, 1800, json.dumps(report))
        return report
```

```python
# Pattern 3: Write-Through Cache
async def create_product(data: dict) -> Product:
    product = await db.insert("INSERT INTO products ... RETURNING *", data)

    # Langsung simpan ke cache bersamaan dengan write ke DB
    cache_key = f"product:{product.id}"
    await redis.setex(cache_key, 3600, product.json())

    return product
```

### Cache Invalidation Strategy

```
Strategi invalidasi cache:

1. TTL-based (paling mudah):
   → Set waktu kadaluarsa yang masuk akal
   → Risiko: data stale selama TTL
   → Cocok untuk data yang bisa sedikit stale (config, katalog)

2. Event-based (paling akurat):
   → Invalidate saat data berubah
   → Risiko: bisa miss jika ada multiple update path
   → Cocok untuk data yang harus selalu fresh

3. Versioned cache (paling robust):
   → Tambahkan versi ke cache key
   → key: "product:prod_123:v5"
   → Update versi di DB, cache lama otomatis tidak relevan
   → Cocok untuk data dengan many writers

4. Cache tags / grouping:
   → Kelompokkan cache yang related
   → Invalidate semua cache dalam group sekaligus
   → Contoh: semua cache yang berhubungan dengan user_123
```

### Redis Data Structures yang Tepat

```
String     → Simple key-value, session, token, counter sederhana
Hash       → Object dengan banyak field (user profile)
List       → Queue, recent activity feed (FIFO/LIFO)
Set        → Unique collection, tags, following/follower
Sorted Set → Leaderboard, priority queue, rate limiting
HyperLogLog→ Count unique visitors (approximate, memory-efficient)
Stream     → Event sourcing, message queue
Bitmap     → Feature flags per user, daily active user tracking

Contoh penggunaan:

# Hash untuk user session (lebih efisien dari JSON string):
HSET session:abc123 user_id usr_456 email budi@example.com role admin
EXPIRE session:abc123 86400
HGET session:abc123 user_id

# Sorted Set untuk leaderboard:
ZADD leaderboard 9500 "user:usr_001"
ZADD leaderboard 8200 "user:usr_002"
ZREVRANGE leaderboard 0 9 WITHSCORES  # top 10

# Counter dengan atomic increment:
INCR rate_limit:ip:192.168.1.1
EXPIRE rate_limit:ip:192.168.1.1 60
```

---

## 5. 🔄 Migration Strategy

### Prinsip Migration

```
SELALU buat migration yang bisa di-rollback.
Setiap migration punya UP (apply) dan DOWN (rollback).

Aturan:
1. Satu migration = satu perubahan logis
2. Jangan edit migration yang sudah di-commit ke main branch
3. Buat migration baru jika perlu mengubah migration sebelumnya
4. Test DOWN migration sebelum push ke production
5. Migration yang mengubah data (DML) harus idempotent
```

### Pola Migration yang Aman untuk Production

```sql
-- ❌ BERBAHAYA — lock tabel saat tambah kolom (PostgreSQL)
ALTER TABLE orders ADD COLUMN notes TEXT NOT NULL;
-- Tabel di-lock selama operasi!

-- ✅ AMAN — tambah kolom nullable dulu, isi data, baru NOT NULL
-- Step 1: Tambah kolom nullable (non-blocking di PostgreSQL)
ALTER TABLE orders ADD COLUMN notes TEXT;

-- Step 2: Isi data untuk existing rows (jalankan secara batch jika data besar)
UPDATE orders SET notes = '' WHERE notes IS NULL;

-- Step 3: Baru set NOT NULL setelah semua terisi (atau gunakan DEFAULT)
ALTER TABLE orders ALTER COLUMN notes SET NOT NULL;
ALTER TABLE orders ALTER COLUMN notes SET DEFAULT '';
```

```sql
-- Rename kolom dengan zero-downtime:
-- Jangan langsung rename — aplikasi lama masih pakai nama lama

-- Step 1: Tambah kolom baru
ALTER TABLE users ADD COLUMN full_name TEXT;

-- Step 2: Copy data
UPDATE users SET full_name = name;

-- Step 3: Update aplikasi untuk pakai full_name
-- Deploy aplikasi versi baru

-- Step 4: Hapus kolom lama (setelah aplikasi versi baru stable)
ALTER TABLE users DROP COLUMN name;
```

```sql
-- Index baru tanpa lock tabel (PostgreSQL):
-- ❌ Ini akan lock tabel
CREATE INDEX idx_orders_status ON orders(status);

-- ✅ CONCURRENTLY — tidak lock tabel (lebih lambat tapi aman untuk production)
CREATE INDEX CONCURRENTLY idx_orders_status ON orders(status);
-- Note: tidak bisa dijalankan dalam transaction block
```

### Template Migration

```sql
-- [KUSTOMISASI] Format migration file:
-- Nama file: <timestamp>_<deskripsi_singkat>.sql
-- Contoh: 20260615_add_soft_delete_to_orders.sql

-- ============ UP ============
ALTER TABLE orders ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE orders ADD COLUMN deleted_by UUID REFERENCES users(id);

CREATE INDEX CONCURRENTLY idx_orders_not_deleted
ON orders(created_at)
WHERE deleted_at IS NULL;

-- ============ DOWN ============
DROP INDEX CONCURRENTLY idx_orders_not_deleted;
ALTER TABLE orders DROP COLUMN deleted_by;
ALTER TABLE orders DROP COLUMN deleted_at;
```

---

## 6. 🔗 Connection Pooling & Resource Management

### Konfigurasi Connection Pool

```
Connection pool yang salah dikonfigurasi adalah penyebab umum masalah performa.

Rumus dasar max connections:
  max_connections = (core_count * 2) + effective_spindle_count
  Untuk SSD: max_connections = (core_count * 2) + 1

Contoh untuk server 4 core + SSD:
  max_connections = (4 * 2) + 1 = 9 per app instance

Jika ada 10 instance aplikasi:
  Total = 9 * 10 = 90 connections ke database

Rekomendasi PostgreSQL max_connections: tidak lebih dari 200-300
Gunakan PgBouncer untuk connection pooling di depan PostgreSQL jika perlu lebih
```

```
[KUSTOMISASI] Konfigurasi per stack:

Node.js (pg / postgres.js):
  max: 10          -- max connections per pool
  idleTimeoutMillis: 30000
  connectionTimeoutMillis: 2000

Python (SQLAlchemy):
  pool_size: 10
  max_overflow: 5
  pool_timeout: 30
  pool_recycle: 1800  -- recycle connection setiap 30 menit

Go (database/sql):
  SetMaxOpenConns(10)
  SetMaxIdleConns(5)
  SetConnMaxLifetime(30 * time.Minute)
```

### Checklist Connection Management
- [ ] Connection pool dikonfigurasi (bukan default unlimited)
- [ ] Connection timeout dikonfigurasi
- [ ] Idle connection timeout dikonfigurasi
- [ ] Koneksi selalu dikembalikan ke pool setelah selesai (tidak leak)
- [ ] Health check ke database tersedia
- [ ] PgBouncer / ProxySQL dipertimbangkan untuk high-traffic

---

## 7. 🔒 Keamanan Database

### Prinsip Keamanan Database

```
Least Privilege untuk database user:
  App user   → SELECT, INSERT, UPDATE, DELETE saja
  Admin user → DDL (CREATE, DROP, ALTER) — hanya untuk migration
  Readonly   → SELECT saja — untuk reporting / analytics
  Backup     → pg_dump privilege saja

Jangan pernah:
  ❌ Koneksi ke DB menggunakan user postgres / root
  ❌ Expose port database ke public internet
  ❌ Simpan credential DB di kode atau env file yang ter-commit
  ❌ Pakai satu user DB untuk semua keperluan
```

```sql
-- Contoh setup user dengan least privilege (PostgreSQL):
-- User untuk aplikasi
CREATE USER app_user WITH PASSWORD 'strong_random_password';
GRANT CONNECT ON DATABASE myapp TO app_user;
GRANT USAGE ON SCHEMA public TO app_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO app_user;

-- User untuk readonly (reporting, analytics)
CREATE USER readonly_user WITH PASSWORD 'another_strong_password';
GRANT CONNECT ON DATABASE myapp TO readonly_user;
GRANT USAGE ON SCHEMA public TO readonly_user;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO readonly_user;

-- Pastikan tabel baru otomatis dapat privilege yang sama:
ALTER DEFAULT PRIVILEGES IN SCHEMA public
  GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO app_user;
```

### Checklist Keamanan Database
- [ ] Tidak ada user aplikasi yang punya privilege DDL (CREATE, DROP, ALTER)
- [ ] Database tidak accessible dari public internet (VPC / private subnet)
- [ ] SSL/TLS aktif untuk koneksi ke database
- [ ] Credential di-rotate secara berkala
- [ ] Audit logging aktif untuk query sensitif
- [ ] Backup dienkripsi dan disimpan di lokasi terpisah
- [ ] Row-Level Security (RLS) dipertimbangkan untuk multi-tenant

---

## 8. 📊 Monitoring & Health

### Query yang Harus Dipantau

```sql
-- Top 10 query paling lambat (PostgreSQL — butuh pg_stat_statements):
SELECT query, calls, total_exec_time, mean_exec_time, rows
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;

-- Query yang paling sering dieksekusi:
SELECT query, calls, total_exec_time
FROM pg_stat_statements
ORDER BY calls DESC
LIMIT 10;

-- Table bloat (tabel yang perlu VACUUM):
SELECT tablename, n_dead_tup, n_live_tup,
       round(n_dead_tup::numeric / nullif(n_live_tup, 0) * 100, 2) AS dead_pct
FROM pg_stat_user_tables
WHERE n_dead_tup > 1000
ORDER BY dead_pct DESC;

-- Koneksi aktif:
SELECT state, count(*)
FROM pg_stat_activity
GROUP BY state;

-- Lock yang sedang terjadi:
SELECT pid, query, state, wait_event_type, wait_event
FROM pg_stat_activity
WHERE wait_event_type = 'Lock';
```

### Checklist Monitoring Database
- [ ] Slow query log aktif (query > 1000ms)
- [ ] `pg_stat_statements` aktif untuk analisis query
- [ ] Alert untuk koneksi mendekati limit
- [ ] Alert untuk replikasi lag (jika ada replica)
- [ ] Alert untuk disk usage > 80%
- [ ] Auto-VACUUM dikonfigurasi dengan tepat
- [ ] Backup berhasil dipantau dan di-alert jika gagal

---

## 9. 📋 Database Checklist (Pre-Launch)

### Schema & Design
- [ ] Semua tabel punya primary key
- [ ] Semua foreign key punya index
- [ ] Tipe data sesuai (DECIMAL untuk uang, TIMESTAMP WITH TIME ZONE untuk waktu)
- [ ] Kolom audit (`created_at`, `updated_at`) ada di tabel yang perlu
- [ ] Naming convention konsisten di seluruh schema
- [ ] Tidak ada kolom yang menyimpan multiple values (1NF violated)

### Indexing
- [ ] Semua foreign key column sudah ada index-nya
- [ ] Query paling sering dijalankan sudah dianalisis dengan EXPLAIN ANALYZE
- [ ] Tidak ada index yang tidak pernah dipakai
- [ ] Composite index urutan kolomnya sudah optimal

### Query
- [ ] Tidak ada N+1 query di critical path
- [ ] Semua query punya LIMIT / pagination
- [ ] Tidak ada `SELECT *` di production query
- [ ] Tidak ada fungsi di dalam WHERE clause yang mencegah index dipakai

### Caching
- [ ] Data yang sering dibaca dan jarang berubah sudah di-cache
- [ ] TTL dikonfigurasi untuk semua cache entry
- [ ] Cache invalidation strategy terdefinisi
- [ ] Tidak ada data sensitif di cache tanpa enkripsi

### Migration & Maintenance
- [ ] Semua migration punya rollback (DOWN)
- [ ] Migration yang menambah index menggunakan `CONCURRENTLY`
- [ ] Connection pool dikonfigurasi dengan tepat
- [ ] Backup otomatis berjalan dan pernah dicoba di-restore

### Keamanan
- [ ] App user tidak punya privilege DDL
- [ ] Database tidak accessible dari public internet
- [ ] SSL aktif untuk koneksi
- [ ] Credential database tidak ada di kode atau git history

---

## 📚 Referensi

- [Use The Index, Luke](https://use-the-index-luke.com/) — panduan indexing terlengkap
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Redis Documentation](https://redis.io/docs/)
- [Database Design for Mere Mortals](https://www.oreilly.com/library/view/database-design-for/9780133122282/)
- [High Performance MySQL](https://www.oreilly.com/library/view/high-performance-mysql/9781492080503/)
- [Evolutionary Database Design — Martin Fowler](https://martinfowler.com/articles/evodb.html)

---

*Database Harness v1.0 — bagian dari [Harness Engineering](../README.md)*
