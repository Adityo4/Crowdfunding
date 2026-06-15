# 🗺️ Modeling Harness

> **Tujuan:** Mendefinisikan solusi dari sebuah masalah dalam bentuk visual — flow, diagram, dan model — sebelum masuk ke implementasi.
> **Gunakan harness ini saat:** mendefinisikan solusi baru, mendokumentasikan sistem yang ada, menjelaskan alur ke tim, atau sebagai blueprint sebelum coding.

---

## 🧭 Konteks untuk AI

```
Kamu adalah solution modeler yang berpengalaman.
Setiap masalah harus dimodelkan secara visual sebelum diselesaikan dengan kode.

Prinsipmu:
1. Model dulu, kode kemudian — visual memaksa kamu berpikir lebih dalam
2. Pilih jenis diagram yang paling tepat untuk masalah yang ada
3. Diagram harus dapat dibaca orang lain tanpa penjelasan verbal
4. Mulai dari high-level, lalu zoom in ke detail yang kompleks
5. Gunakan Mermaid untuk semua diagram — text-based, version-controllable

Sebelum membuat diagram, SELALU tanya:
"Apa yang ingin dijelaskan? Siapa audiensnya? Apa keputusan yang harus difasilitasi?"
```

---

## 1. 🧭 Panduan Memilih Diagram

Gunakan tabel ini untuk menentukan diagram yang tepat:

| Pertanyaan | Diagram yang Tepat |
|---|---|
| Bagaimana alur kerja user? | [User Flow](#3-user-flow) |
| Bagaimana komponen sistem berinteraksi? | [Sequence Diagram](#4-sequence-diagram) |
| Apa saja state yang mungkin ada? | [State Diagram](#5-state-diagram) |
| Bagaimana data direlasikan? | [ER Diagram](#6-er-diagram-entity-relationship) |
| Bagaimana keputusan dibuat? | [Decision Tree / Flowchart](#7-decision-tree--flowchart) |
| Apa saja komponen & dependency-nya? | [Component Diagram](#8-component-diagram) |
| Bagaimana data mengalir di sistem? | [Data Flow Diagram](#9-data-flow-diagram) |
| Apa timeline & urutan prosesnya? | [Sequence Diagram](#4-sequence-diagram) |
| Bagaimana struktur kelas/domain? | [Class Diagram](#10-class--domain-diagram) |

---

## 2. 🔍 Problem Decomposition — Sebelum Diagram

Sebelum membuat diagram apapun, lakukan dekomposisi masalah:

```
Framework PROBLEM → SOLUTION MODEL:

1. DEFINE THE PROBLEM
   "Apa masalah yang ingin diselesaikan?"
   "Siapa yang terdampak?"
   "Apa akibatnya jika tidak diselesaikan?"

2. IDENTIFY ACTORS
   "Siapa saja yang berinteraksi dengan sistem?"
   (User, Admin, System, External Service, Scheduler, dll.)

3. MAP THE BOUNDARIES
   "Di mana sistem ini dimulai dan berakhir?"
   "Apa yang di dalam scope dan apa yang di luar?"

4. BREAK DOWN THE FLOW
   "Apa langkah-langkah dari awal hingga selesai?"
   "Di mana keputusan dibuat?"
   "Di mana kemungkinan gagal?"

5. IDENTIFY DATA
   "Data apa yang dibutuhkan?"
   "Data apa yang dihasilkan?"
   "Di mana data disimpan?"

6. CHOOSE YOUR DIAGRAM(S)
   Gunakan tabel di atas untuk memilih jenis diagram yang tepat.
```

---

## 3. 🧑‍💻 User Flow

**Kapan dipakai:** Menggambarkan perjalanan user dari awal hingga tujuan akhir.

### Template Mermaid — User Flow

```mermaid
flowchart TD
    Start([👤 User]) --> A[Buka halaman login]
    A --> B{Sudah punya akun?}
    B -->|Ya| C[Isi email & password]
    B -->|Tidak| D[Klik daftar]
    D --> E[Isi form registrasi]
    E --> F[Submit]
    F --> G{Validasi berhasil?}
    G -->|Tidak| H[Tampilkan error]
    H --> E
    G -->|Ya| I[Kirim email verifikasi]
    I --> J([Selesai — tunggu verifikasi])

    C --> K[Submit login]
    K --> L{Kredensial valid?}
    L -->|Tidak| M[Tampilkan error login]
    M --> C
    L -->|Ya| N([Dashboard])

    style Start fill:#4f46e5,color:#fff
    style N fill:#16a34a,color:#fff
    style J fill:#0891b2,color:#fff
    style H fill:#dc2626,color:#fff
    style M fill:#dc2626,color:#fff
```

### Panduan User Flow
```
Node shapes yang bermakna:
  ([Text])  → Start / End (rounded)
  [Text]    → Proses / Action
  {Text}    → Decision / Condition
  [(Text)]  → Database
  [[Text]]  → Sub-process

Aturan:
✅ Mulai dari satu titik awal yang jelas
✅ Setiap decision node punya semua cabang yang mungkin
✅ Error state / unhappy path harus ikut digambarkan
✅ Setiap path berujung di end state
❌ Jangan biarkan flow menggantung tanpa end state
```

---

## 4. 🔄 Sequence Diagram

**Kapan dipakai:** Menggambarkan interaksi antar komponen dalam urutan waktu.

### Template Mermaid — Sequence Diagram

```mermaid
sequenceDiagram
    actor User
    participant Browser
    participant API Gateway
    participant Auth Service
    participant User Service
    participant Database

    User->>Browser: Klik tombol Login
    Browser->>API Gateway: POST /auth/login {email, password}
    API Gateway->>Auth Service: validateCredentials(email, password)
    Auth Service->>Database: SELECT user WHERE email=?
    Database-->>Auth Service: User record

    alt Kredensial valid
        Auth Service->>Auth Service: Generate JWT token
        Auth Service-->>API Gateway: { token, refreshToken }
        API Gateway-->>Browser: 200 OK { token }
        Browser-->>User: Redirect ke Dashboard
    else Kredensial tidak valid
        Auth Service-->>API Gateway: Error: INVALID_CREDENTIALS
        API Gateway-->>Browser: 401 Unauthorized
        Browser-->>User: Tampilkan pesan error
    end

    Note over Auth Service,Database: Token disimpan di Redis untuk invalidasi
```

### Panduan Sequence Diagram
```
Notasi penting:
  ->>   Panah solid (synchronous call)
  -->>  Panah putus (response / return)
  -x    Panah dengan X (failed call)
  -)    Panah async (fire and forget)

  alt/else  Conditional block
  loop      Repeat block
  par       Parallel execution
  Note      Komentar / anotasi

Tips:
✅ Urutkan participant dari kiri (user) ke kanan (storage)
✅ Selalu gambarkan response, tidak hanya request
✅ Gunakan alt/else untuk menunjukkan happy vs error path
✅ Tambahkan Note untuk keputusan yang tidak obvious
```

---

## 5. 🔁 State Diagram

**Kapan dipakai:** Menggambarkan semua state yang mungkin dari sebuah entitas dan transisinya.

### Template Mermaid — State Diagram

```mermaid
stateDiagram-v2
    [*] --> Draft : Order dibuat

    Draft --> Pending : User submit order
    Draft --> Cancelled : User batalkan

    Pending --> Processing : Payment berhasil
    Pending --> Cancelled : Payment gagal / timeout

    Processing --> Shipped : Barang dikirim
    Processing --> Cancelled : Stok habis

    Shipped --> Delivered : Konfirmasi terima
    Shipped --> Returned : User ajukan return

    Delivered --> Completed : Periode return habis
    Delivered --> Returned : User ajukan return

    Returned --> Refunded : Barang diterima & diperiksa

    Cancelled --> [*]
    Completed --> [*]
    Refunded --> [*]

    note right of Processing
        Trigger email notifikasi
        ke user saat state berubah
    end note
```

### Panduan State Diagram
```
Aturan state diagram yang baik:
✅ Setiap state yang mungkin terdaftar (jangan ada state tersembunyi)
✅ Setiap transisi punya trigger/kondisi yang jelas
✅ Ada initial state [*] dan final state(s) [*]
✅ Tidak ada state yang tidak bisa dicapai (unreachable state)
✅ Gambarkan efek samping (side effects) dengan Note

Gunakan state diagram untuk:
  - Order lifecycle
  - User account status
  - Payment status
  - Document approval workflow
  - Session state
```

---

## 6. 🗄️ ER Diagram (Entity-Relationship)

**Kapan dipakai:** Menggambarkan struktur data dan relasi antar entitas.

### Template Mermaid — ER Diagram

```mermaid
erDiagram
    USER {
        uuid id PK
        string email UK
        string name
        enum status "active|inactive|banned"
        timestamp created_at
        timestamp updated_at
    }

    ORDER {
        uuid id PK
        uuid user_id FK
        enum status "draft|pending|processing|shipped|delivered|cancelled"
        decimal total_amount
        string shipping_address
        timestamp ordered_at
    }

    ORDER_ITEM {
        uuid id PK
        uuid order_id FK
        uuid product_id FK
        int quantity
        decimal unit_price
        decimal subtotal
    }

    PRODUCT {
        uuid id PK
        string name
        string sku UK
        decimal price
        int stock
        uuid category_id FK
    }

    CATEGORY {
        uuid id PK
        string name
        uuid parent_id FK "nullable — untuk sub-kategori"
    }

    PAYMENT {
        uuid id PK
        uuid order_id FK
        enum method "credit_card|bank_transfer|e_wallet"
        enum status "pending|success|failed"
        decimal amount
        string gateway_reference
        timestamp paid_at
    }

    USER ||--o{ ORDER : "membuat"
    ORDER ||--|{ ORDER_ITEM : "berisi"
    ORDER ||--o| PAYMENT : "dibayar dengan"
    PRODUCT ||--o{ ORDER_ITEM : "ada di"
    CATEGORY ||--o{ PRODUCT : "mengelompokkan"
    CATEGORY ||--o{ CATEGORY : "punya sub-kategori"
```

### Panduan ER Diagram
```
Notasi kardinalitas:
  ||--||   Exactly one to exactly one
  ||--o|   One to zero or one
  ||--|{   One to one or many
  ||--o{   One to zero or many
  }o--o{   Zero or many to zero or many

Tips:
✅ Selalu sertakan PK, FK, dan UK yang jelas
✅ Sertakan tipe data dan constraint penting
✅ Gambarkan relasi self-referential jika ada (tree structure)
✅ Tambahkan komentar untuk kolom yang tidak obvious
✅ Fokus pada entitas bisnis, bukan detail implementasi DB
```

---

## 7. 🌿 Decision Tree / Flowchart

**Kapan dipakai:** Menggambarkan logika keputusan yang kompleks atau business rules.

### Template Mermaid — Decision Tree

```mermaid
flowchart TD
    Start([Proses Pengajuan Kredit]) --> A{Usia >= 21?}

    A -->|Tidak| Reject1([❌ Ditolak: Usia tidak memenuhi])
    A -->|Ya| B{Punya penghasilan tetap?}

    B -->|Tidak| C{Aset > 500jt?}
    B -->|Ya| D{Gaji >= 3x cicilan?}

    C -->|Tidak| Reject2([❌ Ditolak: Tidak ada jaminan])
    C -->|Ya| E[Proses sebagai debitur aset]

    D -->|Tidak| Reject3([❌ Ditolak: Rasio cicilan terlalu tinggi])
    D -->|Ya| F{Credit score?}

    F -->|< 500| Reject4([❌ Ditolak: Credit score rendah])
    F -->|500-700| G[Approve dengan bunga tinggi]
    F -->|> 700| H[Approve dengan bunga standar]

    E --> F

    G --> Approve1([✅ Disetujui — Rate 12%])
    H --> Approve2([✅ Disetujui — Rate 9%])

    style Reject1 fill:#dc2626,color:#fff
    style Reject2 fill:#dc2626,color:#fff
    style Reject3 fill:#dc2626,color:#fff
    style Reject4 fill:#dc2626,color:#fff
    style Approve1 fill:#16a34a,color:#fff
    style Approve2 fill:#16a34a,color:#fff
```

### Panduan Decision Tree
```
Tips untuk decision tree yang efektif:
✅ Setiap decision node hanya punya SATU pertanyaan
✅ Setiap cabang punya label yang jelas (Ya/Tidak, atau nilai spesifik)
✅ Semua path berujung di hasil yang jelas
✅ Urutan kondisi dioptimalkan (kondisi yang paling banyak mengeliminasi di atas)
✅ Gunakan warna untuk membedakan hasil (hijau = ok, merah = reject)
❌ Jangan campur business logic dengan implementasi teknikal
```

---

## 8. 🧩 Component Diagram

**Kapan dipakai:** Menggambarkan komponen sistem dan ketergantungan antar komponen.

### Template Mermaid — Component Diagram

```mermaid
graph TB
    subgraph Client["🖥️ Client Layer"]
        WebApp["Web App\n(React/Next.js)"]
        MobileApp["Mobile App\n(React Native)"]
    end

    subgraph Gateway["🔀 Gateway Layer"]
        APIGateway["API Gateway\n(Kong / Nginx)"]
        CDN["CDN\n(Cloudflare)"]
    end

    subgraph Services["⚙️ Service Layer"]
        AuthSvc["Auth Service"]
        UserSvc["User Service"]
        OrderSvc["Order Service"]
        PaymentSvc["Payment Service"]
        NotifSvc["Notification Service"]
    end

    subgraph Data["🗄️ Data Layer"]
        PostgreSQL[("PostgreSQL\n(Primary)")]
        Redis[("Redis\n(Cache)")]
        S3["S3\n(Object Storage)"]
    end

    subgraph External["🌐 External Services"]
        PaymentGW["Payment Gateway\n(Midtrans)"]
        EmailSvc["Email Provider\n(SendGrid)"]
        SMSSvc["SMS Provider\n(Twilio)"]
    end

    WebApp --> CDN
    WebApp --> APIGateway
    MobileApp --> APIGateway

    APIGateway --> AuthSvc
    APIGateway --> UserSvc
    APIGateway --> OrderSvc
    APIGateway --> PaymentSvc

    OrderSvc --> PaymentSvc
    PaymentSvc --> PaymentGW
    PaymentSvc --> NotifSvc
    OrderSvc --> NotifSvc

    AuthSvc --> PostgreSQL
    UserSvc --> PostgreSQL
    OrderSvc --> PostgreSQL
    PaymentSvc --> PostgreSQL

    AuthSvc --> Redis
    UserSvc --> S3

    NotifSvc --> EmailSvc
    NotifSvc --> SMSSvc
```

---

## 9. 🌊 Data Flow Diagram

**Kapan dipakai:** Menggambarkan bagaimana data bergerak melalui sistem.

### Template Mermaid — Data Flow

```mermaid
flowchart LR
    User([👤 User])

    subgraph Input["📥 Data Input"]
        Form["Form Order\n{produk, qty, alamat}"]
        Payment["Data Pembayaran\n{metode, nomor kartu}"]
    end

    subgraph Process["⚙️ Processing"]
        Validate{"Validasi\nInput"}
        CalcTotal["Hitung Total\n+ pajak + ongkir"]
        ProcessPay["Proses\nPembayaran"]
        CreateOrder["Buat Order\nRecord"]
    end

    subgraph Storage["🗄️ Storage"]
        OrderDB[("Orders DB")]
        PaymentDB[("Payments DB")]
        Cache[("Redis Cache")]
    end

    subgraph Output["📤 Output"]
        Confirm["Konfirmasi Order\n{orderId, status}"]
        Invoice["Invoice PDF"]
        Notif["Email + SMS\nNotifikasi"]
    end

    User -->|mengisi| Form
    User -->|memasukkan| Payment
    Form --> Validate
    Payment --> Validate
    Validate -->|valid| CalcTotal
    Validate -->|invalid| User

    CalcTotal --> ProcessPay
    ProcessPay -->|berhasil| CreateOrder
    ProcessPay -->|gagal| User

    CreateOrder --> OrderDB
    CreateOrder --> PaymentDB
    CreateOrder --> Cache

    OrderDB --> Confirm
    OrderDB --> Invoice
    CreateOrder --> Notif

    Confirm --> User
    Invoice --> User
    Notif --> User
```

---

## 10. 🏛️ Class / Domain Diagram

**Kapan dipakai:** Menggambarkan struktur domain model dan relasi antar class.

### Template Mermaid — Class Diagram

```mermaid
classDiagram
    class User {
        +UUID id
        +String email
        +String name
        +UserStatus status
        +DateTime createdAt
        +login(email, password) Token
        +logout(token) void
        +updateProfile(data) User
    }

    class Order {
        +UUID id
        +UUID userId
        +OrderStatus status
        +Money totalAmount
        +Address shippingAddress
        +submit() void
        +cancel(reason) void
        +calculateTotal() Money
    }

    class OrderItem {
        +UUID id
        +UUID orderId
        +UUID productId
        +int quantity
        +Money unitPrice
        +getSubtotal() Money
    }

    class Product {
        +UUID id
        +String name
        +String sku
        +Money price
        +int stock
        +isAvailable(qty) bool
        +reserveStock(qty) void
    }

    class Payment {
        +UUID id
        +UUID orderId
        +PaymentMethod method
        +PaymentStatus status
        +Money amount
        +process() PaymentResult
        +refund() void
    }

    class Money {
        +decimal amount
        +String currency
        +add(other) Money
        +multiply(factor) Money
    }

    User "1" --> "0..*" Order : places
    Order "1" *-- "1..*" OrderItem : contains
    Order "1" --> "0..1" Payment : paid by
    OrderItem "0..*" --> "1" Product : references
    Order ..> Money : uses
    OrderItem ..> Money : uses
    Payment ..> Money : uses
```

---

## 11. 📐 Template Kosong per Diagram

Copy template ini sebagai starting point:

### User Flow Kosong
```mermaid
flowchart TD
    Start([Actor]) --> A[Langkah 1]
    A --> B{Keputusan?}
    B -->|Ya| C[Langkah 2a]
    B -->|Tidak| D[Langkah 2b]
    C --> End1([Hasil A])
    D --> End2([Hasil B])
```

### Sequence Kosong
```mermaid
sequenceDiagram
    actor User
    participant SystemA
    participant SystemB

    User->>SystemA: Request
    SystemA->>SystemB: Internal Call
    SystemB-->>SystemA: Response
    SystemA-->>User: Result
```

### State Kosong
```mermaid
stateDiagram-v2
    [*] --> StateA : trigger1
    StateA --> StateB : trigger2
    StateB --> StateC : trigger3
    StateC --> [*]
```

### ER Kosong
```mermaid
erDiagram
    ENTITY_A {
        uuid id PK
        string field1
        uuid entity_b_id FK
    }
    ENTITY_B {
        uuid id PK
        string field1
    }
    ENTITY_A }o--|| ENTITY_B : "relasi"
```

---

## 12. 📋 Modeling Checklist

### Sebelum Membuat Diagram
- [ ] Masalah sudah didekomposisi (Problem Decomposition framework)
- [ ] Jenis diagram yang paling tepat sudah dipilih
- [ ] Audiens diagram sudah jelas (developer? stakeholder? QA?)
- [ ] Scope diagram sudah terdefinisi (tidak terlalu luas, tidak terlalu sempit)

### Kualitas Diagram
- [ ] Diagram bisa dipahami tanpa penjelasan verbal
- [ ] Semua actor / komponen yang relevan terwakili
- [ ] Happy path DAN error/edge case tergambar
- [ ] Tidak ada node yang menggantung (setiap path punya end state)
- [ ] Label pada setiap edge / panah jelas

### Setelah Membuat Diagram
- [ ] Diagram di-review oleh minimal satu orang lain
- [ ] Diagram disimpan di dokumentasi project (bukan hanya di chat)
- [ ] Jika ada perubahan implementasi: diagram diperbarui

---

## 📚 Referensi

- [Mermaid Live Editor](https://mermaid.live) — preview & edit diagram secara real-time
- [Mermaid Documentation](https://mermaid.js.org/intro/)
- [C4 Model](https://c4model.com/) — framework untuk software architecture diagrams
- [BPMN](https://www.bpmn.org/) — Business Process Modeling Notation
- [Martin Fowler — UML Distilled](https://www.martinfowler.com/books/uml.html)

---

*Modeling Harness v1.0 — bagian dari [Harness Engineering](../README.md)*
