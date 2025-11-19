# Gaspoll - Ringkasan Alur & Tech Stack

## 🚀 Ringkasan Alur Aplikasi Gaspoll (Final Version)

Gaspoll adalah aplikasi e-commerce khusus driver (Gojek, Grab, Maxim, dll) untuk membeli topup driver, pulsa, dan paket internet. Setiap transaksi memberikan poin atau akumulasi yang dapat ditukar menjadi voucher servis dari bengkel mitra. Bengkel dapat bergabung sebagai mitra untuk menawarkan promo dan menerima voucher yang ditukarkan driver menggunakan QR Code.

---

## 🔄 Alur Kerja Aplikasi

### 1. Registrasi dan Login User (Driver)
- Driver membuat akun dan mengisi data dasar (nama, nomor HP, jenis kendaraan).
- Mendapatkan dashboard:
  - Riwayat transaksi
  - Voucher reward
  - Progress reward

### 2. Driver Membeli Produk (E-Commerce Tanpa Saldo)
Produk digital yang dijual:
- Topup Gojek / Grab / Maxim
- Pulsa
- Paket Internet

Metode pembayaran:
- QRIS
- Virtual Account (VA)
- E-Wallet (Dana, ShopeePay, OVO)

Flow transaksi:
- Pilih produk → Checkout → Pembayaran → Webhook → Order Paid → Proses produk.

### 3. Perolehan Reward
Setiap transaksi menambah poin/akumulasi.
Jika memenuhi syarat → Driver mendapatkan voucher servis (QR Code).

### 4. Onboarding Mitra Bengkel
Mitra diverifikasi admin, lalu mendapatkan dashboard:
- Tambah paket promo
- Kelola voucher
- Scan QR voucher driver
- Laporan voucher redeemed
- Laporan settlement

### 5. Redeem Voucher
Driver menunjukkan QR Code ke bengkel → Mitra scan via dashboard → Sistem validasi → Voucher redeemed → Masuk laporan settlement.

### 6. Settlement
Admin melihat total voucher redeemed → Hitung kompensasi → Transfer pembayaran ke mitra → Mitra melihat laporan di dashboard.

---

# 🧱 Tech Stack Fullstack Golang Echo

## 1. Frontend (SSR Menggunakan Echo)
- Templating: `html/template`
- TailwindCSS untuk styling
- Komponen bergaya shadcn (diadaptasi manual)
- htmx (AJAX ringan)
- Alpine.js (interaksi UI ringan)
- QR Code rendering (PNG/SVG)

## 2. Backend
- Golang + Echo Framework
- PostgreSQL sebagai database utama
- Bun ORM (bun.uptrace.dev)
- Auth:
  - JWT untuk API
  - Session cookies untuk web
- Password hashing: bcrypt
- Validasi: go-playground/validator
- Payment integration:
  - QRIS
  - Virtual Account
  - E-Wallet (Dana, ShopeePay, OVO)
- QR Code generator: skip2/go-qrcode

## 3. Observability
- Uptrace + OpenTelemetry (trace query, log, error)
- Logging: zerolog / slog
- Error monitoring: opsional

## 4. DevOps & Deployment
- Docker + docker-compose (App + PostgreSQL + Uptrace)
- Reverse proxy (Nginx / Traefik)
- CI/CD (GitHub Actions)
- Env config (.env)

---

# 📦 Struktur Project (Singkat)
```
/cmd/server
/internal
  /auth
  /user
  /mitra
  /product
  /order
  /payment
  /reward
  /voucher
  /settlement
  /view
```

---

# 📘 Lisensi
Anda bebas menggunakan dokumen ini untuk kebutuhan Git repository dalam proyek Gaspoll.
