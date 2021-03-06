package club_ws

import (
	"club/dal"
	appError "club/models/error"
	"club/ws/user_ws"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 2624144
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2624144,
	WriteBufferSize: 2624144,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
func (s Subscription) readPump() {
	c := s.conn
	defer func() {
		H.unregister <- s
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		m := message{msg, s.Room}
		H.broadcast <- m
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (s *Subscription) writePump() {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(c *gin.Context) (*int, error) {
	var userId int
	var userModel dal.User

	if value, ok := c.GetQuery("userId"); ok {
		p, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}

		userId = p
	} else {
		userId = 0
		return nil, appError.AppError{Message: "userId is required"}
	}

	if _, ok := user_ws.H.Users[userId]; !ok {
		return nil, appError.AppError{Message: "The user's connection is not exist."}
	}

	clubIdStr := c.Param("clubId")
	w := c.Writer
	r := c.Request
	clubId, err := strconv.Atoi(clubIdStr)
	if err != nil {
		return &userId, err
	}

	theSame, err := userModel.CompareUserAndClub(&userId, &clubId)
	if err != nil {
		return &userId, err
	}
	if !theSame {
		return &userId, appError.AppError{Message: "The user is not in the room."}
	}

	connections := H.rooms[clubId]
	if connections != nil {
		for _, v := range connections {
			if v == userId {
				return &userId, appError.AppError{Message: "Repeat connection."}
			}
		}
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return &userId, err
	}
	con := &connection{send: make(chan []byte, 256), ws: ws}
	s := Subscription{con, clubId, userId}
	H.register <- s
	go s.writePump()
	go s.readPump()

	return nil, nil
}
