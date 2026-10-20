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
				const msg = (e as Error).message;
				let friendly = `Sign in failed: ${msg}`;
				if (msg.includes('Failed to fetch')) {
					friendly = 'Cannot reach the Agni server. Check that it is running, then try again.';
				} else if (
					msg.includes('Wrong credentials') ||
					msg.includes('401') ||
					msg.includes('Invalid')
				) {
					friendly = 'Invalid username or password.';
				}
				update((s) => ({ ...s, loading: false, error: friendly }));
				return false;
			}
		},
		async setup(username: string, password: string) {
			update((s) => ({ ...s, loading: true, error: null }));
			try {
				await api.setup(username, password);
				// After setup, sign in with the new account
				const response = await api.login(username, password);
				set({ user: response.user, loading: false, error: null, setupRequired: false });
				return true;
			} catch (e) {
				const msg = (e as Error).message;
				let friendly = `Setup failed: ${msg}`;
				if (msg.includes('already')) {
					friendly = 'An administrator account already exists. Sign in instead.';
				} else if (msg.includes('weak') || msg.includes('8 characters')) {
					friendly = 'Password must be at least 8 characters.';
				} else if (msg.includes('Failed to fetch')) {
					friendly = 'Cannot reach the Agni server. Check that it is running, then try again.';
				}
				update((s) => ({ ...s, loading: false, error: friendly }));
				return false;
			}
		},
		async logout() {
			try {
				await api.logout();
			} catch {
				// Signing out proceeds even if the server call fails
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
