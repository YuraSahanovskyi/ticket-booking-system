import http from 'k6/http';
import { sleep } from 'k6';
import { Trend, Counter } from 'k6/metrics';

const bookingDuration  = new Trend('booking_duration', true);
const bookingSuccesses = new Counter('booking_successes');
const bookingConflicts = new Counter('booking_conflicts');
const bookingErrors    = new Counter('booking_errors');

export const options = {
    scenarios: {
        booking_load: {
            executor: 'ramping-vus',
            startVUs: 10,
            stages: [
                { duration: '10s', target: 50 },
                { duration: '30s', target: 50 },
                { duration: '10s', target: 0 },
            ],
            gracefulRampDown: '5s',
        },
    },
    thresholds: {
        booking_duration: ['p(95)<600'],
    },
};

const BASE_URL = 'http://localhost/api/v1';

export function setup() {
    const users = [];
    for (let i = 0; i < 200; i++) {
        const email = `testuser_${i}_${Date.now()}@example.com`;
        const reg = http.post(
            `${BASE_URL}/auth/register`,
            JSON.stringify({ email, password: 'secret123' }),
            { headers: { 'Content-Type': 'application/json' } },
        );
        if (reg.status !== 201) continue;

        const login = http.post(
            `${BASE_URL}/auth/login`,
            JSON.stringify({ email, password: 'secret123' }),
            { headers: { 'Content-Type': 'application/json' } },
        );
        if (login.status === 200) {
            users.push({ token: login.json('token'), email });
        }
    }
    console.log(`Created ${users.length} test users`);

    const authHeaders = {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${users[0].token}`,
        },
    };

    const eventsRes = http.get(`${BASE_URL}/events`, authHeaders);
    const events = eventsRes.json();
    if (!events || events.length === 0) {
        throw new Error('No events found — check seed.sql');
    }

    const eventId = events[0].id;
    const seatsRes = http.get(`${BASE_URL}/events/${eventId}/seats`, authHeaders);
    const seats = (seatsRes.json().seats || []).map(s => s.id);
    if (seats.length === 0) {
        throw new Error('No seats found for event — check seed.sql');
    }

    console.log(`Event: ${eventId}, seats: ${seats.length} (${seats.join(', ')})`);

    return { users, seats };
}

export default function (data) {
    const { users, seats } = data;

    const user   = users[__VU % users.length];
    const seatId = seats[Math.floor(Math.random() * seats.length)];

    const res = http.post(
        `${BASE_URL}/bookings/`,
        JSON.stringify({ seat_id: seatId }),
        {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${user.token}`,
            },
        },
    );

    bookingDuration.add(res.timings.duration);

    if (res.status === 201)      bookingSuccesses.add(1);
    else if (res.status === 409) bookingConflicts.add(1);
    else                         bookingErrors.add(1);

    sleep(0.05);
}