<script lang="ts">
	import { page } from '$app/state';

	let { children } = $props();

	const tabs = [
		{ href: '/me', label: 'Overview' },
		{ href: '/me/edu-places', label: 'Education' },
		{ href: '/me/work-places', label: 'Work' },
		{ href: '/me/socials', label: 'Socials' },
		{ href: '/me/key-skills', label: 'Key skills' },
		{ href: '/me/soft-skills', label: 'Soft skills' }
	];

	function isActive(href: string) {
		return page.url.pathname === href;
	}
</script>

<div class="space-y-6">
	<div>
		<h1 class="text-2xl font-bold tracking-tight">My profile</h1>
		<p class="mt-1 text-sm text-zinc-500">Manage your own profile data.</p>
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
