package postgres

import (
	"context"

	"github.com/aK1r4z/workpal/internal/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type db struct {
	pool *pgxpool.Pool

	userStore user.Store
}

func New(ctx context.Context, connString string) (*db, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	db := &db{
		pool: pool,

		userStore: &userStore{pool},
	}

	return db, nil
}

func (db *db) Close() {
	db.pool.Close()
}

func (db *db) UserStore() user.Store {
	return db.userStore
}
