<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import TokoKuPediaLogo from '$lib/components/ui/TokoKuPediaLogo.svelte';
	import TopedMascot from '$lib/components/ui/TopedMascot.svelte';
	import { Store, ShoppingBag, Loader2 } from 'lucide-svelte';

	let email = $state('');
	let password = $state('');
	let error = $state<string | null>(null);
	let isSubmitting = $state(false);

	const redirectTo = $derived(page.url.searchParams.get('redirect') || '');

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		error = null;

		if (!email.trim() || !password) {
			error = 'Email dan password wajib diisi';
			return;
		}

		isSubmitting = true;
		try {
			const res = await auth.login({ email: email.trim(), password });
			toast.success(`Selamat datang kembali, ${res.user.name}!`);

			if (redirectTo) {
				goto(redirectTo);
			} else if (res.user.role === 'seller') {
				goto('/dashboard');
			} else {
				goto('/');
			}
		} catch (err: unknown) {
			const msg =
				err instanceof Error ? err.message : 'Login gagal. Periksa kembali email dan password.';
			error = msg;
			toast.error(msg);
		} finally {
			isSubmitting = false;
		}
	}

	function fillDemoAccount(role: 'buyer' | 'seller') {
		if (role === 'buyer') {
			email = 'buyer@example.com';
			password = 'Password123';
		} else {
			email = 'seller@techgear.com';
			password = 'Password123';
		}
	}

	function handleGoogleAuth() {
		email = 'buyer@example.com';
		password = 'Password123';
		toast.info('Akun Google uji coba siap masuk.');
	}
</script>

<svelte:head>
	<title>Masuk / Login | TokoKuPedia</title>
</svelte:head>

<div class="flex min-h-screen flex-col bg-white">
	<!-- Header Logo -->
	<header class="flex items-center justify-center py-8">
		<a href="/" class="group">
			<TokoKuPediaLogo size="lg" />
		</a>
	</header>

	<!-- Main Content -->
	<main class="mx-auto flex w-full max-w-5xl flex-1 items-center justify-center px-4 py-4 sm:px-6">
		<div class="grid w-full grid-cols-1 items-center gap-12 lg:grid-cols-2">
			<!-- LEFT COLUMN: Toped Mascot -->
			<div class="hidden flex-col items-center text-center lg:flex">
				<TopedMascot size="lg" />
				<h2 class="mt-6 text-2xl font-bold tracking-tight text-[#212121]">
					Jual Beli Mudah Hanya di TokoKuPedia
				</h2>
				<p class="mt-2 text-sm text-[#6D7588]">
					Gabung dan nikmati jutaan promo menarik setiap hari
				</p>
			</div>

			<!-- RIGHT COLUMN: Login Card -->
			<div class="flex justify-center">
				<div
					class="w-full max-w-md rounded-2xl border border-[#E5E7E9] bg-white p-6 shadow-xl sm:p-8"
				>
					<div class="mb-6 flex items-center justify-between">
						<h1 class="text-2xl font-bold tracking-tight text-[#212121]">Masuk</h1>
						<a href="/register" class="text-sm font-bold text-[#00AA5B] hover:underline">
							Daftar
						</a>
					</div>

					<!-- Google Button -->
					<button
						type="button"
						onclick={handleGoogleAuth}
						class="flex h-11 w-full cursor-pointer items-center justify-center gap-2.5 rounded-lg border border-[#E5E7E9] bg-white px-4 text-sm font-semibold text-[#212121] transition-all hover:bg-slate-50 active:scale-98"
					>
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
						<div>
							<label for="login-email" class="block text-xs font-semibold text-[#212121]">
								Nomor HP atau E-mail
							</label>
							<input
								id="login-email"
								type="email"
								bind:value={email}
								placeholder="Nomor HP atau E-mail"
								required
								class="mt-1 h-10 w-full rounded-lg border border-[#E5E7E9] px-3 text-sm text-[#212121] placeholder:text-[#9FA6B0] transition-all focus:border-[#00AA5B] focus:ring-1 focus:ring-[#00AA5B] focus:outline-none"
							/>
						</div>

						<div>
							<div class="flex items-center justify-between">
								<label for="login-password" class="block text-xs font-semibold text-[#212121]">
									Kata Sandi
								</label>
								<a href="#lupa-password" class="text-[11px] font-bold text-[#00AA5B] hover:underline">
									Lupa kata sandi?
								</a>
							</div>
							<input
								id="login-password"
								type="password"
								bind:value={password}
								placeholder="Masukkan kata sandi"
								required
								class="mt-1 h-10 w-full rounded-lg border border-[#E5E7E9] px-3 text-sm text-[#212121] placeholder:text-[#9FA6B0] transition-all focus:border-[#00AA5B] focus:ring-1 focus:ring-[#00AA5B] focus:outline-none"
							/>
						</div>

						<div class="pt-2">
							<button
								type="submit"
								disabled={isSubmitting}
								class="flex h-11 w-full cursor-pointer items-center justify-center rounded-lg bg-[#00AA5B] text-sm font-bold text-white shadow-xs transition-all hover:bg-[#029b53] active:scale-98 disabled:opacity-50"
							>
								{#if isSubmitting}
									<Loader2 class="mr-2 h-4 w-4 animate-spin" />
								{/if}
								<span>Masuk</span>
							</button>
						</div>
					</form>

					<!-- Quick Demo Credentials Fill -->
					<div class="mt-6 rounded-xl border border-slate-100 bg-[#F5F6F7] p-3 text-xs text-[#6D7588]">
						<div class="mb-2 font-bold text-[#212121]">Akun Demo Pengujian:</div>
						<div class="flex gap-2">
							<button
								type="button"
								onclick={() => fillDemoAccount('buyer')}
								class="flex-1 rounded-md border border-[#E5E7E9] bg-white py-1.5 font-medium text-[#212121] hover:bg-slate-50"
							>
								Pembeli (Buyer)
							</button>
							<button
								type="button"
								onclick={() => fillDemoAccount('seller')}
								class="flex-1 rounded-md border border-[#E5E7E9] bg-white py-1.5 font-medium text-[#212121] hover:bg-slate-50"
							>
								Penjual (Seller)
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>
	</main>

	<!-- Footer -->
	<footer class="py-6 text-center text-xs text-[#6D7588]">
		<span>&copy; 2009-2026, PT TokoKuPedia | </span>
		<a href="/#care" class="font-bold text-[#00AA5B] hover:underline">TokoKuPedia Care</a>
	</footer>
</div>
