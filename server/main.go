package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	kitlog "github.com/go-kit/kit/log"
	"github.com/ratedemon/go-rest/config"
	"github.com/ratedemon/go-rest/grpcserver"
	"github.com/ratedemon/go-rest/httpserver"
)

func main() {
	ctx := context.Background()
	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC, "caller", kitlog.DefaultCaller)

	cfg := config.Config{
		HTTPListenerAddress: 8081,
		GRPCListenerAddress: 8082,
		ExpLoginTimeout:     15,
		JWTSecret:           "SrTY3wmw80",
		DB: config.DB{
			User:     "rest_user",
			Password: "rest_password",
			Name:     "rest_db",
			Port:     5432,
		},
		Image: config.Image{
			SideMeasure:     160,
			ImagePrefixPath: "files",
		},
	}

	grpcServer, err := grpcserver.NewGRPCServer(ctx, &cfg, kitlog.With(logger, "type", "grpc server"))
	if err != nil {
		logger.Log("msg", "Failed to create new GRPC server", "err", err)
		os.Exit(1)
	}
	logger.Log("msg", "GRPC server is created")

	s, err := httpserver.NewServer(ctx, &cfg, kitlog.With(logger, "type", "http server"))
	if err != nil {
		logger.Log("msg", "Failed to create new HTTP server", "err", err)
		os.Exit(1)
	}
	logger.Log("msg", "HTTP server is created")

	grpcServerErrCh := make(chan error, 1)
	httpServerErrCh := make(chan error, 1)
	syscalCh := make(chan os.Signal, 1)
	signal.Notify(syscalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		grpcServerErrCh <- grpcServer.Run()
	}()

	go func() {
		httpServerErrCh <- s.Run()
	}()

	select {
	case osSignal := <-syscalCh:
		logger.Log("msg", "Server got signal", "signal", osSignal)
		os.Exit(1)
	case err = <-httpServerErrCh:
		logger.Log("msg", "HTTP server is going down", "err", err)
		os.Exit(1)
	case err = <-grpcServerErrCh:
		logger.Log("msg", "GRPC server is going down", "err", err)
		os.Exit(1)
	}
}
