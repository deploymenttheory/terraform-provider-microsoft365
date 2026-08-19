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
	deletedItems      map[string]map[string]any
	sync.Mutex
}

var mockState = &MockState{
	servicePrincipals: make(map[string]map[string]any),
	deletedItems:      make(map[string]map[string]any),
}

// CleanupMockState cleans up the mock state
func (m *MockState) CleanupMockState() {
	m.Lock()
	defer m.Unlock()
	m.servicePrincipals = make(map[string]map[string]any)
	m.deletedItems = make(map[string]map[string]any)
}

// hasKeyCredential reports whether the service principal's keyCredentials collection
// contains the given keyId. Graph enforces this reference before accepting a
// tokenEncryptionKeyId ("No KeyCredential found with the configured TokenEncryptionKeyId").
func hasKeyCredential(servicePrincipal map[string]any, keyId string) bool {
	credentials, ok := servicePrincipal["keyCredentials"].([]any)
	if !ok {
		return false
	}
	for _, raw := range credentials {
		if credential, ok := raw.(map[string]any); ok && credential["keyId"] == keyId {
			return true
		}
	}
	return false
}

const invalidTokenEncryptionKeyIdResponse = `{"error":{"code":"Request_BadRequest","message":"No KeyCredential found with the configured TokenEncryptionKeyId."}}`

// RegisterServicePrincipalMockResponders registers mock HTTP responders for service principal operations
func RegisterServicePrincipalMockResponders() *MockState {
	mockState.Lock()
	mockState.servicePrincipals = make(map[string]map[string]any)
	mockState.deletedItems = make(map[string]map[string]any)
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

			// Generate mock service principal. Read-only fields mirror the shape of a real
			// GET /servicePrincipals/{id} response for the properties this resource models.
			id := "11111111-1111-1111-1111-111111111111"
			servicePrincipal := map[string]any{
				"@odata.context":            "https://graph.microsoft.com/beta/$metadata#servicePrincipals/$entity",
				"id":                        id,
				"appId":                     appId,
				"accountEnabled":            true,
				"appRoleAssignmentRequired": false,
				"servicePrincipalType":      "Application",
				"tags":                      []string{},
				"appOwnerOrganizationId":    "2cbe968c-9683-4d8a-af06-dab1bb350a04",
				"createdByAppId":            "04b07795-8ddb-461a-bbee-02f9e1bf7b46",
				"keyCredentials":            []any{},
				"passwordCredentials":       []any{},
			}

			// Requests that carry a description simulate a service principal whose backing
			// application already has certificates (e.g. instantiated from a gallery
			// template), so the wire-format mapping of both credential collections is
			// exercised; bare requests keep the empty collections a fresh principal has.
			if _, hasDescription := requestBody["description"]; hasDescription {
				servicePrincipal["keyCredentials"] = []any{
					map[string]any{
						// Base64 of the certificate thumbprint, as the API returns it
						"customKeyIdentifier": "a8NSGsQqlkjIPN1kEpJ8QIe4AgI=",
						"displayName":         "CN=test-signing",
						"endDateTime":         "2027-01-01T00:00:00Z",
						"keyId":               "dddddddd-1111-2222-3333-444444444444",
						"startDateTime":       "2026-01-01T00:00:00Z",
						"type":                "AsymmetricX509Cert",
						"usage":               "Sign",
					},
				}
				servicePrincipal["passwordCredentials"] = []any{
					map[string]any{
						"customKeyIdentifier": "a8NSGsQqlkjIPN1kEpJ8QIe4AgI=",
						"displayName":         "test-secret",
						"endDateTime":         "2027-01-01T00:00:00Z",
						"hint":                "abc",
						"keyId":               "eeeeeeee-1111-2222-3333-444444444444",
						"startDateTime":       "2026-01-01T00:00:00Z",
					},
				}
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
			if notificationEmailAddresses, ok := requestBody["notificationEmailAddresses"].([]any); ok {
				servicePrincipal["notificationEmailAddresses"] = notificationEmailAddresses
			}
			if preferredSingleSignOnMode, ok := requestBody["preferredSingleSignOnMode"].(string); ok {
				servicePrincipal["preferredSingleSignOnMode"] = preferredSingleSignOnMode
			}
			if tags, ok := requestBody["tags"].([]any); ok {
				servicePrincipal["tags"] = tags
			}
			if alternativeNames, ok := requestBody["alternativeNames"].([]any); ok {
				servicePrincipal["alternativeNames"] = alternativeNames
			}
			if samlSingleSignOnSettings, ok := requestBody["samlSingleSignOnSettings"].(map[string]any); ok {
				servicePrincipal["samlSingleSignOnSettings"] = samlSingleSignOnSettings
			}
			if tokenEncryptionKeyId, ok := requestBody["tokenEncryptionKeyId"].(string); ok {
				if !hasKeyCredential(servicePrincipal, tokenEncryptionKeyId) {
					return httpmock.NewStringResponse(400, invalidTokenEncryptionKeyIdResponse), nil
				}
				servicePrincipal["tokenEncryptionKeyId"] = tokenEncryptionKeyId
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
			if notificationEmailAddresses, ok := requestBody["notificationEmailAddresses"].([]any); ok {
				servicePrincipal["notificationEmailAddresses"] = notificationEmailAddresses
			}
			if preferredSingleSignOnMode, ok := requestBody["preferredSingleSignOnMode"].(string); ok {
				servicePrincipal["preferredSingleSignOnMode"] = preferredSingleSignOnMode
			}
			if tags, ok := requestBody["tags"].([]any); ok {
				servicePrincipal["tags"] = tags
			}
			if alternativeNames, ok := requestBody["alternativeNames"].([]any); ok {
				servicePrincipal["alternativeNames"] = alternativeNames
			}
			// An explicit JSON null clears the property (sent when the block is removed)
			if raw, keyPresent := requestBody["samlSingleSignOnSettings"]; keyPresent {
				if raw == nil {
					delete(servicePrincipal, "samlSingleSignOnSettings")
				} else if samlSingleSignOnSettings, ok := raw.(map[string]any); ok {
					servicePrincipal["samlSingleSignOnSettings"] = samlSingleSignOnSettings
				}
			}
			if raw, keyPresent := requestBody["tokenEncryptionKeyId"]; keyPresent {
				if raw == nil {
					delete(servicePrincipal, "tokenEncryptionKeyId")
				} else if tokenEncryptionKeyId, ok := raw.(string); ok {
					if !hasKeyCredential(servicePrincipal, tokenEncryptionKeyId) {
						mockState.Unlock()
						return httpmock.NewStringResponse(400, invalidTokenEncryptionKeyIdResponse), nil
					}
					servicePrincipal["tokenEncryptionKeyId"] = tokenEncryptionKeyId
				}
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
			servicePrincipal, exists := mockState.servicePrincipals[id]
			if !exists {
				mockState.Unlock()
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Service principal not found"}}`), nil
			}

			// Move to deleted items (soft delete)
			servicePrincipal["@odata.type"] = "#microsoft.graph.servicePrincipal"
			mockState.deletedItems[id] = servicePrincipal
			delete(mockState.servicePrincipals, id)
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Get deleted item - GET /directory/deletedItems/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directory/deletedItems/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			deletedItem, exists := mockState.deletedItems[id]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
			}

			return httpmock.NewJsonResponse(200, deletedItem)
		})

	// Permanently delete - DELETE /directory/deletedItems/{id}
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/directory/deletedItems/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			_, exists := mockState.deletedItems[id]
			if !exists {
				mockState.Unlock()
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
			}

			delete(mockState.deletedItems, id)
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	return mockState
}
