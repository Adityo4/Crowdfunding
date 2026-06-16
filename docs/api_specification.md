# API Specification - Crowdfunding Platform

Dokumen ini mendefinisikan spesifikasi API untuk platform Crowdfunding berbasis REST API dengan format JSON sesuai dengan pedoman `api-design.md`.

## Standard Response Format

Setiap response API menggunakan struktur JSON yang konsisten.

### Success Response
```json
{
  "data": { ... },
  "meta": {
    "requestId": "uuid-request-trace"
  }
}
```

### Collection Response (Paginated)
```json
{
  "data": [ ... ],
  "pagination": {
    "page": 1,
    "perPage": 10,
    "total": 100,
    "totalPages": 10,
    "hasNext": true,
    "hasPrev": false
  },
  "meta": {
    "requestId": "uuid-request-trace"
  }
}
```

### Error Response
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Error description message",
    "details": [
      {
        "field": "input_field_name",
        "message": "Validation message"
      }
    ]
  },
  "meta": {
    "requestId": "uuid-request-trace"
  }
}
```

---

## Endpoint List

### 1. Authentication & Users

#### POST `/v1/auth/register`
Mendaftarkan user baru.
- **Request Body:**
  ```json
  {
    "fullName": "John Doe",
    "email": "johndoe@example.com",
    "password": "Password123!",
    "phoneNumber": "08123456789"
  }
  ```
- **Response (201 Created):**
  ```json
  {
    "data": {
      "id": "user-uuid",
      "fullName": "John Doe",
      "email": "johndoe@example.com",
      "role": "user"
    }
  }
  ```

#### POST `/v1/auth/login`
Autentikasi user dan mendapatkan JWT token.
- **Request Body:**
  ```json
  {
    "email": "johndoe@example.com",
    "password": "Password123!"
  }
  ```
- **Response (200 OK):**
  ```json
  {
    "data": {
      "token": "jwt-bearer-token-string",
      "expiresAt": "2026-06-17T20:00:00Z"
    }
  }
  ```

#### GET `/v1/users/me`
Mendapatkan profil user yang sedang login.
- **Headers:** `Authorization: Bearer <token>`
- **Response (200 OK):**
  ```json
  {
    "data": {
      "id": "user-uuid",
      "fullName": "John Doe",
      "email": "johndoe@example.com",
      "role": "user",
      "phoneNumber": "08123456789",
      "createdAt": "2026-06-16T12:00:00Z"
    }
  }
  ```

---

### 2. Categories

#### GET `/v1/categories`
Mendapatkan semua kategori campaign.
- **Response (200 OK):**
  ```json
  {
    "data": [
      {
        "id": "category-uuid",
        "name": "Pendidikan",
        "slug": "pendidikan"
      }
    ]
  }
  ```

---

### 3. Charities (Campaigns)

#### GET `/v1/charities`
Mendapatkan daftar campaign (Mendukung filter, sort, dan pagination).
- **Query Params:**
  - `page` (int, default: 1)
  - `perPage` (int, default: 10)
  - `category` (string, slug kategori)
  - `status` (string: active, completed)
  - `sort` (string: newest, target_amount, urgent)
  - `q` (string, search query)
- **Response (200 OK):**
  ```json
  {
    "data": [
      {
        "id": "charity-uuid",
        "title": "Bantu Pembangunan Sekolah Roboh",
        "slug": "bantu-pembangunan-sekolah-roboh",
        "targetAmount": 150000000.00,
        "currentAmount": 25000000.00,
        "startDate": "2026-06-16",
        "endDate": "2026-07-16",
        "coverImageUrl": "https://storage.local/cover.jpg",
        "organizationName": "Yayasan Pendidikan Bersama",
        "status": "active",
        "isVerified": true,
        "createdAt": "2026-06-16T12:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "perPage": 10,
      "total": 1,
      "totalPages": 1,
      "hasNext": false,
      "hasPrev": false
    }
  }
  ```

#### GET `/v1/charities/{slug}`
Mendapatkan detail campaign berdasarkan slug.
- **Response (200 OK):**
  ```json
  {
    "data": {
      "id": "charity-uuid",
      "title": "Bantu Pembangunan Sekolah Roboh",
      "slug": "bantu-pembangunan-sekolah-roboh",
      "description": "Detail deskripsi lengkap penggalangan dana...",
      "targetAmount": 150000000.00,
      "currentAmount": 25000000.00,
      "startDate": "2026-06-16",
      "endDate": "2026-07-16",
      "coverImageUrl": "https://storage.local/cover.jpg",
      "additionalImages": [
        "https://storage.local/img1.jpg",
        "https://storage.local/img2.jpg"
      ],
      "contactPerson": "Budi Santoso",
      "contactEmail": "budi@schoolfoundation.org",
      "contactPhone": "08123456789",
      "organizationName": "Yayasan Pendidikan Bersama",
      "status": "active",
      "isVerified": true,
      "creator": {
        "id": "user-uuid",
        "fullName": "Budi Santoso"
      },
      "category": {
        "id": "category-uuid",
        "name": "Pendidikan"
      },
      "createdAt": "2026-06-16T12:00:00Z"
    }
  }
  ```

#### POST `/v1/charities`
Membuat campaign baru.
- **Headers:** `Authorization: Bearer <token>`
- **Request Body (Multipart Form-Data):**
  - `title` (text)
  - `description` (text)
  - `categoryId` (uuid)
  - `targetAmount` (number)
  - `startDate` (date, YYYY-MM-DD)
  - `endDate` (date, YYYY-MM-DD)
  - `coverImage` (file)
  - `additionalImages` (files)
  - `contactPerson` (text)
  - `contactEmail` (text)
  - `contactPhone` (text)
  - `organization` (text)
- **Response (201 Created):**
  ```json
  {
    "data": {
      "id": "charity-uuid",
      "title": "Bantu Pembangunan Sekolah Roboh",
      "slug": "bantu-pembangunan-sekolah-roboh",
      "status": "pending"
    }
  }
  ```

---

### 4. Donations

#### POST `/v1/donations`
Membuat donasi baru (Menghasilkan invoice pembayaran).
- **Request Body:**
  ```json
  {
    "charityId": "charity-uuid",
    "amount": 100000.00,
    "paymentMethod": "gopay",
    "anonymous": false,
    "message": "Semoga bermanfaat dan lekas selesai pembangunannya."
  }
  ```
- **Response (201 Created):**
  ```json
  {
    "data": {
      "donationId": "donation-uuid",
      "amount": 100000.00,
      "status": "pending",
      "paymentMethod": "gopay",
      "paymentAction": {
        "qrUrl": "https://payment-gateway.com/qr/123",
        "deeplink": "gopay://pay?code=123"
      }
    }
  }
  ```

#### POST `/v1/donations/callback`
Webhook/callback dari payment gateway untuk update status donasi.
- **Request Body:**
  ```json
  {
    "paymentReference": "pay-ref-123",
    "donationId": "donation-uuid",
    "status": "settlement"
  }
  ```
- **Response (200 OK):**
  ```json
  {
    "data": {
      "donationId": "donation-uuid",
      "status": "paid"
    }
  }
  ```

---

### 5. Articles

#### GET `/v1/articles`
Mendapatkan semua artikel.
- **Response (200 OK):**
  ```json
  {
    "data": [
      {
        "id": "article-uuid",
        "title": "Tips Melakukan Penggalangan Dana Sosial",
        "slug": "tips-melakukan-penggalangan-dana-sosial",
        "coverImageUrl": "https://storage.local/article1.jpg",
        "publishedAt": "2026-06-16T10:00:00Z"
      }
    ]
  }
  ```

#### GET `/v1/articles/{slug}`
Mendapatkan detail artikel berdasarkan slug.
- **Response (200 OK):**
  ```json
  {
    "data": {
      "id": "article-uuid",
      "title": "Tips Melakukan Penggalangan Dana Sosial",
      "slug": "tips-melakukan-penggalangan-dana-sosial",
      "content": "<p>Isi konten artikel lengkap...</p>",
      "coverImageUrl": "https://storage.local/article1.jpg",
      "author": {
        "fullName": "Admin Crowdfunding"
      },
      "publishedAt": "2026-06-16T10:00:00Z"
    }
  }
  ```
