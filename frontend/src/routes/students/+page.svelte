<script lang="ts">
	import { goto } from '$app/navigation';
	import Alert from '$lib/components/Alert.svelte';
	import Button from '$lib/components/Button.svelte';
	import Card from '$lib/components/Card.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	function navigate(offset: number) {
		const params = new URLSearchParams();
		if (offset > 0) params.set('offset', String(offset));
		goto(`/students${params.size ? `?${params}` : ''}`);
	}
</script>

<div class="space-y-6">
	<div>
		<h1 class="text-2xl font-bold tracking-tight">Students</h1>
		<p class="mt-1 text-sm text-zinc-500">All registered students.</p>
	</div>

	{#if data.error}
		<Alert variant="error">{data.error}</Alert>
	{/if}

	{#if data.data}
		<Card>
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
						{#each data.data.data as u (u.id)}
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
							<tr>
								<td colspan="5" class="px-5 py-6 text-center text-zinc-500">No students yet.</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
			{#if data.data.total > 0}
				<div class="mt-4 flex items-center justify-between text-sm text-zinc-600">
					<span>
						{data.offset + 1}–{Math.min(data.offset + data.limit, data.data.total)} of
						{data.data.total}
					</span>
					<div class="flex gap-2">
						<Button
							variant="secondary"
							size="sm"
							onclick={() => navigate(Math.max(0, data.offset - data.limit))}
							disabled={data.offset === 0}
						>
							Previous
						</Button>
						<Button
							variant="secondary"
							size="sm"
							onclick={() => navigate(data.offset + data.limit)}
							disabled={data.offset + data.limit >= data.data.total}
						>
							Next
						</Button>
					</div>
				</div>
			{/if}
		</Card>
	{/if}
</div>
