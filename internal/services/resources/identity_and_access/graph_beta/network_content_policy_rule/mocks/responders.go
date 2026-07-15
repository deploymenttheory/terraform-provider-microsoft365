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
	"github.com/jarcoal/httpmock"
)

const contentPolicyRuleID = "00000000-0000-0000-0000-000000000302"

var mockState struct {
	sync.Mutex
	rules map[string]map[string]any
}

func init() {
	mockState.rules = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("network_content_policy_rule", &ContentPolicyRuleMock{})
}

type ContentPolicyRuleMock struct{}

var _ mocks.MockRegistrar = (*ContentPolicyRuleMock)(nil)

func (m *ContentPolicyRuleMock) RegisterMocks() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/networkaccess/filePolicies/([^/]+)/policyRules$`, m.createResponder())
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/networkaccess/filePolicies/([^/]+)/policyRules/([^/]+)$`, m.getResponder())
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/networkaccess/filePolicies/([^/]+)/policyRules/([^/]+)$`, m.updateResponder())
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/networkaccess/filePolicies/([^/]+)/policyRules/([^/]+)$`, m.deleteResponder())
}

func (m *ContentPolicyRuleMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/networkaccess/filePolicies/([^/]+)/policyRules$`, func(req *http.Request) (*http.Response, error) {
		return fixtureResponse(req, 400, filepath.Join("..", "tests", "responses", "validate_create", "post_content_policy_rule_error.json"))
	})
}

func (m *ContentPolicyRuleMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.rules = make(map[string]map[string]any)
}

func (m *ContentPolicyRuleMock) createResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}
		if body["@odata.type"] != "#microsoft.graph.networkaccess.fileRule" {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid rule type"}}`), nil
		}
		response, err := loadFixture(filepath.Join("..", "tests", "responses", "validate_create", "post_content_policy_rule.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load fixture"}}`), nil
		}
		mergeRule(response, body)
		parentID, _ := requestIDs(req.URL.Path)
		mockState.Lock()
		mockState.rules[parentID+"/"+contentPolicyRuleID] = response
		mockState.Unlock()
		return factories.SuccessResponse(201, response)(req)
	}
}

func (m *ContentPolicyRuleMock) getResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		parentID, ruleID := requestIDs(req.URL.Path)
		mockState.Lock()
		rule, exists := mockState.rules[parentID+"/"+ruleID]
		mockState.Unlock()
		if !exists {
			return fixtureResponse(req, 404, filepath.Join("..", "tests", "responses", "validate_delete", "get_content_policy_rule_not_found.json"))
		}
		return factories.SuccessResponse(200, cloneRule(rule))(req)
	}
}

func (m *ContentPolicyRuleMock) updateResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		parentID, ruleID := requestIDs(req.URL.Path)
		key := parentID + "/" + ruleID
		mockState.Lock()
		rule, exists := mockState.rules[key]
		mockState.Unlock()
		if !exists {
			return fixtureResponse(req, 404, filepath.Join("..", "tests", "responses", "validate_delete", "get_content_policy_rule_not_found.json"))
		}
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}
		updated := cloneRule(rule)
		mergeRule(updated, body)
		mockState.Lock()
		mockState.rules[key] = updated
		mockState.Unlock()
		return factories.EmptySuccessResponse(204)(req)
	}
}

func (m *ContentPolicyRuleMock) deleteResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		parentID, ruleID := requestIDs(req.URL.Path)
		key := parentID + "/" + ruleID
		mockState.Lock()
		_, exists := mockState.rules[key]
		delete(mockState.rules, key)
		mockState.Unlock()
		if !exists {
			return fixtureResponse(req, 404, filepath.Join("..", "tests", "responses", "validate_delete", "get_content_policy_rule_not_found.json"))
		}
		return factories.EmptySuccessResponse(204)(req)
	}
}

func requestIDs(path string) (string, string) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for index, part := range parts {
		if part != "filePolicies" || index+1 >= len(parts) {
			continue
		}
		parentID := parts[index+1]
		if index+3 < len(parts) && parts[index+2] == "policyRules" {
			return parentID, parts[index+3]
		}
		return parentID, ""
	}
	return "", ""
}

func mergeRule(response, body map[string]any) {
	for _, key := range []string{"@odata.type", "name", "description", "action", "priority", "settings", "matchingConditions"} {
		if value, ok := body[key]; ok {
			response[key] = value
		}
	}
}

func cloneRule(input map[string]any) map[string]any {
	result := make(map[string]any, len(input))
	for key, value := range input {
		result[key] = value
	}
	return result
}

func loadFixture(path string) (map[string]any, error) {
	content, err := helpers.ParseJSONFile(path)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func fixtureResponse(req *http.Request, status int, path string) (*http.Response, error) {
	response, err := loadFixture(path)
	if err != nil {
		return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load fixture"}}`), nil
	}
	return factories.SuccessResponse(status, response)(req)
}
