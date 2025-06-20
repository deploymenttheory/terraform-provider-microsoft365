package helpers

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

// ParseJSONFile reads a JSON file and returns its contents as a string.
// This is useful for mocking HTTP responses in tests.
func ParseJSONFile(t *testing.T, path string) string {
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open JSON file %s: %v", path, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("Failed to read JSON file %s: %v", path, err)
	}

	return string(data)
}

// PrettyJSON converts an interface to a pretty-printed JSON string.
func PrettyJSON(data interface{}) (string, error) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ParseJSON parses a JSON string into the provided interface.
func ParseJSON(jsonStr string, target interface{}) error {
	return json.Unmarshal([]byte(jsonStr), target)
}
