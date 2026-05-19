import { env } from '$env/dynamic/private';
import type {
	Company,
	EduPlace,
	EduPlaceRequest,
	Faq,
	Page,
	Skill,
	Social,
	SocialRequest,
	Status,
	University,
	User,
	UserUpdateRequest,
	WorkPlace,
	WorkPlaceRequest,
	WorkPosition
} from '$lib/api/types';

const API_BASE = env.API_BASE_URL ?? 'http://127.0.0.1:8080';

export class ApiError extends Error {
	status: number;
	body: unknown;
	constructor(status: number, body: unknown, message: string) {
		super(message);
		this.status = status;
		this.body = body;
	}
}

export type FetchFn = typeof globalThis.fetch;

export type ServerApi = ReturnType<typeof createServerApi>;

export function createServerApi(bffCookie: string | null, fetchFn: FetchFn) {
	async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
		const headers = new Headers(init.headers);
		if (init.body && !headers.has('content-type')) {
			headers.set('content-type', 'application/json');
		}
		if (bffCookie) {
			headers.set('cookie', `bff.cookie=${bffCookie}`);
		}
		const res = await fetchFn(`${API_BASE}${path}`, { ...init, headers });
		if (res.status === 204) return undefined as T;

		const text = await res.text();
		let body: unknown = null;
		if (text) {
			try {
				body = JSON.parse(text);
			} catch {
				body = { error: text };
			}
		}

		if (!res.ok) {
			const msg =
				body && typeof body === 'object' && body !== null && 'error' in body
					? String((body as { error: unknown }).error)
					: `HTTP ${res.status}`;
			throw new ApiError(res.status, body, msg);
		}
		return body as T;
	}

	return {
		raw: request,
		students: {
			list: (limit = 20, offset = 0) =>
				request<Page<User>>(`/api/students/?limit=${limit}&offset=${offset}`),
			get: (id: string) => request<User>(`/api/students/${id}`),
			eduPlaces: (id: string) => request<EduPlace[]>(`/api/students/${id}/edu-places`),
			workPlaces: (id: string) => request<WorkPlace[]>(`/api/students/${id}/work-places`),
			socials: (id: string) => request<Social[]>(`/api/students/${id}/socials`),
			keySkills: (id: string) => request<Skill[]>(`/api/students/${id}/key-skills`),
			softSkills: (id: string) => request<Skill[]>(`/api/students/${id}/soft-skills`)
		},
		me: {
			get: () => request<User>(`/api/me/`),
			update: (body: UserUpdateRequest) =>
				request<User>(`/api/me/`, { method: 'PUT', body: JSON.stringify(body) }),
			eduPlaces: {
				list: () => request<EduPlace[]>(`/api/me/edu-places`),
				create: (body: EduPlaceRequest) =>
					request<EduPlace>(`/api/me/edu-places`, { method: 'POST', body: JSON.stringify(body) }),
				update: (id: number, body: EduPlaceRequest) =>
					request<EduPlace>(`/api/me/edu-places/${id}`, {
						method: 'PUT',
						body: JSON.stringify(body)
					}),
				remove: (id: number) =>
					request<void>(`/api/me/edu-places/${id}`, { method: 'DELETE' })
			},
			workPlaces: {
				list: () => request<WorkPlace[]>(`/api/me/work-places`),
				create: (body: WorkPlaceRequest) =>
					request<WorkPlace>(`/api/me/work-places`, { method: 'POST', body: JSON.stringify(body) }),
				update: (id: number, body: WorkPlaceRequest) =>
					request<WorkPlace>(`/api/me/work-places/${id}`, {
						method: 'PUT',
						body: JSON.stringify(body)
					}),
				remove: (id: number) =>
					request<void>(`/api/me/work-places/${id}`, { method: 'DELETE' })
			},
			socials: {
				list: () => request<Social[]>(`/api/me/socials`),
				create: (body: SocialRequest) =>
					request<Social>(`/api/me/socials`, { method: 'POST', body: JSON.stringify(body) }),
				update: (id: number, body: SocialRequest) =>
					request<Social>(`/api/me/socials/${id}`, { method: 'PUT', body: JSON.stringify(body) }),
				remove: (id: number) => request<void>(`/api/me/socials/${id}`, { method: 'DELETE' })
			},
			keySkills: {
				list: () => request<Skill[]>(`/api/me/key-skills`),
				add: (skillId: number) =>
					request<Skill>(`/api/me/key-skills/${skillId}`, { method: 'POST' }),
				remove: (skillId: number) =>
					request<void>(`/api/me/key-skills/${skillId}`, { method: 'DELETE' })
			},
			softSkills: {
				list: () => request<Skill[]>(`/api/me/soft-skills`),
				add: (skillId: number) =>
					request<Skill>(`/api/me/soft-skills/${skillId}`, { method: 'POST' }),
				remove: (skillId: number) =>
					request<void>(`/api/me/soft-skills/${skillId}`, { method: 'DELETE' })
			}
		},
		reference: {
			statuses: () => request<Status[]>(`/api/statuses`),
			keySkills: () => request<Skill[]>(`/api/key-skills`),
			softSkills: () => request<Skill[]>(`/api/soft-skills`),
			companies: () => request<Company[]>(`/api/companies`),
			workPositions: () => request<WorkPosition[]>(`/api/work-positions`),
			universities: () => request<University[]>(`/api/universities`),
			faqs: () => request<Faq[]>(`/api/faqs`),
			socialNetworks: () => request<string[]>(`/api/social-networks`),
			eduGrades: () => request<string[]>(`/api/edu-grades`),
			workGrades: () => request<string[]>(`/api/work-grades`),
			eventCategories: () => request<string[]>(`/api/event-categories`)
		}
	};
}

export const API_BASE_URL = API_BASE;
