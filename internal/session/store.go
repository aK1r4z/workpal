package session

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Store interface {
	// 创建用户回话，返回 Token
	//
	// TTL (Time To Live) 代表会话持续时间
	Create(ctx context.Context, userID uuid.UUID, ttl time.Duration) (string, error)

	Get(ctx context.Context, token string) (uuid.UUID, error)

	Delete(ctx context.Context, token string) error
}
