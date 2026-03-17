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

var mockState struct {
	sync.Mutex
	deployments map[string]map[string]any
}

func init() {
	mockState.deployments = make(map[string]map[string]any)

	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	mocks.GlobalRegistry.Register("windows_update_deployment_state", &WindowsUpdatesAutopatchDeploymentStateMock{})
}

type WindowsUpdatesAutopatchDeploymentStateMock struct{}

var _ mocks.MockRegistrar = (*WindowsUpdatesAutopatchDeploymentStateMock)(nil)

func (m *WindowsUpdatesAutopatchDeploymentStateMock) RegisterMocks() {
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/subscribedSkus",
		m.getLicenseResponder())

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deployments/([^/]+)$`,
		m.getDeploymentResponder())

	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deployments/([^/]+)$`,
		m.patchDeploymentStateResponder())
}

func (m *WindowsUpdatesAutopatchDeploymentStateMock) patchDeploymentStateResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := factories.ExtractIDFromURL(req.URL.Path, "/admin/windows/updates/deployments/")

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "patch_state_success.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}

		var response map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
		}

		response["id"] = id

		if state, ok := requestBody["state"].(map[string]any); ok {
			if existingState, ok := response["state"].(map[string]any); ok {
				for k, v := range state {
					existingState[k] = v
				}
				if rv, ok := state["requestedValue"].(string); ok {
					existingState["effectiveValue"] = rv
				}
			}
		}

		mockState.Lock()
		mockState.deployments[id] = response
		mockState.Unlock()

		return factories.SuccessResponse(202, response)(req)
	}
}

func (m *WindowsUpdatesAutopatchDeploymentStateMock) getDeploymentResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := factories.ExtractIDFromURL(req.URL.Path, "/admin/windows/updates/deployments/")

		mockState.Lock()
		deployment, exists := mockState.deployments[id]
		mockState.Unlock()

		if !exists {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_get", "get_deployment_with_state.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}

			var response map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
			}
			response["id"] = id
			return factories.SuccessResponse(200, response)(req)
		}

		return factories.SuccessResponse(200, deployment)(req)
	}
}

func (m *WindowsUpdatesAutopatchDeploymentStateMock) getLicenseResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "license", "get_subscribed_skus_success.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load license mock"}}`), nil
		}

		var response map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse license mock"}}`), nil
		}
		return factories.SuccessResponse(200, response)(req)
	}
}

func (m *WindowsUpdatesAutopatchDeploymentStateMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/subscribedSkus",
		m.getLicenseResponder())

	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deployments/([^/]+)$`,
		factories.ErrorResponse(400, "BadRequest", "Invalid deployment configuration"))
}

func (m *WindowsUpdatesAutopatchDeploymentStateMock) CleanupMockState() {
	mockState.Lock()
	mockState.deployments = make(map[string]map[string]any)
	mockState.Unlock()
}
