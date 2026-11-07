<script lang="ts">
	import type { Product } from '@ecommerce/shared-types';
	import { sellerApi } from '../../api/seller';
	import { toastStore } from '../../stores/toast.svelte';
	import Modal from '../ui/Modal.svelte';
	import Input from '../ui/Input.svelte';
	import Button from '../ui/Button.svelte';
	import { Image as ImageIcon, Sparkles } from 'lucide-svelte';

	interface Props {
		open: boolean;
		product?: Product | null;
		onsave?: (savedProduct: Product) => void;
		onclose?: () => void;
	}

	let { open = $bindable(false), product = null, onsave, onclose }: Props = $props();

	let isEdit = $derived(!!product);

	let name = $state('');
	let price = $state<number | string>('');
	let stock = $state<number | string>('');
	let imageUrl = $state('');
	let description = $state('');

	let submitting = $state(false);
	let errors = $state<{
		name?: string;
		price?: string;
		stock?: string;
		description?: string;
	}>({});

	// Watch open & product changes to populate fields
	$effect(() => {
		if (open) {
			if (product) {
				name = product.name || '';
				price = product.price ?? '';
				stock = product.stock ?? '';
				imageUrl = product.image_url || product.imageUrl || '';
				description = product.description || '';
			} else {
				name = '';
				price = '';
				stock = '';
				imageUrl = '';
				description = '';
			}
			errors = {};
		}
	});

	function validate(): boolean {
		const newErrors: typeof errors = {};
		if (!name.trim()) {
			newErrors.name = 'Nama produk wajib diisi';
		} else if (name.trim().length < 3) {
			newErrors.name = 'Nama produk minimal 3 karakter';
		}

		const numPrice = Number(price);
		if (price === '' || isNaN(numPrice) || numPrice <= 0) {
			newErrors.price = 'Harga wajib diisi dan lebih dari 0';
		}

		const numStock = Number(stock);
		if (stock === '' || isNaN(numStock) || numStock < 0 || !Number.isInteger(numStock)) {
			newErrors.stock = 'Stok wajib berupa angka bulat (minimal 0)';
		}

		if (!description.trim()) {
			newErrors.description = 'Deskripsi produk wajib diisi';
		} else if (description.trim().length < 10) {
			newErrors.description = 'Deskripsi minimal 10 karakter';
		}

		errors = newErrors;
		return Object.keys(newErrors).length === 0;
	}

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		if (!validate()) return;

		submitting = true;
		try {
			const payload = {
				name: name.trim(),
				price: Number(price),
				stock: Number(stock),
				image_url: imageUrl.trim() || undefined,
				description: description.trim()
			};

			let result: Product;
			if (product?.id) {
				result = await sellerApi.updateProduct(product.id, payload);
				toastStore.success('Produk berhasil diperbarui');
			} else {
				result = await sellerApi.createProduct(payload);
				toastStore.success('Produk baru berhasil ditambahkan');
			}

			open = false;
			if (onsave) onsave(result);
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : 'Gagal menyimpan produk';
			toastStore.error(msg);
		} finally {
			submitting = false;
		}
	}

	function handleSampleImage() {
		const samples = [
			'https://images.unsplash.com/photo-1587829741301-dc798b83add3?auto=format&fit=crop&w=800&q=80',
			'https://images.unsplash.com/photo-1615663245857-ac93bb7c39e7?auto=format&fit=crop&w=800&q=80',
			'https://images.unsplash.com/photo-1546868871-7041f2a55e12?auto=format&fit=crop&w=800&q=80',
			'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?auto=format&fit=crop&w=800&q=80'
		];
		imageUrl = samples[Math.floor(Math.random() * samples.length)];
	}
</script>

<Modal
	bind:open
	title={isEdit ? 'Edit Informasi Produk' : 'Tambah Produk Baru'}
	maxWidth="lg"
	{onclose}
>
	<form id="product-form" onsubmit={handleSubmit} class="space-y-4">
		<Input
			label="Nama Produk"
			required
			placeholder="Contoh: Mechanical Keyboard TKL RGB"
			bind:value={name}
			error={errors.name}
		/>

		<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
			<Input
				label="Harga (IDR)"
				required
				type="number"
				min="1"
				placeholder="650000"
				bind:value={price}
				error={errors.price}
			/>

			<Input
				label="Jumlah Stok"
				required
				type="number"
				min="0"
				step="1"
				placeholder="15"
				bind:value={stock}
				error={errors.stock}
			/>
		</div>

		<div>
			<div class="flex items-center justify-between">
				<label for="product-img" class="text-xs font-semibold text-slate-800">
					URL Gambar Produk
				</label>
				<button
					type="button"
					onclick={handleSampleImage}
					class="inline-flex items-center gap-1 text-xs font-medium text-amber-600 hover:text-amber-700"
				>
					<Sparkles class="h-3 w-3" />
					Pakai Gambar Contoh
				</button>
			</div>
			<input
				id="product-img"
				type="url"
				placeholder="https://images.unsplash.com/..."
				bind:value={imageUrl}
				class="mt-1 block w-full rounded-lg border border-slate-300 px-3 py-2 text-sm text-slate-900 placeholder:text-slate-400 focus:border-slate-900 focus:ring-1 focus:ring-slate-900 focus:outline-hidden"
			/>

			<!-- Image Preview -->
			{#if imageUrl}
				<div
					class="mt-2.5 flex items-center gap-3 rounded-lg border border-slate-200 bg-slate-50 p-2"
				>
					<img
						src={imageUrl}
						alt="Preview"
						class="h-12 w-12 rounded-md border border-slate-200 object-cover"
						onerror={(e) => {
							(e.currentTarget as HTMLImageElement).src =
								'https://images.unsplash.com/photo-1523275335684-37898b6baf30?auto=format&fit=crop&w=400&q=80';
						}}
					/>
					<div class="overflow-hidden text-xs text-ellipsis whitespace-nowrap text-slate-500">
						<span class="font-medium text-slate-700">Preview Gambar Aktif</span>
						<p class="truncate text-[11px] text-slate-400">{imageUrl}</p>
					</div>
				</div>
			{:else}
				<div class="mt-2 flex items-center gap-2 text-xs text-slate-400">
					<ImageIcon class="h-3.5 w-3.5" />
					<span>Jika dikosongkan, sistem akan menggunakan gambar placeholder standar.</span>
				</div>
			{/if}
		</div>

		<div class="flex flex-col gap-1.5">
			<label
				for="product-desc"
				class="flex items-center gap-1 text-xs font-semibold text-slate-800"
			>
				Deskripsi Produk
				<span class="text-red-500">*</span>
			</label>
			<textarea
				id="product-desc"
				rows="4"
				placeholder="Jelaskan spesifikasi, keunggulan, kondisi fisik, dan detail garansi produk..."
				bind:value={description}
				class="block w-full rounded-lg border px-3 py-2 text-sm text-slate-900 transition-colors placeholder:text-slate-400 focus:outline-hidden {errors.description
					? 'border-red-500 focus:border-red-500 focus:ring-1 focus:ring-red-500'
					: 'border-slate-300 focus:border-slate-900 focus:ring-1 focus:ring-slate-900'}"
			></textarea>
			{#if errors.description}
				<p class="text-xs text-red-600">{errors.description}</p>
			{/if}
		</div>
	</form>

	{#snippet footer()}
		<div class="flex w-full items-center justify-end gap-2.5">
			<Button
				variant="secondary"
				onclick={() => {
					open = false;
					if (onclose) onclose();
				}}
				disabled={submitting}
			>
				Batal
			</Button>
			<Button type="submit" form="product-form" variant="primary" loading={submitting}>
				{isEdit ? 'Simpan Perubahan' : 'Terbitkan Produk'}
			</Button>
		</div>
	{/snippet}
</Modal>
