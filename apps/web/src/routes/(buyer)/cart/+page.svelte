<script lang="ts">
	import { goto } from '$app/navigation';
	import { cart } from '$lib/stores/cart.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import { formatIDR } from '$lib/utils/format';
	import TokoKuPediaLogo from '$lib/components/ui/TokoKuPediaLogo.svelte';
	import {
		ShoppingBag,
		ArrowRight,
		Trash2,
		ShieldCheck,
		MapPin,
		Minus,
		Plus,
		Heart,
		Ticket,
		ChevronRight,
		Sparkles
	} from 'lucide-svelte';

	let voucherCode = $state('');
	let appliedVoucher = $state<string | null>(null);
	let discountAmount = $state(0);
	let noteInputs = $state<Record<string, string>>({});

	const promoDiscount = $derived(appliedVoucher ? discountAmount : 0);
	const finalTotal = $derived(Math.max(0, cart.totalAmount - promoDiscount));

	async function handleUpdateQuantity(itemId: string, newQty: number) {
		try {
			await cart.updateQuantity(itemId, newQty);
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : 'Gagal memperbarui kuantitas';
			toast.error(msg);
		}
	}

	async function handleRemoveItem(itemId: string) {
		try {
			await cart.removeItem(itemId);
			toast.info('Barang dihapus dari keranjang.');
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : 'Gagal menghapus barang';
			toast.error(msg);
		}
	}

	async function handleClearCart() {
		if (confirm('Kosongkan seluruh keranjang belanja?')) {
			try {
				await cart.clearCart();
				toast.info('Keranjang berhasil dikosongkan.');
			} catch (err: unknown) {
				const msg = err instanceof Error ? err.message : 'Gagal mengosongkan keranjang';
				toast.error(msg);
			}
		}
	}

	function applyPromo() {
		if (voucherCode.trim().toUpperCase() === 'BEBASONGKIR' || voucherCode.trim().toUpperCase() === 'DISKON10') {
			appliedVoucher = voucherCode.trim().toUpperCase();
			discountAmount = 50000;
			toast.success(`Promo "${appliedVoucher}" berhasil digunakan! Hemat Rp50.000`);
		} else {
			toast.error('Kode promo tidak valid. Coba kode: BEBASONGKIR atau DISKON10');
		}
	}

	function handleCheckout() {
		if (cart.items.length === 0) {
			toast.warning('Keranjang belanja kosong');
			return;
		}
		goto('/checkout');
	}
</script>

<svelte:head>
	<title>Keranjang Belanja | TokoKuPedia</title>
</svelte:head>

<div class="min-h-screen bg-[#F5F6F7] pb-16">
	<!-- Cart Top Breadcrumb / Title -->
	<div class="border-b border-[#E5E7E9] bg-white">
		<div class="mx-auto flex max-w-7xl items-center justify-between px-4 py-4 sm:px-6 lg:px-8">
			<div class="flex items-center gap-3">
				<h1 class="text-xl font-bold tracking-tight text-[#212121]">Keranjang</h1>
				<span class="rounded-full bg-emerald-100 px-2.5 py-0.5 text-xs font-bold text-[#00AA5B]">
					{cart.totalItems} Barang
				</span>
			</div>

			{#if cart.items.length > 0}
				<button
					type="button"
					onclick={handleClearCart}
					class="flex items-center gap-1 text-xs font-semibold text-[#EE4D2D] hover:underline"
				>
					<Trash2 class="h-3.5 w-3.5" />
					<span>Hapus Semua</span>
				</button>
			{/if}
		</div>
	</div>

	<div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
		{#if cart.items.length === 0}
			<!-- Empty Cart State -->
			<div
				class="my-8 flex flex-col items-center justify-center rounded-2xl border border-[#E5E7E9] bg-white p-12 text-center shadow-xs"
			>
				<div class="flex h-24 w-24 items-center justify-center rounded-full bg-emerald-50 text-[#00AA5B]">
					<ShoppingBag class="h-12 w-12" />
				</div>
				<h2 class="mt-4 text-xl font-bold text-[#212121]">Wah, keranjang belanjamu kosong</h2>
				<p class="mt-1 text-xs text-[#6D7588]">
					Yuk, isi dengan barang-barang impianmu sekarang!
				</p>
				<a
					href="/"
					class="mt-6 rounded-lg bg-[#00AA5B] px-6 py-2.5 text-sm font-bold text-white shadow-xs transition-colors hover:bg-[#029b53]"
				>
					Mulai Belanja
				</a>
			</div>
		{:else}
			<!-- Cart 2-Column Content -->
			<div class="grid grid-cols-1 gap-8 lg:grid-cols-12">
				<!-- LEFT: Store & Items List (8 Cols) -->
				<div class="space-y-4 lg:col-span-8">
					<!-- Store Card Wrapper -->
					<div class="rounded-2xl border border-[#E5E7E9] bg-white p-5 shadow-xs">
						<!-- Store Header -->
						<div class="flex items-center gap-2 border-b border-[#E5E7E9] pb-3.5">
							<input
								type="checkbox"
								checked={true}
								class="h-4 w-4 rounded border-[#E5E7E9] text-[#00AA5B] focus:ring-[#00AA5B]"
							/>
							<div class="flex items-center gap-1.5 font-bold text-sm text-[#212121]">
								<ShieldCheck class="h-4 w-4 fill-[#00AA5B] text-white" />
								<span>Official Store Partner</span>
							</div>
							<span class="text-xs text-[#6D7588]">• Kota Jakarta Pusat</span>
						</div>

						<!-- Items within Store -->
						<div class="divide-y divide-[#E5E7E9]">
							{#each cart.items as item}
								<div class="py-4">
									<div class="flex items-start gap-3">
										<input
											type="checkbox"
											checked={true}
											class="mt-1 h-4 w-4 rounded border-[#E5E7E9] text-[#00AA5B] focus:ring-[#00AA5B]"
										/>

										<!-- Image Thumbnail -->
										<a
											href="/products/{item.product.id}"
											class="relative block h-20 w-20 shrink-0 overflow-hidden rounded-xl border border-[#E5E7E9] bg-slate-50"
										>
											<img
												src={item.product.image_url}
												alt={item.product.name}
												class="h-full w-full object-cover"
											/>
										</a>

										<!-- Item Detail -->
										<div class="flex-1 space-y-1">
											<a
												href="/products/{item.product.id}"
												class="line-clamp-2 text-xs font-semibold text-[#212121] hover:text-[#00AA5B]"
											>
												{item.product.name}
											</a>

											<div class="text-[11px] text-[#6D7588]">
												Varian: <span class="font-medium text-[#212121]">Default / 512GB</span>
											</div>

											<div class="pt-1">
												<span class="text-sm font-extrabold text-[#212121] tabular-nums">
													{formatIDR(item.product.price)}
												</span>
											</div>
										</div>
									</div>

									<!-- Item Bottom Controls (Note & Stepper) -->
									<div class="mt-3 flex flex-col justify-between gap-3 pl-7 sm:flex-row sm:items-center">
										<!-- Note Input -->
										<div class="flex-1">
											<input
												type="text"
												placeholder="Tulis catatan untuk penjual..."
												bind:value={noteInputs[item.id]}
												class="h-8 w-full max-w-sm rounded-lg border border-[#E5E7E9] px-2.5 text-xs text-[#212121] placeholder:text-[#9FA6B0] focus:border-[#00AA5B] focus:ring-1 focus:ring-[#00AA5B] focus:outline-none"
											/>
										</div>

										<!-- Stepper & Actions -->
										<div class="flex items-center gap-4">
											<button
												type="button"
												onclick={() => toast.info('Disimpan ke Wishlist!')}
												class="text-[#6D7588] hover:text-[#EE4D2D]"
												title="Pindahkan ke Wishlist"
											>
												<Heart class="h-4 w-4" />
											</button>

											<button
												type="button"
												onclick={() => handleRemoveItem(item.id)}
												class="text-[#6D7588] hover:text-[#EE4D2D]"
												title="Hapus barang"
											>
												<Trash2 class="h-4 w-4" />
											</button>

											<!-- Stepper -->
											<div class="flex items-center overflow-hidden rounded-lg border border-[#E5E7E9] bg-white">
												<button
													type="button"
													onclick={() => handleUpdateQuantity(item.id, item.quantity - 1)}
													disabled={item.quantity <= 1}
													class="flex h-7 w-7 items-center justify-center text-[#6D7588] hover:bg-slate-100 disabled:opacity-30"
												>
													<Minus class="h-3 w-3" />
												</button>
												<span class="w-8 text-center text-xs font-bold text-[#212121] tabular-nums">
													{item.quantity}
												</span>
												<button
													type="button"
													onclick={() => handleUpdateQuantity(item.id, item.quantity + 1)}
													disabled={item.quantity >= item.product.stock}
													class="flex h-7 w-7 items-center justify-center text-[#6D7588] hover:bg-slate-100 disabled:opacity-30"
												>
													<Plus class="h-3 w-3" />
												</button>
											</div>
										</div>
									</div>
								</div>
							{/each}
						</div>
					</div>
				</div>

				<!-- RIGHT: Ringkasan Belanja Sticky (4 Cols) -->
				<div class="lg:col-span-4">
					<div class="sticky top-28 space-y-4">
						<!-- Promo Coupon Card -->
						<div class="rounded-2xl border border-[#E5E7E9] bg-white p-4 shadow-xs">
							<div class="flex items-center gap-2 font-bold text-xs text-[#212121]">
								<Ticket class="h-4 w-4 text-[#00AA5B]" />
								<span>Makin hemat pakai promo TokoKuPedia</span>
							</div>

							<div class="mt-3 flex gap-2">
								<input
									type="text"
									placeholder="Masukkan kode promo"
									bind:value={voucherCode}
									class="h-9 flex-1 rounded-lg border border-[#E5E7E9] px-3 text-xs uppercase tracking-wider text-[#212121] placeholder:normal-case placeholder:tracking-normal focus:border-[#00AA5B] focus:ring-1 focus:ring-[#00AA5B] focus:outline-none"
								/>
								<button
									type="button"
									onclick={applyPromo}
									class="rounded-lg bg-[#00AA5B] px-3 text-xs font-bold text-white hover:bg-[#029b53]"
								>
									Terapkan
								</button>
							</div>

							{#if appliedVoucher}
								<div class="mt-2.5 flex items-center justify-between rounded-lg bg-[#F0FDF4] p-2 text-xs font-semibold text-[#00AA5B]">
									<span>Voucher {appliedVoucher} Aktif</span>
									<button
										type="button"
										onclick={() => {
											appliedVoucher = null;
											discountAmount = 0;
										}}
										class="text-[11px] text-[#EE4D2D] hover:underline"
									>
										Batal
									</button>
								</div>
							{/if}
						</div>

						<!-- Ringkasan Belanja Card -->
						<div class="rounded-2xl border border-[#E5E7E9] bg-white p-5 shadow-xs space-y-4">
							<h3 class="text-sm font-bold text-[#212121]">Ringkasan belanja</h3>

							<div class="space-y-2 text-xs">
								<div class="flex justify-between text-[#6D7588]">
									<span>Total Harga ({cart.totalItems} barang)</span>
									<span class="font-semibold text-[#212121] tabular-nums">
										{formatIDR(cart.totalAmount)}
									</span>
								</div>

								{#if promoDiscount > 0}
									<div class="flex justify-between text-[#00AA5B]">
										<span>Total Diskon Promo</span>
										<span class="font-semibold tabular-nums">
											-{formatIDR(promoDiscount)}
										</span>
									</div>
								{/if}

								<div class="flex justify-between text-[#6D7588]">
									<span>Biaya Pengiriman</span>
									<span class="font-bold text-[#00AA5B]">Gratis (Bebas Ongkir)</span>
								</div>

								<div class="border-t border-[#E5E7E9] pt-3 flex justify-between items-baseline">
									<span class="font-bold text-sm text-[#212121]">Total Belanja</span>
									<span class="text-xl font-extrabold text-[#212121] tabular-nums">
										{formatIDR(finalTotal)}
									</span>
								</div>
							</div>

							<!-- Checkout CTA -->
							<button
								type="button"
								onclick={handleCheckout}
								class="flex w-full cursor-pointer items-center justify-center gap-2 rounded-xl bg-[#00AA5B] py-3 text-sm font-bold text-white shadow-xs transition-colors hover:bg-[#029b53] active:scale-98"
							>
								<span>Beli ({cart.totalItems})</span>
								<ArrowRight class="h-4 w-4" />
							</button>
						</div>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>
