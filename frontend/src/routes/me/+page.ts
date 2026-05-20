import type { PageLoad } from './$types';
import { api, ApiError } from '$lib/api/client';

export const load: PageLoad = async () => {
	try {
		const [user, statuses] = await Promise.all([api.me.get(), api.reference.statuses()]);
		return { user, statuses, error: null as string | null };
	} catch (e) {
		const msg = e instanceof ApiError ? e.message : String(e);
		return { user: null, statuses: [], error: msg };
	}
};
