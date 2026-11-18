package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/jarcoal/httpmock"

	providermocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const policyURL = "https://graph.microsoft.com/beta/policies/adminConsentRequestPolicy"

var state struct {
	sync.Mutex
	policy map[string]any
	puts   []map[string]any
}

type AdminConsentRequestPolicyMock struct{}

func init() {
	providermocks.GlobalRegistry.Register("admin_consent_request_policy", &AdminConsentRequestPolicyMock{})
}

func reply(status int, body any) (*http.Response, error) {
	response, err := httpmock.NewJsonResponse(status, body)
	if err != nil {
		return nil, fmt.Errorf("encode admin consent policy mock: %w", err)
	}
	return response, nil
}

func (m *AdminConsentRequestPolicyMock) RegisterMocks() {
	httpmock.RegisterResponder("GET", policyURL, func(_ *http.Request) (*http.Response, error) {
		state.Lock()
		defer state.Unlock()
		return reply(200, state.policy)
	})
	httpmock.RegisterResponder("PUT", policyURL, func(req *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return nil, fmt.Errorf("decode admin consent policy request: %w", err)
		}
		for _, key := range []string{"isEnabled", "notifyReviewers", "remindersEnabled", "requestDurationInDays", "reviewers"} {
			if body[key] == nil {
				return reply(400, map[string]any{"error": map[string]any{"code": "BadRequest", "message": "Missing writable property: " + key}})
			}
		}
		if len(body) != 5 {
			return reply(400, map[string]any{"error": map[string]any{"code": "BadRequest", "message": "Only the five writable properties are accepted"}})
		}
		state.Lock()
		defer state.Unlock()
		version := state.policy["version"].(int) + 1
		state.puts = append(state.puts, body)
		state.policy = make(map[string]any, len(body)+1)
		for key, value := range body {
			state.policy[key] = value
		}
		state.policy["version"] = version
		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *AdminConsentRequestPolicyMock) RegisterErrorMocks() {
	for _, method := range []string{"GET", "PUT"} {
		httpmock.RegisterResponder(method, policyURL, func(_ *http.Request) (*http.Response, error) {
			return reply(403, map[string]any{"error": map[string]any{"code": "Authorization_RequestDenied", "message": "Insufficient privileges"}})
		})
	}
}

func (m *AdminConsentRequestPolicyMock) CleanupMockState() {
	state.Lock()
	defer state.Unlock()
	state.policy = map[string]any{"isEnabled": false, "notifyReviewers": false, "remindersEnabled": false, "requestDurationInDays": 0, "reviewers": []any{}, "version": 0}
	state.puts = nil
}

func (m *AdminConsentRequestPolicyMock) Writes() []map[string]any {
	state.Lock()
	defer state.Unlock()
	return append([]map[string]any(nil), state.puts...)
}
