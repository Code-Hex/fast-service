package server

import (
	"context"
	"net"
	"net/http"
)

// Server represents an HTTP server.
type Server struct {
	server *http.Server
}

// NewServer creates new Server.
func New(h http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Handler: h,
		},
	}
}

// Serve starts accept requests from the given listener. If any returns error.
func (s *Server) Serve(ln net.Listener) error {
	// ErrServerClosed is returned by the Server's Serve
	// after a call to Shutdown or Close, we can ignore it.
	if err := s.server.Serve(ln); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown gracefully shutdown the server without interrupting any
// active connections. If any returns error.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// ServeHTTP for represents http.Handler
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}
