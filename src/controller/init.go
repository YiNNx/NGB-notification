package controller

func ListenAndSave() {
	go savePg()
	go saveRedis()
}
