<script lang="ts">
	import { onMount } from 'svelte';
	import { api, ApiClientError } from '$lib/api/client';
	import Alert from '$lib/components/Alert.svelte';
	import Card from '$lib/components/Card.svelte';

	type EnumBlock = { title: string; values: string[] };

	let blocks = $state<EnumBlock[]>([]);
	let error = $state<string | null>(null);
	let loading = $state(true);

	async function load() {
		loading = true;
		error = null;
		try {
			const [social, edu, work, event] = await Promise.all([
				api.reference.socialNetworks(),
				api.reference.eduGrades(),
				api.reference.workGrades(),
				api.reference.eventCategories()
			]);
			blocks = [
				{ title: 'Social networks', values: social },
				{ title: 'Edu grades', values: edu },
				{ title: 'Work grades', values: work },
				{ title: 'Event categories', values: event }
			];
		} catch (e) {
			error = e instanceof ApiClientError ? e.message : String(e);
		} finally {
			loading = false;
		}
	}

	onMount(load);
</script>

{#if error}<Alert variant="error">{error}</Alert>{/if}

{#if loading}
	<p class="text-sm text-zinc-500">Loading…</p>
{:else}
	<div class="grid gap-4 md:grid-cols-2">
		{#each blocks as b (b.title)}
			<Card title={b.title}>
				{#if b.values.length === 0}
					<p class="text-sm text-zinc-500">Empty.</p>
				{:else}
					<div class="flex flex-wrap gap-1.5">
						{#each b.values as v (v)}
							<span class="rounded-full bg-zinc-100 px-2.5 py-1 font-mono text-xs text-zinc-700">
								{v}
							</span>
						{/each}
					</div>
				{/if}
			</Card>
		{/each}
	</div>
{/if}
