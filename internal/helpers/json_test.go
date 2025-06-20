package helpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseJSONFile(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "json-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	t.Run("valid JSON file", func(t *testing.T) {
		// Create a temporary JSON file
		jsonContent := `{"name": "test", "value": 123}`
		jsonPath := filepath.Join(tempDir, "valid.json")
		err := os.WriteFile(jsonPath, []byte(jsonContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write test JSON file: %v", err)
		}

		// Test parsing the file
		result := ParseJSONFile(t, jsonPath)
		assert.Equal(t, jsonContent, result, "Should return the correct JSON content")
	})

	// Note: We can't effectively test the error case for ParseJSONFile
	// because it calls t.Fatalf() which terminates the test immediately
}

func TestPrettyJSON(t *testing.T) {
	t.Run("valid struct", func(t *testing.T) {
		data := struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		}{
			Name:  "test",
			Value: 123,
		}

		expected := `{
  "name": "test",
  "value": 123
}`

		result, err := PrettyJSON(data)
		assert.NoError(t, err, "Should not return an error for valid struct")
		assert.Equal(t, expected, result, "Should return correctly formatted JSON")
	})

	t.Run("valid map", func(t *testing.T) {
		data := map[string]interface{}{
			"name":    "test",
			"value":   123,
			"enabled": true,
		}

		result, err := PrettyJSON(data)
		assert.NoError(t, err, "Should not return an error for valid map")

		// Since map iteration order is not guaranteed, we'll parse the result back and compare
		var parsed map[string]interface{}
		err = ParseJSON(result, &parsed)
		assert.NoError(t, err, "Should parse the generated JSON")
		assert.Equal(t, data["name"], parsed["name"], "Should have correct name value")
		// When unmarshaling JSON, numbers become float64
		assert.Equal(t, float64(123), parsed["value"], "Should have correct numeric value")
		assert.Equal(t, data["enabled"], parsed["enabled"], "Should have correct boolean value")
	})

	t.Run("invalid data", func(t *testing.T) {
		// Create a circular reference that can't be marshaled to JSON
		type Circular struct {
			Self *Circular
		}
		circular := &Circular{}
		circular.Self = circular

		result, err := PrettyJSON(circular)
		assert.Error(t, err, "Should return an error for data that can't be marshaled")
		assert.Empty(t, result, "Result should be empty for invalid data")
	})
}

func TestParseJSON(t *testing.T) {
	t.Run("valid JSON to struct", func(t *testing.T) {
		jsonStr := `{"name": "test", "value": 123}`
		var result struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		}

		err := ParseJSON(jsonStr, &result)
		assert.NoError(t, err, "Should not return an error for valid JSON")
		assert.Equal(t, "test", result.Name, "Should parse name correctly")
		assert.Equal(t, 123, result.Value, "Should parse value correctly")
	})

	t.Run("valid JSON to map", func(t *testing.T) {
		jsonStr := `{"name": "test", "value": 123, "enabled": true}`
		var result map[string]interface{}

		err := ParseJSON(jsonStr, &result)
		assert.NoError(t, err, "Should not return an error for valid JSON")
		assert.Equal(t, "test", result["name"], "Should parse name correctly")
		assert.Equal(t, float64(123), result["value"], "Should parse value correctly")
		assert.Equal(t, true, result["enabled"], "Should parse boolean correctly")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		jsonStr := `{"name": "test", "value": 123,}` // Invalid JSON with trailing comma
		var result map[string]interface{}

		err := ParseJSON(jsonStr, &result)
		assert.Error(t, err, "Should return an error for invalid JSON")
	})

	t.Run("JSON type mismatch", func(t *testing.T) {
		jsonStr := `{"name": "test", "value": "not-a-number"}`
		var result struct {
			Name  string `json:"name"`
			Value int    `json:"value"` // Expecting an int but getting a string
		}

		err := ParseJSON(jsonStr, &result)
		assert.Error(t, err, "Should return an error for type mismatch")
	})
}
