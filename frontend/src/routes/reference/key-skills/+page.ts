import type { PageLoad } from './$types';
import { api, ApiError } from '$lib/api/client';

export const load: PageLoad = async () => {
	try {
		return { items: await api.reference.keySkills(), error: null as string | null };
	} catch (e) {
		return { items: [], error: e instanceof ApiError ? e.message : String(e) };
	}
};
