package mocks

import (
	"encoding/json"
	"fmt"
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

	RegisterCategoriesMock()
	RegisterGetByIdMock()
	RegisterListAndFilterMocks()
}

func RegisterGetByIdMock() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileApps/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			appId := parts[len(parts)-1]

			switch appId {
			case "00000000-0000-0000-0000-000000000001":
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_app_by_id.json")
				var responseObj map[string]any
				json.Unmarshal([]byte(jsonStr), &responseObj)
				return httpmock.NewJsonResponse(200, responseObj)
			default:
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Mobile app not found"}}`), nil
			}
		})
}

func RegisterListAndFilterMocks() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileApps`,
		func(req *http.Request) (*http.Response, error) {
			queryParams, _ := url.ParseQuery(req.URL.RawQuery)
			filter := queryParams.Get("$filter")

			if filter != "" {
				if strings.Contains(filter, "contains(tolower(publisher)") || strings.Contains(filter, "startswith(publisher") {
					jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_apps_by_publisher.json")
					var responseObj map[string]any
					json.Unmarshal([]byte(jsonStr), &responseObj)
					return httpmock.NewJsonResponse(200, responseObj)
				}
				if strings.Contains(filter, "contains(tolower(displayName)") || strings.Contains(filter, "startswith(displayName") {
					jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_apps_by_display_name.json")
					var responseObj map[string]any
					json.Unmarshal([]byte(jsonStr), &responseObj)
					return httpmock.NewJsonResponse(200, responseObj)
				}
				if strings.Contains(filter, "contains(tolower(developer)") {
					jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_apps_by_developer.json")
					var responseObj map[string]any
					json.Unmarshal([]byte(jsonStr), &responseObj)
					return httpmock.NewJsonResponse(200, responseObj)
				}
			}

			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_apps_all.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		})
}

func RegisterCategoriesMock() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileApps/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/categories`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			appId := parts[len(parts)-2]

			switch appId {
			case "00000000-0000-0000-0000-000000000001":
				responseObj := map[string]any{
					"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/mobileApps/categories",
					"value": []map[string]any{
						{
							"@odata.type": "#microsoft.graph.mobileAppCategory",
							"id":          "cat-001",
							"displayName": "Productivity",
						},
					},
				}
				return httpmock.NewJsonResponse(200, responseObj)
			case "00000000-0000-0000-0000-000000000003":
				responseObj := map[string]any{
					"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/mobileApps/categories",
					"value": []map[string]any{
						{
							"@odata.type": "#microsoft.graph.mobileAppCategory",
							"id":          "cat-002",
							"displayName": "Communication",
						},
					},
				}
				return httpmock.NewJsonResponse(200, responseObj)
			default:
				responseObj := map[string]any{
					"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/mobileApps/categories",
					"value":          []map[string]any{},
				}
				return httpmock.NewJsonResponse(200, responseObj)
			}
		})
}

func (m *MobileAppsMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.mobileApps = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileApps`, func(req *http.Request) (*http.Response, error) {
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

func (m *MobileAppsMock) CleanupMockState() {
	mockState.Lock()
	mockState.mobileApps = make(map[string]map[string]any)
	mockState.Unlock()
}
