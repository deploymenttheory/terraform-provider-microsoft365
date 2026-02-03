package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	applications map[string]map[string]any
}

func init() {
	mockState.applications = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("application", &ApplicationMock{})
}

type ApplicationMock struct{}

var _ mocks.MockRegistrar = (*ApplicationMock)(nil)

func (m *ApplicationMock) RegisterMocks() {
	mockState.Lock()
	mockState.applications = make(map[string]map[string]any)
	mockState.Unlock()

	// 1. Get all applications - GET /applications with query parameters
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/applications", func(req *http.Request) (*http.Response, error) {
		queryParams := req.URL.Query()
		filter := queryParams.Get("$filter")

		// Handle different filter scenarios with specific JSON files
		if filter != "" {
			// Display name filter
			if strings.Contains(filter, "displayName eq 'Test Application'") {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_application_by_display_name.json")
				var responseObj map[string]any
				json.Unmarshal([]byte(jsonStr), &responseObj)
				return httpmock.NewJsonResponse(200, responseObj)
			}

			// App ID filter - exact match
			if filter == "appId eq '12345678-1234-1234-1234-123456789012'" {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_application_by_app_id.json")
				var responseObj map[string]any
				json.Unmarshal([]byte(jsonStr), &responseObj)
				return httpmock.NewJsonResponse(200, responseObj)
			}

			// OData filter queries
			if strings.Contains(filter, "displayName eq 'Test Application' and signInAudience eq 'AzureADMyOrg'") ||
				strings.Contains(filter, "tags/any(t:t eq 'MyCustomTag')") {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_application_odata_filter.json")
				var responseObj map[string]any
				json.Unmarshal([]byte(jsonStr), &responseObj)
				return httpmock.NewJsonResponse(200, responseObj)
			}
		}

		// Default: return empty list for unmocked queries
		responseObj := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#applications",
			"value":          []map[string]any{},
		}
		return httpmock.NewJsonResponse(200, responseObj)
	})

	// 2. Get application by ID - GET /applications/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		appId := parts[len(parts)-1]

		// Return mock response for known IDs
		switch appId {
		case "a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d":
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_application_by_id.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj), nil
		default:
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Application not found"}}`), nil
		}
	})
}

func (m *ApplicationMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.applications = make(map[string]map[string]any)
	mockState.Unlock()

	// Return errors for all operations
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/applications", func(req *http.Request) (*http.Response, error) {
		errorObj := map[string]any{
			"error": map[string]any{
				"code":    "Forbidden",
				"message": "Insufficient privileges to complete the operation.",
			},
		}
		return httpmock.NewJsonResponse(403, errorObj), nil
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/applications/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		errorObj := map[string]any{
			"error": map[string]any{
				"code":    "Forbidden",
				"message": "Insufficient privileges to complete the operation.",
			},
		}
		return httpmock.NewJsonResponse(403, errorObj), nil
	})
}

func (m *ApplicationMock) CleanupMockState() {
	mockState.Lock()
	mockState.applications = make(map[string]map[string]any)
	mockState.Unlock()
}
