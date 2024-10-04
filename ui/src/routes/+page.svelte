<script>
	import { Button } from '$lib/components/ui/button';
	import {
		Card,
		CardContent,
		CardDescription,
		CardFooter,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { ScrollArea } from '$lib/components/ui/scroll-area';
</script>

<svelte:head>
	<title>Home</title>
	<meta name="description" content="Svelte demo app" />
</svelte:head>

<Card class="w-full max-w-2xl mx-auto">
	<CardHeader>
		<CardTitle>Secret Santa Generator</CardTitle>
		<CardDescription
			>Add participants, set exclusions, and generate Secret Santa assignments</CardDescription
		>
	</CardHeader>
	<CardContent>
		<div class="space-y-4">
			<div class="flex space-x-2">
				<Input type="text" placeholder="Enter participant name" />
				<Button>Add</Button>
			</div>
			<ScrollArea class="h-full w-full rounded-md border p-4">
				<!-- {participants.map((participant, index) => (
            <div key={index} class="mb-4 pb-4 border-b last:border-b-0">
              <div class="flex justify-between items-center mb-2">
                <span class="font-semibold">{participant.name}</span>
                <Button variant="destructive" size="sm" onClick={() => removeParticipant(index)}>Remove</Button>
              </div>
              <div class="pl-4">
                <Label class="text-sm font-medium mb-2 block">Exclusions:</Label>
                {participants.filter(p => p.name !== participant.name).map((otherParticipant, otherIndex) => (
                  <div key={otherIndex} class="flex items-center space-x-2 mb-1">
                    <Checkbox
                      id={`${index}-${otherIndex}`}
                      checked={participant.exclusions.includes(otherParticipant.name)}
                      onCheckedChange={() => toggleExclusion(index, otherParticipant.name)}
                    />
                    <Label htmlFor={`${index}-${otherIndex}`}>{otherParticipant.name}</Label>
                  </div>
                ))}
              </div>
            </div>
          ))} -->
			</ScrollArea>
			<!-- {error && <div class="text-red-500">{error}</div>} -->
			<Button class="w-full" disabled>Generate Assignments</Button>
		</div>
	</CardContent>
	<CardFooter>
		<div class="w-full">
			<h3 class="text-lg font-semibold mb-2">Assignments:</h3>
			<!-- {assignments.map((assignment, index) => (
          <div key={index} class="mb-1">
            <strong>{assignment.giver}</strong> gives a gift to <strong>{assignment.receiver}</strong>
          </div>
        ))} -->
		</div>
	</CardFooter>
</Card>
