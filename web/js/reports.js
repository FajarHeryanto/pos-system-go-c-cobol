// ==================================================================
// reports.js — Logic Halaman Laporan
// Sistem Informasi Penjualan (Point of Sales)
// ==================================================================

let step1Done = false;
let step2Done = false;

document.addEventListener('DOMContentLoaded', () => {
    requireAuth();
    loadUserInfo();

    // Coba muat laporan yang sudah ada
    loadReport();
});

function loadUserInfo() {
    const user = getCurrentUser();
    if (user) {
        document.getElementById('userInfo').textContent = `${user.full_name} (${user.role})`;
    }
}

// ================================================================
// STEP 1: Proses Data (C Engine)
// ================================================================
async function runStep1() {
    const btn = document.getElementById('btnStep1');
    const number = document.getElementById('step1Number');
    const status = document.getElementById('step1Status');

    btn.disabled = true;
    btn.innerHTML = '<span class="spinner-border spinner-border-sm me-1"></span>Memproses...';
    number.className = 'step-number processing';
    status.textContent = 'Sedang memproses...';

    try {
        const result = await ProcessAPI.runCEngine();
        const data = result.data;

        step1Done = true;
        number.className = 'step-number success';
        number.innerHTML = '<i class="bi bi-check-lg"></i>';
        status.textContent = `Selesai! ${data.records_processed} record diproses.`;
        status.className = 'step-status text-success';

        btn.innerHTML = '<i class="bi bi-check-circle me-1"></i>Selesai';
        btn.className = 'btn btn-success btn-sm';

        showToast(`C Engine berhasil memproses ${data.records_processed} record`);
    } catch (err) {
        number.className = 'step-number pending';
        status.textContent = 'Gagal: ' + err.message;
        status.className = 'step-status text-danger';

        btn.disabled = false;
        btn.innerHTML = '<i class="bi bi-gear-fill me-1"></i>Proses Data';

        showToast('Gagal: ' + err.message, 'error');
    }
}

// ================================================================
// STEP 2: Generate Laporan (COBOL Reporter)
// ================================================================
async function runStep2() {
    const btn = document.getElementById('btnStep2');
    const number = document.getElementById('step2Number');
    const status = document.getElementById('step2Status');

    btn.disabled = true;
    btn.innerHTML = '<span class="spinner-border spinner-border-sm me-1"></span>Generating...';
    number.className = 'step-number processing';
    status.textContent = 'Sedang generate laporan...';

    try {
        const result = await ProcessAPI.generateReport();

        step2Done = true;
        number.className = 'step-number success';
        number.innerHTML = '<i class="bi bi-check-lg"></i>';
        status.textContent = 'Selesai!';
        status.className = 'step-status text-success';

        btn.innerHTML = '<i class="bi bi-check-circle me-1"></i>Selesai';
        btn.className = 'btn btn-success btn-sm';

        showToast('COBOL Reporter berhasil generate laporan');

        // Otomatis muat laporan
        loadReport();
    } catch (err) {
        number.className = 'step-number pending';
        status.textContent = 'Gagal: ' + err.message;
        status.className = 'step-status text-danger';

        btn.disabled = false;
        btn.innerHTML = '<i class="bi bi-file-earmark-text-fill me-1"></i>Generate Laporan';

        showToast('Gagal: ' + err.message, 'error');
    }
}

// ================================================================
// LOAD REPORT
// ================================================================
async function loadReport() {
    try {
        const result = await ProcessAPI.getIncomeReport();
        const content = result.data.report_content;

        document.getElementById('reportContent').textContent = content;
        document.getElementById('btnPrint').disabled = false;

        // Update step indicators jika laporan sudah ada
        const number1 = document.getElementById('step1Number');
        const number2 = document.getElementById('step2Number');

        number1.className = 'step-number success';
        number1.innerHTML = '<i class="bi bi-check-lg"></i>';
        document.getElementById('step1Status').textContent = 'Sudah diproses';
        document.getElementById('step1Status').className = 'step-status text-success';

        number2.className = 'step-number success';
        number2.innerHTML = '<i class="bi bi-check-lg"></i>';
        document.getElementById('step2Status').textContent = 'Sudah digenerate';
        document.getElementById('step2Status').className = 'step-status text-success';

    } catch (err) {
        // Laporan belum ada — normal, bukan error
        document.getElementById('reportContent').textContent =
            'Laporan belum digenerate.\n\n' +
            'Ikuti langkah-langkah di atas:\n' +
            '  1. Klik "Proses Data" untuk menjalankan C Engine\n' +
            '  2. Klik "Generate Laporan" untuk menjalankan COBOL Reporter\n' +
            '  3. Laporan akan tampil di sini';
        document.getElementById('btnPrint').disabled = true;
    }
}

// ================================================================
// PRINT REPORT
// ================================================================
function printReport() {
    window.print();
}
