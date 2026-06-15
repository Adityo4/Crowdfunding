# Genesis — Design Reference

> File referensi ini digunakan bersama [design.md](../design.md).
> Berisi konteks spesifik project Genesis yang akan di-*lookup* oleh AI saat vibe coding.

---

## Project Context

**Genesis** adalah [deskripsikan sistem kamu di sini — contoh: platform manajemen data engineering yang memungkinkan tim data untuk mendefinisikan, menjadwalkan, dan memonitor pipeline data].

**Tim:** [Nama tim / jumlah engineer]
**Stage:** [Planning | Development | Staging | Production]
**Stack Utama:** [Contoh: Node.js, PostgreSQL, Redis, Kafka]

---

## Key Decisions

| Keputusan | Pilihan | Alasan Singkat |
|---|---|---|
| [Contoh] Database utama | PostgreSQL | ACID transactions, relational data |
| [Contoh] Message queue | Kafka | High throughput, event replay |
| [Contoh] Auth | JWT + OAuth2 | Stateless, integrasi SSO |

---

## Constraints

- **Teknikal:** [Contoh: harus compatible dengan infra existing di AWS]
- **Bisnis:** [Contoh: budget infrastruktur maks $X/bulan]
- **Tim:** [Contoh: tim 3 orang, belum berpengalaman dengan Kubernetes]
- **Waktu:** [Contoh: MVP harus live dalam 3 bulan]

---

## Component Diagram

```
[Gambar atau deskripsi arsitektur komponen utama]

Contoh:
┌──────────┐    REST    ┌───────────────┐    gRPC    ┌──────────────┐
│  Client  │ ─────────▶│   API Gateway │ ──────────▶│  Auth Service│
└──────────┘            └───────┬───────┘            └──────────────┘
                                │
                    ┌───────────┼───────────┐
                    ▼           ▼           ▼
             ┌──────────┐ ┌─────────┐ ┌──────────┐
             │ Service A│ │Service B│ │Service C │
             └────┬─────┘ └────┬────┘ └────┬─────┘
                  │            │            │
                  └────────────┼────────────┘
                               ▼
                        ┌────────────┐
                        │ PostgreSQL │
                        └────────────┘
```

---

## ADR (Architecture Decision Records)

### ADR-001: [Judul Keputusan Pertama]

**Status:** Accepted
**Tanggal:** YYYY-MM-DD

**Konteks:** [Situasi yang membutuhkan keputusan]

**Keputusan:** [Apa yang diputuskan]

**Alasan:** [Mengapa]

**Trade-off:**
- ✅ [Benefit]
- ⚠️ [Downside yang diterima]

---

### ADR-002: [Judul Keputusan Kedua]

**Status:** Proposed
**Tanggal:** YYYY-MM-DD

**Konteks:** [...]

**Keputusan:** [...]

---

## Open Questions

> Pertanyaan desain yang belum dijawab — perbarui saat sudah ada keputusan.

- [ ] [Pertanyaan 1 yang masih open]
- [ ] [Pertanyaan 2 yang masih open]

---

*Genesis Design Reference v0.1 — digunakan bersama [design.md](../design.md)*
