import type { PageServerLoad } from './$types';
import { createServerApi, ApiError } from '$lib/server/api';

export const load: PageServerLoad = async ({ locals, fetch }) => {
	const api = createServerApi(locals.bffCookie, fetch);
	try {
		const [social, edu, work, event] = await Promise.all([
			api.reference.socialNetworks(),
			api.reference.eduGrades(),
			api.reference.workGrades(),
			api.reference.eventCategories()
		]);
		return {
			blocks: [
				{ title: 'Social networks', values: social },
				{ title: 'Edu grades', values: edu },
				{ title: 'Work grades', values: work },
				{ title: 'Event categories', values: event }
			],
			error: null as string | null
		};
	} catch (e) {
		return { blocks: [], error: e instanceof ApiError ? e.message : String(e) };
	}
};
