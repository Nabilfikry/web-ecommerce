# Design System: Simple E-Commerce Platform
**Project ID:** simple-ecommerce-v1  
**Target Platform:** Web (Desktop, Tablet, Mobile Responsive)  
**Implementation Stack:** SvelteKit (Svelte 5) · Tailwind CSS v4 · TypeScript · Storybook  

---

## 1. Visual Theme & Atmosphere

Simple E-Commerce Platform adopts a **Crafted Modern Utilitarian** aesthetic—a balance of Scandinavian restraint, Swiss grid precision, and tactile retail clarity. The visual language prioritizes **unobtrusive product presentation, effortless scanability, and instant buyer trust**.

The atmosphere is **clean, grounded, and airy**. Rather than relying on gratuitous decoration, flashy neon orbs, or trendy multi-colored gradients, the interface commands attention through **flawless typographic hierarchy, disciplined whitespace, crisp hairline borders, and purposeful contrast**. Every UI surface feels physical and intentional: cards sit stable on a calm ground plane, interactive elements possess tactile micro-feedback, and transactional flows (such as cart and checkout) maintain an uncluttered editorial calm that eliminates cognitive friction.

### Key Aesthetic Principles
- **Product-First Canvas:** Surfaces and layout frames remain neutral and quiet so that product photography and information take center stage.
- **Purposeful Tactility:** Flat by default with subtle hairline borders; depth is introduced sparingly via whisper-soft diffused shadows only upon hover or elevation.
- **High-Density Legibility:** Clear typographic rhythm with tabular numerals for prices and quantities, ensuring instant mental calculation.
- **Restraint Over Excess:** Zero decorative gradient blobs, zero gratuitous glassmorphism, zero emoji-as-icons. Every visual detail must pass the purpose test.

---

## 2. Color Palette & Roles

The color system enforces a disciplined **2 Core + 1 Deliberate Accent** philosophy (Anti-Slop R-29). The palette is optimized for minimum 4.5:1 WCAG contrast across all states.

```
┌────────────────────────────────────────────────────────────────────────┐
│                              PALETTE MAP                               │
│                                                                        │
│  [ Canvas #F8FAFC ]   [ Surface #FFFFFF ]   [ Primary #0F172A ]        │
│  [ Border #E2E8F0 ]   [ Muted #64748B   ]   [ Accent  #D97706 ]        │
│  [ Success #059669]   [ Error   #DC2626 ]   [ Info    #2563EB ]        │
└────────────────────────────────────────────────────────────────────────┘
```

### 2.1 Core Neutral Foundation
- **Canvas Base Gray** (`#F8FAFC`) — Primary application background canvas. Provides a soft, glare-free warm-slate tone that reduces eye fatigue compared to harsh pure white.
- **Pure Surface White** (`#FFFFFF`) — Surface layer for product cards, navigation bar, modal dialogs, and slide-over drawers. Creates clean contrast against the canvas.
- **Hairline Border Slate** (`#E2E8F0`) — 1px structural borders and dividers. Creates crisp architectural boundaries without visual heaviness.
- **Deep Obsidian Charcoal** (`#0F172A`) — Primary foreground color for headlines, product titles, bold prices, and primary action buttons. Delivers an authoritative 14.5:1 contrast ratio.
- **Muted Slate Gray** (`#64748B`) — Secondary text for descriptions, seller labels, metadata, and placeholder text. Maintains clear hierarchy without dropping below 4.5:1 contrast.

### 2.2 Brand & Interactive Colors
- **Tailored Deep Slate** (`#0F172A`) — **Primary Action Fill.** Used for primary CTA buttons ("Beli Sekarang", "Checkout", "Tambah Produk"), active tab pills, and key interactive focal points.
- **Warm Spiced Ochre** (`#D97706`) — **Deliberate Accent.** Used selectively for checkout highlights, active step indicators, and "Best Seller" / "Promo" accents. Used with extreme restraint to preserve its visual punch.
- **Deep Slate Hover** (`#1E293B`) — Subtle interactive hover state for primary controls.

### 2.3 Semantic & Lifecycle Status Colors
Reserved exclusively for system feedback, order lifecycle states, and inventory indicators:
- **Success Emerald** (`#059669`) / Tint (`#ECFDF5`) — Order status `completed`, in-stock badges, checkout confirmation.
- **Warning / Pending Amber** (`#D97706`) / Tint (`#FFFBEB`) — Order status `pending`, low stock warnings (`stok < 5`).
- **Processing Cyan-Blue** (`#0284C7`) / Tint (`#F0F9FF`) — Order status `confirmed` & `processing`.
- **Shipped Indigo** (`#4F46E5`) / Tint (`#EEF2FF`) — Order status `shipped`.
- **Destructive Crimson** (`#DC2626`) / Tint (`#FEF2F2`) — Order status `cancelled` & `rejected`, out of stock notices, delete product confirmation.

---

## 3. Typography Rules

**Primary Font Family:** `Plus Jakarta Sans`, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif.  
**Characteristics:** Modern geometric sans-serif engineered for digital screens. Features wide apertures, clean geometric proportions, and humanist warmth that ensure extreme readability at small sizes.

### 3.1 Typographic Hierarchy
| Role | Size | Weight | Tracking | Line Height | Penggunaan |
|---|---|---|---|---|---|
| **Display Title (H1)** | 32px / 2.0rem | Bold (700) | `-0.025em` | 1.25 | Hero banner, judul halaman utama |
| **Section Header (H2)** | 24px / 1.5rem | Semi-Bold (600) | `-0.02em` | 1.3 | Judul katalog, dashboard section |
| **Component Header (H3)**| 18px / 1.125rem | Semi-Bold (600) | `-0.015em` | 1.35 | Nama produk pada card, judul modal |
| **Body Standard** | 15px / 0.9375rem| Regular (400) | `normal` | 1.6 | Deskripsi produk, instruksi form |
| **Body Small / Meta** | 13px / 0.8125rem| Medium (500) | `+0.01em` | 1.5 | Label seller, breadcrumb, tanggal |
| **Price / Number** | 16px - 22px | Semi-Bold (600) | `tight` | 1.2 | `font-variant-numeric: tabular-nums` |
| **Badge / Micro Tag** | 11px / 0.6875rem| Semi-Bold (600) | `+0.03em` | 1.2 | Status pesanan, stock badge (uppercase) |

### 3.2 Numeric Formatting Convention
Seluruh harga produk dan angka kalkulasi mata uang wajib menggunakan fitur **Tabular Figures** (`tabular-nums font-mono-numbers`) agar angka berjajar sejajar secara vertikal pada keranjang dan invoice transaksi:
```css
font-feature-settings: "tnum";
font-variant-numeric: tabular-nums;
```

---

## 4. Component Stylings

### 4.1 Buttons
- **Corner Roundness:** Subtly rounded corners (`8px` / `rounded-lg`). Approachable and contemporary, never full pills.
- **Height & Touch Target:** Minimum `44px` on mobile/desktop (`py-2.5 px-5`).
- **Primary CTA:**
  - Background: Deep Obsidian (`#0F172A`)
  - Text: Pure White (`#FFFFFF`), font weight 600
  - Hover: Background shifts to `#1E293B`, subtle translateY `-1px`
  - Active: Scale `0.98`
  - Focus Ring: `focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:ring-[#0F172A]`
- **Secondary / Outlined:**
  - Background: Transparent
  - Border: 1px solid Hairline Slate (`#E2E8F0`)
  - Text: Deep Obsidian (`#0F172A`), font weight 500
  - Hover: Background `#F8FAFC`, border `#CBD5E1`
- **Destructive:**
  - Background: Crimson Red (`#DC2626`) with Pure White text, or subtle red tint background (`#FEF2F2`) with red text for secondary danger actions.

### 4.2 Cards & Containers
- **Border:** Crisp 1px solid Hairline Slate (`#E2E8F0`).
- **Corner Style:** Gently rounded corners (`12px` / `rounded-xl`).
- **Background:** Pure White (`#FFFFFF`).
- **Shadow Strategy:**
  - Default: **Flat** (no shadow, crisp 1px border definition).
  - Hover (Interactive Cards): Whisper-soft diffused elevation (`box-shadow: 0 4px 16px -2px rgba(15, 23, 42, 0.06)`).
  - Never apply heavy, blurry drop shadows that detach the card from the page plane.

### 4.3 Form Inputs & Controls
- **Border:** 1px solid `#CBD5E1` (Default), smooth transition on interaction.
- **Background:** Pure White (`#FFFFFF`).
- **Corner Style:** Subtly rounded (`8px` / `rounded-lg`).
- **Padding:** `10px 14px` (touch-friendly height ~44px).
- **Focus State:** Border `#0F172A` with `ring-1 ring-[#0F172A]`. Never hide focus rings.
- **Labels:** Visible, explicit labels above the input field in Deep Obsidian (`#0F172A`), text size 13px, font weight 500. Never rely on placeholders as labels.

### 4.4 Product Cards (Core E-Commerce Pattern)
```
┌──────────────────────────────────────────┐
│  [ Aspect 1:1 Image with subtle zoom ]   │
│  [ Floating Stock Badge top-left     ]   │
├──────────────────────────────────────────┤
│  Category / Seller Name  (12px Slate)    │
│  Product Title           (15px SemiBold) │
│  Rp 450.000              (16px Bold)     │
│  [ + Keranjang (Button)              ]   │
└──────────────────────────────────────────┘
```
- **Image Container:** Aspect ratio 1:1 (square), background `#F1F5F9`. Image styled with `object-cover` and subtle `transition-transform duration-300 ease-out group-hover:scale-102`.
- **Badges:** Soft pill badge positioned at top-left (`bg-white/90 backdrop-blur-xs border border-slate-200 text-xs px-2.5 py-1`).
- **Body Padding:** `16px` (`p-4`).
- **Price Tag:** Displayed in prominent bold font with Indonesian currency formatting: `Rp 450.000`.

### 4.5 Navigation & Header
- **Layout:** Sticky top header with solid Pure White background, bordered below with 1px Hairline Slate (`border-b border-slate-200`).
- **Height:** 64px (`h-16`).
- **Search Bar:** Centered search input with instant magnifying glass icon and clear button.
- **Cart Button:** Includes an active count badge in Warm Spiced Ochre (`#D97706`) with pure white text when `item_count > 0`.

### 4.6 Order Status Badges
Status pesanan disajikan dengan pill badge lembut dengan teks kontras tinggi:
- `pending`: Background `#FEF3C7`, Text `#92400E` (Menunggu Konfirmasi)
- `confirmed`: Background `#E0F2FE`, Text `#0369A1` (Dikonfirmasi)
- `processing`: Background `#DBEAFE`, Text `#1E40AF` (Sedang Diproses)
- `shipped`: Background `#E0E7FF`, Text `#3730A3` (Sedang Dikirim)
- `completed`: Background `#D1FAE5`, Text `#065F46` (Selesai)
- `cancelled`: Background `#F1F5F9`, Text `#475569` (Dibatalkan)
- `rejected`: Background `#FEE2E2`, Text `#991B1B` (Ditolak)

---

## 5. Layout & Responsive Principles

### 5.1 Grid & Max-Width
- **Max Content Width:** `1280px` (`max-w-7xl`) centered with responsive horizontal padding.
- **Horizontal Page Margins:**
  - Mobile (<640px): `16px` (`px-4`)
  - Tablet (640px - 1024px): `24px` (`px-6`)
  - Desktop (>1024px): `32px` (`px-8`)
- **Product Grid Columns:**
  - Mobile (<640px): 1 or 2 columns (`grid-cols-2 gap-3`)
  - Tablet (640px - 1024px): 3 columns (`grid-cols-3 gap-5`)
  - Desktop (>1024px): 4 columns (`grid-cols-4 gap-6`)

### 5.2 Whitespace & Vertical Rhythm
- **Micro Spacing:** 4px, 8px, 12px (antara icon dan label, badge padding).
- **Component Spacing:** 16px, 24px (antara field form, padding card).
- **Section Spacing:** 48px to 64px (`py-12` s/d `py-16`) menciptakan ruang nafas lega antar seksi halaman.

### 5.3 Touch Targets & Mobile Ergonomics
- Seluruh elemen interaktif (tombol, input, link navigasi, switch kuantitas) wajib memiliki **bounding box sentuh minimal 44x44px** pada viewport mobile.
- Tombol kuantitas keranjang (`+` dan `-`) berukuran minimal 36x36px dengan padding sentuh 44px untuk mencegah salah tekan (*fat-finger error*).

---

## 6. Anti-Slop & Quality Enforcement (Guardrails)

Patuhi aturan tegas berikut saat membuat komponen dan halaman:

| Larangan (Slop Tells) | Mengapa Dilarang | Solusi Sesuai Design System |
|---|---|---|
| **Gradient Ungu-Biru AI** | Menandakan template generik tanpa identitas produk | Gunakan solid canvas `#F8FAFC` dan Deep Slate `#0F172A` |
| **Glassmorphism Berlebihan** | Mengaburkan hierarki visual dan memberatkan rendering | Gunakan solid Pure White `#FFFFFF` dengan 1px border |
| **Tombol & Card Serba Pil** | Menghilangkan karakter bentuk antar elemen | Terapkan radius terukur: `rounded-lg` (8px) untuk tombol, `rounded-xl` (12px) untuk card |
| **Shadow Tebal Melayang** | Membuat halaman terasa mengambang dan membingungkan mata | Card flat secara default; bayangan tipis hanya muncul saat hover |
| **Emoji Sebagai Icon** | Terlihat amatir dan inkonsisten antar platform OS | Wajib gunakan icon SVG vektor (Lucide Svelte / Heroicons) |
| **Background Grid / Dot Blueprint** | Tekstur palsu yang tidak relevan dengan e-commerce | Biarkan kanvas bersih dan lapang (white space jujur) |
| **Kalkulasi Harga di Frontend** | Celah keamanan dan manipulasi data belanja | Seluruh harga, subtotal, dan total dihitung oleh backend Go |

---

## 7. Tailwind CSS v4 Semantic Tokens

Dalam Tailwind CSS v4 (`apps/web/src/routes/layout.css`), definisikan variabel desain berikut:

```css
@import 'tailwindcss';
@plugin '@tailwindcss/forms';
@plugin '@tailwindcss/typography';

@theme {
  /* Color Tokens */
  --color-canvas: #F8FAFC;
  --color-surface: #FFFFFF;
  --color-border: #E2E8F0;
  --color-border-hover: #CBD5E1;

  --color-foreground: #0F172A;
  --color-foreground-muted: #64748B;

  --color-primary: #0F172A;
  --color-primary-hover: #1E293B;
  --color-primary-contrast: #FFFFFF;

  --color-accent: #D97706;
  --color-accent-hover: #B45309;
  --color-accent-contrast: #FFFFFF;

  --color-success: #059669;
  --color-success-bg: #ECFDF5;
  --color-warning: #D97706;
  --color-warning-bg: #FFFBEB;
  --color-error: #DC2626;
  --color-error-bg: #FEF2F2;

  /* Typography */
  --font-sans: 'Plus Jakarta Sans', system-ui, -apple-system, sans-serif;

  /* Elevation */
  --shadow-subtle: 0 1px 3px 0 rgba(15, 23, 42, 0.05);
  --shadow-hover: 0 4px 16px -2px rgba(15, 23, 42, 0.06);

  /* Radius */
  --radius-sm: 6px;
  --radius-md: 8px;
  --radius-lg: 12px;
  --radius-full: 9999px;
}
```

---

## 8. Prompting & Screen Generation Guidelines

Saat menggunakan alat bantu AI, Stitch, atau membuat halaman baru, gunakan deskriptor desain konsisten berikut:

- **Deskripsi Suasana:** *"A clean, crafted minimalist e-commerce storefront with airy whitespace, structured 4-column product grid, and tactile deep slate accents."*
- **Deskripsi Tombol:** *"Primary button in Deep Obsidian (#0F172A) with subtly rounded corners (8px), white text, and crisp focus ring."*
- **Deskripsi Card:** *"Product card on pure white background with a 1px hairline border (#E2E8F0), 1:1 aspect ratio image container, and whisper-soft lift shadow on hover."*
- **Deskripsi Input:** *"Form input with clear top label, subtle 8px radius, 1px border (#CBD5E1), and deep slate focus highlight."*
- **Deskripsi Status:** *"Pastel tinted pill badge with high-contrast text indicating order status."*

