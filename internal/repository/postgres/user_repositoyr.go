package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/YuraSahanovskyi/booking-system/internal/db/sqlc"
	"github.com/YuraSahanovskyi/booking-system/internal/domain"
)

type UserRepository struct {
	q *sqlc.Queries
}

func NewUserRepository(q *sqlc.Queries) *UserRepository {
	return &UserRepository{q: q}
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) (uuid.UUID, error) {
	id, err := r.q.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
	})
	
	if err != nil {
		var pgErr *pgconn.PgError
		// unique violation (23505)
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return uuid.Nil, domain.ErrUserAlreadyExists
		}
		return uuid.Nil, err
	}
	
	return id, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	res, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &domain.User{
		ID:           res.ID,
		Email:        res.Email,
		PasswordHash: res.PasswordHash,
	}, nil
}

func (r *UserRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	return r.q.ExistsUserByID(ctx, id)
}