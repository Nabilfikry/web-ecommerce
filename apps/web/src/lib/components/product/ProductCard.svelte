<script lang="ts">
	import type { Product } from '@ecommerce/shared-types';
	import { formatIDR } from '../../utils/format';
	import { ShoppingCart, Star, MapPin, ShieldCheck, Image as ImageIcon } from 'lucide-svelte';

	interface Props {
		product: Product;
		onAddToCart?: (product: Product) => void;
	}

	let { product, onAddToCart }: Props = $props();

	let imgError = $state(false);

	const isOutOfStock = $derived(product.stock <= 0);
	const isLowStock = $derived(product.stock > 0 && product.stock <= 5);

	// Generate deterministic rating and sold count based on product id
	const rating = $derived(4.8 + ((product.price % 3) * 0.1));
	const soldCount = $derived(50 + (product.price % 400));
	const hasDiscount = $derived(product.price > 200000);
	const originalPrice = $derived(Math.round(product.price * 1.15));
</script>

<div
	class="group relative flex flex-col overflow-hidden rounded-xl border border-[#E5E7E9] bg-white transition-all duration-200 hover:-translate-y-0.5 hover:shadow-lg"
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

		<!-- Badges on Image (Top Left) -->
		<div class="absolute top-2 left-2 flex flex-col gap-1">
			{#if isOutOfStock}
				<span
					class="rounded bg-[#EE4D2D] px-2 py-0.5 text-[10px] font-bold tracking-wider text-white uppercase shadow-xs"
				>
					Stok Habis
				</span>
			{:else if isLowStock}
				<span
					class="rounded bg-[#FF9800] px-2 py-0.5 text-[10px] font-bold tracking-wider text-white uppercase shadow-xs"
				>
					Sisa {product.stock}
				</span>
			{/if}
		</div>

		<!-- Bebas Ongkir Pill (Bottom Left) -->
		<div class="absolute bottom-2 left-2">
			<span class="inline-flex items-center gap-1 rounded bg-emerald-500/90 px-1.5 py-0.5 text-[9px] font-bold text-white shadow-xs backdrop-blur-xs">
				Bebas Ongkir
			</span>
		</div>
	</a>

	<!-- Product Content Info -->
	<div class="flex flex-1 flex-col justify-between gap-2 p-3">
		<div class="space-y-1">
			<!-- Product Title (2 Lines) -->
			<a
				href="/products/{product.id}"
				class="line-clamp-2 block text-xs leading-snug font-medium text-[#212121] transition-colors group-hover:text-[#00AA5B]"
				title={product.name}
			>
				{product.name}
			</a>

			<!-- Price -->
			<div class="pt-0.5">
				<div class="text-sm font-bold tracking-tight text-[#212121] tabular-nums">
					{formatIDR(product.price)}
				</div>

				{#if hasDiscount}
					<div class="flex items-center gap-1.5 text-[10px]">
						<span class="rounded bg-red-100 px-1 font-bold text-[#EE4D2D]">15%</span>
						<span class="text-[#6D7588] line-through tabular-nums">
							{formatIDR(originalPrice)}
						</span>
					</div>
				{/if}
			</div>

			<!-- Seller Badge & Location -->
			<div class="flex items-center gap-1 pt-1 text-[11px] text-[#6D7588]">
				<span class="flex items-center text-[#00AA5B]" title="Official Store">
					<ShieldCheck class="h-3 w-3 fill-[#00AA5B] text-white" />
				</span>
				<span class="truncate">{product.seller?.name || 'Jakarta Pusat'}</span>
			</div>

			<!-- Rating & Sold Count -->
			<div class="flex items-center gap-1 text-[11px] text-[#6D7588]">
				<Star class="h-3 w-3 fill-[#FFAB00] text-[#FFAB00]" />
				<span class="font-bold text-[#212121]">{rating.toFixed(1)}</span>
				<span class="text-slate-300">•</span>
				<span>Terjual {soldCount}+</span>
			</div>
		</div>

		<!-- Action CTA -->
		{#if onAddToCart}
			<div class="pt-2">
				<button
					type="button"
					onclick={() => onAddToCart(product)}
					disabled={isOutOfStock}
					class="flex w-full cursor-pointer items-center justify-center gap-1 rounded-lg border border-[#00AA5B] py-1.5 text-xs font-bold text-[#00AA5B] transition-colors hover:bg-[#F0FDF4] active:scale-98 disabled:cursor-not-allowed disabled:border-slate-200 disabled:text-slate-300"
				>
					<ShoppingCart class="h-3.5 w-3.5" />
					<span>+ Keranjang</span>
				</button>
			</div>
		{/if}
	</div>
</div>
