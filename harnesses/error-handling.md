# 🚨 Error Handling Harness

> **Tujuan:** Memastikan setiap sistem menangani kegagalan dengan graceful, informatif, dan konsisten — bukan hanya mengoptimalkan happy path.
> **Gunakan harness ini saat:** menulis kode baru, review error handling, mendesain error taxonomy, atau debugging sistem yang tidak informatif saat gagal.

---

## 🧭 Konteks untuk AI

```
Kamu adalah engineer yang memperlakukan error sebagai first-class citizen.
Sistem yang baik bukan sistem yang tidak pernah gagal — tapi sistem yang gagal dengan elegan.

Prinsipmu:
1. Setiap operasi yang bisa gagal HARUS ditangani — tidak ada yang diabaikan diam-diam
2. Error message untuk USER harus human-readable dan actionable
3. Error message untuk DEVELOPER harus technical dan traceable
4. Fail fast, fail loudly — lebih baik crash dengan jelas daripada silent bug
5. Jangan swallow error — catch dan dismiss tanpa logging adalah dosa terbesar

Setiap kali menulis kode, tanya: "Apa yang terjadi jika ini gagal? Apakah kita tahu kenapa?"
```

---

## 1. 🗂️ Error Taxonomy

### Kategorisasi Error

Setiap error harus dikategorikan berdasarkan sifatnya:

```
OPERATIONAL ERROR (Expected — bisa terjadi di production yang sehat):
  → Validasi input user yang salah
  → Resource tidak ditemukan (404)
  → Timeout dari external service
  → Rate limit tercapai
  → Network failure sementara
  Handling: tangani dengan graceful, informasikan ke user

PROGRAMMER ERROR (Bug — tidak boleh terjadi di production):
  → NullPointerException yang tidak diantisipasi
  → Fungsi dipanggil dengan tipe argumen yang salah
  → Logic error / kalkulasi yang salah
  → Array index out of bounds
  Handling: crash cepat, log dengan detail, alert tim, perbaiki kode

EXTERNAL ERROR (Di luar kendali sistem):
  → Third-party API down
  → Database tidak bisa dicapai
  → Disk penuh
  → Memory exhausted
  Handling: retry strategy, fallback, circuit breaker, alert ops
```

### Error Code Structure

Definisikan error code yang konsisten di seluruh sistem:

```
Format yang direkomendasikan:
  <DOMAIN>_<ENTITY>_<REASON>

Contoh:
  AUTH_TOKEN_EXPIRED
  AUTH_TOKEN_INVALID
  USER_NOT_FOUND
  USER_EMAIL_DUPLICATE
  ORDER_PAYMENT_FAILED
  ORDER_ITEM_OUT_OF_STOCK
  VALIDATION_REQUIRED_FIELD_MISSING
  VALIDATION_FORMAT_INVALID
  RATE_LIMIT_EXCEEDED
  EXTERNAL_PAYMENT_GATEWAY_TIMEOUT
```

---

## 2. 🏗️ Error Design Patterns

### Custom Error / Exception Classes

```
[KUSTOMISASI] Sesuaikan dengan bahasa yang digunakan:

Prinsip universal:
- Buat hierarki error yang jelas
- Sertakan error code, message, dan context
- Bedakan operational vs programmer error
```

```javascript
// Contoh — TypeScript/JavaScript:

// Base class
class AppError extends Error {
  constructor(
    public readonly code: string,
    public readonly message: string,
    public readonly statusCode: number,
    public readonly context?: Record<string, unknown>
  ) {
    super(message);
    this.name = this.constructor.name;
  }
}

// Operational errors — bisa terjadi di production normal
class NotFoundError extends AppError {
  constructor(resource: string, id: string) {
    super('RESOURCE_NOT_FOUND', `${resource} with id ${id} not found`, 404, { resource, id });
  }
}

class ValidationError extends AppError {
  constructor(fields: Array<{ field: string; message: string }>) {
    super('VALIDATION_ERROR', 'Input validation failed', 422, { fields });
  }
}

class UnauthorizedError extends AppError {
  constructor(reason: string) {
    super('UNAUTHORIZED', reason, 401);
  }
}

// External errors — gagal karena dependency eksternal
class ExternalServiceError extends AppError {
  constructor(service: string, originalError: Error) {
    super('EXTERNAL_SERVICE_ERROR', `${service} is unavailable`, 503, {
      service,
      originalMessage: originalError.message,
    });
  }
}
```

```python
# Contoh — Python:

class AppError(Exception):
    def __init__(self, code: str, message: str, status_code: int, context: dict = None):
        super().__init__(message)
        self.code = code
        self.status_code = status_code
        self.context = context or {}

class NotFoundError(AppError):
    def __init__(self, resource: str, id: str):
        super().__init__('RESOURCE_NOT_FOUND', f'{resource} {id} not found', 404,
                         {'resource': resource, 'id': id})

class ValidationError(AppError):
    def __init__(self, fields: list):
        super().__init__('VALIDATION_ERROR', 'Input validation failed', 422,
                         {'fields': fields})
```

```go
// Contoh — Go:

type AppError struct {
    Code       string
    Message    string
    StatusCode int
    Context    map[string]any
    Err        error // wrapped original error
}

func (e *AppError) Error() string { return e.Message }
func (e *AppError) Unwrap() error { return e.Err }

func NewNotFoundError(resource, id string) *AppError {
    return &AppError{
        Code:       "RESOURCE_NOT_FOUND",
        Message:    fmt.Sprintf("%s with id %s not found", resource, id),
        StatusCode: 404,
        Context:    map[string]any{"resource": resource, "id": id},
    }
}
```

---

## 3. 🔄 Error Propagation

### Prinsip Propagasi

```
LAYER BAWAH (Repository / Infrastructure):
  → Tangkap error low-level (DB error, network error)
  → Wrap menjadi domain error yang meaningful
  → Sertakan context yang relevan
  → JANGAN expose detail implementasi ke layer atas

LAYER TENGAH (Service / Use Case):
  → Handle operational error yang bisa diselesaikan di sini
  → Propagate error yang harus ditangani layer atas
  → Tambahkan context jika diperlukan

LAYER ATAS (Controller / Handler):
  → Handle semua error yang tidak tertangkap
  → Map error ke HTTP status code / response format
  → Log error dengan severity yang tepat
  → Kembalikan response yang user-friendly
```

### Pola Wrapping Error

```javascript
// ❌ SALAH — error detail implementasi bocor ke atas
async function getUserById(id) {
  const result = await db.query('SELECT * FROM users WHERE id = $1', [id]);
  // Jika gagal, PostgreSQL error langsung muncul ke controller
  return result.rows[0];
}

// ✅ BENAR — wrap error menjadi domain error
async function getUserById(id) {
  try {
    const result = await db.query('SELECT * FROM users WHERE id = $1', [id]);
    if (!result.rows[0]) {
      throw new NotFoundError('User', id);
    }
    return result.rows[0];
  } catch (error) {
    if (error instanceof NotFoundError) throw error; // re-throw domain error
    throw new ExternalServiceError('Database', error); // wrap low-level error
  }
}
```

### Anti-Pattern yang Harus Dihindari

```javascript
// ❌ Silent error — paling berbahaya
try {
  await doSomething();
} catch (error) {
  // tidak ada apa-apa di sini
}

// ❌ Log tapi tidak propagate — error hilang
try {
  await doSomething();
} catch (error) {
  console.log(error); // lanjut seolah tidak ada masalah
}

// ❌ Catch terlalu broad — sembunyikan programmer error
try {
  processData(input);
} catch (error) {
  return null; // menyembunyikan bug!
}

// ✅ Handle yang spesifik, propagate yang lain
try {
  await doSomething();
} catch (error) {
  if (error instanceof NetworkError) {
    return await fallbackStrategy();
  }
  throw error; // programmer error dan unknown error harus propagate
}
```

---

## 4. 📨 Error Response Format

### Format Response Error (API)

Gunakan format yang konsisten di seluruh API (lihat juga [api-design.md](./api-design.md)):

```json
// Single error
{
  "error": {
    "code": "USER_NOT_FOUND",
    "message": "User dengan ID ini tidak ditemukan.",
    "requestId": "req_abc123xyz"
  }
}

// Validation error dengan multiple fields
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Data yang dikirim tidak valid.",
    "details": [
      {
        "field": "email",
        "code": "FORMAT_INVALID",
        "message": "Format email tidak valid."
      },
      {
        "field": "password",
        "code": "TOO_SHORT",
        "message": "Password minimal 8 karakter."
      }
    ],
    "requestId": "req_abc123xyz"
  }
}
```

### Pesan Error yang Baik untuk User

```
Pesan error yang baik memenuhi 3 kriteria:

1. JELAS — User mengerti apa yang salah
   ❌ "Error 422"
   ✅ "Email yang kamu masukkan tidak valid."

2. ACTIONABLE — User tahu apa yang harus dilakukan
   ❌ "Terjadi kesalahan"
   ✅ "Sesi kamu telah berakhir. Silakan login kembali."

3. TIDAK MENAKUTKAN — Tidak expose technical detail
   ❌ "NullPointerException at line 47 in UserService.java"
   ✅ "Sepertinya ada masalah di sistem kami. Tim kami sudah diberitahu."
```

### Status Code yang Tepat

```
Gunakan status code yang semantically correct:

Validasi input gagal          → 422 Unprocessable Entity
Resource tidak ditemukan      → 404 Not Found
Tidak ter-autentikasi         → 401 Unauthorized
Tidak punya akses             → 403 Forbidden
Konflik data (email duplikat) → 409 Conflict
Rate limit                    → 429 Too Many Requests
Error server yang tidak terduga → 500 Internal Server Error
Service dependency down       → 503 Service Unavailable
```

---

## 5. 📝 Logging Error

### Prinsip Logging Error

```
Log HARUS menjawab pertanyaan: "Apa yang terjadi, kapan, di mana, dan kenapa?"

Setiap log error harus berisi:
  - Timestamp (ISO 8601)
  - Severity level (ERROR, WARN, INFO)
  - Error code dan message
  - Stack trace (untuk error yang tidak diantisipasi)
  - Request context (requestId, userId, endpoint)
  - Input / context yang menyebabkan error (tanpa data sensitif)
```

### Severity yang Tepat

```
ERROR  → Sesuatu yang seharusnya tidak terjadi, butuh perhatian
         Contoh: DB connection failed, unhandled exception

WARN   → Sesuatu yang tidak ideal tapi sistem masih berjalan
         Contoh: Retry ke-3 berhasil, response lambat, config deprecated

INFO   → Event penting yang normal
         Contoh: User login berhasil, payment processed

DEBUG  → Detail untuk debugging (nonaktif di production)
         Contoh: SQL query, internal state, function call
```

### Structured Logging

```json
// ✅ Structured log (mudah di-parse dan di-query):
{
  "timestamp": "2026-06-15T09:00:00.000Z",
  "level": "ERROR",
  "code": "EXTERNAL_PAYMENT_GATEWAY_TIMEOUT",
  "message": "Payment gateway tidak merespons dalam 30 detik",
  "requestId": "req_abc123",
  "userId": "usr_456",
  "endpoint": "POST /api/payments",
  "context": {
    "orderId": "ord_789",
    "gatewayName": "stripe",
    "timeoutMs": 30000
  },
  "stack": "TimeoutError: ..."
}

// ❌ Unstructured log (susah di-search dan di-parse):
"ERROR: Payment failed for user 456 at 09:00 because stripe timeout"
```

### Yang TIDAK Boleh Di-log

```
❌ Password atau password hash
❌ Token (JWT, API key, refresh token)
❌ Nomor kartu kredit
❌ Data PII lengkap (NIK, nomor rekening) — log hanya sebagian jika perlu
❌ Secret / credential apapun
```

---

## 6. 🔁 Retry & Fallback Strategy

### Kapan Retry

```
✅ Retry untuk TRANSIENT errors (sementara, bisa membaik sendiri):
   - Network timeout
   - Rate limit (dengan backoff)
   - Service temporary unavailable (503)

❌ JANGAN retry untuk PERMANENT errors:
   - 400 Bad Request (validasi gagal — retry tidak akan berhasil)
   - 401 Unauthorized (retry tidak menyelesaikan masalah auth)
   - 404 Not Found (resource memang tidak ada)
   - 422 Validation Error
```

### Exponential Backoff

```
Rumus: delay = baseDelay * (2 ^ attemptNumber) + jitter

Contoh dengan base 1 detik:
  Attempt 1: tunggu 1s  + jitter
  Attempt 2: tunggu 2s  + jitter
  Attempt 3: tunggu 4s  + jitter
  Attempt 4: tunggu 8s  + jitter
  Max delay: 30s (cap agar tidak terlalu lama)

Jitter (randomness) penting untuk menghindari thundering herd:
  → Semua client retry pada waktu yang sama = DDoS ke server sendiri
```

### Fallback Strategy

```
Ketika semua retry gagal, pertimbangkan fallback:

Stale cache    → Kembalikan data cache lama (meski outdated)
Default value  → Kembalikan nilai default yang aman
Degraded mode  → Jalankan fitur dengan fungsionalitas terbatas
Queue          → Simpan ke queue, proses nanti saat service pulih
Error page     → Informasikan user dengan ramah bahwa fitur sementara tidak tersedia
```

---

## 7. 🌐 Error Handling per Layer

### Controller / Route Handler

```
Tanggung jawab:
✅ Catch semua error yang tidak tertangkap
✅ Map error ke HTTP response yang tepat
✅ Log error dengan context request
✅ Kembalikan format response yang konsisten
❌ JANGAN lakukan bisnis logic di sini

Implementasi (gunakan global error handler / middleware):
- Express.js  → error middleware (4 parameter)
- FastAPI     → exception handler
- Spring Boot → @ControllerAdvice
- Go gin      → custom Recovery middleware
```

### Service / Use Case Layer

```
Tanggung jawab:
✅ Handle operational error yang bisa diselesaikan (retry, fallback)
✅ Throw domain error yang meaningful
✅ Validasi business rule dan throw jika dilanggar
❌ JANGAN handle HTTP concerns (status code, response format)
```

### Repository / Infrastructure Layer

```
Tanggung jawab:
✅ Wrap low-level error menjadi domain error
✅ Sertakan context yang relevan saat wrapping
❌ JANGAN expose raw DB/network error ke layer atas
❌ JANGAN handle business logic
```

---

## 8. 📋 Error Handling Checklist (Pre-Launch)

### Design
- [ ] Error taxonomy terdefinisi (error code, kategori, severity)
- [ ] Custom error class / exception dibuat untuk domain errors
- [ ] Error propagation antar layer terdefinisi dan konsisten
- [ ] Format error response API konsisten di seluruh endpoint

### Implementation
- [ ] Tidak ada silent error (`catch` kosong atau hanya `console.log`)
- [ ] Semua external call punya timeout yang dikonfigurasi
- [ ] Retry dengan exponential backoff ada untuk transient error
- [ ] Fallback strategy ada untuk dependency kritis
- [ ] Error tidak mengekspos stack trace atau implementasi detail ke user

### Logging
- [ ] Semua error di-log dengan context yang cukup (requestId, userId)
- [ ] Severity level digunakan dengan tepat
- [ ] Tidak ada data sensitif di log
- [ ] Log dalam format terstruktur (JSON)

### Monitoring
- [ ] Alert dikonfigurasi untuk error rate yang tinggi
- [ ] Alert dikonfigurasi untuk error type kritis (DB down, payment gagal)
- [ ] Dashboard error tersedia untuk tim

---

## 📚 Referensi

- [Error Handling in Node.js — Joyent](https://web.archive.org/web/20230000000000*/https://www.joyent.com/node-js/production/design/errors)
- [Designing Error Messages — Nielsen Norman Group](https://www.nngroup.com/articles/error-message-guidelines/)
- [Exponential Backoff — AWS](https://docs.aws.amazon.com/general/latest/gr/api-retries.html)
- [Production Best Practices — Express.js](https://expressjs.com/en/advanced/best-practice-performance.html)

---

*Error Handling Harness v1.0 — bagian dari [Harness Engineering](../README.md)*
