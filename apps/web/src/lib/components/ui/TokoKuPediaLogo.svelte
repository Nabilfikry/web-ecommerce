<script lang="ts">
	interface Props {
		size?: 'sm' | 'md' | 'lg' | 'xl';
		showText?: boolean;
		class?: string;
	}

	let { size = 'md', showText = true, class: className = '' }: Props = $props();

	const sizes = {
		sm: { icon: 'w-6 h-6', text: 'text-lg', badge: 'text-[9px]' },
		md: { icon: 'w-8 h-8', text: 'text-2xl', badge: 'text-[10px]' },
		lg: { icon: 'w-10 h-10', text: 'text-3xl', badge: 'text-[11px]' },
		xl: { icon: 'w-12 h-12', text: 'text-4xl', badge: 'text-xs' }
	};
</script>

<div class="flex items-center gap-2 select-none {className}">
	<!-- Tokopedia Mascot Bag Icon -->
	<svg
		class="{sizes[size].icon} shrink-0 transition-transform duration-200 group-hover:scale-105"
		viewBox="0 0 48 48"
		fill="none"
		xmlns="http://www.w3.org/2000/svg"
		aria-hidden="true"
	>
		<!-- Bag Handle -->
		<path
			d="M17 14V11C17 7.13401 20.134 4 24 4C27.866 4 31 7.13401 31 11V14"
			stroke="#008848"
			stroke-width="3.5"
			stroke-linecap="round"
		/>
		<!-- Bag Body (Tokopedia Green) -->
		<rect x="6" y="12" width="36" height="32" rx="7" fill="#00AA5B" />
		<!-- Inner Bag Shading / Fold -->
		<path
			d="M6 19C10 21 16 22 24 22C32 22 38 21 42 19V40C42 42.2091 40.2091 44 38 44H10C7.79086 44 6 42.2091 6 40V19Z"
			fill="#03AC0E"
			fill-opacity="0.3"
		/>
		<!-- Left Big Owl Eye -->
		<circle cx="17.5" cy="24.5" r="5.5" fill="#FFFFFF" />
		<circle cx="18" cy="24.5" r="3.2" fill="#212121" />
		<circle cx="19.2" cy="23.3" r="1.2" fill="#FFFFFF" />
		<!-- Right Big Owl Eye -->
		<circle cx="30.5" cy="24.5" r="5.5" fill="#FFFFFF" />
		<circle cx="30" cy="24.5" r="3.2" fill="#212121" />
		<circle cx="31.2" cy="23.3" r="1.2" fill="#FFFFFF" />
		<!-- Owl Yellow Beak -->
		<path
			d="M24 27L21.5 31H26.5L24 27Z"
			fill="#FFA000"
			stroke="#F57C00"
			stroke-width="0.8"
			stroke-linejoin="round"
		/>
	</svg>

	{#if showText}
		<div class="flex items-baseline">
			<span
				class="font-extrabold tracking-tight text-[#00AA5B] transition-colors {sizes[size].text}"
			>
				toko<span class="text-[#03AC0E]">ku</span>pedia
			</span>
		</div>
	{/if}
</div>

