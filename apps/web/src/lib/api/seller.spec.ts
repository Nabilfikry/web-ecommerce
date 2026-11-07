import { describe, it, expect, beforeEach, vi } from 'vitest';
import { sellerApi } from './seller';

describe('seller api', () => {
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

	it('creates and retrieves seller products', async () => {
		vi.spyOn(globalThis, 'fetch').mockRejectedValue(new Error('Network error'));

		const newProd = await sellerApi.createProduct({
			name: 'Custom Mechanical Keyboard 65%',
			description: 'Compact 65% gasket mount with brass weight',
			price: 1250000,
			stock: 10,
			image_url: 'https://example.com/kb.jpg'
		});

		expect(newProd.id).toBeDefined();
		expect(newProd.name).toBe('Custom Mechanical Keyboard 65%');
		expect(newProd.price).toBe(1250000);
		expect(newProd.stock).toBe(10);

		const { products } = await sellerApi.getSellerProducts();
		expect(products.some((p) => p.name === 'Custom Mechanical Keyboard 65%')).toBe(true);
	});

	it('updates product and stock level', async () => {
		vi.spyOn(globalThis, 'fetch').mockRejectedValue(new Error('Network error'));

		const created = await sellerApi.createProduct({
			name: 'Desk Mat Minimalist',
			description: '900x400mm waterproof desk mat',
			price: 150000,
			stock: 5
		});

		const updated = await sellerApi.updateProduct(created.id, {
			price: 165000,
			name: 'Desk Mat Minimalist Pro'
		});

		expect(updated.price).toBe(165000);
		expect(updated.name).toBe('Desk Mat Minimalist Pro');

		const stockUpdated = await sellerApi.updateStock(created.id, 20);
		expect(stockUpdated.stock).toBe(20);
	});

	it('deletes a product from seller catalog', async () => {
		vi.spyOn(globalThis, 'fetch').mockRejectedValue(new Error('Network error'));

		const created = await sellerApi.createProduct({
			name: 'To be deleted product',
			description: 'Temp item',
			price: 50000,
			stock: 2
		});

		await sellerApi.deleteProduct(created.id);

		const { products } = await sellerApi.getSellerProducts();
		expect(products.some((p) => p.id === created.id)).toBe(false);
	});

	it('retrieves incoming orders and updates status transitions', async () => {
		vi.spyOn(globalThis, 'fetch').mockRejectedValue(new Error('Network error'));

		const { orders } = await sellerApi.getIncomingOrders();
		expect(orders.length).toBeGreaterThan(0);

		const orderToUpdate = orders[0];
		expect(orderToUpdate.id).toBeDefined();

		const confirmed = await sellerApi.updateOrderStatus(orderToUpdate.id, 'confirmed');
		expect(confirmed.status).toBe('confirmed');

		const processing = await sellerApi.updateOrderStatus(orderToUpdate.id, 'processing');
		expect(processing.status).toBe('processing');

		const shipped = await sellerApi.updateOrderStatus(orderToUpdate.id, 'shipped');
		expect(shipped.status).toBe('shipped');
	});
});
