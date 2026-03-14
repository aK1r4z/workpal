package store

import "github.com/aK1r4z/workpal/internal/user"

type Provider interface {
	UserStore() user.Store
}
