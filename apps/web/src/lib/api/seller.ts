import { apiFetch, apiFetchWithMeta, ApiClientError } from './client';
import type { Product, Order, OrderStatus, ApiMeta } from '@ecommerce/shared-types';
import { SAMPLE_PRODUCTS } from './products';

export interface CreateProductInput {
	name: string;
	description: string;
	price: number;
	stock: number;
	image_url?: string;
}

export interface UpdateProductInput {
	name?: string;
	description?: string;
	price?: number;
	stock?: number;
	image_url?: string;
	is_active?: boolean;
}

export interface SellerProductsResponse {
	products: Product[];
	meta: ApiMeta;
}

export interface SellerOrdersResponse {
	orders: Order[];
	meta: ApiMeta;
}

// Local storage fallback for offline demo
const LOCAL_SELLER_PRODUCTS_KEY = 'ecommerce_demo_seller_products';
const LOCAL_ORDERS_KEY = 'ecommerce_demo_orders';

function getLocalSellerProducts(): Product[] {
	if (typeof window === 'undefined') return [...SAMPLE_PRODUCTS];
	try {
		const raw = localStorage.getItem(LOCAL_SELLER_PRODUCTS_KEY);
		if (raw) return JSON.parse(raw);
	} catch {
		// Ignore parse error
	}
	// Seed with default samples
	saveLocalSellerProducts(SAMPLE_PRODUCTS);
	return [...SAMPLE_PRODUCTS];
}

function saveLocalSellerProducts(products: Product[]) {
	if (typeof window === 'undefined') return;
	localStorage.setItem(LOCAL_SELLER_PRODUCTS_KEY, JSON.stringify(products));
}

function getLocalOrders(): Order[] {
	if (typeof window === 'undefined') return [];
	try {
		const raw = localStorage.getItem(LOCAL_ORDERS_KEY);
		if (raw) return JSON.parse(raw);
	} catch {
		// Ignore parse error
	}
	return [];
}

function saveLocalOrders(orders: Order[]) {
	if (typeof window === 'undefined') return;
	localStorage.setItem(LOCAL_ORDERS_KEY, JSON.stringify(orders));
}

export const sellerApi = {
	async getSellerProducts(
		params: { page?: number; limit?: number; search?: string } = {}
	): Promise<SellerProductsResponse> {
		try {
			const res = await apiFetchWithMeta<Product[]>('/products', {
				method: 'GET',
				params: {
					page: params.page ?? 1,
					limit: params.limit ?? 20,
					search: params.search
				}
			});

			return {
				products: res.data || [],
				meta: res.meta || {
					page: params.page ?? 1,
					limit: params.limit ?? 20,
					total: res.data?.length || 0,
					total_pages: 1
				}
			};
		} catch (err: unknown) {
			if (err instanceof ApiClientError && err.code === 'NETWORK_ERROR') {
				let products = getLocalSellerProducts();
				if (params.search) {
					const q = params.search.toLowerCase();
					products = products.filter(
						(p) => p.name.toLowerCase().includes(q) || p.description.toLowerCase().includes(q)
					);
				}
				const limit = params.limit ?? 20;
				const page = params.page ?? 1;
				const start = (page - 1) * limit;
				const paginated = products.slice(start, start + limit);

				return {
					products: paginated,
					meta: {
						page,
						limit,
						total: products.length,
						total_pages: Math.ceil(products.length / limit) || 1
					}
				};
			}
			throw err;
		}
	},

	async createProduct(input: CreateProductInput): Promise<Product> {
		try {
			return await apiFetch<Product>('/products', {
				method: 'POST',
				body: JSON.stringify(input)
			});
		} catch (err: unknown) {
			if (err instanceof ApiClientError && err.code === 'NETWORK_ERROR') {
				const products = getLocalSellerProducts();
				const newProduct: Product = {
					id: `prod-${Date.now()}`,
					seller_id: 'seller-demo',
					name: input.name,
					description: input.description,
					price: input.price,
					stock: input.stock,
					image_url:
						input.image_url ||
						'https://images.unsplash.com/photo-1523275335684-37898b6baf30?auto=format&fit=crop&w=800&q=80',
					is_active: true,
					seller: {
						id: 'seller-demo',
						name: 'Toko Saya'
					},
					created_at: new Date().toISOString()
				};
				products.unshift(newProduct);
				saveLocalSellerProducts(products);
				return newProduct;
			}
			throw err;
		}
	},

	async updateProduct(id: string, input: UpdateProductInput): Promise<Product> {
		try {
			return await apiFetch<Product>(`/products/${id}`, {
				method: 'PUT',
				body: JSON.stringify(input)
			});
		} catch (err: unknown) {
			if (err instanceof ApiClientError && err.code === 'NETWORK_ERROR') {
				const products = getLocalSellerProducts();
				const index = products.findIndex((p) => p.id === id);
				if (index !== -1) {
					products[index] = {
						...products[index],
						...input,
						updated_at: new Date().toISOString()
					};
					saveLocalSellerProducts(products);
					return products[index];
				}
			}
			throw err;
		}
	},

	async updateStock(id: string, stock: number): Promise<Product> {
		try {
			return await apiFetch<Product>(`/products/${id}/stock`, {
				method: 'PATCH',
				body: JSON.stringify({ stock })
			});
		} catch (err: unknown) {
			if (err instanceof ApiClientError && err.code === 'NETWORK_ERROR') {
				const products = getLocalSellerProducts();
				const index = products.findIndex((p) => p.id === id);
				if (index !== -1) {
					products[index].stock = stock;
					products[index].updated_at = new Date().toISOString();
					saveLocalSellerProducts(products);
					return products[index];
				}
			}
			throw err;
		}
	},

	async deleteProduct(id: string): Promise<void> {
		try {
			await apiFetch<void>(`/products/${id}`, {
				method: 'DELETE'
			});
		} catch (err: unknown) {
			if (err instanceof ApiClientError && err.code === 'NETWORK_ERROR') {
				let products = getLocalSellerProducts();
				products = products.filter((p) => p.id !== id);
				saveLocalSellerProducts(products);
				return;
			}
			throw err;
		}
	},

	async getIncomingOrders(
		params: { page?: number; limit?: number; status?: string } = {}
	): Promise<SellerOrdersResponse> {
		try {
			const res = await apiFetchWithMeta<Order[]>('/orders', {
				method: 'GET',
				params: {
					page: params.page ?? 1,
					limit: params.limit ?? 20,
					status: params.status
				}
			});

			return {
				orders: res.data || [],
				meta: res.meta || {
					page: params.page ?? 1,
					limit: params.limit ?? 20,
					total: res.data?.length || 0,
					total_pages: 1
				}
			};
		} catch (err: unknown) {
			if (err instanceof ApiClientError && err.code === 'NETWORK_ERROR') {
				let orders = getLocalOrders();
				// If empty, create realistic demo incoming orders for seller
				if (orders.length === 0) {
					orders = [
						{
							id: 'ORD-DEMO-001',
							buyer_id: 'buyer-001',
							seller_id: 'seller-demo',
							total_amount: 650000,
							totalAmount: 650000,
							status: 'pending',
							shipping_address: {
								recipient_name: 'Budi Santoso',
								phone: '081234567890',
								address: 'Jl. Riau No. 45, Bandung Wetan',
								city: 'Bandung',
								province: 'Jawa Barat',
								postal_code: '40115'
							},
							shippingAddress: {
								recipient_name: 'Budi Santoso',
								phone: '081234567890',
								address: 'Jl. Riau No. 45, Bandung Wetan',
								city: 'Bandung',
								province: 'Jawa Barat',
								postal_code: '40115'
							},
							notes: 'Tolong packing kayu dan bubble wrap tebal',
							items: [
								{
									productId: 'prod-001',
									product_id: 'prod-001',
									productName: 'Mechanical Keyboard TKL RGB',
									product_name: 'Mechanical Keyboard TKL RGB',
									unitPrice: 650000,
									unit_price: 650000,
									quantity: 1,
									subtotal: 650000
								}
							],
							buyer: {
								id: 'buyer-001',
								name: 'Budi Santoso',
								email: 'budi@example.com'
							},
							created_at: new Date(Date.now() - 3600000 * 2).toISOString()
						},
						{
							id: 'ORD-DEMO-002',
							buyer_id: 'buyer-002',
							seller_id: 'seller-demo',
							total_amount: 840000,
							totalAmount: 840000,
							status: 'confirmed',
							shipping_address: {
								recipient_name: 'Siti Rahma',
								phone: '081987654321',
								address: 'Apartemen Sudirman Tower Lt. 12',
								city: 'Jakarta Selatan',
								province: 'DKI Jakarta',
								postal_code: '12190'
							},
							shippingAddress: {
								recipient_name: 'Siti Rahma',
								phone: '081987654321',
								address: 'Apartemen Sudirman Tower Lt. 12',
								city: 'Jakarta Selatan',
								province: 'DKI Jakarta',
								postal_code: '12190'
							},
							items: [
								{
									productId: 'prod-002',
									product_id: 'prod-002',
									productName: 'Wireless Ergonomic Mouse 26K DPI',
									product_name: 'Wireless Ergonomic Mouse 26K DPI',
									unitPrice: 420000,
									unit_price: 420000,
									quantity: 2,
									subtotal: 840000
								}
							],
							buyer: {
								id: 'buyer-002',
								name: 'Siti Rahma',
								email: 'siti@example.com'
							},
							created_at: new Date(Date.now() - 3600000 * 8).toISOString()
						}
					];
					saveLocalOrders(orders);
				}

				if (params.status && params.status !== 'all') {
					orders = orders.filter((o) => o.status === params.status);
				}

				return {
					orders,
					meta: {
						page: params.page ?? 1,
						limit: params.limit ?? 20,
						total: orders.length,
						total_pages: 1
					}
				};
			}
			throw err;
		}
	},

	async updateOrderStatus(id: string, status: OrderStatus): Promise<Order> {
		try {
			return await apiFetch<Order>(`/orders/${id}/status`, {
				method: 'PATCH',
				body: JSON.stringify({ status })
			});
		} catch (err: unknown) {
			if (err instanceof ApiClientError && err.code === 'NETWORK_ERROR') {
				const orders = getLocalOrders();
				const order = orders.find((o) => o.id === id);
				if (order) {
					order.status = status;
					order.updated_at = new Date().toISOString();
					saveLocalOrders(orders);
					return order;
				}
			}
			throw err;
		}
	}
};
