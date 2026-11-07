import { apiFetch, apiFetchWithMeta, ApiClientError } from './client';
import type {
	Order,
	ShippingAddress,
	OrderStatus,
	ApiMeta,
	PayOrderResponse,
	Payment
} from '@ecommerce/shared-types';

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
			const res = await apiFetchWithMeta<Order[]>('/orders', {
				method: 'GET',
				params: {
					page: params.page ?? 1,
					limit: params.limit ?? 10,
					status: params.status
				}
			});

			return {
				orders: res.data || [],
				meta: res.meta || {
					page: params.page ?? 1,
					limit: params.limit ?? 10,
					total: res.data?.length || 0,
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
	},

	async payOrder(id: string): Promise<PayOrderResponse> {
		try {
			return await apiFetch<PayOrderResponse>(`/orders/${id}/pay`, {
				method: 'POST'
			});
		} catch (err: unknown) {
			if (err instanceof ApiClientError && err.code === 'NETWORK_ERROR') {
				const orders = getLocalOrders();
				const order = orders.find((o) => o.id === id);
				const inv = `INV-${Date.now()}`;
				if (order) {
					order.invoice_number = inv;
					order.payment_url = `/orders/${id}?doku_checkout=1&invoice=${inv}&amount=${order.total_amount}`;
					order.payment_status = 'unpaid';
					saveLocalOrders(orders);
				}
				return {
					order_id: id,
					invoice_number: inv,
					amount: order ? (order.total_amount ?? order.totalAmount) : 0,
					payment_url: `/orders/${id}?doku_checkout=1&invoice=${inv}`,
					payment_status: 'unpaid'
				};
			}
			throw err;
		}
	},

	async simulatePayment(
		invoiceNumber: string,
		status: 'SUCCESS' | 'FAILED' | 'EXPIRED' = 'SUCCESS'
	): Promise<void> {
		try {
			await apiFetch<void>('/payments/simulate', {
				method: 'POST',
				body: JSON.stringify({ invoice_number: invoiceNumber, status })
			});
		} catch (err: unknown) {
			if (err instanceof ApiClientError && err.code === 'NETWORK_ERROR') {
				const orders = getLocalOrders();
				const order = orders.find((o) => o.invoice_number === invoiceNumber);
				if (order) {
					if (status === 'SUCCESS') {
						order.payment_status = 'paid';
						if (order.status === 'pending') order.status = 'confirmed';
					} else {
						order.payment_status = status === 'EXPIRED' ? 'expired' : 'failed';
						order.status = 'cancelled';
					}
					saveLocalOrders(orders);
				}
				return;
			}
			throw err;
		}
	},

	async getPaymentStatus(invoiceNumber: string): Promise<Payment> {
		return await apiFetch<Payment>(`/payments/status/${invoiceNumber}`, {
			method: 'GET'
		});
	}
};
