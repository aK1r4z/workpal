package tag

import (
	"context"
	"fmt"
	"unicode"

	"github.com/google/uuid"
)

var (
	ErrInvalidName = fmt.Errorf("invalid name")
)

type service struct {
	tags Store
}

func NewService(
	tags Store,
) *service {
	return &service{
		tags: tags,
	}
}

func (s *service) Create(ctx context.Context, userID uuid.UUID, name string) error {
	// 检测标签名称是否有效
	for _, r := range name {
		if !unicode.IsGraphic(r) {
			return ErrInvalidName
		}
	}

	t := New(userID, name)

	if err := s.tags.Create(ctx, t); err != nil {
		return fmt.Errorf("tags create failed: %w", err)
	}

	return nil
}

func (s *service) Get(ctx context.Context, userID uuid.UUID, name string) (*tag, error) {
	// 检测标签名称是否有效
	for _, r := range name {
		if !unicode.IsGraphic(r) {
			return nil, ErrInvalidName
		}
	}

	t, err := s.tags.Get(ctx, userID, name)
	if err != nil {
		return nil, err
	}

	if t.CreatedBy != userID {
		// you should check and see how the fuck could this error exists.
		return nil, fmt.Errorf("internal/tag/service.go > func (s *service) Get: %w", err)
	}

	return t, nil
}

type ListFilter struct {
	Page  int32
	Limit int32
}

func (s *service) List(ctx context.Context, userID uuid.UUID, filter ListFilter) ([]tag, error) {
	if filter.Limit > 100 {
		// uhh im not sure. should it just be 100?
		filter.Limit = 100
	}

	offset := (filter.Page - 1) * filter.Limit

	return s.tags.List(ctx, userID, offset, filter.Limit)
}

// [FIXME] isn't it too simple? i need to recheck this in the day.
func (s *service) Delete(ctx context.Context, userID uuid.UUID, name string) error {
	return s.tags.Delete(ctx, userID, name)
}
