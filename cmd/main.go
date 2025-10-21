package main

import (
	"errors"
	"github.com/Gitrupesh20/real-time-notification-system/cmd/config"
	"github.com/Gitrupesh20/real-time-notification-system/internal/mq/consumer"
	"github.com/Gitrupesh20/real-time-notification-system/internal/mq/rabbitMq"
	"github.com/Gitrupesh20/real-time-notification-system/internal/server"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	conf := config.LoadConfig()

	mq := rabbitMq.NewRabbitMessageQueue(conf)

	<-time.After(1 * time.Second)
	//setup route and producer
	newRoute := server.NewRoute(&conf, mq)
	handler := newRoute.RegisterRoute()
	if handler == nil {
		log.Fatal("handler is nil")
		return
	}
	//setup consumer and worker

	_ = consumer.NewConsumer(&conf, mq)
	//start

	log.Printf("Server Started in port %s....", conf.Port)
	if err := http.ListenAndServe(":"+conf.Port, handler); err != nil {
		log.Fatal(err)
	}
	log.Println("done")

}

func connectRabbitMq(addr string, queueName string) (*amqp091.Connection, error) {
	if addr == "" {
		return nil, errors.New("addr is empty")
	} else if queueName == "" {
		return nil, errors.New("queue name is empty")
	} else if !strings.HasPrefix(addr, "amqp://") {
		return nil, errors.New("amqp url is invalid, it should start with 'amqp://'")
	}

	conn, err := amqp091.Dial(addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
