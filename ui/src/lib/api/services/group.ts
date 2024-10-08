import { Client } from '$lib/api/client';
import type { Group, GroupExclusion, Member } from '../dto';

const base = `/api/group` as const;
class GroupService {
	readonly client: Client;

	constructor(client: Client) {
		this.client = client;
	}

	async list(): Promise<Group[]> {
		return this.client.request(base);
	}

	async group(id: number): Promise<Group> {
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

	async addMember(groupId: number, name: string): Promise<void> {
		const formData = new FormData();
		formData.append('name', name);
		return this.client.request(`${base}/${groupId}/member`, {
			method: 'POST',
			body: formData
		});
	}

	async addExclusion(groupId: number, memberId: number, name: string): Promise<void> {
		const formData = new FormData();
		formData.append('name', name);
		return this.client.request(`${base}/${groupId}/member/${memberId}/exclusion`, {
			method: 'POST',
			body: formData
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
