<script lang="ts">
	import Alert from '$lib/components/Alert.svelte';
	import Card from '$lib/components/Card.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

{#if data.error}<Alert variant="error">{data.error}</Alert>{/if}

<div class="grid gap-4 md:grid-cols-2">
	{#each data.blocks as b (b.title)}
		<Card title={b.title}>
			{#if b.values.length === 0}
				<p class="text-sm text-zinc-500">Empty.</p>
			{:else}
				<div class="flex flex-wrap gap-1.5">
					{#each b.values as v (v)}
						<span
							class="rounded-full bg-zinc-100 px-2.5 py-1 font-mono text-xs text-zinc-700"
						>
							{v}
						</span>
					{/each}
				</div>
			{/if}
		</Card>
	{/each}
</div>
