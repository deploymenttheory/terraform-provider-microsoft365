package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	commonMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	groups map[string]map[string]interface{}
}

func init() {
	// Initialize mockState
	mockState.groups = make(map[string]map[string]interface{})

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// GroupMock provides mock responses for group operations
type GroupMock struct{}

// RegisterMocks registers HTTP mock responses for group operations
func (m *GroupMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.groups = make(map[string]map[string]interface{})
	mockState.Unlock()

	// Register specific test groups
	registerTestGroups()

	// Register GET for group by ID
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

	// Register GET for listing groups
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/groups(\?.+)?$`,
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			defer mockState.Unlock()

			groups := make([]map[string]interface{}, 0, len(mockState.groups))
			for _, group := range mockState.groups {
				groups = append(groups, group)
			}

			response := map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#groups",
				"value":          groups,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register POST for creating groups
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/groups",
		func(req *http.Request) (*http.Response, error) {
			var groupData map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&groupData)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Validate required fields
			if _, ok := groupData["displayName"].(string); !ok {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"displayName is required"}}`), nil
			}
			if _, ok := groupData["mailNickname"].(string); !ok {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"mailNickname is required"}}`), nil
			}
			if _, ok := groupData["mailEnabled"].(bool); !ok {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"mailEnabled is required"}}`), nil
			}
			if _, ok := groupData["securityEnabled"].(bool); !ok {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"securityEnabled is required"}}`), nil
			}

			// Generate ID if not provided
			if groupData["id"] == nil {
				groupData["id"] = uuid.New().String()
			}

			// Set computed fields
			now := time.Now().Format(time.RFC3339)
			groupData["createdDateTime"] = now

			// Set default values for optional fields if not provided
			if groupData["visibility"] == nil {
				groupData["visibility"] = "Private"
			}
			if groupData["isAssignableToRole"] == nil {
				groupData["isAssignableToRole"] = false
			}
			if groupData["membershipRuleProcessingState"] == nil {
				groupData["membershipRuleProcessingState"] = "Paused"
			}

			// Generate mail if mail-enabled
			if mailEnabled, ok := groupData["mailEnabled"].(bool); ok && mailEnabled {
				if mailNickname, ok := groupData["mailNickname"].(string); ok {
					groupData["mail"] = mailNickname + "@contoso.com"
				}
			}

			// Set security identifier if security-enabled
			if securityEnabled, ok := groupData["securityEnabled"].(bool); ok && securityEnabled {
				groupData["securityIdentifier"] = "S-1-12-1-" + uuid.New().String()
			}

			// Ensure collection fields are initialized
			commonMocks.EnsureField(groupData, "proxyAddresses", []string{})
			if mailEnabled, ok := groupData["mailEnabled"].(bool); ok && mailEnabled {
				if mail, ok := groupData["mail"].(string); ok {
					groupData["proxyAddresses"] = []string{"SMTP:" + mail}
				}
			}

			// Store group in mock state
			groupId := groupData["id"].(string)
			mockState.Lock()
			mockState.groups[groupId] = groupData
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, groupData)
		})

	// Register PATCH for updating groups
	httpmock.RegisterResponder("PATCH", `=~^https://graph.microsoft.com/beta/groups/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			groupId := urlParts[len(urlParts)-1]

			mockState.Lock()
			groupData, exists := mockState.groups[groupId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Group not found"}}`), nil
			}

			var updateData map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&updateData)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Update group data
			mockState.Lock()

			// Special handling for updates that remove fields
			// If we're updating from maximal to minimal, we need to remove fields not in the minimal config
			// Check if this is a minimal update by looking for key indicators
			isMinimalUpdate := false
			if _, hasDisplayName := updateData["displayName"]; hasDisplayName {
				if _, hasDescription := updateData["description"]; !hasDescription {
					isMinimalUpdate = true
				}
			}

			if isMinimalUpdate {
				// Remove fields that are not part of minimal configuration
				fieldsToRemove := []string{
					"description", "groupTypes", "membershipRule", "preferredDataLocation",
					"preferredLanguage", "theme", "classification",
				}

				for _, field := range fieldsToRemove {
					delete(groupData, field)
				}

				// Reset fields to defaults
				groupData["visibility"] = "Private"
				groupData["isAssignableToRole"] = false
				groupData["membershipRuleProcessingState"] = "Paused"
			}

			// Apply the updates
			for k, v := range updateData {
				groupData[k] = v
			}

			// Update mail and proxy addresses if mail-enabled changes
			if mailEnabled, ok := updateData["mailEnabled"].(bool); ok {
				if mailEnabled {
					if mailNickname, ok := groupData["mailNickname"].(string); ok {
						groupData["mail"] = mailNickname + "@contoso.com"
						groupData["proxyAddresses"] = []string{"SMTP:" + mailNickname + "@contoso.com"}
					}
				} else {
					delete(groupData, "mail")
					groupData["proxyAddresses"] = []string{}
				}
			}

			// Update security identifier if security-enabled changes
			if securityEnabled, ok := updateData["securityEnabled"].(bool); ok {
				if securityEnabled {
					groupData["securityIdentifier"] = "S-1-12-1-" + uuid.New().String()
				} else {
					delete(groupData, "securityIdentifier")
				}
			}

			mockState.groups[groupId] = groupData
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, groupData)
		})

	// Register DELETE for removing groups
	httpmock.RegisterResponder("DELETE", `=~^https://graph.microsoft.com/beta/groups/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			groupId := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.groups[groupId]
			if exists {
				delete(mockState.groups, groupId)
			}
			mockState.Unlock()

			// Return 204 No Content for successful deletion
			return httpmock.NewStringResponse(204, ""), nil
		})
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *GroupMock) RegisterErrorMocks() {
	// Register error response for group creation
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/groups",
		factories.ErrorResponse(400, "BadRequest", "Error creating group"))

	// Register error response for group not found
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/groups/not-found-group",
		factories.ErrorResponse(404, "ResourceNotFound", "Group not found"))

	// Register error response for duplicate display name
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/groups",
		func(req *http.Request) (*http.Response, error) {
			var groupData map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&groupData)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			if displayName, ok := groupData["displayName"].(string); ok && displayName == "Error Group" {
				return factories.ErrorResponse(400, "BadRequest", "Group with this displayName already exists")(req)
			}

			// Fallback to normal creation flow
			return nil, nil
		})
}

// registerTestGroups registers predefined test groups
func registerTestGroups() {
	// Minimal group with only required attributes
	minimalGroupId := "00000000-0000-0000-0000-000000000001"
	minimalGroupData := map[string]interface{}{
		"id":                            minimalGroupId,
		"displayName":                   "Minimal Group",
		"mailNickname":                  "minimal.group",
		"mailEnabled":                   false,
		"securityEnabled":               true,
		"visibility":                    "Private",
		"isAssignableToRole":            false,
		"membershipRuleProcessingState": "Paused",
		"createdDateTime":               "2023-01-01T00:00:00Z",
		"proxyAddresses":                []string{},
		"securityIdentifier":            "S-1-12-1-1234567890-1234567890-1234567890-1234567890",
	}

	// Maximal group with all attributes
	maximalGroupId := "00000000-0000-0000-0000-000000000002"
	maximalGroupData := map[string]interface{}{
		"id":                            maximalGroupId,
		"displayName":                   "Maximal Group",
		"description":                   "This is a maximal group configuration for testing",
		"mailNickname":                  "maximal.group",
		"mailEnabled":                   true,
		"securityEnabled":               true,
		"groupTypes":                    []string{"Unified", "DynamicMembership"},
		"visibility":                    "Private",
		"isAssignableToRole":            false,
		"membershipRule":                "user.department -eq \"Engineering\"",
		"membershipRuleProcessingState": "On",
		"preferredDataLocation":         "NAM",
		"preferredLanguage":             "en-US",
		"theme":                         "Blue",
		"classification":                "High",
		"mail":                          "maximal.group@contoso.com",
		"proxyAddresses":                []string{"SMTP:maximal.group@contoso.com"},
		"securityIdentifier":            "S-1-12-1-2345678901-2345678901-2345678901-2345678901",
		"createdDateTime":               "2023-01-01T00:00:00Z",
	}

	// Store groups in mock state
	mockState.Lock()
	mockState.groups[minimalGroupId] = minimalGroupData
	mockState.groups[maximalGroupId] = maximalGroupData
	mockState.Unlock()
}
