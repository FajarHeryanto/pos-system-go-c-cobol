// ==================================================================
// products.js — Logic Halaman Manajemen Produk
// Sistem Informasi Penjualan (Point of Sales)
// ==================================================================

let allProducts = [];
let productModal, deleteModal;
let deleteProductId = null;

document.addEventListener('DOMContentLoaded', () => {
    requireAuth();
    loadUserInfo();

    productModal = new bootstrap.Modal(document.getElementById('productModal'));
    deleteModal = new bootstrap.Modal(document.getElementById('deleteModal'));

    loadProducts();

    // Event listeners
    document.getElementById('searchInput').addEventListener('input', filterProducts);
    document.getElementById('filterCategory').addEventListener('change', filterProducts);
    document.getElementById('productForm').addEventListener('submit', handleSaveProduct);
    document.getElementById('btnConfirmDelete').addEventListener('click', handleDeleteProduct);
});

function loadUserInfo() {
    const user = getCurrentUser();
    if (user) {
        document.getElementById('userInfo').textContent = `${user.full_name} (${user.role})`;
    }
}

// ================================================================
// LOAD & RENDER PRODUCTS
// ================================================================
async function loadProducts() {
    try {
        const result = await ProductAPI.getAll();
        allProducts = result.data || [];
        filterProducts();
    } catch (err) {
        showToast('Gagal memuat produk: ' + err.message, 'error');
    }
}

function filterProducts() {
    const search = document.getElementById('searchInput').value.toLowerCase();
    const category = document.getElementById('filterCategory').value;

    let filtered = allProducts;

    if (search) {
        filtered = filtered.filter(p =>
            p.name.toLowerCase().includes(search) ||
            p.product_id.toLowerCase().includes(search)
        );
    }

    if (category) {
        filtered = filtered.filter(p => p.category === category);
    }

    renderProductTable(filtered);
    document.getElementById('productCount').textContent =
        `Menampilkan ${filtered.length} dari ${allProducts.length} produk`;
}

function renderProductTable(products) {
    const tbody = document.getElementById('productTableBody');

    if (products.length === 0) {
        tbody.innerHTML = '<tr><td colspan="6" class="text-center text-muted py-4">Tidak ada produk ditemukan</td></tr>';
        return;
    }

    tbody.innerHTML = products.map(p => `
        <tr>
            <td><code>${p.product_id}</code></td>
            <td><strong>${p.name}</strong></td>
            <td><span class="badge badge-category ${getCategoryBadge(p.category)}">${p.category}</span></td>
            <td class="text-rupiah">${formatRupiah(p.price)}</td>
            <td>
                <span class="${p.stock <= 10 ? 'text-danger fw-bold' : ''}">${p.stock}</span>
                ${p.stock <= 10 ? '<i class="bi bi-exclamation-triangle text-warning ms-1" title="Stok rendah"></i>' : ''}
            </td>
            <td class="text-center">
                <button class="btn btn-sm btn-outline-primary me-1" onclick="openEditModal('${p.product_id}')" title="Edit">
                    <i class="bi bi-pencil"></i>
                </button>
                <button class="btn btn-sm btn-outline-danger" onclick="openDeleteModal('${p.product_id}', '${p.name}')" title="Hapus">
                    <i class="bi bi-trash"></i>
                </button>
            </td>
        </tr>
    `).join('');
}

// ================================================================
// ADD / EDIT PRODUCT
// ================================================================
function openAddModal() {
    document.getElementById('modalTitle').textContent = 'Tambah Produk Baru';
    document.getElementById('editProductId').value = '';
    document.getElementById('productForm').reset();
    productModal.show();
}

async function openEditModal(id) {
    try {
        const result = await ProductAPI.getById(id);
        const p = result.data;

        document.getElementById('modalTitle').textContent = 'Edit Produk';
        document.getElementById('editProductId').value = p.product_id;
        document.getElementById('prodName').value = p.name;
        document.getElementById('prodCategory').value = p.category;
        document.getElementById('prodPrice').value = p.price;
        document.getElementById('prodStock').value = p.stock;

        productModal.show();
    } catch (err) {
        showToast('Gagal memuat data produk: ' + err.message, 'error');
    }
}

async function handleSaveProduct(e) {
    e.preventDefault();

    const btn = document.getElementById('btnSaveProduct');
    btn.disabled = true;
    btn.innerHTML = '<span class="spinner-border spinner-border-sm me-1"></span>Menyimpan...';

    const productData = {
        name: document.getElementById('prodName').value.trim(),
        category: document.getElementById('prodCategory').value,
        price: parseFloat(document.getElementById('prodPrice').value),
        stock: parseInt(document.getElementById('prodStock').value) || 0,
    };

    const editId = document.getElementById('editProductId').value;

    try {
        if (editId) {
            await ProductAPI.update(editId, productData);
            showToast('Produk berhasil diupdate');
        } else {
            await ProductAPI.create(productData);
            showToast('Produk berhasil ditambahkan');
        }

        productModal.hide();
        loadProducts();
    } catch (err) {
        showToast(err.message, 'error');
    } finally {
        btn.disabled = false;
        btn.innerHTML = '<i class="bi bi-check-lg me-1"></i>Simpan';
    }
}

// ================================================================
// DELETE PRODUCT
// ================================================================
function openDeleteModal(id, name) {
    deleteProductId = id;
    document.getElementById('deleteProductName').textContent = name;
    deleteModal.show();
}

async function handleDeleteProduct() {
    if (!deleteProductId) return;

    const btn = document.getElementById('btnConfirmDelete');
    btn.disabled = true;

    try {
        await ProductAPI.delete(deleteProductId);
        showToast('Produk berhasil dihapus');
        deleteModal.hide();
        loadProducts();
    } catch (err) {
        showToast(err.message, 'error');
    } finally {
        btn.disabled = false;
        deleteProductId = null;
    }
}
