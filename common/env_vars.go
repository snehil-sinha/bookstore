package common

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

// Fetch the MONGODB URI env variable
func getEnvMongoDbUri() string {
	return os.Getenv("MONGODB_URI")
}

// Fetch the MONGODB URI env variable
func getEnvMongoDbName() string {
	return os.Getenv("DB_NAME")
}

// Fetch the LOGPATH env variable
func getLogPath() string {
	return os.Getenv("LOG_PATH")
}

// Fetch CORS allowed origins env variable
func getAllowedOrigins() []string {
	origins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}
	return origins
}

// Fetch CORS allowed methods env variable
func getAllowedMethods() []string {
	return parseEnv("ALLOWED_METHODS")
}

// Fetch CORS allowed headers env variable
func getAllowedHeaders() []string {
	return parseEnv("ALLOWED_HEADERS")
}

// Fetch CORS exposed headers env variable
func getExposedHeaders() []string {
	return parseEnv("EXPOSED_HEADERS")
}

// Fetch CORS max age env variable
func getMaxAge() (time.Duration, error) {
	return time.ParseDuration(os.Getenv("MAX_AGE") + "s") // append s for seconds parsing
}

// Fetch CORS allow credentials env variable
func getAllowCredentials() (bool, error) {
	return strconv.ParseBool(os.Getenv("ALLOW_CREDENTIALS"))
}

// Load application environment specific variables
func LoadEnvSpecificConfigVariables(cfg *Config) (err error) {
	// if cfg.Env != "development" && cfg.Env != "test" {
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

	if allowedOrigins := getAllowedOrigins(); len(allowedOrigins) != 0 {
		cfg.GoBookStore.CORS.ALLOWED_ORIGINS = allowedOrigins
	} else {
		return errors.New("CORS allowed origins not set")
	}

	if allowedHeaders := getAllowedHeaders(); len(allowedHeaders) != 0 {
		cfg.GoBookStore.CORS.ALLOWED_HEADERS = allowedHeaders
	} else {
		return errors.New("CORS allowed headers not set")
	}

	if exposedHeaders := getExposedHeaders(); len(exposedHeaders) != 0 {
		cfg.GoBookStore.CORS.EXPOSED_HEADERS = exposedHeaders
	} else {
		return errors.New("CORS exposed headers not set")
	}

	if allowedMethods := getAllowedMethods(); len(allowedMethods) != 0 {
		cfg.GoBookStore.CORS.ALLOWED_METHOS = allowedMethods
	} else {
		return errors.New("CORS exposed headers not set")
	}

	if allowCredentials, err := getAllowCredentials(); err == nil {
		cfg.GoBookStore.CORS.ALLOW_CREDENTIALS = allowCredentials
	} else {
		return errors.New("CORS allow credentials not set")
	}

	if maxAge, err := getMaxAge(); err == nil {
		cfg.GoBookStore.CORS.MAX_AGE = maxAge
	} else {
		return errors.New("CORS max age not set. error parsing maxage env variable")
	}
	// }
	return
}
