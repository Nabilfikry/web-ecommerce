# 🛍️ TokoKuPedia — Enterprise Full-Stack E-Commerce Platform

[![CI / Quality Checks](https://img.shields.io/badge/Quality%20Checks-Passing-00AA5B?style=for-the-badge&logo=githubactions&logoColor=white)](https://github.com/Nabilfikry/web-ecommerce)
[![Frontend Stack](https://img.shields.io/badge/Frontend-Svelte%205%20%7C%20SvelteKit%20%7C%20Tailwind%20v4-FF3E00?style=for-the-badge&logo=svelte&logoColor=white)](https://svelte.dev)
[![Backend Stack](https://img.shields.io/badge/Backend-Go%201.24%20%7C%20REST%20API%20%7C%20GORM-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Database](https://img.shields.io/badge/Database-PostgreSQL%2016%20(ACID)-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org)
[![Payment Gateway](https://img.shields.io/badge/Payment-DOKU%20Jokul%20Checkout-E31B23?style=for-the-badge&logo=contactlesspayment&logoColor=white)](https://doku.com)
[![Monorepo](https://img.shields.io/badge/Monorepo-Turborepo%20%7C%20Bun-000000?style=for-the-badge&logo=turborepo&logoColor=white)](https://turbo.build)

> **TokoKuPedia** adalah platform e-commerce marketplace skala *enterprise* yang mereplikasi secara presisi (*pixel-perfect clone*) antarmuka dan pengalaman pengguna **Tokopedia**, didukung oleh arsitektur backend performa tinggi berbasis **Go (Golang)**, basis data relasional **PostgreSQL** dengan integritas transaksi ACID, integrasi resmi **Payment Gateway DOKU**, serta arsitektur **Turborepo Monorepo**.

---

## 📌 Executive Summary

Jika Anda meninjau repositori ini untuk keperluan seleksi talenta *Full-Stack / Frontend / Backend Engineer*, berikut adalah ringkasan kemampuan teknis dan standar industri yang diimplementasikan dalam proyek ini:

| Dimensi | Penjelasan Solusi & Standar Rekayasa |
|---|---|
| **Business Impact** | Marketplace multi-seller dengan alur lengkap pembeli (*buyer*) dan penjual (*seller*): katalog dinamis, keranjang terisolasi per toko (BR-010), kalkulasi diskon/voucher instan, *one-click checkout*, integrasi gateway pembayaran nyata, dan dashboard pengelolaan stok penjual. |
| **Frontend Craftsmanship** | Replikasi UI Tokopedia Unify secara presisi (Palet hijau `#00AA5B`, 8 rute identik termasuk PDP 3-kolom Surface Pro 12, split-screen register dengan maskot Toped, keranjang berkelompok, dan pengaturan akun lengkap dengan mode Terang/Gelap/Sistem). Menggunakan **Svelte 5 Runes** (`$state`, `$derived`, `$props`), **Tailwind CSS v4**, dan **TypeScript**. |
| **Backend Architecture** | Backend Go dibangun dengan pola **Clean Architecture / Domain-Driven Design (DDD)** yang modular (`handler` $\rightarrow$ `service` $\rightarrow$ `repository` $\rightarrow$ `model`). Terpisah secara jelas (*separation of concerns*), mudah diuji (*testable*), dan mudah diperluas (*scalable*). |
| **Data Integrity & Concurrency** | **Zero-Race Condition Stock Locking**. Menggunakan transaksi database PostgreSQL dengan klausul `SELECT ... FOR UPDATE` untuk mencegah *overselling* saat pesanan terjadi secara bersamaan (*high-concurrency checkout*). Stok dikembalikan otomatis jika pesanan kedaluwarsa atau dibatalkan. |
| **Payment Gateway Integration** | Integrasi end-to-end **DOKU Jokul Checkout** dengan verifikasi tanda tangan digital kriptografi **HMAC-SHA256**, penanganan webhook notifikasi asinkron, serta mode simulasi built-in untuk pengujian lokal tanpa internet. |
| **Type Safety & Developer DX** | **End-to-End Type Safety** melalui paket `@ecommerce/shared-types` monorepo. Kontrak model data dibagi bersama antara Go REST API dan frontend SvelteKit, menjamin nol *drift* tipe data. |
| **Code Quality & Testing** | Lolos **38 unit & component tests (Vitest)**, pemeriksaan statis TypeScript ketat (`svelte-check` 0 errors), dan `go vet` tanpa temuan keamanan/linting. |

---

## 🎯 Akun Demo Pengujian Cepat

Aplikasi telah dilengkapi seeder data bawaan. Anda dapat langsung menguji tanpa perlu mendaftar dari awal:

| Peran (Role) | Email | Password | Akses & Hak Istimewa |
|---|---|---|---|
| 🛒 **Pembeli (Buyer)** | `buyer@example.com` | `Password123` | Jelajahi katalog, keranjang belanja, checkout DOKU, wishlist, daftar alamat, pengaturan pembayaran, & mode layar. |
| 🏬 **Penjual (Seller)** | `seller@techgear.com` | `Password123` | Dashboard seller (`/dashboard`), manajemen stok produk realtime, tinjau pesanan masuk, dan proses fulfillment. |

> *Catatan: Halaman login juga memiliki tombol "1-Click Fill Demo Account" untuk mengisi kredensial secara instan.*

---

## 🏗️ Arsitektur Sistem

```
                                  ┌───────────────────────────────┐
                                  │      Client Web Browser       │
                                  │ (SvelteKit 5 + Tailwind v4)   │
                                  └──────────────┬────────────────┘
                                                 │ HTTP / REST
                                                 ▼
┌─────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                   Go Backend Microservice                                       │
│                                                                                                 │
│  ┌───────────────────────┐   ┌───────────────────────┐   ┌───────────────────────────────────┐  │
│  │   HTTP Handler Layer  │──▶│     Service Layer     │──▶│         Repository Layer          │  │
│  │   (Gin / Chi REST)    │   │ (Business Logic, ACID)│   │       (GORM + SQL Builder)        │  │
│  └───────────────────────┘   └───────────┬───────────┘   └─────────────────┬─────────────────┘  │
│                                          │                                 │                    │
└──────────────────────────────────────────┼─────────────────────────────────┼────────────────────┘
                                           │ Webhook / Checkout              │ SQL Queries
                                           ▼                                 ▼
                               ┌───────────────────────┐         ┌───────────────────────┐
                               │     DOKU Payment      │         │   PostgreSQL Engine   │
                               │  (Jokul Checkout v1)  │         │  (ACID Stock Locks)   │
                               └───────────────────────┘         └───────────────────────┘
```

---

## 🚀 Panduan Setup Awal (Developer Onboarding)

Panduan langkah demi langkah untuk mengkloning dan menjalankan proyek ini di lingkungan lokal Anda.

### 1. Kebutuhan Sistem (Prerequisites)
Pastikan dependensi berikut telah terpasang di komputer Anda:
- **Node.js** (v20+) atau **Bun** (v1.2+) *(direkomendasikan Bun)*
- **Go (Golang)** (v1.22 atau v1.24+)
- **PostgreSQL** (v15+) atau **Docker**
- **Git**

### 2. Kloning Repositori
```bash
git clone https://github.com/Nabilfikry/web-ecommerce.git
cd web-ecommerce
```

### 3. Instalasi Dependensi Monorepo
Instal seluruh paket workspace (Frontend, Backend dev tools, Shared Types):
```bash
# Menggunakan Bun (Sangat Cepat)
bun install

# Atau menggunakan NPM / PNPM
npm install
```

### 4. Konfigurasi Environment Variable (`.env`)

#### a. Backend Service (`apps/backend-service/.env`)
Salin file template environment backend:
```bash
cp apps/backend-service/.env.example apps/backend-service/.env
```

Pastikan variabel database dan kredensial DOKU sesuai:
```env
PORT=8000
DATABASE_URL=host=localhost user=admin password=admin dbname=ecommerce_db port=5432 sslmode=disable TimeZone=Asia/Jakarta
JWT_SECRET=ecommerce-super-secret-jwt-key-minimum-32-chars-long
REFRESH_TOKEN_SECRET=ecommerce-super-secret-refresh-token-key-32-chars
FRONTEND_URL=http://localhost:5173

# DOKU Payment Gateway (Jokul Checkout Sandbox)
DOKU_CLIENT_ID=BRN-0218-1794027396955
DOKU_SECRET_KEY=SK-Nq9xFKnpsSTRvbExexCz
DOKU_IS_PRODUCTION=false
DOKU_CALLBACK_URL=http://localhost:5173/orders
DOKU_SIMULATE_MODE=true
```

#### b. Frontend Web (`apps/web/.env`)
Salin file template environment frontend:
```bash
cp apps/web/.env.example apps/web/.env
```
Isi konfigurasi URL backend:
```env
PUBLIC_API_URL=http://localhost:8000/api/v1
```

---

### 5. Menjalankan Database PostgreSQL

Jika Anda menggunakan **Docker**, Anda dapat menyalakan instance PostgreSQL dengan satu perintah:
```bash
docker run --name tokokupedia-db \
  -e POSTGRES_USER=admin \
  -e POSTGRES_PASSWORD=admin \
  -e POSTGRES_DB=ecommerce_db \
  -p 5432:5432 \
  -d postgres:16-alpine
```

### 6. Migrasi & Seeding Data Awal

Jalankan migrasi tabel basis data dan populate data dummy (produk, toko, pengguna, alamat):
```bash
# Migrasi skema database
bun run db:migrate

# Seed data dummy lengkap (Pengguna, 12+ Produk, Microsoft Surface Pro 12, Transaksi)
bun run db:seed

# Opsi Cepat: Reset total & fresh seed dalam satu langkah
bun run db:fresh
```

---

### 7. Menjalankan Server Pengembangan (Dev Mode)

Jalankan frontend dan backend secara bersamaan melalui pipeline **Turborepo**:
```bash
bun run dev
```

Layanan akan aktif pada URL berikut:
- 🌐 **Frontend TokoKuPedia:** `http://localhost:5173`
- ⚙️ **Backend Go REST API:** `http://localhost:8000/api/v1`
- 📚 **Swagger / Contract Specs:** [`docs/API_CONTRACT.md`](./docs/API_CONTRACT.md)

*Ingin menjalankan layanan secara terpisah?*
```bash
# Menjalankan frontend web saja
bun run dev:web

# Menjalankan backend Go saja
bun run dev:backend
```

---

## 🧪 Pengujian & Jaminan Mutu (Testing & Quality)

Repositori ini menerapkan standar pengujian otomatis sebelum setiap build rilis:

```bash
# 1. Jalankan Unit & Component Tests (38 Vitest Tests)
bun run test

# 2. Type Checking Ketat (Svelte-Check TypeScript)
bun run check

# 3. Linter & Format Audit (Prettier & ESLint)
bun run lint

# 4. Build Verifikasi Produksi
bun run build
```

---

## 📂 Struktur Direktori Monorepo

```tree
web-ecommerce/
├── apps/
│   ├── web/                         # Frontend SvelteKit 5 + Tailwind v4
│   │   ├── src/
│   │   │   ├── lib/
│   │   │   │   ├── api/             # API Client dengan error fallback
│   │   │   │   ├── components/      # UI Tokopedia (Navbar, Footer, ProductCard, Logo, Toped)
│   │   │   │   ├── stores/          # Reactive stores Svelte 5 (auth, cart, toast)
│   │   │   │   └── utils/           # Helper format mata uang IDR & tanggal
│   │   │   └── routes/              # SvelteKit Routing (Home, PDP, Register, Cart, Settings)
│   │   └── package.json
│   │
│   └── backend-service/             # Backend REST Microservice (Go 1.24)
│       ├── cmd/
│       │   ├── server/              # Entry point HTTP daemon
│       │   └── migrate/             # CLI database migration & seeder
│       ├── internal/
│       │   ├── handler/             # HTTP controller & request validation
│       │   ├── service/             # Business logic & ACID stock locks
│       │   ├── repository/          # PostgreSQL database access layer
│       │   ├── database/            # Migration schema & comprehensive seeder
│       │   └── model/               # GORM entities & domain models
│       ├── pkg/                     # Reusable utilities (JWT, HMAC, Hash)
│       └── go.mod
│
├── packages/
│   └── shared-types/                # Shared TypeScript models (Product, Cart, Order, User)
├── docs/                            # Arsitektur, ERD Database, PRD, dan Panduan API
├── DESIGN.md                        # Pedoman Design System TokoKuPedia Unify
├── turbo.json                       # Konfigurasi Turborepo build pipeline
└── package.json                     # Monorepo root scripts & workspaces
```

---

## 🎨 Halaman & Alur UI TokoKuPedia yang Tersedia

1. **Homepage (`/`)**: Navigasi multi-tier Tokopedia, Kategori Flyout, Search bar dengan tag instan, banner carousel promosi, widget *Top Up & Tagihan*, countdown *Flash Sale 11rb*, dan feed rekomendasi produk.
2. **Pendaftaran & Masuk (`/register` & `/login`)**: Layout split-screen dengan ilustrasi maskot Toped, masuk Google, selector peran Pembeli/Penjual, dan disclaimer privasi.
3. **Product Detail Page (`/[shop]/[slug]` & `/products/[id]`)**: Tampilan 3-kolom persis Tokopedia PDP untuk demo produk **Microsoft Surface Pro 12 Copilot+ PC**, pemilihan varian dinamis, profil toko, estimasi ongkir, dan sticky buy box.
4. **Keranjang Belanja (`/cart`)**: Pengelompokan barang per toko, kupon promo diskon (`BEBASONGKIR`), stepper kuantitas, catatan untuk penjual, dan kalkulasi tagihan sticky.
5. **Daftar Alamat Pengiriman (`/user/settings/address`)**: Kelola alamat rumah/kantor, badge alamat utama, pinpoint lokasi, dan modal penambahan alamat baru.
6. **Metode Pembayaran (`/user/settings/payment`)**: Saldo GoPay & TokoKuPedia, penautan e-wallet OVO, kartu debit/kredit tersimpan, dan daftar virtual account.
7. **Rekening Bank Penarikan (`/user/settings/bank`)**: Kelola akun bank penarikan saldo/refund dengan batas maksimal 3 rekening dan proteksi nama KTP.
8. **Pengaturan Tampilan / Mode Layar (`/user/settings/mode-screen`)**: Switcher interaktif untuk mode Terang (*Light*), Gelap (*Dark*), atau Ikuti Sistem (*System*) dengan penyimpanan permanen `localStorage`.

---

## 📜 Lisensi & Kontributor

- **Pengembang:** [Naabil Fikry](https://github.com/Nabilfikry)
- **Repositori:** [github.com/Nabilfikry/web-ecommerce](https://github.com/Nabilfikry/web-ecommerce)
- **Lisensi:** MIT License — Dibuat untuk demonstrasi rekayasa perangkat lunak skala produksi (*production-grade software engineering*).
