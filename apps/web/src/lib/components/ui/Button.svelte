<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLButtonAttributes } from 'svelte/elements';
	import { Loader2 } from 'lucide-svelte';

	interface Props extends HTMLButtonAttributes {
		variant?: 'primary' | 'secondary' | 'outline' | 'destructive' | 'ghost';
		size?: 'sm' | 'md' | 'lg';
		loading?: boolean;
		children?: Snippet;
	}

	let {
		variant = 'primary',
		size = 'md',
		loading = false,
		disabled = false,
		type = 'button',
		class: className = '',
		children,
		...rest
	}: Props = $props();

	const baseStyles =
		'inline-flex items-center justify-center font-medium transition-all duration-150 rounded-lg focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] disabled:opacity-50 disabled:pointer-events-none cursor-pointer';

	const variants = {
		primary: 'bg-slate-900 text-white hover:bg-slate-800 focus-visible:ring-slate-900 shadow-xs',
		secondary: 'bg-slate-100 text-slate-900 hover:bg-slate-200 focus-visible:ring-slate-400',
		outline:
			'border border-slate-300 bg-transparent text-slate-800 hover:bg-slate-50 focus-visible:ring-slate-400',
		destructive: 'bg-red-600 text-white hover:bg-red-700 focus-visible:ring-red-600 shadow-xs',
		ghost:
			'bg-transparent text-slate-700 hover:bg-slate-100 hover:text-slate-900 focus-visible:ring-slate-400'
	};

	const sizes = {
		sm: 'h-9 px-3 text-xs gap-1.5 min-h-[36px]',
		md: 'h-11 px-5 text-sm gap-2 min-h-[44px]',
		lg: 'h-12 px-6 text-base gap-2.5 min-h-[48px]'
	};
</script>

<button
	{type}
	disabled={disabled || loading}
	class="{baseStyles} {variants[variant]} {sizes[size]} {className}"
	{...rest}
>
	{#if loading}
		<Loader2 class="h-4 w-4 shrink-0 animate-spin" />
	{/if}
	{#if children}
		{@render children()}
	{/if}
</button>
