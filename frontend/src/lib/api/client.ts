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
} from './types';

const API_BASE = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080';

export class ApiClientError extends Error {
	status: number;
	body: unknown;
	constructor(status: number, body: unknown, message: string) {
		super(message);
		this.status = status;
		this.body = body;
	}
}

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
	const headers = new Headers(init.headers);
	if (init.body && !headers.has('content-type')) {
		headers.set('content-type', 'application/json');
	}
	const res = await fetch(`${API_BASE}${path}`, {
		...init,
		headers,
		credentials: 'include'
	});

	if (res.status === 204) return undefined as T;

	const text = await res.text();
	const body = text ? JSON.parse(text) : null;

	if (!res.ok) {
		const msg =
			body && typeof body === 'object' && 'error' in body
				? String((body as { error: unknown }).error)
				: `HTTP ${res.status}`;
		throw new ApiClientError(res.status, body, msg);
	}
	return body as T;
}

export const api = {
	students: {
		list: (limit = 20, offset = 0) =>
			request<Page<User>>(`/api/students/?limit=${limit}&offset=${offset}`),
		eduPlaces: (id: string) => request<EduPlace[]>(`/api/students/${id}/edu-places`),
		workPlaces: (id: string) => request<WorkPlace[]>(`/api/students/${id}/work-places`),
		socials: (id: string) => request<Social[]>(`/api/students/${id}/socials`),
		keySkills: (id: string) => request<Skill[]>(`/api/students/${id}/key-skills`),
		softSkills: (id: string) => request<Skill[]>(`/api/students/${id}/soft-skills`),
		register: () => request<User>(`/api/students/register`, { method: 'POST' })
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
