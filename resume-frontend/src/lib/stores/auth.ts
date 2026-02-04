import { writable, derived } from 'svelte/store';
import type { User } from '$types';
import { api } from '$api/client';

interface AuthState {
	user: User | null;
	loading: boolean;
	initialized: boolean;
}

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>({
		user: null,
		loading: true,
		initialized: false
	});

	return {
		subscribe,

		async initialize() {
			if (!api.isAuthenticated()) {
				set({ user: null, loading: false, initialized: true });
				return;
			}

			try {
				const { user } = await api.getCurrentUser();
				set({ user, loading: false, initialized: true });
			} catch {
				api.clearTokens();
				set({ user: null, loading: false, initialized: true });
			}
		},

		async login(email: string, password: string) {
			update((state) => ({ ...state, loading: true }));
			try {
				const response = await api.login(email, password);
				set({ user: response.user, loading: false, initialized: true });
				return { success: true };
			} catch (error) {
				update((state) => ({ ...state, loading: false }));
				return { success: false, error: (error as Error).message };
			}
		},

		async register(email: string, password: string, name: string) {
			update((state) => ({ ...state, loading: true }));
			try {
				const response = await api.register(email, password, name);
				set({ user: response.user, loading: false, initialized: true });
				return { success: true };
			} catch (error) {
				update((state) => ({ ...state, loading: false }));
				return { success: false, error: (error as Error).message };
			}
		},

		async logout() {
			try {
				await api.logout();
			} finally {
				set({ user: null, loading: false, initialized: true });
			}
		},

		setUser(user: User | null) {
			update((state) => ({ ...state, user }));
		}
	};
}

export const auth = createAuthStore();

export const isAuthenticated = derived(auth, ($auth) => !!$auth.user);
export const currentUser = derived(auth, ($auth) => $auth.user);
export const isLoading = derived(auth, ($auth) => $auth.loading);
