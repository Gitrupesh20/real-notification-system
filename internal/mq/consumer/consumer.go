package consumer

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/Gitrupesh20/real-time-notification-system/cmd/config"
	"github.com/Gitrupesh20/real-time-notification-system/internal/domain"
	interfaces "github.com/Gitrupesh20/real-time-notification-system/internal/interface"
	"github.com/Gitrupesh20/real-time-notification-system/internal/mq/rabbitMq"
	"github.com/Gitrupesh20/real-time-notification-system/internal/services"
)

type CMq interface {
	Consume() error
	Close()
}

//********************************************* Worker Consumer *********************************************

type Consumer struct {
	messageQ     interfaces.MessageQ
	noOfWorker   int
	notification *services.NotificationService
}

func NewConsumer(config *config.Config, mq *rabbitMq.MessageQueue, n *services.NotificationService) *Consumer {
	consumer := &Consumer{messageQ: mq, noOfWorker: config.NoOfWorker, notification: n}

	//for i := 0; i < config.NoOfWorker; i++ {
	go consumer.startWorkerJob()
	//}
	return consumer
}

func (c *Consumer) startWorkerJob() {
	for i := range 20 {
		go c.startWorker(i)
	}
}

func (c *Consumer) startWorker(id int) {
	log.Println("start worker " + strconv.Itoa(id))
	for {
		if !c.messageQ.IsConsumerReady() {
			log.Println("consumer is not ready, retrying after 100ms")
			time.Sleep(time.Millisecond * 100)
			continue
		}

		deliveryChan, err := c.messageQ.Consume()
		if err != nil {
			log.Printf("error while seting up consumer err %v", err)
			time.Sleep(200 * time.Millisecond)
			continue
		}

		for msg := range deliveryChan {
			var msgSchema domain.Message

			err = json.Unmarshal(msg.Body, &msgSchema)
			if err != nil {
				log.Printf("failed to unmarshal message from user %s: %v", string(msg.Body), err)
				msg.Nack(false, false)
				continue
			}
			t := time.Now()
			c.notification.ProcessNotification(context.Background(), &msgSchema)

			log.Printf("total time taken to sent message %v", time.Now().Sub(t))
			if err = msg.Ack(false); err != nil {
				log.Printf("failed to ack message from user %s: %v", msg.Body, err)
				continue
			}
		}

		time.Sleep(200 * time.Millisecond)
	}
	log.Print("consumer is died", id)

}
