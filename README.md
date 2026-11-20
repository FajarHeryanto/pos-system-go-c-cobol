# 🛒 POS System — Sistem Informasi Penjualan (Point of Sales)

**Final Project — Algoritma dan Pemrograman 3**

Sistem POS terintegrasi yang menggabungkan **Go** (web server & REST API), **C** (pengolahan data dengan manajemen memori dinamis), dan **COBOL** (pelaporan bisnis) dalam satu aplikasi berbasis web.

---

## 📁 Struktur Proyek

```
pos-system/
├── go/                          # Backend Go (Web Server + REST API)
│   ├── main.go                  # Entry point, router Gin
│   ├── go.mod                   # Go module definition
│   ├── models/                  # Struct data (User, Product, Transaction)
│   ├── handlers/                # API handler (auth, CRUD, process, report)
│   ├── middleware/               # JWT auth + CORS middleware
│   ├── database/                # SQLite init, migrasi, seeding
│   └── utils/                   # Validator, CSV exporter, hash
├── c/                           # Modul C (Pengolahan Data)
│   ├── processor.h              # Header file (struct + deklarasi fungsi)
│   ├── processor.c              # Implementasi (malloc/calloc/realloc/free)
│   └── Makefile
├── cobol/                       # Modul COBOL (Pelaporan)
│   ├── reporter.cob             # Program GnuCOBOL (3 jenis laporan)
│   └── Makefile
├── web/                         # Frontend Web
│   ├── login.html               # Halaman login
│   ├── index.html               # Dashboard
│   ├── products.html            # Manajemen produk
│   ├── transactions.html        # Transaksi penjualan
│   ├── reports.html             # Laporan COBOL
│   ├── css/style.css            # Custom stylesheet
│   └── js/                      # JavaScript (api.js, dashboard.js, dll)
├── data/                        # Data files (CSV, DB, laporan)
│   ├── transactions.csv         # Go → C
│   ├── processed_data.csv       # C → COBOL
│   ├── report_output.txt        # Output COBOL
│   └── pos_database.db          # SQLite database (auto-generated)
├── Makefile                     # Makefile utama (compile semua)
└── README.md
```

---

## ⚙️ Prasyarat

| Tool | Versi | Keterangan |
|------|-------|------------|
| **Go** | >= 1.21 | Backend web server |
| **GCC** | >= 9.0 | Compiler untuk program C |
| **GnuCOBOL** | >= 3.0 | Compiler untuk program COBOL |
| **Git** | (opsional) | Version control |

### Instalasi Prasyarat

**Windows (dengan MSYS2):**
```bash
# Install Go dari https://go.dev/dl/
# Install MSYS2 dari https://www.msys2.org/

# Di MSYS2 terminal:
pacman -S mingw-w64-x86_64-gcc
pacman -S mingw-w64-x86_64-gnucobol

# Tambahkan ke PATH:
# C:\msys64\mingw64\bin
```

**Ubuntu/Debian:**
```bash
# Go
sudo apt update
sudo apt install golang-go

# GCC
sudo apt install build-essential

# GnuCOBOL
sudo apt install gnucobol
```

**macOS:**
```bash
# Go
brew install go

# GCC (via Xcode Command Line Tools)
xcode-select --install

# GnuCOBOL
brew install gnucobol
```

### Verifikasi Instalasi

```bash
go version        # >= go1.21
gcc --version      # GCC
cobc --version     # GnuCOBOL
```

---

## 🚀 Kompilasi & Menjalankan

### Langkah 1: Kompilasi C Processor

```bash
cd c/
make
# Output: processor (Linux/Mac) atau processor.exe (Windows)
```

Atau manual:
```bash
gcc -Wall -Wextra -std=c99 -o processor processor.c -lm
```

### Langkah 2: Kompilasi COBOL Reporter

```bash
cd cobol/
make
# Output: reporter (Linux/Mac) atau reporter.exe (Windows)
```

Atau manual:
```bash
cobc -x -o reporter reporter.cob
```

### Langkah 3: Setup & Jalankan Go Web Server

```bash
cd go/
go mod tidy      # Download dependencies
go run main.go   # Jalankan server
```

Atau build lalu jalankan:
```bash
cd go/
go build -o pos-system .
cd ..
go/pos-system        # Linux/Mac
go\pos-system.exe    # Windows
```

### Langkah Cepat (Makefile)

```bash
# Dari root proyek (pos-system/):
make          # Compile C + COBOL
make test     # Test pipeline C → COBOL
make run      # Compile semua + jalankan Go server
```

---

## 🌐 Mengakses Aplikasi

Setelah server berjalan, buka browser:

```
http://localhost:8080
```

### Akun Login Default

| Username | Password | Role |
|----------|----------|------|
| `admin` | `admin123` | Admin |
| `andi` | `kasir123` | Kasir |
| `budi` | `kasir123` | Kasir |
| `citra` | `kasir123` | Kasir |
| `dewi` | `kasir123` | Kasir |

---

## 📋 REST API Endpoint

| Method | Endpoint | Deskripsi | Auth |
|--------|----------|-----------|------|
| `POST` | `/api/auth/login` | Login | ❌ |
| `GET` | `/api/products` | Lihat semua produk | ✅ |
| `GET` | `/api/products/:id` | Lihat produk by ID | ✅ |
| `POST` | `/api/products` | Tambah produk baru | ✅ |
| `PUT` | `/api/products/:id` | Update produk | ✅ |
| `DELETE` | `/api/products/:id` | Hapus produk | ✅ |
| `GET` | `/api/transactions` | Lihat semua transaksi | ✅ |
| `POST` | `/api/transactions` | Buat transaksi baru | ✅ |
| `DELETE` | `/api/transactions/:id` | Hapus transaksi | ✅ |
| `POST` | `/api/process` | Jalankan C processor | ✅ |
| `POST` | `/api/reports/generate` | Jalankan COBOL reporter | ✅ |
| `GET` | `/api/reports/income` | Ambil laporan pendapatan | ✅ |
| `GET` | `/api/dashboard` | Data ringkasan dashboard | ✅ |

---

## 🧪 Pengujian End-to-End

### Test 1: Pipeline C + COBOL (tanpa Go)

```bash
# Dari root proyek:

# 1. Pastikan data/transactions.csv berisi data
# 2. Jalankan C processor
./c/processor          # Linux/Mac
c\processor.exe        # Windows

# 3. Periksa output
cat data/processed_data.csv

# 4. Jalankan COBOL reporter
./cobol/reporter       # Linux/Mac
cobol\reporter.exe     # Windows

# 5. Periksa laporan
cat data/report_output.txt
```

### Test 2: Full Web Application

```bash
# 1. Compile C + COBOL
make

# 2. Jalankan Go server
cd go && go run main.go

# 3. Buka browser → http://localhost:8080
# 4. Login dengan admin / admin123
# 5. Di Dashboard: lihat statistik
# 6. Di Produk: tambah/edit/hapus produk
# 7. Di Transaksi: buat transaksi baru
# 8. Di Laporan:
#    a. Klik "Proses Data" (C Engine)
#    b. Klik "Generate Laporan" (COBOL Reporter)
#    c. Lihat laporan yang tampil
#    d. Klik "Print" untuk mencetak
```

### Test 3: API dengan cURL

```bash
# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# Ambil token dari response, lalu:
TOKEN="<token_dari_response>"

# Get semua produk
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/products

# Buat transaksi
curl -X POST http://localhost:8080/api/transactions \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"product_id":"P001","quantity":2,"discount_percent":10,"payment_method":"Cash","cashier_name":"Andi"}'

# Jalankan C processor
curl -X POST -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/process

# Generate laporan COBOL
curl -X POST -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/reports/generate

# Lihat laporan
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/reports/income
```

---

## 🔄 Alur Integrasi

```
Browser (HTML/CSS/JS)
    ↕ HTTP/JSON
Go Web Server (Gin, port 8080)
    ↕ SQLite CRUD
Database (pos_database.db)
    ↓ Export CSV
transactions.csv (Go → C)
    ↓ os/exec
C Processor (malloc/calloc/realloc/free)
    ↓ Write CSV
processed_data.csv (C → COBOL)
    ↓ os/exec
COBOL Reporter (record processing)
    ↓ Write report
report_output.txt (COBOL → Go → Browser)
```

---

## 📊 Data Simulasi

Saat pertama kali dijalankan, database akan otomatis diisi dengan:

- **5 User**: 1 admin + 4 kasir
- **20 Produk**: 5 kategori (Makanan, Minuman, Elektronik, Alat Tulis, Kebutuhan RT)
- **100 Transaksi**: Periode 01-Nov-2026 s/d 19-Nov-2026

Distribusi transaksi:
- Kasir: Andi, Budi, Citra, Dewi (rotasi)
- Metode bayar: Cash 50%, Debit 30%, Credit 20%
- Diskon: 0%, 5%, 10%, 15% (rotasi)

---

## 🛡️ Keamanan

- **Autentikasi**: JWT token (24 jam expiry) + bcrypt password hashing
- **Validasi ganda**: Frontend (JavaScript) + Backend (Go)
- **Keamanan memori C**: NULL check, `strncpy` (bukan `strcpy`), `fgets` (bukan `gets`), free setelah selesai
- **Exit code C**: 0=sukses, 1=file not found, 2=invalid format, 3=memory failure

---

## 👥 Pembagian Tugas Tim

| No | Peran | Tanggung Jawab |
|----|-------|----------------|
| 1 | Go Backend Dev 1 | Server Gin, auth (login/JWT), CRUD produk, database |
| 2 | Go Backend Dev 2 | CRUD transaksi, export CSV, integrasi C & COBOL, dashboard |
| 3 | Frontend Dev | Halaman HTML, Bootstrap, JavaScript (Fetch API) |
| 4 | C Developer | processor.c (memori dinamis, kalkulasi, sorting) |
| 5 | COBOL Developer | reporter.cob (pembacaan record, rekapitulasi, laporan) |
| 6 | QA & Data | Data simulasi, pengujian integrasi end-to-end, dokumentasi |

---

## 📝 Teknologi

| Komponen | Teknologi |
|----------|-----------|
| Backend | Go 1.21+ (Gin framework) |
| Frontend | HTML5 + CSS3 + JavaScript (ES6+) |
| CSS Framework | Bootstrap 5 (CDN) |
| Chart | Chart.js (CDN) |
| Database | SQLite 3 |
| Pengolahan Data | C (GCC) |
| Pelaporan | COBOL (GnuCOBOL) |
| Autentikasi | JWT + bcrypt |
