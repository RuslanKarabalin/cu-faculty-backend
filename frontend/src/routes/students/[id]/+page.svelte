<script lang="ts">
	import { page } from '$app/state';
	import { onMount } from 'svelte';
	import { api, ApiClientError } from '$lib/api/client';
	import type { EduPlace, Skill, Social, User, WorkPlace } from '$lib/api/types';
	import Alert from '$lib/components/Alert.svelte';
	import Card from '$lib/components/Card.svelte';

	const id = $derived(page.params.id ?? '');

	let user = $state<User | null>(null);
	let eduPlaces = $state<EduPlace[]>([]);
	let workPlaces = $state<WorkPlace[]>([]);
	let socials = $state<Social[]>([]);
	let keySkills = $state<Skill[]>([]);
	let softSkills = $state<Skill[]>([]);
	let error = $state<string | null>(null);
	let loading = $state(true);

	async function load() {
		loading = true;
		error = null;
		try {
			// No GET /students/:id in backend — find from the list.
			// Hit the sub-resources directly; they 200 with [] for missing users.
			const [edu, work, soc, ks, ss, list] = await Promise.all([
				api.students.eduPlaces(id),
				api.students.workPlaces(id),
				api.students.socials(id),
				api.students.keySkills(id),
				api.students.softSkills(id),
				api.students.list(100, 0)
			]);
			eduPlaces = edu;
			workPlaces = work;
			socials = soc;
			keySkills = ks;
			softSkills = ss;
			user = list.data.find((u) => u.id === id) ?? null;
		} catch (e) {
			error = e instanceof ApiClientError ? e.message : String(e);
		} finally {
			loading = false;
		}
	}

	onMount(load);
</script>

<div class="space-y-6">
	<a href="/students" class="text-sm text-zinc-500 hover:text-zinc-900">← All students</a>

	{#if loading}
		<p class="text-sm text-zinc-500">Loading…</p>
	{:else if error}
		<Alert variant="error">{error}</Alert>
	{:else}
		<div>
			{#if user}
				<h1 class="text-2xl font-bold tracking-tight">{user.firstName} {user.lastName}</h1>
				<p class="mt-1 text-sm text-zinc-500">
					{user.speciality ?? 'No speciality'} · {user.role}
				</p>
				{#if user.bio}<p class="mt-3 text-sm text-zinc-700">{user.bio}</p>{/if}
			{:else}
				<h1 class="text-2xl font-bold tracking-tight">Student</h1>
				<p class="mt-1 font-mono text-xs text-zinc-500">{id}</p>
			{/if}
		</div>

		<div class="grid gap-4 lg:grid-cols-2">
			<Card title="Education places">
				{#if eduPlaces.length === 0}
					<p class="text-sm text-zinc-500">No education places.</p>
				{:else}
					<ul class="space-y-3 text-sm">
						{#each eduPlaces as e (e.id)}
							<li class="rounded-md border border-zinc-200 bg-zinc-50 p-3">
								<div class="font-medium text-zinc-900">{e.universityName}</div>
								<div class="text-zinc-600">
									{e.specialization} · {e.grade}{e.level ? ` · ${e.level}` : ''}
								</div>
								<div class="text-xs text-zinc-500">
									{e.startYear}–{e.isStudyingNow ? 'now' : (e.endYear ?? '?')}
								</div>
							</li>
						{/each}
					</ul>
				{/if}
			</Card>

			<Card title="Work places">
				{#if workPlaces.length === 0}
					<p class="text-sm text-zinc-500">No work places.</p>
				{:else}
					<ul class="space-y-3 text-sm">
						{#each workPlaces as w (w.id)}
							<li class="rounded-md border border-zinc-200 bg-zinc-50 p-3">
								<div class="font-medium text-zinc-900">{w.companyName}</div>
								<div class="text-zinc-600">{w.position} · {w.grade}</div>
								<div class="text-xs text-zinc-500">
									{w.startYear}–{w.isWorkingNow ? 'now' : (w.endYear ?? '?')}
								</div>
							</li>
						{/each}
					</ul>
				{/if}
			</Card>

			<Card title="Socials">
				{#if socials.length === 0}
					<p class="text-sm text-zinc-500">No socials.</p>
				{:else}
					<ul class="space-y-2 text-sm">
						{#each socials as s (s.id)}
							<li class="flex items-center justify-between">
								<div>
									<span class="font-medium text-zinc-900">{s.social}</span>
									<span class="ml-2 text-zinc-600">{s.link}</span>
								</div>
								{#if s.isPreferred}
									<span class="rounded-full bg-zinc-900 px-2 py-0.5 text-xs text-white">
										preferred
									</span>
								{/if}
							</li>
						{/each}
					</ul>
				{/if}
			</Card>

			<Card title="Skills">
				<div class="mb-3">
					<div class="text-xs font-medium uppercase text-zinc-500">Key skills</div>
					{#if keySkills.length === 0}
						<p class="mt-1 text-sm text-zinc-500">None.</p>
					{:else}
						<div class="mt-2 flex flex-wrap gap-1.5">
							{#each keySkills as s (s.id)}
								<span class="rounded-full bg-zinc-100 px-2.5 py-1 text-xs text-zinc-700">
									{s.name}
								</span>
							{/each}
						</div>
					{/if}
				</div>
				<div>
					<div class="text-xs font-medium uppercase text-zinc-500">Soft skills</div>
					{#if softSkills.length === 0}
						<p class="mt-1 text-sm text-zinc-500">None.</p>
					{:else}
						<div class="mt-2 flex flex-wrap gap-1.5">
							{#each softSkills as s (s.id)}
								<span class="rounded-full bg-zinc-100 px-2.5 py-1 text-xs text-zinc-700">
									{s.name}
								</span>
							{/each}
						</div>
					{/if}
				</div>
			</Card>
		</div>
	{/if}
</div>
