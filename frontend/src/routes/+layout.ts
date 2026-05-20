import { browser } from '$app/environment';
import { hasBffCookie } from '$lib/auth';
import type { LayoutLoad } from './$types';

// Fully static SPA: render only in the browser and prerender the app shell.
export const ssr = false;
export const prerender = true;

export const load: LayoutLoad = () => {
	return { authed: browser ? hasBffCookie() : false };
};
