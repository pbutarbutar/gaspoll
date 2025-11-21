# Gaspoll - Ringkasan Alur & Tech Stack

## 🚀 Ringkasan Alur Aplikasi Gaspoll (Final Version)

Gaspoll adalah aplikasi e-commerce khusus driver (Gojek, Grab, Maxim, dll) untuk membeli topup driver, pulsa, paket internet, dan voucher murah. Setiap transaksi memberikan poin atau akumulasi yang dapat ditukar menjadi voucher servis dari bengkel mitra. Bengkel dapat bergabung sebagai mitra untuk menawarkan promo dan menerima voucher yang ditukarkan driver menggunakan QR Code. Dihalaman utama ada layout page menampilkan mitra mitra yang tergabung yang sudah memiliki promo, ketika 

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
- Pilih produk → Beli → Pembayaran → Proses (QRIS, NO VA, WALLET) → Open link WhatsApp → Status pembelian dikirimkan melalui WhatsApp

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


## Log Perubhab