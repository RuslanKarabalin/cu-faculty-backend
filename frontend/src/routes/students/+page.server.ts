import type { PageServerLoad } from './$types';
import { createServerApi, ApiError } from '$lib/server/api';

const LIMIT = 20;

export const load: PageServerLoad = async ({ locals, fetch, url }) => {
	const offset = Math.max(0, Number(url.searchParams.get('offset') ?? '0') || 0);
	const api = createServerApi(locals.bffCookie, fetch);
	try {
		const data = await api.students.list(LIMIT, offset);
		return { data, offset, limit: LIMIT, error: null as string | null };
	} catch (e) {
		const msg = e instanceof ApiError ? e.message : String(e);
		return { data: null, offset, limit: LIMIT, error: msg };
	}
};
