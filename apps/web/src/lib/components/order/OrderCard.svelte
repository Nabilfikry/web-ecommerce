<script lang="ts">
	import type { Order } from '@ecommerce/shared-types';
	import { formatIDR, formatDate, getStatusLabel } from '../../utils/format';
	import Badge from '../ui/Badge.svelte';
	import { Package, Store, ChevronRight } from 'lucide-svelte';

	interface Props {
		order: Order;
	}

	let { order }: Props = $props();

	const totalItems = $derived(
		order.items?.reduce((sum, item) => sum + (item.quantity || 0), 0) || 0
	);

	const itemsPreview = $derived(
		order.items
			?.map((item) => `${item.productName || item.product_name} (${item.quantity}x)`)
			.join(', ') || 'Detail produk'
	);
</script>

<div
	class="flex flex-col gap-4 rounded-xl border border-slate-200 bg-white p-5 shadow-xs transition-all hover:border-slate-300"
>
	<!-- Order Header -->
	<div
		class="flex flex-wrap items-center justify-between gap-3 border-b border-slate-100 pb-3 text-xs"
	>
		<div class="flex items-center gap-2">
			<div class="flex h-7 w-7 items-center justify-center rounded-md bg-slate-100 text-slate-700">
				<Package class="h-4 w-4" />
			</div>
			<div>
				<span class="block max-w-[180px] truncate font-bold text-slate-900 sm:max-w-xs">
					Pesanan #{order.id}
				</span>
				<span class="block text-[11px] text-slate-400">
					{formatDate(order.createdAt || order.created_at)}
				</span>
			</div>
		</div>

		<div class="flex items-center gap-2">
			<Badge variant={order.status}>
				{getStatusLabel(order.status)}
			</Badge>
		</div>
	</div>

	<!-- Order Body Preview -->
	<div class="flex flex-col items-start justify-between gap-4 sm:flex-row sm:items-center">
		<div class="min-w-0 flex-1 space-y-1">
			<div class="flex items-center gap-1.5 text-xs font-medium text-slate-500">
				<Store class="h-3.5 w-3.5 text-emerald-600" />
				<span>Toko: {order.seller?.name || 'Toko Penjual'}</span>
			</div>
			<p class="truncate text-xs text-slate-600" title={itemsPreview}>
				{itemsPreview}
			</p>
			<span class="text-[11px] text-slate-400">Total {totalItems} barang</span>
		</div>

		<div
			class="flex w-full items-center justify-between gap-4 border-t border-slate-100 pt-2 sm:w-auto sm:justify-end sm:border-t-0 sm:pt-0"
		>
			<div class="text-left sm:text-right">
				<span class="block text-[10px] text-slate-400">Total Tagihan</span>
				<span class="text-base font-bold text-slate-900 tabular-nums">
					{formatIDR(order.totalAmount || order.total_amount)}
				</span>
			</div>

			<a
				href="/orders/{order.id}"
				class="flex h-9 cursor-pointer items-center gap-1 rounded-lg bg-slate-900 px-3.5 text-xs font-semibold text-white transition-colors hover:bg-slate-800"
			>
				<span>Detail</span>
				<ChevronRight class="h-3.5 w-3.5" />
			</a>
		</div>
	</div>
</div>
