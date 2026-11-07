<script lang="ts">
	import { toast } from '$lib/stores/toast.svelte';
	import { Building2, Plus, Trash2, CheckCircle2, ShieldAlert } from 'lucide-svelte';

	interface BankAccount {
		id: string;
		bankName: string;
		accountNumber: string;
		holderName: string;
		isPrimary: boolean;
	}

	let bankAccounts = $state<BankAccount[]>([
		{
			id: 'bank-1',
			bankName: 'Bank Central Asia (BCA)',
			accountNumber: '123-456-7890',
			holderName: 'NAABIL FIKRY',
			isPrimary: true
		},
		{
			id: 'bank-2',
			bankName: 'Bank Mandiri',
			accountNumber: '137-00-1234567-8',
			holderName: 'NAABIL FIKRY',
			isPrimary: false
		}
	]);

	let showModal = $state(false);
	let newBank = $state('Bank Central Asia (BCA)');
	let newAccNumber = $state('');
	let newHolderName = $state('NAABIL FIKRY');

	function handleAddBank(e: SubmitEvent) {
		e.preventDefault();
		if (bankAccounts.length >= 3) {
			toast.error('Maksimal hanya 3 rekening bank yang dapat didaftarkan');
			return;
		}

		if (!newAccNumber.trim()) {
			toast.error('Masukkan nomor rekening bank');
			return;
		}

		const newAcc: BankAccount = {
			id: `bank-${Date.now()}`,
			bankName: newBank,
			accountNumber: newAccNumber.trim(),
			holderName: newHolderName.trim().toUpperCase(),
			isPrimary: bankAccounts.length === 0
		};

		bankAccounts = [...bankAccounts, newAcc];
		toast.success('Rekening bank berhasil ditambahkan!');
		newAccNumber = '';
		showModal = false;
	}

	function removeBank(id: string) {
		if (confirm('Hapus nomor rekening bank ini?')) {
			bankAccounts = bankAccounts.filter((b) => b.id !== id);
			toast.info('Rekening bank berhasil dihapus.');
		}
	}

	function setPrimaryBank(id: string) {
		bankAccounts = bankAccounts.map((b) => ({
			...b,
			isPrimary: b.id === id
		}));
		toast.success('Rekening penarikan utama berhasil diperbarui!');
	}
</script>

<svelte:head>
	<title>Rekening Bank | TokoKuPedia</title>
</svelte:head>

<div class="space-y-6">
	<!-- Page Header -->
	<div class="flex flex-col justify-between gap-4 border-b border-[#E5E7E9] pb-4 sm:flex-row sm:items-center">
		<div>
			<h1 class="text-xl font-bold text-[#212121]">Rekening Bank</h1>
			<p class="text-xs text-[#6D7588]">
				Rekening bank terdaftar untuk penarikan dana saldo TokoKuPedia atau pencairan refund
			</p>
		</div>

		{#if bankAccounts.length < 3}
			<button
				type="button"
				onclick={() => (showModal = true)}
				class="inline-flex cursor-pointer items-center justify-center gap-1.5 rounded-lg bg-[#00AA5B] px-4 py-2 text-xs font-bold text-white shadow-xs transition-colors hover:bg-[#029b53]"
			>
				<Plus class="h-4 w-4" />
				<span>Tambah Rekening</span>
			</button>
		{/if}
	</div>

	<!-- Information Alert Notice -->
	<div class="flex items-start gap-3 rounded-xl border border-emerald-100 bg-[#F0FDF4] p-3.5 text-xs text-[#212121]">
		<Building2 class="h-5 w-5 shrink-0 text-[#00AA5B]" />
		<div class="space-y-0.5 leading-relaxed">
			<span class="font-bold">Maksimal 3 rekening bank yang dapat didaftarkan.</span>
			<p class="text-[11px] text-[#6D7588]">
				Demi keamanan dana Anda, pastikan nama pemilik rekening bank sama dengan nama pada identitas KTP yang terdaftar di akun TokoKuPedia.
			</p>
		</div>
	</div>

	<!-- Bank Accounts List -->
	<div class="space-y-3.5">
		{#each bankAccounts as acc}
			<div
				class="flex items-center justify-between rounded-xl border p-4.5 transition-all {acc.isPrimary
					? 'border-[#00AA5B] bg-[#F0FDF4]/30'
					: 'border-[#E5E7E9] bg-white'}"
			>
				<div class="flex items-center gap-4">
					<div
						class="flex h-11 w-11 items-center justify-center rounded-xl bg-slate-100 font-extrabold text-[#00AA5B]"
					>
						<Building2 class="h-6 w-6" />
					</div>

					<div class="space-y-0.5">
						<div class="flex items-center gap-2">
							<span class="text-sm font-bold text-[#212121]">{acc.bankName}</span>
							{#if acc.isPrimary}
								<span class="rounded bg-emerald-100 px-2 py-0.5 text-[10px] font-bold text-[#00AA5B]">
									Rekening Utama
								</span>
							{/if}
						</div>
						<div class="font-mono text-xs font-bold text-[#212121]">
							{acc.accountNumber}
						</div>
						<div class="text-[11px] text-[#6D7588]">
							a.n. <strong class="text-[#212121]">{acc.holderName}</strong>
						</div>
					</div>
				</div>

				<div class="flex items-center gap-3">
					{#if !acc.isPrimary}
						<button
							type="button"
							onclick={() => setPrimaryBank(acc.id)}
							class="text-xs font-bold text-[#00AA5B] hover:underline"
						>
							Jadikan Utama
						</button>
					{/if}

					<button
						type="button"
						onclick={() => removeBank(acc.id)}
						class="text-[#6D7588] hover:text-[#EE4D2D]"
						title="Hapus Rekening"
					>
						<Trash2 class="h-4 w-4" />
					</button>
				</div>
			</div>
		{/each}
	</div>
</div>

<!-- Modal Tambah Rekening -->
{#if showModal}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
		<div class="w-full max-w-md rounded-2xl border border-[#E5E7E9] bg-white p-6 shadow-2xl">
			<div class="flex items-center justify-between border-b border-[#E5E7E9] pb-3">
				<h3 class="text-base font-bold text-[#212121]">Tambah Rekening Bank</h3>
				<button
					type="button"
					onclick={() => (showModal = false)}
					class="text-[#6D7588] hover:text-[#212121] text-lg font-bold"
				>
					&times;
				</button>
			</div>

			<form onsubmit={handleAddBank} class="mt-4 space-y-3.5">
				<div>
					<label for="bank-select" class="block text-xs font-semibold text-[#212121]">
						Pilih Bank
					</label>
					<select
						id="bank-select"
						bind:value={newBank}
						class="mt-1 h-9 w-full rounded-lg border border-[#E5E7E9] px-3 text-xs text-[#212121] focus:border-[#00AA5B] focus:outline-none"
					>
						<option value="Bank Central Asia (BCA)">Bank Central Asia (BCA)</option>
						<option value="Bank Mandiri">Bank Mandiri</option>
						<option value="Bank Rakyat Indonesia (BRI)">Bank Rakyat Indonesia (BRI)</option>
						<option value="Bank Negara Indonesia (BNI)">Bank Negara Indonesia (BNI)</option>
						<option value="Bank CIMB Niaga">Bank CIMB Niaga</option>
						<option value="Bank Syariah Indonesia (BSI)">Bank Syariah Indonesia (BSI)</option>
						<option value="Bank Danamon">Bank Danamon</option>
					</select>
				</div>

				<div>
					<label for="acc-num" class="block text-xs font-semibold text-[#212121]">
						Nomor Rekening
					</label>
					<input
						id="acc-num"
						type="text"
						bind:value={newAccNumber}
						placeholder="Contoh: 1234567890"
						required
						class="mt-1 h-9 w-full rounded-lg border border-[#E5E7E9] px-3 font-mono text-xs text-[#212121] focus:border-[#00AA5B] focus:outline-none"
					/>
				</div>

				<div>
					<label for="acc-holder" class="block text-xs font-semibold text-[#212121]">
						Nama Pemilik Rekening (Sesuai KTP)
					</label>
					<input
						id="acc-holder"
						type="text"
						bind:value={newHolderName}
						required
						class="mt-1 h-9 w-full rounded-lg border border-[#E5E7E9] px-3 text-xs uppercase text-[#212121] focus:border-[#00AA5B] focus:outline-none"
					/>
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
						Simpan Rekening
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

