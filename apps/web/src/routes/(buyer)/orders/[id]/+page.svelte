<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import type { Order } from '@ecommerce/shared-types';
	import { ordersApi } from '$lib/api/orders';
	import { formatIDR, formatDate, getStatusLabel } from '$lib/utils/format';
	import { toast } from '$lib/stores/toast.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Modal from '$lib/components/ui/Modal.svelte';
	import OrderTimeline from '$lib/components/order/OrderTimeline.svelte';
	import { ArrowLeft, Store, MapPin, CheckCircle, XCircle } from 'lucide-svelte';

	let order = $state<Order | null>(null);
	let loading = $state(true);

	// Cancellation modal state
	let cancelModalOpen = $state(false);
	let cancelReason = $state('Salah memilih varian produk');
	let isCancelling = $state(false);

	// Completion state
	let isCompleting = $state(false);

	const orderId = $derived(page.params.id);

	async function loadOrder() {
		loading = true;
		try {
			if (orderId) {
				order = await ordersApi.getOrderDetail(orderId);
			}
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : 'Pesanan tidak ditemukan';
			toast.error(msg);
		} finally {
			loading = false;
		}
	}

	async function handleCancelOrder() {
		if (!order) return;
		isCancelling = true;
		try {
			const updated = await ordersApi.cancelOrder(order.id, cancelReason);
			order = updated;
			cancelModalOpen = false;
			toast.success('Pesanan berhasil dibatalkan. Stok barang telah dikembalikan.');
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : 'Gagal membatalkan pesanan';
			toast.error(msg);
		} finally {
			isCancelling = false;
		}
	}

	async function handleCompleteOrder() {
		if (!order) return;
		if (!confirm('Apakah Anda yakin barang sudah diterima dengan baik dan pesanan selesai?')) {
			return;
		}

		isCompleting = true;
		try {
			const updated = await ordersApi.completeOrder(order.id);
			order = updated;
			toast.success('Pesanan selesai! Terima kasih telah berbelanja di SimpleShop.');
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : 'Gagal menyelesaikan pesanan';
			toast.error(msg);
		} finally {
			isCompleting = false;
		}
	}

	onMount(() => {
		loadOrder();
	});
</script>

<svelte:head>
	<title>{order ? `Detail Pesanan #${order.id} — SimpleShop` : 'Detail Pesanan — SimpleShop'}</title
	>
</svelte:head>

<div class="mx-auto flex max-w-7xl flex-1 flex-col gap-6 px-4 py-8 sm:px-6 lg:px-8">
	<!-- Navigation Header -->
	<nav class="flex items-center gap-2 text-xs text-slate-500">
		<a href="/orders" class="flex items-center gap-1 transition-colors hover:text-slate-900">
			<ArrowLeft class="h-3.5 w-3.5" />
			<span>Kembali ke Pesanan Saya</span>
		</a>
		<span>/</span>
		<span class="font-medium text-slate-900">Invoice #{order?.id || ''}</span>
	</nav>

	{#if loading}
		<div class="animate-pulse space-y-4 rounded-xl border border-slate-200 bg-white p-8">
			<div class="h-6 w-1/4 rounded bg-slate-200"></div>
			<div class="h-10 w-1/2 rounded bg-slate-200"></div>
			<div class="h-40 w-full rounded bg-slate-200"></div>
		</div>
	{:else if !order}
		<div class="rounded-xl border border-slate-200 bg-white p-6 py-20 text-center">
			<h2 class="text-lg font-bold text-slate-900">Pesanan Tidak Ditemukan</h2>
			<p class="mt-1 text-xs text-slate-500">Kode pesanan tidak terdaftar pada akun Anda.</p>
			<a
				href="/orders"
				class="mt-4 inline-block rounded-lg bg-slate-900 px-4 py-2 text-xs font-semibold text-white hover:bg-slate-800"
			>
				Kembali ke Riwayat Pesanan
			</a>
		</div>
	{:else}
		<!-- Order Invoice Header Card -->
		<div
			class="flex flex-col items-start justify-between gap-4 rounded-xl border border-slate-200 bg-white p-6 shadow-xs sm:flex-row sm:items-center"
		>
			<div>
				<div class="flex items-center gap-3">
					<h1 class="text-xl font-bold tracking-tight text-slate-900 sm:text-2xl">
						Invoice #{order.id}
					</h1>
					<Badge variant={order.status}>
						{getStatusLabel(order.status)}
					</Badge>
				</div>
				<p class="mt-1 text-xs text-slate-500">
					Waktu Pemesanan: <strong>{formatDate(order.createdAt || order.created_at)}</strong>
				</p>
			</div>

			<!-- Action Buttons depending on status -->
			<div class="flex items-center gap-2">
				{#if order.status === 'pending'}
					<Button
						type="button"
						variant="outline"
						onclick={() => (cancelModalOpen = true)}
						class="flex items-center gap-1.5 border-red-200 text-red-600 hover:border-red-300 hover:bg-red-50"
					>
						<XCircle class="h-4 w-4" />
						<span>Batalkan Pesanan</span>
					</Button>
				{:else if order.status === 'shipped'}
					<Button
						type="button"
						variant="primary"
						loading={isCompleting}
						onclick={handleCompleteOrder}
						class="flex items-center gap-1.5 bg-emerald-700 hover:bg-emerald-800"
					>
						<CheckCircle class="h-4 w-4" />
						<span>Konfirmasi Barang Diterima</span>
					</Button>
				{/if}
			</div>
		</div>

		<!-- Status Timeline -->
		<OrderTimeline status={order.status} />

		<!-- Two Column Invoice Breakdown -->
		<div class="grid grid-cols-1 items-start gap-8 lg:grid-cols-3">
			<!-- Left: Items Snapshot Table (2 cols) -->
			<div
				class="space-y-4 rounded-xl border border-slate-200 bg-white p-6 shadow-xs lg:col-span-2"
			>
				<div class="flex items-center justify-between border-b border-slate-100 pb-3">
					<h3 class="text-sm font-bold tracking-wider text-slate-900 uppercase">
						Rincian Produk (Snapshot Transaksi)
					</h3>
					<span class="text-xs text-slate-500">
						{order.items?.length || 0} jenis barang
					</span>
				</div>

				<div class="divide-y divide-slate-100">
					{#each order.items || [] as item (item.id || item.productId || item.product_id)}
						<div class="flex items-center justify-between gap-4 py-3 text-xs">
							<div class="min-w-0 flex-1">
								<span class="block truncate font-semibold text-slate-900">
									{item.productName || item.product_name}
								</span>
								<span class="text-[11px] text-slate-400">
									{item.quantity} x @ {formatIDR(item.unitPrice || item.unit_price)}
								</span>
							</div>

							<span class="shrink-0 font-bold text-slate-900 tabular-nums">
								{formatIDR(item.subtotal)}
							</span>
						</div>
					{/each}
				</div>

				<!-- Financial Breakdown -->
				<div class="space-y-2 border-t border-slate-200 pt-4 text-xs text-slate-600">
					<div class="flex items-center justify-between">
						<span>Subtotal Produk</span>
						<span class="font-semibold text-slate-900 tabular-nums">
							{formatIDR(order.totalAmount || order.total_amount)}
						</span>
					</div>

					<div class="flex items-center justify-between">
						<span>Biaya Pengiriman</span>
						<span class="font-semibold text-emerald-600">Gratis (Simulasi MVP)</span>
					</div>

					<div class="flex items-baseline justify-between border-t border-slate-100 pt-3">
						<span class="text-sm font-bold text-slate-900">Total Pembayaran</span>
						<span class="text-xl font-extrabold text-slate-900 tabular-nums">
							{formatIDR(order.totalAmount || order.total_amount)}
						</span>
					</div>
				</div>
			</div>

			<!-- Right: Shipping Address & Notes (1 col) -->
			<div class="space-y-6">
				<!-- Shipping Address Box -->
				<div class="space-y-3 rounded-xl border border-slate-200 bg-white p-6 shadow-xs">
					<div
						class="flex items-center gap-2 border-b border-slate-100 pb-2 text-xs font-bold tracking-wider text-slate-800 uppercase"
					>
						<MapPin class="h-4 w-4 text-slate-500" />
						<span>Alamat Pengiriman</span>
					</div>

					{#if order.shippingAddress || order.shipping_address}
						{@const addr = order.shippingAddress || order.shipping_address}
						<div class="space-y-1 text-xs text-slate-700">
							<p class="font-bold text-slate-900">{addr.recipient_name}</p>
							<p class="text-slate-500">{addr.phone}</p>
							<p class="mt-1 leading-relaxed text-slate-600">{addr.address}</p>
							<p class="text-slate-500">{addr.city}, {addr.province || ''} {addr.postal_code}</p>
						</div>
					{:else}
						<p class="text-xs text-slate-400">Informasi alamat tidak tersedia.</p>
					{/if}

					{#if order.notes}
						<div class="border-t border-slate-100 pt-3">
							<span class="mb-0.5 block text-[11px] font-semibold text-slate-400"
								>Catatan Pembeli:</span
							>
							<p class="text-xs text-slate-600 italic">"{order.notes}"</p>
						</div>
					{/if}
				</div>

				<!-- Store Info Card -->
				<div
					class="flex items-center gap-3 rounded-xl border border-slate-200 bg-white p-5 shadow-xs"
				>
					<div
						class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-emerald-50 text-emerald-700"
					>
						<Store class="h-5 w-5" />
					</div>
					<div class="min-w-0 flex-1">
						<span class="block text-[10px] font-semibold text-slate-400 uppercase"
							>Penjual Terkait</span
						>
						<span class="block truncate text-xs font-bold text-slate-900">
							{order.seller?.name || 'Toko Mitra Penjual'}
						</span>
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>

<!-- Cancel Order Confirmation Modal -->
<Modal bind:open={cancelModalOpen} title="Konfirmasi Pembatalan Pesanan" maxWidth="md">
	<div class="space-y-4 text-xs">
		<p class="leading-relaxed text-slate-600">
			Apakah Anda yakin ingin membatalkan pesanan ini? Stok produk akan otomatis dipulihkan ke toko
			penjual.
		</p>

		<div class="space-y-1.5">
			<label for="cancel-reason" class="block font-semibold text-slate-800">
				Alasan Pembatalan:
			</label>
			<select
				id="cancel-reason"
				bind:value={cancelReason}
				class="h-10 w-full rounded-lg border border-slate-300 bg-white px-3 text-xs text-slate-900 focus:border-slate-900 focus:outline-none"
			>
				<option value="Salah memilih varian produk">Salah memilih varian produk</option>
				<option value="Ingin mengubah alamat pengiriman">Ingin mengubah alamat pengiriman</option>
				<option value="Menemukan penawaran lain">Menemukan penawaran lain</option>
				<option value="Lainnya">Lainnya</option>
			</select>
		</div>
	</div>

	{#snippet footer()}
		<Button type="button" variant="ghost" onclick={() => (cancelModalOpen = false)}>Batal</Button>
		<Button type="button" variant="destructive" loading={isCancelling} onclick={handleCancelOrder}>
			Konfirmasi Batalkan
		</Button>
	{/snippet}
</Modal>
