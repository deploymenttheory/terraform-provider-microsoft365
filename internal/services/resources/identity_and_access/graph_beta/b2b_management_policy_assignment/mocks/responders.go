package mocks

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/jarcoal/httpmock"
)

// Fixed directory object IDs used by the unit test configurations.
const (
	ServicePrincipalID = "00000000-0000-0000-0000-000000000020"
	ApplicationID      = "00000000-0000-0000-0000-000000000030"
	UserID             = "00000000-0000-0000-0000-000000000040"
)

var directoryObjectTypes = map[string]string{
	ServicePrincipalID: "#microsoft.graph.servicePrincipal",
	ApplicationID:      "#microsoft.graph.application",
	UserID:             "#microsoft.graph.user",
}

var (
	policyRefCollectionPattern = regexp.MustCompile(`/beta/(applications|servicePrincipals)/([0-9a-fA-F-]+)/policies/\$ref$`)
	policyRefItemPattern       = regexp.MustCompile(`/beta/(applications|servicePrincipals)/([0-9a-fA-F-]+)/policies/([0-9a-fA-F-]+)/\$ref$`)
)

var mockState struct {
	sync.Mutex
	// appliesTo maps a policy ID to the directory object IDs it applies to
	appliesTo map[string]map[string]bool
	// calls counts policy reference writes keyed by "METHOD collection", e.g. "POST servicePrincipals"
	calls map[string]int
}

func init() {
	mockState.appliesTo = make(map[string]map[string]bool)
	mockState.calls = make(map[string]int)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("b2b_management_policy_assignment", &B2bManagementPolicyAssignmentMock{})
}

// B2bManagementPolicyAssignmentMock provides mock responses for B2B management policy assignment operations
type B2bManagementPolicyAssignmentMock struct{}

var _ mocks.MockRegistrar = (*B2bManagementPolicyAssignmentMock)(nil)

// RegisterMocks registers HTTP mock responses for B2B management policy assignment operations
func (m *B2bManagementPolicyAssignmentMock) RegisterMocks() {
	// The documented appliesTo/$ref endpoints are read-only in Microsoft Graph
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/policies/b2bManagementPolicies/[0-9a-fA-F-]+/appliesTo/\$ref$`,
		factories.ErrorResponse(400, "Request_BadRequest", "Property 'appliesTo' is read-only and cannot be set."))

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/policies/b2bManagementPolicies/[0-9a-fA-F-]+$`,
		m.getReferencedPolicyResponder())
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/directoryObjects/[0-9a-fA-F-]+$`,
		m.getDirectoryObjectResponder())
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/(applications|servicePrincipals)/[0-9a-fA-F-]+/policies/\$ref$`,
		m.addPolicyReferenceResponder())
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/(applications|servicePrincipals)/[0-9a-fA-F-]+/policies/[0-9a-fA-F-]+/\$ref$`,
		m.removePolicyReferenceResponder())
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/policies/b2bManagementPolicies/[0-9a-fA-F-]+/appliesTo$`,
		m.listAppliesToResponder())
}

func (m *B2bManagementPolicyAssignmentMock) getReferencedPolicyResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		pathParts := strings.Split(req.URL.Path, "/")
		policyID := pathParts[len(pathParts)-1]

		response := map[string]any{
			"@odata.context":        "https://graph.microsoft.com/beta/$metadata#policies/b2bManagementPolicies/$entity",
			"id":                    policyID,
			"displayName":           "test-b2b-management-policy",
			"isOrganizationDefault": false,
			"definition":            []string{`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"BlockedDomains":["example.com"]}}}`},
		}
		return factories.SuccessResponse(200, response)(req)
	}
}

func (m *B2bManagementPolicyAssignmentMock) getDirectoryObjectResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		pathParts := strings.Split(req.URL.Path, "/")
		id := pathParts[len(pathParts)-1]

		odataType, ok := directoryObjectTypes[id]
		if !ok {
			return httpmock.NewStringResponse(404, `{"error":{"code":"Request_ResourceNotFound","message":"Resource '`+id+`' does not exist or one of its queried reference-property objects are not present."}}`), nil
		}

		response := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#directoryObjects/$entity",
			"@odata.type":    odataType,
			"id":             id,
		}
		return factories.SuccessResponse(200, response)(req)
	}
}

func (m *B2bManagementPolicyAssignmentMock) addPolicyReferenceResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		matches := policyRefCollectionPattern.FindStringSubmatch(req.URL.Path)
		if matches == nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid path"}}`), nil
		}
		directoryObjectID := matches[2]

		mockState.Lock()
		mockState.calls["POST "+matches[1]]++
		mockState.Unlock()

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid JSON"}}`), nil
		}

		odataID, ok := requestBody["@odata.id"].(string)
		if !ok || !strings.HasPrefix(odataID, "https://graph.microsoft.com/beta/policies/b2bManagementPolicies/") {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid @odata.id"}}`), nil
		}
		policyID := odataID[strings.LastIndex(odataID, "/")+1:]

		mockState.Lock()
		if mockState.appliesTo[policyID] == nil {
			mockState.appliesTo[policyID] = make(map[string]bool)
		}
		mockState.appliesTo[policyID][directoryObjectID] = true
		mockState.Unlock()

		return factories.EmptySuccessResponse(204)(req)
	}
}

func (m *B2bManagementPolicyAssignmentMock) removePolicyReferenceResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		matches := policyRefItemPattern.FindStringSubmatch(req.URL.Path)
		if matches == nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid path"}}`), nil
		}
		directoryObjectID, policyID := matches[2], matches[3]

		mockState.Lock()
		mockState.calls["DELETE "+matches[1]]++
		exists := mockState.appliesTo[policyID][directoryObjectID]
		delete(mockState.appliesTo[policyID], directoryObjectID)
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"Request_ResourceNotFound","message":"Resource not found"}}`), nil
		}

		return factories.EmptySuccessResponse(204)(req)
	}
}

func (m *B2bManagementPolicyAssignmentMock) listAppliesToResponder() httpmock.Responder {
	return func(req *http.Request) (*http.Response, error) {
		pathParts := strings.Split(req.URL.Path, "/")
		policyID := pathParts[len(pathParts)-2]

		mockState.Lock()
		value := make([]any, 0, len(mockState.appliesTo[policyID]))
		for directoryObjectID := range mockState.appliesTo[policyID] {
			value = append(value, map[string]any{
				"@odata.type": directoryObjectTypes[directoryObjectID],
				"id":          directoryObjectID,
			})
		}
		mockState.Unlock()

		response := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#directoryObjects",
			"value":          value,
		}
		return factories.SuccessResponse(200, response)(req)
	}
}

// CleanupMockState clears the mock state for clean test runs
func (m *B2bManagementPolicyAssignmentMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	clear(mockState.appliesTo)
	clear(mockState.calls)
}

// CallCount returns how many policy reference requests with the given method were sent to the
// given collection ("applications" or "servicePrincipals").
func (m *B2bManagementPolicyAssignmentMock) CallCount(method, collection string) int {
	mockState.Lock()
	defer mockState.Unlock()
	return mockState.calls[method+" "+collection]
}

// RegisterEventualConsistencyMocks overrides the referenced-policy and appliesTo responders to
// simulate Microsoft Entra replication lag around resource creation:
//   - the first policyGet404Count GETs of the referenced policy return 404, then delegate to the
//     normal responder (exercises Create's pre-assignment propagation wait)
//   - after a successful POST, the first staleListCount appliesTo GETs return an empty
//     collection, then delegate to the normal responder (exercises the ConsistencyPredicate)
//
// Call after RegisterMocks.
func (m *B2bManagementPolicyAssignmentMock) RegisterEventualConsistencyMocks(policyGet404Count, staleListCount int) {
	var mu sync.Mutex
	policyGetFailuresRemaining := policyGet404Count
	staleListsRemaining := staleListCount

	getPolicy := m.getReferencedPolicyResponder()
	list := m.listAppliesToResponder()

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/policies/b2bManagementPolicies/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			mu.Lock()
			if policyGetFailuresRemaining > 0 {
				policyGetFailuresRemaining--
				mu.Unlock()
				return httpmock.NewStringResponse(404, `{"error":{"code":"Directory_ObjectNotFound","message":"Unable to read the company information from the directory."}}`), nil
			}
			mu.Unlock()
			return getPolicy(req)
		})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/policies/b2bManagementPolicies/[0-9a-fA-F-]+/appliesTo$`,
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			assignmentExists := len(mockState.appliesTo) > 0
			mockState.Unlock()

			mu.Lock()
			if assignmentExists && staleListsRemaining > 0 {
				staleListsRemaining--
				mu.Unlock()
				response := map[string]any{
					"@odata.context": "https://graph.microsoft.com/beta/$metadata#directoryObjects",
					"value":          []any{},
				}
				return factories.SuccessResponse(200, response)(req)
			}
			mu.Unlock()
			return list(req)
		})
}

// RegisterErrorMocks registers mock responses that simulate error conditions
func (m *B2bManagementPolicyAssignmentMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/(applications|servicePrincipals)/[0-9a-fA-F-]+/policies/\$ref$`,
		factories.ErrorResponse(403, "Authorization_RequestDenied", "Insufficient privileges to complete the operation."))
}
