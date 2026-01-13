package helpers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDebugMode(t *testing.T) {
	t.Run("returns true when M365_DEBUG_MODE is set to true", func(t *testing.T) {
		// Save original value to restore after test
		originalValue := os.Getenv("M365_DEBUG_MODE")
		defer os.Setenv("M365_DEBUG_MODE", originalValue)

		os.Setenv("M365_DEBUG_MODE", "true")
		result := IsDebugMode()
		assert.True(t, result, "IsDebugMode should return true when M365_DEBUG_MODE=true")
	})

	t.Run("returns true when M365_DEBUG_MODE is set to TRUE (case insensitive)", func(t *testing.T) {
		originalValue := os.Getenv("M365_DEBUG_MODE")
		defer os.Setenv("M365_DEBUG_MODE", originalValue)

		os.Setenv("M365_DEBUG_MODE", "TRUE")
		result := IsDebugMode()
		assert.True(t, result, "IsDebugMode should return true when M365_DEBUG_MODE=TRUE")
	})

	t.Run("returns true when M365_DEBUG_MODE is set to TrUe (mixed case)", func(t *testing.T) {
		originalValue := os.Getenv("M365_DEBUG_MODE")
		defer os.Setenv("M365_DEBUG_MODE", originalValue)

		os.Setenv("M365_DEBUG_MODE", "TrUe")
		result := IsDebugMode()
		assert.True(t, result, "IsDebugMode should return true when M365_DEBUG_MODE=TrUe")
	})

	t.Run("returns false when M365_DEBUG_MODE is set to false", func(t *testing.T) {
		originalValue := os.Getenv("M365_DEBUG_MODE")
		defer os.Setenv("M365_DEBUG_MODE", originalValue)

		os.Setenv("M365_DEBUG_MODE", "false")
		result := IsDebugMode()
		assert.False(t, result, "IsDebugMode should return false when M365_DEBUG_MODE=false")
	})

	t.Run("returns false when M365_DEBUG_MODE is not set", func(t *testing.T) {
		originalValue := os.Getenv("M365_DEBUG_MODE")
		defer func() {
			if originalValue != "" {
				os.Setenv("M365_DEBUG_MODE", originalValue)
			} else {
				os.Unsetenv("M365_DEBUG_MODE")
			}
		}()

		os.Unsetenv("M365_DEBUG_MODE")
		result := IsDebugMode()
		assert.False(t, result, "IsDebugMode should return false when M365_DEBUG_MODE is not set")
	})

	t.Run("returns false when M365_DEBUG_MODE is set to empty string", func(t *testing.T) {
		originalValue := os.Getenv("M365_DEBUG_MODE")
		defer os.Setenv("M365_DEBUG_MODE", originalValue)

		os.Setenv("M365_DEBUG_MODE", "")
		result := IsDebugMode()
		assert.False(t, result, "IsDebugMode should return false when M365_DEBUG_MODE is empty")
	})

	t.Run("returns false when M365_DEBUG_MODE is set to invalid value", func(t *testing.T) {
		originalValue := os.Getenv("M365_DEBUG_MODE")
		defer os.Setenv("M365_DEBUG_MODE", originalValue)

		os.Setenv("M365_DEBUG_MODE", "yes")
		result := IsDebugMode()
		assert.False(t, result, "IsDebugMode should return false when M365_DEBUG_MODE has invalid value")
	})

	t.Run("returns false when M365_DEBUG_MODE is set to 1", func(t *testing.T) {
		originalValue := os.Getenv("M365_DEBUG_MODE")
		defer os.Setenv("M365_DEBUG_MODE", originalValue)

		os.Setenv("M365_DEBUG_MODE", "1")
		result := IsDebugMode()
		assert.False(t, result, "IsDebugMode should return false when M365_DEBUG_MODE=1")
	})
}
