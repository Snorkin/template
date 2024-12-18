package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"log"
	"time"
)

const (
	directoryPath = "config"
	name          = "config"
	extension     = "json"
)

var cfg *Config

type Config struct {
	Server      Server
	Logger      Logger
	Jaeger      Jaeger
	Sentry      Sentry
	Postgres    Postgres
	Environment string `validate:"required,oneof=dev stage prod"`
}

type Jaeger struct {
	URL         string `validate:"required"`
	Password    string `validate:"required"`
	Username    string `validate:"required"`
	ServiceName string `validate:"required"`
}

type Sentry struct {
	Enabled bool
	Dsn     string `validate:"required_if=Enabled true"`
}

type Logger struct {
	Level          string `validate:"required,oneof=debug info warn error panic fatal noLevel disabled"`
	SkipFrameCount int    `validate:"required"`
}

type Server struct {
	Http HTTP
	Grpc GRPC
}

type HTTP struct {
	Host string `validate:"required"`
	Port string `validate:"required"`
}

type GRPC struct {
	Host string `validate:"required"`
	Port string `validate:"required"`
}

type Postgres struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	User     string `validate:"required" observer:"ignore"`
	Password string `validate:"required" observer:"ignore"`
	DBName   string `validate:"required"`
	SSLMode  string `validate:"required"`
	PGDriver string `validate:"required"`
	Settings struct {
		MaxOpenConns    int           `validate:"required,min=1"`
		ConnMaxLifetime time.Duration `validate:"required,min=1"`
		MaxIdleConns    int           `validate:"required,min=1"`
		ConnMaxIdleTime time.Duration `validate:"required,min=1"`
	}
}

func GetConfig() *Config {
	return cfg
}

func LoadConfig() {
	var res Config

	v := viper.New()
	v.AddConfigPath(directoryPath)
	v.SetConfigName(name)
	v.SetConfigType(extension)

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Failed to load config from disc, %v", err)
	}

	if err := v.Unmarshal(&res); err != nil {
		log.Fatalf("Unable to parse config into struct, %v", err)
	}

	if err := validator.New().Struct(res); err != nil {
		log.Fatalf("Failed to validate struct, %v", err)
	}

	cfg = &res
}
