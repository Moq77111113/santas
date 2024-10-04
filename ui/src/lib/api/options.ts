export interface RequestOptions extends RequestInit {
	[x: string]: unknown;

	/**
	 * Optional fetch implementation to use.
	 */
	fetch?: typeof fetch;

	/**
	 * Optional headers to include in the request.
	 */
	headers?: Record<string, string>;

	/**
	 *  Body of the request (automatically serialized)
	 */
	body?: BodyInit;

	/**
	 * Query parameters to include in the request.
	 */
	query?: Record<string, unknown>;

	/**
	 * Request key to use for automatic cancellation.
	 */
	requestKey?: `${'GET' | 'POST' | 'PATCH' | 'PUT' | 'DELETE'}:${string}`;
}
