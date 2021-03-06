package user_ws

import (
	"club/dal"
	appError "club/models/error"
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
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// connection is an middleman between the websocket connection and the hub.
type Connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
func (s subscription) readPump() {
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

		m := Message{msg, s.userId}
		H.Broadcast <- m
	}
}

// write writes a message with the given message type and payload.
func (c *Connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (s *subscription) writePump() {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
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
func ServeWs(c *gin.Context) error {
	var userModel dal.User

	w := c.Writer
	r := c.Request

	userIdStr := c.Param("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return err
	}

	userExist, err := userModel.CheckUserExist(&userId)
	if err != nil {
		return err
	}
	if !userExist {
		return appError.AppError{Message: "The user does not exist."}
	}

	conn := H.Users[userId]
	if conn != nil {
		return appError.AppError{Message: "Repeat connection."}
	}

	err = userModel.UpdateDisconnectionTime(&userId, nil)
	if err != nil {
		return err
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	con := &Connection{Send: make(chan []byte, 256), ws: ws}
	s := subscription{con, userId}
	H.register <- s
	go s.writePump()
	go s.readPump()

	return nil
}
