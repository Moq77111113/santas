import { Client } from '$lib/api/client';
import type { EnrichedGroup, Group, GroupExclusion, Member } from '../dto';

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

	async exlusions(id: number): Promise<GroupExclusion[]> {
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

	subscribe(id: number, callback: (exclusions: GroupExclusion[]) => void): () => void {
		if (typeof window === 'undefined') {
			throw new Error('EventSource is not supported');
		}
		const eventSource = new EventSource(`${base}/${id}/events`);
		eventSource.onmessage = (event) => {
			try {
				const exclusions = JSON.parse(event.data);
				callback(exclusions);
			} catch (e) {
				// TODO: handle error
				console.error(e);
			}
		};
		return () => {
			eventSource.close();
		};
	}
}

export default GroupService;
