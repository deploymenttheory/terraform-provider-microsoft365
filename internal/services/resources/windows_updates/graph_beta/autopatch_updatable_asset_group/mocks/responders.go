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
	groups map[string]map[string]any
}

func init() {
	mockState.groups = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_update_updatable_asset_group", &WindowsUpdateUpdatableAssetGroupMock{})
}

type WindowsUpdateUpdatableAssetGroupMock struct{}

var _ mocks.MockRegistrar = (*WindowsUpdateUpdatableAssetGroupMock)(nil)

func (m *WindowsUpdateUpdatableAssetGroupMock) RegisterMocks() {
	mockState.Lock()
	mockState.groups = make(map[string]map[string]any)
	mockState.Unlock()

	m.registerCreateUpdatableAssetGroupResponder()
	m.registerGetUpdatableAssetGroupResponder()
	m.registerGetGroupMembersResponder()
	m.registerAddMembersByIdResponder()
	m.registerRemoveMembersByIdResponder()
	m.registerDeleteUpdatableAssetGroupResponder()
	m.registerListDevicesResponder()
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
				mockState.groups[id] = map[string]any{
					"id":      id,
					"members": map[string]bool{},
				}
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
			groupID := parts[len(parts)-1]

			mockState.Lock()
			group, exists := mockState.groups[groupID]
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
				group = responseObj
			}

			resp, err := httpmock.NewJsonResponse(200, group)
			if err != nil {
				return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
			}
			return resp, nil
		})
}

func (m *WindowsUpdateUpdatableAssetGroupMock) registerGetGroupMembersResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatableAssets/[0-9a-fA-F-]+/microsoft\.graph\.windowsUpdates\.updatableAssetGroup/members`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			// path: .../updatableAssets/{groupId}/microsoft.graph.windowsUpdates.updatableAssetGroup/members
			groupID := parts[len(parts)-3]

			mockState.Lock()
			group, exists := mockState.groups[groupID]
			mockState.Unlock()

			if exists && group != nil {
				members, _ := group["members"].(map[string]bool)
				memberList := make([]map[string]any, 0, len(members))
				for id := range members {
					memberList = append(memberList, map[string]any{
						"@odata.type": "#microsoft.graph.windowsUpdates.azureADDevice",
						"id":          id,
					})
				}
				resp, err := httpmock.NewJsonResponse(200, map[string]any{
					"@odata.context": "https://graph.microsoft.com/beta/$metadata#admin/windows/updates/updatableAssets/members",
					"value":          memberList,
				})
				if err != nil {
					return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
				}
				return resp, nil
			}

			// Fall back to fixture when no in-memory state exists (e.g. import)
			_, filename, _, _ := runtime.Caller(0)
			sourceDir := filepath.Dir(filename)
			responsesPath := filepath.Join(sourceDir, "..", "tests", "responses", constants.TfOperationRead, "get_members_empty.json")

			jsonData, err := os.ReadFile(responsesPath)
			if err != nil {
				return httpmock.NewStringResponse(200, `{"value":[]}`), nil
			}

			var responseObj map[string]any
			if err := json.Unmarshal(jsonData, &responseObj); err != nil {
				return httpmock.NewStringResponse(200, `{"value":[]}`), nil
			}

			resp, err := httpmock.NewJsonResponse(200, responseObj)
			if err != nil {
				return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
			}
			return resp, nil
		})
}

func (m *WindowsUpdateUpdatableAssetGroupMock) registerAddMembersByIdResponder() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatableAssets/[0-9a-fA-F-]+/microsoft\.graph\.windowsUpdates\.addMembersById$`,
		func(req *http.Request) (*http.Response, error) {
			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewJsonResponse(400, map[string]any{
					"error": map[string]any{"code": "BadRequest", "message": "Invalid request body"},
				})
			}

			parts := strings.Split(req.URL.Path, "/")
			groupID := parts[len(parts)-2]
			ids, _ := body["ids"].([]any)

			mockState.Lock()
			if _, ok := mockState.groups[groupID]; !ok {
				mockState.groups[groupID] = map[string]any{
					"id":      groupID,
					"members": map[string]bool{},
				}
			}
			members, _ := mockState.groups[groupID]["members"].(map[string]bool)
			for _, idAny := range ids {
				if id, ok := idAny.(string); ok {
					members[id] = true
				}
			}
			mockState.groups[groupID]["members"] = members
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})
}

func (m *WindowsUpdateUpdatableAssetGroupMock) registerRemoveMembersByIdResponder() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatableAssets/[0-9a-fA-F-]+/microsoft\.graph\.windowsUpdates\.removeMembersById$`,
		func(req *http.Request) (*http.Response, error) {
			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewJsonResponse(400, map[string]any{
					"error": map[string]any{"code": "BadRequest", "message": "Invalid request body"},
				})
			}

			parts := strings.Split(req.URL.Path, "/")
			groupID := parts[len(parts)-2]
			ids, _ := body["ids"].([]any)

			mockState.Lock()
			if group, ok := mockState.groups[groupID]; ok {
				members, _ := group["members"].(map[string]bool)
				for _, idAny := range ids {
					if id, ok := idAny.(string); ok {
						delete(members, id)
					}
				}
				mockState.groups[groupID]["members"] = members
			}
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})
}

func (m *WindowsUpdateUpdatableAssetGroupMock) registerDeleteUpdatableAssetGroupResponder() {
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatableAssets/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			groupID := parts[len(parts)-1]

			mockState.Lock()
			delete(mockState.groups, groupID)
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})
}

func (m *WindowsUpdateUpdatableAssetGroupMock) registerListDevicesResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/devices\?`,
		func(req *http.Request) (*http.Response, error) {
			devices := []map[string]any{
				{"id": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", "deviceId": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"},
				{"id": "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb", "deviceId": "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"},
			}
			resp, err := httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#devices",
				"value":          devices,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
			}
			return resp, nil
		})
}

func (m *WindowsUpdateUpdatableAssetGroupMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.groups = make(map[string]map[string]any)
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
	mockState.groups = make(map[string]map[string]any)
	mockState.Unlock()
}
