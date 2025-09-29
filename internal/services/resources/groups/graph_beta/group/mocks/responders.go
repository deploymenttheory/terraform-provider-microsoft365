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
	groups map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.groups = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("group", &GroupMock{})
}

// GroupMock provides mock responses for group operations
type GroupMock struct{}

// Ensure GroupMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*GroupMock)(nil)

// RegisterMocks registers HTTP mock responses for group operations
func (m *GroupMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.groups = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for listing groups
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/groups",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			groups := make([]map[string]any, 0, len(mockState.groups))
			for _, group := range mockState.groups {
				groups = append(groups, group)
			}
			mockState.Unlock()

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#groups",
				"value":          groups,
			}

			respBody, _ := json.Marshal(response)
			return httpmock.NewStringResponse(200, string(respBody)), nil
		})

	// Register GET for specific group with special test IDs
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/(.+)$`,
		func(req *http.Request) (*http.Response, error) {
			groupID := httpmock.MustGetSubmatch(req, 1)

			// Handle special test IDs with external JSON files
			switch {
			case strings.Contains(groupID, "minimal"):
				jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_read/get_group_minimal.json")
				if err != nil {
					return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
				}
				return httpmock.NewStringResponse(200, jsonStr), nil
			case strings.Contains(groupID, "maximal"):
				jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_read/get_group_maximal.json")
				if err != nil {
					return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
				}
				return httpmock.NewStringResponse(200, jsonStr), nil
			case strings.Contains(groupID, "error"):
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_resource_not_found.json")
				return httpmock.NewStringResponse(404, jsonStr), nil
			}

			mockState.Lock()
			defer mockState.Unlock()

			if group, exists := mockState.groups[groupID]; exists {
				respBody, _ := json.Marshal(group)
				return httpmock.NewStringResponse(200, string(respBody)), nil
			}

			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		})

	// Register POST for creating groups
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/groups",
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Generate a new ID for the group
			newID := uuid.New().String()
			requestBody["id"] = newID

			// Add additional fields that would be set by the API
			requestBody["@odata.context"] = "https://graph.microsoft.com/beta/$metadata#groups/$entity"
			requestBody["createdDateTime"] = "2023-01-01T00:00:00Z"

			// Handle special test cases based on display name
			if displayName, ok := requestBody["displayName"].(string); ok {
				switch {
				case strings.Contains(displayName, "minimal"):
					jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_create/post_group_minimal.json")
					if err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
					}
					return httpmock.NewStringResponse(201, jsonStr), nil
				case strings.Contains(displayName, "maximal"):
					jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_create/post_group_maximal.json")
					if err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
					}
					return httpmock.NewStringResponse(201, jsonStr), nil
				}
			}

			mockState.Lock()
			mockState.groups[newID] = requestBody
			mockState.Unlock()

			respBody, _ := json.Marshal(requestBody)
			return httpmock.NewStringResponse(201, string(respBody)), nil
		})

	// Register PATCH for updating groups
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/groups/(.+)$`,
		func(req *http.Request) (*http.Response, error) {
			groupID := httpmock.MustGetSubmatch(req, 1)

			var requestBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Handle special test cases
			if displayName, ok := requestBody["displayName"].(string); ok {
				switch {
				case strings.Contains(displayName, "minimal_to_maximal"):
					return httpmock.NewStringResponse(204, ""), nil
				case strings.Contains(displayName, "maximal_to_minimal"):
					return httpmock.NewStringResponse(204, ""), nil
				case strings.Contains(displayName, "minimal"):
					return httpmock.NewStringResponse(204, ""), nil
				case strings.Contains(displayName, "maximal"):
					return httpmock.NewStringResponse(204, ""), nil
				}
			}

			mockState.Lock()
			defer mockState.Unlock()

			if group, exists := mockState.groups[groupID]; exists {
				// Update the existing group
				for k, v := range requestBody {
					group[k] = v
				}
				return httpmock.NewStringResponse(204, ""), nil
			}

			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		})

	// Register DELETE for deleting groups
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/groups/(.+)$`,
		func(req *http.Request) (*http.Response, error) {
			groupID := httpmock.MustGetSubmatch(req, 1)

			mockState.Lock()
			defer mockState.Unlock()

			if _, exists := mockState.groups[groupID]; exists {
				delete(mockState.groups, groupID)
				return httpmock.NewStringResponse(204, ""), nil
			}

			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		})
}

// RegisterErrorMocks registers HTTP mock responses that return errors for testing error handling
func (m *GroupMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/groups",
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_invalid_display_name.json")
			return httpmock.NewStringResponse(400, jsonStr), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/(.+)$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_resource_not_found.json")
			return httpmock.NewStringResponse(404, jsonStr), nil
		})
}

// CleanupMockState cleans up the mock state (called after each test)
func (m *GroupMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.groups = make(map[string]map[string]any)
}
