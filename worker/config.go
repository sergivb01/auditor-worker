package worker

import (
	"os"

	"github.com/joho/godotenv"
)

// Config defines the server configuration
type Config struct {
	RedisAddr string

	Production bool

	CCacheEnabled bool

	RemoveOld bool

	TLSCert string
	TLSKey  string
}

func loadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	return &Config{
		RedisAddr:     getEnv("REDIS_ADDRESS", "localhost:6379"),
		Production:    getEnv("PRODUCTION", "false") == "true",
		RemoveOld:     getEnv("REMOVE_OLD", "true") == "true",
		CCacheEnabled: getEnv("CCACHE_ENABLED", "true") == "true",
	}, nil
}

func getEnv(key string, defaultVal string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return defaultVal
}
