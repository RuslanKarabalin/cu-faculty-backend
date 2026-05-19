<script lang="ts">
	import { page } from '$app/state';

	let { children } = $props();

	const tabs = [
		{ href: '/reference', label: 'Overview' },
		{ href: '/reference/statuses', label: 'Statuses' },
		{ href: '/reference/key-skills', label: 'Key skills' },
		{ href: '/reference/soft-skills', label: 'Soft skills' },
		{ href: '/reference/companies', label: 'Companies' },
		{ href: '/reference/work-positions', label: 'Work positions' },
		{ href: '/reference/universities', label: 'Universities' },
		{ href: '/reference/faqs', label: 'FAQs' },
		{ href: '/reference/enums', label: 'Enums' }
	];

	function isActive(href: string) {
		return page.url.pathname === href;
	}
</script>

<div class="space-y-6">
	<div>
		<h1 class="text-2xl font-bold tracking-tight">Reference data</h1>
		<p class="mt-1 text-sm text-zinc-500">Catalogs and enums backing the profile pages.</p>
	</div>

	<nav class="-mx-2 flex flex-wrap gap-1 border-b border-zinc-200 pb-2">
		{#each tabs as { href, label } (href)}
			<a
				{href}
				class="rounded-md px-3 py-1.5 text-sm font-medium {isActive(href)
					? 'bg-zinc-900 text-white'
					: 'text-zinc-600 hover:bg-zinc-100 hover:text-zinc-900'}"
			>
				{label}
			</a>
		{/each}
	</nav>

	{@render children()}
</div>
