# 📋 Progress — Harness Engineering

> Dokumen ini mencatat status pengerjaan harness dan rencana ke depan.
> **Update setiap kali** harness baru selesai dibuat atau ada perubahan signifikan.

---

## 🗓️ Last Updated
`2026-06-16`

---

## ✅ Selesai

- [x] **README.md** — Indeks & panduan penggunaan sistem harness
- [x] **agent.md** — Meta-harness: perilaku AI, learning loop, entry point & context map
- [x] **security.md** — Auth, enkripsi, validasi input, OWASP Top 10
- [x] **design.md** — UI/UX design style: color tokens, tipografi, komponen, animasi, aksesibilitas
- [x] **architecture.md** — System design: pola arsitektur, ADR, scalability, data design, **tech stack section**
- [x] **feature.md** — Template spec fitur siap pakai (problem statement, user stories, technical design, testing)
- [x] **progress.md** — Dokumen tracking todo list & progress harness
- [x] **templates/harness-template.md** — Template kosong untuk harness baru
- [x] **references/genesis-design.md** — Contoh file referensi design project
- [x] **testing.md** — Strategi testing: unit, integration, E2E, TDD, coverage
- [x] **performance.md** — Caching, query optimization, profiling, load testing
- [x] **devops.md** — CI/CD pipeline, Docker, monitoring, deployment strategy
- [x] **api-design.md** — REST/GraphQL best practices, versioning, kontrak API
- [x] **code-quality.md** — Naming convention, struktur kode, code review checklist, refactoring
- [x] **error-handling.md** — Error taxonomy, propagation pattern, logging, retry & fallback
- [x] **modeling.md** — User flow, sequence, state, ER, decision tree, component diagram — Mermaid templates
- [x] **database.md** — Schema design, indexing strategy, Redis caching patterns, migration, query optimization

---

## 🔄 In Progress

- [ ] **Inisialisasi Proyek Crowdfunding**
  - [x] Desain Skema Database (ERD) & DDL (`database/schema.sql`)
  - [x] Desain API Specification (`docs/api_specification.md`)
  - [x] Setup Infrastruktur dengan Docker Compose (PostgreSQL, Redis, MinIO)
  - [x] Inisialisasi Backend Go (Gin & GORM setup)
  - [x] Implementasi API Backend (Gin + GORM)
  - [x] Inisialisasi Frontend Nuxt.js 3 & Tailwind CSS (Slicing Index, Login/Register, Charities List/Detail/Create, Articles List/Detail)

---

## 📌 Backlog

### Harness Baru
*Semua harness yang direncanakan sudah selesai! 🎉*


### Peningkatan Harness yang Ada
- [ ] `security.md` — Tambah contoh implementasi per bahasa (Node, Python, Go)
- [ ] `design.md` — Tambah section animasi Framer Motion & CSS animation comparison
- [ ] `feature.md` — Tambah contoh feature spec yang sudah terisi lengkap
- [ ] `architecture.md` — Tambah section deployment diagram

### Infrastruktur
- [ ] Buat `references/` template yang lebih lengkap untuk design.md
- [ ] Buat script atau CLI helper untuk generate harness baru dari template

---

## 💡 Ideas & Notes

> Ide yang muncul selama pengerjaan — belum tentu masuk backlog.

- Pertimbangkan membuat **harness combinator** — file yang menggabungkan beberapa harness untuk use case spesifik (contoh: `fullstack-app.md` = security + architecture + api-design)
- Bisa ditambahkan tag di setiap harness untuk memudahkan filter (e.g., `#frontend`, `#backend`, `#mobile`)

---

## 📈 Statistik

| Kategori | Jumlah |
|---|---|
| Harness selesai | 14 |
| Harness in progress | 0 |
| Harness backlog | 0 |
| File referensi | 1 |

---

*Progress log — [Harness Engineering](./README.md)*

