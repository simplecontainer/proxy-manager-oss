package configuration

import (
	"os"
)

func New() *Configuration {
	return &Configuration{
		AllowOrigin: getEnv("ALLOW_ORIGIN_FROM"),
		MasterPort:  getEnv("MASTER_PORT"),
		Port:        getEnv("PORT"),
		Environment: getEnv("ENVIRONMENT"),
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)

	if val == "" {
		panic("Missing required environment variable: " + key)
	}

	return val
}
