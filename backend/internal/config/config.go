package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
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

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

type Config struct {
	Addr                      string
	ReadTimeout               time.Duration
	WriteTimeout              time.Duration
	DbConfig                  DbConfig
	JWTSecret                 string
	JWTExpiry                 time.Duration
	EmailVerificationBaseURL  string
	EmailVerificationTokenTTL time.Duration
	SMTP                      SMTPConfig
	CronSchedule              string
	CORSAllowedOrigins        []string
	GoogleOAuthClientID       string
}

func NewConfig(addr string, readTimeout, writeTimeout time.Duration, dbConfig DbConfig, jwtSecret string, jwtExpiry time.Duration, emailVerificationBaseURL string, emailVerificationTokenTTL time.Duration, smtp SMTPConfig, cronSchedule string, corsAllowedOrigins []string, googleOAuthClientID string) *Config {
	return &Config{
		Addr:                      addr,
		ReadTimeout:               readTimeout,
		WriteTimeout:              writeTimeout,
		DbConfig:                  dbConfig,
		JWTSecret:                 jwtSecret,
		JWTExpiry:                 jwtExpiry,
		EmailVerificationBaseURL:  emailVerificationBaseURL,
		EmailVerificationTokenTTL: emailVerificationTokenTTL,
		SMTP:                      smtp,
		CronSchedule:              cronSchedule,
		CORSAllowedOrigins:        corsAllowedOrigins,
		GoogleOAuthClientID:       googleOAuthClientID,
	}
}

func LoadConfig() (*Config, error) {

	_ = godotenv.Load()

	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DATABASE_URL")

	if port == "" || dbUrl == "" {
		return nil, errors.New("PORT or DATABASE_URL is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET is not set")
	}

	jwtExpiryStr := os.Getenv("JWT_EXPIRY")
	if jwtExpiryStr == "" {
		jwtExpiryStr = "24h"
	}
	jwtExpiry, err := time.ParseDuration(jwtExpiryStr)
	if err != nil {
		return nil, fmt.Errorf("JWT_EXPIRY: %w", err)
	}

	baseURL := os.Getenv("EMAIL_VERIFICATION_BASE_URL")
	if baseURL == "" {
		return nil, errors.New("EMAIL_VERIFICATION_BASE_URL is not set")
	}

	tokenTTLStr := os.Getenv("EMAIL_VERIFICATION_TOKEN_TTL")
	if tokenTTLStr == "" {
		tokenTTLStr = "48h"
	}
	tokenTTL, err := time.ParseDuration(tokenTTLStr)
	if err != nil {
		return nil, fmt.Errorf("EMAIL_VERIFICATION_TOKEN_TTL: %w", err)
	}

	smtpPort := 587
	if p := os.Getenv("SMTP_PORT"); p != "" {
		smtpPort, err = strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("SMTP_PORT: %w", err)
		}
	}

	smtp := SMTPConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     smtpPort,
		User:     os.Getenv("SMTP_USER"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
	}

	cronSchedule := os.Getenv("CRON_SCHEDULE")
	if cronSchedule == "" {
		cronSchedule = "5 0 * * *" // 00:05 UTC daily
	}

	corsRaw := os.Getenv("CORS_ALLOWED_ORIGINS")
	var corsAllowedOrigins []string
	if corsRaw != "" {
		for _, o := range strings.Split(corsRaw, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				corsAllowedOrigins = append(corsAllowedOrigins, o)
			}
		}
	}
	if len(corsAllowedOrigins) == 0 {
		// Default: allow local frontend dev servers
		corsAllowedOrigins = []string{"http://localhost:5173", "http://localhost:3001"}
	}

	googleClientID := strings.TrimSpace(os.Getenv("GOOGLE_OAUTH_CLIENT_ID"))

	return NewConfig(
		":"+port,
		time.Second*15,
		time.Second*15,
		*NewDbConfig(dbUrl, 25, 10, time.Hour, time.Minute*10),
		jwtSecret,
		jwtExpiry,
		baseURL,
		tokenTTL,
		smtp,
		cronSchedule,
		corsAllowedOrigins,
		googleClientID,
	), nil
}
