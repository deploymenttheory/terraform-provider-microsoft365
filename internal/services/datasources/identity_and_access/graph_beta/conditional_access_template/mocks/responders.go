package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	templates map[string]map[string]any
}

func init() {
	mockState.templates = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("conditional_access_template", &ConditionalAccessTemplateMock{})
}

type ConditionalAccessTemplateMock struct{}

var _ mocks.MockRegistrar = (*ConditionalAccessTemplateMock)(nil)

func (m *ConditionalAccessTemplateMock) RegisterMocks() {
	mockState.Lock()
	mockState.templates = make(map[string]map[string]any)
	mockState.Unlock()

	// GET /conditionalAccess/templates
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/templates$`, func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_template_list.json")
		var responseObj map[string]any
		json.Unmarshal([]byte(jsonStr), &responseObj)
		resp, err := httpmock.NewJsonResponse(200, responseObj)
		if err != nil {
			return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
		}
		return resp, nil
	})
}

func (m *ConditionalAccessTemplateMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.templates = make(map[string]map[string]any)
	mockState.Unlock()

	// Return errors for all operations
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/templates$`, func(req *http.Request) (*http.Response, error) {
		errorObj := map[string]any{
			"error": map[string]any{
				"code":    "Forbidden",
				"message": "Insufficient privileges to complete the operation.",
			},
		}
		resp, err := httpmock.NewJsonResponse(403, errorObj)
		if err != nil {
			return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
		}
		return resp, nil
	})
}

func (m *ConditionalAccessTemplateMock) CleanupMockState() {
	mockState.Lock()
	mockState.templates = make(map[string]map[string]any)
	mockState.Unlock()
}
