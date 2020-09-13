package api

import (
	"context"

	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"

	"github.com/ratedemon/go-rest/api/auth"
	"github.com/ratedemon/go-rest/api/helper"
)

func registerRoute(ah helper.ApiHandler) []helper.Route {
	return ah.RegisterRoutes()
}

// InitRoutes defines routes and handlers for themselves
func InitRoutes(ctx context.Context, log log.Logger, grpcConn *grpc.ClientConn) []helper.Route {
	routes := []helper.Route{}

	{
		auth := auth.NewAuthHandler(ctx, log, grpcConn)
		routes = append(routes, registerRoute(auth)...)
	}

	return routes
}
