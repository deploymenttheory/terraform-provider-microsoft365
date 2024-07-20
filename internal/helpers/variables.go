package helpers

import "os"

// GetEnvOrDefault fetches an environment variable value or returns a default if not set.
func GetEnvOrDefault(value, envKey string) string {
	if value == "" {
		return os.Getenv(envKey)
	}
	return value
}

// GetEnvOrDefaultInt fetches an environment variable value or returns a default if not set.
func GetEnvOrDefaultBool(value bool, envKey string) bool {
	if value {
		return value
	}
	envValue, exists := os.LookupEnv(envKey)
	if exists && envValue == "true" {
		return true
	}
	return false
}
