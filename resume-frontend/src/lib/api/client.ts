import type {
	AuthResponse,
	Resume,
	ListResumesResponse,
	Template,
	ResumeVersion,
	ResumeData
} from '$types';

const API_BASE = '/api/v1';

class APIClient {
	private accessToken: string | null = null;
	private refreshToken: string | null = null;

	constructor() {
		// Load tokens from localStorage if available (client-side only)
		if (typeof window !== 'undefined') {
			this.accessToken = localStorage.getItem('access_token');
			this.refreshToken = localStorage.getItem('refresh_token');
		}
	}

	setTokens(accessToken: string, refreshToken: string) {
		this.accessToken = accessToken;
		this.refreshToken = refreshToken;
		if (typeof window !== 'undefined') {
			localStorage.setItem('access_token', accessToken);
			localStorage.setItem('refresh_token', refreshToken);
		}
	}

	clearTokens() {
		this.accessToken = null;
		this.refreshToken = null;
		if (typeof window !== 'undefined') {
			localStorage.removeItem('access_token');
			localStorage.removeItem('refresh_token');
		}
	}

	private async request<T>(
		endpoint: string,
		options: RequestInit = {}
	): Promise<T> {
		const headers: Record<string, string> = {
			'Content-Type': 'application/json',
			...(options.headers as Record<string, string>)
		};

		if (this.accessToken) {
			headers['Authorization'] = `Bearer ${this.accessToken}`;
		}

		const response = await fetch(`${API_BASE}${endpoint}`, {
			...options,
			headers
		});

		if (response.status === 401 && this.refreshToken) {
			// Try to refresh the token
			const refreshed = await this.refreshAccessToken();
			if (refreshed) {
				headers['Authorization'] = `Bearer ${this.accessToken}`;
				const retryResponse = await fetch(`${API_BASE}${endpoint}`, {
					...options,
					headers
				});
				if (!retryResponse.ok) {
					throw await this.handleError(retryResponse);
				}
				return retryResponse.json();
			}
		}

		if (!response.ok) {
			throw await this.handleError(response);
		}

		// Handle 204 No Content
		if (response.status === 204) {
			return {} as T;
		}

		return response.json();
	}

	private async handleError(response: Response): Promise<Error> {
		try {
			const data = await response.json();
			return new Error(data.error || 'An error occurred');
		} catch {
			return new Error(`HTTP ${response.status}: ${response.statusText}`);
		}
	}

	private async refreshAccessToken(): Promise<boolean> {
		if (!this.refreshToken) return false;

		try {
			const response = await fetch(`${API_BASE}/auth/refresh`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ refresh_token: this.refreshToken })
			});

			if (!response.ok) {
				this.clearTokens();
				return false;
			}

			const data: AuthResponse = await response.json();
			this.setTokens(data.access_token, data.refresh_token);
			return true;
		} catch {
			this.clearTokens();
			return false;
		}
	}

	// Auth endpoints
	async register(email: string, password: string, name: string): Promise<AuthResponse> {
		const data = await this.request<AuthResponse>('/auth/register', {
			method: 'POST',
			body: JSON.stringify({ email, password, name })
		});
		this.setTokens(data.access_token, data.refresh_token);
		return data;
	}

	async login(email: string, password: string): Promise<AuthResponse> {
		const data = await this.request<AuthResponse>('/auth/login', {
			method: 'POST',
			body: JSON.stringify({ email, password })
		});
		this.setTokens(data.access_token, data.refresh_token);
		return data;
	}

	async logout(): Promise<void> {
		await this.request('/auth/logout', { method: 'POST' });
		this.clearTokens();
	}

	async getCurrentUser() {
		return this.request<{ user: AuthResponse['user'] }>('/auth/me');
	}

	async forgotPassword(email: string): Promise<{ message: string }> {
		return this.request('/auth/forgot-password', {
			method: 'POST',
			body: JSON.stringify({ email })
		});
	}

	async resetPassword(token: string, newPassword: string): Promise<{ message: string }> {
		return this.request('/auth/reset-password', {
			method: 'POST',
			body: JSON.stringify({ token, new_password: newPassword })
		});
	}

	async verifyEmail(token: string): Promise<{ message: string }> {
		return this.request('/auth/verify-email', {
			method: 'POST',
			body: JSON.stringify({ token })
		});
	}

	// Resume endpoints
	async listResumes(page = 1, pageSize = 10): Promise<ListResumesResponse> {
		return this.request(`/resumes?page=${page}&page_size=${pageSize}`);
	}

	async getResume(id: string): Promise<Resume> {
		return this.request(`/resumes/${id}`);
	}

	async createResume(data: {
		title: string;
		template_id?: string;
		data?: ResumeData;
		is_master?: boolean;
	}): Promise<Resume> {
		return this.request('/resumes', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	async updateResume(
		id: string,
		data: Partial<{
			title: string;
			template_id: string;
			data: ResumeData;
			style_config: Record<string, unknown>;
		}>
	): Promise<Resume> {
		return this.request(`/resumes/${id}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		});
	}

	async deleteResume(id: string): Promise<void> {
		return this.request(`/resumes/${id}`, { method: 'DELETE' });
	}

	async duplicateResume(id: string, title?: string): Promise<Resume> {
		return this.request(`/resumes/${id}/duplicate`, {
			method: 'POST',
			body: JSON.stringify({ title })
		});
	}

	async getResumeVersions(id: string): Promise<{ versions: ResumeVersion[] }> {
		return this.request(`/resumes/${id}/versions`);
	}

	async restoreResumeVersion(resumeId: string, versionId: string): Promise<Resume> {
		return this.request(`/resumes/${resumeId}/versions/${versionId}/restore`, {
			method: 'POST'
		});
	}

	async setResumeVisibility(
		id: string,
		isPublic: boolean,
		slug?: string
	): Promise<Resume> {
		return this.request(`/resumes/${id}/visibility`, {
			method: 'PATCH',
			body: JSON.stringify({ is_public: isPublic, slug })
		});
	}

	async getPublicResume(slugOrId: string): Promise<Resume> {
		return this.request(`/share/${slugOrId}`);
	}

	// Template endpoints
	async listTemplates(): Promise<{ templates: Template[] }> {
		return this.request('/templates');
	}

	async getTemplate(id: string): Promise<Template> {
		return this.request(`/templates/${id}`);
	}

	async getTemplatePreview(id: string): Promise<{ id: string; name: string; config: Template['config'] }> {
		return this.request(`/templates/${id}/preview`);
	}

	// AI endpoints
	async extractResumeData(text: string): Promise<ResumeData> {
		const result = await this.request<{ data: ResumeData; confidence: Record<string, number>; warnings?: string[] }>('/ai/extract', {
			method: 'POST',
			body: JSON.stringify({ text })
		});
		return result.data;
	}

	async improveText(text: string, context?: string): Promise<{ improved_text: string }> {
		return this.request('/ai/improve', {
			method: 'POST',
			body: JSON.stringify({ text, context })
		});
	}

	async generateSummary(data: ResumeData): Promise<{ summary: string }> {
		return this.request('/ai/generate-summary', {
			method: 'POST',
			body: JSON.stringify({ data })
		});
	}

	// Export endpoints (return blobs/text, not JSON)
	async exportPDF(resumeId: string, quality: string = 'print'): Promise<Blob> {
		const headers: Record<string, string> = {
			'Content-Type': 'application/json'
		};
		if (this.accessToken) {
			headers['Authorization'] = `Bearer ${this.accessToken}`;
		}

		const response = await fetch(`${API_BASE}/export/pdf`, {
			method: 'POST',
			headers,
			body: JSON.stringify({ resume_id: resumeId, quality })
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ error: 'Export failed' }));
			throw new Error(error.error || 'Failed to export PDF');
		}

		return response.blob();
	}

	async exportTXT(resumeId: string): Promise<string> {
		const headers: Record<string, string> = {
			'Content-Type': 'application/json'
		};
		if (this.accessToken) {
			headers['Authorization'] = `Bearer ${this.accessToken}`;
		}

		const response = await fetch(`${API_BASE}/export/txt`, {
			method: 'POST',
			headers,
			body: JSON.stringify({ resume_id: resumeId })
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ error: 'Export failed' }));
			throw new Error(error.error || 'Failed to export TXT');
		}

		return response.text();
	}

	async exportJSON(resumeId: string): Promise<string> {
		const headers: Record<string, string> = {
			'Content-Type': 'application/json'
		};
		if (this.accessToken) {
			headers['Authorization'] = `Bearer ${this.accessToken}`;
		}

		const response = await fetch(`${API_BASE}/export/json`, {
			method: 'POST',
			headers,
			body: JSON.stringify({ resume_id: resumeId })
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ error: 'Export failed' }));
			throw new Error(error.error || 'Failed to export JSON');
		}

		return response.text();
	}

	// Helper methods
	isAuthenticated(): boolean {
		return !!this.accessToken;
	}
}

export const api = new APIClient();
