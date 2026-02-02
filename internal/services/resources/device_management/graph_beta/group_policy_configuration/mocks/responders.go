package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

type GroupPolicyConfigurationMock struct {
	state   map[string]any
	stateMu sync.RWMutex
}

func init() {
	mocks.GlobalRegistry.Register("group_policy_configuration", &GroupPolicyConfigurationMock{})
}

func (m *GroupPolicyConfigurationMock) RegisterMocks() {
	m.state = make(map[string]any)

	// POST - Create Group Policy Configuration
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations",
		m.createGroupPolicyConfiguration,
	)

	// POST - Assign Group Policy Configuration
	httpmock.RegisterRegexpResponder("POST", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/[0-9a-fA-F-]+/assign`),
		m.assignGroupPolicyConfiguration,
	)

	// GET - Read Group Policy Configuration
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/[0-9a-fA-F-]+$`),
		m.getGroupPolicyConfiguration,
	)

	// GET - Read Group Policy Configuration Assignments
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/[0-9a-fA-F-]+/assignments`),
		m.getGroupPolicyConfigurationAssignments,
	)

	// PATCH - Update Group Policy Configuration
	httpmock.RegisterRegexpResponder("PATCH", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/[0-9a-fA-F-]+$`),
		m.updateGroupPolicyConfiguration,
	)

	// DELETE - Delete Group Policy Configuration
	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/[0-9a-fA-F-]+$`),
		m.deleteGroupPolicyConfiguration,
	)
}

func (m *GroupPolicyConfigurationMock) RegisterErrorMocks() {
	// Error response for all operations
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/groupPolicyConfigurations",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(400, map[string]any{
				"error": map[string]any{
					"code":    "BadRequest",
					"message": "Invalid request body",
				},
			})
		},
	)
}

func (m *GroupPolicyConfigurationMock) CleanupMockState() {
	m.stateMu.Lock()
	defer m.stateMu.Unlock()
	m.state = make(map[string]any)
}

func (m *GroupPolicyConfigurationMock) createGroupPolicyConfiguration(req *http.Request) (*http.Response, error) {
	var requestBody map[string]any
	if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
		return httpmock.NewJsonResponse(400, map[string]any{
			"error": map[string]any{
				"code":    "BadRequest",
				"message": "Invalid request body",
			},
		})
	}

	id := uuid.New().String()

	// Determine which scenario file to load based on displayName
	scenarioFile := determineCreateScenario(requestBody)
	jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_create", scenarioFile))
	if err != nil {
		return httpmock.NewJsonResponse(500, map[string]any{
			"error": map[string]any{
				"code":    "InternalServerError",
				"message": fmt.Sprintf("Failed to load response fixture: %v", err),
			},
		})
	}

	var response map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		return httpmock.NewJsonResponse(500, map[string]any{
			"error": map[string]any{
				"code":    "InternalServerError",
				"message": fmt.Sprintf("Failed to parse response fixture: %v", err),
			},
		})
	}

	response["id"] = id

	m.stateMu.Lock()
	m.state[id] = response
	if assignments, ok := requestBody["assignments"]; ok {
		m.state[id+"_assignments"] = assignments
	}
	m.stateMu.Unlock()

	return httpmock.NewJsonResponse(201, response)
}

func determineCreateScenario(body map[string]any) string {
	displayName, _ := body["displayName"].(string)

	switch {
	case strings.Contains(displayName, "unit-test-003-minimal-assignment"):
		return "post_test_003_minimal_assignment.json"
	case strings.Contains(displayName, "unit-test-004-maximal-assignment"):
		return "post_test_004_maximal_assignment.json"
	case strings.Contains(displayName, "unit-test-005-lifecycle-maximal"):
		return "post_test_006_lifecycle_maximal.json"
	case strings.Contains(displayName, "unit-test-002-maximal"):
		return "post_test_002_maximal.json"
	case strings.Contains(displayName, "unit-test-001-minimal"):
		return "post_test_001_minimal.json"
	default:
		return "post_group_policy_configuration_success.json"
	}
}

func (m *GroupPolicyConfigurationMock) assignGroupPolicyConfiguration(req *http.Request) (*http.Response, error) {
	pathParts := strings.Split(req.URL.Path, "/")
	id := pathParts[len(pathParts)-2]

	var requestBody map[string]any
	if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
		return httpmock.NewJsonResponse(400, map[string]any{
			"error": map[string]any{
				"code":    "BadRequest",
				"message": "Invalid request body",
			},
		})
	}

	m.stateMu.Lock()
	if assignments, ok := requestBody["assignments"].([]any); ok {
		m.state[id+"_assignments"] = assignments
	}
	m.stateMu.Unlock()

	response := map[string]any{
		"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyConfigurations('" + id + "')/Microsoft.Graph.assign",
		"value":          []any{},
	}

	return httpmock.NewJsonResponse(200, response)
}

func (m *GroupPolicyConfigurationMock) getGroupPolicyConfiguration(req *http.Request) (*http.Response, error) {
	pathParts := strings.Split(req.URL.Path, "/")
	id := pathParts[len(pathParts)-1]

	m.stateMu.RLock()
	config, exists := m.state[id]
	m.stateMu.RUnlock()

	if !exists {
		return httpmock.NewJsonResponse(404, map[string]any{
			"error": map[string]any{
				"code":    "NotFound",
				"message": fmt.Sprintf("Resource '%s' does not exist", id),
			},
		})
	}

	return httpmock.NewJsonResponse(200, config)
}

func (m *GroupPolicyConfigurationMock) getGroupPolicyConfigurationAssignments(req *http.Request) (*http.Response, error) {
	pathParts := strings.Split(req.URL.Path, "/")
	id := pathParts[len(pathParts)-2]

	m.stateMu.RLock()
	assignments, exists := m.state[id+"_assignments"]
	m.stateMu.RUnlock()

	if !exists {
		assignments = []any{}
	}

	response := map[string]any{
		"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyConfigurations('" + id + "')/assignments",
		"value":          assignments,
	}

	return httpmock.NewJsonResponse(200, response)
}

func (m *GroupPolicyConfigurationMock) updateGroupPolicyConfiguration(req *http.Request) (*http.Response, error) {
	pathParts := strings.Split(req.URL.Path, "/")
	id := pathParts[len(pathParts)-1]

	m.stateMu.RLock()
	config, exists := m.state[id]
	m.stateMu.RUnlock()

	if !exists {
		return httpmock.NewJsonResponse(404, map[string]any{
			"error": map[string]any{
				"code":    "NotFound",
				"message": fmt.Sprintf("Resource '%s' does not exist", id),
			},
		})
	}

	var updateBody map[string]any
	if err := json.NewDecoder(req.Body).Decode(&updateBody); err != nil {
		return httpmock.NewJsonResponse(400, map[string]any{
			"error": map[string]any{
				"code":    "BadRequest",
				"message": "Invalid request body",
			},
		})
	}

	// Determine which scenario file to load based on displayName
	scenarioFile := determineUpdateScenario(updateBody, config)
	if scenarioFile != "" {
		jsonContent, err := helpers.ParseJSONFile(filepath.Join("..", "tests", "responses", "validate_update", scenarioFile))
		if err == nil {
			var response map[string]any
			if err := json.Unmarshal([]byte(jsonContent), &response); err == nil {
				response["id"] = id
				m.stateMu.Lock()
				m.state[id] = response
				m.stateMu.Unlock()
				return httpmock.NewJsonResponse(200, response)
			}
		}
	}

	// Fallback to dynamic update
	m.stateMu.Lock()
	if configMap, ok := config.(map[string]any); ok {
		for k, v := range updateBody {
			configMap[k] = v
		}
		configMap["lastModifiedDateTime"] = "2024-01-02T00:00:00Z"
		m.state[id] = configMap
	}
	m.stateMu.Unlock()

	return httpmock.NewJsonResponse(200, config)
}

func determineUpdateScenario(body map[string]any, existing any) string {
	displayName, _ := body["displayName"].(string)

	switch {
	case strings.Contains(displayName, "unit-test-005-lifecycle-maximal"):
		return "patch_test_005_lifecycle_maximal.json"
	case strings.Contains(displayName, "unit-test-006-lifecycle-minimal"):
		return "patch_test_006_lifecycle_minimal.json"
	default:
		return ""
	}
}

func (m *GroupPolicyConfigurationMock) deleteGroupPolicyConfiguration(req *http.Request) (*http.Response, error) {
	pathParts := strings.Split(req.URL.Path, "/")
	id := pathParts[len(pathParts)-1]

	m.stateMu.Lock()
	delete(m.state, id)
	delete(m.state, id+"_assignments")
	m.stateMu.Unlock()

	return httpmock.NewJsonResponse(204, nil)
}
