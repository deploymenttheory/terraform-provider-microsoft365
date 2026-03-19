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
	enrolledDevices map[string]map[string]any
}

func init() {
	mockState.enrolledDevices = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_autopatch_device_registration", &WindowsAutopatchDeviceRegistrationMock{})
}

type WindowsAutopatchDeviceRegistrationMock struct{}

var _ mocks.MockRegistrar = (*WindowsAutopatchDeviceRegistrationMock)(nil)

func (m *WindowsAutopatchDeviceRegistrationMock) RegisterMocks() {
	mockState.Lock()
	mockState.enrolledDevices = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/admin/windows/updates/updatableAssets/enrollAssetsById",
		func(req *http.Request) (*http.Response, error) {
			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewJsonResponse(400, map[string]any{
					"error": map[string]any{
						"code":    "BadRequest",
						"message": "Invalid request body",
					},
				})
			}

			updateCategory, ok := body["updateCategory"].(string)
			if !ok {
				return httpmock.NewJsonResponse(400, map[string]any{
					"error": map[string]any{
						"code":    "BadRequest",
						"message": "Missing updateCategory",
					},
				})
			}

			ids, ok := body["ids"].([]any)
			if !ok {
				return httpmock.NewJsonResponse(400, map[string]any{
					"error": map[string]any{
						"code":    "BadRequest",
						"message": "Missing ids",
					},
				})
			}

			mockState.Lock()
			for _, idAny := range ids {
				if id, ok := idAny.(string); ok {
					key := fmt.Sprintf("%s_%s", id, updateCategory)
					mockState.enrolledDevices[key] = map[string]any{
						"id":             id,
						"updateCategory": updateCategory,
						"enrollments": []map[string]any{
							{
								"@odata.type":  "#microsoft.graph.windowsUpdates.updateManagementEnrollment",
								"updateCategory": updateCategory,
							},
						},
					}
				}
			}
			mockState.Unlock()

			return httpmock.NewStringResponse(202, ""), nil
		})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/admin/windows/updates/updatableAssets/unenrollAssetsById",
		func(req *http.Request) (*http.Response, error) {
			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewJsonResponse(400, map[string]any{
					"error": map[string]any{
						"code":    "BadRequest",
						"message": "Invalid request body",
					},
				})
			}

			updateCategory, ok := body["updateCategory"].(string)
			if !ok {
				return httpmock.NewJsonResponse(400, map[string]any{
					"error": map[string]any{
						"code":    "BadRequest",
						"message": "Missing updateCategory",
					},
				})
			}

			ids, ok := body["ids"].([]any)
			if !ok {
				return httpmock.NewJsonResponse(400, map[string]any{
					"error": map[string]any{
						"code":    "BadRequest",
						"message": "Missing ids",
					},
				})
			}

			mockState.Lock()
			for _, idAny := range ids {
				if id, ok := idAny.(string); ok {
					key := fmt.Sprintf("%s_%s", id, updateCategory)
					delete(mockState.enrolledDevices, key)
				}
			}
			mockState.Unlock()

			return httpmock.NewStringResponse(202, ""), nil
		})

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/admin/windows/updates/updatableAssets/microsoft.graph.windowsUpdates.azureADDevice",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			devices := make([]map[string]any, 0, len(mockState.enrolledDevices))
			for _, deviceData := range mockState.enrolledDevices {
				devices = append(devices, map[string]any{
					"@odata.type": "#microsoft.graph.windowsUpdates.azureADDevice",
					"id":          deviceData["id"],
					"enrollments": deviceData["enrollments"],
					"errors":      []any{},
				})
			}
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#admin/windows/updates/updatableAssets",
				"value":          devices,
			})
		})
}

func (m *WindowsAutopatchDeviceRegistrationMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.enrolledDevices = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/admin/windows/updates/updatableAssets/enrollAssetsById",
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_error_scenario.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(400, errObj)
		})

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/admin/windows/updates/updatableAssets/microsoft.graph.windowsUpdates.azureADDevice",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#admin/windows/updates/updatableAssets",
				"value":          []any{},
			})
		})
}

func (m *WindowsAutopatchDeviceRegistrationMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for key := range mockState.enrolledDevices {
		delete(mockState.enrolledDevices, key)
	}
}
