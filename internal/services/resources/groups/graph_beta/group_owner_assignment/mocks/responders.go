package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	groups map[string]map[string]interface{} // groupId -> groupData
	owners map[string]map[string]string      // groupId -> ownerId -> ownerType
}

func init() {
	mockState.groups = make(map[string]map[string]interface{})
	mockState.owners = make(map[string]map[string]string)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

type GroupOwnerAssignmentMock struct{}

func (m *GroupOwnerAssignmentMock) RegisterMocks() {
	mockState.Lock()
	mockState.groups = make(map[string]map[string]interface{})
	mockState.owners = make(map[string]map[string]string)
	mockState.Unlock()

	// Register explicit test groups
	minimalGroupId := "00000000-0000-0000-0000-000000000002"
	maximalGroupId := "00000000-0000-0000-0000-000000000003"
	mockState.Lock()
	mockState.groups[minimalGroupId] = map[string]interface{}{
		"id":          minimalGroupId,
		"displayName": "Minimal Group",
		"groupTypes":  []string{"Unified"},
	}
	mockState.groups[maximalGroupId] = map[string]interface{}{
		"id":          maximalGroupId,
		"displayName": "Maximal Group",
		"groupTypes":  []string{"Unified"},
	}
	mockState.owners[minimalGroupId] = make(map[string]string)
	mockState.owners[maximalGroupId] = make(map[string]string)
	mockState.Unlock()

	// Specific GET for minimal group
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/groups/"+minimalGroupId,
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			groupData := mockState.groups[minimalGroupId]
			mockState.Unlock()
			return httpmock.NewJsonResponse(200, groupData)
		})
	// Specific GET for maximal group
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/groups/"+maximalGroupId,
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			groupData := mockState.groups[maximalGroupId]
			mockState.Unlock()
			return httpmock.NewJsonResponse(200, groupData)
		})

	// GET owners for a group
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/groups/([^/]+)/owners$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			groupId := urlParts[len(urlParts)-2]
			mockState.Lock()
			owners, exists := mockState.owners[groupId]
			mockState.Unlock()
			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Group not found"}}`), nil
			}
			var value []map[string]interface{}
			for ownerId, ownerType := range owners {
				value = append(value, map[string]interface{}{"id": ownerId, "@odata.type": ownerType})
			}
			return httpmock.NewJsonResponse(200, map[string]interface{}{"value": value})
		})

	// POST add owner (add or update for update)
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/groups/([^/]+)/owners/\$ref$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			groupId := urlParts[len(urlParts)-3]
			var ref struct {
				OdataId string `json:"@odata.id"`
			}
			_ = json.NewDecoder(req.Body).Decode(&ref)
			ownerId := ref.OdataId[strings.LastIndex(ref.OdataId, "/")+1:]
			ownerType := "#microsoft.graph.user"
			if strings.Contains(strings.ToLower(ref.OdataId), "serviceprincipal") || ownerId == "00000000-0000-0000-0000-000000000005" {
				ownerType = "#microsoft.graph.servicePrincipal"
			}
			mockState.Lock()
			if mockState.owners[groupId] == nil {
				mockState.owners[groupId] = make(map[string]string)
			}
			mockState.owners[groupId][ownerId] = ownerType
			mockState.Unlock()
			return httpmock.NewStringResponse(204, ""), nil
		})

	// DELETE remove owner
	httpmock.RegisterResponder("DELETE", `=~^https://graph.microsoft.com/beta/groups/([^/]+)/owners/([^/]+)/\$ref$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			groupId := urlParts[len(urlParts)-4]
			ownerId := urlParts[len(urlParts)-2]
			mockState.Lock()
			if owners, exists := mockState.owners[groupId]; exists {
				delete(owners, ownerId)
			}
			mockState.Unlock()
			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register specific GET for test user (minimal owner)
	minimalUserId := "00000000-0000-0000-0000-000000000004"
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/users/"+minimalUserId,
		func(req *http.Request) (*http.Response, error) {
			userData := map[string]interface{}{
				"id":                minimalUserId,
				"@odata.type":       "#microsoft.graph.user",
				"displayName":       "Minimal User",
				"userPrincipalName": "minimal.user@contoso.com",
			}
			return httpmock.NewJsonResponse(200, userData)
		})

	// Register specific GET for test service principal (maximal owner)
	maximalServicePrincipalId := "00000000-0000-0000-0000-000000000005"
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/servicePrincipals/"+maximalServicePrincipalId,
		func(req *http.Request) (*http.Response, error) {
			spData := map[string]interface{}{
				"id":          maximalServicePrincipalId,
				"@odata.type": "#microsoft.graph.servicePrincipal",
				"displayName": "Maximal Service Principal",
			}
			return httpmock.NewJsonResponse(200, spData)
		})
}

func (m *GroupOwnerAssignmentMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/groups/([^/]+)/owners/\$ref$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(400, `{"error":{"code":"RequestConflict","message":"Owner already assigned to group"}}`), nil
		})
}

// SetupImportTest sets up the mock state for an import test
func (m *GroupOwnerAssignmentMock) SetupImportTest(groupId, ownerId, ownerType string) {
	mockState.Lock()
	if mockState.owners[groupId] == nil {
		mockState.owners[groupId] = make(map[string]string)
	}
	mockState.owners[groupId][ownerId] = ownerType
	mockState.Unlock()
}
