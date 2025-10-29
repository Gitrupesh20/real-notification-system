package repo

import (
	"sync"

	"github.com/Gitrupesh20/real-time-notification-system/internal/domain"
)

type SyncUser struct {
	user sync.Map
}

func NewUserRepo() *SyncUser {
	return &SyncUser{user: sync.Map{}}
}

func (u *SyncUser) Load(userId string) (*domain.User, bool) {
	if val, ok := u.user.Load(userId); !ok {
		return nil, false
	} else if conn, ok := val.(*domain.User); !ok {
		//log error
		return nil, false
	} else {
		return conn, true
	}
}

func (u *SyncUser) Store(userId string, conn *domain.User) {
	u.user.Store(userId, conn)
	return
}

func (u *SyncUser) Swap(userId string, conn *domain.User) (*domain.User, bool) {
	oldVal, isExits := u.user.Swap(userId, conn)

	if OldConn, ok := oldVal.(*domain.User); !ok {
		return nil, false
	} else {
		return OldConn, isExits
	}
}

func (u *SyncUser) Delete(userId string) {
	u.user.Delete(userId)
}

func (u *SyncUser) Range(fns func(key, val interface{}) bool) {
	u.user.Range(fns)
}
