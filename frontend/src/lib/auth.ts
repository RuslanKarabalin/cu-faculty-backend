const COOKIE_NAME = 'bff.cookie';
const MAX_AGE_SECONDS = 60 * 60 * 24 * 30;

export function getBffCookie(): string | null {
	if (typeof document === 'undefined') return null;
	const match = document.cookie
		.split('; ')
		.find((c) => c.startsWith(`${COOKIE_NAME}=`));
	return match ? match.slice(COOKIE_NAME.length + 1) : null;
}

export function setBffCookie(value: string) {
	if (typeof document === 'undefined') return;
	document.cookie = `${COOKIE_NAME}=${value}; path=/; max-age=${MAX_AGE_SECONDS}; SameSite=Lax`;
}

export function clearBffCookie() {
	if (typeof document === 'undefined') return;
	document.cookie = `${COOKIE_NAME}=; path=/; max-age=0`;
}
