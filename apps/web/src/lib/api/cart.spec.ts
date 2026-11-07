import { describe, it, expect, beforeEach, vi } from 'vitest';
import { cartApi } from './cart';
import type { Product } from '@ecommerce/shared-types';

describe('cart api', () => {
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

	const sampleProduct: Product = {
		id: 'p-1',
		name: 'Keyboard',
		description: 'Desc',
		price: 500000,
		stock: 10,
		is_active: true
	};

	it('adds item to cart in fallback mode and calculates subtotal correctly', async () => {
		vi.spyOn(globalThis, 'fetch').mockRejectedValue(new Error('Network error'));

		const cart = await cartApi.addItem(sampleProduct, 2);
		expect(cart.items).toHaveLength(1);
		expect(cart.items[0].quantity).toBe(2);
		expect(cart.items[0].subtotal).toBe(1000000);
		expect(cart.total_items).toBe(2);
		expect(cart.total_amount).toBe(1000000);
	});

	it('limits added item quantity by available product stock', async () => {
		vi.spyOn(globalThis, 'fetch').mockRejectedValue(new Error('Network error'));

		const limitedProduct = { ...sampleProduct, stock: 3 };
		const cart = await cartApi.addItem(limitedProduct, 10);
		expect(cart.items[0].quantity).toBe(3);
		expect(cart.items[0].subtotal).toBe(1500000);
	});

	it('updates item quantity and recalculates totals', async () => {
		vi.spyOn(globalThis, 'fetch').mockRejectedValue(new Error('Network error'));

		const initialCart = await cartApi.addItem(sampleProduct, 1);
		const itemId = initialCart.items[0].id;

		const updatedCart = await cartApi.updateItemQuantity(itemId, 4);
		expect(updatedCart.items[0].quantity).toBe(4);
		expect(updatedCart.total_amount).toBe(2000000);
	});

	it('removes item from cart and clears cart', async () => {
		vi.spyOn(globalThis, 'fetch').mockRejectedValue(new Error('Network error'));

		const cartWithItem = await cartApi.addItem(sampleProduct, 1);
		const itemId = cartWithItem.items[0].id;

		const afterRemove = await cartApi.removeItem(itemId);
		expect(afterRemove.items).toHaveLength(0);
		expect(afterRemove.total_items).toBe(0);
	});
});
