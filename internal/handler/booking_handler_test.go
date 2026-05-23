package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/YuraSahanovskyi/booking-system/internal/domain"
	"github.com/YuraSahanovskyi/booking-system/internal/handler"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type mockAuthService struct {
	parseTokenFunc func(ctx context.Context, token string) (uuid.UUID, error)
}

func (m *mockAuthService) Register(ctx context.Context, email, password string) (uuid.UUID, error) {
	return uuid.Nil, nil
}
func (m *mockAuthService) Login(ctx context.Context, email, password string) (string, error) {
	return "", nil
}
func (m *mockAuthService) ParseToken(ctx context.Context, accessToken string) (uuid.UUID, error) {
	if m.parseTokenFunc != nil {
		return m.parseTokenFunc(ctx, accessToken)
	}
	return uuid.New(), nil
}

type mockBookingService struct {
	bookSeatFunc        func(ctx context.Context, userID, seatID uuid.UUID) (*domain.Booking, error)
	getUserBookingsFunc func(ctx context.Context, userID uuid.UUID) ([]domain.BookingWithDetails, error)
	cancelBookingFunc   func(ctx context.Context, userID, bookingID uuid.UUID) error
	confirmPaymentFunc  func(ctx context.Context, orderID uuid.UUID) error
}

func (m *mockBookingService) BookSeat(ctx context.Context, userID, seatID uuid.UUID) (*domain.Booking, error) {
	if m.bookSeatFunc != nil {
		return m.bookSeatFunc(ctx, userID, seatID)
	}
	return &domain.Booking{}, nil
}
func (m *mockBookingService) GetUserBookings(ctx context.Context, userID uuid.UUID) ([]domain.BookingWithDetails, error) {
	if m.getUserBookingsFunc != nil {
		return m.getUserBookingsFunc(ctx, userID)
	}
	return nil, nil
}
func (m *mockBookingService) CancelBooking(ctx context.Context, userID, bookingID uuid.UUID) error {
	if m.cancelBookingFunc != nil {
		return m.cancelBookingFunc(ctx, userID, bookingID)
	}
	return nil
}
func (m *mockBookingService) ConfirmPayment(ctx context.Context, orderID uuid.UUID) error {
	if m.confirmPaymentFunc != nil {
		return m.confirmPaymentFunc(ctx, orderID)
	}
	return nil
}
func (m *mockBookingService) CleanupExpiredBookings(ctx context.Context) (int64, error) {
	return 0, nil
}


func TestBookingHandler_CreateBooking(t *testing.T) {
	userID := uuid.New()

	mockAuth := &mockAuthService{
		parseTokenFunc: func(ctx context.Context, token string) (uuid.UUID, error) {
			if token == "valid-token" {
				return userID, nil
			}
			return uuid.Nil, errors.New("invalid token")
		},
	}

	t.Run("401 Unauthorized - Empty Auth Header", func(t *testing.T) {
		h := handler.NewHandler(mockAuth, nil, nil)
		router := h.Init()

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/bookings/", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "empty auth header")
	})

	t.Run("400 Bad Request - Invalid Seat UUID String", func(t *testing.T) {
		h := handler.NewHandler(mockAuth, nil, nil)
		router := h.Init()

		body, _ := json.Marshal(map[string]string{"seat_id": "not-a-uuid"})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/bookings/", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer valid-token")
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Validation failed")
		assert.Contains(t, rr.Body.String(), "uuid4")
	})

	t.Run("201 Created - Success Booking", func(t *testing.T) {
		seatID := uuid.New()
		bookingID := uuid.New()

		mockBooking := &mockBookingService{
			bookSeatFunc: func(ctx context.Context, uID, sID uuid.UUID) (*domain.Booking, error) {
				assert.Equal(t, userID, uID)
				assert.Equal(t, seatID, sID)
				return &domain.Booking{
					ID:        bookingID,
					UserID:    uID,
					SeatID:    sID,
					Status:    "reserved",
					ExpiresAt: time.Now().Add(15 * time.Minute),
				}, nil
			},
		}

		h := handler.NewHandler(mockAuth, nil, mockBooking)
		router := h.Init()

		body, _ := json.Marshal(map[string]string{"seat_id": seatID.String()})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/bookings/", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer valid-token")
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Contains(t, rr.Body.String(), bookingID.String())
	})
}

func TestBookingHandler_CancelBooking(t *testing.T) {
	userID := uuid.New()
	bookingID := uuid.New()

	mockAuth := &mockAuthService{
		parseTokenFunc: func(ctx context.Context, token string) (uuid.UUID, error) {
			return userID, nil
		},
	}

	t.Run("204 No Content - Success Cancellation", func(t *testing.T) {
		mockBooking := &mockBookingService{
			cancelBookingFunc: func(ctx context.Context, uID, bID uuid.UUID) error {
				assert.Equal(t, userID, uID)
				assert.Equal(t, bookingID, bID)
				return nil
			},
		}

		h := handler.NewHandler(mockAuth, nil, mockBooking)
		router := h.Init()

		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/bookings/"+bookingID.String(), nil)
		req.Header.Set("Authorization", "Bearer valid-token")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
	})
}

