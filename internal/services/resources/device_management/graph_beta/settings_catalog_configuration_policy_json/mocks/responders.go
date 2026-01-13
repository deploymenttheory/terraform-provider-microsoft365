package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/jarcoal/httpmock"
)

// SettingsCatalogConfigurationPolicyJsonMock handles HTTP mocking for Settings Catalog JSON policies
type SettingsCatalogConfigurationPolicyJsonMock struct{}

// getTestDataPath returns the path to the test data directory
func (m *SettingsCatalogConfigurationPolicyJsonMock) getTestDataPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "tests", "responses")
}

// loadJSONResponse loads a JSON response from a file
func (m *SettingsCatalogConfigurationPolicyJsonMock) loadJSONResponse(relativePath string) (map[string]any, error) {
	fullPath := filepath.Join(m.getTestDataPath(), relativePath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", fullPath, err)
	}

	var response map[string]any
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON from %s: %w", fullPath, err)
	}

	return response, nil
}

// RegisterMocks registers all HTTP mock responders for settings catalog JSON policies
func (m *SettingsCatalogConfigurationPolicyJsonMock) RegisterMocks() {
	// Mock successful creation of settings catalog configuration policy
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		func(req *http.Request) (*http.Response, error) {
			response, err := m.loadJSONResponse("validate_create/post_settings_catalog_configuration_policy_json_minimal.json")
			if err != nil {
				return httpmock.NewStringResponse(500, err.Error()), nil
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Mock successful retrieval of settings catalog configuration policy
	policyIDRegex := regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/(.+)$`)
	httpmock.RegisterRegexpResponder("GET", policyIDRegex,
		func(req *http.Request) (*http.Response, error) {
			matches := policyIDRegex.FindStringSubmatch(req.URL.String())
			if len(matches) < 2 {
				return httpmock.NewStringResponse(400, "Invalid URL"), nil
			}

			response, err := m.loadJSONResponse("validate_read/get_settings_catalog_configuration_policy_json_minimal.json")
			if err != nil {
				return httpmock.NewStringResponse(500, err.Error()), nil
			}

			// Update the ID in the response to match the requested ID
			policyID := matches[1]
			response["id"] = policyID

			return httpmock.NewJsonResponse(200, response)
		})

	// Mock successful update of settings catalog configuration policy
	httpmock.RegisterRegexpResponder("PATCH", policyIDRegex,
		func(req *http.Request) (*http.Response, error) {
			matches := policyIDRegex.FindStringSubmatch(req.URL.String())
			if len(matches) < 2 {
				return httpmock.NewStringResponse(400, "Invalid URL"), nil
			}

			response, err := m.loadJSONResponse("validate_update/patch_settings_catalog_configuration_policy_json_updated.json")
			if err != nil {
				return httpmock.NewStringResponse(500, err.Error()), nil
			}

			// Update the ID in the response to match the requested ID
			policyID := matches[1]
			response["id"] = policyID

			return httpmock.NewJsonResponse(200, response)
		})

	// Mock successful deletion of settings catalog configuration policy
	httpmock.RegisterRegexpResponder(constants.TfTfOperationDelete, policyIDRegex,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		})

	fmt.Println("Settings Catalog JSON Policy mocks registered successfully")
}

// RegisterErrorMocks registers HTTP mock responders that simulate error conditions
func (m *SettingsCatalogConfigurationPolicyJsonMock) RegisterErrorMocks() {
	// Mock creation error
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		func(req *http.Request) (*http.Response, error) {
			errorResponse := map[string]any{
				"error": map[string]any{
					"code":    "BadRequest",
					"message": "Invalid request body",
					"details": []map[string]any{
						{
							"code":    "InvalidRequestBody",
							"message": "The request body is invalid or malformed",
							"target":  "settings",
						},
					},
				},
			}
			responseBody, _ := json.Marshal(errorResponse)
			return httpmock.NewBytesResponse(400, responseBody), nil
		})

	// Mock retrieval error (resource not found)
	policyIDRegex := regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/(.+)$`)
	httpmock.RegisterRegexpResponder("GET", policyIDRegex,
		func(req *http.Request) (*http.Response, error) {
			response, err := m.loadJSONResponse("validate_delete/get_settings_catalog_configuration_policy_json_not_found.json")
			if err != nil {
				return httpmock.NewStringResponse(500, err.Error()), nil
			}
			responseBody, _ := json.Marshal(response)
			return httpmock.NewBytesResponse(404, responseBody), nil
		})

	fmt.Println("Settings Catalog JSON Policy error mocks registered successfully")
}

// CleanupMockState cleans up any mock state
func (m *SettingsCatalogConfigurationPolicyJsonMock) CleanupMockState() {
	// No persistent state to clean up
	fmt.Println("Settings Catalog JSON Policy mock state cleaned up")
}
