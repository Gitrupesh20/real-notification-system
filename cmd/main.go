package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Gitrupesh20/real-time-notification-system/cmd/config"
	"github.com/Gitrupesh20/real-time-notification-system/internal/handler"
	"github.com/Gitrupesh20/real-time-notification-system/internal/mq/consumer"
	"github.com/Gitrupesh20/real-time-notification-system/internal/mq/rabbitMq"
	"github.com/Gitrupesh20/real-time-notification-system/internal/repo"
	"github.com/Gitrupesh20/real-time-notification-system/internal/server"
	"github.com/Gitrupesh20/real-time-notification-system/internal/services"
)

func main() {
	conf := config.LoadConfig()

	mq := rabbitMq.NewRabbitMessageQueue(conf)

	userRepo := repo.NewUserRepo()

	userServices := services.NewUserService(userRepo)
	notificaionServices := services.NewNotificationService(mq, userServices)

	wsHandler := handler.NewWS(conf, userServices)
	notificationHandler := handler.NewPushNotificationHandler(notificaionServices)

	routes := server.NewRoute(&conf, wsHandler, notificationHandler)
	handler := routes.RegisterRoute()

	_ = consumer.NewConsumer(&conf, mq, notificaionServices)

	//give some time to goroutine to settel down
	<-time.After(time.Microsecond * 100)

	log.Print("Starting Server at port 8080...")
	err := http.ListenAndServe(":"+conf.Port, handler)
	if err != nil {
		log.Fatal("error while stating server at port 8080", err)
	}

}
