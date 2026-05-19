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
	const res = await fetch(`/proxy${path}`, { ...init, headers });

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
		throw new ApiClientError(res.status, body, msg);
	}
	return body as T;
}

export const api = {
	students: {
		list: (limit = 20, offset = 0) =>
			request<Page<User>>(`/students/?limit=${limit}&offset=${offset}`),
		get: (id: string) => request<User>(`/students/${id}`),
		eduPlaces: (id: string) => request<EduPlace[]>(`/students/${id}/edu-places`),
		workPlaces: (id: string) => request<WorkPlace[]>(`/students/${id}/work-places`),
		socials: (id: string) => request<Social[]>(`/students/${id}/socials`),
		keySkills: (id: string) => request<Skill[]>(`/students/${id}/key-skills`),
		softSkills: (id: string) => request<Skill[]>(`/students/${id}/soft-skills`)
	},
	me: {
		get: () => request<User>(`/me/`),
		update: (body: UserUpdateRequest) =>
			request<User>(`/me/`, { method: 'PUT', body: JSON.stringify(body) }),
		eduPlaces: {
			list: () => request<EduPlace[]>(`/me/edu-places`),
			create: (body: EduPlaceRequest) =>
				request<EduPlace>(`/me/edu-places`, { method: 'POST', body: JSON.stringify(body) }),
			update: (id: number, body: EduPlaceRequest) =>
				request<EduPlace>(`/me/edu-places/${id}`, {
					method: 'PUT',
					body: JSON.stringify(body)
				}),
			remove: (id: number) => request<void>(`/me/edu-places/${id}`, { method: 'DELETE' })
		},
		workPlaces: {
			list: () => request<WorkPlace[]>(`/me/work-places`),
			create: (body: WorkPlaceRequest) =>
				request<WorkPlace>(`/me/work-places`, { method: 'POST', body: JSON.stringify(body) }),
			update: (id: number, body: WorkPlaceRequest) =>
				request<WorkPlace>(`/me/work-places/${id}`, {
					method: 'PUT',
					body: JSON.stringify(body)
				}),
			remove: (id: number) => request<void>(`/me/work-places/${id}`, { method: 'DELETE' })
		},
		socials: {
			list: () => request<Social[]>(`/me/socials`),
			create: (body: SocialRequest) =>
				request<Social>(`/me/socials`, { method: 'POST', body: JSON.stringify(body) }),
			update: (id: number, body: SocialRequest) =>
				request<Social>(`/me/socials/${id}`, { method: 'PUT', body: JSON.stringify(body) }),
			remove: (id: number) => request<void>(`/me/socials/${id}`, { method: 'DELETE' })
		},
		keySkills: {
			list: () => request<Skill[]>(`/me/key-skills`),
			add: (skillId: number) => request<Skill>(`/me/key-skills/${skillId}`, { method: 'POST' }),
			remove: (skillId: number) =>
				request<void>(`/me/key-skills/${skillId}`, { method: 'DELETE' })
		},
		softSkills: {
			list: () => request<Skill[]>(`/me/soft-skills`),
			add: (skillId: number) => request<Skill>(`/me/soft-skills/${skillId}`, { method: 'POST' }),
			remove: (skillId: number) =>
				request<void>(`/me/soft-skills/${skillId}`, { method: 'DELETE' })
		}
	},
	reference: {
		statuses: () => request<Status[]>(`/statuses`),
		keySkills: () => request<Skill[]>(`/key-skills`),
		softSkills: () => request<Skill[]>(`/soft-skills`),
		companies: () => request<Company[]>(`/companies`),
		workPositions: () => request<WorkPosition[]>(`/work-positions`),
		universities: () => request<University[]>(`/universities`),
		faqs: () => request<Faq[]>(`/faqs`),
		socialNetworks: () => request<string[]>(`/social-networks`),
		eduGrades: () => request<string[]>(`/edu-grades`),
		workGrades: () => request<string[]>(`/work-grades`),
		eventCategories: () => request<string[]>(`/event-categories`)
	}
};
