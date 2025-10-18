package handler

import (
	"fmt"
	"github.com/Gitrupesh20/real-time-notification-system/cmd/config"
	"github.com/Gitrupesh20/real-time-notification-system/internal/domain"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
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
	Users   sync.Map
}

func NewWS(config config.Config) *WS {
	ws := &WS{
		config: config,
		upgrade: websocket.Upgrader{
			HandshakeTimeout:  config.HandShakeTimeout * time.Millisecond,
			ReadBufferSize:    config.ReadBufferSize,
			WriteBufferSize:   config.WriteBufferSize,
			EnableCompression: true,
		},
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

	wsConn, err := h.upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error while upgrading conn", err)
		return
	}

	// if exits disconnect older
	if existing, ok := domain.User.Swap(userId, wsConn); ok {
		oldConn := existing.(*websocket.Conn)
		if err = CloseWsConn(oldConn, websocket.ClosePolicyViolation, "New Login Detected"); err != nil {
			log.Println("error closing user conn", err)
		}
	}

	defer func() {
		if current, ok := domain.User.Load(userId); ok {
			if current.(*websocket.Conn) == wsConn {
				domain.User.Delete(userId)
			}
		}
		wsConn.Close()
	}()

	//hold to conn
	for {
		_, _, err = wsConn.ReadMessage()
		if err != nil {
			log.Println("closing connection ", err)
			break
		}
	}
}

func CloseWsConn(conn *websocket.Conn, code int, text string) error {

	errMsg := websocket.FormatCloseMessage(code, text)

	_ = conn.WriteControl(websocket.CloseMessage, errMsg, time.Now().Add(time.Second*2))

	_ = conn.Close()
	return nil
}
