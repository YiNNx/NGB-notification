package main

import (
	"ngb-noti/controller"
	"ngb-noti/model"
	_ "ngb-noti/util"
)

func main() {
	model.Connect()
	defer model.Close()

	controller.InitHandle()
	controller.StartWebSocket()
}
