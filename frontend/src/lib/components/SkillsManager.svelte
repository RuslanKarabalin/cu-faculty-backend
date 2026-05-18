<script lang="ts">
	import { onMount } from 'svelte';
	import { ApiClientError } from '$lib/api/client';
	import type { Skill } from '$lib/api/types';
	import Alert from './Alert.svelte';
	import Button from './Button.svelte';
	import Card from './Card.svelte';
	import Field from './Field.svelte';
	import Input from './Input.svelte';

	type Props = {
		title: string;
		api: {
			list: () => Promise<Skill[]>;
			add: (id: number) => Promise<Skill>;
			remove: (id: number) => Promise<void>;
		};
	};

	let { title, api }: Props = $props();

	let items = $state<Skill[]>([]);
	let error = $state<string | null>(null);
	let loading = $state(true);
	let newId = $state<number | null>(null);
	let submitting = $state(false);

	async function load() {
		loading = true;
		error = null;
		try {
			items = await api.list();
		} catch (e) {
			error = e instanceof ApiClientError ? e.message : String(e);
		} finally {
			loading = false;
		}
	}

	onMount(load);

	async function add(e: Event) {
		e.preventDefault();
		if (!newId) return;
		submitting = true;
		error = null;
		try {
			await api.add(newId);
			newId = null;
			await load();
		} catch (err) {
			error = err instanceof ApiClientError ? err.message : String(err);
		} finally {
			submitting = false;
		}
	}

	async function remove(id: number) {
		if (!confirm('Remove this skill?')) return;
		try {
			await api.remove(id);
			await load();
		} catch (err) {
			error = err instanceof ApiClientError ? err.message : String(err);
		}
	}
</script>

{#if error}<Alert variant="error">{error}</Alert>{/if}

<div class="grid gap-4 lg:grid-cols-3">
	<div class="lg:col-span-2">
		<Card {title}>
			{#if loading}
				<p class="text-sm text-zinc-500">Loading…</p>
			{:else if items.length === 0}
				<p class="text-sm text-zinc-500">No skills yet.</p>
			{:else}
				<div class="flex flex-wrap gap-2">
					{#each items as s (s.id)}
						<span class="inline-flex items-center gap-2 rounded-full bg-zinc-100 py-1 pl-3 pr-1 text-sm text-zinc-700">
							{s.name}
							<button
								type="button"
								onclick={() => remove(s.id)}
								aria-label="Remove {s.name}"
								class="flex h-5 w-5 items-center justify-center rounded-full text-zinc-400 hover:bg-zinc-200 hover:text-zinc-900"
							>
								×
							</button>
						</span>
					{/each}
				</div>
			{/if}
		</Card>
	</div>

	<Card title="Add skill">
		<form onsubmit={add} class="space-y-3">
			<Field label="Skill ID" hint="The backend has no listing endpoint yet — supply a numeric ID.">
				<Input type="number" required min="1" bind:value={newId} />
			</Field>
			<Button type="submit" disabled={submitting || !newId}>Add</Button>
		</form>
	</Card>
</div>
