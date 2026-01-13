package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of manager relationships for consistent responses
var mockState struct {
	sync.Mutex
	managers map[string]string // map[userId]managerId
}

func init() {
	mockState.managers = make(map[string]string)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("user_manager", &UserManagerMock{})
}

// UserManagerMock provides mock responses for user manager operations
type UserManagerMock struct{}

// Ensure UserManagerMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*UserManagerMock)(nil)

// RegisterMocks registers HTTP mock responses for user manager operations
func (m *UserManagerMock) RegisterMocks() {
	mockState.Lock()
	mockState.managers = make(map[string]string)
	mockState.Unlock()

	// Register GET for user manager
	// GET /users/{usersId}/manager
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/users/([a-fA-F0-9\-]+)/manager$`,
		func(req *http.Request) (*http.Response, error) {
			userId := httpmock.MustGetSubmatch(req, 1)

			mockState.Lock()
			managerId, exists := mockState.managers[userId]
			mockState.Unlock()

			if !exists {
				errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")
				return httpmock.NewStringResponse(404, errorResp), nil
			}

			// Load from externalized JSON and update with dynamic managerId
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_read/get_manager_success.json")
			var response map[string]any
			if err := json.Unmarshal([]byte(jsonStr), &response); err == nil {
				response["id"] = managerId
				respBody, _ := json.Marshal(response)
				resp := httpmock.NewStringResponse(200, string(respBody))
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			}

			// Direct return if parsing fails
			resp := httpmock.NewStringResponse(200, jsonStr)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	// Register PUT for adding manager reference
	// PUT /users/{usersId}/manager/$ref
	httpmock.RegisterResponder("PUT", `=~^https://graph\.microsoft\.com/beta/users/([a-fA-F0-9\-]+)/manager/\$ref$`,
		func(req *http.Request) (*http.Response, error) {
			userId := httpmock.MustGetSubmatch(req, 1)

			// Extract manager ID from the request body @odata.id
			var body map[string]string
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
				return httpmock.NewStringResponse(400, errorResp), nil
			}

			odataId := body["@odata.id"]
			if odataId == "" {
				errorResp, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
				return httpmock.NewStringResponse(400, errorResp), nil
			}

			// Extract manager ID from the @odata.id URL
			odataParts := strings.Split(odataId, "/")
			managerId := odataParts[len(odataParts)-1]

			mockState.Lock()
			mockState.managers[userId] = managerId
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register DELETE for removing manager reference
	// DELETE /users/{usersId}/manager/$ref
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph\.microsoft\.com/beta/users/([a-fA-F0-9\-]+)/manager/\$ref$`,
		func(req *http.Request) (*http.Response, error) {
			userId := httpmock.MustGetSubmatch(req, 1)

			mockState.Lock()
			_, exists := mockState.managers[userId]
			if exists {
				delete(mockState.managers, userId)
			}
			mockState.Unlock()

			// Always return 204 - if not exists, it's already deleted
			return httpmock.NewStringResponse(204, ""), nil
		})
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *UserManagerMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.managers = make(map[string]string)
	mockState.Unlock()

	errorBadRequest, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
	errorNotFound, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")

	// Register error response for GET manager
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/users/([a-fA-F0-9\-]+)/manager$`,
		httpmock.NewStringResponder(404, errorNotFound))

	// Register error response for PUT manager reference
	httpmock.RegisterResponder("PUT", `=~^https://graph\.microsoft\.com/beta/users/([a-fA-F0-9\-]+)/manager/\$ref$`,
		httpmock.NewStringResponder(400, errorBadRequest))

	// Register error response for DELETE manager reference
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph\.microsoft\.com/beta/users/([a-fA-F0-9\-]+)/manager/\$ref$`,
		httpmock.NewStringResponder(400, errorBadRequest))
}

// CleanupMockState clears the mock state
func (m *UserManagerMock) CleanupMockState() {
	mockState.Lock()
	mockState.managers = make(map[string]string)
	mockState.Unlock()
}
