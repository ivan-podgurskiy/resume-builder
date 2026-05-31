import { describe, it, expect, vi } from 'vitest';
import {
	generateId,
	formatDate,
	formatDateRange,
	debounce,
	truncate,
	capitalizeFirst,
	slugify
} from './index';

describe('generateId', () => {
	it('returns a non-empty string', () => {
		expect(generateId().length).toBeGreaterThan(0);
	});

	it('returns distinct values across calls', () => {
		const ids = new Set(Array.from({ length: 100 }, () => generateId()));
		expect(ids.size).toBeGreaterThan(90);
	});
});

describe('formatDate', () => {
	it('returns an empty string for empty input', () => {
		expect(formatDate('')).toBe('');
	});

	it('formats an ISO date as "Mon YYYY"', () => {
		expect(formatDate('2023-06-15')).toBe('Jun 2023');
	});
});

describe('formatDateRange', () => {
	it('uses "Present" when isCurrent is true', () => {
		expect(formatDateRange('2020-01-01', undefined, true)).toBe('Jan 2020 - Present');
	});

	it('joins start and end dates', () => {
		expect(formatDateRange('2020-01-01', '2022-03-01')).toBe('Jan 2020 - Mar 2022');
	});

	it('returns only the start date when no end date is given', () => {
		expect(formatDateRange('2020-01-01')).toBe('Jan 2020');
	});
});

describe('truncate', () => {
	it('leaves short strings unchanged', () => {
		expect(truncate('hello', 10)).toBe('hello');
	});

	it('truncates and appends an ellipsis', () => {
		expect(truncate('hello world', 5)).toBe('hello...');
	});
});

describe('capitalizeFirst', () => {
	it('capitalizes only the first character', () => {
		expect(capitalizeFirst('hello')).toBe('Hello');
	});

	it('handles an empty string', () => {
		expect(capitalizeFirst('')).toBe('');
	});
});

describe('slugify', () => {
	it('lowercases and hyphenates', () => {
		expect(slugify('Senior Software Engineer')).toBe('senior-software-engineer');
	});

	it('strips punctuation and collapses separators', () => {
		expect(slugify('  Jane_Doe!! Résumé  ')).toBe('jane-doe-rsum');
	});

	it('trims leading and trailing hyphens', () => {
		expect(slugify('--hello--')).toBe('hello');
	});
});

describe('debounce', () => {
	it('invokes the function only once after rapid calls', () => {
		vi.useFakeTimers();
		const fn = vi.fn();
		const debounced = debounce(fn, 100);

		debounced();
		debounced();
		debounced();
		expect(fn).not.toHaveBeenCalled();

		vi.advanceTimersByTime(100);
		expect(fn).toHaveBeenCalledTimes(1);
		vi.useRealTimers();
	});

	it('passes the latest arguments through', () => {
		vi.useFakeTimers();
		const fn = vi.fn();
		const debounced = debounce(fn, 50);

		debounced('a');
		debounced('b');
		vi.advanceTimersByTime(50);

		expect(fn).toHaveBeenCalledWith('b');
		vi.useRealTimers();
	});
});
