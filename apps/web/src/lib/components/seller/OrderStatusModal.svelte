<script lang="ts">
	import type { Order, OrderStatus } from '@ecommerce/shared-types';
	import { sellerApi } from '../../api/seller';
	import { toastStore } from '../../stores/toast.svelte';
	import { formatIDR, getStatusLabel } from '../../utils/format';
	import Modal from '../ui/Modal.svelte';
	import Button from '../ui/Button.svelte';
	import Badge from '../ui/Badge.svelte';
	import { CheckCircle2, XCircle, PackageCheck, Truck, AlertTriangle } from 'lucide-svelte';

	interface Props {
		open: boolean;
		order?: Order | null;
		onupdated?: (updatedOrder: Order) => void;
		onclose?: () => void;
	}

	let { open = $bindable(false), order = null, onupdated, onclose }: Props = $props();

	let updating = $state(false);
	let rejectMode = $state(false);

	$effect(() => {
		if (open) {
			rejectMode = false;
		}
	});

	async function handleStatusChange(targetStatus: OrderStatus) {
		if (!order?.id) return;

		updating = true;
		try {
			const updated = await sellerApi.updateOrderStatus(order.id, targetStatus);
			toastStore.success(
				`Status pesanan berhasil diubah menjadi "${getStatusLabel(targetStatus)}"`
			);
			open = false;
			if (onupdated) onupdated(updated);
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : 'Gagal mengubah status pesanan';
			toastStore.error(msg);
		} finally {
			updating = false;
		}
	}
</script>

<Modal bind:open title="Tindak Lanjuti Pesanan" maxWidth="md" {onclose}>
	{#if order}
		<div class="space-y-4">
			<!-- Order Header Summary -->
			<div class="rounded-lg border border-slate-200 bg-slate-50 p-3.5">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-xs font-medium text-slate-500">Nomor Pesanan</p>
						<p class="font-mono text-xs font-bold text-slate-900">{order.id}</p>
					</div>
					<Badge status={order.status} />
				</div>

				<div
					class="mt-3 flex justify-between border-t border-slate-200 pt-2 text-xs text-slate-600"
				>
					<span>Total Tagihan:</span>
					<span class="font-bold text-slate-900 tabular-nums">
						{formatIDR(order.total_amount || order.totalAmount || 0)}
					</span>
				</div>

				<div class="mt-1 text-xs text-slate-500">
					Pembeli: <span class="font-medium text-slate-800"
						>{order.buyer?.name || order.shippingAddress?.recipient_name || 'Pembeli'}</span
					>
				</div>
			</div>

			<!-- Actions depending on status -->
			{#if order.status === 'pending'}
				{#if !rejectMode}
					<div class="space-y-3">
						<p class="text-sm text-slate-600">
							Pesanan baru masuk dan menunggu konfirmasi Anda. Apakah Anda ingin menerima pesanan
							ini untuk disiapkan?
						</p>

						<div class="flex flex-col gap-2 pt-2">
							<Button
								variant="primary"
								class="w-full justify-center bg-emerald-600 text-white hover:bg-emerald-700"
								loading={updating}
								onclick={() => handleStatusChange('confirmed')}
							>
								<CheckCircle2 class="h-4 w-4" />
								Terima & Konfirmasi Pesanan
							</Button>

							<button
								type="button"
								onclick={() => (rejectMode = true)}
								class="inline-flex w-full items-center justify-center gap-1.5 rounded-lg border border-slate-200 py-2 text-xs font-semibold text-red-600 transition-colors hover:bg-red-50"
							>
								<XCircle class="h-3.5 w-3.5" />
								Tolak Pesanan Ini
							</button>
						</div>
					</div>
				{:else}
					<div class="space-y-3 rounded-lg border border-red-200 bg-red-50 p-4">
						<div class="flex items-center gap-2 text-sm font-semibold text-red-800">
							<AlertTriangle class="h-4 w-4 text-red-600" />
							Konfirmasi Tolak Pesanan
						</div>
						<p class="text-xs leading-relaxed text-red-700">
							Menolak pesanan akan membatalkan transaksi ini. Stok barang yang dipesan akan otomatis
							dikembalikan ke inventaris produk Anda.
						</p>
						<div class="flex items-center justify-end gap-2 pt-2">
							<button
								type="button"
								onclick={() => (rejectMode = false)}
								class="rounded-lg px-3 py-1.5 text-xs font-medium text-slate-600 hover:bg-slate-200"
							>
								Kembali
							</button>
							<Button
								variant="danger"
								size="sm"
								loading={updating}
								onclick={() => handleStatusChange('rejected')}
							>
								Ya, Tolak Pesanan
							</Button>
						</div>
					</div>
				{/if}
			{:else if order.status === 'confirmed'}
				<div class="space-y-3">
					<p class="text-sm text-slate-600">
						Pesanan sudah dikonfirmasi. Tekan tombol di bawah jika Anda mulai mengemas dan
						menyiapkan pesanan.
					</p>

					<div class="pt-2">
						<Button
							variant="primary"
							class="w-full justify-center"
							loading={updating}
							onclick={() => handleStatusChange('processing')}
						>
							<PackageCheck class="h-4 w-4" />
							Mulai Siapkan Pesanan (Processing)
						</Button>
					</div>
				</div>
			{:else if order.status === 'processing'}
				<div class="space-y-3">
					<p class="text-sm text-slate-600">
						Paket sudah selesai dikemas dan siap diantar ke kurir logistik. Perbarui status menjadi
						telah dikirim.
					</p>

					<div class="pt-2">
						<Button
							variant="primary"
							class="w-full justify-center bg-indigo-600 text-white hover:bg-indigo-700"
							loading={updating}
							onclick={() => handleStatusChange('shipped')}
						>
							<Truck class="h-4 w-4" />
							Konfirmasi Kirim Pesanan (Shipped)
						</Button>
					</div>
				</div>
			{:else}
				<div class="rounded-lg bg-slate-50 p-4 text-center text-sm text-slate-500">
					Pesanan ini berstatus <strong>{getStatusLabel(order.status)}</strong> dan tidak membutuhkan
					tindakan pembaruan lebih lanjut dari penjual.
				</div>
			{/if}
		</div>
	{/if}

	{#snippet footer()}
		<div class="flex w-full items-center justify-end">
			<Button
				variant="secondary"
				onclick={() => {
					open = false;
					if (onclose) onclose();
				}}
				disabled={updating}
			>
				Tutup
			</Button>
		</div>
	{/snippet}
</Modal>
