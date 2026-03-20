package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

	m.registerEntraDeviceValidationMocks()
	m.registerEnrollAssetsResponder()
	m.registerUnenrollAssetsResponder()
	m.registerIndividualAssetGetResponder()
	m.registerUpdatableAssetsListResponder()
}

func (m *WindowsAutopatchDeviceRegistrationMock) registerEntraDeviceValidationMocks() {
	knownTestDeviceIDs := []string{
		"12345678-1234-1234-1234-123456789001",
		"12345678-1234-1234-1234-123456789002",
		"12345678-1234-1234-1234-123456789003",
	}
	for _, deviceID := range knownTestDeviceIDs {
		id := deviceID
		httpmock.RegisterResponder("GET",
			fmt.Sprintf("https://graph.microsoft.com/beta/devices/%s", id),
			httpmock.NewJsonResponderOrPanic(200, map[string]any{
				"@odata.context":  "https://graph.microsoft.com/beta/$metadata#devices/$entity",
				"id":              id,
				"displayName":     fmt.Sprintf("TestDevice-%s", id),
				"operatingSystem": "Windows",
			}))
	}
}

func (m *WindowsAutopatchDeviceRegistrationMock) registerEnrollAssetsResponder() {
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/admin/windows/updates/updatableAssets/microsoft.graph.windowsUpdates.enrollAssetsById",
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
					enrollmentObj := map[string]any{}
					enrollmentObj[updateCategory] = map[string]any{}
					mockState.enrolledDevices[key] = map[string]any{
						"id":             id,
						"updateCategory": updateCategory,
						"enrollment":     enrollmentObj,
					}
				}
			}
			mockState.Unlock()

			return httpmock.NewStringResponse(202, ""), nil
		})
}

func (m *WindowsAutopatchDeviceRegistrationMock) registerUnenrollAssetsResponder() {
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/admin/windows/updates/updatableAssets/microsoft.graph.windowsUpdates.unenrollAssetsById",
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
}

func (m *WindowsAutopatchDeviceRegistrationMock) registerIndividualAssetGetResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatableAssets/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			deviceID := parts[len(parts)-1]

			mockState.Lock()
			var found map[string]any
			for _, data := range mockState.enrolledDevices {
				if id, ok := data["id"].(string); ok && id == deviceID {
					found = data
					break
				}
			}
			mockState.Unlock()

			if found == nil {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
			}

			updateCategory, _ := found["updateCategory"].(string)
			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.type": "#microsoft.graph.windowsUpdates.azureADDevice",
				"id":          deviceID,
				"enrollment":  map[string]any{updateCategory: map[string]any{}},
				"errors":      []any{},
			})
		})
}

func (m *WindowsAutopatchDeviceRegistrationMock) registerUpdatableAssetsListResponder() {
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/admin/windows/updates/updatableAssets",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			devices := make([]map[string]any, 0, len(mockState.enrolledDevices))

			allTestDeviceIDs := []string{
				"12345678-1234-1234-1234-123456789001",
				"12345678-1234-1234-1234-123456789002",
				"12345678-1234-1234-1234-123456789003",
			}

			for _, deviceID := range allTestDeviceIDs {
				enrollment := map[string]any{
					"feature": map[string]any{"enrollmentState": "notEnrolled"},
					"quality": map[string]any{"enrollmentState": "notEnrolled"},
					"driver":  map[string]any{"enrollmentState": "notEnrolled"},
				}

				for _, deviceData := range mockState.enrolledDevices {
					if id, ok := deviceData["id"].(string); ok && id == deviceID {
						if updateCategory, ok := deviceData["updateCategory"].(string); ok {
							enrollment[updateCategory] = map[string]any{"enrollmentState": "enrolled"}
						}
					}
				}

				devices = append(devices, map[string]any{
					"@odata.type": "#microsoft.graph.windowsUpdates.azureADDevice",
					"id":          deviceID,
					"enrollment":  enrollment,
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

	m.registerEntraDeviceValidationMocks()
	m.registerEnrollAssetsErrorResponder()
	m.registerUpdatableAssetsListErrorResponder()
}

func (m *WindowsAutopatchDeviceRegistrationMock) registerEnrollAssetsErrorResponder() {
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/admin/windows/updates/updatableAssets/microsoft.graph.windowsUpdates.enrollAssetsById",
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_error_scenario.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(400, errObj)
		})
}

func (m *WindowsAutopatchDeviceRegistrationMock) registerUpdatableAssetsListErrorResponder() {
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/admin/windows/updates/updatableAssets",
		func(req *http.Request) (*http.Response, error) {
			allTestDeviceIDs := []string{
				"12345678-1234-1234-1234-123456789001",
				"12345678-1234-1234-1234-123456789002",
				"12345678-1234-1234-1234-123456789003",
			}

			devices := make([]map[string]any, 0, len(allTestDeviceIDs))
			for _, deviceID := range allTestDeviceIDs {
				devices = append(devices, map[string]any{
					"@odata.type": "#microsoft.graph.windowsUpdates.azureADDevice",
					"id":          deviceID,
					"enrollment": map[string]any{
						"feature": map[string]any{"enrollmentState": "notEnrolled"},
						"quality": map[string]any{"enrollmentState": "notEnrolled"},
						"driver":  map[string]any{"enrollmentState": "notEnrolled"},
					},
					"errors": []any{},
				})
			}

			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#admin/windows/updates/updatableAssets",
				"value":          devices,
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
