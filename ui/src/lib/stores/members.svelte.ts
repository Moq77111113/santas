import type { GroupExclusion } from '../api/dto';

export function createMemberState(initial: GroupExclusion[]) {
	let members = $state(initial);

	return {
		...members,
		get exclusions() {
			return members;
		},
		set exclusions(value: GroupExclusion[]) {
			members = value;
		}
	};
}

export type MemberState = ReturnType<typeof createMemberState>;
