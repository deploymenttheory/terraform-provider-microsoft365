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
	scripts           []map[string]any
	filterExpectations []string
}

func init() {
	mockState.scripts = make([]map[string]any, 0)
	mockState.filterExpectations = make([]string, 0)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_platform_script_list", &WindowsPlatformScriptListMock{})
}

type WindowsPlatformScriptListMock struct{}

var _ mocks.MockRegistrar = (*WindowsPlatformScriptListMock)(nil)

func (m *WindowsPlatformScriptListMock) RegisterMocks() {
	responsePath := filepath.Join("..", "tests", "responses", "get_all_scripts_success.json")
	mockState.Lock()
	scripts, err := m.loadUnitTestJson(responsePath)
	if err != nil {
		mockState.Unlock()
		panic(fmt.Sprintf("FATAL: Failed to load required mock JSON file: %v", err))
	}
	mockState.scripts = scripts
	mockState.Unlock()

	// Register GET for listing Windows platform scripts
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceManagementScripts",
		func(req *http.Request) (*http.Response, error) {
			return m.handleListRequest(req)
		})
}

func (m *WindowsPlatformScriptListMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.scripts = make([]map[string]any, 0)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceManagementScripts",
		httpmock.NewStringResponder(403, `{"error":{"code":"Forbidden","message":"Forbidden - 403"}}`))
}

func (m *WindowsPlatformScriptListMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.scripts = make([]map[string]any, 0)
	mockState.filterExpectations = make([]string, 0)
}

func (m *WindowsPlatformScriptListMock) SetFilterExpectations(expectations []string) {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.filterExpectations = expectations
}

func (m *WindowsPlatformScriptListMock) handleListRequest(req *http.Request) (*http.Response, error) {
	mockState.Lock()
	allScripts := make([]map[string]any, len(mockState.scripts))
	for i, script := range mockState.scripts {
		scriptCopy := m.copyScript(script)
		allScripts[i] = scriptCopy
	}
	mockState.Unlock()

	filteredScripts := m.applyFilters(allScripts, req.URL.Query())

	// Handle $orderby
	if orderBy := req.URL.Query().Get("$orderby"); orderBy != "" {
		filteredScripts = m.applyOrderBy(filteredScripts, orderBy)
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
	if skip >= len(filteredScripts) {
		filteredScripts = []map[string]any{}
	} else if skip > 0 {
		filteredScripts = filteredScripts[skip:]
	}

	// Apply top and generate nextLink
	var nextLink string
	if top < len(filteredScripts) {
		nextLink = fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/deviceManagementScripts?$skip=%d", skip+top)
		if filter := req.URL.Query().Get("$filter"); filter != "" {
			nextLink += "&$filter=" + url.QueryEscape(filter)
		}
		if orderBy := req.URL.Query().Get("$orderby"); orderBy != "" {
			nextLink += "&$orderby=" + url.QueryEscape(orderBy)
		}
		nextLink += "&$top=" + strconv.Itoa(top)
		filteredScripts = filteredScripts[:top]
	}

	response := map[string]any{
		"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceManagementScripts",
		"value":          filteredScripts,
	}

	if nextLink != "" {
		response["@odata.nextLink"] = nextLink
	}

	resp, err := httpmock.NewJsonResponse(200, response)
	if err != nil {
		return nil, fmt.Errorf("failed to create mock JSON response: %w", err)
	}
	return resp, nil
}

func (m *WindowsPlatformScriptListMock) applyFilters(scripts []map[string]any, query url.Values) []map[string]any {
	filter := query.Get("$filter")
	if filter == "" {
		return scripts
	}

	filtered := make([]map[string]any, 0)
	for _, script := range scripts {
		if m.matchesFilter(script, filter) {
			filtered = append(filtered, script)
		}
	}
	return filtered
}

func (m *WindowsPlatformScriptListMock) matchesFilter(script map[string]any, filter string) bool {
	filter = strings.ToLower(strings.TrimSpace(filter))

	// Handle contains() for displayName
	if strings.Contains(filter, "contains(displayname,") {
		start := strings.Index(filter, "contains(displayname,'") + len("contains(displayname,'")
		end := strings.Index(filter[start:], "')")
		if end > 0 {
			searchTerm := filter[start : start+end]
			displayName := strings.ToLower(fmt.Sprintf("%v", script["displayName"]))
			if !strings.Contains(displayName, searchTerm) {
				return false
			}
		}
	}

	// Handle contains() for fileName
	if strings.Contains(filter, "contains(filename,") {
		start := strings.Index(filter, "contains(filename,'") + len("contains(filename,'")
		end := strings.Index(filter[start:], "')")
		if end > 0 {
			searchTerm := filter[start : start+end]
			fileName := strings.ToLower(fmt.Sprintf("%v", script["fileName"]))
			if !strings.Contains(fileName, searchTerm) {
				return false
			}
		}
	}

	// Handle runAsAccount eq
	if strings.Contains(filter, "runasaccount eq") {
		re := regexp.MustCompile(`runasaccount\s+eq\s+'([^']+)'`)
		matches := re.FindStringSubmatch(filter)
		if len(matches) > 1 {
			expectedAccount := matches[1]
			actualAccount := strings.ToLower(fmt.Sprintf("%v", script["runAsAccount"]))
			if actualAccount != expectedAccount {
				return false
			}
		}
	}

	return true
}

func (m *WindowsPlatformScriptListMock) applyOrderBy(scripts []map[string]any, orderBy string) []map[string]any {
	// Simple ordering implementation (for testing purposes)
	return scripts
}

func (m *WindowsPlatformScriptListMock) copyScript(script map[string]any) map[string]any {
	copy := make(map[string]any)
	for k, v := range script {
		copy[k] = v
	}
	if _, hasODataType := copy["@odata.type"]; !hasODataType {
		copy["@odata.type"] = "#microsoft.graph.deviceManagementScript"
	}
	return copy
}

func (m *WindowsPlatformScriptListMock) loadUnitTestJson(filePath string) ([]map[string]any, error) {
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
