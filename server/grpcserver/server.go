package grpcserver

import (
	"context"
	"fmt"
	"net"

	"github.com/ratedemon/go-rest/config"
	"github.com/ratedemon/go-rest/datastore/db"
	"github.com/ratedemon/go-rest/grpcserver/auth"
	"github.com/ratedemon/go-rest/grpcserver/image"
	"github.com/ratedemon/go-rest/grpcserver/profile"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"
)

// GRPCServer is enity for grpc server
type GRPCServer struct {
	cfg      *config.Config
	listener net.Listener
	server   *grpc.Server
	log      log.Logger

	db *db.DB
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

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%d host=postgres sslmode=disable", cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.Port)
	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	server.db = db.NewDB(dbConn)
	server.log.Log("msg", "Successfully connected to db")

	for _, name := range server.initServices() {
		server.log.Log("msg", "Service successfully added", "service", name)
	}

	return server, nil
}

func (s *GRPCServer) initServices() []string {
	services := []string{}
	{
		authService := auth.NewAuthService(s.cfg, s.log, s.db)
		authService.RegisterService(s.server)
		services = append(services, "auth")
	}
	{
		profileService := profile.NewProfileService(s.cfg, s.log, s.db)
		profileService.RegisterService(s.server)
		services = append(services, "profile")
	}
	{
		imageService := image.NewImageService(s.cfg, s.log, s.db)
		imageService.RegisterService(s.server)
		services = append(services, "image")
	}

	return services
}

func (s *GRPCServer) Run() error {
	s.log.Log("msg", "Starting GRPC server")
	return s.server.Serve(s.listener)
}

func (s *GRPCServer) Shutdown() {
	s.server.Stop()
	s.listener.Close()
}
