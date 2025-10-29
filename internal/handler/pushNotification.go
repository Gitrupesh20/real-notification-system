package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Gitrupesh20/real-time-notification-system/internal/domain"
	"github.com/Gitrupesh20/real-time-notification-system/internal/services"
)

type PushNotificationHandler struct {
	notification *services.NotificationService
}

func NewPushNotificationHandler(notification *services.NotificationService) *PushNotificationHandler {
	return &PushNotificationHandler{
		notification: notification,
	}
}

func (n *PushNotificationHandler) PushNotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rawMsg domain.Message

	if err := json.NewDecoder(r.Body).Decode(&rawMsg); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"bad request, unsupported message"}`))
		return
	}
	log.Printf("Received push notification from message %v", rawMsg)

	rawMsg.Timestamp = time.Now().Unix()

	// now publish msg to MQ

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := n.notification.PushNotification(ctx, &rawMsg)
	if err != nil {
		log.Printf("Failed to publish msg: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))

}
