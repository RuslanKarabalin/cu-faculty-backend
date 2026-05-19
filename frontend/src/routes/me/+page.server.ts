import type { PageServerLoad } from './$types';
import { createServerApi, ApiError } from '$lib/server/api';

export const load: PageServerLoad = async ({ locals, fetch }) => {
	const api = createServerApi(locals.bffCookie, fetch);
	try {
		const [user, statuses] = await Promise.all([api.me.get(), api.reference.statuses()]);
		return { user, statuses, error: null as string | null };
	} catch (e) {
		const msg = e instanceof ApiError ? e.message : String(e);
		return { user: null, statuses: [], error: msg };
	}
};
