import { describe, it, expect, vi, beforeEach } from 'vitest';
import { productsApi, SAMPLE_PRODUCTS } from './products';

describe('products api', () => {
	beforeEach(() => {
		vi.restoreAllMocks();
	});

	it('returns fallback sample products when network error occurs', async () => {
		vi.spyOn(globalThis, 'fetch').mockRejectedValueOnce(new Error('Network error'));

		const res = await productsApi.getProducts({ limit: 4 });
		expect(res.products).toBeDefined();
		expect(res.products.length).toBeGreaterThan(0);
		expect(res.meta.total).toBe(SAMPLE_PRODUCTS.length);
	});

	it('filters products by search query in fallback mode', async () => {
		vi.spyOn(globalThis, 'fetch').mockRejectedValueOnce(new Error('Network error'));

		const res = await productsApi.getProducts({ search: 'Keyboard' });
		expect(res.products.length).toBeGreaterThan(0);
		expect(res.products[0].name).toContain('Keyboard');
	});

	it('sorts products by price in ascending order in fallback mode', async () => {
		vi.spyOn(globalThis, 'fetch').mockRejectedValueOnce(new Error('Network error'));

		const res = await productsApi.getProducts({ sort: 'price_asc' });
		for (let i = 0; i < res.products.length - 1; i++) {
			expect(res.products[i].price).toBeLessThanOrEqual(res.products[i + 1].price);
		}
	});

	it('sorts products by price in descending order in fallback mode', async () => {
		vi.spyOn(globalThis, 'fetch').mockRejectedValueOnce(new Error('Network error'));

		const res = await productsApi.getProducts({ sort: 'price_desc' });
		for (let i = 0; i < res.products.length - 1; i++) {
			expect(res.products[i].price).toBeGreaterThanOrEqual(res.products[i + 1].price);
		}
	});

	it('fetches single product detail by id', async () => {
		vi.spyOn(globalThis, 'fetch').mockRejectedValueOnce(new Error('Network error'));

		const product = await productsApi.getProduct('prod-001');
		expect(product).toBeDefined();
		expect(product.id).toBe('prod-001');
		expect(product.name).toContain('Keyboard');
	});
});
