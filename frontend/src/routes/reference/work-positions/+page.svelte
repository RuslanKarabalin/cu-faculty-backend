<script lang="ts">
	import Alert from '$lib/components/Alert.svelte';
	import ReferenceTable from '$lib/components/ReferenceTable.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

{#if data.error}<Alert variant="error">{data.error}</Alert>{/if}

<ReferenceTable
	title="Work positions"
	items={data.items}
	matches={(p, q) => p.name.toLowerCase().includes(q)}
>
	{#snippet headers()}
		<th class="px-5 py-2 font-medium">ID</th>
		<th class="px-5 py-2 font-medium">Name</th>
	{/snippet}
	{#snippet row(p)}
		<td class="px-5 py-2 font-mono text-xs text-zinc-500">{p.id}</td>
		<td class="px-5 py-2 text-zinc-900">{p.name}</td>
	{/snippet}
</ReferenceTable>
