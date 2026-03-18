package auth

import (
	"context"
	"log"

	"github.com/aK1r4z/workpal/internal/user"
)

type Service struct {
	users user.Store
}

// 创建用户认证服务
func NewService(
	users user.Store,
) *Service {
	return &Service {
		users: users,
	}
}

// 用户注册
func (s *Service) Register(ctx context.Context, username string, password string) (error) {
	auth, err := GenerateAuth(Config, password)
	if err != nil {
		return err
	}

	u, err := s.users.Create(ctx, username, auth)
	if err != nil {
		return err
	}

	log.Printf("User %s(%s): Registered", u.Name, u.ID)

	return nil
}

func (s *Service) Login(ctx context.Context, username string, password string) (bool, error) {
	u, err := s.users.GetByName(ctx, username)
	if err != nil {
		return false, err
	}

	ok, err := VerifyPassword(password, u.Auth)
	if err != nil {
		return false, err
	}

	return ok, nil
}
