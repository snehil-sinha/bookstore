package common

import (
	"os"
	"strings"
)

func parseEnv(envVar string) []string {
	values := strings.Split(os.Getenv(envVar), ",")
	for i := range values {
		values[i] = strings.TrimSpace(values[i])
	}
	return values
}
