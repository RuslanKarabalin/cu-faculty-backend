package config

import (
	"errors"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Addr           string
	CuBaseUrl      string
	AllowedOrigins []string
	PgUsername     string
	PgPassword     string
	PgHost         string
	PgPort         string
	PgDatabase     string
	PgMaxConns     int32
	PgMinConns     int32
	S3Endpoint     string
	S3Region       string
	S3AccessKey    string
	S3SecretKey    string
	S3Bucket       string
	S3UsePathStyle bool
	S3PresignTTL   time.Duration
}

func ReadConfig() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetConfigType("env")
	viper.SetDefault("POSTGRES_MAX_CONNS", 25)
	viper.SetDefault("POSTGRES_MIN_CONNS", 5)
	viper.SetDefault("S3_REGION", "garage")
	viper.SetDefault("S3_USE_PATH_STYLE", true)
	viper.SetDefault("S3_PRESIGN_TTL", 15*time.Minute)

	cfg := &Config{
		Addr:           viper.GetString("APP_PORT"),
		CuBaseUrl:      viper.GetString("CU_BASE_URL"),
		AllowedOrigins: viper.GetStringSlice("ALLOWED_ORIGINS"),
		PgUsername:     viper.GetString("POSTGRES_USER"),
		PgPassword:     viper.GetString("POSTGRES_PASSWORD"),
		PgHost:         viper.GetString("POSTGRES_HOST"),
		PgPort:         viper.GetString("POSTGRES_PORT"),
		PgDatabase:     viper.GetString("POSTGRES_DB"),
		PgMaxConns:     viper.GetInt32("POSTGRES_MAX_CONNS"),
		PgMinConns:     viper.GetInt32("POSTGRES_MIN_CONNS"),
		S3Endpoint:     viper.GetString("S3_ENDPOINT"),
		S3Region:       viper.GetString("S3_REGION"),
		S3AccessKey:    viper.GetString("S3_ACCESS_KEY"),
		S3SecretKey:    viper.GetString("S3_SECRET_KEY"),
		S3Bucket:       viper.GetString("S3_BUCKET"),
		S3UsePathStyle: viper.GetBool("S3_USE_PATH_STYLE"),
		S3PresignTTL:   viper.GetDuration("S3_PRESIGN_TTL"),
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
	if cfg.AllowedOrigins == nil {
		errs = append(errs, errors.New("ALLOWED_ORIGINS is required"))
	}
	if cfg.S3Endpoint == "" {
		errs = append(errs, errors.New("S3_ENDPOINT is required"))
	}
	if cfg.S3AccessKey == "" {
		errs = append(errs, errors.New("S3_ACCESS_KEY is required"))
	}
	if cfg.S3SecretKey == "" {
		errs = append(errs, errors.New("S3_SECRET_KEY is required"))
	}
	if cfg.S3Bucket == "" {
		errs = append(errs, errors.New("S3_BUCKET is required"))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return cfg, nil
}
