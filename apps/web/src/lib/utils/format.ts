/**
 * Formats a number to Indonesian Rupiah (IDR) currency format.
 * Example: 450000 -> "Rp 450.000"
 */
export function formatIDR(amount: number | string | undefined | null): string {
	const numeric = typeof amount === 'string' ? parseFloat(amount) : Number(amount ?? 0);
	if (isNaN(numeric)) return 'Rp 0';

	return new Intl.NumberFormat('id-ID', {
		style: 'currency',
		currency: 'IDR',
		minimumFractionDigits: 0,
		maximumFractionDigits: 0
	})
		.format(numeric)
		.replace(/\s/g, ' ');
}

/**
 * Formats ISO date string to readable Indonesian date.
 * Example: "2026-11-07T03:00:00Z" -> "7 Nov 2026, 10:00"
 */
export function formatDate(
	dateStr: string | undefined | null,
	includeTime: boolean = true
): string {
	if (!dateStr) return '-';
	try {
		const date = new Date(dateStr);
		if (isNaN(date.getTime())) return dateStr;

		return new Intl.DateTimeFormat('id-ID', {
			day: 'numeric',
			month: 'short',
			year: 'numeric',
			...(includeTime
				? {
						hour: '2-digit',
						minute: '2-digit'
					}
				: {})
		}).format(date);
	} catch {
		return dateStr;
	}
}

/**
 * Human readable status labels in Indonesian.
 */
export const STATUS_LABELS: Record<string, string> = {
	pending: 'Menunggu Konfirmasi',
	confirmed: 'Dikonfirmasi',
	processing: 'Sedang Diproses',
	shipped: 'Sedang Dikirim',
	completed: 'Selesai',
	cancelled: 'Dibatalkan',
	rejected: 'Ditolak'
};

export function getStatusLabel(status: string): string {
	return STATUS_LABELS[status] || status;
}
