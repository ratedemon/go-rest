package config

type DB struct {
	User     string
	Password string
	Name     string
	Port     int
}

type Config struct {
	HTTPListenerAddress int
	GRPCListenerAddress int
	ExpLoginTimeout     int
	JWTSecret           string

	DB DB
}
