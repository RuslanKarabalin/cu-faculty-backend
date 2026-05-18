<script lang="ts">
	import './layout.css';
	import { page } from '$app/state';
	import favicon from '$lib/assets/favicon.svg';
	import BffCookieInput from '$lib/components/BffCookieInput.svelte';

	let { children } = $props();

	const links = [
		{ href: '/', label: 'Dashboard' },
		{ href: '/students', label: 'Students' },
		{ href: '/me', label: 'My profile' }
	];

	function isActive(href: string) {
		if (href === '/') return page.url.pathname === '/';
		return page.url.pathname === href || page.url.pathname.startsWith(href + '/');
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>CU Faculty admin</title>
</svelte:head>

<div class="flex min-h-screen bg-zinc-50 text-zinc-900">
	<aside class="hidden w-60 shrink-0 border-r border-zinc-200 bg-white p-4 md:block">
		<div class="mb-6 px-2">
			<div class="text-base font-bold text-zinc-900">CU Faculty</div>
			<div class="text-xs text-zinc-500">Admin panel</div>
		</div>
		<nav class="space-y-1">
			{#each links as { href, label } (href)}
				<a
					{href}
					class="block rounded-md px-3 py-2 text-sm font-medium {isActive(href)
						? 'bg-zinc-100 text-zinc-900'
						: 'text-zinc-600 hover:bg-zinc-50 hover:text-zinc-900'}"
				>
					{label}
				</a>
			{/each}
		</nav>
	</aside>

	<main class="min-w-0 flex-1">
		<header
			class="flex flex-wrap items-center justify-between gap-4 border-b border-zinc-200 bg-white px-6 py-3 md:px-10"
		>
			<nav class="flex gap-4 md:hidden">
				{#each links as { href, label } (href)}
					<a
						{href}
						class="text-sm font-medium {isActive(href)
							? 'text-zinc-900'
							: 'text-zinc-500 hover:text-zinc-900'}"
					>
						{label}
					</a>
				{/each}
			</nav>
			<div class="ml-auto"><BffCookieInput /></div>
		</header>
		<div class="px-6 py-6 md:px-10 md:py-8">{@render children()}</div>
	</main>
</div>
