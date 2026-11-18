package mocks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	providermocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const (
	GrantID       = "opaque-grant-ID_001"
	ClientID      = "11111111-1111-4111-8111-111111111111"
	ResourceID    = "22222222-2222-4222-8222-222222222222"
	PrincipalID   = "33333333-3333-4333-8333-333333333333"
	collectionURL = "https://graph.microsoft.com/beta/oauth2PermissionGrants"
)

var state struct {
	sync.Mutex
	active  map[string]any
	patches []map[string]any
}

type OAuth2PermissionGrantMock struct{}

func init() {
	providermocks.GlobalRegistry.Register("oauth2_permission_grant", &OAuth2PermissionGrantMock{})
}

func readJSON(req *http.Request) (map[string]any, error) {
	content, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("read grant payload: %w", err)
	}
	var body map[string]any
	if err := json.Unmarshal(content, &body); err != nil {
		return nil, fmt.Errorf("decode grant payload: %w", err)
	}
	return body, nil
}

func graphError(status int, message string) (*http.Response, error) {
	response, err := httpmock.NewJsonResponse(status, map[string]any{"error": map[string]any{"code": "Request_BadRequest", "message": message}})
	if err != nil {
		return nil, fmt.Errorf("encode mock grant error: %w", err)
	}
	return response, nil
}

func (m *OAuth2PermissionGrantMock) RegisterMocks() {
	httpmock.RegisterResponder("POST", collectionURL, func(req *http.Request) (*http.Response, error) {
		body, err := readJSON(req)
		if err != nil {
			return nil, err
		}
		state.Lock()
		defer state.Unlock()
		if state.active != nil {
			return graphError(400, "Permission entry already exists.")
		}
		if body["startTime"] == nil || body["expiryTime"] == nil {
			return graphError(400, "Missing property: expiryTime")
		}
		body["id"] = GrantID
		if _, exists := body["principalId"]; !exists {
			body["principalId"] = nil
		}
		state.active = body
		return httpmock.NewJsonResponse(201, body)
	})
	httpmock.RegisterResponder("GET", collectionURL+"/"+GrantID, func(_ *http.Request) (*http.Response, error) {
		state.Lock()
		defer state.Unlock()
		if state.active == nil {
			return graphError(404, "Resource not found")
		}
		return httpmock.NewJsonResponse(200, state.active)
	})
	httpmock.RegisterResponder("PATCH", collectionURL+"/"+GrantID, func(req *http.Request) (*http.Response, error) {
		body, err := readJSON(req)
		if err != nil {
			return nil, err
		}
		state.Lock()
		defer state.Unlock()
		if state.active == nil {
			return graphError(404, "Resource not found")
		}
		for _, key := range []string{"clientId", "resourceId", "principalId", "consentType", "startTime", "expiryTime", "id"} {
			if _, sent := body[key]; sent {
				return graphError(400, "Property cannot be updated: "+key)
			}
		}
		state.patches = append(state.patches, body)
		state.active["scope"] = body["scope"]
		return httpmock.NewStringResponse(204, ""), nil
	})
	httpmock.RegisterResponder("DELETE", collectionURL+"/"+GrantID, func(_ *http.Request) (*http.Response, error) {
		state.Lock()
		defer state.Unlock()
		if state.active == nil {
			return graphError(404, "Resource not found")
		}
		state.active = nil
		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *OAuth2PermissionGrantMock) RegisterErrorMocks() {
	for _, method := range []string{"POST", "GET", "PATCH", "DELETE"} {
		url := collectionURL
		if method != "POST" {
			url += "/" + GrantID
		}
		httpmock.RegisterResponder(method, url, func(_ *http.Request) (*http.Response, error) { return graphError(403, "Forbidden") })
	}
}

func (m *OAuth2PermissionGrantMock) SeedFixture() error {
	fixture, err := helpers.ParseJSONFile("../tests/responses/get_grant.json")
	if err != nil {
		return fmt.Errorf("load grant fixture: %w", err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(fixture), &body); err != nil {
		return fmt.Errorf("decode grant fixture: %w", err)
	}
	state.Lock()
	defer state.Unlock()
	state.active = body
	return nil
}

func (m *OAuth2PermissionGrantMock) Patches() []map[string]any {
	state.Lock()
	defer state.Unlock()
	return append([]map[string]any(nil), state.patches...)
}

func (m *OAuth2PermissionGrantMock) CleanupMockState() {
	state.Lock()
	defer state.Unlock()
	state.active = nil
	state.patches = nil
}
