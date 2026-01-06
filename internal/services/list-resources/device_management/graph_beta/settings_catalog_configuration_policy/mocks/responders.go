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
	mocks.GlobalRegistry.Register("settings_catalog_list", &SettingsCatalogListMock{})
}

type SettingsCatalogListMock struct{}

var _ mocks.MockRegistrar = (*SettingsCatalogListMock)(nil)

func (m *SettingsCatalogListMock) RegisterMocks() {
	responsePath := filepath.Join("..", "tests", "responses", "get_all_policies_success.json")
	mockState.Lock()
	policies, err := m.loadUnitTestJson(responsePath)
	if err != nil {
		mockState.Unlock()
		panic(fmt.Sprintf("FATAL: Failed to load required mock JSON file: %v", err))
	}
	mockState.policies = policies
	mockState.Unlock()

	// Register GET for listing configuration policies
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		func(req *http.Request) (*http.Response, error) {
			return m.handleListRequest(req)
		})

	// Register GET for policy assignments (used by is_assigned_filter)
	// Use a regex responder to match any policy ID in the URL
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[^/]+/assignments`),
		func(req *http.Request) (*http.Response, error) {
			return m.handleAssignmentsRequest(req)
		})
}

func (m *SettingsCatalogListMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.policies = make([]map[string]any, 0)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		httpmock.NewStringResponder(403, `{"error":{"code":"Forbidden","message":"Forbidden - 403"}}`))
}

func (m *SettingsCatalogListMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.policies = make([]map[string]any, 0)
}

func (m *SettingsCatalogListMock) handleAssignmentsRequest(req *http.Request) (*http.Response, error) {
	// Extract policy ID from URL path
	pathParts := strings.Split(req.URL.Path, "/")
	var policyID string
	for i, part := range pathParts {
		if part == "configurationPolicies" && i+1 < len(pathParts) {
			policyID = pathParts[i+1]
			break
		}
	}

	mockState.Lock()
	defer mockState.Unlock()

	// Find the policy and check if it's assigned
	for _, policy := range mockState.policies {
		if id, ok := policy["id"].(string); ok && id == policyID {
			// Check if policy has assignments based on isAssigned field
			if isAssigned, ok := policy["isAssigned"].(bool); ok && isAssigned {
				// Return mock assignment
				response := map[string]any{
					"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies('" + policyID + "')/assignments",
					"value": []map[string]any{
						{
							"@odata.type": "#microsoft.graph.deviceManagementConfigurationPolicyAssignment",
							"id":          "assignment-" + policyID,
							"target":      map[string]any{"@odata.type": "#microsoft.graph.allDevicesAssignmentTarget"},
						},
					},
				}
				return httpmock.NewJsonResponse(200, response)
			}
			// Return empty assignments
			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies('" + policyID + "')/assignments",
				"value":          []map[string]any{},
			}
			return httpmock.NewJsonResponse(200, response)
		}
	}

	return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"Policy not found"}}`), nil
}

func (m *SettingsCatalogListMock) handleListRequest(req *http.Request) (*http.Response, error) {
	mockState.Lock()
	allPolicies := make([]map[string]any, len(mockState.policies))
	for i, policy := range mockState.policies {
		policyCopy := m.copyPolicy(policy)
		allPolicies[i] = policyCopy
	}
	mockState.Unlock()

	filteredPolicies := m.applyFilters(allPolicies, req.URL.Query())

	// Handle $expand
	expandParam := req.URL.Query().Get("$expand")
	if strings.Contains(expandParam, "assignments") {
		// Load policies with assignments
		responsePath := filepath.Join("..", "tests", "responses", "get_policies_with_assignments_success.json")
		policiesWithAssignments, err := m.loadUnitTestJson(responsePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load assignments JSON: %w", err)
		}
		// Use policies with assignments data
		filteredPolicies = policiesWithAssignments
	}

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
		nextLink = fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/configurationPolicies?$skip=%d", skip+top)
		if filter := req.URL.Query().Get("$filter"); filter != "" {
			nextLink += "&$filter=" + url.QueryEscape(filter)
		}
		if orderBy := req.URL.Query().Get("$orderby"); orderBy != "" {
			nextLink += "&$orderby=" + url.QueryEscape(orderBy)
		}
		if expandParam != "" {
			nextLink += "&$expand=" + url.QueryEscape(expandParam)
		}
		nextLink += "&$top=" + strconv.Itoa(top)
		filteredPolicies = filteredPolicies[:top]
	}

	response := map[string]any{
		"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies",
		"value":          filteredPolicies,
	}

	if nextLink != "" {
		response["@odata.nextLink"] = nextLink
	}

	return httpmock.NewJsonResponse(200, response)
}

func (m *SettingsCatalogListMock) applyFilters(policies []map[string]any, query url.Values) []map[string]any {
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

func (m *SettingsCatalogListMock) matchesFilter(policy map[string]any, filter string) bool {
	filter = strings.ToLower(strings.TrimSpace(filter))

	// Handle contains() for name
	if strings.Contains(filter, "contains(name,") {
		start := strings.Index(filter, "contains(name,'") + len("contains(name,'")
		end := strings.Index(filter[start:], "')")
		if end > 0 {
			searchTerm := filter[start : start+end]
			name := strings.ToLower(fmt.Sprintf("%v", policy["name"]))
			if !strings.Contains(name, searchTerm) {
				return false
			}
		}
	}

	// Handle platforms filter
	if strings.Contains(filter, "platforms") {
		var platformsStr string
		switch v := policy["platforms"].(type) {
		case string:
			platformsStr = strings.ToLower(v)
		case []any:
			platformsStr = strings.ToLower(strings.Join(m.anyToStringSlice(v), ","))
		}
		if strings.Contains(filter, "windows10") && !strings.Contains(platformsStr, "windows10") {
			return false
		}
		if strings.Contains(filter, "macos") && !strings.Contains(platformsStr, "macos") {
			return false
		}
	}

	// Handle technologies filter
	if strings.Contains(filter, "technologies") {
		var techStr string
		switch v := policy["technologies"].(type) {
		case string:
			techStr = strings.ToLower(v)
		case []any:
			techStr = strings.ToLower(strings.Join(m.anyToStringSlice(v), ","))
		}
		if strings.Contains(filter, "mdm") && !strings.Contains(techStr, "mdm") {
			return false
		}
	}

	// Handle templateReference/templateFamily filter
	if strings.Contains(filter, "templatereference/templatefamily") {
		if templateRef, ok := policy["templateReference"].(map[string]any); ok {
			if templateFamily, ok := templateRef["templateFamily"].(string); ok {
				templateFamily = strings.ToLower(templateFamily)
				filterLower := strings.ToLower(filter)
				// Extract the family value from the filter (e.g., "eq 'baseline'")
				if strings.Contains(filterLower, "eq '") {
					start := strings.Index(filterLower, "eq '") + 4
					end := strings.Index(filterLower[start:], "'")
					if end > 0 {
						expectedFamily := filterLower[start : start+end]
						if templateFamily != expectedFamily {
							return false
						}
					}
				}
			}
		}
	}

	// Handle isAssigned filter
	if strings.Contains(filter, "isassigned eq true") {
		if isAssigned, ok := policy["isAssigned"].(bool); !ok || !isAssigned {
			return false
		}
	}
	if strings.Contains(filter, "isassigned eq false") {
		if isAssigned, ok := policy["isAssigned"].(bool); ok && isAssigned {
			return false
		}
	}

	// Handle roleScopeTagIds filter
	if strings.Contains(filter, "rolescopetagids/any") {
		start := strings.Index(filter, "r eq '") + len("r eq '")
		end := strings.Index(filter[start:], "'")
		if end > 0 {
			searchTag := filter[start : start+end]
			if tags, ok := policy["roleScopeTagIds"].([]any); ok {
				found := false
				for _, tag := range tags {
					if fmt.Sprintf("%v", tag) == searchTag {
						found = true
						break
					}
				}
				if !found {
					return false
				}
			}
		}
	}

	return true
}

func (m *SettingsCatalogListMock) applyOrderBy(policies []map[string]any, orderBy string) []map[string]any {
	// Simple ordering by name (for testing purposes)
	// In reality, would need proper sorting implementation
	return policies
}

func (m *SettingsCatalogListMock) anyToStringSlice(slice []any) []string {
	result := make([]string, len(slice))
	for i, v := range slice {
		result[i] = fmt.Sprintf("%v", v)
	}
	return result
}

func (m *SettingsCatalogListMock) copyPolicy(policy map[string]any) map[string]any {
	copy := make(map[string]any)
	for k, v := range policy {
		copy[k] = v
	}
	if _, hasODataType := copy["@odata.type"]; !hasODataType {
		copy["@odata.type"] = "#microsoft.graph.deviceManagementConfigurationPolicy"
	}
	return copy
}

func (m *SettingsCatalogListMock) loadUnitTestJson(filePath string) ([]map[string]any, error) {
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
