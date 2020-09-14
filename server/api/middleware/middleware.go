package middleware

import (
	"github.com/go-kit/kit/log"
	"github.com/ratedemon/go-rest/config"
)

// HTTPMiddleware defines middleware entity
type HTTPMiddleware struct {
	cfg *config.Config
	log log.Logger
}

// NewHTTPMiddleware creates new middleware entity
func NewHTTPMiddleware(cfg *config.Config, log log.Logger) *HTTPMiddleware {
	return &HTTPMiddleware{cfg, log}
}
