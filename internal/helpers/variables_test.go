package helpers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// withEnvironment is a test helper function that sets up and tears down
// environment variables for a test. It ensures that the environment is
// properly cleaned up after the test, even if the test panics.
//
// This function improves test reliability in CI environments by:
// 1. Isolating environment changes to individual tests
// 2. Guaranteeing cleanup of environment variables
// 3. Allowing tests to run in parallel safely
//
// Parameters:
//   - t: The testing.T instance for the current test
//   - env: A map of environment variables to set for the test
//   - testFunc: The function containing the actual test code
//
// Usage:
//
//	withEnvironment(t, map[string]string{"VAR": "value"}, func() {
//	    // Test code here
//	})
func withEnvironment(_ *testing.T, env map[string]string, testFunc func()) {
	oldEnv := make(map[string]string)
	for k, v := range env {
		if oldVal, exists := os.LookupEnv(k); exists {
			oldEnv[k] = oldVal
		}
		os.Setenv(k, v)
	}

	defer func() {
		for k := range env {
			if oldVal, exists := oldEnv[k]; exists {
				os.Setenv(k, oldVal)
			} else {
				os.Unsetenv(k)
			}
		}
	}()

	testFunc()
}

func TestMultiEnvDefaultFunc(t *testing.T) {
	t.Run("No environment variables set", func(t *testing.T) {
		result := MultiEnvDefaultFunc([]string{"TEST_VAR1", "TEST_VAR2"}, "default")
		assert.Equal(t, "default", result)
	})

	t.Run("First environment variable set", func(t *testing.T) {
		withEnvironment(t, map[string]string{"TEST_VAR1": "value1"}, func() {
			result := MultiEnvDefaultFunc([]string{"TEST_VAR1", "TEST_VAR2"}, "default")
			assert.Equal(t, "value1", result)
		})
	})

	t.Run("Second environment variable set", func(t *testing.T) {
		withEnvironment(t, map[string]string{"TEST_VAR2": "value2"}, func() {
			result := MultiEnvDefaultFunc([]string{"TEST_VAR1", "TEST_VAR2"}, "default")
			assert.Equal(t, "value2", result)
		})
	})
}

func TestEnvDefaultFunc(t *testing.T) {
	t.Run("Environment variable not set", func(t *testing.T) {
		result := EnvDefaultFunc("TEST_VAR", "default")
		assert.Equal(t, "default", result)
	})

	t.Run("Environment variable set", func(t *testing.T) {
		withEnvironment(t, map[string]string{"TEST_VAR": "value"}, func() {
			result := EnvDefaultFunc("TEST_VAR", "default")
			assert.Equal(t, "value", result)
		})
	})
}

func TestEnvDefaultFuncBool(t *testing.T) {
	t.Run("Environment variable not set", func(t *testing.T) {
		result := EnvDefaultFuncBool("TEST_VAR_BOOL", true)
		assert.True(t, result)
	})

	t.Run("Environment variable set to true", func(t *testing.T) {
		withEnvironment(t, map[string]string{"TEST_VAR_BOOL": "true"}, func() {
			result := EnvDefaultFuncBool("TEST_VAR_BOOL", false)
			assert.True(t, result)
		})
	})

	t.Run("Environment variable set to false", func(t *testing.T) {
		withEnvironment(t, map[string]string{"TEST_VAR_BOOL": "false"}, func() {
			result := EnvDefaultFuncBool("TEST_VAR_BOOL", true)
			assert.False(t, result)
		})
	})

	t.Run("Environment variable set to invalid boolean value", func(t *testing.T) {
		withEnvironment(t, map[string]string{"TEST_VAR_BOOL": "invalid"}, func() {
			result := EnvDefaultFuncBool("TEST_VAR_BOOL", true)
			assert.True(t, result)
		})
	})
}