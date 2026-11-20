// ==================================================================
// handlers/process_handler.go — Handler Eksekusi C Processor
// Sistem Informasi Penjualan (Point of Sales)
// ==================================================================
//
// Endpoint: POST /api/process
//
// Alur:
//   1. Export data transaksi dari SQLite ke data/transactions.csv
//   2. Jalankan binary C processor (c/processor atau c/processor.exe)
//   3. Tangani exit code C: 0=sukses, 1=file not found,
//      2=invalid format, 3=memory failure (Bab 7.2 RPD)
//   4. Kembalikan status ke client
//
// Path Resolution:
//   - Working directory di-set ke rootDir proyek (pos-system/)
//   - Deteksi .exe otomatis jika runtime.GOOS == "windows"
// ==================================================================

package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"

	"pos-system/database"
	"pos-system/utils"
)

// rootDir menyimpan path absolut ke root directory proyek.
// Di-set oleh SetRootDir() dari main.go.
var rootDir string

// SetRootDir menyimpan root directory proyek untuk resolusi path.
func SetRootDir(dir string) {
	rootDir = dir
}

// getBinaryPath mengembalikan path binary dengan ekstensi .exe
// jika dijalankan di Windows.
func getBinaryPath(basePath string) string {
	if runtime.GOOS == "windows" {
		return basePath + ".exe"
	}
	return basePath
}

// RunCProcessor menangani request untuk menjalankan program C processor.
func RunCProcessor(c *gin.Context) {
	// Langkah 1: Export transaksi ke CSV
	csvPath := filepath.Join(rootDir, "data", "transactions.csv")
	recordCount, err := utils.ExportTransactionsToCSV(database.DB, csvPath)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError,
			"Gagal export CSV: "+err.Error())
		return
	}

	if recordCount == 0 {
		utils.RespondError(c, http.StatusBadRequest,
			"Tidak ada data transaksi untuk diproses. Tambahkan transaksi terlebih dahulu.")
		return
	}

	// Langkah 2: Jalankan binary C processor
	processorPath := filepath.Join(rootDir, "c", getBinaryPath("processor"))
	cmd := exec.Command(processorPath)
	cmd.Dir = rootDir // Set working directory ke root proyek

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		// Tangani exit code sesuai Bab 7.2 RPD
		if exitErr, ok := err.(*exec.ExitError); ok {
			code := exitErr.ExitCode()
			var errMsg string
			switch code {
			case 1:
				errMsg = "C Engine Error (code 1): File input tidak ditemukan"
			case 2:
				errMsg = "C Engine Error (code 2): Format data tidak valid"
			case 3:
				errMsg = "C Engine Error (code 3): Gagal alokasi memori"
			default:
				errMsg = fmt.Sprintf("C Engine Error (code %d): %s", code, stderr.String())
			}
			utils.RespondError(c, http.StatusInternalServerError, errMsg)
			return
		}

		utils.RespondError(c, http.StatusInternalServerError,
			"Gagal menjalankan C processor: "+err.Error())
		return
	}

	// Langkah 3: Response sukses
	utils.RespondSuccess(c, http.StatusOK, "Data berhasil diproses oleh C Engine", map[string]interface{}{
		"records_processed": recordCount,
		"output_file":       "processed_data.csv",
		"c_output":          stdout.String(),
	})
}
