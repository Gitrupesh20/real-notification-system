package server

import (
	"fmt"
	"github.com/Gitrupesh20/real-time-notification-system/cmd/config"
	"github.com/Gitrupesh20/real-time-notification-system/internal/handler"
	"github.com/Gitrupesh20/real-time-notification-system/internal/mq/rabbitMq"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type Routes struct {
	config    config.Config
	MessagesQ *rabbitMq.MessageQueue
}

func NewRoute(config *config.Config, mq *rabbitMq.MessageQueue) *Routes {
	return &Routes{
		config:    *config,
		MessagesQ: mq,
	}
}

func (r *Routes) RegisterRoute() http.Handler {
	h := chi.NewRouter()

	ws := handler.NewWS(r.config)
	notify, err := handler.NewPushNotificationHandler(r.MessagesQ)
	if err != nil {
		log.Printf("failed to register push notification handler: %v", err)
		return nil
	}
	//routes
	h.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
		return
	})

	h.HandleFunc("/ws", ws.WsHandler)

	h.Post("/push_notification", notify.PushNotification)

	return h
}
