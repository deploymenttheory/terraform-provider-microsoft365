package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	// auID -> scopedRoleMembershipID -> assignment data
	assignments map[string]map[string]scopedRoleAssignment
}

type scopedRoleAssignment struct {
	ID                 string
	AdministrativeUnit string
	RoleID             string
	MemberID           string
	DisplayName        string
	UserPrincipalName  string
}

// memberDisplayNames maps known test member IDs to display names and UPNs for mock responses
var memberDisplayNames = map[string][2]string{
	"22222222-2222-2222-2222-222222222222": {"Test User 1", "testuser1@example.com"},
	"33333333-3333-3333-3333-333333333333": {"Test User 2", "testuser2@example.com"},
	"55555555-5555-5555-5555-555555555555": {"Test User 3", "testuser3@example.com"},
}

func init() {
	mockState.assignments = make(map[string]map[string]scopedRoleAssignment)
}

type AdministrativeUnitRoleAssignmentMock struct{}

func (m *AdministrativeUnitRoleAssignmentMock) RegisterMocks() {
	mockState.Lock()
	mockState.assignments = make(map[string]map[string]scopedRoleAssignment)
	mockState.Unlock()

	m.registerDirectoryObjectValidationResponder()
	m.registerPostScopedRoleMemberResponder()
	m.registerGetScopedRoleMemberResponder()
	m.registerDeleteScopedRoleMemberResponder()
}

func (m *AdministrativeUnitRoleAssignmentMock) registerDirectoryObjectValidationResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directoryObjects/([^/]+)$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			objectID := parts[len(parts)-1]

			if _, exists := memberDisplayNames[objectID]; !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"Directory object not found"}}`), nil
			}

			response := map[string]any{
				"@odata.type": "#microsoft.graph.user",
				"id":          objectID,
			}
			return httpmock.NewJsonResponse(200, response)
		})
}

func (m *AdministrativeUnitRoleAssignmentMock) registerPostScopedRoleMemberResponder() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/administrativeUnits/([^/]+)/scopedRoleMembers$`,
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

			roleID, _ := requestBody["roleId"].(string)

			var memberID string
			if roleMemberInfo, ok := requestBody["roleMemberInfo"].(map[string]any); ok {
				memberID, _ = roleMemberInfo["id"].(string)
			}

			if roleID == "" || memberID == "" {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"roleId and roleMemberInfo.id are required"}}`), nil
			}

			membershipID := uuid.New().String()
			displayName := ""
			upn := ""
			if info, exists := memberDisplayNames[memberID]; exists {
				displayName = info[0]
				upn = info[1]
			}

			assignment := scopedRoleAssignment{
				ID:                 membershipID,
				AdministrativeUnit: auID,
				RoleID:             roleID,
				MemberID:           memberID,
				DisplayName:        displayName,
				UserPrincipalName:  upn,
			}

			mockState.Lock()
			if mockState.assignments[auID] == nil {
				mockState.assignments[auID] = make(map[string]scopedRoleAssignment)
			}
			mockState.assignments[auID][membershipID] = assignment
			mockState.Unlock()

			response := map[string]any{
				"@odata.context":     "https://graph.microsoft.com/beta/$metadata#scopedRoleMemberships/$entity",
				"id":                 membershipID,
				"administrativeUnitId": auID,
				"roleId":             roleID,
				"roleMemberInfo": map[string]any{
					"id":                memberID,
					"displayName":       displayName,
					"userPrincipalName": upn,
				},
			}
			return httpmock.NewJsonResponse(201, response)
		})
}

func (m *AdministrativeUnitRoleAssignmentMock) registerGetScopedRoleMemberResponder() {
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/administrativeUnits/([^/]+)/scopedRoleMembers/([^/]+)$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			var auID, membershipID string
			for i, part := range parts {
				if part == "administrativeUnits" && i+1 < len(parts) {
					auID = parts[i+1]
				}
				if part == "scopedRoleMembers" && i+1 < len(parts) {
					membershipID = parts[i+1]
				}
			}

			mockState.Lock()
			auAssignments, auExists := mockState.assignments[auID]
			var assignment scopedRoleAssignment
			var exists bool
			if auExists {
				assignment, exists = auAssignments[membershipID]
			}
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, fmt.Sprintf(
					`{"error":{"code":"NotFound","message":"Scoped role membership %s not found"}}`, membershipID,
				)), nil
			}

			response := map[string]any{
				"@odata.context":     "https://graph.microsoft.com/beta/$metadata#scopedRoleMemberships/$entity",
				"id":                 assignment.ID,
				"administrativeUnitId": assignment.AdministrativeUnit,
				"roleId":             assignment.RoleID,
				"roleMemberInfo": map[string]any{
					"id":                assignment.MemberID,
					"displayName":       assignment.DisplayName,
					"userPrincipalName": assignment.UserPrincipalName,
				},
			}
			return httpmock.NewJsonResponse(200, response)
		})
}

func (m *AdministrativeUnitRoleAssignmentMock) registerDeleteScopedRoleMemberResponder() {
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/administrativeUnits/([^/]+)/scopedRoleMembers/([^/]+)$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			var auID, membershipID string
			for i, part := range parts {
				if part == "administrativeUnits" && i+1 < len(parts) {
					auID = parts[i+1]
				}
				if part == "scopedRoleMembers" && i+1 < len(parts) {
					membershipID = parts[i+1]
				}
			}

			mockState.Lock()
			if mockState.assignments[auID] != nil {
				delete(mockState.assignments[auID], membershipID)
			}
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})
}

func (m *AdministrativeUnitRoleAssignmentMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/administrativeUnits/([^/]+)/scopedRoleMembers$`,
		httpmock.NewStringResponder(403, `{"error":{"code":"Forbidden","message":"Insufficient privileges to complete the operation"}}`))

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/administrativeUnits/([^/]+)/scopedRoleMembers/([^/]+)$`,
		httpmock.NewStringResponder(404, `{"error":{"code":"NotFound","message":"Scoped role membership not found"}}`))

	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/administrativeUnits/([^/]+)/scopedRoleMembers/([^/]+)$`,
		httpmock.NewStringResponder(403, `{"error":{"code":"Forbidden","message":"Insufficient privileges to complete the operation"}}`))
}
