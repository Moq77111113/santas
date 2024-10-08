<script lang="ts">
	import type { PageData } from './$types';
	import { Card, CardHeader, CardTitle } from '$lib/components/ui/card';
	type Props = {
		data: PageData;
	};
	const { data }: Props = $props();
</script>

<div class="container mx-auto p-4 flex flex-col">
	<h2 class="text-2xl font-bold mb-6">Groups</h2>
	<div class="flex items-center mb-6">
		<div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
			{#if data.groups}
				{#each data.groups as group}
					<a href={`/g/${group.id}`} data-sveltekit-preload-data="hover">
						<Card class="group hover:shadow-md hover:bg-muted transition-shadow duration-200">
							<CardHeader class="p-4">
								<CardTitle class="flex flex-col items-center text-center">
									<span class="text-sm font-medium">{group.name}</span>
								</CardTitle>
							</CardHeader>
						</Card>
					</a>
				{/each}
			{:else}
				<p>No groups found</p>
			{/if}
		</div>
	</div>
</div>
