# ============================================================
# Makefile Utama — Sistem Informasi Penjualan (Point of Sales)
# Final Project — Algoritma dan Pemrograman 3
# ============================================================
#
# Penggunaan:
#   make        → Kompilasi semua (C, COBOL, Go)
#   make c      → Kompilasi C processor saja
#   make cobol  → Kompilasi COBOL reporter saja
#   make go     → Build Go web server saja
#   make run    → Jalankan Go web server
#   make test   → Jalankan C processor lalu COBOL reporter (tes pipeline)
#   make clean  → Bersihkan semua binary
# ============================================================

# Deteksi OS
ifeq ($(OS),Windows_NT)
    EXT = .exe
    RM = del /Q
    SEP = \\
    MKDIR = if not exist data mkdir data
else
    EXT =
    RM = rm -f
    SEP = /
    MKDIR = mkdir -p data
endif

# Target binary
C_TARGET = c$(SEP)processor$(EXT)
COBOL_TARGET = cobol$(SEP)reporter$(EXT)
GO_TARGET = go$(SEP)pos-system$(EXT)

.PHONY: all c cobol go run test clean

# ============================================================
# Target utama: compile semua
# ============================================================
all: c cobol
	@echo ""
	@echo "============================================================"
	@echo "  Kompilasi selesai!"
	@echo "  C Processor  : $(C_TARGET)"
	@echo "  COBOL Reporter: $(COBOL_TARGET)"
	@echo "============================================================"

# ============================================================
# Kompilasi C processor
# ============================================================
c:
	@echo "--- Kompilasi C Processor ---"
	cd c && $(MAKE)

# ============================================================
# Kompilasi COBOL reporter
# ============================================================
cobol:
	@echo "--- Kompilasi COBOL Reporter ---"
	cd cobol && $(MAKE)

# ============================================================
# Build Go web server
# ============================================================
go:
	@echo "--- Build Go Web Server ---"
	cd go && go mod tidy && go build -o pos-system$(EXT) .

# ============================================================
# Jalankan Go web server
# ============================================================
run: all go
	@echo "--- Menjalankan POS Web Server di http://localhost:8080 ---"
	cd go && .$(SEP)pos-system$(EXT)

# ============================================================
# Test pipeline: C processor -> COBOL reporter
# ============================================================
test: c cobol
	@echo ""
	@echo "--- Menjalankan C Processor ---"
	$(MKDIR)
	.$(SEP)$(C_TARGET)
	@echo ""
	@echo "--- Menjalankan COBOL Reporter ---"
	.$(SEP)$(COBOL_TARGET)
	@echo ""
	@echo "============================================================"
	@echo "  Pipeline test selesai!"
	@echo "  Periksa: data/processed_data.csv"
	@echo "  Periksa: data/report_output.txt"
	@echo "============================================================"

# ============================================================
# Bersihkan semua binary
# ============================================================
clean:
	@echo "--- Membersihkan binary ---"
	cd c && $(MAKE) clean
	cd cobol && $(MAKE) clean
	-$(RM) $(GO_TARGET)
	@echo "[OK] Semua binary berhasil dibersihkan"
