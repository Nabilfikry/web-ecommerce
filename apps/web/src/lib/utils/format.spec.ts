import { describe, it, expect } from 'vitest';
import { formatIDR, formatDate, getStatusLabel } from './format';

describe('format utilities', () => {
	describe('formatIDR', () => {
		it('formats zero correctly', () => {
			expect(formatIDR(0)).toMatch(/Rp\s*0/);
		});

		it('formats standard integer numbers correctly', () => {
			expect(formatIDR(450000)).toMatch(/Rp\s*450\.000/);
			expect(formatIDR(1500000)).toMatch(/Rp\s*1\.500\.000/);
		});

		it('handles string input safely', () => {
			expect(formatIDR('250000')).toMatch(/Rp\s*250\.000/);
		});

		it('handles undefined or null gracefully', () => {
			expect(formatIDR(undefined)).toMatch(/Rp\s*0/);
			expect(formatIDR(null)).toMatch(/Rp\s*0/);
		});
	});

	describe('formatDate', () => {
		it('returns "-" for empty input', () => {
			expect(formatDate(null)).toBe('-');
			expect(formatDate(undefined)).toBe('-');
		});

		it('formats valid ISO dates', () => {
			const formatted = formatDate('2026-11-07T03:00:00Z', false);
			expect(formatted).toContain('2026');
			expect(formatted).toContain('Nov');
		});
	});

	describe('getStatusLabel', () => {
		it('returns correct label for valid statuses', () => {
			expect(getStatusLabel('pending')).toBe('Menunggu Konfirmasi');
			expect(getStatusLabel('confirmed')).toBe('Dikonfirmasi');
			expect(getStatusLabel('shipped')).toBe('Sedang Dikirim');
			expect(getStatusLabel('completed')).toBe('Selesai');
			expect(getStatusLabel('cancelled')).toBe('Dibatalkan');
		});

		it('falls back to raw status if unknown', () => {
			expect(getStatusLabel('unknown_status')).toBe('unknown_status');
		});
	});
});
