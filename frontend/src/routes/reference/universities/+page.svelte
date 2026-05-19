<script lang="ts">
	import Alert from '$lib/components/Alert.svelte';
	import ReferenceTable from '$lib/components/ReferenceTable.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

{#if data.error}
	<Alert variant="error">{data.error}</Alert>
{:else}
	<ReferenceTable
		title="Universities"
		items={data.items}
		matches={(u, q) => u.name.toLowerCase().includes(q) || u.shortName.toLowerCase().includes(q)}
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
{/if}
