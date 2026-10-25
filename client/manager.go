package client

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Client wraps a websocket connection with a send channel.
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// Manager maintains active clients and broadcasts messages.
type Manager struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mu         sync.RWMutex
}

// NewManager creates a new Manager instance.
func NewManager() *Manager {
	return &Manager{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

// Run starts the manager's main loop.
func (m *Manager) Run() {
	for {
		select {
		case client := <-m.register:
			m.mu.Lock()
			m.clients[client] = true
			m.mu.Unlock()
			logrus.Infof("Client registered: %v", client.conn.RemoteAddr())
		case client := <-m.unregister:
			m.mu.Lock()
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				close(client.send)
				_ = client.conn.Close()
				logrus.Infof("Client unregistered: %v", client.conn.RemoteAddr())
			}
			m.mu.Unlock()
		case message := <-m.broadcast:
			m.mu.RLock()
			for client := range m.clients {
				select {
				case client.send <- message:
				default:
					// If the client send buffer is full, drop the client.
					close(client.send)
					delete(m.clients, client)
				}
			}
			m.mu.RUnlock()
		}
	}
}

// ServeClient starts read/write pumps for a new client.
func (m *Manager) ServeClient(conn *websocket.Conn) {
	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}
	m.register <- client

	// Start goroutines for read and write.
	go client.readPump(m)
	go client.writePump()
}

// readPump reads messages from the websocket and forwards them to broadcast.
func (c *Client) readPump(m *Manager) {
	defer func() {
		m.unregister <- c
	}()
	c.conn.SetReadLimit(5120)
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.WithError(err).Warn("Unexpected websocket close")
			}
			break
		}
		m.broadcast <- message
	}
}

// writePump writes messages from the send channel to the websocket.
func (c *Client) writePump() {
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			logrus.WithError(err).Warn("Failed to write message to client")
			break
		}
	}
}