package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv   string
	HTTPPort string

	DBHost, DBPort, DBUser, DBPassword, DBName, DBSSLMode string

	HTTPDialTimeoutSec int
	HTTPTLSTimeoutSec  int
	HTTPReqTimeoutSec  int

	MaxLinkWorkers    int
	MaxLinksPerPage   int
	RequestTimeoutSec int
}

func MustLoad() *Config {
	_ = godotenv.Load()
	return &Config{
		AppEnv:             getenv("APP_ENV", "development"),
		HTTPPort:           getenv("HTTP_PORT", "8080"),
		DBHost:             getenv("DB_HOST", "localhost"),
		DBPort:             getenv("DB_PORT", "5432"),
		DBUser:             getenv("DB_USER", "postgres"),
		DBPassword:         getenv("DB_PASSWORD", "postgres"),
		DBName:             getenv("DB_NAME", "web_analyzer"),
		DBSSLMode:          getenv("DB_SSLMODE", "disable"),
		HTTPDialTimeoutSec: atoi("HTTP_DIAL_TIMEOUT", 5),
		HTTPTLSTimeoutSec:  atoi("HTTP_TLS_TIMEOUT", 5),
		HTTPReqTimeoutSec:  atoi("HTTP_REQ_TIMEOUT", 10),
		MaxLinkWorkers:     atoi("MAX_LINK_WORKERS", 20),
		MaxLinksPerPage:    atoi("MAX_LINKS_PER_PAGE", 500),
		RequestTimeoutSec:  atoi("REQUEST_TIMEOUT_SEC", 30),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func atoi(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}
