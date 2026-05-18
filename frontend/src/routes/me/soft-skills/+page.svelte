<script lang="ts">
	import { api } from '$lib/api/client';
	import SkillsManager from '$lib/components/SkillsManager.svelte';
</script>

<SkillsManager title="My soft skills" api={api.me.softSkills} />
