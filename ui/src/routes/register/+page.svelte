<script lang="ts">
	import api from '@/lib/api/client';
	import Button from '@/lib/components/ui/button/button.svelte';
	import Input from '@/lib/components/ui/input/input.svelte';

	import type { EventHandler } from 'svelte/elements';
	import type { PageData } from './$types';
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
			await api.auth.register(name);
			goto('/');
		} catch (e) {
			error = 'Il semble que ce nom soit déjà pris';
			console.error(e);
		}
	};
</script>

<div class="container mx-auto p-4 flex flex-col justify-center items-center">
	<h2 class="text-2xl font-bold mb-6">Quel est ton petit nom ?</h2>
	<form onsubmit={create} method="post" class="min-w-60 flex flex-col space-y-4">
		<Input onchange={() => (error = '')} type="text" name="name" />
		<Button type="submit">S'enregistrer</Button>
		{#if error}
			<p class="text-red-500">{error}</p>
		{/if}
	</form>
</div>
