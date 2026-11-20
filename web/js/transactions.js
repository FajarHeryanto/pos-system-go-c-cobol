// ==================================================================
// transactions.js — Logic Halaman Transaksi
// Sistem Informasi Penjualan (Point of Sales)
// ==================================================================

let productList = [];
let receiptModal;

document.addEventListener('DOMContentLoaded', () => {
    requireAuth();
    loadUserInfo();

    receiptModal = new bootstrap.Modal(document.getElementById('receiptModal'));

    loadProducts();
    loadTransactions();

    // Event listeners
    document.getElementById('txProduk').addEventListener('change', updatePricePreview);
    document.getElementById('txQty').addEventListener('input', updatePricePreview);
    document.getElementById('txDiskon').addEventListener('input', updatePricePreview);
    document.getElementById('transactionForm').addEventListener('submit', handleCreateTransaction);
});

function loadUserInfo() {
    const user = getCurrentUser();
    if (user) {
        document.getElementById('userInfo').textContent = `${user.full_name} (${user.role})`;
    }
}

// ================================================================
// LOAD PRODUCTS INTO SELECT
// ================================================================
async function loadProducts() {
    try {
        const result = await ProductAPI.getAll();
        productList = result.data || [];

        const select = document.getElementById('txProduk');
        select.innerHTML = '<option value="">Pilih Produk...</option>';

        productList.forEach(p => {
            const option = document.createElement('option');
            option.value = p.product_id;
            option.textContent = `${p.name} — ${formatRupiah(p.price)} (Stok: ${p.stock})`;
            option.dataset.price = p.price;
            option.dataset.stock = p.stock;
            option.dataset.name = p.name;
            if (p.stock <= 0) {
                option.disabled = true;
                option.textContent += ' [HABIS]';
            }
            select.appendChild(option);
        });
    } catch (err) {
        showToast('Gagal memuat produk: ' + err.message, 'error');
    }
}

// ================================================================
// UPDATE PRICE PREVIEW
// ================================================================
function updatePricePreview() {
    const select = document.getElementById('txProduk');
    const selected = select.options[select.selectedIndex];
    const qty = parseInt(document.getElementById('txQty').value) || 0;
    const discount = parseFloat(document.getElementById('txDiskon').value) || 0;

    if (!selected || !selected.value) {
        document.getElementById('txSubtotal').textContent = 'Rp 0';
        document.getElementById('txDiscountAmt').textContent = '- Rp 0';
        document.getElementById('txTotal').textContent = 'Rp 0';
        document.getElementById('txProdukInfo').textContent = '';
        return;
    }

    const price = parseFloat(selected.dataset.price) || 0;
    const stock = parseInt(selected.dataset.stock) || 0;

    document.getElementById('txProdukInfo').textContent =
        `Harga: ${formatRupiah(price)} | Stok tersedia: ${stock}`;

    const subtotal = qty * price;
    const discountAmt = Math.round(subtotal * discount / 100);
    const total = subtotal - discountAmt;

    document.getElementById('txSubtotal').textContent = formatRupiah(subtotal);
    document.getElementById('txDiscountAmt').textContent = '- ' + formatRupiah(discountAmt);
    document.getElementById('txTotal').textContent = formatRupiah(total);
}

// ================================================================
// CREATE TRANSACTION
// ================================================================
async function handleCreateTransaction(e) {
    e.preventDefault();

    const btn = document.getElementById('btnSubmitTx');
    btn.disabled = true;
    btn.innerHTML = '<span class="spinner-border spinner-border-sm me-1"></span>Menyimpan...';

    const user = getCurrentUser();
    const cashierName = user ? user.full_name.split(' ')[0] : 'Admin';

    const paymentMethod = document.querySelector('input[name="txPayment"]:checked').value;

    const txData = {
        product_id: document.getElementById('txProduk').value,
        quantity: parseInt(document.getElementById('txQty').value),
        discount_percent: parseFloat(document.getElementById('txDiskon').value) || 0,
        payment_method: paymentMethod,
        cashier_name: cashierName,
    };

    try {
        const result = await TransactionAPI.create(txData);
        showToast('Transaksi berhasil disimpan!');

        // Tampilkan struk
        showReceipt(result.data);

        // Reset form
        document.getElementById('transactionForm').reset();
        document.getElementById('txProdukInfo').textContent = '';
        updatePricePreview();

        // Reload data
        loadProducts();
        loadTransactions();
    } catch (err) {
        showToast(err.message, 'error');
    } finally {
        btn.disabled = false;
        btn.innerHTML = '<i class="bi bi-check2-circle me-1"></i>Simpan Transaksi';
    }
}

// ================================================================
// LOAD TRANSACTIONS
// ================================================================
async function loadTransactions() {
    try {
        const result = await TransactionAPI.getAll();
        const transactions = result.data || [];

        document.getElementById('txCount').textContent = transactions.length;
        renderTransactionTable(transactions);
    } catch (err) {
        showToast('Gagal memuat transaksi: ' + err.message, 'error');
    }
}

function renderTransactionTable(transactions) {
    const tbody = document.getElementById('txTableBody');

    if (transactions.length === 0) {
        tbody.innerHTML = '<tr><td colspan="7" class="text-center text-muted py-4">Belum ada transaksi</td></tr>';
        return;
    }

    tbody.innerHTML = transactions.map(t => `
        <tr>
            <td><code>${t.transaction_id}</code></td>
            <td>${t.product_name}</td>
            <td>${t.quantity}</td>
            <td class="text-rupiah">${formatRupiah(t.total_price)}</td>
            <td><span class="badge badge-category ${getPaymentBadge(t.payment_method)}">${t.payment_method}</span></td>
            <td class="small">${formatDate(t.transaction_date)}</td>
            <td>
                <button class="btn btn-sm btn-outline-info me-1" onclick='showReceipt(${JSON.stringify(t)})' title="Struk">
                    <i class="bi bi-receipt"></i>
                </button>
                <button class="btn btn-sm btn-outline-danger" onclick="deleteTransaction('${t.transaction_id}')" title="Hapus">
                    <i class="bi bi-trash"></i>
                </button>
            </td>
        </tr>
    `).join('');
}

// ================================================================
// DELETE TRANSACTION
// ================================================================
async function deleteTransaction(id) {
    if (!confirm(`Yakin ingin menghapus transaksi ${id}?`)) return;

    try {
        await TransactionAPI.delete(id);
        showToast('Transaksi berhasil dihapus');
        loadProducts();
        loadTransactions();
    } catch (err) {
        showToast(err.message, 'error');
    }
}

// ================================================================
// RECEIPT / STRUK
// ================================================================
function showReceipt(tx) {
    const subtotal = tx.quantity * tx.unit_price;
    const discountAmt = Math.round(subtotal * tx.discount_percent / 100);

    const receiptHTML = `
        <div class="receipt-header">
            <strong style="font-size: 1.1rem;">POS SYSTEM</strong><br>
            <span>Sistem Informasi Penjualan</span><br>
            <small>${formatDate(tx.transaction_date)}</small>
        </div>
        <div style="margin-bottom: 0.5rem;">
            <strong>No: ${tx.transaction_id}</strong><br>
            Kasir: ${tx.cashier_name}
        </div>
        <div style="border-top: 1px dashed #999; border-bottom: 1px dashed #999; padding: 0.5rem 0; margin: 0.5rem 0;">
            ${tx.product_name}<br>
            ${tx.quantity} x ${formatRupiah(tx.unit_price)}<br>
            <div style="text-align: right;">
                Subtotal: ${formatRupiah(subtotal)}<br>
                ${tx.discount_percent > 0 ? `Diskon ${tx.discount_percent}%: -${formatRupiah(discountAmt)}<br>` : ''}
            </div>
        </div>
        <div style="text-align: right; font-size: 1.1rem;">
            <strong>TOTAL: ${formatRupiah(tx.total_price)}</strong>
        </div>
        <div style="margin-top: 0.5rem;">
            Bayar: ${tx.payment_method}
        </div>
        <div class="receipt-footer">
            Terima kasih atas kunjungan Anda!
        </div>
    `;

    document.getElementById('receiptContent').innerHTML = receiptHTML;
    receiptModal.show();
}
