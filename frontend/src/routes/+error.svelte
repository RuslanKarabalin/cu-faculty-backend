<script lang="ts">
	import { page } from '$app/state';
	import Alert from '$lib/components/Alert.svelte';
	import Button from '$lib/components/Button.svelte';
	import Card from '$lib/components/Card.svelte';
</script>

<div class="space-y-6">
	<div>
		<h1 class="text-2xl font-bold tracking-tight">
			{page.status} · Something went wrong
		</h1>
		<p class="mt-1 text-sm text-zinc-500">
			An unexpected error occurred while rendering this page.
		</p>
	</div>

	{#if page.error?.message}
		<Alert variant="error">{page.error.message}</Alert>
	{/if}

	<Card>
		<div class="flex flex-wrap gap-2">
			<Button onclick={() => location.reload()}>Reload</Button>
			<Button variant="secondary" onclick={() => history.back()}>Go back</Button>
			<a href="/" class="inline-flex items-center text-sm text-zinc-600 underline">Dashboard</a>
		</div>
	</Card>
</div>
