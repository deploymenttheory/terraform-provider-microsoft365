package utilities

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Base64Encode encodes the input string to base64.
func Base64Encode(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("input string is empty")
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(input))
	return encoded, nil
}

// StringToInt converts a string to an integer based on a provided map
func StringToInt(str string, mapping map[string]int) (int, error) {
	if val, exists := mapping[str]; exists {
		return val, nil
	}
	return -1, fmt.Errorf("invalid string: %s. Supported strings: %v", str, mapping)
}

// Utility function to convert a bool to *bool
func BoolPtr(b bool) *bool {
	return &b
}

// Helper function to convert a string to uppercase
func ToUpperCase(s string) string {
	return strings.ToUpper(s)
}

// DownloadImage downloads an image from a given URL and returns it as a byte slice
// with retries
func DownloadImage(url string) ([]byte, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK response code: %d", resp.StatusCode)
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %v", err)
	}

	return imageData, nil
}

// IsDebugMode checks if the M365_DEBUG_MODE environment variable is set to true
func IsDebugMode() bool {
	debugMode := os.Getenv("M365_DEBUG_MODE")
	return strings.ToLower(debugMode) == "true"
}
