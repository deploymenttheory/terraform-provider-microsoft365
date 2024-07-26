package helpers

import (
	"os"
	"strconv"
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
