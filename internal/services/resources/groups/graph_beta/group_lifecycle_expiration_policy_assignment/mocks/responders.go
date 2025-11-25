package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	// Map of policy ID to list of assigned group IDs
	policyAssignments map[string][]string
	// Map of group ID to its expiration date (if assigned to policy)
	groupExpirations map[string]string
	// The mock policy ID (there's only one per tenant)
	policyID string
}

func init() {
	// Initialize mockState
	mockState.policyID = ""
	mockState.policyAssignments = make(map[string][]string)
	mockState.groupExpirations = make(map[string]string)

	// Register with global registry only
	mocks.GlobalRegistry.Register("group_lifecycle_expiration_policy_assignment", &GroupLifecycleExpirationPolicyAssignmentMock{})
}

// GroupLifecycleExpirationPolicyAssignmentMock provides mock responses for policy assignment operations
type GroupLifecycleExpirationPolicyAssignmentMock struct{}

// Ensure GroupLifecycleExpirationPolicyAssignmentMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*GroupLifecycleExpirationPolicyAssignmentMock)(nil)

// RegisterMocks registers HTTP mock responses for policy assignment operations
func (m *GroupLifecycleExpirationPolicyAssignmentMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.policyID = uuid.New().String()
	mockState.policyAssignments = make(map[string][]string)
	mockState.policyAssignments[mockState.policyID] = []string{}
	mockState.groupExpirations = make(map[string]string)
	mockState.Unlock()

	// Register GET for listing policies (to get the tenant's policy)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/groupLifecyclePolicies",
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/get_policy_response.json")

			// Update the policy ID in the JSON response with the mock policy ID
			mockState.Lock()
			defer mockState.Unlock()

			var response map[string]any
			json.Unmarshal([]byte(jsonStr), &response)
			if value, ok := response["value"].([]any); ok && len(value) > 0 {
				if policy, ok := value[0].(map[string]any); ok {
					policy["id"] = mockState.policyID
				}
			}

			respBody, _ := json.Marshal(response)
			resp := httpmock.NewStringResponse(200, string(respBody))
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	// Register POST for addGroup
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/groupLifecyclePolicies/([a-fA-F0-9\-]+)/addGroup`,
		func(req *http.Request) (*http.Response, error) {
			policyID := httpmock.MustGetSubmatch(req, 1)

			var requestBody map[string]string
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_bad_request.json")
				return httpmock.NewStringResponse(400, jsonStr), nil
			}

			groupID := requestBody["groupId"]
			if groupID == "" {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_bad_request.json")
				return httpmock.NewStringResponse(400, jsonStr), nil
			}

			mockState.Lock()
			defer mockState.Unlock()

			// Check if group is already assigned
			alreadyAssigned := false
			for _, assignedGroupID := range mockState.policyAssignments[policyID] {
				if assignedGroupID == groupID {
					alreadyAssigned = true
					break
				}
			}

			if !alreadyAssigned {
				// Add group to policy
				mockState.policyAssignments[policyID] = append(mockState.policyAssignments[policyID], groupID)
				// Set expiration date (180 days from now)
				mockState.groupExpirations[groupID] = "2025-06-01T00:00:00Z"
			}

			// Return success response from JSON file
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_addgroup_success.json")
			resp := httpmock.NewStringResponse(200, jsonStr)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	// Register POST for removeGroup
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/groupLifecyclePolicies/([a-fA-F0-9\-]+)/removeGroup`,
		func(req *http.Request) (*http.Response, error) {
			policyID := httpmock.MustGetSubmatch(req, 1)

			var requestBody map[string]string
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_bad_request.json")
				return httpmock.NewStringResponse(400, jsonStr), nil
			}

			groupID := requestBody["groupId"]
			if groupID == "" {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_bad_request.json")
				return httpmock.NewStringResponse(400, jsonStr), nil
			}

			mockState.Lock()
			defer mockState.Unlock()

			// Remove group from policy
			newAssignments := []string{}
			for _, assignedGroupID := range mockState.policyAssignments[policyID] {
				if assignedGroupID != groupID {
					newAssignments = append(newAssignments, assignedGroupID)
				}
			}
			mockState.policyAssignments[policyID] = newAssignments

			// Remove expiration date
			delete(mockState.groupExpirations, groupID)

			// Return success response from JSON file
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_removegroup_success.json")
			resp := httpmock.NewStringResponse(200, jsonStr)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	// Register GET for specific group (to check if it has expiration date)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/([a-fA-F0-9\-]+)`,
		func(req *http.Request) (*http.Response, error) {
			groupID := httpmock.MustGetSubmatch(req, 1)

			// Handle special test IDs with error JSON
			if strings.Contains(groupID, "error") {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_resource_not_found.json")
				return httpmock.NewStringResponse(404, jsonStr), nil
			}

			mockState.Lock()
			_, hasExpiration := mockState.groupExpirations[groupID]
			policyID := mockState.policyID
			mockState.Unlock()

			// Load appropriate JSON file based on expiration state
			var jsonStr string
			if hasExpiration {
				jsonStr, _ = helpers.ParseJSONFile("../tests/responses/validate_create/get_group_with_expiration.json")
			} else {
				jsonStr, _ = helpers.ParseJSONFile("../tests/responses/validate_create/get_group_without_expiration.json")
			}

			// Update the group ID and policy ID in the JSON response
			var group map[string]any
			json.Unmarshal([]byte(jsonStr), &group)
			group["id"] = groupID

			// If group has expiration, update the policy ID in groupLifecyclePolicies
			if hasExpiration {
				if policies, ok := group["groupLifecyclePolicies"].([]any); ok && len(policies) > 0 {
					if policy, ok := policies[0].(map[string]any); ok {
						policy["id"] = policyID
					}
				}
			}

			respBody, _ := json.Marshal(group)
			resp := httpmock.NewStringResponse(200, string(respBody))
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
// This sets up mocks for testing error handling when a policy exists but operations fail
func (m *GroupLifecycleExpirationPolicyAssignmentMock) RegisterErrorMocks() {
	// Initialize state with a mock policy
	mockState.Lock()
	mockState.policyID = uuid.New().String()
	mockState.policyAssignments = make(map[string][]string)
	mockState.policyAssignments[mockState.policyID] = []string{}
	mockState.groupExpirations = make(map[string]string)
	mockState.Unlock()

	// Register GET for listing policies (to find the single policy) - success response with policy
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/groupLifecyclePolicies",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			defer mockState.Unlock()

			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/get_policy_response.json")

			// Update the policy ID in the JSON response with the mock policy ID
			var response map[string]any
			json.Unmarshal([]byte(jsonStr), &response)
			if value, ok := response["value"].([]any); ok && len(value) > 0 {
				if policy, ok := value[0].(map[string]any); ok {
					policy["id"] = mockState.policyID
				}
			}

			respBody, _ := json.Marshal(response)
			resp := httpmock.NewStringResponse(200, string(respBody))
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	// Mock a 400 Bad Request for addGroup
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/groupLifecyclePolicies/([a-fA-F0-9\-]+)/addGroup$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_bad_request.json")
			resp := httpmock.NewStringResponse(400, jsonStr)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	// Mock successful group GET (return valid M365 group so validation passes)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/([a-fA-F0-9\-]+)`,
		func(req *http.Request) (*http.Response, error) {
			groupID := httpmock.MustGetSubmatch(req, 1)

			// Return a valid M365 group so validation passes
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/get_group_without_expiration.json")
			var group map[string]any
			json.Unmarshal([]byte(jsonStr), &group)
			group["id"] = groupID

			respBody, _ := json.Marshal(group)
			resp := httpmock.NewStringResponse(200, string(respBody))
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})
}

// RegisterNoPolicyErrorMocks registers HTTP mock responses for when no policy exists
func (m *GroupLifecycleExpirationPolicyAssignmentMock) RegisterNoPolicyErrorMocks() {
	// Reset state to have no policy
	mockState.Lock()
	mockState.policyID = ""
	mockState.policyAssignments = make(map[string][]string)
	mockState.groupExpirations = make(map[string]string)
	mockState.Unlock()

	// Mock no policy found - empty value array
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/groupLifecyclePolicies",
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_error/error_no_policy_found.json")
			resp := httpmock.NewStringResponse(200, jsonStr)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})
}

// CleanupMockState clears the mock state
func (m *GroupLifecycleExpirationPolicyAssignmentMock) CleanupMockState() {
	mockState.Lock()
	mockState.policyID = uuid.New().String()
	mockState.policyAssignments = make(map[string][]string)
	mockState.policyAssignments[mockState.policyID] = []string{}
	mockState.groupExpirations = make(map[string]string)
	mockState.Unlock()
}
