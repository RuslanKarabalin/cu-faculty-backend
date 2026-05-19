<script lang="ts">
	import Alert from '$lib/components/Alert.svelte';
	import ReferenceTable from '$lib/components/ReferenceTable.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

{#if data.error}<Alert variant="error">{data.error}</Alert>{/if}

<ReferenceTable
	title="Statuses"
	items={data.items}
	matches={(s, q) => s.content.toLowerCase().includes(q)}
>
	{#snippet headers()}
		<th class="px-5 py-2 font-medium">ID</th>
		<th class="px-5 py-2 font-medium">Content</th>
	{/snippet}
	{#snippet row(s)}
		<td class="px-5 py-2 font-mono text-xs text-zinc-500">{s.id}</td>
		<td class="px-5 py-2 text-zinc-900">{s.content}</td>
	{/snippet}
</ReferenceTable>
