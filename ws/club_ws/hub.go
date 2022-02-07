package club_ws

import (
	"club/dal"
	model "club/models"
	"club/ws/user_ws"
	"encoding/json"
	"log"
)

type message struct {
	data []byte
	room int
}

type Subscription struct {
	conn   *connection
	Room   int
	userId int
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	rooms map[int]map[*connection]int

	// Inbound messages from the connections.
	broadcast chan message

	// Register requests from the connections.
	register chan Subscription

	// Unregister requests from connections.
	unregister chan Subscription

	CloseRoom chan Subscription
}

var H = hub{
	broadcast:  make(chan message),
	register:   make(chan Subscription),
	unregister: make(chan Subscription),
	CloseRoom:  make(chan Subscription),
	rooms:      make(map[int]map[*connection]int),
}

func (h *hub) Run() {
	var clubModel dal.Club

	for {
		select {
		case s := <-h.register:
			connections := h.rooms[s.Room]
			if connections == nil {
				connections = make(map[*connection]int)
				h.rooms[s.Room] = connections
			}

			h.rooms[s.Room][s.conn] = s.userId

			papulation := len(h.rooms[s.Room])
			clubModel.UpdatePopulation(&s.Room, &papulation)
		case s := <-h.unregister:
			connections := h.rooms[s.Room]
			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					delete(connections, s.conn)
					close(s.conn.send)

					papulation := len(h.rooms[s.Room])
					clubModel.UpdatePopulation(&s.Room, &papulation)
				}
			}
		case s := <-h.CloseRoom:
			connections := h.rooms[s.Room]

			for k, userId := range connections {
				userConnection := user_ws.H.Users[userId]
				wsMsg := model.WsMsg{Action: "leave", Message: "The owner has left."}
				msgString, err := json.Marshal(wsMsg)
				if err != nil {
					log.Printf("error: %v", err)
					return
				}
				userConnection.Send <- msgString
				delete(connections, k)
				close(k.send)
			}

			delete(h.rooms, s.Room)
		case m := <-h.broadcast:
			connections := h.rooms[m.room]
			for c := range connections {
				select {
				case c.send <- m.data:
				default:
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.rooms, m.room)
					}
				}
			}
		}
	}
}
