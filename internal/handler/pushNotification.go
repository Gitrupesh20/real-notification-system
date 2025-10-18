package handler

import (
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"time"
)

type PushNotificationHandler struct {
	queueName    string
	producerConn *amqp091.Connection
	channel      *amqp091.Channel
	queue        *amqp091.Queue
}

func NewPushNotificationHandler(prodConn *amqp091.Connection, channel *amqp091.Channel, queue string) (*PushNotificationHandler, error) {
	q, err := channel.QueueDeclare(queue, false, false, false, false, nil)
	if err != nil {
		log.Printf("failed to declare queue %s: %v", queue, err)
		return nil, err
	}

	return &PushNotificationHandler{
		producerConn: prodConn,
		channel:      channel,
		queueName:    queue,
		queue:        &q,
	}, nil
}

func (n *PushNotificationHandler) PushNotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rawMsg Message

	if err := json.NewDecoder(r.Body).Decode(&rawMsg); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"bad request, unsupported message"}`))
		return
	}
	log.Printf("Received push notification from message %v", rawMsg)

	if ok, err := ValidateMessage(&rawMsg); err != nil && !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error bad req\":\"%v\"}", err)))
		return
	}
	rawMsg.Timestamp = time.Now().Unix()

	// now publish msg to MQ
	bytes, err := json.Marshal(&rawMsg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{\"error: internal error\":\"%v\"}", err)))
		return
	}

	err = n.channel.Publish("", n.queue.Name, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        bytes,
	})
	if err != nil {
		log.Printf("Failed to publish msg: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
	time.Sleep(time.Millisecond * 200)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))

}

type Message struct {
	Id        string `json:"id"`
	Type      string `json:"type"`
	Body      string `json:"body"`
	Timestamp int64  `json:"timestamp"`
	UserId    string `json:"userId,omitempty"` // if empty means fan out
}

func ValidateMessage(data *Message) (bool, error) {
	if data == nil {
		return false, fmt.Errorf("data is nil")
	} else if len(data.Id) == 0 {
		return false, fmt.Errorf("id is empty")
	} else if len(data.Type) == 0 {
		return false, fmt.Errorf("type is empty")
	} else if data.Timestamp == 0 {
		return false, fmt.Errorf("timestamp is empty")
	}

	return true, nil
}
