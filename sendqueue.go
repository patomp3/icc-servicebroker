package main

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// QueueService struct...
type Send struct {
	URL       string
	QueueName string
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s", msg, err)
	}
}

// Close for
func (q Send) Close() {
	//q.conn.Close()
	//q.ch.Close()
}

// Connect for
func (q Send) Connect() *amqp.Channel {
	conn, err := amqp.Dial(q.URL)
	//defer conn.Close()
	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ")
		return nil
	}

	ch, err := conn.Channel()
	//defer ch.Close()
	if err != nil {
		failOnError(err, "Failed to open a channel")
		return nil
	}

	return ch
}

// SendMessage to Queue
func (q Send) SendMessage(ch *amqp.Channel, userId string, msgId string, msgType string, msgBody string) bool {

	/*conn, err := amqp.Dial(q.URL)
	defer conn.Close()
	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ")
		return false
	}

	ch, err := conn.Channel()
	defer ch.Close()
	if err != nil {
		failOnError(err, "Failed to open a channel")
		return false
	}*/

	/*q, err := ch.QueueDeclare(
		"hello2", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")*/

	body := msgBody
	err := ch.Publish(
		"",          // exchange
		q.QueueName, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			MessageId:    msgId,
			UserId:       userId,
			ContentType:  msgType,
			Body:         []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	if err != nil {
		failOnError(err, "Failed to publish a message")
		return false
	}

	return true
}
