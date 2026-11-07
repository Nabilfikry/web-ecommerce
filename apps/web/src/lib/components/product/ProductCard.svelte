<script lang="ts">
	import type { Product } from '@ecommerce/shared-types';
	import { formatIDR } from '../../utils/format';
	import { ShoppingCart, Image as ImageIcon } from 'lucide-svelte';

	interface Props {
		product: Product;
		onAddToCart?: (product: Product) => void;
	}

	let { product, onAddToCart }: Props = $props();

	let imgError = $state(false);

	const isOutOfStock = $derived(product.stock <= 0);
	const isLowStock = $derived(product.stock > 0 && product.stock <= 5);
</script>

<div
	class="group relative flex flex-col overflow-hidden rounded-xl border border-slate-200 bg-white transition-all duration-200 hover:border-slate-300 hover:shadow-md"
>
	<!-- 1:1 Aspect Image Container -->
	<a
		href="/products/{product.id}"
		class="relative block aspect-square w-full overflow-hidden bg-slate-100"
	>
		{#if product.image_url && !imgError}
			<img
				src={product.image_url}
				alt={product.name}
				onerror={() => (imgError = true)}
				loading="lazy"
				class="h-full w-full object-cover transition-transform duration-300 ease-out group-hover:scale-105"
			/>
		{:else}
			<div
				class="flex h-full w-full flex-col items-center justify-center bg-slate-50 text-slate-400"
			>
				<ImageIcon class="h-10 w-10 stroke-[1.5]" />
				<span class="mt-1 text-[11px] font-medium">Gambar Produk</span>
			</div>
		{/if}

		<!-- Stock Badge Pill -->
		<div class="absolute top-2.5 left-2.5 flex flex-col gap-1">
			{#if isOutOfStock}
				<span
					class="rounded bg-red-600 px-2 py-0.5 text-[10px] font-bold tracking-wider text-white uppercase shadow-xs"
				>
					Stok Habis
				</span>
			{:else if isLowStock}
				<span
					class="rounded bg-amber-500 px-2 py-0.5 text-[10px] font-bold tracking-wider text-white uppercase shadow-xs"
				>
					Sisa {product.stock}
				</span>
			{/if}
		</div>
	</a>

	<!-- Product Content Info -->
	<div class="flex flex-1 flex-col justify-between gap-3 p-4">
		<div class="space-y-1">
			<!-- Seller Store Name -->
			<div class="flex items-center justify-between text-xs text-slate-500">
				<span class="max-w-[140px] truncate font-medium">
					{product.seller?.name || 'Toko Mitra'}
				</span>
				{#if !isOutOfStock && !isLowStock}
					<span class="text-[11px] font-medium text-emerald-600">Stok: {product.stock}</span>
				{/if}
			</div>

			<!-- Product Title -->
			<a
				href="/products/{product.id}"
				class="line-clamp-2 block text-sm leading-snug font-semibold text-slate-900 transition-colors group-hover:text-slate-700"
				title={product.name}
			>
				{product.name}
			</a>
		</div>

		<!-- Price & Action -->
		<div class="flex items-center justify-between gap-2 border-t border-slate-100 pt-2">
			<div class="flex flex-col">
				<span class="text-[10px] font-medium text-slate-400">Harga</span>
				<span class="text-base font-bold tracking-tight text-slate-900 tabular-nums">
					{formatIDR(product.price)}
				</span>
			</div>

			{#if onAddToCart}
				<button
					type="button"
					onclick={() => onAddToCart(product)}
					disabled={isOutOfStock}
					class="cursor-pointer rounded-lg bg-slate-900 p-2 text-white shadow-xs transition-all hover:bg-slate-800 active:scale-95 disabled:cursor-not-allowed disabled:bg-slate-200 disabled:text-slate-400"
					aria-label="Tambah ke keranjang"
					title={isOutOfStock ? 'Stok produk habis' : 'Tambah ke keranjang'}
				>
					<ShoppingCart class="h-4 w-4" />
				</button>
			{/if}
		</div>
	</div>
</div>
