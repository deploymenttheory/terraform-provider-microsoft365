package helpers

import (
	"encoding/base64"
	"errors"
)

// Base64Encode encodes the input string to base64.
func Base64Encode(input string) (string, error) {
	if input == "" {
		return "", errors.New("input string is empty")
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(input))
	return encoded, nil
}
