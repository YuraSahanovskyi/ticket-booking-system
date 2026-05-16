import { apiRequest } from './api.js';

export async function login(email, password) {
    const { data } = await apiRequest('/auth/login', 'POST', { email, password });
    return data;
}

export async function register(email, password) {
    const { data } = await apiRequest('/auth/register', 'POST', { email, password });
    return data;
}
