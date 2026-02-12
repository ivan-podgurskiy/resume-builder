<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Button, Card } from '$components/ui';
	import { FileText, Layout, Loader2, Shield, Star } from 'lucide-svelte';
	import { isAuthenticated } from '$stores/auth';
	import { resumeStore } from '$stores/resume';
	import { api } from '$lib/api/client';
	import type { Template } from '$types';

	const API_BASE = '/api/v1';

	function getTemplatePreviewUrl(template: Template): string {
		if (template.preview_image_url) {
			return template.preview_image_url;
		}
		return `${API_BASE}/templates/${template.id}/preview-image`;
	}

	let templates: Template[] = [];
	let isLoading = true;
	let templateFilter = 'all';
	let creatingTemplateId: string | null = null;

	const templateCategories = [
		{ id: 'all', label: 'All Templates' },
		{ id: 'modern', label: 'Modern' },
		{ id: 'classic', label: 'Classic' },
		{ id: 'creative', label: 'Creative' },
		{ id: 'minimalist', label: 'Minimalist' },
		{ id: 'tech', label: 'Tech' },
		{ id: 'executive', label: 'Executive' },
		{ id: 'academic', label: 'Academic' }
	];

	$: filteredTemplates =
		templateFilter === 'all'
			? templates
			: templates.filter((t) => t.category === templateFilter);

	async function handleUseTemplate(template: Template) {
		if ($isAuthenticated) {
			creatingTemplateId = template.id;
			try {
				const resume = await resumeStore.createResume({
					title: 'My Resume',
					template_id: template.id
				});
				goto(`/resumes/${resume.id}/edit`);
			} catch (error) {
				console.error('Failed to create resume:', error);
				creatingTemplateId = null;
			}
		} else {
			goto('/signup?template=' + template.id);
		}
	}

	onMount(async () => {
		try {
			const result = await api.listTemplates();
			templates = result.templates;
		} catch {
			templates = [];
		} finally {
			isLoading = false;
		}
	});
</script>

<svelte:head>
	<title>Resume Templates - Resume Builder</title>
	<meta
		name="description"
		content="Choose from 15+ professional resume templates. Modern, classic, creative, and ATS-friendly designs for every industry."
	/>
</svelte:head>

<div class="min-h-screen bg-background">
	<!-- Navigation -->
	<nav class="sticky top-0 z-50 border-b border-border/40 bg-background/80 backdrop-blur-xl">
		<div class="container mx-auto flex h-16 items-center justify-between px-4">
			<a href="/" class="flex items-center space-x-2 transition-opacity hover:opacity-80">
				<FileText class="h-6 w-6 text-primary" />
				<span class="text-xl font-semibold tracking-tight">ResumeBuilder</span>
			</a>
			<div class="flex items-center gap-6">
				<a
					href="/#templates"
					class="text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
				>
					Home
				</a>
				<a
					href="/about"
					class="text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
				>
					About
				</a>
				{#if $isAuthenticated}
					<a href="/dashboard">
						<Button>Dashboard</Button>
					</a>
				{:else}
					<a href="/login" class="text-sm font-medium text-muted-foreground transition-colors hover:text-foreground">
						Login
					</a>
					<a href="/signup">
						<Button>Get Started Free</Button>
					</a>
				{/if}
			</div>
		</div>
	</nav>

	<!-- Page Header -->
	<section class="border-b py-16">
		<div class="container mx-auto px-4">
			<h1 class="font-display mb-4 text-4xl font-normal tracking-tight sm:text-5xl">
				All Resume Templates
			</h1>
			<p class="max-w-2xl text-lg text-muted-foreground">
				Choose from professional templates designed for every industry. Each template is
				ATS-optimized and fully customizable.
			</p>
		</div>
	</section>

	<!-- Category Filter -->
	<section class="sticky top-16 z-40 border-b bg-background/95 py-4 backdrop-blur-sm">
		<div class="container mx-auto px-4">
			<div class="flex flex-wrap gap-2">
				{#each templateCategories as category}
					<button
						type="button"
						class="rounded-full px-4 py-2 text-sm font-medium transition-colors {templateFilter ===
						category.id
							? 'bg-primary text-primary-foreground'
							: 'bg-muted hover:bg-muted/80 text-muted-foreground hover:text-foreground'}"
						on:click={() => (templateFilter = category.id)}
					>
						{category.label}
					</button>
				{/each}
			</div>
		</div>
	</section>

	<!-- Templates Grid -->
	<section class="py-12">
		<div class="container mx-auto px-4">
			{#if isLoading}
				<div class="flex justify-center py-24">
					<Loader2 class="h-12 w-12 animate-spin text-muted-foreground" />
				</div>
			{:else if filteredTemplates.length === 0}
				<div class="py-24 text-center">
					<Layout class="mx-auto mb-4 h-16 w-16 text-muted-foreground/50" />
					<p class="text-muted-foreground">
						{templateFilter === 'all'
							? 'No templates available.'
							: `No ${templateFilter} templates.`}
					</p>
				</div>
			{:else}
				<div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
					{#each filteredTemplates as template (template.id)}
						<Card
							class="group flex flex-col overflow-hidden transition-all duration-300 hover:shadow-xl hover:shadow-primary/5 hover:border-primary/30"
						>
							<!-- Preview -->
							<div class="relative aspect-[8.5/11] overflow-hidden bg-white">
								<img
									src={getTemplatePreviewUrl(template)}
									alt={template.name}
									class="h-full w-full object-cover object-top transition-transform duration-300 group-hover:scale-[1.02]"
								/>
								<!-- Use Template overlay -->
								<div
									class="absolute inset-0 flex items-center justify-center bg-black/40 opacity-0 transition-opacity group-hover:opacity-100"
								>
									<Button
										size="sm"
										class="shadow-lg"
										disabled={creatingTemplateId === template.id}
										on:click={() => handleUseTemplate(template)}
									>
										{#if creatingTemplateId === template.id}
											<Loader2 class="h-4 w-4 animate-spin" />
										{:else}
											Use Template
										{/if}
									</Button>
								</div>
							</div>
							<!-- Info -->
							<div class="flex flex-1 flex-col p-4">
								<div class="flex items-start justify-between gap-2">
									<h3 class="font-semibold">{template.name}</h3>
									{#if template.is_premium}
										<span
											class="flex shrink-0 items-center gap-0.5 rounded bg-amber-100 px-2 py-0.5 text-xs font-medium text-amber-700"
										>
											<Star class="h-3 w-3" />
											Pro
										</span>
									{/if}
								</div>
								<p class="mt-1 line-clamp-2 flex-1 text-sm text-muted-foreground">
									{template.description}
								</p>
								<div class="mt-3 flex items-center gap-2 text-xs text-muted-foreground">
									<span class="flex items-center gap-1">
										<Shield class="h-3.5 w-3.5" />
										ATS {template.ats_score}%
									</span>
								</div>
							</div>
						</Card>
					{/each}
				</div>
			{/if}
		</div>
	</section>
</div>
