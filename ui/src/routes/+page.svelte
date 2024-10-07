<script lang="ts">
	import { Button } from '$lib/components/ui/button';

	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { ScrollArea } from '$lib/components/ui/scroll-area';
	import { type GroupExclusion } from '@/lib/api/dto';

	import { type MemberState } from '@/lib/stores/members.svelte';
	import { getContext } from 'svelte';

	const members = getContext<MemberState>('members');
</script>

<svelte:head>
	<title>Le Noel des Chmoly</title>
	<meta name="description" content="Le Noel des Chmoly" />
</svelte:head>

{#snippet exclusions({ member, excludedMembers }: GroupExclusion)}
	<Label class="text-sm font-medium mb-2 block">Exclusions:</Label>
	<ScrollArea class="max-h-40 w-full rounded-md border-muted mb-2 flex flex-col">
		{#each members.exclusions
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
<div class="min-h-screen p-4 md:p-6 lg:p-8">
	<div class="max-w-lg mx-auto space-y-6">
		<h1 class="text-2xl font-bold text-center mb-6">Secret Santa Generator</h1>

		<div class="p-4 rounded-lg shadow-md">
			<h2 class="text-lg font-semibold mb-3">Add Participants</h2>
			<div class="flex space-x-2">
				<Input type="text" placeholder="Enter participant name" class="flex-grow" />
				<Button>Add</Button>
			</div>
		</div>

		<article class="p-4 rounded-lg shadow-md">
			<h2 class="text-lg font-semibold mb-3">Participants & Exclusions</h2>
			<ScrollArea class="h-[300px] w-full rounded-md border p-4">
				{#each members.exclusions as membersWithExclusions, index}
					<div class="flex justify-between items-center mb-2">
						<span class="font-semibold">{membersWithExclusions.member.name}</span>
						<Button variant="destructive" size="sm">Remove</Button>
					</div>
					<div class="pl-4">
						{@render exclusions(membersWithExclusions)}
					</div>
				{/each}
			</ScrollArea>
			<!-- {error && <div class="text-red-500">{error}</div>} -->
			<Button class="w-full" disabled>Generate Assignments</Button>
		</article>
	</div>
</div>
