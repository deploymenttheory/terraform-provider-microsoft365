package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/jarcoal/httpmock"

	providermocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const (
	AssignmentID  = "opaque-directory-role-assignment-1"
	collectionURL = "https://graph.microsoft.com/beta/roleManagement/directory/roleAssignments"
)

var state struct {
	sync.Mutex
	active map[string]any
}

type DirectoryRoleAssignmentMock struct{}

func init() {
	providermocks.GlobalRegistry.Register("directory_role_assignment", &DirectoryRoleAssignmentMock{})
}

func reply(status int, body any) (*http.Response, error) {
	response, err := httpmock.NewJsonResponse(status, body)
	if err != nil {
		return nil, fmt.Errorf("encode directory role assignment mock: %w", err)
	}
	return response, nil
}

func graphError(status int) (*http.Response, error) {
	return reply(status, map[string]any{"error": map[string]any{"code": "Request_BadRequest", "message": "directory role assignment probe error"}})
}

func (m *DirectoryRoleAssignmentMock) RegisterMocks() {
	httpmock.RegisterResponder("POST", collectionURL, func(req *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return nil, fmt.Errorf("decode directory role assignment request: %w", err)
		}
		state.Lock()
		defer state.Unlock()
		if state.active != nil {
			return graphError(400)
		}
		if body["directoryScopeId"] != "/" || body["appScopeId"] != nil {
			return graphError(400)
		}
		body["id"] = AssignmentID
		state.active = body
		return reply(201, body)
	})
	httpmock.RegisterResponder("GET", collectionURL+"/"+AssignmentID, func(_ *http.Request) (*http.Response, error) {
		state.Lock()
		defer state.Unlock()
		if state.active == nil {
			return graphError(404)
		}
		return reply(200, state.active)
	})
	httpmock.RegisterResponder("DELETE", collectionURL+"/"+AssignmentID, func(_ *http.Request) (*http.Response, error) {
		state.Lock()
		defer state.Unlock()
		state.active = nil
		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *DirectoryRoleAssignmentMock) RegisterErrorMocks() {
	for _, method := range []string{"POST", "GET", "DELETE"} {
		url := collectionURL
		if method != "POST" {
			url += "/" + AssignmentID
		}
		httpmock.RegisterResponder(method, url, func(_ *http.Request) (*http.Response, error) { return graphError(403) })
	}
}

func (m *DirectoryRoleAssignmentMock) CleanupMockState() {
	state.Lock()
	defer state.Unlock()
	state.active = nil
}

func (m *DirectoryRoleAssignmentMock) Seed(scope string) {
	state.Lock()
	defer state.Unlock()
	state.active = map[string]any{
		"id":               AssignmentID,
		"principalId":      "11111111-1111-4111-8111-111111111111",
		"roleDefinitionId": "88d8e3e3-8f55-4a1e-953a-9b9898b8876b",
		"directoryScopeId": scope,
	}
}
