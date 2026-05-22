package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/YuraSahanovskyi/booking-system/internal/domain"
	"github.com/YuraSahanovskyi/booking-system/internal/service"
)

type mockUserRepo struct {
	createFunc     func(ctx context.Context, user *domain.User) (uuid.UUID, error)
	getByEmailFunc func(ctx context.Context, email string) (*domain.User, error)
	existsFunc     func(ctx context.Context, id uuid.UUID) (bool, error) // <-- uuid.UUID тут
}

func (m *mockUserRepo) Create(ctx context.Context, user *domain.User) (uuid.UUID, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, user)
	}
	return uuid.New(), nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.getByEmailFunc != nil {
		return m.getByEmailFunc(ctx, email)
	}
	return &domain.User{}, nil
}

func (m *mockUserRepo) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	if m.existsFunc != nil {
		return m.existsFunc(ctx, id)
	}
	return false, nil
}

func TestAuthService(t *testing.T) {
	signingKey := "test-secret-key-12345"
	tokenTTL := time.Hour

	t.Run("Register - Success", func(t *testing.T) {
		expectedID := uuid.New()
		mockRepo := &mockUserRepo{
			createFunc: func(ctx context.Context, user *domain.User) (uuid.UUID, error) {
				assert.Equal(t, "test@example.com", user.Email)
				assert.NotEmpty(t, user.PasswordHash)
				assert.NotEqual(t, "password123", user.PasswordHash)
				return expectedID, nil
			},
		}

		svc := service.NewAuthService(mockRepo, signingKey, tokenTTL)
		id, err := svc.Register(context.Background(), "test@example.com", "password123")

		assert.NoError(t, err)
		assert.Equal(t, expectedID, id)
	})

	t.Run("Login - Success", func(t *testing.T) {
		userID := uuid.New()
		email := "login@example.com"
		password := "secure-pass"

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		mockRepo := &mockUserRepo{
			getByEmailFunc: func(ctx context.Context, e string) (*domain.User, error) {
				assert.Equal(t, email, e)
				return &domain.User{
					ID:           userID,
					Email:        email,
					PasswordHash: string(hashedPassword),
				}, nil
			},
		}

		svc := service.NewAuthService(mockRepo, signingKey, tokenTTL)
		token, err := svc.Login(context.Background(), email, password)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Login - Invalid Credentials (User Not Found)", func(t *testing.T) {
		mockRepo := &mockUserRepo{
			getByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
				return nil, domain.ErrUserNotFound
			},
		}

		svc := service.NewAuthService(mockRepo, signingKey, tokenTTL)
		token, err := svc.Login(context.Background(), "wrong@example.com", "any-pass")

		assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
		assert.Empty(t, token)
	})

	t.Run("Login - Invalid Credentials (Wrong Password)", func(t *testing.T) {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct-pass"), bcrypt.DefaultCost)

		mockRepo := &mockUserRepo{
			getByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
				return &domain.User{
					PasswordHash: string(hashedPassword),
				}, nil
			},
		}

		svc := service.NewAuthService(mockRepo, signingKey, tokenTTL)
		token, err := svc.Login(context.Background(), "user@example.com", "wrong-pass")

		assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
		assert.Empty(t, token)
	})

	t.Run("ParseToken - Success", func(t *testing.T) {
		userID := uuid.New()
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)

		mockRepo := &mockUserRepo{
			getByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
				return &domain.User{
					ID:           userID,
					PasswordHash: string(hashedPassword),
				}, nil
			},
		}

		svc := service.NewAuthService(mockRepo, signingKey, tokenTTL)

		token, err := svc.Login(context.Background(), "user@example.com", "pass")
		assert.NoError(t, err)

		parsedID, err := svc.ParseToken(context.Background(), token)
		assert.NoError(t, err)
		assert.Equal(t, userID, parsedID)
	})

	t.Run("ParseToken - Invalid Token (Signed with wrong key)", func(t *testing.T) {
		mockRepo := &mockUserRepo{}

		svcWithWrongKey := service.NewAuthService(mockRepo, "completely-different-key", tokenTTL)
		token, _ := svcWithWrongKey.Login(context.Background(), "user@example.com", "pass")

		svc := service.NewAuthService(mockRepo, signingKey, tokenTTL)
		_, err := svc.ParseToken(context.Background(), token)

		assert.Error(t, err)
	})

	t.Run("ParseToken - Expired Token", func(t *testing.T) {
		mockRepo := &mockUserRepo{}

		svcExpired := service.NewAuthService(mockRepo, signingKey, -time.Second)
		token, _ := svcExpired.Login(context.Background(), "user@example.com", "pass")

		svc := service.NewAuthService(mockRepo, signingKey, tokenTTL)
		_, err := svc.ParseToken(context.Background(), token)

		assert.Error(t, err)
	})
}
