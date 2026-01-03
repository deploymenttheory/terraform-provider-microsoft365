package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

type GroupPolicyConfigurationMock struct {
	state   map[string]interface{}
	stateMu sync.RWMutex
}

func init() {
	mocks.GlobalRegistry.Register("group_policy_configuration", &GroupPolicyConfigurationMock{})
}

func (m *GroupPolicyConfigurationMock) RegisterMocks() {
	m.state = make(map[string]interface{})

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
			return httpmock.NewJsonResponse(400, map[string]interface{}{
				"error": map[string]interface{}{
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
	m.state = make(map[string]interface{})
}

func (m *GroupPolicyConfigurationMock) createGroupPolicyConfiguration(req *http.Request) (*http.Response, error) {
	var requestBody map[string]interface{}
	if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
		return httpmock.NewJsonResponse(400, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "BadRequest",
				"message": "Invalid request body",
			},
		})
	}

	id := uuid.New().String()
	response, err := mocks.LoadJSONResponse("tests/responses/validate_create/post_group_policy_configuration_success.json")
	if err != nil {
		return httpmock.NewJsonResponse(500, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "InternalServerError",
				"message": fmt.Sprintf("Failed to load response fixture: %v", err),
			},
		})
	}

	response["id"] = id
	if displayName, ok := requestBody["displayName"].(string); ok {
		response["displayName"] = displayName
	}
	if description, ok := requestBody["description"].(string); ok {
		response["description"] = description
	}
	if roleScopeTagIds, ok := requestBody["roleScopeTagIds"].([]interface{}); ok {
		response["roleScopeTagIds"] = roleScopeTagIds
	}

	m.stateMu.Lock()
	m.state[id] = response
	if assignments, ok := requestBody["assignments"]; ok {
		m.state[id+"_assignments"] = assignments
	}
	m.stateMu.Unlock()

	return httpmock.NewJsonResponse(201, response)
}

func (m *GroupPolicyConfigurationMock) assignGroupPolicyConfiguration(req *http.Request) (*http.Response, error) {
	pathParts := strings.Split(req.URL.Path, "/")
	id := pathParts[len(pathParts)-2]

	var requestBody map[string]interface{}
	if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
		return httpmock.NewJsonResponse(400, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "BadRequest",
				"message": "Invalid request body",
			},
		})
	}

	m.stateMu.Lock()
	if assignments, ok := requestBody["assignments"].([]interface{}); ok {
		m.state[id+"_assignments"] = assignments
	}
	m.stateMu.Unlock()

	response := map[string]interface{}{
		"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyConfigurations('" + id + "')/Microsoft.Graph.assign",
		"value":          []interface{}{},
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
		return httpmock.NewJsonResponse(404, map[string]interface{}{
			"error": map[string]interface{}{
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
		assignments = []interface{}{}
	}

	response := map[string]interface{}{
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
		return httpmock.NewJsonResponse(404, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "NotFound",
				"message": fmt.Sprintf("Resource '%s' does not exist", id),
			},
		})
	}

	var updateBody map[string]interface{}
	if err := json.NewDecoder(req.Body).Decode(&updateBody); err != nil {
		return httpmock.NewJsonResponse(400, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "BadRequest",
				"message": "Invalid request body",
			},
		})
	}

	m.stateMu.Lock()
	if configMap, ok := config.(map[string]interface{}); ok {
		for k, v := range updateBody {
			configMap[k] = v
		}
		m.state[id] = configMap
	}
	m.stateMu.Unlock()

	return httpmock.NewJsonResponse(200, config)
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
