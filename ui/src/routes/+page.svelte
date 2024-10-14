<script lang="ts">
	import type { PageData } from './$types';
	import { Card, CardHeader, CardTitle } from '$lib/components/ui/card';
	import api from '@/lib/api/client';
	import type { EventHandler } from 'svelte/elements';
	import { Input } from '@/lib/components/ui/input';
	import { Button } from '@/lib/components/ui/button';
	import { goto } from '$app/navigation';
	type Props = {
		data: PageData;
	};
	const { data }: Props = $props();

	let error = $state('');
	const create: EventHandler<SubmitEvent, HTMLFormElement> = async (e) => {
		e.preventDefault();
		const form = e.currentTarget;
		const formData = new FormData(form);

		const name = formData.get('name') as string;
		try {
			const group = await api.groups.create(name);
			goto(`/g/${group.id}`);
		} catch (e) {
			error = 'Il semble que ce nom soit déjà pris';
			console.error(e);
		}
	};
</script>

<div class="container mx-auto p-4 flex flex-col space-y-8">
	<article class="flex flex-col space-y-4 items-start">
		<h2 class="text-2xl font-bold">Rejoindre un groupe</h2>
		<div class="flex items-center">
			<div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
				{#if data.groups?.length}
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
					<p>Aucun groupe pour le moment</p>
				{/if}
			</div>
		</div>
	</article>

	<article class="flex flex-col space-y-4 items-start">
		<h3 class="text-xl font-bold">Tu ne trouves pas ton bonheur ? Crée un groupe</h3>

		<form onsubmit={create} method="post" class="min-w-60 flex flex-col space-y-4">
			<Input onchange={() => (error = '')} type="text" name="name" />
			<Button type="submit">Go !</Button>
			{#if error}
				<p class="text-red-500">{error}</p>
			{/if}
		</form>
	</article>
</div>
