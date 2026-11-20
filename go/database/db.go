// ==================================================================
// database/db.go — Koneksi, Migrasi, dan Seeding SQLite
// Sistem Informasi Penjualan (Point of Sales)
// ==================================================================
//
// Fungsi:
//   - InitDB()      : Membuka koneksi SQLite, menjalankan migrasi tabel
//   - migrateDB()   : Membuat tabel users, products, transactions + index
//   - seedDB()      : Auto-seed data awal jika database kosong:
//                      5 user, 20 produk, 100 transaksi simulasi
//
// Schema sesuai Bab 8.1 RPD.
// Data produk sesuai Bab 10.1 RPD.
// Data transaksi sesuai Bab 10.2 RPD (distribusi 4 kasir, 3 metode bayar).
// ==================================================================

package database

import (
	"database/sql"
	"fmt"
	"log"
	"math"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// DB adalah variabel global koneksi database
var DB *sql.DB

// InitDB membuka koneksi ke SQLite dan menjalankan migrasi + seeding.
func InitDB(dbPath string) *sql.DB {
	var err error
	DB, err = sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		log.Fatalf("[FATAL] Gagal membuka database: %v", err)
	}

	// Tes koneksi
	if err = DB.Ping(); err != nil {
		log.Fatalf("[FATAL] Gagal koneksi ke database: %v", err)
	}

	log.Println("[INFO] Koneksi SQLite berhasil:", dbPath)

	// Jalankan migrasi tabel
	migrateDB()

	// Jalankan seeding jika database masih kosong
	seedDB()

	return DB
}

// migrateDB membuat tabel dan index sesuai Bab 8.1 RPD.
func migrateDB() {
	schema := `
	-- Tabel User (untuk login)
	CREATE TABLE IF NOT EXISTS users (
		user_id       TEXT PRIMARY KEY,
		username      TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		full_name     TEXT NOT NULL,
		role          TEXT NOT NULL CHECK(role IN ('admin','kasir'))
	);

	-- Tabel Produk
	CREATE TABLE IF NOT EXISTS products (
		product_id    TEXT PRIMARY KEY,
		name          TEXT NOT NULL,
		category      TEXT NOT NULL,
		price         REAL NOT NULL CHECK(price > 0),
		stock         INTEGER NOT NULL DEFAULT 0 CHECK(stock >= 0),
		created_at    DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Tabel Transaksi
	CREATE TABLE IF NOT EXISTS transactions (
		transaction_id    TEXT PRIMARY KEY,
		product_id        TEXT NOT NULL,
		product_name      TEXT NOT NULL,
		quantity          INTEGER NOT NULL CHECK(quantity > 0),
		unit_price        REAL NOT NULL CHECK(unit_price > 0),
		discount_percent  REAL DEFAULT 0 CHECK(discount_percent >= 0 AND discount_percent <= 100),
		total_price       REAL NOT NULL,
		payment_method    TEXT NOT NULL CHECK(payment_method IN ('Cash','Debit','Credit')),
		transaction_date  TEXT NOT NULL,
		cashier_name      TEXT NOT NULL,
		FOREIGN KEY (product_id) REFERENCES products(product_id)
	);

	-- Index untuk performa query
	CREATE INDEX IF NOT EXISTS idx_trans_date ON transactions(transaction_date);
	CREATE INDEX IF NOT EXISTS idx_trans_product ON transactions(product_id);
	CREATE INDEX IF NOT EXISTS idx_product_category ON products(category);
	`

	_, err := DB.Exec(schema)
	if err != nil {
		log.Fatalf("[FATAL] Gagal migrasi database: %v", err)
	}
	log.Println("[INFO] Migrasi tabel berhasil")
}

// hashPassword menghasilkan bcrypt hash dari password plaintext.
func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("[FATAL] Gagal hash password: %v", err)
	}
	return string(hash)
}

// seedDB melakukan auto-seeding data awal jika tabel masih kosong.
func seedDB() {
	// Cek apakah sudah ada user
	var userCount int
	err := DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		log.Fatalf("[FATAL] Gagal cek tabel users: %v", err)
	}

	if userCount > 0 {
		log.Println("[INFO] Database sudah berisi data, skip seeding")
		return
	}

	log.Println("[INFO] Database kosong, memulai seeding data awal...")

	// ================================================================
	// SEED 1: 5 User (admin + 4 kasir)
	// ================================================================
	seedUsers()

	// ================================================================
	// SEED 2: 20 Produk sesuai Bab 10.1 RPD
	// ================================================================
	seedProducts()

	// ================================================================
	// SEED 3: 100 Transaksi simulasi sesuai Bab 10.2 RPD
	// ================================================================
	seedTransactions()

	log.Println("[INFO] Seeding data awal selesai")
}

// seedUsers menyisipkan 5 user awal.
func seedUsers() {
	adminHash := hashPassword("admin123")
	kasirHash := hashPassword("kasir123")

	users := []struct {
		id, username, hash, fullName, role string
	}{
		{"U001", "admin", adminHash, "Administrator", "admin"},
		{"U002", "andi", kasirHash, "Andi Pratama", "kasir"},
		{"U003", "budi", kasirHash, "Budi Santoso", "kasir"},
		{"U004", "citra", kasirHash, "Citra Dewi", "kasir"},
		{"U005", "dewi", kasirHash, "Dewi Lestari", "kasir"},
	}

	stmt, err := DB.Prepare("INSERT OR IGNORE INTO users (user_id, username, password_hash, full_name, role) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatalf("[FATAL] Gagal prepare statement users: %v", err)
	}
	defer stmt.Close()

	for _, u := range users {
		_, err := stmt.Exec(u.id, u.username, u.hash, u.fullName, u.role)
		if err != nil {
			log.Printf("[WARNING] Gagal insert user %s: %v", u.username, err)
		}
	}
	log.Println("[INFO] Seeding 5 user berhasil")
}

// seedProducts menyisipkan 20 produk sesuai Bab 10.1 RPD.
func seedProducts() {
	products := []struct {
		id, name, category string
		price              float64
		stock              int
	}{
		{"P001", "Nasi Goreng", "Makanan", 25000, 100},
		{"P002", "Mie Ayam", "Makanan", 20000, 80},
		{"P003", "Ayam Geprek", "Makanan", 22000, 60},
		{"P004", "Sate Ayam", "Makanan", 30000, 50},
		{"P005", "Kopi Arabica", "Minuman", 18000, 150},
		{"P006", "Teh Manis", "Minuman", 8000, 200},
		{"P007", "Jus Jeruk", "Minuman", 15000, 100},
		{"P008", "Es Coklat", "Minuman", 12000, 120},
		{"P009", "Mouse Wireless", "Elektronik", 150000, 30},
		{"P010", "Keyboard Mech", "Elektronik", 350000, 20},
		{"P011", "Headset Gaming", "Elektronik", 275000, 25},
		{"P012", "Flash Disk 32GB", "Elektronik", 85000, 50},
		{"P013", "Pulpen", "Alat Tulis", 5000, 300},
		{"P014", "Buku Tulis A5", "Alat Tulis", 8000, 250},
		{"P015", "Penghapus", "Alat Tulis", 3000, 200},
		{"P016", "Pensil 2B", "Alat Tulis", 4000, 300},
		{"P017", "Sabun Cuci", "Kebutuhan RT", 15000, 100},
		{"P018", "Tisu Wajah", "Kebutuhan RT", 12000, 150},
		{"P019", "Sikat Gigi", "Kebutuhan RT", 10000, 100},
		{"P020", "Pasta Gigi", "Kebutuhan RT", 18000, 80},
	}

	stmt, err := DB.Prepare("INSERT OR IGNORE INTO products (product_id, name, category, price, stock) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatalf("[FATAL] Gagal prepare statement products: %v", err)
	}
	defer stmt.Close()

	for _, p := range products {
		_, err := stmt.Exec(p.id, p.name, p.category, p.price, p.stock)
		if err != nil {
			log.Printf("[WARNING] Gagal insert produk %s: %v", p.name, err)
		}
	}
	log.Println("[INFO] Seeding 20 produk berhasil")
}

// seedTransactions menyisipkan 100 transaksi simulasi sesuai Bab 10.2 RPD.
// Distribusi:
//   - Kasir: 4 orang (Andi, Budi, Citra, Dewi) — rotasi
//   - Metode bayar: Cash 50%, Debit 30%, Credit 20%
//   - Diskon: 0%, 5%, 10%, 15% — rotasi
//   - Periode: 2026-11-01 s/d 2026-11-19
func seedTransactions() {
	// Daftar produk referensi (id, name, category, price)
	type prodRef struct {
		id, name, category string
		price              float64
	}
	prodList := []prodRef{
		{"P001", "Nasi Goreng", "Makanan", 25000},
		{"P002", "Mie Ayam", "Makanan", 20000},
		{"P003", "Ayam Geprek", "Makanan", 22000},
		{"P004", "Sate Ayam", "Makanan", 30000},
		{"P005", "Kopi Arabica", "Minuman", 18000},
		{"P006", "Teh Manis", "Minuman", 8000},
		{"P007", "Jus Jeruk", "Minuman", 15000},
		{"P008", "Es Coklat", "Minuman", 12000},
		{"P009", "Mouse Wireless", "Elektronik", 150000},
		{"P010", "Keyboard Mech", "Elektronik", 350000},
		{"P011", "Headset Gaming", "Elektronik", 275000},
		{"P012", "Flash Disk 32GB", "Elektronik", 85000},
		{"P013", "Pulpen", "Alat Tulis", 5000},
		{"P014", "Buku Tulis A5", "Alat Tulis", 8000},
		{"P015", "Penghapus", "Alat Tulis", 3000},
		{"P016", "Pensil 2B", "Alat Tulis", 4000},
		{"P017", "Sabun Cuci", "Kebutuhan RT", 15000},
		{"P018", "Tisu Wajah", "Kebutuhan RT", 12000},
		{"P019", "Sikat Gigi", "Kebutuhan RT", 10000},
		{"P020", "Pasta Gigi", "Kebutuhan RT", 18000},
	}

	cashiers := []string{"Andi", "Budi", "Citra", "Dewi"}

	// Metode bayar: 50 Cash, 30 Debit, 20 Credit (dari 100 transaksi)
	// Disusun dalam pola agar distribusi merata
	paymentPattern := []string{
		"Cash", "Cash", "Debit", "Cash", "Credit",
		"Cash", "Debit", "Cash", "Credit", "Debit",
	}

	discounts := []float64{0, 5, 10, 15}

	// Quantity pattern: variasi 1-10
	qtyPattern := []int{2, 3, 1, 5, 2, 4, 1, 3, 2, 1, 6, 2, 8, 3, 1, 10, 2, 3, 1, 2}

	// Tanggal: 19 hari (01-Nov s/d 19-Nov-2026), distribusi ~5-6 trx/hari
	dates := []string{
		"2026-11-01", "2026-11-02", "2026-11-03", "2026-11-04", "2026-11-05",
		"2026-11-06", "2026-11-07", "2026-11-08", "2026-11-09", "2026-11-10",
		"2026-11-11", "2026-11-12", "2026-11-13", "2026-11-14", "2026-11-15",
		"2026-11-16", "2026-11-17", "2026-11-18", "2026-11-19",
	}

	tx, err := DB.Begin()
	if err != nil {
		log.Fatalf("[FATAL] Gagal begin transaction: %v", err)
	}

	stmt, err := tx.Prepare(`INSERT OR IGNORE INTO transactions 
		(transaction_id, product_id, product_name, quantity, unit_price,
		 discount_percent, total_price, payment_method, transaction_date, cashier_name)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Fatalf("[FATAL] Gagal prepare statement transactions: %v", err)
	}
	defer stmt.Close()

	for i := 0; i < 100; i++ {
		transID := fmt.Sprintf("T%03d", i+1)

		// Rotasi produk (20 produk, loop 5x untuk 100 transaksi)
		prod := prodList[i%len(prodList)]

		// Rotasi kasir
		cashier := cashiers[i%len(cashiers)]

		// Pola metode bayar (10 items, loop 10x)
		payment := paymentPattern[i%len(paymentPattern)]

		// Pola diskon
		discount := discounts[i%len(discounts)]

		// Pola quantity
		qty := qtyPattern[i%len(qtyPattern)]

		// Tanggal: distribusikan ~5-6 per hari ke 19 hari
		dateIdx := i * len(dates) / 100
		if dateIdx >= len(dates) {
			dateIdx = len(dates) - 1
		}
		date := dates[dateIdx]

		// Hitung total_price
		subtotal := float64(qty) * prod.price
		discountAmount := subtotal * discount / 100.0
		totalPrice := math.Round(subtotal - discountAmount)

		_, err := stmt.Exec(transID, prod.id, prod.name, qty, prod.price,
			discount, totalPrice, payment, date, cashier)
		if err != nil {
			log.Printf("[WARNING] Gagal insert transaksi %s: %v", transID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("[FATAL] Gagal commit transaksi seeding: %v", err)
	}
	log.Println("[INFO] Seeding 100 transaksi berhasil")
}
