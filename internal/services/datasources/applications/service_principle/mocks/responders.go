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
	servicePrincipals map[string]map[string]any
}

func init() {
	mockState.servicePrincipals = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("service_principal", &ServicePrincipalMock{})
}

type ServicePrincipalMock struct{}

var _ mocks.MockRegistrar = (*ServicePrincipalMock)(nil)

func (m *ServicePrincipalMock) RegisterMocks() {
	mockState.Lock()
	mockState.servicePrincipals = make(map[string]map[string]any)
	mockState.Unlock()

	// 1. Get all service principals - GET /servicePrincipals
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/servicePrincipals", func(req *http.Request) (*http.Response, error) {
		// Parse query parameters
		queryParams, _ := url.ParseQuery(req.URL.RawQuery)
		
		// Handle different scenarios based on query parameters
		if filter := queryParams.Get("$filter"); filter != "" {
			if strings.Contains(filter, "preferredSingleSignOnMode ne 'notSupported'") {
				// Return filtered results
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_service_principals_odata_filter.json")
				var responseObj map[string]any
				json.Unmarshal([]byte(jsonStr), &responseObj)
				return httpmock.NewJsonResponse(200, responseObj)
			} else if strings.Contains(filter, "appId eq") {
				// Extract appId from filter
				parts := strings.Split(filter, "'")
				if len(parts) >= 2 {
					appId := parts[1]
					if appId == "63e61dc2-f593-4a6f-92b9-92e4d2c03d4f" {
						// Return single service principal for Microsoft Intune SCCM Connector
						jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_service_principal_by_id.json")
						var spObj map[string]any
						json.Unmarshal([]byte(jsonStr), &spObj)
						
						responseObj := map[string]any{
							"@odata.context": "https://graph.microsoft.com/beta/$metadata#servicePrincipals",
							"value":          []map[string]any{spObj},
						}
						return httpmock.NewJsonResponse(200, responseObj)
					}
				}
			}
		}
		
		// Default: return all service principals
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_service_principals_all.json")
		var responseObj map[string]any
		json.Unmarshal([]byte(jsonStr), &responseObj)
		return httpmock.NewJsonResponse(200, responseObj)
	})

	// 2. Get service principal by ID - GET /servicePrincipals/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		spId := parts[len(parts)-1]

		// Return mock response for known IDs
		switch spId {
		case "3b6f95b0-2064-4cc9-b5e5-1ab72af707b3":
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_service_principal_by_id.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		case "ac7ce817-df9d-4bce-aeb2-f006c182508d":
			// Microsoft Intune Service Discovery service principal
			responseObj := map[string]any{
				"id":               "ac7ce817-df9d-4bce-aeb2-f006c182508d",
				"appId":            "9cb77803-d937-493e-9a3b-4b49de3f5a74",
				"appDisplayName":   "Microsoft Intune Service Discovery",
				"displayName":      "Microsoft Intune Service Discovery",
				"accountEnabled":   true,
				"publisherName":    "Microsoft Services",
				"servicePrincipalNames": []string{
					"https://location.manage-beta.microsoft.com",
					"https://location.manage.microsoft.com",
					"https://location.manage.microsoft.us",
					"https://location.manage-test.microsoft.us",
					"https://location.manage-ppe.microsoft.us",
					"9cb77803-d937-493e-9a3b-4b49de3f5a74",
					"https://location.manage-mig.microsoft.com",
				},
				"preferredSingleSignOnMode": nil,
				"signInAudience":           "AzureADMultipleOrgs",
				"tags":                     []string{},
				"verifiedPublisher": map[string]any{
					"displayName":        nil,
					"verifiedPublisherId": nil,
					"addedDateTime":      nil,
				},
			}
			return httpmock.NewJsonResponse(200, responseObj)
		default:
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Service principal not found"}}`), nil
		}
	})

	// 3. Handle OData queries with pagination simulation
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/servicePrincipals\?.*`, func(req *http.Request) (*http.Response, error) {
		queryParams, _ := url.ParseQuery(req.URL.RawQuery)
		
		// Handle $count parameter
		if queryParams.Get("$count") == "true" {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_service_principals_odata_filter.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			responseObj["@odata.count"] = 2
			return httpmock.NewJsonResponse(200, responseObj)
		}
		
		// Handle $orderby parameter
		if orderBy := queryParams.Get("$orderby"); orderBy != "" && strings.Contains(orderBy, "displayName") {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_service_principals_odata_filter.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		}
		
		// Handle $select parameter
		if selectFields := queryParams.Get("$select"); selectFields != "" {
			// Return limited fields based on select
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_service_principals_odata_filter.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		}
		
		// Default OData response
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_service_principals_all.json")
		var responseObj map[string]any
		json.Unmarshal([]byte(jsonStr), &responseObj)
		return httpmock.NewJsonResponse(200, responseObj)
	})
}

func (m *ServicePrincipalMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.servicePrincipals = make(map[string]map[string]any)
	mockState.Unlock()

	// Return errors for all operations
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/servicePrincipals", func(req *http.Request) (*http.Response, error) {
		errorObj := map[string]any{
			"error": map[string]any{
				"code":    "Forbidden",
				"message": "Insufficient privileges to complete the operation.",
			},
		}
		return httpmock.NewJsonResponse(403, errorObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		errorObj := map[string]any{
			"error": map[string]any{
				"code":    "NotFound",
				"message": "Service principal not found",
			},
		}
		return httpmock.NewJsonResponse(404, errorObj)
	})
}

func (m *ServicePrincipalMock) CleanupMockState() {
	mockState.Lock()
	mockState.servicePrincipals = make(map[string]map[string]any)
	mockState.Unlock()
}