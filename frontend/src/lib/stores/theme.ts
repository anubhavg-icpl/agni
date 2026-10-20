import { writable } from 'svelte/store';

export type Theme = 'light' | 'dark';

const KEY = 'agni-theme';

function initial(): Theme {
	if (typeof window === 'undefined') return 'dark';
	const stored = window.localStorage.getItem(KEY);
	if (stored === 'light' || stored === 'dark') return stored;
	return 'dark';
}

function apply(theme: Theme) {
	if (typeof document === 'undefined') return;
	document.documentElement.classList.toggle('dark', theme === 'dark');
}

export const theme = writable<Theme>(initial());

export function setTheme(next: Theme) {
	if (typeof window !== 'undefined') {
		window.localStorage.setItem(KEY, next);
	}
	apply(next);
	theme.set(next);
}

export function toggleTheme() {
	let next: Theme = 'dark';
	theme.update((current) => {
		next = current === 'dark' ? 'light' : 'dark';
		return next;
	});
	setTheme(next);
}
