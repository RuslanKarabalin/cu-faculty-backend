<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { api, ApiError } from '$lib/api/client';
	import type { User, UserUpdateRequest } from '$lib/api/types';
	import Alert from '$lib/components/Alert.svelte';
	import Button from '$lib/components/Button.svelte';
	import Card from '$lib/components/Card.svelte';
	import Field from '$lib/components/Field.svelte';
	import Input from '$lib/components/Input.svelte';
	import Select from '$lib/components/Select.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	let error = $state<string | null>(null);
	let editing = $state(false);
	let submitting = $state(false);

	function emptyForm(): UserUpdateRequest {
		return { photoS3Key: null, bio: null, speciality: null, statusId: null };
	}

	let form = $state(emptyForm());
	let statusIdStr = $state('');

	function fillForm(u: User) {
		const match = data.statuses.find((s) => s.content === u.status);
		form = {
			photoS3Key: u.photoS3Key,
			bio: u.bio,
			speciality: u.speciality,
			statusId: match ? match.id : null
		};
		statusIdStr = match ? String(match.id) : '';
	}

	function startEdit() {
		if (!data.user) return;
		fillForm(data.user);
		editing = true;
	}

	function cancelEdit() {
		editing = false;
		form = emptyForm();
	}

	async function save(e: Event) {
		e.preventDefault();
		submitting = true;
		error = null;
		try {
			const body: UserUpdateRequest = {
				photoS3Key: form.photoS3Key?.trim() ? form.photoS3Key.trim() : null,
				bio: form.bio?.trim() ? form.bio.trim() : null,
				speciality: form.speciality?.trim() ? form.speciality.trim() : null,
				statusId: statusIdStr ? Number(statusIdStr) : null
			};
			await api.me.update(body);
			editing = false;
			await invalidateAll();
		} catch (err) {
			error = err instanceof ApiError ? err.message : String(err);
		} finally {
			submitting = false;
		}
	}

	function formatBirthdate(s: string): string {
		const d = new Date(s);
		return Number.isNaN(d.getTime()) ? s : d.toISOString().slice(0, 10);
	}

	const sections = [
		{ href: '/me/edu-places', title: 'Education', desc: 'Universities, grades, specializations.' },
		{ href: '/me/work-places', title: 'Work', desc: 'Companies, positions, work history.' },
		{ href: '/me/socials', title: 'Socials', desc: 'VK, Telegram, GitHub links.' },
		{ href: '/me/key-skills', title: 'Key skills', desc: 'Add/remove technical skills by ID.' },
		{ href: '/me/soft-skills', title: 'Soft skills', desc: 'Add/remove soft skills by ID.' }
	];
</script>

<div class="space-y-6">
	{#if data.error && !data.user}
		<Alert variant="error">{data.error}</Alert>
	{:else if data.user}
		{#if error}<Alert variant="error">{error}</Alert>{/if}

		{#if !editing}
			<Card title="Profile">
				{#snippet actions()}
					<Button size="sm" variant="secondary" onclick={startEdit}>Edit</Button>
				{/snippet}
				<div class="space-y-3 text-sm">
					<div>
						<div class="text-lg font-semibold text-zinc-900">
							{data.user.firstName}
							{data.user.lastName}
						</div>
						<div class="mt-0.5 text-xs text-zinc-500">
							<span class="font-mono">{data.user.id}</span> · {data.user.role}
						</div>
					</div>

					{#if data.user.bio}
						<p class="text-zinc-700">{data.user.bio}</p>
					{/if}

					<dl class="grid grid-cols-1 gap-x-6 gap-y-2 sm:grid-cols-2">
						<div>
							<dt class="text-xs font-medium uppercase text-zinc-500">Birth date</dt>
							<dd class="text-zinc-900">{formatBirthdate(data.user.birthdate)}</dd>
						</div>
						<div>
							<dt class="text-xs font-medium uppercase text-zinc-500">Speciality</dt>
							<dd class="text-zinc-900">{data.user.speciality ?? '-'}</dd>
						</div>
						<div>
							<dt class="text-xs font-medium uppercase text-zinc-500">Status</dt>
							<dd class="text-zinc-900">{data.user.status ?? '-'}</dd>
						</div>
						<div>
							<dt class="text-xs font-medium uppercase text-zinc-500">Photo</dt>
							<dd class="text-zinc-900">
								{data.user.photoS3Key ?? '-'}
							</dd>
						</div>
					</dl>
				</div>
			</Card>
		{:else}
			<Card title="Edit profile">
				<form onsubmit={save} class="space-y-4">
					<Field label="Bio" hint="Up to 255 characters.">
						<Input
							bind:value={form.bio}
							maxlength={255}
							placeholder="Tell something about yourself"
						/>
					</Field>

					<Field label="Speciality" hint="Up to 63 characters.">
						<Input
							bind:value={form.speciality}
							maxlength={63}
							placeholder="e.g. Backend Engineer"
						/>
					</Field>

					<Field label="Status">
						<Select bind:value={statusIdStr}>
							<option value="">-</option>
							{#each data.statuses as s (s.id)}
								<option value={String(s.id)}>{s.content}</option>
							{/each}
						</Select>
					</Field>

					<Field label="Photo S3 key" hint="Internal storage key.">
						<Input bind:value={form.photoS3Key} maxlength={255} placeholder="users/abc.jpg" />
					</Field>

					<div class="flex gap-2 pt-2">
						<Button type="submit" disabled={submitting}>
							{submitting ? 'Saving…' : 'Save'}
						</Button>
						<Button type="button" variant="secondary" onclick={cancelEdit}>Cancel</Button>
					</div>
				</form>
			</Card>
		{/if}
	{/if}

	<div class="grid gap-4 md:grid-cols-2">
		{#each sections as s (s.href)}
			<a href={s.href} class="block">
				<Card title={s.title}>
					<p class="text-sm text-zinc-600">{s.desc}</p>
				</Card>
			</a>
		{/each}
	</div>
</div>
