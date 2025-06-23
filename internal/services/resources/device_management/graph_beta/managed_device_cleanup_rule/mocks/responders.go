package mocks

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/jarcoal/httpmock"
)

// ManagedDeviceCleanupRuleMock provides mock responses for managed device cleanup rule operations
type ManagedDeviceCleanupRuleMock struct{}

// RegisterMocks registers HTTP mock responses for managed device cleanup rule operations
func (m *ManagedDeviceCleanupRuleMock) RegisterMocks() {
	// GET Read - Basic/Default rule
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/managedDeviceCleanupRules/00000000-0000-0000-0000-000000000001",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]interface{}{
				"id":                                     "00000000-0000-0000-0000-000000000001",
				"displayName":                            "Test Cleanup Rule",
				"description":                            "Test description",
				"deviceCleanupRulePlatformType":          "windows",
				"deviceInactivityBeforeRetirementInDays": 90,
				"lastModifiedDateTime":                   "2023-11-01T10:30:00.0000000Z",
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// GET Read - Updated rule
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/managedDeviceCleanupRules/00000000-0000-0000-0000-000000000003",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]interface{}{
				"id":                                     "00000000-0000-0000-0000-000000000003",
				"displayName":                            "Updated Cleanup Rule",
				"description":                            "Updated description",
				"deviceCleanupRulePlatformType":          "macOS",
				"deviceInactivityBeforeRetirementInDays": 120,
				"lastModifiedDateTime":                   "2023-11-02T15:45:00.0000000Z",
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// GET Read - Minimal rule
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/managedDeviceCleanupRules/00000000-0000-0000-0000-000000000005",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]interface{}{
				"id":                                     "00000000-0000-0000-0000-000000000005",
				"displayName":                            "Minimal Cleanup Rule",
				"deviceCleanupRulePlatformType":          "all",
				"deviceInactivityBeforeRetirementInDays": 30,
				"lastModifiedDateTime":                   "2023-11-01T10:30:00.0000000Z",
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// GET Read - Maximal rule
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/managedDeviceCleanupRules/00000000-0000-0000-0000-000000000006",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]interface{}{
				"id":                                     "00000000-0000-0000-0000-000000000006",
				"displayName":                            "Maximal Cleanup Rule",
				"description":                            "This is a comprehensive cleanup rule with all fields populated",
				"deviceCleanupRulePlatformType":          "ios",
				"deviceInactivityBeforeRetirementInDays": 180,
				"lastModifiedDateTime":                   "2023-11-01T10:30:00.0000000Z",
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// POST Create
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/managedDeviceCleanupRules",
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, "Invalid request body"), nil
			}

			// Generate a response based on the request
			id := "00000000-0000-0000-0000-000000000001" // Default ID
			displayName := requestBody["displayName"].(string)

			// Assign different IDs based on display name for different test cases
			if displayName == "Minimal Cleanup Rule" {
				id = "00000000-0000-0000-0000-000000000005"
			} else if displayName == "Maximal Cleanup Rule" {
				id = "00000000-0000-0000-0000-000000000006"
			} else if displayName == "Updated Cleanup Rule" {
				id = "00000000-0000-0000-0000-000000000003"
			}

			// Create response with the same fields as the request plus ID and dates
			response := requestBody
			response["id"] = id
			response["lastModifiedDateTime"] = "2023-11-01T10:30:00.0000000Z"

			return httpmock.NewJsonResponse(201, response)
		})

	// PATCH Update
	httpmock.RegisterRegexpResponder("PATCH", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/managedDeviceCleanupRules/[0-9a-f-]+`),
		func(req *http.Request) (*http.Response, error) {
			// Parse the request body to get the updated fields
			var requestBody map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, "Invalid request body"), nil
			}

			// Return success with no content
			return httpmock.NewStringResponse(204, ""), nil
		})

	// DELETE
	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/managedDeviceCleanupRules/[0-9a-f-]+`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(204, ""), nil
		})
}

// RegisterErrorMocks registers HTTP mock responses that return errors
func (m *ManagedDeviceCleanupRuleMock) RegisterErrorMocks() {
	// Register mocks that return errors
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/managedDeviceCleanupRules/00000000-0000-0000-0000-000000000001",
		factories.ErrorResponse(403, "Forbidden", "Access denied"))

	// POST Create with error
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/managedDeviceCleanupRules",
		factories.ErrorResponse(403, "Forbidden", "Access denied"))

	// PATCH Update with error
	httpmock.RegisterRegexpResponder("PATCH", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/managedDeviceCleanupRules/[0-9a-f-]+`),
		factories.ErrorResponse(403, "Forbidden", "Access denied"))

	// DELETE with error
	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/managedDeviceCleanupRules/[0-9a-f-]+`),
		factories.ErrorResponse(403, "Forbidden", "Access denied"))
}

func init() {
	// Register with the global registry
	mocks.GlobalRegistry.Register("managed_device_cleanup_rule", &ManagedDeviceCleanupRuleMock{})
}
