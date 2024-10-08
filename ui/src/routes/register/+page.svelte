<script lang="ts">
	import api from '@/lib/api/client';
	import Button from '@/lib/components/ui/button/button.svelte';
	import Input from '@/lib/components/ui/input/input.svelte';

	import type { EventHandler } from 'svelte/elements';
	import type { PageData } from './$types';

	type Props = {
		data: PageData;
	};
	
	const { data }: Props = $props();

	const create: EventHandler<SubmitEvent, HTMLFormElement> = async (e) => {
		e.preventDefault();
		const form = e.currentTarget;
		const formData = new FormData(form);

		const name = formData.get('name') as string;
		await api.auth.register(name);
	};
</script>

<form onsubmit={create} method="post">
	<Input type="text" name="name" />
	<Button type="submit">Create Group</Button>
</form>
