package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	providermocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const (
	ServicePrincipalID = "11111111-1111-4111-8111-111111111111"
	PolicyID           = "22222222-2222-4222-8222-222222222222"
	policyURL          = "https://graph.microsoft.com/beta/servicePrincipals/" + ServicePrincipalID + "/claimsPolicy"
)

var state struct {
	sync.Mutex
	policy  map[string]any
	puts    []map[string]any
	deletes int
}

type ServicePrincipalClaimsPolicyJSONMock struct{}

func init() {
	providermocks.GlobalRegistry.Register("service_principal_claims_policy_json", &ServicePrincipalClaimsPolicyJSONMock{})
}

func reply(status int, body any) (*http.Response, error) {
	response, err := httpmock.NewJsonResponse(status, body)
	if err != nil {
		return nil, fmt.Errorf("encode claims policy mock: %w", err)
	}
	return response, nil
}

func graphError(status int) (*http.Response, error) {
	return reply(status, map[string]any{"error": map[string]any{"code": "Request_BadRequest", "message": "claims policy probe error"}})
}

func (m *ServicePrincipalClaimsPolicyJSONMock) RegisterMocks() {
	httpmock.RegisterResponder("PUT", policyURL, func(req *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return nil, fmt.Errorf("decode claims policy request: %w", err)
		}
		if body["id"] != nil || body["@odata.context"] != nil {
			return graphError(400)
		}
		state.Lock()
		defer state.Unlock()
		state.puts = append(state.puts, body)
		remote := make(map[string]any, len(body)+4)
		for key, value := range body {
			remote[key] = value
		}
		remote["id"] = PolicyID
		remote["@odata.context"] = "https://graph.microsoft.com/beta/$metadata#claimsPolicy"
		if _, ok := remote["audienceOverride"]; !ok {
			remote["audienceOverride"] = nil
		}
		if _, ok := remote["groupFilter"]; !ok {
			remote["groupFilter"] = nil
		}
		state.policy = remote
		return httpmock.NewStringResponse(204, ""), nil
	})
	httpmock.RegisterResponder("GET", policyURL, func(_ *http.Request) (*http.Response, error) {
		state.Lock()
		defer state.Unlock()
		if state.policy == nil {
			return graphError(404)
		}
		return reply(200, state.policy)
	})
	httpmock.RegisterResponder("DELETE", policyURL, func(_ *http.Request) (*http.Response, error) {
		state.Lock()
		defer state.Unlock()
		state.deletes++
		return graphError(405)
	})
}

func (m *ServicePrincipalClaimsPolicyJSONMock) RegisterErrorMocks() {
	for _, method := range []string{"GET", "PUT"} {
		httpmock.RegisterResponder(method, policyURL, func(_ *http.Request) (*http.Response, error) { return graphError(403) })
	}
}

func (m *ServicePrincipalClaimsPolicyJSONMock) CleanupMockState() {
	state.Lock()
	defer state.Unlock()
	state.policy = nil
	state.puts = nil
	state.deletes = 0
}

func (m *ServicePrincipalClaimsPolicyJSONMock) SeedFixture() error {
	raw, err := helpers.ParseJSONFile("../tests/responses/get_policy.json")
	if err != nil {
		return fmt.Errorf("read claims policy fixture: %w", err)
	}
	var body map[string]any
	if err = json.Unmarshal([]byte(raw), &body); err != nil {
		return fmt.Errorf("decode claims policy fixture: %w", err)
	}
	state.Lock()
	defer state.Unlock()
	state.policy = body
	return nil
}

func (m *ServicePrincipalClaimsPolicyJSONMock) Requests() ([]map[string]any, int) {
	state.Lock()
	defer state.Unlock()
	return append([]map[string]any(nil), state.puts...), state.deletes
}
