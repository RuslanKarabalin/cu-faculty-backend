<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		title?: string;
		actions?: Snippet;
		children: Snippet;
		class?: string;
	};

	let { title, actions, children, class: klass = '' }: Props = $props();
</script>

<section class="rounded-lg border border-zinc-200 bg-white shadow-sm {klass}">
	{#if title || actions}
		<header
			class="flex items-center justify-between gap-3 border-b border-zinc-200 px-5 py-4"
		>
			{#if title}<h2 class="text-base font-semibold text-zinc-900">{title}</h2>{/if}
			{#if actions}<div class="flex items-center gap-2">{@render actions()}</div>{/if}
		</header>
	{/if}
	<div class="p-5">{@render children()}</div>
</section>
