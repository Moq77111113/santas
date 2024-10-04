<script lang="ts">
	import { applyAction, enhance } from '$app/forms';
	import api from '@/lib/api/client';
	import Button from '@/lib/components/ui/button/button.svelte';
	import Input from '@/lib/components/ui/input/input.svelte';

	import type { EventHandler } from 'svelte/elements';

	const create: EventHandler<SubmitEvent, HTMLFormElement> = async (e) => {
		const form = e.currentTarget;
		const data = new FormData(form);

		const name = data.get('name') as string;
		await api.groups.addMember(1, name);
	};
</script>

<form on:submit|preventDefault={create} method="post">
	<Input type="text" name="name" />
	<Button type="submit">Create Group</Button>
</form>
