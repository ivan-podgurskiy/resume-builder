<script lang="ts">
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
		Loader2
	} from 'lucide-svelte';
	import { formatDate } from '$utils';
	import { api } from '$lib/api/client';

	let showNewResumeModal = false;
	let showImportModal = false;
	let newResumeTitle = 'My Resume';
	let importText = '';
	let importTitle = 'Imported Resume';
	let isImporting = false;
	let importError = '';
	let activeMenu: string | null = null;

	onMount(() => {
		resumeStore.loadResumes();
	});

	async function createNewResume() {
		try {
			const resume = await resumeStore.createResume({ title: newResumeTitle });
			goto(`/resumes/${resume.id}/edit`);
		} catch (error) {
			console.error('Failed to create resume:', error);
		}
	}

	async function importResume() {
		if (!importText.trim()) {
			importError = 'Please paste your resume text';
			return;
		}

		isImporting = true;
		importError = '';

		try {
			// Extract data using AI
			const extractedData = await api.extractResumeData(importText);

			// Create a new resume with the extracted data
			const resume = await resumeStore.createResume({
				title: importTitle,
				data: extractedData
			});

			showImportModal = false;
			importText = '';
			importTitle = 'Imported Resume';
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
			<Button on:click={() => (showNewResumeModal = true)}>
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
				<Button on:click={() => (showNewResumeModal = true)}>
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
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
		on:click|self={() => (showNewResumeModal = false)}
	>
		<Card class="w-full max-w-md p-6">
			<h2 class="mb-4 text-xl font-semibold">Create New Resume</h2>
			<form on:submit|preventDefault={createNewResume}>
				<div class="mb-4">
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
			</form>
		</Card>
	</div>
{/if}

<!-- Import Resume Modal -->
{#if showImportModal}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
		on:click|self={() => (showImportModal = false)}
	>
		<Card class="w-full max-w-2xl p-6">
			<h2 class="mb-2 text-xl font-semibold">Import Resume</h2>
			<p class="mb-4 text-sm text-muted-foreground">
				Paste your resume text below and our AI will extract the information automatically.
			</p>
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
				<div class="mb-4">
					<label for="importText" class="mb-2 block text-sm font-medium">Resume Content</label>
					<textarea
						id="importText"
						bind:value={importText}
						class="flex min-h-[300px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
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
