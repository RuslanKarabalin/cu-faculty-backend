<script lang="ts">
	import { onMount } from 'svelte';
	import { api, ApiClientError } from '$lib/api/client';
	import { EDU_GRADES, type EduPlace, type EduPlaceRequest } from '$lib/api/types';
	import Alert from '$lib/components/Alert.svelte';
	import Button from '$lib/components/Button.svelte';
	import Card from '$lib/components/Card.svelte';
	import Field from '$lib/components/Field.svelte';
	import Input from '$lib/components/Input.svelte';
	import Select from '$lib/components/Select.svelte';

	let items = $state<EduPlace[]>([]);
	let error = $state<string | null>(null);
	let loading = $state(true);

	function emptyForm(): EduPlaceRequest {
		return {
			universityId: 0,
			grade: 'bachelor',
			level: null,
			specialization: '',
			startYear: new Date().getFullYear(),
			endYear: null,
			isStudyingNow: false
		};
	}

	let form = $state(emptyForm());
	let editingId = $state<number | null>(null);
	let submitting = $state(false);

	async function load() {
		loading = true;
		error = null;
		try {
			items = await api.me.eduPlaces.list();
		} catch (e) {
			error = e instanceof ApiClientError ? e.message : String(e);
		} finally {
			loading = false;
		}
	}

	onMount(load);

	function startEdit(item: EduPlace) {
		editingId = item.id;
		form = {
			// API doesn't return universityId — user must supply on edit.
			universityId: form.universityId || 0,
			grade: item.grade,
			level: item.level,
			specialization: item.specialization,
			startYear: item.startYear,
			endYear: item.endYear,
			isStudyingNow: item.isStudyingNow
		};
	}

	function cancelEdit() {
		editingId = null;
		form = emptyForm();
	}

	async function save(e: Event) {
		e.preventDefault();
		submitting = true;
		error = null;
		try {
			if (editingId !== null) {
				await api.me.eduPlaces.update(editingId, form);
			} else {
				await api.me.eduPlaces.create(form);
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
		if (!confirm('Delete this education place?')) return;
		try {
			await api.me.eduPlaces.remove(id);
			await load();
		} catch (err) {
			error = err instanceof ApiClientError ? err.message : String(err);
		}
	}
</script>

{#if error}<Alert variant="error">{error}</Alert>{/if}

<div class="grid gap-4 lg:grid-cols-3">
	<div class="lg:col-span-2">
		<Card title="My education places">
			{#if loading}
				<p class="text-sm text-zinc-500">Loading…</p>
			{:else if items.length === 0}
				<p class="text-sm text-zinc-500">No education places yet.</p>
			{:else}
				<ul class="space-y-3">
					{#each items as e (e.id)}
						<li class="flex items-start justify-between gap-4 rounded-md border border-zinc-200 bg-zinc-50 p-3">
							<div class="text-sm">
								<div class="font-medium text-zinc-900">{e.universityName}</div>
								<div class="text-zinc-600">
									{e.specialization} · {e.grade}{e.level ? ` · ${e.level}` : ''}
								</div>
								<div class="text-xs text-zinc-500">
									{e.startYear}–{e.isStudyingNow ? 'now' : (e.endYear ?? '?')}
								</div>
							</div>
							<div class="flex gap-2">
								<Button variant="secondary" size="sm" onclick={() => startEdit(e)}>Edit</Button>
								<Button variant="danger" size="sm" onclick={() => remove(e.id)}>Delete</Button>
							</div>
						</li>
					{/each}
				</ul>
			{/if}
		</Card>
	</div>

	<Card title={editingId !== null ? 'Edit education place' : 'Add education place'}>
		<form onsubmit={save} class="space-y-3">
			<Field label="University ID" hint="Numeric ID from universities catalog">
				<Input
					type="number"
					required
					min="1"
					bind:value={form.universityId}
				/>
			</Field>
			<Field label="Grade">
				<Select bind:value={form.grade}>
					{#each EDU_GRADES as g (g)}<option value={g}>{g}</option>{/each}
				</Select>
			</Field>
			<Field label="Level (optional)">
				<Input bind:value={form.level} placeholder="e.g. 3rd year" />
			</Field>
			<Field label="Specialization">
				<Input required bind:value={form.specialization} />
			</Field>
			<div class="grid grid-cols-2 gap-3">
				<Field label="Start year">
					<Input type="number" required bind:value={form.startYear} />
				</Field>
				<Field label="End year">
					<Input type="number" bind:value={form.endYear} disabled={form.isStudyingNow} />
				</Field>
			</div>
			<label class="flex items-center gap-2 text-sm text-zinc-700">
				<input type="checkbox" bind:checked={form.isStudyingNow} class="h-4 w-4 rounded border-zinc-300" />
				Currently studying
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
