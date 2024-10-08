import api from '@/lib/api/client';
import type { LayoutLoad } from './$types';
import { goto } from '$app/navigation';

export const prerender = true;
export const ssr = false;

export const load = (async ({ url }) => {
	const groups = await api.groups.list();

	const me = await api.auth.me().catch(() => null);
	if (!me && url.pathname !== '/register') {
		return goto('/register');
	}

	return { groups, me };
}) satisfies LayoutLoad;
