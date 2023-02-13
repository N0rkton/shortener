package config

import (
	"flag"
	"os"
)

type Cfg struct {
	ServerAddress   *string
	BaseURL         *string
	FileStoragePath *string
	DbAddress       *string
}

var config Cfg

const defaultBaseURL = "http://localhost:8080"

func init() {
	config.ServerAddress = flag.String("a", "localhost:8080", "server address")
	config.BaseURL = flag.String("b", defaultBaseURL, "base URL")
	config.FileStoragePath = flag.String("f", "", "file path")
	config.DbAddress = flag.String("d", "", "data base connection address")
}
func NewConfig() Cfg {
	flag.Parse()
	dbAddressEnv := os.Getenv("DATABASE_DSN")
	//dbAddressEnv = "postgresql://localhost:5432/shvm"
	if dbAddressEnv != "" {
		config.DbAddress = &dbAddressEnv
	}
	serverAddressEnv := os.Getenv("SERVER_ADDRESS")
	if serverAddressEnv != "" {
		config.ServerAddress = &serverAddressEnv
	}
	baseURLEnv := os.Getenv("BASE_URL")
	if baseURLEnv != "" {
		config.BaseURL = &baseURLEnv
	}
	fileStoragePathEnv := os.Getenv("FILE_STORAGE_PATH")
	if fileStoragePathEnv != "" {
		config.FileStoragePath = &fileStoragePathEnv
	}
	return config
}
func GetServerAddress() string {
	return *config.ServerAddress
}
