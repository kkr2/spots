package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

// Config is all config struct
type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Logger   Logger
}

// ServerConfig  struct
type ServerConfig struct {
	AppVersion   string
	Port         string
	Mode         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	//SSL               bool
	CtxDefaultTimeout time.Duration
	Debug             bool
}

// Logger config
type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

// PostgresConfig contains all postgres properties
type PostgresConfig struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
	PostgresqlSSLMode  bool
	PgDriver           string
}

// ServerConfig  struct
type AnotherServerConfing struct {
	AppVersion   string
	Port         string
	Mode         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	CustomTimeout time.Duration
	Debug             bool
}


// LoadConfig file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv ()
	 if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// ParseConfig is used to get confings from file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
