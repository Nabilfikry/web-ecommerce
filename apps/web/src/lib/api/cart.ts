import { apiFetch, ApiClientError, getAuthToken } from './client';
import type { Cart, CartItem, Product } from '@ecommerce/shared-types';

export interface AddCartItemPayload {
	product_id: string;
	quantity: number;
}

export interface UpdateCartItemPayload {
	quantity: number;
}

// Local storage key for offline/guest cart persistence
const LOCAL_CART_KEY = 'ecommerce_demo_cart';

function normalizeCart(cart: Cart): Cart {
	if (!cart || !Array.isArray(cart.items)) {
		return { id: 'local-cart', items: [], total_items: 0, total_amount: 0 };
	}

	const normalizedItems = cart.items.map((item) => {
		const prodId = item.productId || item.product_id || item.product?.id || '';
		const sellerId =
			item.product?.seller_id || item.product?.sellerId || item.product?.seller?.id || '';

		return {
			...item,
			productId: prodId,
			product_id: prodId,
			subtotal: item.subtotal ?? (item.product ? item.product.price * item.quantity : 0),
			product: item.product
				? {
						...item.product,
						seller_id: sellerId,
						sellerId: sellerId
					}
				: item.product
		};
	});

	const totalItems = normalizedItems.reduce((sum, item) => sum + (item.quantity || 0), 0);
	const totalAmount = normalizedItems.reduce((sum, item) => sum + (item.subtotal || 0), 0);

	return {
		...cart,
		items: normalizedItems,
		total_items: totalItems,
		total_amount: totalAmount
	};
}

function getLocalCart(): Cart {
	if (typeof window === 'undefined') {
		return { id: 'local-cart', items: [], total_items: 0, total_amount: 0 };
	}
	try {
		const raw = localStorage.getItem(LOCAL_CART_KEY);
		if (raw) {
			const parsed = JSON.parse(raw);
			return normalizeCart(parsed);
		}
	} catch {
		// Ignore parse error
	}
	return { id: 'local-cart', items: [], total_items: 0, total_amount: 0 };
}

function saveLocalCart(cart: Cart) {
	if (typeof window === 'undefined') return;
	const normalized = normalizeCart(cart);
	localStorage.setItem(LOCAL_CART_KEY, JSON.stringify(normalized));
}

export const cartApi = {
	async getCart(): Promise<Cart> {
		const token = getAuthToken();
		if (!token) {
			return getLocalCart();
		}

		try {
			const res = await apiFetch<Cart>('/cart', {
				method: 'GET'
			});
			if (res && Array.isArray(res.items)) {
				const normalized = normalizeCart(res);
				saveLocalCart(normalized);
				return normalized;
			}
			return getLocalCart();
		} catch (err: unknown) {
			if (
				err instanceof ApiClientError &&
				(err.code === 'NETWORK_ERROR' || err.status === 401 || err.code === 'UNAUTHORIZED')
			) {
				return getLocalCart();
			}
			throw err;
		}
	},

	async addItem(product: Product, quantity = 1): Promise<Cart> {
		const token = getAuthToken();

		if (token) {
			try {
				const res = await apiFetch<Cart>('/cart/items', {
					method: 'POST',
					body: JSON.stringify({
						product_id: product.id,
						quantity
					})
				});

				if (res && Array.isArray(res.items)) {
					const normalized = normalizeCart(res);
					saveLocalCart(normalized);
					return normalized;
				}
				return await this.getCart();
			} catch (err: unknown) {
				if (
					err instanceof ApiClientError &&
					(err.code === 'NETWORK_ERROR' || err.status === 401 || err.code === 'UNAUTHORIZED')
				) {
					// Fall through to local fallback
				} else {
					throw err;
				}
			}
		}

		// Local/Guest Cart Fallback
		const cart = getLocalCart();
		const existingIndex = cart.items.findIndex((i) => (i.productId || i.product_id) === product.id);

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
		return normalizeCart(cart);
	},

	async updateItemQuantity(itemId: string, quantity: number): Promise<Cart> {
		const token = getAuthToken();

		if (token && !itemId.startsWith('cart-item-')) {
			try {
				const res = await apiFetch<Cart>(`/cart/items/${itemId}`, {
					method: 'PUT',
					body: JSON.stringify({ quantity })
				});

				if (res && Array.isArray(res.items)) {
					const normalized = normalizeCart(res);
					saveLocalCart(normalized);
					return normalized;
				}
				return await this.getCart();
			} catch (err: unknown) {
				if (
					err instanceof ApiClientError &&
					(err.code === 'NETWORK_ERROR' || err.status === 401 || err.code === 'UNAUTHORIZED')
				) {
					// Fall through to local fallback
				} else {
					throw err;
				}
			}
		}

		// Local Cart Fallback
		const cart = getLocalCart();
		const item = cart.items.find((i) => i.id === itemId);
		if (item) {
			item.quantity = Math.max(1, Math.min(quantity, item.product.stock));
			item.subtotal = item.quantity * item.product.price;
			saveLocalCart(cart);
		}
		return normalizeCart(cart);
	},

	async removeItem(itemId: string): Promise<Cart> {
		const token = getAuthToken();

		if (token && !itemId.startsWith('cart-item-')) {
			try {
				const res = await apiFetch<Cart>(`/cart/items/${itemId}`, {
					method: 'DELETE'
				});

				if (res && Array.isArray(res.items)) {
					const normalized = normalizeCart(res);
					saveLocalCart(normalized);
					return normalized;
				}
				return await this.getCart();
			} catch (err: unknown) {
				if (
					err instanceof ApiClientError &&
					(err.code === 'NETWORK_ERROR' || err.status === 401 || err.code === 'UNAUTHORIZED')
				) {
					// Fall through to local fallback
				} else {
					throw err;
				}
			}
		}

		// Local Cart Fallback
		const cart = getLocalCart();
		cart.items = cart.items.filter((i) => i.id !== itemId);
		saveLocalCart(cart);
		return normalizeCart(cart);
	},

	async clearCart(): Promise<void> {
		const emptyCart: Cart = { id: 'local-cart', items: [], total_items: 0, total_amount: 0 };
		saveLocalCart(emptyCart);

		const token = getAuthToken();
		if (token) {
			try {
				await apiFetch<void>('/cart', {
					method: 'DELETE'
				});
			} catch (err: unknown) {
				if (
					err instanceof ApiClientError &&
					(err.code === 'NETWORK_ERROR' || err.status === 401 || err.code === 'UNAUTHORIZED')
				) {
					return;
				}
				throw err;
			}
		}
	},

	async syncGuestCartToServer(): Promise<Cart> {
		const token = getAuthToken();
		if (!token) return getLocalCart();

		const localCart = getLocalCart();
		if (localCart.items.length === 0) {
			return await this.getCart();
		}

		// Push each local item to server cart
		for (const item of localCart.items) {
			try {
				await apiFetch<void>('/cart/items', {
					method: 'POST',
					body: JSON.stringify({
						product_id: item.productId || item.product_id,
						quantity: item.quantity
					})
				});
			} catch {
				// Ignore individual sync errors (e.g. out of stock)
			}
		}

		// Fetch updated canonical cart from server and clear local-only keys
		const serverCart = await this.getCart();
		saveLocalCart(serverCart);
		return serverCart;
	}
};
