# 🧪 Testing Harness

> **Tujuan:** Memastikan setiap sistem dibangun dengan strategi testing yang solid, sehingga perubahan bisa dilakukan dengan percaya diri.
> **Gunakan harness ini saat:** memulai project baru, menulis test untuk fitur baru, melakukan test review, atau memperbaiki test yang rapuh.

---

## 🧭 Konteks untuk AI

```
Kamu adalah quality-focused engineer. Setiap kode yang ditulis harus bisa diverifikasi.
Prinsipmu:
1. Test adalah dokumentasi yang bisa dieksekusi — tulis test yang mudah dibaca
2. Test dulu, kode kemudian (TDD) jika memungkinkan — atau minimal test bersamaan
3. Prioritaskan test yang memberikan nilai tertinggi: integration > unit untuk coverage bisnis
4. Test yang lambat, flaky, atau sulit dipahami lebih berbahaya daripada tidak ada test
5. Pikirkan "apa yang bisa salah?" sebelum menulis test happy path

Sebelum menulis test, SELALU tanya: "Test ini melindungi dari bug apa?"
```

---

## 1. 🏗️ Strategi Testing (Test Pyramid)

```
         ╱‾‾‾‾‾‾‾‾‾‾‾‾╲
        ╱   E2E Tests   ╲       ← Sedikit, lambat, tapi verifikasi full flow
       ╱‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾╲
      ╱  Integration Tests  ╲   ← Sedang, verifikasi antar komponen
     ╱‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾╲
    ╱       Unit Tests          ╲  ← Banyak, cepat, verifikasi logika terisolasi
   ╱‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾╲
```

### Rasio yang Direkomendasikan
| Tipe | Porsi | Kecepatan | Kepercayaan |
|---|---|---|---|
| Unit | ~70% | < 1ms/test | Logika terisolasi |
| Integration | ~20% | 10ms–1s/test | Antar komponen |
| E2E | ~10% | 1s–30s/test | Full user flow |

---

## 2. 🔬 Unit Testing

### Prinsip
- Satu test = satu behaviour yang diverifikasi
- Test harus **FIRST**: Fast, Isolated, Repeatable, Self-validating, Timely
- Gunakan **AAA pattern**: Arrange → Act → Assert
- Mock dependency eksternal (DB, API, filesystem)

### Struktur Test yang Baik

```
[KUSTOMISASI] Sesuaikan dengan framework test stack-mu:
- JavaScript/TypeScript → Jest, Vitest
- Python               → pytest, unittest
- Go                   → testing (std), testify
- Java                 → JUnit 5, Mockito
- Rust                 → built-in #[test]
- C#                   → xUnit, NUnit
```

```
// Pola AAA — berlaku di semua bahasa:

Test: "harus mengembalikan total harga dengan diskon"

Arrange: Siapkan data — produk dengan harga 100, diskon 20%
Act:     Panggil fungsi — hitungTotal(produk, diskon)
Assert:  Verifikasi hasil — hasilnya harus 80
```

### Checklist Unit Test
- [ ] Setiap fungsi/method punya test untuk **happy path**
- [ ] Test mencakup **edge case**: input kosong, null, nilai ekstrem
- [ ] Test mencakup **error case**: exception, invalid input
- [ ] Nama test deskriptif: `"harus [expected] ketika [kondisi]"`
- [ ] Tidak ada logika kondisional (if/switch) di dalam test
- [ ] Test tidak bergantung satu sama lain (urutan tidak berpengaruh)
- [ ] Tidak ada `sleep()` atau wait arbitrary di unit test

### Anti-Pattern yang Harus Dihindari
```
❌ Test yang menguji implementasi, bukan behaviour
❌ Test yang bergantung pada urutan eksekusi
❌ Satu test yang memverifikasi banyak hal sekaligus
❌ Mock yang terlalu banyak — tanda bahwa kode terlalu coupled
❌ Test name yang tidak deskriptif: test1(), testA()
```

---

## 3. 🔗 Integration Testing

### Prinsip
- Uji **interaksi antar komponen nyata** — database, queue, external service
- Gunakan **test database** atau container (Docker) — jangan pakai production DB
- Setiap test harus bisa **reset state** — gunakan transaction rollback atau seed ulang

### Apa yang Harus Diuji
- [ ] Operasi database (CRUD) benar-benar tersimpan dan terbaca
- [ ] API endpoint mengembalikan response yang sesuai kontrak
- [ ] Message queue menerima dan memproses pesan dengan benar
- [ ] Integrasi dengan service eksternal (gunakan mock server / test double)
- [ ] Auth flow end-to-end (login → token → akses resource)

### [KUSTOMISASI] Tools Integration Test
```
Database testing:
- Node.js  → Supertest + test DB, Testcontainers
- Python   → pytest + SQLAlchemy + test DB
- Go       → net/http/httptest + testcontainers-go
- Java     → Spring Boot Test + Testcontainers

Mock server (untuk external API):
- WireMock (multi-language)
- MSW - Mock Service Worker (JS)
- responses (Python)
- httptest (Go)
```

### Checklist Integration Test
- [ ] Test database dibersihkan sebelum/sesudah setiap test suite
- [ ] Tidak ada test yang menggunakan data production
- [ ] Semua dependency eksternal di-mock / di-stub
- [ ] Test berjalan di lingkungan CI yang identik dengan lokal

---

## 4. 🌐 End-to-End (E2E) Testing

### Prinsip
- Uji dari perspektif **user yang sesungguhnya**
- Fokus pada **critical user journey** — jangan cover semua skenario
- E2E test harus **stabil** — flaky E2E test lebih merusak daripada tidak ada

### Critical User Journey yang Wajib Di-cover
- [ ] Alur **autentikasi** (register, login, logout)
- [ ] Alur **core business** (transaksi utama / fitur paling penting)
- [ ] Alur **error recovery** (apa yang terjadi jika gagal)

### [KUSTOMISASI] Tools E2E
```
Web (Browser):
- Playwright  → cross-browser, modern, recommended
- Cypress     → developer experience bagus, Chrome-focused
- Puppeteer   → headless Chrome, low-level

API E2E:
- Postman / Newman (CLI runner)
- k6 (juga bisa untuk load test)
- REST-assured (Java)

Mobile:
- Detox (React Native)
- Espresso (Android)
- XCTest (iOS)
```

### Checklist E2E
- [ ] Hanya cover **critical path** — bukan semua fitur
- [ ] Test berjalan di lingkungan staging yang identik dengan production
- [ ] Tidak ada dependency ke data yang bisa berubah (gunakan seed data)
- [ ] Screenshot / video direkam saat test gagal
- [ ] Timeout dikonfigurasi dengan wajar (jangan terlalu pendek atau panjang)

---

## 5. 🔄 Test-Driven Development (TDD)

### Siklus TDD (Red → Green → Refactor)

```
1. RED    → Tulis test yang GAGAL untuk behaviour yang belum ada
2. GREEN  → Tulis kode MINIMUM agar test lulus (boleh jelek dulu)
3. REFACTOR → Perbaiki kode tanpa mengubah behaviour (test tetap hijau)

Ulangi untuk setiap behaviour baru.
```

### Kapan TDD Paling Efektif
```
✅ Business logic yang kompleks
✅ Algoritma dan kalkulasi
✅ State machine / workflow
✅ Parsing dan transformasi data
❌ UI/layout (lebih baik visual testing)
❌ Integrasi dengan third-party yang belum stabil
❌ Prototipe yang masih belum jelas requirement-nya
```

---

## 6. 🎯 Test Coverage

### Target Coverage
```
Bukan tentang angka — tapi tentang confidence:

Minimum yang bermakna: 70% line coverage
Target yang baik:      80% line coverage
Overkill jika dipaksakan: 100% coverage

PENTING: 100% coverage BUKAN berarti bug-free.
Coverage mengukur "kode yang dieksekusi saat test", bukan "behaviour yang benar diuji".
```

### Prioritas Coverage
1. **Business logic / domain layer** → 90%+ wajib
2. **API handler / controller** → 80%+ disarankan
3. **Utility functions** → 90%+ (mudah di-test)
4. **Infrastructure layer** → integration test, bukan unit
5. **UI components** → snapshot test + interaction test

### Checklist Coverage
- [ ] Coverage report dijalankan di CI setiap PR
- [ ] Coverage tidak boleh turun dari threshold yang ditentukan
- [ ] Fokus pada **branch coverage**, bukan hanya line coverage
- [ ] Identifikasi uncovered code dan evaluasi apakah perlu di-test

---

## 7. 🚥 Testing di CI/CD Pipeline

### Gate per Tahap

```
Pre-commit (lokal):
  └── Lint + format check
  └── Unit test (hanya yang terpengaruh perubahan)

Pull Request / Merge Request:
  └── Unit test (semua)
  └── Integration test
  └── Coverage check (tidak boleh turun)
  └── Static analysis / SAST

Staging Deploy:
  └── E2E test (critical path)
  └── Performance test (smoke)

Production Deploy:
  └── Smoke test (minimal, cepat)
  └── Monitoring alert siap
```

### Checklist CI Testing
- [ ] Test berjalan otomatis di setiap PR
- [ ] PR tidak bisa di-merge jika test gagal
- [ ] Coverage threshold dikonfigurasi dan di-enforce
- [ ] Test result / report bisa dilihat tanpa harus pull dan run lokal
- [ ] Flaky test di-track dan diprioritaskan untuk diperbaiki

---

## 8. 🩺 Test Quality & Maintenance

### Tanda Test yang Perlu Diperbaiki
```
🚩 Flaky — kadang lulus kadang gagal tanpa perubahan kode
🚩 Slow  — unit test > 100ms, integration test > 5s
🚩 Brittle — gagal karena perubahan kecil yang tidak relevan
🚩 Unclear — tidak jelas apa yang sedang diuji dari namanya
🚩 Duplicate — beberapa test menguji hal yang persis sama
```

### Checklist Kualitas Test
- [ ] Flaky test langsung di-quarantine dan diperbaiki
- [ ] Test suite selesai dalam waktu yang wajar (unit < 1 menit, semua < 15 menit)
- [ ] Test lama di-review dan di-optimasi
- [ ] Dead test (tidak pernah gagal, tidak pernah melindungi apa-apa) di-hapus

---

## 9. 📋 Testing Checklist (Pre-Launch)

### Unit & Integration
- [ ] Business logic kritis punya unit test ≥ 90% coverage
- [ ] API endpoint punya integration test
- [ ] Error path dan edge case ter-cover
- [ ] Semua test hijau di CI

### E2E
- [ ] Critical user journey ter-cover
- [ ] E2E berjalan di staging environment
- [ ] Tidak ada flaky test di suite E2E

### Non-Functional
- [ ] Performance test dijalankan (minimal smoke)
- [ ] Security test (SAST) lulus
- [ ] Dependency vulnerability scan lulus

---

## 📚 Referensi

- [The Practical Test Pyramid](https://martinfowler.com/articles/practical-test-pyramid.html)
- [Test-Driven Development by Example — Kent Beck](https://www.oreilly.com/library/view/test-driven-development/0321146530/)
- [Growing Object-Oriented Software, Guided by Tests](http://www.growing-object-oriented-software.com/)
- [Testing Trophy (Kent C. Dodds)](https://kentcdodds.com/blog/the-testing-trophy-and-testing-classifications)

---

*Testing Harness v1.0 — bagian dari [Harness Engineering](../README.md)*
