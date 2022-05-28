package mq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"ngb-noti/config"
	"ngb-noti/util/log"
	"time"
)

type Notification struct {
	Time      time.Time
	Uid       int
	Type      int
	ContentId int
	Status    int
}

var mqURL = "amqp://" + config.C.Rabbitmq.User + ":" + config.C.Rabbitmq.Password + "@" + config.C.Rabbitmq.Host + ":" + config.C.Rabbitmq.Port + "/"

func ReceiveFromQueue() {
	conn, err := amqp.Dial(mqURL)
	log.Logger.Error(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	log.Logger.Error(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		config.C.Rabbitmq.ExchangeName, // name
		"direct",                       // type
		true,                           // durable
		false,                          // auto-deleted
		false,                          // internal
		false,                          // no-wait
		nil,                            // arguments
	)
	log.Logger.Error(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	log.Logger.Error(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,                         // queue name
		config.C.Rabbitmq.RoutingKey,   // routing key
		config.C.Rabbitmq.ExchangeName, // exchange
		false,
		nil)
	log.Logger.Error(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	log.Logger.Error(err, "Failed to register a consumer")

	forever := make(chan bool)

	go receive(msgs)

	log.Logger.Info(" Waiting for logs")
	<-forever
}

func receive(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		log.Logger.Printf(" [x] %s", d.Body)

		n := &Notification{}
		err := json.Unmarshal(d.Body, n)
		if err != nil {
			log.Logger.Error(err)
		}
	}
}
