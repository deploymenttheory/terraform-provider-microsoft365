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
	devices map[string]map[string]any
}

func init() {
	mockState.devices = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("device", &DeviceMock{})
}

type DeviceMock struct{}

var _ mocks.MockRegistrar = (*DeviceMock)(nil)

func (m *DeviceMock) RegisterMocks() {
	mockState.Lock()
	mockState.devices = make(map[string]map[string]any)
	mockState.Unlock()

	RegisterGetByObjectIdMock()
	RegisterListAndFilterMocks()
	RegisterRelationshipMocks()
}

func (m *DeviceMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/devices`,
		httpmock.NewStringResponder(403, `{"error":{"code":"Forbidden","message":"Insufficient privileges to complete the operation"}}`))
}

func (m *DeviceMock) CleanupMockState() {
	mockState.Lock()
	mockState.devices = make(map[string]map[string]any)
	mockState.Unlock()
}

func RegisterGetByObjectIdMock() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/devices/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			objectId := parts[len(parts)-1]

			switch objectId {
			case "23ace577-ee29-416f-8566-11c948310bff":
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_device_by_object_id.json")
				var responseObj map[string]any
				json.Unmarshal([]byte(jsonStr), &responseObj)
				return httpmock.NewJsonResponse(200, responseObj)
			case "aaaaaaaa-1111-2222-3333-000000000002":
				responseObj := map[string]any{
					"id":                     "aaaaaaaa-1111-2222-3333-000000000002",
					"deviceId":               "bbbbbbbb-1111-2222-3333-000000000002",
					"displayName":            "DT-TEST-DEVICE-002",
					"operatingSystem":        "Windows",
					"operatingSystemVersion": "10.0.22621.2506",
					"isCompliant":            false,
					"isManaged":              true,
					"manufacturer":           "Microsoft Corporation",
					"model":                  "Surface Laptop 5",
					"trustType":              "AzureAd",
					"accountEnabled":         true,
				}
				return httpmock.NewJsonResponse(200, responseObj)
			case "cccccccc-2222-3333-4444-000000000003":
				responseObj := map[string]any{
					"id":                     "cccccccc-2222-3333-4444-000000000003",
					"deviceId":               "dddddddd-2222-3333-4444-000000000003",
					"displayName":            "DT-TEST-DEVICE-003",
					"operatingSystem":        "Windows",
					"operatingSystemVersion": "10.0.22631.4169",
					"isCompliant":            true,
					"isManaged":              true,
					"manufacturer":           "Lenovo",
					"model":                  "ThinkPad X1 Carbon",
					"trustType":              "AzureAd",
					"accountEnabled":         true,
				}
				return httpmock.NewJsonResponse(200, responseObj)
			default:
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Device not found"}}`), nil
			}
		})
}

func RegisterListAndFilterMocks() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/devices(\?.*)?$`,
		func(req *http.Request) (*http.Response, error) {
			queryParams := req.URL.Query()
			filter := queryParams.Get("$filter")
			count := queryParams.Get("$count")

			// List all devices
			if filter == "" && count == "" {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_devices_all.json")
				var responseObj map[string]any
				json.Unmarshal([]byte(jsonStr), &responseObj)
				return httpmock.NewJsonResponse(200, responseObj)
			}

			// Filter by display name
			if strings.Contains(filter, "displayName eq") {
				if strings.Contains(filter, "DT-000481110457") {
					jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_device_by_display_name.json")
					var responseObj map[string]any
					json.Unmarshal([]byte(jsonStr), &responseObj)
					return httpmock.NewJsonResponse(200, responseObj)
				}
			}

			// Filter by device ID
			if strings.Contains(filter, "deviceId eq") {
				if strings.Contains(filter, "06771871-1375-494e-97f9-ab87ba64edeb") {
					jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_device_by_device_id.json")
					var responseObj map[string]any
					json.Unmarshal([]byte(jsonStr), &responseObj)
					return httpmock.NewJsonResponse(200, responseObj)
				}
			}

			// OData filter for compliant Windows devices
			if strings.Contains(filter, "operatingSystem eq 'Windows'") && strings.Contains(filter, "isCompliant eq true") {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_devices_odata_filter.json")
				var responseObj map[string]any
				json.Unmarshal([]byte(jsonStr), &responseObj)
				return httpmock.NewJsonResponse(200, responseObj)
			}

			// Default empty response for unmatched filters
			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#devices",
				"value":          []any{},
			})
		})
}

func RegisterRelationshipMocks() {
	// Member of
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/devices/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/memberOf$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_device_member_of.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		})

	// Registered owners
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/devices/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/registeredOwners$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_device_registered_owners.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		})

	// Registered users
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/devices/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/registeredUsers$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_device_registered_users.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		})
}
