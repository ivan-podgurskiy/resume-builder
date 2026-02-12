<script lang="ts">
	import { Button, Card } from '$components/ui';
	import { FileText, Sparkles, Download, Shield, Zap, Users } from 'lucide-svelte';
	import { isAuthenticated } from '$stores/auth';

	let heroImageError = false;
</script>

<svelte:head>
	<title>Resume Builder - Create Professional Resumes with AI</title>
	<meta
		name="description"
		content="Create professional, ATS-optimized resumes in minutes with AI assistance. Import your existing resume, get smart suggestions, and export in multiple formats."
	/>
</svelte:head>

<div class="min-h-screen bg-background">
	<!-- Navigation -->
	<nav
		class="sticky top-0 z-50 border-b border-border/40 bg-background/80 backdrop-blur-xl transition-colors"
	>
		<div class="container mx-auto flex h-16 items-center justify-between px-4">
			<a href="/" class="flex items-center space-x-2 transition-opacity hover:opacity-80">
				<FileText class="h-6 w-6 text-primary" />
				<span class="text-xl font-semibold tracking-tight">ResumeBuilder</span>
			</a>
			<div class="flex items-center gap-6">
				<a
					href="#features"
					class="text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
				>
					Features
				</a>
				<a
					href="#templates"
					class="text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
				>
					Templates
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
					<a
						href="/login"
						class="text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
					>
						Login
					</a>
					<a href="/signup">
						<Button>Get Started Free</Button>
					</a>
				{/if}
			</div>
		</div>
	</nav>

	<!-- Hero Section -->
	<section
		class="relative overflow-hidden py-24 lg:py-36"
		style="background: radial-gradient(ellipse 80% 60% at 50% 0%, hsl(var(--primary) / 0.08), transparent 60%), linear-gradient(180deg, hsl(var(--hero-gradient-from)), hsl(var(--hero-gradient-via)), hsl(var(--hero-gradient-to)))"
	>
		<!-- Subtle grid pattern -->
		<div
			class="pointer-events-none absolute inset-0 opacity-[0.02]"
			style="background-image: linear-gradient(to right, currentColor 1px, transparent 1px), linear-gradient(to bottom, currentColor 1px, transparent 1px); background-size: 3rem 3rem"
		></div>

		<div class="container relative mx-auto px-4">
			<div class="flex flex-col items-center gap-12 lg:flex-row lg:items-center lg:justify-between lg:gap-16">
				<!-- Left: Value prop & CTA -->
				<div class="flex-1 text-center lg:text-left">
					<p
						class="mb-4 inline-flex items-center gap-1.5 rounded-full border border-primary/20 bg-primary/5 px-4 py-1.5 text-sm font-medium text-primary"
					>
						<Sparkles class="h-4 w-4" />
						AI-powered resume creation
					</p>
					<h1
						class="font-display mb-6 text-4xl font-normal tracking-tight text-foreground sm:text-5xl lg:text-6xl xl:text-7xl"
					>
						Build a Job-Winning
						<span class="text-primary">Resume</span>
						<br />
						<span class="text-muted-foreground">That Gets Noticed</span>
					</h1>
					<p class="mx-auto mb-10 max-w-xl text-lg leading-relaxed text-muted-foreground lg:mx-0">
						Import your existing resume and let AI extract your info in seconds. Get smart
						suggestions, optimize for ATS, and export to PDF—all in one place.
					</p>
					<div class="flex flex-col items-center gap-4 sm:flex-row lg:justify-start">
						<a href="/signup">
							<Button
								size="lg"
								class="group px-8 shadow-lg shadow-primary/25 transition-all hover:shadow-xl hover:shadow-primary/30"
							>
								<Sparkles class="mr-2 h-5 w-5 transition-transform group-hover:scale-110" />
								Get Started Free
							</Button>
						</a>
						<a href="#features">
							<Button variant="outline" size="lg" class="border-2">
								See How It Works
							</Button>
						</a>
					</div>
					<p class="mt-6 text-sm text-muted-foreground">No credit card required · Free to start</p>
				</div>

				<!-- Right: Resume preview image (add static/images/hero-resume-preview.png) -->
				<div class="flex flex-1 items-center justify-center lg:justify-end">
					<div
						class="relative w-full max-w-md overflow-hidden rounded-lg border border-border bg-card shadow-xl"
						style="aspect-ratio: 8.5/11"
					>
						{#if heroImageError}
							<div
								class="flex h-full w-full items-center justify-center bg-muted/50 text-muted-foreground"
								aria-hidden="true"
							>
								<FileText class="h-24 w-24 opacity-30" />
							</div>
						{:else}
							<img
								src="/images/hero-resume-preview.png"
								alt="Professional resume preview"
								class="h-full w-full object-cover object-top"
								on:error={() => (heroImageError = true)}
							/>
						{/if}
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- Templates Section -->
	<section id="templates" class="scroll-mt-20 border-t py-24 lg:py-32">
		<div class="container mx-auto px-4">
			<div class="mx-auto max-w-2xl text-center">
				<h2 class="font-display mb-4 text-3xl font-normal tracking-tight sm:text-4xl">
					Professional Templates
				</h2>
				<p class="text-muted-foreground">
					Choose from professionally designed templates. Customize colors, fonts, and layouts to
					match your style.
				</p>
			</div>
			<!-- Template grid placeholder - to be populated with API data -->
			<div class="mt-12 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
				{#each [1, 2, 3] as _}
					<Card class="aspect-[3/4] bg-muted/50"></Card>
				{/each}
			</div>
			<div class="mt-8 text-center">
				<a href="/signup">
					<Button variant="outline">Browse All Templates</Button>
				</a>
			</div>
		</div>
	</section>

	<!-- Features Section -->
	<section id="features" class="scroll-mt-20 border-t py-24 lg:py-32">
		<div class="container mx-auto px-4">
			<div class="mx-auto max-w-2xl text-center">
				<h2 class="font-display mb-4 text-3xl font-normal tracking-tight sm:text-4xl">
					Everything You Need
				</h2>
				<p class="text-muted-foreground">
					From AI extraction to ATS optimization—built for modern job seekers.
				</p>
			</div>

			<!-- Feature grid -->
			<div class="mt-16 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
				<Card
					class="group flex flex-col p-6 transition-all duration-300 hover:shadow-lg hover:shadow-primary/5 hover:border-primary/20"
				>
					<div
						class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-primary/10 text-primary transition-colors group-hover:bg-primary/20"
					>
						<Sparkles class="h-6 w-6" />
					</div>
					<h3 class="mb-2 text-xl font-semibold">AI-Powered Extraction</h3>
					<p class="text-muted-foreground">
						Upload your existing resume and let AI extract all your information automatically.
						Save hours of manual data entry.
					</p>
				</Card>

				<Card class="group flex flex-col p-6 transition-all duration-300 hover:shadow-md hover:border-primary/20">
					<div
						class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-primary/10 text-primary transition-colors group-hover:bg-primary/20"
					>
						<Zap class="h-6 w-6" />
					</div>
					<h3 class="mb-2 text-lg font-semibold">Smart Suggestions</h3>
					<p class="text-muted-foreground">
						Get AI-powered suggestions to transform weak bullet points into impactful achievements.
					</p>
				</Card>

				<Card class="group flex flex-col p-6 transition-all duration-300 hover:shadow-md hover:border-primary/20">
					<div
						class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-primary/10 text-primary transition-colors group-hover:bg-primary/20"
					>
						<Shield class="h-6 w-6" />
					</div>
					<h3 class="mb-2 text-lg font-semibold">ATS Optimization</h3>
					<p class="text-muted-foreground">
						Ensure your resume passes Applicant Tracking Systems with our built-in compatibility checker.
					</p>
				</Card>

				<Card class="group flex flex-col p-6 transition-all duration-300 hover:shadow-md hover:border-primary/20">
					<div
						class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-primary/10 text-primary transition-colors group-hover:bg-primary/20"
					>
						<FileText class="h-6 w-6" />
					</div>
					<h3 class="mb-2 text-lg font-semibold">Professional Templates</h3>
					<p class="text-muted-foreground">
						10+ professionally designed templates. Customize colors, fonts, and layouts.
					</p>
				</Card>

				<Card class="group flex flex-col p-6 transition-all duration-300 hover:shadow-md hover:border-primary/20">
					<div
						class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-primary/10 text-primary transition-colors group-hover:bg-primary/20"
					>
						<Download class="h-6 w-6" />
					</div>
					<h3 class="mb-2 text-lg font-semibold">Multiple Exports</h3>
					<p class="text-muted-foreground">
						Export as PDF, DOCX, or plain text. Perfect for any application requirement.
					</p>
				</Card>

				<Card class="group flex flex-col p-6 transition-all duration-300 hover:shadow-md hover:border-primary/20">
					<div
						class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-primary/10 text-primary transition-colors group-hover:bg-primary/20"
					>
						<Users class="h-6 w-6" />
					</div>
					<h3 class="mb-2 text-lg font-semibold">Multiple Versions</h3>
					<p class="text-muted-foreground">
						Create tailored versions for different jobs. Version history keeps your edits safe.
					</p>
				</Card>
			</div>
		</div>
	</section>

	<!-- How It Works -->
	<section class="border-t bg-muted/30 py-24 lg:py-32">
		<div class="container mx-auto px-4">
			<div class="mx-auto max-w-2xl text-center">
				<h2 class="font-display mb-4 text-3xl font-normal tracking-tight sm:text-4xl">
					How It Works
				</h2>
				<p class="text-muted-foreground">Three simple steps to your perfect resume.</p>
			</div>

			<div class="relative mt-16">
				<!-- Connection line (desktop) -->
				<div
					class="absolute left-1/2 top-12 hidden h-[calc(100%-6rem)] w-px -translate-x-1/2 bg-gradient-to-b from-primary/40 via-primary/20 to-transparent md:block"
				></div>

				<div class="grid gap-12 md:grid-cols-3 md:gap-8">
					<div class="relative flex flex-col items-center text-center md:items-start md:text-left">
						<div
							class="mb-6 flex h-14 w-14 shrink-0 items-center justify-center rounded-2xl bg-primary text-2xl font-bold text-primary-foreground shadow-lg shadow-primary/25"
						>
							1
						</div>
						<h3 class="mb-2 text-xl font-semibold">Import or Start Fresh</h3>
						<p class="text-muted-foreground">
							Upload your existing resume or start from a blank template. Our AI will extract
							your information in seconds.
						</p>
					</div>

					<div class="relative flex flex-col items-center text-center md:items-start md:text-left">
						<div
							class="mb-6 flex h-14 w-14 shrink-0 items-center justify-center rounded-2xl bg-primary text-2xl font-bold text-primary-foreground shadow-lg shadow-primary/25"
						>
							2
						</div>
						<h3 class="mb-2 text-xl font-semibold">Edit & Enhance</h3>
						<p class="text-muted-foreground">
							Use our intuitive editor with AI suggestions to refine your content. Get real-time
							preview as you make changes.
						</p>
					</div>

					<div class="relative flex flex-col items-center text-center md:items-start md:text-left">
						<div
							class="mb-6 flex h-14 w-14 shrink-0 items-center justify-center rounded-2xl bg-primary text-2xl font-bold text-primary-foreground shadow-lg shadow-primary/25"
						>
							3
						</div>
						<h3 class="mb-2 text-xl font-semibold">Export & Apply</h3>
						<p class="text-muted-foreground">
							Download your polished resume in your preferred format and start applying to
							your dream jobs.
						</p>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- CTA Section -->
	<section
		class="relative overflow-hidden border-t py-24"
		style="background: linear-gradient(135deg, hsl(var(--primary)), hsl(221.2 83.2% 45%))"
	>
		<div
			class="pointer-events-none absolute inset-0 opacity-10"
			style="background-image: radial-gradient(circle at 20% 80%, white 1px, transparent 1px); background-size: 2rem 2rem"
		></div>

		<div class="container relative mx-auto px-4 text-center">
			<h2 class="font-display mb-4 text-3xl font-normal tracking-tight text-white sm:text-4xl">
				Ready to Build Your Perfect Resume?
			</h2>
			<p class="mx-auto mb-10 max-w-xl text-lg text-white/90">
				Join thousands of job seekers who have landed their dream jobs with professionally
				crafted resumes.
			</p>
			<a href="/signup">
				<Button
					size="lg"
					variant="secondary"
					class="bg-white px-8 font-semibold text-primary shadow-xl transition-all hover:bg-white/95 hover:shadow-2xl"
				>
					Get Started Free
				</Button>
			</a>
		</div>
	</section>

	<!-- Footer -->
	<footer class="border-t py-12">
		<div class="container mx-auto px-4">
			<div class="flex flex-col items-center justify-between gap-6 md:flex-row">
				<a href="/" class="flex items-center space-x-2 transition-opacity hover:opacity-80">
					<FileText class="h-5 w-5 text-primary" />
					<span class="font-semibold">ResumeBuilder</span>
				</a>
				<div class="flex gap-8 text-sm text-muted-foreground">
					<a href="/about" class="transition-colors hover:text-foreground">About</a>
					<a href="/privacy" class="transition-colors hover:text-foreground">Privacy</a>
					<a href="/terms" class="transition-colors hover:text-foreground">Terms</a>
				</div>
				<p class="text-sm text-muted-foreground">
					© {new Date().getFullYear()} ResumeBuilder
				</p>
			</div>
		</div>
	</footer>
</div>
