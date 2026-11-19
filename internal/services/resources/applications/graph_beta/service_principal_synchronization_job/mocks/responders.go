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
	JobID              = "scim.opaque-job.with-dots"
	collectionURL      = "https://graph.microsoft.com/beta/servicePrincipals/" + ServicePrincipalID + "/synchronization/jobs"
)

var state struct {
	sync.Mutex
	active                  map[string]any
	starts, pauses, creates int
}

type ServicePrincipalSynchronizationJobMock struct{}

func init() {
	providermocks.GlobalRegistry.Register("service_principal_synchronization_job", &ServicePrincipalSynchronizationJobMock{})
}

func reply(status int, body any) (*http.Response, error) {
	response, err := httpmock.NewJsonResponse(status, body)
	if err != nil {
		return nil, fmt.Errorf("encode synchronization job mock: %w", err)
	}
	return response, nil
}

func graphError(status int) (*http.Response, error) {
	return reply(status, map[string]any{"error": map[string]any{"code": "Request_BadRequest", "message": "synchronization job probe error"}})
}

func fixture() (map[string]any, error) {
	raw, err := helpers.ParseJSONFile("../tests/responses/get_job.json")
	if err != nil {
		return nil, fmt.Errorf("read synchronization job fixture: %w", err)
	}
	var body map[string]any
	if err = json.Unmarshal([]byte(raw), &body); err != nil {
		return nil, fmt.Errorf("decode synchronization job fixture: %w", err)
	}
	return body, nil
}

func (m *ServicePrincipalSynchronizationJobMock) RegisterMocks() {
	httpmock.RegisterResponder("POST", collectionURL, func(req *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return nil, fmt.Errorf("decode synchronization job request: %w", err)
		}
		if len(body) != 1 || body["templateId"] == nil {
			return graphError(400)
		}
		value, err := fixture()
		if err != nil {
			return nil, err
		}
		value["templateId"] = body["templateId"]
		state.Lock()
		defer state.Unlock()
		if state.active != nil {
			return graphError(409)
		}
		state.active = value
		state.creates++
		return reply(201, value)
	})
	httpmock.RegisterResponder("GET", collectionURL+"/"+JobID, func(_ *http.Request) (*http.Response, error) {
		state.Lock()
		defer state.Unlock()
		if state.active == nil {
			return graphError(404)
		}
		return reply(200, state.active)
	})
	for _, action := range []string{"start", "pause"} {
		httpmock.RegisterResponder("POST", collectionURL+"/"+JobID+"/"+action, func(_ *http.Request) (*http.Response, error) {
			state.Lock()
			defer state.Unlock()
			if state.active == nil {
				return graphError(404)
			}
			schedule := state.active["schedule"].(map[string]any)
			if action == "start" {
				schedule["state"] = "Active"
				state.starts++
			} else {
				schedule["state"] = "Paused"
				state.pauses++
			}
			return httpmock.NewStringResponse(204, ""), nil
		})
	}
	httpmock.RegisterResponder("DELETE", collectionURL+"/"+JobID, func(_ *http.Request) (*http.Response, error) {
		state.Lock()
		defer state.Unlock()
		state.active = nil
		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *ServicePrincipalSynchronizationJobMock) RegisterErrorMocks() {
	for _, method := range []string{"GET", "DELETE", "POST"} {
		url := collectionURL
		if method != "POST" {
			url += "/" + JobID
		}
		httpmock.RegisterResponder(method, url, func(_ *http.Request) (*http.Response, error) { return graphError(403) })
	}
}

func (m *ServicePrincipalSynchronizationJobMock) CleanupMockState() {
	state.Lock()
	defer state.Unlock()
	state.active = nil
	state.starts = 0
	state.pauses = 0
	state.creates = 0
}

func (m *ServicePrincipalSynchronizationJobMock) SeedFixture() error {
	value, err := fixture()
	if err != nil {
		return err
	}
	state.Lock()
	defer state.Unlock()
	state.active = value
	return nil
}

func (m *ServicePrincipalSynchronizationJobMock) Actions() (int, int, int) {
	state.Lock()
	defer state.Unlock()
	return state.creates, state.starts, state.pauses
}
