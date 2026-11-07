<script lang="ts">
	import { goto } from '$app/navigation';
	import { cart } from '$lib/stores/cart.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import { formatIDR } from '$lib/utils/format';
	import CartItemRow from '$lib/components/cart/CartItemRow.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { ShoppingBag, ArrowRight, Trash2, AlertTriangle, ShieldCheck } from 'lucide-svelte';

	async function handleUpdateQuantity(itemId: string, newQty: number) {
		try {
			await cart.updateQuantity(itemId, newQty);
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : 'Gagal memperbarui kuantitas barang';
			toast.error(msg);
		}
	}

	async function handleRemoveItem(itemId: string) {
		try {
			await cart.removeItem(itemId);
			toast.info('Item dihapus dari keranjang.');
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : 'Gagal menghapus item';
			toast.error(msg);
		}
	}

	async function handleClearCart() {
		if (confirm('Apakah Anda yakin ingin mengosongkan seluruh keranjang belanja?')) {
			try {
				await cart.clearCart();
				toast.info('Keranjang berhasil dikosongkan.');
			} catch (err: unknown) {
				const msg = err instanceof Error ? err.message : 'Gagal mengosongkan keranjang';
				toast.error(msg);
			}
		}
	}

	async function handleKeepOnlyFirstSeller() {
		const firstProduct = cart.items[0]?.product;
		const firstSellerId =
			firstProduct?.seller_id || firstProduct?.sellerId || firstProduct?.seller?.id;
		if (firstSellerId) {
			await cart.keepOnlySeller(firstSellerId);
			toast.success('Hanya produk dari toko ini yang dipertahankan.');
		}
	}
</script>

<svelte:head>
	<title>Keranjang Belanja — SimpleShop</title>
</svelte:head>

<div class="mx-auto flex max-w-7xl flex-1 flex-col gap-6 px-4 py-8 sm:px-6 lg:px-8">
	<!-- Page Header -->
	<div
		class="flex flex-col items-start justify-between gap-4 border-b border-slate-200 pb-4 sm:flex-row sm:items-center"
	>
		<div>
			<h1 class="text-2xl font-bold tracking-tight text-slate-900 sm:text-3xl">
				Keranjang Belanja
			</h1>
			<p class="mt-0.5 text-xs text-slate-500 sm:text-sm">
				Tinjau barang yang ingin Anda beli sebelum melakukan proses checkout.
			</p>
		</div>

		{#if cart.items.length > 0}
			<button
				type="button"
				onclick={handleClearCart}
				class="flex cursor-pointer items-center gap-1.5 text-xs font-semibold text-slate-500 transition-colors hover:text-red-600"
			>
				<Trash2 class="h-3.5 w-3.5" />
				<span>Kosongkan Keranjang</span>
			</button>
		{/if}
	</div>

	{#if cart.items.length === 0}
		<!-- Empty Cart State -->
		<div
			class="my-4 flex flex-col items-center justify-center gap-3 rounded-xl border border-slate-200 bg-white py-20 text-center"
		>
			<div
				class="flex h-16 w-16 items-center justify-center rounded-full bg-slate-100 text-slate-400"
			>
				<ShoppingBag class="h-8 w-8 stroke-[1.5]" />
			</div>
			<h2 class="text-lg font-bold text-slate-900">Keranjang Belanja Kosong</h2>
			<p class="max-w-sm text-xs text-slate-500">
				Anda belum menambahkan produk apa pun ke keranjang. Jelajahi katalog dan temukan barang
				impian Anda.
			</p>
			<a
				href="/products"
				class="mt-2 inline-flex items-center gap-2 rounded-lg bg-slate-900 px-5 py-2.5 text-xs font-semibold text-white shadow-xs transition-colors hover:bg-slate-800"
			>
				<span>Mulai Belanja</span>
				<ArrowRight class="h-3.5 w-3.5" />
			</a>
		</div>
	{:else}
		<!-- Single-Seller Warning Notice (BR-010 Rule) -->
		{#if cart.hasMultipleSellers}
			<div
				class="flex flex-col gap-3 rounded-xl border border-amber-200 bg-amber-50 p-4 text-amber-900 sm:flex-row sm:items-center sm:justify-between"
			>
				<div class="flex items-start gap-3">
					<AlertTriangle class="mt-0.5 h-5 w-5 shrink-0 text-amber-600" />
					<div class="space-y-1 text-xs">
						<h4 class="font-bold text-amber-950">Peringatan: Produk Berasal dari Toko Berbeda</h4>
						<p class="leading-relaxed">
							Sesuai aturan bisnis sistem, satu pesanan checkout hanya dapat memuat produk dari
							<strong>satu penjual</strong>.
						</p>
					</div>
				</div>

				{#if cart.items[0]?.product}
					<button
						type="button"
						onclick={handleKeepOnlyFirstSeller}
						class="shrink-0 cursor-pointer rounded-lg bg-amber-600 px-3 py-1.5 text-xs font-semibold text-white shadow-xs transition-colors hover:bg-amber-700 active:scale-95"
					>
						Hanya simpan toko {cart.items[0].product.seller?.name || 'pertama'}
					</button>
				{/if}
			</div>
		{/if}

		<!-- Cart Content: Items List + Order Summary -->
		<div class="grid grid-cols-1 items-start gap-8 lg:grid-cols-3">
			<!-- Items Column (2/3 width) -->
			<div class="space-y-3 lg:col-span-2">
				{#each cart.items as item (item.id)}
					<CartItemRow {item} onUpdateQuantity={handleUpdateQuantity} onRemove={handleRemoveItem} />
				{/each}
			</div>

			<!-- Summary Column (1/3 width) -->
			<div
				class="sticky top-20 space-y-5 rounded-xl border border-slate-200 bg-white p-6 shadow-xs"
			>
				<h3
					class="border-b border-slate-100 pb-3 text-sm font-bold tracking-wider text-slate-900 uppercase"
				>
					Ringkasan Belanja
				</h3>

				<div class="space-y-3 text-xs text-slate-600">
					<div class="flex items-center justify-between">
						<span>Total Jumlah Barang</span>
						<span class="font-semibold text-slate-900 tabular-nums">{cart.totalItems} barang</span>
					</div>

					<div class="flex items-center justify-between">
						<span>Ongkos Kirim</span>
						<span class="font-semibold text-emerald-600">Gratis (Simulasi MVP)</span>
					</div>

					<div class="flex items-baseline justify-between border-t border-slate-100 pt-3">
						<span class="text-sm font-bold text-slate-900">Total Pembayaran</span>
						<span class="text-xl font-extrabold tracking-tight text-slate-900 tabular-nums">
							{formatIDR(cart.totalAmount)}
						</span>
					</div>
				</div>

				<div class="pt-2">
					<Button
						type="button"
						variant="primary"
						disabled={cart.hasMultipleSellers || cart.items.length === 0}
						onclick={() => goto('/checkout')}
						class="flex w-full items-center justify-center gap-2"
					>
						<span>Lanjut ke Checkout</span>
						<ArrowRight class="h-4 w-4" />
					</Button>
				</div>

				<div class="flex items-center justify-center gap-2 pt-2 text-[11px] text-slate-400">
					<ShieldCheck class="h-4 w-4 text-emerald-600" />
					<span>Stok divalidasi secara instan dan aman</span>
				</div>
			</div>
		</div>
	{/if}
</div>
