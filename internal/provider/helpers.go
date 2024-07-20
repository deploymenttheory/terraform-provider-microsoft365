package provider

import "os"

// getEnvOrDefault fetches an environment variable value or returns a default if not set.
func getEnvOrDefault(value, envKey string) string {
	if value == "" {
		return os.Getenv(envKey)
	}
	return value
}

// getEnvOrDefaultInt fetches an environment variable value or returns a default if not set.
func getEnvOrDefaultBool(value bool, envKey string) bool {
	if value {
		return value
	}
	envValue, exists := os.LookupEnv(envKey)
	if exists && envValue == "true" {
		return true
	}
	return false
}
