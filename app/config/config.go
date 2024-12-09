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

type Config struct {
	Server      Server
	Logger      Logger
	Jaeger      Jaeger
	Postgres    Postgres
	Environment string `validate:"required,oneof=dev stage prod"`
}

type Jaeger struct {
	URL         string `validate:"required"`
	Password    string `validate:"required"`
	Username    string `validate:"required"`
	ServiceName string `validate:"required"`
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
	User     string `validate:"required"`
	Password string `validate:"required"`
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

func LoadConfig() (Config, error) {
	var res Config

	v := viper.New()
	v.AddConfigPath(directoryPath)
	v.SetConfigName(name)
	v.SetConfigType(extension)

	if err := v.ReadInConfig(); err != nil {
		return res, err
	}

	if err := v.Unmarshal(&res); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
		return res, err
	}

	if err := validator.New().Struct(res); err != nil {
		return res, err
	}

	return res, nil
}
