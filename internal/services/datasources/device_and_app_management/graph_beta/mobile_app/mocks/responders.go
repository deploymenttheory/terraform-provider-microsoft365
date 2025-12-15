package mocks

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	mobileApps map[string]map[string]any
}

func init() {
	mockState.mobileApps = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("mobile_apps", &MobileAppsMock{})
}

type MobileAppsMock struct{}

var _ mocks.MockRegistrar = (*MobileAppsMock)(nil)

func (m *MobileAppsMock) RegisterMocks() {
	mockState.Lock()
	mockState.mobileApps = make(map[string]map[string]any)
	mockState.Unlock()

	// 1. Get all mobile apps - GET /deviceAppManagement/mobileApps
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps", func(req *http.Request) (*http.Response, error) {
		// Parse query parameters
		queryParams, _ := url.ParseQuery(req.URL.RawQuery)

		// Handle different scenarios based on query parameters
		if filter := queryParams.Get("$filter"); filter != "" {
			// Handle OData filters
			if strings.Contains(filter, "startswith(publisher") {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_apps_by_publisher.json")
				var responseObj map[string]any
				json.Unmarshal([]byte(jsonStr), &responseObj)
				return httpmock.NewJsonResponse(200, responseObj)
			}
			if strings.Contains(filter, "startswith(displayName") {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_apps_by_display_name.json")
				var responseObj map[string]any
				json.Unmarshal([]byte(jsonStr), &responseObj)
				return httpmock.NewJsonResponse(200, responseObj)
			}
		}

		// Check for $expand parameter
		if expand := queryParams.Get("$expand"); expand != "" {
			// Return response with categories expanded
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_apps_all.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		}

		// Default: return all mobile apps
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_apps_all.json")
		var responseObj map[string]any
		json.Unmarshal([]byte(jsonStr), &responseObj)
		return httpmock.NewJsonResponse(200, responseObj)
	})

	// 2. Get mobile app by ID - GET /deviceAppManagement/mobileApps/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileApps/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		appId := parts[len(parts)-1]

		// Return mock response for known IDs
		switch appId {
		case "00000000-0000-0000-0000-000000000001": // Microsoft Edge
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_app_by_id.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		default:
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Mobile app not found"}}`), nil
		}
	})

	// 3. Handle OData queries with pagination simulation
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileApps\?.*`, func(req *http.Request) (*http.Response, error) {
		queryParams, _ := url.ParseQuery(req.URL.RawQuery)

		// Handle $top parameter
		if top := queryParams.Get("$top"); top != "" {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_apps_all.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		}

		// Handle $orderby parameter
		if orderBy := queryParams.Get("$orderby"); orderBy != "" && strings.Contains(orderBy, "displayName") {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_apps_all.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		}

		// Handle $select parameter
		if selectFields := queryParams.Get("$select"); selectFields != "" {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_apps_all.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		}

		// Default OData response
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_apps_all.json")
		var responseObj map[string]any
		json.Unmarshal([]byte(jsonStr), &responseObj)
		return httpmock.NewJsonResponse(200, responseObj)
	})
}

func (m *MobileAppsMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.mobileApps = make(map[string]map[string]any)
	mockState.Unlock()

	// Return errors for all operations
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceAppManagement/mobileApps", func(req *http.Request) (*http.Response, error) {
		errorObj := map[string]any{
			"error": map[string]any{
				"code":    "Forbidden",
				"message": "Insufficient privileges to complete the operation.",
			},
		}
		return httpmock.NewJsonResponse(403, errorObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileApps/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		errorObj := map[string]any{
			"error": map[string]any{
				"code":    "NotFound",
				"message": "Mobile app not found",
			},
		}
		return httpmock.NewJsonResponse(404, errorObj)
	})
}

func (m *MobileAppsMock) CleanupMockState() {
	mockState.Lock()
	mockState.mobileApps = make(map[string]map[string]any)
	mockState.Unlock()
}
