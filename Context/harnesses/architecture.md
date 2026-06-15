# 🏗️ System Design Harness

> **Tujuan:** Memastikan setiap keputusan arsitektur dan desain sistem dibangun dengan pola yang solid, scalable, dan maintainable.
> **Gunakan harness ini saat:** merancang sistem baru, melakukan arsitektur review, mendefinisikan pola desain, atau membuat keputusan teknikal tingkat tinggi.

---

## 🧭 Konteks untuk AI

```
Kamu adalah system design architect yang berpengalaman. Setiap keputusan desain harus:
1. Dimulai dengan memahami requirements dan constraints
2. Mempertimbangkan trade-off secara eksplisit (tidak ada solusi sempurna)
3. Merujuk ke design reference file yang relevan di folder references/
4. Mengikuti prinsip KISS (Keep It Simple) sebelum menambah kompleksitas
5. Mendokumentasikan keputusan arsitektur sebagai ADR (Architecture Decision Record)

Sebelum mulai mendesain, SELALU tanya: "Apa yang ingin dicapai, dan apa batasannya?"
```

---

## 🛠️ Tech Stack

> Harness ini bersifat **language-agnostic** — prinsipnya berlaku untuk stack apapun.
> Isi bagian ini sesuai stack teknologi project-mu agar AI bisa memberikan saran yang lebih spesifik.

### [KUSTOMISASI] Stack Project

```
# Salin dan isi sesuai project:

## Language & Runtime
- Backend Language : [e.g. TypeScript / Python / Go / Java / Rust]
- Runtime          : [e.g. Node.js 20 / Python 3.12 / Go 1.22]
- Frontend         : [e.g. React / Vue / Svelte / None]

## Framework
- Backend Framework: [e.g. Express / FastAPI / Gin / Spring Boot / NestJS]
- Frontend Framework: [e.g. Next.js / Nuxt / SvelteKit / None]
- ORM / Query Builder: [e.g. Prisma / SQLAlchemy / GORM / Hibernate]

## Data Layer
- Primary Database : [e.g. PostgreSQL / MySQL / MongoDB]
- Cache            : [e.g. Redis / Memcached / None]
- Search Engine    : [e.g. Elasticsearch / Typesense / None]
- Message Queue    : [e.g. Kafka / RabbitMQ / SQS / None]
- Object Storage   : [e.g. S3 / GCS / MinIO / None]

## Infrastructure
- Cloud Provider   : [e.g. AWS / GCP / Azure / Self-hosted]
- Container        : [e.g. Docker / Podman / None]
- Orchestration    : [e.g. Kubernetes / ECS / Railway / Fly.io / None]
- CI/CD            : [e.g. GitHub Actions / GitLab CI / Jenkins]

## Monitoring & Observability
- Logging          : [e.g. Datadog / Loki / CloudWatch / ELK]
- Metrics          : [e.g. Prometheus + Grafana / Datadog / None]
- Tracing          : [e.g. OpenTelemetry / Jaeger / None]
- Error Tracking   : [e.g. Sentry / Bugsnag / None]
```

### Panduan Pemilihan Stack (Language-Agnostic)

| Kebutuhan | Pertimbangan Utama |
|---|---|
| **High concurrency** | Go, Rust, atau Node.js (event loop) lebih efisien dari thread-per-request |
| **Rapid prototyping** | Python / TypeScript — ekosistem library besar, iterasi cepat |
| **Type safety** | TypeScript, Go, Rust, Java/Kotlin — kurangi runtime error |
| **Data science / ML** | Python — ekosistem paling matang (numpy, pandas, torch) |
| **Mobile** | Swift (iOS), Kotlin (Android), atau Flutter (cross-platform) |
| **Embedded / Performa kritis** | Rust, C++ |
| **Enterprise / Long-term** | Java, Kotlin, C# — mature, banyak tooling enterprise |

### Prinsip Pemilihan (Berlaku di Semua Stack)
- **Pilih yang tim sudah kuasai** — produktivitas > "teknologi terbaik" secara teoritis
- **Pertimbangkan ekosistem** — library, komunitas, dokumentasi yang tersedia
- **Sesuaikan dengan kebutuhan scaling** — monolith dulu, pecah jika terbukti perlu
- **Avoid resume-driven development** — jangan pilih teknologi hanya karena trendy


## 1. 🎯 Prinsip Desain Sistem


### Prinsip Utama
- **Separation of Concerns** — setiap komponen punya tanggung jawab tunggal yang jelas
- **High Cohesion, Low Coupling** — komponen yang berhubungan erat disatukan, yang tidak terkait dipisahkan
- **Design for Failure** — asumsikan setiap komponen bisa gagal, rancang recovery-nya
- **Evolutionary Architecture** — desain yang bisa berkembang tanpa full rewrite
- **Prefer Boring Technology** — gunakan teknologi proven sebelum yang baru dan belum terbukti

### Trade-off yang Harus Selalu Dipertimbangkan

| Dimensi | Pilihan A | vs | Pilihan B |
|---|---|---|---|
| Konsistensi | Strong Consistency | ↔ | Eventual Consistency |
| Ketersediaan | High Availability | ↔ | Strong Consistency (CAP) |
| Performa | Low Latency | ↔ | High Throughput |
| Kompleksitas | Microservices | ↔ | Monolith |
| Biaya | Self-hosted | ↔ | Managed Service |

---

## 2. 🗺️ Proses Desain Sistem

### Langkah 1 — Klarifikasi Requirements

Sebelum menulis satu baris kode atau diagram, jawab:

```
Functional Requirements:
- Apa yang sistem HARUS bisa lakukan?
- User journey utamanya seperti apa?

Non-Functional Requirements:
- Berapa expected users / requests per second?
- Berapa target latency yang dapat diterima?
- Berapa SLA availability yang dibutuhkan? (99.9%? 99.99%?)
- Apakah ada constraint regulasi / compliance?
- Berapa budget infrastruktur?

Out of Scope:
- Apa yang TIDAK perlu dibangun sekarang?
```

### Langkah 2 — High-Level Design

```
1. Gambar komponen utama (bukan detail implementasi)
2. Tentukan bagaimana komponen berkomunikasi (sync/async)
3. Identifikasi data store yang dibutuhkan
4. Tentukan boundary antar komponen
```

### Langkah 3 — Deep Dive per Komponen

Fokus pada komponen yang:
- Paling kritikal untuk bisnis
- Paling berisiko secara teknikal
- Paling tidak familiar untuk tim

### Langkah 4 — Identifikasi Bottleneck & Single Points of Failure

```
Tanya untuk setiap komponen:
- "Apa yang terjadi jika komponen ini down?"
- "Apa yang terjadi jika load 10x dari sekarang?"
- "Apa yang terjadi jika network antara A dan B putus?"
```

---

## 3. 🧩 Pola Arsitektur

### Layered Architecture
```
┌─────────────────────────┐
│    Presentation Layer   │ ← API / UI
├─────────────────────────┤
│    Application Layer    │ ← Business Logic / Use Cases
├─────────────────────────┤
│      Domain Layer       │ ← Entities, Domain Rules
├─────────────────────────┤
│  Infrastructure Layer   │ ← DB, External Services, Cache
└─────────────────────────┘

Aturan: Dependency hanya boleh mengarah ke bawah (inward)
```

### Event-Driven Architecture
```
Producer → [Event Bus / Message Queue] → Consumer(s)

Kapan dipakai:
✅ Decoupling antar service
✅ Async processing yang berat
✅ Audit trail dan event sourcing
❌ Jika butuh response synchronous yang simple
```

### CQRS (Command Query Responsibility Segregation)
```
Write Path: Command → Command Handler → Write DB
Read Path:  Query  → Query Handler  → Read DB (bisa read replica)

Kapan dipakai:
✅ Read/write load sangat berbeda
✅ Perlu optimasi read yang kompleks
❌ Overkill untuk CRUD sederhana
```

### Strangler Fig Pattern (Migrasi Legacy)
```
Legacy System ──→ Facade/Proxy ──→ New System
                      ↑
               Routing berdasarkan fitur
               yang sudah dimigrasikan

Gunakan saat migrasi bertahap dari sistem lama
```

---

## 4. 📊 Data Design

### Pemilihan Database

| Kebutuhan | Rekomendasi | Contoh |
|---|---|---|
| Relational data, ACID transactions | PostgreSQL | User, Order, Payment |
| Document store, flexible schema | MongoDB | Product catalog, CMS |
| Key-value, caching | Redis | Session, rate limiting |
| Full-text search | Elasticsearch | Search fitur |
| Time-series data | InfluxDB / TimescaleDB | Metrics, IoT |
| Graph data | Neo4j | Social network, recommendation |
| File / object storage | S3 / GCS | Image, video, dokumen |

### Prinsip Data Modeling
- [ ] Desain schema berdasarkan **access pattern**, bukan hanya struktur data
- [ ] Tentukan **data ownership** yang jelas — satu service = satu source of truth
- [ ] Rencanakan **data migration strategy** sejak awal
- [ ] Pertimbangkan **soft delete** vs hard delete untuk data audit
- [ ] Tentukan **retention policy** untuk data lama

---

## 5. 🔄 Komunikasi Antar Komponen

### Synchronous (Request-Response)
```
REST API       → Simple, widely understood, stateless
gRPC           → High performance, strongly typed, streaming
GraphQL        → Flexible queries, reduce over-fetching

Kapan: Butuh response langsung, operasi sederhana
```

### Asynchronous (Message-Based)
```
Message Queue  → Point-to-point (RabbitMQ, SQS)
Event Streaming → Pub/Sub, replay (Kafka, Kinesis)
Webhook        → Push notification ke external system

Kapan: Long-running task, decoupling, resilience
```

### Checklist Komunikasi
- [ ] Definisikan **API contract** sebelum implementasi (API-first)
- [ ] Implementasi **retry dengan exponential backoff** untuk external calls
- [ ] Gunakan **circuit breaker** untuk mencegah cascade failure
- [ ] Tentukan **timeout** yang eksplisit untuk setiap external call
- [ ] Dokumentasikan **error response format** yang konsisten

---

## 6. ⚡ Scalability & Reliability

### Strategi Scaling

```
Vertical Scaling   → Tambah resource (CPU/RAM) — batas hardware
Horizontal Scaling → Tambah instance — butuh stateless design

Read Scaling:
- Read replica database
- CDN untuk static asset
- Caching layer (Redis, Memcached)

Write Scaling:
- Database sharding
- CQRS
- Event sourcing
```

### Checklist Reliability
- [ ] **Health check endpoint** tersedia untuk semua service
- [ ] **Graceful shutdown** diimplementasi (drain request sebelum stop)
- [ ] **Circuit breaker** untuk dependency eksternal
- [ ] **Idempotency** untuk operasi yang bisa di-retry
- [ ] **Distributed tracing** untuk debug lintas service
- [ ] **SLO (Service Level Objective)** terdefinisi dan dimonitor

---

## 7. 📝 Architecture Decision Records (ADR)

Setiap keputusan arsitektur signifikan harus didokumentasikan sebagai ADR.

### Template ADR

```markdown
# ADR-[NNN]: [Judul Keputusan]

**Status:** [Proposed | Accepted | Deprecated | Superseded]
**Tanggal:** YYYY-MM-DD
**Deciders:** [Nama/Tim yang terlibat]

## Konteks
[Apa situasi atau masalah yang membutuhkan keputusan ini?]

## Keputusan
[Apa keputusan yang diambil?]

## Alasan
[Mengapa keputusan ini dipilih?]

## Konsekuensi
**Positif:**
- [Benefit dari keputusan ini]

**Negatif / Trade-off:**
- [Downside yang diterima]

## Alternatif yang Dipertimbangkan
- [Opsi A] — ditolak karena [alasan]
- [Opsi B] — ditolak karena [alasan]
```

---

## 8. 📋 Architecture Review Checklist

Gunakan sebelum finalisasi desain:

### Requirements
- [ ] Functional requirements terdokumentasi dan disetujui stakeholder
- [ ] Non-functional requirements (latency, throughput, availability) terdefinisi
- [ ] Scope dan out-of-scope jelas

### Arsitektur
- [ ] Single points of failure teridentifikasi dan ada mitigation plan
- [ ] Data ownership per komponen/service sudah jelas
- [ ] Strategi komunikasi antar komponen terdefinisi (sync/async)
- [ ] Error handling dan retry strategy sudah dirancang

### Data
- [ ] Schema database dirancang berdasarkan access pattern
- [ ] Strategi backup dan recovery terdefinisi
- [ ] Data migration plan tersedia jika mengubah schema yang ada

### Operasional
- [ ] Monitoring dan alerting sudah direncanakan
- [ ] Deployment strategy terdefinisi (blue-green, canary, rolling)
- [ ] Runbook untuk incident response tersedia
- [ ] ADR untuk keputusan signifikan sudah ditulis

### Reference Check
- [ ] File referensi project di `references/<designname>-design.md` sudah diperbarui
- [ ] Keputusan baru sudah ditambahkan ke ADR list di file referensi

---

## 📚 Referensi

- [System Design Primer](https://github.com/donnemartin/system-design-primer)
- [Architecture Decision Records](https://adr.github.io/)
- [Martin Fowler — Architecture Guide](https://martinfowler.com/architecture/)
- [The Twelve-Factor App](https://12factor.net/)
- [CAP Theorem](https://en.wikipedia.org/wiki/CAP_theorem)

---

*Architecture Harness v1.0 — bagian dari [Harness Engineering](../README.md)*
