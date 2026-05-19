<script lang="ts">
	import { onMount } from 'svelte';
	import { api, ApiClientError } from '$lib/api/client';
	import type { University } from '$lib/api/types';
	import Alert from '$lib/components/Alert.svelte';
	import ReferenceTable from '$lib/components/ReferenceTable.svelte';

	let items = $state<University[]>([]);
	let error = $state<string | null>(null);
	let loading = $state(true);

	async function load() {
		loading = true;
		error = null;
		try {
			items = await api.reference.universities();
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
	title="Universities"
	{items}
	{loading}
	matches={(u, q) =>
		u.name.toLowerCase().includes(q) || u.shortName.toLowerCase().includes(q)}
>
	{#snippet headers()}
		<th class="px-5 py-2 font-medium">ID</th>
		<th class="px-5 py-2 font-medium">Short name</th>
		<th class="px-5 py-2 font-medium">Full name</th>
	{/snippet}
	{#snippet row(u)}
		<td class="px-5 py-2 align-top font-mono text-xs text-zinc-500">{u.id}</td>
		<td class="whitespace-nowrap px-5 py-2 align-top font-medium text-zinc-900">{u.shortName}</td>
		<td class="px-5 py-2 align-top text-zinc-700">{u.name}</td>
	{/snippet}
</ReferenceTable>
