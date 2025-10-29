package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Gitrupesh20/real-time-notification-system/internal/domain"
	interfaces "github.com/Gitrupesh20/real-time-notification-system/internal/interface"
)

type NotificationService struct {
	mq   interfaces.MessageQ
	user *UserService
}

func NewNotificationService(mq interfaces.MessageQ, user *UserService) *NotificationService {
	return &NotificationService{
		mq:   mq,
		user: user,
	}
}

func (n *NotificationService) PushNotification(ctx context.Context, payload *domain.Message) error {
	if payload == nil {
		return domain.ErrMessageIsNil
	} else if isValid, err := payload.Validate(); !isValid || err != nil {
		return err
	}

	bytePayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error while marshaling message err: %v", err)
	}

	return n.mq.Publish(ctx, bytePayload)
}

func (m *NotificationService) ProcessNotification(ctx context.Context, message *domain.Message) {
	m.user.BroadcastMessageToAll(message)
}
