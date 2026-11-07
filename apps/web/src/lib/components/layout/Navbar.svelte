<script lang="ts">
	import { auth } from '../../stores/auth.svelte';
	import { cart } from '../../stores/cart.svelte';
	import { toast } from '../../stores/toast.svelte';
	import { goto } from '$app/navigation';
	import TokoKuPediaLogo from '../ui/TokoKuPediaLogo.svelte';
	import {
		ShoppingCart,
		Search,
		LogOut,
		Store,
		MapPin,
		ChevronDown,
		Smartphone,
		Laptop,
		Tv,
		Headphones,
		HeartHandshake,
		Home as HomeIcon,
		HelpCircle,
		Tag,
		User,
		Settings,
		CreditCard,
		Building2,
		SunMoon,
		Menu,
		X
	} from 'lucide-svelte';

	interface Props {
		cartItemCount?: number;
	}

	let { cartItemCount = 0 }: Props = $props();

	let searchQuery = $state('');
	let categoryOpen = $state(false);
	let userMenuOpen = $state(false);
	let locationModalOpen = $state(false);
	let selectedCity = $state('Jakarta Pusat');
	let mobileMenuOpen = $state(false);

	const categories = [
		{ name: 'Handphone & Tablet', icon: Smartphone, href: '/products?category=gadget' },
		{ name: 'Komputer & Laptop', icon: Laptop, href: '/products?category=pc' },
		{ name: 'Elektronik', icon: Tv, href: '/products?category=electronics' },
		{ name: 'Audio & Aksesoris', icon: Headphones, href: '/products?category=audio' },
		{ name: 'Rumah Tangga & Setup', icon: HomeIcon, href: '/products?category=setup' },
		{ name: 'Keuangan & Tagihan', icon: CreditCard, href: '/#topup' }
	];

	const trendingSearches = [
		'Microsoft Surface Pro 12',
		'Mechanical Keyboard',
		'Laptop Stand',
		'Speaker Studio 50W',
		'Desk Mat Felt'
	];

	const cities = [
		'Jakarta Pusat',
		'Jakarta Selatan',
		'Jakarta Barat',
		'Bandung',
		'Surabaya',
		'Medan',
		'Semarang',
		'Yogyakarta'
	];

	function handleSearch(e?: SubmitEvent) {
		if (e) e.preventDefault();
		if (searchQuery.trim()) {
			categoryOpen = false;
			goto(`/products?search=${encodeURIComponent(searchQuery.trim())}`);
		} else {
			goto('/products');
		}
	}

	function handleQuickSearch(tag: string) {
		searchQuery = tag;
		goto(`/products?search=${encodeURIComponent(tag)}`);
	}

	function selectLocation(city: string) {
		selectedCity = city;
		locationModalOpen = false;
		toast.success(`Alamat pengiriman diubah ke ${city}`);
	}

	async function handleLogout() {
		try {
			await auth.logout();
			toast.success('Berhasil keluar akun');
			userMenuOpen = false;
			mobileMenuOpen = false;
			goto('/');
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : 'Gagal logout';
			toast.error(msg);
		}
	}
</script>

<header class="sticky top-0 z-50 w-full border-b border-[#E5E7E9] bg-white text-[#212121] shadow-xs">
	<!-- Top App Bar -->
	<div class="border-b border-slate-100 bg-[#F5F6F7] text-xs text-[#6D7588]">
		<div class="mx-auto flex h-8 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
			<!-- Mobile App Link -->
			<div class="flex items-center gap-2">
				<a
					href="/#promo"
					class="inline-flex items-center gap-1.5 rounded-full bg-[#00AA5B] px-2.5 py-0.5 text-[11px] font-semibold text-white transition-opacity hover:opacity-90"
				>
					<Smartphone class="h-3 w-3" />
					<span>Gratis Ongkir + Banyak Promo di aplikasi</span>
					<span class="text-[10px]">&rsaquo;</span>
				</a>
			</div>

			<!-- Secondary Links -->
			<div class="hidden items-center gap-6 md:flex">
				<a href="/#about" class="transition-colors hover:text-[#00AA5B]">Tentang TokoKuPedia</a>
				{#if auth.isSeller}
					<a href="/dashboard" class="font-semibold text-[#00AA5B] hover:underline"
						>Dashboard Seller</a
					>
				{:else}
					<a href="/register?role=seller" class="transition-colors hover:text-[#00AA5B]"
						>Mulai Berjualan</a
					>
				{/if}
				<a href="/#promo" class="transition-colors hover:text-[#00AA5B]">Promo</a>
				<a href="/#care" class="flex items-center gap-1 transition-colors hover:text-[#00AA5B]">
					<HelpCircle class="h-3 w-3" />
					<span>TokoKuPedia Care</span>
				</a>
			</div>
		</div>
	</div>

	<!-- Main Header Bar -->
	<div class="mx-auto max-w-7xl px-4 py-3 sm:px-6 lg:px-8">
		<div class="flex items-center justify-between gap-4">
			<!-- Logo -->
			<a href="/" class="group flex shrink-0 items-center gap-2">
				<TokoKuPediaLogo size="md" />
			</a>

			<!-- Kategori Dropdown Trigger -->
			<div class="relative hidden lg:block">
				<button
					type="button"
					onclick={() => (categoryOpen = !categoryOpen)}
					class="flex items-center gap-1 rounded-lg px-2.5 py-1.5 text-sm font-medium text-[#212121] transition-colors hover:bg-slate-100"
				>
					<span>Kategori</span>
					<ChevronDown
						class="h-4 w-4 text-[#6D7588] transition-transform duration-200 {categoryOpen
							? 'rotate-180'
							: ''}"
					/>
				</button>

				<!-- Categories Flyout Menu -->
				{#if categoryOpen}
					<div
						class="absolute top-full left-0 z-50 mt-2 w-72 rounded-xl border border-[#E5E7E9] bg-white p-3 shadow-xl"
					>
						<div class="mb-2 px-2 text-xs font-bold text-[#6D7588] uppercase tracking-wider">
							Kategori Pilihan
						</div>
						<div class="space-y-1">
							{#each categories as cat}
								<a
									href={cat.href}
									onclick={() => (categoryOpen = false)}
									class="flex items-center gap-3 rounded-lg px-3 py-2 text-sm text-[#212121] transition-colors hover:bg-[#F0FDF4] hover:text-[#00AA5B]"
								>
									<cat.icon class="h-4 w-4 text-[#00AA5B]" />
									<span>{cat.name}</span>
								</a>
							{/each}
						</div>
					</div>
				{/if}
			</div>

			<!-- Search Bar -->
			<div class="flex-1">
				<form onsubmit={handleSearch} class="relative flex w-full items-center">
					<div class="relative w-full">
						<input
							type="search"
							bind:value={searchQuery}
							placeholder="Cari di TokoKuPedia (contoh: Microsoft Surface, Keyboard)"
							class="h-10 w-full rounded-lg border border-[#E5E7E9] bg-white pr-10 pl-10 text-sm text-[#212121] placeholder:text-[#9FA6B0] transition-all focus:border-[#00AA5B] focus:ring-1 focus:ring-[#00AA5B] focus:outline-none"
						/>
						<Search class="pointer-events-none absolute top-3 left-3 h-4 w-4 text-[#9FA6B0]" />
						<button
							type="submit"
							class="absolute top-1.5 right-1.5 flex h-7 items-center rounded-md bg-[#00AA5B] px-3 text-xs font-bold text-white transition-colors hover:bg-[#029b53]"
						>
							Cari
						</button>
					</div>
				</form>
			</div>

			<!-- Cart Icon Button -->
			<a
				href="/cart"
				class="relative flex h-10 w-10 shrink-0 items-center justify-center rounded-lg text-[#212121] transition-colors hover:bg-slate-100"
				aria-label="Keranjang Belanja"
			>
				<ShoppingCart class="h-5 w-5" />
				{#if (cartItemCount || cart.totalItems) > 0}
					<span
						class="absolute top-1 right-1 flex h-4 min-w-[16px] items-center justify-center rounded-full bg-[#EE4D2D] px-1 text-[10px] font-extrabold text-white"
					>
						{(cartItemCount || cart.totalItems) > 99 ? '99+' : (cartItemCount || cart.totalItems)}
					</span>
				{/if}
			</a>

			<!-- Divider -->
			<div class="hidden h-6 w-px bg-[#E5E7E9] sm:block"></div>

			<!-- Auth Buttons / User Profile -->
			{#if auth.isAuthenticated}
				<!-- Logged In User Dropdown -->
				<div class="relative">
					<button
						type="button"
						onclick={() => (userMenuOpen = !userMenuOpen)}
						class="flex items-center gap-2 rounded-lg p-1.5 transition-colors hover:bg-slate-100"
					>
						<div
							class="flex h-8 w-8 items-center justify-center rounded-full bg-[#00AA5B] font-bold text-white shadow-xs"
						>
							{auth.user?.name ? auth.user.name[0].toUpperCase() : 'U'}
						</div>
						<span class="hidden max-w-[120px] truncate text-sm font-semibold text-[#212121] sm:block">
							{auth.user?.name || 'Pengguna'}
						</span>
						<ChevronDown class="h-3.5 w-3.5 text-[#6D7588]" />
					</button>

					{#if userMenuOpen}
						<div
							class="absolute right-0 z-50 mt-2 w-64 rounded-xl border border-[#E5E7E9] bg-white p-2 shadow-xl"
						>
							<div class="border-b border-[#E5E7E9] p-3">
								<div class="font-bold text-[#212121]">{auth.user?.name}</div>
								<div class="text-xs text-[#6D7588]">{auth.user?.email}</div>
								<div class="mt-2 inline-flex items-center gap-1 rounded bg-emerald-50 px-2 py-0.5 text-[10px] font-bold text-[#00AA5B]">
									{auth.isSeller ? 'Akun Penjual' : 'Akun Pembeli'}
								</div>
							</div>

							<div class="space-y-0.5 py-1 text-sm text-[#212121]">
								<a
									href="/user/settings/address"
									onclick={() => (userMenuOpen = false)}
									class="flex items-center gap-2.5 rounded-lg px-3 py-2 transition-colors hover:bg-slate-100"
								>
									<MapPin class="h-4 w-4 text-[#6D7588]" />
									<span>Daftar Alamat</span>
								</a>

								<a
									href="/user/settings/payment"
									onclick={() => (userMenuOpen = false)}
									class="flex items-center gap-2.5 rounded-lg px-3 py-2 transition-colors hover:bg-slate-100"
								>
									<CreditCard class="h-4 w-4 text-[#6D7588]" />
									<span>Metode Pembayaran</span>
								</a>

								<a
									href="/user/settings/bank"
									onclick={() => (userMenuOpen = false)}
									class="flex items-center gap-2.5 rounded-lg px-3 py-2 transition-colors hover:bg-slate-100"
								>
									<Building2 class="h-4 w-4 text-[#6D7588]" />
									<span>Rekening Bank</span>
								</a>

								<a
									href="/user/settings/mode-screen"
									onclick={() => (userMenuOpen = false)}
									class="flex items-center gap-2.5 rounded-lg px-3 py-2 transition-colors hover:bg-slate-100"
								>
									<SunMoon class="h-4 w-4 text-[#6D7588]" />
									<span>Mode Tampilan</span>
								</a>

								<a
									href="/orders"
									onclick={() => (userMenuOpen = false)}
									class="flex items-center gap-2.5 rounded-lg px-3 py-2 transition-colors hover:bg-slate-100"
								>
									<Tag class="h-4 w-4 text-[#6D7588]" />
									<span>Pesanan Saya</span>
								</a>

								{#if auth.isSeller}
									<a
										href="/dashboard"
										onclick={() => (userMenuOpen = false)}
										class="flex items-center gap-2.5 rounded-lg px-3 py-2 font-semibold text-[#00AA5B] transition-colors hover:bg-emerald-50"
									>
										<Store class="h-4 w-4 text-[#00AA5B]" />
										<span>Dashboard Seller</span>
									</a>
								{/if}
							</div>

							<div class="border-t border-[#E5E7E9] pt-1">
								<button
									type="button"
									onclick={handleLogout}
									class="flex w-full items-center gap-2.5 rounded-lg px-3 py-2 text-sm text-[#EE4D2D] transition-colors hover:bg-red-50"
								>
									<LogOut class="h-4 w-4" />
									<span>Keluar Akun</span>
								</button>
							</div>
						</div>
					{/if}
				</div>
			{:else}
				<!-- Guest Auth CTAs (Masuk & Daftar) -->
				<div class="flex items-center gap-2">
					<a
						href="/login"
						class="flex h-9 items-center justify-center rounded-lg border border-[#00AA5B] bg-white px-4 text-sm font-bold text-[#00AA5B] transition-colors hover:bg-[#F0FDF4]"
					>
						Masuk
					</a>
					<a
						href="/register"
						class="flex h-9 items-center justify-center rounded-lg bg-[#00AA5B] px-4 text-sm font-bold text-white shadow-xs transition-colors hover:bg-[#029b53]"
					>
						Daftar
					</a>
				</div>
			{/if}
		</div>

		<!-- Sub Header Bar: Location & Quick Keywords -->
		<div class="mt-2.5 flex items-center justify-between text-xs text-[#6D7588]">
			<!-- Delivery Location Selector -->
			<div class="relative">
				<button
					type="button"
					onclick={() => (locationModalOpen = !locationModalOpen)}
					class="flex items-center gap-1 font-medium transition-colors hover:text-[#00AA5B]"
				>
					<MapPin class="h-3.5 w-3.5 text-[#00AA5B]" />
					<span>Dikirim ke</span>
					<strong class="font-bold text-[#212121]">{selectedCity}</strong>
					<ChevronDown class="h-3 w-3" />
				</button>

				{#if locationModalOpen}
					<div
						class="absolute top-full left-0 z-50 mt-1.5 w-60 rounded-xl border border-[#E5E7E9] bg-white p-3 shadow-xl"
					>
						<div class="mb-2 text-xs font-bold text-[#212121]">Pilih Kota Pengiriman</div>
						<div class="space-y-1">
							{#each cities as city}
								<button
									type="button"
									onclick={() => selectLocation(city)}
									class="flex w-full items-center justify-between rounded-lg px-2.5 py-1.5 text-left text-xs transition-colors {selectedCity ===
									city
										? 'bg-[#F0FDF4] font-bold text-[#00AA5B]'
										: 'text-[#212121] hover:bg-slate-100'}"
								>
									<span>{city}</span>
									{#if selectedCity === city}
										<span>✓</span>
									{/if}
								</button>
							{/each}
						</div>
					</div>
				{/if}
			</div>

			<!-- Trending search tags -->
			<div class="hidden items-center gap-3 overflow-x-auto lg:flex">
				{#each trendingSearches as tag}
					<button
						type="button"
						onclick={() => handleQuickSearch(tag)}
						class="whitespace-nowrap transition-colors hover:text-[#00AA5B]"
					>
						{tag}
					</button>
				{/each}
			</div>
		</div>
	</div>
</header>
