package mocks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"text/template"
)

// EnsureField ensures that a field exists in a data map
// If the field doesn't exist, it will be initialized with the provided default value
// This is particularly useful for ensuring collection fields are initialized in mock responses
func EnsureField(data map[string]any, fieldName string, defaultValue any) {
	if data[fieldName] == nil {
		data[fieldName] = defaultValue
	}
}

// LoadTerraformConfigFile reads a terraform configuration file from the localised
// acceptance test directory in services/<resource>/<resource_type>/tests/terraform/acceptance
// folder.
// Deprecated: Use LoadUnitTerraformConfig for unit tests instead
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
// Deprecated: Use LoadUnitTerraformConfig for unit tests instead
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
// Deprecated: Use LoadAcceptanceDependencyConfig instead
func LoadCentralizedTerraformConfig(filepath string) string {
	content, err := os.ReadFile(filepath)
	if err != nil {
		// Fallback to empty string if file cannot be read
		return ""
	}
	return string(content)
}

// LoadUnitTerraformConfig loads a terraform configuration file for unit tests from the caller's local directory.
// This function is specifically designed for unit tests and only looks in tests/terraform/unit/.
// It uses runtime.Caller to determine the calling test file's location for proper path resolution.
func LoadUnitTerraformConfig(filename string) string {
	// Get the caller's directory to construct correct relative path
	_, callerFile, _, ok := runtime.Caller(1)
	if !ok {
		fmt.Printf("ERROR: LoadUnitTerraformConfig failed to get caller information for file: %s\n", filename)
		return ""
	}

	// Construct path relative to the caller's directory
	callerDir := filepath.Dir(callerFile)
	configPath := filepath.Join(callerDir, "tests", "terraform", "unit", filename)

	fmt.Printf("DEBUG: LoadUnitTerraformConfig loading - filename=%s, callerDir=%s, configPath=%s\n",
		filename, callerDir, configPath)

	content, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("ERROR: LoadUnitTerraformConfig failed to read terraform config file at %s: %v\n", configPath, err)
		return ""
	}

	if len(content) == 0 {
		fmt.Printf("WARN: LoadUnitTerraformConfig terraform config file at %s is empty\n", configPath)
		return ""
	}

	fmt.Printf("DEBUG: LoadUnitTerraformConfig successfully loaded config - filename=%s, contentSize=%d\n",
		filename, len(content))

	return string(content)
}

// LoadUnitTerraformTemplate loads and processes a terraform template file for unit tests.
// This function is specifically designed for unit tests and only looks in tests/terraform/unit/.
// It uses runtime.Caller to determine the calling test file's location for proper path resolution.
func LoadUnitTerraformTemplate(filename string, data any) string {
	// Get the caller's directory to construct correct relative path
	_, callerFile, _, ok := runtime.Caller(1)
	if !ok {
		fmt.Printf("ERROR: LoadUnitTerraformTemplate failed to get caller information for file: %s\n", filename)
		return ""
	}

	// Construct path relative to the caller's directory
	callerDir := filepath.Dir(callerFile)
	configPath := filepath.Join(callerDir, "tests", "terraform", "unit", filename)

	fmt.Printf("DEBUG: LoadUnitTerraformTemplate loading - filename=%s, callerDir=%s, configPath=%s\n",
		filename, callerDir, configPath)

	content, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("ERROR: LoadUnitTerraformTemplate failed to read terraform template file at %s: %v\n", configPath, err)
		return ""
	}

	if len(content) == 0 {
		fmt.Printf("WARN: LoadUnitTerraformTemplate terraform template file at %s is empty\n", configPath)
		return ""
	}

	tmpl, err := template.New("terraform").Parse(string(content))
	if err != nil {
		fmt.Printf("ERROR: LoadUnitTerraformTemplate failed to parse template at %s: %v\n", configPath, err)
		return ""
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		fmt.Printf("ERROR: LoadUnitTerraformTemplate failed to execute template at %s: %v\n", configPath, err)
		return ""
	}

	result := buf.String()
	fmt.Printf("DEBUG: LoadUnitTerraformTemplate successfully processed template - filename=%s, resultSize=%d\n",
		filename, len(result))

	return result
}

// LoadTerraformTemplateFile reads a terraform template file and applies the provided data
// This is used for
func LoadTerraformTemplateFile(filename string, data any) string {
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
func LoadJSONResponse(filepath string) (map[string]any, error) {
	jsonData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var response map[string]any
	err = json.Unmarshal(jsonData, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
