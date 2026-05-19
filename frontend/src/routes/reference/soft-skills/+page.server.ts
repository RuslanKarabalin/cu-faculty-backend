import type { PageServerLoad } from './$types';
import { createServerApi, ApiError } from '$lib/server/api';

export const load: PageServerLoad = async ({ locals, fetch }) => {
	const api = createServerApi(locals.bffCookie, fetch);
	try {
		return { items: await api.reference.softSkills(), error: null as string | null };
	} catch (e) {
		return { items: [], error: e instanceof ApiError ? e.message : String(e) };
	}
};
