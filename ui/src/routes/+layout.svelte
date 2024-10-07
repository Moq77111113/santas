<script lang="ts">
	import api from '$lib/api';
	import { onDestroy, onMount, setContext } from 'svelte';
	import '../app.css';
	import { createMemberState } from '@/lib/stores/members.svelte';

	const members = createMemberState([]);

	onMount(async () => {
		members.exclusions = await api.groups.exlusions(1);
	});

	setContext('members', members);
	const { children } = $props();

	const unsub = api.groups.subscribe(1, (exc) => {
		members.exclusions = exc;
	});

	onDestroy(() => {
		unsub();
	});
</script>

<div class="app">
	<main>
		{@render children()}
	</main>
</div>
