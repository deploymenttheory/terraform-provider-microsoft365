package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	connections map[string]map[string]any
}

func init() {
	mockState.connections = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_update_updatable_asset_group", &WindowsUpdateUpdatableAssetGroupMock{})
}

type WindowsUpdateUpdatableAssetGroupMock struct{}

var _ mocks.MockRegistrar = (*WindowsUpdateUpdatableAssetGroupMock)(nil)

func (m *WindowsUpdateUpdatableAssetGroupMock) RegisterMocks() {
	mockState.Lock()
	mockState.connections = make(map[string]map[string]any)
	mockState.Unlock()

	m.registerCreateUpdatableAssetGroupResponder()
	m.registerGetUpdatableAssetGroupResponder()
	m.registerDeleteUpdatableAssetGroupResponder()
}

func (m *WindowsUpdateUpdatableAssetGroupMock) registerCreateUpdatableAssetGroupResponder() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatableAssets$`,
		func(req *http.Request) (*http.Response, error) {
			_, filename, _, _ := runtime.Caller(0)
			sourceDir := filepath.Dir(filename)
			responsesPath := filepath.Join(sourceDir, "..", "tests", "responses", constants.TfOperationCreate, "post_updatable_asset_group.json")

			jsonData, err := os.ReadFile(responsesPath)
			if err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load response: %s"}}`, err.Error())), nil
			}

			var responseObj map[string]any
			if err := json.Unmarshal(jsonData, &responseObj); err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse response: %s"}}`, err.Error())), nil
			}

			mockState.Lock()
			if id, ok := responseObj["id"].(string); ok {
				mockState.connections[id] = responseObj
			}
			mockState.Unlock()

			resp, err := httpmock.NewJsonResponse(200, responseObj)
			if err != nil {
				return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
			}
			return resp, nil
		})
}

func (m *WindowsUpdateUpdatableAssetGroupMock) registerGetUpdatableAssetGroupResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatableAssets/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			connectionId := parts[len(parts)-1]

			mockState.Lock()
			connection, exists := mockState.connections[connectionId]
			mockState.Unlock()

			if !exists {
				_, filename, _, _ := runtime.Caller(0)
				sourceDir := filepath.Dir(filename)
				responsesPath := filepath.Join(sourceDir, "..", "tests", "responses", constants.TfOperationRead, "get_updatable_asset_group.json")

				jsonData, err := os.ReadFile(responsesPath)
				if err != nil {
					return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"Updatable asset group not found"}}`), nil
				}

				var responseObj map[string]any
				if err := json.Unmarshal(jsonData, &responseObj); err != nil {
					return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse response: %s"}}`, err.Error())), nil
				}
				connection = responseObj
			}

			resp, err := httpmock.NewJsonResponse(200, connection)
			if err != nil {
				return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
			}
			return resp, nil
		})
}

func (m *WindowsUpdateUpdatableAssetGroupMock) registerDeleteUpdatableAssetGroupResponder() {
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatableAssets/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			connectionId := parts[len(parts)-1]

			mockState.Lock()
			delete(mockState.connections, connectionId)
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})
}

func (m *WindowsUpdateUpdatableAssetGroupMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.connections = make(map[string]map[string]any)
	mockState.Unlock()

	m.registerCreateUpdatableAssetGroupErrorResponder()
	m.registerGetUpdatableAssetGroupErrorResponder()
}

func (m *WindowsUpdateUpdatableAssetGroupMock) registerCreateUpdatableAssetGroupErrorResponder() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatableAssets`,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(403, map[string]any{
				"error": map[string]any{
					"code":    "Forbidden",
					"message": "Insufficient privileges to complete the operation.",
				},
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
			}
			return resp, nil
		})
}

func (m *WindowsUpdateUpdatableAssetGroupMock) registerGetUpdatableAssetGroupErrorResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatableAssets`,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(403, map[string]any{
				"error": map[string]any{
					"code":    "Forbidden",
					"message": "Insufficient privileges to complete the operation.",
				},
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
			}
			return resp, nil
		})
}

func (m *WindowsUpdateUpdatableAssetGroupMock) CleanupMockState() {
	mockState.Lock()
	mockState.connections = make(map[string]map[string]any)
	mockState.Unlock()
}
