import { Client } from '$lib/api/client';
import type { EnrichedGroup, ExclusionsEvent, Group, GroupConfig, GroupConfigEvent, GroupExclusion, Member } from '../dto';

const base = `/api/group` as const;
class GroupService {
	readonly client: Client;

	constructor(client: Client) {
		this.client = client;
	}

	async list(): Promise<EnrichedGroup[]> {
		return this.client.request(base);
	}

	async group(id: number): Promise<EnrichedGroup> {
		return this.client.request(`${base}/${id}`);
	}

	async members(id: number): Promise<Member[]> {
		return this.client.request(`${base}/${id}/members`);
	}

	async config(id: number): Promise<GroupConfig> {
		return this.client.request(`${base}/${id}/config`);
	}

	async exclusions(id: number): Promise<GroupExclusion[]> {
		return this.client.request(`${base}/${id}/exclusion`);
	}

	async create(name: string): Promise<Group> {
		const formData = new FormData();
		formData.append('name', name);
		return this.client.request(base, {
			method: 'POST',
			body: formData
		});
	}

	async join(groupId: number): Promise<void> {
		return this.client.request(`${base}/${groupId}/join`, {
			method: 'POST'
		});
	}

	async removeMember(groupId: number, memberId: number): Promise<void> {
		return this.client.request(`${base}/${groupId}/member/${memberId}`, {
			method: 'DELETE'
		});
	}

	async addExclusion(groupId: number, memberId: number, excludeId: number): Promise<void> {
		const formData = new FormData();
		formData.append('memberId', `${excludeId}`);
		return this.client.request(`${base}/${groupId}/member/${memberId}/exclusion`, {
			method: 'POST',
			body: formData
		});
	}

	async removeExclusion(groupId: number, memberId: number, excludeId: number): Promise<void> {
		return this.client.request(`${base}/${groupId}/member/${memberId}/exclusion/${excludeId}`, {
			method: 'DELETE'
		});
	}

	async getSantas(id: number): Promise<Member[]> {
		return this.client.request(`${base}/${id}/santas`);
	}



	 
	subscribe(id: number, on: (ev: ExclusionsEvent | GroupConfigEvent) => void): () => void {
		return this.client.subscribe(`${base}/${id}/events`, on);
	}
}

export default GroupService;
