package mocks

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	users        []map[string]any
	devices      []map[string]any
	groupMembers map[string][]map[string]any // keyed by group ID
}

func init() {
	mockState.users = []map[string]any{}
	mockState.devices = []map[string]any{}
	mockState.groupMembers = make(map[string][]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("guid_list_sharder", &GuidListSharderMock{})
}

type GuidListSharderMock struct{}

var _ mocks.MockRegistrar = (*GuidListSharderMock)(nil)

func (m *GuidListSharderMock) RegisterMocks() {
	mockState.Lock()
	// Generate mock users (30 users for testing)
	mockState.users = make([]map[string]any, 30)
	for i := 0; i < 30; i++ {
		mockState.users[i] = map[string]any{
			"id":                fmt.Sprintf("user-%02d-0000-0000-0000-000000000000", i),
			"displayName":       fmt.Sprintf("Test User %d", i),
			"userPrincipalName": fmt.Sprintf("testuser%d@contoso.com", i),
			"accountEnabled":    true,
		}
	}

	// Generate mock devices (24 devices for testing)
	mockState.devices = make([]map[string]any, 24)
	for i := 0; i < 24; i++ {
		mockState.devices[i] = map[string]any{
			"id":              fmt.Sprintf("device-%02d-0000-0000-0000-000000000000", i),
			"displayName":     fmt.Sprintf("Test Device %d", i),
			"operatingSystem": "Windows",
		}
	}

	// Generate mock group members for test groups (20 members each)
	testGroupIds := []string{
		"group-test-0000-0000-0000-000000000000",
		"12345678-1234-1234-1234-123456789abc", // Used in TF test files
	}
	for _, groupId := range testGroupIds {
		mockState.groupMembers[groupId] = make([]map[string]any, 20)
		for i := 0; i < 20; i++ {
			mockState.groupMembers[groupId][i] = map[string]any{
				"id":          fmt.Sprintf("member-%02d-0000-0000-0000-000000000000", i),
				"displayName": fmt.Sprintf("Test Member %d", i),
				"@odata.type": "#microsoft.graph.user",
			}
		}
	}
	mockState.Unlock()

	// List users - GET /users
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/users`, m.handleListUsers)

	// List devices - GET /devices
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/devices`, m.handleListDevices)

	// List group members - GET /groups/{id}/members
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/[^/]+/members`, m.handleListGroupMembers)

	// Get group by ID - GET /groups/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/[0-9a-fA-F-]+$`, m.handleGetGroup)
}

func (m *GuidListSharderMock) RegisterErrorMocks() {
	// Error scenarios for testing
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/users`, httpmock.NewStringResponder(500, `{"error":{"code":"InternalServerError","message":"Internal server error"}}`))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/devices`, httpmock.NewStringResponder(500, `{"error":{"code":"InternalServerError","message":"Internal server error"}}`))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/[^/]+/members`, httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Group not found"}}`))
}

func (m *GuidListSharderMock) CleanupMockState() {
	mockState.Lock()
	mockState.users = []map[string]any{}
	mockState.devices = []map[string]any{}
	mockState.groupMembers = make(map[string][]map[string]any)
	mockState.Unlock()
}

// handleListUsers returns paginated user list with optional filtering
func (m *GuidListSharderMock) handleListUsers(req *http.Request) (*http.Response, error) {
	mockState.Lock()
	defer mockState.Unlock()

	// Check for filter query param
	filter := req.URL.Query().Get("$filter")
	users := mockState.users

	// Apply filter if present
	if filter != "" {
		filteredUsers := []map[string]any{}
		// Simple filter support for accountEnabled
		if strings.Contains(filter, "accountEnabled eq true") {
			for _, user := range users {
				if enabled, ok := user["accountEnabled"].(bool); ok && enabled {
					filteredUsers = append(filteredUsers, user)
				}
			}
			users = filteredUsers
		}
	}

	// For simplicity, return all users in one page (no pagination for tests)
	response := map[string]any{
		"@odata.context": "https://graph.microsoft.com/beta/$metadata#users",
		"value":          users,
	}

	return httpmock.NewJsonResponse(200, response)
}

// handleListDevices returns paginated device list with optional filtering
func (m *GuidListSharderMock) handleListDevices(req *http.Request) (*http.Response, error) {
	mockState.Lock()
	defer mockState.Unlock()

	devices := mockState.devices

	// For simplicity, return all devices in one page (no pagination for tests)
	response := map[string]any{
		"@odata.context": "https://graph.microsoft.com/beta/$metadata#devices",
		"value":          devices,
	}

	return httpmock.NewJsonResponse(200, response)
}

// handleListGroupMembers returns group members with optional filtering
func (m *GuidListSharderMock) handleListGroupMembers(req *http.Request) (*http.Response, error) {
	mockState.Lock()
	defer mockState.Unlock()

	// Extract group ID from path
	parts := strings.Split(req.URL.Path, "/")
	groupId := ""
	for i, part := range parts {
		if part == "groups" && i+1 < len(parts) {
			groupId = parts[i+1]
			break
		}
	}

	members, exists := mockState.groupMembers[groupId]
	if !exists {
		return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Group not found"}}`), nil
	}

	// For simplicity, return all members in one page (no pagination for tests)
	response := map[string]any{
		"@odata.context": fmt.Sprintf("https://graph.microsoft.com/beta/$metadata#groups('%s')/members", groupId),
		"value":          members,
	}

	return httpmock.NewJsonResponse(200, response)
}

// handleGetGroup returns a mock group by ID
func (m *GuidListSharderMock) handleGetGroup(req *http.Request) (*http.Response, error) {
	// Extract group ID from path
	parts := strings.Split(req.URL.Path, "/")
	groupId := parts[len(parts)-1]

	response := map[string]any{
		"id":              groupId,
		"displayName":     "Test Group",
		"mailEnabled":     false,
		"securityEnabled": true,
	}

	return httpmock.NewJsonResponse(200, response)
}
