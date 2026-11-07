import type { User } from '@ecommerce/shared-types';
import { authApi, type LoginPayload, type RegisterPayload } from '../api/auth';
import { getAuthToken, setAuthToken } from '../api/client';

class AuthStore {
	user = $state<User | null>(null);
	token = $state<string | null>(null);
	isLoading = $state<boolean>(true);
	isInitialized = $state<boolean>(false);

	isAuthenticated = $derived(!!this.user);
	isBuyer = $derived(this.user?.role === 'buyer');
	isSeller = $derived(this.user?.role === 'seller');

	async init() {
		if (this.isInitialized) return;
		this.isLoading = true;

		try {
			const savedToken = getAuthToken();
			if (savedToken) {
				this.token = savedToken;
				const currentUser = await authApi.getMe();
				this.user = currentUser;
			}
		} catch {
			// Token expired or invalid
			this.token = null;
			this.user = null;
			setAuthToken(null);
		} finally {
			this.isLoading = false;
			this.isInitialized = true;
		}
	}

	async login(payload: LoginPayload) {
		this.isLoading = true;
		try {
			const res = await authApi.login(payload);
			this.token = res.access_token;
			this.user = res.user;
			return res;
		} finally {
			this.isLoading = false;
		}
	}

	async register(payload: RegisterPayload) {
		this.isLoading = true;
		try {
			return await authApi.register(payload);
		} finally {
			this.isLoading = false;
		}
	}

	async logout() {
		this.isLoading = true;
		try {
			await authApi.logout();
		} finally {
			this.user = null;
			this.token = null;
			setAuthToken(null);
			this.isLoading = false;
		}
	}

	setUser(user: User | null) {
		this.user = user;
	}
}

export const auth = new AuthStore();
export const authStore = auth;
