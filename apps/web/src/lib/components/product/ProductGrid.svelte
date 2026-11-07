<script lang="ts">
	import type { Product } from '@ecommerce/shared-types';
	import ProductCard from './ProductCard.svelte';
	import { PackageSearch } from 'lucide-svelte';

	interface Props {
		products: Product[];
		loading?: boolean;
		onAddToCart?: (product: Product) => void;
		onResetFilter?: () => void;
	}

	let { products, loading = false, onAddToCart, onResetFilter }: Props = $props();
</script>

{#if loading}
	<!-- Shimmer Skeleton Grid -->
	<div class="grid grid-cols-2 gap-4 sm:gap-6 md:grid-cols-3 lg:grid-cols-4">
		{#each [1, 2, 3, 4, 5, 6, 7, 8] as skeletonId (skeletonId)}
			<div class="animate-pulse overflow-hidden rounded-xl border border-slate-200 bg-white">
				<div class="aspect-square w-full bg-slate-200"></div>
				<div class="space-y-3 p-4">
					<div class="h-3 w-1/3 rounded bg-slate-200"></div>
					<div class="h-4 w-3/4 rounded bg-slate-200"></div>
					<div class="h-5 w-1/2 rounded bg-slate-200 pt-2"></div>
				</div>
			</div>
		{/each}
	</div>
{:else if products.length === 0}
	<!-- Empty State -->
	<div
		class="mx-auto my-8 flex max-w-lg flex-col items-center gap-3 rounded-xl border border-slate-200 bg-white px-4 py-16 text-center"
	>
		<div
			class="flex h-12 w-12 items-center justify-center rounded-full bg-slate-100 text-slate-500"
		>
			<PackageSearch class="h-6 w-6 stroke-[1.5]" />
		</div>
		<h3 class="text-base font-bold text-slate-900">Produk Tidak Ditemukan</h3>
		<p class="max-w-sm text-xs text-slate-500">
			Tidak ada produk yang sesuai dengan kriteria pencarian atau filter yang dipilih.
		</p>
		{#if onResetFilter}
			<button
				type="button"
				onclick={onResetFilter}
				class="mt-2 cursor-pointer rounded-lg bg-slate-900 px-4 py-2 text-xs font-semibold text-white transition-colors hover:bg-slate-800"
			>
				Reset Filter Pencarian
			</button>
		{/if}
	</div>
{:else}
	<!-- Active Products Grid -->
	<div class="grid grid-cols-2 gap-4 sm:gap-6 md:grid-cols-3 lg:grid-cols-4">
		{#each products as product (product.id)}
			<ProductCard {product} {onAddToCart} />
		{/each}
	</div>
{/if}
