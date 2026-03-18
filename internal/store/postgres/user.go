package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/aK1r4z/workpal/internal/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound      = fmt.Errorf("not found")
	ErrAlreadyExists = fmt.Errorf("user already exists")
)

type userStore struct {
	pool *pgxpool.Pool
}

// 传入的用户名和认证串应已通过安全验证
func (s *userStore) Create(ctx context.Context, name string, auth string) (*user.User, error) {
	// 开启事务，自动回滚
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("open tx failed: %w", err)
	}
	defer func() {
		tx.Rollback(ctx) // 不知道这里自动回滚是否会出错，如果出错请改为手动回滚
	}()

	// 查询语句
	const query = `
		insert into users (
			name, auth,
			nickname, email,
			status, created_at, updated_at, last_login,
			failed_login_attempt, locked_until,
			deleted_at
		)
		values (
			$1, $2,
			$3, $4,
			$5, $6, $7, $8,
			$9, $10,
			$11
		)
		returning id
		;
	`

	// 创建用户，插入到表中，返回标识符
	u := user.New(name, auth)
	err = tx.QueryRow(ctx, query,
		u.Name, u.Auth,
		u.Nickname, u.Email,
		u.Status, u.CreatedAt, u.UpdatedAt, u.LastLogin,
		u.FailedLoginAttempt, u.LockedUntil,
		u.DeletedAt,
	).Scan(
		&u.ID,
	)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == "23505" {
				return nil, ErrAlreadyExists
			}
		}
		return nil, err
	}

	// 提交事务
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx failed: %w", err)
	}
	return u, nil
}

func (s *userStore) Get(ctx context.Context, id uuid.UUID) (*user.User, error) {
	// 查询语句
	const query = `
		select
			id, name, auth,
			nickname, email,
			status, created_at, updated_at, last_login,
			failed_login_attempt, locked_until,
			deleted_at
		from users u
		where u.id = $1
		;
	`

	// 查询并扫描到实体中
	u := &user.User{}
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Name, &u.Auth,
		&u.Nickname, &u.Email,
		&u.Status, &u.CreatedAt, &u.UpdatedAt, &u.LastLogin,
		&u.FailedLoginAttempt, &u.LockedUntil,
		&u.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user %s: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("scan user %s: %w", id, err)
	}

	return u, nil
}

func (s *userStore) GetByName(ctx context.Context, username string) (*user.User, error) {
	// 查询语句
	const query = `
		select
			id, name, auth,
			nickname, email,
			status, created_at, updated_at, last_login,
			failed_login_attempt, locked_until,
			deleted_at
		from users u
		where u.name = $1
		;
	`

	// 查询并扫描到实体中
	u := &user.User{}
	err := s.pool.QueryRow(ctx, query, username).Scan(
		&u.ID, &u.Name, &u.Auth,
		&u.Nickname, &u.Email,
		&u.Status, &u.CreatedAt, &u.UpdatedAt, &u.LastLogin,
		&u.FailedLoginAttempt, &u.LockedUntil,
		&u.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user %s: %w", username, ErrNotFound)
		}
		return nil, fmt.Errorf("scan user %s: %w", username, err)
	}

	return u, nil
}
