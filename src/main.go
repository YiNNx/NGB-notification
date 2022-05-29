package main

import (
	"ngb-noti/controller"
	"ngb-noti/model"
	_ "ngb-noti/mq"
	_ "ngb-noti/util"
	_ "ngb-noti/util/log"
)

func main() {
	model.Connect()
	defer model.Close()

	controller.StartWebSocket()
}
