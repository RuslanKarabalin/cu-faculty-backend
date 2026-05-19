<script lang="ts">
	import Card from '$lib/components/Card.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<div class="space-y-6">
	<a href="/students" class="text-sm text-zinc-500 hover:text-zinc-900">← All students</a>

	<div>
		<h1 class="text-2xl font-bold tracking-tight">{data.user.firstName} {data.user.lastName}</h1>
		<p class="mt-1 text-sm text-zinc-500">
			{data.user.speciality ?? 'No speciality'} · {data.user.role}
		</p>
		{#if data.user.bio}<p class="mt-3 text-sm text-zinc-700">{data.user.bio}</p>{/if}
	</div>

	<div class="grid gap-4 lg:grid-cols-2">
		<Card title="Education places">
			{#if data.eduPlaces.length === 0}
				<p class="text-sm text-zinc-500">No education places.</p>
			{:else}
				<ul class="space-y-3 text-sm">
					{#each data.eduPlaces as e (e.id)}
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
			{#if data.workPlaces.length === 0}
				<p class="text-sm text-zinc-500">No work places.</p>
			{:else}
				<ul class="space-y-3 text-sm">
					{#each data.workPlaces as w (w.id)}
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
			{#if data.socials.length === 0}
				<p class="text-sm text-zinc-500">No socials.</p>
			{:else}
				<ul class="space-y-2 text-sm">
					{#each data.socials as s (s.id)}
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
				{#if data.keySkills.length === 0}
					<p class="mt-1 text-sm text-zinc-500">None.</p>
				{:else}
					<div class="mt-2 flex flex-wrap gap-1.5">
						{#each data.keySkills as s (s.id)}
							<span class="rounded-full bg-zinc-100 px-2.5 py-1 text-xs text-zinc-700">
								{s.name}
							</span>
						{/each}
					</div>
				{/if}
			</div>
			<div>
				<div class="text-xs font-medium uppercase text-zinc-500">Soft skills</div>
				{#if data.softSkills.length === 0}
					<p class="mt-1 text-sm text-zinc-500">None.</p>
				{:else}
					<div class="mt-2 flex flex-wrap gap-1.5">
						{#each data.softSkills as s (s.id)}
							<span class="rounded-full bg-zinc-100 px-2.5 py-1 text-xs text-zinc-700">
								{s.name}
							</span>
						{/each}
					</div>
				{/if}
			</div>
		</Card>
	</div>
</div>
