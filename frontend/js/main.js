import { login, register } from './auth.js';
import { getEvents, getEventDetails } from './events.js';
import { getBookings, createBooking, cancelBooking } from './bookings.js';
import { getRemainingTime, formatTime } from './timers.js';
import { ApiError } from './api.js';

document.addEventListener('alpine:init', () => {
    Alpine.data('ticketApp', () => ({
        view: 'events',
        token: localStorage.getItem('token') || null,
        events: [],
        currentEvent: null,
        currentSeats: [],
        bookings: [],
        selectedSeat: null,
        bookingCreated: false,
        authForm: { email: '', password: '' },
        notification: null,
        notificationTimeout: null,
        timers: {},

        init() {
            this.fetchEvents();
            if (this.token) {
                this.fetchBookings();
            }
        },

        get isAuthenticated() {
            return !!this.token;
        },

        showNotification(message, type = 'error') {
            this.notification = { message, type };
            if (this.notificationTimeout) clearTimeout(this.notificationTimeout);
            this.notificationTimeout = setTimeout(() => {
                this.notification = null;
            }, 4000);
        },

        handleApiError(err) {
            if (err instanceof ApiError) {
                if (err.status === 401 && this.isAuthenticated) {
                    this.logout(false);
                    this.showNotification('Session expired or unauthorized. Please log in.', 'warning');
                } else if (err.status === 409) {
                    this.showNotification('Conflict: Seat is already booked or unavailable.', 'warning');
                } else {
                    this.showNotification(err.message, 'error');
                }
            } else {
                this.showNotification('An unexpected error occurred.', 'error');
            }
        },

        async login() {
            try {
                const data = await login(this.authForm.email, this.authForm.password);
                this.token = data.token;
                localStorage.setItem('token', data.token);
                this.authForm = { email: '', password: '' };
                this.showNotification('Logged in successfully', 'success');
                this.navigate('events');
            } catch (err) {
                this.handleApiError(err);
            }
        },

        async register() {
            try {
                await register(this.authForm.email, this.authForm.password);
                this.showNotification('Registration successful. Please log in.', 'success');
                this.navigate('login');
            } catch (err) {
                this.handleApiError(err);
            }
        },

        logout(showNotify = true) {
            this.token = null;
            localStorage.removeItem('token');
            this.bookings = [];
            for (let id in this.timers) clearInterval(this.timers[id].interval);
            this.timers = {};

            if (this.view === 'bookings') {
                this.navigate('events');
            }
            if (showNotify) {
                this.showNotification('Logged out successfully', 'success');
            }
        },

        navigate(newView) {
            this.view = newView;
            window.scrollTo(0, 0);
            if (newView === 'events') {
                this.fetchEvents();
            } else if (newView === 'bookings') {
                this.fetchBookings();
            }
        },

        async fetchEvents() {
            try {
                this.events = await getEvents() || [];
            } catch (err) {
                this.events = [];
                this.handleApiError(err);
            }
        },

        async fetchEventDetails(id) {
            try {
                const data = await getEventDetails(id);

                this.currentEvent = data;
                this.currentSeats = data.seats || [];

                this.selectedSeat = null;
                this.bookingCreated = false;

                this.navigate('event');
            } catch (err) {
                this.handleApiError(err);
            }
        },

        selectSeat(seat) {
            if (!seat.is_available) return;

            if (this.selectedSeat?.id === seat.id) {
                this.selectedSeat = null;
                this.bookingCreated = false;
                return;
            }

            this.selectedSeat = seat;
            this.bookingCreated = false;
        },

        scrollToSeatMap() {
            const el = document.getElementById('seat-selection-panel');

            if (!el) return;

            window.scrollTo({
                top: el.offsetTop - 90,
                behavior: 'smooth'
            });
        },

        async bookSeat() {
            if (!this.token) {
                this.showNotification(
                    "Please login to book tickets",
                    "warning"
                );
                return;
            }

            if (!this.selectedSeat) {
                return;
            }

            try {
                await createBooking(
                    this.selectedSeat.id,
                    this.token
                );

                this.bookingCreated = true;

                this.showNotification(
                    "Seat reserved successfully!",
                    "success"
                );

                await this.fetchBookings();

                const data = await getEventDetails(
                    this.currentEvent.id
                );

                this.currentSeats = data.seats || [];

            } catch (err) {
                this.handleApiError(err);
            }
        },

        async fetchBookings() {
            if (!this.isAuthenticated) return;
            try {
                this.bookings = await getBookings(this.token) || [];
                this.setupTimers();
            } catch (err) {
                this.handleApiError(err);
            }
        },

        async cancelBooking(id) {
            if (confirm('Are you sure you want to cancel this booking?')) {
                try {
                    await cancelBooking(id, this.token);
                    this.showNotification('Booking cancelled successfully', 'success');
                    await this.fetchBookings();
                } catch (err) {
                    this.handleApiError(err);
                }
            }
        },

        setupTimers() {
            for (let id in this.timers) {
                if (this.timers[id].interval) {
                    clearInterval(this.timers[id].interval);
                }
            }

            const newTimers = {};

            this.bookings.forEach(booking => {
                if (booking.status === 'reserved' && booking.expires_at) {
                    const expiresAt = new Date(booking.expires_at).getTime();
                    newTimers[booking.id] = {
                        remaining: getRemainingTime(expiresAt),
                        interval: null
                    };

                    newTimers[booking.id].interval = setInterval(() => {
                        const remaining = getRemainingTime(expiresAt);
                        if (remaining <= 0) {
                            clearInterval(newTimers[booking.id].interval);
                            this.fetchBookings();
                        } else {
                            this.timers[booking.id].remaining = remaining;
                        }
                    }, 1000);
                }
            });

            this.timers = newTimers;
        },

        formatTime
    }));
});
