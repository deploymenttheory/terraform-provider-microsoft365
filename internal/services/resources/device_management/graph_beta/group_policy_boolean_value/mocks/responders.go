package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	configurations     map[string]map[string]any
	definitionValues   map[string]map[string]any
	presentationValues map[string][]map[string]any
	definitions        map[string]map[string]any
	presentations      map[string][]map[string]any
}

func init() {
	// Initialize mockState
	mockState.configurations = make(map[string]map[string]any)
	mockState.definitionValues = make(map[string]map[string]any)
	mockState.presentationValues = make(map[string][]map[string]any)
	mockState.definitions = make(map[string]map[string]any)
	mockState.presentations = make(map[string][]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("group_policy_boolean_value", &GroupPolicyBooleanValueMock{})
}

// GroupPolicyBooleanValueMock provides mock responses for Group Policy Boolean Value operations
type GroupPolicyBooleanValueMock struct{}

// Ensure GroupPolicyBooleanValueMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*GroupPolicyBooleanValueMock)(nil)

// RegisterMocks registers HTTP mock responses for Group Policy Boolean Value operations
func (m *GroupPolicyBooleanValueMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.configurations = make(map[string]map[string]any)
	mockState.definitionValues = make(map[string]map[string]any)
	mockState.presentationValues = make(map[string][]map[string]any)
	mockState.definitions = make(map[string]map[string]any)
	mockState.presentations = make(map[string][]map[string]any)
	mockState.Unlock()

	// Register GET for listing group policy configurations
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			configs := make([]map[string]any, 0, len(mockState.configurations))
			for _, config := range mockState.configurations {
				configCopy := make(map[string]any)
				for k, v := range config {
					configCopy[k] = v
				}
				configs = append(configs, configCopy)
			}
			mockState.Unlock()

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyConfigurations",
				"value":          configs,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for individual group policy configuration
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			mockState.Lock()
			config, exists := mockState.configurations[id]
			mockState.Unlock()

			if !exists {
				return httpmock.NewJsonResponse(404, map[string]any{
					"error": map[string]any{
						"code":    "ResourceNotFound",
						"message": "Group policy configuration not found",
					},
				})
			}

			return httpmock.NewJsonResponse(200, config)
		})

	// Register GET for definition values
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/[^/]+/definitionValues$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			configID := urlParts[4]

			mockState.Lock()
			values := make([]map[string]any, 0)
			for _, defValue := range mockState.definitionValues {
				if defValue["configId"] == configID {
					values = append(values, defValue)
				}
			}
			mockState.Unlock()

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyConfigurations('" + configID + "')/definitionValues",
				"value":          values,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for individual definition value
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/[^/]+/definitionValues/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			defValueID := urlParts[len(urlParts)-1]

			mockState.Lock()
			defValue, exists := mockState.definitionValues[defValueID]
			mockState.Unlock()

			if !exists {
				return httpmock.NewJsonResponse(404, map[string]any{
					"error": map[string]any{
						"code":    "ResourceNotFound",
						"message": "Definition value not found",
					},
				})
			}

			return httpmock.NewJsonResponse(200, defValue)
		})

	// Register GET for presentation values
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/[^/]+/definitionValues/[^/]+/presentationValues$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			defValueID := urlParts[6]

			mockState.Lock()
			presValues := mockState.presentationValues[defValueID]
			if presValues == nil {
				presValues = []map[string]any{}
			}
			mockState.Unlock()

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyConfigurations/definitionValues/presentationValues",
				"value":          presValues,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for group policy definitions (for ID resolution)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyDefinitions$`,
		func(req *http.Request) (*http.Response, error) {
			filter := req.URL.Query().Get("$filter")

			mockState.Lock()
			definitions := make([]map[string]any, 0)
			for _, def := range mockState.definitions {
				// Simple filter matching based on displayName
				if filter != "" && strings.Contains(filter, "displayName eq") {
					displayName := extractDisplayNameFromFilter(filter)
					if defDisplayName, ok := def["displayName"].(string); ok && defDisplayName == displayName {
						definitions = append(definitions, def)
					}
				} else {
					definitions = append(definitions, def)
				}
			}
			mockState.Unlock()

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyDefinitions",
				"value":          definitions,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for presentations of a definition
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyDefinitions/[^/]+/presentations$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			defID := urlParts[4]

			mockState.Lock()
			presentations := mockState.presentations[defID]
			if presentations == nil {
				presentations = []map[string]any{}
			}
			mockState.Unlock()

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyDefinitions('" + defID + "')/presentations",
				"value":          presentations,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register POST for updateDefinitionValues
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/[^/]+/updateDefinitionValues$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			configID := urlParts[4]

			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewJsonResponse(400, map[string]any{
					"error": map[string]any{
						"code":    "BadRequest",
						"message": "Invalid request body",
					},
				})
			}

			mockState.Lock()
			defer mockState.Unlock()

			// Handle added values
			if added, ok := requestBody["added"].([]any); ok && len(added) > 0 {
				for _, item := range added {
					if addedValue, ok := item.(map[string]any); ok {
						defValueID := uuid.New().String()
						presValues := []map[string]any{}

						if presentationValues, ok := addedValue["presentationValues"].([]any); ok {
							for _, pv := range presentationValues {
								if presValue, ok := pv.(map[string]any); ok {
									presValueID := uuid.New().String()
									presValue["id"] = presValueID
									presValues = append(presValues, presValue)
								}
							}
						}

						defValue := map[string]any{
							"@odata.type":          "#microsoft.graph.groupPolicyDefinitionValue",
							"id":                   defValueID,
							"enabled":              addedValue["enabled"],
							"configId":             configID,
							"createdDateTime":      "2024-01-01T00:00:00Z",
							"lastModifiedDateTime": "2024-01-01T00:00:00Z",
						}

						mockState.definitionValues[defValueID] = defValue
						mockState.presentationValues[defValueID] = presValues
					}
				}
			}

			// Handle updated values
			if updated, ok := requestBody["updated"].([]any); ok && len(updated) > 0 {
				for _, item := range updated {
					if updatedValue, ok := item.(map[string]any); ok {
						if defValueID, ok := updatedValue["id"].(string); ok {
							if existingDefValue, exists := mockState.definitionValues[defValueID]; exists {
								existingDefValue["enabled"] = updatedValue["enabled"]
								existingDefValue["lastModifiedDateTime"] = "2024-01-02T00:00:00Z"

								if presentationValues, ok := updatedValue["presentationValues"].([]any); ok {
									presValues := []map[string]any{}
									for _, pv := range presentationValues {
										if presValue, ok := pv.(map[string]any); ok {
											if presValueID, ok := presValue["id"].(string); ok {
												presValue["id"] = presValueID
											}
											presValues = append(presValues, presValue)
										}
									}
									mockState.presentationValues[defValueID] = presValues
								}
							}
						}
					}
				}
			}

			// Handle deleted IDs
			if deletedIDs, ok := requestBody["deletedIds"].([]any); ok && len(deletedIDs) > 0 {
				for _, id := range deletedIDs {
					if idStr, ok := id.(string); ok {
						delete(mockState.definitionValues, idStr)
						delete(mockState.presentationValues, idStr)
					}
				}
			}

			return httpmock.NewStringResponse(204, ""), nil
		})
}

// CleanupMockState clears the mock state for clean test runs
func (m *GroupPolicyBooleanValueMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()

	for id := range mockState.configurations {
		delete(mockState.configurations, id)
	}
	for id := range mockState.definitionValues {
		delete(mockState.definitionValues, id)
	}
	for id := range mockState.presentationValues {
		delete(mockState.presentationValues, id)
	}
	for id := range mockState.definitions {
		delete(mockState.definitions, id)
	}
	for id := range mockState.presentations {
		delete(mockState.presentations, id)
	}
}

// RegisterErrorMocks registers mock responses that simulate error conditions
func (m *GroupPolicyBooleanValueMock) RegisterErrorMocks() {
	// Reset the state when registering error mocks
	mockState.Lock()
	mockState.configurations = make(map[string]map[string]any)
	mockState.definitionValues = make(map[string]map[string]any)
	mockState.presentationValues = make(map[string][]map[string]any)
	mockState.definitions = make(map[string]map[string]any)
	mockState.presentations = make(map[string][]map[string]any)
	mockState.Unlock()

	// Register error response for updateDefinitionValues
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/[^/]+/updateDefinitionValues$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_error.json")
			var errorResponse map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
			return httpmock.NewJsonResponse(400, errorResponse)
		})
}

// SetupMockConfiguration sets up a mock group policy configuration for testing
func (m *GroupPolicyBooleanValueMock) SetupMockConfiguration(configID, displayName string) {
	mockState.Lock()
	defer mockState.Unlock()

	config := map[string]any{
		"@odata.type":          "#microsoft.graph.groupPolicyConfiguration",
		"id":                   configID,
		"displayName":          displayName,
		"description":          "Test configuration",
		"createdDateTime":      "2024-01-01T00:00:00Z",
		"lastModifiedDateTime": "2024-01-01T00:00:00Z",
	}

	mockState.configurations[configID] = config
}

// SetupMockDefinition sets up a mock group policy definition for testing
func (m *GroupPolicyBooleanValueMock) SetupMockDefinition(defID, displayName, classType, categoryPath string, presentations []map[string]any) {
	mockState.Lock()
	defer mockState.Unlock()

	definition := map[string]any{
		"@odata.type":  "#microsoft.graph.groupPolicyDefinition",
		"id":           defID,
		"displayName":  displayName,
		"classType":    classType,
		"categoryPath": categoryPath,
	}

	mockState.definitions[defID] = definition
	mockState.presentations[defID] = presentations
}

// extractDisplayNameFromFilter extracts display name from filter string
func extractDisplayNameFromFilter(filter string) string {
	// Simple extraction: "displayName eq 'value'" -> "value"
	parts := strings.Split(filter, "'")
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}
