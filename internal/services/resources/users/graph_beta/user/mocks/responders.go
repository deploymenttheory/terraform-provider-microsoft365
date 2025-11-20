package mocks

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	commonMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	users map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.users = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// UserMock provides mock responses for user operations
type UserMock struct{}

// RegisterMocks registers HTTP mock responses for user operations
func (m *UserMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.users = make(map[string]map[string]any)
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

			users := make([]map[string]any, 0, len(mockState.users))
			for _, user := range mockState.users {
				users = append(users, user)
			}

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#users",
				"value":          users,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register POST for creating users
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/users",
		func(req *http.Request) (*http.Response, error) {
			var userData map[string]any
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
			if passwordProfile, ok := userData["passwordProfile"].(map[string]any); !ok || passwordProfile["password"] == nil {
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
			commonMocks.EnsureField(userData, "identities", []map[string]any{})
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

			var updateData map[string]any
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
	// Reset the state when registering error mocks
	mockState.Lock()
	mockState.users = make(map[string]map[string]any)
	mockState.Unlock()

	// Register error response for user creation - always return error
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/users",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid user data"}}`), nil
		})

	// Register error response for user not found
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/users/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"User not found"}}`), nil
		})

	// Register error response for DELETE
	httpmock.RegisterResponder("DELETE", `=~^https://graph.microsoft.com/beta/users/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"User not found"}}`), nil
		})
}

// loadFixture loads a JSON fixture file from the tests/responses directory
func loadFixture(filename string) (map[string]any, error) {
	// Path relative to the mocks directory: ../tests/responses/
	fixturesPath := filepath.Join("tests", "responses", filename)
	data, err := os.ReadFile(fixturesPath)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// registerTestUsers registers predefined test users from JSON fixtures
func registerTestUsers() {
	// Load minimal user from fixture
	minimalUserData, err := loadFixture("user_minimal.json")
	if err != nil {
		// Fallback to inline data if fixture loading fails
		minimalUserData = map[string]any{
			"id":                "00000000-0000-0000-0000-000000000001",
			"displayName":       "Minimal User",
			"userPrincipalName": "minimal.user@deploymenttheory.com",
			"accountEnabled":    true,
			"passwordProfile": map[string]any{
				"password":                             "SecureP@ssw0rd123!",
				"forceChangePasswordNextSignIn":        false,
				"forceChangePasswordNextSignInWithMfa": false,
			},
			"createdDateTime": "2023-01-01T00:00:00Z",
			"businessPhones":  []any{},
			"identities":      []any{},
			"imAddresses":     []any{},
			"otherMails":      []any{},
			"proxyAddresses":  []any{},
		}
	}

	// Load maximal user from fixture
	maximalUserData, err := loadFixture("user_maximal.json")
	if err != nil {
		// Fallback to inline data if fixture loading fails
		maximalUserData = map[string]any{
			"id":                "00000000-0000-0000-0000-000000000002",
			"displayName":       "Maximal User",
			"userPrincipalName": "maximal.user@deploymenttheory.com",
			"accountEnabled":    true,
			"givenName":         "Maximal",
			"surname":           "User",
			"mail":              "maximal.user@deploymenttheory.com",
			"mailNickname":      "maximal.user",
			"jobTitle":          "Senior Developer",
			"department":        "Engineering",
			"companyName":       "Contoso Ltd",
			"officeLocation":    "Building A",
			"city":              "Redmond",
			"state":             "WA",
			"country":           "US",
			"postalCode":        "98052",
			"usageLocation":     "US",
			"businessPhones":    []any{"+1 425-555-0100"},
			"mobilePhone":       "+1 425-555-0101",
			"passwordProfile": map[string]any{
				"password":                             "SecureP@ssw0rd123!",
				"forceChangePasswordNextSignIn":        false,
				"forceChangePasswordNextSignInWithMfa": false,
			},
			"identities": []any{
				map[string]any{
					"signInType":       "emailAddress",
					"issuer":           "DeploymentTheory.onmicrosoft.com",
					"issuerAssignedId": "maximal.user@deploymenttheory.com",
				},
			},
			"otherMails":        []any{"maximal.user.other@deploymenttheory.com"},
			"showInAddressList": true,
			"createdDateTime":    "2023-01-01T00:00:00Z",
			"imAddresses":        []any{},
		}
	}

	minimalUserId := minimalUserData["id"].(string)
	maximalUserId := maximalUserData["id"].(string)

	// Store users in mock state
	mockState.Lock()
	mockState.users[minimalUserId] = minimalUserData
	mockState.users[maximalUserId] = maximalUserData
	mockState.Unlock()
}

// CleanupMockState clears the mock state
func (m *UserMock) CleanupMockState() {
	mockState.Lock()
	mockState.users = make(map[string]map[string]any)
	mockState.Unlock()
}
