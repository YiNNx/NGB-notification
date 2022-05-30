package controller

import (
	"encoding/json"
	"ngb-noti/model"
	"ngb-noti/mq"
	"ngb-noti/util/log"
	"strconv"
)

func HandlePostgres() {
	tx := model.BeginTx()
	defer tx.Close()
	for {
		n := <-mq.PgChan
		notification := &model.Notification{
			Time:      n.Time,
			Uid:       n.Uid,
			Type:      n.Type,
			ContentId: n.ContentId,
			Status:    n.Status,
		}
		log.Logger.Debug("pg receive:", n.Uid)
		if err := model.Insert(notification); err != nil {
			tx.Rollback()
			log.Logger.Error(err)
		}
	}
}

func HandleRedis() {
	for {
		n := <-mq.RedisChan
		log.Logger.Debug("redis receive:", n.Uid)
		data, _ := json.Marshal(n)
		list := []string{string(data)}
		if err := model.RedisPush("notification_"+strconv.Itoa(n.Uid), list); err != nil {
			log.Logger.Error(err)
		}
	}
}
