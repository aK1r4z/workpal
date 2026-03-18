package user

import (
	"time"

	"github.com/google/uuid"
)

type UserStatus string

const (
	UserActive UserStatus = "active"
)

// 用户类型
type User struct {
	ID   uuid.UUID `json:"id" db:"id"`     // 用户标识符
	Name string    `json:"name" db:"name"` // 用户名
	Auth string    `json:"auth" db:"auth"` // 用户认证串

	Nickname string  `json:"nickname" db:"nickname"` // 用户昵称
	Email    *string `json:"email" db:"email"`       // 用户邮箱

	Status    UserStatus `json:"status" db:"status"`         // 用户状态
	CreatedAt time.Time  `json:"created_at" db:"created_at"` // 用户注册时间
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"` // 用户上次信息更新时间
	LastLogin time.Time  `json:"last_login" db:"last_login"` // 用户上次登录时间，等于注册时间则用户从未登录

	FailedLoginAttempt uint8     `json:"failed_login_attempt" db:"failed_login_attempt"` // 用户连续登录失败次数
	LockedUntil        time.Time `json:"locked_until" db:"locked_until"`                 // 用户账号解封时间

	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"` // 用户账号注销时间，为空则未注销
}

// 仅使用用户名和认证串创建用户实体，传入的用户名和认证串应已通过安全验证
func New(name string, auth string) *User {
	now := time.Now()

	return &User{
		ID:   uuid.Nil,
		Name: name,
		Auth: auth,

		Nickname: name,
		Email:    nil,

		Status:    UserActive,
		CreatedAt: now,
		UpdatedAt: now,
		LastLogin: now,

		FailedLoginAttempt: 0,
		LockedUntil:        now,

		DeletedAt: nil,
	}
}
