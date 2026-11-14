package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const (
	workflowsURL           = "https://graph.microsoft.com/beta/identityGovernance/lifecycleWorkflows/workflows"
	guidPattern            = `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`
	workflowItemPattern    = `=~^https://graph\.microsoft\.com/beta/identityGovernance/lifecycleWorkflows/workflows/` + guidPattern + `$`
	createNewVersionURL    = `=~^https://graph\.microsoft\.com/beta/identityGovernance/lifecycleWorkflows/workflows/` + guidPattern + `/microsoft\.graph\.identityGovernance\.createNewVersion$`
	deletedWorkflowPattern = `=~^https://graph\.microsoft\.com/beta/identityGovernance/lifecycleWorkflows/deletedItems/workflows/` + guidPattern + `$`
	notFoundBody           = `{"error":{"code":"Resource not found","message":""}}`
	forbiddenBody          = `{"error":{"code":"Forbidden","message":"Insufficient privileges to complete the operation."}}`
	patchRejectedMessage   = "Only displayname description isenabled and isschedulingenabled can be patched"
)

var mockState struct {
	sync.Mutex
	workflows        map[string]map[string]any
	deletedWorkflows map[string]map[string]any
}

// basicProperties are the only properties Graph accepts on PATCH.
var basicProperties = map[string]bool{"displayName": true, "description": true, "isEnabled": true, "isSchedulingEnabled": true}

func init() {
	mockState.workflows = make(map[string]map[string]any)
	mockState.deletedWorkflows = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, notFoundBody))
	mocks.GlobalRegistry.Register("identity_and_access_lifecycle_workflow", &LifecycleWorkflowMock{})
}

// LifecycleWorkflowMock provides stateful httpmock responders for lifecycle workflows
type LifecycleWorkflowMock struct{}

var _ mocks.MockRegistrar = (*LifecycleWorkflowMock)(nil)

// getJSONFileForDisplayName picks the create response fixture from the scenario prefix of the
// display name, for example "unit-test-lw001: ..." loads post_lifecycle_workflow_lw001_success.json.
func getJSONFileForDisplayName(displayName string) string {
	match := regexp.MustCompile(`(?i)\blw\d{3}\b`).FindString(displayName)
	if match == "" {
		return ""
	}
	return fmt.Sprintf("post_lifecycle_workflow_%s_success.json", strings.ToLower(match))
}

func errorResponse(status int, code, message string) *http.Response {
	resp, _ := httpmock.NewJsonResponse(status, map[string]any{"error": map[string]any{"code": code, "message": message}})
	return resp
}

// invalidCollectionResponse mirrors Graph, which rejects null collections in a workflow body and a
// task whose arguments property is missing or null ("Task arguments missing or invalid").
func invalidCollectionResponse(workflow map[string]any) *http.Response {
	if value, ok := workflow["tasks"]; ok && value == nil {
		return errorResponse(400, "BAD_REQUEST", "tasks must not be null")
	}
	tasks, _ := workflow["tasks"].([]any)
	for _, item := range tasks {
		task, _ := item.(map[string]any)
		if value, ok := task["arguments"]; !ok || value == nil {
			return errorResponse(400, "BAD_REQUEST", "Task arguments missing or invalid")
		}
	}
	return nil
}

// modifiedTimestamp returns a distinct lastModifiedDateTime for each write to a workflow.
func modifiedTimestamp(revision int) string {
	return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(revision) * time.Minute).Format("2006-01-02T15:04:05.0000000Z")
}

// applyWorkflowBody copies the writable properties of a request onto a stored workflow and
// fills in the defaults Graph adds: empty description, numbered tasks with new IDs, task
// category from the workflow, isEnabled/continueOnError false and empty arguments.
func applyWorkflowBody(stored map[string]any, body map[string]any) {
	for _, key := range []string{"category", "displayName", "description", "isEnabled", "isSchedulingEnabled", "executionConditions"} {
		if value, ok := body[key]; ok {
			stored[key] = value
		}
	}
	if _, ok := stored["description"]; !ok || stored["description"] == nil {
		stored["description"] = ""
	}
	for _, key := range []string{"isEnabled", "isSchedulingEnabled"} {
		if _, ok := stored[key]; !ok {
			stored[key] = false
		}
	}

	requested, _ := body["tasks"].([]any)
	tasks := make([]any, 0, len(requested))
	for i, item := range requested {
		source, _ := item.(map[string]any)
		task := map[string]any{
			"id":                uuid.New().String(),
			"category":          stored["category"],
			"description":       "",
			"isEnabled":         false,
			"continueOnError":   false,
			"executionSequence": i + 1,
			"arguments":         []any{},
		}
		for _, key := range []string{"taskDefinitionId", "displayName", "description", "isEnabled", "continueOnError", "arguments"} {
			if value, ok := source[key]; ok && value != nil {
				task[key] = value
			}
		}
		tasks = append(tasks, task)
	}
	stored["tasks"] = tasks
}

func workflowIdFromPath(path string, offsetFromEnd int) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1-offsetFromEnd]
}

func (m *LifecycleWorkflowMock) RegisterMocks() {
	mockState.Lock()
	mockState.workflows = make(map[string]map[string]any)
	mockState.deletedWorkflows = make(map[string]map[string]any)
	mockState.Unlock()

	// Create workflow - POST /identityGovernance/lifecycleWorkflows/workflows
	httpmock.RegisterResponder("POST", workflowsURL, func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return errorResponse(400, "BAD_REQUEST", "Invalid request body"), nil
		}
		if resp := invalidCollectionResponse(requestBody); resp != nil {
			return resp, nil
		}
		if _, ok := requestBody["category"].(string); !ok {
			return errorResponse(400, "BAD_REQUEST", "Missing or invalid parameters (Parameter 'Category')"), nil
		}

		displayName, ok := requestBody["displayName"].(string)
		if !ok {
			return errorResponse(400, "BAD_REQUEST", "displayName is required"), nil
		}

		jsonFileName := getJSONFileForDisplayName(displayName)
		if jsonFileName == "" {
			return errorResponse(400, "BAD_REQUEST", fmt.Sprintf("Unable to determine JSON file for displayName: %s", displayName)), nil
		}

		jsonData, err := os.ReadFile(filepath.Join("tests", "responses", "validate_create", jsonFileName))
		if err != nil {
			return errorResponse(500, "InternalServerError", fmt.Sprintf("Failed to load JSON response file: %s", err.Error())), nil
		}

		var responseObj map[string]any
		if err := json.Unmarshal(jsonData, &responseObj); err != nil {
			return errorResponse(500, "InternalServerError", fmt.Sprintf("Failed to parse JSON response: %s", err.Error())), nil
		}

		newId := uuid.New().String()
		responseObj["id"] = newId
		responseObj["version"] = 1
		responseObj["revision"] = 0
		responseObj["lastModifiedDateTime"] = modifiedTimestamp(0)
		applyWorkflowBody(responseObj, requestBody)

		mockState.Lock()
		mockState.workflows[newId] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// Get workflow - GET /identityGovernance/lifecycleWorkflows/workflows/{id}
	httpmock.RegisterResponder("GET", workflowItemPattern, func(req *http.Request) (*http.Response, error) {
		workflowId := workflowIdFromPath(req.URL.Path, 0)

		mockState.Lock()
		workflow, exists := mockState.workflows[workflowId]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}

		return httpmock.NewJsonResponse(200, workflow)
	})

	// Update basic properties - PATCH /identityGovernance/lifecycleWorkflows/workflows/{id}
	httpmock.RegisterResponder("PATCH", workflowItemPattern, func(req *http.Request) (*http.Response, error) {
		workflowId := workflowIdFromPath(req.URL.Path, 0)

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return errorResponse(400, "BAD_REQUEST", "Invalid request body"), nil
		}
		for key := range requestBody {
			if !basicProperties[key] && !strings.HasPrefix(key, "@odata") {
				return errorResponse(400, "BAD_REQUEST", patchRejectedMessage), nil
			}
		}

		mockState.Lock()
		defer mockState.Unlock()

		workflow, exists := mockState.workflows[workflowId]
		if !exists {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}

		for key, value := range requestBody {
			if basicProperties[key] {
				workflow[key] = value
			}
		}
		revision, _ := workflow["revision"].(int)
		workflow["revision"] = revision + 1
		workflow["lastModifiedDateTime"] = modifiedTimestamp(revision + 1)

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Create new version - POST /identityGovernance/lifecycleWorkflows/workflows/{id}/microsoft.graph.identityGovernance.createNewVersion
	httpmock.RegisterResponder("POST", createNewVersionURL, func(req *http.Request) (*http.Response, error) {
		workflowId := workflowIdFromPath(req.URL.Path, 1)

		var requestBody struct {
			Workflow map[string]any `json:"workflow"`
		}
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil || requestBody.Workflow == nil {
			return errorResponse(400, "BAD_REQUEST", "Missing or invalid parameters (Parameter 'workflow')"), nil
		}
		if _, ok := requestBody.Workflow["category"].(string); !ok {
			return errorResponse(400, "BAD_REQUEST", "Missing or invalid parameters (Parameter 'Category')"), nil
		}
		if resp := invalidCollectionResponse(requestBody.Workflow); resp != nil {
			return resp, nil
		}

		mockState.Lock()
		defer mockState.Unlock()

		workflow, exists := mockState.workflows[workflowId]
		if !exists {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}

		applyWorkflowBody(workflow, requestBody.Workflow)
		version, _ := workflow["version"].(int)
		workflow["version"] = version + 1
		revision, _ := workflow["revision"].(int)
		workflow["revision"] = revision + 1
		workflow["lastModifiedDateTime"] = modifiedTimestamp(revision + 1)

		return httpmock.NewJsonResponse(200, workflow)
	})

	// Delete workflow - DELETE /identityGovernance/lifecycleWorkflows/workflows/{id}
	// Moves the workflow to deletedItems, as Graph does.
	httpmock.RegisterResponder("DELETE", workflowItemPattern, func(req *http.Request) (*http.Response, error) {
		workflowId := workflowIdFromPath(req.URL.Path, 0)

		mockState.Lock()
		defer mockState.Unlock()

		if workflow, exists := mockState.workflows[workflowId]; exists {
			workflow["deletedDateTime"] = modifiedTimestamp(1000)
			mockState.deletedWorkflows[workflowId] = workflow
			delete(mockState.workflows, workflowId)
			return httpmock.NewStringResponse(204, ""), nil
		}

		return httpmock.NewStringResponse(404, notFoundBody), nil
	})

	// Get deleted workflow - GET /identityGovernance/lifecycleWorkflows/deletedItems/workflows/{id}
	httpmock.RegisterResponder("GET", deletedWorkflowPattern, func(req *http.Request) (*http.Response, error) {
		workflowId := workflowIdFromPath(req.URL.Path, 0)

		mockState.Lock()
		workflow, exists := mockState.deletedWorkflows[workflowId]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}
		return httpmock.NewJsonResponse(200, workflow)
	})

	// Permanently delete workflow - DELETE /identityGovernance/lifecycleWorkflows/deletedItems/workflows/{id}
	httpmock.RegisterResponder("DELETE", deletedWorkflowPattern, func(req *http.Request) (*http.Response, error) {
		workflowId := workflowIdFromPath(req.URL.Path, 0)

		mockState.Lock()
		defer mockState.Unlock()

		if _, exists := mockState.deletedWorkflows[workflowId]; exists {
			delete(mockState.deletedWorkflows, workflowId)
			return httpmock.NewStringResponse(204, ""), nil
		}

		return httpmock.NewStringResponse(404, notFoundBody), nil
	})
}

func (m *LifecycleWorkflowMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", workflowsURL, httpmock.NewStringResponder(403, forbiddenBody))
	httpmock.RegisterResponder("GET", workflowItemPattern, httpmock.NewStringResponder(404, notFoundBody))
	httpmock.RegisterResponder("PATCH", workflowItemPattern, httpmock.NewStringResponder(400, `{"error":{"code":"BAD_REQUEST","message":"Invalid request"}}`))
	httpmock.RegisterResponder("POST", createNewVersionURL, httpmock.NewStringResponder(400, `{"error":{"code":"BAD_REQUEST","message":"Invalid request"}}`))
	httpmock.RegisterResponder("DELETE", workflowItemPattern, httpmock.NewStringResponder(404, notFoundBody))
}

func (m *LifecycleWorkflowMock) CleanupMockState() {
	mockState.Lock()
	mockState.workflows = make(map[string]map[string]any)
	mockState.deletedWorkflows = make(map[string]map[string]any)
	mockState.Unlock()
}
