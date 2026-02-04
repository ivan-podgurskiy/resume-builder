<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button, Card, Input, Label } from '$components/ui';
	import { auth } from '$stores/auth';
	import { FileText } from 'lucide-svelte';

	let email = '';
	let password = '';
	let error = '';
	let loading = false;

	async function handleSubmit() {
		error = '';
		loading = true;

		const result = await auth.login(email, password);

		if (result.success) {
			goto('/dashboard');
		} else {
			error = result.error || 'Failed to login';
		}

		loading = false;
	}
</script>

<svelte:head>
	<title>Login - Resume Builder</title>
</svelte:head>

<div class="flex min-h-screen items-center justify-center bg-muted/50 px-4">
	<Card class="w-full max-w-md p-8">
		<div class="mb-8 text-center">
			<a href="/" class="inline-flex items-center space-x-2">
				<FileText class="h-8 w-8 text-primary" />
				<span class="text-2xl font-bold">ResumeBuilder</span>
			</a>
			<h1 class="mt-6 text-2xl font-semibold">Welcome back</h1>
			<p class="mt-2 text-sm text-muted-foreground">
				Enter your credentials to access your account
			</p>
		</div>

		<form on:submit|preventDefault={handleSubmit} class="space-y-4">
			{#if error}
				<div class="rounded-md bg-destructive/15 p-3 text-sm text-destructive">
					{error}
				</div>
			{/if}

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
				<div class="flex items-center justify-between">
					<Label htmlFor="password">Password</Label>
					<a href="/forgot-password" class="text-sm text-primary hover:underline">
						Forgot password?
					</a>
				</div>
				<Input
					id="password"
					type="password"
					placeholder="••••••••"
					bind:value={password}
					required
				/>
			</div>

			<Button type="submit" class="w-full" disabled={loading}>
				{loading ? 'Signing in...' : 'Sign in'}
			</Button>
		</form>

		<p class="mt-6 text-center text-sm text-muted-foreground">
			Don't have an account?
			<a href="/signup" class="font-medium text-primary hover:underline">Sign up</a>
		</p>
	</Card>
</div>
