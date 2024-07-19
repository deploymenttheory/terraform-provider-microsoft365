package helpers

import (
	"encoding/base64"
	"fmt"
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
