package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/YuraSahanovskyi/booking-system/internal/handler"
)

func TestPaymentHandler_Webhook(t *testing.T) {
	orderID := uuid.New()

	t.Run("200 OK - Successful Payment Confirmation", func(t *testing.T) {
		mockBooking := &mockBookingService{
			confirmPaymentFunc: func(ctx context.Context, oID uuid.UUID) error {
				assert.Equal(t, orderID, oID)
				return nil
			},
		}

		h := handler.NewHandler(nil, nil, mockBooking)
		router := h.Init()

		body, _ := json.Marshal(map[string]interface{}{
			"order_id": orderID.String(),
			"status":   "success",
		})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/payments/webhook", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "payment confirmed")
	})

	t.Run("200 OK - Non-success status ignored", func(t *testing.T) {
		mockBooking := &mockBookingService{
			confirmPaymentFunc: func(ctx context.Context, oID uuid.UUID) error {
				t.Fatal("ConfirmPayment не повинен викликатися, якщо статус не success")
				return nil
			},
		}

		h := handler.NewHandler(nil, nil, mockBooking)
		router := h.Init()

		body, _ := json.Marshal(map[string]interface{}{
			"order_id": orderID.String(),
			"status":   "failed",
		})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/payments/webhook", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "non-success status ignored")
	})

	t.Run("400 Bad Request - Invalid Payload", func(t *testing.T) {
		h := handler.NewHandler(nil, nil, nil)
		router := h.Init()

		body, _ := json.Marshal(map[string]string{"random_field": "hello"})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/payments/webhook", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid payment payload")
	})
}