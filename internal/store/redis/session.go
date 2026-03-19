package redis

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// 生成一个安全的 32 字节随机的 base64 字符串
func generateToken() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("generate random bytes failed: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

type sessionStore struct {
	client *redis.Client
}

func (s *sessionStore) Create(ctx context.Context, userID uuid.UUID, ttl time.Duration) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", fmt.Errorf("generate token failed: %w", err)
	}

	key := "session: " + token
	err = s.client.Set(ctx, key, userID.String(), ttl).Err()
	if err != nil {
		return "", fmt.Errorf("redis set failed: %w", err)
	}

	return token, nil
}

func (s *sessionStore) Get(ctx context.Context, token string) (uuid.UUID, error) {
	key := "session: " + token
	result, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return uuid.Nil, fmt.Errorf("redis get failed: %w", err)
	}

	userID, err := uuid.Parse(result)
	if err != nil {
		return uuid.Nil, fmt.Errorf("uuid parse failed: %w", err)
	}

	return userID, nil
}

func (s *sessionStore) Delete(ctx context.Context, token string) error {
	key := "session: " + token
	err := s.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("redis del failed: %w", err)
	}

	return nil
}
