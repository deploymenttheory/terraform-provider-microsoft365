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
	attributeSets map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.attributeSets = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("attribute_set", &AttributeSetMock{})
}

// AttributeSetMock provides mock responses for Attribute Set operations
type AttributeSetMock struct{}

// Ensure AttributeSetMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*AttributeSetMock)(nil)

// RegisterMocks registers HTTP mock responses for Attribute Set operations
func (m *AttributeSetMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.attributeSets = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for listing Attribute Sets (used in validation)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/directory/attributeSets",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			sets := make([]map[string]any, 0, len(mockState.attributeSets))
			for _, set := range mockState.attributeSets {
				// Ensure @odata.type is present
				setCopy := make(map[string]any)
				for k, v := range set {
					setCopy[k] = v
				}
				if _, hasODataType := setCopy["@odata.type"]; !hasODataType {
					setCopy["@odata.type"] = "#microsoft.graph.attributeSet"
				}
				sets = append(sets, setCopy)
			}
			mockState.Unlock()

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#directory/attributeSets",
				"value":          sets,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for individual Attribute Set
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directory/attributeSets/[a-zA-Z0-9_-]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			mockState.Lock()
			set, exists := mockState.attributeSets[id]
			mockState.Unlock()

			if !exists {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_attribute_set_not_found.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			// Create response copy
			setCopy := make(map[string]any)
			for k, v := range set {
				setCopy[k] = v
			}
			if _, hasODataType := setCopy["@odata.type"]; !hasODataType {
				setCopy["@odata.type"] = "#microsoft.graph.attributeSet"
			}

			return httpmock.NewJsonResponse(200, setCopy)
		})

	// Register POST for creating Attribute Set
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/directory/attributeSets",
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Load the base response template
			jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_create/post_attribute_set.json")
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
			if description, ok := requestBody["description"].(string); ok {
				response["description"] = description
			} else {
				delete(response, "description")
			}
			if maxAttrs, ok := requestBody["maxAttributesPerSet"].(float64); ok {
				response["maxAttributesPerSet"] = int32(maxAttrs)
			} else {
				delete(response, "maxAttributesPerSet")
			}

			// Store in mock state
			mockState.Lock()
			mockState.attributeSets[response["id"].(string)] = response
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, response)
		})

	// Register PATCH for updating Attribute Set
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/directory/attributeSets/[a-zA-Z0-9_-]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_attribute_set_error.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(400, errorResponse)
			}

			// Load update template
			jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_update/get_attribute_set_updated.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}
			var updatedSet map[string]any
			if err := json.Unmarshal([]byte(jsonStr), &updatedSet); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
			}

			mockState.Lock()
			set, exists := mockState.attributeSets[id]
			if !exists {
				mockState.Unlock()
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_attribute_set_not_found.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			// Start with existing data
			for k, v := range set {
				updatedSet[k] = v
			}

			// Apply updates from request body
			for k, v := range requestBody {
				updatedSet[k] = v
			}

			// Store updated state
			mockState.attributeSets[id] = updatedSet
			mockState.Unlock()

			return factories.SuccessResponse(200, updatedSet)(req)
		})

	// Register DELETE for Attribute Set
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/directory/attributeSets/[a-zA-Z0-9_-]+$`,
		func(req *http.Request) (*http.Response, error) {
			// Extract ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			id := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.attributeSets[id]
			if exists {
				delete(mockState.attributeSets, id)
			}
			mockState.Unlock()

			if !exists {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_attribute_set_not_found.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			return httpmock.NewStringResponse(204, ""), nil
		})
}

// CleanupMockState clears the mock state for clean test runs
func (m *AttributeSetMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()

	// Clear all stored Attribute Sets
	for id := range mockState.attributeSets {
		delete(mockState.attributeSets, id)
	}
}

// RegisterErrorMocks registers mock responses that simulate error conditions
func (m *AttributeSetMock) RegisterErrorMocks() {
	// Reset the state when registering error mocks
	mockState.Lock()
	mockState.attributeSets = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for listing Attribute Sets (needed for validation)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/directory/attributeSets",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#directory/attributeSets",
				"value":          []map[string]any{}, // Empty list for error scenarios
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Register error response for creating Attribute Set with invalid data
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/directory/attributeSets",
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_attribute_set_error.json")
			var errorResponse map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
			return httpmock.NewJsonResponse(400, errorResponse)
		})

	// Register error response for Attribute Set not found
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directory/attributeSets/[a-zA-Z0-9_-]+$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_attribute_set_not_found.json")
			var errorResponse map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
			return httpmock.NewJsonResponse(404, errorResponse)
		})
}
