<script lang="ts">
	import type { OrderStatus } from '@ecommerce/shared-types';
	import { CheckCircle2, Clock, Package, Truck, Check, AlertCircle } from 'lucide-svelte';

	interface Props {
		status: OrderStatus;
	}

	let { status }: Props = $props();

	const steps = [
		{ key: 'pending', label: 'Menunggu Konfirmasi', icon: Clock },
		{ key: 'confirmed', label: 'Dikonfirmasi', icon: CheckCircle2 },
		{ key: 'processing', label: 'Diproses', icon: Package },
		{ key: 'shipped', label: 'Dikirim', icon: Truck },
		{ key: 'completed', label: 'Selesai', icon: Check }
	];

	const statusOrder: Record<string, number> = {
		pending: 0,
		confirmed: 1,
		processing: 2,
		shipped: 3,
		completed: 4
	};

	const currentStepIndex = $derived(status in statusOrder ? statusOrder[status] : -1);

	const isTerminalCancel = $derived(status === 'cancelled' || status === 'rejected');
</script>

<div class="rounded-xl border border-slate-200 bg-white p-6">
	<h3 class="mb-6 text-xs font-bold tracking-wider text-slate-800 uppercase">
		Status Perjalanan Pesanan
	</h3>

	{#if isTerminalCancel}
		<div class="flex items-start gap-3 rounded-lg border border-red-200 bg-red-50 p-4 text-red-800">
			<AlertCircle class="mt-0.5 h-5 w-5 shrink-0 text-red-600" />
			<div>
				<h4 class="text-sm font-bold">
					{status === 'cancelled'
						? 'Pesanan Dibatalkan oleh Pembeli'
						: 'Pesanan Ditolak oleh Penjual'}
				</h4>
				<p class="mt-0.5 text-xs text-red-700">
					Sesuai aturan bisnis, stok produk yang sempat dipesan telah dikembalikan ke inventaris
					toko secara otomatis.
				</p>
			</div>
		</div>
	{:else}
		<!-- Horizontal Stepper on Desktop, Vertical on Mobile -->
		<div class="relative flex flex-col justify-between gap-6 md:flex-row md:gap-2">
			{#each steps as step, index (step.key)}
				{@const Icon = step.icon}
				{@const isPast = index < currentStepIndex}
				{@const isCurrent = index === currentStepIndex}

				<div
					class="relative z-10 flex flex-1 items-center gap-3 text-left md:flex-col md:gap-2 md:text-center"
				>
					<!-- Step Icon Indicator -->
					<div
						class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full text-xs font-bold transition-colors {isCurrent
							? 'bg-slate-900 text-white shadow-xs ring-4 ring-slate-100'
							: isPast
								? 'bg-emerald-600 text-white'
								: 'border border-slate-200 bg-slate-100 text-slate-400'}"
					>
						{#if isPast}
							<Check class="h-4 w-4" />
						{:else}
							<Icon class="h-4 w-4" />
						{/if}
					</div>

					<!-- Step Label -->
					<div>
						<p
							class="text-xs leading-tight font-semibold {isCurrent
								? 'font-bold text-slate-900'
								: isPast
									? 'text-emerald-800'
									: 'text-slate-400'}"
						>
							{step.label}
						</p>
						{#if isCurrent}
							<span
								class="mt-0.5 inline-block rounded bg-emerald-50 px-1.5 py-0.5 text-[10px] font-medium text-emerald-700"
							>
								Status Terkini
							</span>
						{/if}
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
