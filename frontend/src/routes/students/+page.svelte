<script lang="ts">
	import { onMount } from 'svelte';
	import { api, ApiClientError } from '$lib/api/client';
	import type { Page, User } from '$lib/api/types';
	import Alert from '$lib/components/Alert.svelte';
	import Button from '$lib/components/Button.svelte';
	import Card from '$lib/components/Card.svelte';

	const LIMIT = 20;
	let offset = $state(0);
	let data = $state<Page<User> | null>(null);
	let error = $state<string | null>(null);
	let loading = $state(false);

	async function load() {
		loading = true;
		error = null;
		try {
			data = await api.students.list(LIMIT, offset);
		} catch (e) {
			error = e instanceof ApiClientError ? e.message : String(e);
		} finally {
			loading = false;
		}
	}

	onMount(load);

	$effect(() => {
		void offset;
		load();
	});

	function next() {
		if (data && offset + LIMIT < data.total) offset += LIMIT;
	}
	function prev() {
		if (offset > 0) offset = Math.max(0, offset - LIMIT);
	}
</script>

<div class="space-y-6">
	<div>
		<h1 class="text-2xl font-bold tracking-tight">Students</h1>
		<p class="mt-1 text-sm text-zinc-500">All registered students.</p>
	</div>

	{#if error}
		<Alert variant="error">{error}</Alert>
	{/if}

	<Card>
		{#if loading && !data}
			<p class="text-sm text-zinc-500">Loading…</p>
		{:else if data}
			<div class="-mx-5 overflow-x-auto">
				<table class="w-full text-left text-sm">
					<thead class="border-y border-zinc-200 bg-zinc-50 text-xs uppercase text-zinc-500">
						<tr>
							<th class="px-5 py-2 font-medium">Name</th>
							<th class="px-5 py-2 font-medium">Speciality</th>
							<th class="px-5 py-2 font-medium">Status</th>
							<th class="px-5 py-2 font-medium">Role</th>
							<th class="px-5 py-2"></th>
						</tr>
					</thead>
					<tbody class="divide-y divide-zinc-100">
						{#each data.data as u (u.id)}
							<tr class="hover:bg-zinc-50">
								<td class="px-5 py-3 font-medium text-zinc-900">
									{u.firstName} {u.lastName}
								</td>
								<td class="px-5 py-3 text-zinc-600">{u.speciality ?? '—'}</td>
								<td class="px-5 py-3 text-zinc-600">{u.status ?? '—'}</td>
								<td class="px-5 py-3 text-zinc-600">{u.role}</td>
								<td class="px-5 py-3 text-right">
									<a href="/students/{u.id}" class="text-sm font-medium text-zinc-900 underline">
										View
									</a>
								</td>
							</tr>
						{:else}
							<tr><td colspan="5" class="px-5 py-6 text-center text-zinc-500">No students yet.</td></tr>
						{/each}
					</tbody>
				</table>
			</div>
			<div class="mt-4 flex items-center justify-between text-sm text-zinc-600">
				<span>{Math.min(offset + 1, data.total)}–{Math.min(offset + LIMIT, data.total)} of {data.total}</span>
				<div class="flex gap-2">
					<Button variant="secondary" size="sm" onclick={prev} disabled={offset === 0}>Previous</Button>
					<Button variant="secondary" size="sm" onclick={next} disabled={offset + LIMIT >= data.total}>Next</Button>
				</div>
			</div>
		{/if}
	</Card>
</div>
