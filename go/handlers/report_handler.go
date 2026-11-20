// ==================================================================
// handlers/report_handler.go — Handler Eksekusi COBOL Reporter
// Sistem Informasi Penjualan (Point of Sales)
// ==================================================================
//
// Endpoint:
//   POST /api/reports/generate  → RunCOBOLReporter
//   GET  /api/reports/income    → GetIncomeReport
//
// Alur generate:
//   1. Jalankan binary COBOL reporter (cobol/reporter atau .exe)
//   2. COBOL membaca data/processed_data.csv
//   3. COBOL menulis laporan ke data/report_output.txt
//   4. Kembalikan status ke client
//
// Alur baca laporan:
//   1. Baca file data/report_output.txt
//   2. Kembalikan isi laporan sebagai string ke client
// ==================================================================

package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"pos-system/utils"
)

// RunCOBOLReporter menangani request untuk menjalankan program COBOL reporter.
func RunCOBOLReporter(c *gin.Context) {
	// Cek apakah processed_data.csv sudah ada
	processedCSV := filepath.Join(rootDir, "data", "processed_data.csv")
	if _, err := os.Stat(processedCSV); os.IsNotExist(err) {
		utils.RespondError(c, http.StatusBadRequest,
			"File processed_data.csv belum ada. Jalankan C Processor terlebih dahulu.")
		return
	}

	// Jalankan binary COBOL reporter
	reporterPath := filepath.Join(rootDir, "cobol", getBinaryPath("reporter"))
	cmd := exec.Command(reporterPath)
	cmd.Dir = rootDir // Set working directory ke root proyek

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code := exitErr.ExitCode()
			var errMsg string
			switch code {
			case 1:
				errMsg = "COBOL Reporter Error (code 1): File input tidak ditemukan"
			default:
				errMsg = fmt.Sprintf("COBOL Reporter Error (code %d): %s", code, stderr.String())
			}
			utils.RespondError(c, http.StatusInternalServerError, errMsg)
			return
		}

		utils.RespondError(c, http.StatusInternalServerError,
			"Gagal menjalankan COBOL reporter: "+err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "Laporan berhasil digenerate oleh COBOL Reporter", map[string]interface{}{
		"output_file":  "report_output.txt",
		"cobol_output": stdout.String(),
	})
}

// GetIncomeReport membaca dan mengembalikan isi file laporan.
func GetIncomeReport(c *gin.Context) {
	reportPath := filepath.Join(rootDir, "data", "report_output.txt")

	content, err := os.ReadFile(reportPath)
	if err != nil {
		if os.IsNotExist(err) {
			utils.RespondError(c, http.StatusNotFound,
				"Laporan belum digenerate. Jalankan COBOL Reporter terlebih dahulu.")
			return
		}
		utils.RespondError(c, http.StatusInternalServerError,
			"Gagal membaca file laporan: "+err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "", map[string]string{
		"report_content": string(content),
	})
}
