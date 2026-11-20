# ============================================================
# Makefile — Modul Bahasa C (processor)
# Sistem Informasi Penjualan (Point of Sales)
# ============================================================

CC = gcc
CFLAGS = -Wall -Wextra -std=c99
LDFLAGS = -lm

# Deteksi OS untuk ekstensi binary
ifeq ($(OS),Windows_NT)
    EXT = .exe
    RM = del /Q
else
    EXT =
    RM = rm -f
endif

TARGET = processor$(EXT)
SOURCES = processor.c
HEADERS = processor.h

.PHONY: all clean

all: $(TARGET)

$(TARGET): $(SOURCES) $(HEADERS)
	$(CC) $(CFLAGS) -o $(TARGET) $(SOURCES) $(LDFLAGS)
	@echo "[OK] Kompilasi C processor berhasil: $(TARGET)"

clean:
	$(RM) $(TARGET)
	@echo "[OK] Clean berhasil"
