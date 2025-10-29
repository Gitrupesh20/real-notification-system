package rabbitMq

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/Gitrupesh20/real-time-notification-system/cmd/config"
	error2 "github.com/Gitrupesh20/real-time-notification-system/pkg/error"
	"github.com/rabbitmq/amqp091-go"
)

type MessageQueue struct {
	queueName       string
	conn            *amqp091.Connection
	notifyConnClose chan *amqp091.Error
	reConnInterval  time.Duration //conn retry
	reSendInterval  time.Duration
	done            chan bool // status for conn close or not
	consumer
	producer
}

type consumer struct {
	mu            sync.Mutex
	ch            *amqp091.Channel
	q             amqp091.Queue
	notifyChClose chan *amqp091.Error
	isReady       bool
}

type producer struct {
	mu            sync.Mutex
	ch            *amqp091.Channel
	q             amqp091.Queue
	notifyChClose chan *amqp091.Error
	isReady       bool
}

func NewRabbitMessageQueue(config config.Config) *MessageQueue {
	mq := &MessageQueue{
		queueName:      config.MqQueueName,
		done:           make(chan bool, 1),
		reConnInterval: time.Millisecond * 100,
		reSendInterval: time.Millisecond * 100,
	}

	go mq.reconnectConnection(config.MqAddr, mq.queueName)

	return mq
}
func (m *MessageQueue) reconnectConnection(addr, quque string) {
	for {
		_, err := m.connect(addr)
		if err != nil {
			log.Println("Error connecting to RabbitMQ:", err)
			log.Println("Reconnecting to RabbitMQ at", addr)
			select {
			case <-m.done:
				log.Printf("reconnectConnection: connection closed")
				break
			case <-time.After(m.reConnInterval):
				continue
			}
		}
		isShutDown := make(chan bool, 1)

		go m.consumerChannelInit(isShutDown)
		go m.producerChannelInit(isShutDown)

		select {
		case is := <-isShutDown:
			if is {
				log.Printf("Mq shutting down!")
				break
			}
			continue
		}

	}
}

func (m *MessageQueue) connect(addr string) (*amqp091.Connection, error) {
	conn, err := amqp091.Dial(addr)
	if err != nil {
		return nil, error2.NoConnectionError
	}
	log.Println("Connected to RabbitMQ at", addr, err)
	m.conn = conn
	m.notifyConnClose = make(chan *amqp091.Error, 1)
	m.conn.NotifyClose(m.notifyConnClose)
	return conn, nil
}

func (m *MessageQueue) consumerChannelInit(isShutDown chan bool) {

	for {
		m.consumer.mu.Lock()
		m.consumer.isReady = false
		m.consumer.mu.Unlock()

		ch, e, err := m.channelInit()
		if err != nil {
			select {
			case <-m.done: //if conn is shutdown, restart
				log.Printf("consumerChannelInit: Connection closed,")
				isShutDown <- true
				return
			case <-m.notifyConnClose:
				log.Printf("consumerChannelInit: Connection closed, Reconnecting to RabbitMQ")
				isShutDown <- false
				return
			case <-time.After(m.reConnInterval):
				continue
			}
		}
		m.consumer.changeConsumerChannel(ch, *e)
		m.consumer.mu.Lock()
		m.consumer.isReady = true
		m.consumer.mu.Unlock()
		log.Print("consumer setup done!")

		//monitor the channel and conn
		select {
		case <-m.done:
			log.Printf("consumerChannelInit: Connection closed")
			isShutDown <- true
			return
		case <-m.notifyConnClose:
			log.Printf("consumerChannelInit: Connection closed, Reconnecting to RabbitMQ")
			isShutDown <- false
			return
		case <-m.consumer.notifyChClose:
			log.Printf("consumerChannelInit: Connection closed, Reconnecting to RabbitMQ")
			continue
		}
	}
}

func (m *MessageQueue) producerChannelInit(isShutDown chan bool) {

	for {
		m.producer.mu.Lock()
		m.producer.isReady = false
		m.producer.mu.Unlock()

		ch, e, err := m.channelInit()
		if err != nil {
			select {
			case <-m.done:
				log.Printf("producerChannelInit: Connection closed,")
				isShutDown <- true
				return
			case <-m.notifyConnClose:
				log.Printf("producerChannelInit: Connection closed, Reconnecting to RabbitMQ")
				isShutDown <- false
				return
			case <-time.After(m.reSendInterval):
				continue
			}
		}
		m.producer.changeProducerChannel(ch, *e)
		m.producer.mu.Lock()
		m.producer.isReady = true
		m.producer.mu.Unlock()
		log.Printf("setUp!")

		select {
		case <-m.done:
			log.Printf("producerChannelInit: Connection closed")
			isShutDown <- true
			return
		case <-m.notifyConnClose:
			log.Printf("producerChannelInit: Connection closed, Reconnecting to RabbitMQ")
			isShutDown <- false
			return
		case <-m.producer.notifyChClose:
			log.Printf("producerChannelInit: Connection closed, Reconnecting to RabbitMQ")
			continue
		}
	}
}

func (c *consumer) changeConsumerChannel(ch *amqp091.Channel, queue amqp091.Queue) {
	// store conn instance
	c.ch = ch
	c.q = queue
	c.notifyChClose = make(chan *amqp091.Error, 1)
	c.ch.NotifyClose(c.notifyChClose)

}

func (p *producer) changeProducerChannel(ch *amqp091.Channel, queue amqp091.Queue) {
	// store conn instance
	p.ch = ch
	p.q = queue
	p.notifyChClose = make(chan *amqp091.Error, 1)
	p.ch.NotifyClose(p.notifyChClose)

}

// channelInit is core func that create channel and queue
func (m *MessageQueue) channelInit() (*amqp091.Channel, *amqp091.Queue, error) {

	ch, err := m.conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	q, err := ch.QueueDeclare(m.queueName, false, false, false, false, nil)
	if err != nil {
		return nil, nil, err
	}
	return ch, &q, nil
}

func (m *MessageQueue) Publish(ctx context.Context, data []byte) error {

	m.producer.mu.Lock()
	if !m.producer.isReady {
		m.producer.mu.Unlock()
		return errors.New("publish fail: not ready")
	}
	m.producer.mu.Unlock()

	if data == nil {
		log.Printf("Warning Publish: empty data")
		return errors.New("publish empty data")
	}
	for {
		err := m.producer.directPublish(ctx, data)
		if err != nil {
			select {
			case <-m.done:
				log.Printf("Publish: Connection closed")
				return nil
			case <-time.After(m.reSendInterval):
				continue
			}
		}

		return nil
	}
}

func (p *producer) directPublish(ctx context.Context, data []byte) error {
	p.mu.Lock()
	if !p.isReady {
		p.mu.Unlock()
		return error2.NoConnectionError
	}
	p.mu.Unlock()

	err := p.ch.PublishWithContext(ctx, "", p.q.Name, false, false, amqp091.Publishing{
		Body:        data,
		ContentType: "application/json",
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *consumer) Consume() (<-chan amqp091.Delivery, error) {
	c.mu.Lock()
	if !c.isReady {
		c.mu.Unlock()
		return nil, errors.New("consume fail: not ready")
	}
	c.mu.Unlock()

	if err := c.ch.Qos(10, 0, false); err != nil {
		return nil, err
	}

	delivery, err := c.ch.Consume(c.q.Name, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("consume started")
	return delivery, nil
}

func (m *MessageQueue) Close(ctx context.Context) error {
	m.consumer.mu.Lock()
	m.consumer.isReady = false
	defer m.consumer.mu.Unlock()
	m.producer.mu.Lock()
	m.producer.isReady = false
	defer m.producer.mu.Unlock()

	close(m.done)

	m.conn.Close()

	m.consumer.ch.Close()
	m.producer.ch.Close()

	return nil
}

func (m *MessageQueue) IsConsumerReady() bool {
	m.consumer.mu.Lock()
	defer m.consumer.mu.Unlock()
	return m.consumer.isReady
}
