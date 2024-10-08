<script lang="ts">
	import { Button } from '$lib/components/ui/button';

	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Label } from '$lib/components/ui/label';
	import { ScrollArea } from '$lib/components/ui/scroll-area';
	import { type GroupExclusion } from '@/lib/api/dto';

	import {
		Accordion,
		AccordionContent,
		AccordionItem,
		AccordionTrigger
	} from '$lib/components/ui/accordion';
	import type { PageData } from './$types';

	import api from '@/lib/api/client';
	import { onDestroy } from 'svelte';

	type Props = {
		data: PageData;
	};
	const { data }: Props = $props();

	let groupWithExclusions = $state(data.groupWithExclusions);
	const unsubsribe = api.groups.subscribe(data.id, (exc) => {
		groupWithExclusions = exc;
	});

	onDestroy(() => {
		unsubsribe();
	});
</script>

<svelte:head>
	<title>Le Noel des Chmoly</title>
	<meta name="description" content="Le Noel des Chmoly" />
</svelte:head>

{#snippet exclusions({ member, excludedMembers }: GroupExclusion)}
	<Label class="text-sm font-medium mb-2 block">Exclusions:</Label>
	<ScrollArea class="max-h-40 w-full rounded-md border-muted mb-2 flex flex-col">
		{#each groupWithExclusions
			.filter((m) => m.member.name !== member.name)
			.sort((a, b) => a.member.name.localeCompare(b.member.name)) as { member: otherMember }}
			<div class="flex items-center space-x-2 mb-1">
				<Checkbox
					id={`${member.id}-${otherMember.id}`}
					checked={excludedMembers.some((_) => _.id === otherMember.id)}
				/>
				<Label for={`${member.id}-${otherMember.id}`}>{otherMember.name}</Label>
			</div>
		{/each}
	</ScrollArea>
{/snippet}

<div class="p-4 md:p-6 lg:p-8">
	<div class="mx-auto space-y-6">
		<h1 class="text-2xl font-bold text-center mb-6">Secret Santa Generator</h1>

		<article class="p-4 rounded-lg shadow-md flex flex-col space-y-4">
			<h2 class="text-lg font-semibold mb-3">Participants & Exclusions</h2>
			<ScrollArea class="h-[300px] md:h-[600px] w-full rounded-md border p-4">
				<Accordion multiple class="w-full">
					{#each groupWithExclusions as member}
						<AccordionItem value={`${member.member.id}`}>
							<AccordionTrigger class="text-left w-full"
								><div class="flex justify-between items-center w-full">
									<span>{member.member.name}</span>
									<Button variant="destructive" size="sm">Remove</Button>
								</div>
							</AccordionTrigger>
							<AccordionContent
								><div class="pl-4 mt-2">
									{@render exclusions(member)}
								</div>
							</AccordionContent>
						</AccordionItem>
					{/each}
				</Accordion>
			</ScrollArea>
			<!-- {error && <div class="text-red-500">{error}</div>} -->
			<Button class="w-full" disabled>Generate Assignments</Button>
		</article>
	</div>
</div>
