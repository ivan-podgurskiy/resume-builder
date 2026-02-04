import { writable, derived } from 'svelte/store';
import type { Resume, ResumeData } from '$types';
import { api } from '$api/client';

interface ResumeState {
	resumes: Resume[];
	currentResume: Resume | null;
	loading: boolean;
	saving: boolean;
	error: string | null;
}

function createResumeStore() {
	const { subscribe, set, update } = writable<ResumeState>({
		resumes: [],
		currentResume: null,
		loading: false,
		saving: false,
		error: null
	});

	let saveTimeout: ReturnType<typeof setTimeout> | null = null;

	return {
		subscribe,

		async loadResumes(page = 1, pageSize = 10) {
			update((state) => ({ ...state, loading: true, error: null }));
			try {
				const response = await api.listResumes(page, pageSize);
				update((state) => ({
					...state,
					resumes: response.resumes,
					loading: false
				}));
				return response;
			} catch (error) {
				update((state) => ({
					...state,
					loading: false,
					error: (error as Error).message
				}));
				throw error;
			}
		},

		async loadResume(id: string) {
			update((state) => ({ ...state, loading: true, error: null }));
			try {
				const resume = await api.getResume(id);
				update((state) => ({
					...state,
					currentResume: resume,
					loading: false
				}));
				return resume;
			} catch (error) {
				update((state) => ({
					...state,
					loading: false,
					error: (error as Error).message
				}));
				throw error;
			}
		},

		async createResume(data: {
			title: string;
			template_id?: string;
			data?: ResumeData;
		}) {
			update((state) => ({ ...state, loading: true, error: null }));
			try {
				const resume = await api.createResume(data);
				update((state) => ({
					...state,
					resumes: [resume, ...state.resumes],
					currentResume: resume,
					loading: false
				}));
				return resume;
			} catch (error) {
				update((state) => ({
					...state,
					loading: false,
					error: (error as Error).message
				}));
				throw error;
			}
		},

		async updateResume(id: string, data: Partial<Resume>) {
			update((state) => ({ ...state, saving: true, error: null }));
			try {
				const resume = await api.updateResume(id, data);
				update((state) => ({
					...state,
					currentResume: resume,
					resumes: state.resumes.map((r) => (r.id === id ? resume : r)),
					saving: false
				}));
				return resume;
			} catch (error) {
				update((state) => ({
					...state,
					saving: false,
					error: (error as Error).message
				}));
				throw error;
			}
		},

		// Debounced auto-save (5 seconds)
		scheduleAutoSave(id: string, data: Partial<Resume>) {
			if (saveTimeout) {
				clearTimeout(saveTimeout);
			}

			update((state) => ({ ...state, saving: true }));

			saveTimeout = setTimeout(async () => {
				try {
					const resume = await api.updateResume(id, data);
					update((state) => ({
						...state,
						currentResume: resume,
						resumes: state.resumes.map((r) => (r.id === id ? resume : r)),
						saving: false
					}));
				} catch (error) {
					update((state) => ({
						...state,
						saving: false,
						error: (error as Error).message
					}));
				}
			}, 5000);
		},

		// Update local state immediately (optimistic)
		updateLocalResume(data: Partial<ResumeData>) {
			update((state) => {
				if (!state.currentResume) return state;
				return {
					...state,
					currentResume: {
						...state.currentResume,
						data: {
							...state.currentResume.data,
							...data
						}
					}
				};
			});
		},

		async deleteResume(id: string) {
			update((state) => ({ ...state, loading: true, error: null }));
			try {
				await api.deleteResume(id);
				update((state) => ({
					...state,
					resumes: state.resumes.filter((r) => r.id !== id),
					currentResume: state.currentResume?.id === id ? null : state.currentResume,
					loading: false
				}));
			} catch (error) {
				update((state) => ({
					...state,
					loading: false,
					error: (error as Error).message
				}));
				throw error;
			}
		},

		async duplicateResume(id: string, title?: string) {
			update((state) => ({ ...state, loading: true, error: null }));
			try {
				const resume = await api.duplicateResume(id, title);
				update((state) => ({
					...state,
					resumes: [resume, ...state.resumes],
					loading: false
				}));
				return resume;
			} catch (error) {
				update((state) => ({
					...state,
					loading: false,
					error: (error as Error).message
				}));
				throw error;
			}
		},

		setCurrentResume(resume: Resume | null) {
			update((state) => ({ ...state, currentResume: resume }));
		},

		clearError() {
			update((state) => ({ ...state, error: null }));
		}
	};
}

export const resumeStore = createResumeStore();

export const currentResume = derived(resumeStore, ($store) => $store.currentResume);
export const resumes = derived(resumeStore, ($store) => $store.resumes);
export const isLoading = derived(resumeStore, ($store) => $store.loading);
export const isSaving = derived(resumeStore, ($store) => $store.saving);
export const resumeError = derived(resumeStore, ($store) => $store.error);
