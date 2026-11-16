package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const (
	// IntunePolicyID is the policy ID of the Microsoft Intune mobility management application
	IntunePolicyID = "0000000a-0000-0000-c000-000000000000"

	deviceCollection = "mobileDeviceManagementPolicies"
	appCollection    = "mobileAppManagementPolicies"

	notFoundBody      = `{"error":{"code":"Request_ResourceNotFound","message":"Resource does not exist or one of its queried reference-property objects are not present."}}`
	forbiddenBody     = `{"error":{"code":"Authorization_RequestDenied","message":"Insufficient privileges to complete the operation."}}`
	badRequestBody    = `{"error":{"code":"Request_BadRequest","message":"Invalid request"}}`
	appliesToBodyFmt  = `{"error":{"code":"Request_BadRequest","message":"%s"}}`
	duplicateRefBody  = `{"error":{"code":"Request_BadRequest","message":"One or more added object references already exist for 'includedGroups'."}}`
	groupsPageSize    = 1
	odataContextGroup = "https://graph.microsoft.com/beta/$metadata#groups(id)"
)

// policyURLPattern matches /policies/{collection}/{id} and captures the collection and policy ID
const policyURLPattern = `^https://graph\.microsoft\.com/beta/policies/(mobileDeviceManagementPolicies|mobileAppManagementPolicies)/([0-9a-fA-F-]{36})`

var (
	policyItemURL   = regexp.MustCompile(policyURLPattern + `(\?.*)?$`)
	includedListURL = regexp.MustCompile(policyURLPattern + `/includedGroups(\?.*)?$`)
	includedRefURL  = regexp.MustCompile(policyURLPattern + `/includedGroups/(\$|%24)ref(\?.*)?$`)
	groupRefURL     = regexp.MustCompile(policyURLPattern + `/includedGroups/([0-9a-fA-F-]{36})/(\$|%24)ref(\?.*)?$`)
)

// policyState is the mock tenant copy of one mobility management policy
type policyState struct {
	body   map[string]any
	groups []string
}

var mockState struct {
	sync.Mutex
	policies   map[string]*policyState
	refPosts   int
	refDeletes int
	patches    []map[string]any
}

func init() {
	resetState()
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, notFoundBody))
	mocks.GlobalRegistry.Register("identity_and_access_mobility_management_policy", &MobilityManagementPolicyMock{})
}

type MobilityManagementPolicyMock struct{}

var _ mocks.MockRegistrar = (*MobilityManagementPolicyMock)(nil)

func policyKey(collection, id string) string {
	return collection + "/" + strings.ToLower(id)
}

// loadResponse reads a JSON response fixture from tests/responses/validate_get
func loadResponse(fileName string) map[string]any {
	jsonData, err := os.ReadFile(filepath.Join("tests", "responses", "validate_get", fileName))
	if err != nil {
		panic(fmt.Sprintf("read fixture %s: %v", fileName, err))
	}
	var responseObj map[string]any
	if err := json.Unmarshal(jsonData, &responseObj); err != nil {
		panic(fmt.Sprintf("parse fixture %s: %v", fileName, err))
	}
	return responseObj
}

// resetState seeds the Microsoft Intune MDM and MAM policies, which the service creates in every Intune tenant
func resetState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.policies = map[string]*policyState{}
	mockState.refPosts = 0
	mockState.refDeletes = 0
	mockState.patches = nil
}

func seedPolicies() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.policies[policyKey(deviceCollection, IntunePolicyID)] = &policyState{body: loadResponse("get_mobile_device_management_policy_intune.json"), groups: []string{}}
	mockState.policies[policyKey(appCollection, IntunePolicyID)] = &policyState{body: loadResponse("get_mobile_app_management_policy_intune.json"), groups: []string{}}
}

// RefPostCount returns how many includedGroups $ref POST requests were received
func RefPostCount() int {
	mockState.Lock()
	defer mockState.Unlock()
	return mockState.refPosts
}

// RefDeleteCount returns how many includedGroups $ref DELETE requests were received
func RefDeleteCount() int {
	mockState.Lock()
	defer mockState.Unlock()
	return mockState.refDeletes
}

// PatchBodies returns every PATCH body received, in order
func PatchBodies() []map[string]any {
	mockState.Lock()
	defer mockState.Unlock()
	return slices.Clone(mockState.patches)
}

func jsonError(status int, body string) (*http.Response, error) {
	return httpmock.NewStringResponse(status, body), nil
}

func (m *MobilityManagementPolicyMock) RegisterMocks() {
	resetState()
	seedPolicies()

	// GET /policies/{collection}/{id}
	httpmock.RegisterRegexpResponder("GET", policyItemURL, func(req *http.Request) (*http.Response, error) {
		match := policyItemURL.FindStringSubmatch(req.URL.String())
		mockState.Lock()
		defer mockState.Unlock()
		policy, ok := mockState.policies[policyKey(match[1], match[2])]
		if !ok {
			return jsonError(404, notFoundBody)
		}
		return httpmock.NewJsonResponse(200, policy.body)
	})

	// PATCH /policies/{collection}/{id}
	// Mirrors Graph: appliesTo cannot be sent with other properties and selected cannot be written directly.
	httpmock.RegisterRegexpResponder("PATCH", policyItemURL, func(req *http.Request) (*http.Response, error) {
		match := policyItemURL.FindStringSubmatch(req.URL.String())
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return jsonError(400, badRequestBody)
		}

		mockState.Lock()
		defer mockState.Unlock()
		mockState.patches = append(mockState.patches, requestBody)

		policy, ok := mockState.policies[policyKey(match[1], match[2])]
		if !ok {
			return jsonError(404, notFoundBody)
		}

		if appliesTo, hasAppliesTo := requestBody["appliesTo"]; hasAppliesTo {
			if len(requestBody) > 1 {
				return jsonError(400, fmt.Sprintf(appliesToBodyFmt, "appliesTo cannot be updated together with other properties."))
			}
			if appliesTo == "selected" {
				return jsonError(400, fmt.Sprintf(appliesToBodyFmt, "appliesTo selected is set by adding included groups."))
			}
			policy.groups = []string{}
		}
		if _, ok := requestBody["isMdmEnrollmentDuringRegistrationDisabled"]; ok && match[1] == appCollection {
			return jsonError(400, badRequestBody)
		}

		for key, value := range requestBody {
			policy.body[key] = value
		}

		return httpmock.NewStringResponse(204, ""), nil
	})

	// GET /policies/{collection}/{id}/includedGroups - one group per page so paging is exercised
	httpmock.RegisterRegexpResponder("GET", includedListURL, func(req *http.Request) (*http.Response, error) {
		match := includedListURL.FindStringSubmatch(req.URL.String())
		mockState.Lock()
		defer mockState.Unlock()
		policy, ok := mockState.policies[policyKey(match[1], match[2])]
		if !ok {
			return jsonError(404, notFoundBody)
		}

		start := 0
		if skip := req.URL.Query().Get("$skiptoken"); skip != "" {
			if _, err := fmt.Sscanf(skip, "%d", &start); err != nil {
				return jsonError(400, badRequestBody)
			}
		}
		end := min(start+groupsPageSize, len(policy.groups))

		values := make([]map[string]any, 0)
		for _, id := range policy.groups[start:end] {
			values = append(values, map[string]any{"id": id})
		}
		response := map[string]any{"@odata.context": odataContextGroup, "value": values}
		if end < len(policy.groups) {
			response["@odata.nextLink"] = fmt.Sprintf("https://graph.microsoft.com/beta/policies/%s/%s/includedGroups?$select=id&$skiptoken=%d", match[1], match[2], end)
		}
		return httpmock.NewJsonResponse(200, response)
	})

	// POST /policies/{collection}/{id}/includedGroups/$ref - adding a group sets appliesTo to selected
	httpmock.RegisterRegexpResponder("POST", includedRefURL, func(req *http.Request) (*http.Response, error) {
		match := includedRefURL.FindStringSubmatch(req.URL.String())
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return jsonError(400, badRequestBody)
		}
		odataID, _ := requestBody["@odata.id"].(string)
		groupID := odataID[strings.LastIndex(odataID, "/")+1:]
		if !strings.HasPrefix(odataID, "https://graph.microsoft.com/") || len(groupID) != 36 {
			return jsonError(400, badRequestBody)
		}

		mockState.Lock()
		defer mockState.Unlock()
		mockState.refPosts++
		policy, ok := mockState.policies[policyKey(match[1], match[2])]
		if !ok {
			return jsonError(404, notFoundBody)
		}
		if slices.Contains(policy.groups, groupID) {
			return jsonError(400, duplicateRefBody)
		}
		policy.groups = append(policy.groups, groupID)
		policy.body["appliesTo"] = "selected"
		return httpmock.NewStringResponse(204, ""), nil
	})

	// DELETE /policies/{collection}/{id}/includedGroups/{groupId}/$ref
	httpmock.RegisterRegexpResponder("DELETE", groupRefURL, func(req *http.Request) (*http.Response, error) {
		match := groupRefURL.FindStringSubmatch(req.URL.String())
		mockState.Lock()
		defer mockState.Unlock()
		mockState.refDeletes++
		policy, ok := mockState.policies[policyKey(match[1], match[2])]
		if !ok {
			return jsonError(404, notFoundBody)
		}
		index := slices.Index(policy.groups, match[3])
		if index < 0 {
			return jsonError(404, notFoundBody)
		}
		policy.groups = slices.Delete(policy.groups, index, index+1)
		return httpmock.NewStringResponse(204, ""), nil
	})
}

// RegisterErrorMocks registers a 403 on every read and a 400 on every write
func (m *MobilityManagementPolicyMock) RegisterErrorMocks() {
	resetState()

	httpmock.RegisterRegexpResponder("GET", policyItemURL, httpmock.NewStringResponder(403, forbiddenBody))
	httpmock.RegisterRegexpResponder("PATCH", policyItemURL, httpmock.NewStringResponder(400, badRequestBody))
	httpmock.RegisterRegexpResponder("GET", includedListURL, httpmock.NewStringResponder(403, forbiddenBody))
	httpmock.RegisterRegexpResponder("POST", includedRefURL, httpmock.NewStringResponder(400, badRequestBody))
	httpmock.RegisterRegexpResponder("DELETE", groupRefURL, httpmock.NewStringResponder(400, badRequestBody))
}

// RegisterBadRequestOnWriteMocks registers working reads with a 400 on every PATCH
func (m *MobilityManagementPolicyMock) RegisterBadRequestOnWriteMocks() {
	m.RegisterMocks()

	httpmock.RegisterRegexpResponder("PATCH", policyItemURL, httpmock.NewStringResponder(400, badRequestBody))
}

// RegisterNotFoundMocks registers a 404 for a policy ID that does not exist
func (m *MobilityManagementPolicyMock) RegisterNotFoundMocks() {
	resetState()

	httpmock.RegisterRegexpResponder("GET", policyItemURL, httpmock.NewStringResponder(404, notFoundBody))
	httpmock.RegisterRegexpResponder("PATCH", policyItemURL, httpmock.NewStringResponder(404, notFoundBody))
	httpmock.RegisterRegexpResponder("GET", includedListURL, httpmock.NewStringResponder(404, notFoundBody))
}

func (m *MobilityManagementPolicyMock) CleanupMockState() {
	resetState()
}
