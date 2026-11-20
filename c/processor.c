/*
 * ============================================================
 * processor.c — Modul Pengolahan Data Transaksi (Bahasa C)
 * Sistem Informasi Penjualan (Point of Sales)
 * Final Project — Algoritma dan Pemrograman 3
 * ============================================================
 *
 * Fungsi Utama:
 *   1. Membaca data transaksi dari data/transactions.csv
 *   2. Menggunakan alokasi memori dinamis (malloc, calloc, realloc, free)
 *   3. Menghitung subtotal, diskon, dan net_total
 *   4. Melakukan sorting berdasarkan tanggal
 *   5. Menulis hasil ke data/processed_data.csv
 *
 * Exit Code:
 *   0 = Sukses
 *   1 = File input tidak ditemukan
 *   2 = Format data tidak valid
 *   3 = Gagal alokasi memori
 *
 * Kompilasi:
 *   gcc -o processor processor.c -lm
 * ============================================================
 */

#include "processor.h"

/* ============================================================
 * IMPLEMENTASI FUNGSI UTILITAS
 * ============================================================ */

/*
 * trim_line()
 * Menghapus karakter \r dan \n di akhir string.
 * Mencegah masalah parsing CSV akibat perbedaan line ending
 * antara Windows (CRLF) dan Unix (LF).
 */
void trim_line(char *line) {
    int len;
    if (line == NULL) {
        return;
    }
    len = (int)strlen(line);
    while (len > 0 && (line[len - 1] == '\n' || line[len - 1] == '\r')) {
        line[len - 1] = '\0';
        len--;
    }
}

/*
 * get_csv_field()
 * Mengekstrak field ke-n (0-based) dari baris CSV.
 * Delimiter: koma (,).
 * Menggunakan strncpy untuk mencegah buffer overflow.
 *
 * Parameter:
 *   line        - baris CSV yang akan diparsing
 *   field_index - indeks field yang diinginkan (0-based)
 *   output      - buffer untuk menyimpan hasil
 *   max_len     - ukuran maksimum buffer output
 *
 * Return:
 *   0 = sukses, -1 = field tidak ditemukan
 */
int get_csv_field(const char *line, int field_index, char *output, int max_len) {
    int current_field = 0;
    int start = 0;
    int i = 0;
    int field_len;

    if (line == NULL || output == NULL || max_len <= 0) {
        return -1;
    }

    while (line[i] != '\0') {
        if (line[i] == ',') {
            if (current_field == field_index) {
                field_len = i - start;
                if (field_len >= max_len) {
                    field_len = max_len - 1;
                }
                strncpy(output, line + start, field_len);
                output[field_len] = '\0';
                return 0;
            }
            current_field++;
            start = i + 1;
        }
        i++;
    }

    /* Jika field terakhir (tidak diakhiri koma) */
    if (current_field == field_index) {
        field_len = i - start;
        if (field_len >= max_len) {
            field_len = max_len - 1;
        }
        strncpy(output, line + start, field_len);
        output[field_len] = '\0';
        return 0;
    }

    return -1; /* field_index di luar jangkauan */
}

/*
 * count_data_lines()
 * Menghitung jumlah baris data di file CSV, tidak termasuk header.
 * Menggunakan fgets() untuk mencegah buffer overflow.
 *
 * Return:
 *   Jumlah baris data (>= 0) jika sukses
 *   -1 jika file tidak dapat dibuka
 */
int count_data_lines(const char *filename) {
    FILE *fp;
    char buffer[MAX_LINE_LENGTH];
    int count = 0;

    fp = fopen(filename, "r");
    if (fp == NULL) {
        return -1;
    }

    /* Lewati baris header */
    if (fgets(buffer, MAX_LINE_LENGTH, fp) == NULL) {
        fclose(fp);
        return 0;
    }

    /* Hitung baris data */
    while (fgets(buffer, MAX_LINE_LENGTH, fp) != NULL) {
        trim_line(buffer);
        /* Abaikan baris kosong */
        if (strlen(buffer) > 0) {
            count++;
        }
    }

    fclose(fp);
    return count;
}

/* ============================================================
 * IMPLEMENTASI FUNGSI ALOKASI MEMORI
 * ============================================================ */

/*
 * read_transactions()
 * Membaca seluruh data transaksi dari file CSV.
 *
 * Proses:
 *   1. Buka file dan hitung jumlah baris data
 *   2. Alokasikan memori dengan malloc() untuk array Transaction
 *   3. Baca ulang file, skip header, parse setiap baris
 *   4. Validasi jumlah field per baris (harus 10 field)
 *
 * Memory Safety:
 *   - Validasi NULL setelah malloc()
 *   - Gunakan fgets() (bukan gets()) untuk mencegah buffer overflow
 *   - Gunakan strncpy() (bukan strcpy()) untuk mencegah buffer overflow
 *
 * Parameter:
 *   filename   - path file input CSV
 *   count      - pointer untuk menyimpan jumlah record yang terbaca
 *   error_code - pointer untuk menyimpan kode error jika gagal
 *
 * Return:
 *   Pointer ke array Transaction, atau NULL jika gagal
 */
Transaction* read_transactions(const char *filename, int *count, int *error_code) {
    FILE *fp;
    char line[MAX_LINE_LENGTH];
    char field_buf[MAX_FIELD_LENGTH];
    int num_lines;
    int i;
    Transaction *transactions;

    *count = 0;
    *error_code = EXIT_SUCCESS_CODE;

    /* Hitung jumlah baris data terlebih dahulu */
    num_lines = count_data_lines(filename);
    if (num_lines < 0) {
        fprintf(stderr, "[ERROR] File input tidak ditemukan: %s\n", filename);
        *error_code = EXIT_FILE_NOT_FOUND;
        return NULL;
    }
    if (num_lines == 0) {
        fprintf(stderr, "[ERROR] File input kosong (tidak ada data): %s\n", filename);
        *error_code = EXIT_INVALID_FORMAT;
        return NULL;
    }

    /* Alokasi memori dengan malloc() */
    transactions = (Transaction *)malloc(sizeof(Transaction) * num_lines);
    if (transactions == NULL) {
        fprintf(stderr, "[ERROR] Gagal alokasi memori untuk %d transaksi\n", num_lines);
        *error_code = EXIT_MEMORY_FAILURE;
        return NULL;
    }
    fprintf(stdout, "[INFO] malloc(): Berhasil mengalokasikan memori untuk %d transaksi "
                    "(%lu bytes)\n", num_lines,
                    (unsigned long)(sizeof(Transaction) * num_lines));

    /* Buka file untuk membaca data */
    fp = fopen(filename, "r");
    if (fp == NULL) {
        fprintf(stderr, "[ERROR] Gagal membuka file: %s\n", filename);
        free(transactions);
        transactions = NULL;
        *error_code = EXIT_FILE_NOT_FOUND;
        return NULL;
    }

    /* Lewati baris header */
    if (fgets(line, MAX_LINE_LENGTH, fp) == NULL) {
        fprintf(stderr, "[ERROR] File kosong, tidak ada header\n");
        fclose(fp);
        free(transactions);
        transactions = NULL;
        *error_code = EXIT_INVALID_FORMAT;
        return NULL;
    }

    /* Baca setiap baris data */
    i = 0;
    while (fgets(line, MAX_LINE_LENGTH, fp) != NULL && i < num_lines) {
        trim_line(line);

        /* Abaikan baris kosong */
        if (strlen(line) == 0) {
            continue;
        }

        /* Parse field 0: transaction_id */
        if (get_csv_field(line, 0, field_buf, MAX_FIELD_LENGTH) != 0) {
            fprintf(stderr, "[ERROR] Format tidak valid pada baris %d: field transaction_id\n", i + 2);
            fclose(fp);
            free(transactions);
            transactions = NULL;
            *error_code = EXIT_INVALID_FORMAT;
            return NULL;
        }
        strncpy(transactions[i].transaction_id, field_buf, sizeof(transactions[i].transaction_id) - 1);
        transactions[i].transaction_id[sizeof(transactions[i].transaction_id) - 1] = '\0';

        /* Parse field 1: product_id */
        if (get_csv_field(line, 1, field_buf, MAX_FIELD_LENGTH) != 0) {
            fprintf(stderr, "[ERROR] Format tidak valid pada baris %d: field product_id\n", i + 2);
            fclose(fp);
            free(transactions);
            transactions = NULL;
            *error_code = EXIT_INVALID_FORMAT;
            return NULL;
        }
        strncpy(transactions[i].product_id, field_buf, sizeof(transactions[i].product_id) - 1);
        transactions[i].product_id[sizeof(transactions[i].product_id) - 1] = '\0';

        /* Parse field 2: product_name */
        if (get_csv_field(line, 2, field_buf, MAX_FIELD_LENGTH) != 0) {
            fprintf(stderr, "[ERROR] Format tidak valid pada baris %d: field product_name\n", i + 2);
            fclose(fp);
            free(transactions);
            transactions = NULL;
            *error_code = EXIT_INVALID_FORMAT;
            return NULL;
        }
        strncpy(transactions[i].product_name, field_buf, sizeof(transactions[i].product_name) - 1);
        transactions[i].product_name[sizeof(transactions[i].product_name) - 1] = '\0';

        /* Parse field 3: category */
        if (get_csv_field(line, 3, field_buf, MAX_FIELD_LENGTH) != 0) {
            fprintf(stderr, "[ERROR] Format tidak valid pada baris %d: field category\n", i + 2);
            fclose(fp);
            free(transactions);
            transactions = NULL;
            *error_code = EXIT_INVALID_FORMAT;
            return NULL;
        }
        strncpy(transactions[i].category, field_buf, sizeof(transactions[i].category) - 1);
        transactions[i].category[sizeof(transactions[i].category) - 1] = '\0';

        /* Parse field 4: quantity */
        if (get_csv_field(line, 4, field_buf, MAX_FIELD_LENGTH) != 0) {
            fprintf(stderr, "[ERROR] Format tidak valid pada baris %d: field quantity\n", i + 2);
            fclose(fp);
            free(transactions);
            transactions = NULL;
            *error_code = EXIT_INVALID_FORMAT;
            return NULL;
        }
        transactions[i].quantity = atoi(field_buf);

        /* Parse field 5: unit_price */
        if (get_csv_field(line, 5, field_buf, MAX_FIELD_LENGTH) != 0) {
            fprintf(stderr, "[ERROR] Format tidak valid pada baris %d: field unit_price\n", i + 2);
            fclose(fp);
            free(transactions);
            transactions = NULL;
            *error_code = EXIT_INVALID_FORMAT;
            return NULL;
        }
        transactions[i].unit_price = atof(field_buf);

        /* Parse field 6: discount_percent */
        if (get_csv_field(line, 6, field_buf, MAX_FIELD_LENGTH) != 0) {
            fprintf(stderr, "[ERROR] Format tidak valid pada baris %d: field discount_percent\n", i + 2);
            fclose(fp);
            free(transactions);
            transactions = NULL;
            *error_code = EXIT_INVALID_FORMAT;
            return NULL;
        }
        transactions[i].discount_percent = atof(field_buf);

        /* Parse field 7: payment_method */
        if (get_csv_field(line, 7, field_buf, MAX_FIELD_LENGTH) != 0) {
            fprintf(stderr, "[ERROR] Format tidak valid pada baris %d: field payment_method\n", i + 2);
            fclose(fp);
            free(transactions);
            transactions = NULL;
            *error_code = EXIT_INVALID_FORMAT;
            return NULL;
        }
        strncpy(transactions[i].payment_method, field_buf, sizeof(transactions[i].payment_method) - 1);
        transactions[i].payment_method[sizeof(transactions[i].payment_method) - 1] = '\0';

        /* Parse field 8: transaction_date */
        if (get_csv_field(line, 8, field_buf, MAX_FIELD_LENGTH) != 0) {
            fprintf(stderr, "[ERROR] Format tidak valid pada baris %d: field transaction_date\n", i + 2);
            fclose(fp);
            free(transactions);
            transactions = NULL;
            *error_code = EXIT_INVALID_FORMAT;
            return NULL;
        }
        strncpy(transactions[i].transaction_date, field_buf, sizeof(transactions[i].transaction_date) - 1);
        transactions[i].transaction_date[sizeof(transactions[i].transaction_date) - 1] = '\0';

        /* Parse field 9: cashier_name */
        if (get_csv_field(line, 9, field_buf, MAX_FIELD_LENGTH) != 0) {
            fprintf(stderr, "[ERROR] Format tidak valid pada baris %d: field cashier_name\n", i + 2);
            fclose(fp);
            free(transactions);
            transactions = NULL;
            *error_code = EXIT_INVALID_FORMAT;
            return NULL;
        }
        strncpy(transactions[i].cashier_name, field_buf, sizeof(transactions[i].cashier_name) - 1);
        transactions[i].cashier_name[sizeof(transactions[i].cashier_name) - 1] = '\0';

        i++;
    }

    fclose(fp);
    *count = i;
    fprintf(stdout, "[INFO] Berhasil membaca %d record dari %s\n", i, filename);
    return transactions;
}

/*
 * allocate_processed_data()
 * Mengalokasikan array ProcessedData menggunakan calloc().
 * calloc() menginisialisasi semua byte ke 0 (lebih aman daripada malloc).
 *
 * Parameter:
 *   count - jumlah elemen yang akan dialokasikan
 *
 * Return:
 *   Pointer ke array ProcessedData, atau NULL jika gagal
 */
ProcessedData* allocate_processed_data(int count) {
    ProcessedData *processed;

    processed = (ProcessedData *)calloc(count, sizeof(ProcessedData));
    if (processed == NULL) {
        fprintf(stderr, "[ERROR] calloc(): Gagal mengalokasikan memori untuk %d ProcessedData\n", count);
        return NULL;
    }
    fprintf(stdout, "[INFO] calloc(): Berhasil mengalokasikan memori untuk %d ProcessedData "
                    "(%lu bytes)\n", count,
                    (unsigned long)(sizeof(ProcessedData) * count));
    return processed;
}

/*
 * resize_transactions()
 * Memperbesar ukuran array Transaction menggunakan realloc().
 * Mendemonstrasikan penggunaan realloc() sesuai kurikulum.
 *
 * Memory Safety:
 *   - Simpan pointer lama sebelum realloc untuk mencegah memory leak
 *     jika realloc gagal.
 *
 * Parameter:
 *   arr      - pointer ke array Transaction yang sudah ada
 *   new_size - ukuran baru (jumlah elemen)
 *
 * Return:
 *   Pointer baru ke array Transaction, atau NULL jika gagal
 *   (pointer lama tetap valid jika gagal)
 */
Transaction* resize_transactions(Transaction *arr, int new_size) {
    Transaction *new_arr;

    new_arr = (Transaction *)realloc(arr, sizeof(Transaction) * new_size);
    if (new_arr == NULL) {
        fprintf(stderr, "[ERROR] realloc(): Gagal memperbesar array ke %d elemen\n", new_size);
        /* arr masih valid, pemanggil yang bertanggung jawab membebaskannya */
        return NULL;
    }
    fprintf(stdout, "[INFO] realloc(): Berhasil memperbesar array ke %d elemen "
                    "(%lu bytes)\n", new_size,
                    (unsigned long)(sizeof(Transaction) * new_size));
    return new_arr;
}

/* ============================================================
 * IMPLEMENTASI FUNGSI PENGOLAHAN DATA
 * ============================================================ */

/*
 * calculate_totals()
 * Menghitung subtotal, discount_amount, dan net_total untuk setiap transaksi.
 *
 * Formula:
 *   subtotal        = quantity * unit_price
 *   discount_amount = subtotal * discount_percent / 100
 *   net_total       = subtotal - discount_amount
 *
 * Menyalin field-field yang diperlukan dari Transaction ke ProcessedData.
 */
void calculate_totals(Transaction *trans, ProcessedData *processed, int count) {
    int i;
    double subtotal;
    double discount_amount;
    double net_total;

    for (i = 0; i < count; i++) {
        /* Salin field-field yang diperlukan (menggunakan strncpy untuk keamanan) */
        strncpy(processed[i].transaction_id, trans[i].transaction_id,
                sizeof(processed[i].transaction_id) - 1);
        processed[i].transaction_id[sizeof(processed[i].transaction_id) - 1] = '\0';

        strncpy(processed[i].product_name, trans[i].product_name,
                sizeof(processed[i].product_name) - 1);
        processed[i].product_name[sizeof(processed[i].product_name) - 1] = '\0';

        strncpy(processed[i].category, trans[i].category,
                sizeof(processed[i].category) - 1);
        processed[i].category[sizeof(processed[i].category) - 1] = '\0';

        processed[i].quantity = trans[i].quantity;

        strncpy(processed[i].payment_method, trans[i].payment_method,
                sizeof(processed[i].payment_method) - 1);
        processed[i].payment_method[sizeof(processed[i].payment_method) - 1] = '\0';

        strncpy(processed[i].transaction_date, trans[i].transaction_date,
                sizeof(processed[i].transaction_date) - 1);
        processed[i].transaction_date[sizeof(processed[i].transaction_date) - 1] = '\0';

        /* Hitung nilai-nilai finansial */
        subtotal = (double)trans[i].quantity * trans[i].unit_price;
        discount_amount = subtotal * trans[i].discount_percent / 100.0;
        net_total = subtotal - discount_amount;

        processed[i].subtotal = subtotal;
        processed[i].discount_amount = discount_amount;
        processed[i].net_total = net_total;
    }

    fprintf(stdout, "[INFO] Kalkulasi subtotal, diskon, dan net_total selesai untuk %d record\n", count);
}

/*
 * compare_by_date()
 * Fungsi perbandingan untuk qsort().
 * Membandingkan dua ProcessedData berdasarkan transaction_date.
 * Karena format tanggal adalah YYYY-MM-DD, perbandingan string
 * secara leksikografis sudah memberikan urutan kronologis yang benar.
 */
static int compare_by_date(const void *a, const void *b) {
    const ProcessedData *pa = (const ProcessedData *)a;
    const ProcessedData *pb = (const ProcessedData *)b;
    return strcmp(pa->transaction_date, pb->transaction_date);
}

/*
 * sort_by_date()
 * Mengurutkan array ProcessedData berdasarkan tanggal (ascending).
 * Menggunakan qsort() dari stdlib.h.
 */
void sort_by_date(ProcessedData *data, int count) {
    if (data == NULL || count <= 0) {
        return;
    }
    qsort(data, count, sizeof(ProcessedData), compare_by_date);
    fprintf(stdout, "[INFO] Data berhasil diurutkan berdasarkan tanggal (ascending)\n");
}

/* ============================================================
 * IMPLEMENTASI FUNGSI OUTPUT
 * ============================================================ */

/*
 * write_processed_csv()
 * Menulis array ProcessedData ke file CSV.
 *
 * Format output (sesuai RPD Bab 4.3):
 *   transaction_id,product_name,category,quantity,subtotal,
 *   discount_amount,net_total,payment_method,transaction_date
 *
 * Nilai numerik menggunakan %.0f (tanpa desimal) agar COBOL
 * dapat mem-parsing dengan benar menggunakan FUNCTION NUMVAL.
 */
void write_processed_csv(const char *filename, ProcessedData *data, int count) {
    FILE *fp;
    int i;

    fp = fopen(filename, "w");
    if (fp == NULL) {
        fprintf(stderr, "[ERROR] Gagal membuka file output: %s\n", filename);
        return;
    }

    /* Tulis header CSV */
    fprintf(fp, "transaction_id,product_name,category,quantity,subtotal,"
                "discount_amount,net_total,payment_method,transaction_date\n");

    /* Tulis setiap baris data */
    for (i = 0; i < count; i++) {
        fprintf(fp, "%s,%s,%s,%d,%.0f,%.0f,%.0f,%s,%s\n",
                data[i].transaction_id,
                data[i].product_name,
                data[i].category,
                data[i].quantity,
                data[i].subtotal,
                data[i].discount_amount,
                data[i].net_total,
                data[i].payment_method,
                data[i].transaction_date);
    }

    fclose(fp);
    fprintf(stdout, "[INFO] Berhasil menulis %d record ke %s\n", count, filename);
}

/* ============================================================
 * IMPLEMENTASI FUNGSI PEMBEBASAN MEMORI
 * ============================================================ */

/*
 * free_all_memory()
 * Membebaskan semua memori yang telah dialokasikan secara dinamis.
 *
 * Memory Safety:
 *   - Cek NULL sebelum free (mencegah double free)
 *   - Set pointer ke NULL setelah free (mencegah dangling pointer)
 *   - Menggunakan pointer-to-pointer agar bisa set ke NULL dari dalam fungsi
 */
void free_all_memory(Transaction **trans, ProcessedData **processed) {
    if (trans != NULL && *trans != NULL) {
        free(*trans);
        *trans = NULL;
        fprintf(stdout, "[INFO] free(): Memori Transaction berhasil dibebaskan\n");
    }

    if (processed != NULL && *processed != NULL) {
        free(*processed);
        *processed = NULL;
        fprintf(stdout, "[INFO] free(): Memori ProcessedData berhasil dibebaskan\n");
    }
}

/* ============================================================
 * FUNGSI MAIN — PROGRAM UTAMA
 * ============================================================ */

/*
 * main()
 * Alur eksekusi utama program C processor:
 *
 *   1. Baca data dari transactions.csv (malloc)
 *   2. Alokasikan array ProcessedData (calloc)
 *   3. Demonstrasi realloc (perbesar kapasitas)
 *   4. Hitung subtotal, diskon, net_total
 *   5. Urutkan berdasarkan tanggal
 *   6. Tulis hasil ke processed_data.csv
 *   7. Bebaskan semua memori (free)
 *
 * Exit Code:
 *   0 = Sukses (semua proses berhasil)
 *   1 = File input tidak ditemukan
 *   2 = Format data CSV tidak valid
 *   3 = Gagal alokasi memori
 */
int main(void) {
    Transaction *transactions = NULL;
    ProcessedData *processed = NULL;
    Transaction *resized = NULL;
    int count = 0;
    int error_code = EXIT_SUCCESS_CODE;

    fprintf(stdout, "============================================================\n");
    fprintf(stdout, "  POS System — C Data Processor\n");
    fprintf(stdout, "  Membaca : %s\n", INPUT_FILE);
    fprintf(stdout, "  Menulis : %s\n", OUTPUT_FILE);
    fprintf(stdout, "============================================================\n\n");

    /* ---- LANGKAH 1: Baca data transaksi (malloc) ---- */
    fprintf(stdout, "--- Langkah 1: Membaca data transaksi ---\n");
    transactions = read_transactions(INPUT_FILE, &count, &error_code);
    if (transactions == NULL) {
        /* error_code sudah di-set oleh read_transactions() */
        return error_code;
    }

    /* ---- LANGKAH 2: Alokasikan ProcessedData (calloc) ---- */
    fprintf(stdout, "\n--- Langkah 2: Mengalokasikan ProcessedData ---\n");
    processed = allocate_processed_data(count);
    if (processed == NULL) {
        free(transactions);
        transactions = NULL;
        return EXIT_MEMORY_FAILURE;
    }

    /* ---- LANGKAH 3: Demonstrasi realloc ---- */
    fprintf(stdout, "\n--- Langkah 3: Demonstrasi realloc ---\n");
    /*
     * Demonstrasi realloc: perbesar kapasitas array Transaction
     * menjadi 2x lipat (simulasi kasus data bertambah).
     * Setelah demonstrasi, kembalikan ke ukuran asli.
     */
    resized = resize_transactions(transactions, count * 2);
    if (resized != NULL) {
        transactions = resized;
        fprintf(stdout, "[INFO] Demonstrasi realloc berhasil (kapasitas: %d -> %d)\n",
                count, count * 2);
        /* Kembalikan ke ukuran asli untuk efisiensi memori */
        resized = resize_transactions(transactions, count);
        if (resized != NULL) {
            transactions = resized;
            fprintf(stdout, "[INFO] Kapasitas dikembalikan ke %d elemen\n", count);
        }
    } else {
        fprintf(stderr, "[WARNING] realloc gagal, melanjutkan dengan array asli\n");
    }

    /* ---- LANGKAH 4: Hitung subtotal, diskon, net_total ---- */
    fprintf(stdout, "\n--- Langkah 4: Menghitung subtotal, diskon, net_total ---\n");
    calculate_totals(transactions, processed, count);

    /* ---- LANGKAH 5: Sorting berdasarkan tanggal ---- */
    fprintf(stdout, "\n--- Langkah 5: Sorting berdasarkan tanggal ---\n");
    sort_by_date(processed, count);

    /* ---- LANGKAH 6: Tulis hasil ke processed_data.csv ---- */
    fprintf(stdout, "\n--- Langkah 6: Menulis hasil ke file output ---\n");
    write_processed_csv(OUTPUT_FILE, processed, count);

    /* ---- LANGKAH 7: Bebaskan semua memori ---- */
    fprintf(stdout, "\n--- Langkah 7: Membebaskan memori ---\n");
    free_all_memory(&transactions, &processed);

    fprintf(stdout, "\n============================================================\n");
    fprintf(stdout, "  Proses selesai dengan sukses (exit code: 0)\n");
    fprintf(stdout, "============================================================\n");

    return EXIT_SUCCESS_CODE;
}
