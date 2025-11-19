# Gaspoll

Aplikasi e-commerce untuk driver ojol (Gojek/Grab/Maxim) membeli topup driver, pulsa, dan paket internet, lengkap dengan reward voucher servis bengkel mitra. Dibangun dengan Golang + Echo (SSR).

## Fitur ringkas
- Registrasi/login driver, session cookie + JWT helper.
- Katalog produk digital (topup, pulsa, paket data) dengan checkout cepat (htmx).
- Pembayaran multi-channel (QRIS, VA, e-wallet) melalui service stub yang bisa diganti PSP asli.
- Reward poin/akumulasi order → terbit voucher QR untuk bengkel mitra.
- Dashboard driver (riwayat order + progress reward), dashboard mitra untuk redeem voucher, dan layar settlement admin.
- Observability hook menggunakan Uptrace + OpenTelemetry.

## Stack
- Backend: Go 1.21, Echo, Bun ORM (PostgreSQL), zerolog.
- Frontend: SSR `html/template`, Tailwind (CDN), htmx, Alpine.js, komponen ringan bergaya shadcn.
- Auth: gorilla/sessions (cookie), JWT helper untuk API, bcrypt hashing.
- Utilities: go-playground/validator, skip2/go-qrcode.

## Struktur
```
cmd/server          # Entrypoint
internal/config     # Env loader
internal/server     # Echo setup + routes
internal/{user,product,order,payment,reward,voucher,mitra,settlement}
internal/view       # Templating engine + embedded templates
web/assets          # Static assets (Tailwind helpers)
```

## Menjalankan lokal
```bash
cp .env.example .env
go mod tidy          # download dependency
go run ./cmd/server  # listen di :1323
```

## Docker Compose
```bash
docker compose up --build
# App: http://localhost:1323
# Uptrace: http://localhost:14318
```

## Catatan
- Service pembayaran saat ini berupa stub (`internal/payment`) untuk simulasi QRIS/VA/E-Wallet. Ganti dengan integrasi PSP pilihan Anda.
- Skema database mengikuti model Bun (tabel users, products, orders, vouchers, mitras, settlements). Tambahkan migrasi sesuai kebutuhan.
