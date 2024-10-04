import type { RequestOptions } from './options';
import GroupService from './services/group';

export class Client {
	baseUrl: string;

	readonly groups: GroupService;

	private cancelControllers = new Map<string, AbortController>();
	private allowAutomaticCancellation = true;

	constructor(baseUrl: string) {
		this.baseUrl = baseUrl;

		this.groups = new GroupService(this);
	}

	async request<T = unknown>(path: string, options: RequestOptions = {}): Promise<T> {
		const url = `${this.baseUrl}${path}`;

		const fetchImplementation = options.fetch || fetch;

		const response = await fetchImplementation(url, this.serializeRequestOptions(path, options));

		let data = {};
		try {
			data = await response.json();
		} catch {
			// ignore
		}

		if (response.status >= 400) {
			throw new Error('An error occurred');
		}

		return data as T;
	}

	autoCancelRequests(allow: boolean): Client {
		this.allowAutomaticCancellation = !!allow;
		return this;
	}

	private serializeRequestOptions(path: string, options: RequestOptions = {}): RequestOptions {
		const enhancedOptions = { ...options };

		if (!enhancedOptions.headers?.['Content-Type'] && !this.isFormData(enhancedOptions.body)) {
			enhancedOptions.headers = {
				...enhancedOptions.headers,
				'Content-Type': 'application/json'
			};
		}

		if (enhancedOptions.body) {
			enhancedOptions.body = this.serializeBody(enhancedOptions.body);
		}

		if (this.allowAutomaticCancellation && options.requestKey) {
			const key = enhancedOptions.requestKey || (`${options.method || 'GET'}:${path}` as const);
			delete enhancedOptions.requestKey;

			this.cancelRequest(key);

			const controller = new AbortController();
			this.cancelControllers.set(key, controller);
			enhancedOptions.signal = controller.signal;
		}

		return enhancedOptions;
	}

	private isFormData(body: RequestOptions['body']): body is FormData {
		return body instanceof FormData;
	}

	private serializeBody(body: BodyInit): BodyInit {
		if (this.isFormData(body)) {
			return body;
		}

		return JSON.stringify(body);
	}

	private cancelRequest(key: string): void {
		const controller = this.cancelControllers.get(key);
		if (controller) {
			controller.abort();
			this.cancelControllers.delete(key);
		}
	}
}

const api = new Client('http://localhost:3456');

export default api;
