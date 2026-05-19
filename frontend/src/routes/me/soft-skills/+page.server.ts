import type { PageServerLoad } from './$types';
import { createServerApi, ApiError } from '$lib/server/api';

export const load: PageServerLoad = async ({ locals, fetch }) => {
	const api = createServerApi(locals.bffCookie, fetch);
	try {
		const [items, all] = await Promise.all([
			api.me.softSkills.list(),
			api.reference.softSkills()
		]);
		return { items, all, error: null as string | null };
	} catch (e) {
		const msg = e instanceof ApiError ? e.message : String(e);
		return { items: [], all: [], error: msg };
	}
};
