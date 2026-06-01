<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { api } from '$lib/api/client';
	import { Button } from '$components/ui';
	import ResumeDocument from '$components/ResumeDocument.svelte';
	import { FileText, Loader2, Printer, ArrowRight, FileQuestion } from 'lucide-svelte';
	import type { Resume, Template } from '$types';

	let resume: Resume | null = null;
	let template: Template | null = null;
	let isLoading = true;
	let notFound = false;

	$: slug = $page.params.slug ?? '';

	$: fullName = resume
		? `${resume.data.personal_info.first_name ?? ''} ${resume.data.personal_info.last_name ?? ''}`.trim()
		: '';
	$: pageTitle = fullName ? `${fullName} — Resume` : 'Shared Resume';

	onMount(async () => {
		try {
			const result = await api.getPublicResume(slug);
			resume = result;
			template = result.template ?? null;
		} catch {
			notFound = true;
		} finally {
			isLoading = false;
		}
	});

	function print() {
		window.print();
	}
</script>

<svelte:head>
	<title>{pageTitle}</title>
	<meta name="robots" content="noindex" />
	{#if resume}
		<meta property="og:title" content={pageTitle} />
		{#if resume.data.personal_info.title}
			<meta property="og:description" content={resume.data.personal_info.title} />
		{/if}
	{/if}
</svelte:head>

<div class="flex min-h-screen flex-col bg-muted/40">
	<!-- Top bar (hidden when printing) -->
	<header class="sticky top-0 z-10 border-b bg-background/80 backdrop-blur print:hidden">
		<div class="container mx-auto flex h-16 items-center justify-between px-4">
			<a href="/" class="flex items-center space-x-2 transition-opacity hover:opacity-80">
				<FileText class="h-6 w-6 text-primary" />
				<span class="text-xl font-semibold tracking-tight">ResumeBuilder</span>
			</a>
			{#if resume}
				<div class="flex items-center gap-2">
					<Button variant="outline" size="sm" on:click={print}>
						<Printer class="mr-2 h-4 w-4" />
						Print
					</Button>
					<a href="/signup">
						<Button size="sm">
							Build your own
							<ArrowRight class="ml-2 h-4 w-4" />
						</Button>
					</a>
				</div>
			{/if}
		</div>
	</header>

	<main class="flex-1">
		{#if isLoading}
			<div class="flex min-h-[60vh] items-center justify-center">
				<div class="text-center">
					<Loader2 class="mx-auto mb-4 h-8 w-8 animate-spin text-primary" />
					<p class="text-muted-foreground">Loading resume…</p>
				</div>
			</div>
		{:else if notFound || !resume}
			<div class="flex min-h-[60vh] items-center justify-center px-4">
				<div class="max-w-md text-center">
					<div class="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-full bg-muted">
						<FileQuestion class="h-7 w-7 text-muted-foreground" />
					</div>
					<h1 class="text-2xl font-semibold">Resume not available</h1>
					<p class="mt-2 text-muted-foreground">
						This resume is private or the link is no longer valid.
					</p>
					<a href="/" class="mt-6 inline-block">
						<Button>Go to homepage</Button>
					</a>
				</div>
			</div>
		{:else}
			<div class="mx-auto max-w-[850px] px-4 py-8 print:p-0">
				<div class="print:shadow-none">
					<ResumeDocument {resume} {template} />
				</div>

				<!-- Footer CTA (hidden when printing) -->
				<div class="mt-8 flex flex-col items-center gap-3 text-center print:hidden">
					<p class="text-sm text-muted-foreground">
						Made with ResumeBuilder — create an ATS-optimized resume in minutes.
					</p>
					<a href="/signup">
						<Button variant="outline">
							Get started free
							<ArrowRight class="ml-2 h-4 w-4" />
						</Button>
					</a>
				</div>
			</div>
		{/if}
	</main>
</div>
