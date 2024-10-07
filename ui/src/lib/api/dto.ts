type Group = {
	id: number;
	name: string;
};

type Member = {
	id: number;
	name: string;
};

type GroupExclusion = {
	member: Member;
	excludedMembers: Member[];
};

export type { Group, Member, GroupExclusion };
