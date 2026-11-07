<script lang="ts">
	import { page } from '$app/state';
	import { LayoutDashboard, Package, ShoppingBag, Plus, Store } from 'lucide-svelte';
	import Button from '../ui/Button.svelte';

	interface Props {
		onAddProduct?: () => void;
	}

	let { onAddProduct }: Props = $props();

	const currentPath = $derived(page.url.pathname);

	const tabs = [
		{
			label: 'Ringkasan',
			href: '/dashboard',
			icon: LayoutDashboard,
			exact: true
		},
		{
			label: 'Produk Saya',
			href: '/dashboard/products',
			icon: Package,
			exact: false
		},
		{
			label: 'Pesanan Masuk',
			href: '/dashboard/orders',
			icon: ShoppingBag,
			exact: false
		}
	];

	function isTabActive(tabHref: string, exact: boolean) {
		if (exact) {
			return currentPath === tabHref;
		}
		return currentPath.startsWith(tabHref);
	}
</script>

<div class="mb-8 border-b border-slate-200 bg-white shadow-xs">
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<!-- Header Row -->
		<div class="flex flex-col gap-4 py-6 md:flex-row md:items-center md:justify-between">
			<div class="flex items-center gap-3">
				<div
					class="flex h-11 w-11 items-center justify-center rounded-xl bg-slate-900 text-white shadow-xs"
				>
					<Store class="h-6 w-6 text-amber-400" />
				</div>
				<div>
					<div class="flex items-center gap-2">
						<h1 class="text-xl font-bold tracking-tight text-slate-900 sm:text-2xl">
							Dashboard Penjual
						</h1>
						<span
							class="inline-flex items-center rounded-full bg-amber-50 px-2 py-0.5 text-xs font-semibold text-amber-800 ring-1 ring-amber-600/20 ring-inset"
						>
							Seller Mode
						</span>
					</div>
					<p class="mt-0.5 text-xs text-slate-500 sm:text-sm">
						Kelola katalog produk, pantau stok, dan proses pengiriman pesanan pembeli
					</p>
				</div>
			</div>

			{#if onAddProduct}
				<div class="flex items-center gap-3">
					<Button variant="primary" onclick={onAddProduct} class="shadow-sm">
						<Plus class="h-4 w-4" />
						Tambah Produk Baru
					</Button>
				</div>
			{/if}
		</div>

		<!-- Navigation Tabs -->
		<nav class="-mb-px flex space-x-6 overflow-x-auto" aria-label="Tabs Penjual">
			{#each tabs as tab (tab.href)}
				{@const active = isTabActive(tab.href, tab.exact)}
				{@const Icon = tab.icon}
				<a
					href={tab.href}
					class="inline-flex items-center gap-2 border-b-2 px-1 py-3.5 text-sm font-medium whitespace-nowrap transition-colors {active
						? 'border-slate-950 font-semibold text-slate-950'
						: 'border-transparent text-slate-500 hover:border-slate-300 hover:text-slate-700'}"
				>
					<Icon class="h-4 w-4 {active ? 'text-slate-950' : 'text-slate-400'}" />
					{tab.label}
				</a>
			{/each}
		</nav>
	</div>
</div>
