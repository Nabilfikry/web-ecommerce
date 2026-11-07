<script lang="ts">
	import { page } from '$app/state';
	import Button from '$lib/components/ui/Button.svelte';
	import { ArrowLeft, Home, ShoppingBag, AlertCircle } from 'lucide-svelte';

	const status = $derived(page.status);
	const message = $derived(
		page.error?.message || 'Terjadi kendala yang tidak terduga pada sistem.'
	);
	const is404 = $derived(status === 404);
</script>

<svelte:head>
	<title>{is404 ? '404 - Halaman Tidak Ditemukan' : `${status} - Terjadi Kesalahan`}</title>
</svelte:head>

<div class="flex min-h-[75vh] items-center justify-center px-4 py-16 sm:px-6 lg:px-8">
	<div class="w-full max-w-md text-center">
		<!-- Icon Container -->
		<div
			class="mx-auto flex h-16 w-16 items-center justify-center rounded-2xl {is404
				? 'border border-amber-200 bg-amber-50 text-amber-600'
				: 'border border-red-200 bg-red-50 text-red-600'}"
		>
			<AlertCircle class="h-8 w-8" />
		</div>

		<!-- Status Code Badge -->
		<div class="mt-4">
			<span
				class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-bold tracking-wider uppercase {is404
					? 'bg-amber-100 text-amber-800'
					: 'bg-red-100 text-red-800'}"
			>
				Error {status}
			</span>
		</div>

		<!-- Title -->
		<h1 class="mt-3 text-2xl font-bold tracking-tight text-slate-900 sm:text-3xl">
			{is404 ? 'Halaman Tidak Ditemukan' : 'Terjadi Kesalahan Sistem'}
		</h1>

		<!-- Message -->
		<p class="mt-2 text-sm leading-relaxed text-slate-600">
			{#if is404}
				Halaman atau tautan yang Anda tuju mungkin sudah dipindahkan, dihapus, atau alamat URL salah
				ketik.
			{:else}
				{message}
			{/if}
		</p>

		<!-- Action Buttons -->
		<div class="mt-8 flex flex-col gap-2.5 sm:flex-row sm:justify-center">
			<Button variant="primary" href="/">
				<Home class="h-4 w-4" />
				Kembali ke Beranda
			</Button>

			<Button variant="secondary" href="/products">
				<ShoppingBag class="h-4 w-4" />
				Katalog Produk
			</Button>
		</div>

		<!-- Back trigger -->
		<div class="mt-6">
			<button
				type="button"
				onclick={() => history.back()}
				class="inline-flex items-center gap-1.5 text-xs font-semibold text-slate-500 transition-colors hover:text-slate-800"
			>
				<ArrowLeft class="h-3.5 w-3.5" />
				Kembali ke halaman sebelumnya
			</button>
		</div>
	</div>
</div>
