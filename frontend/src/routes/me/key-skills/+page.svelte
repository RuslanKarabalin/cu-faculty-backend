<script lang="ts">
	import { api } from '$lib/api/client';
	import SkillsManager from '$lib/components/SkillsManager.svelte';
</script>

<SkillsManager title="My key skills" api={api.me.keySkills} listAll={api.reference.keySkills} />
