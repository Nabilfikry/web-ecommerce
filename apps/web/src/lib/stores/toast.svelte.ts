export interface Toast {
	id: string;
	message: string;
	type: 'success' | 'error' | 'info' | 'warning';
	duration?: number;
}

class ToastStore {
	toasts = $state<Toast[]>([]);

	add(message: string, type: Toast['type'] = 'info', duration = 3500): string {
		const id = Math.random().toString(36).substring(2, 9);
		const newToast: Toast = { id, message, type, duration };

		this.toasts = [...this.toasts, newToast];

		if (duration > 0) {
			setTimeout(() => {
				this.remove(id);
			}, duration);
		}

		return id;
	}

	remove(id: string) {
		this.toasts = this.toasts.filter((t) => t.id !== id);
	}

	success(message: string, duration?: number) {
		return this.add(message, 'success', duration);
	}

	error(message: string, duration = 5000) {
		return this.add(message, 'error', duration);
	}

	info(message: string, duration?: number) {
		return this.add(message, 'info', duration);
	}

	warning(message: string, duration?: number) {
		return this.add(message, 'warning', duration);
	}
}

export const toast = new ToastStore();
