<script lang="ts">
	import type { Product } from '@ecommerce/shared-types';
	import { sellerApi } from '../../api/seller';
	import { toastStore } from '../../stores/toast.svelte';
	import Modal from '../ui/Modal.svelte';
	import Button from '../ui/Button.svelte';
	import { Minus, Plus, AlertCircle } from 'lucide-svelte';

	interface Props {
		open: boolean;
		product?: Product | null;
		onsaved?: (updatedProduct: Product) => void;
		onclose?: () => void;
	}

	let { open = $bindable(false), product = null, onsaved, onclose }: Props = $props();

	let stock = $state<number>(0);
	let submitting = $state(false);

	$effect(() => {
		if (open && product) {
			stock = product.stock ?? 0;
		}
	});

	function adjustStock(delta: number) {
		const next = stock + delta;
		if (next >= 0) {
			stock = next;
		}
	}

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		if (!product?.id) return;

		submitting = true;
		try {
			const res = await sellerApi.updateStock(product.id, stock);
			toastStore.success(`Stok produk berhasil diperbarui menjadi ${stock}`);
			open = false;
			if (onsaved) onsaved(res);
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : 'Gagal memperbarui stok';
			toastStore.error(msg);
		} finally {
			submitting = false;
		}
	}
</script>

<Modal bind:open title="Atur Ketersediaan Stok" maxWidth="sm" {onclose}>
	{#if product}
		<form id="quick-stock-form" onsubmit={handleSubmit} class="space-y-4">
			<div class="rounded-lg border border-slate-200 bg-slate-50 p-3">
				<p class="text-xs font-medium text-slate-500">Produk</p>
				<p class="line-clamp-1 text-sm font-semibold text-slate-900">{product.name}</p>
				<p class="mt-1 text-xs text-slate-500">
					Stok saat ini: <span class="font-bold text-slate-800">{product.stock} unit</span>
				</p>
			</div>

			<div>
				<label for="stock-input" class="mb-2 block text-xs font-semibold text-slate-800">
					Jumlah Stok Baru
				</label>

				<div class="flex items-center gap-2">
					<button
						type="button"
						onclick={() => adjustStock(-10)}
						class="flex h-9 w-10 items-center justify-center rounded-lg border border-slate-300 bg-white text-xs font-semibold text-slate-700 hover:bg-slate-100"
						title="-10 unit"
					>
						-10
					</button>
					<button
						type="button"
						onclick={() => adjustStock(-1)}
						class="flex h-9 w-9 items-center justify-center rounded-lg border border-slate-300 bg-white text-slate-700 hover:bg-slate-100"
						title="-1 unit"
					>
						<Minus class="h-4 w-4" />
					</button>

					<input
						id="stock-input"
						type="number"
						min="0"
						step="1"
						bind:value={stock}
						class="h-9 w-full rounded-lg border border-slate-300 px-3 text-center text-base font-semibold text-slate-900 tabular-nums focus:border-slate-900 focus:ring-1 focus:ring-slate-900 focus:outline-hidden"
					/>

					<button
						type="button"
						onclick={() => adjustStock(1)}
						class="flex h-9 w-9 items-center justify-center rounded-lg border border-slate-300 bg-white text-slate-700 hover:bg-slate-100"
						title="+1 unit"
					>
						<Plus class="h-4 w-4" />
					</button>
					<button
						type="button"
						onclick={() => adjustStock(10)}
						class="flex h-9 w-10 items-center justify-center rounded-lg border border-slate-300 bg-white text-xs font-semibold text-slate-700 hover:bg-slate-100"
						title="+10 unit"
					>
						+10
					</button>
				</div>
			</div>

			{#if stock === 0}
				<div
					class="flex items-start gap-2 rounded-lg border border-amber-200 bg-amber-50 p-3 text-xs text-amber-800"
				>
					<AlertCircle class="mt-0.5 h-4 w-4 shrink-0 text-amber-600" />
					<span>Stok 0 akan menandai produk sebagai <strong>Habis</strong> di katalog pembeli.</span
					>
				</div>
			{/if}
		</form>
	{/if}

	{#snippet footer()}
		<div class="flex w-full items-center justify-end gap-2.5">
			<Button
				variant="secondary"
				onclick={() => {
					open = false;
					if (onclose) onclose();
				}}
				disabled={submitting}
			>
				Batal
			</Button>
			<Button type="submit" form="quick-stock-form" variant="primary" loading={submitting}>
				Simpan Stok
			</Button>
		</div>
	{/snippet}
</Modal>
