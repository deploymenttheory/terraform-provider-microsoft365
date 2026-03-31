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

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	deployments map[string]map[string]any
}

func init() {
	mockState.deployments = make(map[string]map[string]any)

	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	mocks.GlobalRegistry.Register("windows_update_deployment", &WindowsUpdateDeploymentMock{})
}

type WindowsUpdateDeploymentMock struct{}

var _ mocks.MockRegistrar = (*WindowsUpdateDeploymentMock)(nil)

func (m *WindowsUpdateDeploymentMock) RegisterMocks() {
	m.registerGetLicenseResponder()
	m.registerCreateDeploymentResponder()
	m.registerGetDeploymentResponder()
	m.registerUpdateDeploymentResponder()
	m.registerDeleteDeploymentResponder()
}

func (m *WindowsUpdateDeploymentMock) registerGetLicenseResponder() {
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/subscribedSkus",
		m.getLicenseResponder())
}

func (m *WindowsUpdateDeploymentMock) registerCreateDeploymentResponder() {
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/admin/windows/updates/deployments",
		m.createDeploymentResponder())
}

func (m *WindowsUpdateDeploymentMock) registerGetDeploymentResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deployments/([^/]+)$`,
		m.getDeploymentResponder())
}

func (m *WindowsUpdateDeploymentMock) registerUpdateDeploymentResponder() {
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deployments/([^/]+)$`,
		m.updateDeploymentResponder())
}

func (m *WindowsUpdateDeploymentMock) registerDeleteDeploymentResponder() {
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deployments/([^/]+)$`,
		m.deleteDeploymentResponder())
}

func (m *WindowsUpdateDeploymentMock) createDeploymentResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_deployment_success.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}

		var response map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
		}

		id := uuid.New().String()
		response["id"] = id

		if content, ok := requestBody["content"]; ok {
			response["content"] = content
		}
		if settings, ok := requestBody["settings"]; ok {
			response["settings"] = settings
		} else {
			response["settings"] = nil
		}
		if audience, ok := requestBody["audience"]; ok {
			response["audience"] = audience
		}
		if state, ok := requestBody["state"].(map[string]any); ok {
			if responseState, ok := response["state"].(map[string]any); ok {
				for k, v := range state {
					responseState[k] = v
				}
			} else {
				response["state"] = state
			}
		}

		mockState.Lock()
		mockState.deployments[id] = response
		mockState.Unlock()

		return factories.SuccessResponse(201, response)(req)
	}
}

func (m *WindowsUpdateDeploymentMock) getDeploymentResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := factories.ExtractIDFromURL(req.URL.Path, "/admin/windows/updates/deployments/")

		mockState.Lock()
		deployment, exists := mockState.deployments[id]
		mockState.Unlock()

		if !exists {
			switch {
			case strings.Contains(id, "minimal"):
				jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_get", "get_deployment_minimal.json"))
				if err != nil {
					return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
				}

				var response map[string]any
				if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
					return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
				}
				response["id"] = id
				return factories.SuccessResponse(200, response)(req)
			default:
				jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_deployment_not_found.json"))
				if err == nil {
					var errorResponse map[string]any
					if json.Unmarshal([]byte(jsonContent), &errorResponse) == nil {
						return httpmock.NewJsonResponse(404, errorResponse)
					}
				}
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
			}
		}

		return factories.SuccessResponse(200, deployment)(req)
	}
}

func (m *WindowsUpdateDeploymentMock) updateDeploymentResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := factories.ExtractIDFromURL(req.URL.Path, "/admin/windows/updates/deployments/")

		mockState.Lock()
		deployment, exists := mockState.deployments[id]
		mockState.Unlock()

		if !exists {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_deployment_not_found.json"))
			if err == nil {
				var errorResponse map[string]any
				if json.Unmarshal([]byte(jsonContent), &errorResponse) == nil {
					return httpmock.NewJsonResponse(404, errorResponse)
				}
			}
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_update", "get_deployment_updated.json"))
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
		}

		var updatedDeployment map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &updatedDeployment); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse mock response"}}`), nil
		}

		for k, v := range deployment {
			updatedDeployment[k] = v
		}

		for k, v := range requestBody {
			updatedDeployment[k] = v
		}

		updatedDeployment["lastModifiedDateTime"] = "2024-01-01T12:00:00Z"

		mockState.Lock()
		mockState.deployments[id] = updatedDeployment
		mockState.Unlock()

		return factories.SuccessResponse(202, updatedDeployment)(req)
	}
}

func (m *WindowsUpdateDeploymentMock) deleteDeploymentResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := factories.ExtractIDFromURL(req.URL.Path, "/admin/windows/updates/deployments/")

		mockState.Lock()
		_, exists := mockState.deployments[id]
		if exists {
			delete(mockState.deployments, id)
		}
		mockState.Unlock()

		if !exists {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_delete", "get_deployment_not_found.json"))
			if err == nil {
				var errorResponse map[string]any
				if json.Unmarshal([]byte(jsonContent), &errorResponse) == nil {
					return httpmock.NewJsonResponse(404, errorResponse)
				}
			}
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return factories.EmptySuccessResponse(202)(req)
	}
}

func (m *WindowsUpdateDeploymentMock) getLicenseResponder() httpmock.Responder {
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

func (m *WindowsUpdateDeploymentMock) RegisterErrorMocks() {
	m.registerGetLicenseResponder()
	m.registerCreateDeploymentErrorResponder()
	m.registerGetDeploymentErrorResponder()
	m.registerUpdateDeploymentErrorResponder()
	m.registerDeleteDeploymentErrorResponder()
}

func (m *WindowsUpdateDeploymentMock) registerCreateDeploymentErrorResponder() {
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/admin/windows/updates/deployments",
		func(req *http.Request) (*http.Response, error) {
			jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", "post_deployment_error.json"))
			if err == nil {
				var errorResponse map[string]any
				if json.Unmarshal([]byte(jsonContent), &errorResponse) == nil {
					return httpmock.NewJsonResponse(400, errorResponse)
				}
			}
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`), nil
		})
}

func (m *WindowsUpdateDeploymentMock) registerGetDeploymentErrorResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deployments/error-id$`,
		factories.ErrorResponse(403, "Forbidden", "Access denied"))
}

func (m *WindowsUpdateDeploymentMock) registerUpdateDeploymentErrorResponder() {
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deployments/error-id$`,
		factories.ErrorResponse(500, "InternalServerError", "Internal server error"))
}

func (m *WindowsUpdateDeploymentMock) registerDeleteDeploymentErrorResponder() {
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deployments/error-id$`,
		factories.ErrorResponse(409, "Conflict", "Deployment is in use"))
}

func (m *WindowsUpdateDeploymentMock) CleanupMockState() {
	mockState.Lock()
	mockState.deployments = make(map[string]map[string]any)
	mockState.Unlock()
}
