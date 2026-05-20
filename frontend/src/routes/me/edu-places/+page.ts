import type { PageLoad } from './$types';
import { api, ApiError } from '$lib/api/client';

export const load: PageLoad = async () => {
	try {
		const [items, universities] = await Promise.all([
			api.me.eduPlaces.list(),
			api.reference.universities()
		]);
		return { items, universities, error: null as string | null };
	} catch (e) {
		const msg = e instanceof ApiError ? e.message : String(e);
		return { items: [], universities: [], error: msg };
	}
};
