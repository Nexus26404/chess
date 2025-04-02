interface RequestOptions {
    method?: string;
    headers?: Record<string, string>;
    body?: any;
}

class Http {
    private baseUrl: string = '';

    constructor(baseUrl: string = '') {
        this.baseUrl = baseUrl;
    }

    private async request(endpoint: string, options: RequestOptions = {}) {
        const url = `${this.baseUrl}${endpoint}`;
        const headers = {
            'Content-Type': 'application/json',
            ...options.headers
        };

        const response = await fetch(url, {
            ...options,
            headers,
            body: options.body ? JSON.stringify(options.body) : undefined
        });

        const data = await response.json();
        if (!response.ok) {
            throw new Error(data.error || 'Request failed');
        }
        return data;
    }

    async post(endpoint: string, data: any) {
        return this.request(endpoint, {
            method: 'POST',
            body: data
        });
    }

    async get(endpoint: string) {
        return this.request(endpoint, {
            method: 'GET'
        });
    }
}

export const http = new Http();
