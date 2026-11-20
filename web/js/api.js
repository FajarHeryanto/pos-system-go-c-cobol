// ==================================================================
// api.js — Fetch API Wrapper dengan JWT Authentication
// Sistem Informasi Penjualan (Point of Sales)
// ==================================================================
//
// Fitur:
//   - Auto-attach Authorization: Bearer <token> dari localStorage
//   - Auto-redirect ke /login jika menerima 401 Unauthorized
//   - Objek API: ProductAPI, TransactionAPI, ProcessAPI, DashboardAPI
//
// Sesuai Bab 6.4 RPD.
// ==================================================================

const API_BASE = '/api';

// ================================================================
// CORE: Fetch wrapper dengan auth token
// ================================================================
async function apiRequest(endpoint, method = 'GET', body = null) {
    const token = localStorage.getItem('pos_token');
    const headers = {
        'Content-Type': 'application/json',
    };
    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    const options = { method, headers };
    if (body) {
        options.body = JSON.stringify(body);
    }

    try {
        const response = await fetch(`${API_BASE}${endpoint}`, options);
        const data = await response.json();

        // Auto-redirect ke login jika 401
        if (response.status === 401) {
            localStorage.removeItem('pos_token');
            localStorage.removeItem('pos_user');
            window.location.href = '/login';
            throw new Error('Sesi telah berakhir. Silakan login ulang.');
        }

        if (!response.ok) {
            const errMsg = data.message || `Request gagal (${response.status})`;
            throw new Error(errMsg);
        }

        return data;
    } catch (err) {
        if (err.name === 'TypeError' && err.message.includes('fetch')) {
            throw new Error('Tidak dapat terhubung ke server. Pastikan server berjalan.');
        }
        throw err;
    }
}

// ================================================================
// AUTH API
// ================================================================
const AuthAPI = {
    login: (username, password) =>
        apiRequest('/auth/login', 'POST', { username, password }),
};

// ================================================================
// PRODUCT API
// ================================================================
const ProductAPI = {
    getAll:    ()            => apiRequest('/products'),
    getById:   (id)          => apiRequest(`/products/${id}`),
    create:    (product)     => apiRequest('/products', 'POST', product),
    update:    (id, product) => apiRequest(`/products/${id}`, 'PUT', product),
    delete:    (id)          => apiRequest(`/products/${id}`, 'DELETE'),
};

// ================================================================
// TRANSACTION API
// ================================================================
const TransactionAPI = {
    getAll:    ()    => apiRequest('/transactions'),
    create:    (trx) => apiRequest('/transactions', 'POST', trx),
    delete:    (id)  => apiRequest(`/transactions/${id}`, 'DELETE'),
};

// ================================================================
// PROCESS & REPORT API
// ================================================================
const ProcessAPI = {
    runCEngine:       () => apiRequest('/process', 'POST'),
    generateReport:   () => apiRequest('/reports/generate', 'POST'),
    getIncomeReport:  () => apiRequest('/reports/income'),
};

// ================================================================
// DASHBOARD API
// ================================================================
const DashboardAPI = {
    getStats: () => apiRequest('/dashboard'),
};

// ================================================================
// UTILITY: Cek apakah user sudah login
// ================================================================
function isLoggedIn() {
    return !!localStorage.getItem('pos_token');
}

function getCurrentUser() {
    const user = localStorage.getItem('pos_user');
    return user ? JSON.parse(user) : null;
}

function logout() {
    localStorage.removeItem('pos_token');
    localStorage.removeItem('pos_user');
    window.location.href = '/login';
}

// Cek auth pada halaman yang memerlukan login
function requireAuth() {
    if (!isLoggedIn()) {
        window.location.href = '/login';
    }
}

// ================================================================
// UTILITY: Format Rupiah
// ================================================================
function formatRupiah(amount) {
    if (amount === null || amount === undefined) return 'Rp 0';
    return 'Rp ' + Math.round(amount).toLocaleString('id-ID');
}

// ================================================================
// UTILITY: Format tanggal
// ================================================================
function formatDate(dateStr) {
    if (!dateStr) return '-';
    const months = ['Jan','Feb','Mar','Apr','Mei','Jun',
                    'Jul','Agu','Sep','Okt','Nov','Des'];
    const parts = dateStr.split('-');
    if (parts.length !== 3) return dateStr;
    const day = parseInt(parts[2], 10);
    const monthIdx = parseInt(parts[1], 10) - 1;
    return `${day}-${months[monthIdx]}-${parts[0]}`;
}

// ================================================================
// UTILITY: Get badge class untuk kategori
// ================================================================
function getCategoryBadge(category) {
    const map = {
        'Makanan': 'badge-makanan',
        'Minuman': 'badge-minuman',
        'Elektronik': 'badge-elektronik',
        'Alat Tulis': 'badge-alat-tulis',
        'Kebutuhan RT': 'badge-kebutuhan-rt',
    };
    return map[category] || 'bg-secondary';
}

function getPaymentBadge(method) {
    const map = {
        'Cash': 'badge-cash',
        'Debit': 'badge-debit',
        'Credit': 'badge-credit',
    };
    return map[method] || 'bg-secondary';
}

// ================================================================
// UTILITY: Show Toast Notification
// ================================================================
function showToast(message, type = 'success') {
    // Buat container jika belum ada
    let container = document.querySelector('.toast-container');
    if (!container) {
        container = document.createElement('div');
        container.className = 'toast-container';
        document.body.appendChild(container);
    }

    const bgClass = type === 'success' ? 'bg-success' :
                    type === 'error' ? 'bg-danger' :
                    type === 'warning' ? 'bg-warning text-dark' : 'bg-info';

    const toastHTML = `
        <div class="toast show align-items-center text-white ${bgClass} border-0 mb-2" role="alert">
            <div class="d-flex">
                <div class="toast-body">${message}</div>
                <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
            </div>
        </div>`;

    container.insertAdjacentHTML('beforeend', toastHTML);

    const toast = container.lastElementChild;
    setTimeout(() => {
        toast.classList.remove('show');
        setTimeout(() => toast.remove(), 300);
    }, 4000);

    // Close button
    toast.querySelector('.btn-close').addEventListener('click', () => {
        toast.classList.remove('show');
        setTimeout(() => toast.remove(), 300);
    });
}
