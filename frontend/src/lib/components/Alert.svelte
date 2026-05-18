<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		variant?: 'error' | 'info';
		children: Snippet;
	};

	let { variant = 'info', children }: Props = $props();

	const styles = $derived(
		variant === 'error'
			? 'border-red-200 bg-red-50 text-red-800'
			: 'border-blue-200 bg-blue-50 text-blue-800'
	);
</script>

<div class="rounded-md border px-4 py-3 text-sm {styles}">
	{@render children()}
</div>
