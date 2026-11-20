/*
 * ============================================================
 * processor.h — Header File untuk Modul Pengolahan Data C
 * Sistem Informasi Penjualan (Point of Sales)
 * Final Project — Algoritma dan Pemrograman 3
 * ============================================================
 *
 * Tanggung Jawab:
 *   - Mendefinisikan struktur data untuk transaksi dan hasil olahan
 *   - Mendeklarasikan fungsi-fungsi utama pengolahan data
 *   - Menetapkan konstanta exit code dan ukuran buffer
 *
 * Exit Code:
 *   0 = Sukses
 *   1 = File input tidak ditemukan
 *   2 = Format data tidak valid
 *   3 = Gagal alokasi memori
 * ============================================================
 */

#ifndef PROCESSOR_H
#define PROCESSOR_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>

/* ============================================================
 * KONSTANTA
 * ============================================================ */

#define EXIT_SUCCESS_CODE     0
#define EXIT_FILE_NOT_FOUND   1
#define EXIT_INVALID_FORMAT   2
#define EXIT_MEMORY_FAILURE   3

#define MAX_LINE_LENGTH       512
#define MAX_FIELD_LENGTH      100

#define INPUT_FILE            "data/transactions.csv"
#define OUTPUT_FILE           "data/processed_data.csv"

/* ============================================================
 * STRUKTUR DATA
 * ============================================================ */

/* Struct untuk menyimpan satu transaksi (dari transactions.csv) */
typedef struct {
    char   transaction_id[10];
    char   product_id[10];
    char   product_name[50];
    char   category[30];
    int    quantity;
    double unit_price;
    double discount_percent;
    char   payment_method[10];
    char   transaction_date[12];
    char   cashier_name[30];
} Transaction;

/* Struct untuk menyimpan data yang sudah diproses (output ke processed_data.csv) */
typedef struct {
    char   transaction_id[10];
    char   product_name[50];
    char   category[30];
    int    quantity;
    double subtotal;          /* unit_price * quantity */
    double discount_amount;   /* subtotal * discount_percent / 100 */
    double net_total;         /* subtotal - discount_amount */
    char   payment_method[10];
    char   transaction_date[12];
} ProcessedData;

/* Struct untuk statistik per kategori */
typedef struct {
    char   category[30];
    int    total_transactions;
    double total_revenue;
    double average_per_transaction;
} CategoryStats;

/* ============================================================
 * DEKLARASI FUNGSI — UTILITAS
 * ============================================================ */

/* Menghapus karakter newline (\n, \r) dari akhir string */
void trim_line(char *line);

/* Mengekstrak field ke-n (0-based) dari baris CSV berdasarkan delimiter koma */
int get_csv_field(const char *line, int field_index, char *output, int max_len);

/* Menghitung jumlah baris data (minus header) di file CSV */
int count_data_lines(const char *filename);

/* ============================================================
 * DEKLARASI FUNGSI — ALOKASI MEMORI
 * ============================================================ */

/* Membaca file transactions.csv, mengalokasikan array Transaction via malloc.
 * Mengembalikan pointer ke array Transaction, mengisi *count dengan jumlah record.
 * Mengembalikan NULL jika gagal (exit code di-set via parameter error_code). */
Transaction* read_transactions(const char *filename, int *count, int *error_code);

/* Mengalokasikan array ProcessedData via calloc.
 * Mengembalikan NULL jika gagal. */
ProcessedData* allocate_processed_data(int count);

/* Memperbesar array Transaction menggunakan realloc.
 * Mengembalikan pointer baru atau NULL jika gagal. */
Transaction* resize_transactions(Transaction *arr, int new_size);

/* ============================================================
 * DEKLARASI FUNGSI — PENGOLAHAN DATA
 * ============================================================ */

/* Menghitung subtotal, discount_amount, dan net_total untuk setiap transaksi */
void calculate_totals(Transaction *trans, ProcessedData *processed, int count);

/* Mengurutkan array ProcessedData berdasarkan transaction_date (ascending, YYYY-MM-DD) */
void sort_by_date(ProcessedData *data, int count);

/* ============================================================
 * DEKLARASI FUNGSI — OUTPUT
 * ============================================================ */

/* Menulis array ProcessedData ke file CSV output */
void write_processed_csv(const char *filename, ProcessedData *data, int count);

/* ============================================================
 * DEKLARASI FUNGSI — PEMBEBASAN MEMORI
 * ============================================================ */

/* Membebaskan semua memori yang dialokasikan dan set pointer ke NULL */
void free_all_memory(Transaction **trans, ProcessedData **processed);

#endif /* PROCESSOR_H */
