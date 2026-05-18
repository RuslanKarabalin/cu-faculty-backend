<script lang="ts">
	import Card from '$lib/components/Card.svelte';

	const sections = [
		{ href: '/me/edu-places', title: 'Education', desc: 'Universities, grades, specializations.' },
		{ href: '/me/work-places', title: 'Work', desc: 'Companies, positions, work history.' },
		{ href: '/me/socials', title: 'Socials', desc: 'VK, Telegram, GitHub links.' },
		{ href: '/me/key-skills', title: 'Key skills', desc: 'Add/remove technical skills by ID.' },
		{ href: '/me/soft-skills', title: 'Soft skills', desc: 'Add/remove soft skills by ID.' }
	];
</script>

<div class="grid gap-4 md:grid-cols-2">
	{#each sections as s (s.href)}
		<a href={s.href} class="block">
			<Card title={s.title}>
				<p class="text-sm text-zinc-600">{s.desc}</p>
			</Card>
		</a>
	{/each}
</div>
