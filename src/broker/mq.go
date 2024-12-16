package broker

import (
	"context"
	"fmt"
	"github.com/konradkrasno/ragserver/environment"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const MQEndpointFormat = "%s://%s:%s@%s:%s%s"

var (
	forever chan struct{}
)

type MQBroker struct {
	MQEndpoint string
}

func NewMQBroker(envs *environment.Envs) *MQBroker {
	return &MQBroker{
		MQEndpoint: fmt.Sprintf(
			MQEndpointFormat,
			envs.RabbitMQProtocol,
			envs.RabbitMQUsername,
			envs.RabbitMQPassword,
			envs.RabbitMQHost,
			envs.RabbitMQPort,
			envs.RabbitMQVHost,
		),
	}
}

func (mq *MQBroker) Publish(exchangeName, sessionId string, data []byte) error {
	conn, ch, err := mq.connect()
	if err != nil {
		return err
	}
	defer ch.Close()
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		exchangeName,
		getSessionRoutingKey(sessionId),
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         data,
		})
	if err != nil {
		return err
	}

	log.Printf(" [x] sent message")
	return nil
}

func (mq *MQBroker) Listen(exchangeName, sessionId string, process func([]byte) error) error {
	conn, ch, err := mq.connect()
	if err != nil {
		return err
	}
	defer ch.Close()
	defer conn.Close()

	q, err := ch.QueueDeclare(
		fmt.Sprintf("%s-session-queue", sessionId),
		false,
		true,
		false,
		false,
		nil)
	if err != nil {
		return err
	}
	err = ch.QueueBind(q.Name, getSessionRoutingKey(sessionId), exchangeName, false, nil)

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go mq.handleMessage(msgs, process)

	log.Printf(" [*] waiting for messages...")
	<-forever

	return nil
}

func (mq *MQBroker) connect() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(mq.MQEndpoint)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	err = ch.Qos(1, 0, false)
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
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

func getSessionRoutingKey(sessionId string) string {
	return fmt.Sprintf("session.%s", sessionId)
}
