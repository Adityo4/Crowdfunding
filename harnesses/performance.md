# ⚡ Performance Harness

> **Tujuan:** Memastikan sistem yang dibangun cepat, efisien, dan bisa menangani beban yang diharapkan.
> **Gunakan harness ini saat:** optimasi fitur yang lambat, load testing, arsitektur caching, atau performance review.

---

## 🧭 Konteks untuk AI

```
Kamu adalah performance-aware engineer. Setiap keputusan teknis mempertimbangkan dampaknya ke latensi dan throughput.
Prinsipmu:
1. JANGAN optimasi sebelum ada data — "premature optimization is the root of all evil"
2. Ukur dulu, baru optimasi — profiling sebelum asumsi
3. Optimasi di level yang paling berdampak: algoritma > query > caching > infra
4. Setiap optimasi harus punya benchmark sebelum dan sesudah
5. Performa adalah fitur — degradasi performa = bug

Sebelum optimasi, SELALU tanya: "Di mana bottleneck-nya? Data apa yang membuktikannya?"
```

---

## 1. 📏 Metrik Performa

### Metrik Utama yang Harus Dipantau

| Metrik | Definisi | Target Umum |
|---|---|---|
| **Latency (p50)** | 50% request selesai dalam waktu ini | < 100ms (API), < 3s (web) |
| **Latency (p95)** | 95% request selesai dalam waktu ini | < 500ms (API), < 5s (web) |
| **Latency (p99)** | 99% request selesai dalam waktu ini | < 1s (API), < 10s (web) |
| **Throughput (RPS)** | Request per second yang bisa ditangani | Sesuai kebutuhan bisnis |
| **Error Rate** | % request yang error | < 0.1% |
| **Availability** | % uptime sistem | ≥ 99.9% (3 nines) |
| **TTFB** | Time to First Byte (web) | < 200ms |
| **LCP** | Largest Contentful Paint (web) | < 2.5s |
| **CLS** | Cumulative Layout Shift (web) | < 0.1 |

### Core Web Vitals (Frontend)
```
LCP  (Largest Contentful Paint)  → ukur loading performa → target < 2.5s
FID  (First Input Delay)         → ukur interaktivitas   → target < 100ms
CLS  (Cumulative Layout Shift)   → ukur visual stability → target < 0.1
INP  (Interaction to Next Paint) → ukur responsivitas    → target < 200ms
```

---

## 2. 🔍 Profiling & Menemukan Bottleneck

### Proses Investigasi

```
1. Reproduce — pastikan masalah bisa direproduksi secara konsisten
2. Measure   — ukur baseline sebelum perubahan apapun
3. Profile   — gunakan profiler untuk menemukan hot path
4. Isolate   — identifikasi bottleneck spesifik
5. Optimize  — perbaiki satu hal sekaligus
6. Verify    — ukur ulang, pastikan ada improvement
```

### [KUSTOMISASI] Tools Profiling
```
Backend Profiling:
- Node.js  → clinic.js, 0x, --prof flag, Chrome DevTools
- Python   → cProfile, py-spy, line_profiler
- Go       → pprof (built-in), go tool trace
- Java     → JProfiler, async-profiler, JVM Flight Recorder
- Rust     → perf, flamegraph

Database Profiling:
- PostgreSQL → EXPLAIN ANALYZE, pg_stat_statements, pgBadger
- MySQL      → EXPLAIN, slow query log, Performance Schema
- MongoDB    → explain(), db.currentOp(), Atlas Performance Advisor

Frontend Profiling:
- Chrome DevTools Performance tab
- Lighthouse (CLI atau browser extension)
- WebPageTest
- React DevTools Profiler (khusus React)
```

### Checklist Profiling
- [ ] Profiling dilakukan di environment yang mendekati production (bukan lokal dev)
- [ ] Baseline diukur sebelum optimasi apapun
- [ ] Hasil profiling disimpan untuk perbandingan
- [ ] Bottleneck diidentifikasi berdasarkan data, bukan asumsi

---

## 3. 🗃️ Query & Database Optimization

### Prinsip
- **Query yang lambat** adalah bottleneck paling umum — selalu cek slow query log
- **Index dengan tepat** — terlalu sedikit = lambat, terlalu banyak = write lambat
- **N+1 query** adalah musuh utama performa di ORM

### Checklist Query
- [ ] Semua query yang sering dipakai sudah di-`EXPLAIN ANALYZE`
- [ ] Index ada di kolom yang dipakai untuk `WHERE`, `JOIN`, `ORDER BY`
- [ ] Tidak ada **N+1 query** (gunakan eager loading / join)
- [ ] Query yang berat pakai **pagination** — tidak pernah `SELECT *` tanpa limit
- [ ] Query yang sangat berat dijalankan **async** atau di background job
- [ ] Connection pool dikonfigurasi dengan benar (tidak terlalu kecil atau besar)

### Pola N+1 dan Solusinya

```
❌ N+1 Problem:
   SELECT * FROM orders          → 1 query
   For each order:
     SELECT * FROM users WHERE id = order.user_id  → N query

✅ Solusi — Eager Loading / JOIN:
   SELECT orders.*, users.* FROM orders
   JOIN users ON orders.user_id = users.id  → 1 query
```

### Index Strategy
```
Kapan tambah index:
✅ Kolom yang sering ada di WHERE clause
✅ Kolom yang dipakai untuk JOIN
✅ Kolom yang dipakai untuk ORDER BY atau GROUP BY
✅ Foreign key

Kapan HINDARI index:
❌ Tabel yang sangat kecil (< 1000 rows)
❌ Kolom dengan cardinality rendah (e.g. boolean, enum sedikit)
❌ Tabel yang sangat sering di-INSERT/UPDATE (index mahal di write)
```

---

## 4. 🚀 Caching Strategy

### Hierarki Cache (dari tercepat ke terlambat)

```
L1: In-memory (dalam proses)   → nanoseconds  → Map, LRU cache lokal
L2: Distributed cache (Redis)  → microseconds → Data yang dishare antar instance
L3: CDN / Edge cache           → milliseconds → Static asset, response publik
L4: Database query cache       → milliseconds → Hasil query yang mahal
```

### Strategi Caching

| Pattern | Kapan Dipakai | Risiko |
|---|---|---|
| **Cache-aside** | Data yang dibaca sering, ditulis jarang | Cache miss pertama tetap lambat |
| **Write-through** | Data harus selalu konsisten | Write jadi lebih lambat |
| **Write-behind** | Write throughput tinggi, toleran delay | Data bisa hilang jika cache crash |
| **Read-through** | Sederhana, library yang handle | Kurang kontrol |

### Checklist Caching
- [ ] **Cache key** unik dan predictable — sertakan versi jika perlu
- [ ] **TTL (Time to Live)** dikonfigurasi — tidak ada cache yang hidup selamanya
- [ ] **Cache invalidation** strategy terdefinisi — kapan cache di-clear?
- [ ] **Cache stampede** dicegah (gunakan mutex / probabilistic early expiration)
- [ ] Data sensitif (password, token) TIDAK di-cache
- [ ] Cache miss rate dipantau — terlalu tinggi berarti caching tidak efektif

### [KUSTOMISASI] Tools Caching
```
Distributed Cache:
- Redis     → paling populer, support berbagai data structure
- Memcached → lebih simpel, hanya key-value

CDN:
- Cloudflare → populer, edge network luas
- AWS CloudFront
- Fastly

In-Process:
- Node.js → node-cache, lru-cache
- Python  → cachetools, functools.lru_cache
- Go      → ristretto, groupcache
- Java    → Caffeine, Guava Cache
```

---

## 5. 🌐 Frontend Performance

### Checklist Asset Optimization
- [ ] **Image** di-compress dan pakai format modern (WebP, AVIF)
- [ ] **Lazy loading** untuk gambar dan komponen di below-the-fold
- [ ] **Code splitting** — tidak load semua JS sekaligus
- [ ] **Tree shaking** aktif — hapus kode yang tidak dipakai
- [ ] **Minification** untuk JS, CSS, HTML di production build
- [ ] **Gzip / Brotli compression** aktif di web server

### Checklist Loading Strategy
- [ ] Critical CSS di-inline, sisanya di-load async
- [ ] Font di-preload, gunakan `font-display: swap`
- [ ] Resource hints: `<link rel="preconnect">` untuk domain third-party
- [ ] Static asset pakai **CDN** dengan cache header yang tepat
- [ ] **Service Worker** untuk offline support dan caching (jika relevan)

### Checklist Bundle
- [ ] Bundle size dipantau — alert jika naik signifikan
- [ ] Third-party library dipilih dengan pertimbangan ukuran
- [ ] `import { spesifik }` bukan `import seluruhLibrary`
- [ ] Dynamic import untuk fitur yang tidak selalu dipakai

---

## 6. 🏋️ Load Testing

### Jenis Load Test

| Tipe | Tujuan | Kapan |
|---|---|---|
| **Load Test** | Verifikasi performa di beban normal | Sebelum launch, setiap major release |
| **Stress Test** | Temukan batas sistem | Saat sizing infrastruktur |
| **Spike Test** | Simulasi lonjakan tiba-tiba | Sebelum event / promo besar |
| **Soak Test** | Deteksi memory leak & degradasi panjang | Setelah load test normal |

### [KUSTOMISASI] Tools Load Test
```
- k6          → JavaScript-based, developer-friendly, recommended
- Apache JMeter → GUI-based, populer di enterprise
- Locust       → Python-based, mudah untuk skenario kompleks
- Gatling      → Scala-based, report bagus
- Artillery    → YAML-based, mudah untuk API testing
```

### Checklist Load Test
- [ ] Skenario test merepresentasikan traffic production yang sesungguhnya
- [ ] Baseline performa terdokumentasi sebelum load test
- [ ] Target RPS dan latency terdefinisi sebelum test
- [ ] Load test dijalankan di environment yang identik dengan production
- [ ] Resource utilization (CPU, RAM, DB connection) dipantau selama test
- [ ] Hasil load test terdokumentasi dan dibandingkan dengan baseline

---

## 7. 📊 Performance Monitoring (Production)

### Checklist Monitoring
- [ ] **APM (Application Performance Monitoring)** aktif di production
- [ ] Alert untuk p99 latency di atas threshold
- [ ] Alert untuk error rate di atas threshold
- [ ] Dashboard performa tersedia dan mudah dibaca tim
- [ ] **Distributed tracing** aktif untuk debug lintas service
- [ ] **Database slow query log** aktif dan dipantau

### [KUSTOMISASI] Tools Monitoring
```
APM:
- Datadog APM
- New Relic
- Elastic APM
- Sentry Performance

Metrics + Dashboards:
- Prometheus + Grafana (self-hosted)
- Datadog
- CloudWatch (AWS)

Distributed Tracing:
- OpenTelemetry (vendor-neutral, recommended)
- Jaeger (self-hosted)
- Zipkin
```

---

## 8. 📋 Performance Checklist (Pre-Launch)

### Backend
- [ ] Slow query log diaktifkan dan tidak ada query > 1s
- [ ] N+1 query tidak ditemukan di critical path
- [ ] Caching strategy terdefinisi untuk data yang sering diakses
- [ ] Load test dijalankan dan hasilnya memuaskan
- [ ] Connection pool dikonfigurasi dengan benar

### Frontend
- [ ] Lighthouse score: Performance ≥ 85
- [ ] LCP < 2.5s, CLS < 0.1
- [ ] Bundle size dalam batas yang wajar
- [ ] Gambar dioptimasi dan lazy-loaded
- [ ] CDN aktif untuk static asset

### Infrastructure
- [ ] Auto-scaling dikonfigurasi
- [ ] Resource limit (CPU, RAM) dikonfigurasi di container
- [ ] Monitoring dan alerting aktif

---

## 📚 Referensi

- [Web Vitals](https://web.dev/vitals/)
- [High Performance Browser Networking — Ilya Grigorik](https://hpbn.co/)
- [Use The Index, Luke (Database)](https://use-the-index-luke.com/)
- [k6 Documentation](https://k6.io/docs/)
- [Systems Performance — Brendan Gregg](https://www.brendangregg.com/systems-performance.html)

---

*Performance Harness v1.0 — bagian dari [Harness Engineering](../README.md)*
