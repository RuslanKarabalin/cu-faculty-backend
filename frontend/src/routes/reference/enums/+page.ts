import type { PageLoad } from './$types';
import { api, ApiError } from '$lib/api/client';

export const load: PageLoad = async () => {
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
