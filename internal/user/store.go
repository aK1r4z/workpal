package user

import (
	"context"

	"github.com/google/uuid"
)

type Store interface {
	Create(ctx context.Context, u *User) error

	Get(ctx context.Context, id uuid.UUID) (*User, error)
	GetByName(ctx context.Context, username string) (*User, error)

	Delete(ctx context.Context, id uuid.UUID) error
}
