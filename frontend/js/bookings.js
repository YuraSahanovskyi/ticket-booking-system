import { apiRequest } from './api.js';

export async function getBookings(token) {
    const { data } = await apiRequest('/bookings', 'GET', null, token);
    return data;
}

export async function createBooking(seatId, token) {
    const { data } = await apiRequest('/bookings', 'POST', { seat_id: seatId }, token);
    return data;
}

export async function cancelBooking(id, token) {
    const { data } = await apiRequest(`/bookings/${id}`, 'DELETE', null, token);
    return data;
}
