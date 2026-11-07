import { describe, it, expect, beforeEach, vi } from 'vitest';
import { ordersApi } from './orders';

describe('orders api', () => {
	const storageMock = (() => {
		let store: Record<string, string> = {};
		return {
			getItem: (key: string) => store[key] || null,
			setItem: (key: string, value: string) => {
				store[key] = value.toString();
			},
			removeItem: (key: string) => {
				delete store[key];
			},
			clear: () => {
				store = {};
			}
		};
	})();

	beforeEach(() => {
		vi.stubGlobal('window', {});
		vi.stubGlobal('localStorage', storageMock);
		storageMock.clear();
		vi.restoreAllMocks();
	});

	it('creates order via checkout and persists locally in fallback mode', async () => {
		vi.spyOn(globalThis, 'fetch').mockRejectedValue(new Error('Network error'));

		const newOrder = await ordersApi.checkout({
			shipping_address: {
				recipient_name: 'Naabil',
				phone: '0812345678',
				address: 'Jl. Merdeka No. 1',
				city: 'Bandung',
				postal_code: '40115'
			},
			notes: 'Packing rapi'
		});

		expect(newOrder.id).toBeDefined();
		expect(newOrder.status).toBe('pending');
		expect(newOrder.shippingAddress.recipient_name).toBe('Naabil');
	});

	it('retrieves orders list and filters by status in fallback mode', async () => {
		vi.spyOn(globalThis, 'fetch').mockRejectedValue(new Error('Network error'));

		// Create an order first
		await ordersApi.checkout({
			shipping_address: {
				recipient_name: 'Naabil',
				phone: '0812345678',
				address: 'Jl. Merdeka',
				city: 'Bandung',
				postal_code: '40115'
			}
		});

		const res = await ordersApi.getOrders();
		expect(res.orders.length).toBeGreaterThan(0);
	});

	it('cancels pending order and transitions status to cancelled', async () => {
		vi.spyOn(globalThis, 'fetch').mockRejectedValue(new Error('Network error'));

		const order = await ordersApi.checkout({
			shipping_address: {
				recipient_name: 'Naabil',
				phone: '0812345678',
				address: 'Jl. Merdeka',
				city: 'Bandung',
				postal_code: '40115'
			}
		});

		const cancelled = await ordersApi.cancelOrder(order.id, 'Ubah pesanan');
		expect(cancelled.status).toBe('cancelled');
	});

	it('completes shipped order and transitions status to completed', async () => {
		vi.spyOn(globalThis, 'fetch').mockRejectedValue(new Error('Network error'));

		const order = await ordersApi.checkout({
			shipping_address: {
				recipient_name: 'Naabil',
				phone: '0812345678',
				address: 'Jl. Merdeka',
				city: 'Bandung',
				postal_code: '40115'
			}
		});

		const completed = await ordersApi.completeOrder(order.id);
		expect(completed.status).toBe('completed');
	});
});
