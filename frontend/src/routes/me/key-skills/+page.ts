import type { PageLoad } from './$types';
import { api, ApiError } from '$lib/api/client';

export const load: PageLoad = async () => {
	try {
		const [items, all] = await Promise.all([api.me.keySkills.list(), api.reference.keySkills()]);
		return { items, all, error: null as string | null };
	} catch (e) {
		const msg = e instanceof ApiError ? e.message : String(e);
		return { items: [], all: [], error: msg };
	}
};
