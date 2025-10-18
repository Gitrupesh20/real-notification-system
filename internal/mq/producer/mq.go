package producer

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

// ClientMq is producer that will produce the msq to mq
type ClientMq interface {
	Sent(message []byte) error
	Close()
}

type RabbitMq struct {
	address   string
	queueName string
	conn      *amqp091.Connection
}

func NewRabbitMq(address string, queueName string) ClientMq {
	client := &RabbitMq{
		address:   address,
		queueName: queueName,
	}

	go client.connect()
	return client
}

func (r *RabbitMq) connect() {

	conn, err := amqp091.Dial(r.address)
	if err != nil {
		log.Print("rabbitmq connect error:", err)
		return
	}
	r.conn = conn
}

func (r *RabbitMq) Sent(message []byte) error {

	ch, err := r.conn.Channel()
	if err != nil {
		log.Print("while posting msg rabbitmq sent error:", err)
		return err
	}
	q, err := ch.QueueDeclare(r.queueName, true, false, false, false, nil)
	if err != nil {
		log.Print("while posting msg rabbitmq queue declare error:", err)
		return err
	}
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	err = ch.PublishWithContext(ctx, "", q.Name, false, false, amqp091.Publishing{
		Body:        message,
		ContentType: "text/plain",
	})
	if err != nil {
		log.Print("while posting msg rabbitmq queue publish error:", err)
		return err
	}

	return nil
}

func (r *RabbitMq) Close() {
	r.conn.Close()
}
