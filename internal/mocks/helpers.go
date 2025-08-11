package mocks

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"text/template"
)

// EnsureField ensures that a field exists in a data map
// If the field doesn't exist, it will be initialized with the provided default value
// This is particularly useful for ensuring collection fields are initialized in mock responses
func EnsureField(data map[string]interface{}, fieldName string, defaultValue interface{}) {
	if data[fieldName] == nil {
		data[fieldName] = defaultValue
	}
}

// LoadTerraformConfigFile reads a terraform configuration file from the localised
// acceptance test directory in services/<resource>/<resource_type>/tests/terraform/acceptance
// folder.
// Deprecated: Use LoadLocalTerraformConfig instead for better clarity
func LoadTerraformConfigFile(filename string) string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "acceptance", filename))
	if err != nil {
		// Fallback to empty string if file cannot be read
		return ""
	}
	return string(content)
}

// LoadLocalTerraformConfig reads a terraform configuration file from the local test directory
// in services/<resource>/<resource_type>/tests/terraform/acceptance folder.
// Falls back to unit directory if not found in acceptance.
func LoadLocalTerraformConfig(filename string) string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "acceptance", filename))
	if err != nil {
		// Try loading from unit directory as fallback
		content, err = os.ReadFile(filepath.Join("tests", "terraform", "unit", filename))
		if err != nil {
			// Fallback to empty string if file cannot be read
			return ""
		}
	}
	return string(content)
}

// LoadCentralizedTerraformConfig reads a terraform configuration file from centralized dependencies
// This is used for loading shared dependency files from the acceptance/terraform_dependencies directory
func LoadCentralizedTerraformConfig(filepath string) string {
	content, err := os.ReadFile(filepath)
	if err != nil {
		// Fallback to empty string if file cannot be read
		return ""
	}
	return string(content)
}

// LoadTerraformTemplateFile reads a terraform template file and applies the provided data
// This is used for
func LoadTerraformTemplateFile(filename string, data interface{}) string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "acceptance", filename))
	if err != nil {
		// Fallback to empty string if file cannot be read
		return ""
	}

	tmpl, err := template.New("terraform").Parse(string(content))
	if err != nil {
		return ""
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return ""
	}

	return buf.String()
}

// LoadJSONResponse loads a JSON response file and returns its contents
// This is a common utility function used by mock responders to load test response data
func LoadJSONResponse(filepath string) (map[string]interface{}, error) {
	jsonData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	err = json.Unmarshal(jsonData, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
