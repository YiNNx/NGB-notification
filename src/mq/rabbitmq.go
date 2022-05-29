package mq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"ngb-noti/config"
	"ngb-noti/util"
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

var mqURL = "amqp://" +
	config.C.Rabbitmq.User + ":" +
	config.C.Rabbitmq.Password + "@" +
	config.C.Rabbitmq.Host + ":" +
	config.C.Rabbitmq.Port + "/"

func ReceiveFromQueue() {
	conn, err := amqp.Dial(mqURL)
	if err != nil {
		log.Logger.Error(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Logger.Error(err)
	}
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
	if err != nil {
		log.Logger.Error(err)
	}

	WsQueue, err := ch.QueueDeclare(
		"ws",  // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Logger.Error(err)
	}
	for i, _ := range config.C.Rabbitmq.WsRoutingKey {
		err = ch.QueueBind(
			WsQueue.Name,                      // queue name
			config.C.Rabbitmq.WsRoutingKey[i], // routing key
			config.C.Rabbitmq.ExchangeName,    // exchange
			false,
			nil)
		if err != nil {
			log.Logger.Error(err)
		}
	}
	wsDelivery, err := ch.Consume(
		WsQueue.Name, // queue
		"",           // consumer
		true,         // auto ack
		false,        // exclusive
		false,        // no local
		false,        // no wait
		nil,          // args
	)
	if err != nil {
		log.Logger.Error(err)
	}

	EmailQueue, err := ch.QueueDeclare(
		"email", // name
		false,   // durable
		false,   // delete when usused
		true,    // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Logger.Error(err)
	}
	for i, _ := range config.C.Rabbitmq.EmailRoutingKey {
		err = ch.QueueBind(
			EmailQueue.Name,                      // queue name
			config.C.Rabbitmq.EmailRoutingKey[i], // routing key
			config.C.Rabbitmq.ExchangeName,       // exchange
			false,
			nil)
		if err != nil {
			log.Logger.Error(err)
		}
	}
	emailDelivery, err := ch.Consume(
		EmailQueue.Name, // queue
		"",              // consumer
		true,            // auto ack
		false,           // exclusive
		false,           // no local
		false,           // no wait
		nil,             // args
	)
	if err != nil {
		log.Logger.Error(err)
	}

	wait := make(chan bool)
	go receiveToWs(wsDelivery)
	go receiveToEmail(emailDelivery)
	<-wait
}

func receiveToWs(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		log.Logger.Printf(" receiveNoTi: %s", d.Body)

		n := &Notification{}
		err := json.Unmarshal(d.Body, n)
		if err != nil {
			log.Logger.Error(err)
		}

		PgChan <- n

		client := util.ConnectClient(n.Uid)
		if client == nil {
			RedisChan <- n
		} else {
			client.Send <- n
		}
	}
}

func receiveToEmail(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		log.Logger.Printf(" receiveNoTi: %s", d.Body)

		n := &Notification{}
		err := json.Unmarshal(d.Body, n)
		if err != nil {
			log.Logger.Error(err)
		}
		//util.EmailPool()
	}
}
