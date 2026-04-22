package api

import (
	"net/http"

	"github.com/ansh-singh/game-engine/internal/engine"
)

type Server struct {
	httpServer *http.Server
	engine     engine.GameEngine
}

func NewServer(addr string, ge engine.GameEngine) *Server {
	mux := http.NewServeMux()
	s := &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		engine: ge,
	}
	mux.Handle("/submit", NewSubmitHandler(ge))
	mux.Handle("/metrics", NewMetricsHandler(ge))
	return s
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Close() error {
	return s.httpServer.Close()
}
