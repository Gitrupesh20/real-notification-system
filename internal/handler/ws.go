package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Gitrupesh20/real-time-notification-system/cmd/config"
	"github.com/Gitrupesh20/real-time-notification-system/internal/domain"
	"github.com/Gitrupesh20/real-time-notification-system/internal/services"
	"github.com/gorilla/websocket"
)

// checkOrigin check for the allowed origin
func (h *WS) checkOrigin(r *http.Request) bool {
	if h.config.Mode == 0 {
		return true
	}
	//check origin
	origin := r.Header.Get("Origin")
	if origin == "" {
		if h.config.AllowNoOrigin {
			return true
		}
		return false
	}
	//check of https or http
	schema := r.URL.Scheme
	if schema == "" {
		if u, err := url.Parse(origin); err == nil {
			schema = u.Scheme
		} else {
			return false
		}
	}

	if strings.EqualFold(schema, "https") && h.config.AllowWsWithSSL && h.config.IsAllowOrigin(origin) {
		return true
	} else if strings.EqualFold(schema, "http") && !h.config.AllowWsWithSSL && h.config.IsAllowOrigin(origin) {
		return true
	}

	return false
}

type WS struct {
	config  config.Config
	upgrade websocket.Upgrader
	user    *services.UserService
}

func NewWS(config config.Config, userServices *services.UserService) *WS {
	ws := &WS{
		config: config,
		upgrade: websocket.Upgrader{
			HandshakeTimeout:  config.HandShakeTimeout * time.Millisecond,
			ReadBufferSize:    config.ReadBufferSize,
			WriteBufferSize:   config.WriteBufferSize,
			EnableCompression: true,
		},
		user: userServices,
	}

	ws.upgrade.CheckOrigin = ws.checkOrigin
	return ws
}

//var Users sync.Map

func (h *WS) WsHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("req come to WsHandler")
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("userId is required"))
		return
	}

	fmt.Println("req for conn is come userId:", userId)

	rawWsconn, err := h.upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error while upgrading conn", err)
		return
	}
	wsConn := domain.NewUser(userId, rawWsconn)

	// if exits disconnect older
	if _, err := h.user.ConnectUser(wsConn); err != nil {
		//bypass the error as swap store new user
		log.Printf("new user login %v", err)
	}

	go wsConn.WriteData() // start writer worker for each user

	defer func() {
		log.Print("inside defer close fns ws connn")
		if err := h.user.DisconnectUser(userId, wsConn); err != nil {
			log.Printf("error while closing conn err: %v", err)
			rawWsconn.Close() // close forcefully
		}
	}()

	//hold to conn
	for {
		_, _, err = wsConn.Conn.ReadMessage()
		if err != nil {
			log.Println("closing connection ", err)
			break
		}
	}
}
