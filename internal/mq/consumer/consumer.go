package consumer

import (
	"encoding/json"
	"github.com/Gitrupesh20/real-time-notification-system/cmd/config"
	"github.com/Gitrupesh20/real-time-notification-system/internal/domain"
	"github.com/Gitrupesh20/real-time-notification-system/internal/handler"
	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type CMq interface {
	Consume() error
	Close()
}

//********************************************* Worker Consumer *********************************************

type Consumer struct {
	config     *config.Config
	mqConn     *amqp091.Connection
	mqChannel  *amqp091.Channel
	noOfWorker int
	delivery   <-chan amqp091.Delivery // received only cannel
	messageCh  chan handler.Message
}

func NewConsumer(config *config.Config, mqConn *amqp091.Connection, mqChannel *amqp091.Channel) *Consumer {
	consumer := &Consumer{config: config, mqConn: mqConn, mqChannel: mqChannel, noOfWorker: config.NoOfWorker}

	//for i := 0; i < config.NoOfWorker; i++ {
	go consumer.startWorker()
	//}
	return consumer
}

func (c *Consumer) startWorker() {
	log.Println("start worker job")

	for {
		select {
		case msg := <-c.delivery:
			log.Printf("message recevied %v", msg)

			log.Printf("received message from user %s", msg.Body)
			var msgSchema handler.Message
			err := json.Unmarshal(msg.Body, &msgSchema)
			if err != nil {
				log.Printf("failed to unmarshal message from user %s: %v", string(msg.Body), err)
				continue
			}
			sendMessageToAll(msgSchema)
		case <-time.After(time.Millisecond * 100):
			log.Println("worker stopped for 100ms so that consumer can setup")
			continue
		}
		//select {
		//case msg := <-c.messageCh:
		//	sendMessageToAll(msg)
		//}
	}
	log.Println("end worker job")
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

func (c *Consumer) Consume() error {
	dch, err := c.mqChannel.Consume(c.config.MqQueueName, "", false, false, false, false, nil)
	if err != nil {
		log.Printf("failed to consume message from user %s: %v", c.config.MqQueueName, err)
		return err
	}

	//select {
	//case msg := <-dch:
	//	var msgSchema handler.Message
	//	log.Printf("received message from user %s", msg.Body)
	//	err = json.Unmarshal(msg.Body, &msgSchema)
	//	if err != nil {
	//		log.Printf("failed to unmarshal message from user %s: %v", string(msg.Body), err)
	//	}
	//	c.messageCh <- msgSchema
	//}
	c.delivery = dch
	//for msg := range dch {
	//
	//	c.messageCh <- msgSchema
	//}
	log.Println("end consume message from user")
	return nil
}
