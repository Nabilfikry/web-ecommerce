<script lang="ts">
	import type { Snippet } from 'svelte';
	import { X } from 'lucide-svelte';

	interface Props {
		open: boolean;
		title?: string;
		maxWidth?: 'sm' | 'md' | 'lg' | 'xl';
		onclose?: () => void;
		children?: Snippet;
		footer?: Snippet;
	}

	let {
		open = $bindable(false),
		title,
		maxWidth = 'md',
		onclose,
		children,
		footer
	}: Props = $props();

	function handleClose() {
		open = false;
		if (onclose) onclose();
	}

	function handleBackdrop(e: MouseEvent) {
		if (e.target === e.currentTarget) {
			handleClose();
		}
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape' && open) {
			handleClose();
		}
	}

	const widthClasses = {
		sm: 'max-w-sm',
		md: 'max-w-md',
		lg: 'max-w-lg',
		xl: 'max-w-xl'
	};
</script>

<svelte:window onkeydown={handleKeydown} />

{#if open}
	<div
		class="animate-in fade-in fixed inset-0 z-50 flex items-center justify-center bg-slate-900/40 p-4 backdrop-blur-xs duration-150"
		onclick={handleBackdrop}
		onkeydown={handleKeydown}
		role="dialog"
		aria-modal="true"
		tabindex="-1"
	>
		<div
			class="w-full {widthClasses[
				maxWidth
			]} flex max-h-[90vh] flex-col overflow-hidden rounded-xl border border-slate-200 bg-white shadow-xl"
		>
			{#if title}
				<div class="flex items-center justify-between border-b border-slate-100 px-5 py-4">
					<h3 class="text-base font-semibold text-slate-900">{title}</h3>
					<button
						type="button"
						onclick={handleClose}
						aria-label="Tutup dialog"
						class="cursor-pointer rounded-md p-1 text-slate-400 transition-colors hover:bg-slate-100 hover:text-slate-600"
					>
						<X class="h-4 w-4" />
					</button>
				</div>
			{/if}

			<div class="flex-1 overflow-y-auto p-5">
				{#if children}
					{@render children()}
				{/if}
			</div>

			{#if footer}
				<div
					class="flex items-center justify-end gap-2.5 border-t border-slate-100 bg-slate-50 px-5 py-3.5"
				>
					{@render footer()}
				</div>
			{/if}
		</div>
	</div>
{/if}
