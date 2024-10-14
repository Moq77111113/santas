import api from '@/lib/api/client';
import type { PageLoad } from './$types';
export const load = (async () => {
	const groups = await api.groups.list();

	return { groups };
}) satisfies PageLoad;
