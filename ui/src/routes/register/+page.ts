import { goto } from '$app/navigation';
import type { PageLoad } from './$types';
import api from '@/lib/api/client';

export const load = (async ({ url }) => {
	const qs = url.searchParams;
	const groupId = qs.get('groupId');

	const me = await api.auth.me().catch(() => null);
	if (me) {
		return goto('/');
	}
	if (!groupId) {
		return {
			groupId: null
		};
	}

	const group = await api.groups.group(parseInt(groupId)).catch(() => null);

	return {
		groupId: group?.id || null
	};
}) satisfies PageLoad;
