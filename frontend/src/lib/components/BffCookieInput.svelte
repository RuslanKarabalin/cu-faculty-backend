<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { clearBffCookie, setBffCookie } from '$lib/auth';
	import Button from './Button.svelte';
	import Input from './Input.svelte';

	type Props = { authed: boolean };
	let { authed }: Props = $props();

	let value = $state('');
	let busy = $state(false);

	async function save() {
		const v = value.trim();
		if (!v) return;
		busy = true;
		try {
			await setBffCookie(v);
			value = '';
			await invalidateAll();
		} finally {
			busy = false;
		}
	}

	async function clear() {
		busy = true;
		try {
			await clearBffCookie();
			await invalidateAll();
		} finally {
			busy = false;
		}
	}
</script>

<div class="flex items-center gap-2">
	<span
		class="hidden text-xs font-medium md:inline {authed ? 'text-green-700' : 'text-zinc-500'}"
		title={authed ? 'bff.cookie set' : 'bff.cookie is not set'}
	>
		● bff.cookie {authed ? 'set' : 'not set'}
	</span>
	<Input
		type="password"
		placeholder="Paste bff.cookie value"
		bind:value
		class="w-56"
		autocomplete="off"
	/>
	<Button size="sm" onclick={save} disabled={busy || !value.trim()}>Set</Button>
	{#if authed}
		<Button size="sm" variant="secondary" onclick={clear} disabled={busy}>Clear</Button>
	{/if}
</div>
