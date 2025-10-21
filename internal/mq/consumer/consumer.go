package consumer

import (
	"encoding/json"
	"github.com/Gitrupesh20/real-time-notification-system/cmd/config"
	"github.com/Gitrupesh20/real-time-notification-system/internal/domain"
	"github.com/Gitrupesh20/real-time-notification-system/internal/handler"
	"github.com/Gitrupesh20/real-time-notification-system/internal/mq/rabbitMq"
	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"strconv"
)

type CMq interface {
	Consume() error
	Close()
}

//********************************************* Worker Consumer *********************************************

type Consumer struct {
	config     *config.Config
	messageQ   *rabbitMq.MessageQueue
	noOfWorker int
	delivery   <-chan amqp091.Delivery // received only cannel
	messageCh  chan handler.Message
}

func NewConsumer(config *config.Config, mq *rabbitMq.MessageQueue) *Consumer {
	consumer := &Consumer{config: config, messageQ: mq, noOfWorker: config.NoOfWorker}

	//for i := 0; i < config.NoOfWorker; i++ {
	go consumer.startWorkerJob()
	//}
	return consumer
}

func (c *Consumer) startWorkerJob() {
	for i := 0; i < c.noOfWorker; i++ {
		go c.startWorker(i)
	}
}

func (c *Consumer) startWorker(id int) {
	log.Println("start worker " + strconv.Itoa(id))
	deliveryChan, err := c.messageQ.Consume()
	if err != nil {
		log.Println("consume error", err)
		return
	}
	for {
		select {
		case msg := <-deliveryChan:
			var msgSchema handler.Message
			log.Printf("received message from user %s", msg.Body)
			err = json.Unmarshal(msg.Body, &msgSchema)
			if err != nil {
				log.Printf("failed to unmarshal message from user %s: %v", string(msg.Body), err)
				continue
			}
			sendMessageToAll(msgSchema)
			if err = msg.Ack(false); err != nil {
				log.Printf("failed to ack message from user %s: %v", msg.Body, err)
				continue
			}
		}
	}
}
func sendMessageToAll(msg handler.Message) {

	log.Printf("sending message to all users")

	domain.User.Range(func(key, value interface{}) bool {
		userId := key.(string)
		wsConn := value.(*websocket.Conn)

		err := wsConn.WriteJSON(msg)
		if err != nil {
			log.Printf("failed to send message to user %s: %v", userId, err)
			return false
		}
		return true
	})
	log.Printf("sended message to all users")
}
