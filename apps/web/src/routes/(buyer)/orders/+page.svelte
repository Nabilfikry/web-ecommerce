<script lang="ts">
	import { onMount } from 'svelte';
	import type { Order } from '@ecommerce/shared-types';
	import { ordersApi } from '$lib/api/orders';
	import OrderCard from '$lib/components/order/OrderCard.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import { PackageSearch, ShoppingBag } from 'lucide-svelte';

	let orders = $state<Order[]>([]);
	let loading = $state(true);
	let activeStatus = $state('all');

	const tabs = [
		{ id: 'all', label: 'Semua Pesanan' },
		{ id: 'pending', label: 'Menunggu' },
		{ id: 'processing', label: 'Diproses' },
		{ id: 'shipped', label: 'Dikirim' },
		{ id: 'completed', label: 'Selesai' },
		{ id: 'cancelled', label: 'Dibatalkan' }
	];

	async function loadOrders() {
		loading = true;
		try {
			const res = await ordersApi.getOrders({
				status: activeStatus === 'all' ? undefined : activeStatus
			});
			orders = res.orders;
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : 'Gagal memuat daftar pesanan';
			toast.error(msg);
		} finally {
			loading = false;
		}
	}

	function handleTabChange(statusId: string) {
		activeStatus = statusId;
		loadOrders();
	}

	onMount(() => {
		loadOrders();
	});
</script>

<svelte:head>
	<title>Pesanan Saya — SimpleShop</title>
</svelte:head>

<div class="mx-auto flex max-w-7xl flex-1 flex-col gap-6 px-4 py-8 sm:px-6 lg:px-8">
	<!-- Page Header -->
	<div class="flex flex-col gap-1 pb-2">
		<h1 class="text-2xl font-bold tracking-tight text-slate-900 sm:text-3xl">
			Riwayat Pesanan Saya
		</h1>
		<p class="text-xs text-slate-500 sm:text-sm">
			Pantau status transaksi, konfirmasi penerimaan barang, atau lihat invoice pembelian.
		</p>
	</div>

	<!-- Status Tabs -->
	<div class="flex items-center gap-2 overflow-x-auto border-b border-slate-200 pb-1">
		{#each tabs as tab (tab.id)}
			<button
				type="button"
				onclick={() => handleTabChange(tab.id)}
				class="cursor-pointer rounded-lg px-3.5 py-2 text-xs font-semibold whitespace-nowrap transition-colors {activeStatus ===
				tab.id
					? 'bg-slate-900 text-white shadow-xs'
					: 'text-slate-600 hover:bg-slate-100 hover:text-slate-900'}"
			>
				{tab.label}
			</button>
		{/each}
	</div>

	<!-- Orders List -->
	{#if loading}
		<div class="space-y-4">
			{#each [1, 2, 3] as idx (idx)}
				<div class="animate-pulse space-y-4 rounded-xl border border-slate-200 bg-white p-6">
					<div class="h-4 w-1/4 rounded bg-slate-200"></div>
					<div class="h-5 w-3/4 rounded bg-slate-200"></div>
					<div class="h-4 w-1/3 rounded bg-slate-200"></div>
				</div>
			{/each}
		</div>
	{:else if orders.length === 0}
		<div
			class="my-4 flex flex-col items-center gap-3 rounded-xl border border-slate-200 bg-white py-20 text-center"
		>
			<div
				class="flex h-14 w-14 items-center justify-center rounded-full bg-slate-100 text-slate-400"
			>
				<PackageSearch class="h-7 w-7 stroke-[1.5]" />
			</div>
			<h3 class="text-base font-bold text-slate-900">Belum Ada Pesanan</h3>
			<p class="max-w-sm text-xs text-slate-500">
				{activeStatus === 'all'
					? 'Anda belum pernah melakukan checkout pesanan di SimpleShop.'
					: 'Tidak ada pesanan dengan filter status yang dipilih.'}
			</p>
			<a
				href="/products"
				class="mt-2 inline-flex items-center gap-2 rounded-lg bg-slate-900 px-4 py-2 text-xs font-semibold text-white transition-colors hover:bg-slate-800"
			>
				<ShoppingBag class="h-3.5 w-3.5" />
				<span>Mulai Belanja</span>
			</a>
		</div>
	{:else}
		<div class="space-y-4">
			{#each orders as order (order.id)}
				<OrderCard {order} />
			{/each}
		</div>
	{/if}
</div>
