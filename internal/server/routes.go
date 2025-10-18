package server

import (
	"fmt"
	"github.com/Gitrupesh20/real-time-notification-system/cmd/config"
	"github.com/Gitrupesh20/real-time-notification-system/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
)

type Routes struct {
	config    config.Config
	mqConn    *amqp091.Connection
	mqChannel *amqp091.Channel
}

func NewRoute(config *config.Config, mqConn *amqp091.Connection, mqCh *amqp091.Channel) *Routes {
	return &Routes{
		config:    *config,
		mqConn:    mqConn,
		mqChannel: mqCh,
	}
}

func (r *Routes) RegisterRoute() http.Handler {
	h := chi.NewRouter()

	ws := handler.NewWS(r.config)
	notify, err := handler.NewPushNotificationHandler(r.mqConn, r.mqChannel, r.config.MqQueueName)
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
