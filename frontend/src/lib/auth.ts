export async function setBffCookie(value: string): Promise<void> {
	const res = await fetch('/auth/cookie', {
		method: 'POST',
		headers: { 'content-type': 'application/json' },
		body: JSON.stringify({ value })
	});
	if (!res.ok) throw new Error(`Failed to set cookie (${res.status})`);
}

export async function clearBffCookie(): Promise<void> {
	const res = await fetch('/auth/cookie', { method: 'DELETE' });
	if (!res.ok) throw new Error(`Failed to clear cookie (${res.status})`);
}
