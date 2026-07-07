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
	policies map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.policies = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("operation_approval_policy", &OperationApprovalPolicyMock{})
}

// OperationApprovalPolicyMock provides mock responses for operation approval policy operations
type OperationApprovalPolicyMock struct{}

// Ensure OperationApprovalPolicyMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*OperationApprovalPolicyMock)(nil)

// RegisterMocks sets up all the mock HTTP responders for operation approval policy operations
// This implements the MockRegistrar interface
func (m *OperationApprovalPolicyMock) RegisterMocks() {
	// POST /deviceManagement/operationApprovalPolicies - Create
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/operationApprovalPolicies",
		m.createPolicyResponder())

	// GET /deviceManagement/operationApprovalPolicies/{id} - Read
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/operationApprovalPolicies/([^/]+)$`,
		m.getPolicyResponder())

	// PATCH /deviceManagement/operationApprovalPolicies/{id} - Update
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/operationApprovalPolicies/([^/]+)$`,
		m.updatePolicyResponder())

	// DELETE /deviceManagement/operationApprovalPolicies/{id} - Delete
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/operationApprovalPolicies/([^/]+)$`,
		m.deletePolicyResponder())
}

// createPolicyResponder handles POST requests to create operation approval policies
func (m *OperationApprovalPolicyMock) createPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		// Load base response from JSON file
		jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_operation_approval_policy.json"))
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
		if policyType, ok := requestBody["policyType"]; ok {
			response["policyType"] = policyType
		}
		if policyPlatform, ok := requestBody["policyPlatform"]; ok {
			response["policyPlatform"] = policyPlatform
		}
		if policySet, ok := requestBody["policySet"]; ok {
			response["policySet"] = policySet
		}
		if approverGroupIds, ok := requestBody["approverGroupIds"]; ok {
			response["approverGroupIds"] = approverGroupIds
		}

		// Store in mock state
		mockState.Lock()
		mockState.policies[id] = response
		mockState.Unlock()

		return factories.SuccessResponse(201, response)(req)
	}
}

// getPolicyResponder handles GET requests to retrieve operation approval policies
func (m *OperationApprovalPolicyMock) getPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/operationApprovalPolicies/")

		mockState.Lock()
		policy, exists := mockState.policies[id]
		mockState.Unlock()

		if !exists {
			// Check for special test IDs
			switch {
			case strings.Contains(id, "minimal"):
				jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_operation_approval_policy_minimal.json"))
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
				jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "get_operation_approval_policy_maximal.json"))
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
				jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_operation_approval_policy_not_found.json"))
				if err == nil {
					var errorResponse map[string]any
					if json.Unmarshal([]byte(jsonContent), &errorResponse) == nil {
						return httpmock.NewJsonResponse(404, errorResponse)
					}
				}
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
			}
		}

		return factories.SuccessResponse(200, policy)(req)
	}
}

// updatePolicyResponder handles PATCH requests to update operation approval policies
func (m *OperationApprovalPolicyMock) updatePolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/operationApprovalPolicies/")

		mockState.Lock()
		policy, exists := mockState.policies[id]
		mockState.Unlock()

		if !exists {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_operation_approval_policy_not_found.json"))
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
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		// Load update template
		jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_update", "get_operation_approval_policy_updated.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}

		var updatedPolicy map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &updatedPolicy); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
		}

		// Start with existing data
		for k, v := range policy {
			updatedPolicy[k] = v
		}

		// Apply updates from request body
		for k, v := range requestBody {
			updatedPolicy[k] = v
		}

		// Update lastModifiedDateTime
		updatedPolicy["lastModifiedDateTime"] = "2024-01-01T12:00:00Z"

		// Store updated version
		mockState.Lock()
		mockState.policies[id] = updatedPolicy
		mockState.Unlock()

		return factories.SuccessResponse(200, updatedPolicy)(req)
	}
}

// deletePolicyResponder handles DELETE requests to remove operation approval policies
func (m *OperationApprovalPolicyMock) deletePolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		// Extract ID from URL
		id := factories.ExtractIDFromURL(req.URL.Path, "/deviceManagement/operationApprovalPolicies/")

		mockState.Lock()
		_, exists := mockState.policies[id]
		if exists {
			delete(mockState.policies, id)
		}
		mockState.Unlock()

		if !exists {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_operation_approval_policy_not_found.json"))
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
func (m *OperationApprovalPolicyMock) RegisterErrorMocks() {
	// POST - Create error
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/operationApprovalPolicies",
		func(req *http.Request) (*http.Response, error) {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_operation_approval_policy_error.json"))
			if err == nil {
				var errorResponse map[string]any
				if json.Unmarshal([]byte(jsonContent), &errorResponse) == nil {
					return httpmock.NewJsonResponse(400, errorResponse)
				}
			}
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`), nil
		})

	// GET - Read error (simulates not found or access denied)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/operationApprovalPolicies/error-id$`,
		factories.ErrorResponse(403, "Forbidden", "Access denied"))

	// PATCH - Update error
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/operationApprovalPolicies/error-id$`,
		factories.ErrorResponse(500, "InternalServerError", "Internal server error"))

	// DELETE - Delete error
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/operationApprovalPolicies/error-id$`,
		factories.ErrorResponse(409, "Conflict", "Operation approval policy is in use"))
}

// CleanupMockState clears all stored mock state
func (m *OperationApprovalPolicyMock) CleanupMockState() {
	mockState.Lock()
	mockState.policies = make(map[string]map[string]any)
	mockState.Unlock()
}
