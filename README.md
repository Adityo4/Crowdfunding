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

## ⚡ Prompt Pembuka

> `agent.md` sudah punya **Context Map** yang menentukan harness apa yang perlu dibaca berdasarkan task.
> Kamu tidak perlu tentukan harness secara manual — biarkan model yang memutuskan.

---

### Variant A — AI dengan File Access ✅ (Recommended)
*Untuk: Cursor, Windsurf, AI agent dengan akses filesystem*

**1. Prompt Pertama (Inisialisasi):**
```
Baca harnesses/agent.md terlebih dahulu untuk memahami standar, workflow, dan Context Map proyek ini. Konfirmasi jika sudah selesai membaca.
```

**2. Prompt Kedua dan Seterusnya (Task-driven):**
```
TASK: [deskripsikan apa yang ingin dikerjakan]
```

**Kenapa pendekatan ini jauh lebih baik:**
* **Tanpa Token Leak:** AI hanya akan memanggil/membaca harness tambahan (seperti `database.md`, `architecture.md`, dll.) secara dinamis saat dibutuhkan oleh `TASK` Anda, berkat **Context Map** yang ada di `agent.md`.
* **Konteks Terjaga:** Model langsung terbiasa dengan kepribadian/vibe coding sesuai arahan di `agent.md` sejak interaksi pertama.

---

### Variant B — Chat Interface (Paste Manual)
*Untuk: Claude.ai, ChatGPT, Gemini web — tidak bisa baca file lokal*

**1. Prompt Pertama (Inisialisasi):**
```
Kamu adalah Harness-Aware Engineering Agent. Baca agent.md berikut untuk memahami workflow proyek ini:

---
[paste isi agent.md di sini]
---

Konfirmasi jika kamu sudah paham dan siap menerima TASK pertama.
```

**2. Prompt Kedua dan Seterusnya (Task-driven):**
```
TASK: [deskripsikan apa yang ingin dikerjakan]
```

> Berdasarkan `TASK` yang diberikan nanti, AI akan otomatis meminta Anda mem-paste file harness yang relevan (misalnya: *"Saya butuh harnesses/database.md untuk menyelesaikan task ini, tolong paste isinya"*). Anda tidak perlu menentukan secara manual.


---

## 🚀 Cara Penggunaan

### 1. Pilih template prompt yang sesuai task
Gunakan tabel di atas atau lihat **Context Map** di [agent.md](./harnesses/agent.md).

### 2. Paste harness yang relevan
Copy isi file harness dan tempel di dalam prompt — **jangan paste semua harness sekaligus**.

### 3. Kustomisasi bagian `[KUSTOMISASI]`
Setiap harness punya bagian bertanda `[KUSTOMISASI]` — isi sesuai stack teknologi project.

### 4. Untuk design project spesifik
Tambahkan file dari [references/](./harnesses/references/) khusus untuk `design.md`.

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
