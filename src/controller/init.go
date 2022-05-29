package controller

func init() {
	go HandlePostgres()
	go HandleRedis()
}
