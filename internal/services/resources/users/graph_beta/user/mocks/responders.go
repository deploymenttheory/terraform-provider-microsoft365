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
	users map[string]map[string]interface{}
}

func init() {
	// Initialize mockState
	mockState.users = make(map[string]map[string]interface{})

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// UserMock provides mock responses for user operations
type UserMock struct{}

// RegisterMocks registers HTTP mock responses for user operations
func (m *UserMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.users = make(map[string]map[string]interface{})
	mockState.Unlock()

	// Register specific test users
	registerTestUsers()

	// Register GET for user by ID
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/users/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			userId := urlParts[len(urlParts)-1]

			mockState.Lock()
			userData, exists := mockState.users[userId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"User not found"}}`), nil
			}

			return httpmock.NewJsonResponse(200, userData)
		})

	// Register GET for listing users
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/users(\?.+)?$`,
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			defer mockState.Unlock()

			users := make([]map[string]interface{}, 0, len(mockState.users))
			for _, user := range mockState.users {
				users = append(users, user)
			}

			response := map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#users",
				"value":          users,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register POST for creating users
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/users",
		func(req *http.Request) (*http.Response, error) {
			var userData map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&userData)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Validate required fields
			if _, ok := userData["displayName"].(string); !ok {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"displayName is required"}}`), nil
			}
			if _, ok := userData["userPrincipalName"].(string); !ok {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"userPrincipalName is required"}}`), nil
			}
			if passwordProfile, ok := userData["passwordProfile"].(map[string]interface{}); !ok || passwordProfile["password"] == nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"passwordProfile with password is required"}}`), nil
			}

			// Generate ID if not provided
			if userData["id"] == nil {
				userData["id"] = uuid.New().String()
			}

			// Set computed fields
			now := time.Now().Format(time.RFC3339)
			userData["createdDateTime"] = now

			// Ensure collection fields are initialized
			commonMocks.EnsureField(userData, "businessPhones", []string{})
			commonMocks.EnsureField(userData, "identities", []map[string]interface{}{})
			commonMocks.EnsureField(userData, "imAddresses", []string{})
			commonMocks.EnsureField(userData, "otherMails", []string{})
			commonMocks.EnsureField(userData, "proxyAddresses", []string{})

			// Store user in mock state
			userId := userData["id"].(string)
			mockState.Lock()
			mockState.users[userId] = userData
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, userData)
		})

	// Register PATCH for updating users
	httpmock.RegisterResponder("PATCH", `=~^https://graph.microsoft.com/beta/users/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			userId := urlParts[len(urlParts)-1]

			mockState.Lock()
			userData, exists := mockState.users[userId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"User not found"}}`), nil
			}

			var updateData map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&updateData)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Update user data
			mockState.Lock()

			// Special handling for updates that remove fields
			// If we're updating from maximal to minimal, we need to remove fields not in the minimal config
			// Check if this is a minimal update by looking for key indicators
			isMinimalUpdate := false
			if _, hasDisplayName := updateData["displayName"]; hasDisplayName {
				if _, hasGivenName := updateData["givenName"]; !hasGivenName {
					isMinimalUpdate = true
				}
			}

			if isMinimalUpdate {
				// Remove fields that are not part of minimal configuration
				fieldsToRemove := []string{
					"givenName", "surname", "jobTitle", "department", "companyName",
					"officeLocation", "city", "state", "country", "postalCode",
					"mobilePhone", "mail", "mailNickname", "usageLocation",
				}

				for _, field := range fieldsToRemove {
					delete(userData, field)
				}

				// Reset collections to empty
				userData["businessPhones"] = []string{}
				userData["otherMails"] = []string{}
				userData["proxyAddresses"] = []string{}
			}

			// Apply the updates
			for k, v := range updateData {
				userData[k] = v
			}

			mockState.users[userId] = userData
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, userData)
		})

	// Register DELETE for removing users
	httpmock.RegisterResponder("DELETE", `=~^https://graph.microsoft.com/beta/users/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			userId := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.users[userId]
			if exists {
				delete(mockState.users, userId)
			}
			mockState.Unlock()

			// Return 204 No Content for successful deletion
			return httpmock.NewStringResponse(204, ""), nil
		})
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *UserMock) RegisterErrorMocks() {
	// Register error response for user creation
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/users",
		factories.ErrorResponse(400, "BadRequest", "Error creating user"))

	// Register error response for user not found
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/users/not-found-user",
		factories.ErrorResponse(404, "ResourceNotFound", "User not found"))

	// Register error response for duplicate user principal name
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/users",
		func(req *http.Request) (*http.Response, error) {
			var userData map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&userData)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			if upn, ok := userData["userPrincipalName"].(string); ok && upn == "duplicate@contoso.com" {
				return factories.ErrorResponse(400, "BadRequest", "User with this userPrincipalName already exists")(req)
			}

			// Fallback to normal creation flow
			return nil, nil
		})
}

// registerTestUsers registers predefined test users
func registerTestUsers() {
	// Minimal user with only required attributes
	minimalUserId := "00000000-0000-0000-0000-000000000001"
	minimalUserData := map[string]interface{}{
		"id":                minimalUserId,
		"displayName":       "Minimal User",
		"userPrincipalName": "minimal.user@contoso.com",
		"accountEnabled":    true,
		"passwordProfile": map[string]interface{}{
			"password":                             "SecureP@ssw0rd!",
			"forceChangePasswordNextSignIn":        false,
			"forceChangePasswordNextSignInWithMfa": false,
		},
		"createdDateTime": "2023-01-01T00:00:00Z",
		"businessPhones":  []string{},
		"identities":      []map[string]interface{}{},
		"imAddresses":     []string{},
		"otherMails":      []string{},
		"proxyAddresses":  []string{},
	}

	// Maximal user with all attributes
	maximalUserId := "00000000-0000-0000-0000-000000000002"
	maximalUserData := map[string]interface{}{
		"id":                maximalUserId,
		"displayName":       "Maximal User",
		"userPrincipalName": "maximal.user@contoso.com",
		"accountEnabled":    true,
		"givenName":         "Maximal",
		"surname":           "User",
		"mail":              "maximal.user@contoso.com",
		"mailNickname":      "maxuser",
		"jobTitle":          "Senior Developer",
		"department":        "Engineering",
		"companyName":       "Contoso Ltd",
		"officeLocation":    "Building A",
		"city":              "Redmond",
		"state":             "WA",
		"country":           "US",
		"postalCode":        "98052",
		"usageLocation":     "US",
		"businessPhones":    []string{"+1 425-555-0100"},
		"mobilePhone":       "+1 425-555-0101",
		"passwordProfile": map[string]interface{}{
			"password":                             "SecureP@ssw0rd!",
			"forceChangePasswordNextSignIn":        true,
			"forceChangePasswordNextSignInWithMfa": false,
		},
		"identities": []map[string]interface{}{
			{
				"signInType":       "emailAddress",
				"issuer":           "contoso.com",
				"issuerAssignedId": "maximal.user@contoso.com",
			},
		},
		"otherMails":      []string{"maximal.user.other@contoso.com"},
		"proxyAddresses":  []string{"SMTP:maximal.user@contoso.com"},
		"createdDateTime": "2023-01-01T00:00:00Z",
		"imAddresses":     []string{},
	}

	// Store users in mock state
	mockState.Lock()
	mockState.users[minimalUserId] = minimalUserData
	mockState.users[maximalUserId] = maximalUserData
	mockState.Unlock()
}

// Helper function to ensure collection fields exist
func ensureCollectionField(data map[string]interface{}, fieldName string, defaultValue interface{}) {
	if data[fieldName] == nil {
		data[fieldName] = defaultValue
	}
}
