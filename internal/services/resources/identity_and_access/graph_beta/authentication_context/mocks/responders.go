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
	authenticationContexts map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.authenticationContexts = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("authentication_context", &AuthenticationContextMock{})
}

// AuthenticationContextMock provides mock responses for Authentication Context operations
type AuthenticationContextMock struct{}

// Ensure AuthenticationContextMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*AuthenticationContextMock)(nil)

// RegisterMocks registers HTTP mock responses for Authentication Context operations
func (m *AuthenticationContextMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.authenticationContexts = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for listing Authentication Contexts (used in validation)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/identity/conditionalAccess/authenticationContextClassReferences",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			contexts := make([]map[string]any, 0, len(mockState.authenticationContexts))
			for _, context := range mockState.authenticationContexts {
				// Ensure @odata.type is present
				contextCopy := make(map[string]any)
				for k, v := range context {
					contextCopy[k] = v
				}
				if _, hasODataType := contextCopy["@odata.type"]; !hasODataType {
					contextCopy["@odata.type"] = "#microsoft.graph.authenticationContextClassReference"
				}
				contexts = append(contexts, contextCopy)
			}
			mockState.Unlock()

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#identity/conditionalAccess/authenticationContextClassReferences",
				"value":          contexts,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for individual Authentication Context
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/authenticationContextClassReferences/c[0-9]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			mockState.Lock()
			context, exists := mockState.authenticationContexts[id]
			mockState.Unlock()

			if !exists {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_authentication_context_not_found.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			// Create response copy
			contextCopy := make(map[string]any)
			for k, v := range context {
				contextCopy[k] = v
			}
			if _, hasODataType := contextCopy["@odata.type"]; !hasODataType {
				contextCopy["@odata.type"] = "#microsoft.graph.authenticationContextClassReference"
			}

			return httpmock.NewJsonResponse(200, contextCopy)
		})

	// Register POST for creating Authentication Context
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/identity/conditionalAccess/authenticationContextClassReferences",
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Load the base response template
			jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_create/post_authentication_context.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var response map[string]any
			if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
			}

			// Use the ID from the request
			if id, ok := requestBody["id"].(string); ok {
				response["id"] = id
			}
			if displayName, ok := requestBody["displayName"].(string); ok {
				response["displayName"] = displayName
			}
			if description, ok := requestBody["description"].(string); ok {
				response["description"] = description
			}
			if isAvailable, ok := requestBody["isAvailable"].(bool); ok {
				response["isAvailable"] = isAvailable
			}

			// Store in mock state
			mockState.Lock()
			mockState.authenticationContexts[response["id"].(string)] = response
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, response)
		})

	// Register PATCH for updating Authentication Context
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/authenticationContextClassReferences/c[0-9]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_authentication_context_error.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(400, errorResponse)
			}

			// Load update template
			jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_update/get_authentication_context_updated.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var updatedContext map[string]any
			if err := json.Unmarshal([]byte(jsonStr), &updatedContext); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
			}

			mockState.Lock()
			context, exists := mockState.authenticationContexts[id]
			if !exists {
				mockState.Unlock()
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_authentication_context_not_found.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			// Start with existing data
			for k, v := range context {
				updatedContext[k] = v
			}

			// Apply updates from request body
			for k, v := range requestBody {
				updatedContext[k] = v
			}

			// Store updated state
			mockState.authenticationContexts[id] = updatedContext
			mockState.Unlock()

			return factories.SuccessResponse(200, updatedContext)(req)
		})

	// Register DELETE for Authentication Context
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/authenticationContextClassReferences/c[0-9]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.authenticationContexts[id]
			if exists {
				delete(mockState.authenticationContexts, id)
			}
			mockState.Unlock()

			if !exists {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_authentication_context_not_found.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			return httpmock.NewStringResponse(204, ""), nil
		})
}

// CleanupMockState clears the mock state for clean test runs
func (m *AuthenticationContextMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()

	// Clear all stored Authentication Contexts
	for id := range mockState.authenticationContexts {
		delete(mockState.authenticationContexts, id)
	}
}

// RegisterErrorMocks registers mock responses that simulate error conditions
func (m *AuthenticationContextMock) RegisterErrorMocks() {
	// Reset the state when registering error mocks
	mockState.Lock()
	mockState.authenticationContexts = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for listing Authentication Contexts (needed for validation)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/identity/conditionalAccess/authenticationContextClassReferences",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#identity/conditionalAccess/authenticationContextClassReferences",
				"value":          []map[string]any{}, // Empty list for error scenarios
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Register error response for creating Authentication Context with invalid data
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/identity/conditionalAccess/authenticationContextClassReferences",
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_authentication_context_error.json")
			var errorResponse map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
			return httpmock.NewJsonResponse(400, errorResponse)
		})

	// Register error response for Authentication Context not found
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/authenticationContextClassReferences/c[0-9]+$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_authentication_context_not_found.json")
			var errorResponse map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
			return httpmock.NewJsonResponse(404, errorResponse)
		})
}