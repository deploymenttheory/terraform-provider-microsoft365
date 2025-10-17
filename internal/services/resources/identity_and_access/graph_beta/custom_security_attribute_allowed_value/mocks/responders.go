package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	// Key is definitionId/allowedValueId
	allowedValues map[string]map[string]any
	// Track count per definition for validation
	definitionCounts map[string]int
}

func init() {
	// Initialize mockState
	mockState.allowedValues = make(map[string]map[string]any)
	mockState.definitionCounts = make(map[string]int)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("custom_security_attribute_allowed_value", &CustomSecurityAttributeAllowedValueMock{})
}

// CustomSecurityAttributeAllowedValueMock provides mock responses for Allowed Value operations
type CustomSecurityAttributeAllowedValueMock struct{}

// Ensure CustomSecurityAttributeAllowedValueMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*CustomSecurityAttributeAllowedValueMock)(nil)

// RegisterMocks registers HTTP mock responses for Allowed Value operations
func (m *CustomSecurityAttributeAllowedValueMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.allowedValues = make(map[string]map[string]any)
	mockState.definitionCounts = make(map[string]int)
	mockState.Unlock()

	// Register GET for listing Allowed Values (used in validation)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directory/customSecurityAttributeDefinitions/[a-zA-Z0-9_-]+/allowedValues$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract definition ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			definitionId := urlParts[len(urlParts)-2]

			mockState.Lock()
			values := make([]map[string]any, 0)
			for key, value := range mockState.allowedValues {
				// Check if this allowed value belongs to the requested definition
				if strings.HasPrefix(key, definitionId+"/") {
					valueCopy := make(map[string]any)
					for k, v := range value {
						valueCopy[k] = v
					}
					values = append(values, valueCopy)
				}
			}
			mockState.Unlock()

			response := map[string]any{
				"@odata.context": fmt.Sprintf("https://graph.microsoft.com/beta/$metadata#directory/customSecurityAttributeDefinitions('%s')/allowedValues", definitionId),
				"value":          values,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for individual Allowed Value
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directory/customSecurityAttributeDefinitions/[a-zA-Z0-9_-]+/allowedValues/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract definition ID and value ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			definitionId := urlParts[len(urlParts)-3]
			valueId := urlParts[len(urlParts)-1]
			key := fmt.Sprintf("%s/%s", definitionId, valueId)

			mockState.Lock()
			value, exists := mockState.allowedValues[key]
			mockState.Unlock()

			if !exists {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/get_allowed_value_not_found.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			// Create response copy
			valueCopy := make(map[string]any)
			for k, v := range value {
				valueCopy[k] = v
			}

			return httpmock.NewJsonResponse(200, valueCopy)
		})

	// Register POST for creating Allowed Value
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/directory/customSecurityAttributeDefinitions/[a-zA-Z0-9_-]+/allowedValues$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract definition ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			definitionId := urlParts[len(urlParts)-2]

			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Check if we've hit the limit for this definition
			mockState.Lock()
			count := mockState.definitionCounts[definitionId]
			mockState.Unlock()

			if count >= 100 {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Cannot create new allowed value: maximum limit of 100 allowed values reached"}}`), nil
			}

			// Load the base response template
			jsonStr, err := helpers.ParseJSONFile("../tests/responses/post_allowed_value.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var response map[string]any
			if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
			}

			// Use the data from the request
			if id, ok := requestBody["id"].(string); ok {
				response["id"] = id
			}
			if isActive, ok := requestBody["isActive"].(bool); ok {
				response["isActive"] = isActive
			}

			// Store in mock state
			key := fmt.Sprintf("%s/%s", definitionId, response["id"].(string))
			mockState.Lock()
			mockState.allowedValues[key] = response
			mockState.definitionCounts[definitionId]++
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, response)
		})

	// Register PATCH for updating Allowed Value
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/directory/customSecurityAttributeDefinitions/[a-zA-Z0-9_-]+/allowedValues/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract definition ID and value ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			definitionId := urlParts[len(urlParts)-3]
			valueId := urlParts[len(urlParts)-1]
			key := fmt.Sprintf("%s/%s", definitionId, valueId)

			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/post_allowed_value_error.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(400, errorResponse)
			}

			// Load update template
			jsonStr, err := helpers.ParseJSONFile("../tests/responses/get_allowed_value_updated.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var updatedValue map[string]any
			if err := json.Unmarshal([]byte(jsonStr), &updatedValue); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
			}

			mockState.Lock()
			value, exists := mockState.allowedValues[key]
			if !exists {
				mockState.Unlock()
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/get_allowed_value_not_found.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			// Start with existing data
			for k, v := range value {
				updatedValue[k] = v
			}

			// Apply updates from request body
			for k, v := range requestBody {
				updatedValue[k] = v
			}

			// Store updated state
			mockState.allowedValues[key] = updatedValue
			mockState.Unlock()

			return factories.SuccessResponse(200, updatedValue)(req)
		})
}

// CleanupMockState clears the mock state for clean test runs
func (m *CustomSecurityAttributeAllowedValueMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()

	// Clear all stored Allowed Values
	for key := range mockState.allowedValues {
		delete(mockState.allowedValues, key)
	}
	for key := range mockState.definitionCounts {
		delete(mockState.definitionCounts, key)
	}
}

// RegisterErrorMocks registers mock responses that simulate error conditions
func (m *CustomSecurityAttributeAllowedValueMock) RegisterErrorMocks() {
	// Reset the state when registering error mocks
	mockState.Lock()
	mockState.allowedValues = make(map[string]map[string]any)
	mockState.definitionCounts = make(map[string]int)
	mockState.Unlock()

	// Register GET for listing Allowed Values (needed for validation)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directory/customSecurityAttributeDefinitions/[a-zA-Z0-9_-]+/allowedValues$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract definition ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			definitionId := urlParts[len(urlParts)-2]

			response := map[string]any{
				"@odata.context": fmt.Sprintf("https://graph.microsoft.com/beta/$metadata#directory/customSecurityAttributeDefinitions('%s')/allowedValues", definitionId),
				"value":          []map[string]any{}, // Empty list for error scenarios
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Register error response for creating Allowed Value with invalid data
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/directory/customSecurityAttributeDefinitions/[a-zA-Z0-9_-]+/allowedValues$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/post_allowed_value_error.json")
			var errorResponse map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
			return httpmock.NewJsonResponse(400, errorResponse)
		})

	// Register error response for Allowed Value not found
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directory/customSecurityAttributeDefinitions/[a-zA-Z0-9_-]+/allowedValues/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/get_allowed_value_not_found.json")
			var errorResponse map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
			return httpmock.NewJsonResponse(404, errorResponse)
		})
}
