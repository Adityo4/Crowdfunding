# 🔩 Harness Engineering — Vibe Code Template System

> Kumpulan harness reusable untuk membangun sistem yang aman, scalable, dan maintainable.  
> Setiap harness adalah *blueprint* terstruktur yang bisa langsung digunakan sebagai konteks AI coding (vibe coding).

---

## 📁 Struktur

```
harness-engineering/
├── README.md                    ← Kamu lagi baca ini
├── harnesses/
│   ├── agent.md                 ← 🚪 ENTRY POINT — baca ini pertama
│   ├── security.md              ← Keamanan sistem (auth, enkripsi, input validation, dll)
│   ├── design.md                ← UI/UX design style (token, tipografi, komponen, a11y)
│   ├── architecture.md          ← Arsitektur sistem (pola, ADR, scalability)
│   ├── feature.md               ← Template spec fitur siap pakai
│   ├── testing.md               ← Strategi testing (unit, integration, E2E, TDD)
│   ├── performance.md           ← Optimasi performa & load testing
│   ├── devops.md                ← CI/CD, Docker, deployment, observability
│   ├── api-design.md            ← Desain API & kontrak
│   ├── code-quality.md          ← Naming, struktur kode, code review, refactoring
│   ├── error-handling.md        ← Error taxonomy, propagation, logging, retry
│   ├── modeling.md              ← Flow, diagram, sequence, state, ER, Mermaid templates
│   ├── database.md              ← Schema design, indexing, caching, migration, query optimization
│   └── references/              ← File referensi spesifik per project
│       └── <designname>-design.md
└── templates/
    └── harness-template.md      ← Template kosong untuk membuat harness baru
```

---

## 🚀 Cara Penggunaan

### 1. Pilih harness yang relevan
Sesuaikan dengan aspek sistem yang sedang kamu bangun.

### 2. Paste ke konteks AI
Copy isi harness yang dipilih dan sertakan sebagai **system prompt** atau **konteks awal** saat vibe coding.

### 3. Kombinasikan beberapa harness
Harness bisa dikombinasikan. Contoh untuk membangun backend API:
```
security.md + api-design.md + testing.md
```

### 4. Kustomisasi
Setiap harness punya bagian `[KUSTOMISASI]` — sesuaikan dengan stack teknologi proyekmu.

---

## 🎯 Prinsip Harness Engineering

| Prinsip | Deskripsi |
|---|---|
| **Reusable** | Satu harness bisa dipakai di banyak proyek |
| **Composable** | Harness bisa digabungkan sesuai kebutuhan |
| **Opinionated** | Punya best practice yang jelas, bukan sekedar checklist |
| **AI-Ready** | Ditulis agar mudah dipahami oleh LLM saat vibe coding |
| **Living Document** | Harness terus diperbarui seiring perkembangan industri |

---

## 📋 Harness yang Tersedia

| Harness | Status | Deskripsi Singkat |
|---|---|---|
| [agent.md](./harnesses/agent.md) | ✅ Ready | 🚪 Entry point — perilaku AI, learning loop, context map |
| [security.md](./harnesses/security.md) | ✅ Ready | Auth, enkripsi, validasi input, OWASP |
| [design.md](./harnesses/design.md) | ✅ Ready | Color tokens, tipografi, komponen, animasi, a11y |
| [architecture.md](./harnesses/architecture.md) | ✅ Ready | Pola arsitektur, ADR, scalability, data design, tech stack |
| [feature.md](./harnesses/feature.md) | ✅ Ready | Template spec fitur siap pakai |
| [testing.md](./harnesses/testing.md) | ✅ Ready | Unit, integration, E2E, TDD, coverage strategy |
| [performance.md](./harnesses/performance.md) | ✅ Ready | Caching, query optimization, profiling, load testing |
| [devops.md](./harnesses/devops.md) | ✅ Ready | CI/CD, Docker, deployment strategy, observability |
| [api-design.md](./harnesses/api-design.md) | ✅ Ready | REST/GraphQL best practices, versioning, kontrak API |
| [code-quality.md](./harnesses/code-quality.md) | ✅ Ready | Naming, struktur kode, code review, refactoring |
| [error-handling.md](./harnesses/error-handling.md) | ✅ Ready | Error taxonomy, propagation, logging, retry strategy |
| [modeling.md](./harnesses/modeling.md) | ✅ Ready | User flow, sequence, state, ER, decision tree — Mermaid templates |
| [database.md](./harnesses/database.md) | ✅ Ready | Schema design, indexing, caching Redis, migration, query optimization |

---

## 🛠️ Membuat Harness Baru

Gunakan [harness-template.md](./templates/harness-template.md) sebagai titik awal.

---

*Dibuat dengan ❤️ menggunakan prinsip Harness Engineering*
