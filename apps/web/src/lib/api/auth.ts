import { apiFetch, setAuthToken } from './client';
import type { User, LoginResponse, UserRole } from '@ecommerce/shared-types';

export interface RegisterPayload {
	name: string;
	email: string;
	password: string;
	role: UserRole;
}

export interface LoginPayload {
	email: string;
	password: string;
}

export const authApi = {
	async register(payload: RegisterPayload): Promise<User> {
		return apiFetch<User>('/auth/register', {
			method: 'POST',
			body: JSON.stringify(payload)
		});
	},

	async login(payload: LoginPayload): Promise<LoginResponse> {
		const res = await apiFetch<LoginResponse>('/auth/login', {
			method: 'POST',
			body: JSON.stringify(payload)
		});
		if (res?.access_token) {
			setAuthToken(res.access_token);
		}
		return res;
	},

	async getMe(): Promise<User> {
		return apiFetch<User>('/auth/me', {
			method: 'GET'
		});
	},

	async logout(): Promise<void> {
		try {
			await apiFetch<void>('/auth/logout', {
				method: 'POST'
			});
		} finally {
			setAuthToken(null);
		}
	}
};
