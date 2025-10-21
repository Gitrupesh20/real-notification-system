package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Gitrupesh20/real-time-notification-system/internal/mq/rabbitMq"
	"log"
	"net/http"
	"time"
)

type PushNotificationHandler struct {
	queueName string
	messageQ  *rabbitMq.MessageQueue
}

func NewPushNotificationHandler(mq *rabbitMq.MessageQueue) (*PushNotificationHandler, error) {

	return &PushNotificationHandler{
		messageQ: mq,
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = n.messageQ.Publish(ctx, bytes)
	if err != nil {
		log.Printf("Failed to publish msg: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

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
