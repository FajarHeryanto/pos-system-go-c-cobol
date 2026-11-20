// ==================================================================
// main.go — Entry Point Go Web Server
// Sistem Informasi Penjualan (Point of Sales)
// Final Project — Algoritma dan Pemrograman 3
// ==================================================================
//
// Fungsi:
//   - Mendeteksi root directory proyek secara otomatis
//   - Inisialisasi database SQLite (migrasi + seeding)
//   - Konfigurasi router Gin (halaman web + REST API)
//   - Menyajikan static files (HTML/CSS/JS) dari folder web/
//   - Menjalankan web server di http://localhost:8080
//
// Path Resolution:
//   Program ini bisa dijalankan dari:
//   1. Root proyek: cd pos-system && go run go/main.go
//   2. Folder go/ : cd pos-system/go && go run main.go
//   3. Binary     : ./go/pos-system (dari root) atau ./pos-system (dari go/)
//   Mekanisme deteksi: cek keberadaan folder "data/" relatif terhadap
//   working directory, lalu cek parent directory.
//
// Kompilasi:
//   cd go/ && go build -o pos-system .
//
// Eksekusi:
//   cd pos-system/ && go/pos-system (Linux/Mac)
//   cd pos-system\ && go\pos-system.exe (Windows)
//   ATAU
//   cd pos-system/go && go run main.go
// ==================================================================

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"pos-system/database"
	"pos-system/handlers"
	"pos-system/middleware"
)

// detectRootDir mendeteksi root directory proyek (pos-system/).
// Strategi:
//  1. Cek apakah folder "data" ada di working directory → root = cwd
//  2. Cek apakah folder "data" ada di parent directory → root = parent
//  3. Cek apakah folder "data" ada di directory executable → root = exec dir
//  4. Fallback: gunakan working directory
func detectRootDir() string {
	// Strategi 1: cek dari working directory
	cwd, err := os.Getwd()
	if err == nil {
		if isProjectRoot(cwd) {
			log.Println("[INFO] Root directory (cwd):", cwd)
			return cwd
		}

		// Strategi 2: cek parent directory (jika dijalankan dari go/)
		parent := filepath.Dir(cwd)
		if isProjectRoot(parent) {
			log.Println("[INFO] Root directory (parent):", parent)
			return parent
		}
	}

	// Strategi 3: cek dari lokasi executable
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		if isProjectRoot(execDir) {
			log.Println("[INFO] Root directory (exec):", execDir)
			return execDir
		}

		execParent := filepath.Dir(execDir)
		if isProjectRoot(execParent) {
			log.Println("[INFO] Root directory (exec parent):", execParent)
			return execParent
		}
	}

	// Fallback: gunakan cwd
	log.Println("[WARNING] Tidak dapat mendeteksi root directory, menggunakan cwd:", cwd)
	return cwd
}

// isProjectRoot mengecek apakah sebuah directory adalah root proyek pos-system.
// Kriteria: harus ada subfolder "data" dan minimal salah satu dari "c" atau "cobol".
func isProjectRoot(dir string) bool {
	dataDir := filepath.Join(dir, "data")
	cDir := filepath.Join(dir, "c")
	cobolDir := filepath.Join(dir, "cobol")

	if info, err := os.Stat(dataDir); err == nil && info.IsDir() {
		// Minimal harus ada folder c/ atau cobol/ juga
		if info2, err2 := os.Stat(cDir); err2 == nil && info2.IsDir() {
			return true
		}
		if info3, err3 := os.Stat(cobolDir); err3 == nil && info3.IsDir() {
			return true
		}
	}
	return false
}

func main() {
	// ================================================================
	// LANGKAH 1: Deteksi root directory proyek
	// ================================================================
	rootDir := detectRootDir()
	handlers.SetRootDir(rootDir)

	// Pastikan folder data/ ada
	dataDir := filepath.Join(rootDir, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("[FATAL] Gagal membuat folder data: %v", err)
	}

	// ================================================================
	// LANGKAH 2: Inisialisasi Database
	// ================================================================
	dbPath := filepath.Join(rootDir, "data", "pos_database.db")
	db := database.InitDB(dbPath)
	defer db.Close()

	// ================================================================
	// LANGKAH 3: Setup Gin Router
	// ================================================================
	r := gin.Default()

	// CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Serve static files (website) dari folder web/
	webDir := filepath.Join(rootDir, "web")
	r.Static("/css", filepath.Join(webDir, "css"))
	r.Static("/js", filepath.Join(webDir, "js"))
	r.Static("/assets", filepath.Join(webDir, "assets"))
	r.LoadHTMLGlob(filepath.Join(webDir, "*.html"))

	// ================================================================
	// LANGKAH 4: Halaman Web (HTML)
	// ================================================================
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	r.GET("/products", func(c *gin.Context) {
		c.HTML(200, "products.html", nil)
	})
	r.GET("/transactions", func(c *gin.Context) {
		c.HTML(200, "transactions.html", nil)
	})
	r.GET("/reports", func(c *gin.Context) {
		c.HTML(200, "reports.html", nil)
	})
	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})

	// ================================================================
	// LANGKAH 5: REST API Routes (sesuai Bab 3 RPD)
	// ================================================================
	api := r.Group("/api")
	{
		// === Public Route: Login ===
		api.POST("/auth/login", handlers.Login)

		// === Protected Routes (memerlukan JWT token) ===
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// --- Produk ---
			protected.GET("/products", handlers.GetAllProducts)
			protected.GET("/products/:id", handlers.GetProduct)
			protected.POST("/products", handlers.CreateProduct)
			protected.PUT("/products/:id", handlers.UpdateProduct)
			protected.DELETE("/products/:id", handlers.DeleteProduct)

			// --- Transaksi ---
			protected.GET("/transactions", handlers.GetAllTransactions)
			protected.POST("/transactions", handlers.CreateTransaction)
			protected.DELETE("/transactions/:id", handlers.DeleteTransaction)

			// --- Proses & Laporan ---
			protected.POST("/process", handlers.RunCProcessor)
			protected.POST("/reports/generate", handlers.RunCOBOLReporter)
			protected.GET("/reports/income", handlers.GetIncomeReport)

			// --- Dashboard ---
			protected.GET("/dashboard", handlers.GetDashboardStats)
		}
	}

	// ================================================================
	// LANGKAH 6: Jalankan Server
	// ================================================================
	log.Println("============================================================")
	log.Println("  POS System — Go Web Server")
	log.Println("  Root Directory :", rootDir)
	log.Println("  Database       :", dbPath)
	log.Println("  Web Directory  :", webDir)
	log.Println("  Server         : http://localhost:8080")
	log.Println("============================================================")
	log.Println("")
	log.Println("  Login Default:")
	log.Println("    Admin : admin / admin123")
	log.Println("    Kasir : andi / kasir123 (atau budi, citra, dewi)")
	log.Println("")

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("[FATAL] Gagal menjalankan server: %v", err)
	}
}
