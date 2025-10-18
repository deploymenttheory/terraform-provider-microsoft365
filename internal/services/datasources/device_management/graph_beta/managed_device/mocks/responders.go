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
	managedDevices map[string]map[string]any
}

func init() {
	mockState.managedDevices = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("managed_device", &ManagedDeviceMock{})
}

type ManagedDeviceMock struct{}

var _ mocks.MockRegistrar = (*ManagedDeviceMock)(nil)

func (m *ManagedDeviceMock) RegisterMocks() {
	mockState.Lock()
	mockState.managedDevices = make(map[string]map[string]any)
	mockState.Unlock()

	// 1. Get all managed devices - GET /deviceManagement/managedDevices
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/managedDevices", func(req *http.Request) (*http.Response, error) {
		queryParams, _ := url.ParseQuery(req.URL.RawQuery)

		// Handle different scenarios based on query parameters
		if filter := queryParams.Get("$filter"); filter != "" {
			if strings.Contains(filter, "complianceState eq 'compliant'") ||
				strings.Contains(filter, "operatingSystem eq 'Windows'") {
				// Return filtered results
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_managed_devices_odata_filter.json")
				var responseObj map[string]any
				json.Unmarshal([]byte(jsonStr), &responseObj)
				return httpmock.NewJsonResponse(200, responseObj)
			}
		}

		// Default: return all managed devices
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_managed_devices_all.json")
		var responseObj map[string]any
		json.Unmarshal([]byte(jsonStr), &responseObj)
		return httpmock.NewJsonResponse(200, responseObj)
	})

	// 2. Get managed device by ID - GET /deviceManagement/managedDevices/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/managedDevices/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		deviceId := parts[len(parts)-1]

		// Return mock response for known IDs
		switch deviceId {
		case "00000000-0000-0000-0000-000000000001":
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_managed_device_by_id.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		case "00000000-0000-0000-0000-000000000002":
			// Return second device
			responseObj := map[string]any{
				"id":                        "00000000-0000-0000-0000-000000000002",
				"userId":                    "user2@contoso.com",
				"deviceName":                "DESKTOP-WIN-002",
				"managedDeviceOwnerType":    "company",
				"operatingSystem":           "Windows",
				"complianceState":           "noncompliant",
				"osVersion":                 "10.0.22621.2506",
				"serialNumber":              "SN-WIN-002",
				"manufacturer":              "Microsoft Corporation",
				"model":                     "Surface Laptop 5",
				"userPrincipalName":         "user2@contoso.com",
				"enrolledDateTime":          "2024-02-20T14:15:00Z",
				"lastSyncDateTime":          "2024-10-18T08:45:00Z",
				"azureADDeviceId":           "aaaaaaaa-0000-0000-0000-000000000002",
				"deviceRegistrationState":   "registered",
				"deviceCategoryDisplayName": "Corporate",
				"isEncrypted":               true,
				"managementAgent":           "mdm",
			}
			return httpmock.NewJsonResponse(200, responseObj)
		case "00000000-0000-0000-0000-000000000003":
			// Return third device
			responseObj := map[string]any{
				"id":                        "00000000-0000-0000-0000-000000000003",
				"userId":                    "user3@contoso.com",
				"deviceName":                "LAPTOP-WIN-003",
				"managedDeviceOwnerType":    "company",
				"operatingSystem":           "Windows",
				"complianceState":           "compliant",
				"osVersion":                 "10.0.22631.4169",
				"serialNumber":              "SN-WIN-003",
				"manufacturer":              "Lenovo",
				"model":                     "ThinkPad X1 Carbon",
				"userPrincipalName":         "user3@contoso.com",
				"enrolledDateTime":          "2024-03-10T11:20:00Z",
				"lastSyncDateTime":          "2024-10-18T09:30:00Z",
				"azureADDeviceId":           "aaaaaaaa-0000-0000-0000-000000000003",
				"deviceRegistrationState":   "registered",
				"deviceCategoryDisplayName": "Corporate",
				"isEncrypted":               true,
				"managementAgent":           "mdm",
			}
			return httpmock.NewJsonResponse(200, responseObj)
		default:
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Managed device not found"}}`), nil
		}
	})

	// 3. Handle OData queries with pagination simulation
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/managedDevices\?.*`, func(req *http.Request) (*http.Response, error) {
		queryParams, _ := url.ParseQuery(req.URL.RawQuery)

		// Handle $count parameter
		if queryParams.Get("$count") == "true" {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_managed_devices_odata_filter.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			responseObj["@odata.count"] = 2
			return httpmock.NewJsonResponse(200, responseObj)
		}

		// Handle $orderby parameter
		if orderBy := queryParams.Get("$orderby"); orderBy != "" && (strings.Contains(orderBy, "deviceName") || strings.Contains(orderBy, "lastSyncDateTime")) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_managed_devices_odata_filter.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		}

		// Handle $select parameter
		if selectFields := queryParams.Get("$select"); selectFields != "" {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_managed_devices_odata_filter.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		}

		// Default OData response
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_managed_devices_all.json")
		var responseObj map[string]any
		json.Unmarshal([]byte(jsonStr), &responseObj)
		return httpmock.NewJsonResponse(200, responseObj)
	})
}

func (m *ManagedDeviceMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.managedDevices = make(map[string]map[string]any)
	mockState.Unlock()

	// Return errors for all operations
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/managedDevices", func(req *http.Request) (*http.Response, error) {
		errorObj := map[string]any{
			"error": map[string]any{
				"code":    "Forbidden",
				"message": "Insufficient privileges to complete the operation.",
			},
		}
		return httpmock.NewJsonResponse(403, errorObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/managedDevices/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		errorObj := map[string]any{
			"error": map[string]any{
				"code":    "NotFound",
				"message": "Managed device not found",
			},
		}
		return httpmock.NewJsonResponse(404, errorObj)
	})
}

func (m *ManagedDeviceMock) CleanupMockState() {
	mockState.Lock()
	mockState.managedDevices = make(map[string]map[string]any)
	mockState.Unlock()
}
