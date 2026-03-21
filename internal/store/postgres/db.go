package postgres

import (
	"context"

	"github.com/aK1r4z/workpal/internal/tag"
	"github.com/aK1r4z/workpal/internal/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type db struct {
	pool *pgxpool.Pool

	userStore user.Store
	tagStore tag.Store
}

func New(ctx context.Context, connString string) (*db, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	d := &db{
		pool: pool,

		userStore: &userStore{pool},
		tagStore: &tagStore{pool},
	}

	return d, nil
}

func (d *db) Close() {
	d.pool.Close()
}

func (d *db) UserStore() user.Store {
	return d.userStore
}

func (d *db) TagStore() tag.Store {
	return d.tagStore
}
