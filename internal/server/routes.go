package server

import (
	"fmt"
	"net/http"

	"github.com/Gitrupesh20/real-time-notification-system/cmd/config"
	"github.com/Gitrupesh20/real-time-notification-system/internal/handler"
	"github.com/go-chi/chi/v5"
)

type Routes struct {
	config       config.Config
	ws           *handler.WS
	notificatoin *handler.PushNotificationHandler
}

func NewRoute(config *config.Config, ws *handler.WS, notificatoin *handler.PushNotificationHandler) *Routes {
	return &Routes{
		config:       *config,
		ws:           ws,
		notificatoin: notificatoin,
	}
}

func (r *Routes) RegisterRoute() http.Handler {
	h := chi.NewRouter()

	//routes
	h.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
		return
	})

	h.HandleFunc("/ws", r.ws.WsHandler)

	h.Post("/push_notification", r.notificatoin.PushNotification)

	return h
}
