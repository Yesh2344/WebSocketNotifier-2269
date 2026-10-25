package server

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/websocket-notifier/client"
	"github.com/example/websocket-notifier/config"
	"github.com/gorilla/websocket"
)

// TestWebSocketEcho verifies that a message sent by one client is broadcast to another.
func TestWebSocketEcho(t *testing.T) {
	// Setup configuration with a random endpoint.
	cfg := &config.Config{
		Port:       0, // unused for httptest
		WSEndpoint: "/ws",
	}

	manager := client.NewManager()
	go manager.Run()

	srv := New(cfg, manager)

	ts := httptest.NewServer(http.HandlerFunc(srv.handleWebSocket))
	defer ts.Close()

	// Build websocket URL.
	wsURL := "ws" + ts.URL[len("http"):]+cfg.WSEndpoint
	dialer := websocket.DefaultDialer

	// Connect first client.
	c1, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect client 1: %v", err)
	}
	defer c1.Close()

	// Connect second client.
	c2, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect client 2: %v", err)
	}
	defer c2.Close()

	// Send a message from client 1.
	msg := []byte(`{"type":"test","data":{"msg":"hello"}}`)
	if err := c1.WriteMessage(websocket.TextMessage, msg); err != nil {
		t.Fatalf("Client 1 write error: %v", err)
	}

	// Expect client 2 to receive the same message.
	c2.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, received, err := c2.ReadMessage()
	if err != nil {
		t.Fatalf("Client 2 read error: %v", err)
	}
	if string(received) != string(msg) {
		t.Fatalf("Expected %s, got %s", msg, received)
	}
}