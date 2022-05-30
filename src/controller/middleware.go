package controller

import (
	"context"
	"github.com/gorilla/websocket"
	"net/http"
	"ngb-noti/util"
	"ngb-noti/util/log"
	"strings"
)

var upgrade = websocket.Upgrader{}

type contextValue map[string]interface{}

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := r.Context().Value("ws").(contextValue)["ws_connection"].(*websocket.Conn)

		token := strings.Replace(r.Header["Authorization"][0], "Bearer ", "", -1)
		claims, err := util.ParseToken(token)
		if err != nil || claims == nil {
			errorMessage(c, err)
			return
		}
		data := contextValue{
			"claims": claims,
		}
		ctx := context.WithValue(r.Context(), "jwt", data)
		// next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func WsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrade.Upgrade(w, r, nil)
		if err != nil {
			log.Logger.Error("upgrade:", err)
			return
		}
		data := contextValue{
			"ws_connection": c,
		}
		ctx := context.WithValue(r.Context(), "ws", data)
		// next handler
		next.ServeHTTP(w, r.WithContext(ctx))
		defer c.Close()
	})
}
