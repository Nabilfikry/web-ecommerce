import { apiFetch, apiFetchWithMeta, ApiClientError } from './client';
import type { Product, ApiMeta } from '@ecommerce/shared-types';

export interface GetProductsParams {
	page?: number;
	limit?: number;
	search?: string;
	sort?: 'latest' | 'price_asc' | 'price_desc';
}

export interface ProductsResponse {
	products: Product[];
	meta: ApiMeta;
}

// High quality fallback sample products matching Tokopedia Unify & PRD
export const SAMPLE_PRODUCTS: Product[] = [
	{
		id: 'microsoft-surface-pro-12-copilot-pc',
		seller_id: 'seller-jullindo',
		name: 'Microsoft Surface Pro 12" Copilot+ PC, Snapdragon X Plus (8 Core) Ram 16GB SSD 256/512GB 12-inch Touchscreen - 512GB Platinum',
		description:
			'Microsoft Surface Pro 12" Copilot+ PC Generasi Terbaru bertenaga Qualcomm Snapdragon X Plus (8 Core), NPU AI 45 TOPS, RAM 16GB LPDDR5X, SSD 512GB NVMe. Layar 12 inch PixelSense Flow OLED 120Hz Touchscreen dengan dukungan Surface Slim Pen. Dilengkapi Windows 11 Home & Microsoft Office Home & Student asli. Garansi Resmi 1 Tahun.',
		price: 24500000,
		stock: 12,
		image_url:
			'https://images.unsplash.com/photo-1544244015-0df4b3ffc6b0?auto=format&fit=crop&w=800&q=80',
		is_active: true,
		seller: {
			id: 'seller-jullindo',
			name: 'Jullindo Official Store'
		},
		created_at: '2026-11-06T12:00:00Z'
	},
	{
		id: 'prod-001',
		seller_id: 'seller-tech-01',
		name: 'Mechanical Keyboard TKL RGB',
		description:
			'Keyboard mechanical tenkeyless dengan Gateron Yellow Pro switches, double-shot PBT keycaps, dan hotswappable PCB 5-pin. Sangat nyaman untuk mengetik dan gaming.',
		price: 650000,
		stock: 14,
		image_url:
			'https://images.unsplash.com/photo-1587829741301-dc798b83add3?auto=format&fit=crop&w=800&q=80',
		is_active: true,
		seller: {
			id: 'seller-tech-01',
			name: 'TechGear Store'
		},
		created_at: '2026-11-05T08:00:00Z'
	},
	{
		id: 'prod-002',
		seller_id: 'seller-tech-01',
		name: 'Wireless Ergonomic Mouse 26K DPI',
		description:
			'Mouse nirkabel ultra-lightweight 58 gram dengan sensor optik 26.000 DPI, koneksi 2.4GHz tri-mode, dan baterai tahan hingga 80 jam pemakaian.',
		price: 420000,
		stock: 25,
		image_url:
			'https://images.unsplash.com/photo-1615663245857-ac93bb7c39e7?auto=format&fit=crop&w=800&q=80',
		is_active: true,
		seller: {
			id: 'seller-tech-01',
			name: 'TechGear Store'
		},
		created_at: '2026-11-05T09:30:00Z'
	},
	{
		id: 'prod-003',
		seller_id: 'seller-audio-02',
		name: 'Active Noise Cancelling Headphones',
		description:
			'Headphone over-ear dengan driver 40mm beryllium, peredam bising aktif hibrida, dan codec audio resolusi tinggi. Nyaman digunakan seharian.',
		price: 890000,
		stock: 8,
		image_url:
			'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?auto=format&fit=crop&w=800&q=80',
		is_active: true,
		seller: {
			id: 'seller-audio-02',
			name: 'Sonic Sound Lab'
		},
		created_at: '2026-11-06T11:00:00Z'
	},
	{
		id: 'prod-004',
		seller_id: 'seller-tech-01',
		name: 'Minimalist Desk Mat Felt & Leather',
		description:
			'Alas meja kerja premium ukuran 90x40cm berbahan kulit sintetis anti-air dan felt lembut. Melindungi permukaan meja sekaligus estetis.',
		price: 185000,
		stock: 3,
		image_url:
			'https://images.unsplash.com/photo-1527864550417-7fd91fc51a46?auto=format&fit=crop&w=800&q=80',
		is_active: true,
		seller: {
			id: 'seller-tech-01',
			name: 'TechGear Store'
		},
		created_at: '2026-11-06T14:20:00Z'
	},
	{
		id: 'prod-005',
		seller_id: 'seller-audio-02',
		name: 'Studio Monitor Speakers 4-inch',
		description:
			'Sepasang speaker monitor aktif 50W dengan input optik, Bluetooth 5.0, dan output subwoofer terpisah. Suara jernih dan bass bertenaga.',
		price: 1250000,
		stock: 0,
		image_url:
			'https://images.unsplash.com/photo-1545454675-3531b543be5d?auto=format&fit=crop&w=800&q=80',
		is_active: true,
		seller: {
			id: 'seller-audio-02',
			name: 'Sonic Sound Lab'
		},
		created_at: '2026-11-07T01:10:00Z'
	},
	{
		id: 'prod-006',
		seller_id: 'seller-tech-01',
		name: 'USB-C GaN Fast Charger 100W',
		description:
			'Kepala charger GaN 4-port (3x USB-C, 1x USB-A) dengan dukungan Power Delivery 3.0 dan PPS. Mampu mengisi laptop dan smartphone sekaligus.',
		price: 345000,
		stock: 19,
		image_url:
			'https://images.unsplash.com/photo-1583863788434-e58a36330cf0?auto=format&fit=crop&w=800&q=80',
		is_active: true,
		seller: {
			id: 'seller-tech-01',
			name: 'TechGear Store'
		},
		created_at: '2026-11-07T02:00:00Z'
	}
];

export const productsApi = {
	async getProducts(params: GetProductsParams = {}): Promise<ProductsResponse> {
		try {
			const res = await apiFetchWithMeta<Product[]>('/products', {
				method: 'GET',
				params: {
					page: params.page ?? 1,
					limit: params.limit ?? 12,
					search: params.search,
					sort: params.sort ?? 'latest'
				}
			});

			return {
				products: res.data || [],
				meta: res.meta || {
					page: params.page ?? 1,
					limit: params.limit ?? 12,
					total: res.data?.length || 0,
					total_pages: 1
				}
			};
		} catch (err: unknown) {
			// Fallback to local sample products when backend is not running
			if (err instanceof ApiClientError && err.code === 'NETWORK_ERROR') {
				let filtered = [...SAMPLE_PRODUCTS];

				if (params.search) {
					const q = params.search.toLowerCase();
					filtered = filtered.filter(
						(p) => p.name.toLowerCase().includes(q) || p.description.toLowerCase().includes(q)
					);
				}

				if (params.sort === 'price_asc') {
					filtered.sort((a, b) => a.price - b.price);
				} else if (params.sort === 'price_desc') {
					filtered.sort((a, b) => b.price - a.price);
				}

				return {
					products: filtered,
					meta: {
						page: 1,
						limit: params.limit ?? 12,
						total: filtered.length,
						total_pages: Math.ceil(filtered.length / (params.limit ?? 12)) || 1
					}
				};
			}
			throw err;
		}
	},

	async getProduct(id: string): Promise<Product> {
		if (
			id === 'microsoft-surface-pro-12-copilot-pc' ||
			id.includes('surface-pro-12') ||
			id.includes('1731633667861546606')
		) {
			const found = SAMPLE_PRODUCTS.find((p) => p.id === 'microsoft-surface-pro-12-copilot-pc');
			if (found) return found;
		}

		try {
			return await apiFetch<Product>(`/products/${id}`, {
				method: 'GET'
			});
		} catch (err: unknown) {
			const found = SAMPLE_PRODUCTS.find((p) => p.id === id || id.includes(p.id));
			if (found) return found;
			throw err;
		}
	}
};
