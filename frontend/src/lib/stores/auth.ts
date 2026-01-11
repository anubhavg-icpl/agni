import { writable, derived } from 'svelte/store';
import { api, type User } from '$lib/api/client';

interface AuthState {
	user: User | null;
	loading: boolean;
	error: string | null;
	setupRequired: boolean;
}

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>({
		user: null,
		loading: true,
		error: null,
		setupRequired: false
	});

	return {
		subscribe,
		async init() {
			update((s) => ({ ...s, loading: true, error: null }));
			try {
				const status = await api.getAuthStatus();
				if (status.setup_required) {
					set({ user: null, loading: false, error: null, setupRequired: true });
					return;
				}

				const token = api.getToken();
				if (token) {
					const user = await api.getMe();
					set({ user, loading: false, error: null, setupRequired: false });
				} else {
					set({ user: null, loading: false, error: null, setupRequired: false });
				}
			} catch {
				set({ user: null, loading: false, error: null, setupRequired: false });
			}
		},
		async login(username: string, password: string) {
			update((s) => ({ ...s, loading: true, error: null }));
			try {
				const response = await api.login(username, password);
				set({ user: response.user, loading: false, error: null, setupRequired: false });
				return true;
			} catch (e) {
				update((s) => ({ ...s, loading: false, error: (e as Error).message }));
				return false;
			}
		},
		async setup(username: string, password: string) {
			update((s) => ({ ...s, loading: true, error: null }));
			try {
				await api.setup(username, password);
				// After setup, login
				const response = await api.login(username, password);
				set({ user: response.user, loading: false, error: null, setupRequired: false });
				return true;
			} catch (e) {
				update((s) => ({ ...s, loading: false, error: (e as Error).message }));
				return false;
			}
		},
		async logout() {
			try {
				await api.logout();
			} catch {
				// Ignore
			}
			api.setToken(null);
			set({ user: null, loading: false, error: null, setupRequired: false });
		},
		clearError() {
			update((s) => ({ ...s, error: null }));
		}
	};
}

export const auth = createAuthStore();
export const isAuthenticated = derived(auth, ($auth) => $auth.user !== null);
