<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { api, ApiError } from '$lib/api/client';
	import { EDU_GRADES, type EduPlace, type EduPlaceRequest } from '$lib/api/types';
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

	type FormState = Omit<EduPlaceRequest, 'universityId'> & { universityIdStr: string };

	function emptyForm(): FormState {
		return {
			universityIdStr: '',
			grade: 'bachelor',
			level: null,
			specialization: '',
			startYear: new Date().getFullYear(),
			endYear: null,
			isStudyingNow: false
		};
	}

	let form = $state<FormState>(emptyForm());
	let editingId = $state<number | null>(null);
	let submitting = $state(false);

	function startEdit(item: EduPlace) {
		editingId = item.id;
		form = {
			universityIdStr: String(item.universityId),
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

	$effect(() => {
		if (form.isStudyingNow) form.endYear = null;
	});

	async function save(e: Event) {
		e.preventDefault();
		const universityId = Number(form.universityIdStr);
		if (!universityId) return;
		submitting = true;
		mutationError = null;
		try {
			const body: EduPlaceRequest = {
				universityId,
				grade: form.grade,
				level: form.level?.trim() ? form.level.trim() : null,
				specialization: form.specialization.trim(),
				startYear: form.startYear,
				endYear: form.isStudyingNow ? null : form.endYear,
				isStudyingNow: form.isStudyingNow
			};
			if (editingId !== null) {
				await api.me.eduPlaces.update(editingId, body);
			} else {
				await api.me.eduPlaces.create(body);
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
		if (!confirm('Delete this education place?')) return;
		try {
			await api.me.eduPlaces.remove(id);
			await invalidateAll();
		} catch (err) {
			mutationError = err instanceof ApiError ? err.message : String(err);
		}
	}
</script>

{#if error}<Alert variant="error">{error}</Alert>{/if}

<div class="grid gap-4 lg:grid-cols-3">
	<div class="lg:col-span-2">
		<Card title="My education places">
			{#if data.items.length === 0}
				<p class="text-sm text-zinc-500">No education places yet.</p>
			{:else}
				<ul class="space-y-3">
					{#each data.items as e (e.id)}
						<li
							class="flex items-start justify-between gap-4 rounded-md border border-zinc-200 bg-zinc-50 p-3"
						>
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
			<Field label="University">
				<Select bind:value={form.universityIdStr} required>
					<option value="">Select a university…</option>
					{#each data.universities as u (u.id)}
						<option value={String(u.id)}>{u.shortName} — {u.name}</option>
					{/each}
				</Select>
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
				<input
					type="checkbox"
					bind:checked={form.isStudyingNow}
					class="h-4 w-4 rounded border-zinc-300"
				/>
				Currently studying
			</label>
			<div class="flex gap-2 pt-2">
				<Button type="submit" disabled={submitting || !form.universityIdStr}>
					{editingId !== null ? 'Save' : 'Add'}
				</Button>
				{#if editingId !== null}
					<Button type="button" variant="secondary" onclick={cancelEdit}>Cancel</Button>
				{/if}
			</div>
		</form>
	</Card>
</div>
