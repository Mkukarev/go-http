package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type HTTPServer struct {
	IdleTimeout  time.Duration
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Database struct {
	DatabaseURL        string `envconfig:"DATABASE_URL" required:"true"`
	LogLevel           string `envconfig:"DATABASE_LOG_LEVEL" default:"warn"`
	MaxOpenConnections int    `envconfig:"DATABASE_MAX_OPEN_CONNECTIONS" default:"10"`
}

type TodoGRPCServer struct {
	TodoGRPCServerUrl string
}

type Configuration struct {
	HTTPServer
	Database
	TodoGRPCServer
}

const envPrefix = ""

func toInt(s string) int {
	v, err := strconv.Atoi(s)

	// TODO: временное решение, нужно подумать как сделать лучше
	if err != nil {
		return 1
	}

	return v
}

func duration(s string) time.Duration {
	return time.Duration(toInt(s)) * time.Second
}

func Load() (Configuration, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	server := HTTPServer{
		IdleTimeout:  duration(os.Getenv("HTTP_SERVER_IDLE_TIMEOUT")),
		Port:         toInt(os.Getenv("HTTP_SERVER_PORT")),
		ReadTimeout:  duration(os.Getenv("HTTP_SERVER_READ_TIMEOUT")),
		WriteTimeout: duration(os.Getenv("HTTP_SERVER_WRITE_TIMEOUT")),
	}

	db := Database{
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		LogLevel:           os.Getenv("DATABASE_LOG_LEVEL"),
		MaxOpenConnections: toInt(os.Getenv("DATABASE_MAX_OPEN_CONNECTIONS")),
	}

	rpc := TodoGRPCServer{
		TodoGRPCServerUrl: os.Getenv("GRPC_SERVER_URL"),
	} 

	var cfg Configuration = Configuration{
		HTTPServer: server,
		Database:   db,
		TodoGRPCServer: rpc,
	}

	return cfg, nil
}
