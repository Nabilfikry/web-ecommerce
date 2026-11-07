<script lang="ts">
	import { cart } from '$lib/stores/cart.svelte';
	import { auth } from '$lib/stores/auth.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import { formatIDR } from '$lib/utils/format';
	import { ordersApi } from '$lib/api/orders';
	import { goto } from '$app/navigation';
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import {
		ShieldCheck,
		ArrowLeft,
		Store,
		AlertTriangle,
		CheckCircle,
		ArrowRight
	} from 'lucide-svelte';

	let recipientName = $state(auth.user?.name || '');
	let phone = $state('');
	let address = $state('');
	let city = $state('Bandung');
	let province = $state('Jawa Barat');
	let postalCode = $state('40257');
	let notes = $state('');

	let isSubmitting = $state(false);
	let error = $state<string | null>(null);

	$effect(() => {
		if (auth.user?.name && !recipientName) {
			recipientName = auth.user.name;
		}
	});

	const sellerName = $derived(cart.items[0]?.product?.seller?.name || 'Toko Penjual');

	async function handleCheckout(e: SubmitEvent) {
		e.preventDefault();
		error = null;

		if (!auth.isAuthenticated) {
			toast.warning('Silakan masuk terlebih dahulu dengan akun pembeli untuk melakukan checkout.');
			goto('/login?redirect=/checkout');
			return;
		}

		if (auth.user?.role !== 'buyer') {
			error = 'Akun penjual tidak dapat melakukan pembelian. Harap gunakan akun pembeli.';
			toast.error(error);
			return;
		}

		if (cart.items.length === 0) {
			toast.error('Keranjang belanja Anda kosong.');
			goto('/products');
			return;
		}

		if (cart.hasMultipleSellers) {
			error =
				'Satu pesanan hanya boleh memuat barang dari satu penjual. Mohon periksa kembali keranjang Anda.';
			toast.error(error);
			return;
		}

		if (
			!recipientName.trim() ||
			!phone.trim() ||
			!address.trim() ||
			!city.trim() ||
			!postalCode.trim()
		) {
			error = 'Seluruh data alamat pengiriman bertanda bintang (*) wajib dilengkapi.';
			return;
		}

		isSubmitting = true;
		try {
			const order = await ordersApi.checkout({
				shipping_address: {
					recipient_name: recipientName.trim(),
					phone: phone.trim(),
					address: address.trim(),
					city: city.trim(),
					province: province.trim(),
					postal_code: postalCode.trim()
				},
				notes: notes.trim() || undefined
			});

			// Clear cart after successful transaction
			await cart.clearCart();
			toast.success('Pesanan berhasil dibuat! Stok produk telah dikunci.');
			goto(`/orders/${order.id}`);
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : 'Checkout gagal diproses oleh sistem.';
			error = msg;
			toast.error(msg);
		} finally {
			isSubmitting = false;
		}
	}
</script>

<svelte:head>
	<title>Checkout Pesanan — SimpleShop</title>
</svelte:head>

<div class="mx-auto flex max-w-7xl flex-1 flex-col gap-6 px-4 py-8 sm:px-6 lg:px-8">
	<!-- Navigation Header -->
	<nav class="flex items-center gap-2 text-xs text-slate-500">
		<a href="/cart" class="flex items-center gap-1 transition-colors hover:text-slate-900">
			<ArrowLeft class="h-3.5 w-3.5" />
			<span>Kembali ke Keranjang</span>
		</a>
		<span>/</span>
		<span class="font-medium text-slate-900">Checkout Pembayaran</span>
	</nav>

	<div class="flex flex-col gap-1 border-b border-slate-200 pb-4">
		<h1 class="text-2xl font-bold tracking-tight text-slate-900 sm:text-3xl">
			Checkout & Informasi Pengiriman
		</h1>
		<p class="text-xs text-slate-500 sm:text-sm">
			Lengkapi alamat penerima dan konfirmasi rincian pesanan Anda.
		</p>
	</div>

	<!-- Auth Status Callout for Guest -->
	{#if !auth.isLoading && !auth.isAuthenticated}
		<div
			class="flex flex-col gap-3 rounded-xl border border-blue-200 bg-blue-50 p-4 text-xs text-blue-900 sm:flex-row sm:items-center sm:justify-between"
		>
			<div class="flex items-center gap-2">
				<ShieldCheck class="h-4 w-4 shrink-0 text-blue-600" />
				<span
					>Anda belum masuk ke akun pembeli. Silakan masuk untuk menyelesaikan pesanan Anda.</span
				>
			</div>
			<a
				href="/login?redirect=/checkout"
				class="inline-block shrink-0 rounded-lg bg-blue-600 px-3.5 py-1.5 text-center font-semibold text-white shadow-xs transition-colors hover:bg-blue-700"
			>
				Masuk Akun Pembeli
			</a>
		</div>
	{/if}

	<!-- Seller Warning if logged in as seller -->
	{#if !auth.isLoading && auth.isSeller}
		<div
			class="flex items-center gap-2 rounded-xl border border-amber-200 bg-amber-50 p-4 text-xs text-amber-900"
		>
			<AlertTriangle class="h-4 w-4 shrink-0 text-amber-600" />
			<span>
				Anda sedang masuk sebagai <strong>Seller</strong>. Checkout hanya dapat diproses menggunakan
				akun <strong>Buyer (Pembeli)</strong>.
			</span>
		</div>
	{/if}

	{#if error}
		<div
			class="flex items-center gap-2 rounded-xl border border-red-200 bg-red-50 p-4 text-xs text-red-800"
		>
			<AlertTriangle class="h-4 w-4 shrink-0 text-red-600" />
			<span>{error}</span>
		</div>
	{/if}

	{#if cart.items.length === 0}
		<!-- Empty Cart State in Checkout -->
		<div
			class="my-4 flex flex-col items-center justify-center gap-3 rounded-xl border border-slate-200 bg-white py-16 text-center"
		>
			<h2 class="text-lg font-bold text-slate-900">Keranjang Belanja Kosong</h2>
			<p class="max-w-sm text-xs text-slate-500">
				Pilih produk terlebih dahulu sebelum dapat melanjutkan ke proses checkout.
			</p>
			<a
				href="/products"
				class="mt-2 inline-flex items-center gap-2 rounded-lg bg-slate-900 px-5 py-2.5 text-xs font-semibold text-white shadow-xs transition-colors hover:bg-slate-800"
			>
				<span>Mulai Belanja</span>
				<ArrowRight class="h-3.5 w-3.5" />
			</a>
		</div>
	{:else}
		<!-- Checkout Form & Summary Grid -->
		<form onsubmit={handleCheckout} class="grid grid-cols-1 items-start gap-8 lg:grid-cols-3">
			<!-- Left: Shipping Address & Notes (2 cols) -->
			<div class="space-y-6 lg:col-span-2">
				<!-- Seller Info Card -->
				<div
					class="flex items-center justify-between rounded-xl border border-slate-200 bg-white p-5"
				>
					<div class="flex items-center gap-2.5">
						<div
							class="flex h-9 w-9 items-center justify-center rounded-lg bg-emerald-50 text-emerald-700"
						>
							<Store class="h-5 w-5" />
						</div>
						<div>
							<span class="block text-[10px] font-medium tracking-wider text-slate-400 uppercase"
								>Toko Penjual</span
							>
							<span class="block text-sm font-bold text-slate-900">{sellerName}</span>
						</div>
					</div>

					<span class="rounded-full bg-emerald-50 px-2.5 py-1 text-xs font-medium text-emerald-600">
						Single-Seller Checkout
					</span>
				</div>

				<!-- Shipping Address Details -->
				<div class="space-y-4 rounded-xl border border-slate-200 bg-white p-6 shadow-xs">
					<h3
						class="border-b border-slate-100 pb-3 text-sm font-bold tracking-wider text-slate-900 uppercase"
					>
						Alamat Tujuan Pengiriman
					</h3>

					<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
						<Input
							label="Nama Penerima"
							placeholder="contoh: Naabil Fikry"
							bind:value={recipientName}
							required
						/>

						<Input
							label="Nomor Telepon / WhatsApp"
							type="tel"
							placeholder="contoh: 081234567890"
							bind:value={phone}
							required
						/>
					</div>

					<Input
						label="Alamat Lengkap"
						placeholder="Nama jalan, nomor rumah, RT/RW, kelurahan, kecamatan"
						bind:value={address}
						required
					/>

					<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
						<Input
							label="Kota / Kabupaten"
							placeholder="contoh: Bandung"
							bind:value={city}
							required
						/>

						<Input label="Provinsi" placeholder="contoh: Jawa Barat" bind:value={province} />

						<Input label="Kode Pos" placeholder="contoh: 40257" bind:value={postalCode} required />
					</div>

					<div class="pt-2">
						<label for="order-notes" class="mb-1.5 block text-xs font-semibold text-slate-800">
							Catatan untuk Penjual (Opsional)
						</label>
						<textarea
							id="order-notes"
							rows="2"
							bind:value={notes}
							placeholder="contoh: Mohon bubble wrap ekstra tebal ya min..."
							class="w-full rounded-lg border border-slate-300 bg-white p-3 text-xs text-slate-900 transition-colors placeholder:text-slate-400 focus:border-slate-900 focus:ring-1 focus:ring-slate-900 focus:outline-none"
						></textarea>
					</div>
				</div>
			</div>

			<!-- Right: Order Items Snapshot & Total (1 col) -->
			<div
				class="sticky top-20 space-y-5 rounded-xl border border-slate-200 bg-white p-6 shadow-xs"
			>
				<h3
					class="border-b border-slate-100 pb-3 text-sm font-bold tracking-wider text-slate-900 uppercase"
				>
					Ringkasan Pesanan ({cart.totalItems} barang)
				</h3>

				<!-- Items list snapshot -->
				<div class="max-h-60 space-y-3 overflow-y-auto pr-1">
					{#each cart.items as item (item.id)}
						<div class="flex items-center justify-between gap-3 text-xs">
							<div class="flex-1 truncate">
								<span class="block truncate font-semibold text-slate-900">{item.product?.name}</span
								>
								<span class="text-[11px] text-slate-400"
									>{item.quantity}x @ {formatIDR(item.product?.price || 0)}</span
								>
							</div>
							<span class="shrink-0 font-bold text-slate-900 tabular-nums">
								{formatIDR(item.subtotal)}
							</span>
						</div>
					{/each}
				</div>

				<div class="space-y-2.5 border-t border-slate-100 pt-4 text-xs text-slate-600">
					<div class="flex items-center justify-between">
						<span>Subtotal Barang</span>
						<span class="font-semibold text-slate-900 tabular-nums"
							>{formatIDR(cart.totalAmount)}</span
						>
					</div>
					<div class="flex items-center justify-between">
						<span>Biaya Pengiriman</span>
						<span class="font-semibold text-emerald-600">Gratis (MVP)</span>
					</div>
					<div class="flex items-center justify-between">
						<span>Metode Pembayaran</span>
						<span class="font-semibold text-blue-600">DOKU Payment Gateway</span>
					</div>
					<div class="flex items-baseline justify-between border-t border-slate-100 pt-3">
						<span class="text-sm font-bold text-slate-900">Total Tagihan</span>
						<span class="text-xl font-extrabold tracking-tight text-slate-900 tabular-nums">
							{formatIDR(cart.totalAmount)}
						</span>
					</div>
				</div>

				<div class="pt-2">
					<Button
						type="submit"
						variant="primary"
						loading={isSubmitting}
						disabled={isSubmitting ||
							cart.items.length === 0 ||
							cart.hasMultipleSellers ||
							auth.isSeller}
						class="flex w-full items-center justify-center gap-2"
					>
						<CheckCircle class="h-4 w-4" />
						<span>Konfirmasi & Buat Pesanan</span>
					</Button>
				</div>

				<div class="flex items-center justify-center gap-2 pt-2 text-[11px] text-slate-400">
					<ShieldCheck class="h-4 w-4 shrink-0 text-emerald-600" />
					<span>Transaksi diproses atomik melalui PostgreSQL ACID</span>
				</div>
			</div>
		</form>
	{/if}
</div>
