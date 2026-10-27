# Product Requirement Document (PRD)
## Sistem Point of Sales (POS) Coffee Shop

**Tech Stack:** Go (Golang), MySQL, Midtrans Snap  
**Versi Dokumen:** 1.0  
**Tanggal:** 27 Oktober 2026  
**Status:** Approved / Final  

```text
PRD.md
├── Deskripsi Produk
├── User Roles
│   ├── Owner
│   └── Cashier
├── Fitur Owner
│   ├── Manajemen Produk (+ acceptance criteria)
│   ├── Manajemen Kategori (+ acceptance criteria)
│   ├── Manajemen Stok (+ acceptance criteria)
│   ├── Manajemen Meja (+ acceptance criteria)
│   ├── Manajemen User Cashier (+ acceptance criteria)
│   ├── Dashboard Laporan (+ acceptance criteria)
│   ├── Manajemen Promo (+ acceptance criteria)
│   └── Export CSV (+ acceptance criteria)
├── Fitur Cashier
│   ├── Shift Management (+ acceptance criteria)
│   ├── Buat Transaksi (+ acceptance criteria)
│   ├── Apply Promo (+ acceptance criteria)
│   ├── Checkout Midtrans (+ acceptance criteria)
│   └── Riwayat Transaksi (+ acceptance criteria)
├── Business Rules
└── Out of Scope
```

---

## 1. Deskripsi Produk

### 1.1 Ringkasan Produk
Sistem Point of Sales (POS) Coffee Shop adalah sistem aplikasi kasir dan operasional berbasis web/API yang dirancang untuk mengotomatisasi pencatatan penjualan, manajemen inventaris stok, pembagian shift kerja kasir, serta analitik performa kedai kopi secara terpusat.

### 1.2 Target Pengguna & Segmen
Produk ini ditujukan bagi **Coffee Shop Skala Menengah** dengan profil operasional:
- Melayani pelanggan santap di tempat (*dine-in*) berbasis meja dan pesan bawa pulang (*takeaway*).
- Memiliki 1–2 kasir operasional yang bekerja dalam sistem pergantian shift (*shift-based*).
- Membutuhkan sistem pembayaran fleksibel: tunai (*cash*) serta pembayaran digital modern (QRIS, GoPay, ShopeePay, Virtual Account via Midtrans Snap).
- Pemilik bisnis (*Owner*) membutuhkan visibilitas penuh atas stok bahan/menu dan analitik pendapatan harian/mingguan/bulanan dari jarak jauh.

---

## 2. User Roles

Sistem menerapkan Role-Based Access Control (RBAC) dengan dua peran utama:

### 2.1 Owner
Pemilik kedai kopi yang memiliki otoritas administratif penuh atas sistem:
- Mengelola katalog produk, kategori, harga, dan resep dasar/stok.
- Mengatur tata kelola meja dan akun staf kasir.
- Menentukan strategi promo dan diskon penjualan.
- Mengakses dasbor analitik bisnis serta mengekspor data laporan keuangan ke format CSV.

### 2.2 Cashier
Staf operasional kasir yang bertugas di garda depan (*front-of-house*):
- Bertanggung jawab atas modal kas laci kerja melalui proses buka dan tutup shift.
- Memproses pesanan pelanggan (*dine-in* maupun *takeaway*).
- Mengaplikasikan promo aktif dan memproses pembayaran (Tunai atau Midtrans Snap).
- Melihat rangkuman dan riwayat transaksi pada sesi shift aktifnya.

#### Matriks Hak Akses (RBAC)

| Modul / Kemampuan | Owner | Cashier |
| :--- | :---: | :---: |
| Autentikasi & Akun Pribadi | Ya | Ya |
| Manajemen User Cashier | Full Akses (CRUD) | Tidak |
| Manajemen Produk & Kategori | Full Akses (CRUD, Soft Delete) | View Menu Aktif |
| Manajemen Stok & Adjustment | Full Akses (Koreksi & Audit Log) | View Sisa Stok |
| Manajemen Meja | Full Akses (CRUD Meja) | View & Pilih Meja |
| Manajemen Promo & Diskon | Full Akses (CRUD Promo) | Apply Promo Aktif |
| Dashboard Laporan & Analitik | Full Akses | Tidak |
| Export Laporan CSV | Full Akses | Tidak |
| Shift Management (Buka/Tutup Kas) | View Rekap Shift | Eksekusi Shift Sendiri |
| Pembuatan & Pembayaran Transaksi | View / Audit Transaksi | Eksekusi Transaksi |
| Riwayat Transaksi Shift Aktif | Semua Transaksi | Shift Sendiri |

---

## 3. Fitur Owner

### 3.1 Manajemen Produk
Mengelola seluruh daftar makanan dan minuman yang dijual di kedai kopi.
- **Deskripsi:**
  - Operasi CRUD (Create, Read, Update, Delete) data produk (nama, SKU, deskripsi, harga jual, harga dasar/HPP).
  - Unggah foto produk (format JPEG/PNG/WebP, max 2MB).
  - Pengaturan status produk (`aktif` / `nonaktif`). Produk nonaktif tidak akan muncul di layar kasir.
  - Penghapusan produk menggunakan **Soft Delete** (`deleted_at`), sehingga data transaksi lampau tetap utuh.
- **Acceptance Criteria:**
  - [ ] Owner dapat menambah, mengedit, dan me-nonaktifkan produk.
  - [ ] Produk yang di-*soft delete* hilang dari katalog kasir namun tetap tercatat pada riwayat transaksi historis.
  - [ ] Validasi menolak input bila nama kosong, harga jual $\le 0$, atau format foto tidak sesuai.

### 3.2 Manajemen Kategori
Mengelompokkan produk ke dalam kategori menu (misal: *Signature Coffee*, *Manual Brew*, *Tea & Non-Coffee*, *Pastry*, *Main Course*).
- **Deskripsi:** CRUD kategori produk, penetapan urutan tampilan (*sort order*), dan status aktif kategori.
- **Acceptance Criteria:**
  - [ ] Owner dapat membuat, mengubah nama, dan menghapus kategori.
  - [ ] Kategori yang masih menampung produk aktif tidak dapat dihapus sebelum produk dialihkan atau dinonaktifkan.

### 3.3 Manajemen Stok
Kontrol inventori produk jadi dan pelacakan riwayat mutasi stok.
- **Deskripsi:**
  - Input stok awal produk dan penyesuaian stok manual (*stock adjustment* untuk barang masuk, selisih fisik, atau *waste*/rusak).
  - Pengaturan ambang batas stok minimum (*low stock alert*).
  - Pencatatan otomatis setiap pergerakan stok ke dalam tabel riwayat mutasi stok (*stock movement log*) mencakup: timestamp, tipe (`IN`, `OUT`, `SALE`, `ADJUSTMENT`), kuantitas, stok sebelum & sesudah, serta user penanggung jawab.
- **Acceptance Criteria:**
  - [ ] Setiap pembayaran transaksi sukses otomatis memotong stok dan tercatat sebagai tipe `SALE`.
  - [ ] Riwayat pergerakan stok dapat difilter berdasarkan rentang tanggal dan nama produk.
  - [ ] Penyesuaian stok manual wajib menyertakan alasan keterangan (misal: "tumpah / expired").

### 3.4 Manajemen Meja
Pengaturan denah dan ketersediaan nomor meja untuk pelanggan santap di tempat (*dine-in*).
- **Deskripsi:**
  - Tambah, ubah, dan hapus nomor/identitas meja serta kapasitas kursi (contoh: "Meja 01", "Meja 02 (Outdoor)").
  - Pemantauan status meja secara *real-time* (`Available`, `Occupied`).
- **Acceptance Criteria:**
  - [ ] Owner dapat menambah dan mengatur daftar meja; nomor meja harus unik.
  - [ ] Meja yang sedang dalam status `Occupied` tidak dapat dihapus.

### 3.5 Manajemen User Cashier
Pengelolaan akun staf kasir yang bertugas di coffee shop.
- **Deskripsi:**
  - Pendaftaran kasir baru (Nama Lengkap, Username/Email, Password hash bcrypt).
  - Pengubahan status kasir (`aktif` / `nonaktif`).
  - Fitur reset password staf kasir oleh Owner.
- **Acceptance Criteria:**
  - [ ] Akun kasir yang dinonaktifkan tidak dapat login ke aplikasi.
  - [ ] Password kasir tersimpan secara aman dalam format hash bcrypt di database MySQL.

### 3.6 Dashboard Laporan
Visualisasi metrik penjualan dan kesehatan finansial kedai secara ringkas dan informatif.
- **Deskripsi:**
  - Ringkasan pendapatan (*Total Revenue*): Harian, Mingguan, dan Bulanan.
  - Grafik tren volume pesanan dan nominal penjualan per periode waktu.
  - Daftar Produk Terlaris (*Top 5 / Top 10 Best Selling Products* berdasarkan kuantitas dan omzet).
  - Rekap total transaksi dan pendapatan yang dihasilkan oleh masing-masing kasir.
- **Acceptance Criteria:**
  - [ ] Angka pendapatan pada dashboard hanya menghitung transaksi yang berstatus lunas (`PAID` / `SETTLED`).
  - [ ] Kalkulasi waktu mengikuti zona waktu kedai lokal (WIB/UTC+7).

### 3.7 Manajemen Promo
Pengaturan program promosi dan potongan harga.
- **Deskripsi:**
  - Pembuatan promo bertipe **Persentase** (contoh: 15% dengan batas potongan maksimum Rp 20.000).
  - Pembuatan promo bertipe **Nominal Tetap** (contoh: Potongan langsung Rp 10.000).
  - Konfigurasi aturan promo: minimal pembelian (*minimum spend*), tanggal berlaku awal & akhir, serta kuota penggunaan.
  - Pengaktifan dan penonaktifan promo secara instan (*toggle active*).
- **Acceptance Criteria:**
  - [ ] Promo yang berada di luar rentang tanggal aktif otomatis tidak dapat digunakan kasir.
  - [ ] Diskon persentase yang melebihi batas nominal maksimum akan otomatis dipotong sesuai nilai batas maksimum tersebut.

### 3.8 Export CSV
Fitur unduh data mentah transaksi untuk rekonsiliasi akuntansi eksternal.
- **Deskripsi:**
  - Filter ekspor berdasarkan rentang tanggal (*start date* s/d *end date*).
  - File `.csv` yang diunduh memuat rincian: No. Invoice, Waktu Transaksi, Nama Kasir, Tipe Pesanan (Dine-in / Takeaway), No. Meja, Subtotal, Nilai Diskon, Total Akhir, Metode Bayar, dan Status Transaksi.
- **Acceptance Criteria:**
  - [ ] File CSV ter-generate secara presisi dengan separator standar dan format UTF-8.
  - [ ] Data yang diekspor sesuai dengan rentang filter tanggal yang dipilih.

---

## 4. Fitur Cashier

### 4.1 Shift Management
Pencatatan akuntabilitas perputaran uang tunai pada laci kasir (*cash drawer*).
- **Deskripsi:**
  - **Buka Shift (*Open Shift*):** Kasir menginput modal kas awal (*starting float*) sebelum memulai transaksi harian.
  - **Tutup Shift (*Close Shift*):** Kasir menghitung fisik uang tunai di laci pada akhir shift, melihat rekapan total penjualan kas vs nontunai dari sistem, serta mencatat selisih (*variance over/short*).
- **Acceptance Criteria:**
  - [ ] Kasir diblokir dari membuat transaksi baru apabila belum melakukan *Open Shift*.
  - [ ] Satu user kasir hanya dapat membuka satu shift aktif pada satu waktu.
  - [ ] Sistem menghitung dan mendokumentasikan selisih antara uang tunai sistem vs hitungan fisik kasir saat *Close Shift*.

### 4.2 Buat Transaksi
Alur cepat pencatatan pesanan menu pelanggan.
- **Deskripsi:**
  - Pemilihan tipe pesanan: *Dine-In* (wajib memilih meja berstatus `Available`) atau *Takeaway*.
  - Pemilihan menu produk aktif dengan dukungan filter kategori dan pencarian nama.
  - Penyesuaian kuantitas (*quantity*) dan penambahan catatan khusus per item (misal: "Oatmilk", "Less Sugar").
  - Tampilan ringkasan pesanan (*order summary*): subtotal item, potongan diskon, dan grand total.
- **Acceptance Criteria:**
  - [ ] Status meja *dine-in* yang dipilih otomatis berubah menjadi `Occupied`.
  - [ ] Sistem menolak penambahan item jika kuantitas yang dipesan melebihi sisa stok yang tersedia.

### 4.3 Apply Promo
Penerapan diskon resmi ke dalam keranjang transaksi belanja.
- **Deskripsi:**
  - Memilih promo dari daftar promosi yang sedang berlaku dan memenuhi kriteria belanja.
  - Menghapus kembali promo dari pesanan jika pelanggan membatalkan penggunaan voucher sebelum pembayaran.
- **Acceptance Criteria:**
  - [ ] Promo berhasil diaplikasikan jika nilai subtotal memenuhi syarat *minimum spend*.
  - [ ] Nilai diskon mengurangi grand total tagihan secara akurat.

### 4.4 Checkout Midtrans
Penyelesaian transaksi pembayaran melalui Midtrans Snap maupun Tunai.
- **Deskripsi:**
  - **Pembayaran Nontunai:** Sistem memanggil Midtrans Snap API untuk memunculkan modal/token pembayaran digital (QRIS, GoPay, ShopeePay, Virtual Account).
  - Notifikasi webhook dari Midtrans mengonfirmasi pembayaran lunas secara asinkron dan memperbarui status pesanan menjadi `PAID` / `SETTLED`.
  - **Pembayaran Tunai:** Kasir memasukkan nominal uang yang diterima dari pelanggan, sistem menghitung nominal kembalian (*change*), dan transaksi langsung berstatus `PAID`.
- **Acceptance Criteria:**
  - [ ] Snap token/modal berhasil di-generate dengan rincian invoice dan nilai total tagihan yang sesuai.
  - [ ] Stok produk resmi dipotong dan status meja kembali menjadi `Available` sesaat setelah webhook Midtrans `settlement`/`capture` berhasil diterima.
  - [ ] Webhook bersifat *idempotent* (pemanggilan berulang tidak menyebabkan duplikasi pengurangan stok).

### 4.5 Riwayat Transaksi
Pemantauan dan verifikasi transaksi yang diproses selama shift kasir berjalan.
- **Deskripsi:**
  - Menampilkan daftar transaksi yang dibuat pada sesi shift aktif kasir tersebut (No. Invoice, Waktu, Meja, Metode Pembayaran, Total, Status).
  - Fitur pencarian cepat berdasarkan nomor invoice.
  - Melihat rincian detail pesanan dan cetak ulang tanda terima/struk transaksi.
- **Acceptance Criteria:**
  - [ ] Kasir hanya dapat melihat riwayat transaksi yang dibuat dalam sesi shift aktifnya sendiri.
  - [ ] Kasir dapat membuka detail tiap invoice untuk melihat rincian item pesanan.

---

## 5. Business Rules

Sistem wajib menegakkan aturan integritas bisnis berikut:

1. **Shift Kasir:**
   - **Cashier harus buka shift (*Open Shift*)** dengan mencatat modal kas awal sebelum bisa membuat transaksi baru.
   - Setiap transaksi penjualan wajib terikat secara permanen ke `shift_id` yang sedang aktif.

2. **Integritas & Pemotongan Stok:**
   - **Stok hanya berkurang setelah webhook Midtrans confirmed (`settlement` / `capture`) diterima.** Sistem tidak boleh mengurangi stok ketika pesanan masih dalam tahap draft atau berstatus `pending`. Untuk transaksi tunai, stok berkurang saat kasir memproses pembayaran lunas.
   - Stok fisik produk tidak boleh bernilai negatif ($< 0$).
   - Jika pembayaran digital kedaluwarsa (*expired*) atau gagal (*cancel/deny*), transaksi dibatalkan, status meja dilepaskan, dan stok tetap utuh.

3. **Imutabilitas Transaksi Pasca-Checkout:**
   - **Order yang sudah checkout tidak bisa diubah itemnya** (dilarang menambah, menghapus, atau mengganti kuantitas produk serta catatan pesanan). Jika terjadi perubahan sebelum pembayaran lunas, transaksi harus dibatalkan dan dibuatkan transaksi baru.

4. **Ketentuan Promo & Diskon:**
   - **Satu transaksi hanya bisa menggunakan satu promo** (tidak diperkenankan menggabungkan multi-promo dalam satu invoice).
   - Nilai potongan promo tidak boleh melebihi nilai subtotal tagihan (grand total $\ge 0$).
   - Diskon persentase dibatasi oleh batas nominal maksimal (*max discount*) jika diatur.

5. **Siklus Hidup Meja (*Table Lifecycle*):**
   - Meja yang berstatus `Occupied` tidak dapat dipilih untuk pesanan *dine-in* baru lainnya.
   - Status meja otomatis kembali menjadi `Available` setelah transaksi selesai dibayar lunas (`PAID` / `SETTLED`) atau jika pesanan dibatalkan.

6. **Keamanan & Konsistensi Midtrans:**
   - Server Go wajib memvalidasi *Signature Key* webhook Midtrans:  
     `SHA512(order_id + status_code + gross_amount + ServerKey)`
   - Handler webhook harus bersifat *idempotent* untuk mencegah pemotongan stok ganda akibat pengiriman ulang notifikasi webhook.

7. **Integritas Data Transaksional (*Soft Delete*):**
   - Data produk, kategori, dan meja yang memiliki riwayat transaksi tidak boleh di-*hard delete* dari basis data MySQL. Seluruh entitas dilindungi menggunakan mekanisme `deleted_at`.

---

## 6. Out of Scope

Fitur-fitur berikut berada di luar cakupan pengembangan versi 1.0 ini:

1. **Multi-Outlet / Multi-Cabang:** Sistem saat ini hanya dikhususkan untuk 1 gerai fisik (*single-tenant*).
2. **Loyalty System & Customer Membership:** Belum ada sistem poin reward, level membership pelanggan, atau integrasi CRM.
3. **Split Bill & Merge Bill:** Tidak menyediakan pemisahan satu tagihan meja ke banyak pembayaran terpisah atau penggabungan beberapa meja menjadi satu tagihan.
4. **Driver Raw Thermal ESC/POS Socket:** Sistem tidak mengontrol hardware printer via koneksi soket raw; pencetakan struk menggunakan standar Web Print API browser.
5. **Integrasi Online Delivery:** Tidak terintegrasi secara langsung dengan API GoFood, GrabFood, atau ShopeeFood.
6. **Manajemen Bahan Baku Mentah (Recipe / COGS BOM):** Inventaris stok dicatat per unit produk jadi (contoh: 1 Cup Americano), bukan berbasis pengurangan gramatur biji kopi atau mililiter sirup.
7. **Kitchen Display System (KDS) Real-time WebSocket:** Tidak menyediakan antarmuka interaktif live-stream di stasiun barista dapur.
