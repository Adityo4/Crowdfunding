# 🔌 API Design Harness

> **Tujuan:** Memastikan API yang dibangun konsisten, intuitif, dan mudah digunakan oleh consumer.
> **Gunakan harness ini saat:** mendesain API baru, melakukan API review, mendefinisikan kontrak API, atau migrasi API.

---

## 🧭 Konteks untuk AI

```
Kamu adalah API-first engineer yang berpengalaman. API adalah produk — developer adalah user-nya.
Prinsipmu:
1. API-first — desain dan dokumentasikan kontrak sebelum implementasi
2. Consistency over cleverness — naming dan struktur yang konsisten lebih penting dari yang "pintar"
3. Design for the consumer — pikirkan bagaimana API ini akan dipakai, bukan hanya diimplementasikan
4. Backward compatibility is a promise — breaking change = versi baru
5. Error message adalah dokumentasi — pesan error harus actionable

Sebelum desain API, SELALU tanya: "Siapa yang akan memanggil API ini dan untuk apa?"
```

---

## 1. 🏗️ Pemilihan Paradigma API

### Kapan Pakai REST

```
✅ Resource-based operations (CRUD)
✅ Public API yang harus mudah dipahami
✅ Caching penting (HTTP caching native)
✅ Tim consumer tidak diketahui (public API)
✅ Simple request-response pattern
```

### Kapan Pakai GraphQL

```
✅ Data requirements bervariasi antar client (mobile vs web)
✅ Menghindari over-fetching dan under-fetching
✅ Banyak relasi antar data yang perlu di-query fleksibel
✅ Consumer diketahui dan bisa berkolaborasi (BFF pattern)
❌ Caching lebih kompleks daripada REST
❌ File upload lebih rumit
```

### Kapan Pakai gRPC

```
✅ Internal service-to-service (bukan public API)
✅ Performa sangat kritis (binary protocol, HTTP/2)
✅ Streaming dua arah dibutuhkan
✅ Strong typing penting (protobuf)
❌ Tidak cocok untuk browser-facing API langsung
```

### Kapan Pakai WebSocket / SSE

```
WebSocket:
✅ Real-time bidirectional (chat, collaborative editing, game)

SSE (Server-Sent Events):
✅ Real-time server-to-client satu arah (live feed, notifikasi)
✅ Lebih simpel dari WebSocket jika hanya butuh server push
```

---

## 2. 🌐 REST API Design

### Naming Convention

```
Resource:  Gunakan NOUN jamak, bukan verb
  ✅ /users          ❌ /getUser
  ✅ /orders         ❌ /createOrder
  ✅ /products       ❌ /productList

Hirarki:   Representasikan relasi dengan nesting (max 2 level)
  ✅ /users/{id}/orders        ← order milik user tertentu
  ✅ /orders/{id}/items        ← item dalam order
  ❌ /users/{id}/orders/{id}/items/{id}   ← terlalu dalam

Action (verb):  Jika benar-benar tidak bisa di-resource-kan
  ✅ /orders/{id}/cancel       ← action spesifik
  ✅ /auth/login
  ✅ /payments/{id}/refund
```

### HTTP Method yang Tepat

| Method | Tujuan | Idempotent | Body |
|---|---|---|---|
| **GET** | Ambil data | ✅ | ❌ |
| **POST** | Buat resource baru | ❌ | ✅ |
| **PUT** | Ganti seluruh resource | ✅ | ✅ |
| **PATCH** | Update sebagian resource | ❌ | ✅ |
| **DELETE** | Hapus resource | ✅ | ❌ |

### HTTP Status Code yang Tepat

```
2xx — Sukses
  200 OK             → GET, PUT, PATCH berhasil
  201 Created        → POST berhasil membuat resource baru (sertakan Location header)
  204 No Content     → DELETE berhasil, tidak ada response body

3xx — Redirect
  301 Moved Permanently → URL resource berubah permanen
  304 Not Modified      → Client bisa pakai cache

4xx — Client Error
  400 Bad Request    → Input tidak valid
  401 Unauthorized   → Tidak ada / invalid authentication
  403 Forbidden      → Authenticated tapi tidak punya akses
  404 Not Found      → Resource tidak ditemukan
  409 Conflict       → Konflik state (e.g. email sudah ada)
  422 Unprocessable  → Validasi bisnis gagal (berbeda dari 400)
  429 Too Many Req.  → Rate limit tercapai

5xx — Server Error
  500 Internal Error → Unexpected server error
  502 Bad Gateway    → Upstream service error
  503 Unavailable    → Service down / maintenance
  504 Gateway Timeout → Upstream timeout
```

### Response Format yang Konsisten

```json
// ✅ Success Response
{
  "data": {
    "id": "usr_123",
    "name": "Budi Santoso",
    "email": "budi@example.com",
    "createdAt": "2026-06-15T09:00:00Z"
  },
  "meta": {
    "requestId": "req_abc123"
  }
}

// ✅ Collection Response (dengan pagination)
{
  "data": [ ... ],
  "pagination": {
    "page": 1,
    "perPage": 20,
    "total": 150,
    "totalPages": 8,
    "hasNext": true,
    "hasPrev": false
  },
  "meta": {
    "requestId": "req_abc123"
  }
}

// ✅ Error Response
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Request validation failed",
    "details": [
      {
        "field": "email",
        "message": "Email tidak valid"
      },
      {
        "field": "password",
        "message": "Password minimal 8 karakter"
      }
    ]
  },
  "meta": {
    "requestId": "req_abc123"
  }
}
```

### Checklist REST Design
- [ ] Nama resource menggunakan noun jamak
- [ ] HTTP method digunakan sesuai semantiknya
- [ ] Status code tepat dan konsisten
- [ ] Response format konsisten di seluruh API
- [ ] Error response informatif dan actionable
- [ ] `requestId` / `traceId` ada di setiap response (untuk debugging)

---

## 3. 📄 Pagination, Filtering & Sorting

### Strategi Pagination

```
Offset Pagination (simpel, tapi ada masalah di data yang berubah):
  GET /users?page=2&perPage=20

Cursor Pagination (recommended untuk data besar atau real-time):
  GET /users?cursor=eyJpZCI6MTIzfQ==&limit=20
  Response: { "data": [...], "nextCursor": "...", "hasMore": true }

Keyed Pagination (untuk sort berdasarkan unique field):
  GET /users?afterId=123&limit=20
```

### Filtering & Sorting

```
Filtering:
  GET /products?category=electronics&minPrice=100&maxPrice=500
  GET /orders?status=pending&createdAfter=2026-01-01

Sorting:
  GET /products?sort=price&order=asc
  GET /users?sort=-createdAt        ← prefix minus untuk descending

Search:
  GET /products?q=laptop
```

### Field Selection (untuk GraphQL-like di REST)

```
Hanya kirim field yang dibutuhkan:
  GET /users?fields=id,name,email

Nested fields:
  GET /orders?include=user,items
```

---

## 4. 🔖 API Versioning

### Strategi Versioning

```
URL Path (paling eksplisit, recommended untuk public API):
  /v1/users
  /v2/users

Header (lebih clean URL, tapi kurang visible):
  Accept: application/vnd.myapp.v2+json

Query Parameter (kurang direkomendasikan):
  /users?version=2
```

### Kapan Buat Versi Baru

```
Breaking Changes (WAJIB versi baru):
  ❌ Hapus field dari response
  ❌ Ubah tipe data field yang ada
  ❌ Ubah nama field
  ❌ Hapus endpoint
  ❌ Ubah behaviour yang sudah ada

Non-Breaking Changes (AMAN, tidak perlu versi baru):
  ✅ Tambah field baru di response
  ✅ Tambah endpoint baru
  ✅ Tambah query parameter opsional baru
  ✅ Tambah HTTP method baru di resource yang ada
```

### Deprecation Policy

```
1. Announce: Informasikan deprecation di dokumentasi dan response header
   Deprecation: Sat, 01 Jan 2027 00:00:00 GMT
   Sunset: Mon, 01 Jul 2027 00:00:00 GMT
   Link: <https://api.example.com/v2/users>; rel="successor-version"

2. Grace Period: Berikan minimal 6 bulan sebelum sunset
3. Sunset: Matikan versi lama sesuai jadwal yang diumumkan
```

---

## 5. 🔐 Keamanan API

### Autentikasi

```
Public API:
  → API Key di header: X-API-Key: key123
  → OAuth 2.0 untuk akses atas nama user

Internal API:
  → JWT (service-to-service)
  → mTLS untuk keamanan tinggi

User-facing API:
  → JWT Bearer token: Authorization: Bearer <token>
  → OAuth 2.0 + OIDC untuk SSO
```

### Checklist Keamanan API
- [ ] Semua endpoint di-autentikasi kecuali yang eksplisit publik
- [ ] Rate limiting aktif (lihat [security.md](./security.md) untuk detail)
- [ ] Input divalidasi di server-side
- [ ] Response tidak mengekspos data internal atau stack trace
- [ ] CORS dikonfigurasi dengan whitelist domain
- [ ] API key / token dirotasi secara berkala

---

## 6. 📚 Dokumentasi API

### Prinsip
- Dokumentasi adalah bagian dari API — bukan afterthought
- Dokumentasi harus bisa di-try langsung (interactive)
- Update dokumentasi bersamaan dengan update kode

### Standard Dokumentasi

```
OpenAPI / Swagger (REST):
  → Tulis spec di openapi.yaml / openapi.json
  → Generate dari kode (code-first) atau tulis manual (design-first)
  → Host dengan Swagger UI atau Redoc

GraphQL:
  → Schema self-documenting via introspection
  → Tambah deskripsi di setiap type dan field
  → Host GraphiQL atau GraphQL Playground

gRPC:
  → Protobuf schema sebagai dokumentasi
  → Tambahkan komentar di .proto files
  → Generate documentation dengan protoc-gen-doc
```

### Yang Harus Ada di Dokumentasi

```
Per Endpoint:
  ✅ Deskripsi singkat apa yang dilakukan
  ✅ Request: method, URL, header, path params, query params, body (dengan tipe)
  ✅ Response: semua kemungkinan status code dan contoh response
  ✅ Error codes dan artinya
  ✅ Contoh request dan response (curl / code snippet)
  ✅ Rate limit jika berlaku

Di Level API:
  ✅ Authentication guide
  ✅ Base URL per environment
  ✅ Changelog / version history
  ✅ Deprecation notices
```

### Checklist Dokumentasi
- [ ] Semua endpoint terdokumentasi (OpenAPI spec atau equivalent)
- [ ] Contoh request dan response tersedia
- [ ] Error codes terdokumentasi dengan penjelasan yang jelas
- [ ] Dokumentasi bisa diakses dan di-try secara interaktif
- [ ] Dokumentasi di-update bersamaan dengan setiap perubahan API
- [ ] Changelog terdokumentasi

---

## 7. 🧪 API Testing

### Level Testing API

```
Contract Test:  Verifikasi API sesuai kontrak (OpenAPI spec)
Unit Test:      Test logic di handler / resolver
Integration:    Test endpoint dengan database nyata
E2E:            Test full user flow via API
```

### Checklist API Testing
- [ ] **Contract test** memverifikasi response sesuai OpenAPI spec
- [ ] Semua endpoint punya integration test
- [ ] Happy path ter-cover
- [ ] Error path ter-cover (400, 401, 403, 404, 422, 500)
- [ ] Rate limit behavior di-test
- [ ] Pagination di-test (edge case: halaman pertama, terakhir, empty)

### [KUSTOMISASI] Tools API Testing
```
Contract Testing:
- Dredd        → test API vs OpenAPI spec
- Pact         → consumer-driven contract testing
- Schemathesis  → property-based testing dari OpenAPI spec

Integration & E2E:
- Postman / Newman (CLI)
- Insomnia
- REST-assured (Java)
- httpx + pytest (Python)
- Supertest (Node.js)
```

---

## 8. 📋 API Design Checklist (Pre-Launch)

### Design
- [ ] Resource naming konsisten (noun jamak, snake_case atau camelCase dipilih satu)
- [ ] HTTP method digunakan sesuai semantik
- [ ] Status code konsisten di seluruh API
- [ ] Response format seragam (termasuk error format)
- [ ] Versioning strategy terdefinisi

### Dokumentasi
- [ ] OpenAPI spec / GraphQL schema tersedia dan up-to-date
- [ ] Semua endpoint terdokumentasi dengan contoh
- [ ] Error codes terdokumentasi
- [ ] Authentication guide tersedia
- [ ] Changelog tersedia

### Keamanan & Stabilitas
- [ ] Semua endpoint ter-autentikasi (kecuali yang eksplisit publik)
- [ ] Rate limiting aktif
- [ ] Input validation di semua endpoint
- [ ] Tidak ada breaking change dari versi sebelumnya (atau sudah versi baru)

---

## 📚 Referensi

- [REST API Design Rulebook](https://www.oreilly.com/library/view/rest-api-design/9781449317904/)
- [OpenAPI Specification](https://spec.openapis.org/oas/latest.html)
- [GraphQL Best Practices](https://graphql.org/learn/best-practices/)
- [Microsoft REST API Guidelines](https://github.com/microsoft/api-guidelines)
- [Stripe API Design](https://stripe.com/docs/api) ← contoh REST API terbaik

---

*API Design Harness v1.0 — bagian dari [Harness Engineering](../README.md)*
