<script lang="ts">
	import Card from '$lib/components/Card.svelte';
</script>

<div class="space-y-6">
	<div>
		<h1 class="text-2xl font-bold tracking-tight">Dashboard</h1>
		<p class="mt-1 text-sm text-zinc-500">
			Admin panel for the CU Faculty backend.
		</p>
	</div>

	<div class="grid gap-4 md:grid-cols-2">
		<a href="/students" class="block">
			<Card title="Students">
				<p class="text-sm text-zinc-600">
					Browse all registered students and view their profiles, education and work history,
					socials and skills.
				</p>
			</Card>
		</a>
		<a href="/me" class="block">
			<Card title="My profile">
				<p class="text-sm text-zinc-600">
					Manage your own profile data: education places, work places, socials, key skills and
					soft skills.
				</p>
			</Card>
		</a>
	</div>
</div>
