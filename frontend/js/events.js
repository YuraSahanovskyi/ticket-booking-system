import { apiRequest } from './api.js';

export async function getEvents() {
    const { data } = await apiRequest('/events');
    return data;
}

export async function getEventDetails(id) {
    const { data } = await apiRequest(`/events/${id}/seats`);
    return data;
}
