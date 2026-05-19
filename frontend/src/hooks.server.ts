import type { Handle } from '@sveltejs/kit';
import { BFF_COOKIE_NAME } from '$lib/server/auth';

export const handle: Handle = async ({ event, resolve }) => {
	event.locals.bffCookie = event.cookies.get(BFF_COOKIE_NAME) ?? null;
	return resolve(event);
};
