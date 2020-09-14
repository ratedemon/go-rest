package httpserver

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/ratedemon/go-rest/api/middleware"

	"github.com/ratedemon/go-rest/api/helper"
	"github.com/ratedemon/go-rest/config"
	"google.golang.org/grpc"

	"github.com/ratedemon/go-rest/api"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
)

// Server is entity for HTTP server
type Server struct {
	ctx    context.Context
	cancel context.CancelFunc
	cfg    *config.Config
	logger log.Logger

	router     *mux.Router
	listener   net.Listener
	httpServer *http.Server

	grpcConn *grpc.ClientConn
}

// NewServer creates new HTTP server
func NewServer(ctx context.Context, cfg *config.Config, logger log.Logger) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.HTTPListenerAddress))
	if err != nil {
		return nil, fmt.Errorf("Failed to create TCP listener: %v", err)
	}

	ctx, cancel := context.WithCancel(ctx)

	httpServer := &http.Server{}

	grpcConn, err := grpc.Dial(
		fmt.Sprintf(":%d", cfg.GRPCListenerAddress),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to dial GRPC server: %v", err)
	}

	router := mux.NewRouter()
	jwtMiddleware := middleware.NewHTTPMiddleware(cfg, logger)
	router = router.PathPrefix("/api").Subrouter()
	router.Use(jwtMiddleware.JWTAuthentication)
	http.Handle("/", router)

	routes := api.InitRoutes(ctx, logger, grpcConn)
	for _, r := range routes {
		router.HandleFunc(r.Path, helper.HandleWrapper(r.HandleFunc)).Methods(r.Method)
		logger.Log("msg", "Handler is registered", "path", r.Path, "method", r.Method)
	}

	return &Server{
		ctx:        ctx,
		cancel:     cancel,
		cfg:        cfg,
		logger:     logger,
		router:     router,
		listener:   listener,
		httpServer: httpServer,
		grpcConn:   grpcConn,
	}, nil
}

func (s *Server) Run() error {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()

	s.logger.Log("msg", "Starting HTTP server")
	return s.httpServer.Serve(s.listener)
}
