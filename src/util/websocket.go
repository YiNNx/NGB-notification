package util

import (
	"github.com/gorilla/websocket"

	"ngb-noti/mq"
	"ngb-noti/util/log"
)

type Hub struct {
	// Registered clients.
	clients map[int]*Client

	// Inbound notifications from the clients.
	broadcast chan *mq.Notification

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *mq.Notification, 100),
		register:   make(chan *Client, 100),
		unregister: make(chan *Client, 100),
		clients:    map[int]*Client{},
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			chatID := client.user
			h.clients[chatID] = client
		case client := <-h.unregister:
			chatID := client.user
			if _, ok := h.clients[chatID]; ok {
				delete(h.clients, chatID)
				close(client.Send)
			}
		case notification := <-h.broadcast:
			chatID := notification.Uid
			client := h.clients[chatID]
			select {
			case client.Send <- notification:
			default:
				close(client.Send)
				delete(h.clients, chatID)
			}

		}
	}
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	// The websocket connection.
	conn *websocket.Conn
	user int
	Send chan *mq.Notification
}

func (c *Client) WriteNotification() {
	defer func() {
		c.conn.Close()
	}()
	for {
		ntf, ok := <-c.Send
		if !ok {
			// The hub closed the channel.
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		err := c.conn.WriteJSON(ntf)
		if err != nil {
			log.Logger.Error(err)
		}
	}
}

func GetClient(user int, ws *websocket.Conn) *Client {
	chatID := user
	if res, ok := hub.clients[chatID]; ok {
		client := &Client{
			conn: ws,
			user: user,
			Send: res.Send,
		}
		return client
	}
	client := &Client{
		conn: ws,
		user: user,
		Send: make(chan *mq.Notification, 100),
	}
	hub.register <- client
	return client
}

func (c *Client) WriteOfflineNotification(offlineNoti []string) {
	defer c.conn.Close()
	for i := range offlineNoti {
		err := c.conn.WriteJSON(offlineNoti[i])
		if err != nil {
			log.Logger.Error(err)
		}
	}
}

func ConnectClient(user int) *Client {
	chatID := user
	if res, ok := hub.clients[chatID]; ok {
		return res
	}
	return nil
}
