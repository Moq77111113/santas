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

export type { EnrichedGroup, Group, GroupConfig, GroupExclusion, Member };

