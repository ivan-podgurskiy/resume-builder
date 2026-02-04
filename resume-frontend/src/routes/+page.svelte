<script lang="ts">
	import { Button, Card } from '$components/ui';
	import { FileText, Sparkles, Download, Shield, Zap, Users } from 'lucide-svelte';
	import { auth, isAuthenticated } from '$stores/auth';
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
	<nav class="border-b">
		<div class="container mx-auto flex h-16 items-center justify-between px-4">
			<a href="/" class="flex items-center space-x-2">
				<FileText class="h-6 w-6 text-primary" />
				<span class="text-xl font-bold">ResumeBuilder</span>
			</a>
			<div class="flex items-center space-x-4">
				{#if $isAuthenticated}
					<a href="/dashboard">
						<Button>Dashboard</Button>
					</a>
				{:else}
					<a href="/login" class="text-sm text-muted-foreground hover:text-foreground">Login</a>
					<a href="/signup">
						<Button>Get Started Free</Button>
					</a>
				{/if}
			</div>
		</div>
	</nav>

	<!-- Hero Section -->
	<section class="py-20 lg:py-32">
		<div class="container mx-auto px-4 text-center">
			<h1 class="mb-6 text-4xl font-bold tracking-tight sm:text-5xl lg:text-6xl">
				Create Your Perfect Resume
				<span class="text-primary">in Minutes</span>
			</h1>
			<p class="mx-auto mb-8 max-w-2xl text-lg text-muted-foreground">
				AI-powered resume builder that helps you create professional, ATS-optimized resumes.
				Import your existing resume, get smart suggestions, and stand out from the crowd.
			</p>
			<div class="flex flex-col items-center justify-center gap-4 sm:flex-row">
				<a href="/signup">
					<Button size="lg" class="px-8">
						<Sparkles class="mr-2 h-5 w-5" />
						Get Started Free
					</Button>
				</a>
				<a href="#features">
					<Button variant="outline" size="lg">See How It Works</Button>
				</a>
			</div>
			<p class="mt-4 text-sm text-muted-foreground">No credit card required</p>
		</div>
	</section>

	<!-- Features Section -->
	<section id="features" class="border-t bg-muted/50 py-20">
		<div class="container mx-auto px-4">
			<h2 class="mb-12 text-center text-3xl font-bold">Everything You Need</h2>
			<div class="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
				<Card class="p-6">
					<Sparkles class="mb-4 h-10 w-10 text-primary" />
					<h3 class="mb-2 text-xl font-semibold">AI-Powered Extraction</h3>
					<p class="text-muted-foreground">
						Upload your existing resume and let AI extract all your information automatically.
						Save hours of manual data entry.
					</p>
				</Card>

				<Card class="p-6">
					<Zap class="mb-4 h-10 w-10 text-primary" />
					<h3 class="mb-2 text-xl font-semibold">Smart Content Suggestions</h3>
					<p class="text-muted-foreground">
						Get AI-powered suggestions to improve your content. Transform weak bullet points
						into impactful achievements.
					</p>
				</Card>

				<Card class="p-6">
					<Shield class="mb-4 h-10 w-10 text-primary" />
					<h3 class="mb-2 text-xl font-semibold">ATS Optimization</h3>
					<p class="text-muted-foreground">
						Ensure your resume passes Applicant Tracking Systems with our built-in
						compatibility checker and keyword optimizer.
					</p>
				</Card>

				<Card class="p-6">
					<FileText class="mb-4 h-10 w-10 text-primary" />
					<h3 class="mb-2 text-xl font-semibold">Professional Templates</h3>
					<p class="text-muted-foreground">
						Choose from 10+ professionally designed templates. Customize colors, fonts,
						and layouts to match your style.
					</p>
				</Card>

				<Card class="p-6">
					<Download class="mb-4 h-10 w-10 text-primary" />
					<h3 class="mb-2 text-xl font-semibold">Multiple Export Formats</h3>
					<p class="text-muted-foreground">
						Export your resume as PDF, DOCX, or plain text. Perfect for any application
						requirement.
					</p>
				</Card>

				<Card class="p-6">
					<Users class="mb-4 h-10 w-10 text-primary" />
					<h3 class="mb-2 text-xl font-semibold">Multiple Versions</h3>
					<p class="text-muted-foreground">
						Create multiple versions of your resume tailored to different job applications.
						Version history keeps your edits safe.
					</p>
				</Card>
			</div>
		</div>
	</section>

	<!-- How It Works -->
	<section class="py-20">
		<div class="container mx-auto px-4">
			<h2 class="mb-12 text-center text-3xl font-bold">How It Works</h2>
			<div class="grid gap-8 md:grid-cols-3">
				<div class="text-center">
					<div
						class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-primary text-xl font-bold text-primary-foreground"
					>
						1
					</div>
					<h3 class="mb-2 text-xl font-semibold">Import or Start Fresh</h3>
					<p class="text-muted-foreground">
						Upload your existing resume or start from a blank template. Our AI will extract
						your information in seconds.
					</p>
				</div>

				<div class="text-center">
					<div
						class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-primary text-xl font-bold text-primary-foreground"
					>
						2
					</div>
					<h3 class="mb-2 text-xl font-semibold">Edit & Enhance</h3>
					<p class="text-muted-foreground">
						Use our intuitive editor with AI suggestions to refine your content. Get real-time
						preview as you make changes.
					</p>
				</div>

				<div class="text-center">
					<div
						class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-primary text-xl font-bold text-primary-foreground"
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
	</section>

	<!-- CTA Section -->
	<section class="border-t bg-primary py-20 text-primary-foreground">
		<div class="container mx-auto px-4 text-center">
			<h2 class="mb-4 text-3xl font-bold">Ready to Build Your Perfect Resume?</h2>
			<p class="mx-auto mb-8 max-w-xl opacity-90">
				Join thousands of job seekers who have landed their dream jobs with professionally
				crafted resumes.
			</p>
			<a href="/signup">
				<Button size="lg" variant="secondary" class="px-8">
					Get Started Free
				</Button>
			</a>
		</div>
	</section>

	<!-- Footer -->
	<footer class="border-t py-12">
		<div class="container mx-auto px-4">
			<div class="flex flex-col items-center justify-between gap-4 md:flex-row">
				<div class="flex items-center space-x-2">
					<FileText class="h-5 w-5 text-primary" />
					<span class="font-semibold">ResumeBuilder</span>
				</div>
				<div class="flex gap-6 text-sm text-muted-foreground">
					<a href="/about" class="hover:text-foreground">About</a>
					<a href="/privacy" class="hover:text-foreground">Privacy</a>
					<a href="/terms" class="hover:text-foreground">Terms</a>
				</div>
				<p class="text-sm text-muted-foreground">
					© {new Date().getFullYear()} ResumeBuilder. All rights reserved.
				</p>
			</div>
		</div>
	</footer>
</div>
