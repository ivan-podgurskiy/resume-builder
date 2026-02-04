<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { Button } from '$components/ui';
	import { auth, isAuthenticated, currentUser } from '$stores/auth';
	import { FileText, LayoutDashboard, Settings, LogOut, User } from 'lucide-svelte';

	let showUserMenu = false;

	onMount(() => {
		// Check auth status after store initializes
		const unsubscribe = auth.subscribe((state) => {
			if (state.initialized && !state.user) {
				goto('/login');
			}
		});

		return unsubscribe;
	});

	async function handleLogout() {
		await auth.logout();
		goto('/');
	}
</script>

{#if $isAuthenticated}
	<div class="flex min-h-screen flex-col">
		<!-- Top Navigation -->
		<nav class="border-b bg-background">
			<div class="flex h-16 items-center justify-between px-6">
				<a href="/dashboard" class="flex items-center space-x-2">
					<FileText class="h-6 w-6 text-primary" />
					<span class="text-xl font-bold">ResumeBuilder</span>
				</a>

				<div class="flex items-center space-x-4">
					<div class="relative">
						<button
							class="flex items-center space-x-2 rounded-md px-3 py-2 hover:bg-muted"
							on:click={() => (showUserMenu = !showUserMenu)}
						>
							<div
								class="flex h-8 w-8 items-center justify-center rounded-full bg-primary text-primary-foreground"
							>
								{$currentUser?.name?.charAt(0).toUpperCase() || 'U'}
							</div>
							<span class="text-sm font-medium">{$currentUser?.name}</span>
						</button>

						{#if showUserMenu}
							<div
								class="absolute right-0 top-full z-50 mt-2 w-48 rounded-md border bg-background py-1 shadow-lg"
							>
								<a
									href="/settings"
									class="flex items-center px-4 py-2 text-sm hover:bg-muted"
									on:click={() => (showUserMenu = false)}
								>
									<Settings class="mr-2 h-4 w-4" />
									Settings
								</a>
								<button
									class="flex w-full items-center px-4 py-2 text-sm text-destructive hover:bg-muted"
									on:click={handleLogout}
								>
									<LogOut class="mr-2 h-4 w-4" />
									Sign out
								</button>
							</div>
						{/if}
					</div>
				</div>
			</div>
		</nav>

		<!-- Main Content -->
		<main class="flex-1">
			<slot />
		</main>
	</div>
{:else}
	<div class="flex min-h-screen items-center justify-center">
		<div class="text-center">
			<div class="mb-4 h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent mx-auto"></div>
			<p class="text-muted-foreground">Loading...</p>
		</div>
	</div>
{/if}
