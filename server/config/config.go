package config

type DB struct {
	User     string
	Password string
	Name     string
	Port     int
}

type Image struct {
	SideMeasure     int
	ImagePrefixPath string
}

type Config struct {
	HTTPListenerAddress int
	GRPCListenerAddress int
	ExpLoginTimeout     int
	JWTSecret           string

	DB DB

	Image Image
}
