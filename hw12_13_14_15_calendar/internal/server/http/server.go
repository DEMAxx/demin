package internalhttp

import (
	"context"
	"errors"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	logger     Logger
	app        Application
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type LoggerConf struct {
	Level string
}

type Application interface { // TODO
}

func NewServer(logger Logger, hostAndPort string, app Application) *Server {
	mux := http.NewServeMux()

	return &Server{
		httpServer: &http.Server{
			Addr:    hostAndPort,
			Handler: mux,
		},
		logger: logger,
		app:    app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("Starting HTTP server...")

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("HTTP server ListenAndServe: " + err.Error())
		}
	}()

	<-ctx.Done()
	return s.Stop(ctx)
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping HTTP server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		s.logger.Error("HTTP server Shutdown: " + err.Error())
		return err
	}

	s.logger.Info("HTTP server stopped")

	return nil
}

// TODO
