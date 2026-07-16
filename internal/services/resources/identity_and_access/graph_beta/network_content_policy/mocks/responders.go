package mocks

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/jarcoal/httpmock"
)

const contentPolicyID = "00000000-0000-0000-0000-000000000201"

var mockState struct {
	sync.Mutex
	contentPolicies map[string]map[string]any
}

func init() {
	mockState.contentPolicies = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("network_content_policy", &ContentPolicyMock{})
}

type ContentPolicyMock struct{}

var _ mocks.MockRegistrar = (*ContentPolicyMock)(nil)

func (m *ContentPolicyMock) RegisterMocks() {
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/networkaccess/filePolicies", m.createContentPolicyResponder())
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/networkaccess/filePolicies/([^/]+)$`, m.getContentPolicyResponder())
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/networkaccess/filePolicies/([^/]+)$`, m.updateContentPolicyResponder())
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/networkaccess/filePolicies/([^/]+)$`, m.deleteContentPolicyResponder())
}

func (m *ContentPolicyMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/networkaccess/filePolicies", func(req *http.Request) (*http.Response, error) {
		return jsonFixtureResponse(req, 400, filepath.Join("..", "tests", "responses", "validate_create", "post_content_policy_error.json"))
	})
}

func (m *ContentPolicyMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.contentPolicies = make(map[string]map[string]any)
}

func (m *ContentPolicyMock) createContentPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}
		if _, exists := requestBody["policyRules"]; exists {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"policyRules must not be managed by this resource"}}`), nil
		}

		response, err := loadJSONFixture(filepath.Join("..", "tests", "responses", "validate_create", "post_content_policy.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}
		mergeContentPolicyRequest(response, requestBody)

		mockState.Lock()
		mockState.contentPolicies[contentPolicyID] = response
		mockState.Unlock()

		return factories.SuccessResponse(201, response)(req)
	}
}

func (m *ContentPolicyMock) getContentPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := factories.ExtractIDFromURL(req.URL.Path, "/networkaccess/filePolicies/")

		mockState.Lock()
		policy, exists := mockState.contentPolicies[id]
		mockState.Unlock()
		if !exists {
			return jsonFixtureResponse(req, 404, filepath.Join("..", "tests", "responses", "validate_delete", "get_content_policy_not_found.json"))
		}

		return factories.SuccessResponse(200, cloneMap(policy))(req)
	}
}

func (m *ContentPolicyMock) updateContentPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := factories.ExtractIDFromURL(req.URL.Path, "/networkaccess/filePolicies/")

		mockState.Lock()
		policy, exists := mockState.contentPolicies[id]
		mockState.Unlock()
		if !exists {
			return jsonFixtureResponse(req, 404, filepath.Join("..", "tests", "responses", "validate_delete", "get_content_policy_not_found.json"))
		}

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}
		if _, exists := requestBody["policyRules"]; exists {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"policyRules must not be managed by this resource"}}`), nil
		}

		updated := cloneMap(policy)
		mergeContentPolicyRequest(updated, requestBody)
		updated["lastModifiedDateTime"] = "2026-07-14T13:15:48.7486334Z"

		mockState.Lock()
		mockState.contentPolicies[id] = updated
		mockState.Unlock()

		return factories.SuccessResponse(200, updated)(req)
	}
}

func (m *ContentPolicyMock) deleteContentPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := factories.ExtractIDFromURL(req.URL.Path, "/networkaccess/filePolicies/")

		mockState.Lock()
		_, exists := mockState.contentPolicies[id]
		if exists {
			delete(mockState.contentPolicies, id)
		}
		mockState.Unlock()
		if !exists {
			return jsonFixtureResponse(req, 404, filepath.Join("..", "tests", "responses", "validate_delete", "get_content_policy_not_found.json"))
		}

		return factories.EmptySuccessResponse(204)(req)
	}
}

func mergeContentPolicyRequest(response map[string]any, requestBody map[string]any) {
	for _, key := range []string{"name", "description", "settings"} {
		if value, ok := requestBody[key]; ok {
			response[key] = value
		}
	}
}

func cloneMap(input map[string]any) map[string]any {
	output := make(map[string]any, len(input))
	for key, value := range input {
		output[key] = value
	}
	return output
}

func loadJSONFixture(path string) (map[string]any, error) {
	jsonContent, err := helpers.ParseJSONFile(path)
	if err != nil {
		return nil, err
	}
	var response map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		return nil, err
	}
	return response, nil
}

func jsonFixtureResponse(req *http.Request, status int, path string) (*http.Response, error) {
	response, err := loadJSONFixture(path)
	if err != nil {
		return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
	}
	return factories.SuccessResponse(status, response)(req)
}
