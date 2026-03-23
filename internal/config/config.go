package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Addr       string
	pgUsername string
	pgPassword string
	pgHost     string
	pgPort     string
	pgBasename string
}

func ReadConfig() *Config {
	cfg := &Config{}

	viper.AutomaticEnv()
	viper.SetConfigType("env")

	cfg.Addr = viper.GetString("APP_PORT")
	cfg.pgUsername = viper.GetString("POSTGRES_USER")
	cfg.pgPassword = viper.GetString("POSTGRES_PASSWORD")
	cfg.pgHost = viper.GetString("POSTGRES_HOST")
	cfg.pgPort = viper.GetString("POSTGRES_PORT")
	cfg.pgBasename = viper.GetString("POSTGRES_DB")
	return cfg
}

func (c *Config) GetPostgresUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		c.pgUsername,
		c.pgPassword,
		c.pgHost,
		c.pgPort,
		c.pgBasename,
	)
}
