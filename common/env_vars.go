package common

import (
	"errors"
	"os"
)

// Fetch the APP_ENV env variable
func GetAppEnv() string {
	return os.Getenv("APP_ENV")
}

// Fetch the MONGODB URI env variable
func getEnvMongoDbUri() string {
	return os.Getenv("MONGODB_URI")
}

// Fetch the MONGODB name env variable
func getEnvMongoDbName() string {
	return os.Getenv("DB_NAME")
}

// Fetch the LOGPATH env variable
func getLogPath() string {
	return os.Getenv("LOG_PATH")
}

// Fetch the allowed origins regex env variable
func GetAllowedOriginsRegex() string {
	return os.Getenv("ALLOWED_ORIGINS_REGEX")
}

// Load application environment specific variables
func LoadEnvSpecificConfigVariables(cfg *Config) (err error) {
	if cfg.Env != "development" && cfg.Env != "test" {
		if mongoDBUri := getEnvMongoDbUri(); mongoDBUri != "" {
			cfg.GoBookStore.URI = mongoDBUri
		} else {
			return errors.New("mongodb uri environment variable not set")
		}

		if mongoDbName := getEnvMongoDbName(); mongoDbName != "" {
			cfg.GoBookStore.DB = mongoDbName
		} else {
			return errors.New("mongodb name env variable not set")
		}

		if logPath := getLogPath(); logPath != "" {
			cfg.GoBookStore.LOGPATH = logPath
		} else {
			return errors.New("logpath env variable not set")
		}
	}
	return
}
