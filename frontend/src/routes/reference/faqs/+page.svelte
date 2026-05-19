<script lang="ts">
	import { onMount } from 'svelte';
	import { api, ApiClientError } from '$lib/api/client';
	import type { Faq } from '$lib/api/types';
	import Alert from '$lib/components/Alert.svelte';
	import Card from '$lib/components/Card.svelte';
	import Input from '$lib/components/Input.svelte';

	let items = $state<Faq[]>([]);
	let error = $state<string | null>(null);
	let loading = $state(true);
	let query = $state('');

	async function load() {
		loading = true;
		error = null;
		try {
			items = await api.reference.faqs();
		} catch (e) {
			error = e instanceof ApiClientError ? e.message : String(e);
		} finally {
			loading = false;
		}
	}

	onMount(load);

	const filtered = $derived.by(() => {
		const q = query.trim().toLowerCase();
		if (!q) return items;
		return items.filter(
			(f) => f.question.toLowerCase().includes(q) || f.answer.toLowerCase().includes(q)
		);
	});
</script>

{#if error}<Alert variant="error">{error}</Alert>{/if}

<Card title="FAQs">
	{#snippet actions()}
		<div class="w-56"><Input bind:value={query} placeholder="Filter…" /></div>
	{/snippet}

	{#if loading}
		<p class="text-sm text-zinc-500">Loading…</p>
	{:else if items.length === 0}
		<p class="text-sm text-zinc-500">No FAQs.</p>
	{:else if filtered.length === 0}
		<p class="text-sm text-zinc-500">No matches.</p>
	{:else}
		<ul class="space-y-3">
			{#each filtered as f (f.id)}
				<li class="rounded-md border border-zinc-200 bg-zinc-50 p-3">
					<div class="text-sm font-medium text-zinc-900">{f.question}</div>
					<div class="mt-1 text-sm text-zinc-700">{f.answer}</div>
					<div class="mt-1 font-mono text-xs text-zinc-400">id {f.id}</div>
				</li>
			{/each}
		</ul>
		<div class="mt-3 text-xs text-zinc-500">{filtered.length} of {items.length}</div>
	{/if}
</Card>
