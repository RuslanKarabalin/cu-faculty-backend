declare global {
	namespace App {
		interface Locals {
			bffCookie: string | null;
		}
		interface PageData {
			authed: boolean;
		}
	}
}

export {};
