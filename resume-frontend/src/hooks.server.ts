import type { Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	const backendUrl = process.env.BACKEND_URL || 'http://localhost:8080';

	if (event.url.pathname.startsWith('/api')) {
		const targetUrl = new URL(event.request.url);
		targetUrl.protocol = new URL(backendUrl).protocol;
		targetUrl.host = new URL(backendUrl).host;

		const headers = new Headers(event.request.headers);
		headers.set('Host', targetUrl.host);

		const response = await fetch(targetUrl.toString(), {
			method: event.request.method,
			headers,
			body: event.request.body,
			// @ts-expect-error duplex is required for streaming body in Node
			duplex: 'half'
		});
		return response;
	}

	return resolve(event);
};
