import type { RequestHandler } from './$types';
import { API_BASE_URL } from '$lib/server/api';

const PASSTHROUGH_METHODS = new Set(['GET', 'POST', 'PUT', 'PATCH', 'DELETE']);

const forward: RequestHandler = async ({ params, request, locals, fetch }) => {
	if (!PASSTHROUGH_METHODS.has(request.method)) {
		return new Response(null, { status: 405 });
	}

	const target = new URL(request.url);
	const upstreamUrl = `${API_BASE_URL}/api/${params.path}${target.search}`;

	const headers = new Headers();
	const contentType = request.headers.get('content-type');
	if (contentType) headers.set('content-type', contentType);
	if (locals.bffCookie) {
		headers.set('cookie', `bff.cookie=${locals.bffCookie}`);
	}

	const init: RequestInit = { method: request.method, headers };
	if (request.method !== 'GET' && request.method !== 'DELETE') {
		init.body = await request.arrayBuffer();
	}

	const res = await fetch(upstreamUrl, init);
	const responseHeaders = new Headers();
	const upstreamCT = res.headers.get('content-type');
	if (upstreamCT) responseHeaders.set('content-type', upstreamCT);

	return new Response(res.body, { status: res.status, headers: responseHeaders });
};

export const GET = forward;
export const POST = forward;
export const PUT = forward;
export const PATCH = forward;
export const DELETE = forward;
