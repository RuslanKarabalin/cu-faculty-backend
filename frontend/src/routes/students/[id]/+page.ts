import type { PageLoad } from './$types';
import { api, ApiError } from '$lib/api/client';
import { error } from '@sveltejs/kit';

export const prerender = false;

export const load: PageLoad = async ({ params }) => {
	try {
		const [user, eduPlaces, workPlaces, socials, keySkills, softSkills] = await Promise.all([
			api.students.get(params.id),
			api.students.eduPlaces(params.id),
			api.students.workPlaces(params.id),
			api.students.socials(params.id),
			api.students.keySkills(params.id),
			api.students.softSkills(params.id)
		]);
		return { user, eduPlaces, workPlaces, socials, keySkills, softSkills };
	} catch (e) {
		if (e instanceof ApiError && e.status === 404) {
			throw error(404, 'Student not found');
		}
		const msg = e instanceof ApiError ? e.message : String(e);
		throw error(e instanceof ApiError ? e.status : 500, msg);
	}
};
