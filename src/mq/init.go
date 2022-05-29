package mq

var PgChan = make(chan *Notification, 100)
var RedisChan = make(chan *Notification, 100)

func init() {
	go ReceiveFromQueue()
}
