<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { api, ApiError } from '$lib/api/client';
	import { SOCIAL_NETWORKS, type Social, type SocialRequest } from '$lib/api/types';
	import Alert from '$lib/components/Alert.svelte';
	import Button from '$lib/components/Button.svelte';
	import Card from '$lib/components/Card.svelte';
	import Field from '$lib/components/Field.svelte';
	import Input from '$lib/components/Input.svelte';
	import Select from '$lib/components/Select.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	let mutationError = $state<string | null>(null);
	const error = $derived(mutationError ?? data.error);

	function emptyForm(): SocialRequest {
		return { social: 'telegram', link: '', isPreferred: false };
	}

	let form = $state(emptyForm());
	let editingId = $state<number | null>(null);
	let submitting = $state(false);

	function startEdit(item: Social) {
		editingId = item.id;
		form = { social: item.social, link: item.link, isPreferred: item.isPreferred };
	}

	function cancelEdit() {
		editingId = null;
		form = emptyForm();
	}

	async function save(e: Event) {
		e.preventDefault();
		submitting = true;
		mutationError = null;
		try {
			if (editingId !== null) {
				await api.me.socials.update(editingId, form);
			} else {
				await api.me.socials.create(form);
			}
			cancelEdit();
			await invalidateAll();
		} catch (err) {
			mutationError = err instanceof ApiError ? err.message : String(err);
		} finally {
			submitting = false;
		}
	}

	async function remove(id: number) {
		if (!confirm('Delete this social?')) return;
		try {
			await api.me.socials.remove(id);
			await invalidateAll();
		} catch (err) {
			mutationError = err instanceof ApiError ? err.message : String(err);
		}
	}
</script>

{#if error}<Alert variant="error">{error}</Alert>{/if}

<div class="grid gap-4 lg:grid-cols-3">
	<div class="lg:col-span-2">
		<Card title="My socials">
			{#if data.items.length === 0}
				<p class="text-sm text-zinc-500">No socials yet.</p>
			{:else}
				<ul class="space-y-2">
					{#each data.items as s (s.id)}
						<li
							class="flex items-center justify-between gap-4 rounded-md border border-zinc-200 bg-zinc-50 px-3 py-2 text-sm"
						>
							<div>
								<span class="font-medium text-zinc-900">{s.social}</span>
								<span class="ml-2 text-zinc-600">{s.link}</span>
								{#if s.isPreferred}
									<span class="ml-2 rounded-full bg-zinc-900 px-2 py-0.5 text-xs text-white">
										preferred
									</span>
								{/if}
							</div>
							<div class="flex gap-2">
								<Button variant="secondary" size="sm" onclick={() => startEdit(s)}>Edit</Button>
								<Button variant="danger" size="sm" onclick={() => remove(s.id)}>Delete</Button>
							</div>
						</li>
					{/each}
				</ul>
			{/if}
		</Card>
	</div>

	<Card title={editingId !== null ? 'Edit social' : 'Add social'}>
		<form onsubmit={save} class="space-y-3">
			<Field label="Network">
				<Select bind:value={form.social}>
					{#each SOCIAL_NETWORKS as n (n)}<option value={n}>{n}</option>{/each}
				</Select>
			</Field>
			<Field label="Link or username">
				<Input required bind:value={form.link} maxlength={127} />
			</Field>
			<label class="flex items-center gap-2 text-sm text-zinc-700">
				<input
					type="checkbox"
					bind:checked={form.isPreferred}
					class="h-4 w-4 rounded border-zinc-300"
				/>
				Preferred social
			</label>
			<div class="flex gap-2 pt-2">
				<Button type="submit" disabled={submitting}>
					{editingId !== null ? 'Save' : 'Add'}
				</Button>
				{#if editingId !== null}
					<Button type="button" variant="secondary" onclick={cancelEdit}>Cancel</Button>
				{/if}
			</div>
		</form>
	</Card>
</div>
