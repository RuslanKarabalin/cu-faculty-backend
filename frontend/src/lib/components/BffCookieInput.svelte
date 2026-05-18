<script lang="ts">
	import { onMount } from 'svelte';
	import { clearBffCookie, getBffCookie, setBffCookie } from '$lib/auth';
	import Button from './Button.svelte';
	import Input from './Input.svelte';

	let value = $state('');
	let currentLen = $state(0);

	function refresh() {
		const current = getBffCookie();
		currentLen = current?.length ?? 0;
	}

	onMount(refresh);

	function save() {
		if (!value.trim()) return;
		setBffCookie(value.trim());
		value = '';
		refresh();
	}

	function clear() {
		clearBffCookie();
		refresh();
	}
</script>

<div class="flex items-center gap-2">
	<span
		class="hidden text-xs font-medium md:inline {currentLen > 0 ? 'text-green-700' : 'text-zinc-500'}"
		title={currentLen > 0 ? `bff.cookie set (${currentLen} chars)` : 'bff.cookie is not set'}
	>
		● bff.cookie {currentLen > 0 ? `set · ${currentLen} chars` : 'not set'}
	</span>
	<Input
		type="password"
		placeholder="Paste bff.cookie value"
		bind:value
		class="w-56"
		autocomplete="off"
	/>
	<Button size="sm" onclick={save} disabled={!value.trim()}>Set</Button>
	{#if currentLen > 0}
		<Button size="sm" variant="secondary" onclick={clear}>Clear</Button>
	{/if}
</div>
