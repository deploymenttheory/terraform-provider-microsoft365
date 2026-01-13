package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of managed device cleanup rules for consistent responses
var mockState struct {
	sync.Mutex
	managedDeviceCleanupRules map[string]map[string]any
}

func init() {
	mockState.managedDeviceCleanupRules = make(map[string]map[string]any)

	// Default 404 for unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("managed_device_cleanup_rule", &ManagedDeviceCleanupRuleMock{})
}

// ManagedDeviceCleanupRuleMock provides mock responses for managed device cleanup rule operations
type ManagedDeviceCleanupRuleMock struct{}

// Ensure ManagedDeviceCleanupRuleMock implements MockRegistrar
var _ mocks.MockRegistrar = (*ManagedDeviceCleanupRuleMock)(nil)

// RegisterMocks registers HTTP mock responses for managed device cleanup rule operations
func (m *ManagedDeviceCleanupRuleMock) RegisterMocks() {
	mockState.Lock()
	mockState.managedDeviceCleanupRules = make(map[string]map[string]any)
	mockState.Unlock()

	// List rules
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/managedDeviceCleanupRules",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			rules := make([]map[string]any, 0, len(mockState.managedDeviceCleanupRules))
			for _, r := range mockState.managedDeviceCleanupRules {
				rules = append(rules, r)
			}
			mockState.Unlock()

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/managedDeviceCleanupRules",
				"value":          rules,
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Get rule by ID
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/managedDeviceCleanupRules/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			rule, exists := mockState.managedDeviceCleanupRules[id]
			mockState.Unlock()

			if !exists {
				switch {
				case strings.Contains(id, "minimal"):
					jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_read/get_managed_device_cleanup_rule_minimal.json")
					if err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
					}
					var response map[string]any
					if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
					}
					response["id"] = id
					return factories.SuccessResponse(200, response)(req)
				case strings.Contains(id, "maximal"):
					jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_read/get_managed_device_cleanup_rule_maximal.json")
					if err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
					}
					var response map[string]any
					if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
					}
					response["id"] = id
					return factories.SuccessResponse(200, response)(req)
				default:
					jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_resource_not_found.json")
					var errorResponse map[string]any
					_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
					return httpmock.NewJsonResponse(404, errorResponse)
				}
			}

			// Ensure assignments or other optional arrays are returned consistently if needed (none for this resource)
			return httpmock.NewJsonResponse(200, rule)
		})

	// Create rule
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/managedDeviceCleanupRules",
		func(req *http.Request) (*http.Response, error) {
			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Enforce uniqueness by deviceCleanupRulePlatformType (API returns 500 on duplicate)
			platform := ""
			if v, ok := body["deviceCleanupRulePlatformType"].(string); ok {
				platform = v
			}
			mockState.Lock()
			for _, existing := range mockState.managedDeviceCleanupRules {
				if p, ok := existing["deviceCleanupRulePlatformType"].(string); ok && p == platform {
					mockState.Unlock()
					jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_duplicate_platform.json")
					var errorResponse map[string]any
					_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
					return httpmock.NewJsonResponse(500, errorResponse)
				}
			}
			mockState.Unlock()

			id := uuid.New().String()

			rule := map[string]any{
				"id":                                     id,
				"displayName":                            body["displayName"],
				"description":                            body["description"],
				"deviceCleanupRulePlatformType":          body["deviceCleanupRulePlatformType"],
				"deviceInactivityBeforeRetirementInDays": body["deviceInactivityBeforeRetirementInDays"],
				"lastModifiedDateTime":                   "2024-01-01T00:00:00Z",
			}

			mockState.Lock()
			mockState.managedDeviceCleanupRules[id] = rule
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, rule)
		})

	// Update rule
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/managedDeviceCleanupRules/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_invalid_display_name.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(400, errorResponse)
			}

			mockState.Lock()
			existing, exists := mockState.managedDeviceCleanupRules[id]
			if !exists {
				mockState.Unlock()
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_resource_not_found.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}
			// Load update template and merge
			if jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_update/patch_managed_device_cleanup_rule_minimal.json"); err == nil {
				var base map[string]any
				if json.Unmarshal([]byte(jsonStr), &base) == nil {
					for k, v := range existing {
						base[k] = v
					}
					for k, v := range body {
						base[k] = v
					}
					existing = base
				}
			} else {
				for k, v := range body {
					existing[k] = v
				}
			}
			// Simulate server updating last modified time
			existing["lastModifiedDateTime"] = "2024-01-02T00:00:00Z"
			mockState.managedDeviceCleanupRules[id] = existing
			mockState.Unlock()

			return factories.SuccessResponse(200, existing)(req)
		})

	// Delete rule
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/managedDeviceCleanupRules/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			if _, exists := mockState.managedDeviceCleanupRules[id]; exists {
				delete(mockState.managedDeviceCleanupRules, id)
				mockState.Unlock()
				return httpmock.NewStringResponse(204, ""), nil
			}
			mockState.Unlock()
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		})
}

// RegisterErrorMocks registers error responses to simulate failures
func (m *ManagedDeviceCleanupRuleMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.managedDeviceCleanupRules = make(map[string]map[string]any)
	mockState.Unlock()

	// Create error
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/managedDeviceCleanupRules",
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_invalid_display_name.json")
			var errorResponse map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
			return httpmock.NewJsonResponse(400, errorResponse)
		})

	// Get by id not found
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/managedDeviceCleanupRules/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_resource_not_found.json")
			var errorResponse map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
			return httpmock.NewJsonResponse(404, errorResponse)
		})
}

// CleanupMockState clears the mock state for clean test runs
func (m *ManagedDeviceCleanupRuleMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.managedDeviceCleanupRules {
		delete(mockState.managedDeviceCleanupRules, id)
	}
}
