import type { Cart, Product } from '@ecommerce/shared-types';
import { cartApi } from '../api/cart';

class CartStore {
	cart = $state<Cart>({
		id: '',
		items: [],
		total_items: 0,
		total_amount: 0
	});

	loading = $state(false);
	initialized = $state(false);

	items = $derived(this.cart.items);

	totalItems = $derived(this.cart.items.reduce((sum, item) => sum + (item.quantity || 0), 0));

	totalAmount = $derived(this.cart.items.reduce((sum, item) => sum + (item.subtotal || 0), 0));

	// Detect if cart contains items from multiple sellers (BR-010 violation check)
	sellerIds = $derived(
		Array.from(
			new Set(
				this.cart.items
					.map(
						(item) => item.product?.seller_id || item.product?.sellerId || item.product?.seller?.id
					)
					.filter(Boolean)
			)
		)
	);

	hasMultipleSellers = $derived(this.sellerIds.length > 1);

	async init() {
		if (this.initialized) return;
		await this.fetchCart();
		this.initialized = true;
	}

	async fetchCart() {
		this.loading = true;
		try {
			const res = await cartApi.getCart();
			this.cart = res;
		} finally {
			this.loading = false;
		}
	}

	async addItem(product: Product, quantity = 1) {
		this.loading = true;
		try {
			const updated = await cartApi.addItem(product, quantity);
			this.cart = updated;
			return updated;
		} finally {
			this.loading = false;
		}
	}

	async updateQuantity(itemId: string, quantity: number) {
		this.loading = true;
		try {
			const updated = await cartApi.updateItemQuantity(itemId, quantity);
			this.cart = updated;
			return updated;
		} finally {
			this.loading = false;
		}
	}

	async removeItem(itemId: string) {
		this.loading = true;
		try {
			const updated = await cartApi.removeItem(itemId);
			this.cart = updated;
			return updated;
		} finally {
			this.loading = false;
		}
	}

	async clearCart() {
		this.loading = true;
		try {
			await cartApi.clearCart();
			this.cart = { id: '', items: [], total_items: 0, total_amount: 0 };
		} finally {
			this.loading = false;
		}
	}

	async syncGuestCart() {
		this.loading = true;
		try {
			const updated = await cartApi.syncGuestCartToServer();
			this.cart = updated;
			return updated;
		} finally {
			this.loading = false;
		}
	}

	async keepOnlySeller(targetSellerId: string) {
		this.loading = true;
		try {
			const itemsToRemove = this.cart.items.filter((item) => {
				const sid = item.product?.seller_id || item.product?.sellerId || item.product?.seller?.id;
				return sid !== targetSellerId;
			});

			for (const it of itemsToRemove) {
				await cartApi.removeItem(it.id);
			}

			const updated = await cartApi.getCart();
			this.cart = updated;
			return updated;
		} finally {
			this.loading = false;
		}
	}

	reset() {
		this.cart = { id: '', items: [], total_items: 0, total_amount: 0 };
	}
}

export const cart = new CartStore();
