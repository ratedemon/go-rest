package grpcserver

import (
	"context"
	"fmt"
	"net"

	"github.com/ratedemon/go-rest/config"
	"github.com/ratedemon/go-rest/grpcserver/auth"

	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	cfg      *config.Config
	listener net.Listener
	server   *grpc.Server
	log      log.Logger
}

// NewGRPCServer creates new GRPC server
func NewGRPCServer(ctx context.Context, cfg *config.Config, log log.Logger) (*GRPCServer, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCListenerAddress))
	if err != nil {
		return nil, err
	}

	server := &GRPCServer{
		cfg:      cfg,
		listener: listener,
		server:   grpc.NewServer(),
		log:      log,
	}

	for _, name := range server.initServices() {
		server.log.Log("msg", "Service successfully added", "service", name)
	}

	return server, nil
}

func (s *GRPCServer) initServices() []string {
	services := []string{}
	{
		authService := auth.NewAuthService(s.cfg, s.log)
		authService.RegisterService(s.server)
		services = append(services, "auth")
	}

	return services
}

func (s *GRPCServer) Run() error {
	s.log.Log("msg", "Starting GRPC server")
	return s.server.Serve(s.listener)
}
