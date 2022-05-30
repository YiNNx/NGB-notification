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

func StartWebSocket() {
	http.Handle("/notification", WsMiddleware(JwtMiddleware(http.HandlerFunc(ConnectWs))))
	log.Logger.Fatal(http.ListenAndServe(config.C.Ws.Addr, nil))
}

func ConnectWs(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value("jwt").(contextValue)["claims"].(*util.JwtUserClaims).Id
	c := r.Context().Value("ws").(contextValue)["ws_connection"].(*websocket.Conn)
	client := util.GetClient(uid, c)

	offlineN, err := PullOfflineNotification(uid)
	if err != nil {
		errorMessage(c, err)
	}
	if offlineN != nil {
		if err := client.WriteOfflineNotification(offlineN); err != nil {
			errorMessage(c, err)
		}
	}

	wait := make(chan bool)
	go client.ReceiveNotification()
	<-wait
}

func errorMessage(c *websocket.Conn, err error) {
	c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()))
}

func PullOfflineNotification(uid int) ([]string, error) {
	key := "notification_" + strconv.Itoa(uid)
	offlineNotification, err := model.RedisLRange(key)
	if err != nil {
		return nil, err
	}
	if err := model.RedisDelete(key); err != nil {
		return nil, err
	}
	return offlineNotification, nil
}
