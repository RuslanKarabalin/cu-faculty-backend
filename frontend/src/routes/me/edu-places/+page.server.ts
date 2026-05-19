import type { PageServerLoad } from './$types';
import { createServerApi, ApiError } from '$lib/server/api';

export const load: PageServerLoad = async ({ locals, fetch }) => {
	const api = createServerApi(locals.bffCookie, fetch);
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
