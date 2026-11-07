<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import type { Product } from '@ecommerce/shared-types';
	import { productsApi } from '$lib/api/products';
	import ProductGrid from '$lib/components/product/ProductGrid.svelte';
	import ProductFilter from '$lib/components/product/ProductFilter.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import { cart } from '$lib/stores/cart.svelte';
	import { ChevronLeft, ChevronRight } from 'lucide-svelte';

	let products = $state<Product[]>([]);
	let loading = $state(true);
	let totalItems = $state(0);
	let totalPages = $state(1);

	// Synchronize with URL search params
	let search = $state(page.url.searchParams.get('search') || '');
	let sort = $state<'latest' | 'price_asc' | 'price_desc'>(
		(page.url.searchParams.get('sort') as 'latest' | 'price_asc' | 'price_desc') || 'latest'
	);
	let currentPage = $state(Number(page.url.searchParams.get('page')) || 1);

	async function loadProducts() {
		loading = true;
		try {
			const res = await productsApi.getProducts({
				page: currentPage,
				limit: 12,
				search: search.trim() || undefined,
				sort
			});
			products = res.products;
			totalItems = res.meta.total;
			totalPages = Math.max(1, res.meta.total_pages);
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : 'Gagal memuat katalog produk';
			toast.error(msg);
		} finally {
			loading = false;
		}
	}

	function handleFilterChange(filters: {
		search: string;
		sort: 'latest' | 'price_asc' | 'price_desc';
	}) {
		search = filters.search;
		sort = filters.sort;
		currentPage = 1;
		updateUrl();
		loadProducts();
	}

	function handleResetFilter() {
		search = '';
		sort = 'latest';
		currentPage = 1;
		updateUrl();
		loadProducts();
	}

	function handlePageChange(newPage: number) {
		if (newPage < 1 || newPage > totalPages) return;
		currentPage = newPage;
		updateUrl();
		loadProducts();
		window.scrollTo({ top: 0, behavior: 'smooth' });
	}

	function updateUrl() {
		const parts: string[] = [];
		if (search) parts.push(`search=${encodeURIComponent(search)}`);
		if (sort !== 'latest') parts.push(`sort=${encodeURIComponent(sort)}`);
		if (currentPage > 1) parts.push(`page=${currentPage}`);

		const query = parts.join('&');
		goto(`/products${query ? `?${query}` : ''}`);
	}

	onMount(() => {
		loadProducts();
	});

	async function handleAddToCart(product: Product) {
		try {
			await cart.addItem(product, 1);
			toast.success(`"${product.name}" berhasil ditambahkan ke keranjang belanja.`);
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : 'Gagal menambahkan produk ke keranjang';
			toast.error(msg);
		}
	}
</script>

<svelte:head>
	<title>Katalog Produk — SimpleShop</title>
</svelte:head>

<div class="mx-auto flex max-w-7xl flex-1 flex-col gap-6 px-4 py-8 sm:px-6 lg:px-8">
	<!-- Page Header -->
	<div class="flex flex-col gap-1">
		<h1 class="text-2xl font-bold tracking-tight text-slate-900 sm:text-3xl">Katalog Produk</h1>
		<p class="text-xs text-slate-500 sm:text-sm">
			Temukan produk kebutuhanmu langsung dari seller terpercaya.
		</p>
	</div>

	<!-- Filter Component -->
	<ProductFilter bind:search bind:sort {totalItems} onchange={handleFilterChange} />

	<!-- Product Grid -->
	<ProductGrid
		{products}
		{loading}
		onAddToCart={handleAddToCart}
		onResetFilter={handleResetFilter}
	/>

	<!-- Pagination -->
	{#if totalPages > 1 && !loading}
		<div class="mt-auto flex items-center justify-center gap-2 border-t border-slate-200 pt-6">
			<button
				type="button"
				disabled={currentPage <= 1}
				onclick={() => handlePageChange(currentPage - 1)}
				class="flex h-9 cursor-pointer items-center gap-1 rounded-lg border border-slate-300 bg-white px-3 text-xs font-medium text-slate-700 transition-colors hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-40"
			>
				<ChevronLeft class="h-3.5 w-3.5" />
				<span>Sebelumnya</span>
			</button>

			<span class="px-3 text-xs text-slate-600">
				Halaman <strong>{currentPage}</strong> dari <strong>{totalPages}</strong>
			</span>

			<button
				type="button"
				disabled={currentPage >= totalPages}
				onclick={() => handlePageChange(currentPage + 1)}
				class="flex h-9 cursor-pointer items-center gap-1 rounded-lg border border-slate-300 bg-white px-3 text-xs font-medium text-slate-700 transition-colors hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-40"
			>
				<span>Berikutnya</span>
				<ChevronRight class="h-3.5 w-3.5" />
			</button>
		</div>
	{/if}
</div>
