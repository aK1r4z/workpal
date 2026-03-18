package user

import (
	"context"

	"github.com/google/uuid"
)

type Store interface {
	Create(ctx context.Context, name string, auth string) (*User, error)

	Get(ctx context.Context, id uuid.UUID) (*User, error)
	GetByName(ctx context.Context, username string) (*User, error)
}
