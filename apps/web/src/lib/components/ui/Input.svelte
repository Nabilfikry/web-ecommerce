<script lang="ts">
	import type { HTMLInputAttributes } from 'svelte/elements';

	interface Props extends HTMLInputAttributes {
		label?: string;
		error?: string;
		helperText?: string;
		value?: string | number;
	}

	let {
		label,
		error,
		helperText,
		id = Math.random().toString(36).substring(2, 9),
		type = 'text',
		value = $bindable(''),
		placeholder = '',
		disabled = false,
		required = false,
		class: className = '',
		...rest
	}: Props = $props();
</script>

<div class="flex w-full flex-col gap-1.5">
	{#if label}
		<label for={id} class="flex items-center gap-1 text-xs font-semibold text-slate-800">
			{label}
			{#if required}
				<span class="text-red-500">*</span>
			{/if}
		</label>
	{/if}

	<div class="relative">
		<input
			{id}
			{type}
			bind:value
			{placeholder}
			{disabled}
			{required}
			aria-invalid={!!error}
			aria-describedby={error ? `${id}-error` : helperText ? `${id}-helper` : undefined}
			class="h-11 w-full rounded-lg border bg-white px-3.5 text-sm text-slate-900 transition-colors duration-150 placeholder:text-slate-400 focus:ring-2 focus:ring-offset-1 focus:outline-none disabled:cursor-not-allowed disabled:bg-slate-50 disabled:text-slate-500 {error
				? 'border-red-400 focus:border-red-500 focus:ring-red-200'
				: 'border-slate-300 focus:border-slate-900 focus:ring-slate-900/10'} {className}"
			{...rest}
		/>
	</div>

	{#if error}
		<p id="{id}-error" class="text-xs font-medium text-red-600">{error}</p>
	{:else if helperText}
		<p id="{id}-helper" class="text-xs text-slate-500">{helperText}</p>
	{/if}
</div>
