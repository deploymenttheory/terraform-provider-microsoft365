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
	mocks.GlobalRegistry.Register("windows_update_updatable_asset_group_assignment", &WindowsUpdateUpdatableAssetGroupAssignmentMock{})
}

type WindowsUpdateUpdatableAssetGroupAssignmentMock struct{}

var _ mocks.MockRegistrar = (*WindowsUpdateUpdatableAssetGroupAssignmentMock)(nil)

func (m *WindowsUpdateUpdatableAssetGroupAssignmentMock) RegisterMocks() {
	mockState.Lock()
	mockState.groups = make(map[string]map[string]any)
	mockState.Unlock()

	m.registerListUpdatableAssetsResponder()
	m.registerListDevicesResponder()
	m.registerAddMembersByIdResponder()
	m.registerRemoveMembersByIdResponder()
	m.registerGetGroupMembersResponder()
}

func (m *WindowsUpdateUpdatableAssetGroupAssignmentMock) registerListUpdatableAssetsResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatableAssets\?`,
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			groupIDs := make([]map[string]any, 0, len(mockState.groups))
			for id := range mockState.groups {
				groupIDs = append(groupIDs, map[string]any{"id": id})
			}
			mockState.Unlock()

			// Always include the well-known test group ID so validation passes
			found := false
			for _, g := range groupIDs {
				if g["id"] == "d4e5f6a7-4567-8901-defa-d4e5f6a7b8c9" {
					found = true
					break
				}
			}
			if !found {
				groupIDs = append(groupIDs, map[string]any{"id": "d4e5f6a7-4567-8901-defa-d4e5f6a7b8c9"})
			}

			resp, err := httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#admin/windows/updates/updatableAssets",
				"value":          groupIDs,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
			}
			return resp, nil
		})
}

func (m *WindowsUpdateUpdatableAssetGroupAssignmentMock) registerListDevicesResponder() {
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

func (m *WindowsUpdateUpdatableAssetGroupAssignmentMock) registerAddMembersByIdResponder() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatableAssets/[0-9a-fA-F-]+/microsoft\.graph\.windowsUpdates\.addMembersById$`,
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

			parts := strings.Split(req.URL.Path, "/")
			// path: /beta/admin/windows/updates/updatableAssets/{groupId}/microsoft.graph.windowsUpdates.addMembersById
			// groupId is at index len(parts)-2
			groupId := parts[len(parts)-2]

			ids, _ := body["ids"].([]any)

			mockState.Lock()
			if _, ok := mockState.groups[groupId]; !ok {
				mockState.groups[groupId] = map[string]any{
					"members": map[string]bool{},
				}
			}
			members, _ := mockState.groups[groupId]["members"].(map[string]bool)
			for _, idAny := range ids {
				if id, ok := idAny.(string); ok {
					members[id] = true
				}
			}
			mockState.groups[groupId]["members"] = members
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})
}

func (m *WindowsUpdateUpdatableAssetGroupAssignmentMock) registerRemoveMembersByIdResponder() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatableAssets/[0-9a-fA-F-]+/microsoft\.graph\.windowsUpdates\.removeMembersById$`,
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

			parts := strings.Split(req.URL.Path, "/")
			groupId := parts[len(parts)-2]

			ids, _ := body["ids"].([]any)

			mockState.Lock()
			if group, ok := mockState.groups[groupId]; ok {
				members, _ := group["members"].(map[string]bool)
				for _, idAny := range ids {
					if id, ok := idAny.(string); ok {
						delete(members, id)
					}
				}
				mockState.groups[groupId]["members"] = members
			}
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})
}

func (m *WindowsUpdateUpdatableAssetGroupAssignmentMock) registerGetGroupMembersResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatableAssets/[0-9a-fA-F-]+/microsoft\.graph\.windowsUpdates\.updatableAssetGroup/members`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			// path: .../updatableAssets/{groupId}/microsoft.graph.windowsUpdates.updatableAssetGroup/members
			// [... "updatableAssets", groupId, "microsoft.graph...", "members"]
			// groupId is at index len(parts)-3
			groupId := parts[len(parts)-3]

			mockState.Lock()
			group, exists := mockState.groups[groupId]
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

			// Fall back to fixture file when no in-memory state is present
			_, filename, _, _ := runtime.Caller(0)
			sourceDir := filepath.Dir(filename)
			responsesPath := filepath.Join(sourceDir, "..", "tests", "responses", constants.TfOperationRead, "get_updatable_asset_group_assignment_members.json")

			jsonData, err := os.ReadFile(responsesPath)
			if err != nil {
				return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"Members not found"}}`), nil
			}

			var responseObj map[string]any
			if err := json.Unmarshal(jsonData, &responseObj); err != nil {
				return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse response: %s"}}`, err.Error())), nil
			}

			resp, err := httpmock.NewJsonResponse(200, responseObj)
			if err != nil {
				return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
			}
			return resp, nil
		})
}

func (m *WindowsUpdateUpdatableAssetGroupAssignmentMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.groups = make(map[string]map[string]any)
	mockState.Unlock()

	m.registerListUpdatableAssetsErrorResponder()
	m.registerListDevicesErrorResponder()
	m.registerAddMembersByIdErrorResponder()
}

func (m *WindowsUpdateUpdatableAssetGroupAssignmentMock) registerListUpdatableAssetsErrorResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatableAssets\?`,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#admin/windows/updates/updatableAssets",
				"value":          []map[string]any{{"id": "d4e5f6a7-4567-8901-defa-d4e5f6a7b8c9"}},
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
			}
			return resp, nil
		})
}

func (m *WindowsUpdateUpdatableAssetGroupAssignmentMock) registerListDevicesErrorResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/devices\?`,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#devices",
				"value":          []map[string]any{{"id": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", "deviceId": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"}},
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
			}
			return resp, nil
		})
}

func (m *WindowsUpdateUpdatableAssetGroupAssignmentMock) registerAddMembersByIdErrorResponder() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/admin/windows/updates/updatableAssets/[0-9a-fA-F-]+/microsoft\.graph\.windowsUpdates\.addMembersById$`,
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

func (m *WindowsUpdateUpdatableAssetGroupAssignmentMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.groups = make(map[string]map[string]any)
}
