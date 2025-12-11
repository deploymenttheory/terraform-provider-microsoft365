package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	assignments map[string]bool // key: agentInstanceId/agentCollectionId
}

func init() {
	mockState.assignments = make(map[string]bool)
	mocks.GlobalRegistry.Register("agent_collection_assignment", &AgentCollectionAssignmentMock{})
}

type AgentCollectionAssignmentMock struct{}

var _ mocks.MockRegistrar = (*AgentCollectionAssignmentMock)(nil)

func (m *AgentCollectionAssignmentMock) RegisterMocks() {
	mockState.Lock()
	mockState.assignments = make(map[string]bool)
	mockState.Unlock()

	// Validate agent instance exists - GET /agentRegistry/agentInstances
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentInstances`,
		func(req *http.Request) (*http.Response, error) {
			jsonContent, err := helpers.ParseJSONFile("../tests/responses/validate_list/list_agent_instances.json")
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load JSON response file"}}`), nil
			}

			var responseObj map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &responseObj); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}

			return httpmock.NewJsonResponse(200, responseObj)
		})

	// Create assignment - POST /agentRegistry/agentCollections/{id}/members/$ref
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentCollections/[0-9a-fA-F-]+/members/\$ref$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			// Path: /beta/agentRegistry/agentCollections/{id}/members/$ref
			agentCollectionId := parts[4]

			// Store assignment using fixed agent instance ID from test
			key := "11111111-1111-1111-1111-111111111111/" + agentCollectionId

			mockState.Lock()
			mockState.assignments[key] = true
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	// List members - GET /agentRegistry/agentCollections/{id}/members
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentCollections/[0-9a-fA-F-]+/members$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			agentCollectionId := parts[4]

			key := "11111111-1111-1111-1111-111111111111/" + agentCollectionId

			mockState.Lock()
			exists := mockState.assignments[key]
			mockState.Unlock()

			var jsonFile string
			if exists {
				jsonFile = "../tests/responses/validate_read/list_members_with_assignment.json"
			} else {
				jsonFile = "../tests/responses/validate_read/list_members_empty.json"
			}

			jsonContent, err := helpers.ParseJSONFile(jsonFile)
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load JSON response file"}}`), nil
			}

			var responseObj map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &responseObj); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse JSON response"}}`), nil
			}

			return httpmock.NewJsonResponse(200, responseObj)
		})

	// Delete assignment - DELETE /agentRegistry/agentCollections/{id}/members/{instanceId}
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentCollections/[0-9a-fA-F-]+/members/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			agentCollectionId := parts[4]
			agentInstanceId := parts[6]

			key := agentInstanceId + "/" + agentCollectionId

			mockState.Lock()
			delete(mockState.assignments, key)
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})
}

func (m *AgentCollectionAssignmentMock) RegisterErrorMocks() {
	errorBadRequest, _ := helpers.ParseJSONFile("../tests/responses/error_bad_request.json")
	errorNotFound, _ := helpers.ParseJSONFile("../tests/responses/error_not_found.json")

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentInstances`,
		httpmock.NewStringResponder(400, errorBadRequest))
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentCollections/[0-9a-fA-F-]+/members/\$ref$`,
		httpmock.NewStringResponder(400, errorBadRequest))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentCollections/[0-9a-fA-F-]+/members$`,
		httpmock.NewStringResponder(404, errorNotFound))
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/agentRegistry/agentCollections/[0-9a-fA-F-]+/members/[0-9a-fA-F-]+$`,
		httpmock.NewStringResponder(400, errorBadRequest))
}

func (m *AgentCollectionAssignmentMock) CleanupMockState() {
	mockState.Lock()
	mockState.assignments = make(map[string]bool)
	mockState.Unlock()
}
