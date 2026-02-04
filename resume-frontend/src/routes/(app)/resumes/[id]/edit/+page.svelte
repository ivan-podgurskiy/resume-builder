<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Button, Card, Input, Label, Textarea } from '$components/ui';
	import { resumeStore, currentResume, isSaving } from '$stores/resume';
	import {
		ArrowLeft,
		Save,
		Download,
		Eye,
		Settings,
		Plus,
		Trash2,
		GripVertical,
		Sparkles,
		Check,
		FileText,
		FileJson,
		Loader2,
		X
	} from 'lucide-svelte';
	import { generateId } from '$utils';
	import type { Experience, Education, ResumeData } from '$types';
	import { api } from '$lib/api/client';

	let activeSection = 'personal';
	let saveStatus = 'saved';
	let showExportModal = false;
	let showPreviewModal = false;
	let isExporting = false;
	let exportError = '';

	$: resumeId = $page.params.id;

	onMount(async () => {
		try {
			await resumeStore.loadResume(resumeId);
		} catch (error) {
			console.error('Failed to load resume:', error);
			goto('/dashboard');
		}
	});

	// Auto-save when data changes
	function handleChange() {
		if ($currentResume) {
			saveStatus = 'saving';
			resumeStore.scheduleAutoSave(resumeId, { data: $currentResume.data });
		}
	}

	$: if (!$isSaving && saveStatus === 'saving') {
		saveStatus = 'saved';
	}

	function updatePersonalInfo(field: string, value: string) {
		if ($currentResume) {
			resumeStore.updateLocalResume({
				personal_info: {
					...$currentResume.data.personal_info,
					[field]: value
				}
			});
			handleChange();
		}
	}

	function updateSummary(value: string) {
		resumeStore.updateLocalResume({ summary: value });
		handleChange();
	}

	function addExperience() {
		if ($currentResume) {
			const currentExperience = $currentResume.data.experience || [];
			const newExp: Experience = {
				id: generateId(),
				company: '',
				position: '',
				location: '',
				start_date: '',
				end_date: '',
				is_current: false,
				description: '',
				achievements: [],
				technologies: [],
				order: currentExperience.length
			};
			resumeStore.updateLocalResume({
				experience: [...currentExperience, newExp]
			});
			handleChange();
		}
	}

	function updateExperience(index: number, field: string, value: string | boolean) {
		if ($currentResume) {
			const updated = [...($currentResume.data.experience || [])];
			updated[index] = { ...updated[index], [field]: value };
			resumeStore.updateLocalResume({ experience: updated });
			handleChange();
		}
	}

	function removeExperience(index: number) {
		if ($currentResume) {
			const updated = ($currentResume.data.experience || []).filter((_, i) => i !== index);
			resumeStore.updateLocalResume({ experience: updated });
			handleChange();
		}
	}

	function addEducation() {
		if ($currentResume) {
			const currentEducation = $currentResume.data.education || [];
			const newEdu: Education = {
				id: generateId(),
				institution: '',
				degree: '',
				field_of_study: '',
				location: '',
				start_date: '',
				end_date: '',
				is_current: false,
				gpa: '',
				honors: '',
				description: '',
				order: currentEducation.length
			};
			resumeStore.updateLocalResume({
				education: [...currentEducation, newEdu]
			});
			handleChange();
		}
	}

	function updateEducation(index: number, field: string, value: string | boolean) {
		if ($currentResume) {
			const updated = [...($currentResume.data.education || [])];
			updated[index] = { ...updated[index], [field]: value };
			resumeStore.updateLocalResume({ education: updated });
			handleChange();
		}
	}

	function removeEducation(index: number) {
		if ($currentResume) {
			const updated = ($currentResume.data.education || []).filter((_, i) => i !== index);
			resumeStore.updateLocalResume({ education: updated });
			handleChange();
		}
	}

	function updateSkills(type: 'technical' | 'soft', value: string) {
		if ($currentResume) {
			const skills = value.split(',').map((s) => s.trim()).filter(Boolean);
			if (type === 'technical') {
				resumeStore.updateLocalResume({
					skills: {
						...$currentResume.data.skills,
						technical: skills.map((name) => ({ name }))
					}
				});
			} else {
				resumeStore.updateLocalResume({
					skills: {
						...$currentResume.data.skills,
						soft: skills
					}
				});
			}
			handleChange();
		}
	}

	const sections = [
		{ id: 'personal', label: 'Personal Info' },
		{ id: 'summary', label: 'Summary' },
		{ id: 'experience', label: 'Experience' },
		{ id: 'education', label: 'Education' },
		{ id: 'skills', label: 'Skills' }
	];

	function downloadFile(content: Blob | string, filename: string, mimeType: string) {
		const blob = content instanceof Blob ? content : new Blob([content], { type: mimeType });
		const url = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = url;
		a.download = filename;
		document.body.appendChild(a);
		a.click();
		document.body.removeChild(a);
		URL.revokeObjectURL(url);
	}

	async function exportPDF() {
		if (!$currentResume) return;
		isExporting = true;
		exportError = '';
		try {
			const blob = await api.exportPDF(resumeId);
			downloadFile(blob, `${$currentResume.title}.pdf`, 'application/pdf');
			showExportModal = false;
		} catch (error) {
			console.error('Failed to export PDF:', error);
			exportError = error instanceof Error ? error.message : 'Failed to export PDF';
		} finally {
			isExporting = false;
		}
	}

	async function exportTXT() {
		if (!$currentResume) return;
		isExporting = true;
		exportError = '';
		try {
			const text = await api.exportTXT(resumeId);
			downloadFile(text, `${$currentResume.title}.txt`, 'text/plain');
			showExportModal = false;
		} catch (error) {
			console.error('Failed to export TXT:', error);
			exportError = error instanceof Error ? error.message : 'Failed to export TXT';
		} finally {
			isExporting = false;
		}
	}

	async function exportJSON() {
		if (!$currentResume) return;
		isExporting = true;
		exportError = '';
		try {
			const json = await api.exportJSON(resumeId);
			downloadFile(json, `${$currentResume.title}.json`, 'application/json');
			showExportModal = false;
		} catch (error) {
			console.error('Failed to export JSON:', error);
			exportError = error instanceof Error ? error.message : 'Failed to export JSON';
		} finally {
			isExporting = false;
		}
	}

	function openPreview() {
		showPreviewModal = true;
	}
</script>

<svelte:head>
	<title>{$currentResume?.title || 'Edit Resume'} - Resume Builder</title>
</svelte:head>

{#if $currentResume}
	<div class="flex h-[calc(100vh-4rem)] flex-col">
		<!-- Editor Header -->
		<div class="flex items-center justify-between border-b bg-background px-4 py-3">
			<div class="flex items-center gap-4">
				<a href="/dashboard" class="rounded p-1 hover:bg-muted">
					<ArrowLeft class="h-5 w-5" />
				</a>
				<input
					type="text"
					value={$currentResume.title}
					class="border-0 bg-transparent text-lg font-semibold focus:outline-none focus:ring-0"
					on:change={(e) => {
						resumeStore.updateResume(resumeId, { title: e.currentTarget.value });
					}}
				/>
				<span class="flex items-center text-sm text-muted-foreground">
					{#if $isSaving}
						<span class="mr-1 h-2 w-2 animate-pulse rounded-full bg-yellow-500"></span>
						Saving...
					{:else}
						<Check class="mr-1 h-4 w-4 text-green-500" />
						Saved
					{/if}
				</span>
			</div>
			<div class="flex items-center gap-2">
				<Button variant="outline" size="sm" on:click={openPreview}>
					<Eye class="mr-2 h-4 w-4" />
					Preview
				</Button>
				<Button variant="outline" size="sm" on:click={() => (showExportModal = true)}>
					<Download class="mr-2 h-4 w-4" />
					Export
				</Button>
			</div>
		</div>

		<!-- Editor Content -->
		<div class="flex flex-1 overflow-hidden">
			<!-- Left Panel: Editor -->
			<div class="w-[400px] overflow-y-auto border-r bg-muted/30">
				<!-- Section Tabs -->
				<div class="flex border-b bg-background">
					{#each sections as section}
						<button
							class="flex-1 border-b-2 px-3 py-3 text-sm font-medium transition-colors {activeSection ===
							section.id
								? 'border-primary text-primary'
								: 'border-transparent text-muted-foreground hover:text-foreground'}"
							on:click={() => (activeSection = section.id)}
						>
							{section.label}
						</button>
					{/each}
				</div>

				<!-- Section Content -->
				<div class="p-4">
					{#if activeSection === 'personal'}
						<div class="space-y-4">
							<div class="grid grid-cols-2 gap-4">
								<div>
									<Label htmlFor="firstName">First Name</Label>
									<Input
										id="firstName"
										value={$currentResume.data.personal_info.first_name}
										on:input={(e) => updatePersonalInfo('first_name', e.currentTarget.value)}
										placeholder="John"
									/>
								</div>
								<div>
									<Label htmlFor="lastName">Last Name</Label>
									<Input
										id="lastName"
										value={$currentResume.data.personal_info.last_name}
										on:input={(e) => updatePersonalInfo('last_name', e.currentTarget.value)}
										placeholder="Doe"
									/>
								</div>
							</div>
							<div>
								<Label htmlFor="title">Professional Title</Label>
								<Input
									id="title"
									value={$currentResume.data.personal_info.title}
									on:input={(e) => updatePersonalInfo('title', e.currentTarget.value)}
									placeholder="Senior Software Engineer"
								/>
							</div>
							<div>
								<Label htmlFor="email">Email</Label>
								<Input
									id="email"
									type="email"
									value={$currentResume.data.personal_info.email}
									on:input={(e) => updatePersonalInfo('email', e.currentTarget.value)}
									placeholder="john@example.com"
								/>
							</div>
							<div>
								<Label htmlFor="phone">Phone</Label>
								<Input
									id="phone"
									value={$currentResume.data.personal_info.phone}
									on:input={(e) => updatePersonalInfo('phone', e.currentTarget.value)}
									placeholder="+1 (555) 123-4567"
								/>
							</div>
							<div>
								<Label htmlFor="location">Location</Label>
								<Input
									id="location"
									value={$currentResume.data.personal_info.location}
									on:input={(e) => updatePersonalInfo('location', e.currentTarget.value)}
									placeholder="San Francisco, CA"
								/>
							</div>
							<div>
								<Label htmlFor="linkedin">LinkedIn</Label>
								<Input
									id="linkedin"
									value={$currentResume.data.personal_info.linkedin || ''}
									on:input={(e) => updatePersonalInfo('linkedin', e.currentTarget.value)}
									placeholder="linkedin.com/in/johndoe"
								/>
							</div>
							<div>
								<Label htmlFor="website">Website</Label>
								<Input
									id="website"
									value={$currentResume.data.personal_info.website || ''}
									on:input={(e) => updatePersonalInfo('website', e.currentTarget.value)}
									placeholder="johndoe.com"
								/>
							</div>
						</div>
					{:else if activeSection === 'summary'}
						<div class="space-y-4">
							<div>
								<div class="mb-2 flex items-center justify-between">
									<Label htmlFor="summary">Professional Summary</Label>
									<Button variant="ghost" size="sm">
										<Sparkles class="mr-1 h-4 w-4" />
										Improve with AI
									</Button>
								</div>
								<Textarea
									id="summary"
									value={$currentResume.data.summary}
									on:input={(e) => updateSummary(e.currentTarget.value)}
									placeholder="Write a brief summary of your professional background and career goals..."
									rows={6}
								/>
							</div>
						</div>
					{:else if activeSection === 'experience'}
						<div class="space-y-4">
							{#each ($currentResume.data.experience || []) as exp, index (exp.id)}
								<Card class="p-4">
									<div class="mb-4 flex items-center justify-between">
										<div class="flex items-center gap-2">
											<GripVertical class="h-4 w-4 cursor-move text-muted-foreground" />
											<span class="font-medium">{exp.position || 'New Position'}</span>
										</div>
										<button
											class="rounded p-1 text-muted-foreground hover:bg-destructive/10 hover:text-destructive"
											on:click={() => removeExperience(index)}
										>
											<Trash2 class="h-4 w-4" />
										</button>
									</div>
									<div class="space-y-3">
										<div>
											<Label>Position</Label>
											<Input
												value={exp.position}
												on:input={(e) => updateExperience(index, 'position', e.currentTarget.value)}
												placeholder="Software Engineer"
											/>
										</div>
										<div>
											<Label>Company</Label>
											<Input
												value={exp.company}
												on:input={(e) => updateExperience(index, 'company', e.currentTarget.value)}
												placeholder="Acme Inc."
											/>
										</div>
										<div class="grid grid-cols-2 gap-3">
											<div>
												<Label>Start Date</Label>
												<Input
													type="month"
													value={exp.start_date}
													on:input={(e) =>
														updateExperience(index, 'start_date', e.currentTarget.value)}
												/>
											</div>
											<div>
												<Label>End Date</Label>
												<Input
													type="month"
													value={exp.end_date}
													on:input={(e) =>
														updateExperience(index, 'end_date', e.currentTarget.value)}
													disabled={exp.is_current}
												/>
											</div>
										</div>
										<label class="flex items-center gap-2 text-sm">
											<input
												type="checkbox"
												checked={exp.is_current}
												on:change={(e) =>
													updateExperience(index, 'is_current', e.currentTarget.checked)}
												class="rounded border-input"
											/>
											Currently working here
										</label>
										<div>
											<div class="mb-1 flex items-center justify-between">
												<Label>Description</Label>
												<Button variant="ghost" size="sm" class="h-6 px-2 text-xs">
													<Sparkles class="mr-1 h-3 w-3" />
													Improve
												</Button>
											</div>
											<Textarea
												value={exp.description}
												on:input={(e) =>
													updateExperience(index, 'description', e.currentTarget.value)}
												placeholder="Describe your responsibilities and achievements..."
												rows={4}
											/>
										</div>
									</div>
								</Card>
							{/each}
							<Button variant="outline" class="w-full" on:click={addExperience}>
								<Plus class="mr-2 h-4 w-4" />
								Add Experience
							</Button>
						</div>
					{:else if activeSection === 'education'}
						<div class="space-y-4">
							{#each ($currentResume.data.education || []) as edu, index (edu.id)}
								<Card class="p-4">
									<div class="mb-4 flex items-center justify-between">
										<div class="flex items-center gap-2">
											<GripVertical class="h-4 w-4 cursor-move text-muted-foreground" />
											<span class="font-medium">{edu.degree || 'New Education'}</span>
										</div>
										<button
											class="rounded p-1 text-muted-foreground hover:bg-destructive/10 hover:text-destructive"
											on:click={() => removeEducation(index)}
										>
											<Trash2 class="h-4 w-4" />
										</button>
									</div>
									<div class="space-y-3">
										<div>
											<Label>Institution</Label>
											<Input
												value={edu.institution}
												on:input={(e) =>
													updateEducation(index, 'institution', e.currentTarget.value)}
												placeholder="University of California"
											/>
										</div>
										<div>
											<Label>Degree</Label>
											<Input
												value={edu.degree}
												on:input={(e) => updateEducation(index, 'degree', e.currentTarget.value)}
												placeholder="Bachelor of Science"
											/>
										</div>
										<div>
											<Label>Field of Study</Label>
											<Input
												value={edu.field_of_study}
												on:input={(e) =>
													updateEducation(index, 'field_of_study', e.currentTarget.value)}
												placeholder="Computer Science"
											/>
										</div>
										<div class="grid grid-cols-2 gap-3">
											<div>
												<Label>Start Date</Label>
												<Input
													type="month"
													value={edu.start_date}
													on:input={(e) =>
														updateEducation(index, 'start_date', e.currentTarget.value)}
												/>
											</div>
											<div>
												<Label>End Date</Label>
												<Input
													type="month"
													value={edu.end_date}
													on:input={(e) =>
														updateEducation(index, 'end_date', e.currentTarget.value)}
												/>
											</div>
										</div>
									</div>
								</Card>
							{/each}
							<Button variant="outline" class="w-full" on:click={addEducation}>
								<Plus class="mr-2 h-4 w-4" />
								Add Education
							</Button>
						</div>
					{:else if activeSection === 'skills'}
						<div class="space-y-4">
							<div>
								<Label htmlFor="technicalSkills">Technical Skills</Label>
								<Textarea
									id="technicalSkills"
									value={($currentResume.data.skills.technical || []).map((s) => s.name).join(', ')}
									on:input={(e) => updateSkills('technical', e.currentTarget.value)}
									placeholder="JavaScript, TypeScript, React, Node.js..."
									rows={3}
								/>
								<p class="mt-1 text-xs text-muted-foreground">Separate skills with commas</p>
							</div>
							<div>
								<Label htmlFor="softSkills">Soft Skills</Label>
								<Textarea
									id="softSkills"
									value={$currentResume.data.skills.soft?.join(', ') || ''}
									on:input={(e) => updateSkills('soft', e.currentTarget.value)}
									placeholder="Leadership, Communication, Problem-solving..."
									rows={3}
								/>
								<p class="mt-1 text-xs text-muted-foreground">Separate skills with commas</p>
							</div>
						</div>
					{/if}
				</div>
			</div>

			<!-- Right Panel: Live Preview -->
			<div class="flex-1 overflow-y-auto bg-muted/50 p-8">
				<div class="mx-auto max-w-[850px]">
					<!-- Preview Controls -->
					<div class="mb-4 flex items-center justify-between">
						<span class="text-sm text-muted-foreground">Live Preview</span>
						<div class="flex items-center gap-2">
							<Button variant="outline" size="sm">
								<Settings class="mr-2 h-4 w-4" />
								Change Template
							</Button>
						</div>
					</div>

					<!-- Resume Preview -->
					<div
						class="aspect-[8.5/11] rounded-lg border bg-white p-8 shadow-lg"
						style="font-family: Inter, sans-serif;"
					>
						<!-- Header -->
						<div class="mb-6 border-b pb-4 text-center">
							<h1 class="text-2xl font-bold text-gray-900">
								{$currentResume.data.personal_info.first_name}
								{$currentResume.data.personal_info.last_name}
							</h1>
							{#if $currentResume.data.personal_info.title}
								<p class="mt-1 text-lg text-gray-600">
									{$currentResume.data.personal_info.title}
								</p>
							{/if}
							<div class="mt-2 flex flex-wrap justify-center gap-3 text-sm text-gray-500">
								{#if $currentResume.data.personal_info.email}
									<span>{$currentResume.data.personal_info.email}</span>
								{/if}
								{#if $currentResume.data.personal_info.phone}
									<span>{$currentResume.data.personal_info.phone}</span>
								{/if}
								{#if $currentResume.data.personal_info.location}
									<span>{$currentResume.data.personal_info.location}</span>
								{/if}
							</div>
						</div>

						<!-- Summary -->
						{#if $currentResume.data.summary}
							<div class="mb-6">
								<h2 class="mb-2 text-sm font-bold uppercase tracking-wide text-gray-900">
									Professional Summary
								</h2>
								<p class="text-sm text-gray-700">{$currentResume.data.summary}</p>
							</div>
						{/if}

						<!-- Experience -->
						{#if ($currentResume.data.experience || []).length > 0}
							<div class="mb-6">
								<h2 class="mb-3 text-sm font-bold uppercase tracking-wide text-gray-900">
									Experience
								</h2>
								{#each ($currentResume.data.experience || []) as exp}
									<div class="mb-3">
										<div class="flex items-start justify-between">
											<div>
												<h3 class="font-semibold text-gray-900">{exp.position}</h3>
												<p class="text-sm text-gray-600">{exp.company}</p>
											</div>
											<span class="text-sm text-gray-500">
												{exp.start_date} - {exp.is_current ? 'Present' : exp.end_date}
											</span>
										</div>
										{#if exp.description}
											<p class="mt-1 text-sm text-gray-700">{exp.description}</p>
										{/if}
									</div>
								{/each}
							</div>
						{/if}

						<!-- Education -->
						{#if ($currentResume.data.education || []).length > 0}
							<div class="mb-6">
								<h2 class="mb-3 text-sm font-bold uppercase tracking-wide text-gray-900">
									Education
								</h2>
								{#each ($currentResume.data.education || []) as edu}
									<div class="mb-2">
										<div class="flex items-start justify-between">
											<div>
												<h3 class="font-semibold text-gray-900">
													{edu.degree} in {edu.field_of_study}
												</h3>
												<p class="text-sm text-gray-600">{edu.institution}</p>
											</div>
											<span class="text-sm text-gray-500">
												{edu.start_date} - {edu.end_date}
											</span>
										</div>
									</div>
								{/each}
							</div>
						{/if}

						<!-- Skills -->
						{#if ($currentResume.data.skills.technical || []).length > 0}
							<div>
								<h2 class="mb-2 text-sm font-bold uppercase tracking-wide text-gray-900">
									Skills
								</h2>
								<div class="flex flex-wrap gap-2">
									{#each ($currentResume.data.skills.technical || []) as skill}
										<span class="rounded bg-gray-100 px-2 py-1 text-sm text-gray-700">
											{skill.name}
										</span>
									{/each}
								</div>
							</div>
						{/if}
					</div>
				</div>
			</div>
		</div>
	</div>
{:else}
	<div class="flex h-[calc(100vh-4rem)] items-center justify-center">
		<div class="text-center">
			<div
				class="mx-auto mb-4 h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"
			></div>
			<p class="text-muted-foreground">Loading resume...</p>
		</div>
	</div>
{/if}

<!-- Export Modal -->
{#if showExportModal}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
		on:click|self={() => (showExportModal = false)}
		on:keydown={(e) => e.key === 'Escape' && (showExportModal = false)}
		role="dialog"
		tabindex="-1"
	>
		<Card class="w-full max-w-md p-6">
			<div class="mb-4 flex items-center justify-between">
				<h2 class="text-xl font-semibold">Export Resume</h2>
				<button
					class="rounded p-1 hover:bg-muted"
					on:click={() => (showExportModal = false)}
				>
					<X class="h-5 w-5" />
				</button>
			</div>
			<p class="mb-6 text-sm text-muted-foreground">
				Choose a format to download your resume.
			</p>
			{#if exportError}
				<div class="mb-4 rounded-md bg-destructive/10 p-3 text-sm text-destructive">
					{exportError}
				</div>
			{/if}
			<div class="space-y-3">
				<button
					class="flex w-full items-center gap-4 rounded-lg border p-4 text-left transition-colors hover:bg-muted disabled:opacity-50"
					on:click={exportPDF}
					disabled={isExporting}
				>
					<div class="flex h-10 w-10 items-center justify-center rounded-lg bg-red-100 text-red-600">
						{#if isExporting}
							<Loader2 class="h-5 w-5 animate-spin" />
						{:else}
							<FileText class="h-5 w-5" />
						{/if}
					</div>
					<div>
						<div class="font-medium">PDF Document</div>
						<div class="text-sm text-muted-foreground">Best for printing and sharing</div>
					</div>
				</button>
				<button
					class="flex w-full items-center gap-4 rounded-lg border p-4 text-left transition-colors hover:bg-muted disabled:opacity-50"
					on:click={exportTXT}
					disabled={isExporting}
				>
					<div class="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-100 text-blue-600">
						<FileText class="h-5 w-5" />
					</div>
					<div>
						<div class="font-medium">Plain Text</div>
						<div class="text-sm text-muted-foreground">ATS-friendly format</div>
					</div>
				</button>
				<button
					class="flex w-full items-center gap-4 rounded-lg border p-4 text-left transition-colors hover:bg-muted disabled:opacity-50"
					on:click={exportJSON}
					disabled={isExporting}
				>
					<div class="flex h-10 w-10 items-center justify-center rounded-lg bg-green-100 text-green-600">
						<FileJson class="h-5 w-5" />
					</div>
					<div>
						<div class="font-medium">JSON Data</div>
						<div class="text-sm text-muted-foreground">For backup or import elsewhere</div>
					</div>
				</button>
			</div>
		</Card>
	</div>
{/if}

<!-- Preview Modal -->
{#if showPreviewModal && $currentResume}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/80 p-8"
		on:click|self={() => (showPreviewModal = false)}
		on:keydown={(e) => e.key === 'Escape' && (showPreviewModal = false)}
		role="dialog"
		tabindex="-1"
	>
		<!-- Close button fixed in top-right corner -->
		<button
			class="absolute right-4 top-4 z-50 rounded-full bg-white p-2 shadow-lg hover:bg-gray-100"
			on:click={() => (showPreviewModal = false)}
		>
			<X class="h-5 w-5" />
		</button>
		<div class="relative max-h-full w-full max-w-4xl overflow-auto">
			<div
				class="mx-auto aspect-[8.5/11] max-w-[850px] rounded-lg bg-white p-12 shadow-2xl"
				style="font-family: Inter, sans-serif;"
			>
				<!-- Header -->
				<div class="mb-8 border-b-2 border-gray-200 pb-6 text-center">
					<h1 class="text-3xl font-bold text-gray-900">
						{$currentResume.data.personal_info.first_name}
						{$currentResume.data.personal_info.last_name}
					</h1>
					{#if $currentResume.data.personal_info.title}
						<p class="mt-2 text-xl text-gray-600">
							{$currentResume.data.personal_info.title}
						</p>
					{/if}
					<div class="mt-3 flex flex-wrap justify-center gap-4 text-sm text-gray-500">
						{#if $currentResume.data.personal_info.email}
							<span>{$currentResume.data.personal_info.email}</span>
						{/if}
						{#if $currentResume.data.personal_info.phone}
							<span>{$currentResume.data.personal_info.phone}</span>
						{/if}
						{#if $currentResume.data.personal_info.location}
							<span>{$currentResume.data.personal_info.location}</span>
						{/if}
					</div>
				</div>

				<!-- Summary -->
				{#if $currentResume.data.summary}
					<div class="mb-8">
						<h2 class="mb-3 text-lg font-bold uppercase tracking-wide text-gray-900">
							Professional Summary
						</h2>
						<p class="text-gray-700 leading-relaxed">{$currentResume.data.summary}</p>
					</div>
				{/if}

				<!-- Experience -->
				{#if ($currentResume.data.experience || []).length > 0}
					<div class="mb-8">
						<h2 class="mb-4 text-lg font-bold uppercase tracking-wide text-gray-900">
							Experience
						</h2>
						{#each ($currentResume.data.experience || []) as exp}
							<div class="mb-4">
								<div class="flex items-start justify-between">
									<div>
										<h3 class="text-lg font-semibold text-gray-900">{exp.position}</h3>
										<p class="text-gray-600">{exp.company}</p>
									</div>
									<span class="text-sm text-gray-500">
										{exp.start_date} - {exp.is_current ? 'Present' : exp.end_date}
									</span>
								</div>
								{#if exp.description}
									<p class="mt-2 text-gray-700">{exp.description}</p>
								{/if}
							</div>
						{/each}
					</div>
				{/if}

				<!-- Education -->
				{#if ($currentResume.data.education || []).length > 0}
					<div class="mb-8">
						<h2 class="mb-4 text-lg font-bold uppercase tracking-wide text-gray-900">
							Education
						</h2>
						{#each ($currentResume.data.education || []) as edu}
							<div class="mb-3">
								<div class="flex items-start justify-between">
									<div>
										<h3 class="font-semibold text-gray-900">
											{edu.degree} in {edu.field_of_study}
										</h3>
										<p class="text-gray-600">{edu.institution}</p>
									</div>
									<span class="text-sm text-gray-500">
										{edu.start_date} - {edu.end_date}
									</span>
								</div>
							</div>
						{/each}
					</div>
				{/if}

				<!-- Skills -->
				{#if ($currentResume.data.skills.technical || []).length > 0}
					<div>
						<h2 class="mb-3 text-lg font-bold uppercase tracking-wide text-gray-900">
							Skills
						</h2>
						<div class="flex flex-wrap gap-2">
							{#each ($currentResume.data.skills.technical || []) as skill}
								<span class="rounded-full bg-gray-100 px-3 py-1 text-sm text-gray-700">
									{skill.name}
								</span>
							{/each}
						</div>
					</div>
				{/if}
			</div>
		</div>
	</div>
{/if}
