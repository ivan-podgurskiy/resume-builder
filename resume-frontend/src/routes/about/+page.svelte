<script lang="ts">
	import { Button } from '$components/ui';
	import { FileText, Target, Heart, Users } from 'lucide-svelte';
	import { isAuthenticated } from '$stores/auth';
</script>

<svelte:head>
	<title>About Us - Resume Builder</title>
	<meta
		name="description"
		content="Learn about Resume Builder - our mission to help job seekers create professional, ATS-optimized resumes with AI assistance."
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
	<section class="py-16 lg:py-24">
		<div class="container mx-auto px-4 text-center">
			<h1 class="mb-6 text-4xl font-bold tracking-tight sm:text-5xl">About ResumeBuilder</h1>
			<p class="mx-auto max-w-2xl text-lg text-muted-foreground">
				We're on a mission to help job seekers create professional, standout resumes that land
				interviews.
			</p>
		</div>
	</section>

	<!-- Mission Section -->
	<section class="border-t bg-muted/50 py-16">
		<div class="container mx-auto px-4">
			<div class="mx-auto max-w-3xl">
				<div class="mb-12 text-center">
					<Target class="mx-auto mb-4 h-12 w-12 text-primary" />
					<h2 class="mb-4 text-3xl font-bold">Our Mission</h2>
					<p class="text-lg text-muted-foreground">
						Job searching can be stressful. We believe that creating a professional resume
						shouldn't add to that stress. Our mission is to empower every job seeker with the
						tools they need to present their best selves to potential employers.
					</p>
				</div>

				<div class="grid gap-8 md:grid-cols-2">
					<div class="rounded-lg border bg-card p-6">
						<Heart class="mb-4 h-8 w-8 text-primary" />
						<h3 class="mb-2 text-xl font-semibold">What We Believe</h3>
						<p class="text-muted-foreground">
							Everyone deserves access to professional resume tools. We've built ResumeBuilder to
							be accessible, intuitive, and powerful enough to help anyone create a resume that
							gets noticed.
						</p>
					</div>

					<div class="rounded-lg border bg-card p-6">
						<Users class="mb-4 h-8 w-8 text-primary" />
						<h3 class="mb-2 text-xl font-semibold">Who We Help</h3>
						<p class="text-muted-foreground">
							From fresh graduates to seasoned professionals, career changers to executives, we
							help job seekers at every stage create resumes that showcase their unique value.
						</p>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- Story Section -->
	<section class="py-16">
		<div class="container mx-auto px-4">
			<div class="mx-auto max-w-3xl">
				<h2 class="mb-6 text-center text-3xl font-bold">Our Story</h2>
				<div class="space-y-4 text-muted-foreground">
					<p>
						ResumeBuilder was born from a simple frustration: why is it so hard to create a great
						resume? Traditional resume builders were clunky, templates looked outdated, and none of
						them understood the modern job market.
					</p>
					<p>
						We set out to build something better. By combining the latest in AI technology with
						deep knowledge of hiring practices, we created a tool that not only makes resume
						building easy but actually helps you create a better resume than you could on your own.
					</p>
					<p>
						Today, thousands of job seekers use ResumeBuilder to create professional, ATS-optimized
						resumes that help them land their dream jobs. And we're just getting started.
					</p>
				</div>
			</div>
		</div>
	</section>

	<!-- CTA Section -->
	<section class="border-t bg-primary py-16 text-primary-foreground">
		<div class="container mx-auto px-4 text-center">
			<h2 class="mb-4 text-3xl font-bold">Ready to Get Started?</h2>
			<p class="mx-auto mb-8 max-w-xl opacity-90">
				Join thousands of job seekers who have transformed their job search with ResumeBuilder.
			</p>
			<a href="/signup">
				<Button size="lg" variant="secondary" class="px-8">Create Your Resume</Button>
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
