// Package config describes supported flags and environmental vars.
package config

import (
	"flag"
	"os"
	"strings"
)

// Cfg config of the app
type Cfg struct {
	ServerAddress   string
	BaseURL         *string
	FileStoragePath *string
	DBAddress       *string
	EnableHTTPS     *bool
}

var config Cfg

const defaultBaseURL = "http://localhost:8080"

func init() {
	config.ServerAddress = *flag.String("a", "localhost:8080", "server address")
	config.BaseURL = flag.String("b", defaultBaseURL, "base URL")
	config.FileStoragePath = flag.String("f", "", "file path")
	config.DBAddress = flag.String("d", "", "data base connection address")
	config.EnableHTTPS = flag.Bool("s", false, "enable HTTPS")
}

// NewConfig - new default config based on flag or environmental variables. Env variables prioritise flag.
func NewConfig() Cfg {
	flag.Parse()
	dbAddressEnv := os.Getenv("DATABASE_DSN")
	//dbAddressEnv = "postgresql://localhost:5432/shvm"
	if dbAddressEnv != "" {
		config.DBAddress = &dbAddressEnv
	}
	serverAddressEnv := os.Getenv("SERVER_ADDRESS")
	if serverAddressEnv != "" {
		config.ServerAddress = serverAddressEnv
	}
	baseURLEnv := os.Getenv("BASE_URL")
	if baseURLEnv != "" {
		config.BaseURL = &baseURLEnv
	}
	fileStoragePathEnv := os.Getenv("FILE_STORAGE_PATH")
	if fileStoragePathEnv != "" {
		config.FileStoragePath = &fileStoragePathEnv
	}

	enableHTTPSEnv := os.Getenv("ENABLE_HTTPS")
	enableHTTPSEnv = "true"
	if enableHTTPSEnv == "true" {
		*config.EnableHTTPS = true
		*config.BaseURL = strings.Replace(defaultBaseURL, "http", "https", 1)
	}
	return config
}

// GetServerAddress - returns server address from config
func GetServerAddress() string {
	return config.ServerAddress
}

// GetEnableHTTPS - returns setting for HTTPS server
func GetEnableHTTPS() bool {
	return *config.EnableHTTPS
}
