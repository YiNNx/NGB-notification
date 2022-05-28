package main

import (
	"ngb-noti/model"
	"ngb-noti/mq"
)

func main() {
	model.Connect()
	defer model.Close()

	mq.ReceiveFromQueue()
}
