package mocks

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

type WindowsUpdateDeploymentAudienceMock struct{}

var (
	mockState = struct {
		sync.Mutex
		audiences map[string]map[string]any
	}{
		audiences: make(map[string]map[string]any),
	}
)

func (m *WindowsUpdateDeploymentAudienceMock) RegisterMocks() {
	m.registerCreateAudienceResponder()
	m.registerGetAudienceResponder()
	m.registerDeleteAudienceResponder()
}

func (m *WindowsUpdateDeploymentAudienceMock) registerCreateAudienceResponder() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deploymentAudiences$`,
		m.createAudienceResponder())
}

func (m *WindowsUpdateDeploymentAudienceMock) registerGetAudienceResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deploymentAudiences/([^/]+)$`,
		m.getAudienceResponder())
}

func (m *WindowsUpdateDeploymentAudienceMock) registerDeleteAudienceResponder() {
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deploymentAudiences/([^/]+)$`,
		m.deleteAudienceResponder())
}

func (m *WindowsUpdateDeploymentAudienceMock) createAudienceResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		id := uuid.New().String()

		// Load response template from JSON file
		jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_create/post_audience_success.json")
		if err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load response JSON"}}`), nil
		}

		var response map[string]any
		if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
			return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse response JSON"}}`), nil
		}

		// Override with generated ID
		response["id"] = id

		mockState.Lock()
		mockState.audiences[id] = response
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, response)
	}
}

func (m *WindowsUpdateDeploymentAudienceMock) getAudienceResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := factories.ExtractIDFromURL(req.URL.Path, "/admin/windows/updates/deploymentAudiences/")

		mockState.Lock()
		audience, exists := mockState.audiences[id]
		mockState.Unlock()

		if !exists {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/error_handling/not_found.json")
			var errorResponse map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
			return httpmock.NewJsonResponse(404, errorResponse)
		}

		return httpmock.NewJsonResponse(200, audience)
	}
}

func (m *WindowsUpdateDeploymentAudienceMock) deleteAudienceResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		id := factories.ExtractIDFromURL(req.URL.Path, "/admin/windows/updates/deploymentAudiences/")

		mockState.Lock()
		_, exists := mockState.audiences[id]
		if exists {
			delete(mockState.audiences, id)
		}
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return httpmock.NewStringResponse(204, ""), nil
	}
}

func (m *WindowsUpdateDeploymentAudienceMock) RegisterErrorMocks() {
	m.registerCreateAudienceErrorResponder()
	m.registerGetAudienceErrorResponder()
	m.registerDeleteAudienceErrorResponder()
}

func (m *WindowsUpdateDeploymentAudienceMock) registerCreateAudienceErrorResponder() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deploymentAudiences$`,
		factories.ErrorResponse(400, "BadRequest", "Invalid request"))
}

func (m *WindowsUpdateDeploymentAudienceMock) registerGetAudienceErrorResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deploymentAudiences/error-id$`,
		factories.ErrorResponse(404, "ResourceNotFound", "Resource not found"))
}

func (m *WindowsUpdateDeploymentAudienceMock) registerDeleteAudienceErrorResponder() {
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/deploymentAudiences/error-id$`,
		factories.ErrorResponse(409, "Conflict", "Audience is in use"))
}

func (m *WindowsUpdateDeploymentAudienceMock) CleanupMockState() {
	mockState.Lock()
	mockState.audiences = make(map[string]map[string]any)
	mockState.Unlock()
}
