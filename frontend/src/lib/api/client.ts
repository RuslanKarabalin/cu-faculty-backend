import { ApiError, buildEndpoints, makeRequest, type Request } from './endpoints';

const request: Request = (path, init = {}) =>
	makeRequest(`/api${path}`, { credentials: 'same-origin', ...init }, fetch);

export const api = buildEndpoints(request);
export { ApiError };
