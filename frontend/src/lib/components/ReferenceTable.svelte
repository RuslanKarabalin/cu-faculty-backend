<script lang="ts" generics="T extends { id: number }">
	import type { Snippet } from 'svelte';
	import Card from './Card.svelte';
	import Input from './Input.svelte';

	type Props = {
		title: string;
		items: T[];
		matches: (item: T, q: string) => boolean;
		row: Snippet<[T]>;
		headers?: Snippet;
		placeholder?: string;
	};

	let { title, items, matches, row, headers, placeholder = 'Filter…' }: Props = $props();

	let query = $state('');

	const filtered = $derived(
		query.trim() ? items.filter((it) => matches(it, query.trim().toLowerCase())) : items
	);
</script>

<Card {title}>
	{#snippet actions()}
		<div class="w-56">
			<Input bind:value={query} {placeholder} />
		</div>
	{/snippet}

	<div class="-mx-5 overflow-x-auto">
		<table class="w-full text-left text-sm">
			{#if headers}
				<thead class="border-y border-zinc-200 bg-zinc-50 text-xs uppercase text-zinc-500">
					<tr>{@render headers()}</tr>
				</thead>
			{/if}
			<tbody class="divide-y divide-zinc-100">
				{#each filtered as item (item.id)}
					<tr class="hover:bg-zinc-50">{@render row(item)}</tr>
				{:else}
					<tr>
						<td class="px-5 py-6 text-center text-zinc-500" colspan="99">
							{items.length === 0 ? 'No items.' : 'No matches.'}
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
	<div class="mt-3 text-xs text-zinc-500">
		{filtered.length}
		of {items.length}
	</div>
</Card>
