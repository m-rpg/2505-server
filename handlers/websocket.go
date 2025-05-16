package handlers

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

type Client struct {
	conn     *websocket.Conn
	userID   uint
	username string
	send     chan []byte
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

var (
	clients    = make(map[*Client]bool)
	broadcast  = make(chan []byte)
	register   = make(chan *Client)
	unregister = make(chan *Client)
	mutex      sync.Mutex
)

func HandleWebSocket(c *gin.Context) {
	// Get user info from JWT token
	userID := c.GetUint("user_id")
	username := c.GetString("username")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &Client{
		conn:     conn,
		userID:   userID,
		username: username,
		send:     make(chan []byte, 256),
	}

	register <- client

	// Start goroutines for reading and writing
	go readPump(client)
	go writePump(client)
}

func readPump(c *Client) {
	defer func() {
		unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		broadcast <- message
	}
}

func writePump(c *Client) {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func broadcastMessage(message []byte) {
	mutex.Lock()
	defer mutex.Unlock()

	for client := range clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(clients, client)
		}
	}
}

func init() {
	go func() {
		for {
			select {
			case client := <-register:
				mutex.Lock()
				clients[client] = true
				mutex.Unlock()

			case client := <-unregister:
				mutex.Lock()
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
				}
				mutex.Unlock()

			case message := <-broadcast:
				broadcastMessage(message)
			}
		}
	}()
}
