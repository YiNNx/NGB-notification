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

// Transaction

var tx *pg.Tx

type Transaction struct {
	Tx    *pg.Tx
	abort bool
}

func BeginTx() *Transaction {
	tx, _ = db.Begin()
	trans := &Transaction{
		Tx:    tx,
		abort: false,
	}
	return trans
}

func (trans *Transaction) Rollback() {
	err := trans.Tx.Rollback()
	if err != nil {
		log.Logger.Error("tx-close failed:" + err.Error())
	}
	trans.abort = true
	tx = nil
}

func (trans *Transaction) Close() {
	if trans.abort == false {
		err := trans.Tx.Commit()
		if err != nil {
			log.Logger.Error("tx-commit failed:" + err.Error())
		}
	}
	err := trans.Tx.Close()
	if err != nil {
		log.Logger.Error("tx-close failed:" + err.Error())
	}
	tx = nil
}

// Some shared model functions

func Insert(m interface{}) error {
	_, err := tx.Model(m).Insert()
	return err
}
