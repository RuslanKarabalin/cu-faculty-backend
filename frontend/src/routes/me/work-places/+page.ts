import type { PageLoad } from './$types';
import { api, ApiError } from '$lib/api/client';

export const load: PageLoad = async () => {
	try {
		const items = await api.me.workPlaces.list();
		return { items, error: null as string | null };
	} catch (e) {
		const msg = e instanceof ApiError ? e.message : String(e);
		return { items: [], error: msg };
	}
};
