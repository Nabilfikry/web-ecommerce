<script lang="ts">
	import { toast } from '../../stores/toast.svelte';
	import { CheckCircle2, AlertCircle, AlertTriangle, Info, X } from 'lucide-svelte';

	const icons = {
		success: CheckCircle2,
		error: AlertCircle,
		warning: AlertTriangle,
		info: Info
	};

	const styles = {
		success: 'bg-emerald-50 text-emerald-900 border-emerald-200',
		error: 'bg-red-50 text-red-900 border-red-200',
		warning: 'bg-amber-50 text-amber-900 border-amber-200',
		info: 'bg-slate-900 text-white border-slate-800'
	};

	const iconStyles = {
		success: 'text-emerald-600',
		error: 'text-red-600',
		warning: 'text-amber-600',
		info: 'text-slate-300'
	};
</script>

<div
	class="pointer-events-none fixed right-4 bottom-4 z-50 flex w-full max-w-sm flex-col gap-2 p-2"
	aria-live="polite"
>
	{#each toast.toasts as item (item.id)}
		{@const Icon = icons[item.type]}
		<div
			class="pointer-events-auto flex items-start gap-3 rounded-lg border p-3.5 text-sm shadow-lg transition-all duration-200 {styles[
				item.type
			]}"
			role="alert"
		>
			<Icon class="mt-0.5 h-4 w-4 shrink-0 {iconStyles[item.type]}" />
			<div class="flex-1 font-medium">{item.message}</div>
			<button
				type="button"
				onclick={() => toast.remove(item.id)}
				class="cursor-pointer rounded p-0.5 text-current opacity-60 hover:opacity-100"
				aria-label="Tutup notifikasi"
			>
				<X class="h-3.5 w-3.5" />
			</button>
		</div>
	{/each}
</div>
