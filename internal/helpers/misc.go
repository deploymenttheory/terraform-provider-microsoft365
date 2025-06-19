package helpers

import (
	"os"
	"strings"
)

// IsDebugMode checks if the M365_DEBUG_MODE environment variable is set to true
func IsDebugMode() bool {
	debugMode := os.Getenv("M365_DEBUG_MODE")
	return strings.ToLower(debugMode) == "true"
}
