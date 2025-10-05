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
	mobileAppCatalogPackages map[string]map[string]any
}

func init() {
	mockState.mobileAppCatalogPackages = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("mobile_app_catalog_packages", &MobileAppCatalogPackagesMock{})
}

type MobileAppCatalogPackagesMock struct{}

var _ mocks.MockRegistrar = (*MobileAppCatalogPackagesMock)(nil)

func (m *MobileAppCatalogPackagesMock) RegisterMocks() {
	mockState.Lock()
	mockState.mobileAppCatalogPackages = make(map[string]map[string]any)
	mockState.Unlock()

	// 1. Get all mobile app catalog packages - GET /deviceAppManagement/mobileAppCatalogPackages
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceAppManagement/mobileAppCatalogPackages", func(req *http.Request) (*http.Response, error) {
		// Parse query parameters
		queryParams, _ := url.ParseQuery(req.URL.RawQuery)

		// Handle different scenarios based on query parameters
		if filter := queryParams.Get("$filter"); filter != "" {
			if strings.Contains(filter, "productId eq") {
				// Extract productId from filter
				parts := strings.Split(filter, "'")
				if len(parts) >= 2 {
					productId := parts[1]
					if productId == "3a6307ef-6991-faf1-01e1-35e1557287aa" {
						// Return single package for 7-Zip
						jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_app_catalog_package_by_id.json")
						var packageObj map[string]any
						json.Unmarshal([]byte(jsonStr), &packageObj)

						responseObj := map[string]any{
							"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/mobileAppCatalogPackages",
							"value":          []map[string]any{packageObj},
						}
						return httpmock.NewJsonResponse(200, responseObj)
					}
				}
			} else if strings.Contains(filter, "productDisplayName eq") {
				// Handle product name filter
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_app_catalog_packages_odata_filter.json")
				var responseObj map[string]any
				json.Unmarshal([]byte(jsonStr), &responseObj)
				return httpmock.NewJsonResponse(200, responseObj)
			} else if strings.Contains(filter, "publisherDisplayName eq") {
				// Handle publisher name filter
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_app_catalog_packages_odata_filter.json")
				var responseObj map[string]any
				json.Unmarshal([]byte(jsonStr), &responseObj)
				return httpmock.NewJsonResponse(200, responseObj)
			}
		}

		// Default: return all mobile app catalog packages
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_app_catalog_packages_all.json")
		var responseObj map[string]any
		json.Unmarshal([]byte(jsonStr), &responseObj)
		return httpmock.NewJsonResponse(200, responseObj)
	})

	// 2. Get mobile app catalog package by ID - GET /deviceAppManagement/mobileAppCatalogPackages/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileAppCatalogPackages/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		packageId := parts[len(parts)-1]

		// Return mock response for known IDs
		switch packageId {
		case "5af1ade9-6966-3608-7e04-848252e29681":
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_app_catalog_package_by_id.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		default:
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Mobile app catalog package not found"}}`), nil
		}
	})

	// 3. Handle OData queries with pagination simulation
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileAppCatalogPackages\?.*`, func(req *http.Request) (*http.Response, error) {
		queryParams, _ := url.ParseQuery(req.URL.RawQuery)

		// Handle $count parameter
		if queryParams.Get("$count") == "true" {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_app_catalog_packages_odata_filter.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			responseObj["@odata.count"] = 2
			return httpmock.NewJsonResponse(200, responseObj)
		}

		// Handle $orderby parameter
		if orderBy := queryParams.Get("$orderby"); orderBy != "" && strings.Contains(orderBy, "productDisplayName") {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_app_catalog_packages_odata_filter.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		}

		// Handle $select parameter
		if selectFields := queryParams.Get("$select"); selectFields != "" {
			// Return limited fields based on select
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_app_catalog_packages_odata_filter.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		}

		// Handle $search parameter
		if search := queryParams.Get("$search"); search != "" {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_app_catalog_packages_odata_filter.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		}

		// Default OData response
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_mobile_app_catalog_packages_all.json")
		var responseObj map[string]any
		json.Unmarshal([]byte(jsonStr), &responseObj)
		return httpmock.NewJsonResponse(200, responseObj)
	})
}

func (m *MobileAppCatalogPackagesMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.mobileAppCatalogPackages = make(map[string]map[string]any)
	mockState.Unlock()

	// Return errors for all operations
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceAppManagement/mobileAppCatalogPackages", func(req *http.Request) (*http.Response, error) {
		errorObj := map[string]any{
			"error": map[string]any{
				"code":    "Forbidden",
				"message": "Insufficient privileges to complete the operation.",
			},
		}
		return httpmock.NewJsonResponse(403, errorObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceAppManagement/mobileAppCatalogPackages/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		errorObj := map[string]any{
			"error": map[string]any{
				"code":    "NotFound",
				"message": "Mobile app catalog package not found",
			},
		}
		return httpmock.NewJsonResponse(404, errorObj)
	})
}

func (m *MobileAppCatalogPackagesMock) CleanupMockState() {
	mockState.Lock()
	mockState.mobileAppCatalogPackages = make(map[string]map[string]any)
	mockState.Unlock()
}
