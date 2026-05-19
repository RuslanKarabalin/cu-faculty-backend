<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { ApiError } from '$lib/api/client';
	import type { Skill } from '$lib/api/types';
	import Alert from './Alert.svelte';
	import Button from './Button.svelte';
	import Card from './Card.svelte';
	import Field from './Field.svelte';
	import Select from './Select.svelte';

	type Props = {
		title: string;
		items: Skill[];
		all: Skill[];
		api: {
			add: (id: number) => Promise<Skill>;
			remove: (id: number) => Promise<void>;
		};
	};

	let { title, items, all, api }: Props = $props();

	let error = $state<string | null>(null);
	let selectedId = $state('');
	let submitting = $state(false);

	const available = $derived.by(() => {
		const taken = new Set(items.map((s) => s.id));
		return all.filter((s) => !taken.has(s.id));
	});

	async function add(e: Event) {
		e.preventDefault();
		const id = Number(selectedId);
		if (!id) return;
		submitting = true;
		error = null;
		try {
			await api.add(id);
			selectedId = '';
			await invalidateAll();
		} catch (err) {
			error = err instanceof ApiError ? err.message : String(err);
		} finally {
			submitting = false;
		}
	}

	async function remove(id: number) {
		if (!confirm('Remove this skill?')) return;
		try {
			await api.remove(id);
			await invalidateAll();
		} catch (err) {
			error = err instanceof ApiError ? err.message : String(err);
		}
	}
</script>

{#if error}<Alert variant="error">{error}</Alert>{/if}

<div class="grid gap-4 lg:grid-cols-3">
	<div class="lg:col-span-2">
		<Card {title}>
			{#if items.length === 0}
				<p class="text-sm text-zinc-500">No skills yet.</p>
			{:else}
				<div class="flex flex-wrap gap-2">
					{#each items as s (s.id)}
						<span
							class="inline-flex items-center gap-2 rounded-full bg-zinc-100 py-1 pl-3 pr-1 text-sm text-zinc-700"
						>
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
			<Field label="Skill">
				<Select bind:value={selectedId} disabled={available.length === 0}>
					<option value="">
						{available.length === 0 ? 'No more skills to add' : 'Select a skill…'}
					</option>
					{#each available as s (s.id)}
						<option value={String(s.id)}>{s.name}</option>
					{/each}
				</Select>
			</Field>
			<Button type="submit" disabled={submitting || !selectedId}>Add</Button>
		</form>
	</Card>
</div>
