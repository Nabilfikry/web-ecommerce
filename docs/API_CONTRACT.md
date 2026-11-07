# REST API v1 Specification & Contract

**Base URL:** `http://localhost:8080`  
**API Prefix:** `/api/v1`  
**Protocol:** HTTP/1.1 or HTTP/2 over JSON  
**Authentication:** JWT Bearer Token (`Authorization: Bearer <token>`) & `HttpOnly` Refresh Cookie  
**Time Format:** ISO 8601 UTC (`YYYY-MM-DDTHH:mm:ssZ`)  

---

## 1. General Conventions

### 1.1 HTTP Headers
- **Client Request:**
  - `Content-Type: application/json` (untuk request yang memiliki body)
  - `Accept: application/json`
  - `Authorization: Bearer <access_token>` (untuk endpoint terproteksi)
- **Cookies:**
  - `refresh_token=<token>; HttpOnly; Secure; SameSite=Strict; Path=/api/v1/auth`

### 1.2 Standard Success Response
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": {},
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 50,
    "total_pages": 5
  }
}
```
*(Properti `meta` bersifat opsional, hanya muncul pada endpoint yang mendukung paginasi).*

### 1.3 Standard Error Response
```json
{
  "success": false,
  "message": "Validation failed on the provided input",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": [
      {
        "field": "price",
        "message": "Price must be greater than or equal to 0"
      }
    ]
  }
}
```

### 1.4 Error Code Glossary
| Error Code | HTTP Status | Keterangan |
|---|:---:|---|
| `VALIDATION_ERROR` | 400 | Payload request tidak lolos validasi tipe/format/panjang |
| `INVALID_CREDENTIALS` | 401 | Email atau password tidak sesuai saat login |
| `UNAUTHORIZED` | 401 | Token Bearer tidak disertakan atau format invalid |
| `TOKEN_EXPIRED` | 401 | Masa berlaku access token JWT telah habis |
| `FORBIDDEN` | 403 | Role pengguna tidak memiliki hak mengakses resource ini |
| `RESOURCE_NOT_FOUND` | 404 | Entitas data yang dicari tidak ditemukan |
| `EMAIL_ALREADY_EXISTS` | 409 | Email sudah digunakan akun lain |
| `INSUFFICIENT_STOCK` | 409 | Stok produk tidak mencukupi kuantitas checkout |
| `CART_SELLER_MISMATCH` | 409 | Upaya checkout/menambahkan produk dari seller yang berbeda |
| `INVALID_ORDER_STATUS` | 409 | Status pesanan tidak dapat ditransisikan ke status yang diminta |
| `INTERNAL_SERVER_ERROR`| 500 | Terjadi kegagalan sistem internal pada server |

---

## 2. API Authorization Matrix

| Endpoint | Method | Public | Buyer | Seller | Catatan Kepemilikan |
|---|:---:|:---:|:---:|:---:|---|
| `/auth/register` | POST | ✓ | ✓ | ✓ | Terbuka untuk umum |
| `/auth/login` | POST | ✓ | ✓ | ✓ | Terbuka untuk umum |
| `/auth/refresh` | POST | ✓ | ✓ | ✓ | Memerlukan cookie refresh token |
| `/auth/logout` | POST | - | ✓ | ✓ | Pengguna terotentikasi |
| `/auth/me` | GET | - | ✓ | ✓ | Pengguna terotentikasi |
| `/products` | GET | ✓ | ✓ | ✓ | Katalog umum produk aktif |
| `/products/:id` | GET | ✓ | ✓ | ✓ | Detail spesifik produk |
| `/seller/products` | GET | - | - | ✓ | Hanya produk milik seller aktif |
| `/products` | POST | - | - | ✓ | Seller membuat produk baru |
| `/products/:id` | PUT | - | - | ✓ | Hanya pemilik produk |
| `/products/:id` | DELETE | - | - | ✓ | Hanya pemilik produk (soft delete) |
| `/products/:id/stock` | PATCH | - | - | ✓ | Hanya pemilik produk |
| `/cart` | GET | - | ✓ | - | Keranjang aktif buyer |
| `/cart/items` | POST | - | ✓ | - | Tambah produk ke keranjang |
| `/cart/items/:id` | PATCH | - | ✓ | - | Ubah kuantitas item keranjang |
| `/cart/items/:id` | DELETE | - | ✓ | - | Hapus item dari keranjang |
| `/cart` | DELETE | - | ✓ | - | Kosongkan seluruh keranjang |
| `/orders` | POST | - | ✓ | - | Checkout keranjang menjadi pesanan |
| `/orders` | GET | - | ✓ | ✓ | Buyer: miliknya; Seller: toko miliknya |
| `/orders/:id` | GET | - | ✓ | ✓ | Hanya pihak yang terlibat (buyer/seller) |
| `/orders/:id/cancel` | PATCH | - | ✓ | - | Buyer pemilik pesanan (hanya pending) |
| `/orders/:id/status` | PATCH | - | - | ✓ | Seller penerima pesanan |
| `/orders/:id/complete` | PATCH | - | ✓ | - | Buyer pemilik pesanan (hanya shipped) |

---

## 3. Endpoints Detail

### 3.1 Authentication & Profile

#### A. Register User
- **Method / Path:** `POST /api/v1/auth/register`
- **Access:** Public
- **Request Body:**
```json
{
  "name": "Naabil Fikry",
  "email": "naabil@example.com",
  "password": "SecurePassword123",
  "role": "buyer"
}
```
- **Validation Rules:**
  - `name`: string, wajib, 2-100 karakter.
  - `email`: string, format email valid, wajib.
  - `password`: string, minimal 8 karakter.
  - `role`: string enum (`buyer` | `seller`), wajib.
- **Response 201 Created:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Naabil Fikry",
    "email": "naabil@example.com",
    "role": "buyer",
    "created_at": "2026-11-07T02:00:00Z"
  }
}
```

#### B. Login
- **Method / Path:** `POST /api/v1/auth/login`
- **Access:** Public
- **Request Body:**
```json
{
  "email": "naabil@example.com",
  "password": "SecurePassword123"
}
```
- **Response 200 OK:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 900,
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Naabil Fikry",
      "email": "naabil@example.com",
      "role": "buyer"
    }
  }
}
```
*(Header `Set-Cookie` disertakan server untuk menyimpan refresh token).*

#### C. Refresh Access Token
- **Method / Path:** `POST /api/v1/auth/refresh`
- **Access:** Public (dengan HttpOnly cookie `refresh_token`)
- **Response 200 OK:**
```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 900
  }
}
```

#### D. Logout
- **Method / Path:** `POST /api/v1/auth/logout`
- **Access:** Authenticated
- **Response 200 OK / 204 No Content:**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

#### E. Get Current User Profile
- **Method / Path:** `GET /api/v1/auth/me`
- **Access:** Authenticated
- **Response 200 OK:**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Naabil Fikry",
    "email": "naabil@example.com",
    "role": "buyer",
    "created_at": "2026-11-07T02:00:00Z"
  }
}
```

---

### 3.2 Products API

#### A. List Active Catalog Products
- **Method / Path:** `GET /api/v1/products`
- **Access:** Public
- **Query Parameters:**
  - `page`: int (default `1`)
  - `limit`: int (default `10`, max `50`)
  - `search`: string (nama produk, opsional)
  - `sort`: string enum (`latest` | `price_asc` | `price_desc`, default `latest`)
- **Response 200 OK:**
```json
{
  "success": true,
  "data": [
    {
      "id": "e6a2b85c-6b3a-4f51-b8ef-13b7d1883311",
      "name": "Mechanical Keyboard TKL RGB",
      "description": "Gateron Yellow Switch, PBT Keycaps",
      "price": 650000,
      "stock": 14,
      "image_url": "https://example.com/keyboard.webp",
      "seller": {
        "id": "a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d",
        "name": "Tech Gear Official"
      },
      "created_at": "2026-11-06T10:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "total_pages": 1
  }
}
```

#### B. Get Product Detail
- **Method / Path:** `GET /api/v1/products/:id`
- **Access:** Public
- **Response 200 OK:**
```json
{
  "success": true,
  "data": {
    "id": "e6a2b85c-6b3a-4f51-b8ef-13b7d1883311",
    "name": "Mechanical Keyboard TKL RGB",
    "description": "Gateron Yellow Switch, PBT Keycaps",
    "price": 650000,
    "stock": 14,
    "image_url": "https://example.com/keyboard.webp",
    "is_active": true,
    "seller": {
      "id": "a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d",
      "name": "Tech Gear Official",
      "email": "seller@techgear.com"
    },
    "created_at": "2026-11-06T10:00:00Z",
    "updated_at": "2026-11-06T10:00:00Z"
  }
}
```

#### C. Create Product
- **Method / Path:** `POST /api/v1/products`
- **Access:** Authenticated (Role: `seller`)
- **Request Body:**
```json
{
  "name": "Wireless Gaming Mouse",
  "description": "26K DPI, Ultra-lightweight 55g",
  "price": 420000,
  "stock": 25,
  "image_url": "https://example.com/mouse.webp"
}
```
- **Response 201 Created:**
```json
{
  "success": true,
  "message": "Product created successfully",
  "data": {
    "id": "789e0123-4567-89ab-cdef-0123456789ab",
    "seller_id": "a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d",
    "name": "Wireless Gaming Mouse",
    "price": 420000,
    "stock": 25,
    "is_active": true,
    "created_at": "2026-11-07T02:15:00Z"
  }
}
```

#### D. Update Product
- **Method / Path:** `PUT /api/v1/products/:id`
- **Access:** Authenticated (Role: `seller`, Pemilik Produk)
- **Request Body:** Sama seperti Create Product
- **Response 200 OK:**
```json
{
  "success": true,
  "message": "Product updated successfully",
  "data": { ... }
}
```

#### E. Soft-Delete Product
- **Method / Path:** `DELETE /api/v1/products/:id`
- **Access:** Authenticated (Role: `seller`, Pemilik Produk)
- **Response 200 OK / 204 No Content:**
```json
{
  "success": true,
  "message": "Product deleted successfully"
}
```

#### F. Quick Update Stock
- **Method / Path:** `PATCH /api/v1/products/:id/stock`
- **Access:** Authenticated (Role: `seller`, Pemilik Produk)
- **Request Body:**
```json
{
  "stock": 50
}
```
- **Response 200 OK:**
```json
{
  "success": true,
  "message": "Stock updated successfully",
  "data": {
    "id": "789e0123-4567-89ab-cdef-0123456789ab",
    "stock": 50
  }
}
```

---

### 3.3 Shopping Cart API

#### A. Get Cart
- **Method / Path:** `GET /api/v1/cart`
- **Access:** Authenticated (Role: `buyer`)
- **Response 200 OK:**
```json
{
  "success": true,
  "data": {
    "id": "cart-uuid-112233",
    "items": [
      {
        "id": "item-uuid-445566",
        "product": {
          "id": "e6a2b85c-6b3a-4f51-b8ef-13b7d1883311",
          "name": "Mechanical Keyboard TKL RGB",
          "price": 650000,
          "stock": 14,
          "image_url": "https://example.com/keyboard.webp",
          "seller_id": "a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d"
        },
        "quantity": 2,
        "subtotal": 1300000
      }
    ],
    "total_items": 2,
    "total_amount": 1300000
  }
}
```

#### B. Add Item to Cart
- **Method / Path:** `POST /api/v1/cart/items`
- **Access:** Authenticated (Role: `buyer`)
- **Request Body:**
```json
{
  "product_id": "e6a2b85c-6b3a-4f51-b8ef-13b7d1883311",
  "quantity": 2
}
```
- **Response 201 Created:**
```json
{
  "success": true,
  "message": "Item added to cart",
  "data": {
    "cart_id": "cart-uuid-112233",
    "product_id": "e6a2b85c-6b3a-4f51-b8ef-13b7d1883311",
    "quantity": 2
  }
}
```

#### C. Update Cart Item Quantity
- **Method / Path:** `PATCH /api/v1/cart/items/:id`
- **Access:** Authenticated (Role: `buyer`)
- **Request Body:**
```json
{
  "quantity": 3
}
```
- **Response 200 OK:**
```json
{
  "success": true,
  "message": "Cart item updated successfully"
}
```

#### D. Delete Cart Item
- **Method / Path:** `DELETE /api/v1/cart/items/:id`
- **Access:** Authenticated (Role: `buyer`)
- **Response 200 OK / 204 No Content:**
```json
{
  "success": true,
  "message": "Cart item removed successfully"
}
```

#### E. Clear Cart
- **Method / Path:** `DELETE /api/v1/cart`
- **Access:** Authenticated (Role: `buyer`)
- **Response 200 OK / 204 No Content:**
```json
{
  "success": true,
  "message": "Cart cleared successfully"
}
```

---

### 3.4 Orders & Checkout API

#### A. Checkout
- **Method / Path:** `POST /api/v1/orders`
- **Access:** Authenticated (Role: `buyer`)
- **Request Body:**
```json
{
  "shipping_address": {
    "recipient_name": "Naabil Fikry",
    "phone": "081234567890",
    "address": "Jl. Telekomunikasi No. 1",
    "city": "Bandung",
    "province": "Jawa Barat",
    "postal_code": "40257"
  },
  "notes": "Mohon bubble wrap tebal ya min"
}
```
*(Catatan: Harga, seller, dan barang diambil otomatis dari database cart aktif).*
- **Response 201 Created:**
```json
{
  "success": true,
  "message": "Order created successfully",
  "data": {
    "id": "order-uuid-998877",
    "buyer_id": "550e8400-e29b-41d4-a716-446655440000",
    "seller_id": "a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d",
    "total_amount": 1300000,
    "status": "pending",
    "shipping_address": {
      "recipient_name": "Naabil Fikry",
      "phone": "081234567890",
      "address": "Jl. Telekomunikasi No. 1",
      "city": "Bandung",
      "province": "Jawa Barat",
      "postal_code": "40257"
    },
    "items": [
      {
        "product_id": "e6a2b85c-6b3a-4f51-b8ef-13b7d1883311",
        "product_name": "Mechanical Keyboard TKL RGB",
        "unit_price": 650000,
        "quantity": 2,
        "subtotal": 1300000
      }
    ],
    "created_at": "2026-11-07T03:00:00Z"
  }
}
```

#### B. List Orders
- **Method / Path:** `GET /api/v1/orders`
- **Access:** Authenticated (`buyer` / `seller`)
- **Query Parameters:**
  - `page`: int (default `1`)
  - `limit`: int (default `10`)
  - `status`: string (filter status, opsional)
- **Response 200 OK:**
```json
{
  "success": true,
  "data": [
    {
      "id": "order-uuid-998877",
      "total_amount": 1300000,
      "status": "pending",
      "buyer": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "Naabil Fikry"
      },
      "created_at": "2026-11-07T03:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "total_pages": 1
  }
}
```

#### C. Get Order Detail
- **Method / Path:** `GET /api/v1/orders/:id`
- **Access:** Authenticated (Buyer pemilik atau Seller penerima)
- **Response 200 OK:**
```json
{
  "success": true,
  "data": {
    "id": "order-uuid-998877",
    "status": "pending",
    "total_amount": 1300000,
    "shipping_address": {
      "recipient_name": "Naabil Fikry",
      "phone": "081234567890",
      "address": "Jl. Telekomunikasi No. 1",
      "city": "Bandung",
      "province": "Jawa Barat",
      "postal_code": "40257"
    },
    "notes": "Mohon bubble wrap tebal ya min",
    "items": [
      {
        "id": "oi-uuid-1",
        "product_id": "e6a2b85c-6b3a-4f51-b8ef-13b7d1883311",
        "product_name": "Mechanical Keyboard TKL RGB",
        "unit_price": 650000,
        "quantity": 2,
        "subtotal": 1300000
      }
    ],
    "created_at": "2026-11-07T03:00:00Z",
    "updated_at": "2026-11-07T03:00:00Z"
  }
}
```

#### D. Cancel Order (Buyer)
- **Method / Path:** `PATCH /api/v1/orders/:id/cancel`
- **Access:** Authenticated (Role: `buyer`, Pemilik Pesanan)
- **Request Body:**
```json
{
  "reason": "Ingin mengubah alamat pengiriman"
}
```
- **Response 200 OK:**
```json
{
  "success": true,
  "message": "Order cancelled successfully. Product stock restored.",
  "data": {
    "id": "order-uuid-998877",
    "status": "cancelled"
  }
}
```

#### E. Update Order Status (Seller)
- **Method / Path:** `PATCH /api/v1/orders/:id/status`
- **Access:** Authenticated (Role: `seller`, Penerima Pesanan)
- **Request Body:**
```json
{
  "status": "confirmed"
}
```
*(Valid target statuses: `confirmed`, `processing`, `shipped`, `rejected`)*.  
*Jika status diubah ke `rejected`, stok produk otomatis dikembalikan dalam transaksi database.*
- **Response 200 OK:**
```json
{
  "success": true,
  "message": "Order status updated to confirmed",
  "data": {
    "id": "order-uuid-998877",
    "status": "confirmed"
  }
}
```

#### F. Complete Order (Buyer)
- **Method / Path:** `PATCH /api/v1/orders/:id/complete`
- **Access:** Authenticated (Role: `buyer`, Pemilik Pesanan)
- **Prerequisite:** Status pesanan saat ini harus `shipped`.
- **Response 200 OK:**
```json
{
  "success": true,
  "message": "Order confirmed as completed. Thank you!",
  "data": {
    "id": "order-uuid-998877",
    "status": "completed"
  }
}
```

