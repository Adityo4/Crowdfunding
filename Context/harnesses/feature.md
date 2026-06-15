# 🚀 Feature Harness

> **Tujuan:** Template siap pakai untuk mendefinisikan, merancang, dan mengimplementasikan sebuah fitur secara terstruktur.
> **Gunakan harness ini saat:** memulai pengerjaan fitur baru, menulis spec fitur, atau melakukan feature review.

---

## 🧭 Konteks untuk AI

```
Kamu adalah product engineer yang berpengalaman. Saat membangun sebuah fitur:
1. Pahami "mengapa" fitur ini dibutuhkan sebelum "bagaimana" cara membuatnya
2. Pecah fitur menjadi task-task kecil yang independen dan bisa di-deliver secara bertahap
3. Selalu pikirkan edge case dan error state, bukan hanya happy path
4. Implementasi harus language-agnostic dalam prinsip — sesuaikan dengan stack yang digunakan
5. Tulis kode yang bisa dibaca manusia lain, bukan hanya mesin

Sebelum mulai coding, SELALU tanya: "Apa definisi 'done' untuk fitur ini?"
```

---

## 📋 Feature Spec Template

> **Cara pakai:** Duplikat bagian di bawah ini untuk setiap fitur baru. Isi semua bagian sebelum mulai coding.

---

### 🏷️ [NAMA FITUR]

**Epic/Modul:** `[nama epic atau modul induk]`
**Status:** `[ ] Draft` | `[ ] Review` | `[ ] Approved` | `[ ] In Progress` | `[ ] Done`
**Priority:** `[ ] Critical` | `[ ] High` | `[ ] Medium` | `[ ] Low`
**Estimasi:** `[X hari / X minggu]`
**Owner:** `[nama / tim]`

---

#### 📌 Problem Statement

> Jelaskan masalah yang ingin diselesaikan. Gunakan format "Sebagai [user], saya kesulitan [masalah], sehingga [dampak]."

```
Sebagai [tipe user],
saya [masalah atau kebutuhan yang dirasakan],
sehingga [dampak negatif jika tidak diselesaikan].
```

---

#### 🎯 Goals & Success Criteria

**Goals:**
- [ ] [Tujuan spesifik 1 yang ingin dicapai fitur ini]
- [ ] [Tujuan spesifik 2]
- [ ] [Tujuan spesifik 3]

**Definition of Done:**
- [ ] [Kriteria terukur 1 — contoh: user bisa login dalam < 3 langkah]
- [ ] [Kriteria terukur 2 — contoh: response API < 200ms]
- [ ] [Kriteria terukur 3 — contoh: test coverage ≥ 80%]

**Out of Scope:**
- [Hal yang TIDAK termasuk dalam fitur ini]
- [Hal lain yang di-exclude]

---

#### 👤 User Stories

```
Story 1:
Sebagai [tipe user],
saya ingin [aksi],
agar [manfaat yang didapat].

Acceptance Criteria:
- GIVEN [kondisi awal]
  WHEN [aksi dilakukan]
  THEN [hasil yang diharapkan]

- GIVEN [kondisi lain]
  WHEN [aksi lain]
  THEN [hasil lain]
```

```
Story 2:
Sebagai [tipe user],
saya ingin [aksi],
agar [manfaat].

Acceptance Criteria:
- GIVEN [kondisi]
  WHEN [aksi]
  THEN [hasil]
```

---

#### 🗺️ User Flow

```
[Deskripsikan atau gambarkan alur user dari awal hingga akhir]

Contoh:
User buka halaman → Klik tombol X → Isi form → Submit → Muncul konfirmasi → Selesai

Atau dalam format diagram:
[Start] → [Langkah 1] → [Langkah 2] → [Decision?]
                                            ├── Yes → [Langkah 3] → [End]
                                            └── No  → [Error State] → [Retry]
```

---

#### 🔧 Technical Design

**Komponen / Module yang Terlibat:**
- `[NamaKomponen]` — [peran dan tanggung jawabnya]
- `[NamaModule]` — [peran dan tanggung jawabnya]

**Perubahan pada Data Model:**
```
[Deskripsikan perubahan schema/model jika ada]

Contoh:
Tabel users: tambah kolom `last_login_at` (timestamp, nullable)
Tabel sessions: buat tabel baru dengan kolom id, user_id, token, expires_at
```

**API / Interface yang Dibuat atau Diubah:**
```
[Metode] [Endpoint/Fungsi]
  Input:  [parameter]
  Output: [return value / response]
  Error:  [kemungkinan error]

Contoh:
POST /api/auth/login
  Input:  { email: string, password: string }
  Output: { token: string, user: UserObject }
  Error:  401 Invalid credentials | 429 Rate limit exceeded
```

**Dependency & Integrasi:**
- [Library / service eksternal yang digunakan]
- [Service internal lain yang dipanggil]

---

#### ⚠️ Edge Cases & Error States

| Skenario | Perilaku yang Diharapkan |
|---|---|
| [Edge case 1 — contoh: input kosong] | [Tampilkan pesan error spesifik] |
| [Edge case 2 — contoh: network timeout] | [Retry otomatis / tampilkan fallback] |
| [Edge case 3 — contoh: user tidak punya akses] | [Redirect ke halaman unauthorized] |
| [Edge case 4 — contoh: data duplikat] | [Tampilkan konfirmasi atau merge] |

---

#### 🧪 Testing Plan

**Unit Tests:**
- [ ] [Fungsi/komponen yang perlu unit test]
- [ ] [Fungsi/komponen lain]

**Integration Tests:**
- [ ] [Skenario integrasi yang perlu diuji]

**Manual Test Cases:**
- [ ] Happy path: [langkah-langkah test case utama]
- [ ] Error path: [langkah-langkah test case error]
- [ ] Edge case: [langkah-langkah test case edge]

---

#### 📊 Metrics & Monitoring

**Metrics yang Perlu Dipantau:**
- [Metric 1 — contoh: jumlah login per hari]
- [Metric 2 — contoh: error rate endpoint ini]

**Alert yang Perlu Ditambahkan:**
- [ ] Alert jika [kondisi abnormal 1]
- [ ] Alert jika [kondisi abnormal 2]

---

#### 🔗 Referensi & Asset

- **Figma / Design:** [link]
- **PRD / Dokumen:** [link]
- **Ticket:** [link ke Jira / Linear / GitHub Issue]
- **Related Features:** [link ke feature lain yang berhubungan]

---

#### 📝 Catatan & Keputusan

> Catat keputusan penting, perubahan arah, atau hal-hal yang perlu diingat selama pengerjaan.

- `[tanggal]` — [catatan atau keputusan yang dibuat]
- `[tanggal]` — [perubahan scope atau arah]

---

## 🗂️ Daftar Fitur Aktif

> Gunakan tabel ini sebagai index fitur yang sedang dan sudah dikerjakan.

| Fitur | Status | Priority | Owner | Link |
|---|---|---|---|---|
| [Nama Fitur 1] | In Progress | High | [nama] | [#section] |
| [Nama Fitur 2] | Done | Medium | [nama] | [#section] |
| [Nama Fitur 3] | Draft | Low | [nama] | [#section] |

---

*Feature Harness v1.0 — bagian dari [Harness Engineering](../README.md)*
