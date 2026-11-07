import { describe, it, expect, beforeEach, vi } from 'vitest';
import { getAuthToken, setAuthToken, ApiClientError } from './client';
import { authApi } from './auth';

describe('auth api & client', () => {
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
		// Mock window and localStorage for node test environment
		vi.stubGlobal('window', {});
		vi.stubGlobal('localStorage', storageMock);
		storageMock.clear();
		vi.restoreAllMocks();
	});

	it('manages token storage in localStorage', () => {
		expect(getAuthToken()).toBeNull();
		setAuthToken('sample-jwt-token');
		expect(getAuthToken()).toBe('sample-jwt-token');
		setAuthToken(null);
		expect(getAuthToken()).toBeNull();
	});

	it('creates ApiClientError with code and status', () => {
		const err = new ApiClientError('Invalid input', 400, 'VALIDATION_ERROR', [
			{ field: 'email', message: 'Email tidak valid' }
		]);

		expect(err.message).toBe('Invalid input');
		expect(err.status).toBe(400);
		expect(err.code).toBe('VALIDATION_ERROR');
		expect(err.details).toHaveLength(1);
	});

	it('invokes login endpoint and saves token on success', async () => {
		const mockResponse = {
			success: true,
			data: {
				access_token: 'test-token-123',
				token_type: 'Bearer',
				expires_in: 900,
				user: {
					id: 'u-1',
					name: 'Test Buyer',
					email: 'test@example.com',
					role: 'buyer' as const
				}
			}
		};

		vi.spyOn(globalThis, 'fetch').mockResolvedValueOnce({
			ok: true,
			status: 200,
			json: async () => mockResponse
		} as Response);

		const result = await authApi.login({
			email: 'test@example.com',
			password: 'password123'
		});

		expect(result.access_token).toBe('test-token-123');
		expect(result.user.name).toBe('Test Buyer');
		expect(getAuthToken()).toBe('test-token-123');
	});
});
