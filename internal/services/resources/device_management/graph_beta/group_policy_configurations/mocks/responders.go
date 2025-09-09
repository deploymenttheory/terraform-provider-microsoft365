package mocks

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	groupPolicyConfigurations map[string]map[string]interface{}
}

func init() {
	// Initialize mockState
	mockState.groupPolicyConfigurations = make(map[string]map[string]interface{})

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("group_policy_configuration", &GroupPolicyConfigurationMock{})
}

// GroupPolicyConfigurationMock provides mock responses for group policy configuration operations
type GroupPolicyConfigurationMock struct{}

// Ensure GroupPolicyConfigurationMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*GroupPolicyConfigurationMock)(nil)

// RegisterMocks sets up all the mock HTTP responders for group policy configuration operations
// This implements the MockRegistrar interface
func (m *GroupPolicyConfigurationMock) RegisterMocks() {
	// POST /deviceManagement/groupPolicyConfigurations - Create
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations",
		m.createGroupPolicyConfigurationResponder())

	// GET /deviceManagement/groupPolicyConfigurations/{id} - Read
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/([^/]+)$`,
		m.getGroupPolicyConfigurationResponder())

	// PATCH /deviceManagement/groupPolicyConfigurations/{id} - Update
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/([^/]+)$`,
		m.updateGroupPolicyConfigurationResponder())

	// DELETE /deviceManagement/groupPolicyConfigurations/{id} - Delete
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/([^/]+)$`,
		m.deleteGroupPolicyConfigurationResponder())

	// POST /deviceManagement/groupPolicyConfigurations/{id}/updateDefinitionValues - Update Definition Values
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/([^/]+)/updateDefinitionValues$`,
		m.updateDefinitionValuesResponder())

	// GET /deviceManagement/groupPolicyConfigurations/{id}/definitionValues - Read Definition Values
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/([^/]+)/definitionValues`,
		m.getDefinitionValuesResponder())
}

// RegisterErrorMocks sets up error response mocks for testing error scenarios
func (m *GroupPolicyConfigurationMock) RegisterErrorMocks() {
	// POST /deviceManagement/groupPolicyConfigurations - Create Error
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations",
		factories.ErrorResponse(400, "BadRequest", "Invalid request"))

	// GET /deviceManagement/groupPolicyConfigurations/{id} - Read Error
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/([^/]+)$`,
		factories.ErrorResponse(404, "ResourceNotFound", "Resource not found"))
}

// CleanupMockState cleans up the mock state for testing
func (m *GroupPolicyConfigurationMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.groupPolicyConfigurations = make(map[string]map[string]interface{})
}

// createGroupPolicyConfigurationResponder handles POST requests to create group policy configurations
func (m *GroupPolicyConfigurationMock) createGroupPolicyConfigurationResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		// Generate a new UUID for the created resource
		id := uuid.New().String()

		// Load base response from JSON file
		response, err := mocks.LoadJSONResponse(filepath.Join("tests", "responses", "validate_create", "post_group_policy_configuration.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}

		// Generate a new ID for the created resource
		response["id"] = id

		// Update response with request data
		if displayName, ok := requestBody["displayName"]; ok {
			response["displayName"] = displayName
		}
		if description, ok := requestBody["description"]; ok {
			response["description"] = description
		}
		if roleScopeTagIds, ok := requestBody["roleScopeTagIds"]; ok {
			response["roleScopeTagIds"] = roleScopeTagIds
		}

		// Handle definitionValues if provided
		if definitionValues, ok := requestBody["definitionValues"]; ok {
			// If the request contains definitionValues, use the maximal response template
			maximalResponse, err := mocks.LoadJSONResponse(filepath.Join("tests", "responses", "validate_create", "get_group_policy_configuration_maximal.json"))
			if err == nil {
				// Merge the maximal response data
				for key, value := range maximalResponse {
					if key != "id" { // Don't overwrite the generated ID
						response[key] = value
					}
				}
			}
			// Override with actual request data
			response["definitionValues"] = definitionValues
		}

		// Store in mock state
		mockState.Lock()
		mockState.groupPolicyConfigurations[id] = response
		mockState.Unlock()

		return factories.SuccessResponse(201, response)(req)
	}
}

// getGroupPolicyConfigurationResponder handles GET requests to retrieve group policy configurations
func (m *GroupPolicyConfigurationMock) getGroupPolicyConfigurationResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/groupPolicyConfigurations/")

		mockState.Lock()
		groupPolicyConfiguration, exists := mockState.groupPolicyConfigurations[id]
		mockState.Unlock()

		if !exists {
			return factories.ErrorResponse(404, "ResourceNotFound", "Resource not found")(req)
		}

		return factories.SuccessResponse(200, groupPolicyConfiguration)(req)
	}
}

// updateGroupPolicyConfigurationResponder handles PATCH requests to update group policy configurations
func (m *GroupPolicyConfigurationMock) updateGroupPolicyConfigurationResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/groupPolicyConfigurations/")

		var requestBody map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return factories.ErrorResponse(400, "BadRequest", "Invalid JSON")(req)
		}

		mockState.Lock()
		resourceData, exists := mockState.groupPolicyConfigurations[id]
		if !exists {
			mockState.Unlock()
			return factories.ErrorResponse(404, "ResourceNotFound", "Resource not found")(req)
		}

		// Update the resource data
		for key, value := range requestBody {
			resourceData[key] = value
		}
		resourceData["lastModifiedDateTime"] = "2024-01-02T00:00:00Z"

		mockState.groupPolicyConfigurations[id] = resourceData
		mockState.Unlock()

		return factories.SuccessResponse(200, resourceData)(req)
	}
}

// deleteGroupPolicyConfigurationResponder handles DELETE requests to remove group policy configurations
func (m *GroupPolicyConfigurationMock) deleteGroupPolicyConfigurationResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/groupPolicyConfigurations/")

		mockState.Lock()
		_, exists := mockState.groupPolicyConfigurations[id]
		if !exists {
			mockState.Unlock()
			return factories.ErrorResponse(404, "ResourceNotFound", "Resource not found")(req)
		}

		delete(mockState.groupPolicyConfigurations, id)
		mockState.Unlock()

		return factories.EmptySuccessResponse(204)(req)
	}
}

// updateDefinitionValuesResponder handles POST requests to update definition values
func (m *GroupPolicyConfigurationMock) updateDefinitionValuesResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/groupPolicyConfigurations/")

		var requestBody map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return factories.ErrorResponse(400, "BadRequest", "Invalid JSON")(req)
		}

		mockState.Lock()
		_, exists := mockState.groupPolicyConfigurations[id]
		if !exists {
			mockState.Unlock()
			return factories.ErrorResponse(404, "ResourceNotFound", "Resource not found")(req)
		}
		mockState.Unlock()

		// Return success for definition values update
		return factories.SuccessResponse(200, map[string]interface{}{"value": "Success"})(req)
	}
}

// getDefinitionValuesResponder handles GET requests to retrieve definition values
func (m *GroupPolicyConfigurationMock) getDefinitionValuesResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/groupPolicyConfigurations/")

		mockState.Lock()
		_, exists := mockState.groupPolicyConfigurations[id]
		mockState.Unlock()

		if !exists {
			return factories.ErrorResponse(404, "ResourceNotFound", "Resource not found")(req)
		}

		// Return empty definition values for now
		response := map[string]interface{}{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyConfigurations('" + id + "')/definitionValues",
			"value":          []interface{}{},
		}

		return factories.SuccessResponse(200, response)(req)
	}
}
