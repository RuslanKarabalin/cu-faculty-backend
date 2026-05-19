<script lang="ts">
	import Card from '$lib/components/Card.svelte';

	const sections = [
		{ href: '/reference/statuses', title: 'Statuses', desc: 'User-facing status options.' },
		{ href: '/reference/key-skills', title: 'Key skills', desc: 'Technical skills catalog.' },
		{ href: '/reference/soft-skills', title: 'Soft skills', desc: 'Soft skills catalog.' },
		{ href: '/reference/companies', title: 'Companies', desc: 'Company suggestions for work history.' },
		{
			href: '/reference/work-positions',
			title: 'Work positions',
			desc: 'Common job titles for work history.'
		},
		{
			href: '/reference/universities',
			title: 'Universities',
			desc: 'University catalog used by education places.'
		},
		{ href: '/reference/faqs', title: 'FAQs', desc: 'Frequently asked questions.' },
		{
			href: '/reference/enums',
			title: 'Enums',
			desc: 'Social networks, edu grades, work grades, event categories.'
		}
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
