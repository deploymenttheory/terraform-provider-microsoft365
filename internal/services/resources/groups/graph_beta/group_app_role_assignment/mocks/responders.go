package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	groups             map[string]map[string]any
	appRoleAssignments map[string]map[string]any // groupId -> assignmentId -> assignmentData
}

func init() {
	// Initialize mockState
	mockState.groups = make(map[string]map[string]any)
	mockState.appRoleAssignments = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// GroupAppRoleAssignmentMock provides mock responses for group app role assignment operations
type GroupAppRoleAssignmentMock struct{}

// RegisterMocks registers HTTP mock responses for group app role assignment operations
func (m *GroupAppRoleAssignmentMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.groups = make(map[string]map[string]any)
	mockState.appRoleAssignments = make(map[string]map[string]any)
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
	mockState.appRoleAssignments[baseGroupId] = make(map[string]any)
	mockState.Unlock()

	// Register GET for group data
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/[^/]+$`,
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

	// Register GET for app role assignments
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/[^/]+/appRoleAssignments$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			groupId := urlParts[len(urlParts)-2]

			mockState.Lock()
			assignmentsMap, exists := mockState.appRoleAssignments[groupId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Group not found"}}`), nil
			}

			// Convert assignments map to array for response
			assignments := make([]map[string]any, 0, len(assignmentsMap))
			for _, assignmentData := range assignmentsMap {
				if ad, ok := assignmentData.(map[string]any); ok {
					assignments = append(assignments, ad)
				}
			}

			response := map[string]any{
				"@odata.context": fmt.Sprintf("https://graph.microsoft.com/beta/$metadata#groups('%s')/appRoleAssignments", groupId),
				"value":          assignments,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for specific app role assignment
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/[^/]+/appRoleAssignments/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			groupId := urlParts[len(urlParts)-3]
			assignmentId := urlParts[len(urlParts)-1]

			mockState.Lock()
			assignmentsMap, groupExists := mockState.appRoleAssignments[groupId]
			mockState.Unlock()

			if !groupExists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Group not found"}}`), nil
			}

			mockState.Lock()
			assignmentData, assignmentExists := assignmentsMap[assignmentId]
			mockState.Unlock()

			if !assignmentExists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"App role assignment not found"}}`), nil
			}

			return httpmock.NewJsonResponse(200, assignmentData)
		})

	// Register POST for creating app role assignments
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/groups/[^/]+/appRoleAssignments$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			groupId := urlParts[len(urlParts)-2]

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

			// Extract required fields
			principalId, hasPrincipalId := requestBody["principalId"].(string)
			resourceId, hasResourceId := requestBody["resourceId"].(string)
			appRoleId, hasAppRoleId := requestBody["appRoleId"].(string)

			if !hasPrincipalId || !hasResourceId || !hasAppRoleId {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Missing required fields: principalId, resourceId, or appRoleId"}}`), nil
			}

			// Generate a unique assignment ID
			assignmentId := fmt.Sprintf("%s-%s", principalId, appRoleId)

			// Create assignment data
			assignmentData := map[string]any{
				"id":                   assignmentId,
				"principalId":          principalId,
				"resourceId":           resourceId,
				"appRoleId":            appRoleId,
				"principalDisplayName": fmt.Sprintf("Group %s", principalId[0:8]),
				"resourceDisplayName":  fmt.Sprintf("ServicePrincipal %s", resourceId[0:8]),
				"principalType":        "Group",
				"creationTimestamp":    time.Now().UTC().Format(time.RFC3339),
			}

			// Handle specific test IDs to return expected display names
			if principalId == "00000000-0000-0000-0000-000000000002" {
				assignmentData["principalDisplayName"] = "Minimal Group"
			} else if principalId == "00000000-0000-0000-0000-000000000003" {
				assignmentData["principalDisplayName"] = "Maximal Group"
			}

			if resourceId == "00000000-0000-0000-0000-000000000010" {
				assignmentData["resourceDisplayName"] = "Microsoft Graph"
			} else if resourceId == "00000000-0000-0000-0000-000000000011" {
				assignmentData["resourceDisplayName"] = "SharePoint Online"
			}

			// Add assignment to group
			mockState.Lock()
			if mockState.appRoleAssignments[groupId] == nil {
				mockState.appRoleAssignments[groupId] = make(map[string]any)
			}
			mockState.appRoleAssignments[groupId][assignmentId] = assignmentData
			mockState.Unlock()

			// Return created assignment
			return httpmock.NewJsonResponse(201, assignmentData)
		})

	// Register DELETE for removing app role assignments
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph\.microsoft\.com/beta/groups/[^/]+/appRoleAssignments/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			groupId := urlParts[len(urlParts)-3]
			assignmentId := urlParts[len(urlParts)-1]

			mockState.Lock()
			assignmentsMap, groupExists := mockState.appRoleAssignments[groupId]
			mockState.Unlock()

			if !groupExists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Group not found"}}`), nil
			}

			mockState.Lock()
			_, assignmentExists := assignmentsMap[assignmentId]
			if assignmentExists {
				delete(mockState.appRoleAssignments[groupId], assignmentId)
			}
			mockState.Unlock()

			// Return success (204 No Content is typical for this operation)
			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register specific group and assignment IDs for testing
	registerSpecificGroupMocks()
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *GroupAppRoleAssignmentMock) RegisterErrorMocks() {
	// Register error response for app role assignment
	errorGroupId := "99999999-9999-9999-9999-999999999999"
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/groups/"+errorGroupId+"/appRoleAssignments",
		factories.ErrorResponse(400, "BadRequest", "Error assigning app role to group"))

	// Register GET for error group to ensure it exists but will fail on assignment
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
	// Minimal group with no assignments (Security group)
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
	mockState.appRoleAssignments[minimalGroupId] = make(map[string]any)
	mockState.Unlock()

	// Maximal group with assignments (Security group)
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
	mockState.appRoleAssignments[maximalGroupId] = make(map[string]any)
	mockState.Unlock()

	// Service Principals for testing
	servicePrincipalId1 := "00000000-0000-0000-0000-000000000010"
	servicePrincipalId2 := "00000000-0000-0000-0000-000000000011"

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

	// Register specific GET for service principals (these would normally be separate resources)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/servicePrincipals/"+servicePrincipalId1,
		func(req *http.Request) (*http.Response, error) {
			spData := map[string]any{
				"id":          servicePrincipalId1,
				"displayName": "Microsoft Graph",
				"appId":       "00000003-0000-0000-c000-000000000000",
			}
			return httpmock.NewJsonResponse(200, spData)
		})

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/servicePrincipals/"+servicePrincipalId2,
		func(req *http.Request) (*http.Response, error) {
			spData := map[string]any{
				"id":          servicePrincipalId2,
				"displayName": "SharePoint Online",
				"appId":       "00000003-0000-0ff1-ce00-000000000000",
			}
			return httpmock.NewJsonResponse(200, spData)
		})
}

// CleanupMockState cleans up the mock state
func (m *GroupAppRoleAssignmentMock) CleanupMockState() {
	mockState.Lock()
	mockState.groups = make(map[string]map[string]any)
	mockState.appRoleAssignments = make(map[string]map[string]any)
	mockState.Unlock()
}

// SetupImportTest sets up the mock state for an import test
func (m *GroupAppRoleAssignmentMock) SetupImportTest(groupId, assignmentId string) {
	assignmentData := map[string]any{
		"id":                   assignmentId,
		"principalId":          groupId,
		"resourceId":           "00000000-0000-0000-0000-000000000010",
		"appRoleId":            "00000000-0000-0000-0000-000000000000",
		"principalDisplayName": "Minimal Group",
		"resourceDisplayName":  "Microsoft Graph",
		"principalType":        "Group",
		"creationTimestamp":    time.Now().UTC().Format(time.RFC3339),
	}

	// Add the assignment to the group in the mock state
	mockState.Lock()
	if mockState.appRoleAssignments[groupId] == nil {
		mockState.appRoleAssignments[groupId] = make(map[string]any)
	}
	mockState.appRoleAssignments[groupId][assignmentId] = assignmentData
	mockState.Unlock()
}
