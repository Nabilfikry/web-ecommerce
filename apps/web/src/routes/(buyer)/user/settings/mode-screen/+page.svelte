<script lang="ts">
	import { onMount } from 'svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import { Sun, Moon, Laptop, CheckCircle2 } from 'lucide-svelte';

	type ThemeMode = 'light' | 'dark' | 'system';

	let currentMode = $state<ThemeMode>('system');

	onMount(() => {
		const saved = (localStorage.getItem('tokokupedia_theme_mode') as ThemeMode) || 'system';
		currentMode = saved;
	});

	function applyTheme(mode: ThemeMode) {
		currentMode = mode;
		localStorage.setItem('tokokupedia_theme_mode', mode);

		if (
			mode === 'dark' ||
			(mode === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches)
		) {
			document.documentElement.classList.add('dark');
		} else {
			document.documentElement.classList.remove('dark');
		}

		const modeNames = {
			light: 'Mode Terang (Light Mode)',
			dark: 'Mode Gelap (Dark Mode)',
			system: 'Mode Otomatis (Ikuti Sistem)'
		};
		toast.success(`${modeNames[mode]} berhasil diterapkan!`);
	}
</script>

<svelte:head>
	<title>Mode Tampilan | TokoKuPedia</title>
</svelte:head>

<div class="space-y-6">
	<!-- Page Header -->
	<div class="border-b border-[#E5E7E9] pb-4">
		<h1 class="text-xl font-bold text-[#212121]">Mode Tampilan</h1>
		<p class="text-xs text-[#6D7588]">
			Sesuaikan tema visual aplikasi TokoKuPedia untuk kenyamanan mata Anda
		</p>
	</div>

	<!-- Theme Selection Cards Grid -->
	<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
		<!-- 1. LIGHT MODE -->
		<button
			type="button"
			onclick={() => applyTheme('light')}
			class="group relative flex cursor-pointer flex-col overflow-hidden rounded-2xl border-2 text-left transition-all p-4 {currentMode ===
			'light'
				? 'border-[#00AA5B] bg-[#F0FDF4]/30 ring-2 ring-[#00AA5B]/20'
				: 'border-[#E5E7E9] bg-white hover:border-slate-300'}"
		>
			<!-- Visual Screen Preview Mockup -->
			<div class="h-28 w-full overflow-hidden rounded-xl border border-slate-200 bg-white p-2.5 shadow-2xs">
				<div class="flex items-center justify-between border-b border-slate-100 pb-1.5">
					<div class="h-2 w-12 rounded bg-[#00AA5B]"></div>
					<div class="h-2 w-4 rounded-full bg-slate-300"></div>
				</div>
				<div class="mt-2 space-y-1.5">
					<div class="h-2.5 w-3/4 rounded bg-slate-200"></div>
					<div class="h-2 w-1/2 rounded bg-slate-100"></div>
				</div>
				<div class="mt-3 grid grid-cols-2 gap-1.5">
					<div class="h-8 rounded bg-slate-50 border border-slate-100"></div>
					<div class="h-8 rounded bg-slate-50 border border-slate-100"></div>
				</div>
			</div>

			<div class="mt-4 flex items-center justify-between">
				<div class="flex items-center gap-2">
					<Sun class="h-4 w-4 text-amber-500" />
					<span class="text-sm font-bold text-[#212121]">Terang</span>
				</div>
				{#if currentMode === 'light'}
					<span class="rounded-full bg-[#00AA5B] p-0.5 text-white">
						<CheckCircle2 class="h-4 w-4" />
					</span>
				{/if}
			</div>

			<p class="mt-1.5 text-[11px] leading-relaxed text-[#6D7588]">
				Tampilan putih cerah dengan kontras tinggi yang nyaman di siang hari.
			</p>
		</button>

		<!-- 2. DARK MODE -->
		<button
			type="button"
			onclick={() => applyTheme('dark')}
			class="group relative flex cursor-pointer flex-col overflow-hidden rounded-2xl border-2 text-left transition-all p-4 {currentMode ===
			'dark'
				? 'border-[#00AA5B] bg-[#F0FDF4]/30 ring-2 ring-[#00AA5B]/20'
				: 'border-[#E5E7E9] bg-white hover:border-slate-300'}"
		>
			<!-- Visual Screen Preview Mockup -->
			<div class="h-28 w-full overflow-hidden rounded-xl border border-slate-700 bg-[#121212] p-2.5 shadow-2xs">
				<div class="flex items-center justify-between border-b border-slate-800 pb-1.5">
					<div class="h-2 w-12 rounded bg-[#00AA5B]"></div>
					<div class="h-2 w-4 rounded-full bg-slate-600"></div>
				</div>
				<div class="mt-2 space-y-1.5">
					<div class="h-2.5 w-3/4 rounded bg-slate-700"></div>
					<div class="h-2 w-1/2 rounded bg-slate-800"></div>
				</div>
				<div class="mt-3 grid grid-cols-2 gap-1.5">
					<div class="h-8 rounded bg-slate-900 border border-slate-800"></div>
					<div class="h-8 rounded bg-slate-900 border border-slate-800"></div>
				</div>
			</div>

			<div class="mt-4 flex items-center justify-between">
				<div class="flex items-center gap-2">
					<Moon class="h-4 w-4 text-indigo-400" />
					<span class="text-sm font-bold text-[#212121]">Gelap</span>
				</div>
				{#if currentMode === 'dark'}
					<span class="rounded-full bg-[#00AA5B] p-0.5 text-white">
						<CheckCircle2 class="h-4 w-4" />
					</span>
				{/if}
			</div>

			<p class="mt-1.5 text-[11px] leading-relaxed text-[#6D7588]">
				Tampilan gelap elegan untuk kenyamanan mata di ruangan redup dan hemat daya baterai.
			</p>
		</button>

		<!-- 3. SYSTEM DEFAULT -->
		<button
			type="button"
			onclick={() => applyTheme('system')}
			class="group relative flex cursor-pointer flex-col overflow-hidden rounded-2xl border-2 text-left transition-all p-4 {currentMode ===
			'system'
				? 'border-[#00AA5B] bg-[#F0FDF4]/30 ring-2 ring-[#00AA5B]/20'
				: 'border-[#E5E7E9] bg-white hover:border-slate-300'}"
		>
			<!-- Visual Screen Preview Mockup (Split Light & Dark) -->
			<div class="relative h-28 w-full overflow-hidden rounded-xl border border-slate-300 p-2.5 shadow-2xs">
				<div class="absolute inset-0 grid grid-cols-2">
					<div class="bg-white p-2">
						<div class="h-2 w-10 rounded bg-[#00AA5B]"></div>
						<div class="mt-2 h-2 w-3/4 rounded bg-slate-200"></div>
					</div>
					<div class="bg-[#121212] p-2">
						<div class="h-2 w-10 rounded bg-[#00AA5B]"></div>
						<div class="mt-2 h-2 w-3/4 rounded bg-slate-700"></div>
					</div>
				</div>
			</div>

			<div class="mt-4 flex items-center justify-between">
				<div class="flex items-center gap-2">
					<Laptop class="h-4 w-4 text-blue-500" />
					<span class="text-sm font-bold text-[#212121]">Ikuti Sistem</span>
				</div>
				{#if currentMode === 'system'}
					<span class="rounded-full bg-[#00AA5B] p-0.5 text-white">
						<CheckCircle2 class="h-4 w-4" />
					</span>
				{/if}
			</div>

			<p class="mt-1.5 text-[11px] leading-relaxed text-[#6D7588]">
				Menyesuaikan tema otomatis dengan pengaturan sistem operasi perangkat Anda.
			</p>
		</button>
	</div>
</div>

