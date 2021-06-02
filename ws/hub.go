package ws

import (
	"club/model"
	"club/service/club_service"
	"log"
)

type message struct {
	data []byte
	room int
}

type subscription struct {
	conn   *connection
	room   int
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
	register chan subscription

	// Unregister requests from connections.
	unregister chan subscription
}

var H = hub{
	broadcast:  make(chan message),
	register:   make(chan subscription),
	unregister: make(chan subscription),
	rooms:      make(map[int]map[*connection]int),
}

func (h *hub) Run() {
	var clubModel model.Club

	for {
		select {
		case s := <-h.register:
			connections := h.rooms[s.room]
			if connections == nil {
				connections = make(map[*connection]int)
				h.rooms[s.room] = connections
			}

			h.rooms[s.room][s.conn] = s.userId

			papulation := len(h.rooms[s.room])
			clubModel.UpdatePopulation(&s.room, &papulation)
		case s := <-h.unregister:
			connections := h.rooms[s.room]
			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					isOwner, err := clubModel.CheckOwnerByUserId(&s.userId)
					if err != nil {
						log.Printf("error: %v", err.Error())
					}

					err = club_service.Leave(&s.userId)
					if err != nil {
						log.Printf("error: %v", err.Error())
					}

					if isOwner == true {
						for k := range connections {
							delete(connections, k)
							close(k.send)
						}

						delete(h.rooms, s.room)
					} else {
						delete(connections, s.conn)
						close(s.conn.send)

						if len(connections) == 0 {
							delete(h.rooms, s.room)
						}
						papulation := len(h.rooms[s.room])
						clubModel.UpdatePopulation(&s.room, &papulation)
					}
				}
			}
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
