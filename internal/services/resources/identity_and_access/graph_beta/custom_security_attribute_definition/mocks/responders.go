package mocks

import (
	"encoding/json"
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
	definitions map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.definitions = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("custom_security_attribute_definition", &CustomSecurityAttributeDefinitionMock{})
}

// CustomSecurityAttributeDefinitionMock provides mock responses for Custom Security Attribute Definition operations
type CustomSecurityAttributeDefinitionMock struct{}

// Ensure CustomSecurityAttributeDefinitionMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*CustomSecurityAttributeDefinitionMock)(nil)

// RegisterMocks registers HTTP mock responses for Custom Security Attribute Definition operations
func (m *CustomSecurityAttributeDefinitionMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.definitions = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for listing Custom Security Attribute Definitions
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/directory/customSecurityAttributeDefinitions",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			defs := make([]map[string]any, 0, len(mockState.definitions))
			for _, def := range mockState.definitions {
				// Ensure @odata.type is present
				defCopy := make(map[string]any)
				for k, v := range def {
					defCopy[k] = v
				}
				if _, hasODataType := defCopy["@odata.type"]; !hasODataType {
					defCopy["@odata.type"] = "#microsoft.graph.customSecurityAttributeDefinition"
				}
				defs = append(defs, defCopy)
			}
			mockState.Unlock()

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#directory/customSecurityAttributeDefinitions",
				"value":          defs,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for individual Custom Security Attribute Definition
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directory/customSecurityAttributeDefinitions/[a-zA-Z0-9_-]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			mockState.Lock()
			def, exists := mockState.definitions[id]
			mockState.Unlock()

			if !exists {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_definition_not_found.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			// Create response copy
			defCopy := make(map[string]any)
			for k, v := range def {
				defCopy[k] = v
			}
			if _, hasODataType := defCopy["@odata.type"]; !hasODataType {
				defCopy["@odata.type"] = "#microsoft.graph.customSecurityAttributeDefinition"
			}

			return httpmock.NewJsonResponse(200, defCopy)
		})

	// Register POST for creating Custom Security Attribute Definition
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/directory/customSecurityAttributeDefinitions",
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Load the base response template
			jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_create/post_definition.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var response map[string]any
			if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
			}

			// Build the ID from attribute_set and name
			attributeSet := requestBody["attributeSet"].(string)
			name := requestBody["name"].(string)
			id := attributeSet + "_" + name

			response["id"] = id
			response["attributeSet"] = attributeSet
			response["name"] = name

			// Copy request fields to response
			if description, ok := requestBody["description"].(string); ok {
				response["description"] = description
			} else {
				delete(response, "description")
			}
			if typeVal, ok := requestBody["type"].(string); ok {
				response["type"] = typeVal
			}
			if status, ok := requestBody["status"].(string); ok {
				response["status"] = status
			}
			if isCollection, ok := requestBody["isCollection"].(bool); ok {
				response["isCollection"] = isCollection
			}
			if isSearchable, ok := requestBody["isSearchable"].(bool); ok {
				response["isSearchable"] = isSearchable
			}
			if usePredefined, ok := requestBody["usePreDefinedValuesOnly"].(bool); ok {
				response["usePreDefinedValuesOnly"] = usePredefined
			}

			// Store in mock state
			mockState.Lock()
			mockState.definitions[id] = response
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, response)
		})

	// Register PATCH for updating Custom Security Attribute Definition
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/directory/customSecurityAttributeDefinitions/[a-zA-Z0-9_-]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_definition_error.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(400, errorResponse)
			}

			mockState.Lock()
			def, exists := mockState.definitions[id]
			if !exists {
				mockState.Unlock()
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_definition_not_found.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			// Start with existing data
			updatedDef := make(map[string]any)
			for k, v := range def {
				updatedDef[k] = v
			}

			// Apply updates from request body
			for k, v := range requestBody {
				updatedDef[k] = v
			}

			// Store updated state
			mockState.definitions[id] = updatedDef
			mockState.Unlock()

			return factories.SuccessResponse(200, updatedDef)(req)
		})
}

// CleanupMockState clears the mock state for clean test runs
func (m *CustomSecurityAttributeDefinitionMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()

	// Clear all stored definitions
	for id := range mockState.definitions {
		delete(mockState.definitions, id)
	}
}

// RegisterErrorMocks registers mock responses that simulate error conditions
func (m *CustomSecurityAttributeDefinitionMock) RegisterErrorMocks() {
	// Reset the state when registering error mocks
	mockState.Lock()
	mockState.definitions = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for listing Custom Security Attribute Definitions (needed for validation)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/directory/customSecurityAttributeDefinitions",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#directory/customSecurityAttributeDefinitions",
				"value":          []map[string]any{}, // Empty list for error scenarios
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Register error response for creating Custom Security Attribute Definition with invalid data
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/directory/customSecurityAttributeDefinitions",
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_definition_error.json")
			var errorResponse map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
			return httpmock.NewJsonResponse(400, errorResponse)
		})

	// Register error response for Custom Security Attribute Definition not found
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directory/customSecurityAttributeDefinitions/[a-zA-Z0-9_-]+$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_definition_not_found.json")
			var errorResponse map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
			return httpmock.NewJsonResponse(404, errorResponse)
		})
}
