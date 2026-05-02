package domain

import "errors"

var (
	//auth
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user with this email already exists")

	ErrInvalidCredentials = errors.New("invalid credentials")
    ErrInvalidToken       = errors.New("invalid token")
    ErrUnexpectedSigning  = errors.New("unexpected signing method")

	//events
	ErrEventNotFound = errors.New("event not found")

	//bookings
	ErrBookingNotFound   = errors.New("booking not found")
	ErrSeatAlreadyBooked = errors.New("seat is already booked")
	ErrBookingCannotBeCanceled = errors.New("booking cannot be canceled")
)
