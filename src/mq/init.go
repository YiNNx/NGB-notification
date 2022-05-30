package mq

import (
	"ngb-noti/util"
)

var PgChan = make(chan *util.Notification, 100)
var RedisChan = make(chan *util.Notification, 100)

func init() {
	go ReceiveFromQueue()
}
