package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	_defaultShutdownTimeout = 5 * time.Second
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultAddr            = ":3456"
)

type Server struct {
	server          *http.Server
	shutdownTimeout time.Duration
	notifier        chan error
}

func New(h http.Handler, opts ...Option) *Server {
	s := &Server{
		server: &http.Server{
			Handler:      h,
			ReadTimeout:  _defaultReadTimeout,
			WriteTimeout: _defaultWriteTimeout,
			Addr:         _defaultAddr,
		},
		shutdownTimeout: 5 * time.Second,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Run() {
	fmt.Printf("Server is running on %s\n", s.server.Addr)
	go func() {
		s.notifier <- s.server.ListenAndServe()
		close(s.notifier)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notifier
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
