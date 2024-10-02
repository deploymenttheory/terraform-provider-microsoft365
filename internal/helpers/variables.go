package helpers

import (
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MultiEnvDefaultFunc is a helper function that returns the value of the first
// environment variable in the given list that returns a non-empty value. If
// none of the environment variables return a value, the default value is
// returned.
func MultiEnvDefaultFunc(ks []string, defaultValue string) string {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return defaultValue
}

// EnvDefaultFunc is a helper function that returns the value of the
// given environment variable, if one exists, or the default value
// otherwise.
func EnvDefaultFunc(k string, defaultValue string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return defaultValue
}

// EnvDefaultFuncBool is a helper function that returns the boolean value of the
// given environment variable, if one exists, or the default value otherwise.
func EnvDefaultFuncBool(k string, defaultValue bool) bool {
	if v := os.Getenv(k); v != "" {
		b, err := strconv.ParseBool(v)
		if err == nil {
			return b
		}
	}
	return defaultValue
}

// EnvDefaultFuncInt64Value is a helper function that returns the types.Int64Value
// of the given environment variable, if one exists, or the default value otherwise.
func EnvDefaultFuncInt64Value(k string, defaultValue types.Int64) types.Int64 {
	if v := os.Getenv(k); v != "" {
		i, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return types.Int64Value(i)
		}
	}
	return defaultValue
}

// EnvDefaultFuncStringList is a helper function that returns a slice of strings
// based on the environment variable (if set) or the provided default value.
// The environment variable should be a comma-separated string.
func EnvDefaultFuncStringList(k string, defaultValue []string) []string {
	if v := os.Getenv(k); v != "" {
		elements := strings.Split(v, ",")
		var result []string
		for _, element := range elements {
			trimmed := strings.TrimSpace(element)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}
