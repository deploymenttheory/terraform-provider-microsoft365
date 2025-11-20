package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	policies map[string]map[string]any
}

func init() {
	mockState.policies = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("conditional_access_policy", &ConditionalAccessPolicyMock{})
}

type ConditionalAccessPolicyMock struct{}

var _ mocks.MockRegistrar = (*ConditionalAccessPolicyMock)(nil)

// getJSONFileForDisplayName determines which JSON file to load based on the policy's display name
func getJSONFileForDisplayName(displayName string) string {
	// Extract policy ID prefix (e.g., "CAD001", "CAP002", "CAL003", "CAU004")
	var policyPrefix string
	displayNameUpper := strings.ToUpper(displayName)

	switch {
	case strings.HasPrefix(displayNameUpper, "CAD"):
		policyPrefix = strings.ToLower(displayNameUpper[:6]) // e.g., "cad001"
	case strings.HasPrefix(displayNameUpper, "CAP"):
		policyPrefix = strings.ToLower(displayNameUpper[:6]) // e.g., "cap002"
	case strings.HasPrefix(displayNameUpper, "CAL"):
		policyPrefix = strings.ToLower(displayNameUpper[:6]) // e.g., "cal003"
	case strings.HasPrefix(displayNameUpper, "CAU"):
		// Handle CAU001A specially
		if strings.HasPrefix(displayNameUpper, "CAU001A") {
			policyPrefix = "cau001a"
		} else {
			policyPrefix = strings.ToLower(displayNameUpper[:6]) // e.g., "cau004"
		}
	default:
		return ""
	}

	return fmt.Sprintf("post_conditional_access_policy_%s_success.json", policyPrefix)
}

// deduplicateArraysInResponse recursively deduplicates arrays in a map to match Terraform's Set behavior
func deduplicateArraysInResponse(obj map[string]any) {
	for key, value := range obj {
		switch v := value.(type) {
		case []any:
			// Deduplicate the array
			deduplicated := deduplicateArray(v)
			// Keep empty arrays as empty arrays to match Terraform's Set behavior
			obj[key] = deduplicated
		case map[string]any:
			// Recursively process nested objects
			deduplicateArraysInResponse(v)
		}
	}
}

// deduplicateArray removes duplicate elements from an array
func deduplicateArray(arr []any) []any {
	seen := make(map[string]bool)
	// Initialize as empty slice (not nil) so it marshals to [] instead of null
	result := make([]any, 0)

	for _, item := range arr {
		// Convert item to string for comparison
		var key string
		switch v := item.(type) {
		case string:
			key = v
		case map[string]any:
			// For objects, recursively deduplicate their arrays too
			deduplicateArraysInResponse(v)
			// For deduplication, we'll keep all objects as they likely have different content
			result = append(result, v)
			continue
		default:
			// For other types, convert to JSON string for comparison
			jsonBytes, _ := json.Marshal(v)
			key = string(jsonBytes)
		}

		if !seen[key] {
			seen[key] = true
			result = append(result, item)
		}
	}

	return result
}

func (m *ConditionalAccessPolicyMock) RegisterMocks() {
	mockState.Lock()
	mockState.policies = make(map[string]map[string]any)
	mockState.Unlock()

	// Register mock dependencies for unit tests
	m.registerMockGroups()
	m.registerMockNamedLocations()
	m.registerMockRoleDefinitions()
	m.registerMockServicePrincipals()

	// Create conditional access policy - POST /identity/conditionalAccess/policies
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/identity/conditionalAccess/policies", func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		// Determine which JSON file to load based on displayName
		displayName, ok := requestBody["displayName"].(string)
		if !ok {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"displayName is required"}}`), nil
		}

		jsonFileName := getJSONFileForDisplayName(displayName)
		if jsonFileName == "" {
			return httpmock.NewStringResponse(400, fmt.Sprintf(`{"error":{"code":"BadRequest","message":"Unable to determine JSON file for displayName: %s"}}`, displayName)), nil
		}

		// Load JSON response from file
		responsesPath := filepath.Join("tests", "responses", "validate_create", jsonFileName)
		jsonData, err := os.ReadFile(responsesPath)
		if err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to load JSON response file: %s"}}`, err.Error())), nil
		}

		// Parse the JSON response
		var responseObj map[string]any
		if err := json.Unmarshal(jsonData, &responseObj); err != nil {
			return httpmock.NewStringResponse(500, fmt.Sprintf(`{"error":{"code":"InternalServerError","message":"Failed to parse JSON response: %s"}}`, err.Error())), nil
		}

		// Generate a UUID for the new resource and update the response
		newId := uuid.New().String()
		responseObj["id"] = newId

		// Deduplicate arrays in the response to match Terraform's Set behavior
		deduplicateArraysInResponse(responseObj)

		// Store in mock state
		mockState.Lock()
		mockState.policies[newId] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// Get conditional access policy - GET /identity/conditionalAccess/policies/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/policies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		policyId := parts[len(parts)-1]

		mockState.Lock()
		policy, exists := mockState.policies[policyId]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		// Deduplicate arrays before returning
		deduplicateArraysInResponse(policy)

		// Return the stored policy data
		return httpmock.NewJsonResponse(200, policy)
	})

	// Update conditional access policy - PATCH /identity/conditionalAccess/policies/{id}
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/policies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		policyId := parts[len(parts)-1]

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		policy, exists := mockState.policies[policyId]
		if !exists {
			mockState.Unlock()
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		// Update fields from request
		for key, value := range requestBody {
			policy[key] = value
		}
		policy["modifiedDateTime"] = "2024-01-02T00:00:00Z"
		mockState.policies[policyId] = policy
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Delete conditional access policy - DELETE /identity/conditionalAccess/policies/{id}
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/policies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		policyId := parts[len(parts)-1]

		mockState.Lock()
		_, exists := mockState.policies[policyId]
		if exists {
			delete(mockState.policies, policyId)
		}
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *ConditionalAccessPolicyMock) RegisterErrorMocks() {
	// Error scenarios for testing
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/identity/conditionalAccess/policies", httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/policies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/policies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/policies/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

func (m *ConditionalAccessPolicyMock) CleanupMockState() {
	mockState.Lock()
	mockState.policies = make(map[string]map[string]any)
	mockState.Unlock()
}

// registerMockGroups registers mock group resources for unit tests
func (m *ConditionalAccessPolicyMock) registerMockGroups() {
	// Mock group creation - returns a UUID for any group name
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/groups", func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		newId := uuid.New().String()
		response := map[string]any{
			"id":          newId,
			"displayName": requestBody["displayName"],
		}

		return httpmock.NewJsonResponse(201, response)
	})

	// Mock group GET - return a valid response for any group ID
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		groupId := parts[len(parts)-1]

		response := map[string]any{
			"id":          groupId,
			"displayName": "Mock Group",
		}

		return httpmock.NewJsonResponse(200, response)
	})
}

// registerMockNamedLocations registers mock named location resources for unit tests
func (m *ConditionalAccessPolicyMock) registerMockNamedLocations() {
	// Mock named location creation - returns a UUID for any named location
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/identity/conditionalAccess/namedLocations", func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		newId := uuid.New().String()
		response := map[string]any{
			"id":          newId,
			"displayName": requestBody["displayName"],
			"@odata.type": requestBody["@odata.type"],
		}

		// Copy relevant fields based on type
		if odataType, ok := requestBody["@odata.type"].(string); ok {
			if odataType == "#microsoft.graph.ipNamedLocation" {
				if isTrusted, ok := requestBody["isTrusted"]; ok {
					response["isTrusted"] = isTrusted
				}
				if ipRanges, ok := requestBody["ipRanges"]; ok {
					response["ipRanges"] = ipRanges
				}
			} else if odataType == "#microsoft.graph.countryNamedLocation" {
				if countriesAndRegions, ok := requestBody["countriesAndRegions"]; ok {
					response["countriesAndRegions"] = countriesAndRegions
				}
				if countryLookupMethod, ok := requestBody["countryLookupMethod"]; ok {
					response["countryLookupMethod"] = countryLookupMethod
				}
			}
		}

		return httpmock.NewJsonResponse(201, response)
	})

	// Mock named location LIST - return mock named locations for validation
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/conditionalAccess/namedLocations`, func(req *http.Request) (*http.Response, error) {
		response := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#conditionalAccess/namedLocations",
			"value": []map[string]any{
				{
					"id":          "44444444-4444-4444-4444-444444444444",
					"displayName": "Mock Named Location 1",
					"@odata.type": "#microsoft.graph.ipNamedLocation",
					"isTrusted":   true,
				},
				{
					"id":          "55555555-5555-5555-5555-555555555555",
					"displayName": "Mock Named Location 2",
					"@odata.type": "#microsoft.graph.ipNamedLocation",
					"isTrusted":   false,
				},
			},
		}

		return httpmock.NewJsonResponse(200, response)
	})

	// Mock named location GET - return a valid response for any named location ID
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/identity/conditionalAccess/namedLocations/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		locationId := parts[len(parts)-1]

		response := map[string]any{
			"id":          locationId,
			"displayName": "Mock Named Location",
			"@odata.type": "#microsoft.graph.ipNamedLocation",
			"isTrusted":   false,
			"ipRanges": []map[string]any{
				{
					"@odata.type": "#microsoft.graph.iPv4CidrRange",
					"cidrAddress": "192.168.1.0/24",
				},
			},
		}

		return httpmock.NewJsonResponse(200, response)
	})
}

// registerMockRoleDefinitions registers mock role definition resources for unit tests
func (m *ConditionalAccessPolicyMock) registerMockRoleDefinitions() {
	// Mock role definitions list - returns role definitions for validation (no query params)
	// These are real Azure AD built-in role template IDs
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/roleManagement/directory/roleDefinitions", func(req *http.Request) (*http.Response, error) {
		response := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#roleManagement/directory/roleDefinitions",
			"value": []map[string]any{
				{"id": "9b895d92-2cd3-44c7-9d02-a6ac2d5ea5c3", "displayName": "Application Administrator", "description": "Mock role"},
				{"id": "cf1c38e5-3621-4004-a7cb-879624dced7c", "displayName": "Application Developer", "description": "Mock role"},
				{"id": "c4e39bd9-1100-46d3-8c65-fb160da0071f", "displayName": "Attack Payload Author", "description": "Mock role"},
				{"id": "25a516ed-2fa0-40ea-a2d0-12923a21473a", "displayName": "Attack Simulation Administrator", "description": "Mock role"},
				{"id": "aaf43236-0c0d-4d5f-883a-6955382ac081", "displayName": "Attribute Assignment Administrator", "description": "Mock role"},
				{"id": "b0f54661-2d74-4c50-afa3-1ec803f12efe", "displayName": "Attribute Assignment Reader", "description": "Mock role"},
				{"id": "158c047a-c907-4556-b7ef-446551a6b5f7", "displayName": "Cloud Application Administrator", "description": "Mock role"},
				{"id": "7698a772-787b-4ac8-901f-60d6b08affd2", "displayName": "Cloud App Security Administrator", "description": "Mock role"},
				{"id": "17315797-102d-40b4-93e0-432062caca18", "displayName": "Compliance Administrator", "description": "Mock role"},
				{"id": "b1be1c3e-b65d-4f19-8427-f6fa0d97feb9", "displayName": "Conditional Access Administrator", "description": "Mock role"},
				{"id": "9360feb5-f418-4baa-8175-e2a00bac4301", "displayName": "Directory Writers", "description": "Mock role"},
				{"id": "29232cdf-9323-42fd-ade2-1d097af3e4de", "displayName": "Exchange Administrator", "description": "Mock role"},
				{"id": "f2ef992c-3afb-46b9-b7cf-a126ee74c451", "displayName": "Global Administrator", "description": "Mock role"},
				{"id": "62e90394-69f5-4237-9190-012177145e10", "displayName": "Global Reader", "description": "Mock role"},
				{"id": "729827e3-9c14-49f7-bb1b-9608f156bbb8", "displayName": "Helpdesk Administrator", "description": "Mock role"},
				{"id": "8ac3fc64-6eca-42ea-9e69-59f4c7b60eb2", "displayName": "Hybrid Identity Administrator", "description": "Mock role"},
				{"id": "3a2c62db-5318-420d-8d74-23affee5d9d5", "displayName": "Intune Administrator", "description": "Mock role"},
				{"id": "744ec460-397e-42ad-a462-8b3f9747a02c", "displayName": "License Administrator", "description": "Mock role"},
				{"id": "966707d0-3269-4727-9be2-8c3a10f19b9d", "displayName": "Password Administrator", "description": "Mock role"},
				{"id": "7be44c8a-adaf-4e2a-84d6-ab2649e08a13", "displayName": "Privileged Authentication Administrator", "description": "Mock role"},
				{"id": "e8611ab8-c189-46e8-94e1-60213ab1f814", "displayName": "Privileged Role Administrator", "description": "Mock role"},
				{"id": "194ae4cb-b126-40b2-bd5b-6091b380977d", "displayName": "Security Administrator", "description": "Mock role"},
				{"id": "5f2222b1-57c3-48ba-8ad5-d4759f1fde6f", "displayName": "Security Operator", "description": "Mock role"},
				{"id": "5d6b6bb7-de71-4623-b4af-96380a352509", "displayName": "Security Reader", "description": "Mock role"},
				{"id": "f28a1f50-f6e7-4571-818b-6a12f2af6b6c", "displayName": "SharePoint Administrator", "description": "Mock role"},
				{"id": "69091246-20e8-4a56-aa4d-066075b2a7a8", "displayName": "Teams Administrator", "description": "Mock role"},
				{"id": "fe930be7-5e62-47db-91af-98c3a49a38b1", "displayName": "User Administrator", "description": "Mock role"},
			},
		}

		return httpmock.NewJsonResponse(200, response)
	})

	// Mock role definitions list - returns role definitions for data sources (with query params)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/roleManagement/directory/roleDefinitions\?`, func(req *http.Request) (*http.Response, error) {
		// Generate a mock role definition ID
		roleId := uuid.New().String()

		response := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#roleManagement/directory/roleDefinitions",
			"value": []map[string]any{
				{
					"id":          roleId,
					"displayName": "Mock Role",
					"description": "Mock role for unit tests",
				},
			},
		}

		return httpmock.NewJsonResponse(200, response)
	})
}

// registerMockServicePrincipals registers mock service principal resources for unit tests
func (m *ConditionalAccessPolicyMock) registerMockServicePrincipals() {
	// Mock service principals list - returns service principals for data sources
	// These are real Microsoft service principal app IDs
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/servicePrincipals\?`, func(req *http.Request) (*http.Response, error) {
		response := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#servicePrincipals",
			"value": []map[string]any{
				{"id": uuid.New().String(), "appId": "9cdead84-a844-4324-93f2-b2e6bb768d07", "displayName": "Windows 365"},
				{"id": uuid.New().String(), "appId": "0af06dc6-e4b5-4f28-818e-e78e62d137a5", "displayName": "Microsoft SharePoint"},
				{"id": uuid.New().String(), "appId": "270efc09-cd0d-444b-a71f-39af4910ec45", "displayName": "Windows Cloud PC"},
				{"id": uuid.New().String(), "appId": "00000002-0000-0ff1-ce00-000000000000", "displayName": "Office 365 Exchange Online"},
				{"id": uuid.New().String(), "appId": "00000003-0000-0ff1-ce00-000000000000", "displayName": "Microsoft Office 365 Portal"},
				{"id": uuid.New().String(), "appId": "a4a365df-50f1-4397-bc59-1a1564b8bb9c", "displayName": "Windows Cloud Login"},
			},
		}

		return httpmock.NewJsonResponse(200, response)
	})
}
