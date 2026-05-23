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

	"github.com/YuraSahanovskyi/booking-system/internal/domain"
	"github.com/YuraSahanovskyi/booking-system/internal/handler"
)

func TestAuthHandler_Register(t *testing.T) {
	t.Run("201 Created - Success Registration", func(t *testing.T) {
		expectedID := uuid.New()

		extendedMock := &configurableMockAuth{
			registerFunc: func(ctx context.Context, email, password string) (uuid.UUID, error) {
				assert.Equal(t, "new@example.com", email)
				return expectedID, nil
			},
		}

		h := handler.NewHandler(extendedMock, nil, nil)
		router := h.Init()

		body, _ := json.Marshal(map[string]string{
			"email":    "new@example.com",
			"password": "password123",
		})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Contains(t, rr.Body.String(), expectedID.String())
	})
}

func TestAuthHandler_Login(t *testing.T) {
	t.Run("200 OK - Success Login", func(t *testing.T) {
		expectedToken := "jwt-access-token"
		extendedMock := &configurableMockAuth{
			loginFunc: func(ctx context.Context, email, password string) (string, error) {
				return expectedToken, nil
			},
		}

		h := handler.NewHandler(extendedMock, nil, nil)
		router := h.Init()

		body, _ := json.Marshal(map[string]string{
			"email":    "user@example.com",
			"password": "correct-password",
		})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), expectedToken)
	})

	t.Run("401 Unauthorized - Invalid Credentials", func(t *testing.T) {
		extendedMock := &configurableMockAuth{
			loginFunc: func(ctx context.Context, email, password string) (string, error) {
				return "", domain.ErrInvalidCredentials
			},
		}

		h := handler.NewHandler(extendedMock, nil, nil)
		router := h.Init()

		body, _ := json.Marshal(map[string]string{
			"email":    "user@example.com",
			"password": "wrong-password",
		})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}

type configurableMockAuth struct {
	registerFunc func(ctx context.Context, email, password string) (uuid.UUID, error)
	loginFunc    func(ctx context.Context, email, password string) (string, error)
}

func (m *configurableMockAuth) Register(ctx context.Context, email, password string) (uuid.UUID, error) {
	if m.registerFunc != nil {
		return m.registerFunc(ctx, email, password)
	}
	return uuid.Nil, nil
}
func (m *configurableMockAuth) Login(ctx context.Context, email, password string) (string, error) {
	if m.loginFunc != nil {
		return m.loginFunc(ctx, email, password)
	}
	return "", nil
}
func (m *configurableMockAuth) ParseToken(ctx context.Context, accessToken string) (uuid.UUID, error) {
	return uuid.New(), nil
}
