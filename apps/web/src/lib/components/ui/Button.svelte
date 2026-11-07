<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLButtonAttributes, HTMLAnchorAttributes } from 'svelte/elements';
	import { Loader2 } from 'lucide-svelte';

	type ButtonVariant = 'primary' | 'secondary' | 'outline' | 'destructive' | 'danger' | 'ghost';

	interface Props extends HTMLButtonAttributes {
		variant?: ButtonVariant;
		size?: 'sm' | 'md' | 'lg';
		loading?: boolean;
		href?: string;
		children?: Snippet;
	}

	let {
		variant = 'primary',
		size = 'md',
		loading = false,
		disabled = false,
		type = 'button',
		href,
		class: className = '',
		children,
		...rest
	}: Props = $props();

	const baseStyles =
		'inline-flex items-center justify-center font-bold transition-all duration-150 rounded-lg focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] disabled:opacity-50 disabled:pointer-events-none cursor-pointer select-none';

	const normalizedVariant = $derived(variant === 'danger' ? 'destructive' : variant);

	const variants: Record<'primary' | 'secondary' | 'outline' | 'destructive' | 'ghost', string> = {
		primary: 'bg-[#00AA5B] text-white hover:bg-[#029b53] focus-visible:ring-[#00AA5B] shadow-xs',
		secondary:
			'bg-[#F0FDF4] text-[#00AA5B] hover:bg-[#DCFCE7] focus-visible:ring-[#00AA5B] border border-emerald-200',
		outline:
			'border border-[#00AA5B] bg-white text-[#00AA5B] hover:bg-[#F0FDF4] focus-visible:ring-[#00AA5B]',
		destructive:
			'bg-[#EE4D2D] text-white hover:bg-[#D73819] focus-visible:ring-[#EE4D2D] shadow-xs',
		ghost:
			'bg-transparent text-[#6D7588] hover:bg-slate-100 hover:text-[#212121] focus-visible:ring-slate-300'
	};

	const sizes = {
		sm: 'h-8 px-3 text-xs gap-1.5 min-h-[32px]',
		md: 'h-10 px-4 text-sm gap-2 min-h-[40px]',
		lg: 'h-12 px-6 text-base gap-2.5 min-h-[48px]'
	};
</script>

{#if href}
	<a
		{href}
		class="{baseStyles} {variants[normalizedVariant]} {sizes[size]} {className}"
		{...rest as HTMLAnchorAttributes}
	>
		{#if loading}
			<Loader2 class="h-4 w-4 animate-spin" />
		{/if}
		{#if children}
			{@render children()}
		{/if}
	</a>
{:else}
	<button
		{type}
		disabled={disabled || loading}
		class="{baseStyles} {variants[normalizedVariant]} {sizes[size]} {className}"
		{...rest}
	>
		{#if loading}
			<Loader2 class="h-4 w-4 animate-spin" />
		{/if}
		{#if children}
			{@render children()}
		{/if}
	</button>
{/if}
