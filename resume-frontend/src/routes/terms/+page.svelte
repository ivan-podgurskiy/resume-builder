<script lang="ts">
	import { Button } from '$components/ui';
	import { FileText } from 'lucide-svelte';
	import { isAuthenticated } from '$stores/auth';
</script>

<svelte:head>
	<title>Terms of Service - Resume Builder</title>
	<meta
		name="description"
		content="Terms of Service for Resume Builder - Read our terms and conditions for using our resume building service."
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
				<h1 class="mb-8 text-4xl font-bold">Terms of Service</h1>
				<p class="mb-8 text-muted-foreground">Last updated: February 4, 2026</p>

				<div class="prose prose-neutral dark:prose-invert max-w-none space-y-8">
					<section>
						<h2 class="mb-4 text-2xl font-semibold">1. Acceptance of Terms</h2>
						<p class="text-muted-foreground">
							By accessing or using ResumeBuilder, you agree to be bound by these Terms of Service
							and all applicable laws and regulations. If you do not agree with any of these terms,
							you are prohibited from using or accessing this service.
						</p>
					</section>

					<section>
						<h2 class="mb-4 text-2xl font-semibold">2. Description of Service</h2>
						<p class="text-muted-foreground">
							ResumeBuilder is an online platform that allows users to create, edit, and export
							professional resumes. Our service includes AI-powered features for content extraction
							and suggestions, multiple export formats, and template customization.
						</p>
					</section>

					<section>
						<h2 class="mb-4 text-2xl font-semibold">3. User Accounts</h2>
						<p class="mb-4 text-muted-foreground">To use our service, you must:</p>
						<ul class="list-disc space-y-2 pl-6 text-muted-foreground">
							<li>Be at least 16 years of age</li>
							<li>Provide accurate and complete registration information</li>
							<li>Maintain the security of your account credentials</li>
							<li>Notify us immediately of any unauthorized use of your account</li>
							<li>Accept responsibility for all activities that occur under your account</li>
						</ul>
					</section>

					<section>
						<h2 class="mb-4 text-2xl font-semibold">4. User Content</h2>
						<p class="mb-4 text-muted-foreground">
							You retain ownership of all content you create or upload to ResumeBuilder. By using
							our service, you grant us a limited license to:
						</p>
						<ul class="list-disc space-y-2 pl-6 text-muted-foreground">
							<li>Store and process your content to provide our services</li>
							<li>Use your content to improve our AI features (in anonymized form only)</li>
							<li>Create backups for service continuity</li>
						</ul>
						<p class="mt-4 text-muted-foreground">
							You are solely responsible for the accuracy and legality of the content you include in
							your resumes.
						</p>
					</section>

					<section>
						<h2 class="mb-4 text-2xl font-semibold">5. Acceptable Use</h2>
						<p class="mb-4 text-muted-foreground">You agree not to:</p>
						<ul class="list-disc space-y-2 pl-6 text-muted-foreground">
							<li>Use the service for any unlawful purpose</li>
							<li>Include false, misleading, or fraudulent information in your resumes</li>
							<li>
								Attempt to gain unauthorized access to our systems or other users' accounts
							</li>
							<li>Interfere with or disrupt the service or servers</li>
							<li>Use automated systems to access the service without permission</li>
							<li>Resell or redistribute the service without authorization</li>
						</ul>
					</section>

					<section>
						<h2 class="mb-4 text-2xl font-semibold">6. Intellectual Property</h2>
						<p class="text-muted-foreground">
							The ResumeBuilder service, including its original content, features, and
							functionality, is owned by ResumeBuilder and is protected by international copyright,
							trademark, patent, trade secret, and other intellectual property laws. Our templates
							and designs are licensed for personal use only.
						</p>
					</section>

					<section>
						<h2 class="mb-4 text-2xl font-semibold">7. Service Availability</h2>
						<p class="text-muted-foreground">
							We strive to maintain high availability of our service but do not guarantee
							uninterrupted access. We may modify, suspend, or discontinue any part of the service
							at any time. We will provide reasonable notice of any significant changes that affect
							your use of the service.
						</p>
					</section>

					<section>
						<h2 class="mb-4 text-2xl font-semibold">8. Limitation of Liability</h2>
						<p class="text-muted-foreground">
							ResumeBuilder and its affiliates shall not be liable for any indirect, incidental,
							special, consequential, or punitive damages resulting from your use of or inability to
							use the service. This includes, but is not limited to, damages for loss of profits,
							data, or other intangible losses, even if we have been advised of the possibility of
							such damages.
						</p>
					</section>

					<section>
						<h2 class="mb-4 text-2xl font-semibold">9. Disclaimer of Warranties</h2>
						<p class="text-muted-foreground">
							The service is provided "as is" and "as available" without warranties of any kind,
							either express or implied. We do not warrant that the service will be error-free,
							secure, or uninterrupted. We do not guarantee that your use of the service will result
							in employment or job interviews.
						</p>
					</section>

					<section>
						<h2 class="mb-4 text-2xl font-semibold">10. Termination</h2>
						<p class="text-muted-foreground">
							We may terminate or suspend your account and access to the service immediately,
							without prior notice or liability, for any reason, including breach of these Terms.
							Upon termination, your right to use the service will immediately cease. You may
							request deletion of your data in accordance with our Privacy Policy.
						</p>
					</section>

					<section>
						<h2 class="mb-4 text-2xl font-semibold">11. Changes to Terms</h2>
						<p class="text-muted-foreground">
							We reserve the right to modify or replace these Terms at any time. If a revision is
							material, we will provide at least 30 days' notice prior to any new terms taking
							effect. What constitutes a material change will be determined at our sole discretion.
						</p>
					</section>

					<section>
						<h2 class="mb-4 text-2xl font-semibold">12. Governing Law</h2>
						<p class="text-muted-foreground">
							These Terms shall be governed by and construed in accordance with the laws of the
							jurisdiction in which ResumeBuilder operates, without regard to its conflict of law
							provisions.
						</p>
					</section>

					<section>
						<h2 class="mb-4 text-2xl font-semibold">13. Contact Us</h2>
						<p class="text-muted-foreground">
							If you have any questions about these Terms, please contact us at
							legal@resumebuilder.com.
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
