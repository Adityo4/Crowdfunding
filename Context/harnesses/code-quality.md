# 🧹 Code Quality Harness

> **Tujuan:** Memastikan setiap baris kode yang ditulis bersih, mudah dibaca, dan maintainable — bukan hanya sekadar functional.
> **Gunakan harness ini saat:** menulis kode baru, melakukan code review, refactoring, atau menetapkan standar kode project.

---

## 🧭 Konteks untuk AI

```
Kamu adalah engineer yang peduli dengan kualitas kode jangka panjang.
Kode dibaca jauh lebih sering daripada ditulis — optimasi untuk pembaca, bukan penulis.

Prinsipmu:
1. Kode yang bersih lebih penting dari kode yang pintar
2. Nama yang baik adalah dokumentasi terbaik
3. Fungsi kecil yang jelas lebih baik dari fungsi besar yang kompleks
4. Jika perlu komentar untuk menjelaskan WHAT, namanya kurang baik
5. Komentar untuk menjelaskan WHY — kode untuk menjelaskan WHAT dan HOW

Setiap kali menulis kode, tanya: "Akankah engineer lain memahami ini 6 bulan kemudian?"
```

---

## 1. ✍️ Naming Convention

### Prinsip
- Nama harus **mengungkapkan intent** — bukan cara kerja internal
- Hindari abbreviasi kecuali yang sudah universal (`id`, `url`, `i`)
- Konsisten dalam satu codebase — pilih satu konvensi dan ikuti

### Aturan Penamaan per Jenis

```
Variabel & Parameter:
  ✅ userAccountBalance   ❌ bal, x, temp, data
  ✅ isEmailVerified      ❌ flag, check, status
  ✅ filteredProducts     ❌ arr, list, result

Fungsi / Method:
  ✅ getUserById()        ❌ get(), fetch(), process()
  ✅ calculateTotalPrice()❌ calc(), doStuff()
  ✅ isValidEmail()       ❌ emailCheck(), validate()
  → Gunakan verb + noun yang deskriptif
  → Fungsi boolean mulai dengan: is, has, can, should, will

Kelas / Type:
  ✅ UserRepository       ❌ UserRepo, UR, Helper
  ✅ PaymentProcessor     ❌ PaymentHandler2, Manager
  → Gunakan noun yang spesifik, hindari generic: Manager, Helper, Util

Konstanta:
  ✅ MAX_RETRY_ATTEMPTS   ❌ maxRetry, MAX, three
  ✅ DEFAULT_TIMEOUT_MS   ❌ timeout, TIMEOUT
  → SCREAMING_SNAKE_CASE untuk nilai yang benar-benar konstanta

File / Module:
  ✅ user-repository.ts   ❌ userRepo.ts, UR.ts
  ✅ calculate-tax.ts     ❌ taxHelper.ts, utils.ts
  → kebab-case untuk file, cocok untuk semua OS
```

### [KUSTOMISASI] Konvensi per Bahasa
```
JavaScript/TypeScript:
  variabel & fungsi → camelCase
  kelas & tipe      → PascalCase
  konstanta         → SCREAMING_SNAKE_CASE
  file              → kebab-case

Python:
  variabel & fungsi → snake_case
  kelas             → PascalCase
  konstanta         → SCREAMING_SNAKE_CASE
  file/modul        → snake_case

Go:
  exported          → PascalCase
  unexported        → camelCase
  konstanta         → PascalCase (exported) / camelCase (unexported)

Java/Kotlin:
  variabel & fungsi → camelCase
  kelas             → PascalCase
  konstanta         → SCREAMING_SNAKE_CASE
  paket             → lowercase.dots
```

---

## 2. 🧩 Fungsi & Metode

### Prinsip Single Responsibility

```
Setiap fungsi harus bisa dijelaskan dalam SATU kalimat tanpa kata "dan".

❌ "Fungsi ini memvalidasi input, menyimpan ke database, DAN mengirim email"
   → Ini 3 fungsi berbeda yang disatukan

✅ "Fungsi ini memvalidasi format email"
✅ "Fungsi ini menyimpan user ke database"
✅ "Fungsi ini mengirim welcome email ke user baru"
```

### Ukuran Fungsi

```
Panduan ukuran (bukan aturan kaku):
  Ideal:    5–15 baris
  Wajar:    15–30 baris
  Review:   > 30 baris → pertimbangkan pecah jadi sub-fungsi
  Refactor: > 50 baris → hampir pasti perlu dipecah

Bukan soal jumlah baris — tapi soal level abstraksi:
Satu fungsi = satu level abstraksi
```

### Level Abstraksi yang Konsisten

```javascript
// ❌ SALAH — mencampur level abstraksi
async function processOrder(order) {
  // High-level
  await validateOrder(order);

  // Low-level (seharusnya di dalam saveOrder())
  const query = 'INSERT INTO orders (id, user_id, total) VALUES ($1, $2, $3)';
  await db.query(query, [order.id, order.userId, order.total]);

  // High-level lagi
  await sendConfirmationEmail(order);
}

// ✅ BENAR — konsisten di satu level abstraksi
async function processOrder(order) {
  await validateOrder(order);
  await saveOrder(order);
  await sendConfirmationEmail(order);
}
```

### Checklist Fungsi
- [ ] Fungsi punya satu tanggung jawab yang jelas
- [ ] Nama fungsi menjelaskan apa yang dilakukan, bukan bagaimana
- [ ] Parameter maksimal 3–4 — lebih dari itu, pertimbangkan object/struct
- [ ] Tidak ada side effect tersembunyi (fungsi melakukan lebih dari yang namanya janjikan)
- [ ] Return type konsisten (tidak kadang return nilai, kadang null, kadang throw)

---

## 3. 🏗️ Struktur & Organisasi Kode

### Prinsip

```
Kode yang berhubungan harus berdekatan.
Kode yang tidak berhubungan harus terpisah.

Urutan yang baik dalam sebuah file:
  1. Import / dependency
  2. Konstanta & konfigurasi
  3. Tipe / interface / schema
  4. Fungsi helper / private
  5. Fungsi utama / public / export
  6. Default export (jika ada)
```

### Struktur Folder

```
Pilih satu pendekatan dan konsisten:

Pendekatan 1 — By Feature (direkomendasikan untuk aplikasi kompleks):
  src/
  ├── features/
  │   ├── auth/
  │   │   ├── auth.controller.ts
  │   │   ├── auth.service.ts
  │   │   ├── auth.repository.ts
  │   │   └── auth.test.ts
  │   └── payment/
  │       ├── payment.controller.ts
  │       └── ...

Pendekatan 2 — By Layer (lebih familiar, cocok untuk project kecil):
  src/
  ├── controllers/
  ├── services/
  ├── repositories/
  └── models/
```

### Checklist Struktur
- [ ] File yang berhubungan berada di folder yang sama
- [ ] Tidak ada "god file" yang berisi segalanya (utils.ts, helpers.ts yang ratusan baris)
- [ ] Import path konsisten (absolute atau relative, pilih satu)
- [ ] Tidak ada circular dependency antar module

---

## 4. 💬 Komentar & Dokumentasi Inline

### Komentar yang BAIK vs BURUK

```javascript
// ❌ BURUK — menjelaskan WHAT (sudah jelas dari kode)
// Increment i by 1
i++;

// ❌ BURUK — informasi yang sudah basi / tidak akurat
// TODO: fix this later (ditulis 2 tahun lalu, tidak pernah di-fix)

// ✅ BAIK — menjelaskan WHY (tidak bisa dibaca dari kode)
// Gunakan setTimeout 0 untuk memastikan DOM sudah di-render
// sebelum kita query element-nya (browser rendering quirk)
setTimeout(() => queryElement(), 0);

// ✅ BAIK — menjelaskan keputusan yang tidak obvious
// Sengaja tidak pakai index untuk kolom ini —
// tabel ini sangat sering di-INSERT dan jarang di-query by this field
```

### JSDoc / Docstring (untuk public API)

```javascript
/**
 * Menghitung total harga setelah diskon dan pajak.
 *
 * @param subtotal - Harga sebelum diskon, dalam rupiah
 * @param discountPercent - Persentase diskon (0-100)
 * @param taxRate - Rate pajak (0-1, contoh: 0.11 untuk 11%)
 * @returns Total harga final dalam rupiah
 * @throws {ValidationError} Jika subtotal negatif atau discount > 100
 *
 * @example
 * calculateTotal(100000, 10, 0.11) // → 99900
 */
function calculateTotal(subtotal, discountPercent, taxRate) { ... }
```

### Checklist Komentar
- [ ] Tidak ada komentar yang hanya menjelaskan ulang kode
- [ ] Setiap TODO memiliki konteks dan idealnya nomor tiket
- [ ] Public API / fungsi yang diekspor punya docstring
- [ ] Kode yang kompleks atau tidak obvious punya komentar WHY

---

## 5. 🔍 Code Review Checklist

### Sebelum Submit PR / MR

**Correctness**
- [ ] Kode melakukan apa yang dimaksudkan
- [ ] Edge case ditangani (null, empty, batas nilai)
- [ ] Error ditangani dengan benar (lihat [error-handling.md](./error-handling.md))
- [ ] Tidak ada logic yang salah atau off-by-one error

**Readability**
- [ ] Nama variabel, fungsi, dan kelas deskriptif
- [ ] Fungsi tidak terlalu panjang atau melakukan terlalu banyak hal
- [ ] Komentar ada di tempat yang tepat (WHY, bukan WHAT)
- [ ] Tidak ada kode yang di-comment out tanpa penjelasan

**Maintainability**
- [ ] Tidak ada duplikasi kode yang tidak perlu (DRY)
- [ ] Tidak ada magic number / magic string tanpa named constant
- [ ] Dependency injection digunakan (bukan hardcode dependency)
- [ ] Kode tidak over-engineered untuk requirement yang ada

**Test**
- [ ] Ada test untuk perubahan yang dibuat
- [ ] Test mencakup happy path dan error path
- [ ] Test tidak rapuh (tidak bergantung pada implementasi detail)

**Security & Performance**
- [ ] Tidak ada credential atau secret yang ter-commit
- [ ] Input dari user divalidasi
- [ ] Tidak ada N+1 query yang diperkenalkan
- [ ] Tidak ada resource leak (koneksi, file handle yang tidak ditutup)

---

## 6. 🔄 Refactoring

### Kapan Refactor

```
✅ Saat kamu menyentuh kode itu (Boy Scout Rule: tinggalkan lebih bersih dari sebelum)
✅ Saat test sudah ada dan hijau (refactor dengan safety net)
✅ Sebelum menambah fitur baru ke kode yang messy
❌ Jangan refactor dan tambah fitur sekaligus (buat perubahan terpisah)
❌ Jangan refactor kode yang tidak punya test (terlalu berisiko)
```

### Code Smells yang Paling Umum

```
Long Method          → Fungsi > 30 baris, pertimbangkan pecah
Large Class          → Kelas dengan terlalu banyak tanggung jawab
Duplicate Code       → Kode yang sama di beberapa tempat → extract
Long Parameter List  → > 4 parameter → bungkus dalam object
Magic Numbers        → 86400, 0.15, 7 tanpa nama → buat konstanta
Dead Code            → Kode yang tidak pernah dieksekusi → hapus
Feature Envy         → Method lebih sering pakai data class lain → pindahkan
God Class/File       → Satu file/kelas yang tahu segalanya → pecah
```

### Refactoring yang Aman

```
Langkah yang aman:
1. Pastikan ada test yang cover kode yang akan di-refactor
2. Buat perubahan kecil dan incremental
3. Jalankan test setelah setiap perubahan
4. Commit setelah setiap refactoring yang stabil
5. Jangan ubah behaviour — hanya struktur kode
```

---

## 7. 🛠️ Linting & Formatting

### Prinsip
- Formatting adalah **zero-opinion dengan tools** — jangan debat tab vs spasi, serahkan ke formatter
- Linting menangkap **bug potensial**, bukan hanya style

### [KUSTOMISASI] Tools per Ekosistem
```
JavaScript/TypeScript:
  Formatter → Prettier
  Linter    → ESLint + plugin yang relevan
  Type      → TypeScript strict mode

Python:
  Formatter → Black, isort (untuk import)
  Linter    → Ruff (menggantikan flake8 + pylint), mypy (type check)

Go:
  Formatter → gofmt (built-in, tidak perlu dikonfigurasi)
  Linter    → golangci-lint

Java:
  Formatter → google-java-format, Spotless
  Linter    → Checkstyle, SpotBugs, PMD

Rust:
  Formatter → rustfmt (built-in)
  Linter    → clippy (built-in)
```

### Checklist Linting
- [ ] Linter dikonfigurasi dan berjalan di CI
- [ ] Formatter dikonfigurasi — tidak ada debat style manual
- [ ] Pre-commit hook menjalankan lint + format otomatis
- [ ] Tidak ada `// eslint-disable` atau `# noqa` tanpa komentar penjelasan
- [ ] TypeScript / type checker dijalankan di CI (jika berlaku)

---

## 8. 📋 Code Quality Checklist (Pre-Review)

### Naming & Structure
- [ ] Semua nama variabel, fungsi, kelas deskriptif dan konsisten
- [ ] Tidak ada magic number atau magic string
- [ ] Kode diorganisir dengan baik (by feature atau by layer, konsisten)

### Functions & Logic
- [ ] Setiap fungsi punya satu tanggung jawab
- [ ] Tidak ada fungsi yang terlalu panjang (> 50 baris perlu dievaluasi)
- [ ] Level abstraksi konsisten dalam satu fungsi

### Comments & Docs
- [ ] Komentar menjelaskan WHY, bukan WHAT
- [ ] Public API punya docstring
- [ ] Tidak ada TODO yang terbengkalai tanpa konteks

### Tooling
- [ ] Linter berjalan tanpa error
- [ ] Formatter sudah dijalankan
- [ ] Type checker lulus (jika berlaku)

---

## 📚 Referensi

- [Clean Code — Robert C. Martin](https://www.oreilly.com/library/view/clean-code-a/9780136083238/)
- [Refactoring — Martin Fowler](https://refactoring.com/)
- [The Pragmatic Programmer](https://pragprog.com/titles/tpp20/the-pragmatic-programmer-20th-anniversary-edition/)
- [Google Style Guides](https://google.github.io/styleguide/)

---

*Code Quality Harness v1.0 — bagian dari [Harness Engineering](../README.md)*
