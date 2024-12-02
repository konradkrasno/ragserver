package broker

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	forever chan struct{}
)

type MQBroker struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewMQBroker(queueUrl string) (*MQBroker, error) {
	conn, err := amqp.Dial(queueUrl)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = channel.Qos(1, 0, false)
	if err != nil {
		return nil, err
	}

	return &MQBroker{
		Conn:    conn,
		Channel: channel,
	}, nil
}

func (mq *MQBroker) Close() {
	mq.Channel.Close()
	mq.Conn.Close()
}

func (mq *MQBroker) Publish(queueName string, data []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := mq.Channel.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         data,
		})
	if err != nil {
		log.Println(err)
	} else {
		log.Printf(" [x] sent message")
	}
}

func (mq *MQBroker) Listen(queueName string, process func([]byte) error) {
	msgs, err := mq.Channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
		return
	}

	go mq.handleMessage(msgs, process)

	log.Printf(" [*] waiting for messages...")
	<-forever
}

func (mq *MQBroker) handleMessage(msgs <-chan amqp.Delivery, process func([]byte) error) {
	for d := range msgs {
		log.Println(" [x] received a message")

		err := process(d.Body)
		if err != nil {
			log.Println(err)
			log.Println(" [x] rejecting message")
			nack(d)
		} else {
			ack(d)
		}

		log.Println(" [x] processed a message")
	}
}

func ack(d amqp.Delivery) {
	err := d.Ack(true)
	if err != nil {
		log.Println(err)
	}
}

func nack(d amqp.Delivery) {
	err := d.Ack(false)
	if err != nil {
		log.Println(err)
	}
}
