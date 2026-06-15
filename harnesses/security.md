# 🔐 Security Harness

> **Tujuan:** Memastikan sistem yang dibangun memiliki lapisan keamanan yang solid sejak tahap desain.  
> **Gunakan harness ini saat:** membangun backend API, sistem autentikasi, menangani data sensitif, atau review keamanan kode.

---

## 🧭 Konteks untuk AI

```
Kamu adalah security-aware engineer. Setiap keputusan koding harus mempertimbangkan aspek keamanan.
Ikuti prinsip "Secure by Default" — sistem harus aman tanpa konfigurasi tambahan.
Jika ada trade-off antara keamanan dan kemudahan, selalu prioritaskan keamanan kecuali ada alasan bisnis yang kuat.
```

---

## 1. 🔑 Autentikasi & Otorisasi

### Prinsip
- Gunakan **JWT** dengan expiry pendek (15 menit access token, 7 hari refresh token)
- Implementasikan **refresh token rotation** — setiap refresh menghasilkan token baru
- Gunakan **OAuth 2.0 / OIDC** untuk integrasi third-party
- Terapkan **Least Privilege** — user hanya punya akses minimum yang dibutuhkan

### Checklist Autentikasi
- [ ] Password di-hash dengan `bcrypt` (cost factor ≥ 12) atau `argon2id`
- [ ] Implementasi **rate limiting** pada endpoint login (maks 5 percobaan/menit)
- [ ] **Account lockout** setelah N percobaan gagal berurutan
- [ ] **Multi-Factor Authentication (MFA)** tersedia untuk akun sensitif
- [ ] Session di-invalidasi saat logout (blacklist token atau gunakan short-lived token)
- [ ] Tidak ada credential di URL (gunakan header `Authorization: Bearer <token>`)

### Checklist Otorisasi
- [ ] Validasi otorisasi di **server-side**, bukan hanya client-side
- [ ] Gunakan **RBAC** (Role-Based Access Control) atau **ABAC** untuk sistem kompleks
- [ ] Setiap endpoint terdefinisi permission yang dibutuhkan secara eksplisit
- [ ] **Ownership check** — user hanya bisa mengakses resource miliknya
- [ ] Audit log untuk setiap aksi yang mengubah data sensitif

### [KUSTOMISASI] Stack
```
# Ganti sesuai stack teknologimu:
- Node.js    → jsonwebtoken, bcrypt, passport.js
- Python     → PyJWT, passlib, python-jose
- Go         → golang-jwt/jwt, bcrypt (std)
- Java/Spring → Spring Security, jjwt
```

---

## 2. 🛡️ Validasi Input & Sanitasi

### Prinsip
- **Never trust user input** — validasi semua input di server-side
- Gunakan **whitelist validation** bukan blacklist
- Sanitasi sebelum rendering, bukan sebelum menyimpan

### Checklist
- [ ] Validasi **tipe data**, **panjang**, dan **format** untuk semua field input
- [ ] Gunakan **parameterized queries** / **prepared statements** — tidak ada raw SQL
- [ ] Sanitasi output HTML untuk mencegah **XSS**
- [ ] Validasi **file upload**: ekstensi, MIME type, ukuran, dan scan konten
- [ ] Batasi ukuran payload request (body, header)
- [ ] Validasi dan sanitasi **redirect URL** (open redirect prevention)

### Contoh Pola yang Benar

```javascript
// ❌ SALAH — rentan SQL Injection
const user = await db.query(`SELECT * FROM users WHERE email = '${email}'`);

// ✅ BENAR — parameterized query
const user = await db.query('SELECT * FROM users WHERE email = $1', [email]);
```

```javascript
// ❌ SALAH — rentan XSS
element.innerHTML = userInput;

// ✅ BENAR — escaped output
element.textContent = userInput;
// atau gunakan library: DOMPurify.sanitize(userInput)
```

---

## 3. 🔒 Enkripsi & Manajemen Secret

### Prinsip
- **Encrypt at rest** untuk data sensitif (PII, credential, data finansial)
- **Encrypt in transit** — selalu HTTPS/TLS 1.2+
- Secret tidak pernah ada di kode atau version control

### Checklist Enkripsi
- [ ] TLS 1.2+ untuk semua komunikasi jaringan
- [ ] Sertifikat SSL valid dan auto-renew (gunakan Let's Encrypt / ACM)
- [ ] Data sensitif di database di-enkripsi menggunakan **AES-256-GCM**
- [ ] Enkripsi field-level untuk data ultra-sensitif (nomor kartu, password hint)
- [ ] Backup data juga di-enkripsi

### Checklist Secret Management
- [ ] **Tidak ada hardcoded secret** di kode (API key, password, connection string)
- [ ] Gunakan `.env` untuk development, **Secret Manager** untuk production
  - AWS: Secrets Manager / Parameter Store
  - GCP: Secret Manager
  - Azure: Key Vault
  - Self-hosted: HashiCorp Vault
- [ ] Rotasi secret secara berkala (minimal setiap 90 hari)
- [ ] `.env` masuk `.gitignore` — verifikasi dengan `git status`
- [ ] Scan repository dari secret yang bocor (gunakan `git-secrets` atau `truffleHog`)

### [KUSTOMISASI] Tools Rekomendasi
```
# Secret scanning:
- truffleHog3        → scan git history
- detect-secrets     → pre-commit hook
- GitHub Secret Scanning → otomatis untuk repo GitHub

# Secret management:
- dotenv (dev)
- HashiCorp Vault (self-hosted)
- AWS Secrets Manager (cloud)
```

---

## 4. 🌐 Keamanan API & Network

### Checklist API Security
- [ ] **CORS** dikonfigurasi ketat — whitelist domain yang diizinkan, bukan `*`
- [ ] **Rate limiting** pada semua endpoint publik
- [ ] **API versioning** untuk perubahan breaking change
- [ ] Tidak expose stack trace atau internal error detail ke client
- [ ] **HTTPS Only** — redirect HTTP ke HTTPS, gunakan HSTS header
- [ ] Validasi **Content-Type** header pada request yang membawa body

### Security Headers Wajib

```http
# Tambahkan header ini di semua response:
Strict-Transport-Security: max-age=31536000; includeSubDomains; preload
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Content-Security-Policy: default-src 'self'
Referrer-Policy: strict-origin-when-cross-origin
Permissions-Policy: geolocation=(), microphone=(), camera=()
```

### Checklist Rate Limiting

```
Strategi per endpoint:
- Login / Register     → 5 req/menit per IP
- Password reset       → 3 req/jam per email
- API publik           → 100 req/menit per API key
- API internal         → 1000 req/menit per service
```

---

## 5. 🗄️ Keamanan Database

### Checklist
- [ ] Gunakan **dedicated database user** per aplikasi (bukan root/admin)
- [ ] Database user hanya punya permission yang dibutuhkan (SELECT, INSERT, UPDATE — bukan DROP)
- [ ] **Connection pooling** dengan batas koneksi yang wajar
- [ ] Database tidak expose ke public internet (gunakan VPC / private subnet)
- [ ] **Backup otomatis** dengan retensi minimal 30 hari
- [ ] Audit log untuk query sensitif
- [ ] Nonaktifkan fitur database yang tidak dipakai

---

## 6. 🪵 Logging & Monitoring Keamanan

### Prinsip
- Log **siapa, apa, kapan, dari mana** untuk setiap aksi sensitif
- **Jangan log data sensitif** (password, token, nomor kartu)
- Implementasi alerting untuk anomali

### Checklist Logging
- [ ] Log semua **login attempt** (sukses dan gagal) beserta IP
- [ ] Log semua **perubahan permission** dan aksi admin
- [ ] Log semua **akses ke data sensitif**
- [ ] Log **error 4xx dan 5xx** dengan context yang cukup
- [ ] Tidak ada PII / credential di log
- [ ] Log disimpan di tempat yang aman dan tidak bisa dimodifikasi

### Checklist Monitoring
- [ ] Alert untuk **brute force** (banyak login gagal dari satu IP)
- [ ] Alert untuk **unusual traffic patterns** (spike mendadak)
- [ ] Alert untuk **error rate** di atas threshold
- [ ] **Health check endpoint** yang tidak expose informasi sensitif
- [ ] Integrasi dengan SIEM jika diperlukan

---

## 7. 🔄 Dependency & Supply Chain Security

### Checklist
- [ ] Audit dependencies secara rutin (`npm audit`, `pip-audit`, `govulncheck`)
- [ ] Gunakan **lockfile** (`package-lock.json`, `poetry.lock`, `go.sum`)
- [ ] Pin versi dependency di production
- [ ] Update dependency dengan strategi yang terencana (jangan sembarangan bump major)
- [ ] Scan container image untuk vulnerability (`trivy`, `snyk`)
- [ ] Verifikasi integritas package (checksum)

### [KUSTOMISASI] Tools per Ekosistem
```
Node.js   → npm audit, Snyk, Socket.dev
Python    → pip-audit, Safety, Bandit (SAST)
Go        → govulncheck, Nancy
Java      → OWASP Dependency-Check
Container → Trivy, Grype, Docker Scout
```

---

## 8. 🧪 Security Testing

### Jenis Testing

| Tipe | Tools | Kapan Dijalankan |
|---|---|---|
| **SAST** (Static Analysis) | SonarQube, Semgrep, Bandit | Setiap commit / PR |
| **DAST** (Dynamic Analysis) | OWASP ZAP, Burp Suite | Staging environment |
| **Dependency Scan** | Snyk, npm audit | CI pipeline |
| **Secret Scan** | TruffleHog, detect-secrets | Pre-commit hook |
| **Penetration Test** | Manual / third party | Sebelum major release |

### Checklist CI/CD Security Gate
- [ ] SAST scan wajib pass sebelum merge ke main
- [ ] Dependency vulnerability check di setiap PR
- [ ] Secret scan di pre-commit hook
- [ ] Container scan sebelum deploy ke production

---

## 9. 🚨 OWASP Top 10 — Quick Reference

| # | Ancaman | Mitigasi Utama |
|---|---|---|
| A01 | Broken Access Control | RBAC, ownership check, deny by default |
| A02 | Cryptographic Failures | HTTPS, enkripsi at-rest, hash yang kuat |
| A03 | Injection | Parameterized query, input validation |
| A04 | Insecure Design | Threat modeling, security requirement |
| A05 | Security Misconfiguration | Hardening, remove defaults, auto-scan config |
| A06 | Vulnerable Components | Dependency audit, update rutin |
| A07 | Auth Failures | MFA, rate limit, session management |
| A08 | Software & Data Integrity | Checksum, signed packages, SBOM |
| A09 | Logging Failures | Centralized logging, alerting |
| A10 | SSRF | Whitelist URL, block internal IPs |

---

## 10. 📋 Security Review Checklist (Pre-Launch)

Gunakan ini sebelum deploy ke production:

### Autentikasi & Akses
- [ ] Semua endpoint ter-autentikasi kecuali yang memang publik
- [ ] MFA tersedia untuk admin dan akun sensitif
- [ ] Refresh token rotation berjalan dengan benar
- [ ] Session invalidation saat logout berfungsi

### Data & Enkripsi
- [ ] Tidak ada secret/credential di kode atau environment file yang ter-commit
- [ ] Data sensitif di-enkripsi at-rest
- [ ] HTTPS aktif di semua endpoint production
- [ ] HSTS header aktif

### Input & Output
- [ ] Semua input tervalidasi di server-side
- [ ] Tidak ada raw query string (SQL injection prevention)
- [ ] Output di-sanitasi untuk mencegah XSS
- [ ] File upload divalidasi dan di-scan

### Infrastruktur
- [ ] Database tidak exposed ke internet
- [ ] Security headers dikonfigurasi di web server / API gateway
- [ ] CORS dikonfigurasi dengan whitelist yang ketat
- [ ] Rate limiting aktif di semua endpoint publik

### Monitoring
- [ ] Logging aktif dan tidak menyimpan data sensitif
- [ ] Alert untuk login anomali aktif
- [ ] Health check endpoint tidak expose informasi sistem

---

## 📚 Referensi

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [OWASP Cheat Sheet Series](https://cheatsheetseries.owasp.org/)
- [NIST Cybersecurity Framework](https://www.nist.gov/cyberframework)
- [CWE/SANS Top 25](https://cwe.mitre.org/top25/)

---

*Security Harness v1.0 — bagian dari [Harness Engineering](../README.md)*
