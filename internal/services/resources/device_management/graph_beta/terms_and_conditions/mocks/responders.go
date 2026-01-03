package mocks

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	termsAndConditions map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.termsAndConditions = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("terms_and_conditions", &TermsAndConditionsMock{})
}

// TermsAndConditionsMock provides mock responses for terms and conditions operations
type TermsAndConditionsMock struct{}

// Ensure TermsAndConditionsMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*TermsAndConditionsMock)(nil)

// RegisterMocks sets up all the mock HTTP responders for terms and conditions operations
// This implements the MockRegistrar interface
func (m *TermsAndConditionsMock) RegisterMocks() {
	// POST /deviceManagement/termsAndConditions - Create
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/termsAndConditions",
		m.createTermsAndConditionsResponder())

	// GET /deviceManagement/termsAndConditions/{id} - Read
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/termsAndConditions/([^/]+)$`,
		m.getTermsAndConditionsResponder())

	// PATCH /deviceManagement/termsAndConditions/{id} - Update
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/termsAndConditions/([^/]+)$`,
		m.updateTermsAndConditionsResponder())

	// DELETE /deviceManagement/termsAndConditions/{id} - Delete
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/termsAndConditions/([^/]+)$`,
		m.deleteTermsAndConditionsResponder())

	// GET /deviceManagement/termsAndConditions/{id}/assignments - Get assignments
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/termsAndConditions/([^/]+)/assignments$`,
		m.getTermsAndConditionsAssignmentsResponder())

	// POST /deviceManagement/termsAndConditions/{id}/assignments - Create assignment
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/termsAndConditions/([^/]+)/assignments$`,
		m.createTermsAndConditionsAssignmentResponder())
}

// createTermsAndConditionsResponder handles POST requests to create terms and conditions
func (m *TermsAndConditionsMock) createTermsAndConditionsResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		// Load base response from JSON file - use minimal if no description provided
		var jsonContent string
		var err error
		if description, hasDesc := requestBody["description"]; hasDesc && description != "" {
			jsonContent, err = helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_terms_and_conditions_maximal.json"))
		} else {
			jsonContent, err = helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_terms_and_conditions_minimal.json"))
		}

		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}

		var response map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
		}

		// Generate a new ID for the created resource
		id := uuid.New().String()
		response["id"] = id

		// Update response with request data
		if displayName, ok := requestBody["displayName"]; ok {
			response["displayName"] = displayName
		}
		if description, ok := requestBody["description"]; ok {
			response["description"] = description
		}
		if title, ok := requestBody["title"]; ok {
			response["title"] = title
		}
		if bodyText, ok := requestBody["bodyText"]; ok {
			response["bodyText"] = bodyText
		}
		if acceptanceStatement, ok := requestBody["acceptanceStatement"]; ok {
			response["acceptanceStatement"] = acceptanceStatement
		}
		if version, ok := requestBody["version"]; ok {
			response["version"] = version
		}
		if roleScopeTagIds, ok := requestBody["roleScopeTagIds"]; ok {
			response["roleScopeTagIds"] = roleScopeTagIds
		}

		// Store in mock state
		mockState.Lock()
		mockState.termsAndConditions[id] = response
		mockState.Unlock()

		return factories.SuccessResponse(201, response)(req)
	}
}

// getTermsAndConditionsResponder handles GET requests to retrieve terms and conditions
func (m *TermsAndConditionsMock) getTermsAndConditionsResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/termsAndConditions/")

		mockState.Lock()
		termsAndConditions, exists := mockState.termsAndConditions[id]
		mockState.Unlock()

		if !exists {
			// Check for special test IDs
			switch {
			case strings.Contains(id, "minimal"):
				jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_terms_and_conditions_minimal.json"))
				if err != nil {
					return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
				}
				var response map[string]any
				if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
					return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
				}
				response["id"] = id
				return factories.SuccessResponse(200, response)(req)
			case strings.Contains(id, "maximal"):
				jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_terms_and_conditions_maximal.json"))
				if err != nil {
					return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
				}
				var response map[string]any
				if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
					return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
				}
				response["id"] = id
				return factories.SuccessResponse(200, response)(req)
			default:
				jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_terms_and_conditions_not_found.json"))
				if err == nil {
					var errorResponse map[string]any
					if json.Unmarshal([]byte(jsonContent), &errorResponse) == nil {
						return httpmock.NewJsonResponse(404, errorResponse)
					}
				}
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
			}
		}

		return factories.SuccessResponse(200, termsAndConditions)(req)
	}
}

// updateTermsAndConditionsResponder handles PATCH requests to update terms and conditions
func (m *TermsAndConditionsMock) updateTermsAndConditionsResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/termsAndConditions/")

		mockState.Lock()
		termsAndConditions, exists := mockState.termsAndConditions[id]
		mockState.Unlock()

		if !exists {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_terms_and_conditions_not_found.json"))
			if err == nil {
				var errorResponse map[string]any
				if json.Unmarshal([]byte(jsonContent), &errorResponse) == nil {
					return httpmock.NewJsonResponse(404, errorResponse)
				}
			}
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewJsonResponse(400, map[string]any{
				"error": map[string]any{
					"code":    "BadRequest",
					"message": "Invalid JSON",
				},
			})
		}

		// Load update template
		jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_update", "get_terms_and_conditions_updated.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}

		var updatedTermsAndConditions map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &updatedTermsAndConditions); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
		}

		// Start with existing data
		for k, v := range termsAndConditions {
			updatedTermsAndConditions[k] = v
		}

		// Apply updates from request body
		for k, v := range requestBody {
			updatedTermsAndConditions[k] = v
		}

		// Update lastModifiedDateTime
		updatedTermsAndConditions["lastModifiedDateTime"] = "2024-01-01T12:00:00Z"

		// Store updated version
		mockState.Lock()
		mockState.termsAndConditions[id] = updatedTermsAndConditions
		mockState.Unlock()

		return factories.SuccessResponse(200, updatedTermsAndConditions)(req)
	}
}

// deleteTermsAndConditionsResponder handles DELETE requests to remove terms and conditions
func (m *TermsAndConditionsMock) deleteTermsAndConditionsResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/termsAndConditions/")

		mockState.Lock()
		_, exists := mockState.termsAndConditions[id]
		if exists {
			delete(mockState.termsAndConditions, id)
		}
		mockState.Unlock()

		if !exists {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_terms_and_conditions_not_found.json"))
			if err == nil {
				var errorResponse map[string]any
				if json.Unmarshal([]byte(jsonContent), &errorResponse) == nil {
					return httpmock.NewJsonResponse(404, errorResponse)
				}
			}
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return factories.EmptySuccessResponse(204)(req)
	}
}

// RegisterErrorMocks sets up mock responders that return errors for testing error scenarios
// This implements the MockRegistrar interface
func (m *TermsAndConditionsMock) RegisterErrorMocks() {
	// POST - Create error
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/termsAndConditions",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(400, map[string]any{
				"error": map[string]any{
					"code":    "BadRequest",
					"message": "Invalid terms and conditions data",
				},
			})
		})

	// GET - Read error (simulates not found or access denied)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/termsAndConditions/error-id$`,
		factories.ErrorResponse(403, "Forbidden", "Access denied"))

	// PATCH - Update error
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/termsAndConditions/error-id$`,
		factories.ErrorResponse(500, "InternalServerError", "Internal server error"))

	// DELETE - Delete error
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/termsAndConditions/error-id$`,
		factories.ErrorResponse(409, "Conflict", "Terms and conditions is in use"))
}

// CleanupMockState clears all stored mock state
func (m *TermsAndConditionsMock) CleanupMockState() {
	mockState.Lock()
	mockState.termsAndConditions = make(map[string]map[string]any)
	mockState.Unlock()
}

// GetMockTermsAndConditionsData returns sample terms and conditions data for testing
func (m *TermsAndConditionsMock) GetMockTermsAndConditionsData() map[string]any {
	jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_terms_and_conditions_maximal.json"))
	if err != nil {
		// Fallback to hardcoded response if file loading fails
		return map[string]any{
			"id":                   "test-terms-and-conditions-id",
			"displayName":          "Test Terms and Conditions",
			"description":          "Test terms and conditions for unit testing",
			"title":                "Company Terms and Conditions",
			"bodyText":             "These are the terms and conditions that users must accept.",
			"acceptanceStatement":  "I accept the terms and conditions",
			"version":              1,
			"roleScopeTagIds":      []string{"0"},
			"createdDateTime":      "2024-01-01T00:00:00Z",
			"lastModifiedDateTime": "2024-01-01T00:00:00Z",
		}
	}

	var response map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		// Fallback to hardcoded response if parsing fails
		return map[string]any{
			"id":                   "test-terms-and-conditions-id",
			"displayName":          "Test Terms and Conditions",
			"description":          "Test terms and conditions for unit testing",
			"title":                "Company Terms and Conditions",
			"bodyText":             "These are the terms and conditions that users must accept.",
			"acceptanceStatement":  "I accept the terms and conditions",
			"version":              1,
			"roleScopeTagIds":      []string{"0"},
			"createdDateTime":      "2024-01-01T00:00:00Z",
			"lastModifiedDateTime": "2024-01-01T00:00:00Z",
		}
	}
	return response
}

// getTermsAndConditionsAssignmentsResponder handles GET requests to retrieve assignments
func (m *TermsAndConditionsMock) getTermsAndConditionsAssignmentsResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// For unit tests, return empty assignments collection
		response := map[string]any{
			"value": []any{},
		}
		return factories.SuccessResponse(200, response)(req)
	}
}

// createTermsAndConditionsAssignmentResponder handles POST requests to create assignments
func (m *TermsAndConditionsMock) createTermsAndConditionsAssignmentResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		// Create a mock assignment response
		response := map[string]any{
			"id":     uuid.New().String(),
			"target": requestBody["target"],
		}

		return factories.SuccessResponse(201, response)(req)
	}
}

// GetMockTermsAndConditionsMinimalData returns minimal terms and conditions data for testing
func (m *TermsAndConditionsMock) GetMockTermsAndConditionsMinimalData() map[string]any {
	jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_terms_and_conditions_minimal.json"))
	if err != nil {
		// Fallback to hardcoded response if file loading fails
		return map[string]any{
			"id":                   "test-minimal-terms-and-conditions-id",
			"displayName":          "Test Minimal Terms and Conditions",
			"title":                "Simple Terms",
			"bodyText":             "Basic terms and conditions.",
			"acceptanceStatement":  "I agree",
			"version":              1,
			"roleScopeTagIds":      []string{"0"},
			"createdDateTime":      "2024-01-01T00:00:00Z",
			"lastModifiedDateTime": "2024-01-01T00:00:00Z",
		}
	}

	var response map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		// Fallback to hardcoded response if parsing fails
		return map[string]any{
			"id":                   "test-minimal-terms-and-conditions-id",
			"displayName":          "Test Minimal Terms and Conditions",
			"title":                "Simple Terms",
			"bodyText":             "Basic terms and conditions.",
			"acceptanceStatement":  "I agree",
			"version":              1,
			"roleScopeTagIds":      []string{"0"},
			"createdDateTime":      "2024-01-01T00:00:00Z",
			"lastModifiedDateTime": "2024-01-01T00:00:00Z",
		}
	}
	return response
}
