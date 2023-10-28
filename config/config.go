package config

import "time"

type Config struct {
	// Server
	Server Server `env:"SERVER"`

	// Postgres
	Postgres Postgres `env:"POSTGRES"`

	// Logger
	Logger Logger `env:"LOGGER"`
}

type Server struct {
	Port    int           `env:"SERVER_PORT" env-default:"8080"`
	Timeout time.Duration `env:"SERVER_TIMEOUT" env-default:"30s"`
}

type Postgres struct {
	URL string `env:"POSTGRES_URL" env-default:"postgresql://admin:admin@localhost:5434/interview?sslmode=disable&timezone=UTC"`
}

type Logger struct {
	FilePath string `env:"LOGGER_FILE_PATH" env-default:""`
}

func NewConfig() *Config {
	return &Config{
		Server: Server{},
		Logger: Logger{},
	}
}
