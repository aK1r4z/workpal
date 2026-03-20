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
	Create(ctx context.Context, tag *tag) error

	Get(ctx context.Context, userID uuid.UUID, name string) (*tag, error)
	List(ctx context.Context, userID uuid.UUID, offset int32, limit int32) ([]tag, error)

	Delete(ctx context.Context, userID uuid.UUID, name string) error
}
