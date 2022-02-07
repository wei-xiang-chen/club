package user_ws

import (
	"club/dal"
	"log"
	"time"
)

type Message struct {
	Data []byte
	User int
}

type subscription struct {
	conn   *Connection
	userId int
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	Users map[int]*Connection

	// Inbound messages from the connections.
	Broadcast chan Message

	// Register requests from the connections.
	register chan subscription

	// Unregister requests from connections.
	unregister chan subscription
}

var H = hub{
	Broadcast:  make(chan Message),
	register:   make(chan subscription),
	unregister: make(chan subscription),
	Users:      make(map[int]*Connection),
}

func (h *hub) Run() {

	var userModel dal.User

	for {
		select {
		case s := <-h.register:
			connection := h.Users[s.userId]
			if connection == nil {
				h.Users[s.userId] = s.conn
			}
		case s := <-h.unregister:
			connection := h.Users[s.userId]
			if connection != nil {
				delete(h.Users, s.userId)
				close(s.conn.Send)
			}

			currentTime := time.Now()
			err := userModel.UpdateDisconnectionTime(&s.userId, &currentTime)
			if err != nil {
				log.Printf("error: %v", err)
			}
		case m := <-h.Broadcast:
			connection := h.Users[m.User]
			if connection != nil {
				select {
				case connection.Send <- m.Data:
				default:
					close(connection.Send)
					delete(h.Users, m.User)
				}
			}
		}
	}
}
