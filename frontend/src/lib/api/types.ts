export type User = {
	id: string;
	photoS3Key: string | null;
	firstName: string;
	lastName: string;
	bio: string | null;
	birthdate: string;
	speciality: string | null;
	status: string | null;
	role: 'user' | 'admin';
};

export type Page<T> = {
	data: T[];
	total: number;
	limit: number;
	offset: number;
};

export type EduPlace = {
	id: number;
	universityName: string;
	grade: string;
	level: string | null;
	specialization: string;
	startYear: number;
	endYear: number | null;
	isStudyingNow: boolean;
};

export type EduPlaceRequest = {
	universityId: number;
	grade: string;
	level: string | null;
	specialization: string;
	startYear: number;
	endYear: number | null;
	isStudyingNow: boolean;
};

export type WorkPlace = {
	id: number;
	companyName: string;
	grade: string;
	position: string;
	startYear: number;
	endYear: number | null;
	isWorkingNow: boolean;
};

export type WorkPlaceRequest = {
	companyName: string;
	grade: string;
	position: string;
	startYear: number;
	endYear: number | null;
	isWorkingNow: boolean;
};

export type Social = {
	id: number;
	social: 'vk' | 'telegram' | 'github';
	link: string;
	isPreferred: boolean;
};

export type SocialRequest = {
	social: 'vk' | 'telegram' | 'github';
	link: string;
	isPreferred: boolean;
};

export type Skill = {
	id: number;
	name: string;
};

export type Status = {
	id: number;
	content: string;
};

export type UserUpdateRequest = {
	photoS3Key: string | null;
	bio: string | null;
	speciality: string | null;
	statusId: number | null;
};

export type Company = {
	id: number;
	name: string;
};

export type WorkPosition = {
	id: number;
	name: string;
};

export type University = {
	id: number;
	name: string;
	shortName: string;
};

export type Faq = {
	id: number;
	question: string;
	answer: string;
};

export type ApiError = {
	error: string;
};

export const EDU_GRADES = ['bachelor', 'master', 'specialist', 'postgraduate'] as const;
export const WORK_GRADES = [
	'intern',
	'junior',
	'junior_plus',
	'middle',
	'middle_plus',
	'senior',
	'staff',
	'principal',
	'lead',
	'head'
] as const;
export const SOCIAL_NETWORKS = ['vk', 'telegram', 'github'] as const;
