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
		class="h-30 border-b bg-accent text-muted border-muted-foreground flex items-center justify-between"
	>
		<a href="/" class="h-full w-full flex items-center justify-center">
			<h1 class="text-2xl font-bold">Secret Santa !</h1>
		</a>

		<!-- 
		{#if data.me}
		
			<div class="flex items-center space-x-2">
				<span class="text-sm">Hello, {data.me.name}</span>
			</div>
		{/if}
	 -->
	</header>
	<main class="flex-grow">
		{@render children()}
	</main>
	<footer class="h-20 border-t border-muted flex items-center justify-center">
		<p class="text-sm">© 2024 Secret Santa</p>
	</footer>
</div>
