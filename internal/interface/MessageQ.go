package interfaces

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
)

type MessageQ interface {
	Publish(ctx context.Context, data []byte) error
	Consume() (<-chan amqp091.Delivery, error)
	Close(ctx context.Context) error
	IsConsumerReady() bool
}
