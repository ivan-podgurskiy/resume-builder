<script lang="ts">
	import { cn } from '$utils';

	export let value = '';
	export let placeholder = '';
	export let disabled = false;
	export let id = '';
	export let name = '';
	export let type: 'text' | 'email' | 'password' | 'number' | 'tel' | 'url' | 'search' | 'month' = 'text';

	let className = '';
	export { className as class };

	let inputElement: HTMLInputElement;

	function handleInput(event: Event) {
		const target = event.target as HTMLInputElement;
		value = target.value;
	}
</script>

{#if type === 'password'}
	<input
		type="password"
		{id}
		{name}
		{placeholder}
		{disabled}
		{value}
		on:input={handleInput}
		on:input
		on:change
		on:blur
		on:focus
		class={cn(
			'flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50',
			className
		)}
		{...$$restProps}
	/>
{:else if type === 'email'}
	<input
		type="email"
		{id}
		{name}
		{placeholder}
		{disabled}
		{value}
		on:input={handleInput}
		on:input
		on:change
		on:blur
		on:focus
		class={cn(
			'flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50',
			className
		)}
		{...$$restProps}
	/>
{:else if type === 'month'}
	<input
		type="month"
		{id}
		{name}
		{placeholder}
		{disabled}
		{value}
		on:input={handleInput}
		on:input
		on:change
		on:blur
		on:focus
		class={cn(
			'flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50',
			className
		)}
		{...$$restProps}
	/>
{:else}
	<input
		type="text"
		{id}
		{name}
		{placeholder}
		{disabled}
		{value}
		on:input={handleInput}
		on:input
		on:change
		on:blur
		on:focus
		class={cn(
			'flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50',
			className
		)}
		{...$$restProps}
	/>
{/if}
