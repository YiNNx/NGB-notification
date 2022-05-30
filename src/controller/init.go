package controller

func InitHandle() {
	go HandlePostgres()
	go HandleRedis()
}
