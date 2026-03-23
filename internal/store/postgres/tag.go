package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/aK1r4z/workpal/internal/tag"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type tagStore struct {
	conn *pgxpool.Pool
}

func (s *tagStore) Create(ctx context.Context, tag *tag.Tag) error {
	const query = `
		insert into tags (
			name,
			created_by, created_at, deleted_at
		)
		values (
			$1,
			$2, $3, $4
		)
		returning id
		;
	`

	err := s.conn.QueryRow(ctx, query, tag.Name, tag.CreatedBy, tag.CreatedAt, tag.DeletedAt).Scan(&tag.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *tagStore) Get(ctx context.Context, userID uuid.UUID, name string) (*tag.Tag, error) {
	const query = `
		select
			id, name,
			created_by, created_at, deleted_at
		from tags
		where deleted_at is null
		and created_by = $1
		and name = $2
		;
	`

	t := &tag.Tag{}
	err := s.conn.QueryRow(ctx, query, userID, name).Scan(
		&t.ID, &t.Name,
		&t.CreatedBy, &t.CreatedAt, &t.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, tag.ErrNotFound
		}
		return nil, err
	}

	return t, nil
}

func (s *tagStore) List(ctx context.Context, userID uuid.UUID, limit int32, offset int32) ([]tag.Tag, error) {
	const query = `
		select
			id, name,
			created_by, created_at, deleted_at
		from tags
		where deleted_at is null
		and created_by = $1
		order by created_at
		limit $2 offset $3
	`

	rows, err := s.conn.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []tag.Tag{}
	for rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}

		t := tag.Tag{}
		err := rows.Scan(
			&t.ID, &t.Name,
			&t.CreatedBy, &t.CreatedAt, &t.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		tags = append(tags, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (s *tagStore) Delete(ctx context.Context, userID uuid.UUID, name string) error {
	const query = `
		update tags
		set deleted_at = $3
		where deleted_at is null
		and created_by = $1
		and name = $2
		;
	`

	_, err := s.conn.Exec(ctx, query, userID, name, time.Now())
	if err != nil {
		return err
	}

	return nil
}
