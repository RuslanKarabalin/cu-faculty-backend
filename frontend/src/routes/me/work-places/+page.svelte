<script lang="ts">
	import { onMount } from 'svelte';
	import { api, ApiClientError } from '$lib/api/client';
	import { WORK_GRADES, type WorkPlace, type WorkPlaceRequest } from '$lib/api/types';
	import Alert from '$lib/components/Alert.svelte';
	import Button from '$lib/components/Button.svelte';
	import Card from '$lib/components/Card.svelte';
	import Field from '$lib/components/Field.svelte';
	import Input from '$lib/components/Input.svelte';
	import Select from '$lib/components/Select.svelte';

	let items = $state<WorkPlace[]>([]);
	let error = $state<string | null>(null);
	let loading = $state(true);

	function emptyForm(): WorkPlaceRequest {
		return {
			companyName: '',
			grade: 'junior',
			position: '',
			startYear: new Date().getFullYear(),
			endYear: null,
			isWorkingNow: false
		};
	}

	let form = $state(emptyForm());
	let editingId = $state<number | null>(null);
	let submitting = $state(false);

	async function load() {
		loading = true;
		error = null;
		try {
			items = await api.me.workPlaces.list();
		} catch (e) {
			error = e instanceof ApiClientError ? e.message : String(e);
		} finally {
			loading = false;
		}
	}

	onMount(load);

	function startEdit(item: WorkPlace) {
		editingId = item.id;
		form = {
			companyName: item.companyName,
			grade: item.grade,
			position: item.position,
			startYear: item.startYear,
			endYear: item.endYear,
			isWorkingNow: item.isWorkingNow
		};
	}

	function cancelEdit() {
		editingId = null;
		form = emptyForm();
	}

	$effect(() => {
		if (form.isWorkingNow) form.endYear = null;
	});

	async function save(e: Event) {
		e.preventDefault();
		submitting = true;
		error = null;
		try {
			const body: WorkPlaceRequest = {
				...form,
				endYear: form.isWorkingNow ? null : form.endYear
			};
			if (editingId !== null) {
				await api.me.workPlaces.update(editingId, body);
			} else {
				await api.me.workPlaces.create(body);
			}
			cancelEdit();
			await load();
		} catch (err) {
			error = err instanceof ApiClientError ? err.message : String(err);
		} finally {
			submitting = false;
		}
	}

	async function remove(id: number) {
		if (!confirm('Delete this work place?')) return;
		try {
			await api.me.workPlaces.remove(id);
			await load();
		} catch (err) {
			error = err instanceof ApiClientError ? err.message : String(err);
		}
	}
</script>

{#if error}<Alert variant="error">{error}</Alert>{/if}

<div class="grid gap-4 lg:grid-cols-3">
	<div class="lg:col-span-2">
		<Card title="My work places">
			{#if loading}
				<p class="text-sm text-zinc-500">Loading…</p>
			{:else if items.length === 0}
				<p class="text-sm text-zinc-500">No work places yet.</p>
			{:else}
				<ul class="space-y-3">
					{#each items as w (w.id)}
						<li class="flex items-start justify-between gap-4 rounded-md border border-zinc-200 bg-zinc-50 p-3">
							<div class="text-sm">
								<div class="font-medium text-zinc-900">{w.companyName}</div>
								<div class="text-zinc-600">{w.position} · {w.grade}</div>
								<div class="text-xs text-zinc-500">
									{w.startYear}–{w.isWorkingNow ? 'now' : (w.endYear ?? '?')}
								</div>
							</div>
							<div class="flex gap-2">
								<Button variant="secondary" size="sm" onclick={() => startEdit(w)}>Edit</Button>
								<Button variant="danger" size="sm" onclick={() => remove(w.id)}>Delete</Button>
							</div>
						</li>
					{/each}
				</ul>
			{/if}
		</Card>
	</div>

	<Card title={editingId !== null ? 'Edit work place' : 'Add work place'}>
		<form onsubmit={save} class="space-y-3">
			<Field label="Company name">
				<Input required bind:value={form.companyName} />
			</Field>
			<Field label="Position">
				<Input required bind:value={form.position} />
			</Field>
			<Field label="Grade">
				<Select bind:value={form.grade}>
					{#each WORK_GRADES as g (g)}<option value={g}>{g}</option>{/each}
				</Select>
			</Field>
			<div class="grid grid-cols-2 gap-3">
				<Field label="Start year">
					<Input type="number" required bind:value={form.startYear} />
				</Field>
				<Field label="End year">
					<Input type="number" bind:value={form.endYear} disabled={form.isWorkingNow} />
				</Field>
			</div>
			<label class="flex items-center gap-2 text-sm text-zinc-700">
				<input type="checkbox" bind:checked={form.isWorkingNow} class="h-4 w-4 rounded border-zinc-300" />
				Currently working here
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
