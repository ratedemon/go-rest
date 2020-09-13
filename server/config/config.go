package config

type Config struct {
	HTTPListenerAddress int
	GRPCListenerAddress int
	ExpLoginTimeout     int
	JWTSecret           string
}
