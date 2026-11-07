<script lang="ts">
	import { auth } from '../../stores/auth.svelte';
	import { toast } from '../../stores/toast.svelte';
	import { goto } from '$app/navigation';
	import { ShoppingCart, Search, LogOut, Store, Package, Menu, X } from 'lucide-svelte';

	interface Props {
		cartItemCount?: number;
	}

	let { cartItemCount = 0 }: Props = $props();

	let searchQuery = $state('');
	let mobileMenuOpen = $state(false);
	let userMenuOpen = $state(false);

	function handleSearch(e: SubmitEvent) {
		e.preventDefault();
		if (searchQuery.trim()) {
			goto(`/products?search=${encodeURIComponent(searchQuery.trim())}`);
		} else {
			goto('/products');
		}
	}

	async function handleLogout() {
		try {
			await auth.logout();
			toast.success('Berhasil keluar akun');
			userMenuOpen = false;
			mobileMenuOpen = false;
			goto('/login');
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : 'Gagal logout';
			toast.error(msg);
		}
	}
</script>

<header class="sticky top-0 z-40 border-b border-slate-200 bg-white">
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<div class="flex h-16 items-center justify-between gap-4">
			<!-- Logo -->
			<div class="flex items-center gap-6">
				<a href="/" class="group flex items-center gap-2 text-slate-900">
					<div
						class="flex h-9 w-9 items-center justify-center rounded-lg bg-slate-900 text-lg font-bold text-white shadow-xs transition-colors group-hover:bg-slate-800"
					>
						E
					</div>
					<div class="flex flex-col">
						<span class="text-base leading-tight font-bold tracking-tight">SimpleShop</span>
						<span class="text-[10px] font-medium tracking-wide text-slate-500 uppercase"
							>Learning MVP</span
						>
					</div>
				</a>

				<!-- Desktop Navigation Links -->
				<nav class="hidden items-center gap-1 md:flex">
					<a
						href="/products"
						class="rounded-md px-3 py-1.5 text-sm font-medium text-slate-600 transition-colors hover:bg-slate-100 hover:text-slate-900"
					>
						Katalog Produk
					</a>

					{#if auth.isBuyer}
						<a
							href="/orders"
							class="rounded-md px-3 py-1.5 text-sm font-medium text-slate-600 transition-colors hover:bg-slate-100 hover:text-slate-900"
						>
							Pesanan Saya
						</a>
					{/if}

					{#if auth.isSeller}
						<a
							href="/dashboard"
							class="flex items-center gap-1.5 rounded-md bg-emerald-50 px-3 py-1.5 text-sm font-semibold text-emerald-700 transition-colors hover:bg-emerald-100"
						>
							<Store class="h-3.5 w-3.5" />
							Dashboard Seller
						</a>
					{/if}
				</nav>
			</div>

			<!-- Search Bar -->
			<form
				onsubmit={handleSearch}
				class="relative mx-2 hidden max-w-md flex-1 items-center sm:flex"
			>
				<Search class="pointer-events-none absolute left-3 h-4 w-4 text-slate-400" />
				<input
					type="search"
					placeholder="Cari produk impianmu..."
					bind:value={searchQuery}
					class="h-10 w-full rounded-lg border border-slate-300 bg-slate-50 pr-4 pl-9 text-sm text-slate-900 transition-all placeholder:text-slate-400 focus:border-slate-900 focus:bg-white focus:ring-1 focus:ring-slate-900 focus:outline-none"
				/>
			</form>

			<!-- Action Icons / Profile -->
			<div class="flex items-center gap-2">
				<!-- Cart (Buyer or Guest) -->
				{#if !auth.isSeller}
					<a
						href="/cart"
						class="relative cursor-pointer rounded-lg p-2 text-slate-700 transition-colors hover:bg-slate-100 hover:text-slate-900"
						aria-label="Keranjang Belanja"
					>
						<ShoppingCart class="h-5 w-5" />
						{#if cartItemCount > 0}
							<span
								class="absolute -top-1 -right-1 flex h-5 min-w-[20px] items-center justify-center rounded-full bg-amber-600 px-1 text-[11px] font-bold text-white tabular-nums shadow-xs"
							>
								{cartItemCount > 99 ? '99+' : cartItemCount}
							</span>
						{/if}
					</a>
				{/if}

				<!-- User Auth -->
				{#if auth.isAuthenticated && auth.user}
					<div class="relative">
						<button
							type="button"
							onclick={() => (userMenuOpen = !userMenuOpen)}
							class="flex cursor-pointer items-center gap-2 rounded-lg p-1.5 text-left transition-colors hover:bg-slate-100"
							aria-expanded={userMenuOpen}
						>
							<div
								class="flex h-8 w-8 items-center justify-center rounded-full bg-slate-800 text-xs font-semibold text-white"
							>
								{auth.user.name.charAt(0).toUpperCase()}
							</div>
							<div class="hidden flex-col lg:flex">
								<span class="text-xs leading-tight font-semibold text-slate-900">
									{auth.user.name}
								</span>
								<span class="text-[10px] leading-tight text-slate-500 capitalize">
									{auth.user.role}
								</span>
							</div>
						</button>

						<!-- Dropdown Menu -->
						{#if userMenuOpen}
							<div
								class="animate-in fade-in absolute right-0 z-50 mt-1 w-52 rounded-lg border border-slate-200 bg-white py-1 text-sm shadow-lg duration-100"
							>
								<div class="border-b border-slate-100 px-3.5 py-2">
									<p class="truncate font-semibold text-slate-900">{auth.user.name}</p>
									<p class="truncate text-xs text-slate-500">{auth.user.email}</p>
									<span
										class="mt-1 inline-block rounded px-2 py-0.5 text-[10px] font-semibold uppercase {auth
											.user.role === 'seller'
											? 'bg-emerald-100 text-emerald-800'
											: 'bg-blue-100 text-blue-800'}"
									>
										Role: {auth.user.role}
									</span>
								</div>

								{#if auth.isBuyer}
									<a
										href="/orders"
										onclick={() => (userMenuOpen = false)}
										class="flex items-center gap-2 px-3.5 py-2 text-slate-700 hover:bg-slate-50 hover:text-slate-900"
									>
										<Package class="h-4 w-4 text-slate-400" />
										Riwayat Pesanan
									</a>
								{:else if auth.isSeller}
									<a
										href="/dashboard"
										onclick={() => (userMenuOpen = false)}
										class="flex items-center gap-2 px-3.5 py-2 text-slate-700 hover:bg-slate-50 hover:text-slate-900"
									>
										<Store class="h-4 w-4 text-emerald-600" />
										Dashboard Penjual
									</a>
								{/if}

								<div class="my-1 border-t border-slate-100"></div>

								<button
									type="button"
									onclick={handleLogout}
									class="flex w-full cursor-pointer items-center gap-2 px-3.5 py-2 text-left text-red-600 hover:bg-red-50"
								>
									<LogOut class="h-4 w-4 text-red-500" />
									Keluar (Logout)
								</button>
							</div>
						{/if}
					</div>
				{:else}
					<div class="hidden items-center gap-2 sm:flex">
						<a
							href="/login"
							class="rounded-lg px-3 py-2 text-xs font-semibold text-slate-700 transition-colors hover:bg-slate-100 hover:text-slate-900"
						>
							Masuk
						</a>
						<a
							href="/register"
							class="rounded-lg bg-slate-900 px-3 py-2 text-xs font-semibold text-white shadow-xs transition-colors hover:bg-slate-800"
						>
							Daftar Akun
						</a>
					</div>
				{/if}

				<!-- Mobile Menu Toggle Button -->
				<button
					type="button"
					onclick={() => (mobileMenuOpen = !mobileMenuOpen)}
					class="rounded-lg p-2 text-slate-600 hover:bg-slate-100 hover:text-slate-900 md:hidden"
					aria-label="Toggle Menu"
				>
					{#if mobileMenuOpen}
						<X class="h-5 w-5" />
					{:else}
						<Menu class="h-5 w-5" />
					{/if}
				</button>
			</div>
		</div>

		<!-- Mobile Search (Visible on small screens) -->
		<div class="pb-3 sm:hidden">
			<form onsubmit={handleSearch} class="relative flex items-center">
				<Search class="pointer-events-none absolute left-3 h-4 w-4 text-slate-400" />
				<input
					type="search"
					placeholder="Cari produk..."
					bind:value={searchQuery}
					class="h-9 w-full rounded-lg border border-slate-300 bg-slate-50 pr-3 pl-9 text-xs text-slate-900 focus:bg-white focus:outline-none"
				/>
			</form>
		</div>
	</div>

	<!-- Mobile Menu Drawer -->
	{#if mobileMenuOpen}
		<div class="space-y-1 border-t border-slate-200 bg-white px-4 pt-2 pb-4 md:hidden">
			<a
				href="/products"
				onclick={() => (mobileMenuOpen = false)}
				class="block rounded-md px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100"
			>
				Katalog Produk
			</a>

			{#if auth.isBuyer}
				<a
					href="/orders"
					onclick={() => (mobileMenuOpen = false)}
					class="block rounded-md px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100"
				>
					Pesanan Saya
				</a>
			{/if}

			{#if auth.isSeller}
				<a
					href="/dashboard"
					onclick={() => (mobileMenuOpen = false)}
					class="block rounded-md bg-emerald-50 px-3 py-2 text-sm font-semibold text-emerald-800"
				>
					Dashboard Seller
				</a>
			{/if}

			<div class="my-2 border-t border-slate-100 pt-2"></div>

			{#if !auth.isAuthenticated}
				<div class="grid grid-cols-2 gap-2 pt-1">
					<a
						href="/login"
						onclick={() => (mobileMenuOpen = false)}
						class="rounded-lg border border-slate-300 py-2 text-center text-sm font-medium text-slate-800"
					>
						Masuk
					</a>
					<a
						href="/register"
						onclick={() => (mobileMenuOpen = false)}
						class="rounded-lg bg-slate-900 py-2 text-center text-sm font-medium text-white"
					>
						Daftar
					</a>
				</div>
			{:else}
				<button
					type="button"
					onclick={handleLogout}
					class="flex w-full items-center gap-2 rounded-md px-3 py-2 text-left text-sm font-medium text-red-600 hover:bg-red-50"
				>
					<LogOut class="h-4 w-4" />
					Keluar ({auth.user?.name})
				</button>
			{/if}
		</div>
	{/if}
</header>
