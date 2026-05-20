import { ApiError, buildEndpoints, makeRequest, type Request } from './endpoints';

// Same-origin: nginx serves the static app and proxies /api/* to the Go backend.
// The bff.cookie set by $lib/auth is sent automatically with same-origin requests.
const request: Request = (path, init = {}) =>
	makeRequest(`/api${path}`, { credentials: 'same-origin', ...init }, fetch);

export const api = buildEndpoints(request);
export { ApiError };
