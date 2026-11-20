// ==================================================================
// dashboard.js — Logic Halaman Dashboard
// Sistem Informasi Penjualan (Point of Sales)
// ==================================================================

document.addEventListener('DOMContentLoaded', () => {
    requireAuth();
    loadUserInfo();
    loadDashboard();
});

function loadUserInfo() {
    const user = getCurrentUser();
    if (user) {
        document.getElementById('userInfo').textContent =
            `${user.full_name} (${user.role})`;
    }
}

async function loadDashboard() {
    try {
        const result = await DashboardAPI.getStats();
        const stats = result.data;

        // Update stat cards
        document.getElementById('statProducts').textContent = stats.total_products;
        document.getElementById('statTransactions').textContent = stats.total_transactions;
        document.getElementById('statRevenue').textContent = formatRupiah(stats.total_revenue);
        document.getElementById('statCashiers').textContent = stats.total_cashiers;

        // Render chart
        renderSalesChart(stats.category_stats);

        // Render top products
        renderTopProducts(stats.top_products);

        // Render recent transactions
        renderRecentTransactions(stats.recent_transactions);

    } catch (err) {
        showToast('Gagal memuat dashboard: ' + err.message, 'error');
    }
}

function renderSalesChart(categoryStats) {
    const ctx = document.getElementById('salesChart').getContext('2d');

    const labels = categoryStats.map(c => c.category);
    const revenues = categoryStats.map(c => c.revenue);
    const counts = categoryStats.map(c => c.count);

    const colors = ['#ef4444', '#3b82f6', '#f59e0b', '#10b981', '#6366f1'];

    new Chart(ctx, {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [
                {
                    label: 'Pendapatan (Rp)',
                    data: revenues,
                    backgroundColor: colors.map(c => c + 'cc'),
                    borderColor: colors,
                    borderWidth: 1,
                    borderRadius: 4,
                    yAxisID: 'y',
                },
                {
                    label: 'Jumlah Transaksi',
                    data: counts,
                    type: 'line',
                    borderColor: '#6366f1',
                    backgroundColor: '#6366f1',
                    tension: 0.3,
                    pointRadius: 5,
                    yAxisID: 'y1',
                }
            ]
        },
        options: {
            responsive: true,
            interaction: {
                mode: 'index',
                intersect: false,
            },
            plugins: {
                legend: { position: 'bottom' },
                tooltip: {
                    callbacks: {
                        label: function(context) {
                            if (context.datasetIndex === 0) {
                                return 'Pendapatan: ' + formatRupiah(context.raw);
                            }
                            return 'Transaksi: ' + context.raw;
                        }
                    }
                }
            },
            scales: {
                y: {
                    type: 'linear',
                    display: true,
                    position: 'left',
                    ticks: {
                        callback: (val) => formatRupiah(val)
                    }
                },
                y1: {
                    type: 'linear',
                    display: true,
                    position: 'right',
                    grid: { drawOnChartArea: false },
                }
            }
        }
    });
}

function renderTopProducts(products) {
    const container = document.getElementById('topProductsList');

    if (!products || products.length === 0) {
        container.innerHTML = '<div class="text-center p-3 text-muted">Belum ada data</div>';
        return;
    }

    container.innerHTML = products.map((p, i) => `
        <div class="list-group-item d-flex justify-content-between align-items-center px-3 py-2">
            <div class="d-flex align-items-center gap-2">
                <span class="badge ${i === 0 ? 'bg-warning text-dark' : 'bg-secondary'} rounded-pill">${i + 1}</span>
                <span>${p.product_name}</span>
            </div>
            <div class="text-end">
                <div class="fw-bold small">${p.total_qty} terjual</div>
                <div class="text-muted small">${formatRupiah(p.total_revenue)}</div>
            </div>
        </div>
    `).join('');
}

function renderRecentTransactions(transactions) {
    const tbody = document.getElementById('recentTransBody');

    if (!transactions || transactions.length === 0) {
        tbody.innerHTML = '<tr><td colspan="5" class="text-center text-muted py-3">Belum ada transaksi</td></tr>';
        return;
    }

    tbody.innerHTML = transactions.map(t => `
        <tr>
            <td><code>${t.transaction_id}</code></td>
            <td>${t.product_name}</td>
            <td>${t.quantity}</td>
            <td class="text-rupiah">${formatRupiah(t.total_price)}</td>
            <td>${formatDate(t.transaction_date)}</td>
        </tr>
    `).join('');
}
