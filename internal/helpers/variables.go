package helpers

import "os"

// GetEnvOrDefault fetches an environment variable value or returns a default if not set.
func GetEnvOrDefault(value, envKey string) string {
	if value == "" {
		if envValue, exists := os.LookupEnv(envKey); exists {
			return envValue
		}
	}
	return value
}

// GetEnvOrDefaultBool fetches an environment variable value or returns a default if not set.
func GetEnvOrDefaultBool(value bool, envKey string) bool {
	if !value {
		if envValue, exists := os.LookupEnv(envKey); exists {
			return envValue == "true"
		}
	}
	return value
}
