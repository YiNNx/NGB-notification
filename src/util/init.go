package util

var hub = newHub()

func init() {
	go hub.run()
}
