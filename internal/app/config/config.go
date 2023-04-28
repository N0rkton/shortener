// Package config describes supported flags and environmental vars.
package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"
)

// JSONConfig - struct to unmarshall json config file (less priority than env variable and flags)
type JSONConfig struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDsn     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
}

// Cfg config of the app
type Cfg struct {
	ServerAddress   string
	BaseURL         *string
	FileStoragePath *string
	DatabaseDsn     *string
	EnableHTTPS     *bool
	ConfigFile      *string
}

var config Cfg

const defaultBaseURL = "http://localhost:8080"
const defaultServerAddress = "localhost:8080"
const certFile = "cmd/shortener/certificate/localhost.crt"
const keyFile = "cmd/shortener/certificate/localhost.key"

func init() {
	config.ServerAddress = *flag.String("a", defaultServerAddress, "server address")
	config.BaseURL = flag.String("b", defaultBaseURL, "base URL")
	config.FileStoragePath = flag.String("f", "", "file path")
	config.DatabaseDsn = flag.String("d", "", "data base connection address")
	config.EnableHTTPS = flag.Bool("s", false, "enable HTTPS")
	config.ConfigFile = flag.String("c", "", "path to config file")
}

// NewConfig - new default config based on flag or environmental variables. Env variables prioritise flag.
func NewConfig() Cfg {
	flag.Parse()
	var jsonConfig JSONConfig
	configFileEnv := os.Getenv("CONFIG")
	//configFileEnv = "internal/app/config/config_test.json"
	if configFileEnv != "" {
		config.ConfigFile = &configFileEnv
	}
	if *config.ConfigFile != "" {
		file, err := os.ReadFile(*config.ConfigFile)
		if err != nil {
			log.Print(err)
		}

		err = json.Unmarshal(file, &jsonConfig)
		if err != nil {
			log.Print(err)
		}

	}
	dbAddressEnv := os.Getenv("DATABASE_DSN")
	//dbAddressEnv = "postgresql://localhost:5432/shvm"
	if dbAddressEnv != "" {
		config.DatabaseDsn = &dbAddressEnv
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
	//enableHTTPSEnv = "true"
	if enableHTTPSEnv == "true" {
		*config.EnableHTTPS = true
		*config.BaseURL = strings.Replace(*config.BaseURL, "http", "https", 1)
	}
	if config.ServerAddress == defaultServerAddress && jsonConfig.ServerAddress != "" {
		config.ServerAddress = jsonConfig.ServerAddress
	}
	if *config.BaseURL == defaultBaseURL && jsonConfig.BaseURL != "" {
		*config.BaseURL = jsonConfig.BaseURL
	}
	if *config.DatabaseDsn == "" && jsonConfig.DatabaseDsn != "" {
		*config.DatabaseDsn = jsonConfig.DatabaseDsn
	}
	if *config.FileStoragePath == "" && jsonConfig.FileStoragePath != "" {
		*config.FileStoragePath = jsonConfig.FileStoragePath
	}
	if !*config.EnableHTTPS && jsonConfig.EnableHTTPS {
		*config.EnableHTTPS = jsonConfig.EnableHTTPS
		*config.BaseURL = strings.Replace(*config.BaseURL, "http", "https", 1)
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

// GetCertFile - return path to .crt file
func GetCertFile() string {
	return certFile
}

// GetKeyFile - returns path to .key file
func GetKeyFile() string {
	return keyFile
}
