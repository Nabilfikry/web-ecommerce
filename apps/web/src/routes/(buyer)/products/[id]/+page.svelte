<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import type { Product } from '@ecommerce/shared-types';
	import { productsApi } from '$lib/api/products';
	import { formatIDR } from '$lib/utils/format';
	import { toast } from '$lib/stores/toast.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import {
		ShoppingCart,
		Store,
		ShieldCheck,
		Truck,
		Minus,
		Plus,
		ArrowLeft,
		Image as ImageIcon
	} from 'lucide-svelte';

	let product = $state<Product | null>(null);
	let loading = $state(true);
	let quantity = $state(1);
	let imgError = $state(false);

	const productId = $derived(page.params.id);

	onMount(async () => {
		try {
			if (productId) {
				product = await productsApi.getProduct(productId);
			}
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : 'Produk tidak ditemukan';
			toast.error(msg);
		} finally {
			loading = false;
		}
	});

	function incrementQuantity() {
		if (!product) return;
		if (quantity < product.stock) {
			quantity += 1;
		} else {
			toast.warning(`Maksimal pembelian ${product.stock} barang`);
		}
	}

	function decrementQuantity() {
		if (quantity > 1) {
			quantity -= 1;
		}
	}

	function handleAddToCart() {
		if (!product) return;
		toast.success(`Berhasil menambahkan ${quantity}x "${product.name}" ke keranjang!`);
	}

	function handleBuyNow() {
		if (!product) return;
		handleAddToCart();
		goto('/checkout');
	}
</script>

<svelte:head>
	<title>{product ? `${product.name} — SimpleShop` : 'Detail Produk — SimpleShop'}</title>
</svelte:head>

<div class="mx-auto flex max-w-7xl flex-1 flex-col gap-6 px-4 py-8 sm:px-6 lg:px-8">
	<!-- Breadcrumb / Back Navigation -->
	<nav class="flex items-center gap-2 text-xs text-slate-500">
		<a href="/products" class="flex items-center gap-1 transition-colors hover:text-slate-900">
			<ArrowLeft class="h-3.5 w-3.5" />
			<span>Kembali ke Katalog</span>
		</a>
		<span>/</span>
		<span class="max-w-xs truncate font-medium text-slate-900">{product?.name || 'Detail'}</span>
	</nav>

	{#if loading}
		<!-- Loading Skeleton -->
		<div
			class="grid animate-pulse grid-cols-1 gap-8 rounded-xl border border-slate-200 bg-white p-6 md:grid-cols-2 lg:gap-12"
		>
			<div class="aspect-square w-full rounded-lg bg-slate-200"></div>
			<div class="space-y-4">
				<div class="h-4 w-1/4 rounded bg-slate-200"></div>
				<div class="h-8 w-3/4 rounded bg-slate-200"></div>
				<div class="h-6 w-1/3 rounded bg-slate-200"></div>
				<div class="h-20 w-full rounded bg-slate-200"></div>
			</div>
		</div>
	{:else if !product}
		<!-- Not Found State -->
		<div class="rounded-xl border border-slate-200 bg-white p-6 py-16 text-center">
			<h2 class="text-lg font-bold text-slate-900">Produk Tidak Ditemukan</h2>
			<p class="mt-1 text-xs text-slate-500">
				Produk mungkin telah dihapus atau tidak tersedia saat ini.
			</p>
			<a
				href="/products"
				class="mt-4 inline-block rounded-lg bg-slate-900 px-4 py-2 text-xs font-semibold text-white hover:bg-slate-800"
			>
				Kembali ke Katalog
			</a>
		</div>
	{:else}
		<!-- Product Detail Content -->
		<div
			class="grid grid-cols-1 gap-8 rounded-xl border border-slate-200 bg-white p-6 shadow-xs sm:p-8 md:grid-cols-2 lg:gap-12"
		>
			<!-- Product Image Section -->
			<div class="flex flex-col gap-3">
				<div
					class="relative aspect-square w-full overflow-hidden rounded-xl border border-slate-200 bg-slate-100"
				>
					{#if product.image_url && !imgError}
						<img
							src={product.image_url}
							alt={product.name}
							onerror={() => (imgError = true)}
							class="h-full w-full object-cover"
						/>
					{:else}
						<div
							class="flex h-full w-full flex-col items-center justify-center bg-slate-50 text-slate-400"
						>
							<ImageIcon class="h-16 w-16 stroke-[1.5]" />
							<span class="mt-2 text-xs font-medium">Gambar Tidak Tersedia</span>
						</div>
					{/if}

					<!-- Floating Stock Badge -->
					{#if product.stock <= 0}
						<div
							class="absolute top-3 left-3 rounded bg-red-600 px-3 py-1 text-xs font-bold text-white uppercase shadow-sm"
						>
							Stok Habis
						</div>
					{:else if product.stock <= 5}
						<div
							class="absolute top-3 left-3 rounded bg-amber-500 px-3 py-1 text-xs font-bold text-white uppercase shadow-sm"
						>
							Tersisa {product.stock} Unit
						</div>
					{/if}
				</div>
			</div>

			<!-- Product Info & Actions -->
			<div class="flex flex-col justify-between gap-6">
				<div class="space-y-4">
					<!-- Seller Info Pill -->
					<div
						class="inline-flex items-center gap-2 rounded-full border border-slate-200 bg-slate-100 px-3 py-1 text-xs font-medium text-slate-700"
					>
						<Store class="h-3.5 w-3.5 text-emerald-600" />
						<span>Penjual: <strong>{product.seller?.name || 'Toko Mitra'}</strong></span>
					</div>

					<!-- Title -->
					<h1 class="text-2xl leading-snug font-bold tracking-tight text-slate-900 sm:text-3xl">
						{product.name}
					</h1>

					<!-- Price & Stock -->
					<div class="flex items-baseline justify-between border-y border-slate-100 pt-2 pb-4">
						<div>
							<span class="mb-0.5 block text-xs text-slate-400">Harga Produk</span>
							<span class="text-3xl font-extrabold tracking-tight text-slate-900 tabular-nums">
								{formatIDR(product.price)}
							</span>
						</div>

						<div class="text-right">
							<span class="mb-0.5 block text-xs text-slate-400">Ketersediaan</span>
							{#if product.stock > 0}
								<span
									class="rounded-full bg-emerald-50 px-2.5 py-1 text-xs font-semibold text-emerald-600"
								>
									Tersedia: {product.stock} barang
								</span>
							{:else}
								<span class="rounded-full bg-red-50 px-2.5 py-1 text-xs font-semibold text-red-600">
									Stok Habis
								</span>
							{/if}
						</div>
					</div>

					<!-- Description -->
					<div class="space-y-1.5">
						<h3 class="text-xs font-bold tracking-wider text-slate-900 uppercase">
							Deskripsi Produk
						</h3>
						<p class="text-sm leading-relaxed whitespace-pre-line text-slate-600">
							{product.description ||
								'Penjual tidak menyertakan deskripsi khusus untuk produk ini.'}
						</p>
					</div>
				</div>

				<!-- Purchase Controls -->
				<div class="space-y-4 border-t border-slate-100 pt-6">
					<!-- Quantity Stepper -->
					<div class="flex items-center justify-between">
						<span class="text-xs font-semibold text-slate-800">Jumlah Pembelian</span>

						<div
							class="flex items-center overflow-hidden rounded-lg border border-slate-300 bg-white"
						>
							<button
								type="button"
								onclick={decrementQuantity}
								disabled={quantity <= 1 || product.stock <= 0}
								class="flex h-10 w-10 cursor-pointer items-center justify-center text-slate-600 transition-colors hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-30"
								aria-label="Kurangi kuantitas"
							>
								<Minus class="h-3.5 w-3.5" />
							</button>

							<span class="w-12 text-center text-sm font-semibold text-slate-900 tabular-nums">
								{quantity}
							</span>

							<button
								type="button"
								onclick={incrementQuantity}
								disabled={quantity >= product.stock || product.stock <= 0}
								class="flex h-10 w-10 cursor-pointer items-center justify-center text-slate-600 transition-colors hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-30"
								aria-label="Tambah kuantitas"
							>
								<Plus class="h-3.5 w-3.5" />
							</button>
						</div>
					</div>

					<!-- Subtotal Preview -->
					<div class="flex items-center justify-between pt-1 text-xs text-slate-500">
						<span>Subtotal ({quantity} barang)</span>
						<span class="text-sm font-bold text-slate-900 tabular-nums">
							{formatIDR(product.price * quantity)}
						</span>
					</div>

					<!-- CTAs -->
					<div class="grid grid-cols-1 gap-3 pt-2 sm:grid-cols-2">
						<Button
							type="button"
							variant="outline"
							disabled={product.stock <= 0}
							onclick={handleAddToCart}
							class="flex w-full items-center justify-center gap-2"
						>
							<ShoppingCart class="h-4 w-4" />
							<span>Tambah ke Keranjang</span>
						</Button>

						<Button
							type="button"
							variant="primary"
							disabled={product.stock <= 0}
							onclick={handleBuyNow}
							class="w-full"
						>
							<span>Beli Sekarang</span>
						</Button>
					</div>

					<!-- Trust / Guarantee Callouts -->
					<div
						class="grid grid-cols-2 gap-3 border-t border-slate-100 pt-4 text-[11px] text-slate-500"
					>
						<div class="flex items-center gap-2">
							<ShieldCheck class="h-4 w-4 shrink-0 text-emerald-600" />
							<span>Validasi stok terjamin sebelum checkout</span>
						</div>
						<div class="flex items-center gap-2">
							<Truck class="h-4 w-4 shrink-0 text-blue-600" />
							<span>Diproses langsung oleh toko penjual</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>
