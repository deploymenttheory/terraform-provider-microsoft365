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
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/roleManagement/directory/roleDefinitions", func(req *http.Request) (*http.Response, error) {
		response := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#roleManagement/directory/roleDefinitions",
			"value": []map[string]any{
				{
					"id":          "55555555-5555-5555-5555-555555555555",
					"displayName": "Global Administrator",
					"description": "Mock Global Administrator role",
				},
				{
					"id":          "55555555-5555-5555-5555-555555555556",
					"displayName": "Security Administrator",
					"description": "Mock Security Administrator role",
				},
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
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/servicePrincipals\?`, func(req *http.Request) (*http.Response, error) {
		// Generate a mock service principal app ID
		appId := uuid.New().String()

		response := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#servicePrincipals",
			"value": []map[string]any{
				{
					"id":          uuid.New().String(),
					"appId":       appId,
					"displayName": "Mock Service Principal",
				},
			},
		}

		return httpmock.NewJsonResponse(200, response)
	})
}
