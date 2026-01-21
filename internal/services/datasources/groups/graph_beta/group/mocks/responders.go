package mocks

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
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
	mocks.GlobalRegistry.Register("groups_group", &GroupMock{})
}

type GroupMock struct{}

var _ mocks.MockRegistrar = (*GroupMock)(nil)

func (m *GroupMock) RegisterMocks() {
	mockState.Lock()
	mockState.groups = make(map[string]map[string]any)
	mockState.Unlock()

	// 1. Get group by ID - GET /groups/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		groupId := parts[len(parts)-1]

		// Return mock response for known IDs
		switch groupId {
		case "00000000-0000-0000-0000-000000000001":
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_get/get_group_by_id.json")
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		case "00000000-0000-0000-0000-000000000002":
			// Return second group (Microsoft 365 group)
			responseObj := map[string]any{
				"id":                    "00000000-0000-0000-0000-000000000002",
				"displayName":           "Finance Team",
				"description":           "Finance Department Group",
				"mailNickname":          "finance",
				"mailEnabled":           true,
				"securityEnabled":       false,
				"groupTypes":            []string{"Unified"},
				"visibility":            "Public",
				"isAssignableToRole":    false,
				"mail":                  "finance@contoso.com",
				"proxyAddresses":        []string{"SMTP:finance@contoso.com"},
				"createdDateTime":       "2024-01-15T10:30:00Z",
				"onPremisesSyncEnabled": nil,
			}
			return httpmock.NewJsonResponse(200, responseObj)
		case "00000000-0000-0000-0000-000000000003":
			// Return third group (Dynamic group)
			responseObj := map[string]any{
				"id":                            "00000000-0000-0000-0000-000000000003",
				"displayName":                   "All Windows Devices",
				"description":                   "Dynamic group for all Windows devices",
				"mailNickname":                  "all-windows-devices",
				"mailEnabled":                   false,
				"securityEnabled":               true,
				"groupTypes":                    []string{"DynamicMembership"},
				"visibility":                    "Private",
				"isAssignableToRole":            false,
				"membershipRule":                "(device.deviceOSType -eq \"Windows\")",
				"membershipRuleProcessingState": "On",
				"createdDateTime":               "2024-02-01T14:00:00Z",
				"onPremisesSyncEnabled":         nil,
			}
			return httpmock.NewJsonResponse(200, responseObj)
		default:
			return httpmock.NewStringResponse(404, `{"error":{"code":"Request_ResourceNotFound","message":"Resource not found"}}`), nil
		}
	})

	// 2. List groups with filter - GET /groups?$filter=...
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups\?`, func(req *http.Request) (*http.Response, error) {
		queryParams, _ := url.ParseQuery(req.URL.RawQuery)
		filter := queryParams.Get("$filter")

		// Parse the filter to determine what to return
		if strings.Contains(filter, "displayName eq 'IT Security Team'") {
			// Return single group
			responseObj := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#groups",
				"value": []map[string]any{
					{
						"id":                    "00000000-0000-0000-0000-000000000001",
						"displayName":           "IT Security Team",
						"description":           "Security group for IT department members",
						"mailNickname":          "it-security",
						"mailEnabled":           false,
						"securityEnabled":       true,
						"groupTypes":            []any{},
						"visibility":            "Private",
						"isAssignableToRole":    false,
						"createdDateTime":       "2024-01-10T09:00:00Z",
						"onPremisesSyncEnabled": nil,
					},
				},
			}
			return httpmock.NewJsonResponse(200, responseObj)
		} else if strings.Contains(filter, "mailNickname eq 'finance'") {
			// Return finance group
			responseObj := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#groups",
				"value": []map[string]any{
					{
						"id":                    "00000000-0000-0000-0000-000000000002",
						"displayName":           "Finance Team",
						"description":           "Finance Department Group",
						"mailNickname":          "finance",
						"mailEnabled":           true,
						"securityEnabled":       false,
						"groupTypes":            []string{"Unified"},
						"visibility":            "Public",
						"isAssignableToRole":    false,
						"mail":                  "finance@contoso.com",
						"proxyAddresses":        []string{"SMTP:finance@contoso.com"},
						"createdDateTime":       "2024-01-15T10:30:00Z",
						"onPremisesSyncEnabled": nil,
					},
				},
			}
			return httpmock.NewJsonResponse(200, responseObj)
		}

		// Default empty response
		responseObj := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#groups",
			"value":          []any{},
		}
		return httpmock.NewJsonResponse(200, responseObj)
	})

	// 3. Get group members - GET /groups/{id}/members
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/[0-9a-fA-F-]+/members`, func(req *http.Request) (*http.Response, error) {
		responseObj := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#directoryObjects",
			"value":          []any{},
		}
		return httpmock.NewJsonResponse(200, responseObj)
	})

	// 4. Get group owners - GET /groups/{id}/owners
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/[0-9a-fA-F-]+/owners`, func(req *http.Request) (*http.Response, error) {
		responseObj := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#directoryObjects",
			"value":          []any{},
		}
		return httpmock.NewJsonResponse(200, responseObj)
	})
}

func (m *GroupMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.groups = make(map[string]map[string]any)
	mockState.Unlock()

	// Return errors for all operations
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups`, func(req *http.Request) (*http.Response, error) {
		errorObj := map[string]any{
			"error": map[string]any{
				"code":    "Forbidden",
				"message": "Insufficient privileges to complete the operation.",
			},
		}
		return httpmock.NewJsonResponse(403, errorObj)
	})
}

func (m *GroupMock) CleanupMockState() {
	mockState.Lock()
	mockState.groups = make(map[string]map[string]any)
	mockState.Unlock()
}
