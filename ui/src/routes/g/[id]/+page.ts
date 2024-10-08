import api from '@/lib/api/client';
import type { PageLoad } from './$types';
import { redirect } from '@sveltejs/kit';

export const prerender = false;
export const ssr = false;

export const load = (async ({ params }) => {
	const id = parseInt(params.id);
	if (isNaN(id)) {
		redirect(307, '/');
	}
	const groupWithExclusions = await api.groups.exlusions(id);

	return {
		groupWithExclusions,
        id
	};
}) satisfies PageLoad;
