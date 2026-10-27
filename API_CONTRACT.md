# API Contract Documentation
## Sistem Point of Sales (POS) Coffee Shop

**Dokumen Referensi:** [PRD.md](file:///home/djan/Desktop/web-coffee/PRD.md) & [SCHEMA.md](file:///home/djan/Desktop/web-coffee/SCHEMA.md)  
**Versi API:** v1  
**Tanggal:** 27 Oktober 2026  
**Status:** Approved / Acuan Implementasi Backend & Frontend  

---

## 1. Base URL & Versioning

- **Base URL:** `http://localhost:8080/api/v1` (Development) / `https://api.pos-coffeeshop.com/api/v1` (Production)
- **Content-Type:** `application/json` (Kecuali endpoint multipart/form-data untuk upload foto produk dan download CSV).

---

## 2. Autentikasi (Authentication)

Setiap request ke endpoint terproteksi wajib menyertakan token JWT (*Bearer Token*) pada HTTP Request Header:

```http
Authorization: Bearer <access_token>
```

- Token diterbitkan setelah login berhasil melalui endpoint `/api/v1/auth/login`.
- Payload JWT memuat klaim:
  - `user_id`: UUID pengguna
  - `name`: Nama lengkap
  - `email`: Alamat email
  - `role`: `'owner'` atau `'cashier'`
  - `exp`: Waktu kedaluwarsa token (misal: 24 jam)

Jika header tidak ada, tidak valid, atau kedaluwarsa, server mengembalikan respon `401 Unauthorized`.

---

## 3. Format Respon Standar (Standard Response Format)

Untuk menjaga konsistensi integrasi antara frontend dan backend, seluruh respon JSON wajib mengikuti struktur berikut:

### 3.1 Format Respon

#### 1. Respon Sukses (Single Object / Operasi)
Digunakan untuk operasi seperti detail item, update, hapus, atau aksi non-list.
```json
{
  "success": true,
  "message": "Operasi berhasil dilakukan",
  "data": { ... }
}
```

#### 2. Respon Sukses dengan List (Disertai Pagination `meta`)
Digunakan untuk endpoint yang mengembalikan daftar data dengan paginasi.
```json
{
  "success": true,
  "message": "Berhasil mengambil data produk",
  "data": [ ... ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 45,
    "total_pages": 3
  }
}
```

#### 3. Respon Error Validasi (HTTP 422)
Digunakan saat payload atau parameter request gagal memenuhi validasi field formulir/skema.
```json
{
  "success": false,
  "message": "Validasi gagal",
  "errors": {
    "name": "Nama produk wajib diisi",
    "price": "Harga harus lebih dari 0"
  }
}
```

#### 4. Respon Error Umum (HTTP 400, 401, 403, 404, 429, 500)
Digunakan untuk kegagalan alur bisnis, autentikasi, otorisasi, data tidak ditemukan, rate limiting, atau internal server error.
```json
{
  "success": false,
  "message": "Sesi shift kasir belum dibuka"
}
```

---

### 3.2 HTTP Status Code yang Digunakan

| HTTP Status | Definisi & Penggunaan |
| :---: | :--- |
| **`200`** | **Sukses GET dan PUT** (serta aksi PATCH / DELETE yang mengembalikan response body) |
| **`201`** | **Sukses POST (create)** saat entitas baru berhasil dibuat |
| **`400`** | **Request tidak valid** (pelanggaran aturan bisnis, malformed JSON, atau logika alur salah) |
| **`401`** | **Tidak ada token atau token expired** (gagal verifikasi identitas pengguna) |
| **`403`** | **Token valid tapi tidak punya akses** (misal: Kasir mengakses fitur khusus Owner) |
| **`404`** | **Data tidak ditemukan** (ID resource tidak ada dalam database) |
| **`422`** | **Validasi gagal** (input request tidak lolos validasi tipe/format field) |
| **`429`** | **Rate limit** (terlalu banyak request dalam jangka waktu tertentu) |
| **`500`** | **Server error** (kesalahan internal server atau kegagalan koneksi database) |

---

## 4. Daftar Endpoint API

### 4.1 Auth

#### POST `/api/v1/auth/login`
- **Akses:** Public
- **Deskripsi:** Autentikasi akun pengguna (Owner atau Kasir).
- **Request Body:**
```json
{
  "email": "kasir1@coffeeshop.com",
  "password": "Password123!"
}
```
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "c1f20b33-1f74-4b52-bbf8-2b89ad9211d1",
      "name": "Budi Santoso",
      "email": "kasir1@coffeeshop.com",
      "role": "cashier",
      "is_active": true
    }
  }
}
```
- **Kemungkinan Respon Error:**
  - `401 Unauthorized` (Email atau password salah)
  - `403 Forbidden` (Akun kasir berstatus nonaktif: `is_active: false`)
  - `422 Unprocessable Entity` (Email/password kosong atau format tidak valid)

---

#### GET `/api/v1/auth/me`
- **Akses:** Owner, Cashier
- **Deskripsi:** Mendapatkan data profil pengguna yang sedang login berdasarkan token JWT.
- **Request:** Header `Authorization: Bearer <token>`
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Profil berhasil dimuat",
  "data": {
    "id": "c1f20b33-1f74-4b52-bbf8-2b89ad9211d1",
    "name": "Budi Santoso",
    "email": "kasir1@coffeeshop.com",
    "role": "cashier"
  }
}
```
- **Kemungkinan Respon Error:**
  - `401 Unauthorized` (Token tidak valid / expired)

---

#### POST `/api/v1/auth/logout`
- **Akses:** Owner, Cashier
- **Deskripsi:** Mengakhiri sesi pengguna (invalidasi token di sisi client / token blacklist jika diaktifkan).
- **Request:** Header `Authorization: Bearer <token>`
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Logout berhasil"
}
```

---

### 4.2 Categories (Owner)

#### GET `/api/v1/categories`
- **Akses:** Owner, Cashier
- **Query Params:**
  - `is_active` (boolean, opsional): Filter kategori aktif saja (misal untuk katalog kasir `?is_active=true`)
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Daftar kategori berhasil diambil",
  "data": [
    {
      "id": "78a9c2b4-e40d-4071-b0db-6e1be87eb8f1",
      "name": "Espresso Based",
      "slug": "espresso-based",
      "description": "Kopi berbasis ekstrak espresso murni",
      "sort_order": 1,
      "status": "active"
    }
  ]
}
```

---

#### POST `/api/v1/categories`
- **Akses:** Owner
- **Request Body:**
```json
{
  "name": "Manual Brew",
  "description": "Kopi seduh manual V60, Japanese, Aeropress",
  "sort_order": 2,
  "status": "active"
}
```
- **Response Sukses (201 Created):**
```json
{
  "success": true,
  "message": "Kategori berhasil dibuat",
  "data": {
    "id": "5f1a9a83-a97e-4091-a185-9b2f6efba983",
    "name": "Manual Brew",
    "slug": "manual-brew",
    "description": "Kopi seduh manual V60, Japanese, Aeropress",
    "sort_order": 2,
    "status": "active",
    "created_at": "2026-10-27T10:00:00Z"
  }
}
```
- **Kemungkinan Respon Error:**
  - `403 Forbidden` (User bukan Owner)
  - `409 Conflict` (Slug / Nama kategori sudah ada)
  - `422 Unprocessable Entity` (Nama wajib diisi atau validasi field gagal)

---

#### GET `/api/v1/categories/:id`
- **Akses:** Owner
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Detail kategori berhasil diambil",
  "data": {
    "id": "5f1a9a83-a97e-4091-a185-9b2f6efba983",
    "name": "Manual Brew",
    "slug": "manual-brew",
    "description": "Kopi seduh manual",
    "sort_order": 2,
    "status": "active"
  }
}
```
- **Kemungkinan Respon Error:**
  - `404 Not Found` (Kategori tidak ditemukan)

---

#### PUT `/api/v1/categories/:id`
- **Akses:** Owner
- **Request Body:**
```json
{
  "name": "Manual Brew Specialty",
  "description": "Kopi seduh manual filter spesial",
  "sort_order": 2,
  "status": "active"
}
```
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Kategori berhasil diperbarui",
  "data": {
    "id": "5f1a9a83-a97e-4091-a185-9b2f6efba983",
    "name": "Manual Brew Specialty",
    "slug": "manual-brew-specialty",
    "description": "Kopi seduh manual filter spesial",
    "sort_order": 2,
    "status": "active"
  }
}
```
- **Kemungkinan Respon Error:**
  - `404 Not Found` (Kategori tidak ditemukan)

---

#### DELETE `/api/v1/categories/:id`
- **Akses:** Owner
- **Deskripsi:** Soft delete kategori produk.
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Kategori berhasil dihapus"
}
```
- **Kemungkinan Respon Error:**
  - `400 Bad Request` (Kategori masih memiliki produk aktif yang tertaut)
  - `404 Not Found` (Kategori tidak ditemukan)

---

### 4.3 Products (Owner)

#### GET `/api/v1/products`
- **Akses:** Owner, Cashier
- **Query Params:**
  - `category_id` (string UUID, opsional)
  - `status` (string: `active`/`inactive`, opsional)
  - `search` (string, opsional)
  - `page` (int, default: 1)
  - `per_page` (int, default: 20)
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Berhasil mengambil data produk",
  "data": [
    {
      "id": "e4a2e582-7f72-46a2-9d33-3dcb295a643b",
      "category_id": "78a9c2b4-e40d-4071-b0db-6e1be87eb8f1",
      "category_name": "Espresso Based",
      "sku": "COF-LAT-001",
      "name": "Iced Caffe Latte",
      "description": "Double shot espresso dengan fresh milk dingin",
      "price": 2800000,
      "cost_price": 1200000,
      "image_url": "https://cdn.coffeeshop.com/products/iced-latte.webp",
      "current_stock": 45,
      "min_stock_alert": 10,
      "status": "active"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 45,
    "total_pages": 3
  }
}
```

---

#### POST `/api/v1/products`
- **Akses:** Owner
- **Request Format:** `multipart/form-data` (atau `application/json` jika upload gambar terpisah)
- **Request Form / Body:**
```json
{
  "category_id": "78a9c2b4-e40d-4071-b0db-6e1be87eb8f1",
  "sku": "COF-CAP-002",
  "name": "Hot Cappuccino",
  "description": "Espresso dengan steamed milk foam tebal",
  "price": 2500000,
  "cost_price": 1100000,
  "image_url": "https://cdn.coffeeshop.com/products/cappuccino.webp",
  "initial_stock": 50,
  "min_stock_alert": 5,
  "status": "active"
}
```
- **Response Sukses (201 Created):**
```json
{
  "success": true,
  "message": "Produk berhasil ditambahkan",
  "data": {
    "id": "d8b3941a-1d57-41eb-8785-cfc09f3032ef",
    "sku": "COF-CAP-002",
    "name": "Hot Cappuccino",
    "price": 2500000,
    "cost_price": 1100000,
    "current_stock": 50,
    "status": "active"
  }
}
```
- **Kemungkinan Respon Error:**
  - `400 Bad Request` (Kategori tidak ditemukan)
  - `403 Forbidden` (Bukan Owner)
  - `409 Conflict` (SKU sudah terdaftar)
  - `422 Unprocessable Entity` (Validasi input gagal)
```json
{
  "success": false,
  "message": "Validasi gagal",
  "errors": {
    "name": "Nama produk wajib diisi",
    "price": "Harga harus lebih dari 0"
  }
}
```

---

#### GET `/api/v1/products/:id`
- **Akses:** Owner, Cashier
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Detail produk berhasil diambil",
  "data": {
    "id": "d8b3941a-1d57-41eb-8785-cfc09f3032ef",
    "category_id": "78a9c2b4-e40d-4071-b0db-6e1be87eb8f1",
    "sku": "COF-CAP-002",
    "name": "Hot Cappuccino",
    "description": "Espresso dengan steamed milk foam tebal",
    "price": 2500000,
    "cost_price": 1100000,
    "image_url": "https://cdn.coffeeshop.com/products/cappuccino.webp",
    "current_stock": 50,
    "min_stock_alert": 5,
    "status": "active"
  }
}
```
- **Kemungkinan Respon Error:**
  - `404 Not Found` (Produk tidak ditemukan)

---

#### PUT `/api/v1/products/:id`
- **Akses:** Owner
- **Request Body:**
```json
{
  "category_id": "78a9c2b4-e40d-4071-b0db-6e1be87eb8f1",
  "sku": "COF-CAP-002",
  "name": "Hot Cappuccino Premium",
  "description": "Espresso blend arabica dengan steamed milk foam tebal",
  "price": 2700000,
  "cost_price": 1200000,
  "image_url": "https://cdn.coffeeshop.com/products/cappuccino.webp",
  "min_stock_alert": 10,
  "status": "active"
}
```
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Produk berhasil diperbarui",
  "data": {
    "id": "d8b3941a-1d57-41eb-8785-cfc09f3032ef",
    "name": "Hot Cappuccino Premium",
    "price": 2700000
  }
}
```
- **Kemungkinan Respon Error:**
  - `400 Bad Request` (Harga $\le 0$, format salah)
  - `404 Not Found` (Produk tidak ditemukan)

---

#### PATCH `/api/v1/products/:id/status`
- **Akses:** Owner
- **Request Body:**
```json
{
  "status": "inactive"
}
```
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Status produk berhasil diubah",
  "data": {
    "id": "d8b3941a-1d57-41eb-8785-cfc09f3032ef",
    "status": "inactive"
  }
}
```

---

#### DELETE `/api/v1/products/:id`
- **Akses:** Owner
- **Deskripsi:** Soft delete produk (`deleted_at` terisi).
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Produk berhasil dihapus"
}
```
- **Kemungkinan Respon Error:**
  - `404 Not Found` (Produk tidak ditemukan)

---

### 4.4 Stock (Owner)

#### GET `/api/v1/stock/movements`
- **Akses:** Owner
- **Query Params:**
  - `product_id` (string UUID, opsional)
  - `movement_type` (`IN`/`OUT`/`SALE`/`ADJUSTMENT`, opsional)
  - `start_date` (ISO date, opsional: `2026-10-01`)
  - `end_date` (ISO date, opsional: `2026-10-27`)
  - `page` (int, default: 1)
  - `per_page` (int, default: 20)
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Riwayat mutasi stok berhasil dimuat",
  "data": [
    {
      "id": "f2a1b920-5bb2-40cf-a734-7c5ef1132890",
      "product_id": "e4a2e582-7f72-46a2-9d33-3dcb295a643b",
      "product_name": "Iced Caffe Latte",
      "user_name": "Owner Admin",
      "movement_type": "ADJUSTMENT",
      "quantity": 5,
      "stock_before": 45,
      "stock_after": 40,
      "reference_type": "waste",
      "notes": "Susu basi dan botol tumpah",
      "created_at": "2026-10-27T08:30:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 52,
    "total_pages": 3
  }
}
```

---

#### POST `/api/v1/stock/adjustments`
- **Akses:** Owner
- **Deskripsi:** Penyesuaian stok manual (tambah stok masuk, selisih fisik, atau *waste*).
- **Request Body:**
```json
{
  "product_id": "e4a2e582-7f72-46a2-9d33-3dcb295a643b",
  "movement_type": "ADJUSTMENT",
  "quantity": 10,
  "action": "DECREASE",
  "reference_type": "waste",
  "notes": "Barang kedaluwarsa"
}
```
*(Catatan: `action` bernilai `INCREASE` atau `DECREASE`)*
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Penyesuaian stok berhasil disimpan",
  "data": {
    "product_id": "e4a2e582-7f72-46a2-9d33-3dcb295a643b",
    "stock_before": 50,
    "stock_after": 40,
    "current_stock": 40
  }
}
```
- **Kemungkinan Respon Error:**
  - `400 Bad Request` (Pengurangan stok menyebabkan kuantitas $< 0$)
  - `404 Not Found` (Produk tidak ditemukan)

---

#### GET `/api/v1/stock/low-stock`
- **Akses:** Owner
- **Deskripsi:** Mendapatkan daftar produk yang stok fisiknya $\le$ `min_stock_alert`.
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Daftar produk dengan stok menipis berhasil diambil",
  "data": [
    {
      "product_id": "d8b3941a-1d57-41eb-8785-cfc09f3032ef",
      "name": "Hot Cappuccino",
      "current_stock": 3,
      "min_stock_alert": 5
    }
  ]
}
```

---

### 4.5 Tables (Owner)

#### GET `/api/v1/tables`
- **Akses:** Owner, Cashier
- **Query Params:**
  - `status` (`available`/`occupied`/`reserved`, opsional)
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Daftar meja berhasil diambil",
  "data": [
    {
      "id": "a90b4d4f-56f0-4660-845f-4a3cebf27e01",
      "table_number": "Meja 01",
      "capacity": 4,
      "status": "available"
    },
    {
      "id": "b11c5e5a-67a1-4771-956a-5b4dfca38f12",
      "table_number": "Meja 02 (VIP)",
      "capacity": 6,
      "status": "occupied"
    }
  ]
}
```

---

#### POST `/api/v1/tables`
- **Akses:** Owner
- **Request Body:**
```json
{
  "table_number": "Meja 03",
  "capacity": 2
}
```
- **Response Sukses (201 Created):**
```json
{
  "success": true,
  "message": "Meja berhasil ditambahkan",
  "data": {
    "id": "c22d6f6b-78b2-4882-a67b-6c5efdb49a23",
    "table_number": "Meja 03",
    "capacity": 2,
    "status": "available"
  }
}
```
- **Kemungkinan Respon Error:**
  - `409 Conflict` (Nomor meja sudah ada)
  - `422 Unprocessable Entity` (Nomor meja kosong atau kapasitas $\le 0$)

---

#### PUT `/api/v1/tables/:id`
- **Akses:** Owner
- **Request Body:**
```json
{
  "table_number": "Meja 03 (Smoking)",
  "capacity": 4
}
```
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Data meja berhasil diperbarui",
  "data": {
    "id": "c22d6f6b-78b2-4882-a67b-6c5efdb49a23",
    "table_number": "Meja 03 (Smoking)",
    "capacity": 4,
    "status": "available"
  }
}
```
- **Kemungkinan Respon Error:**
  - `404 Not Found` (Meja tidak ditemukan)

---

#### DELETE `/api/v1/tables/:id`
- **Akses:** Owner
- **Deskripsi:** Soft delete meja kedai.
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Meja berhasil dihapus"
}
```
- **Kemungkinan Respon Error:**
  - `400 Bad Request` (Meja sedang dalam status `occupied` dengan order aktif)
  - `404 Not Found` (Meja tidak ditemukan)

---

### 4.6 Users/Cashier Management (Owner)

#### GET `/api/v1/users`
- **Akses:** Owner
- **Query Params:**
  - `role` (`cashier`/`owner`, opsional)
  - `is_active` (boolean, opsional)
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Daftar pengguna berhasil dimuat",
  "data": [
    {
      "id": "c1f20b33-1f74-4b52-bbf8-2b89ad9211d1",
      "name": "Budi Santoso",
      "email": "kasir1@coffeeshop.com",
      "role": "cashier",
      "is_active": true,
      "created_at": "2026-10-27T08:00:00Z"
    }
  ]
}
```

---

#### POST `/api/v1/users`
- **Akses:** Owner
- **Request Body:**
```json
{
  "name": "Siti Nurhaliza",
  "email": "kasir2@coffeeshop.com",
  "password": "Password123!",
  "role": "cashier",
  "is_active": true
}
```
- **Response Sukses (201 Created):**
```json
{
  "success": true,
  "message": "Pengguna kasir berhasil dibuat",
  "data": {
    "id": "e33f7a7c-89c3-4993-b78c-7d6f0ec50b34",
    "name": "Siti Nurhaliza",
    "email": "kasir2@coffeeshop.com",
    "role": "cashier",
    "is_active": true
  }
}
```
- **Kemungkinan Respon Error:**
  - `409 Conflict` (Email sudah digunakan pengguna lain)
  - `422 Unprocessable Entity` (Validasi format email atau panjang password gagal)

---

#### GET `/api/v1/users/:id`
- **Akses:** Owner
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Detail pengguna berhasil diambil",
  "data": {
    "id": "e33f7a7c-89c3-4993-b78c-7d6f0ec50b34",
    "name": "Siti Nurhaliza",
    "email": "kasir2@coffeeshop.com",
    "role": "cashier",
    "is_active": true
  }
}
```

---

#### PUT `/api/v1/users/:id`
- **Akses:** Owner
- **Request Body:**
```json
{
  "name": "Siti Nurhaliza S.E.",
  "email": "kasir2_update@coffeeshop.com"
}
```
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Data pengguna berhasil diperbarui",
  "data": {
    "id": "e33f7a7c-89c3-4993-b78c-7d6f0ec50b34",
    "name": "Siti Nurhaliza S.E.",
    "email": "kasir2_update@coffeeshop.com"
  }
}
```

---

#### PATCH `/api/v1/users/:id/status`
- **Akses:** Owner
- **Request Body:**
```json
{
  "is_active": false
}
```
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Status pengguna berhasil diubah",
  "data": {
    "id": "e33f7a7c-89c3-4993-b78c-7d6f0ec50b34",
    "is_active": false
  }
}
```

---

#### PUT `/api/v1/users/:id/reset-password`
- **Akses:** Owner
- **Request Body:**
```json
{
  "new_password": "NewSecretPassword2026!"
}
```
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Password pengguna berhasil direset"
}
```

---

#### DELETE `/api/v1/users/:id`
- **Akses:** Owner
- **Deskripsi:** Soft delete akun pengguna.
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Pengguna berhasil dihapus"
}
```
- **Kemungkinan Respon Error:**
  - `400 Bad Request` (Owner tidak dapat menghapus akun miliknya sendiri)
  - `404 Not Found` (Pengguna tidak ditemukan)

---

### 4.7 Promos (Owner)

#### GET `/api/v1/promos`
- **Akses:** Owner
- **Query Params:**
  - `status` (`active`/`inactive`, opsional)
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Daftar promo berhasil diambil",
  "data": [
    {
      "id": "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
      "code": "SENINKOPI",
      "name": "Diskon Senin Semangat 20%",
      "discount_type": "percentage",
      "discount_value": 20,
      "max_discount": 2000000,
      "min_spend": 5000000,
      "start_date": "2026-10-01T00:00:00Z",
      "end_date": "2026-12-31T23:59:59Z",
      "quota": 100,
      "used_count": 14,
      "status": "active"
    }
  ]
}
```

---

#### POST `/api/v1/promos`
- **Akses:** Owner
- **Request Body:**
```json
{
  "code": "HEMAT10K",
  "name": "Potongan Langsung 10 Ribu",
  "discount_type": "fixed_amount",
  "discount_value": 1000000,
  "max_discount": null,
  "min_spend": 4000000,
  "start_date": "2026-10-27T00:00:00Z",
  "end_date": "2026-11-30T23:59:59Z",
  "quota": 50,
  "status": "active"
}
```
- **Response Sukses (201 Created):**
```json
{
  "success": true,
  "message": "Promo berhasil dibuat",
  "data": {
    "id": "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
    "code": "HEMAT10K",
    "name": "Potongan Langsung 10 Ribu",
    "status": "active"
  }
}
```
- **Kemungkinan Respon Error:**
  - `409 Conflict` (Kode promo sudah pernah didaftarkan)
  - `422 Unprocessable Entity` (Tanggal mulai > tanggal selesai, atau nilai diskon $\le 0$)

---

#### GET `/api/v1/promos/:id`
- **Akses:** Owner
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Detail promo berhasil diambil",
  "data": {
    "id": "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
    "code": "HEMAT10K",
    "name": "Potongan Langsung 10 Ribu",
    "discount_type": "fixed_amount",
    "discount_value": 1000000,
    "max_discount": null,
    "min_spend": 4000000,
    "start_date": "2026-10-27T00:00:00Z",
    "end_date": "2026-11-30T23:59:59Z",
    "quota": 50,
    "used_count": 0,
    "status": "active"
  }
}
```

---

#### PUT `/api/v1/promos/:id`
- **Akses:** Owner
- **Request Body:**
```json
{
  "name": "Potongan Spesial 10 Ribu",
  "discount_type": "fixed_amount",
  "discount_value": 1000000,
  "max_discount": null,
  "min_spend": 5000000,
  "start_date": "2026-10-27T00:00:00Z",
  "end_date": "2026-12-31T23:59:59Z",
  "quota": 100,
  "status": "active"
}
```
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Promo berhasil diperbarui",
  "data": {
    "id": "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
    "name": "Potongan Spesial 10 Ribu"
  }
}
```

---

#### PATCH `/api/v1/promos/:id/status`
- **Akses:** Owner
- **Request Body:**
```json
{
  "status": "inactive"
}
```
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Status promo berhasil diperbarui",
  "data": {
    "id": "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
    "status": "inactive"
  }
}
```

---

#### DELETE `/api/v1/promos/:id`
- **Akses:** Owner
- **Deskripsi:** Soft delete promo.
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Promo berhasil dihapus"
}
```

---

### 4.8 Reports (Owner)

#### GET `/api/v1/reports/dashboard`
- **Akses:** Owner
- **Deskripsi:** Menghitung total revenue agregat, tren pesanan, produk terlaris, dan performa per kasir (hanya menghitung transaksi berstatus `paid`).
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Data analitik dashboard berhasil dimuat",
  "data": {
    "revenue": {
      "daily": 450000000,
      "weekly": 3150000000,
      "monthly": 12600000000
    },
    "sales_trend": [
      {
        "date": "2026-10-25",
        "total_revenue": 420000000,
        "total_orders": 85
      },
      {
        "date": "2026-10-26",
        "total_revenue": 450000000,
        "total_orders": 92
      }
    ],
    "top_selling_products": [
      {
        "product_id": "e4a2e582-7f72-46a2-9d33-3dcb295a643b",
        "product_name": "Iced Caffe Latte",
        "total_quantity_sold": 180,
        "total_revenue": 504000000
      },
      {
        "product_id": "d8b3941a-1d57-41eb-8785-cfc09f3032ef",
        "product_name": "Hot Cappuccino",
        "total_quantity_sold": 125,
        "total_revenue": 312500000
      }
    ],
    "cashier_performance": [
      {
        "cashier_id": "c1f20b33-1f74-4b52-bbf8-2b89ad9211d1",
        "cashier_name": "Budi Santoso",
        "total_transactions": 65,
        "total_revenue": 310000000
      },
      {
        "cashier_id": "e33f7a7c-89c3-4993-b78c-7d6f0ec50b34",
        "cashier_name": "Siti Nurhaliza",
        "total_transactions": 27,
        "total_revenue": 140000000
      }
    ]
  }
}
```

---

#### GET `/api/v1/reports/export/csv`
- **Akses:** Owner
- **Query Params:**
  - `start_date` (format: `YYYY-MM-DD`, required)
  - `end_date` (format: `YYYY-MM-DD`, required)
- **Response Sukses (200 OK):**
  - **Headers:**
    - `Content-Type: text/csv; charset=utf-8`
    - `Content-Disposition: attachment; filename="sales-report-2026-10-01-to-2026-10-27.csv"`
  - **Body Format:**
```csv
Invoice Number,Date Time,Cashier,Order Type,Table Number,Subtotal,Discount,Tax,Total Amount,Payment Method,Status
INV-20261027-0001,2026-10-27 10:15:00,Budi Santoso,dine_in,Meja 01,53000,5000,0,48000,midtrans_snap,paid
INV-20261027-0002,2026-10-27 10:45:00,Budi Santoso,takeaway,-,28000,0,0,28000,cash,paid
```
- **Kemungkinan Respon Error:**
  - `400 Bad Request` (Format tanggal salah atau rentang tanggal terbalik)

---

### 4.9 Shifts (Cashier)

#### POST `/api/v1/shifts/open`
- **Akses:** Cashier
- **Deskripsi:** Membuka sesi shift kasir dengan memasukkan modal uang kas laci pertama kali.
- **Request Body:**
```json
{
  "starting_cash": 20000000,
  "notes": "Modal awal kas shift pagi"
}
```
*(Nilai `starting_cash` Rp 200.000 disimpan sebagai `20000000` sen)*
- **Response Sukses (201 Created):**
```json
{
  "success": true,
  "message": "Shift berhasil dibuka",
  "data": {
    "id": "99e1a2b3-c4d5-4e6f-a7b8-c9d0e1f2a3b4",
    "user_id": "c1f20b33-1f74-4b52-bbf8-2b89ad9211d1",
    "starting_cash": 20000000,
    "start_time": "2026-10-27T08:00:00Z",
    "status": "open"
  }
}
```
- **Kemungkinan Respon Error:**
  - `400 Bad Request` (Modal awal kas negatif $< 0$)
  - `409 Conflict` (Kasir sudah memiliki shift berstatus `open` yang belum ditutup)

---

#### POST `/api/v1/shifts/close`
- **Akses:** Cashier
- **Deskripsi:** Menutup sesi shift kasir yang sedang aktif dan memasukkan hitungan fisik uang kas.
- **Request Body:**
```json
{
  "actual_ending_cash": 65000000,
  "notes": "Penutupan shift pagi, selisih pas"
}
```
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Shift berhasil ditutup",
  "data": {
    "id": "99e1a2b3-c4d5-4e6f-a7b8-c9d0e1f2a3b4",
    "starting_cash": 20000000,
    "total_cash_sales": 45000000,
    "expected_ending_cash": 65000000,
    "actual_ending_cash": 65000000,
    "cash_difference": 0,
    "end_time": "2026-10-27T16:00:00Z",
    "status": "closed"
  }
}
```
- **Kemungkinan Respon Error:**
  - `400 Bad Request` (Tidak ada shift aktif yang terbuka)

---

#### GET `/api/v1/shifts/current`
- **Akses:** Cashier
- **Deskripsi:** Mengetahui status shift aktif kasir yang sedang login.
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Informasi shift aktif berhasil diambil",
  "data": {
    "id": "99e1a2b3-c4d5-4e6f-a7b8-c9d0e1f2a3b4",
    "starting_cash": 20000000,
    "start_time": "2026-10-27T08:00:00Z",
    "status": "open",
    "current_cash_sales": 32000000,
    "current_non_cash_sales": 56000000
  }
}
```
- **Kemungkinan Respon Error:**
  - `404 Not Found` (Kasir belum membuka shift)

---

#### GET `/api/v1/shifts`
- **Akses:** Owner
- **Query Params:**
  - `user_id` (string UUID, opsional)
  - `status` (`open`/`closed`, opsional)
  - `page` (int, default: 1)
  - `per_page` (int, default: 10)
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Daftar shift berhasil diambil",
  "data": [
    {
      "id": "99e1a2b3-c4d5-4e6f-a7b8-c9d0e1f2a3b4",
      "cashier_name": "Budi Santoso",
      "start_time": "2026-10-27T08:00:00Z",
      "end_time": "2026-10-27T16:00:00Z",
      "starting_cash": 20000000,
      "expected_ending_cash": 65000000,
      "actual_ending_cash": 65000000,
      "cash_difference": 0,
      "status": "closed"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 10,
    "total": 15,
    "total_pages": 2
  }
}
```

---

### 4.10 Orders (Cashier)

#### GET `/api/v1/orders/promos/active`
- **Akses:** Cashier
- **Deskripsi:** Mendapatkan daftar promo yang sedang aktif dan valid untuk diaplikasikan ke pesanan kasir.
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Daftar promo aktif berhasil diambil",
  "data": [
    {
      "id": "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
      "code": "SENINKOPI",
      "name": "Diskon Senin Semangat 20%",
      "discount_type": "percentage",
      "discount_value": 20,
      "max_discount": 2000000,
      "min_spend": 5000000
    }
  ]
}
```

---

#### POST `/api/v1/orders`
- **Akses:** Cashier
- **Deskripsi:** Membuat invoice order baru.
- **Aturan Bisnis Ditegakkan:**
  - Kasir **wajib** memiliki shift yang sedang `open`.
  - Jika `order_type: dine_in`, `table_id` wajib disertakan dan harus berstatus `available`.
  - Promo maksimal 1 kupon dan wajib memenuhi syarat `min_spend`.
  - Kuantitas item tidak boleh melebihi sisa `current_stock`.
- **Request Body:**
```json
{
  "order_type": "dine_in",
  "table_id": "a90b4d4f-56f0-4660-845f-4a3cebf27e01",
  "promo_code": "SENINKOPI",
  "customer_notes": "Meja dekat jendela",
  "items": [
    {
      "product_id": "e4a2e582-7f72-46a2-9d33-3dcb295a643b",
      "quantity": 2,
      "notes": "Satu less ice, satu normal"
    },
    {
      "product_id": "d8b3941a-1d57-41eb-8785-cfc09f3032ef",
      "quantity": 1,
      "notes": "Extra hot"
    }
  ]
}
```
- **Response Sukses (201 Created):**
```json
{
  "success": true,
  "message": "Pesanan berhasil dibuat",
  "data": {
    "id": "00f1a2b3-4c5d-6e7f-8a9b-0c1d2e3f4a5b",
    "invoice_number": "INV-20261027-0001",
    "order_type": "dine_in",
    "table_number": "Meja 01",
    "subtotal": 8100000,
    "discount_amount": 1620000,
    "tax_amount": 0,
    "total_amount": 6480000,
    "status": "pending",
    "items": [
      {
        "product_id": "e4a2e582-7f72-46a2-9d33-3dcb295a643b",
        "product_name": "Iced Caffe Latte",
        "price": 2800000,
        "quantity": 2,
        "subtotal": 5600000,
        "notes": "Satu less ice, satu normal"
      },
      {
        "product_id": "d8b3941a-1d57-41eb-8785-cfc09f3032ef",
        "product_name": "Hot Cappuccino",
        "price": 2500000,
        "quantity": 1,
        "subtotal": 2500000,
        "notes": "Extra hot"
      }
    ],
    "created_at": "2026-10-27T10:30:00Z"
  }
}
```
- **Kemungkinan Respon Error:**
  - `400 Bad Request` (Kasir belum buka shift: `"ERR_NO_ACTIVE_SHIFT"`)
  - `400 Bad Request` (Stok produk tidak mencukupi: `"ERR_INSUFFICIENT_STOCK"`)
  - `400 Bad Request` (Subtotal tidak memenuhi min_spend promo)
  - `409 Conflict` (Meja sudah `occupied`)
  - `422 Unprocessable Entity` (Item kosong atau kuantitas $\le 0$)

---

#### GET `/api/v1/orders/shift/current`
- **Akses:** Cashier
- **Deskripsi:** Mendapatkan daftar transaksi yang dibuat pada shift kasir yang aktif saat ini.
- **Query Params:**
  - `search` (nomor invoice, opsional)
  - `status` (`pending`/`paid`/`cancelled`, opsional)
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Riwayat transaksi shift berhasil diambil",
  "data": [
    {
      "id": "00f1a2b3-4c5d-6e7f-8a9b-0c1d2e3f4a5b",
      "invoice_number": "INV-20261027-0001",
      "table_number": "Meja 01",
      "order_type": "dine_in",
      "total_amount": 6480000,
      "payment_method": "midtrans_snap",
      "status": "pending",
      "created_at": "2026-10-27T10:30:00Z"
    }
  ]
}
```

---

#### GET `/api/v1/orders/:id`
- **Akses:** Cashier, Owner
- **Deskripsi:** Mengambil detail lengkap nota transaksi.
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Detail order berhasil diambil",
  "data": {
    "id": "00f1a2b3-4c5d-6e7f-8a9b-0c1d2e3f4a5b",
    "invoice_number": "INV-20261027-0001",
    "shift_id": "99e1a2b3-c4d5-4e6f-a7b8-c9d0e1f2a3b4",
    "cashier_name": "Budi Santoso",
    "table_number": "Meja 01",
    "order_type": "dine_in",
    "subtotal": 8100000,
    "discount_amount": 1620000,
    "total_amount": 6480000,
    "payment_method": "midtrans_snap",
    "status": "pending",
    "items": [
      {
        "product_name": "Iced Caffe Latte",
        "price": 2800000,
        "quantity": 2,
        "subtotal": 5600000,
        "notes": "Satu less ice, satu normal"
      }
    ],
    "created_at": "2026-10-27T10:30:00Z"
  }
}
```
- **Kemungkinan Respon Error:**
  - `404 Not Found` (Order tidak ditemukan)

---

#### POST `/api/v1/orders/:id/cancel`
- **Akses:** Cashier, Owner
- **Deskripsi:** Membatalkan pesanan yang belum dibayar (`pending`). Status meja akan dilepaskan kembali menjadi `available`. Sesuai aturan bisnis, order yang sudah dibayar tidak dapat dibatalkan di endpoint ini.
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Pesanan berhasil dibatalkan"
}
```
- **Kemungkinan Respon Error:**
  - `400 Bad Request` (Pesanan sudah berstatus `paid` dan tidak dapat dibatalkan)

---

### 4.11 Payments

#### POST `/api/v1/payments/cash`
- **Akses:** Cashier
- **Deskripsi:** Eksekusi pembayaran tunai (*cash*). Sistem memvalidasi nominal uang kas diterima, menghitung uang kembalian (*change*), memotong stok produk secara permanen, melepaskan status meja menjadi `available`, dan mengubah status order menjadi `paid`.
- **Request Body:**
```json
{
  "order_id": "00f1a2b3-4c5d-6e7f-8a9b-0c1d2e3f4a5b",
  "cash_received": 7000000
}
```
*(Total tagihan order: `6480000`, kas diterima: `7000000`)*
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Pembayaran tunai berhasil",
  "data": {
    "payment_id": "a89c1b2d-3e4f-5a6b-7c8d-9e0f1a2b3c4d",
    "order_id": "00f1a2b3-4c5d-6e7f-8a9b-0c1d2e3f4a5b",
    "invoice_number": "INV-20261027-0001",
    "total_amount": 6480000,
    "cash_received": 7000000,
    "cash_change": 520000,
    "payment_method": "cash",
    "status": "paid",
    "settlement_time": "2026-10-27T10:35:00Z"
  }
}
```
- **Kemungkinan Respon Error:**
  - `400 Bad Request` (Uang diterima kurang dari total belanja: `"cash_received < total_amount"`)
  - `400 Bad Request` (Order sudah lunas atau sudah dibatalkan)

---

#### POST `/api/v1/payments/midtrans/charge`
- **Akses:** Cashier
- **Deskripsi:** Menginisialisasi pembayaran nontunai dengan membuat Midtrans Snap Token dan Redirect URL.
- **Aturan Bisnis:** Pada tahap ini stok **belum** terpotong, menunggu webhook status `settlement`.
- **Request Body:**
```json
{
  "order_id": "00f1a2b3-4c5d-6e7f-8a9b-0c1d2e3f4a5b"
}
```
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Snap token Midtrans berhasil dibuat",
  "data": {
    "payment_id": "d12e3f4a-5b6c-7d8e-9f0a-1b2c3d4e5f6a",
    "order_id": "00f1a2b3-4c5d-6e7f-8a9b-0c1d2e3f4a5b",
    "snap_token": "8b51d8b2-4d2c-4976-904d-17e997087bc9",
    "redirect_url": "https://app.sandbox.midtrans.com/snap/v2/vtweb/8b51d8b2-4d2c-4976-904d-17e997087bc9",
    "amount": 6480000,
    "status": "pending"
  }
}
```
- **Kemungkinan Respon Error:**
  - `400 Bad Request` (Order sudah lunas atau dibatalkan)
  - `502 Bad Gateway` (Gagal berkomunikasi dengan Midtrans API Server)

---

#### POST `/api/v1/payments/midtrans/webhook`
- **Akses:** Public (Midtrans HTTP Notification Callback)
- **Deskripsi:** Webhook penerima status transaksi dari payment gateway Midtrans.
- **Aturan Bisnis & Idempotensi:**
  1. Backend memvalidasi `signature_key`:  
     `SHA512(order_id + status_code + gross_amount + ServerKey)`
  2. Jika valid dan status transaksi adalah `settlement` atau `capture`:
     - Periksa apakah order sudah diproses sebelumnya (idempotency check). Jika sudah, segera return `200 OK` tanpa mutasi stok ganda.
     - Jika baru: ubah status order menjadi `paid`, potong stok produk dengan tipe `SALE` di `stock_movements`, dan ubah status meja kembali menjadi `available`.
  3. Jika status transaksi `expire`, `cancel`, atau `deny`: ubah status order menjadi `cancelled`/`expired` dan bebaskan meja menjadi `available` (stok tidak dipotong).
- **Request Body (Midtrans Webhook Payload):**
```json
{
  "transaction_time": "2026-10-27 10:40:00",
  "transaction_status": "settlement",
  "transaction_id": "c7112028-3e4b-47e2-a7d6-3e58913b8f14",
  "status_message": "midtrans payment notification",
  "status_code": "200",
  "signature_key": "9c1a5b42d4...",
  "payment_type": "qris",
  "order_id": "00f1a2b3-4c5d-6e7f-8a9b-0c1d2e3f4a5b",
  "gross_amount": "64800.00",
  "currency": "IDR"
}
```
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Notifikasi webhook Midtrans berhasil diproses"
}
```
- **Kemungkinan Respon Error:**
  - `400 Bad Request` (Invalid Signature Key atau format payload tidak valid)
  - `404 Not Found` (Order ID tidak terdaftar)

---

#### GET `/api/v1/payments/order/:order_id`
- **Akses:** Cashier, Owner
- **Deskripsi:** Mengecek status pembayaran pesanan secara langsung (polling status kasir jika webhook terlambat).
- **Response Sukses (200 OK):**
```json
{
  "success": true,
  "message": "Status pembayaran berhasil diambil",
  "data": {
    "order_id": "00f1a2b3-4c5d-6e7f-8a9b-0c1d2e3f4a5b",
    "invoice_number": "INV-20261027-0001",
    "payment_method": "midtrans_snap",
    "amount": 6480000,
    "payment_type": "qris",
    "status": "settlement",
    "settlement_time": "2026-10-27T10:40:00Z"
  }
}
```
- **Kemungkinan Respon Error:**
  - `404 Not Found` (Pembayaran belum diinisiasi untuk order tersebut)
