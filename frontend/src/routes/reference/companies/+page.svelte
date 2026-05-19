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
		title="Companies"
		items={data.items}
		matches={(c, q) => c.name.toLowerCase().includes(q)}
	>
		{#snippet headers()}
			<th class="px-5 py-2 font-medium">ID</th>
			<th class="px-5 py-2 font-medium">Name</th>
		{/snippet}
		{#snippet row(c)}
			<td class="px-5 py-2 font-mono text-xs text-zinc-500">{c.id}</td>
			<td class="px-5 py-2 text-zinc-900">{c.name}</td>
		{/snippet}
	</ReferenceTable>
{/if}
