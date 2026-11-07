<script lang="ts">
	import { onMount } from 'svelte';
	import type { Product } from '@ecommerce/shared-types';
	import { sellerApi } from '$lib/api/seller';
	import { formatIDR, formatDate } from '$lib/utils/format';
	import { toastStore } from '$lib/stores/toast.svelte';
	import SellerNav from '$lib/components/seller/SellerNav.svelte';
	import ProductFormModal from '$lib/components/seller/ProductFormModal.svelte';
	import QuickStockModal from '$lib/components/seller/QuickStockModal.svelte';
	import Modal from '$lib/components/ui/Modal.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { Search, Plus, Edit3, Trash2, AlertTriangle, PackageOpen, Sliders } from 'lucide-svelte';

	let loading = $state(true);
	let products = $state<Product[]>([]);
	let searchQuery = $state('');
	let stockFilter = $state<'all' | 'low' | 'out'>('all');

	// Modals
	let isProductModalOpen = $state(false);
	let editingProduct = $state<Product | null>(null);

	let isQuickStockOpen = $state(false);
	let stockTargetProduct = $state<Product | null>(null);

	let isDeleteModalOpen = $state(false);
	let deletingProduct = $state<Product | null>(null);
	let isDeleting = $state(false);

	async function loadProducts() {
		loading = true;
		try {
			const res = await sellerApi.getSellerProducts({ limit: 100 });
			products = res.products;
		} catch (err: unknown) {
			console.error('Failed to load products:', err);
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		loadProducts();
	});

	// Filtered Products
	const filteredProducts = $derived(
		products.filter((p) => {
			const matchesSearch =
				!searchQuery.trim() ||
				p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				p.description.toLowerCase().includes(searchQuery.toLowerCase());

			if (!matchesSearch) return false;

			if (stockFilter === 'low') {
				return (p.stock ?? 0) > 0 && (p.stock ?? 0) <= 5;
			}
			if (stockFilter === 'out') {
				return (p.stock ?? 0) === 0;
			}
			return true;
		})
	);

	function openCreateModal() {
		editingProduct = null;
		isProductModalOpen = true;
	}

	function openEditModal(prod: Product) {
		editingProduct = prod;
		isProductModalOpen = true;
	}

	function openStockModal(prod: Product) {
		stockTargetProduct = prod;
		isQuickStockOpen = true;
	}

	function openDeleteModal(prod: Product) {
		deletingProduct = prod;
		isDeleteModalOpen = true;
	}

	async function confirmDelete() {
		if (!deletingProduct) return;
		isDeleting = true;
		try {
			await sellerApi.deleteProduct(deletingProduct.id);
			toastStore.success(`Produk "${deletingProduct.name}" berhasil dihapus`);
			products = products.filter((p) => p.id !== deletingProduct!.id);
			isDeleteModalOpen = false;
			deletingProduct = null;
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : 'Gagal menghapus produk';
			toastStore.error(msg);
		} finally {
			isDeleting = false;
		}
	}

	function handleProductSaved(savedProduct: Product) {
		const existingIndex = products.findIndex((p) => p.id === savedProduct.id);
		if (existingIndex !== -1) {
			products[existingIndex] = savedProduct;
		} else {
			products = [savedProduct, ...products];
		}
	}

	function handleStockSaved(updatedProduct: Product) {
		products = products.map((p) =>
			p.id === updatedProduct.id ? { ...p, stock: updatedProduct.stock } : p
		);
	}
</script>

<svelte:head>
	<title>Kelola Produk - Seller Dashboard</title>
</svelte:head>

<SellerNav onAddProduct={openCreateModal} />

<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
	<!-- Control Bar: Search & Stock Filters -->
	<div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
		<!-- Search Field -->
		<div class="relative w-full sm:max-w-xs">
			<Search class="absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-slate-400" />
			<input
				type="search"
				placeholder="Cari nama produk..."
				bind:value={searchQuery}
				class="h-9 w-full rounded-lg border border-slate-300 bg-white pr-3 pl-9 text-xs text-slate-900 placeholder:text-slate-400 focus:border-slate-900 focus:ring-1 focus:ring-slate-900 focus:outline-hidden"
			/>
		</div>

		<!-- Stock Filter Tabs -->
		<div
			class="flex items-center gap-1.5 overflow-x-auto rounded-lg border border-slate-200 bg-white p-1 text-xs"
		>
			<button
				type="button"
				onclick={() => (stockFilter = 'all')}
				class="rounded-md px-3 py-1 font-medium transition-colors {stockFilter === 'all'
					? 'bg-slate-900 text-white shadow-xs'
					: 'text-slate-600 hover:bg-slate-100'}"
			>
				Semua ({products.length})
			</button>
			<button
				type="button"
				onclick={() => (stockFilter = 'low')}
				class="rounded-md px-3 py-1 font-medium transition-colors {stockFilter === 'low'
					? 'bg-slate-900 text-white shadow-xs'
					: 'text-slate-600 hover:bg-slate-100'}"
			>
				Stok Rendah (&le;5)
			</button>
			<button
				type="button"
				onclick={() => (stockFilter = 'out')}
				class="rounded-md px-3 py-1 font-medium transition-colors {stockFilter === 'out'
					? 'bg-slate-900 text-white shadow-xs'
					: 'text-slate-600 hover:bg-slate-100'}"
			>
				Habis (0)
			</button>
		</div>
	</div>

	<!-- Products Table -->
	<div class="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-xs">
		{#if loading}
			<div class="space-y-4 p-8">
				{#each [1, 2, 3, 4, 5] as row (row)}
					<div class="h-14 animate-pulse rounded-lg bg-slate-100"></div>
				{/each}
			</div>
		{:else if filteredProducts.length === 0}
			<div class="px-6 py-16 text-center">
				<PackageOpen class="mx-auto h-12 w-12 text-slate-300" />
				<h3 class="mt-4 text-base font-semibold text-slate-900">Tidak ada produk ditemukan</h3>
				<p class="mt-1 text-xs text-slate-500">
					{searchQuery
						? `Tidak ada produk yang cocok dengan pencarian "${searchQuery}"`
						: 'Mulai tambahkan produk pertama Anda ke etalase toko'}
				</p>
				<div class="mt-5">
					<Button variant="primary" size="sm" onclick={openCreateModal}>
						<Plus class="h-4 w-4" />
						Tambah Produk Baru
					</Button>
				</div>
			</div>
		{:else}
			<div class="overflow-x-auto">
				<table class="w-full text-left text-xs">
					<thead
						class="border-b border-slate-200 bg-slate-50/75 text-[11px] font-semibold tracking-wider text-slate-500 uppercase"
					>
						<tr>
							<th scope="col" class="px-6 py-3.5">Produk</th>
							<th scope="col" class="px-6 py-3.5">Harga Satuan</th>
							<th scope="col" class="px-6 py-3.5">Stok Tersedia</th>
							<th scope="col" class="px-6 py-3.5">Tanggal</th>
							<th scope="col" class="px-6 py-3.5 text-right">Aksi</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-100">
						{#each filteredProducts as product (product.id)}
							<tr class="transition-colors hover:bg-slate-50/60">
								<!-- Product Info Column -->
								<td class="px-6 py-4">
									<div class="flex items-center gap-3">
										<img
											src={product.image_url ||
												product.imageUrl ||
												'https://images.unsplash.com/photo-1523275335684-37898b6baf30?auto=format&fit=crop&w=200&q=80'}
											alt={product.name}
											class="h-12 w-12 shrink-0 rounded-lg border border-slate-200 bg-slate-100 object-cover"
											onerror={(e) => {
												(e.currentTarget as HTMLImageElement).src =
													'https://images.unsplash.com/photo-1523275335684-37898b6baf30?auto=format&fit=crop&w=200&q=80';
											}}
										/>
										<div class="max-w-xs min-w-0">
											<p class="truncate text-sm font-semibold text-slate-900">{product.name}</p>
											<p class="mt-0.5 line-clamp-1 text-slate-500">{product.description}</p>
											<span class="font-mono text-[10px] text-slate-400">ID: {product.id}</span>
										</div>
									</div>
								</td>

								<!-- Price Column -->
								<td class="px-6 py-4 text-sm font-bold text-slate-900 tabular-nums">
									{formatIDR(product.price)}
								</td>

								<!-- Stock Column -->
								<td class="px-6 py-4">
									<div class="flex items-center gap-2">
										<span
											class="inline-flex items-center rounded-md px-2 py-0.5 text-xs font-semibold tabular-nums {product.stock ===
											0
												? 'bg-red-50 text-red-700 ring-1 ring-red-600/20 ring-inset'
												: product.stock <= 5
													? 'bg-amber-50 text-amber-700 ring-1 ring-amber-600/20 ring-inset'
													: 'bg-emerald-50 text-emerald-700 ring-1 ring-emerald-600/20 ring-inset'}"
										>
											{product.stock} unit
										</span>
										<button
											type="button"
											onclick={() => openStockModal(product)}
											class="rounded-md border border-slate-200 p-1 text-slate-500 transition-colors hover:border-slate-400 hover:text-slate-800"
											title="Ubah Stok Cepat"
										>
											<Sliders class="h-3.5 w-3.5" />
										</button>
									</div>
								</td>

								<!-- Date Column -->
								<td class="px-6 py-4 whitespace-nowrap text-slate-500">
									{product.created_at || product.createdAt
										? formatDate(product.created_at || product.createdAt || '')
										: '-'}
								</td>

								<!-- Action Column -->
								<td class="px-6 py-4 text-right whitespace-nowrap">
									<div class="inline-flex items-center gap-1.5">
										<button
											type="button"
											onclick={() => openEditModal(product)}
											class="inline-flex items-center gap-1 rounded-md border border-slate-200 px-2.5 py-1 text-xs font-medium text-slate-700 transition-colors hover:bg-slate-100"
										>
											<Edit3 class="h-3.5 w-3.5 text-slate-500" />
											Edit
										</button>
										<button
											type="button"
											onclick={() => openDeleteModal(product)}
											class="inline-flex items-center gap-1 rounded-md border border-slate-200 px-2.5 py-1 text-xs font-medium text-red-600 transition-colors hover:bg-red-50"
										>
											<Trash2 class="h-3.5 w-3.5" />
											Hapus
										</button>
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</div>
</div>

<!-- Product Form Modal (Add & Edit) -->
<ProductFormModal
	bind:open={isProductModalOpen}
	product={editingProduct}
	onsave={handleProductSaved}
/>

<!-- Quick Stock Modal -->
<QuickStockModal
	bind:open={isQuickStockOpen}
	product={stockTargetProduct}
	onsaved={handleStockSaved}
/>

<!-- Delete Confirmation Modal -->
<Modal bind:open={isDeleteModalOpen} title="Hapus Produk" maxWidth="sm">
	{#if deletingProduct}
		<div class="space-y-3">
			<div class="flex items-center gap-2 text-sm font-semibold text-red-600">
				<AlertTriangle class="h-4 w-4" />
				Konfirmasi Penghapusan
			</div>
			<p class="text-xs leading-relaxed text-slate-600">
				Apakah Anda yakin ingin menghapus produk <strong>"{deletingProduct.name}"</strong>? Produk
				yang dihapus tidak akan lagi muncul dalam etalase pembeli.
			</p>
		</div>
	{/if}

	{#snippet footer()}
		<div class="flex w-full items-center justify-end gap-2">
			<Button
				variant="secondary"
				size="sm"
				disabled={isDeleting}
				onclick={() => {
					isDeleteModalOpen = false;
					deletingProduct = null;
				}}
			>
				Batal
			</Button>
			<Button variant="danger" size="sm" loading={isDeleting} onclick={confirmDelete}>
				Ya, Hapus Produk
			</Button>
		</div>
	{/snippet}
</Modal>
