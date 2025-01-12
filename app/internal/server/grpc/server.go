package grpc

import (
	"example/config"
	"example/internal/server/grpc/middleware"
	"example/pkg/observer/logger"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"net"
)

type Server struct {
	db   *sqlx.DB
	grpc *grpc.Server
}

func NewServer(
	db *sqlx.DB,
) *Server {
	return &Server{
		grpc: grpc.NewServer(grpc.ChainUnaryInterceptor(
			middleware.Start,
			recovery.UnaryServerInterceptor([]recovery.Option{
				recovery.WithRecoveryHandler(func(p interface{}) (err error) {
					logger.Log.Errorp("Recovered from panic", "error", p)
					return status.Errorf(codes.Internal, "internal error")
				})}...),
		)),
		db: db,
	}
}

func (s *Server) Run() error {
	//register GRPC servers
	s.MapHandlers()

	cfg := config.GetConfig()

	if cfg.Environment == "dev" {
		reflection.Register(s.grpc)
	}

	l, err := net.Listen("tcp", cfg.Server.Grpc.Host+":"+cfg.Server.Grpc.Port)
	if err != nil {
		return err
	}

	go func() {
		logger.Log.Infoa("GRPC Server started", cfg.Server.Grpc)
		if err := s.grpc.Serve(l); err != nil {
			logger.Log.Errorf("Error starting GRPC server: %s", err)
		}
	}()

	return nil
}

func (s *Server) Shutdown() {
	s.grpc.GracefulStop()
	logger.Log.Infop("GRPC server resolved")
}
