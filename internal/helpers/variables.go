package helpers

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// GetValueOrEnv retrieves a string value from either the schema or an environment variable.
// It prioritizes the schema value if it's not null, unknown, or empty.
// If no value is found in either source, it logs an error and returns an empty string.
func GetValueOrEnv(ctx context.Context, schemaValue types.String, envVar string) string {
	if !schemaValue.IsNull() && !schemaValue.IsUnknown() {
		value := schemaValue.ValueString()
		if value != "" {
			return value
		}
	}

	envValue := os.Getenv(envVar)
	if envValue != "" {
		return envValue
	}

	tflog.Error(ctx, "Value not found in schema or environment variable", map[string]interface{}{
		"envVar": envVar,
	})

	return ""
}

// GetValueOrEnvBool retrieves a boolean value from either the schema or an environment variable.
// It prioritizes the schema value if it's not null or unknown.
// If using the environment variable, it considers "true" or "1" as true values.
// It logs the retrieved value and its source for debugging purposes.
func GetValueOrEnvBool(ctx context.Context, schemaValue types.Bool, envVar string) bool {
	if !schemaValue.IsNull() && !schemaValue.IsUnknown() {
		value := schemaValue.ValueBool()

		return value
	}

	envValue := os.Getenv(envVar)
	value := envValue == "true" || envValue == "1"

	return value
}
