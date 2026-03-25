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
	auMembers map[string]map[string]string // auID -> memberID -> odataType
}

// memberTypeMap maps known test member IDs to their @odata.type
var memberTypeMap = map[string]string{
	"22222222-2222-2222-2222-222222222222": "#microsoft.graph.user",
	"33333333-3333-3333-3333-333333333333": "#microsoft.graph.user",
	"44444444-4444-4444-4444-444444444444": "#microsoft.graph.group",
	"55555555-5555-5555-5555-555555555555": "#microsoft.graph.user",
}

func init() {
	mockState.auMembers = make(map[string]map[string]string)
}

type AdministrativeUnitMembershipMock struct{}

func (m *AdministrativeUnitMembershipMock) RegisterMocks() {
	mockState.Lock()
	mockState.auMembers = make(map[string]map[string]string)
	mockState.Unlock()

	m.registerValidationResponders()
	m.registerAddMemberResponder()
	m.registerGetMembersResponder()
	m.registerRemoveMemberResponder()
}

func (m *AdministrativeUnitMembershipMock) registerValidationResponders() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/users$`,
		func(req *http.Request) (*http.Response, error) {
			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#users",
				"value": []map[string]any{
					{"@odata.type": "#microsoft.graph.user", "id": "22222222-2222-2222-2222-222222222222", "userPrincipalName": "testuser1@example.com", "displayName": "Test User 1"},
					{"@odata.type": "#microsoft.graph.user", "id": "33333333-3333-3333-3333-333333333333", "userPrincipalName": "testuser2@example.com", "displayName": "Test User 2"},
					{"@odata.type": "#microsoft.graph.user", "id": "55555555-5555-5555-5555-555555555555", "userPrincipalName": "testuser3@example.com", "displayName": "Test User 3"},
				},
			}
			return httpmock.NewJsonResponse(200, response)
		})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups$`,
		func(req *http.Request) (*http.Response, error) {
			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#groups",
				"value": []map[string]any{
					{"@odata.type": "#microsoft.graph.group", "id": "44444444-4444-4444-4444-444444444444", "displayName": "Test Group", "mailEnabled": false, "securityEnabled": true},
				},
			}
			return httpmock.NewJsonResponse(200, response)
		})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/devices$`,
		func(req *http.Request) (*http.Response, error) {
			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#devices",
				"value":          []map[string]any{},
			}
			return httpmock.NewJsonResponse(200, response)
		})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directoryObjects/([^/]+)$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			objectID := parts[len(parts)-1]

			odataType, exists := memberTypeMap[objectID]
			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"Directory object not found"}}`), nil
			}

			response := map[string]any{
				"@odata.type": odataType,
				"id":          objectID,
			}
			return httpmock.NewJsonResponse(200, response)
		})
}

func (m *AdministrativeUnitMembershipMock) registerAddMemberResponder() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/directory/administrativeUnits/([^/]+)/members/\$ref$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			var auID string
			for i, part := range parts {
				if part == "administrativeUnits" && i+1 < len(parts) {
					auID = parts[i+1]
					break
				}
			}

			var requestBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Extract member ID from @odata.id: "https://graph.microsoft.com/beta/directoryObjects/{memberID}"
			odataID, ok := requestBody["@odata.id"].(string)
			if !ok || odataID == "" {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"@odata.id is required"}}`), nil
			}
			urlParts := strings.Split(odataID, "/")
			memberID := urlParts[len(urlParts)-1]

			odataType := memberTypeMap[memberID]
			if odataType == "" {
				odataType = "#microsoft.graph.directoryObject"
			}

			mockState.Lock()
			if mockState.auMembers[auID] == nil {
				mockState.auMembers[auID] = make(map[string]string)
			}
			mockState.auMembers[auID][memberID] = odataType
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})
}

func (m *AdministrativeUnitMembershipMock) registerGetMembersResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directory/administrativeUnits/([^/]+)/members$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			var auID string
			for i, part := range parts {
				if part == "administrativeUnits" && i+1 < len(parts) {
					auID = parts[i+1]
					break
				}
			}

			mockState.Lock()
			members := mockState.auMembers[auID]
			memberList := make([]map[string]any, 0, len(members))
			for memberID, odataType := range members {
				memberList = append(memberList, map[string]any{
					"@odata.type": odataType,
					"id":          memberID,
				})
			}
			mockState.Unlock()

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#directoryObjects",
				"value":          memberList,
			}
			return httpmock.NewJsonResponse(200, response)
		})
}

func (m *AdministrativeUnitMembershipMock) registerRemoveMemberResponder() {
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/directory/administrativeUnits/([^/]+)/members/([^/]+)/\$ref$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			var auID, memberID string
			for i, part := range parts {
				if part == "administrativeUnits" && i+1 < len(parts) {
					auID = parts[i+1]
				}
				if part == "members" && i+1 < len(parts) {
					memberID = parts[i+1]
				}
			}

			mockState.Lock()
			if mockState.auMembers[auID] != nil {
				delete(mockState.auMembers[auID], memberID)
			}
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})
}

func (m *AdministrativeUnitMembershipMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/directory/administrativeUnits/([^/]+)/members/\$ref$`,
		httpmock.NewStringResponder(403, `{"error":{"code":"Forbidden","message":"Insufficient privileges"}}`))

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directory/administrativeUnits/([^/]+)/members$`,
		httpmock.NewStringResponder(404, `{"error":{"code":"NotFound","message":"Administrative unit not found"}}`))

	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/directory/administrativeUnits/([^/]+)/members/([^/]+)/\$ref$`,
		httpmock.NewStringResponder(403, `{"error":{"code":"Forbidden","message":"Insufficient privileges"}}`))
}
