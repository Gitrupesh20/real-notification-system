package interfaces

import (
	"github.com/Gitrupesh20/real-time-notification-system/internal/domain"
)

type UserStore interface {
	Store(userId string, conn *domain.User)
	Load(userId string) (*domain.User, bool)
	Delete(userId string)
	Swap(userId string, conn *domain.User) (*domain.User, bool)
	Range(fns func(key, val interface{}) bool)
}
