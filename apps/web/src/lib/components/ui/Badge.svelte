<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { OrderStatus } from '@ecommerce/shared-types';

	type Variant =
		OrderStatus | 'default' | 'secondary' | 'outline' | 'accent' | 'success' | 'danger';

	interface Props {
		variant?: Variant;
		size?: 'sm' | 'md';
		class?: string;
		children?: Snippet;
	}

	let { variant = 'default', size = 'md', class: className = '', children }: Props = $props();

	const variantStyles: Record<Variant, string> = {
		default: 'bg-slate-100 text-slate-800 border-slate-200',
		secondary: 'bg-slate-200 text-slate-900 border-transparent',
		outline: 'bg-transparent text-slate-700 border-slate-300',
		accent: 'bg-amber-500 text-white border-transparent shadow-xs',
		success: 'bg-emerald-100 text-emerald-800 border-emerald-200',
		danger: 'bg-red-100 text-red-800 border-red-200',
		// Order Statuses
		pending: 'bg-amber-100 text-amber-900 border-amber-300',
		confirmed: 'bg-sky-100 text-sky-900 border-sky-200',
		processing: 'bg-blue-100 text-blue-900 border-blue-200',
		shipped: 'bg-indigo-100 text-indigo-900 border-indigo-200',
		completed: 'bg-emerald-100 text-emerald-900 border-emerald-300',
		cancelled: 'bg-slate-100 text-slate-700 border-slate-300',
		rejected: 'bg-red-100 text-red-900 border-red-300'
	};

	const sizeStyles = {
		sm: 'text-[11px] px-2 py-0.5 tracking-wider uppercase font-semibold',
		md: 'text-xs px-2.5 py-1 font-medium'
	};
</script>

<span
	class="inline-flex items-center justify-center rounded-md border font-sans select-none {variantStyles[
		variant
	] || variantStyles.default} {sizeStyles[size]} {className}"
>
	{#if children}
		{@render children()}
	{/if}
</span>
