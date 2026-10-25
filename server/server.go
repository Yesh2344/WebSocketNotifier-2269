package server

import (
	"fmt"
	"net/http"

	"github.com/example/websocket-notifier/client"
	"github.com/example/websocket-notifier/config"
	"github.com/example/websocket-notifier/logger"
	"github.com/gorilla/websocket"
)

// Server bundles HTTP handling and the client manager.
type Server struct {
	cfg     *config.Config
	manager *client.Manager
	upgrader websocket.Upgrader
}

// New creates a new Server instance.
func New(cfg *config.Config, manager *client.Manager) *Server {
	return &Server{
		cfg:     cfg,
		manager: manager,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			// Allow all origins for demo purposes; tighten in production.
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

// Start runs the HTTP server.
func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc(s.cfg.WSEndpoint, s.handleWebSocket)

	addr := fmt.Sprintf(":%d", s.cfg.Port)
	logger.L().Infof("Starting server on %s", addr)
	return http.ListenAndServe(addr, mux)
}

// handleWebSocket upgrades the connection and hands it to the client manager.
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.L().WithError(err).Error("Failed to upgrade connection")
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	s.manager.ServeClient(conn)
}