package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	groups       map[string]map[string]any
	groupMembers map[string]map[string]any // groupId -> memberId -> memberData
}

func init() {
	// Initialize mockState
	mockState.groups = make(map[string]map[string]any)
	mockState.groupMembers = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// GroupMemberAssignmentMock provides mock responses for group member assignment operations
type GroupMemberAssignmentMock struct{}

// RegisterMocks registers HTTP mock responses for group member assignment operations
func (m *GroupMemberAssignmentMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.groups = make(map[string]map[string]any)
	mockState.groupMembers = make(map[string]map[string]any)
	mockState.Unlock()

	// Initialize base group data
	baseGroupId := "00000000-0000-0000-0000-000000000001"
	baseGroupData := map[string]any{
		"id":              baseGroupId,
		"displayName":     "Base Group",
		"description":     "Base test group",
		"groupTypes":      []string{"Unified"},
		"mailEnabled":     true,
		"securityEnabled": false,
	}

	mockState.Lock()
	mockState.groups[baseGroupId] = baseGroupData
	mockState.groupMembers[baseGroupId] = make(map[string]any)
	mockState.Unlock()

	// Register GET for group data
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/groups/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			groupId := urlParts[len(urlParts)-1]

			mockState.Lock()
			groupData, exists := mockState.groups[groupId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Group not found"}}`), nil
			}

			return httpmock.NewJsonResponse(200, groupData)
		})

	// Register GET for group members
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/groups/[^/]+/members$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			groupId := urlParts[len(urlParts)-2]

			mockState.Lock()
			membersMap, exists := mockState.groupMembers[groupId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Group not found"}}`), nil
			}

			// Convert members map to array for response
			members := make([]map[string]any, 0, len(membersMap))
			for _, memberData := range membersMap {
				if md, ok := memberData.(map[string]any); ok {
					members = append(members, md)
				}
			}

			response := map[string]any{
				"@odata.context": fmt.Sprintf("https://graph.microsoft.com/beta/$metadata#directoryObjects"),
				"value":          members,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for specific member
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/groups/[^/]+/members/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			groupId := urlParts[len(urlParts)-3]
			memberId := urlParts[len(urlParts)-1]

			mockState.Lock()
			membersMap, groupExists := mockState.groupMembers[groupId]
			mockState.Unlock()

			if !groupExists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Group not found"}}`), nil
			}

			mockState.Lock()
			memberData, memberExists := membersMap[memberId]
			mockState.Unlock()

			if !memberExists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Member not found"}}`), nil
			}

			return httpmock.NewJsonResponse(200, memberData)
		})

	// Register POST for adding members
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/groups/[^/]+/members/\$ref$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			groupId := urlParts[len(urlParts)-3]

			mockState.Lock()
			_, groupExists := mockState.groups[groupId]
			mockState.Unlock()

			if !groupExists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Group not found"}}`), nil
			}

			// Parse request body
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Extract member ID from @odata.id
			odataId, ok := requestBody["@odata.id"].(string)
			if !ok {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Missing @odata.id"}}`), nil
			}

			// Extract member ID from URL
			odataParts := strings.Split(odataId, "/")
			memberType := odataParts[len(odataParts)-2] // users, groups, etc.
			memberId := odataParts[len(odataParts)-1]

			// Create member data based on type
			var memberData map[string]any

			// Check if this is a specific test member ID
			if memberId == "00000000-0000-0000-0000-000000000004" {
				// Minimal user
				memberData = map[string]any{
					"id":                memberId,
					"@odata.type":       "#microsoft.graph.user",
					"displayName":       "Minimal User",
					"userPrincipalName": "minimal.user@contoso.com",
				}
			} else if memberId == "00000000-0000-0000-0000-000000000005" {
				// Maximal group
				memberData = map[string]any{
					"id":              memberId,
					"@odata.type":     "#microsoft.graph.group",
					"displayName":     "Maximal Group",
					"description":     "A nested group for testing",
					"mailEnabled":     false,
					"securityEnabled": true,
				}
			} else {
				// Generic member based on type
				switch memberType {
				case "users":
					memberData = map[string]any{
						"id":                memberId,
						"@odata.type":       "#microsoft.graph.user",
						"displayName":       fmt.Sprintf("User %s", memberId[0:8]),
						"userPrincipalName": fmt.Sprintf("user-%s@contoso.com", memberId[0:8]),
					}
				case "directoryObjects":
					// For directoryObjects, we need to determine the actual type
					// For our tests, assume it's a user
					memberData = map[string]any{
						"id":                memberId,
						"@odata.type":       "#microsoft.graph.user",
						"displayName":       fmt.Sprintf("User %s", memberId[0:8]),
						"userPrincipalName": fmt.Sprintf("user-%s@contoso.com", memberId[0:8]),
					}
				case "groups":
					memberData = map[string]any{
						"id":              memberId,
						"@odata.type":     "#microsoft.graph.group",
						"displayName":     fmt.Sprintf("Group %s", memberId[0:8]),
						"description":     "A nested group",
						"mailEnabled":     false,
						"securityEnabled": true,
					}
				case "devices":
					memberData = map[string]any{
						"id":          memberId,
						"@odata.type": "#microsoft.graph.device",
						"displayName": fmt.Sprintf("Device %s", memberId[0:8]),
					}
				case "servicePrincipals":
					memberData = map[string]any{
						"id":          memberId,
						"@odata.type": "#microsoft.graph.servicePrincipal",
						"displayName": fmt.Sprintf("ServicePrincipal %s", memberId[0:8]),
					}
				case "contacts":
					memberData = map[string]any{
						"id":          memberId,
						"@odata.type": "#microsoft.graph.orgContact",
						"displayName": fmt.Sprintf("Contact %s", memberId[0:8]),
					}
				default:
					return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Unsupported member type"}}`), nil
				}
			}

			// Add member to group
			mockState.Lock()
			if mockState.groupMembers[groupId] == nil {
				mockState.groupMembers[groupId] = make(map[string]any)
			}
			mockState.groupMembers[groupId][memberId] = memberData
			mockState.Unlock()

			// Return success (204 No Content is typical for this operation)
			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register DELETE for removing members
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph.microsoft.com/beta/groups/[^/]+/members/[^/]+/\$ref$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			groupId := urlParts[len(urlParts)-4]
			memberId := urlParts[len(urlParts)-2]

			mockState.Lock()
			membersMap, groupExists := mockState.groupMembers[groupId]
			mockState.Unlock()

			if !groupExists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Group not found"}}`), nil
			}

			mockState.Lock()
			_, memberExists := membersMap[memberId]
			if memberExists {
				delete(mockState.groupMembers[groupId], memberId)
			}
			mockState.Unlock()

			// Return success (204 No Content is typical for this operation)
			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register specific group and member IDs for testing
	registerSpecificGroupMocks()
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *GroupMemberAssignmentMock) RegisterErrorMocks() {
	// Register error response for member assignment
	errorGroupId := "99999999-9999-9999-9999-999999999999"
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/groups/"+errorGroupId+"/members/$ref",
		factories.ErrorResponse(400, "BadRequest", "Error adding member to group"))

	// Register GET for error group to ensure it exists but will fail on member assignment
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/groups/"+errorGroupId,
		func(req *http.Request) (*http.Response, error) {
			groupData := map[string]any{
				"id":              errorGroupId,
				"displayName":     "Error Group",
				"description":     "Group that generates errors",
				"mailEnabled":     false,
				"securityEnabled": true,
			}
			return httpmock.NewJsonResponse(200, groupData)
		})

	// Register error response for group not found
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/groups/not-found-group",
		factories.ErrorResponse(404, "ResourceNotFound", "Group not found"))
}

// registerSpecificGroupMocks registers mocks for specific test scenarios
func registerSpecificGroupMocks() {
	// Minimal group with no members (Security group)
	minimalGroupId := "00000000-0000-0000-0000-000000000002"
	minimalGroupData := map[string]any{
		"id":              minimalGroupId,
		"displayName":     "Minimal Group",
		"description":     "Minimal test group",
		"mailEnabled":     false,
		"securityEnabled": true,
	}

	mockState.Lock()
	mockState.groups[minimalGroupId] = minimalGroupData
	mockState.groupMembers[minimalGroupId] = make(map[string]any)
	mockState.Unlock()

	// Maximal group with multiple members (Security group)
	maximalGroupId := "00000000-0000-0000-0000-000000000003"
	maximalGroupData := map[string]any{
		"id":              maximalGroupId,
		"displayName":     "Maximal Group",
		"description":     "Maximal test group",
		"mailEnabled":     false,
		"securityEnabled": true,
	}

	mockState.Lock()
	mockState.groups[maximalGroupId] = maximalGroupData
	mockState.groupMembers[maximalGroupId] = make(map[string]any)
	mockState.Unlock()

	// Add test user for minimal configuration
	minimalUserId := "00000000-0000-0000-0000-000000000004"
	minimalUserData := map[string]any{
		"id":                minimalUserId,
		"@odata.type":       "#microsoft.graph.user",
		"displayName":       "Minimal User",
		"userPrincipalName": "minimal.user@contoso.com",
	}

	// Add test group for maximal configuration (Security group)
	maximalMemberId := "00000000-0000-0000-0000-000000000005"
	maximalMemberData := map[string]any{
		"id":              maximalMemberId,
		"@odata.type":     "#microsoft.graph.group",
		"displayName":     "Maximal Group",
		"description":     "A nested group for testing",
		"mailEnabled":     false,
		"securityEnabled": true,
	}

	// Register specific GET for these groups
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/groups/"+minimalGroupId,
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			groupData := mockState.groups[minimalGroupId]
			mockState.Unlock()
			return httpmock.NewJsonResponse(200, groupData)
		})

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/groups/"+maximalGroupId,
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			groupData := mockState.groups[maximalGroupId]
			mockState.Unlock()
			return httpmock.NewJsonResponse(200, groupData)
		})

	// Register specific GET for members
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/users/"+minimalUserId,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, minimalUserData)
		})

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/groups/"+maximalMemberId,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, maximalMemberData)
		})

	// Register specific GET for validating security group member
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/groups/"+maximalMemberId,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, maximalMemberData)
		})
}

// SetupImportTest sets up the mock state for an import test
func (m *GroupMemberAssignmentMock) SetupImportTest(groupId, memberId string) {
	minimalUserData := map[string]any{
		"id":                memberId,
		"@odata.type":       "#microsoft.graph.user",
		"displayName":       "Minimal User",
		"userPrincipalName": "minimal.user@contoso.com",
	}

	// Add the member to the group in the mock state
	mockState.Lock()
	if mockState.groupMembers[groupId] == nil {
		mockState.groupMembers[groupId] = make(map[string]any)
	}
	mockState.groupMembers[groupId][memberId] = minimalUserData
	mockState.Unlock()
}
