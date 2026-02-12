<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { Button, Card, Input, Label } from '$components/ui';
	import { auth } from '$stores/auth';
	import { FileText } from 'lucide-svelte';

	let name = '';
	let email = '';
	let password = '';
	let confirmPassword = '';
	let error = '';
	let loading = false;

	async function handleSubmit() {
		error = '';

		if (password !== confirmPassword) {
			error = 'Passwords do not match';
			return;
		}

		if (password.length < 8) {
			error = 'Password must be at least 8 characters';
			return;
		}

		loading = true;

		const result = await auth.register(email, password, name);

		if (result.success) {
			const templateId = $page.url.searchParams.get('template');
			goto(templateId ? `/dashboard?template=${templateId}` : '/dashboard');
		} else {
			error = result.error || 'Failed to create account';
		}

		loading = false;
	}
</script>

<svelte:head>
	<title>Sign Up - Resume Builder</title>
</svelte:head>

<div class="flex min-h-screen items-center justify-center bg-muted/50 px-4">
	<Card class="w-full max-w-md p-8">
		<div class="mb-8 text-center">
			<a href="/" class="inline-flex items-center space-x-2">
				<FileText class="h-8 w-8 text-primary" />
				<span class="text-2xl font-bold">ResumeBuilder</span>
			</a>
			<h1 class="mt-6 text-2xl font-semibold">Create your account</h1>
			<p class="mt-2 text-sm text-muted-foreground">
				Start building your professional resume today
			</p>
		</div>

		<form on:submit|preventDefault={handleSubmit} class="space-y-4">
			{#if error}
				<div class="rounded-md bg-destructive/15 p-3 text-sm text-destructive">
					{error}
				</div>
			{/if}

			<div class="space-y-2">
				<Label htmlFor="name">Full Name</Label>
				<Input id="name" type="text" placeholder="John Doe" bind:value={name} required />
			</div>

			<div class="space-y-2">
				<Label htmlFor="email">Email</Label>
				<Input
					id="email"
					type="email"
					placeholder="name@example.com"
					bind:value={email}
					required
				/>
			</div>

			<div class="space-y-2">
				<Label htmlFor="password">Password</Label>
				<Input
					id="password"
					type="password"
					placeholder="••••••••"
					bind:value={password}
					required
				/>
				<p class="text-xs text-muted-foreground">Must be at least 8 characters</p>
			</div>

			<div class="space-y-2">
				<Label htmlFor="confirmPassword">Confirm Password</Label>
				<Input
					id="confirmPassword"
					type="password"
					placeholder="••••••••"
					bind:value={confirmPassword}
					required
				/>
			</div>

			<Button type="submit" class="w-full" disabled={loading}>
				{loading ? 'Creating account...' : 'Create account'}
			</Button>
		</form>

		<p class="mt-6 text-center text-sm text-muted-foreground">
			Already have an account?
			<a href="/login" class="font-medium text-primary hover:underline">Sign in</a>
		</p>

		<p class="mt-4 text-center text-xs text-muted-foreground">
			By creating an account, you agree to our
			<a href="/terms" class="underline">Terms of Service</a>
			and
			<a href="/privacy" class="underline">Privacy Policy</a>.
		</p>
	</Card>
</div>
