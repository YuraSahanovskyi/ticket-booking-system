package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // ніколи не відправляємо пароль у JSON
	CreatedAt    time.Time `json:"created_at"`
}
