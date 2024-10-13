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

	import { goto } from '$app/navigation';
	import api from '@/lib/api/client';
	import { onDestroy } from 'svelte';

	type Props = {
		data: PageData;
	};
	const { data }: Props = $props();

	let groupWithExclusions = $state(data.groupWithExclusions);
	const unSubsribe = api.groups.subscribe(data.id, (exc) => {
		groupWithExclusions = exc;
	});

	const isOwner = $derived(data.group.owner.id === data.me?.id);
	const isMember = $derived(groupWithExclusions.some((_) => _.member.id === data.me?.id));
	const join = async () => {
		await api.groups.join(data.id);
	};

	const remove = async (id: number) => {
		await api.groups.removeMember(data.id, id);
	};

	const leave = async () => {
		if (!data.me) return;
		await api.groups.removeMember(data.id, data.me.id);
		goto('/');
	};

	const toggleExclusion = async (value: boolean, id: number, excludeId: number) => {
		if (value) {
			return await api.groups.addExclusion(data.id, id, excludeId);
		}

		return await api.groups.removeExclusion(data.id, id, excludeId);
	};

	const memberExclusionCounts = $derived(
		groupWithExclusions.reduce(
			(acc, { member }) => {
				acc[member.id] =
					groupWithExclusions.filter((_) => _.excludedMembers.some((m) => m.id === member.id))
						.length < data.config.maxMemberExclusions;
				return acc;
			},
			{} as Record<number, boolean>
		)
	);

	const s = async () => {
		const c = await api.groups.getSantas(data.id);
		console.log(c);
	};

	onDestroy(() => {
		unSubsribe();
	});
</script>

<svelte:head>
	<title>Le Noel des Chmoly</title>
	<meta name="description" content={`Le noel de ${data.group.name}`} />
</svelte:head>

{#snippet exclusions({ member, excludedMembers }: GroupExclusion)}
	<Label class="text-sm font-medium mb-2 block">Exclusions:</Label>
	<div
		class=" w-full rounded-md border-muted mb-2 grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4"
	>
		{#each groupWithExclusions
			.filter((m) => m.member.name !== member.name)
			.sort((a, b) => a.member.name.localeCompare(b.member.name)) as { member: otherMember }}
			<div class="flex items-center space-x-2 space-y-2 p-1">
				<Checkbox
					disabled={!isMember ||
						(!excludedMembers.some((_) => _.id === otherMember.id) &&
							(!memberExclusionCounts[otherMember.id] ||
								excludedMembers.length >= data.config.maxMemberExclusions))}
					id={`${member.id}-${otherMember.id}`}
					checked={excludedMembers.some((_) => _.id === otherMember.id)}
					onCheckedChange={(v) => {
						if (typeof v !== 'boolean') return;
						toggleExclusion(v, member.id, otherMember.id);
					}}
				/>
				<Label for={`${member.id}-${otherMember.id}`}>{otherMember.name}</Label>
			</div>
		{/each}
	</div>
{/snippet}

<div class="container mx-auto p-4 md:p-6 lg:p-8">
	<div class="mx-auto space-y-6">
		<h1 class="text-2xl font-bold text-center mb-6">Le Noël de {data.group.name}</h1>
		<Button variant="secondary" on:click={s}>Générer les pères Noël</Button>
		{#if !groupWithExclusions.some((_) => _.member.id === data.me?.id)}
			<Button on:click={join}>Rejoindre</Button>
		{:else}
			<Button variant="destructive" on:click={leave}>Quitter</Button>
		{/if}

		<article class=" rounded-lg shadow-md flex flex-col space-y-4">
			<h2 class="text-lg font-semibold mb-3">Participants & Exclusions</h2>
			<span class="text-md">{groupWithExclusions.length} Membres</span>
			<ScrollArea class="h-[300px] md:h-[600px] w-full rounded-md border p-4">
				<Accordion multiple class="w-full">
					{#each groupWithExclusions as member}
						<AccordionItem value={`${member.member.id}`}>
							<AccordionTrigger class="text-left w-full">
								<div class="flex justify-between items-center w-full mx-2">
									<span>{member.member.name}</span>
									{#if data.group.owner.id === data.me?.id}
										<Button
											variant="destructive"
											size="sm"
											onclick={() => {
												remove(member.member.id);
											}}>Retirer</Button
										>
									{/if}
								</div>
							</AccordionTrigger>
							{#if groupWithExclusions.filter((m) => m.member.name !== member.member.name).length}
								<AccordionContent
									><div class="m-2">
										{@render exclusions(member)}
									</div>
								</AccordionContent>
							{/if}
						</AccordionItem>
					{/each}
				</Accordion>
			</ScrollArea>

			<Button class="w-full" disabled={!isOwner}>Générer les pères Noël</Button>
		</article>
	</div>
</div>
