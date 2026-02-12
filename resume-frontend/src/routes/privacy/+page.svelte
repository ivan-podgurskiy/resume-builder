<script lang="ts">
	import { Button } from '$components/ui';
	import { FileText } from 'lucide-svelte';
	import { isAuthenticated } from '$stores/auth';
</script>

<svelte:head>
	<title>Privacy Policy - Resume Builder</title>
	<meta
		name="description"
		content="Privacy Policy for Resume Builder - Learn how we collect, use, and protect your personal information."
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

	<!-- Content -->
	<section class="py-16">
		<div class="container mx-auto px-4">
			<div class="mx-auto max-w-3xl">
				<h1 class="mb-8 text-4xl font-bold">Privacy Policy</h1>
				<p class="mb-8 text-muted-foreground">Last updated: February 4, 2026</p>

				<div class="prose prose-neutral dark:prose-invert max-w-none space-y-8">
					<section>
						<h2 class="mb-4 text-2xl font-semibold">1. Introduction</h2>
						<p class="text-muted-foreground">
							At ResumeBuilder, we take your privacy seriously. This Privacy Policy explains how we
							collect, use, disclose, and safeguard your information when you use our resume
							building service. Please read this privacy policy carefully.
						</p>
					</section>

					<section>
						<h2 class="mb-4 text-2xl font-semibold">2. Information We Collect</h2>
						<p class="mb-4 text-muted-foreground">
							We collect information that you provide directly to us, including:
						</p>
						<ul class="list-disc space-y-2 pl-6 text-muted-foreground">
							<li>
								<strong>Account Information:</strong> Name, email address, and password when you create
								an account.
							</li>
							<li>
								<strong>Resume Content:</strong> Personal information, work history, education, skills,
								and other details you include in your resumes.
							</li>
							<li>
								<strong>Uploaded Files:</strong> Documents you upload for AI extraction, such as existing
								resumes or CVs.
							</li>
							<li>
								<strong>Usage Data:</strong> Information about how you interact with our service, including
								features used and time spent.
							</li>
						</ul>
					</section>

					<section>
						<h2 class="mb-4 text-2xl font-semibold">3. How We Use Your Information</h2>
						<p class="mb-4 text-muted-foreground">We use the information we collect to:</p>
						<ul class="list-disc space-y-2 pl-6 text-muted-foreground">
							<li>Provide, maintain, and improve our resume building services</li>
							<li>Process and store your resumes securely</li>
							<li>Provide AI-powered suggestions and content improvements</li>
							<li>Send you technical notices, updates, and support messages</li>
							<li>Respond to your comments, questions, and customer service requests</li>
							<li>Analyze usage patterns to improve our service</li>
						</ul>
					</section>

					<section>
						<h2 class="mb-4 text-2xl font-semibold">4. Data Security</h2>
						<p class="text-muted-foreground">
							We implement appropriate technical and organizational security measures to protect
							your personal information. Your resume data is encrypted both in transit and at rest.
							However, no method of transmission over the Internet is 100% secure, and we cannot
							guarantee absolute security.
						</p>
					</section>

					<section>
						<h2 class="mb-4 text-2xl font-semibold">5. Data Retention</h2>
						<p class="text-muted-foreground">
							We retain your personal information for as long as your account is active or as needed
							to provide you services. You can delete your account and associated data at any time
							through your account settings. We may retain certain information as required by law or
							for legitimate business purposes.
						</p>
					</section>

					<section>
						<h2 class="mb-4 text-2xl font-semibold">6. Third-Party Services</h2>
						<p class="text-muted-foreground">
							We may use third-party services to help us operate our service, including cloud
							hosting providers and AI service providers. These third parties have access to your
							information only to perform specific tasks on our behalf and are obligated not to
							disclose or use it for any other purpose.
						</p>
					</section>

					<section>
						<h2 class="mb-4 text-2xl font-semibold">7. Your Rights</h2>
						<p class="mb-4 text-muted-foreground">You have the right to:</p>
						<ul class="list-disc space-y-2 pl-6 text-muted-foreground">
							<li>Access the personal information we hold about you</li>
							<li>Request correction of inaccurate information</li>
							<li>Request deletion of your account and data</li>
							<li>Export your resume data in standard formats</li>
							<li>Opt out of marketing communications</li>
						</ul>
					</section>

					<section>
						<h2 class="mb-4 text-2xl font-semibold">8. Cookies</h2>
						<p class="text-muted-foreground">
							We use cookies and similar tracking technologies to track activity on our service and
							hold certain information. Cookies are files with a small amount of data which may
							include an anonymous unique identifier. You can instruct your browser to refuse all
							cookies or to indicate when a cookie is being sent.
						</p>
					</section>

					<section>
						<h2 class="mb-4 text-2xl font-semibold">9. Changes to This Policy</h2>
						<p class="text-muted-foreground">
							We may update our Privacy Policy from time to time. We will notify you of any changes
							by posting the new Privacy Policy on this page and updating the "Last updated" date.
							You are advised to review this Privacy Policy periodically for any changes.
						</p>
					</section>

					<section>
						<h2 class="mb-4 text-2xl font-semibold">10. Contact Us</h2>
						<p class="text-muted-foreground">
							If you have any questions about this Privacy Policy, please contact us at
							privacy@resumebuilder.com.
						</p>
					</section>
				</div>
			</div>
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
