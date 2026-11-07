<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import { Store, ShoppingBag } from 'lucide-svelte';

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
				goto('/products');
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
</script>

<svelte:head>
	<title>Masuk Akun — SimpleShop</title>
</svelte:head>

<div class="flex flex-1 items-center justify-center p-4 sm:p-6 lg:p-8">
	<div class="w-full max-w-md rounded-xl border border-slate-200 bg-white p-6 shadow-sm sm:p-8">
		<div class="mb-6 text-center">
			<h1 class="text-2xl font-bold tracking-tight text-slate-900">Masuk ke Akun</h1>
			<p class="mt-1 text-xs text-slate-500">
				Masukkan kredensial akun untuk melanjutkan transaksi atau kelola toko.
			</p>
		</div>

		{#if error}
			<div class="mb-5 rounded-lg border border-red-200 bg-red-50 p-3 text-xs text-red-700">
				{error}
			</div>
		{/if}

		<form onsubmit={handleSubmit} class="space-y-4">
			<Input
				label="Alamat Email"
				type="email"
				placeholder="contoh: nama@email.com"
				bind:value={email}
				required
				autocomplete="email"
			/>

			<Input
				label="Kata Sandi (Password)"
				type="password"
				placeholder="Minimal 8 karakter"
				bind:value={password}
				required
				autocomplete="current-password"
			/>

			<Button type="submit" variant="primary" loading={isSubmitting} class="w-full">
				Masuk Sekarang
			</Button>
		</form>

		<!-- Demo Credentials Helper -->
		<div class="mt-6 border-t border-slate-100 pt-5">
			<p
				class="mb-2.5 text-center text-[11px] font-semibold tracking-wider text-slate-400 uppercase"
			>
				Uji Coba Cepat (Demo Preset)
			</p>
			<div class="grid grid-cols-2 gap-2">
				<button
					type="button"
					onclick={() => fillDemoAccount('buyer')}
					class="flex cursor-pointer items-center justify-center gap-1.5 rounded-lg border border-slate-200 p-2 text-xs text-slate-700 transition-colors hover:bg-slate-50"
				>
					<ShoppingBag class="h-3.5 w-3.5 text-blue-600" />
					<span>Sebagai Buyer</span>
				</button>
				<button
					type="button"
					onclick={() => fillDemoAccount('seller')}
					class="flex cursor-pointer items-center justify-center gap-1.5 rounded-lg border border-slate-200 p-2 text-xs text-slate-700 transition-colors hover:bg-slate-50"
				>
					<Store class="h-3.5 w-3.5 text-emerald-600" />
					<span>Sebagai Seller</span>
				</button>
			</div>
		</div>

		<div class="mt-6 text-center text-xs text-slate-500">
			Belum memiliki akun?
			<a href="/register" class="font-semibold text-slate-900 hover:underline"> Daftar di sini </a>
		</div>
	</div>
</div>
