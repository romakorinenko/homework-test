package internalhttp

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/app"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/configs"
)

type Server struct {
	app    *app.App
	server *http.Server
}

func NewServer(config configs.HTTP, app *app.App) *Server {
	return &Server{
		app: app,
		server: &http.Server{
			Addr:         fmt.Sprintf(":%s", config.Port),
			Handler:      loggingMiddleware(helloHandler(), app.GetLogger()),
			WriteTimeout: config.Timeout,
			ReadTimeout:  config.Timeout,
		},
	}
}

func (s *Server) Start() error {
	s.app.GetLogger().Info("start http server", slog.String("addr", s.server.Addr))
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.app.GetLogger().Info("stop http server")
	return s.server.Shutdown(ctx)
}

func helloHandler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /hello", http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		_, _ = rw.Write([]byte("hello world!"))
		rw.WriteHeader(http.StatusOK)
	}))

	return mux
}
