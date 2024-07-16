package provider

import "os"

// getEnvOrDefault fetches an environment variable value or returns a default if not set.
func getEnvOrDefault(value, envVar string) string {
	if value == "" {
		return os.Getenv(envVar)
	}
	return value
}
