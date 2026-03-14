package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/aK1r4z/workpal/internal/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = fmt.Errorf("not found")

type userStore struct {
	pool *pgxpool.Pool
}

func (s *userStore) Get(ctx context.Context, id uuid.UUID) (*user.User, error) {
	const query = `
	select
		id, name
	from
		"users" u
	where
		u.id = $1
	;
	`

	u := user.User{}

	err := s.pool.QueryRow(ctx, query, id).Scan(&u.ID, &u.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user %s: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("scan user %s: %w", id, err)
	}

	return &u, nil
}
