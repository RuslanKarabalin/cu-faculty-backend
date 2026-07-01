import { browser } from '$app/environment';

export const BFF_COOKIE_NAME = 'bff.cookie';
const MAX_AGE = 60 * 60 * 24 * 30;

export function hasBffCookie(): boolean {
	return readCookie(BFF_COOKIE_NAME) !== null;
}

export async function setBffCookie(value: string): Promise<void> {
	if (!browser) return;
	const v = value.trim();
	if (!v) throw new Error('value required');
	const secure = location.protocol === 'https:' ? '; secure' : '';
	document.cookie = `${BFF_COOKIE_NAME}=${encodeURIComponent(v)}; path=/; max-age=${MAX_AGE}; samesite=lax${secure}`;
}

export async function clearBffCookie(): Promise<void> {
	if (!browser) return;
	document.cookie = `${BFF_COOKIE_NAME}=; path=/; max-age=0; samesite=lax`;
}

function readCookie(name: string): string | null {
	if (!browser) return null;
	const prefix = `${name}=`;
	for (const part of document.cookie.split('; ')) {
		if (part.startsWith(prefix)) return decodeURIComponent(part.slice(prefix.length));
	}
	return null;
}
