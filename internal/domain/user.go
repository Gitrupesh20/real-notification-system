package domain

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type User struct {
	Id      string
	Conn    *websocket.Conn
	Message chan Message
	mu      sync.RWMutex
}

func NewUser(userId string, conn *websocket.Conn) *User {
	u := &User{
		Id:      userId,
		Conn:    conn,
		Message: make(chan Message, 10),
	}

	return u
}

func (u *User) WriteData() {
	for msg := range u.Message {
		err := u.Conn.WriteJSON(msg)
		if err != nil {
			log.Printf("error while sending message error: %v", err)
		}
	}
}

func (u *User) Close() {
	u.mu.Lock()
	defer u.mu.Unlock()
	close(u.Message)
	log.Print("closing channel for user ", u.Id)
	return
}

func (u *User) CloseWsConn(code int, text string) error {
	log.Print("sending closing message")
	errMsg := websocket.FormatCloseMessage(code, text)
	u.mu.Lock()
	_ = u.Conn.WriteControl(websocket.CloseMessage, errMsg, time.Now().Add(time.Second*2))
	u.mu.Unlock()
	u.Conn.Close()
	return nil
}
