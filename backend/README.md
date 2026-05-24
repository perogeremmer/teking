# Ayomendaki Admin

> Self-hosted open trip management software untuk para open trip operator.

Dibangun dengan **Go** + **SQLite** + **HTMX** + **Tailwind CSS**. Single binary, zero dependencies runtime. Tinggal download, jalanin, dan pakai.

## Fitur

- **Dashboard** — statistik trip, booking, revenue, kuota alert
- **Manajemen Trip** — CRUD trip dengan jadwal, rute, harga
- **Jadwal Trip** — Calendar view custom dengan multi-schedule overlap
- **Manajemen Booking** — daftar booking, update status (confirm/complete/cancel)
- **Meeting Points** — kelola titik kumpul dengan Leaflet map
- **Paket & Fasilitas** — kelola paket perjalanan dengan fasilitas
- **Laporan Revenue** — grafik pendapatan per bulan dan per trip
- **Profil Operator** — atur informasi dan kontak

## Tech Stack

| Layer | Teknologi |
|-------|-----------|
| Backend | Go 1.22+ (stdlib `net/http`) |
| Database | SQLite (via `modernc.org/sqlite`, pure Go) |
| Templating | Go `html/template` |
| CSS | Tailwind CSS (standalone CLI) |
| Frontend | HTMX + vanilla JS seminimal mungkin |
| Charts | Chart.js |
| Calendar | Custom vanilla JS (dayjs) |
| Map | Leaflet.js + Nominatim geocoding |

## Cara Pakai

### Prerequisites

- Go 1.22+
- Make (opsional)

### Development

```bash
# Clone repo
git clone https://github.com/ayomendaki/ayomendaki-admin.git
cd ayomendaki-admin

# Install dependencies
go mod tidy

# Build Tailwind CSS
make css

# Jalankan development server
make dev
```

### Production (Self-Hosted)

```bash
# Build single binary
make build

# Jalankan (SQLite database akan dibuat otomatis)
./ayomendaki-admin

# Default port: 8080
# Default credentials: admin / admin123 (ganti setelah login pertama)
```

### Environment Variables

| Variable | Default | Deskripsi |
|----------|---------|-----------|
| `PORT` | `8080` | Port server |
| `DB_PATH` | `data/ayomendaki.db` | Path SQLite database |

## Struktur Proyek

```
backend/
├── main.go               # Entry point
├── internal/             # Go package internal
│   ├── server/           # HTTP server, router, middleware
│   ├── handler/          # HTTP handlers per fitur
│   ├── model/            # Struct definitions + DB queries
│   ├── database/         # SQLite init + seeder
│   └── auth/             # Autentikasi (bcrypt + session)
├── web/                  # Frontend assets (di-embed ke binary)
│   ├── templates/        # Go html/template files
│   └── static/           # CSS, JS, images
├── Makefile              # Build commands
└── tailwind.config.js
```

## Lisensi

MIT — silakan digunakan, dimodifikasi, dan didistribusikan kembali.
