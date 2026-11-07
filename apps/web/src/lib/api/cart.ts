import { apiFetch, ApiClientError } from './client';
import type { Cart, CartItem, Product } from '@ecommerce/shared-types';

export interface AddCartItemPayload {
	product_id: string;
	quantity: number;
}

export interface UpdateCartItemPayload {
	quantity: number;
}

// Local storage key for offline/demo cart persistence
const LOCAL_CART_KEY = 'ecommerce_demo_cart';

function getLocalCart(): Cart {
	if (typeof window === 'undefined') {
		return { id: 'local-cart', items: [], total_items: 0, total_amount: 0 };
	}
	try {
		const raw = localStorage.getItem(LOCAL_CART_KEY);
		if (raw) return JSON.parse(raw);
	} catch {
		// Ignore parse error
	}
	return { id: 'local-cart', items: [], total_items: 0, total_amount: 0 };
}

function saveLocalCart(cart: Cart) {
	if (typeof window === 'undefined') return;
	cart.total_items = cart.items.reduce((sum, item) => sum + item.quantity, 0);
	cart.total_amount = cart.items.reduce((sum, item) => sum + item.subtotal, 0);
	localStorage.setItem(LOCAL_CART_KEY, JSON.stringify(cart));
}

export const cartApi = {
	async getCart(): Promise<Cart> {
		try {
			const res = await apiFetch<Cart>('/cart', {
				method: 'GET'
			});
			if (res && res.items) return res;
			return getLocalCart();
		} catch (err: unknown) {
			if (err instanceof ApiClientError && err.code === 'NETWORK_ERROR') {
				return getLocalCart();
			}
			throw err;
		}
	},

	async addItem(product: Product, quantity = 1): Promise<Cart> {
		try {
			await apiFetch<void>('/cart/items', {
				method: 'POST',
				body: JSON.stringify({
					product_id: product.id,
					quantity
				})
			});
			return await this.getCart();
		} catch (err: unknown) {
			if (err instanceof ApiClientError && err.code === 'NETWORK_ERROR') {
				const cart = getLocalCart();
				const existingIndex = cart.items.findIndex(
					(i) => i.productId === product.id || i.product_id === product.id
				);

				if (existingIndex > -1) {
					const existing = cart.items[existingIndex];
					const newQty = existing.quantity + quantity;
					existing.quantity = Math.min(newQty, product.stock);
					existing.subtotal = existing.quantity * product.price;
				} else {
					const newItem: CartItem = {
						id: `cart-item-${Date.now()}`,
						productId: product.id,
						product_id: product.id,
						product,
						quantity: Math.min(quantity, product.stock),
						subtotal: Math.min(quantity, product.stock) * product.price
					};
					cart.items.push(newItem);
				}

				saveLocalCart(cart);
				return cart;
			}
			throw err;
		}
	},

	async updateItemQuantity(itemId: string, quantity: number): Promise<Cart> {
		try {
			await apiFetch<void>(`/cart/items/${itemId}`, {
				method: 'PATCH',
				body: JSON.stringify({ quantity })
			});
			return await this.getCart();
		} catch (err: unknown) {
			if (err instanceof ApiClientError && err.code === 'NETWORK_ERROR') {
				const cart = getLocalCart();
				const item = cart.items.find((i) => i.id === itemId);
				if (item) {
					item.quantity = Math.max(1, quantity);
					item.subtotal = item.quantity * item.product.price;
					saveLocalCart(cart);
				}
				return cart;
			}
			throw err;
		}
	},

	async removeItem(itemId: string): Promise<Cart> {
		try {
			await apiFetch<void>(`/cart/items/${itemId}`, {
				method: 'DELETE'
			});
			return await this.getCart();
		} catch (err: unknown) {
			if (err instanceof ApiClientError && err.code === 'NETWORK_ERROR') {
				const cart = getLocalCart();
				cart.items = cart.items.filter((i) => i.id !== itemId);
				saveLocalCart(cart);
				return cart;
			}
			throw err;
		}
	},

	async clearCart(): Promise<void> {
		try {
			await apiFetch<void>('/cart', {
				method: 'DELETE'
			});
		} catch (err: unknown) {
			if (err instanceof ApiClientError && err.code === 'NETWORK_ERROR') {
				const cart = { id: 'local-cart', items: [], total_items: 0, total_amount: 0 };
				saveLocalCart(cart);
				return;
			}
			throw err;
		}
	}
};
