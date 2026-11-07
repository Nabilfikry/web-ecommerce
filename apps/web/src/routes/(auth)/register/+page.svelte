<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import { goto } from '$app/navigation';
	import TokoKuPediaLogo from '$lib/components/ui/TokoKuPediaLogo.svelte';
	import TopedMascot from '$lib/components/ui/TopedMascot.svelte';
	import type { UserRole } from '@ecommerce/shared-types';
	import { ShoppingBag, Store, Loader2 } from 'lucide-svelte';

	let name = $state('');
	let email = $state('');
	let password = $state('');
	let role = $state<UserRole>('buyer');
	let error = $state<string | null>(null);
	let isSubmitting = $state(false);

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		error = null;

		if (!name.trim() || !email.trim() || !password) {
			error = 'Seluruh kolom wajib diisi';
			return;
		}

		if (password.length < 8) {
			error = 'Password minimal terdiri dari 8 karakter';
			return;
		}

		isSubmitting = true;
		try {
			await auth.register({
				name: name.trim(),
				email: email.trim(),
				password,
				role
			});
			toast.success('Pendaftaran akun berhasil! Silakan masuk.');
			goto('/login');
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : 'Pendaftaran gagal. Silakan coba kembali.';
			error = msg;
			toast.error(msg);
		} finally {
			isSubmitting = false;
		}
	}

	function handleGoogleAuth() {
		toast.info('Autentikasi Google SDK disimulasikan untuk lingkungan pengujian.');
		name = 'Pengguna TokoKuPedia';
		email = 'user.google@tokokupedia.com';
		password = 'Password123';
	}
</script>

<svelte:head>
	<title>Daftar Sekarang | TokoKuPedia</title>
</svelte:head>

<div class="flex min-h-screen flex-col bg-white">
	<!-- Centered Tokopedia Logo Header -->
	<header class="flex items-center justify-center py-8">
		<a href="/" class="group">
			<TokoKuPediaLogo size="lg" />
		</a>
	</header>

	<!-- Main 2-Column Content -->
	<main class="mx-auto flex w-full max-w-5xl flex-1 items-center justify-center px-4 py-4 sm:px-6">
		<div class="grid w-full grid-cols-1 items-center gap-12 lg:grid-cols-2">
			<!-- LEFT COLUMN: Toped Mascot & Slogan -->
			<div class="hidden flex-col items-center text-center lg:flex">
				<TopedMascot size="lg" />
				<h2 class="mt-6 text-2xl font-bold tracking-tight text-[#212121]">
					Jual Beli Mudah Hanya di TokoKuPedia
				</h2>
				<p class="mt-2 text-sm text-[#6D7588]">
					Gabung dan rasakan kemudahan bertransaksi di TokoKuPedia
				</p>
			</div>

			<!-- RIGHT COLUMN: Registration Card -->
			<div class="flex justify-center">
				<div
					class="w-full max-w-md rounded-2xl border border-[#E5E7E9] bg-white p-6 shadow-xl sm:p-8"
				>
					<div class="mb-6 text-center">
						<h1 class="text-2xl font-bold tracking-tight text-[#212121]">Daftar Sekarang</h1>
						<p class="mt-1 text-xs text-[#6D7588]">
							Sudah punya akun TokoKuPedia?
							<a href="/login" class="font-bold text-[#00AA5B] hover:underline">Masuk</a>
						</p>
					</div>

					<!-- Google Sign In Button -->
					<button
						type="button"
						onclick={handleGoogleAuth}
						class="flex h-11 w-full cursor-pointer items-center justify-center gap-2.5 rounded-lg border border-[#E5E7E9] bg-white px-4 text-sm font-semibold text-[#212121] transition-all hover:bg-slate-50 active:scale-98"
					>
						<!-- Google G Logo SVG -->
						<svg class="h-4 w-4" viewBox="0 0 24 24">
							<path
								fill="#4285F4"
								d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
							/>
							<path
								fill="#34A853"
								d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
							/>
							<path
								fill="#FBBC05"
								d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.06H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.94l2.85-2.22.81-.63z"
							/>
							<path
								fill="#EA4335"
								d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.06l3.66 2.84c.87-2.6 3.3-4.52 6.16-4.52z"
							/>
						</svg>
						<span>Google</span>
					</button>

					<!-- Divider -->
					<div class="relative my-6 flex items-center justify-center">
						<div class="w-full border-t border-[#E5E7E9]"></div>
						<span class="absolute bg-white px-3 text-xs text-[#6D7588]">atau</span>
					</div>

					{#if error}
						<div class="mb-4 rounded-lg bg-red-50 p-3 text-xs text-[#EE4D2D]">
							{error}
						</div>
					{/if}

					<!-- Form -->
					<form onsubmit={handleSubmit} class="space-y-4">
						<!-- Name -->
						<div>
							<label for="reg-name" class="block text-xs font-semibold text-[#212121]">
								Nama Lengkap
							</label>
							<input
								id="reg-name"
								type="text"
								bind:value={name}
								placeholder="Contoh: Naabil Fikry"
								required
								class="mt-1 h-10 w-full rounded-lg border border-[#E5E7E9] px-3 text-sm text-[#212121] transition-all placeholder:text-[#9FA6B0] focus:border-[#00AA5B] focus:ring-1 focus:ring-[#00AA5B] focus:outline-none"
							/>
						</div>

						<!-- Email / Phone -->
						<div>
							<label for="reg-email" class="block text-xs font-semibold text-[#212121]">
								Nomor HP atau E-mail
							</label>
							<input
								id="reg-email"
								type="email"
								bind:value={email}
								placeholder="Nomor HP atau E-mail"
								required
								class="mt-1 h-10 w-full rounded-lg border border-[#E5E7E9] px-3 text-sm text-[#212121] transition-all placeholder:text-[#9FA6B0] focus:border-[#00AA5B] focus:ring-1 focus:ring-[#00AA5B] focus:outline-none"
							/>
							<span class="mt-1 block text-[11px] text-[#6D7588]">
								Contoh: email@tokokupedia.com
							</span>
						</div>

						<!-- Password -->
						<div>
							<label for="reg-password" class="block text-xs font-semibold text-[#212121]">
								Password
							</label>
							<input
								id="reg-password"
								type="password"
								bind:value={password}
								placeholder="Minimal 8 karakter"
								required
								class="mt-1 h-10 w-full rounded-lg border border-[#E5E7E9] px-3 text-sm text-[#212121] transition-all placeholder:text-[#9FA6B0] focus:border-[#00AA5B] focus:ring-1 focus:ring-[#00AA5B] focus:outline-none"
							/>
						</div>

						<!-- Role Selection (Buyer / Seller) -->
						<div class="pt-1">
							<span class="block text-xs font-semibold text-[#212121]">Tipe Akun</span>
							<div class="mt-1.5 grid grid-cols-2 gap-2">
								<button
									type="button"
									onclick={() => (role = 'buyer')}
									class="flex cursor-pointer items-center justify-center gap-1.5 rounded-lg border p-2 text-xs font-semibold transition-all {role ===
									'buyer'
										? 'border-[#00AA5B] bg-[#F0FDF4] text-[#00AA5B]'
										: 'border-[#E5E7E9] text-[#6D7588] hover:border-slate-300'}"
								>
									<ShoppingBag class="h-3.5 w-3.5" />
									<span>Pembeli</span>
								</button>
								<button
									type="button"
									onclick={() => (role = 'seller')}
									class="flex cursor-pointer items-center justify-center gap-1.5 rounded-lg border p-2 text-xs font-semibold transition-all {role ===
									'seller'
										? 'border-[#00AA5B] bg-[#F0FDF4] text-[#00AA5B]'
										: 'border-[#E5E7E9] text-[#6D7588] hover:border-slate-300'}"
								>
									<Store class="h-3.5 w-3.5" />
									<span>Penjual</span>
								</button>
							</div>
						</div>

						<!-- Submit Button -->
						<div class="pt-2">
							<button
								type="submit"
								disabled={isSubmitting}
								class="flex h-11 w-full cursor-pointer items-center justify-center rounded-lg bg-[#00AA5B] text-sm font-bold text-white shadow-xs transition-all hover:bg-[#029b53] active:scale-98 disabled:opacity-50"
							>
								{#if isSubmitting}
									<Loader2 class="mr-2 h-4 w-4 animate-spin" />
								{/if}
								<span>Daftar</span>
							</button>
						</div>

						<!-- Terms Disclaimer -->
						<p class="pt-2 text-center text-[11px] leading-relaxed text-[#6D7588]">
							Dengan mendaftar, saya menyetujui
							<a href="/#terms" class="font-bold text-[#00AA5B] hover:underline">
								Syarat & Ketentuan
							</a>
							serta
							<a href="/#privacy" class="font-bold text-[#00AA5B] hover:underline">
								Kebijakan Privasi TokoKuPedia
							</a>.
						</p>
					</form>
				</div>
			</div>
		</div>
	</main>

	<!-- Footer Bar -->
	<footer class="py-6 text-center text-xs text-[#6D7588]">
		<span>&copy; 2009-2026, PT TokoKuPedia | </span>
		<a href="/#care" class="font-bold text-[#00AA5B] hover:underline">TokoKuPedia Care</a>
	</footer>
</div>
