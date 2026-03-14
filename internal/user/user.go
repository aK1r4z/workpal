package user

import "github.com/google/uuid"

// 用户类型
type User struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Name string    `json:"name" db:"name"` // 用户名
}
