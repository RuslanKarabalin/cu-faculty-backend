import { dev } from '$app/environment';
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { BFF_COOKIE_NAME, BFF_COOKIE_MAX_AGE } from '$lib/server/auth';

export const POST: RequestHandler = async ({ request, cookies }) => {
	const body = (await request.json().catch(() => null)) as { value?: string } | null;
	const value = body?.value?.trim();
	if (!value) {
		return json({ error: 'value required' }, { status: 400 });
	}
	cookies.set(BFF_COOKIE_NAME, value, {
		path: '/',
		httpOnly: true,
		sameSite: 'lax',
		secure: !dev,
		maxAge: BFF_COOKIE_MAX_AGE
	});
	return json({ ok: true });
};

export const DELETE: RequestHandler = async ({ cookies }) => {
	cookies.delete(BFF_COOKIE_NAME, { path: '/' });
	return json({ ok: true });
};
