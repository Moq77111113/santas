
type Group = {
	id: number;
	name: string;
};

type EnrichedGroup = Group & {
	owner: Member;
};

type GroupConfig = {
	maxMemberExclusions: number
}

type Member = {
	id: number;
	name: string;
};

type GroupExclusion = {
	member: Member;
	excludedMembers: Member[];
};

type Event<K extends string, T = unknown> = {
	type: K;
	data: T;
};


type KeepAliveEvent = Event<"keep-alive",null>;

type ExclusionsEvent = Event<"exclusions",GroupExclusion[]>

type GroupConfigEvent = Event<"config", GroupConfig>;


type Events = KeepAliveEvent | ExclusionsEvent | GroupConfigEvent;
export type { EnrichedGroup, Event, Events, ExclusionsEvent, Group, GroupConfig, GroupConfigEvent, GroupExclusion, KeepAliveEvent, Member };

