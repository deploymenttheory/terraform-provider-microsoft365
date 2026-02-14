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
	users []map[string]any
}

func init() {
	mockState.users = make([]map[string]any, 0)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("user_list", &UserListMock{})
}

type UserListMock struct{}

var _ mocks.MockRegistrar = (*UserListMock)(nil)

func (m *UserListMock) RegisterMocks() {
	responsePath := filepath.Join("..", "tests", "responses", "get_all_users_success.json")
	mockState.Lock()
	users, err := m.loadUnitTestJson(responsePath)
	if err != nil {
		mockState.Unlock()
		panic(fmt.Sprintf("FATAL: Failed to load required mock JSON file: %v", err))
	}
	mockState.users = users
	mockState.Unlock()

	// Register GET for listing users
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/users",
		func(req *http.Request) (*http.Response, error) {
			return m.handleListRequest(req)
		})
}

func (m *UserListMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.users = make([]map[string]any, 0)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/users",
		httpmock.NewStringResponder(403, `{"error":{"code":"Forbidden","message":"Forbidden - 403"}}`))
}

func (m *UserListMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.users = make([]map[string]any, 0)
}

func (m *UserListMock) handleListRequest(req *http.Request) (*http.Response, error) {
	mockState.Lock()
	allUsers := make([]map[string]any, len(mockState.users))
	for i, user := range mockState.users {
		userCopy := m.copyUser(user)
		allUsers[i] = userCopy
	}
	mockState.Unlock()

	filteredUsers := m.applyFilters(allUsers, req.URL.Query())

	// Handle $orderby
	if orderBy := req.URL.Query().Get("$orderby"); orderBy != "" {
		filteredUsers = m.applyOrderBy(filteredUsers, orderBy)
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
	if skip >= len(filteredUsers) {
		filteredUsers = []map[string]any{}
	} else if skip > 0 {
		filteredUsers = filteredUsers[skip:]
	}

	// Apply top and generate nextLink
	var nextLink string
	if top < len(filteredUsers) {
		nextLink = fmt.Sprintf("https://graph.microsoft.com/beta/users?$skip=%d", skip+top)
		if filter := req.URL.Query().Get("$filter"); filter != "" {
			nextLink += "&$filter=" + url.QueryEscape(filter)
		}
		if orderBy := req.URL.Query().Get("$orderby"); orderBy != "" {
			nextLink += "&$orderby=" + url.QueryEscape(orderBy)
		}
		nextLink += "&$top=" + strconv.Itoa(top)
		filteredUsers = filteredUsers[:top]
	}

	response := map[string]any{
		"@odata.context": "https://graph.microsoft.com/beta/$metadata#users",
		"value":          filteredUsers,
	}

	if nextLink != "" {
		response["@odata.nextLink"] = nextLink
	}

	return httpmock.NewJsonResponse(200, response)
}

func (m *UserListMock) applyFilters(users []map[string]any, query url.Values) []map[string]any {
	filter := query.Get("$filter")
	if filter == "" {
		return users
	}

	filtered := make([]map[string]any, 0)
	for _, user := range users {
		if m.matchesFilter(user, filter) {
			filtered = append(filtered, user)
		}
	}
	return filtered
}

func (m *UserListMock) matchesFilter(user map[string]any, filter string) bool {
	filter = strings.ToLower(strings.TrimSpace(filter))

	// Handle startsWith() for displayName
	if strings.Contains(filter, "startswith(displayname,") {
		start := strings.Index(filter, "startswith(displayname,'") + len("startswith(displayname,'")
		end := strings.Index(filter[start:], "')")
		if end > 0 {
			searchTerm := filter[start : start+end]
			displayName := strings.ToLower(fmt.Sprintf("%v", user["displayName"]))
			if !strings.HasPrefix(displayName, searchTerm) {
				return false
			}
		}
	}

	// Handle startsWith() for userPrincipalName
	if strings.Contains(filter, "startswith(userprincipalname,") {
		start := strings.Index(filter, "startswith(userprincipalname,'") + len("startswith(userprincipalname,'")
		end := strings.Index(filter[start:], "')")
		if end > 0 {
			searchTerm := filter[start : start+end]
			upn := strings.ToLower(fmt.Sprintf("%v", user["userPrincipalName"]))
			if !strings.HasPrefix(upn, searchTerm) {
				return false
			}
		}
	}

	// Handle accountEnabled filter
	if strings.Contains(filter, "accountenabled eq") {
		re := regexp.MustCompile(`accountenabled\s+eq\s+(true|false)`)
		matches := re.FindStringSubmatch(filter)
		if len(matches) > 1 {
			expectedEnabled := matches[1] == "true"
			actualEnabled, _ := user["accountEnabled"].(bool)
			if actualEnabled != expectedEnabled {
				return false
			}
		}
	}

	// Handle userType filter
	if strings.Contains(filter, "usertype eq") {
		re := regexp.MustCompile(`usertype\s+eq\s+'([^']+)'`)
		matches := re.FindStringSubmatch(filter)
		if len(matches) > 1 {
			expectedType := matches[1]
			actualType := strings.ToLower(fmt.Sprintf("%v", user["userType"]))
			if actualType != expectedType {
				return false
			}
		}
	}

	return true
}

func (m *UserListMock) applyOrderBy(users []map[string]any, orderBy string) []map[string]any {
	// Simple ordering by displayName (for testing purposes)
	// In reality, would need proper sorting implementation
	return users
}

func (m *UserListMock) copyUser(user map[string]any) map[string]any {
	copy := make(map[string]any)
	for k, v := range user {
		copy[k] = v
	}
	if _, hasODataType := copy["@odata.type"]; !hasODataType {
		copy["@odata.type"] = "#microsoft.graph.user"
	}
	return copy
}

func (m *UserListMock) loadUnitTestJson(filePath string) ([]map[string]any, error) {
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
