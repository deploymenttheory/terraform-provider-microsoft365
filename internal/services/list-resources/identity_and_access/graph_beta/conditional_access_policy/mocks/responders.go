package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	policies []map[string]any
}

func init() {
	mockState.policies = make([]map[string]any, 0)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("conditional_access_policy_list", &ConditionalAccessPolicyListMock{})
}

type ConditionalAccessPolicyListMock struct{}

var _ mocks.MockRegistrar = (*ConditionalAccessPolicyListMock)(nil)

func (m *ConditionalAccessPolicyListMock) RegisterMocks() {
	responsePath := filepath.Join("..", "tests", "responses", "get_all_policies_success.json")
	mockState.Lock()
	policies, err := m.loadUnitTestJson(responsePath)
	if err != nil {
		mockState.Unlock()
		panic(fmt.Sprintf("FATAL: Failed to load required mock JSON file: %v", err))
	}
	mockState.policies = policies
	mockState.Unlock()

	// Register GET for listing Conditional Access policies
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/identity/conditionalAccess/policies",
		func(req *http.Request) (*http.Response, error) {
			return m.handleListRequest(req)
		})
}

func (m *ConditionalAccessPolicyListMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.policies = make([]map[string]any, 0)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/identity/conditionalAccess/policies",
		httpmock.NewStringResponder(403, `{"error":{"code":"Forbidden","message":"Forbidden - 403"}}`))
}

func (m *ConditionalAccessPolicyListMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.policies = make([]map[string]any, 0)
}

func (m *ConditionalAccessPolicyListMock) handleListRequest(req *http.Request) (*http.Response, error) {
	mockState.Lock()
	allPolicies := make([]map[string]any, len(mockState.policies))
	for i, policy := range mockState.policies {
		policyCopy := m.copyPolicy(policy)
		allPolicies[i] = policyCopy
	}
	mockState.Unlock()

	filteredPolicies := m.applyFilters(allPolicies, req.URL.Query())

	// Handle $orderby
	if orderBy := req.URL.Query().Get("$orderby"); orderBy != "" {
		filteredPolicies = m.applyOrderBy(filteredPolicies, orderBy)
	}

	// Handle pagination
	top := 100 // Default page size
	if topStr := req.URL.Query().Get("$top"); topStr != "" {
		if topVal, err := strconv.Atoi(topStr); err == nil && topVal > 0 {
			top = topVal
		}
	}

	skip := 0
	if skipStr := req.URL.Query().Get("$skip"); skipStr != "" {
		if skipVal, err := strconv.Atoi(skipStr); err == nil && skipVal >= 0 {
			skip = skipVal
		}
	}

	// Apply skip
	if skip >= len(filteredPolicies) {
		filteredPolicies = []map[string]any{}
	} else if skip > 0 {
		filteredPolicies = filteredPolicies[skip:]
	}

	// Apply top and generate nextLink
	var nextLink string
	if top < len(filteredPolicies) {
		nextLink = fmt.Sprintf("https://graph.microsoft.com/beta/identity/conditionalAccess/policies?$skip=%d", skip+top)
		if filter := req.URL.Query().Get("$filter"); filter != "" {
			nextLink += "&$filter=" + url.QueryEscape(filter)
		}
		if orderBy := req.URL.Query().Get("$orderby"); orderBy != "" {
			nextLink += "&$orderby=" + url.QueryEscape(orderBy)
		}
		nextLink += "&$top=" + strconv.Itoa(top)
		filteredPolicies = filteredPolicies[:top]
	}

	response := map[string]any{
		"@odata.context": "https://graph.microsoft.com/beta/$metadata#identity/conditionalAccess/policies",
		"value":          filteredPolicies,
	}

	if nextLink != "" {
		response["@odata.nextLink"] = nextLink
	}

	return httpmock.NewJsonResponse(200, response)
}

func (m *ConditionalAccessPolicyListMock) applyFilters(policies []map[string]any, query url.Values) []map[string]any {
	filter := query.Get("$filter")
	if filter == "" {
		return policies
	}

	filtered := make([]map[string]any, 0)
	for _, policy := range policies {
		if m.matchesFilter(policy, filter) {
			filtered = append(filtered, policy)
		}
	}
	return filtered
}

func (m *ConditionalAccessPolicyListMock) matchesFilter(policy map[string]any, filter string) bool {
	filter = strings.ToLower(strings.TrimSpace(filter))

	// Handle contains() for displayName
	if strings.Contains(filter, "contains(displayname,") {
		start := strings.Index(filter, "contains(displayname,'") + len("contains(displayname,'")
		end := strings.Index(filter[start:], "')")
		if end > 0 {
			searchTerm := filter[start : start+end]
			displayName := strings.ToLower(fmt.Sprintf("%v", policy["displayName"]))
			if !strings.Contains(displayName, searchTerm) {
				return false
			}
		}
	}

	// Handle state filter
	if strings.Contains(filter, "state eq") {
		// Extract state value from filter
		re := regexp.MustCompile(`state\s+eq\s+'([^']+)'`)
		matches := re.FindStringSubmatch(filter)
		if len(matches) > 1 {
			expectedState := matches[1]
			actualState := strings.ToLower(fmt.Sprintf("%v", policy["state"]))
			if actualState != expectedState {
				return false
			}
		}
	}

	return true
}

func (m *ConditionalAccessPolicyListMock) applyOrderBy(policies []map[string]any, orderBy string) []map[string]any {
	// Simple ordering by displayName (for testing purposes)
	// In reality, would need proper sorting implementation
	return policies
}

func (m *ConditionalAccessPolicyListMock) copyPolicy(policy map[string]any) map[string]any {
	copy := make(map[string]any)
	for k, v := range policy {
		copy[k] = v
	}
	if _, hasODataType := copy["@odata.type"]; !hasODataType {
		copy["@odata.type"] = "#microsoft.graph.conditionalAccessPolicy"
	}
	return copy
}

func (m *ConditionalAccessPolicyListMock) loadUnitTestJson(filePath string) ([]map[string]any, error) {
	content, err := helpers.ParseJSONFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load JSON file %s: %w", filePath, err)
	}

	var response struct {
		Value []map[string]any `json:"value"`
	}
	if err := json.Unmarshal([]byte(content), &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON from %s: %w", filePath, err)
	}

	return response.Value, nil
}
