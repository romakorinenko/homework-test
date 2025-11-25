package grpcinternal

import (
	"log/slog"
	"net"

	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/app"
	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	Address string
	log     app.Logger
	server  *grpc.Server
	app     *app.App
}

func NewServer(host, port string, logger app.Logger, app *app.App) *Server {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			GetRequestInfo(logger),
		),
	)
	reflection.Register(server)

	return &Server{
		Address: net.JoinHostPort(host, port),
		log:     logger,
		app:     app,
		server:  server,
	}
}

func (s *Server) Start() error {
	s.app.GetLogger().Info("start http server", slog.String("addr", s.Address))
	listener, err := net.Listen("tcp", s.Address)
	if err != nil {
		panic(err)
	}

	pb.RegisterEventServiceServer(s.server, &EventService{app: s.app})
	if err := s.server.Serve(listener); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	s.server.GracefulStop()
	s.log.Info("server grpc stopped")
	return nil
}
