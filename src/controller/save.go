package controller

import (
	"encoding/json"
	"ngb-noti/model"
	"ngb-noti/mq"
	"ngb-noti/util/log"
	"strconv"
)

func savePg() {
	for {
		n := <-mq.PgChan
		notification := &model.Notification{
			Time:      n.Time,
			Uid:       n.Uid,
			Type:      n.Type,
			ContentId: n.ContentId,
			Status:    n.Status,
		}
		if err := model.Insert(notification); err != nil {
			log.Logger.Error(err)
		}
	}
}

func saveRedis() {
	for {
		n := <-mq.RedisChan
		data, _ := json.Marshal(n)
		list := []string{string(data)}
		if err := model.RedisPush("notification_"+strconv.Itoa(n.Uid), list); err != nil {
			log.Logger.Error(err)
		}
	}
}
