<script lang="ts">
	import type { Resume, Template } from '$types';

	export let resume: Resume;
	export let template: Template | null = null;

	$: data = resume.data;

	// Template-based styles (mirrors the live preview in the editor)
	$: templateStyles = template?.config?.style || null;
	$: templateColors = templateStyles?.colors || {
		primary: '#1a1a1a',
		secondary: '#4a5568',
		text: '#1a1a1a',
		background: '#ffffff',
		accent: '#3b82f6'
	};
	$: templateTypography = templateStyles?.typography || {
		heading_font: 'Inter',
		body_font: 'Inter',
		base_font_size: 14,
		line_height: 1.5
	};
	$: templateSpacing = templateStyles?.spacing || {
		margins: 'normal',
		section_gap: 24,
		element_gap: 12,
		page_padding: 32
	};
	$: templateColumns = template?.config?.layout?.columns ?? 1;
	$: layoutVariant = template?.config?.layout?.layout_variant ?? 'standard';
	$: isSidebarDarkLayout = layoutVariant === 'sidebar_dark';
	$: isTwoColumnLayout = templateColumns === 2 && !isSidebarDarkLayout;

	function getPadding(margins: string): string {
		switch (margins) {
			case 'narrow':
				return '24px';
			case 'wide':
				return '48px';
			default:
				return '32px';
		}
	}
</script>

<div
	class="aspect-[8.5/11] rounded-lg border shadow-lg"
	style="
		font-family: {templateTypography.body_font}, sans-serif;
		font-size: {templateTypography.base_font_size}px;
		line-height: {templateTypography.line_height};
		background-color: {templateColors.background};
		color: {templateColors.text};
		padding: {getPadding(templateSpacing.margins)};
	"
>
	{#if isSidebarDarkLayout}
		<!-- Executive Sidebar: section labels on left align with content on right -->
		<div class="grid grid-cols-[28%_1fr] overflow-hidden rounded-lg">
			<!-- Row 0: Photo | Name + Title -->
			<div class="bg-gray-800 p-4">
				{#if data.personal_info.photo_url}
					<div class="mx-auto h-16 w-16 overflow-hidden rounded-full border-2 border-white/30">
						<img
							src={data.personal_info.photo_url}
							alt=""
							class="h-full w-full object-cover"
						/>
					</div>
				{:else}
					<div class="mx-auto flex h-16 w-16 items-center justify-center rounded-full border-2 border-white/30 bg-gray-700 text-white">
						<span class="text-2xl font-bold text-white/70">
							{data.personal_info.first_name?.charAt(0) || '?'}
							{data.personal_info.last_name?.charAt(0) || ''}
						</span>
					</div>
				{/if}
			</div>
			<div
				class="relative overflow-hidden px-6 pt-4 pb-2"
				style="background-color: {templateColors.background};"
			>
				<div
					class="absolute -top-5 -right-5 h-[120px] w-[120px] opacity-50"
					style="background: radial-gradient(circle at 100% 0%, {templateColors.primary} 0%, transparent 70%);"
				></div>
				<h1
					class="relative text-xl font-bold"
					style="font-family: {templateTypography.heading_font}, sans-serif; color: {templateColors.text};"
				>
					{data.personal_info.first_name}
					{data.personal_info.last_name}
				</h1>
				{#if data.personal_info.title}
					<p class="relative mt-1 text-sm" style="color: {templateColors.primary};">
						{data.personal_info.title}
					</p>
				{/if}
			</div>
			<!-- Row 1: Profil label | Profile content -->
			<div class="border-t border-white/30 bg-gray-800 py-2 pl-4 text-xs font-bold uppercase tracking-widest text-white">
				Profil
			</div>
			<div
				class="border-t px-6 py-3"
				style="background-color: {templateColors.background}; border-color: {templateColors.accent}20;"
			>
				{#if data.summary}
					<p
						class="whitespace-pre-line border-l-2 pl-3 text-sm"
						style="border-color: {templateColors.primary}; color: {templateColors.text};"
					>
						{data.summary}
					</p>
				{:else}
					<p class="text-sm" style="color: {templateColors.secondary}; font-style: italic;">—</p>
				{/if}
			</div>
			<!-- Row 2: Experience label | Experience content -->
			<div class="border-t border-white/30 bg-gray-800 py-2 pl-4 text-xs font-bold uppercase tracking-widest text-white">
				Experience
			</div>
			<div
				class="border-t px-6 py-3"
				style="background-color: {templateColors.background}; border-color: {templateColors.accent}20;"
			>
				{#if (data.experience || []).length > 0}
					{#each data.experience || [] as exp}
						<div class="mb-3 last:mb-0">
							<div class="flex items-start justify-between gap-2">
								<div>
									<h3 class="font-semibold" style="color: {templateColors.text};">{exp.position}</h3>
									<p class="text-sm" style="color: {templateColors.secondary};">{exp.company}</p>
								</div>
								<span class="shrink-0 text-xs" style="color: {templateColors.secondary};">
									{exp.start_date} - {exp.is_current ? 'Present' : exp.end_date}
								</span>
							</div>
							{#if exp.description}
								<p class="mt-1 whitespace-pre-line text-sm" style="color: {templateColors.secondary};">
									{exp.description}
								</p>
							{/if}
						</div>
					{/each}
				{:else}
					<p class="text-sm" style="color: {templateColors.secondary}; font-style: italic;">—</p>
				{/if}
			</div>
			<!-- Row 3: Education label | Education content -->
			<div class="border-t border-white/30 bg-gray-800 py-2 pl-4 text-xs font-bold uppercase tracking-widest text-white">
				Education
			</div>
			<div
				class="border-t px-6 py-3"
				style="background-color: {templateColors.background}; border-color: {templateColors.accent}20;"
			>
				{#if (data.education || []).length > 0}
					{#each data.education || [] as edu}
						<div class="mb-3 last:mb-0">
							<div class="flex items-start justify-between gap-2">
								<div>
									<h3 class="font-semibold" style="color: {templateColors.text};">
										{edu.degree}
										{#if edu.field_of_study}
											in {edu.field_of_study}
										{/if}
									</h3>
									<p class="text-sm" style="color: {templateColors.secondary};">{edu.institution}</p>
								</div>
								<span class="shrink-0 text-xs" style="color: {templateColors.secondary};">
									{edu.start_date} - {edu.end_date}
								</span>
							</div>
						</div>
					{/each}
				{:else}
					<p class="text-sm" style="color: {templateColors.secondary}; font-style: italic;">—</p>
				{/if}
			</div>
		</div>
	{:else if isTwoColumnLayout}
		<div class="grid h-full grid-cols-[1fr_2fr] gap-4">
			<div class="flex flex-col gap-4 border-r pr-4" style="border-color: {templateColors.accent}30;">
				<div>
					<h1 class="text-xl font-bold" style="font-family: {templateTypography.heading_font}, sans-serif; color: {templateColors.primary};">
						{data.personal_info.first_name} {data.personal_info.last_name}
					</h1>
					{#if data.personal_info.title}
						<p class="mt-0.5 text-sm" style="color: {templateColors.secondary};">
							{data.personal_info.title}
						</p>
					{/if}
				</div>
				<div class="flex flex-col gap-1 text-xs" style="color: {templateColors.secondary};">
					{#if data.personal_info.email}<span>{data.personal_info.email}</span>{/if}
					{#if data.personal_info.phone}<span>{data.personal_info.phone}</span>{/if}
					{#if data.personal_info.location}<span>{data.personal_info.location}</span>{/if}
				</div>
				{#if (data.skills.technical || []).length > 0}
					<div>
						<h2 class="mb-1.5 text-xs font-bold uppercase" style="color: {templateColors.accent};">Skills</h2>
						<div class="flex flex-wrap gap-1">
							{#each data.skills.technical || [] as skill}
								<span class="rounded px-1.5 py-0.5 text-xs" style="background-color: {templateColors.accent}20; color: {templateColors.primary};">
									{skill.name}
								</span>
							{/each}
						</div>
					</div>
				{/if}
				{#if (data.education || []).length > 0}
					<div>
						<h2 class="mb-1.5 text-xs font-bold uppercase" style="color: {templateColors.accent};">Education</h2>
						{#each data.education || [] as edu}
							<div class="mb-1">
								<div class="font-semibold" style="color: {templateColors.primary};">{edu.degree}</div>
								<div class="text-xs" style="color: {templateColors.secondary};">{edu.institution}</div>
								<div class="text-xs" style="color: {templateColors.secondary};">{edu.start_date} - {edu.end_date}</div>
							</div>
						{/each}
					</div>
				{/if}
			</div>
			<div class="min-w-0 overflow-hidden">
				{#if data.summary}
					<div class="mb-4">
						<h2 class="mb-1.5 text-xs font-bold uppercase" style="color: {templateColors.accent};">Professional Summary</h2>
						<p class="whitespace-pre-line text-sm">{data.summary}</p>
					</div>
				{/if}
				{#if (data.experience || []).length > 0}
					<div>
						<h2 class="mb-2 text-xs font-bold uppercase" style="color: {templateColors.accent};">Experience</h2>
						{#each data.experience || [] as exp}
							<div class="mb-3">
								<div class="flex items-start justify-between gap-2">
									<div>
										<h3 class="font-semibold" style="color: {templateColors.primary};">{exp.position}</h3>
										<p class="text-sm" style="color: {templateColors.secondary};">{exp.company}</p>
									</div>
									<span class="shrink-0 text-xs" style="color: {templateColors.secondary};">
										{exp.start_date} - {exp.is_current ? 'Present' : exp.end_date}
									</span>
								</div>
								{#if exp.description}
									<p class="mt-1 whitespace-pre-line text-sm">{exp.description}</p>
								{/if}
							</div>
						{/each}
					</div>
				{/if}
			</div>
		</div>
	{:else}
		<div class="mb-6 pb-4 text-center" style="border-bottom: 2px solid {templateColors.accent};">
			<h1 class="text-2xl font-bold" style="font-family: {templateTypography.heading_font}, sans-serif; color: {templateColors.primary};">
				{data.personal_info.first_name} {data.personal_info.last_name}
			</h1>
			{#if data.personal_info.title}
				<p class="mt-1 text-lg" style="color: {templateColors.secondary};">
					{data.personal_info.title}
				</p>
			{/if}
			<div class="mt-2 flex flex-wrap justify-center gap-3 text-sm" style="color: {templateColors.secondary};">
				{#if data.personal_info.email}<span>{data.personal_info.email}</span>{/if}
				{#if data.personal_info.phone}<span>{data.personal_info.phone}</span>{/if}
				{#if data.personal_info.location}<span>{data.personal_info.location}</span>{/if}
			</div>
		</div>
		{#if data.summary}
			<div style="margin-bottom: {templateSpacing.section_gap}px;">
				<h2 class="mb-2 text-sm font-bold uppercase tracking-wide" style="color: {templateColors.accent};">Professional Summary</h2>
				<p class="whitespace-pre-line text-sm">{data.summary}</p>
			</div>
		{/if}
		{#if (data.experience || []).length > 0}
			<div style="margin-bottom: {templateSpacing.section_gap}px;">
				<h2 class="mb-3 text-sm font-bold uppercase tracking-wide" style="color: {templateColors.accent};">Experience</h2>
				{#each data.experience || [] as exp}
					<div style="margin-bottom: {templateSpacing.element_gap}px;">
						<div class="flex items-start justify-between">
							<div>
								<h3 class="font-semibold" style="color: {templateColors.primary};">{exp.position}</h3>
								<p class="text-sm" style="color: {templateColors.secondary};">{exp.company}</p>
							</div>
							<span class="text-sm" style="color: {templateColors.secondary};">
								{exp.start_date} - {exp.is_current ? 'Present' : exp.end_date}
							</span>
						</div>
						{#if exp.description}
							<p class="mt-1 whitespace-pre-line text-sm">{exp.description}</p>
						{/if}
					</div>
				{/each}
			</div>
		{/if}
		{#if (data.education || []).length > 0}
			<div style="margin-bottom: {templateSpacing.section_gap}px;">
				<h2 class="mb-3 text-sm font-bold uppercase tracking-wide" style="color: {templateColors.accent};">Education</h2>
				{#each data.education || [] as edu}
					<div style="margin-bottom: {templateSpacing.element_gap}px;">
						<div class="flex items-start justify-between">
							<div>
								<h3 class="font-semibold" style="color: {templateColors.primary};">
									{edu.degree} in {edu.field_of_study}
								</h3>
								<p class="text-sm" style="color: {templateColors.secondary};">{edu.institution}</p>
							</div>
							<span class="text-sm" style="color: {templateColors.secondary};">
								{edu.start_date} - {edu.end_date}
							</span>
						</div>
					</div>
				{/each}
			</div>
		{/if}
		{#if (data.skills.technical || []).length > 0}
			<div>
				<h2 class="mb-2 text-sm font-bold uppercase tracking-wide" style="color: {templateColors.accent};">Skills</h2>
				<div class="flex flex-wrap gap-2">
					{#each data.skills.technical || [] as skill}
						<span class="rounded px-2 py-1 text-sm" style="background-color: {templateColors.accent}20; color: {templateColors.primary};">
							{skill.name}
						</span>
					{/each}
				</div>
			</div>
		{/if}
	{/if}
</div>
