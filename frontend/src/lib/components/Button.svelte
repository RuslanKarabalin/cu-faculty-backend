<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLButtonAttributes } from 'svelte/elements';

	type Variant = 'primary' | 'secondary' | 'danger' | 'ghost';
	type Size = 'sm' | 'md';

	type Props = HTMLButtonAttributes & {
		variant?: Variant;
		size?: Size;
		children: Snippet;
	};

	let { variant = 'primary', size = 'md', class: klass = '', children, ...rest }: Props = $props();

	const variants: Record<Variant, string> = {
		primary: 'bg-zinc-900 text-white hover:bg-zinc-800 disabled:bg-zinc-400',
		secondary: 'bg-white text-zinc-900 border border-zinc-300 hover:bg-zinc-50',
		danger: 'bg-red-600 text-white hover:bg-red-700',
		ghost: 'bg-transparent text-zinc-700 hover:bg-zinc-100'
	};
	const sizes: Record<Size, string> = {
		sm: 'h-8 px-3 text-sm',
		md: 'h-10 px-4 text-sm'
	};
</script>

<button
	{...rest}
	class="inline-flex items-center justify-center rounded-md font-medium transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-zinc-900 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-60 {variants[
		variant
	]} {sizes[size]} {klass}"
>
	{@render children()}
</button>
