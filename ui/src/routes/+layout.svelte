<script>
	import api from '$lib/api';
	import { onDestroy } from 'svelte';
	import '../app.css';

	const unsub = api.groups.subscribe(1, (groups) => {
		console.log(groups);
	});

	onDestroy(() => {
		unsub();
	});
</script>

<div class="app">
	<main>
		<slot />
	</main>
</div>
