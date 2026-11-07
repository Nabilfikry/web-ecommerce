<script lang="ts">
	import type { Snippet } from 'svelte';
	import { authStore } from '$lib/stores/auth.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { ShieldAlert, ShoppingBag } from 'lucide-svelte';

	interface Props {
		children?: Snippet;
	}

	let { children }: Props = $props();

	// Guard verification
	let isSeller = $derived(authStore.isAuthenticated && authStore.user?.role === 'seller');
</script>

{#if !authStore.isAuthenticated}
	<div class="mx-auto max-w-md px-4 py-24 text-center">
		<div
			class="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl border border-amber-200 bg-amber-50 text-amber-600"
		>
			<ShieldAlert class="h-7 w-7" />
		</div>
		<h2 class="mt-4 text-xl font-bold text-slate-900">Akses Dibatasi</h2>
		<p class="mt-2 text-sm text-slate-600">
			Halaman ini khusus untuk Mitra Penjual. Silakan masuk terlebih dahulu dengan akun Penjual
			Anda.
		</p>
		<div class="mt-6 flex flex-col gap-2.5 sm:flex-row sm:justify-center">
			<Button variant="primary" href="/login">Masuk ke Akun</Button>
			<Button variant="secondary" href="/register">Daftar Akun Penjual</Button>
		</div>
	</div>
{:else if !isSeller}
	<div class="mx-auto max-w-md px-4 py-24 text-center">
		<div
			class="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl border border-red-200 bg-red-50 text-red-600"
		>
			<ShieldAlert class="h-7 w-7" />
		</div>
		<h2 class="mt-4 text-xl font-bold text-slate-900">Peran Akun Tidak Sesuai</h2>
		<p class="mt-2 text-sm text-slate-600">
			Anda saat ini masuk sebagai <strong>Pembeli</strong> ({authStore.user?.email}). Dashboard ini
			memerlukan akun dengan role <strong>Penjual (Seller)</strong>.
		</p>
		<div class="mt-6 flex flex-col gap-2.5 sm:flex-row sm:justify-center">
			<Button variant="primary" href="/products">
				<ShoppingBag class="h-4 w-4" />
				Belanja di Katalog
			</Button>
			<Button
				variant="secondary"
				onclick={() => {
					authStore.logout();
					window.location.href = '/login';
				}}
			>
				Ganti Akun
			</Button>
		</div>
	</div>
{:else}
	<div class="min-h-screen bg-slate-50/60 pb-16">
		{@render children?.()}
	</div>
{/if}
