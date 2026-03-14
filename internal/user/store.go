package user

import (
	"context"

	"github.com/google/uuid"
)

type Store interface {
	Get(context.Context, uuid.UUID) (*User, error)
}
