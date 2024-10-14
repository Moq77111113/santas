import { Client } from '$lib/api/client';
import type { Member } from '../dto';

const base = `/api/auth` as const;
class AuthService {
	readonly client: Client;

	constructor(client: Client) {
		this.client = client;
	}

	me(): Promise<Member> {
		return this.client.request(`${base}/me`);
	}

	async register(name: string): Promise<void> {
		const formData = new FormData();
		formData.append('name', name);
		await this.client.request(`${base}/register`, {
			method: 'POST',
			body: formData
		});
	}
}

export default AuthService;
