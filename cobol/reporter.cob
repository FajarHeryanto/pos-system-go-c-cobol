      ******************************************************************
      * reporter.cob — Modul Pelaporan Bisnis (COBOL)
      * Sistem Informasi Penjualan (Point of Sales)
      * Final Project — Algoritma dan Pemrograman 3
      *
      * Tanggung Jawab:
      *   1. Membaca data dari data/processed_data.csv
      *   2. Merekap data per kategori, per tanggal, per produk
      *   3. Menulis 3 jenis laporan ke data/report_output.txt:
      *      A. Rekapitulasi Pendapatan Per Kategori + Grand Total
      *      B. Laporan Penjualan Harian
      *      C. Laporan Produk Terlaris Top 10
      *
      * Kompilasi: cobc -x -o reporter reporter.cob
      ******************************************************************
       IDENTIFICATION DIVISION.
       PROGRAM-ID. POS-REPORTER.
       AUTHOR. TIM-POS.
       DATE-WRITTEN. 2026-11-19.

       ENVIRONMENT DIVISION.
       INPUT-OUTPUT SECTION.
       FILE-CONTROL.
           SELECT INPUT-FILE
               ASSIGN TO "data/processed_data.csv"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-INPUT-STATUS.
           SELECT OUTPUT-FILE
               ASSIGN TO "data/report_output.txt"
               ORGANIZATION IS LINE SEQUENTIAL
               FILE STATUS IS WS-OUTPUT-STATUS.

       DATA DIVISION.
       FILE SECTION.

       FD  INPUT-FILE.
       01  INPUT-RECORD                  PIC X(300).

       FD  OUTPUT-FILE.
       01  OUTPUT-RECORD                 PIC X(120).

       WORKING-STORAGE SECTION.

      ******************************************************************
      * STATUS FLAGS
      ******************************************************************
       01  WS-INPUT-STATUS               PIC XX.
       01  WS-OUTPUT-STATUS              PIC XX.
       01  WS-EOF-FLAG                   PIC 9 VALUE 0.
           88 WS-EOF                     VALUE 1.

      ******************************************************************
      * PARSED CSV FIELDS (dari processed_data.csv: 9 kolom)
      * transaction_id, product_name, category, quantity,
      * subtotal, discount_amount, net_total, payment_method,
      * transaction_date
      ******************************************************************
       01  WS-CSV-FIELDS.
           05  WS-CSV-TRANS-ID           PIC X(10).
           05  WS-CSV-PROD-NAME          PIC X(50).
           05  WS-CSV-CATEGORY           PIC X(30).
           05  WS-CSV-QTY-STR            PIC X(20).
           05  WS-CSV-SUBTOTAL-STR       PIC X(20).
           05  WS-CSV-DISC-AMT-STR       PIC X(20).
           05  WS-CSV-NET-TOTAL-STR      PIC X(20).
           05  WS-CSV-PAY-METHOD         PIC X(10).
           05  WS-CSV-TRANS-DATE         PIC X(12).

      ******************************************************************
      * PARSED NUMERIC VALUES
      ******************************************************************
       01  WS-PARSED-QTY                 PIC 9(6)      VALUE 0.
       01  WS-PARSED-NET-TOTAL           PIC 9(12)V99  VALUE 0.

      ******************************************************************
      * TABEL KATEGORI (maks 20 kategori)
      ******************************************************************
       01  WS-CAT-TABLE.
           05  WS-CAT-COUNT              PIC 99 VALUE 0.
           05  WS-CAT-ENTRY OCCURS 20 TIMES.
               10  WS-CAT-NAME           PIC X(30).
               10  WS-CAT-TRX-COUNT      PIC 9(6) VALUE 0.
               10  WS-CAT-TOTAL-REV      PIC 9(12)V99 VALUE 0.
               10  WS-CAT-AVG-REV        PIC 9(12)V99 VALUE 0.

      ******************************************************************
      * TABEL TANGGAL (maks 31 tanggal)
      ******************************************************************
       01  WS-DATE-TABLE.
           05  WS-DATE-COUNT             PIC 99 VALUE 0.
           05  WS-DATE-ENTRY OCCURS 31 TIMES.
               10  WS-DT-DATE            PIC X(12).
               10  WS-DT-TRX-COUNT       PIC 9(6) VALUE 0.
               10  WS-DT-TOTAL-REV       PIC 9(12)V99 VALUE 0.

      ******************************************************************
      * TABEL PRODUK (maks 100 produk)
      ******************************************************************
       01  WS-PROD-TABLE.
           05  WS-PROD-COUNT             PIC 9(3) VALUE 0.
           05  WS-PROD-ENTRY OCCURS 100 TIMES.
               10  WS-PR-NAME            PIC X(50).
               10  WS-PR-QTY-SOLD        PIC 9(6) VALUE 0.
               10  WS-PR-TOTAL-REV       PIC 9(12)V99 VALUE 0.

      ******************************************************************
      * GRAND TOTAL
      ******************************************************************
       01  WS-GRAND-TRX                  PIC 9(6) VALUE 0.
       01  WS-GRAND-REVENUE              PIC 9(12)V99 VALUE 0.
       01  WS-GRAND-AVG                  PIC 9(12)V99 VALUE 0.

      ******************************************************************
      * PERIODE (min/max tanggal)
      ******************************************************************
       01  WS-MIN-DATE                   PIC X(12) VALUE "9999-99-99".
       01  WS-MAX-DATE                   PIC X(12) VALUE "0000-00-00".

      ******************************************************************
      * FORMATTING — ANGKA
      * Menggunakan PIC dengan koma sebagai separator ribuan lalu
      * INSPECT REPLACING ALL "," BY "." untuk format Indonesia.
      * PIC koma (,) = separator ribuan dalam COBOL standard.
      ******************************************************************
       01  WS-FMT-NUM                    PIC ZZZ,ZZZ,ZZ9.
       01  WS-FMT-BIG                    PIC Z,ZZZ,ZZZ,ZZ9.
       01  WS-FMT-TRX                    PIC ZZZZ9.

      ******************************************************************
      * OUTPUT LINE BUFFER
      ******************************************************************
       01  WS-OUT-LINE                   PIC X(120).
       01  WS-SEPARATOR-LINE             PIC X(63) VALUE
           "---------------------------------------------"
           "------------------".
       01  WS-DOUBLE-LINE                PIC X(63) VALUE
           "============================================="
           "==================".

      ******************************************************************
      * SORTING TEMP VARIABLES
      ******************************************************************
       01  WS-SORT-TEMP-DATE             PIC X(12).
       01  WS-SORT-TEMP-TRX              PIC 9(6).
       01  WS-SORT-TEMP-REV              PIC 9(12)V99.
       01  WS-SORT-TEMP-NAME             PIC X(50).
       01  WS-SORT-TEMP-QTY              PIC 9(6).

      ******************************************************************
      * DATE CONVERSION — YYYY-MM-DD to DD-Mon-YYYY
      ******************************************************************
       01  WS-CONV-YEAR                  PIC X(4).
       01  WS-CONV-MONTH                 PIC X(2).
       01  WS-CONV-DAY                   PIC X(2).
       01  WS-CONV-MON-NAME              PIC X(3).
       01  WS-CONV-RESULT                PIC X(12).
       01  WS-CONV-RESULT-MIN            PIC X(12).

      ******************************************************************
      * LOOP COUNTERS
      ******************************************************************
       01  WS-I                          PIC 9(3) VALUE 0.
       01  WS-J                          PIC 9(3) VALUE 0.
       01  WS-FOUND                      PIC 9 VALUE 0.
       01  WS-TOP-LIMIT                  PIC 99 VALUE 10.
       01  WS-RECORD-COUNT               PIC 9(6) VALUE 0.
       01  WS-HEADER-SKIPPED             PIC 9 VALUE 0.
       01  WS-TRIMMED-NAME               PIC X(50).
       01  WS-NUM-INDEX                  PIC 99.

      ******************************************************************
      * PROCEDURE DIVISION
      ******************************************************************
       PROCEDURE DIVISION.
       0000-MAIN.
           PERFORM 1000-INIT
           PERFORM 2000-READ-DATA
           PERFORM 3000-CALC-AVERAGES
           PERFORM 4000-SORT-DATES
           PERFORM 4100-SORT-PRODUCTS
           PERFORM 5000-WRITE-REPORT
           PERFORM 9000-CLEANUP
           MOVE 0 TO RETURN-CODE
           STOP RUN.

      ******************************************************************
      * 1000-INIT — Buka file input dan output
      ******************************************************************
       1000-INIT.
           OPEN INPUT INPUT-FILE
           IF WS-INPUT-STATUS NOT = "00"
               DISPLAY "[ERROR] Gagal membuka processed_data.csv"
                   " Status: " WS-INPUT-STATUS
               MOVE 1 TO RETURN-CODE
               STOP RUN
           END-IF

           OPEN OUTPUT OUTPUT-FILE
           IF WS-OUTPUT-STATUS NOT = "00"
               DISPLAY "[ERROR] Gagal membuka report_output.txt"
                   " Status: " WS-OUTPUT-STATUS
               CLOSE INPUT-FILE
               MOVE 1 TO RETURN-CODE
               STOP RUN
           END-IF

           DISPLAY "[INFO] File input dan output berhasil dibuka".

      ******************************************************************
      * 2000-READ-DATA — Baca semua record, akumulasi data
      ******************************************************************
       2000-READ-DATA.
           MOVE 0 TO WS-EOF-FLAG
           MOVE 0 TO WS-HEADER-SKIPPED
           PERFORM UNTIL WS-EOF
               READ INPUT-FILE INTO INPUT-RECORD
                   AT END
                       SET WS-EOF TO TRUE
                   NOT AT END
                       IF WS-HEADER-SKIPPED = 0
                           MOVE 1 TO WS-HEADER-SKIPPED
                       ELSE
                           PERFORM 2100-PARSE-RECORD
                       END-IF
               END-READ
           END-PERFORM
           DISPLAY "[INFO] Total record terbaca: " WS-RECORD-COUNT.

      ******************************************************************
      * 2100-PARSE-RECORD — Parse satu baris CSV dan akumulasi
      ******************************************************************
       2100-PARSE-RECORD.
           INITIALIZE WS-CSV-FIELDS

           UNSTRING INPUT-RECORD DELIMITED BY ","
               INTO WS-CSV-TRANS-ID
                    WS-CSV-PROD-NAME
                    WS-CSV-CATEGORY
                    WS-CSV-QTY-STR
                    WS-CSV-SUBTOTAL-STR
                    WS-CSV-DISC-AMT-STR
                    WS-CSV-NET-TOTAL-STR
                    WS-CSV-PAY-METHOD
                    WS-CSV-TRANS-DATE
           END-UNSTRING

      *    Konversi string numerik menggunakan FUNCTION NUMVAL
           COMPUTE WS-PARSED-QTY =
               FUNCTION NUMVAL(WS-CSV-QTY-STR)
           COMPUTE WS-PARSED-NET-TOTAL =
               FUNCTION NUMVAL(WS-CSV-NET-TOTAL-STR)

           ADD 1 TO WS-RECORD-COUNT

      *    Trim trailing spaces dari nama produk
           MOVE FUNCTION TRIM(WS-CSV-PROD-NAME TRAILING)
               TO WS-TRIMMED-NAME

      *    Akumulasi per kategori
           PERFORM 2200-ACCUM-CATEGORY

      *    Akumulasi per tanggal
           PERFORM 2300-ACCUM-DATE

      *    Akumulasi per produk
           PERFORM 2400-ACCUM-PRODUCT

      *    Track min/max tanggal untuk periode laporan
           IF WS-CSV-TRANS-DATE < WS-MIN-DATE
               MOVE WS-CSV-TRANS-DATE TO WS-MIN-DATE
           END-IF
           IF WS-CSV-TRANS-DATE > WS-MAX-DATE
               MOVE WS-CSV-TRANS-DATE TO WS-MAX-DATE
           END-IF.

      ******************************************************************
      * 2200-ACCUM-CATEGORY — Cari/tambah kategori, akumulasi
      ******************************************************************
       2200-ACCUM-CATEGORY.
           MOVE 0 TO WS-FOUND
           PERFORM VARYING WS-I FROM 1 BY 1
               UNTIL WS-I > WS-CAT-COUNT OR WS-FOUND = 1
               IF FUNCTION TRIM(WS-CAT-NAME(WS-I) TRAILING) =
                  FUNCTION TRIM(WS-CSV-CATEGORY TRAILING)
                   ADD 1 TO WS-CAT-TRX-COUNT(WS-I)
                   ADD WS-PARSED-NET-TOTAL
                       TO WS-CAT-TOTAL-REV(WS-I)
                   MOVE 1 TO WS-FOUND
               END-IF
           END-PERFORM

           IF WS-FOUND = 0
               ADD 1 TO WS-CAT-COUNT
               MOVE WS-CSV-CATEGORY
                   TO WS-CAT-NAME(WS-CAT-COUNT)
               MOVE 1 TO WS-CAT-TRX-COUNT(WS-CAT-COUNT)
               MOVE WS-PARSED-NET-TOTAL
                   TO WS-CAT-TOTAL-REV(WS-CAT-COUNT)
           END-IF.

      ******************************************************************
      * 2300-ACCUM-DATE — Cari/tambah tanggal, akumulasi
      ******************************************************************
       2300-ACCUM-DATE.
           MOVE 0 TO WS-FOUND
           PERFORM VARYING WS-I FROM 1 BY 1
               UNTIL WS-I > WS-DATE-COUNT OR WS-FOUND = 1
               IF FUNCTION TRIM(WS-DT-DATE(WS-I) TRAILING) =
                  FUNCTION TRIM(WS-CSV-TRANS-DATE TRAILING)
                   ADD 1 TO WS-DT-TRX-COUNT(WS-I)
                   ADD WS-PARSED-NET-TOTAL
                       TO WS-DT-TOTAL-REV(WS-I)
                   MOVE 1 TO WS-FOUND
               END-IF
           END-PERFORM

           IF WS-FOUND = 0
               ADD 1 TO WS-DATE-COUNT
               MOVE WS-CSV-TRANS-DATE
                   TO WS-DT-DATE(WS-DATE-COUNT)
               MOVE 1 TO WS-DT-TRX-COUNT(WS-DATE-COUNT)
               MOVE WS-PARSED-NET-TOTAL
                   TO WS-DT-TOTAL-REV(WS-DATE-COUNT)
           END-IF.

      ******************************************************************
      * 2400-ACCUM-PRODUCT — Cari/tambah produk, akumulasi
      ******************************************************************
       2400-ACCUM-PRODUCT.
           MOVE 0 TO WS-FOUND
           PERFORM VARYING WS-I FROM 1 BY 1
               UNTIL WS-I > WS-PROD-COUNT OR WS-FOUND = 1
               IF FUNCTION TRIM(WS-PR-NAME(WS-I) TRAILING) =
                  WS-TRIMMED-NAME
                   ADD WS-PARSED-QTY
                       TO WS-PR-QTY-SOLD(WS-I)
                   ADD WS-PARSED-NET-TOTAL
                       TO WS-PR-TOTAL-REV(WS-I)
                   MOVE 1 TO WS-FOUND
               END-IF
           END-PERFORM

           IF WS-FOUND = 0
               ADD 1 TO WS-PROD-COUNT
               MOVE WS-TRIMMED-NAME
                   TO WS-PR-NAME(WS-PROD-COUNT)
               MOVE WS-PARSED-QTY
                   TO WS-PR-QTY-SOLD(WS-PROD-COUNT)
               MOVE WS-PARSED-NET-TOTAL
                   TO WS-PR-TOTAL-REV(WS-PROD-COUNT)
           END-IF.

      ******************************************************************
      * 3000-CALC-AVERAGES — Hitung rata-rata dan grand total
      ******************************************************************
       3000-CALC-AVERAGES.
           MOVE 0 TO WS-GRAND-TRX
           MOVE 0 TO WS-GRAND-REVENUE

           PERFORM VARYING WS-I FROM 1 BY 1
               UNTIL WS-I > WS-CAT-COUNT
               ADD WS-CAT-TRX-COUNT(WS-I) TO WS-GRAND-TRX
               ADD WS-CAT-TOTAL-REV(WS-I) TO WS-GRAND-REVENUE
               IF WS-CAT-TRX-COUNT(WS-I) > 0
                   DIVIDE WS-CAT-TOTAL-REV(WS-I)
                       BY WS-CAT-TRX-COUNT(WS-I)
                       GIVING WS-CAT-AVG-REV(WS-I)
                       ROUNDED
               END-IF
           END-PERFORM

           IF WS-GRAND-TRX > 0
               DIVIDE WS-GRAND-REVENUE BY WS-GRAND-TRX
                   GIVING WS-GRAND-AVG
                   ROUNDED
           END-IF

           DISPLAY "[INFO] Grand Total: " WS-GRAND-TRX
               " transaksi, Rp " WS-GRAND-REVENUE.

      ******************************************************************
      * 4000-SORT-DATES — Bubble sort tanggal ascending (YYYY-MM-DD)
      ******************************************************************
       4000-SORT-DATES.
           PERFORM VARYING WS-I FROM 1 BY 1
               UNTIL WS-I >= WS-DATE-COUNT
               PERFORM VARYING WS-J FROM 1 BY 1
                   UNTIL WS-J > WS-DATE-COUNT - WS-I
                   IF WS-DT-DATE(WS-J) > WS-DT-DATE(WS-J + 1)
                       MOVE WS-DT-DATE(WS-J)
                           TO WS-SORT-TEMP-DATE
                       MOVE WS-DT-TRX-COUNT(WS-J)
                           TO WS-SORT-TEMP-TRX
                       MOVE WS-DT-TOTAL-REV(WS-J)
                           TO WS-SORT-TEMP-REV

                       MOVE WS-DT-DATE(WS-J + 1)
                           TO WS-DT-DATE(WS-J)
                       MOVE WS-DT-TRX-COUNT(WS-J + 1)
                           TO WS-DT-TRX-COUNT(WS-J)
                       MOVE WS-DT-TOTAL-REV(WS-J + 1)
                           TO WS-DT-TOTAL-REV(WS-J)

                       MOVE WS-SORT-TEMP-DATE
                           TO WS-DT-DATE(WS-J + 1)
                       MOVE WS-SORT-TEMP-TRX
                           TO WS-DT-TRX-COUNT(WS-J + 1)
                       MOVE WS-SORT-TEMP-REV
                           TO WS-DT-TOTAL-REV(WS-J + 1)
                   END-IF
               END-PERFORM
           END-PERFORM.

      ******************************************************************
      * 4100-SORT-PRODUCTS — Bubble sort produk by qty descending
      ******************************************************************
       4100-SORT-PRODUCTS.
           PERFORM VARYING WS-I FROM 1 BY 1
               UNTIL WS-I >= WS-PROD-COUNT
               PERFORM VARYING WS-J FROM 1 BY 1
                   UNTIL WS-J > WS-PROD-COUNT - WS-I
                   IF WS-PR-QTY-SOLD(WS-J) <
                      WS-PR-QTY-SOLD(WS-J + 1)
                       MOVE WS-PR-NAME(WS-J)
                           TO WS-SORT-TEMP-NAME
                       MOVE WS-PR-QTY-SOLD(WS-J)
                           TO WS-SORT-TEMP-QTY
                       MOVE WS-PR-TOTAL-REV(WS-J)
                           TO WS-SORT-TEMP-REV

                       MOVE WS-PR-NAME(WS-J + 1)
                           TO WS-PR-NAME(WS-J)
                       MOVE WS-PR-QTY-SOLD(WS-J + 1)
                           TO WS-PR-QTY-SOLD(WS-J)
                       MOVE WS-PR-TOTAL-REV(WS-J + 1)
                           TO WS-PR-TOTAL-REV(WS-J)

                       MOVE WS-SORT-TEMP-NAME
                           TO WS-PR-NAME(WS-J + 1)
                       MOVE WS-SORT-TEMP-QTY
                           TO WS-PR-QTY-SOLD(WS-J + 1)
                       MOVE WS-SORT-TEMP-REV
                           TO WS-PR-TOTAL-REV(WS-J + 1)
                   END-IF
               END-PERFORM
           END-PERFORM.

      ******************************************************************
      * 5000-WRITE-REPORT — Tulis ketiga bagian laporan
      ******************************************************************
       5000-WRITE-REPORT.
           PERFORM 5100-WRITE-CATEGORY-REPORT
           PERFORM 5200-WRITE-DAILY-REPORT
           PERFORM 5300-WRITE-PRODUCT-REPORT
           DISPLAY "[INFO] Laporan berhasil ditulis ke "
               "data/report_output.txt".

      ******************************************************************
      * 5100-WRITE-CATEGORY-REPORT — Bagian A: Rekapitulasi Kategori
      ******************************************************************
       5100-WRITE-CATEGORY-REPORT.
      *    === Header Utama ===
           WRITE OUTPUT-RECORD FROM WS-DOUBLE-LINE

           MOVE SPACES TO WS-OUT-LINE
           STRING "        LAPORAN PENDAPATAN - "
               "POINT OF SALES SYSTEM"
               DELIMITED BY SIZE INTO WS-OUT-LINE
           WRITE OUTPUT-RECORD FROM WS-OUT-LINE

      *    Konversi dan tulis periode
           PERFORM 6000-CONVERT-MIN-DATE
           PERFORM 6100-CONVERT-MAX-DATE

           MOVE SPACES TO WS-OUT-LINE
           STRING "        Periode: "
               DELIMITED BY SIZE
               WS-CONV-RESULT-MIN DELIMITED BY "  "
               " s/d " DELIMITED BY SIZE
               WS-CONV-RESULT DELIMITED BY "  "
               INTO WS-OUT-LINE
           END-STRING
           WRITE OUTPUT-RECORD FROM WS-OUT-LINE

           WRITE OUTPUT-RECORD FROM WS-DOUBLE-LINE

      *    Baris kosong
           MOVE SPACES TO WS-OUT-LINE
           WRITE OUTPUT-RECORD FROM WS-OUT-LINE

      *    Sub-header kategori
           MOVE SPACES TO WS-OUT-LINE
           MOVE "REKAPITULASI PER KATEGORI:" TO WS-OUT-LINE
           WRITE OUTPUT-RECORD FROM WS-OUT-LINE

           WRITE OUTPUT-RECORD FROM WS-SEPARATOR-LINE

           MOVE SPACES TO WS-OUT-LINE
           STRING "Kategori             "
               "Jml Trx   Total Penjualan   Rata-rata"
               DELIMITED BY SIZE INTO WS-OUT-LINE
           WRITE OUTPUT-RECORD FROM WS-OUT-LINE

           WRITE OUTPUT-RECORD FROM WS-SEPARATOR-LINE

      *    Detail per kategori
           PERFORM VARYING WS-I FROM 1 BY 1
               UNTIL WS-I > WS-CAT-COUNT
               PERFORM 5110-WRITE-CAT-LINE
           END-PERFORM

      *    Garis pemisah sebelum Grand Total
           WRITE OUTPUT-RECORD FROM WS-SEPARATOR-LINE

      *    Grand Total
           PERFORM 5120-WRITE-GRAND-TOTAL

      *    Penutup bagian A
           WRITE OUTPUT-RECORD FROM WS-DOUBLE-LINE.

      ******************************************************************
      * 5110-WRITE-CAT-LINE — Tulis satu baris detail kategori
      * Format: Kategori(1:20) JmlTrx(21:9) Rp Total(35:17) Rp Avg
      ******************************************************************
       5110-WRITE-CAT-LINE.
           MOVE SPACES TO WS-OUT-LINE

      *    Kolom 1: Nama kategori (posisi 1-20)
           MOVE WS-CAT-NAME(WS-I)(1:20) TO WS-OUT-LINE(1:20)

      *    Kolom 2: Jumlah transaksi (posisi 22-30)
           MOVE WS-CAT-TRX-COUNT(WS-I) TO WS-FMT-TRX
           MOVE WS-FMT-TRX TO WS-OUT-LINE(22:9)

      *    Kolom 3: Total Penjualan (posisi 35-50)
           MOVE WS-CAT-TOTAL-REV(WS-I) TO WS-FMT-BIG
           MOVE "Rp" TO WS-OUT-LINE(35:2)
           MOVE WS-FMT-BIG TO WS-OUT-LINE(37:15)

      *    Kolom 4: Rata-rata (posisi 55-68)
           MOVE WS-CAT-AVG-REV(WS-I) TO WS-FMT-NUM
           MOVE "Rp" TO WS-OUT-LINE(55:2)
           MOVE WS-FMT-NUM TO WS-OUT-LINE(57:11)

           WRITE OUTPUT-RECORD FROM WS-OUT-LINE.

      ******************************************************************
      * 5120-WRITE-GRAND-TOTAL — Tulis baris Grand Total
      ******************************************************************
       5120-WRITE-GRAND-TOTAL.
           MOVE SPACES TO WS-OUT-LINE

      *    Label
           MOVE "GRAND TOTAL" TO WS-OUT-LINE(1:20)

      *    Jumlah transaksi
           MOVE WS-GRAND-TRX TO WS-FMT-TRX
           MOVE WS-FMT-TRX TO WS-OUT-LINE(22:9)

      *    Grand total revenue
           MOVE WS-GRAND-REVENUE TO WS-FMT-BIG
           MOVE "Rp" TO WS-OUT-LINE(35:2)
           MOVE WS-FMT-BIG TO WS-OUT-LINE(37:15)

      *    Rata-rata
           MOVE WS-GRAND-AVG TO WS-FMT-NUM
           MOVE "Rp" TO WS-OUT-LINE(55:2)
           MOVE WS-FMT-NUM TO WS-OUT-LINE(57:11)

           WRITE OUTPUT-RECORD FROM WS-OUT-LINE.

      ******************************************************************
      * 5200-WRITE-DAILY-REPORT — Bagian B: Penjualan Harian
      ******************************************************************
       5200-WRITE-DAILY-REPORT.
      *    Baris kosong pemisah
           MOVE SPACES TO WS-OUT-LINE
           WRITE OUTPUT-RECORD FROM WS-OUT-LINE

      *    Header
           WRITE OUTPUT-RECORD FROM WS-DOUBLE-LINE

           MOVE SPACES TO WS-OUT-LINE
           STRING "        LAPORAN PENJUALAN HARIAN"
               DELIMITED BY SIZE INTO WS-OUT-LINE
           WRITE OUTPUT-RECORD FROM WS-OUT-LINE

           WRITE OUTPUT-RECORD FROM WS-DOUBLE-LINE

      *    Kolom header
           MOVE SPACES TO WS-OUT-LINE
           STRING "Tanggal         "
               "Jml Trx   Total Penjualan"
               DELIMITED BY SIZE INTO WS-OUT-LINE
           WRITE OUTPUT-RECORD FROM WS-OUT-LINE

           WRITE OUTPUT-RECORD FROM WS-SEPARATOR-LINE

      *    Detail per tanggal
           PERFORM VARYING WS-I FROM 1 BY 1
               UNTIL WS-I > WS-DATE-COUNT
               PERFORM 5210-WRITE-DATE-LINE
           END-PERFORM

      *    Penutup bagian B
           WRITE OUTPUT-RECORD FROM WS-DOUBLE-LINE.

      ******************************************************************
      * 5210-WRITE-DATE-LINE — Tulis satu baris penjualan harian
      ******************************************************************
       5210-WRITE-DATE-LINE.
           MOVE SPACES TO WS-OUT-LINE

      *    Konversi tanggal YYYY-MM-DD ke DD-Mon-YYYY
           MOVE WS-DT-DATE(WS-I)(1:4) TO WS-CONV-YEAR
           MOVE WS-DT-DATE(WS-I)(6:2) TO WS-CONV-MONTH
           MOVE WS-DT-DATE(WS-I)(9:2) TO WS-CONV-DAY
           PERFORM 6200-GET-MONTH-NAME

           INITIALIZE WS-CONV-RESULT
           STRING WS-CONV-DAY DELIMITED BY SIZE
               "-" DELIMITED BY SIZE
               WS-CONV-MON-NAME DELIMITED BY SIZE
               "-" DELIMITED BY SIZE
               WS-CONV-YEAR DELIMITED BY SIZE
               INTO WS-CONV-RESULT
           END-STRING

      *    Kolom 1: Tanggal (posisi 1-16)
           MOVE WS-CONV-RESULT(1:11) TO WS-OUT-LINE(1:16)

      *    Kolom 2: Jumlah transaksi (posisi 17-25)
           MOVE WS-DT-TRX-COUNT(WS-I) TO WS-FMT-TRX
           MOVE WS-FMT-TRX TO WS-OUT-LINE(17:9)

      *    Kolom 3: Total penjualan (posisi 31-48)
           MOVE WS-DT-TOTAL-REV(WS-I) TO WS-FMT-BIG
           MOVE "Rp" TO WS-OUT-LINE(31:2)
           MOVE WS-FMT-BIG TO WS-OUT-LINE(33:15)

           WRITE OUTPUT-RECORD FROM WS-OUT-LINE.

      ******************************************************************
      * 5300-WRITE-PRODUCT-REPORT — Bagian C: Produk Terlaris Top 10
      ******************************************************************
       5300-WRITE-PRODUCT-REPORT.
      *    Baris kosong pemisah
           MOVE SPACES TO WS-OUT-LINE
           WRITE OUTPUT-RECORD FROM WS-OUT-LINE

      *    Header
           WRITE OUTPUT-RECORD FROM WS-DOUBLE-LINE

           MOVE SPACES TO WS-OUT-LINE
           STRING "        LAPORAN PRODUK TERLARIS"
               DELIMITED BY SIZE INTO WS-OUT-LINE
           WRITE OUTPUT-RECORD FROM WS-OUT-LINE

           WRITE OUTPUT-RECORD FROM WS-DOUBLE-LINE

      *    Kolom header
           MOVE SPACES TO WS-OUT-LINE
           STRING "No  Produk              "
               "Qty Terjual   Total Pendapatan"
               DELIMITED BY SIZE INTO WS-OUT-LINE
           WRITE OUTPUT-RECORD FROM WS-OUT-LINE

           WRITE OUTPUT-RECORD FROM WS-SEPARATOR-LINE

      *    Tentukan batas top (min antara 10 dan jumlah produk)
           IF WS-PROD-COUNT < WS-TOP-LIMIT
               MOVE WS-PROD-COUNT TO WS-TOP-LIMIT
           END-IF

      *    Detail top produk
           PERFORM VARYING WS-I FROM 1 BY 1
               UNTIL WS-I > WS-TOP-LIMIT
               PERFORM 5310-WRITE-PROD-LINE
           END-PERFORM

      *    Penutup bagian C
           WRITE OUTPUT-RECORD FROM WS-DOUBLE-LINE.

      ******************************************************************
      * 5310-WRITE-PROD-LINE — Tulis satu baris produk terlaris
      ******************************************************************
       5310-WRITE-PROD-LINE.
           MOVE SPACES TO WS-OUT-LINE

      *    Kolom 1: Nomor urut (posisi 1-3)
           MOVE WS-I TO WS-NUM-INDEX
           MOVE WS-NUM-INDEX TO WS-OUT-LINE(1:3)

      *    Kolom 2: Nama produk (posisi 5-24, 20 karakter)
           MOVE WS-PR-NAME(WS-I)(1:20)
               TO WS-OUT-LINE(5:20)

      *    Kolom 3: Qty terjual (posisi 25-37)
           MOVE WS-PR-QTY-SOLD(WS-I) TO WS-FMT-TRX
           MOVE WS-FMT-TRX TO WS-OUT-LINE(25:13)

      *    Kolom 4: Total pendapatan (posisi 41-58)
           MOVE WS-PR-TOTAL-REV(WS-I) TO WS-FMT-BIG
           MOVE "Rp" TO WS-OUT-LINE(41:2)
           MOVE WS-FMT-BIG TO WS-OUT-LINE(43:15)

           WRITE OUTPUT-RECORD FROM WS-OUT-LINE.

      ******************************************************************
      * 6000-CONVERT-MIN-DATE — Konversi min date ke DD-Mon-YYYY
      ******************************************************************
       6000-CONVERT-MIN-DATE.
           MOVE WS-MIN-DATE(1:4) TO WS-CONV-YEAR
           MOVE WS-MIN-DATE(6:2) TO WS-CONV-MONTH
           MOVE WS-MIN-DATE(9:2) TO WS-CONV-DAY
           PERFORM 6200-GET-MONTH-NAME
           INITIALIZE WS-CONV-RESULT-MIN
           STRING WS-CONV-DAY DELIMITED BY SIZE
               "-" DELIMITED BY SIZE
               WS-CONV-MON-NAME DELIMITED BY SIZE
               "-" DELIMITED BY SIZE
               WS-CONV-YEAR DELIMITED BY SIZE
               INTO WS-CONV-RESULT-MIN
           END-STRING.

      ******************************************************************
      * 6100-CONVERT-MAX-DATE — Konversi max date ke DD-Mon-YYYY
      ******************************************************************
       6100-CONVERT-MAX-DATE.
           MOVE WS-MAX-DATE(1:4) TO WS-CONV-YEAR
           MOVE WS-MAX-DATE(6:2) TO WS-CONV-MONTH
           MOVE WS-MAX-DATE(9:2) TO WS-CONV-DAY
           PERFORM 6200-GET-MONTH-NAME
           INITIALIZE WS-CONV-RESULT
           STRING WS-CONV-DAY DELIMITED BY SIZE
               "-" DELIMITED BY SIZE
               WS-CONV-MON-NAME DELIMITED BY SIZE
               "-" DELIMITED BY SIZE
               WS-CONV-YEAR DELIMITED BY SIZE
               INTO WS-CONV-RESULT
           END-STRING.

      ******************************************************************
      * 6200-GET-MONTH-NAME — Konversi nomor bulan ke nama singkat
      ******************************************************************
       6200-GET-MONTH-NAME.
           EVALUATE WS-CONV-MONTH
               WHEN "01" MOVE "Jan" TO WS-CONV-MON-NAME
               WHEN "02" MOVE "Feb" TO WS-CONV-MON-NAME
               WHEN "03" MOVE "Mar" TO WS-CONV-MON-NAME
               WHEN "04" MOVE "Apr" TO WS-CONV-MON-NAME
               WHEN "05" MOVE "Mei" TO WS-CONV-MON-NAME
               WHEN "06" MOVE "Jun" TO WS-CONV-MON-NAME
               WHEN "07" MOVE "Jul" TO WS-CONV-MON-NAME
               WHEN "08" MOVE "Agu" TO WS-CONV-MON-NAME
               WHEN "09" MOVE "Sep" TO WS-CONV-MON-NAME
               WHEN "10" MOVE "Okt" TO WS-CONV-MON-NAME
               WHEN "11" MOVE "Nov" TO WS-CONV-MON-NAME
               WHEN "12" MOVE "Des" TO WS-CONV-MON-NAME
               WHEN OTHER MOVE "???" TO WS-CONV-MON-NAME
           END-EVALUATE.

      ******************************************************************
      * 9000-CLEANUP — Tutup semua file
      ******************************************************************
       9000-CLEANUP.
           CLOSE INPUT-FILE
           CLOSE OUTPUT-FILE
           DISPLAY "[INFO] Semua file berhasil ditutup".
