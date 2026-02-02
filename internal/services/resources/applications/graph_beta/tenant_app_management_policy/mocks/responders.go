package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sync"

	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	policy map[string]any
}

func init() {
	mockState.policy = make(map[string]any)
}

// TenantAppManagementPolicyMock contains methods to register mocks for the Tenant App Management Policy resource
type TenantAppManagementPolicyMock struct{}

// RegisterMocks registers HTTP mocks for the Tenant App Management Policy resource
func (m *TenantAppManagementPolicyMock) RegisterMocks() {
	// Initialize default state
	mockState.Lock()
	mockState.policy = map[string]any{
		"@odata.context": "https://graph.microsoft.com/beta/$metadata#policies/defaultAppManagementPolicy/$entity",
		"id":             "00000000-0000-0000-0000-000000000000",
		"displayName":    "Default app management tenant policy",
		"description":    "Default tenant policy that enforces app management restrictions on applications and service principals.",
		"isEnabled":      false,
		"applicationRestrictions": map[string]any{
			"passwordCredentials": []any{},
			"keyCredentials":      []any{},
		},
		"servicePrincipalRestrictions": map[string]any{
			"passwordCredentials": []any{},
			"keyCredentials":      []any{},
		},
	}
	mockState.Unlock()

	// Register GET mock for the resource
	httpmock.RegisterResponder(
		"GET",
		"https://graph.microsoft.com/beta/policies/defaultAppManagementPolicy",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			defer mockState.Unlock()

			responseJSON, err := json.Marshal(mockState.policy)
			if err != nil {
				return httpmock.NewStringResponse(http.StatusInternalServerError, `{"error": {"message": "Failed to create response"}}`), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, string(responseJSON)), nil
		},
	)

	// Register PATCH mock for updating the resource
	httpmock.RegisterResponder(
		"PATCH",
		"https://graph.microsoft.com/beta/policies/defaultAppManagementPolicy",
		func(req *http.Request) (*http.Response, error) {
			// Parse the request body to get the updated values
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(http.StatusBadRequest, `{"error": {"message": "Invalid request body"}}`), nil
			}

			mockState.Lock()
			defer mockState.Unlock()

			// Update the mock state with the request body
			if displayName, ok := requestBody["displayName"].(string); ok {
				mockState.policy["displayName"] = displayName
			}

			if description, ok := requestBody["description"].(string); ok {
				mockState.policy["description"] = description
			}

			if isEnabled, ok := requestBody["isEnabled"].(bool); ok {
				mockState.policy["isEnabled"] = isEnabled
			}

			if appRestrictions, ok := requestBody["applicationRestrictions"].(map[string]any); ok {
				// Add state field to all credentials
				if pwdCreds, ok := appRestrictions["passwordCredentials"].([]any); ok {
					for _, cred := range pwdCreds {
						if credMap, ok := cred.(map[string]any); ok && credMap["state"] == nil {
							credMap["state"] = "enabled"
						}
					}
				}
				if keyCreds, ok := appRestrictions["keyCredentials"].([]any); ok {
					for _, cred := range keyCreds {
						if credMap, ok := cred.(map[string]any); ok && credMap["state"] == nil {
							credMap["state"] = "enabled"
						}
					}
				}
				mockState.policy["applicationRestrictions"] = appRestrictions
			}

			if spRestrictions, ok := requestBody["servicePrincipalRestrictions"].(map[string]any); ok {
				// Add state field to all credentials
				if pwdCreds, ok := spRestrictions["passwordCredentials"].([]any); ok {
					for _, cred := range pwdCreds {
						if credMap, ok := cred.(map[string]any); ok && credMap["state"] == nil {
							credMap["state"] = "enabled"
						}
					}
				}
				if keyCreds, ok := spRestrictions["keyCredentials"].([]any); ok {
					for _, cred := range keyCreds {
						if credMap, ok := cred.(map[string]any); ok && credMap["state"] == nil {
							credMap["state"] = "enabled"
						}
					}
				}
				mockState.policy["servicePrincipalRestrictions"] = spRestrictions
			}

			// PATCH returns 204 No Content per Microsoft docs
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		},
	)
}

// RegisterErrorMocks registers HTTP mocks that return errors
func (m *TenantAppManagementPolicyMock) RegisterErrorMocks() {
	// Register GET mock that returns an error
	httpmock.RegisterRegexpResponder(
		"GET",
		regexp.MustCompile(`https://graph.microsoft.com/beta/policies/defaultAppManagementPolicy`),
		httpmock.NewStringResponder(
			http.StatusForbidden,
			`{
				"error": {
					"code": "Forbidden",
					"message": "Access denied. You do not have permission to perform this action or access this resource.",
					"innerError": {
						"date": "2024-01-01T12:00:00",
						"request-id": "00000000-0000-0000-0000-000000000000"
					}
				}
			}`,
		),
	)

	// Register PATCH mock that returns an error
	httpmock.RegisterRegexpResponder(
		"PATCH",
		regexp.MustCompile(`https://graph.microsoft.com/beta/policies/defaultAppManagementPolicy`),
		httpmock.NewStringResponder(
			http.StatusBadRequest,
			`{
				"error": {
					"code": "BadRequest",
					"message": "Invalid value specified for property 'isEnabled' of resource 'tenantAppManagementPolicy'.",
					"innerError": {
						"date": "2024-01-01T12:00:00",
						"request-id": "00000000-0000-0000-0000-000000000000"
					}
				}
			}`,
		),
	)
}

// CleanupMockState cleans up the mock state
func (m *TenantAppManagementPolicyMock) CleanupMockState() {
	mockState.Lock()
	mockState.policy = map[string]any{
		"@odata.context": "https://graph.microsoft.com/beta/$metadata#policies/defaultAppManagementPolicy/$entity",
		"id":             "00000000-0000-0000-0000-000000000000",
		"displayName":    "Default app management tenant policy",
		"description":    "Default tenant policy that enforces app management restrictions on applications and service principals.",
		"isEnabled":      false,
		"applicationRestrictions": map[string]any{
			"passwordCredentials": []any{},
			"keyCredentials":      []any{},
		},
		"servicePrincipalRestrictions": map[string]any{
			"passwordCredentials": []any{},
			"keyCredentials":      []any{},
		},
	}
	mockState.Unlock()
}

// GetMockResponse returns a mock response for the given scenario
func (m *TenantAppManagementPolicyMock) GetMockResponse(scenario string) string {
	if scenario == "error" {
		return `{
			"error": {
				"code": "NotFound",
				"message": "Resource not found",
				"innerError": {
					"date": "2024-01-01T12:00:00",
					"request-id": "00000000-0000-0000-0000-000000000000"
				}
			}
		}`
	}

	return fmt.Sprintf(`{
		"@odata.context": "https://graph.microsoft.com/beta/$metadata#policies/defaultAppManagementPolicy/$entity",
		"id": "00000000-0000-0000-0000-000000000000",
		"displayName": "Default app management tenant policy",
		"description": "Default tenant policy that enforces app management restrictions on applications and service principals.",
		"isEnabled": true,
		"applicationRestrictions": {
			"passwordCredentials": [
				{
					"restrictionType": "passwordLifetime",
					"state": "enabled",
					"restrictForAppsCreatedAfterDateTime": "2024-01-01T00:00:00Z",
					"maxLifetime": "P90D"
				}
			],
			"keyCredentials": [],
			"identifierUris": {}
		},
		"servicePrincipalRestrictions": {
			"passwordCredentials": [],
			"keyCredentials": []
		}
	}`)
}
