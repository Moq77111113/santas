import api from '@/lib/api/client';
import { redirect } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const prerender = false;
export const ssr = false;

export const load = (async ({ params }) => {
	const id = parseInt(params.id);
	if (isNaN(id)) {
		redirect(307, '/');
	}

	const groupWithExclusions = (await api.groups.exclusions(id)) || [];
	const group = await api.groups.group(id);
	const config = await api.groups.config(id);
	return {
		groupWithExclusions,
		id,
		group, 
		config
	};
}) satisfies PageLoad;
