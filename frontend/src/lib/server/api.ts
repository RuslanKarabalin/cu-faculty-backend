import { env } from '$env/dynamic/private';
import {
	ApiError,
	buildEndpoints,
	makeRequest,
	type Endpoints,
	type FetchFn,
	type Request
} from '$lib/api/endpoints';

const API_BASE = env.API_BASE_URL ?? 'http://127.0.0.1:8080';

export { ApiError };
export type ServerApi = Endpoints;

export function createServerApi(bffCookie: string | null, fetchFn: FetchFn): Endpoints {
	const request: Request = (path, init = {}) => {
		const headers = new Headers(init.headers);
		if (bffCookie) headers.set('cookie', `bff.cookie=${bffCookie}`);
		return makeRequest(`${API_BASE}/api${path}`, { ...init, headers }, fetchFn);
	};
	return buildEndpoints(request);
}

export const API_BASE_URL = API_BASE;
