<script lang="ts">
	import { onMount } from 'svelte';
	import type { Product, Order } from '@ecommerce/shared-types';
	import { sellerApi } from '$lib/api/seller';
	import { formatIDR, formatDate } from '$lib/utils/format';
	import SellerNav from '$lib/components/seller/SellerNav.svelte';
	import ProductFormModal from '$lib/components/seller/ProductFormModal.svelte';
	import QuickStockModal from '$lib/components/seller/QuickStockModal.svelte';
	import OrderStatusModal from '$lib/components/seller/OrderStatusModal.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import {
		Package,
		ShoppingBag,
		Clock,
		TrendingUp,
		AlertTriangle,
		ArrowRight,
		SlidersHorizontal,
		PlusCircle,
		CheckCircle2
	} from 'lucide-svelte';

	let loading = $state(true);
	let products = $state<Product[]>([]);
	let orders = $state<Order[]>([]);

	// Modals state
	let isProductModalOpen = $state(false);
	let selectedProduct = $state<Product | null>(null);

	let isQuickStockOpen = $state(false);
	let stockProduct = $state<Product | null>(null);

	let isOrderStatusOpen = $state(false);
	let selectedOrder = $state<Order | null>(null);

	async function loadDashboardData() {
		loading = true;
		try {
			const [prodRes, orderRes] = await Promise.all([
				sellerApi.getSellerProducts({ limit: 100 }),
				sellerApi.getIncomingOrders({ limit: 50 })
			]);
			products = prodRes.products;
			orders = orderRes.orders;
		} catch (e) {
			console.error('Failed to load seller dashboard data:', e);
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		loadDashboardData();
	});

	// Calculated Stats
	const totalProducts = $derived(products.length);

	const pendingOrdersCount = $derived(
		orders.filter((o) => o.status === 'pending' || o.status === 'confirmed').length
	);

	const shippedOrdersCount = $derived(orders.filter((o) => o.status === 'shipped').length);

	const totalCompletedRevenue = $derived(
		orders
			.filter((o) => o.status === 'completed')
			.reduce((sum, o) => sum + (o.total_amount || o.totalAmount || 0), 0)
	);

	const lowStockProducts = $derived(products.filter((p) => (p.stock ?? 0) <= 5).slice(0, 5));

	const recentOrders = $derived(orders.slice(0, 5));

	function openAddProduct() {
		selectedProduct = null;
		isProductModalOpen = true;
	}

	function openQuickStock(p: Product) {
		stockProduct = p;
		isQuickStockOpen = true;
	}

	function openOrderAction(o: Order) {
		selectedOrder = o;
		isOrderStatusOpen = true;
	}

	function handleProductSaved() {
		loadDashboardData();
	}

	function handleStockSaved(updated: Product) {
		products = products.map((p) => (p.id === updated.id ? { ...p, stock: updated.stock } : p));
	}

	function handleOrderUpdated(updated: Order) {
		orders = orders.map((o) => (o.id === updated.id ? { ...o, status: updated.status } : o));
	}
</script>

<svelte:head>
	<title>Ringkasan Toko - Seller Dashboard</title>
</svelte:head>

<SellerNav onAddProduct={openAddProduct} />

<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
	{#if loading}
		<div class="space-y-6">
			<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
				{#each [1, 2, 3, 4] as item (item)}
					<div class="h-28 animate-pulse rounded-xl bg-slate-200/80"></div>
				{/each}
			</div>
			<div class="h-64 animate-pulse rounded-xl bg-slate-200/80"></div>
		</div>
	{:else}
		<!-- Stats Grid -->
		<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
			<!-- Total Products -->
			<div class="rounded-xl border border-slate-200 bg-white p-5 shadow-xs">
				<div class="flex items-center justify-between">
					<span class="text-xs font-semibold tracking-wider text-slate-500 uppercase"
						>Katalog Produk</span
					>
					<div
						class="flex h-9 w-9 items-center justify-center rounded-lg bg-slate-100 text-slate-700"
					>
						<Package class="h-5 w-5" />
					</div>
				</div>
				<div class="mt-4 flex items-baseline gap-2">
					<span class="text-2xl font-bold tracking-tight text-slate-900 tabular-nums">
						{totalProducts}
					</span>
					<span class="text-xs text-slate-500">item aktif</span>
				</div>
				<div class="mt-3 border-t border-slate-100 pt-2 text-xs">
					<a href="/dashboard/products" class="font-medium text-slate-900 hover:underline">
						Kelola produk &rarr;
					</a>
				</div>
			</div>

			<!-- Pending / Processing Orders -->
			<div class="rounded-xl border border-slate-200 bg-white p-5 shadow-xs">
				<div class="flex items-center justify-between">
					<span class="text-xs font-semibold tracking-wider text-slate-500 uppercase"
						>Perlu Diproses</span
					>
					<div
						class="flex h-9 w-9 items-center justify-center rounded-lg bg-amber-50 text-amber-600"
					>
						<Clock class="h-5 w-5" />
					</div>
				</div>
				<div class="mt-4 flex items-baseline gap-2">
					<span class="text-2xl font-bold tracking-tight text-slate-900 tabular-nums">
						{pendingOrdersCount}
					</span>
					<span class="text-xs text-slate-500">pesanan</span>
				</div>
				<div class="mt-3 border-t border-slate-100 pt-2 text-xs">
					<a href="/dashboard/orders" class="font-medium text-amber-700 hover:underline">
						Lihat antrean pesanan &rarr;
					</a>
				</div>
			</div>

			<!-- Shipped Orders -->
			<div class="rounded-xl border border-slate-200 bg-white p-5 shadow-xs">
				<div class="flex items-center justify-between">
					<span class="text-xs font-semibold tracking-wider text-slate-500 uppercase"
						>Dalam Pengiriman</span
					>
					<div
						class="flex h-9 w-9 items-center justify-center rounded-lg bg-indigo-50 text-indigo-600"
					>
						<ShoppingBag class="h-5 w-5" />
					</div>
				</div>
				<div class="mt-4 flex items-baseline gap-2">
					<span class="text-2xl font-bold tracking-tight text-slate-900 tabular-nums">
						{shippedOrdersCount}
					</span>
					<span class="text-xs text-slate-500">paket di jalan</span>
				</div>
				<div class="mt-3 border-t border-slate-100 pt-2 text-xs text-slate-500">
					Menunggu diterima pembeli
				</div>
			</div>

			<!-- Completed Revenue -->
			<div class="rounded-xl border border-slate-200 bg-white p-5 shadow-xs">
				<div class="flex items-center justify-between">
					<span class="text-xs font-semibold tracking-wider text-slate-500 uppercase"
						>Omset Pesanan Selesai</span
					>
					<div
						class="flex h-9 w-9 items-center justify-center rounded-lg bg-emerald-50 text-emerald-600"
					>
						<TrendingUp class="h-5 w-5" />
					</div>
				</div>
				<div class="mt-4 flex items-baseline gap-2">
					<span class="text-2xl font-bold tracking-tight text-slate-900 tabular-nums">
						{formatIDR(totalCompletedRevenue)}
					</span>
				</div>
				<div class="mt-3 border-t border-slate-100 pt-2 text-xs text-slate-500">
					Total bersih transaksi tuntas
				</div>
			</div>
		</div>

		<!-- Low Stock Alert Banner (If Any) -->
		{#if lowStockProducts.length > 0}
			<div class="mt-8 rounded-xl border border-amber-200 bg-amber-50/70 p-5">
				<div class="flex items-center justify-between gap-4">
					<div class="flex items-center gap-3">
						<div
							class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-amber-100 text-amber-700"
						>
							<AlertTriangle class="h-5 w-5" />
						</div>
						<div>
							<h3 class="text-sm font-bold text-amber-900">Perhatian: Stok Produk Menipis</h3>
							<p class="text-xs text-amber-700">
								Ada {lowStockProducts.length} produk dengan sisa stok 5 unit atau kurang. Segera lakukan
								restok agar penjualan tetap lancar.
							</p>
						</div>
					</div>
					<Button variant="secondary" size="sm" href="/dashboard/products">Lihat Semua Stok</Button>
				</div>

				<div class="mt-4 grid grid-cols-1 gap-2.5 sm:grid-cols-2 lg:grid-cols-3">
					{#each lowStockProducts as prod (prod.id)}
						<div
							class="flex items-center justify-between rounded-lg border border-amber-200/80 bg-white px-3 py-2 text-xs"
						>
							<div class="min-w-0 pr-2">
								<p class="truncate font-semibold text-slate-900">{prod.name}</p>
								<p class="text-slate-500">
									Sisa: <strong class={prod.stock === 0 ? 'text-red-600' : 'text-amber-700'}
										>{prod.stock} unit</strong
									>
								</p>
							</div>
							<button
								type="button"
								onclick={() => openQuickStock(prod)}
								class="shrink-0 rounded-md bg-slate-100 px-2 py-1 font-semibold text-slate-700 transition-colors hover:bg-slate-200"
							>
								Atur Stok
							</button>
						</div>
					{/each}
				</div>
			</div>
		{/if}

		<!-- Two-Column Section: Recent Orders & Quick Actions -->
		<div class="mt-8 grid grid-cols-1 gap-8 lg:grid-cols-3">
			<!-- Recent Orders Table (2 Cols) -->
			<div class="rounded-xl border border-slate-200 bg-white shadow-xs lg:col-span-2">
				<div class="flex items-center justify-between border-b border-slate-200 px-6 py-4">
					<div>
						<h2 class="text-base font-bold text-slate-900">Pesanan Masuk Terbaru</h2>
						<p class="text-xs text-slate-500">
							Daftar transaksi yang perlu dipantau atau ditindaklanjuti
						</p>
					</div>
					<a
						href="/dashboard/orders"
						class="inline-flex items-center gap-1 text-xs font-semibold text-slate-900 hover:underline"
					>
						Semua Pesanan
						<ArrowRight class="h-3.5 w-3.5" />
					</a>
				</div>

				{#if recentOrders.length === 0}
					<div class="px-6 py-12 text-center text-sm text-slate-500">
						Belum ada pesanan masuk saat ini.
					</div>
				{:else}
					<div class="divide-y divide-slate-100">
						{#each recentOrders as order (order.id)}
							<div
								class="flex flex-col gap-3 px-6 py-4 transition-colors hover:bg-slate-50/50 sm:flex-row sm:items-center sm:justify-between"
							>
								<div class="space-y-1">
									<div class="flex items-center gap-2">
										<span class="font-mono text-xs font-bold text-slate-900">{order.id}</span>
										<Badge status={order.status} />
									</div>
									<p class="text-xs text-slate-500">
										Pembeli: <span class="font-medium text-slate-700"
											>{order.buyer?.name ||
												order.shippingAddress?.recipient_name ||
												'Pembeli'}</span
										>
										{#if order.created_at || order.createdAt}
											· {formatDate(order.created_at || order.createdAt || '')}
										{/if}
									</p>
									<p class="text-xs text-slate-600">
										{order.items?.length || 1} item ·
										<span class="font-bold text-slate-900 tabular-nums"
											>{formatIDR(order.total_amount || order.totalAmount || 0)}</span
										>
									</p>
								</div>

								<div class="flex items-center gap-2">
									<Button variant="secondary" size="sm" onclick={() => openOrderAction(order)}>
										Tindak Lanjuti
									</Button>
								</div>
							</div>
						{/each}
					</div>
				{/if}
			</div>

			<!-- Quick Actions & Guidelines (1 Col) -->
			<div class="space-y-6">
				<!-- Quick Actions Card -->
				<div class="rounded-xl border border-slate-200 bg-white p-5 shadow-xs">
					<h3 class="text-sm font-bold text-slate-900">Aksi Cepat</h3>
					<div class="mt-4 space-y-2.5">
						<button
							type="button"
							onclick={openAddProduct}
							class="flex w-full items-center justify-between rounded-lg border border-slate-200 p-3 text-left transition-all hover:border-slate-400 hover:bg-slate-50"
						>
							<div class="flex items-center gap-2.5">
								<PlusCircle class="h-4 w-4 text-amber-600" />
								<div>
									<p class="text-xs font-semibold text-slate-900">Tambah Produk Baru</p>
									<p class="text-[11px] text-slate-500">Input judul, deskripsi, harga & stok</p>
								</div>
							</div>
							<ArrowRight class="h-4 w-4 text-slate-400" />
						</button>

						<a
							href="/dashboard/orders?status=pending"
							class="flex w-full items-center justify-between rounded-lg border border-slate-200 p-3 text-left transition-all hover:border-slate-400 hover:bg-slate-50"
						>
							<div class="flex items-center gap-2.5">
								<Clock class="h-4 w-4 text-amber-600" />
								<div>
									<p class="text-xs font-semibold text-slate-900">Konfirmasi Pesanan Baru</p>
									<p class="text-[11px] text-slate-500">Terima pesanan yang menunggu</p>
								</div>
							</div>
							<ArrowRight class="h-4 w-4 text-slate-400" />
						</a>

						<a
							href="/products"
							class="flex w-full items-center justify-between rounded-lg border border-slate-200 p-3 text-left transition-all hover:border-slate-400 hover:bg-slate-50"
						>
							<div class="flex items-center gap-2.5">
								<SlidersHorizontal class="h-4 w-4 text-slate-600" />
								<div>
									<p class="text-xs font-semibold text-slate-900">Katalog Pembeli</p>
									<p class="text-[11px] text-slate-500">Lihat tampilan etalase pembeli</p>
								</div>
							</div>
							<ArrowRight class="h-4 w-4 text-slate-400" />
						</a>
					</div>
				</div>

				<!-- Workflow Helper Note -->
				<div class="rounded-xl border border-slate-200 bg-white p-5 shadow-xs">
					<h3 class="text-xs font-bold tracking-wider text-slate-500 uppercase">
						SOP Pemrosesan Pesanan
					</h3>
					<ul class="mt-3 space-y-2 text-xs text-slate-600">
						<li class="flex items-start gap-2">
							<CheckCircle2 class="mt-0.5 h-3.5 w-3.5 shrink-0 text-emerald-600" />
							<span
								><strong>Pending:</strong> Pesanan masuk. Periksa stok fisik lalu konfirmasi pesanan.</span
							>
						</li>
						<li class="flex items-start gap-2">
							<CheckCircle2 class="mt-0.5 h-3.5 w-3.5 shrink-0 text-emerald-600" />
							<span
								><strong>Processing:</strong> Kemas barang dengan aman (bubble wrap / kardus).</span
							>
						</li>
						<li class="flex items-start gap-2">
							<CheckCircle2 class="mt-0.5 h-3.5 w-3.5 shrink-0 text-emerald-600" />
							<span
								><strong>Shipped:</strong> Serahkan ke pihak ekspedisi dan perbarui status ke dikirim.</span
							>
						</li>
					</ul>
				</div>
			</div>
		</div>
	{/if}
</div>

<!-- Modals -->
<ProductFormModal
	bind:open={isProductModalOpen}
	product={selectedProduct}
	onsave={handleProductSaved}
/>

<QuickStockModal bind:open={isQuickStockOpen} product={stockProduct} onsaved={handleStockSaved} />

<OrderStatusModal
	bind:open={isOrderStatusOpen}
	order={selectedOrder}
	onupdated={handleOrderUpdated}
/>
