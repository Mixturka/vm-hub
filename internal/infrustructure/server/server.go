package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
	"vm-hub/internal/config"
	"vm-hub/views/templates"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpServer *http.Server
	router     *chi.Router
}

func NewServer(r *chi.Router) *Server {
	s := &Server{router: r}
	return s
}

func (s *Server) Start(config *config.Config) error {
	slog.Info("Http server listening", "address", config.ListenAddr)
	s.httpServer = &http.Server{Addr: config.ListenAddr, Handler: *s.router}
	err := s.httpServer.ListenAndServe()

	if err != http.ErrServerClosed {
		slog.Error("Http server stopped unexpected")
	}

	return nil
}

func (s *Server) Shutdown() {
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to shutdown http server correctly %s", err.Error()))
		} else {
			s.httpServer = nil
		}
	}
}

func SetupRouter() chi.Router {
	r := chi.NewRouter()
	r.Handle("/styles/*", http.StripPrefix("/styles/", http.FileServer(http.Dir("views/styles"))))
	r.Get("/", templ.Handler(templates.Index()).ServeHTTP)

	return r
}
