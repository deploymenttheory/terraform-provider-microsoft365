package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/jarcoal/httpmock"
)

// MockState holds the mock service principal state
type MockState struct {
	servicePrincipals map[string]map[string]any
	sync.Mutex
}

var mockState = &MockState{
	servicePrincipals: make(map[string]map[string]any),
}

// CleanupMockState cleans up the mock state
func (m *MockState) CleanupMockState() {
	m.Lock()
	defer m.Unlock()
	m.servicePrincipals = make(map[string]map[string]any)
}

// RegisterServicePrincipalMockResponders registers mock HTTP responders for service principal operations
func RegisterServicePrincipalMockResponders() *MockState {
	mockState.Lock()
	mockState.servicePrincipals = make(map[string]map[string]any)
	mockState.Unlock()

	// Create service principal - POST /servicePrincipals
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/servicePrincipals$`,
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			appId, ok := requestBody["appId"].(string)
			if !ok {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"appId is required"}}`), nil
			}

			// Generate mock service principal
			id := "11111111-1111-1111-1111-111111111111"
			servicePrincipal := map[string]any{
				"@odata.context":            "https://graph.microsoft.com/beta/$metadata#servicePrincipals/$entity",
				"id":                        id,
				"appId":                     appId,
				"displayName":               "Test Service Principal",
				"accountEnabled":            true,
				"appRoleAssignmentRequired": false,
				"servicePrincipalType":      "Application",
				"servicePrincipalNames":     []string{appId, "https://test.app"},
				"signInAudience":            "AzureADMyOrg",
				"tags":                      []string{},
			}

			// Apply optional fields from request
			if accountEnabled, ok := requestBody["accountEnabled"].(bool); ok {
				servicePrincipal["accountEnabled"] = accountEnabled
			}
			if appRoleAssignmentRequired, ok := requestBody["appRoleAssignmentRequired"].(bool); ok {
				servicePrincipal["appRoleAssignmentRequired"] = appRoleAssignmentRequired
			}
			if description, ok := requestBody["description"].(string); ok {
				servicePrincipal["description"] = description
			}
			if loginUrl, ok := requestBody["loginUrl"].(string); ok {
				servicePrincipal["loginUrl"] = loginUrl
			}
			if notes, ok := requestBody["notes"].(string); ok {
				servicePrincipal["notes"] = notes
			}
			if notificationEmailAddresses, ok := requestBody["notificationEmailAddresses"].([]interface{}); ok {
				servicePrincipal["notificationEmailAddresses"] = notificationEmailAddresses
			}
			if preferredSingleSignOnMode, ok := requestBody["preferredSingleSignOnMode"].(string); ok {
				servicePrincipal["preferredSingleSignOnMode"] = preferredSingleSignOnMode
			}
			if tags, ok := requestBody["tags"].([]interface{}); ok {
				servicePrincipal["tags"] = tags
			}

			mockState.Lock()
			mockState.servicePrincipals[id] = servicePrincipal
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, servicePrincipal)
		})

	// Get service principal - GET /servicePrincipals/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			servicePrincipal, exists := mockState.servicePrincipals[id]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Service principal not found"}}`), nil
			}

			return httpmock.NewJsonResponse(200, servicePrincipal)
		})

	// Update service principal - PATCH /servicePrincipals/{id}
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			var requestBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			mockState.Lock()
			servicePrincipal, exists := mockState.servicePrincipals[id]
			if !exists {
				mockState.Unlock()
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Service principal not found"}}`), nil
			}

			// Update fields
			if accountEnabled, ok := requestBody["accountEnabled"].(bool); ok {
				servicePrincipal["accountEnabled"] = accountEnabled
			}
			if appRoleAssignmentRequired, ok := requestBody["appRoleAssignmentRequired"].(bool); ok {
				servicePrincipal["appRoleAssignmentRequired"] = appRoleAssignmentRequired
			}
			if description, ok := requestBody["description"].(string); ok {
				servicePrincipal["description"] = description
			}
			if loginUrl, ok := requestBody["loginUrl"].(string); ok {
				servicePrincipal["loginUrl"] = loginUrl
			}
			if notes, ok := requestBody["notes"].(string); ok {
				servicePrincipal["notes"] = notes
			}
			if notificationEmailAddresses, ok := requestBody["notificationEmailAddresses"].([]interface{}); ok {
				servicePrincipal["notificationEmailAddresses"] = notificationEmailAddresses
			}
			if preferredSingleSignOnMode, ok := requestBody["preferredSingleSignOnMode"].(string); ok {
				servicePrincipal["preferredSingleSignOnMode"] = preferredSingleSignOnMode
			}
			if tags, ok := requestBody["tags"].([]interface{}); ok {
				servicePrincipal["tags"] = tags
			}

			mockState.servicePrincipals[id] = servicePrincipal
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Delete service principal - DELETE /servicePrincipals/{id}
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/servicePrincipals/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			_, exists := mockState.servicePrincipals[id]
			if !exists {
				mockState.Unlock()
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Service principal not found"}}`), nil
			}

			delete(mockState.servicePrincipals, id)
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	return mockState
}
