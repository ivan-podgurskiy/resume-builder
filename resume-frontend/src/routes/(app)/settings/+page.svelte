<script lang="ts">
	import { Button, Card, Input, Label } from '$components/ui';
	import { currentUser } from '$stores/auth';
	import { User, CreditCard, Bell, Shield } from 'lucide-svelte';

	let activeTab = 'profile';

	const tabs = [
		{ id: 'profile', label: 'Profile', icon: User },
		{ id: 'billing', label: 'Billing', icon: CreditCard },
		{ id: 'notifications', label: 'Notifications', icon: Bell },
		{ id: 'security', label: 'Security', icon: Shield }
	];
</script>

<svelte:head>
	<title>Settings - Resume Builder</title>
</svelte:head>

<div class="container mx-auto px-6 py-8">
	<h1 class="mb-8 text-3xl font-bold">Settings</h1>

	<div class="flex gap-8">
		<!-- Sidebar -->
		<div class="w-48 space-y-1">
			{#each tabs as tab}
				<button
					class="flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm font-medium transition-colors {activeTab ===
					tab.id
						? 'bg-primary text-primary-foreground'
						: 'text-muted-foreground hover:bg-muted hover:text-foreground'}"
					on:click={() => (activeTab = tab.id)}
				>
					<svelte:component this={tab.icon} class="h-4 w-4" />
					{tab.label}
				</button>
			{/each}
		</div>

		<!-- Content -->
		<div class="flex-1">
			{#if activeTab === 'profile'}
				<Card class="p-6">
					<h2 class="mb-6 text-xl font-semibold">Profile Settings</h2>
					<form class="space-y-4">
						<div>
							<Label htmlFor="name">Full Name</Label>
							<Input id="name" value={$currentUser?.name || ''} />
						</div>
						<div>
							<Label htmlFor="email">Email Address</Label>
							<Input id="email" type="email" value={$currentUser?.email || ''} disabled />
							<p class="mt-1 text-xs text-muted-foreground">
								Contact support to change your email address
							</p>
						</div>
						<Button>Save Changes</Button>
					</form>
				</Card>
			{:else if activeTab === 'billing'}
				<Card class="p-6">
					<h2 class="mb-6 text-xl font-semibold">Billing & Subscription</h2>
					<div class="mb-6 rounded-lg bg-muted p-4">
						<div class="flex items-center justify-between">
							<div>
								<p class="font-medium">Current Plan</p>
								<p class="text-2xl font-bold capitalize">{$currentUser?.subscription_tier}</p>
							</div>
							{#if $currentUser?.subscription_tier === 'free'}
								<Button>Upgrade to Pro</Button>
							{/if}
						</div>
					</div>
					{#if $currentUser?.subscription_tier === 'free'}
						<div class="space-y-4">
							<h3 class="font-medium">Pro Features Include:</h3>
							<ul class="space-y-2 text-sm text-muted-foreground">
								<li>• Unlimited resumes</li>
								<li>• All premium templates</li>
								<li>• Unlimited AI improvements</li>
								<li>• PDF export without watermark</li>
								<li>• Priority support</li>
							</ul>
						</div>
					{/if}
				</Card>
			{:else if activeTab === 'notifications'}
				<Card class="p-6">
					<h2 class="mb-6 text-xl font-semibold">Notification Preferences</h2>
					<div class="space-y-4">
						<label class="flex items-center justify-between">
							<div>
								<p class="font-medium">Email Notifications</p>
								<p class="text-sm text-muted-foreground">Receive updates about your resumes</p>
							</div>
							<input type="checkbox" class="rounded border-input" checked />
						</label>
						<label class="flex items-center justify-between">
							<div>
								<p class="font-medium">Marketing Emails</p>
								<p class="text-sm text-muted-foreground">Tips and best practices</p>
							</div>
							<input type="checkbox" class="rounded border-input" />
						</label>
					</div>
				</Card>
			{:else if activeTab === 'security'}
				<Card class="p-6">
					<h2 class="mb-6 text-xl font-semibold">Security Settings</h2>
					<div class="space-y-6">
						<div>
							<h3 class="mb-2 font-medium">Change Password</h3>
							<form class="space-y-4">
								<div>
									<Label htmlFor="currentPassword">Current Password</Label>
									<Input id="currentPassword" type="password" />
								</div>
								<div>
									<Label htmlFor="newPassword">New Password</Label>
									<Input id="newPassword" type="password" />
								</div>
								<div>
									<Label htmlFor="confirmPassword">Confirm New Password</Label>
									<Input id="confirmPassword" type="password" />
								</div>
								<Button>Update Password</Button>
							</form>
						</div>
						<hr />
						<div>
							<h3 class="mb-2 font-medium text-destructive">Danger Zone</h3>
							<p class="mb-4 text-sm text-muted-foreground">
								Once you delete your account, there is no going back. Please be certain.
							</p>
							<Button variant="destructive">Delete Account</Button>
						</div>
					</div>
				</Card>
			{/if}
		</div>
	</div>
</div>
