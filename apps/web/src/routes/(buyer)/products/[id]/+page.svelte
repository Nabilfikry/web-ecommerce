<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import type { Product } from '@ecommerce/shared-types';
	import { productsApi, SAMPLE_PRODUCTS } from '$lib/api/products';
	import { formatIDR } from '$lib/utils/format';
	import { toast } from '$lib/stores/toast.svelte';
	import { cart } from '$lib/stores/cart.svelte';
	import {
		ShoppingCart,
		ShieldCheck,
		Truck,
		Minus,
		Plus,
		ArrowLeft,
		Star,
		Heart,
		Share2,
		MessageSquare,
		Check,
		MapPin,
		BadgeCheck,
		ChevronRight,
		Image as ImageIcon
	} from 'lucide-svelte';

	let product = $state<Product | null>(null);
	let loading = $state(true);
	let quantity = $state(1);
	let isAddingToCart = $state(false);
	let isBuyingNow = $state(false);
	let isWishlisted = $state(false);
	let activeTab = $state<'detail' | 'spesifikasi' | 'info'>('detail');

	// Variants
	let selectedColor = $state('Platinum');
	let selectedStorage = $state('512GB');

	const colors = ['Platinum', 'Black'];
	const storages = ['256GB', '512GB'];

	const galleryImages = [
		'https://images.unsplash.com/photo-1544244015-0df4b3ffc6b0?auto=format&fit=crop&w=800&q=80',
		'https://images.unsplash.com/photo-1527864550417-7fd91fc51a46?auto=format&fit=crop&w=800&q=80',
		'https://images.unsplash.com/photo-1587829741301-dc798b83add3?auto=format&fit=crop&w=800&q=80',
		'https://images.unsplash.com/photo-1518770660439-4636190af475?auto=format&fit=crop&w=800&q=80'
	];
	let activeImageIdx = $state(0);

	const productId = $derived(page.params.id);

	onMount(async () => {
		try {
			if (productId) {
				product = await productsApi.getProduct(productId);
			} else {
				product = SAMPLE_PRODUCTS[0];
			}
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : 'Produk tidak ditemukan';
			toast.error(msg);
			product = SAMPLE_PRODUCTS[0];
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

	async function handleAddToCart() {
		if (!product) return;
		if (product.stock <= 0) {
			toast.error('Stok produk habis');
			return;
		}

		isAddingToCart = true;
		try {
			await cart.addItem(product, quantity);
			toast.success(`Berhasil menambahkan ${quantity}x ke keranjang!`);
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : 'Gagal menambahkan produk ke keranjang';
			toast.error(msg);
		} finally {
			isAddingToCart = false;
		}
	}

	async function handleBuyNow() {
		if (!product) return;
		if (product.stock <= 0) {
			toast.error('Stok produk habis');
			return;
		}

		isBuyingNow = true;
		try {
			await cart.addItem(product, quantity);
			goto('/cart');
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : 'Gagal memproses pembelian langsung';
			toast.error(msg);
		} finally {
			isBuyingNow = false;
		}
	}

	function toggleWishlist() {
		isWishlisted = !isWishlisted;
		if (isWishlisted) {
			toast.success('Berhasil ditambahkan ke Wishlist!');
		} else {
			toast.info('Dihapus dari Wishlist.');
		}
	}

	function handleShare() {
		if (navigator.share) {
			navigator.share({
				title: product?.name || 'TokoKuPedia',
				url: window.location.href
			});
		} else {
			navigator.clipboard.writeText(window.location.href);
			toast.success('Link produk berhasil disalin!');
		}
	}

	function handleChat() {
		toast.info('Fitur Chat Penjual siap menghubungkan Anda dengan Jullindo Official!');
	}
</script>

<svelte:head>
	<title>{product ? `${product.name} | TokoKuPedia` : 'Detail Produk | TokoKuPedia'}</title>
</svelte:head>

<div class="min-h-screen bg-white">
	<!-- Breadcrumb Bar -->
	<div class="border-b border-[#E5E7E9] bg-white">
		<div class="mx-auto flex max-w-7xl items-center gap-1.5 px-4 py-3 text-xs text-[#6D7588] sm:px-6 lg:px-8">
			<a href="/" class="hover:text-[#00AA5B]">Home</a>
			<ChevronRight class="h-3.5 w-3.5" />
			<a href="/products" class="hover:text-[#00AA5B]">Komputer & Laptop</a>
			<ChevronRight class="h-3.5 w-3.5" />
			<a href="/products" class="hover:text-[#00AA5B]">Tablet & PC</a>
			<ChevronRight class="h-3.5 w-3.5" />
			<span class="max-w-[200px] truncate text-[#212121] font-medium sm:max-w-md">
				{product?.name || 'Microsoft Surface Pro 12'}
			</span>
		</div>
	</div>

	{#if loading}
		<div class="mx-auto flex max-w-7xl items-center justify-center py-24">
			<div class="flex flex-col items-center gap-3">
				<div class="h-10 w-10 animate-spin rounded-full border-3 border-[#00AA5B] border-t-transparent"></div>
				<span class="text-xs font-semibold text-[#6D7588]">Memuat informasi produk TokoKuPedia...</span>
			</div>
		</div>
	{:else if product}
		<!-- Main Product Content: 3 Column Layout -->
		<div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
			<div class="grid grid-cols-1 gap-8 lg:grid-cols-12">
				<!-- COLUMN 1: Image Gallery (4 Cols) -->
				<div class="lg:col-span-4">
					<div class="sticky top-28 space-y-4">
						<!-- Main Square Image Display -->
						<div
							class="relative aspect-square w-full overflow-hidden rounded-2xl border border-[#E5E7E9] bg-[#F5F6F7]"
						>
							<img
								src={galleryImages[activeImageIdx] || product.image_url}
								alt={product.name}
								class="h-full w-full object-cover transition-transform duration-300 hover:scale-105"
							/>
							<div class="absolute top-3 left-3 flex flex-col gap-1.5">
								<span
									class="inline-flex items-center gap-1 rounded bg-[#00AA5B] px-2 py-0.5 text-[10px] font-bold text-white shadow-xs"
								>
									<ShieldCheck class="h-3 w-3" />
									Official Store
								</span>
							</div>
						</div>

						<!-- Thumbnails Row -->
						<div class="flex items-center gap-3 overflow-x-auto pb-1">
							{#each galleryImages as img, idx}
								<button
									type="button"
									onclick={() => (activeImageIdx = idx)}
									class="relative h-16 w-16 shrink-0 overflow-hidden rounded-lg border-2 transition-all {activeImageIdx ===
									idx
										? 'border-[#00AA5B] ring-2 ring-[#00AA5B]/20'
										: 'border-[#E5E7E9] hover:border-slate-300'}"
								>
									<img src={img} alt="Thumbnail {idx + 1}" class="h-full w-full object-cover" />
								</button>
							{/each}
						</div>
					</div>
				</div>

				<!-- COLUMN 2: Center Details (5 Cols) -->
				<div class="space-y-6 lg:col-span-5">
					<!-- Title & Rating -->
					<div class="space-y-2 border-b border-[#E5E7E9] pb-4">
						<h1 class="text-xl font-bold leading-snug text-[#212121] sm:text-2xl">
							{product.name}
						</h1>

						<div class="flex flex-wrap items-center gap-3 text-xs text-[#6D7588]">
							<div class="flex items-center gap-1">
								<Star class="h-4 w-4 fill-[#FFAB00] text-[#FFAB00]" />
								<span class="font-bold text-[#212121]">4.9</span>
								<span>(45 rating)</span>
							</div>
							<span class="text-slate-300">•</span>
							<span>Terjual <strong class="font-bold text-[#212121]">100+</strong></span>
							<span class="text-slate-300">•</span>
							<a href="#diskusi" class="hover:text-[#00AA5B]">Diskusi (12)</a>
						</div>

						<!-- Price Display -->
						<div class="pt-3">
							<div class="text-3xl font-extrabold tracking-tight text-[#212121] tabular-nums">
								{formatIDR(product.price)}
							</div>
							<div class="mt-1 flex items-center gap-2">
								<span class="rounded bg-red-100 px-1.5 py-0.5 text-xs font-bold text-[#EE4D2D]">
									15%
								</span>
								<span class="text-sm text-[#6D7588] line-through tabular-nums">
									{formatIDR(Math.round(product.price * 1.15))}
								</span>
							</div>
						</div>
					</div>

					<!-- Variant Selectors -->
					<div class="space-y-4 border-b border-[#E5E7E9] pb-5">
						<!-- Color Variant -->
						<div>
							<div class="mb-2 text-xs font-bold text-[#212121]">
								Pilih Warna: <span class="font-normal text-[#6D7588]">{selectedColor}</span>
							</div>
							<div class="flex flex-wrap gap-2">
								{#each colors as c}
									<button
										type="button"
										onclick={() => (selectedColor = c)}
										class="rounded-lg border px-3.5 py-1.5 text-xs font-medium transition-all {selectedColor ===
										c
											? 'border-[#00AA5B] bg-[#F0FDF4] font-bold text-[#00AA5B] ring-1 ring-[#00AA5B]'
											: 'border-[#E5E7E9] text-[#212121] hover:border-slate-400'}"
									>
										{c}
									</button>
								{/each}
							</div>
						</div>

						<!-- Storage Variant -->
						<div>
							<div class="mb-2 text-xs font-bold text-[#212121]">
								Pilih Kapasitas: <span class="font-normal text-[#6D7588]">{selectedStorage}</span>
							</div>
							<div class="flex flex-wrap gap-2">
								{#each storages as s}
									<button
										type="button"
										onclick={() => (selectedStorage = s)}
										class="rounded-lg border px-3.5 py-1.5 text-xs font-medium transition-all {selectedStorage ===
										s
											? 'border-[#00AA5B] bg-[#F0FDF4] font-bold text-[#00AA5B] ring-1 ring-[#00AA5B]'
											: 'border-[#E5E7E9] text-[#212121] hover:border-slate-400'}"
									>
										{s}
									</button>
								{/each}
							</div>
						</div>
					</div>

					<!-- Seller / Store Card -->
					<div class="flex items-center justify-between rounded-xl border border-[#E5E7E9] p-3.5">
						<div class="flex items-center gap-3">
							<div
								class="flex h-11 w-11 items-center justify-center rounded-full bg-emerald-100 font-bold text-[#00AA5B]"
							>
								JO
							</div>
							<div>
								<div class="flex items-center gap-1.5">
									<span class="text-sm font-bold text-[#212121]">
										{product.seller?.name || 'Jullindo Official'}
									</span>
									<ShieldCheck class="h-4 w-4 fill-[#00AA5B] text-white" />
								</div>
								<div class="flex items-center gap-2 text-xs text-[#6D7588]">
									<span class="text-emerald-600 font-medium">Online 2 jam lalu</span>
									<span>•</span>
									<span class="flex items-center gap-0.5">
										<MapPin class="h-3 w-3" />
										Jakarta Pusat
									</span>
								</div>
							</div>
						</div>
						<button
							type="button"
							class="rounded-lg border border-[#00AA5B] px-3.5 py-1.5 text-xs font-bold text-[#00AA5B] hover:bg-[#F0FDF4]"
						>
							+ Ikuti
						</button>
					</div>

					<!-- Shipping Estimate -->
					<div class="space-y-2 rounded-xl border border-[#E5E7E9] bg-[#F5F6F7]/50 p-3.5 text-xs">
						<div class="flex items-center gap-2 font-bold text-[#212121]">
							<Truck class="h-4 w-4 text-[#00AA5B]" />
							<span>Pengiriman</span>
						</div>
						<div class="text-[#6D7588]">
							Dikirim dari <strong class="text-[#212121]">Kota Jakarta Pusat</strong> ke <strong class="text-[#212121]">Jakarta Pusat</strong>
						</div>
						<div class="flex items-center justify-between border-t border-[#E5E7E9] pt-2">
							<span>Ongkos Kirim Reguler:</span>
							<strong class="font-bold text-[#00AA5B]">Gratis (Bebas Ongkir)</strong>
						</div>
					</div>

					<!-- Tabs: Detail, Spesifikasi, Info Penting -->
					<div class="space-y-4 pt-2">
						<div class="flex border-b border-[#E5E7E9]">
							<button
								type="button"
								onclick={() => (activeTab = 'detail')}
								class="border-b-2 px-4 py-2.5 text-sm font-bold transition-colors {activeTab ===
								'detail'
									? 'border-[#00AA5B] text-[#00AA5B]'
									: 'border-transparent text-[#6D7588] hover:text-[#212121]'}"
							>
								Detail
							</button>
							<button
								type="button"
								onclick={() => (activeTab = 'spesifikasi')}
								class="border-b-2 px-4 py-2.5 text-sm font-bold transition-colors {activeTab ===
								'spesifikasi'
									? 'border-[#00AA5B] text-[#00AA5B]'
									: 'border-transparent text-[#6D7588] hover:text-[#212121]'}"
							>
								Spesifikasi
							</button>
							<button
								type="button"
								onclick={() => (activeTab = 'info')}
								class="border-b-2 px-4 py-2.5 text-sm font-bold transition-colors {activeTab ===
								'info'
									? 'border-[#00AA5B] text-[#00AA5B]'
									: 'border-transparent text-[#6D7588] hover:text-[#212121]'}"
							>
								Info Penting
							</button>
						</div>

						{#if activeTab === 'detail'}
							<div class="space-y-3 text-xs leading-relaxed text-[#212121]">
								<div class="grid grid-cols-2 gap-2 text-xs">
									<div>
										<span class="text-[#6D7588]">Kondisi: </span>
										<span class="font-semibold">Baru</span>
									</div>
									<div>
										<span class="text-[#6D7588]">Min. Pemesanan: </span>
										<span class="font-semibold">1 Buah</span>
									</div>
									<div>
										<span class="text-[#6D7588]">Etalase: </span>
										<span class="font-semibold text-[#00AA5B]">Microsoft Surface</span>
									</div>
								</div>
								<p class="pt-2 text-justify">
									{product.description}
								</p>
							</div>
						{:else if activeTab === 'spesifikasi'}
							<div class="space-y-2 text-xs">
								<div class="flex justify-between border-b border-slate-100 py-1.5">
									<span class="text-[#6D7588]">Processor</span>
									<span class="font-semibold text-[#212121]">Snapdragon X Plus (8 Core)</span>
								</div>
								<div class="flex justify-between border-b border-slate-100 py-1.5">
									<span class="text-[#6D7588]">RAM</span>
									<span class="font-semibold text-[#212121]">16GB LPDDR5X</span>
								</div>
								<div class="flex justify-between border-b border-slate-100 py-1.5">
									<span class="text-[#6D7588]">Storage</span>
									<span class="font-semibold text-[#212121]">{selectedStorage} SSD NVMe</span>
								</div>
								<div class="flex justify-between border-b border-slate-100 py-1.5">
									<span class="text-[#6D7588]">Layar</span>
									<span class="font-semibold text-[#212121]">12-inch Touchscreen PixelSense 120Hz</span>
								</div>
								<div class="flex justify-between border-b border-slate-100 py-1.5">
									<span class="text-[#6D7588]">Sistem Operasi</span>
									<span class="font-semibold text-[#212121]">Windows 11 Copilot+ PC</span>
								</div>
							</div>
						{:else}
							<div class="space-y-2 text-xs text-[#6D7588]">
								<p>• Produk 100% Original bergaransi resmi.</p>
								<p>• Harap rekam video unboxing saat membuka paket untuk klaim asuransi pengiriman.</p>
								<p>• Klaim garansi dapat diproses di service center resmi di seluruh Indonesia.</p>
							</div>
						{/if}
					</div>
				</div>

				<!-- COLUMN 3: Sticky Buy Box (3 Cols) -->
				<div class="lg:col-span-3">
					<div class="sticky top-28 space-y-4 rounded-2xl border border-[#E5E7E9] bg-white p-4 shadow-sm">
						<h3 class="text-sm font-bold text-[#212121]">Atur jumlah dan catatan</h3>

						<!-- Selected Variant Preview -->
						<div class="flex items-center gap-2 rounded-lg bg-slate-50 p-2 text-xs">
							<img src={galleryImages[0]} alt="variant" class="h-9 w-9 rounded object-cover" />
							<div class="truncate">
								<div class="font-semibold text-[#212121]">{selectedColor}, {selectedStorage}</div>
								<div class="text-[10px] text-[#6D7588]">Stok tersedia: {product.stock}</div>
							</div>
						</div>

						<!-- Quantity Counter -->
						<div class="flex items-center justify-between pt-2">
							<div class="flex items-center overflow-hidden rounded-lg border border-[#E5E7E9] bg-white">
								<button
									type="button"
									onclick={decrementQuantity}
									disabled={quantity <= 1}
									class="flex h-8 w-8 items-center justify-center text-[#6D7588] transition-colors hover:bg-slate-100 disabled:opacity-30"
								>
									<Minus class="h-3.5 w-3.5" />
								</button>
								<span class="w-10 text-center text-xs font-bold text-[#212121] tabular-nums">
									{quantity}
								</span>
								<button
									type="button"
									onclick={incrementQuantity}
									disabled={quantity >= product.stock}
									class="flex h-8 w-8 items-center justify-center text-[#6D7588] transition-colors hover:bg-slate-100 disabled:opacity-30"
								>
									<Plus class="h-3.5 w-3.5" />
								</button>
							</div>

							<div class="text-right">
								<span class="text-[11px] text-[#6D7588]">Sisa stok: </span>
								<strong class="text-xs font-bold text-[#212121]">{product.stock}</strong>
							</div>
						</div>

						<!-- Subtotal -->
						<div class="flex items-baseline justify-between border-t border-[#E5E7E9] pt-3">
							<span class="text-xs text-[#6D7588]">Subtotal:</span>
							<span class="text-lg font-extrabold text-[#212121] tabular-nums">
								{formatIDR(product.price * quantity)}
							</span>
						</div>

						<!-- CTA Buttons -->
						<div class="space-y-2 pt-2">
							<!-- + Keranjang -->
							<button
								type="button"
								onclick={handleAddToCart}
								disabled={product.stock <= 0 || isAddingToCart}
								class="flex w-full cursor-pointer items-center justify-center gap-2 rounded-lg border-2 border-[#00AA5B] bg-white py-2.5 text-sm font-bold text-[#00AA5B] transition-colors hover:bg-[#F0FDF4] active:scale-98 disabled:opacity-50"
							>
								<ShoppingCart class="h-4 w-4" />
								<span>+ Keranjang</span>
							</button>

							<!-- Beli Langsung -->
							<button
								type="button"
								onclick={handleBuyNow}
								disabled={product.stock <= 0 || isBuyingNow}
								class="flex w-full cursor-pointer items-center justify-center gap-2 rounded-lg bg-[#00AA5B] py-2.5 text-sm font-bold text-white shadow-xs transition-colors hover:bg-[#029b53] active:scale-98 disabled:opacity-50"
							>
								<span>Beli Langsung</span>
							</button>
						</div>

						<!-- Action Icons: Wishlist, Chat, Share -->
						<div class="flex items-center justify-around border-t border-[#E5E7E9] pt-3 text-xs text-[#6D7588]">
							<button
								type="button"
								onclick={handleChat}
								class="flex items-center gap-1 transition-colors hover:text-[#00AA5B]"
							>
								<MessageSquare class="h-3.5 w-3.5" />
								<span>Chat</span>
							</button>

							<button
								type="button"
								onclick={toggleWishlist}
								class="flex items-center gap-1 transition-colors {isWishlisted ? 'text-[#EE4D2D]' : 'hover:text-[#EE4D2D]'}"
							>
								<Heart class="h-3.5 w-3.5 {isWishlisted ? 'fill-[#EE4D2D]' : ''}" />
								<span>Wishlist</span>
							</button>

							<button
								type="button"
								onclick={handleShare}
								class="flex items-center gap-1 transition-colors hover:text-[#00AA5B]"
							>
								<Share2 class="h-3.5 w-3.5" />
								<span>Share</span>
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>
