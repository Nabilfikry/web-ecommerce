<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import type { Product } from '@ecommerce/shared-types';
	import { productsApi, SAMPLE_PRODUCTS } from '$lib/api/products';
	import ProductCard from '$lib/components/product/ProductCard.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import { cart } from '$lib/stores/cart.svelte';
	import {
		ChevronLeft,
		ChevronRight,
		Smartphone,
		Zap,
		Wifi,
		Globe,
		MoreHorizontal,
		ArrowRight,
		Flame,
		Clock,
		ShieldCheck,
		Laptop,
		Tag,
		Sparkles
	} from 'lucide-svelte';

	let products = $state<Product[]>([]);
	let loading = $state(true);

	// Hero Carousel
	let currentSlide = $state(0);
	const heroSlides = [
		{
			id: 1,
			title: 'WAR TRAKTIRAN TOTAL 2 MILIAR',
			subtitle: 'MISI TOPED Berhadiah iPhone 18 Pro & Flash Sale Serba 11rb!',
			ctaText: 'Serbu Promo',
			href: '#flash-sale',
			bgGradient: 'from-[#00AA5B] via-[#03AC0E] to-[#008848]',
			badge: 'PROMO GUNCANG'
		},
		{
			id: 2,
			title: 'Microsoft Surface Pro 12" Copilot+ PC',
			subtitle: 'Snapdragon X Plus, RAM 16GB, SSD 512GB - Ekstra Cashback Rp2.000.000!',
			ctaText: 'Beli Sekarang',
			href: '/jullindo/microsoft-surface-pro-12-copilot-pc-snapdragon-x-plus-8-core-ram-16gb-ssd-256-512gb-12-inch-touchscreen-1731633667861546606',
			bgGradient: 'from-[#0F172A] via-[#1E293B] to-[#334155]',
			badge: 'OFFICIAL STORE LAUNCH'
		},
		{
			id: 3,
			title: 'SUPER GADGET & SETUP FESTIVAL',
			subtitle: 'Diskon hingga 70% + Bebas Ongkir tanpa minimum belanja ke seluruh Indonesia!',
			ctaText: 'Cek Katalog',
			href: '/products',
			bgGradient: 'from-[#0284C7] via-[#0369A1] to-[#075985]',
			badge: 'BEBAS ONGKIR'
		}
	];

	// Top Up & Tagihan Widget
	let topUpTab = $state<'pulsa' | 'paket' | 'listrik' | 'roaming'>('pulsa');
	let phoneNumber = $state('');
	let selectedNominal = $state('50000');
	let isBuyingTopUp = $state(false);

	// Category Chips
	const quickCategories = [
		{ name: 'Kategori', icon: Tag, href: '/products' },
		{ name: 'Handphone & Tablet', icon: Smartphone, href: '/products?category=gadget' },
		{ name: 'Top-Up & Tagihan', icon: Zap, href: '#topup' },
		{ name: 'Elektronik', icon: Zap, href: '/products?category=electronics' },
		{ name: 'Komputer & Laptop', icon: Laptop, href: '/products?category=pc' },
		{ name: 'Keuangan', icon: Sparkles, href: '#topup' }
	];

	// Active tab in feed
	let feedTab = $state<'semua' | 'gadget' | 'pc' | 'audio'>('semua');

	// Countdown Timer
	let hours = $state(2);
	let minutes = $state(45);
	let seconds = $state(18);

	let slideInterval: NodeJS.Timeout | undefined;
	let timerInterval: NodeJS.Timeout | undefined;

	onMount(async () => {
		// Auto slide hero carousel
		slideInterval = setInterval(() => {
			currentSlide = (currentSlide + 1) % heroSlides.length;
		}, 5000);

		// Countdown tick
		timerInterval = setInterval(() => {
			if (seconds > 0) {
				seconds -= 1;
			} else if (minutes > 0) {
				minutes -= 1;
				seconds = 59;
			} else if (hours > 0) {
				hours -= 1;
				minutes = 59;
				seconds = 59;
			}
		}, 1000);

		try {
			const res = await productsApi.getProducts({ limit: 12, sort: 'latest' });
			products = res.products.length > 0 ? res.products : SAMPLE_PRODUCTS;
		} catch (err: unknown) {
			products = SAMPLE_PRODUCTS;
		} finally {
			loading = false;
		}
	});

	onDestroy(() => {
		if (slideInterval) clearInterval(slideInterval);
		if (timerInterval) clearInterval(timerInterval);
	});

	function handleAddToCart(product: Product) {
		cart.addItem(product, 1);
		toast.success(`"${product.name}" ditambahkan ke keranjang belanja!`);
	}

	function handleTopUpSubmit() {
		if (!phoneNumber.trim()) {
			toast.error('Masukkan nomor telepon tujuan');
			return;
		}
		isBuyingTopUp = true;
		setTimeout(() => {
			isBuyingTopUp = false;
			toast.success(
				`Pesanan ${topUpTab.toUpperCase()} untuk ${phoneNumber} senilai Rp${Number(selectedNominal).toLocaleString('id-ID')} diproses!`
			);
			phoneNumber = '';
		}, 800);
	}
</script>

<svelte:head>
	<title>Situs Jual Beli Online Terlengkap, Mudah & Aman | TokoKuPedia</title>
</svelte:head>

<div class="space-y-6 pb-16">
	<!-- 1. HERO CAROUSEL SECTION -->
	<section class="mx-auto max-w-7xl px-4 pt-4 sm:px-6 lg:px-8">
		<div class="relative overflow-hidden rounded-2xl shadow-sm">
			<!-- Slide Container -->
			<div
				class="flex transition-transform duration-500 ease-out"
				style="transform: translateX(-{currentSlide * 100}%);"
			>
				{#each heroSlides as slide}
					<div
						class="relative flex h-64 w-full shrink-0 items-center justify-between bg-gradient-to-r {slide.bgGradient} p-8 sm:h-80 sm:p-12 text-white"
					>
						<!-- Left Decorative Text -->
						<div class="max-w-xl space-y-3">
							<span
								class="inline-block rounded-full bg-white/20 px-3 py-1 text-[11px] font-bold tracking-wider uppercase backdrop-blur-xs"
							>
								{slide.badge}
							</span>
							<h2 class="text-2xl font-extrabold tracking-tight sm:text-4xl sm:leading-tight">
								{slide.title}
							</h2>
							<p class="text-xs text-white/90 sm:text-sm">
								{slide.subtitle}
							</p>
							<div class="pt-2">
								<a
									href={slide.href}
									class="inline-flex items-center gap-2 rounded-lg bg-white px-5 py-2.5 text-xs font-bold text-[#00AA5B] shadow-sm transition-transform hover:scale-105 active:scale-95"
								>
									<span>{slide.ctaText}</span>
									<ArrowRight class="h-4 w-4" />
								</a>
							</div>
						</div>

						<!-- Right Illustration Mascot Graphic -->
						<div class="hidden pr-8 sm:block">
							<div class="relative flex h-48 w-48 items-center justify-center">
								<div class="absolute inset-0 rounded-full bg-white/10 blur-xl"></div>
								<div class="relative text-center">
									<div class="text-6xl font-extrabold tracking-tighter text-yellow-300 drop-shadow-md">
										2 Miliar
									</div>
									<div class="mt-1 font-bold text-xs tracking-wider uppercase text-white/90">
										Traktiran TokoKuPedia
									</div>
								</div>
							</div>
						</div>
					</div>
				{/each}
			</div>

			<!-- Prev / Next Controls -->
			<button
				type="button"
				onclick={() =>
					(currentSlide = (currentSlide - 1 + heroSlides.length) % heroSlides.length)}
				class="absolute top-1/2 left-4 -translate-y-1/2 flex h-9 w-9 items-center justify-center rounded-full bg-white/80 text-[#212121] shadow-md transition-colors hover:bg-white"
				aria-label="Previous slide"
			>
				<ChevronLeft class="h-5 w-5" />
			</button>

			<button
				type="button"
				onclick={() => (currentSlide = (currentSlide + 1) % heroSlides.length)}
				class="absolute top-1/2 right-4 -translate-y-1/2 flex h-9 w-9 items-center justify-center rounded-full bg-white/80 text-[#212121] shadow-md transition-colors hover:bg-white"
				aria-label="Next slide"
			>
				<ChevronRight class="h-5 w-5" />
			</button>

			<!-- Indicators Dots & Promo Link -->
			<div class="absolute bottom-4 left-6 flex items-center gap-1.5">
				{#each heroSlides as _, idx}
					<button
						type="button"
						onclick={() => (currentSlide = idx)}
						class="h-2 rounded-full transition-all {currentSlide === idx
							? 'w-6 bg-white'
							: 'w-2 bg-white/50'}"
						aria-label="Ke slide {idx + 1}"
					></button>
				{/each}
			</div>

			<a
				href="#promo"
				class="absolute bottom-4 right-6 text-[11px] font-semibold text-white/90 hover:underline"
			>
				Lihat Promo Lainnya
			</a>
		</div>
	</section>

	<!-- 2. KATEGORI POPULER & TOP UP TAGIHAN WIDGET (Matching home.png) -->
	<section class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<div class="grid grid-cols-1 gap-6 lg:grid-cols-12">
			<!-- Left Banner: Kategori Populer -->
			<div class="lg:col-span-6">
				<div class="mb-3 flex items-center justify-between">
					<h3 class="text-lg font-bold text-[#212121]">Kategori Populer</h3>
				</div>
				<div
					class="relative flex h-48 items-center justify-between overflow-hidden rounded-2xl bg-gradient-to-r from-[#00AA5B] to-[#03AC0E] p-6 text-white shadow-xs"
				>
					<div class="max-w-[240px] space-y-2">
						<h4 class="text-base font-extrabold leading-snug">
							Yuk, belanja di TokoKuPedia
						</h4>
						<p class="text-xs text-white/90">
							Barang lengkap dari beragam kategori pilihan
						</p>
						<div class="pt-2">
							<a
								href="/products"
								class="inline-flex items-center rounded-lg bg-white px-4 py-2 text-xs font-bold text-[#00AA5B] transition-transform hover:scale-105 active:scale-95"
							>
								Cek Sekarang
							</a>
						</div>
					</div>

					<!-- Cute Bag Graphics -->
					<div class="pr-2">
						<div class="flex h-24 w-24 items-center justify-center rounded-2xl bg-white/20 backdrop-blur-xs">
							<Smartphone class="h-12 w-12 text-white" />
						</div>
					</div>
				</div>
			</div>

			<!-- Right Widget: Top Up & Tagihan -->
			<div class="lg:col-span-6">
				<div class="mb-3 flex items-center justify-between">
					<h3 class="text-lg font-bold text-[#212121]">Top Up & Tagihan</h3>
					<a href="#topup" class="text-xs font-bold text-[#00AA5B] hover:underline">
						Lihat Semua
					</a>
				</div>

				<div class="rounded-2xl border border-[#E5E7E9] bg-white p-4 shadow-xs">
					<!-- Tabs -->
					<div class="flex border-b border-[#E5E7E9]">
						<button
							type="button"
							onclick={() => (topUpTab = 'pulsa')}
							class="border-b-2 px-4 py-2 text-xs font-bold transition-colors {topUpTab ===
							'pulsa'
								? 'border-[#00AA5B] text-[#00AA5B]'
								: 'border-transparent text-[#6D7588] hover:text-[#212121]'}"
						>
							Pulsa
						</button>
						<button
							type="button"
							onclick={() => (topUpTab = 'paket')}
							class="border-b-2 px-4 py-2 text-xs font-bold transition-colors {topUpTab ===
							'paket'
								? 'border-[#00AA5B] text-[#00AA5B]'
								: 'border-transparent text-[#6D7588] hover:text-[#212121]'}"
						>
							Paket Data
						</button>
						<button
							type="button"
							onclick={() => (topUpTab = 'listrik')}
							class="border-b-2 px-4 py-2 text-xs font-bold transition-colors {topUpTab ===
							'listrik'
								? 'border-[#00AA5B] text-[#00AA5B]'
								: 'border-transparent text-[#6D7588] hover:text-[#212121]'}"
						>
							Listrik PLN
						</button>
						<button
							type="button"
							onclick={() => (topUpTab = 'roaming')}
							class="border-b-2 px-4 py-2 text-xs font-bold transition-colors {topUpTab ===
							'roaming'
								? 'border-[#00AA5B] text-[#00AA5B]'
								: 'border-transparent text-[#6D7588] hover:text-[#212121]'}"
						>
							Roaming
						</button>
					</div>

					<!-- Form Inputs -->
					<div class="mt-4 grid grid-cols-1 gap-3 sm:grid-cols-12">
						<!-- Phone / Target Number -->
						<div class="sm:col-span-6">
							<label for="topup-phone" class="block text-[11px] font-bold text-[#6D7588]">
								Nomor Telepon
							</label>
							<input
								id="topup-phone"
								type="text"
								bind:value={phoneNumber}
								placeholder="Masukan Nomor (08xx)"
								class="mt-1 h-10 w-full rounded-lg border border-[#E5E7E9] px-3 text-xs text-[#212121] placeholder:text-[#9FA6B0] focus:border-[#00AA5B] focus:ring-1 focus:ring-[#00AA5B] focus:outline-none"
							/>
						</div>

						<!-- Nominal Selector -->
						<div class="sm:col-span-4">
							<label for="topup-nominal" class="block text-[11px] font-bold text-[#6D7588]">
								Nominal
							</label>
							<select
								id="topup-nominal"
								bind:value={selectedNominal}
								class="mt-1 h-10 w-full rounded-lg border border-[#E5E7E9] px-3 text-xs text-[#212121] focus:border-[#00AA5B] focus:ring-1 focus:ring-[#00AA5B] focus:outline-none"
							>
								<option value="25000">Rp 25.000</option>
								<option value="50000">Rp 50.000</option>
								<option value="100000">Rp 100.000</option>
								<option value="200000">Rp 200.000</option>
							</select>
						</div>

						<!-- Action Button -->
						<div class="flex items-end sm:col-span-2">
							<button
								type="button"
								onclick={handleTopUpSubmit}
								disabled={isBuyingTopUp}
								class="h-10 w-full rounded-lg bg-[#00AA5B] text-xs font-bold text-white transition-colors hover:bg-[#029b53] disabled:opacity-50"
							>
								Beli
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>

		<!-- Category Quick Pills Row -->
		<div class="mt-6 flex items-center gap-2 overflow-x-auto pb-2">
			{#each quickCategories as item}
				<a
					href={item.href}
					class="flex shrink-0 items-center gap-2 rounded-xl border border-[#E5E7E9] bg-white px-3.5 py-2 text-xs font-semibold text-[#212121] shadow-2xs transition-all hover:border-[#00AA5B] hover:text-[#00AA5B]"
				>
					<item.icon class="h-4 w-4 text-[#00AA5B]" />
					<span>{item.name}</span>
				</a>
			{/each}
		</div>
	</section>

	<!-- 3. FLASH SALE / TRAKTIRAN SECTION (With Live Countdown) -->
	<section id="flash-sale" class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<div class="rounded-2xl border border-red-100 bg-gradient-to-r from-red-500 via-rose-500 to-red-600 p-5 text-white shadow-md">
			<div class="flex flex-col items-start justify-between gap-4 sm:flex-row sm:items-center">
				<div class="flex flex-wrap items-center gap-3">
					<div class="flex items-center gap-1.5 rounded-full bg-white/20 px-3 py-1 font-bold text-xs uppercase backdrop-blur-xs">
						<Flame class="h-4 w-4 fill-yellow-300 text-yellow-300" />
						<span>Flash Sale Serba 11rb</span>
					</div>

					<div class="flex items-center gap-1.5 text-xs">
						<Clock class="h-3.5 w-3.5" />
						<span>Berakhir dalam:</span>
						<span class="rounded bg-black/40 px-2 py-0.5 font-mono font-bold">
							{hours.toString().padStart(2, '0')}:{minutes.toString().padStart(2, '0')}:{seconds.toString().padStart(2, '0')}
						</span>
					</div>
				</div>

				<a
					href="/products"
					class="text-xs font-bold text-white/90 hover:underline flex items-center gap-1"
				>
					<span>Lihat Semua Flash Sale</span>
					<ChevronRight class="h-4 w-4" />
				</a>
			</div>
		</div>
	</section>

	<!-- 4. REKOMENDASI UNTUK ANDA (Tokopedia Product Feed) -->
	<section class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<!-- Section Header & Filter Tabs -->
		<div class="mb-5 flex flex-col items-start justify-between gap-3 sm:flex-row sm:items-center">
			<div>
				<h2 class="text-xl font-bold tracking-tight text-[#212121]">
					Rekomendasi Untuk Anda
				</h2>
				<p class="text-xs text-[#6D7588]">
					Pilihan produk terbaik dengan promo menarik setiap hari
				</p>
			</div>

			<!-- Filter Tabs -->
			<div class="flex items-center gap-1 overflow-x-auto rounded-xl border border-[#E5E7E9] bg-white p-1 shadow-2xs">
				<button
					type="button"
					onclick={() => (feedTab = 'semua')}
					class="rounded-lg px-3 py-1.5 text-xs font-semibold transition-colors {feedTab ===
					'semua'
						? 'bg-[#00AA5B] text-white'
						: 'text-[#6D7588] hover:text-[#212121]'}"
				>
					Semua
				</button>
				<button
					type="button"
					onclick={() => (feedTab = 'pc')}
					class="rounded-lg px-3 py-1.5 text-xs font-semibold transition-colors {feedTab ===
					'pc'
						? 'bg-[#00AA5B] text-white'
						: 'text-[#6D7588] hover:text-[#212121]'}"
				>
					Komputer & PC
				</button>
				<button
					type="button"
					onclick={() => (feedTab = 'gadget')}
					class="rounded-lg px-3 py-1.5 text-xs font-semibold transition-colors {feedTab ===
					'gadget'
						? 'bg-[#00AA5B] text-white'
						: 'text-[#6D7588] hover:text-[#212121]'}"
				>
					Gadget & Audio
				</button>
			</div>
		</div>

		<!-- Product Grid (Tokopedia Cards) -->
		{#if loading}
			<div class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6">
				{#each Array(6) as _}
					<div class="aspect-square animate-pulse rounded-xl bg-slate-200"></div>
				{/each}
			</div>
		{:else}
			<div class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6">
				{#each products as product}
					<ProductCard {product} onAddToCart={handleAddToCart} />
				{/each}
			</div>

			<div class="mt-8 flex justify-center">
				<a
					href="/products"
					class="inline-flex items-center gap-2 rounded-xl border-2 border-[#00AA5B] bg-white px-8 py-3 text-sm font-bold text-[#00AA5B] transition-all hover:bg-[#F0FDF4] active:scale-95"
				>
					<span>Muat Lebih Banyak Produk</span>
					<ArrowRight class="h-4 w-4" />
				</a>
			</div>
		{/if}
	</section>
</div>
