# 🚢 DevOps Harness

> **Tujuan:** Memastikan proses build, test, deploy, dan monitoring berjalan otomatis, repeatable, dan reliable.
> **Gunakan harness ini saat:** setup CI/CD pipeline, konfigurasi container, merancang deployment strategy, atau membangun observability.

---

## 🧭 Konteks untuk AI

```
Kamu adalah DevOps/platform engineer yang berpengalaman.
Prinsipmu:
1. Automate everything — jika dilakukan lebih dari sekali, otomasi
2. Infrastructure as Code — semua konfigurasi infra ada di version control
3. Fail fast — deteksi masalah seawal mungkin di pipeline
4. Deploy kecil dan sering — lebih aman daripada deploy besar dan jarang
5. Observability first — sistem yang tidak bisa diobservasi tidak bisa di-debug

Sebelum setup infra, SELALU tanya: "Bagaimana cara rollback jika gagal?"
```

---

## 1. 🔄 CI/CD Pipeline

### Prinsip
- Pipeline harus **deterministik** — hasil yang sama untuk input yang sama
- **Fail fast** — jalankan check yang paling cepat dan paling sering gagal duluan
- Setiap commit ke main branch harus **deployable**
- Pipeline harus bisa dijalankan lokal (sebisa mungkin)

### Struktur Pipeline yang Direkomendasikan

```
Push / PR
  │
  ├─▶ [Lint & Format]      → Cepat, ~30 detik
  │     └── Code style, type check
  │
  ├─▶ [Unit Test]          → Cepat, ~1-2 menit
  │     └── Coverage check
  │
  ├─▶ [Security Scan]      → Paralel, ~2 menit
  │     └── SAST, secret scan, dependency audit
  │
  ├─▶ [Build]              → ~2-5 menit
  │     └── Compile, bundle, Docker build
  │
  ├─▶ [Integration Test]   → ~5-10 menit
  │     └── DB test, API test
  │
  └─▶ [Deploy Staging]     → Jika semua lulus
        │
        └─▶ [E2E Test]     → ~5-15 menit
              │
              └─▶ [Deploy Production]  ← Manual approval atau auto
```

### [KUSTOMISASI] CI/CD Platform
```
- GitHub Actions   → populer, terintegrasi GitHub, YAML-based
- GitLab CI/CD     → powerful, self-hosted option
- Jenkins          → mature, plugin ecosystem besar
- CircleCI         → cepat, developer experience bagus
- Bitbucket Pipelines → jika pakai Atlassian stack
- Azure DevOps     → jika di ekosistem Microsoft
```

### Checklist CI Pipeline
- [ ] Pipeline berjalan otomatis di setiap push dan PR
- [ ] PR tidak bisa di-merge jika pipeline gagal
- [ ] Pipeline selesai dalam waktu yang wajar (< 15 menit untuk PR)
- [ ] Cache dependency untuk mempercepat pipeline
- [ ] Secret/credential disimpan di secret manager CI, bukan di kode
- [ ] Pipeline bisa di-debug secara lokal (gunakan `act` untuk GitHub Actions)

---

## 2. 🐳 Containerization (Docker)

### Prinsip
- **One process per container** — pisahkan concern antar container
- Image harus **minimal** — gunakan base image yang ramping
- Container harus **stateless** — state disimpan di luar container (DB, volume)
- Image harus **immutable** — tag dengan versi spesifik, bukan `latest`

### Dockerfile Best Practices

```dockerfile
# ✅ BENAR — multi-stage build untuk image yang kecil

# Stage 1: Build
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production   # install dependency dulu (manfaatkan layer cache)
COPY . .
RUN npm run build

# Stage 2: Runtime (hanya artifact yang dibutuhkan)
FROM node:20-alpine AS runtime
WORKDIR /app
RUN addgroup -S appgroup && adduser -S appuser -G appgroup  # non-root user
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
USER appuser                   # jalankan sebagai non-root
EXPOSE 3000
CMD ["node", "dist/index.js"]
```

### Checklist Docker
- [ ] Multi-stage build dipakai untuk image production
- [ ] Base image menggunakan versi spesifik (bukan `latest`)
- [ ] `.dockerignore` dikonfigurasi (exclude `node_modules`, `.git`, dsb)
- [ ] Container berjalan sebagai **non-root user**
- [ ] Health check dikonfigurasi di Dockerfile atau compose
- [ ] Image di-scan untuk vulnerability sebelum deploy (`trivy`, `docker scout`)
- [ ] Secret tidak di-hardcode di Dockerfile atau image layer

### Docker Compose (Development)

```yaml
# docker-compose.yml — untuk development lokal
services:
  app:
    build: .
    ports:
      - "3000:3000"
    environment:
      - NODE_ENV=development
    volumes:
      - .:/app                 # mount source untuk hot reload
      - /app/node_modules      # jangan override node_modules di container
    depends_on:
      db:
        condition: service_healthy

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: myapp_dev
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password  # dev only — jangan di production!
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U user"]
      interval: 5s
      timeout: 5s
      retries: 5
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

---

## 3. 🌍 Deployment Strategy

### Strategi Deployment

| Strategi | Cara Kerja | Kapan Dipakai |
|---|---|---|
| **Recreate** | Stop lama, start baru — downtime ada | Dev/staging, bukan production |
| **Rolling Update** | Ganti instance satu per satu | Default di Kubernetes, zero-downtime |
| **Blue-Green** | Siapkan environment baru, switch traffic | Zero-downtime, rollback cepat |
| **Canary** | Kirim sebagian kecil traffic ke versi baru | Rilis berisiko tinggi, uji di production |
| **Feature Flag** | Deploy kode, kontrol via flag | Rilis bertahap tanpa deploy ulang |

### Blue-Green Deployment

```
Production Traffic
        │
        ▼
  [Load Balancer]
        │
   ┌────┴────┐
   │         │
[Blue]    [Green]    ← Green = versi baru, sudah ready
(v1.0)    (v1.1)
   │         │
   └────┬────┘
        │
   Switch traffic → Green (jika OK)
   Rollback       → Blue (jika ada masalah)
```

### Checklist Deployment
- [ ] **Rollback plan** terdefinisi dan pernah dicoba
- [ ] Database migration bisa di-rollback (reversible migration)
- [ ] **Health check** endpoint tersedia dan dikonfigurasi di load balancer
- [ ] **Graceful shutdown** diimplementasi (drain existing request)
- [ ] Deployment tidak dilakukan di jam sibuk (kecuali urgent)
- [ ] Monitoring dipantau aktif saat dan setelah deployment
- [ ] Changelog terdokumentasi untuk setiap deployment

---

## 4. 🏗️ Infrastructure as Code (IaC)

### Prinsip
- **Semua infra di version control** — tidak ada konfigurasi manual di server
- Infra harus **idempotent** — apply berulang kali menghasilkan state yang sama
- Gunakan **environment yang terpisah** (dev, staging, production)
- **Least privilege** untuk semua IAM role dan service account

### [KUSTOMISASI] Tools IaC
```
Provisioning:
- Terraform    → multi-cloud, paling populer
- Pulumi       → IaC dengan bahasa pemrograman (TypeScript, Python, Go)
- AWS CDK      → khusus AWS, pakai bahasa pemrograman
- Bicep        → khusus Azure

Configuration Management:
- Ansible      → agentless, YAML-based
- Chef/Puppet  → agent-based, mature

Container Orchestration:
- Kubernetes   → standar industri, kompleks tapi powerful
- Docker Swarm → lebih simpel dari K8s
- ECS (AWS)    → managed, lebih mudah dari K8s jika di AWS
- Cloud Run (GCP) / App Service (Azure) → serverless container
```

### Checklist IaC
- [ ] Semua resource infra terdefinisi di kode (tidak ada yang manual)
- [ ] State file IaC disimpan di remote backend (S3, GCS, Terraform Cloud)
- [ ] Environment (dev/staging/prod) dipisah dengan jelas
- [ ] Perubahan infra melalui PR review sebelum di-apply
- [ ] Secret dikelola oleh secret manager, bukan di IaC file

---

## 5. 📊 Observability

### Tiga Pilar Observability

```
LOGS     → Apa yang terjadi? (events, error messages)
METRICS  → Seberapa sering / berapa banyak? (count, gauge, histogram)
TRACES   → Di mana waktunya habis? (distributed request tracing)
```

### Checklist Logging
- [ ] Log level dikonfigurasi per environment (DEBUG di dev, INFO/WARN di prod)
- [ ] Semua error di-log dengan stack trace dan context yang cukup
- [ ] Log tidak menyimpan data sensitif (password, token, PII)
- [ ] Log terstruktur (JSON format) untuk memudahkan parsing
- [ ] Log di-ship ke centralized logging (bukan hanya di lokal server)
- [ ] Log retention policy terdefinisi

### Checklist Metrics
- [ ] **Business metrics** dipantau (jumlah order, user aktif, dsb)
- [ ] **Technical metrics** dipantau (CPU, RAM, DB connection, queue depth)
- [ ] **SLO metrics** dipantau (availability, latency)
- [ ] Dashboard tersedia dan mudah dibaca oleh seluruh tim
- [ ] Alert dikonfigurasi untuk anomali yang penting

### Checklist Distributed Tracing
- [ ] Trace ID diteruskan antar service (via header)
- [ ] Span dibuat untuk operasi yang penting (DB query, external API call)
- [ ] Sampling rate dikonfigurasi (100% untuk error, ~10% untuk normal traffic)

### [KUSTOMISASI] Tools Observability
```
Logging:
- ELK Stack (Elasticsearch + Logstash + Kibana) → self-hosted
- Loki + Grafana → lebih ringan, recommended untuk self-hosted
- Datadog Logs   → managed, terintegrasi dengan APM
- CloudWatch Logs (AWS)

Metrics:
- Prometheus + Grafana → standar open-source, recommended
- Datadog Metrics
- CloudWatch Metrics (AWS)

Tracing:
- OpenTelemetry (instrumentasi) + Jaeger/Tempo (backend)
- Datadog APM
- AWS X-Ray
```

---

## 6. 🚨 Incident Management

### Severity Level

| Level | Dampak | Response Time | Contoh |
|---|---|---|---|
| P1 — Critical | Sistem down, semua user terdampak | < 15 menit | Production outage |
| P2 — High | Fitur utama tidak berfungsi | < 1 jam | Payment gagal |
| P3 — Medium | Fitur non-kritikal terdampak | < 4 jam | Laporan export error |
| P4 — Low | Minor issue, ada workaround | < 24 jam | UI glitch |

### Runbook Template

```markdown
# Runbook: [Nama Incident / Alert]

## Gejala
[Apa yang terlihat/dialami user?]

## Kemungkinan Penyebab
1. [Penyebab 1]
2. [Penyebab 2]

## Langkah Diagnosis
1. Cek [metric/log ini] di [tool ini]
2. Jalankan [command ini] untuk verifikasi
3. Cek [service ini] statusnya

## Langkah Penanganan
### Jika penyebab A:
1. [Langkah 1]
2. [Langkah 2]

### Jika penyebab B:
1. [Langkah 1]

## Rollback (jika karena deployment)
1. [Langkah rollback]

## Eskalasi
- Jika tidak selesai dalam [X menit] → hubungi [nama/tim]
```

### Checklist Post-Incident
- [ ] Incident timeline terdokumentasi
- [ ] Root cause teridentifikasi
- [ ] Action items terdefinisi dengan owner dan deadline
- [ ] Blameless post-mortem dilakukan dalam 48-72 jam

---

## 7. 📋 DevOps Checklist (Pre-Launch)

### CI/CD
- [ ] Pipeline berjalan otomatis dan semua gate lulus
- [ ] Secret tidak ada di kode atau pipeline config
- [ ] Rollback plan terdefinisi dan sudah dicoba

### Container & Deployment
- [ ] Image di-scan dan bersih dari vulnerability kritis
- [ ] Container berjalan sebagai non-root
- [ ] Health check dikonfigurasi
- [ ] Graceful shutdown diimplementasi

### Infrastructure
- [ ] Semua infra terdefinisi sebagai kode
- [ ] Environment staging identik dengan production
- [ ] Backup dan restore pernah dicoba

### Observability
- [ ] Log, metric, dan trace aktif di production
- [ ] Alert dikonfigurasi untuk semua SLO
- [ ] Dashboard tersedia untuk on-call engineer
- [ ] Runbook tersedia untuk alert yang paling sering terjadi

---

## 📚 Referensi

- [The DevOps Handbook](https://itrevolution.com/product/the-devops-handbook/)
- [Site Reliability Engineering (Google)](https://sre.google/sre-book/table-of-contents/)
- [12 Factor App](https://12factor.net/)
- [OpenTelemetry](https://opentelemetry.io/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)

---

*DevOps Harness v1.0 — bagian dari [Harness Engineering](../README.md)*
