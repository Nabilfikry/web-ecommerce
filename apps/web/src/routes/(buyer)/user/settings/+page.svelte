<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import { User, Mail, Phone, Calendar, ShieldCheck } from 'lucide-svelte';

	let name = $state(auth.user?.name || 'Naabil Fikry');
	let email = $state(auth.user?.email || 'buyer@tokokupedia.com');
	let phone = $state('0812-3456-7890');
	let birthDate = $state('1998-05-15');
	let gender = $state<'pria' | 'wanita'>('pria');

	function handleSaveProfile(e: SubmitEvent) {
		e.preventDefault();
		toast.success('Biodata diri berhasil diperbarui!');
	}
</script>

<svelte:head>
	<title>Biodata Diri | TokoKuPedia</title>
</svelte:head>

<div class="space-y-6">
	<!-- Page Header -->
	<div class="border-b border-[#E5E7E9] pb-4">
		<h1 class="text-xl font-bold text-[#212121]">Biodata Diri</h1>
		<p class="text-xs text-[#6D7588]">
			Kelola informasi profil akun Anda untuk keamanan dan kenyamanan bertransaksi di TokoKuPedia
		</p>
	</div>

	<form onsubmit={handleSaveProfile} class="max-w-lg space-y-4">
		<div>
			<label for="bio-name" class="block text-xs font-semibold text-[#212121]">Nama Lengkap</label>
			<input
				id="bio-name"
				type="text"
				bind:value={name}
				required
				class="mt-1 h-10 w-full rounded-lg border border-[#E5E7E9] px-3 text-xs text-[#212121] focus:border-[#00AA5B] focus:outline-none"
			/>
		</div>

		<div>
			<label for="bio-email" class="block text-xs font-semibold text-[#212121]">Alamat E-mail</label>
			<input
				id="bio-email"
				type="email"
				bind:value={email}
				required
				class="mt-1 h-10 w-full rounded-lg border border-[#E5E7E9] px-3 text-xs text-[#212121] focus:border-[#00AA5B] focus:outline-none"
			/>
		</div>

		<div>
			<label for="bio-phone" class="block text-xs font-semibold text-[#212121]">Nomor HP</label>
			<input
				id="bio-phone"
				type="tel"
				bind:value={phone}
				required
				class="mt-1 h-10 w-full rounded-lg border border-[#E5E7E9] px-3 text-xs text-[#212121] focus:border-[#00AA5B] focus:outline-none"
			/>
		</div>

		<div class="grid grid-cols-2 gap-3">
			<div>
				<label for="bio-date" class="block text-xs font-semibold text-[#212121]">Tanggal Lahir</label>
				<input
					id="bio-date"
					type="date"
					bind:value={birthDate}
					class="mt-1 h-10 w-full rounded-lg border border-[#E5E7E9] px-3 text-xs text-[#212121] focus:border-[#00AA5B] focus:outline-none"
				/>
			</div>

			<div>
				<span class="block text-xs font-semibold text-[#212121]">Jenis Kelamin</span>
				<div class="mt-1.5 flex items-center gap-4 text-xs">
					<label class="flex items-center gap-1.5 cursor-pointer">
						<input
							type="radio"
							name="gender"
							value="pria"
							checked={gender === 'pria'}
							onchange={() => (gender = 'pria')}
							class="text-[#00AA5B] focus:ring-[#00AA5B]"
						/>
						<span>Pria</span>
					</label>
					<label class="flex items-center gap-1.5 cursor-pointer">
						<input
							type="radio"
							name="gender"
							value="wanita"
							checked={gender === 'wanita'}
							onchange={() => (gender = 'wanita')}
							class="text-[#00AA5B] focus:ring-[#00AA5B]"
						/>
						<span>Wanita</span>
					</label>
				</div>
			</div>
		</div>

		<div class="pt-4 border-t border-[#E5E7E9]">
			<button
				type="submit"
				class="rounded-lg bg-[#00AA5B] px-6 py-2.5 text-xs font-bold text-white shadow-xs hover:bg-[#029b53]"
			>
				Simpan Perubahan
			</button>
		</div>
	</form>
</div>

