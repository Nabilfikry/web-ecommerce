<script lang="ts">
	import { page } from '$app/state';
	import { auth } from '$lib/stores/auth.svelte';
	import type { LayoutProps } from './$types';
	import {
		User,
		MapPin,
		CreditCard,
		Building2,
		SunMoon,
		ShieldCheck,
		Bell,
		Camera,
		Wallet,
		ChevronRight
	} from 'lucide-svelte';

	let { children }: LayoutProps = $props();

	const currentPath = $derived(page.url.pathname);

	const navItems = [
		{ label: 'Biodata Diri', href: '/user/settings/profile', icon: User },
		{ label: 'Daftar Alamat', href: '/user/settings/address', icon: MapPin },
		{ label: 'Metode Pembayaran', href: '/user/settings/payment', icon: CreditCard },
		{ label: 'Rekening Bank', href: '/user/settings/bank', icon: Building2 },
		{ label: 'Mode Tampilan', href: '/user/settings/mode-screen', icon: SunMoon }
	];
</script>

<div class="min-h-screen bg-[#F5F6F7] pb-16">
	<!-- Breadcrumb Bar -->
	<div class="border-b border-[#E5E7E9] bg-white">
		<div class="mx-auto flex max-w-7xl items-center gap-1.5 px-4 py-3 text-xs text-[#6D7588] sm:px-6 lg:px-8">
			<a href="/" class="hover:text-[#00AA5B]">Home</a>
			<ChevronRight class="h-3.5 w-3.5" />
			<span class="text-[#212121] font-medium">Pengaturan Akun</span>
		</div>
	</div>

	<!-- Container -->
	<div class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
		<div class="grid grid-cols-1 gap-8 lg:grid-cols-12">
			<!-- LEFT SIDEBAR: User Profile & Nav Menu (4 Cols) -->
			<div class="lg:col-span-4">
				<div class="space-y-4">
					<!-- User Summary Card -->
					<div class="rounded-2xl border border-[#E5E7E9] bg-white p-5 shadow-xs">
						<div class="flex items-center gap-4">
							<div class="relative">
								<div
									class="flex h-14 w-14 items-center justify-center rounded-full bg-[#00AA5B] text-xl font-bold text-white shadow-xs"
								>
									{auth.user?.name ? auth.user.name[0].toUpperCase() : 'N'}
								</div>
								<button
									type="button"
									class="absolute -bottom-1 -right-1 flex h-6 w-6 items-center justify-center rounded-full border border-white bg-slate-100 text-[#6D7588] hover:bg-slate-200"
									title="Ubah Foto"
								>
									<Camera class="h-3.5 w-3.5" />
								</button>
							</div>

							<div class="truncate">
								<div class="text-sm font-bold text-[#212121]">
									{auth.user?.name || 'Naabil Fikry'}
								</div>
								<div class="text-xs text-[#6D7588]">
									{auth.user?.email || 'buyer@tokokupedia.com'}
								</div>
								<div class="mt-1 flex items-center gap-1 text-[11px] font-semibold text-[#00AA5B]">
									<ShieldCheck class="h-3.5 w-3.5" />
									<span>Akun Terverifikasi</span>
								</div>
							</div>
						</div>

						<!-- Saldo & GoPay Quick Info -->
						<div class="mt-4 rounded-xl border border-slate-100 bg-[#F5F6F7] p-3 text-xs">
							<div class="flex items-center justify-between">
								<div class="flex items-center gap-1.5 text-[#6D7588]">
									<Wallet class="h-3.5 w-3.5 text-[#00AA5B]" />
									<span>Saldo TokoKuPedia:</span>
								</div>
								<span class="font-bold text-[#212121]">Rp1.500.000</span>
							</div>
							<div class="mt-1.5 flex items-center justify-between border-t border-slate-200/60 pt-1.5">
								<span class="text-[#6D7588]">GoPay:</span>
								<span class="font-bold text-[#00AA5B]">Terhubung (Rp350.000)</span>
							</div>
						</div>
					</div>

					<!-- Settings Navigation Menu -->
					<div class="rounded-2xl border border-[#E5E7E9] bg-white p-2 shadow-xs">
						<div class="space-y-1">
							{#each navItems as item}
								{@const isActive = currentPath === item.href || (item.href.includes('profile') && currentPath === '/user/settings')}
								<a
									href={item.href}
									class="flex items-center justify-between rounded-xl px-3.5 py-2.5 text-xs font-semibold transition-all {isActive
										? 'bg-[#F0FDF4] text-[#00AA5B]'
										: 'text-[#6D7588] hover:bg-slate-100 hover:text-[#212121]'}"
								>
									<div class="flex items-center gap-3">
										<item.icon class="h-4 w-4 {isActive ? 'text-[#00AA5B]' : 'text-[#6D7588]'}" />
										<span>{item.label}</span>
									</div>
									<ChevronRight class="h-3.5 w-3.5 opacity-60" />
								</a>
							{/each}
						</div>
					</div>
				</div>
			</div>

			<!-- RIGHT MAIN CONTENT (8 Cols) -->
			<div class="lg:col-span-8">
				<div class="rounded-2xl border border-[#E5E7E9] bg-white p-6 shadow-xs">
					{@render children()}
				</div>
			</div>
		</div>
	</div>
</div>

