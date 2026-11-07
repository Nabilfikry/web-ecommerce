<script lang="ts">
	import { Search, ArrowUpDown, X } from 'lucide-svelte';

	interface Props {
		search?: string;
		sort?: 'latest' | 'price_asc' | 'price_desc';
		totalItems?: number;
		onchange?: (filters: { search: string; sort: 'latest' | 'price_asc' | 'price_desc' }) => void;
	}

	let {
		search = $bindable(''),
		sort = $bindable('latest'),
		totalItems = 0,
		onchange
	}: Props = $props();

	function handleSearchInput(e: Event) {
		const target = e.target as HTMLInputElement;
		search = target.value;
		triggerChange();
	}

	function handleSortChange(e: Event) {
		const target = e.target as HTMLSelectElement;
		sort = target.value as 'latest' | 'price_asc' | 'price_desc';
		triggerChange();
	}

	function handleClearSearch() {
		search = '';
		triggerChange();
	}

	function triggerChange() {
		if (onchange) {
			onchange({ search, sort });
		}
	}
</script>

<div
	class="flex flex-col items-center justify-between gap-3 rounded-xl border border-slate-200 bg-white p-4 shadow-xs sm:flex-row"
>
	<!-- Search Filter -->
	<div class="relative w-full sm:max-w-md">
		<Search
			class="pointer-events-none absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-slate-400"
		/>
		<input
			type="text"
			value={search}
			oninput={handleSearchInput}
			placeholder="Cari berdasarkan nama produk..."
			class="h-10 w-full rounded-lg border border-slate-300 bg-slate-50 pr-8 pl-9 text-xs text-slate-900 transition-colors placeholder:text-slate-400 focus:border-slate-900 focus:bg-white focus:ring-1 focus:ring-slate-900 focus:outline-none"
		/>
		{#if search}
			<button
				type="button"
				onclick={handleClearSearch}
				class="absolute top-1/2 right-2.5 -translate-y-1/2 cursor-pointer rounded p-1 text-slate-400 hover:text-slate-600"
				aria-label="Bersihkan pencarian"
			>
				<X class="h-3.5 w-3.5" />
			</button>
		{/if}
	</div>

	<!-- Sort Dropdown & Total Count -->
	<div class="flex w-full items-center justify-between gap-3 sm:w-auto sm:justify-end">
		<span class="text-xs whitespace-nowrap text-slate-500">
			Menampilkan <strong class="font-semibold text-slate-900">{totalItems}</strong> produk
		</span>

		<div class="flex items-center gap-1.5">
			<ArrowUpDown class="h-3.5 w-3.5 text-slate-400" />
			<select
				value={sort}
				onchange={handleSortChange}
				class="h-10 cursor-pointer rounded-lg border border-slate-300 bg-slate-50 px-2.5 text-xs text-slate-800 focus:border-slate-900 focus:outline-none"
			>
				<option value="latest">Terbaru</option>
				<option value="price_asc">Harga: Terendah</option>
				<option value="price_desc">Harga: Tertinggi</option>
			</select>
		</div>
	</div>
</div>
