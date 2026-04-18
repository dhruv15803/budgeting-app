package config

import (
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type DbConfig struct {
	Url             string
	MaxOpenConns    int
	MaxIdleConns    int
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

func NewDbConfig(url string, maxOpenConns, maxIdleConns int, maxConnLifetime, maxConnIdleTime time.Duration) *DbConfig {
	return &DbConfig{
		Url:             url,
		MaxOpenConns:    maxOpenConns,
		MaxIdleConns:    maxIdleConns,
		MaxConnLifetime: maxConnLifetime,
		MaxConnIdleTime: maxConnIdleTime,
	}
}

type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	DbConfig     DbConfig
}

func NewConfig(addr string, readTimeout, writeTimeout time.Duration, dbConfig DbConfig) *Config {
	return &Config{
		Addr:         addr,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		DbConfig:     dbConfig,
	}
}

func LoadConfig() (*Config, error) {

	_ = godotenv.Load()

	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DATABASE_URL")

	if port == "" || dbUrl == "" {
		return nil, errors.New("PORT or DATABASE_URL is not set")
	}

	return NewConfig(":"+port, time.Second*15, time.Second*15, *NewDbConfig(dbUrl, 25, 10, time.Hour, time.Minute*10)), nil
}
