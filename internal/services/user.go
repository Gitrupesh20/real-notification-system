package services

import (
	"log"

	"github.com/Gitrupesh20/real-time-notification-system/internal/domain"
	interfaces "github.com/Gitrupesh20/real-time-notification-system/internal/interface"
	"github.com/gorilla/websocket"
)

type UserService struct {
	store interfaces.UserStore
}

func NewUserService(userRepo interfaces.UserStore) *UserService {
	return &UserService{
		store: userRepo,
	}
}

func (u *UserService) ConnectUser(user *domain.User) (*domain.User, error) {
	old, ok := u.store.Swap(user.Id, user)
	if ok {
		if err := old.CloseWsConn(websocket.ClosePolicyViolation, "New Login Detected"); err != nil {
			old.Conn.Close()
			return nil, err
		}
		return old, nil
	}

	return nil, domain.ErrUserNotFound
}

func (u *UserService) DisconnectUser(userId string, user *domain.User) error {
	if userId == "" {
		return domain.ErrInvalidUserID
	}

	if current, ok := u.store.Load(userId); ok {
		if current == user {
			log.Print("both conn same deleting")
			u.store.Delete(userId)
		}
	} else {
		return domain.ErrUserNotFound
	}

	user.Close()

	return nil
}

func (u *UserService) BroadcastMessageToAll(data *domain.Message) {
	u.store.Range(func(key, value interface{}) bool {
		if userId, ok := key.(string); ok {
			data.UserId = userId
		} else {
			log.Printf("Caution: user id is not of type string!!!")
			return true
		}

		u.BroadcastMessage(data)
		return true
	})
}

func (u *UserService) BroadcastMessage(data *domain.Message) {
	user, ok := u.store.Load(data.UserId)
	if ok {
		// send msg to worker
		user.Message <- *data
	}
	log.Printf("user %v is ofline, Note: implement SMTP for offline user", data.UserId)
}
