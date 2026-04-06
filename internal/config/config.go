package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	Addr           string
	CuBaseUrl      string
	AllowedOrigins string
	PgUsername     string
	PgPassword     string
	PgHost         string
	PgPort         string
	PgDatabase     string
}

func ReadConfig() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetConfigType("env")

	cfg := &Config{
		Addr:           viper.GetString("APP_PORT"),
		CuBaseUrl:      viper.GetString("CU_BASE_URL"),
		AllowedOrigins: viper.GetString("ALLOWED_ORIGINS"),
		PgUsername:     viper.GetString("POSTGRES_USER"),
		PgPassword:     viper.GetString("POSTGRES_PASSWORD"),
		PgHost:         viper.GetString("POSTGRES_HOST"),
		PgPort:         viper.GetString("POSTGRES_PORT"),
		PgDatabase:     viper.GetString("POSTGRES_DB"),
	}

	var errs []error
	if cfg.Addr == "" {
		errs = append(errs, errors.New("APP_PORT is required"))
	}
	if cfg.CuBaseUrl == "" {
		errs = append(errs, errors.New("CU_BASE_URL is required"))
	}
	if cfg.PgUsername == "" {
		errs = append(errs, errors.New("POSTGRES_USER is required"))
	}
	if cfg.PgPassword == "" {
		errs = append(errs, errors.New("POSTGRES_PASSWORD is required"))
	}
	if cfg.PgHost == "" {
		errs = append(errs, errors.New("POSTGRES_HOST is required"))
	}
	if cfg.PgPort == "" {
		errs = append(errs, errors.New("POSTGRES_PORT is required"))
	}
	if cfg.PgDatabase == "" {
		errs = append(errs, errors.New("POSTGRES_DB is required"))
	}
	if cfg.AllowedOrigins == "" {
		errs = append(errs, errors.New("ALLOWED_ORIGINS is required"))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return cfg, nil
}
