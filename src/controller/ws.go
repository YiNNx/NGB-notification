package controller

import (
	"github.com/gorilla/websocket"
	"net/http"
	"ngb-noti/config"
	"ngb-noti/model"
	"ngb-noti/util"
	"ngb-noti/util/log"
	"strconv"
)

func ConnectWs(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value("jwt").(contextValue)["claims"].(*util.JwtUserClaims).Id
	c := r.Context().Value("ws").(contextValue)["ws_connection"].(*websocket.Conn)
	client := util.GetClient(uid, c)

	offlineN := PullOfflineNotification(uid)
	if offlineN != nil {
		client.WriteOfflineNotification(offlineN)
	}

	wait := make(chan bool)
	go client.WriteNotification()
	<-wait
}

func PullOfflineNotification(uid int) []string {
	key := "notification_" + strconv.Itoa(uid)
	offlineNotification, err := model.RedisLRange(key)
	if err != nil {
		log.Logger.Error(err)
		return nil
	}
	if err := model.RedisDelete(key); err != nil {
		log.Logger.Error(err)
		return nil
	}
	return offlineNotification
}

func StartWebSocket() {
	http.Handle("/notification", WsMiddleware(JwtMiddleware(http.HandlerFunc(ConnectWs))))
	log.Logger.Fatal(http.ListenAndServe(config.C.App.Addr, nil))
}
