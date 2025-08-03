package mocks

import (
	"bytes"
	"text/template"
	"os"
	"path/filepath"
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
func LoadTerraformConfigFile(filename string) string {
	content, err := os.ReadFile(filepath.Join("tests", "terraform", "acceptance", filename))
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
