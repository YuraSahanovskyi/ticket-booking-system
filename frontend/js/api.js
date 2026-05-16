export const API_BASE = '/api/v1';

export class ApiError extends Error {
    constructor(status, message, errors = null) {
        super(message);
        this.status = status;
        this.errors = errors;
    }
}

export async function apiRequest(endpoint, method = 'GET', body = null, token = null) {
    const headers = { 'Content-Type': 'application/json' };
    if (token) headers['Authorization'] = `Bearer ${token}`;

    const options = { method, headers };
    if (body) options.body = JSON.stringify(body);

    try {
        const response = await fetch(`${API_BASE}${endpoint}`, options);
        if (response.status === 204) return { data: true };

        let data = null;
        if (response.headers.get("content-type")?.includes("application/json")) {
            data = await response.json();
        }

        if (!response.ok) {
            throw new ApiError(response.status, data?.message || `Error ${response.status}`, data?.errors);
        }

        return { data };
    } catch (error) {
        if (error instanceof ApiError) throw error;
        throw new ApiError(0, 'Network error. Please check your connection.');
    }
}
