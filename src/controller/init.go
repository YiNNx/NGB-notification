package controller

var hub = newHub()

func init() {

	go hub.run()
}
