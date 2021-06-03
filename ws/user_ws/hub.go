package user_ws

import (
	"club/model"
	"time"
)

type message struct {
	data []byte
	user int
}

type subscription struct {
	conn   *connection
	userId int
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	users map[int]*connection

	// Inbound messages from the connections.
	broadcast chan message

	// Register requests from the connections.
	register chan subscription

	// Unregister requests from connections.
	unregister chan subscription
}

var H = hub{
	broadcast:  make(chan message),
	register:   make(chan subscription),
	unregister: make(chan subscription),
	users:      make(map[int]*connection),
}

func (h *hub) Run() {
	var userModel model.User

	for {
		select {
		case s := <-h.register:
			connection := h.users[s.userId]
			if connection == nil {
				h.users[s.userId] = s.conn
			}
		case s := <-h.unregister:
			connection := h.users[s.userId]
			if connection != nil {
				time.Sleep(time.Second * 2)

				delete(h.users, s.userId)
				close(s.conn.send)
				userModel.DeleteUserById(&s.userId)
			}
		case m := <-h.broadcast:
			connection := h.users[m.user]
			select {
			case connection.send <- m.data:
			default:
				close(connection.send)
				delete(h.users, m.user)
			}
		}
	}
}
