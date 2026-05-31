<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Button, Card } from '$components/ui';
	import { resumeStore, resumes, isLoading } from '$stores/resume';
	import { currentUser } from '$stores/auth';
	import {
		Plus,
		FileText,
		Upload,
		MoreVertical,
		Copy,
		Trash2,
		ExternalLink,
		Clock,
		Loader2,
		File as FileIcon,
		X,
		Layout,
		Check
	} from 'lucide-svelte';
	import { formatDate } from '$utils';
	import { api } from '$lib/api/client';
	import type { Template } from '$types';

	const API_BASE = '/api/v1';

	function getTemplatePreviewUrl(template: Template): string {
		if (template.preview_image_url) return template.preview_image_url;
		return `${API_BASE}/templates/${template.id}/preview-image`;
	}

	let showNewResumeModal = false;
	let showImportModal = false;
	let newResumeTitle = 'My Resume';
	let importText = '';
	let importTitle = 'Imported Resume';
	let isImporting = false;
	let importError = '';
	let activeMenu: string | null = null;

	// Template picker state
	let templates: Template[] = [];
	let selectedTemplateId: string | null = null;
	let isLoadingTemplates = false;

	// File upload state
	let importMode: 'text' | 'file' = 'file';
	let selectedFile: File | null = null;
	let isDragging = false;
	let fileInputRef: HTMLInputElement;

	const supportedFormats = ['.pdf', '.docx', '.doc', '.txt'];
	const acceptedFileTypes = supportedFormats.join(',');

	async function loadTemplates() {
		if (templates.length === 0) {
			isLoadingTemplates = true;
			try {
				const result = await api.listTemplates();
				templates = result.templates;
			} catch {
				templates = [];
			} finally {
				isLoadingTemplates = false;
			}
		}
	}

	onMount(() => {
		resumeStore.loadResumes();

		// Handle ?template= param (e.g. from signup after selecting template)
		const templateId = $page.url.searchParams.get('template');
		if (templateId) {
			selectedTemplateId = templateId;
			showNewResumeModal = true;
			loadTemplates();
			// Clear the param from URL without navigating
			window.history.replaceState({}, '', '/dashboard');
		}
	});

	async function openNewResumeModal() {
		showNewResumeModal = true;
		selectedTemplateId = null;
		await loadTemplates();
	}

	async function createNewResume() {
		try {
			const resume = await resumeStore.createResume({
				title: newResumeTitle,
				template_id: selectedTemplateId || undefined
			});
			showNewResumeModal = false;
			goto(`/resumes/${resume.id}/edit`);
		} catch (error) {
			console.error('Failed to create resume:', error);
		}
	}

	function handleFileSelect(event: Event) {
		const input = event.target as HTMLInputElement;
		if (input.files && input.files.length > 0) {
			validateAndSetFile(input.files[0]);
		}
	}

	function handleDrop(event: DragEvent) {
		event.preventDefault();
		isDragging = false;
		
		if (event.dataTransfer?.files && event.dataTransfer.files.length > 0) {
			validateAndSetFile(event.dataTransfer.files[0]);
		}
	}

	function handleDragOver(event: DragEvent) {
		event.preventDefault();
		isDragging = true;
	}

	function handleDragLeave() {
		isDragging = false;
	}

	function validateAndSetFile(file: File) {
		importError = '';
		
		// Check file extension
		const ext = '.' + file.name.split('.').pop()?.toLowerCase();
		if (!supportedFormats.includes(ext)) {
			importError = `Unsupported file type. Please upload: ${supportedFormats.join(', ')}`;
			return;
		}
		
		// Check file size (10MB max)
		if (file.size > 10 * 1024 * 1024) {
			importError = 'File too large. Maximum size is 10MB.';
			return;
		}
		
		selectedFile = file;
		
		// Auto-set title from filename
		const nameWithoutExt = file.name.replace(/\.[^/.]+$/, '');
		if (importTitle === 'Imported Resume') {
			importTitle = nameWithoutExt;
		}
	}

	function removeFile() {
		selectedFile = null;
		if (fileInputRef) {
			fileInputRef.value = '';
		}
	}

	function formatFileSize(bytes: number): string {
		if (bytes < 1024) return bytes + ' B';
		if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
		return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
	}

	async function importResume() {
		if (importMode === 'file') {
			if (!selectedFile) {
				importError = 'Please select a file to upload';
				return;
			}
		} else {
			if (!importText.trim()) {
				importError = 'Please paste your resume text';
				return;
			}
		}

		isImporting = true;
		importError = '';

		try {
			// Extract data using AI
			let extractedData;
			if (importMode === 'file' && selectedFile) {
				extractedData = await api.extractResumeFromFile(selectedFile);
			} else {
				extractedData = await api.extractResumeData(importText);
			}

			// Create a new resume with the extracted data
			const resume = await resumeStore.createResume({
				title: importTitle,
				data: extractedData
			});

			showImportModal = false;
			importText = '';
			importTitle = 'Imported Resume';
			selectedFile = null;
			goto(`/resumes/${resume.id}/edit`);
		} catch (error) {
			console.error('Failed to import resume:', error);
			importError = error instanceof Error ? error.message : 'Failed to import resume. Please try again.';
		} finally {
			isImporting = false;
		}
	}

	async function duplicateResume(id: string, title: string) {
		try {
			await resumeStore.duplicateResume(id, `${title} (Copy)`);
			activeMenu = null;
		} catch (error) {
			console.error('Failed to duplicate resume:', error);
		}
	}

	async function deleteResume(id: string) {
		if (confirm('Are you sure you want to delete this resume?')) {
			try {
				await resumeStore.deleteResume(id);
				activeMenu = null;
			} catch (error) {
				console.error('Failed to delete resume:', error);
			}
		}
	}
</script>

<svelte:head>
	<title>Dashboard - Resume Builder</title>
</svelte:head>

<div class="container mx-auto px-6 py-8">
	<!-- Header -->
	<div class="mb-8 flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold">My Resumes</h1>
			<p class="mt-1 text-muted-foreground">
				Create, edit, and manage your professional resumes
			</p>
		</div>
		<div class="flex gap-3">
			<Button variant="outline" on:click={() => (showImportModal = true)}>
				<Upload class="mr-2 h-4 w-4" />
				Import
			</Button>
			<Button on:click={openNewResumeModal}>
				<Plus class="mr-2 h-4 w-4" />
				New Resume
			</Button>
		</div>
	</div>

	<!-- Resume Grid -->
	{#if $isLoading}
		<div class="flex items-center justify-center py-20">
			<div class="text-center">
				<div
					class="mx-auto mb-4 h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"
				></div>
				<p class="text-muted-foreground">Loading your resumes...</p>
			</div>
		</div>
	{:else if $resumes.length === 0}
		<!-- Empty State -->
		<Card class="p-12 text-center">
			<FileText class="mx-auto mb-4 h-16 w-16 text-muted-foreground/50" />
			<h2 class="mb-2 text-xl font-semibold">No resumes yet</h2>
			<p class="mb-6 text-muted-foreground">
				Create your first resume to get started on your job search journey.
			</p>
			<div class="flex justify-center gap-3">
				<Button variant="outline" on:click={() => (showImportModal = true)}>
					<Upload class="mr-2 h-4 w-4" />
					Import Existing
				</Button>
				<Button on:click={openNewResumeModal}>
					<Plus class="mr-2 h-4 w-4" />
					Create New
				</Button>
			</div>
		</Card>
	{:else}
		<div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
			{#each $resumes as resume (resume.id)}
				<Card class="group relative transition-shadow hover:shadow-md">
					<a href="/resumes/{resume.id}/edit" class="block p-4">
						<div class="flex items-center gap-3">
							<div class="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
								<FileText class="h-5 w-5 text-primary" />
							</div>
							<div class="min-w-0 flex-1">
								<h3 class="truncate font-medium">{resume.title}</h3>
								{#if resume.data?.personal_info?.title}
									<p class="truncate text-sm text-muted-foreground">
										{resume.data.personal_info.title}
									</p>
								{/if}
								<p class="mt-0.5 flex items-center text-xs text-muted-foreground">
									<Clock class="mr-1 h-3 w-3" />
									{formatDate(resume.updated_at)}
								</p>
							</div>
						</div>
					</a>

					<!-- Actions Menu -->
					<div class="absolute right-2 top-2">
						<button
							class="rounded p-1.5 opacity-0 transition-opacity hover:bg-muted group-hover:opacity-100"
							on:click|stopPropagation={() =>
								(activeMenu = activeMenu === resume.id ? null : resume.id)}
						>
							<MoreVertical class="h-4 w-4" />
						</button>

						{#if activeMenu === resume.id}
							<div
								class="absolute right-0 top-full z-50 mt-1 w-40 rounded-md border bg-background py-1 shadow-lg"
							>
								<a
									href="/resumes/{resume.id}/edit"
									class="flex items-center px-3 py-2 text-sm hover:bg-muted"
								>
									<FileText class="mr-2 h-4 w-4" />
									Edit
								</a>
								<button
									class="flex w-full items-center px-3 py-2 text-sm hover:bg-muted"
									on:click={() => duplicateResume(resume.id, resume.title)}
								>
									<Copy class="mr-2 h-4 w-4" />
									Duplicate
								</button>
								{#if resume.is_public}
									<a
										href="/share/{resume.public_slug || resume.id}"
										target="_blank"
										class="flex items-center px-3 py-2 text-sm hover:bg-muted"
									>
										<ExternalLink class="mr-2 h-4 w-4" />
										View Public
									</a>
								{/if}
								<hr class="my-1" />
								<button
									class="flex w-full items-center px-3 py-2 text-sm text-destructive hover:bg-muted"
									on:click={() => deleteResume(resume.id)}
								>
									<Trash2 class="mr-2 h-4 w-4" />
									Delete
								</button>
							</div>
						{/if}
					</div>
				</Card>
			{/each}

		</div>
	{/if}
</div>

<!-- New Resume Modal -->
{#if showNewResumeModal}
	<!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
		on:click|self={() => (showNewResumeModal = false)}
		on:keydown={(e) => e.key === 'Escape' && (showNewResumeModal = false)}
		role="dialog"
		tabindex="-1"
	>
		<Card class="flex max-h-[90vh] w-full max-w-4xl flex-col">
			<div class="flex items-center justify-between border-b p-4">
				<h2 class="text-xl font-semibold">Create New Resume</h2>
				<button class="rounded p-1 hover:bg-muted" on:click={() => (showNewResumeModal = false)}>
					<X class="h-5 w-5" />
				</button>
			</div>

			<form on:submit|preventDefault={createNewResume} class="flex flex-1 flex-col overflow-hidden">
				<!-- Template selection -->
				<div class="border-b p-4">
					<p class="mb-3 text-sm font-medium">Choose a template (optional)</p>
					{#if isLoadingTemplates}
						<div class="flex justify-center py-8">
							<Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
						</div>
					{:else if templates.length > 0}
						<div class="grid max-h-56 grid-cols-3 gap-3 overflow-y-auto sm:grid-cols-4">
							<button
								type="button"
								class="relative overflow-hidden rounded-lg border-2 transition-all {selectedTemplateId ===
								null
									? 'border-primary ring-2 ring-primary ring-offset-2'
									: 'border-transparent hover:border-muted-foreground/20'}"
								on:click={() => (selectedTemplateId = null)}
							>
								<div class="flex aspect-[8.5/11] items-center justify-center bg-muted">
									<Layout class="h-8 w-8 text-muted-foreground" />
								</div>
								<span class="block truncate p-2 text-center text-xs font-medium">Default</span>
							</button>
							{#each templates as template (template.id)}
								<button
									type="button"
									class="relative overflow-hidden rounded-lg border-2 transition-all {selectedTemplateId ===
									template.id
										? 'border-primary ring-2 ring-primary ring-offset-2'
										: 'border-transparent hover:border-muted-foreground/20'}"
									on:click={() => (selectedTemplateId = template.id)}
								>
									<div class="aspect-[8.5/11] overflow-hidden rounded-t-lg bg-white">
										<img
											src={getTemplatePreviewUrl(template)}
											alt={template.name}
											class="h-full w-full object-cover object-top"
										/>
									</div>
									<span class="block truncate p-2 text-center text-xs font-medium">{template.name}</span>
									{#if selectedTemplateId === template.id}
										<div class="absolute right-1 top-1 rounded-full bg-primary p-0.5">
											<Check class="h-3 w-3 text-primary-foreground" />
										</div>
									{/if}
								</button>
							{/each}
						</div>
					{/if}
				</div>

				<!-- Title and actions -->
				<div class="flex flex-1 flex-col justify-between gap-4 p-4">
					<div>
						<label for="title" class="mb-2 block text-sm font-medium">Resume Title</label>
						<input
							id="title"
							type="text"
							bind:value={newResumeTitle}
							class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
							placeholder="e.g., Software Engineer Resume"
						/>
					</div>
					<div class="flex justify-end gap-3">
						<Button variant="outline" type="button" on:click={() => (showNewResumeModal = false)}>
							Cancel
						</Button>
						<Button type="submit">Create Resume</Button>
					</div>
				</div>
			</form>
		</Card>
	</div>
{/if}

<!-- Import Resume Modal -->
{#if showImportModal}
	<!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
		on:click|self={() => (showImportModal = false)}
		on:keydown={(e) => e.key === 'Escape' && (showImportModal = false)}
		role="dialog"
		tabindex="-1"
	>
		<Card class="w-full max-w-2xl p-6">
			<h2 class="mb-2 text-xl font-semibold">Import Resume</h2>
			<p class="mb-4 text-sm text-muted-foreground">
				Upload a file or paste text and our AI will extract the information automatically.
			</p>
			
			<!-- Import Mode Tabs -->
			<div class="mb-4 flex border-b">
				<button
					type="button"
					class="px-4 py-2 text-sm font-medium transition-colors {importMode === 'file' 
						? 'border-b-2 border-primary text-primary' 
						: 'text-muted-foreground hover:text-foreground'}"
					on:click={() => (importMode = 'file')}
				>
					<FileIcon class="mr-2 inline-block h-4 w-4" />
					Upload File
				</button>
				<button
					type="button"
					class="px-4 py-2 text-sm font-medium transition-colors {importMode === 'text' 
						? 'border-b-2 border-primary text-primary' 
						: 'text-muted-foreground hover:text-foreground'}"
					on:click={() => (importMode = 'text')}
				>
					<FileText class="mr-2 inline-block h-4 w-4" />
					Paste Text
				</button>
			</div>

			<form on:submit|preventDefault={importResume}>
				<div class="mb-4">
					<label for="importTitle" class="mb-2 block text-sm font-medium">Resume Title</label>
					<input
						id="importTitle"
						type="text"
						bind:value={importTitle}
						class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
						placeholder="e.g., My Imported Resume"
					/>
				</div>

				{#if importMode === 'file'}
					<!-- File Upload Area -->
					<div class="mb-4">
						<span class="mb-2 block text-sm font-medium">Resume File</span>
						<!-- svelte-ignore a11y-no-static-element-interactions -->
						<div
							class="relative rounded-lg border-2 border-dashed transition-colors {isDragging 
								? 'border-primary bg-primary/5' 
								: 'border-muted-foreground/25 hover:border-muted-foreground/50'}"
							on:drop={handleDrop}
							on:dragover={handleDragOver}
							on:dragleave={handleDragLeave}
						>
							{#if selectedFile}
								<!-- Selected File Display -->
								<div class="flex items-center gap-3 p-4">
									<div class="flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
										<FileIcon class="h-6 w-6 text-primary" />
									</div>
									<div class="min-w-0 flex-1">
										<p class="truncate font-medium">{selectedFile.name}</p>
										<p class="text-sm text-muted-foreground">{formatFileSize(selectedFile.size)}</p>
									</div>
									<button
										type="button"
										class="rounded-full p-1 hover:bg-muted"
										on:click={removeFile}
										disabled={isImporting}
									>
										<X class="h-5 w-5" />
									</button>
								</div>
							{:else}
								<!-- Drop Zone -->
								<div class="flex flex-col items-center justify-center p-8 text-center">
									<Upload class="mb-3 h-10 w-10 text-muted-foreground" />
									<p class="mb-1 text-sm font-medium">
										Drag and drop your resume here
									</p>
									<p class="mb-3 text-xs text-muted-foreground">
										or click to browse
									</p>
									<p class="text-xs text-muted-foreground">
										Supported formats: PDF, DOCX, DOC, TXT (max 10MB)
									</p>
								</div>
								<input
									bind:this={fileInputRef}
									type="file"
									accept={acceptedFileTypes}
									class="absolute inset-0 cursor-pointer opacity-0"
									on:change={handleFileSelect}
									disabled={isImporting}
								/>
							{/if}
						</div>
					</div>
				{:else}
					<!-- Text Paste Area -->
					<div class="mb-4">
						<label for="importText" class="mb-2 block text-sm font-medium">Resume Content</label>
						<textarea
							id="importText"
							bind:value={importText}
							class="flex min-h-[250px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
							placeholder="Paste your resume text here...

Example:
John Doe
Software Engineer
john@example.com | (555) 123-4567

EXPERIENCE
Senior Developer at Tech Corp (2020-Present)
- Led development of microservices architecture
- Mentored team of 5 junior developers

EDUCATION
BS Computer Science, State University (2016)"
							disabled={isImporting}
						></textarea>
					</div>
				{/if}

				{#if importError}
					<div class="mb-4 rounded-md bg-destructive/10 p-3 text-sm text-destructive">
						{importError}
					</div>
				{/if}
				<div class="flex justify-end gap-3">
					<Button
						variant="outline"
						type="button"
						on:click={() => {
							showImportModal = false;
							importError = '';
							importText = '';
							selectedFile = null;
						}}
						disabled={isImporting}
					>
						Cancel
					</Button>
					<Button type="submit" disabled={isImporting}>
						{#if isImporting}
							<Loader2 class="mr-2 h-4 w-4 animate-spin" />
							Extracting...
						{:else}
							<Upload class="mr-2 h-4 w-4" />
							Import Resume
						{/if}
					</Button>
				</div>
			</form>
		</Card>
	</div>
{/if}
