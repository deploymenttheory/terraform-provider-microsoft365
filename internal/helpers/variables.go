package helpers

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// GetEnvOrDefault fetches an environment variable value or returns a default if not set.
func GetEnvOrDefault(ctx context.Context, value, envKey string) string {
	if value == "" || value == "null" {
		if envValue, exists := os.LookupEnv(envKey); exists {
			tflog.Debug(ctx, "Using environment variable value", map[string]interface{}{
				"envKey": envKey,
				"value":  envValue,
			})
			return envValue
		}
	}
	tflog.Debug(ctx, "Using provided value", map[string]interface{}{
		"envKey": envKey,
		"value":  value,
	})
	return value
}

// GetEnvOrDefaultBool fetches an environment variable value or returns a default if not set.
func GetEnvOrDefaultBool(ctx context.Context, value bool, envKey string) bool {
	if !value {
		if envValue, exists := os.LookupEnv(envKey); exists {
			result := envValue == "true"
			tflog.Debug(ctx, "Using environment variable value", map[string]interface{}{
				"envKey": envKey,
				"value":  result,
			})
			return result
		}
	}
	tflog.Debug(ctx, "Using provided value", map[string]interface{}{
		"envKey": envKey,
		"value":  value,
	})
	return value
}
