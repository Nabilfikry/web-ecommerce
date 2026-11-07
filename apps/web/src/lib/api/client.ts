import type { ApiResponse, ApiErrorDetail } from '@ecommerce/shared-types';

export class ApiClientError extends Error {
	status: number;
	code?: string;
	details?: ApiErrorDetail[];

	constructor(message: string, status: number, code?: string, details?: ApiErrorDetail[]) {
		super(message);
		this.name = 'ApiClientError';
		this.status = status;
		this.code = code;
		this.details = details;
	}
}

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

export function getAuthToken(): string | null {
	if (typeof window === 'undefined') return null;
	return localStorage.getItem('ecommerce_access_token');
}

export function setAuthToken(token: string | null): void {
	if (typeof window === 'undefined') return;
	if (token) {
		localStorage.setItem('ecommerce_access_token', token);
	} else {
		localStorage.removeItem('ecommerce_access_token');
	}
}

export interface RequestOptions extends RequestInit {
	params?: Record<string, string | number | boolean | undefined | null>;
}

export async function apiFetch<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
	const { params, headers = {}, ...rest } = options;

	let url = `${API_BASE_URL}${endpoint.startsWith('/') ? endpoint : `/${endpoint}`}`;

	if (params) {
		const searchParams = new URLSearchParams();
		for (const [key, val] of Object.entries(params)) {
			if (val !== undefined && val !== null && val !== '') {
				searchParams.append(key, String(val));
			}
		}
		const queryString = searchParams.toString();
		if (queryString) {
			url += `${url.includes('?') ? '&' : '?'}${queryString}`;
		}
	}

	const reqHeaders: Record<string, string> = {
		Accept: 'application/json',
		...(headers as Record<string, string>)
	};

	if (rest.body && !(rest.body instanceof FormData)) {
		reqHeaders['Content-Type'] = 'application/json';
	}

	const token = getAuthToken();
	if (token && !reqHeaders['Authorization']) {
		reqHeaders['Authorization'] = `Bearer ${token}`;
	}

	let response: Response;
	try {
		response = await fetch(url, {
			...rest,
			headers: reqHeaders,
			credentials: 'include' // include cookies for refresh token
		});
	} catch (err: unknown) {
		const msg = err instanceof Error ? err.message : 'Koneksi ke backend gagal';
		throw new ApiClientError(`Gagal menghubungi server (${msg})`, 0, 'NETWORK_ERROR');
	}

	// 204 No Content
	if (response.status === 204) {
		return null as unknown as T;
	}

	let jsonResult: ApiResponse<T>;
	try {
		jsonResult = await response.json();
	} catch {
		if (!response.ok) {
			throw new ApiClientError(
				`Server error (${response.status} ${response.statusText})`,
				response.status
			);
		}
		return null as unknown as T;
	}

	if (!response.ok || jsonResult.success === false) {
		const errorObj = typeof jsonResult.error === 'object' ? jsonResult.error : undefined;
		const errorCode =
			errorObj?.code || (typeof jsonResult.error === 'string' ? jsonResult.error : undefined);
		const errorDetails = errorObj?.details;
		const message = jsonResult.message || 'Terjadi kesalahan pada sistem';

		throw new ApiClientError(message, response.status, errorCode, errorDetails);
	}

	return (jsonResult.data !== undefined ? jsonResult.data : jsonResult) as T;
}
