<script lang="ts">
	import { type Snippet } from 'svelte';
	import '../app.css';
	import type { PageData } from './$types';

	type Props = {
		children: Snippet;
		data: PageData;
	};
	const { children, data }: Props = $props();
</script>

<div class="min-h-dvh flex flex-col">
	<header
		class="h-20 border-b border-muted-foreground flex justify-between p-4 md:p-6 lg:p-8 bg-primary text-primary-foreground"
	>
		<a href="/" class="flex items-center justify-center">
			<h1 class="text-2xl font-bold">Secret Santa !</h1>
		</a>
		{#if data.me}
			<span class="text-xs rounded-xl self-end bg-secondary p-2">
				{data.me.name}
			</span>
		{/if}
	</header>
	<main class="flex flex-grow">
		{@render children()}
	</main>
	<footer class="h-20 border-t border-muted flex items-center justify-center">
		<p class="text-sm">© 2024 Secret Santa</p>
	</footer>
</div>
