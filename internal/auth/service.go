package auth

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aK1r4z/workpal/internal/session"
	"github.com/aK1r4z/workpal/internal/user"
)

// Regexp
var (
	//  首位字母，后接字母数字下划线，总长度 3-16
	regexpUsername = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{2,15}$`)

	// 可打印 ASCII 字符，总长度 8-64
	regexpPassword = regexp.MustCompile(`^[ -~]{8,64}$`)
)

// Error
var (
	ErrInvalidUsername = fmt.Errorf("invalid username")
	ErrInvalidPassword = fmt.Errorf("invalid password")
)

// 检测用户名和密码是否有效
func IsValidCredential(username string, password string) error {
	if !regexpUsername.MatchString(username) {
		return ErrInvalidUsername
	}
	if !regexpPassword.MatchString(password) {
		return ErrInvalidPassword
	}
	return nil
}

type service struct {
	users    user.Store
	sessions session.Store
}

// 创建用户认证服务
func NewService(
	users user.Store,
	sessions session.Store,
) *service {
	return &service{
		users:    users,
		sessions: sessions,
	}
}

// 用户注册
func (s *service) Register(ctx context.Context, username string, password string) error {
	if err := IsValidCredential(username, password); err != nil {
		return err
	}

	auth, err := GenerateAuth(Config, password)
	if err != nil {
		return err
	}

	u := user.New(username, auth)
	if err := s.users.Create(ctx, u); err != nil {
		return err
	}

	log.Printf("user %s(%s): registered", u.Name, u.ID)

	return nil
}

// 登录，返回 Token
func (s *service) Login(ctx context.Context, username string, password string) (string, error) {
	if err := IsValidCredential(username, password); err != nil {
		return "", err
	}

	u, err := s.users.GetByName(ctx, username)
	if err != nil {
		return "", fmt.Errorf("get user by name failed: %w", err)
	}

	// 验证密码
	if err := VerifyPassword(password, u.Auth); err != nil {
		return "", err
	}

	// 创建会话
	token, err := s.sessions.Create(ctx, u.ID, 24*time.Hour)
	if err != nil {
		return "", fmt.Errorf("session create failed: %w", err)
	}

	return token, nil
}
