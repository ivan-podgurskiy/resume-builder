<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { auth } from '$stores/auth';

	onMount(() => {
		auth.initialize();
	});
</script>

<slot />
