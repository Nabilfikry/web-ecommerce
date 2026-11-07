<script lang="ts">
	import { toast } from '$lib/stores/toast.svelte';
	import { Plus, Search, MapPin, Check, Trash2, Edit3, Share2 } from 'lucide-svelte';

	interface Address {
		id: string;
		label: string;
		recipientName: string;
		phone: string;
		fullAddress: string;
		city: string;
		postalCode: string;
		isPrimary: boolean;
		hasPinpoint: boolean;
	}

	let addresses = $state<Address[]>([
		{
			id: 'addr-1',
			label: 'Rumah',
			recipientName: 'Naabil Fikry',
			phone: '0812-3456-7890',
			fullAddress:
				'Jl. Kebon Sirih No. 17A, RT 02 / RW 03, Kel. Gambir, Kec. Gambir, Jakarta Pusat',
			city: 'Jakarta Pusat',
			postalCode: '10110',
			isPrimary: true,
			hasPinpoint: true
		},
		{
			id: 'addr-2',
			label: 'Kantor',
			recipientName: 'Naabil Fikry (Office)',
			phone: '0812-9876-5432',
			fullAddress:
				'Sudirman Central Business District (SCBD), Gedung Menara Mandiri Lt. 14, Jl. Jend. Sudirman Kav. 54-55',
			city: 'Jakarta Selatan',
			postalCode: '12190',
			isPrimary: false,
			hasPinpoint: true
		}
	]);

	let searchQuery = $state('');
	let showModal = $state(false);

	// Form inputs
	let newLabel = $state('Rumah');
	let newRecipient = $state('');
	let newPhone = $state('');
	let newFullAddress = $state('');
	let newCity = $state('Jakarta Pusat');
	let newPostal = $state('');
	let newIsPrimary = $state(false);

	const filteredAddresses = $derived(
		addresses.filter(
			(a) =>
				a.recipientName.toLowerCase().includes(searchQuery.toLowerCase()) ||
				a.label.toLowerCase().includes(searchQuery.toLowerCase()) ||
				a.fullAddress.toLowerCase().includes(searchQuery.toLowerCase())
		)
	);

	function setPrimary(id: string) {
		addresses = addresses.map((a) => ({
			...a,
			isPrimary: a.id === id
		}));
		toast.success('Alamat utama berhasil diperbarui!');
	}

	function deleteAddress(id: string) {
		if (confirm('Hapus alamat ini dari daftar pengiriman?')) {
			addresses = addresses.filter((a) => a.id !== id);
			toast.info('Alamat berhasil dihapus.');
		}
	}

	function handleAddAddress(e: SubmitEvent) {
		e.preventDefault();
		if (!newRecipient.trim() || !newPhone.trim() || !newFullAddress.trim()) {
			toast.error('Mohon lengkapi seluruh data alamat');
			return;
		}

		const newAddr: Address = {
			id: `addr-${Date.now()}`,
			label: newLabel,
			recipientName: newRecipient.trim(),
			phone: newPhone.trim(),
			fullAddress: newFullAddress.trim(),
			city: newCity,
			postalCode: newPostal.trim() || '10110',
			isPrimary: newIsPrimary || addresses.length === 0,
			hasPinpoint: true
		};

		if (newAddr.isPrimary) {
			addresses = addresses.map((a) => ({ ...a, isPrimary: false }));
		}

		addresses = [newAddr, ...addresses];
		toast.success('Alamat baru berhasil ditambahkan!');

		// Reset form
		newRecipient = '';
		newPhone = '';
		newFullAddress = '';
		newPostal = '';
		showModal = false;
	}
</script>

<svelte:head>
	<title>Daftar Alamat | TokoKuPedia</title>
</svelte:head>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex flex-col justify-between gap-4 border-b border-[#E5E7E9] pb-4 sm:flex-row sm:items-center">
		<div>
			<h1 class="text-xl font-bold text-[#212121]">Daftar Alamat</h1>
			<p class="text-xs text-[#6D7588]">
				Kelola alamat pengiriman belanjaan Anda untuk kemudahan bertransaksi
			</p>
		</div>

		<button
			type="button"
			onclick={() => (showModal = true)}
			class="inline-flex cursor-pointer items-center justify-center gap-1.5 rounded-lg bg-[#00AA5B] px-4 py-2 text-xs font-bold text-white shadow-xs transition-colors hover:bg-[#029b53]"
		>
			<Plus class="h-4 w-4" />
			<span>Tambah Alamat Baru</span>
		</button>
	</div>

	<!-- Search & Filter Bar -->
	<div class="relative">
		<input
			type="search"
			bind:value={searchQuery}
			placeholder="Tulis nama penerima atau alamat..."
			class="h-10 w-full rounded-lg border border-[#E5E7E9] pr-4 pl-9 text-xs text-[#212121] placeholder:text-[#9FA6B0] focus:border-[#00AA5B] focus:ring-1 focus:ring-[#00AA5B] focus:outline-none"
		/>
		<Search class="pointer-events-none absolute top-3 left-3 h-4 w-4 text-[#9FA6B0]" />
	</div>

	<!-- Addresses List -->
	<div class="space-y-4">
		{#if filteredAddresses.length === 0}
			<div class="rounded-xl border border-dashed border-[#E5E7E9] p-8 text-center text-xs text-[#6D7588]">
				Tidak ada alamat yang sesuai dengan pencarian Anda.
			</div>
		{:else}
			{#each filteredAddresses as addr}
				<div
					class="rounded-xl border p-4.5 transition-all {addr.isPrimary
						? 'border-[#00AA5B] bg-[#F0FDF4]/30'
						: 'border-[#E5E7E9] bg-white'}"
				>
					<div class="flex items-start justify-between">
						<div class="flex items-center gap-2">
							<span class="rounded bg-slate-100 px-2 py-0.5 text-[11px] font-bold text-[#212121]">
								{addr.label}
							</span>
							{#if addr.isPrimary}
								<span class="rounded bg-emerald-100 px-2 py-0.5 text-[11px] font-bold text-[#00AA5B]">
									Utama
								</span>
							{/if}
						</div>

						<div class="flex items-center gap-2">
							<button
								type="button"
								onclick={() => toast.info('Link alamat disalin untuk dibagikan')}
								class="text-[#6D7588] hover:text-[#00AA5B]"
								title="Bagikan Alamat"
							>
								<Share2 class="h-4 w-4" />
							</button>
							<button
								type="button"
								onclick={() => deleteAddress(addr.id)}
								class="text-[#6D7588] hover:text-[#EE4D2D]"
								title="Hapus Alamat"
							>
								<Trash2 class="h-4 w-4" />
							</button>
						</div>
					</div>

					<div class="mt-2.5 space-y-1">
						<div class="text-sm font-bold text-[#212121]">
							{addr.recipientName}
							<span class="ml-2 font-normal text-xs text-[#6D7588]">{addr.phone}</span>
						</div>
						<p class="text-xs leading-relaxed text-[#212121]">
							{addr.fullAddress}
						</p>
						<div class="flex items-center gap-2 pt-1 text-[11px] text-[#6D7588]">
							<span>{addr.city}, {addr.postalCode}</span>
							{#if addr.hasPinpoint}
								<span class="flex items-center gap-1 font-semibold text-[#00AA5B]">
									<MapPin class="h-3 w-3" />
									Sudah Pinpoint
								</span>
							{/if}
						</div>
					</div>

					<div class="mt-4 flex items-center justify-between border-t border-[#E5E7E9] pt-3">
						<div class="flex items-center gap-3 text-xs">
							<button
								type="button"
								onclick={() => toast.info('Form ubah alamat dibuka')}
								class="font-bold text-[#00AA5B] hover:underline"
							>
								Ubah Alamat
							</button>
						</div>

						{#if !addr.isPrimary}
							<button
								type="button"
								onclick={() => setPrimary(addr.id)}
								class="rounded-lg border border-[#00AA5B] px-3 py-1 text-xs font-bold text-[#00AA5B] hover:bg-[#F0FDF4]"
							>
								Jadikan Alamat Utama
							</button>
						{/if}
					</div>
				</div>
			{/each}
		{/if}
	</div>
</div>

<!-- Modal Tambah Alamat Baru -->
{#if showModal}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
		<div class="w-full max-w-lg rounded-2xl border border-[#E5E7E9] bg-white p-6 shadow-2xl">
			<div class="flex items-center justify-between border-b border-[#E5E7E9] pb-3">
				<h3 class="text-base font-bold text-[#212121]">Tambah Alamat Baru</h3>
				<button
					type="button"
					onclick={() => (showModal = false)}
					class="text-[#6D7588] hover:text-[#212121] text-lg font-bold"
				>
					&times;
				</button>
			</div>

			<form onsubmit={handleAddAddress} class="mt-4 space-y-3.5">
				<div>
					<label for="addr-label" class="block text-xs font-semibold text-[#212121]">
						Label Alamat
					</label>
					<input
						id="addr-label"
						type="text"
						bind:value={newLabel}
						placeholder="Contoh: Rumah, Kantor, Kos"
						required
						class="mt-1 h-9 w-full rounded-lg border border-[#E5E7E9] px-3 text-xs text-[#212121] focus:border-[#00AA5B] focus:outline-none"
					/>
				</div>

				<div class="grid grid-cols-2 gap-3">
					<div>
						<label for="addr-name" class="block text-xs font-semibold text-[#212121]">
							Nama Penerima
						</label>
						<input
							id="addr-name"
							type="text"
							bind:value={newRecipient}
							placeholder="Nama penerima paket"
							required
							class="mt-1 h-9 w-full rounded-lg border border-[#E5E7E9] px-3 text-xs text-[#212121] focus:border-[#00AA5B] focus:outline-none"
						/>
					</div>

					<div>
						<label for="addr-phone" class="block text-xs font-semibold text-[#212121]">
							Nomor HP
						</label>
						<input
							id="addr-phone"
							type="tel"
							bind:value={newPhone}
							placeholder="08xxxxxxxxxx"
							required
							class="mt-1 h-9 w-full rounded-lg border border-[#E5E7E9] px-3 text-xs text-[#212121] focus:border-[#00AA5B] focus:outline-none"
						/>
					</div>
				</div>

				<div>
					<label for="addr-full" class="block text-xs font-semibold text-[#212121]">
						Alamat Lengkap
					</label>
					<textarea
						id="addr-full"
						rows="3"
						bind:value={newFullAddress}
						placeholder="Nama jalan, nomor gedung, RT/RW, kelurahan"
						required
						class="mt-1 w-full rounded-lg border border-[#E5E7E9] p-2.5 text-xs text-[#212121] focus:border-[#00AA5B] focus:outline-none"
					></textarea>
				</div>

				<div class="grid grid-cols-2 gap-3">
					<div>
						<label for="addr-city" class="block text-xs font-semibold text-[#212121]">
							Kota
						</label>
						<input
							id="addr-city"
							type="text"
							bind:value={newCity}
							required
							class="mt-1 h-9 w-full rounded-lg border border-[#E5E7E9] px-3 text-xs text-[#212121] focus:border-[#00AA5B] focus:outline-none"
						/>
					</div>
					<div>
						<label for="addr-post" class="block text-xs font-semibold text-[#212121]">
							Kode Pos
						</label>
						<input
							id="addr-post"
							type="text"
							bind:value={newPostal}
							placeholder="10110"
							class="mt-1 h-9 w-full rounded-lg border border-[#E5E7E9] px-3 text-xs text-[#212121] focus:border-[#00AA5B] focus:outline-none"
						/>
					</div>
				</div>

				<div class="flex items-center gap-2 pt-1">
					<input
						type="checkbox"
						id="isPrimary"
						bind:checked={newIsPrimary}
						class="h-4 w-4 rounded border-[#E5E7E9] text-[#00AA5B] focus:ring-[#00AA5B]"
					/>
					<label for="isPrimary" class="text-xs text-[#212121]">
						Jadikan sebagai alamat utama
					</label>
				</div>

				<div class="flex justify-end gap-2 pt-3 border-t border-[#E5E7E9]">
					<button
						type="button"
						onclick={() => (showModal = false)}
						class="rounded-lg border border-[#E5E7E9] px-4 py-2 text-xs font-semibold text-[#6D7588] hover:bg-slate-50"
					>
						Batal
					</button>
					<button
						type="submit"
						class="rounded-lg bg-[#00AA5B] px-5 py-2 text-xs font-bold text-white hover:bg-[#029b53]"
					>
						Simpan Alamat
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

