package user_ws

type Message struct {
	Data []byte
	User int
}

type subscription struct {
	conn   *connection
	userId int
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	Users map[int]*connection

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
	Users:      make(map[int]*connection),
}

func (h *hub) Run() {

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
				close(s.conn.send)
			}
		case m := <-h.Broadcast:
			connection := h.Users[m.User]
			if connection != nil {
				select {
				case connection.send <- m.Data:
				default:
					close(connection.send)
					delete(h.Users, m.User)
				}
			}
		}
	}
}
