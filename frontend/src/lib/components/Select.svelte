<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLSelectAttributes } from 'svelte/elements';

	type Props = HTMLSelectAttributes & {
		value?: string;
		children: Snippet;
	};

	let { class: klass = '', value = $bindable(), children, ...rest }: Props = $props();
</script>

<select
	bind:value
	{...rest}
	class="h-10 w-full rounded-md border border-zinc-300 bg-white px-3 text-sm text-zinc-900 focus:border-zinc-900 focus:outline-none focus:ring-1 focus:ring-zinc-900 {klass}"
>
	{@render children()}
</select>
