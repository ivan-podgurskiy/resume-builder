import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

export function generateId(): string {
	return Math.random().toString(36).substring(2, 15);
}

export function formatDate(dateString: string): string {
	if (!dateString) return '';
	const date = new Date(dateString);
	return date.toLocaleDateString('en-US', {
		month: 'short',
		year: 'numeric'
	});
}

export function formatDateRange(startDate: string, endDate?: string, isCurrent?: boolean): string {
	const start = formatDate(startDate);
	if (isCurrent) {
		return `${start} - Present`;
	}
	if (endDate) {
		return `${start} - ${formatDate(endDate)}`;
	}
	return start;
}

export function debounce<T extends (...args: unknown[]) => unknown>(
	func: T,
	wait: number
): (...args: Parameters<T>) => void {
	let timeout: ReturnType<typeof setTimeout> | null = null;

	return function (this: unknown, ...args: Parameters<T>) {
		if (timeout) {
			clearTimeout(timeout);
		}
		timeout = setTimeout(() => {
			func.apply(this, args);
		}, wait);
	};
}

export function truncate(str: string, length: number): string {
	if (str.length <= length) return str;
	return str.slice(0, length) + '...';
}

export function capitalizeFirst(str: string): string {
	return str.charAt(0).toUpperCase() + str.slice(1);
}

export function slugify(str: string): string {
	return str
		.toLowerCase()
		.trim()
		.replace(/[^\w\s-]/g, '')
		.replace(/[\s_-]+/g, '-')
		.replace(/^-+|-+$/g, '');
}
