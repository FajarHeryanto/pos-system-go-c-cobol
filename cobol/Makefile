# ============================================================
# Makefile — Modul COBOL (reporter)
# Sistem Informasi Penjualan (Point of Sales)
# ============================================================

COBC = cobc
COBCFLAGS = -x

# Deteksi OS untuk ekstensi binary
ifeq ($(OS),Windows_NT)
    EXT = .exe
    RM = del /Q
else
    EXT =
    RM = rm -f
endif

TARGET = reporter$(EXT)
SOURCES = reporter.cob

.PHONY: all clean

all: $(TARGET)

$(TARGET): $(SOURCES)
	$(COBC) $(COBCFLAGS) -o $(TARGET) $(SOURCES)
	@echo "[OK] Kompilasi COBOL reporter berhasil: $(TARGET)"

clean:
	$(RM) $(TARGET)
	@echo "[OK] Clean berhasil"
