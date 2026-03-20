package tag

import (
	"time"

	"github.com/google/uuid"
)

type tag struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Name string    `json:"name" db:"name"` // 标签名

	CreatedBy uuid.UUID  `json:"created_by" db:"created_by"` // 创建者用户标识符
	CreatedAt time.Time  `json:"created_at" db:"created_at"` // 创建时间
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"` // 删除时间，为空则未删除
}

func New(userID uuid.UUID, name string) *tag {
	now := time.Now()

	return &tag{
		Name: name,

		CreatedBy: userID,
		CreatedAt: now,
		DeletedAt: nil,
	}
}
