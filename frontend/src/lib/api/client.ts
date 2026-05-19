import { ApiError, buildEndpoints, makeRequest, type Request } from './endpoints';

const request: Request = (path, init = {}) => makeRequest(`/proxy${path}`, init, fetch);

export const api = buildEndpoints(request);
export { ApiError };
