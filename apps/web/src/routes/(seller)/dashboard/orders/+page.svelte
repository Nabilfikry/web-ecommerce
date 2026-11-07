<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import type { Order } from '@ecommerce/shared-types';
	import { sellerApi } from '$lib/api/seller';
	import { formatIDR, formatDate } from '$lib/utils/format';
	import SellerNav from '$lib/components/seller/SellerNav.svelte';
	import OrderStatusModal from '$lib/components/seller/OrderStatusModal.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { Inbox, MapPin, Phone, PackageCheck, Truck, CheckCircle2, XCircle } from 'lucide-svelte';

	let loading = $state(true);
	let orders = $state<Order[]>([]);
	let selectedStatus = $state<string>('all');

	// Status action modal
	let isStatusModalOpen = $state(false);
	let actionTargetOrder = $state<Order | null>(null);

	const statusFilters = [
		{ id: 'all', label: 'Semua Pesanan' },
		{ id: 'pending', label: 'Menunggu Konfirmasi' },
		{ id: 'confirmed', label: 'Dikonfirmasi' },
		{ id: 'processing', label: 'Sedang Diproses' },
		{ id: 'shipped', label: 'Dikirim' },
		{ id: 'completed', label: 'Selesai' },
		{ id: 'rejected', label: 'Ditolak / Batal' }
	];

	async function loadOrders() {
		loading = true;
		try {
			const res = await sellerApi.getIncomingOrders({ limit: 100 });
			orders = res.orders;
		} catch (err: unknown) {
			console.error('Failed to load incoming orders:', err);
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		// Check URL query param for status filter
		const urlStatus = page.url.searchParams.get('status');
		if (urlStatus) {
			selectedStatus = urlStatus;
		}
		loadOrders();
	});

	const filteredOrders = $derived(
		orders.filter((o) => {
			if (selectedStatus === 'all') return true;
			if (selectedStatus === 'rejected') {
				return o.status === 'rejected' || o.status === 'cancelled';
			}
			return o.status === selectedStatus;
		})
	);

	function openStatusModal(order: Order) {
		actionTargetOrder = order;
		isStatusModalOpen = true;
	}

	function handleOrderUpdated(updated: Order) {
		orders = orders.map((o) => (o.id === updated.id ? { ...o, status: updated.status } : o));
	}
</script>

<svelte:head>
	<title>Pesanan Masuk - Seller Dashboard</title>
</svelte:head>

<SellerNav />

<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
	<!-- Status Tabs -->
	<div class="mb-6 overflow-x-auto pb-2">
		<div
			class="flex items-center gap-1.5 rounded-xl border border-slate-200 bg-white p-1 text-xs whitespace-nowrap shadow-xs"
		>
			{#each statusFilters as filter (filter.id)}
				<button
					type="button"
					onclick={() => (selectedStatus = filter.id)}
					class="rounded-lg px-3 py-1.5 font-medium transition-colors {selectedStatus === filter.id
						? 'bg-slate-900 font-semibold text-white shadow-xs'
						: 'text-slate-600 hover:bg-slate-100 hover:text-slate-900'}"
				>
					{filter.label}
				</button>
			{/each}
		</div>
	</div>

	<!-- Orders Content -->
	{#if loading}
		<div class="space-y-4">
			{#each [1, 2, 3] as item (item)}
				<div class="h-44 animate-pulse rounded-xl bg-slate-200/80"></div>
			{/each}
		</div>
	{:else if filteredOrders.length === 0}
		<div class="rounded-xl border border-slate-200 bg-white px-6 py-16 text-center shadow-xs">
			<Inbox class="mx-auto h-12 w-12 text-slate-300" />
			<h3 class="mt-4 text-base font-semibold text-slate-900">Tidak ada pesanan</h3>
			<p class="mt-1 text-xs text-slate-500">
				{selectedStatus !== 'all'
					? `Tidak ada pesanan masuk dengan status "${selectedStatus}".`
					: 'Belum ada pesanan yang masuk ke toko Anda saat ini.'}
			</p>
			{#if selectedStatus !== 'all'}
				<div class="mt-4">
					<Button variant="secondary" size="sm" onclick={() => (selectedStatus = 'all')}>
						Lihat Semua Pesanan
					</Button>
				</div>
			{/if}
		</div>
	{:else}
		<div class="space-y-4">
			{#each filteredOrders as order (order.id)}
				{@const addr = order.shippingAddress || order.shipping_address}
				<div
					class="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-xs transition-colors hover:border-slate-300"
				>
					<!-- Order Header -->
					<div
						class="flex flex-wrap items-center justify-between gap-3 border-b border-slate-100 bg-slate-50/60 px-5 py-3 text-xs"
					>
						<div class="flex items-center gap-3">
							<span class="font-mono text-sm font-bold text-slate-900">#{order.id}</span>
							<Badge status={order.status} />
						</div>

						<div class="text-slate-500">
							{#if order.created_at || order.createdAt}
								Waktu Transaksi: <span class="font-medium text-slate-700"
									>{formatDate(order.created_at || order.createdAt || '')}</span
								>
							{/if}
						</div>
					</div>

					<!-- Order Body -->
					<div class="p-5">
						<div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
							<!-- Buyer & Shipping Info (1 Col) -->
							<div
								class="space-y-2 rounded-lg border border-slate-100 bg-slate-50/50 p-3.5 text-xs"
							>
								<div class="flex items-center gap-1.5 font-semibold text-slate-900">
									<MapPin class="h-3.5 w-3.5 text-slate-500" />
									<span>Tujuan Pengiriman</span>
								</div>

								<div class="space-y-1 text-slate-600">
									<p class="font-medium text-slate-800">
										{addr?.recipient_name || order.buyer?.name || 'Penerima'}
									</p>
									{#if addr?.phone}
										<p class="flex items-center gap-1 text-slate-500">
											<Phone class="h-3 w-3" />
											{addr.phone}
										</p>
									{/if}
									<p class="leading-relaxed text-slate-700">
										{addr?.address || 'Alamat tidak disertakan'}
									</p>
									<p class="text-slate-500">
										{addr?.city || ''}
										{addr?.postal_code ? `(${addr.postal_code})` : ''}
									</p>
								</div>

								{#if order.notes}
									<div class="mt-2 border-t border-slate-200/80 pt-2 text-[11px] text-amber-800">
										<strong>Catatan Pembeli:</strong> "{order.notes}"
									</div>
								{/if}
							</div>

							<!-- Items Breakdown (2 Cols) -->
							<div class="space-y-3 lg:col-span-2">
								<div class="divide-y divide-slate-100 border-b border-slate-100 pb-3">
									{#each order.items || [] as item (item.productId || item.product_id || item.id || item.productName)}
										<div class="flex items-center justify-between py-2 text-xs">
											<div class="min-w-0 pr-3">
												<p class="truncate font-semibold text-slate-900">
													{item.productName || item.product_name || 'Produk'}
												</p>
												<p class="text-slate-500">
													{item.quantity} x {formatIDR(item.unitPrice || item.unit_price || 0)}
												</p>
											</div>
											<div class="font-bold text-slate-900 tabular-nums">
												{formatIDR(
													item.subtotal || item.quantity * (item.unitPrice || item.unit_price || 0)
												)}
											</div>
										</div>
									{/each}
								</div>

								<!-- Total & Action Row -->
								<div
									class="flex flex-col gap-3 pt-1 sm:flex-row sm:items-center sm:justify-between"
								>
									<div>
										<span class="text-xs text-slate-500">Total Pembayaran:</span>
										<p class="text-lg font-bold text-slate-900 tabular-nums">
											{formatIDR(order.total_amount || order.totalAmount || 0)}
										</p>
									</div>

									<!-- Action CTA based on status -->
									<div class="flex items-center gap-2">
										{#if order.status === 'pending'}
											<Button
												variant="primary"
												size="sm"
												class="bg-emerald-600 text-white hover:bg-emerald-700"
												onclick={() => openStatusModal(order)}
											>
												<CheckCircle2 class="h-3.5 w-3.5" />
												Terima Pesanan
											</Button>
											<Button variant="danger" size="sm" onclick={() => openStatusModal(order)}>
												<XCircle class="h-3.5 w-3.5" />
												Tolak
											</Button>
										{:else if order.status === 'confirmed'}
											<Button variant="primary" size="sm" onclick={() => openStatusModal(order)}>
												<PackageCheck class="h-3.5 w-3.5" />
												Mulai Siapkan (Process)
											</Button>
										{:else if order.status === 'processing'}
											<Button
												variant="primary"
												size="sm"
												class="bg-indigo-600 text-white hover:bg-indigo-700"
												onclick={() => openStatusModal(order)}
											>
												<Truck class="h-3.5 w-3.5" />
												Kirim Pesanan (Ship)
											</Button>
										{:else if order.status === 'shipped'}
											<div
												class="inline-flex items-center gap-1.5 rounded-lg border border-indigo-200 bg-indigo-50 px-3 py-1.5 text-xs font-semibold text-indigo-700"
											>
												<Truck class="h-3.5 w-3.5" />
												<span>Dalam Pengiriman</span>
											</div>
										{:else if order.status === 'completed'}
											<div
												class="inline-flex items-center gap-1.5 rounded-lg border border-emerald-200 bg-emerald-50 px-3 py-1.5 text-xs font-semibold text-emerald-700"
											>
												<CheckCircle2 class="h-3.5 w-3.5" />
												<span>Transaksi Tuntas</span>
											</div>
										{:else}
											<div
												class="inline-flex items-center gap-1.5 rounded-lg border border-red-200 bg-red-50 px-3 py-1.5 text-xs font-semibold text-red-700"
											>
												<XCircle class="h-3.5 w-3.5" />
												<span>Pesanan Dibatalkan / Ditolak</span>
											</div>
										{/if}
									</div>
								</div>
							</div>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<!-- Order Status Management Modal -->
<OrderStatusModal
	bind:open={isStatusModalOpen}
	order={actionTargetOrder}
	onupdated={handleOrderUpdated}
/>
