import { apiFetch, ApiClientError } from './client';
import type { Order, ShippingAddress, OrderStatus, ApiMeta } from '@ecommerce/shared-types';

export interface CheckoutPayload {
	shipping_address: ShippingAddress;
	notes?: string;
}

export interface GetOrdersParams {
	page?: number;
	limit?: number;
	status?: string;
}

export interface OrdersResponse {
	orders: Order[];
	meta: ApiMeta;
}

// Local mock orders for demo & offline testing
const LOCAL_ORDERS_KEY = 'ecommerce_demo_orders';

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

export const ordersApi = {
	async checkout(payload: CheckoutPayload): Promise<Order> {
		try {
			return await apiFetch<Order>('/orders', {
				method: 'POST',
				body: JSON.stringify(payload)
			});
		} catch (err: unknown) {
			if (err instanceof ApiClientError && err.code === 'NETWORK_ERROR') {
				// Local checkout simulation
				const existingOrders = getLocalOrders();
				const newOrder: Order = {
					id: `order-${Date.now()}`,
					buyerId: 'demo-buyer',
					sellerId: 'seller-tech-01',
					totalAmount: 1070000,
					total_amount: 1070000,
					status: 'pending',
					shippingAddress: payload.shipping_address,
					shipping_address: payload.shipping_address,
					notes: payload.notes || null,
					items: [
						{
							productId: 'prod-001',
							productName: 'Mechanical Keyboard TKL RGB',
							unitPrice: 650000,
							quantity: 1,
							subtotal: 650000
						},
						{
							productId: 'prod-002',
							productName: 'Wireless Ergonomic Mouse 26K DPI',
							unitPrice: 420000,
							quantity: 1,
							subtotal: 420000
						}
					],
					buyer: {
						id: 'demo-buyer',
						name: 'Pembeli Demo',
						email: 'buyer@example.com'
					},
					seller: {
						id: 'seller-tech-01',
						name: 'TechGear Store',
						email: 'seller@techgear.com'
					},
					created_at: new Date().toISOString()
				};

				existingOrders.unshift(newOrder);
				saveLocalOrders(existingOrders);
				return newOrder;
			}
			throw err;
		}
	},

	async getOrders(params: GetOrdersParams = {}): Promise<OrdersResponse> {
		try {
			const res = await apiFetch<Order[]>('/orders', {
				method: 'GET',
				params: {
					page: params.page ?? 1,
					limit: params.limit ?? 10,
					status: params.status
				}
			});

			return {
				orders: res || [],
				meta: {
					page: params.page ?? 1,
					limit: params.limit ?? 10,
					total: res?.length || 0,
					total_pages: 1
				}
			};
		} catch (err: unknown) {
			if (err instanceof ApiClientError && err.code === 'NETWORK_ERROR') {
				let orders = getLocalOrders();
				if (params.status && params.status !== 'all') {
					orders = orders.filter((o) => o.status === params.status);
				}
				return {
					orders,
					meta: {
						page: 1,
						limit: params.limit ?? 10,
						total: orders.length,
						total_pages: 1
					}
				};
			}
			throw err;
		}
	},

	async getOrderDetail(id: string): Promise<Order> {
		try {
			return await apiFetch<Order>(`/orders/${id}`, {
				method: 'GET'
			});
		} catch (err: unknown) {
			if (err instanceof ApiClientError && err.code === 'NETWORK_ERROR') {
				const orders = getLocalOrders();
				const found = orders.find((o) => o.id === id);
				if (found) return found;
			}
			throw err;
		}
	},

	async cancelOrder(id: string, reason?: string): Promise<Order> {
		try {
			return await apiFetch<Order>(`/orders/${id}/cancel`, {
				method: 'PATCH',
				body: JSON.stringify({ reason })
			});
		} catch (err: unknown) {
			if (err instanceof ApiClientError && err.code === 'NETWORK_ERROR') {
				const orders = getLocalOrders();
				const order = orders.find((o) => o.id === id);
				if (order) {
					order.status = 'cancelled';
					saveLocalOrders(orders);
					return order;
				}
			}
			throw err;
		}
	},

	async completeOrder(id: string): Promise<Order> {
		try {
			return await apiFetch<Order>(`/orders/${id}/complete`, {
				method: 'PATCH'
			});
		} catch (err: unknown) {
			if (err instanceof ApiClientError && err.code === 'NETWORK_ERROR') {
				const orders = getLocalOrders();
				const order = orders.find((o) => o.id === id);
				if (order) {
					order.status = 'completed';
					saveLocalOrders(orders);
					return order;
				}
			}
			throw err;
		}
	},

	async updateStatus(id: string, status: OrderStatus): Promise<Order> {
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
					saveLocalOrders(orders);
					return order;
				}
			}
			throw err;
		}
	}
};
