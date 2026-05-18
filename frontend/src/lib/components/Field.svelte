<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		label: string;
		hint?: string;
		children: Snippet;
	};

	let { label, hint, children }: Props = $props();
</script>

<label class="block space-y-1.5">
	<span class="text-sm font-medium text-zinc-700">{label}</span>
	{@render children()}
	{#if hint}<span class="text-xs text-zinc-500">{hint}</span>{/if}
</label>
