<script lang="ts">
	import { toast } from '$lib/stores/toast.svelte';
	import { CreditCard, Wallet, ShieldCheck, Plus, Trash2, CheckCircle2 } from 'lucide-svelte';

	interface SavedCard {
		id: string;
		type: 'Visa' | 'Mastercard';
		bank: string;
		maskedNumber: string;
		expiry: string;
		isDefault: boolean;
	}

	let savedCards = $state<SavedCard[]>([
		{
			id: 'card-1',
			type: 'Visa',
			bank: 'BCA',
			maskedNumber: '•••• •••• •••• 4242',
			expiry: '08/28',
			isDefault: true
		},
		{
			id: 'card-2',
			type: 'Mastercard',
			bank: 'Bank Mandiri',
			maskedNumber: '•••• •••• •••• 8812',
			expiry: '12/26',
			isDefault: false
		}
	]);

	let showAddCardModal = $state(false);
	let newCardBank = $state('BCA');
	let newCardNumber = $state('');
	let newCardExpiry = $state('');
	let newCardHolder = $state('');

	function handleAddCard(e: SubmitEvent) {
		e.preventDefault();
		if (!newCardNumber.trim() || !newCardExpiry.trim() || !newCardHolder.trim()) {
			toast.error('Lengkapi nomor kartu, masa berlaku, dan nama pemilik');
			return;
		}

		const cleanNum = newCardNumber.replace(/\s+/g, '');
		const last4 = cleanNum.slice(-4) || '9999';

		const newCard: SavedCard = {
			id: `card-${Date.now()}`,
			type: cleanNum.startsWith('4') ? 'Visa' : 'Mastercard',
			bank: newCardBank,
			maskedNumber: `•••• •••• •••• ${last4}`,
			expiry: newCardExpiry,
			isDefault: savedCards.length === 0
		};

		savedCards = [...savedCards, newCard];
		toast.success('Kartu debit/kredit berhasil ditambahkan!');
		newCardNumber = '';
		newCardExpiry = '';
		newCardHolder = '';
		showAddCardModal = false;
	}

	function removeCard(id: string) {
		if (confirm('Hapus kartu ini dari metode pembayaran tersimpan?')) {
			savedCards = savedCards.filter((c) => c.id !== id);
			toast.info('Kartu berhasil dihapus.');
		}
	}

	function setDefaultCard(id: string) {
		savedCards = savedCards.map((c) => ({
			...c,
			isDefault: c.id === id
		}));
		toast.success('Kartu pembayaran utama diperbarui!');
	}
</script>

<svelte:head>
	<title>Metode Pembayaran | TokoKuPedia</title>
</svelte:head>

<div class="space-y-6">
	<!-- Page Header -->
	<div class="border-b border-[#E5E7E9] pb-4">
		<h1 class="text-xl font-bold text-[#212121]">Metode Pembayaran</h1>
		<p class="text-xs text-[#6D7588]">
			Kelola metode pembayaran instan, kartu tersimpan, dan dompet digital di TokoKuPedia
		</p>
	</div>

	<!-- SECTION 1: Dompet Digital & Instan -->
	<div class="space-y-3">
		<h2 class="text-sm font-bold text-[#212121]">Dompet Digital & Saldo Instan</h2>

		<div class="space-y-2.5">
			<!-- GoPay -->
			<div class="flex items-center justify-between rounded-xl border border-[#E5E7E9] bg-white p-4">
				<div class="flex items-center gap-3">
					<div class="flex h-10 w-10 items-center justify-center rounded-lg bg-[#00AED6]/10 font-black text-[#00AED6] text-xs">
						GoPay
					</div>
					<div>
						<div class="flex items-center gap-2">
							<span class="text-xs font-bold text-[#212121]">GoPay</span>
							<span class="rounded bg-emerald-100 px-1.5 py-0.2 text-[10px] font-bold text-[#00AA5B]">
								Terhubung
							</span>
						</div>
						<div class="text-xs text-[#6D7588]">
							Saldo: <strong class="text-[#212121]">Rp350.000</strong>
						</div>
					</div>
				</div>

				<button
					type="button"
					onclick={() => toast.info('Menu pengelolaan akun GoPay')}
					class="rounded-lg border border-[#00AA5B] px-3 py-1.5 text-xs font-bold text-[#00AA5B] hover:bg-[#F0FDF4]"
				>
					Kelola
				</button>
			</div>

			<!-- TokoKuPedia Saldo -->
			<div class="flex items-center justify-between rounded-xl border border-[#E5E7E9] bg-white p-4">
				<div class="flex items-center gap-3">
					<div class="flex h-10 w-10 items-center justify-center rounded-lg bg-[#00AA5B]/10 text-[#00AA5B]">
						<Wallet class="h-5 w-5" />
					</div>
					<div>
						<div class="flex items-center gap-2">
							<span class="text-xs font-bold text-[#212121]">Saldo TokoKuPedia</span>
							<span class="rounded bg-emerald-100 px-1.5 py-0.2 text-[10px] font-bold text-[#00AA5B]">
								Aktif
							</span>
						</div>
						<div class="text-xs text-[#6D7588]">
							Saldo: <strong class="text-[#212121]">Rp1.500.000</strong>
						</div>
					</div>
				</div>

				<div class="flex gap-2">
					<button
						type="button"
						onclick={() => toast.success('Form pengisian saldo TokoKuPedia dibuka')}
						class="rounded-lg bg-[#00AA5B] px-3 py-1.5 text-xs font-bold text-white hover:bg-[#029b53]"
					>
						Isi Saldo
					</button>
				</div>
			</div>

			<!-- OVO -->
			<div class="flex items-center justify-between rounded-xl border border-[#E5E7E9] bg-white p-4">
				<div class="flex items-center gap-3">
					<div class="flex h-10 w-10 items-center justify-center rounded-lg bg-[#4C3494]/10 font-black text-[#4C3494] text-xs">
						OVO
					</div>
					<div>
						<div class="text-xs font-bold text-[#212121]">OVO Cash & Points</div>
						<div class="text-xs text-[#6D7588]">Hubungkan untuk kemudahan transaksi 1-klik</div>
					</div>
				</div>

				<button
					type="button"
					onclick={() => toast.info('Fitur penghubungan OVO dikirim ke nomor HP')}
					class="rounded-lg border border-[#00AA5B] px-3 py-1.5 text-xs font-bold text-[#00AA5B] hover:bg-[#F0FDF4]"
				>
					Hubungkan
				</button>
			</div>
		</div>
	</div>

	<!-- SECTION 2: Kartu Debit / Kredit Tersimpan -->
	<div class="space-y-3 pt-4 border-t border-[#E5E7E9]">
		<div class="flex items-center justify-between">
			<h2 class="text-sm font-bold text-[#212121]">Kartu Debit / Kredit Tersimpan</h2>
			<button
				type="button"
				onclick={() => (showAddCardModal = true)}
				class="inline-flex items-center gap-1 text-xs font-bold text-[#00AA5B] hover:underline"
			>
				<Plus class="h-3.5 w-3.5" />
				<span>Tambah Kartu</span>
			</button>
		</div>

		<div class="space-y-3">
			{#each savedCards as card}
				<div class="flex items-center justify-between rounded-xl border border-[#E5E7E9] bg-white p-4">
					<div class="flex items-center gap-3.5">
						<div class="flex h-10 w-14 items-center justify-center rounded border border-[#E5E7E9] bg-slate-50 font-bold text-xs text-[#212121]">
							{card.type}
						</div>
						<div>
							<div class="flex items-center gap-2">
								<span class="text-xs font-bold text-[#212121]">{card.bank} ({card.type})</span>
								{#if card.isDefault}
									<span class="rounded bg-emerald-100 px-1.5 py-0.2 text-[10px] font-bold text-[#00AA5B]">
										Utama
									</span>
								{/if}
							</div>
							<div class="font-mono text-xs text-[#6D7588]">{card.maskedNumber} • Exp {card.expiry}</div>
						</div>
					</div>

					<div class="flex items-center gap-3">
						{#if !card.isDefault}
							<button
								type="button"
								onclick={() => setDefaultCard(card.id)}
								class="text-xs font-bold text-[#00AA5B] hover:underline"
							>
								Jadikan Utama
							</button>
						{/if}
						<button
							type="button"
							onclick={() => removeCard(card.id)}
							class="text-[#6D7588] hover:text-[#EE4D2D]"
							title="Hapus Kartu"
						>
							<Trash2 class="h-4 w-4" />
						</button>
					</div>
				</div>
			{/each}
		</div>
	</div>

	<!-- SECTION 3: Virtual Account Favorit -->
	<div class="space-y-3 pt-4 border-t border-[#E5E7E9]">
		<h2 class="text-sm font-bold text-[#212121]">Virtual Account Favorit</h2>
		<div class="grid grid-cols-2 gap-3 sm:grid-cols-4">
			<div class="rounded-xl border border-[#E5E7E9] bg-white p-3 text-center">
				<div class="font-bold text-xs text-blue-700">BCA VA</div>
				<div class="mt-1 text-[11px] text-[#6D7588]">Otomatis Terverifikasi</div>
			</div>
			<div class="rounded-xl border border-[#E5E7E9] bg-white p-3 text-center">
				<div class="font-bold text-xs text-amber-700">Mandiri VA</div>
				<div class="mt-1 text-[11px] text-[#6D7588]">Otomatis Terverifikasi</div>
			</div>
			<div class="rounded-xl border border-[#E5E7E9] bg-white p-3 text-center">
				<div class="font-bold text-xs text-blue-600">BRI VA (BRIVA)</div>
				<div class="mt-1 text-[11px] text-[#6D7588]">Otomatis Terverifikasi</div>
			</div>
			<div class="rounded-xl border border-[#E5E7E9] bg-white p-3 text-center">
				<div class="font-bold text-xs text-orange-600">BNI VA</div>
				<div class="mt-1 text-[11px] text-[#6D7588]">Otomatis Terverifikasi</div>
			</div>
		</div>
	</div>
</div>

<!-- Modal Tambah Kartu Baru -->
{#if showAddCardModal}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
		<div class="w-full max-w-md rounded-2xl border border-[#E5E7E9] bg-white p-6 shadow-2xl">
			<div class="flex items-center justify-between border-b border-[#E5E7E9] pb-3">
				<h3 class="text-base font-bold text-[#212121]">Tambah Kartu Debit / Kredit</h3>
				<button
					type="button"
					onclick={() => (showAddCardModal = false)}
					class="text-[#6D7588] hover:text-[#212121] text-lg font-bold"
				>
					&times;
				</button>
			</div>

			<form onsubmit={handleAddCard} class="mt-4 space-y-3">
				<div>
					<label for="card-bank" class="block text-xs font-semibold text-[#212121]">
						Nama Bank Penerbit
					</label>
					<select
						id="card-bank"
						bind:value={newCardBank}
						class="mt-1 h-9 w-full rounded-lg border border-[#E5E7E9] px-3 text-xs text-[#212121] focus:border-[#00AA5B] focus:outline-none"
					>
						<option value="BCA">Bank BCA</option>
						<option value="Bank Mandiri">Bank Mandiri</option>
						<option value="BRI">Bank BRI</option>
						<option value="BNI">Bank BNI</option>
						<option value="CIMB Niaga">Bank CIMB Niaga</option>
					</select>
				</div>

				<div>
					<label for="card-num" class="block text-xs font-semibold text-[#212121]">
						Nomor Kartu (16 Digit)
					</label>
					<input
						id="card-num"
						type="text"
						maxlength="19"
						bind:value={newCardNumber}
						placeholder="4xxx xxxx xxxx xxxx"
						required
						class="mt-1 h-9 w-full rounded-lg border border-[#E5E7E9] px-3 font-mono text-xs text-[#212121] focus:border-[#00AA5B] focus:outline-none"
					/>
				</div>

				<div class="grid grid-cols-2 gap-3">
					<div>
						<label for="card-exp" class="block text-xs font-semibold text-[#212121]">
							Masa Berlaku (MM/YY)
						</label>
						<input
							id="card-exp"
							type="text"
							maxlength="5"
							bind:value={newCardExpiry}
							placeholder="12/28"
							required
							class="mt-1 h-9 w-full rounded-lg border border-[#E5E7E9] px-3 text-xs text-[#212121] focus:border-[#00AA5B] focus:outline-none"
						/>
					</div>

					<div>
						<label for="card-cvv" class="block text-xs font-semibold text-[#212121]">
							CVV
						</label>
						<input
							id="card-cvv"
							type="password"
							maxlength="3"
							placeholder="123"
							required
							class="mt-1 h-9 w-full rounded-lg border border-[#E5E7E9] px-3 text-xs text-[#212121] focus:border-[#00AA5B] focus:outline-none"
						/>
					</div>
				</div>

				<div>
					<label for="card-holder" class="block text-xs font-semibold text-[#212121]">
						Nama Pemilik Kartu
					</label>
					<input
						id="card-holder"
						type="text"
						bind:value={newCardHolder}
						placeholder="Sesuai nama tertera di kartu"
						required
						class="mt-1 h-9 w-full rounded-lg border border-[#E5E7E9] px-3 text-xs text-[#212121] focus:border-[#00AA5B] focus:outline-none"
					/>
				</div>

				<div class="flex justify-end gap-2 pt-3 border-t border-[#E5E7E9]">
					<button
						type="button"
						onclick={() => (showAddCardModal = false)}
						class="rounded-lg border border-[#E5E7E9] px-4 py-2 text-xs font-semibold text-[#6D7588] hover:bg-slate-50"
					>
						Batal
					</button>
					<button
						type="submit"
						class="rounded-lg bg-[#00AA5B] px-5 py-2 text-xs font-bold text-white hover:bg-[#029b53]"
					>
						Simpan Kartu
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

