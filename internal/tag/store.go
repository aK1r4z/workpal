package tag

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrNotFound = fmt.Errorf("not found")
)

type Store interface {
	Create(ctx context.Context, tag *Tag) error

	Get(ctx context.Context, userID uuid.UUID, name string) (*Tag, error)
	List(ctx context.Context, userID uuid.UUID, limit int32, offset int32) ([]Tag, error)

	Delete(ctx context.Context, userID uuid.UUID, name string) error
}
