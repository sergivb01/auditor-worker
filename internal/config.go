package worker

import (
	"os"

	"github.com/joho/godotenv"
)

// Config defines the server configuration
type Config struct {
	Listen     string
	Production bool
	TLSCert    string
	TLSKey     string
}

func loadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	return &Config{
		Listen:     getEnv("LISTEN_ADDR", "localhost:3000"),
		Production: getEnv("PRODUCTION", "false") == "true",
		TLSCert:    getEnv("TLS_CERT", "C:\\Users\\Sergi\\Desktop\\acmecopy\\certs\\certificate.pem"),
		TLSKey:     getEnv("TLS_KEY", "C:\\Users\\Sergi\\Desktop\\acmecopy\\certs\\key.pem"),
	}, nil
}

func getEnv(key string, defaultVal string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return defaultVal
}
