package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"ngb-noti/config"
	"ngb-noti/model"
	"ngb-noti/mq"
	"ngb-noti/util"
	"ngb-noti/util/log"
	"strconv"
	"strings"
)

var addr = config.C.App.Addr

var upgrader = websocket.Upgrader{} // use default options

func NotificationWs(w http.ResponseWriter, r *http.Request) {
	log.Logger.Info(strings.Trim(r.Header["Authorization"][0], "Bearer "))
	claims, err := util.ParseToken(strings.Trim(r.Header["Authorization"][0], "Bearer "))
	if err != nil || claims == nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Logger.Error("upgrade:", err)
		return
	}
	defer c.Close()

	user := claims.Id
	client := util.GetClient(user, c)

	offlineNoti, err := model.RedisPull("notification_" + strconv.Itoa(user))
	if err != nil {
		log.Logger.Error("upgrade:", err)
		return
	}
	if offlineNoti != nil {
		client.WriteOfflineNotification(offlineNoti)
	}

	wait := make(chan bool)
	go client.WriteNotification()
	<-wait
}

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
		if err := model.Insert(notification); err != nil {
			tx.Rollback()
			log.Logger.Error(err)
		}
	}
}

func HandleRedis() {
	for {
		n := <-mq.RedisChan
		data, _ := json.Marshal(n)
		list := []string{string(data)}
		if err := model.RedisPush("notification_"+strconv.Itoa(n.Uid), list); err != nil {
			log.Logger.Error(err)
		}
	}
}

func StartWebSocket() {
	http.HandleFunc("/notification", NotificationWs)
	log.Logger.Fatal(http.ListenAndServe(config.C.App.Addr, nil))
}
