# System Architecture & Technical Specification

**Project Architecture:** Monorepo (Turborepo + Bun)  
**Backend:** Go 1.26 (Modular Monolith, Standard Library / pgx)  
**Frontend:** SvelteKit (Svelte 5, Tailwind CSS v4, TypeScript, Storybook 10)  
**Shared Library:** `@ecommerce/shared-types` (TypeScript contracts)  
**Database:** PostgreSQL 15+  

---

## 1. High-Level Architecture

Sistem menggunakan pendekatan **Client-Server Architecture** dengan backend berbentuk **Modular Monolith**.

```mermaid
flowchart TD
    subgraph Client ["Client Tier (Browser)"]
        UI["SvelteKit Web Application (apps/web)<br/>• Svelte 5 Runes<br/>• Tailwind CSS v4<br/>• Storybook Component Library"]
    end

    subgraph MonorepoLib ["Contract Tier"]
        ST["Shared Types (packages/shared-types)<br/>TypeScript API Interfaces & DTOs"]
    end

    subgraph Server ["Server Tier (apps/backend-service)"]
        API["Go HTTP Server (:8080)<br/>net/http Router & Middleware"]
        subgraph Modules ["Modular Monolith Modules"]
            AuthMod["auth module"]
            UserMod["user module"]
            ProdMod["product module"]
            CartMod["cart module"]
            OrderMod["order module (ACID Tx)"]
        end
    end

    subgraph Data ["Data Tier"]
        DB[("PostgreSQL 15+<br/>Users, Products, Carts, Orders, Tokens")]
    end

    UI -->|HTTP / JSON REST API| API
    ST -.->|Type-Safe Imports| UI
    API --> Modules
    Modules -->|pgx / database/sql Connection Pool| DB
```

### 1.1 Architectural Decisions & Guarantees
1. **Frontend-Backend Decoupling:** Frontend tidak pernah berinteraksi langsung dengan database. Seluruh akses data dan mutasi state melewati REST API.
2. **Single Source of Truth for Calculations:** Frontend tidak dipercaya untuk kalkulasi finansial (`unit_price`, `subtotal`, `total_amount`). Backend menghitung ulang seluruh total belanja berdasarkan harga yang tersimpan di database pada saat transaksi.
3. **Strict Concurrency Protection:** Proses checkout melakukan penguncian baris produk (`SELECT ... FOR UPDATE OF products`) di dalam transaksi database PostgreSQL untuk mencegah *overselling* (stok minus akibat concurrent requests).
4. **Historical Price Snapshot:** Tabel `order_items` mencatat snapshot `product_name` dan `unit_price` saat pesanan tercipta. Revisi harga atau nama di katalog produk di masa depan tidak mempengaruhi data historis transaksi.
5. **Soft Delete Policy:** Penghapusan produk oleh penjual tidak menjalankan `DELETE` fisik pada database melainkan mengisi kolom `deleted_at`, menjaga integritas relasi foreign key pada tabel `order_items`.

---

## 2. Monorepo Project Structure

Struktur direktori aktual repositori:

```tree
e-commerce/
├── apps/
│   ├── backend-service/          # Go Backend Service
│   │   ├── cmd/
│   │   │   └── server/
│   │   │       └── main.go       # Entrypoint server Go
│   │   ├── configs/              # Konfigurasi environment (DB, JWT, Port)
│   │   ├── internal/             # Domain modules (Modular Monolith)
│   │   │   ├── auth/             # Handler, Service, Repository, DTO
│   │   │   ├── user/             # Profile & User Management
│   │   │   ├── product/          # Product CRUD & Catalog
│   │   │   ├── cart/             # Shopping Cart
│   │   │   ├── order/            # Checkout & Order Lifecycle (ACID)
│   │   │   ├── database/         # Postgres pool connection & migrations
│   │   │   ├── middleware/       # JWT Auth, RBAC, CORS, Recovery, Logger
│   │   │   └── shared/           # Response helpers, custom errors, validator
│   │   ├── migrations/           # File SQL migrasi database (.up.sql / .down.sql)
│   │   ├── pkg/                  # Utilitas publik (hash, token, uuid)
│   │   ├── go.mod
│   │   ├── go.sum
│   │   └── package.json          # Turbo workspace script untuk backend
│   │
│   └── web/                      # SvelteKit Frontend
│       ├── src/
│       │   ├── lib/
│       │   │   ├── api/          # API Client (auth, products, cart, orders)
│       │   │   ├── components/   # Reusable UI components
│       │   │   │   ├── ui/       # Button, Input, Modal, Badge, Toast
│       │   │   │   ├── product/  # ProductCard, ProductGrid, PriceTag
│       │   │   │   ├── cart/     # CartDrawer, CartItemRow, CartSummary
│       │   │   │   └── order/    # OrderTimeline, OrderCard, StatusBadge
│       │   │   ├── stores/       # Svelte 5 state / runes
│       │   │   └── types/        # Extended frontend-specific types
│       │   ├── routes/
│       │   │   ├── (auth)/       # /login, /register
│       │   │   ├── (buyer)/      # /products, /products/[id], /cart, /checkout, /orders
│       │   │   ├── (seller)/     # /dashboard, /dashboard/products, /dashboard/orders
│       │   │   ├── +layout.svelte
│       │   │   └── +page.svelte  # Landing page
│       │   ├── stories/          # Storybook component stories
│       │   └── app.html
│       ├── package.json
│       ├── vite.config.ts
│       └── svelte.config.js
│
├── packages/
│   └── shared-types/             # Shared TypeScript interfaces
│       ├── src/
│       │   └── index.ts          # API response, User, Product, Cart, Order types
│       ├── package.json
│       └── tsconfig.json
│
├── docs/                         # Project Documentation
│   ├── PRD.md                    # Product Requirements Document
│   ├── DATABASE_ERD.md           # ERD & Database Design
│   ├── API_CONTRACT.md           # REST API Contract & Specifications
│   ├── ARCHITECTURE.md           # This document
│   ├── FEATURE_DEVELOPMENT_GUIDE.md # Panduan pengerjaan penambahan fitur
│   └── README.md                 # Index dokumentasi
│
├── bin/                          # Output binary Go
├── docker-compose.yml            # PostgreSQL & Local tooling
├── go.work                       # Go workspace configuration
├── turbo.json                    # Turborepo task pipeline
├── package.json                  # Root monorepo scripts
└── bun.lock                      # Bun lockfile
```

---

## 3. Go Backend Modular Monolith Design

Setiap modul di dalam `apps/backend-service/internal/` menerapkan prinsip **Separation of Concerns** dengan pembagian layer yang konsisten:

```mermaid
flowchart LR
    Request([HTTP Request]) --> Router[Mux / Router]
    Router --> Middleware[Auth / RBAC / CORS]
    Middleware --> Handler[Handler Layer<br/>Decode & Validate HTTP]
    Handler --> Service[Service Layer<br/>Business Logic & Transactions]
    Service --> Repo[Repository Layer<br/>SQL & DB Queries]
    Repo --> DB[(PostgreSQL)]
    
    Handler -.-> DTO[DTO & Response Formatter]
    Service -.-> Model[Domain Models]
```

### 3.1 Tanggung Jawab Tiap Layer:
1. **Handler Layer (`handler.go`):**
   - Mengambil parameter HTTP (path params, query params, request body).
   - Melakukan validasi format dasar.
   - Memanggil Service.
   - Mengembalikan HTTP Response menggunakan standard response JSON (`shared/response`).
2. **Service Layer (`service.go`):**
   - Menjalankan seluruh aturan bisnis (misal: verifikasi kuantitas vs stok, memastikan seller sama, transisi status pesanan).
   - Mengelola koordinasi transaksi database (`BeginTx`, `Commit`, `Rollback`).
   - Tidak boleh mengimpor paket HTTP (`net/http`).
3. **Repository Layer (`repository.go`):**
   - Menjalankan query SQL mentah / prepared statement menggunakan driver PostgreSQL (`pgx` atau `database/sql`).
   - Menerima parameter `context.Context` dan optional `*sql.Tx` untuk transaksi.
   - Memetakan hasil baris tabel menjadi domain `model.go`.
4. **DTO Layer (`dto.go`):**
   - Struktur data khusus request payload dan response payload yang ditransfer melalui jaringan.

---

## 4. Frontend Architecture (SvelteKit + Svelte 5)

Frontend dibangun dengan **SvelteKit** memanfaatkan fitur reaktivitas modern **Svelte 5 Runes** (`$state`, `$derived`, `$props`, `$effect`).

### 4.1 Organisasi Routing & Akses:
- **Public Routes:**
  - `/`: Landing page & katalog pembuka.
  - `/products`: Katalog produk publik.
  - `/products/[id]`: Detail produk.
  - `/login`: Form masuk akun.
  - `/register`: Form pendaftaran akun baru.
- **Buyer Routes (`(buyer)` route group):**
  - `/cart`: Keranjang belanja interaktif.
  - `/checkout`: Form pengisian alamat & ringkasan pembayaran.
  - `/orders`: Riwayat pesanan pembeli.
  - `/orders/[id]`: Detail dan pelacakan status pesanan.
- **Seller Routes (`(seller)` route group):**
  - `/dashboard`: Ringkasan toko & statistik sederhana.
  - `/dashboard/products`: CRUD produk dan pengaturan stok.
  - `/dashboard/orders`: Daftar pesanan masuk dan pemrosesan status.

### 4.2 Route Guards:
- Menggunakan `hooks.server.ts` pada SvelteKit untuk memverifikasi session cookie/token JWT sebelum render halaman.
- Halaman Seller secara otomatis me-redirect Buyer ke `/` jika mencoba mengakses `/dashboard`.

---

## 5. Security & Authentication Architecture

```mermaid
sequenceDiagram
    autonumber
    actor User as Client (SvelteKit)
    participant API as Backend (Go API)
    participant DB as PostgreSQL Database

    User->>API: POST /api/v1/auth/login { email, password }
    API->>DB: Query user by LOWER(email)
    DB-->>API: User record + password_hash
    API->>API: Verify password with Argon2id/Bcrypt
    API->>API: Generate Access Token (JWT, 15 min expiry)
    API->>API: Generate Refresh Token (Random string, 7 days)
    API->>DB: Insert hashed refresh token into refresh_tokens table
    API-->>User: 200 OK + Body: { access_token } + Set-Cookie: refresh_token (HttpOnly)

    Note over User, API: Request terproteksi berikutnya:
    User->>API: GET /api/v1/cart (Header: Authorization: Bearer <access_token>)
    API->>API: Validate JWT signature & claims (UserID, Role)
    API-->>User: 200 OK + Cart Data
```

### 5.1 Security Checklist:
- **Password Hashing:** Menggunakan **Argon2id** (atau Bcrypt dengan minimum cost factor 10).
- **Token Lifespan:** Access token berlaku 15 menit. Refresh token berlaku 7 hari.
- **Cookie Flags:** `HttpOnly`, `Secure` (saat HTTPS/Production), `SameSite=Strict`, `Path=/api/v1/auth`.
- **SQL Injection Prevention:** 100% query database menggunakan parameter binding (`$1, $2, ...`). Tidak ada perangkaian string (*string concatenation*) pada query SQL.
- **CORS Policy:** Whitelist domain frontend (misal `http://localhost:5173` saat pengembangan lokal).

---

## 6. Local Development Environment (Docker Compose)

File `docker-compose.yml` di-root project menyediakan instance PostgreSQL yang konsisten untuk pengembangan lokal:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:16-alpine
    container_name: ecommerce-postgres
    restart: always
    environment:
      POSTGRES_USER: ecommerce_user
      POSTGRES_PASSWORD: ecommerce_password
      POSTGRES_DB: ecommerce_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

Environment variables backend (`apps/backend-service/.env`):
```env
PORT=8080
DATABASE_URL=postgres://ecommerce_user:ecommerce_password@localhost:5432/ecommerce_db?sslmode=disable
JWT_SECRET=super-secret-key-change-in-production-min-32-chars
REFRESH_TOKEN_SECRET=another-super-secret-key-for-refresh-token
FRONTEND_URL=http://localhost:5173
```

