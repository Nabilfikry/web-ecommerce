# Design System: TokoKuPedia (Tokopedia Unify UI Clone)
**Project ID:** tokokupedia-v1  
**Target Platform:** Web (Desktop, Tablet, Mobile Responsive)  
**Implementation Stack:** SvelteKit (Svelte 5) · Tailwind CSS v4 · TypeScript · Storybook  
**Design Reference:** Tokopedia Unify Design System (www.tokopedia.com)

---

## 1. Visual Theme & Atmosphere

**TokoKuPedia** adopts the iconic **Tokopedia Unify** visual design language: an approachable, trustworthy, high-clarity Indonesian e-commerce marketplace aesthetic centered around the signature **Tokopedia Green (`#00AA5B`)**, crisp white surfaces (`#FFFFFF`), light canvas grounds (`#F5F6F7`), and cheerful mascot cues ("Toped" owl bag).

### Key Aesthetic Principles
- **Vibrant Yet Structured:** Tokopedia Green (`#00AA5B` / `#03AC0E`) serves as the core brand beacon across buttons, badges, links, and accents.
- **High-Density Commerce Grids:** High information density optimized for rapid scanning: product thumbnails with square 1:1 aspect ratio, Official Store green shield badges, bold rupiah prices, discount percentages, ratings, and sold counts.
- **Familiar Indonesian E-Commerce Conventions:** Top App promotional banner, multi-tier navigation with "Kategori" flyout, prominent search input with instant tag suggestions, "Top Up & Tagihan" utility widget, and location selector ("Dikirim ke Jakarta Pusat").
- **Multi-Screen Mode Support:** Fully supports Light Mode, Dark Mode (`#121212`), and System Mode with instant persistence via localStorage.

---

## 2. Color Palette & Roles

```
┌────────────────────────────────────────────────────────────────────────┐
│                        TOKOKUPEDIA COLOR TOKENS                        │
│                                                                        │
│  [ Canvas #F5F6F7 ]   [ Surface #FFFFFF ]   [ Primary #00AA5B ]        │
│  [ Border #E5E7E9 ]   [ Muted #6D7588   ]   [ Hover   #03AC0E ]        │
│  [ Discount #EE4D2D ] [ Rating #FFAB00  ]   [ Soft Tint #F0FDF4 ]      │
└────────────────────────────────────────────────────────────────────────┘
```

### 2.1 Brand & Action Colors
- **Tokopedia Green GN500** (`#00AA5B`) — Primary action color for CTA buttons ("Daftar", "Beli", "+ Keranjang", "Checkout"), active navigation pills, and brand logo.
- **Tokopedia Green Dark GN600** (`#03AC0E` / `#008848`) — Hover and focus-visible state for primary interactive elements.
- **Tokopedia Mint Tint GN50** (`#F0FDF4` / `#E8F5E9`) — Soft background for active variant chips, success alerts, and secondary action hover states.

### 2.2 Foundation & Neutral Colors
- **Canvas Base Gray NN50** (`#F5F6F7`) — Soft application background reducing glare.
- **Pure Surface White NN0** (`#FFFFFF`) — Cards, modals, drawers, and header containers.
- **Foreground Obsidian NN950** (`#212121`) — Primary typography for headlines, product names, and tabular prices.
- **Muted Slate NN600** (`#6D7588`) — Secondary text for metadata, sellers, ratings, and breadcrumbs.
- **Hairline Border NN300** (`#E5E7E9`) — 1px clean architectural boundaries.

### 2.3 Semantic & Commercial Status Colors
- **Discount Red RN500** (`#EE4D2D` / `#F94D63`) — Discount tags (`15%`), flash sale banner badges, and cart notification counts.
- **Rating Star Yellow YN500** (`#FFAB00`) — Star ratings and points.
- **Bebas Ongkir Emerald** (`#10B981`) — Free shipping badges and logistics tags.

---

## 3. Typographic Hierarchy

**Font Family:** `'Open Sauce One'`, `'Plus Jakarta Sans'`, system-ui, -apple-system, sans-serif.

| Role | Size | Weight | Line Height | Usage |
|---|---|---|---|---|
| **H1 Display** | 24px - 32px | Bold (700/800) | 1.25 | Section headers, Hero titles, PDP title |
| **Price Bold** | 18px - 30px | Extra Bold (800) | 1.1 | Tabular figures (`tabular-nums font-extrabold`) |
| **Card Title** | 13px - 14px | Medium / SemiBold | 1.35 | Product card name (`line-clamp-2`) |
| **Body Meta** | 12px - 13px | Regular / Medium | 1.4 | Seller name, location, reviews, tabs |
| **Micro Badge** | 10px - 11px | Bold (700) | 1.2 | Official Store, Discount %, Stock warning |

---

## 4. Replicated Tokopedia Pages

### 4.1 Homepage (`/`)
- **Top Utility Header:** App download banner (`Gratis Ongkir + Banyak Promo di aplikasi >`), utility links (Tentang, Mulai Berjualan, Promo, TokoKuPedia Care).
- **Main Sticky Header:** Green TokoKuPedia logo with owl mascot, "Kategori" dropdown with icon grid, expansive search bar with instant tag suggestions, shopping cart badge, and auth buttons (Masuk & Daftar).
- **Hero Promotional Banner Carousel:** Auto-sliding high-impact campaign slides with indicator dots and prev/next controls.
- **Kategori Populer & Top Up Tagihan:** Multi-tab widget for Pulsa, Paket Data, Listrik PLN, and Roaming with instant nominal selector.
- **Category Quick Chips:** Horizontal scrollable category pills with icons.
- **Flash Sale Serba 11rb:** Red promotion banner with active live countdown timer (`hh:mm:ss`).
- **Rekomendasi Untuk Anda Feed:** 6-column grid of Tokopedia product cards with Official Store badges, star ratings, and location tags.

### 4.2 Register Page (`/register`)
- **Header:** Centered TokoKuPedia wordmark with logo.
- **Two-Column Split Layout:**
  - Left: Toped Mascot illustration in front of marketplace storefront with parachutes and slogan: *"Jual Beli Mudah Hanya di TokoKuPedia"*.
  - Right: Floating card with *"Daftar Sekarang"*, Google authentication button, separator divider, input fields, role selector (Pembeli/Penjual), solid green button *"Daftar"*, and legal disclaimers.

### 4.3 Product Detail Page (PDP) (`/[shop]/[slug]` & `/products/[id]`)
- **Breadcrumbs:** Full hierarchical breadcrumbs trail.
- **Three-Column Desktop Layout:**
  - Column 1: Gallery with 1:1 image viewport, Official Store badge, and thumbnail carousel with active green outline ring.
  - Column 2: Product title, rating stars (`4.9 (45 rating)`), sold count, large bold price with discount percentage, interactive variant chips (Color & Storage), seller store card, and tabs (Detail, Spesifikasi, Info Penting).
  - Column 3: Sticky Buy Box card with quantity stepper, stock availability, subtotal calculation, `+ Keranjang` (green outline), `Beli Langsung` (solid green), and quick actions (Chat, Wishlist, Share).

### 4.4 Cart Page (`/cart`)
- **Minimal Header:** TokoKuPedia logo + "Keranjang" indicator.
- **Left Column:** Store grouping card with store checkbox, Official Store shield, item checkbox, thumbnail, variant badge, price, note input, and quantity stepper.
- **Right Column:** Sticky "Ringkasan Belanja" card with coupon code input ("BEBASONGKIR"), item subtotal, discounts, free shipping badge, and solid green checkout button.

### 4.5 User Settings (`/user/settings/...`)
- **Shared Settings Layout:**
  - Left Sidebar: User profile avatar with photo change trigger, verified badge, TokoKuPedia balance, GoPay linked balance, and navigation menu links.
- **Address Management (`/user/settings/address`):**
  - Search address bar, "+ Tambah Alamat Baru" modal form, primary address green badge, recipient details, and actions (Jadikan Utama, Ubah, Hapus).
- **Payment Methods (`/user/settings/payment`):**
  - Linked e-wallets (GoPay, Saldo TokoKuPedia, OVO), saved debit/credit cards (Visa/Mastercard) with "+ Tambah Kartu" modal, and favorite virtual accounts.
- **Bank Accounts (`/user/settings/bank`):**
  - Registered bank accounts list (BCA, Mandiri, BRI), primary withdrawal badge, security compliance disclaimer, and "+ Tambah Rekening" modal.
- **Screen Mode / Mode Tampilan (`/user/settings/mode-screen`):**
  - Visual selection cards for Terang (Light Mode), Gelap (Dark Mode), and Ikuti Sistem (System Default) with instant reactive theme switching and localStorage persistence.
