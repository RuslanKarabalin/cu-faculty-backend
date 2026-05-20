import type { PageLoad } from './$types';
import { api, ApiError } from '$lib/api/client';

const LIMIT = 20;

export const load: PageLoad = async ({ url }) => {
	const offset = Math.max(0, Number(url.searchParams.get('offset') ?? '0') || 0);
	try {
		const data = await api.students.list(LIMIT, offset);
		return { data, offset, limit: LIMIT, error: null as string | null };
	} catch (e) {
		const msg = e instanceof ApiError ? e.message : String(e);
		return { data: null, offset, limit: LIMIT, error: msg };
	}
};
