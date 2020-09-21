package config

import (
	"errors"
	"os"
	"strconv"
)

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

	DB *DB

	Image *Image
}

// NewConfig creates new config for entire system
func NewConfig() (*Config, error) {
	httpListener, err := strconv.Atoi(os.Getenv("HTTP_LISTENER_PORT"))
	if err != nil {
		return nil, err
	}
	grpcListener, err := strconv.Atoi(os.Getenv("GRPC_LISTENER_PORT"))
	if err != nil {
		return nil, err
	}
	loginTimeout, err := strconv.Atoi(os.Getenv("LOGIN_TIMEOUT"))
	if err != nil {
		return nil, err
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("'JWT_SECRET' must be filled")
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		return nil, errors.New("'DB_USER' must be filled")
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		return nil, errors.New("'DB_PASSWORD' must be filled")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return nil, errors.New("'DB_NAME' must be filled")
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, err
	}

	db := &DB{
		User:     dbUser,
		Password: dbPassword,
		Name:     dbName,
		Port:     dbPort,
	}

	// sm, err := strconv.Atoi(os.Getenv("SIDE_MEASURE"))
	// if err != nil {
	// 	return nil, err
	// }

	prefixPath := os.Getenv("IMAGE_PREFIX_PATH")
	if prefixPath == "" {
		return nil, errors.New("'IMAGE_PREFIX_PATH' must be filled")
	}

	img := &Image{
		SideMeasure:     160,
		ImagePrefixPath: prefixPath,
	}

	return &Config{
		HTTPListenerAddress: httpListener,
		GRPCListenerAddress: grpcListener,
		ExpLoginTimeout:     loginTimeout,
		JWTSecret:           jwtSecret,
		DB:                  db,
		Image:               img,
	}, nil
}
