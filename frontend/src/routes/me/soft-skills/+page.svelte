<script lang="ts">
	import { api } from '$lib/api/client';
	import Alert from '$lib/components/Alert.svelte';
	import SkillsManager from '$lib/components/SkillsManager.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

{#if data.error}<Alert variant="error">{data.error}</Alert>{/if}

<SkillsManager
	title="My soft skills"
	items={data.items}
	all={data.all}
	api={api.me.softSkills}
/>
