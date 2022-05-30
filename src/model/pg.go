package model

import (
	"github.com/go-pg/pg/v10"
	"ngb-noti/config"
	"ngb-noti/util/log"
)

var db *pg.DB

// Connect database
func Connect() *pg.DB {
	db = pg.Connect(&pg.Options{
		Addr:     config.C.Postgresql.Host + ":" + config.C.Postgresql.Port,
		User:     config.C.Postgresql.User,
		Password: config.C.Postgresql.Password,
		Database: config.C.Postgresql.Dbname,
	})
	var n int
	if _, err := db.QueryOne(pg.Scan(&n), "SELECT 1"); err != nil {
		log.Logger.Error("Postgresql-connection failed")
	}
	log.Logger.Info("Postgresql connected")
	return db
}

// Close database
func Close() {
	if err := db.Close(); err != nil {
		log.Logger.Panic("Postgresql-close failed")
	}
}

func Insert(m interface{}) error {
	_, err := db.Model(m).Insert()
	return err
}
