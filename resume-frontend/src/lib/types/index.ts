// User types
export interface User {
	id: string;
	email: string;
	name: string;
	subscription_tier: 'free' | 'pro' | 'enterprise';
	subscription_status: 'active' | 'canceled' | 'past_due' | 'trialing';
	email_verified: boolean;
	created_at: string;
	updated_at: string;
	last_login_at?: string;
}

export interface AuthResponse {
	user: User;
	access_token: string;
	refresh_token: string;
	expires_in: number;
}

// Resume types
export interface Resume {
	id: string;
	user_id: string;
	title: string;
	is_master: boolean;
	template_id?: string;
	data: ResumeData;
	style_config?: Record<string, unknown>;
	is_public: boolean;
	public_slug?: string;
	created_at: string;
	updated_at: string;
	template?: Template;
}

export interface ResumeData {
	personal_info: PersonalInfo;
	summary: string;
	experience: Experience[];
	education: Education[];
	skills: Skills;
	certifications?: Certification[];
	projects?: Project[];
	publications?: Publication[];
	awards?: Award[];
	volunteer?: Volunteer[];
	custom_sections?: CustomSection[];
}

export interface PersonalInfo {
	first_name: string;
	last_name: string;
	title: string;
	email: string;
	phone: string;
	location: string;
	website?: string;
	linkedin?: string;
	github?: string;
	photo_url?: string;
}

export interface Experience {
	id: string;
	company: string;
	position: string;
	location?: string;
	start_date: string;
	end_date?: string;
	is_current: boolean;
	description: string;
	achievements?: string[];
	technologies?: string[];
	order: number;
}

export interface Education {
	id: string;
	institution: string;
	degree: string;
	field_of_study: string;
	location?: string;
	start_date: string;
	end_date?: string;
	is_current: boolean;
	gpa?: string;
	honors?: string;
	description?: string;
	order: number;
}

export interface Skills {
	technical: Skill[];
	soft?: string[];
	languages?: Language[];
}

export interface Skill {
	name: string;
	level?: 'beginner' | 'intermediate' | 'advanced' | 'expert';
}

export interface Language {
	name: string;
	proficiency: 'native' | 'fluent' | 'advanced' | 'intermediate' | 'basic';
}

export interface Certification {
	id: string;
	name: string;
	issuer: string;
	issue_date: string;
	expiry_date?: string;
	credential_id?: string;
	url?: string;
	order: number;
}

export interface Project {
	id: string;
	name: string;
	description: string;
	url?: string;
	repo_url?: string;
	technologies?: string[];
	start_date?: string;
	end_date?: string;
	highlights?: string[];
	order: number;
}

export interface Publication {
	id: string;
	title: string;
	publisher: string;
	date: string;
	url?: string;
	summary?: string;
	order: number;
}

export interface Award {
	id: string;
	title: string;
	issuer: string;
	date: string;
	description?: string;
	order: number;
}

export interface Volunteer {
	id: string;
	organization: string;
	role: string;
	start_date: string;
	end_date?: string;
	is_current: boolean;
	description?: string;
	order: number;
}

export interface CustomSection {
	id: string;
	title: string;
	items: CustomSectionItem[];
	order: number;
}

export interface CustomSectionItem {
	id: string;
	title: string;
	subtitle?: string;
	date?: string;
	description?: string;
	order: number;
}

// Template types
export interface Template {
	id: string;
	name: string;
	category: TemplateCategory;
	description: string;
	preview_image_url: string;
	config: TemplateConfig;
	is_premium: boolean;
	ats_score: number;
	best_for: string[];
	is_active: boolean;
}

export type TemplateCategory =
	| 'modern'
	| 'classic'
	| 'creative'
	| 'tech'
	| 'executive'
	| 'academic'
	| 'minimalist';

export interface TemplateConfig {
	layout: LayoutConfig;
	style: StyleConfig;
	section: SectionConfig;
}

export interface LayoutConfig {
	columns: number;
	photo_position: 'left' | 'right' | 'top' | 'none';
	section_order: string[];
	page_break: 'auto' | 'avoid' | 'force';
}

export interface StyleConfig {
	colors: ColorConfig;
	typography: TypographyConfig;
	spacing: SpacingConfig;
	decoration: DecorationConfig;
}

export interface ColorConfig {
	primary: string;
	secondary: string;
	text: string;
	background: string;
	accent: string;
}

export interface TypographyConfig {
	heading_font: string;
	body_font: string;
	base_font_size: number;
	line_height: number;
	letter_spacing: number;
}

export interface SpacingConfig {
	margins: 'narrow' | 'normal' | 'wide';
	section_gap: number;
	element_gap: number;
	page_padding: number;
}

export interface DecorationConfig {
	dividers: 'line' | 'dots' | 'none';
	bullet_style: 'disc' | 'circle' | 'square' | 'dash';
	use_icons: boolean;
	border_style: 'none' | 'thin' | 'medium';
}

export interface SectionConfig {
	personal: SectionSettings;
	summary: SectionSettings;
	experience: SectionSettings;
	education: SectionSettings;
	skills: SectionSettings;
	projects: SectionSettings;
	certifications: SectionSettings;
}

export interface SectionSettings {
	enabled: boolean;
	label: string;
	icon?: string;
	layout_type: 'list' | 'grid' | 'timeline';
}

// Resume Version
export interface ResumeVersion {
	id: string;
	resume_id: string;
	version_number: number;
	data_snapshot: ResumeData;
	change_description?: string;
	created_at: string;
	created_by: string;
}

// API Response types
export interface ListResumesResponse {
	resumes: Resume[];
	total: number;
	page: number;
	pages: number;
}

export interface APIError {
	error: string;
	details?: string;
}
