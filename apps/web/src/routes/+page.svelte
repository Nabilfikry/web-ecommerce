<script lang="ts">
	import { onMount } from 'svelte';
	import type { Product } from '@ecommerce/shared-types';
	import { productsApi } from '$lib/api/products';
	import ProductGrid from '$lib/components/product/ProductGrid.svelte';
	import { ArrowRight, ShieldCheck, Zap, Store, ShoppingBag } from 'lucide-svelte';
	import { toast } from '$lib/stores/toast.svelte';

	let featuredProducts = $state<Product[]>([]);
	let loading = $state(true);

	onMount(async () => {
		try {
			const res = await productsApi.getProducts({ limit: 8, sort: 'latest' });
			featuredProducts = res.products;
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : 'Gagal memuat produk';
			toast.error(msg);
		} finally {
			loading = false;
		}
	});

	function handleAddToCart(product: Product) {
		toast.success(`"${product.name}" ditambahkan ke keranjang belanja!`);
	}
</script>

<svelte:head>
	<title>SimpleShop — Platform E-Commerce Modern & Terpercaya</title>
</svelte:head>

<!-- Hero Section -->
<section class="border-b border-slate-200 bg-white">
	<div
		class="mx-auto flex max-w-7xl flex-col items-center px-4 py-12 text-center sm:px-6 sm:py-20 lg:px-8"
	>
		<!-- Tech Badge -->
		<div
			class="mb-6 inline-flex items-center gap-2 rounded-full border border-slate-200 bg-slate-100 px-3 py-1 text-xs font-semibold text-slate-800"
		>
			<span class="h-2 w-2 animate-pulse rounded-full bg-emerald-500"></span>
			<span>SvelteKit 5 · Golang REST API · PostgreSQL Monorepo</span>
		</div>

		<!-- Hero Headline -->
		<h1
			class="max-w-4xl text-3xl leading-tight font-bold tracking-tight text-slate-900 sm:text-5xl sm:leading-none lg:text-6xl"
		>
			E-Commerce Modern dengan Transaksi Terjamin.
		</h1>

		<!-- Subtitle -->
		<p class="mt-5 max-w-2xl text-sm leading-relaxed text-slate-600 sm:text-base">
			Jelajahi beragam produk pilihan dari penjual terpercaya. Nikmati proses checkout instan,
			validasi stok atomik, dan pelacakan pesanan real-time.
		</p>

		<!-- CTAs -->
		<div class="mt-8 flex w-full flex-col items-center gap-3 sm:w-auto sm:flex-row">
			<a
				href="/products"
				class="flex w-full cursor-pointer items-center justify-center gap-2 rounded-lg bg-slate-900 px-6 py-3 text-sm font-semibold text-white shadow-xs transition-all hover:bg-slate-800 sm:w-auto"
			>
				<span>Jelajahi Katalog</span>
				<ArrowRight class="h-4 w-4" />
			</a>
			<a
				href="/register"
				class="flex w-full cursor-pointer items-center justify-center gap-2 rounded-lg border border-slate-300 bg-transparent px-6 py-3 text-sm font-semibold text-slate-800 transition-colors hover:bg-slate-50 sm:w-auto"
			>
				<Store class="h-4 w-4 text-emerald-600" />
				<span>Buka Toko Penjual</span>
			</a>
		</div>

		<!-- Value Props Checklist -->
		<div
			class="mt-12 grid w-full max-w-4xl grid-cols-1 gap-6 border-t border-slate-100 pt-8 sm:grid-cols-3"
		>
			<div class="flex items-center gap-3 text-left">
				<div
					class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-emerald-50 text-emerald-700"
				>
					<ShieldCheck class="h-5 w-5" />
				</div>
				<div>
					<h4 class="text-xs font-bold text-slate-900">ACID Transaction</h4>
					<p class="text-[11px] text-slate-500">
						Stok divalidasi dan dikunci secara atomik saat checkout
					</p>
				</div>
			</div>

			<div class="flex items-center gap-3 text-left">
				<div
					class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-blue-50 text-blue-700"
				>
					<Zap class="h-5 w-5" />
				</div>
				<div>
					<h4 class="text-xs font-bold text-slate-900">Super Cepat</h4>
					<p class="text-[11px] text-slate-500">
						Katalog responsif didukung SvelteKit 5 dan Go HTTP
					</p>
				</div>
			</div>

			<div class="flex items-center gap-3 text-left">
				<div
					class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-amber-50 text-amber-700"
				>
					<ShoppingBag class="h-5 w-5" />
				</div>
				<div>
					<h4 class="text-xs font-bold text-slate-900">Role-Based Access</h4>
					<p class="text-[11px] text-slate-500">
						Akses storefront untuk Buyer dan Dashboard untuk Seller
					</p>
				</div>
			</div>
		</div>
	</div>
</section>

<!-- Featured Catalog Section -->
<section class="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
	<div class="mb-8 flex items-center justify-between">
		<div>
			<h2 class="text-xl font-bold tracking-tight text-slate-900 sm:text-2xl">Produk Terbaru</h2>
			<p class="mt-0.5 text-xs text-slate-500">
				Koleksi barang yang siap dikirim langsung oleh penjual
			</p>
		</div>

		<a
			href="/products"
			class="group flex items-center gap-1 text-xs font-semibold text-slate-800 hover:text-slate-900"
		>
			<span>Lihat Semua</span>
			<ArrowRight class="h-3.5 w-3.5 transition-transform group-hover:translate-x-0.5" />
		</a>
	</div>

	<ProductGrid products={featuredProducts} {loading} onAddToCart={handleAddToCart} />
</section>
