<script lang="ts">
	import { onMount } from 'svelte';
	import { api, ApiClientError } from '$lib/api/client';
	import type { Skill } from '$lib/api/types';
	import Alert from '$lib/components/Alert.svelte';
	import ReferenceTable from '$lib/components/ReferenceTable.svelte';

	let items = $state<Skill[]>([]);
	let error = $state<string | null>(null);
	let loading = $state(true);

	async function load() {
		loading = true;
		error = null;
		try {
			items = await api.reference.keySkills();
		} catch (e) {
			error = e instanceof ApiClientError ? e.message : String(e);
		} finally {
			loading = false;
		}
	}

	onMount(load);
</script>

{#if error}<Alert variant="error">{error}</Alert>{/if}

<ReferenceTable
	title="Key skills"
	{items}
	{loading}
	matches={(s, q) => s.name.toLowerCase().includes(q)}
>
	{#snippet headers()}
		<th class="px-5 py-2 font-medium">ID</th>
		<th class="px-5 py-2 font-medium">Name</th>
	{/snippet}
	{#snippet row(s)}
		<td class="px-5 py-2 font-mono text-xs text-zinc-500">{s.id}</td>
		<td class="px-5 py-2 text-zinc-900">{s.name}</td>
	{/snippet}
</ReferenceTable>
