package redis

import (
	"context"

	"github.com/aK1r4z/workpal/internal/session"
	"github.com/redis/go-redis/v9"
)

type db struct {
	client *redis.Client

	sessionStore session.Store
}

func New(ctx context.Context, addr string) (*db, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	d := &db{
		client: rdb,

		sessionStore: &sessionStore{rdb},
	}

	return d, nil
}

func (d *db) SessionStore() session.Store {
	return d.sessionStore
}
