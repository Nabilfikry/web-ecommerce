<script lang="ts">
	import type { CartItem } from '@ecommerce/shared-types';
	import { formatIDR } from '../../utils/format';
	import { Minus, Plus, Trash2, Image as ImageIcon } from 'lucide-svelte';

	interface Props {
		item: CartItem;
		onUpdateQuantity: (itemId: string, newQty: number) => void;
		onRemove: (itemId: string) => void;
	}

	let { item, onUpdateQuantity, onRemove }: Props = $props();

	let imgError = $state(false);

	function decrement() {
		if (item.quantity > 1) {
			onUpdateQuantity(item.id, item.quantity - 1);
		}
	}

	function increment() {
		if (item.quantity < item.product.stock) {
			onUpdateQuantity(item.id, item.quantity + 1);
		}
	}
</script>

<div
	class="flex flex-col items-start justify-between gap-4 rounded-xl border border-slate-200 bg-white p-4 sm:flex-row sm:items-center"
>
	<!-- Product Info & Thumbnail -->
	<div class="flex min-w-0 flex-1 items-center gap-3.5">
		<a
			href="/products/{item.product.id}"
			class="relative block h-16 w-16 shrink-0 overflow-hidden rounded-lg border border-slate-200 bg-slate-100"
		>
			{#if item.product.image_url && !imgError}
				<img
					src={item.product.image_url}
					alt={item.product.name}
					onerror={() => (imgError = true)}
					class="h-full w-full object-cover"
				/>
			{:else}
				<div class="flex h-full w-full items-center justify-center text-slate-400">
					<ImageIcon class="h-6 w-6 stroke-[1.5]" />
				</div>
			{/if}
		</a>

		<div class="min-w-0 flex-1 space-y-0.5">
			<span class="block text-[11px] font-medium text-slate-500">
				{item.product.seller?.name || 'Toko Mitra'}
			</span>
			<a
				href="/products/{item.product.id}"
				class="block truncate text-sm font-semibold text-slate-900 transition-colors hover:text-slate-700"
				title={item.product.name}
			>
				{item.product.name}
			</a>
			<span class="block text-xs font-bold text-slate-900 tabular-nums">
				{formatIDR(item.product.price)}
			</span>
		</div>
	</div>

	<!-- Quantity & Actions -->
	<div
		class="flex w-full items-center justify-between gap-6 border-t border-slate-100 pt-2 sm:w-auto sm:justify-end sm:border-t-0 sm:pt-0"
	>
		<!-- Quantity Stepper -->
		<div class="flex items-center overflow-hidden rounded-lg border border-slate-300 bg-white">
			<button
				type="button"
				onclick={decrement}
				disabled={item.quantity <= 1}
				class="flex h-8 w-8 cursor-pointer items-center justify-center text-slate-600 transition-colors hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-30"
				aria-label="Kurangi kuantitas"
			>
				<Minus class="h-3 w-3" />
			</button>

			<span class="w-9 text-center text-xs font-semibold text-slate-900 tabular-nums">
				{item.quantity}
			</span>

			<button
				type="button"
				onclick={increment}
				disabled={item.quantity >= item.product.stock}
				class="flex h-8 w-8 cursor-pointer items-center justify-center text-slate-600 transition-colors hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-30"
				aria-label="Tambah kuantitas"
			>
				<Plus class="h-3 w-3" />
			</button>
		</div>

		<!-- Subtotal -->
		<div class="min-w-[100px] text-right">
			<span class="block text-[10px] text-slate-400">Subtotal</span>
			<span class="text-sm font-bold text-slate-900 tabular-nums">
				{formatIDR(item.subtotal)}
			</span>
		</div>

		<!-- Remove Button -->
		<button
			type="button"
			onclick={() => onRemove(item.id)}
			class="cursor-pointer rounded-lg p-2 text-slate-400 transition-colors hover:bg-red-50 hover:text-red-600"
			aria-label="Hapus dari keranjang"
			title="Hapus item"
		>
			<Trash2 class="h-4 w-4" />
		</button>
	</div>
</div>
