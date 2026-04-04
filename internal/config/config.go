package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Addr           string
	CuBaseUrl      string
	AllowedOrigins string
	pgUsername     string
	pgPassword     string
	pgHost         string
	pgPort         string
	pgBasename     string
}

func ReadConfig() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetConfigType("env")
	viper.SetDefault("ALLOWED_ORIGINS", "*")

	cfg := &Config{
		Addr:           viper.GetString("APP_PORT"),
		CuBaseUrl:      viper.GetString("CU_BASE_URL"),
		AllowedOrigins: viper.GetString("ALLOWED_ORIGINS"),
		pgUsername:     viper.GetString("POSTGRES_USER"),
		pgPassword:     viper.GetString("POSTGRES_PASSWORD"),
		pgHost:         viper.GetString("POSTGRES_HOST"),
		pgPort:         viper.GetString("POSTGRES_PORT"),
		pgBasename:     viper.GetString("POSTGRES_DB"),
	}

	var errs []error
	if cfg.Addr == "" {
		errs = append(errs, errors.New("APP_PORT is required"))
	}
	if cfg.CuBaseUrl == "" {
		errs = append(errs, errors.New("CU_BASE_URL is required"))
	}
	if cfg.pgUsername == "" {
		errs = append(errs, errors.New("POSTGRES_USER is required"))
	}
	if cfg.pgPassword == "" {
		errs = append(errs, errors.New("POSTGRES_PASSWORD is required"))
	}
	if cfg.pgHost == "" {
		errs = append(errs, errors.New("POSTGRES_HOST is required"))
	}
	if cfg.pgPort == "" {
		errs = append(errs, errors.New("POSTGRES_PORT is required"))
	}
	if cfg.pgBasename == "" {
		errs = append(errs, errors.New("POSTGRES_DB is required"))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return cfg, nil
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
