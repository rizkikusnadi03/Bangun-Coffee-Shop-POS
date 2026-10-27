# Coffee Shop Point of Sales (POS) API

Sistem backend Point of Sales (POS) berbasis RESTful API untuk kedai kopi (coffee shop) skala menengah menggunakan bahasa pemrograman **Go (Golang)**, basis data **MySQL**, dan payment gateway **Midtrans Snap**.

## 📌 Fitur Utama

- **Role Owner:** Manajemen produk, kategori, inventori stok & pergerakan stok, tata kelola meja, akun kasir, promo/diskon, dasbor laporan penjualan, dan ekspor CSV.
- **Role Cashier:** Manajemen shift kerja kasir (buka/tutup kas fisik), pembuatan transaksi pesanan (dine-in/takeaway), klaim voucher promo, serta checkout tunai maupun nontunai via Midtrans Snap.

## 🛠 Tech Stack

- **Language:** Go (Golang)
- **Database:** MySQL 8.x
- **Cache / Session:** Redis
- **Payment Gateway:** Midtrans Snap (QRIS, VA, E-Wallet)
- **Architecture:** Clean / Layered Architecture (`handler` $\rightarrow$ `service` $\rightarrow$ `repository`)

## 📁 Dokumentasi Terkait

- [PRD.md](file:///home/djan/Desktop/web-coffee/PRD.md) — Product Requirement Document
- [SCHEMA.md](file:///home/djan/Desktop/web-coffee/SCHEMA.md) — Desain Skema Basis Data MySQL & Integritas Relasi
- [API_CONTRACT.md](file:///home/djan/Desktop/web-coffee/API_CONTRACT.md) — Kontrak API RESTful v1
- [postman_collection.json](file:///home/djan/Desktop/web-coffee/postman_collection.json) — Postman Collection & Environment Presets
