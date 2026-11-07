<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import { goto } from '$app/navigation';
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import type { UserRole } from '@ecommerce/shared-types';
	import { ShoppingBag, Store } from 'lucide-svelte';

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
</script>

<svelte:head>
	<title>Daftar Akun Baru — SimpleShop</title>
</svelte:head>

<div class="flex flex-1 items-center justify-center p-4 sm:p-6 lg:p-8">
	<div class="w-full max-w-md rounded-xl border border-slate-200 bg-white p-6 shadow-sm sm:p-8">
		<div class="mb-6 text-center">
			<h1 class="text-2xl font-bold tracking-tight text-slate-900">Buat Akun Baru</h1>
			<p class="mt-1 text-xs text-slate-500">
				Pilih peran akun dan mulai gunakan SimpleShop Platform.
			</p>
		</div>

		{#if error}
			<div class="mb-5 rounded-lg border border-red-200 bg-red-50 p-3 text-xs text-red-700">
				{error}
			</div>
		{/if}

		<form onsubmit={handleSubmit} class="space-y-4">
			<!-- Role Selector -->
			<div class="space-y-1.5">
				<span class="block text-xs font-semibold text-slate-800">Peran Akun (Role)</span>
				<div class="grid grid-cols-2 gap-2.5">
					<button
						type="button"
						onclick={() => (role = 'buyer')}
						class="flex cursor-pointer flex-col gap-1 rounded-lg border p-3 text-left transition-all {role ===
						'buyer'
							? 'border-slate-900 bg-slate-50 ring-1 ring-slate-900'
							: 'border-slate-200 hover:border-slate-300'}"
					>
						<div class="flex items-center gap-1.5 text-xs font-semibold text-slate-900">
							<ShoppingBag class="h-4 w-4 text-blue-600" />
							Pembeli (Buyer)
						</div>
						<p class="text-[11px] leading-tight text-slate-500">Cari, belanja, & lacak pesanan</p>
					</button>

					<button
						type="button"
						onclick={() => (role = 'seller')}
						class="flex cursor-pointer flex-col gap-1 rounded-lg border p-3 text-left transition-all {role ===
						'seller'
							? 'border-slate-900 bg-slate-50 ring-1 ring-slate-900'
							: 'border-slate-200 hover:border-slate-300'}"
					>
						<div class="flex items-center gap-1.5 text-xs font-semibold text-slate-900">
							<Store class="h-4 w-4 text-emerald-600" />
							Penjual (Seller)
						</div>
						<p class="text-[11px] leading-tight text-slate-500">Kelola produk & proses pesanan</p>
					</button>
				</div>
			</div>

			<Input
				label="Nama Lengkap"
				type="text"
				placeholder="contoh: Naabil Fikry"
				bind:value={name}
				required
				autocomplete="name"
			/>

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
				helperText="Minimal 8 karakter untuk keamanan akun."
				required
				autocomplete="new-password"
			/>

			<Button type="submit" variant="primary" loading={isSubmitting} class="w-full">
				Daftar Sebagai {role === 'seller' ? 'Penjual' : 'Pembeli'}
			</Button>
		</form>

		<div class="mt-6 text-center text-xs text-slate-500">
			Sudah memiliki akun?
			<a href="/login" class="font-semibold text-slate-900 hover:underline"> Masuk di sini </a>
		</div>
	</div>
</div>
