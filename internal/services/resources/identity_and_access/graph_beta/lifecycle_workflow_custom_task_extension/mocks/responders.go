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
	customTaskExtensionsURL        = "https://graph.microsoft.com/beta/identityGovernance/lifecycleWorkflows/customTaskExtensions"
	guidPattern                    = `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`
	customTaskExtensionItemPattern = `=~^https://graph\.microsoft\.com/beta/identityGovernance/lifecycleWorkflows/customTaskExtensions/` + guidPattern + `$`
	notFoundBody                   = `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`
	forbiddenBody                  = `{"error":{"code":"Forbidden","message":"Insufficient privileges to complete the operation."}}`
	badRequestBody                 = `{"error":{"code":"BAD_REQUEST","message":"Invalid model state"}}`
	logicAppTriggerODataType       = "#microsoft.graph.logicAppTriggerEndpointConfiguration"
)

var mockState struct {
	sync.Mutex
	extensions map[string]map[string]any
}

func init() {
	mockState.extensions = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, notFoundBody))
	mocks.GlobalRegistry.Register("identity_and_access_lifecycle_workflow_custom_task_extension", &LifecycleWorkflowCustomTaskExtensionMock{})
}

// LifecycleWorkflowCustomTaskExtensionMock provides stateful httpmock responders for custom task extensions
type LifecycleWorkflowCustomTaskExtensionMock struct{}

var _ mocks.MockRegistrar = (*LifecycleWorkflowCustomTaskExtensionMock)(nil)

// getJSONFileForDisplayName picks the create response fixture from the scenario prefix of the
// display name, for example "unit-test-ctx001-minimal" loads post_custom_task_extension_ctx001_success.json.
func getJSONFileForDisplayName(displayName string) string {
	match := regexp.MustCompile(`(?i)\bctx\d{3}\b`).FindString(displayName)
	if match == "" {
		return ""
	}
	return fmt.Sprintf("post_custom_task_extension_%s_success.json", strings.ToLower(match))
}

func errorResponse(status int, code, message string) *http.Response {
	resp, _ := httpmock.NewJsonResponse(status, map[string]any{"error": map[string]any{"code": code, "message": message}})
	return resp
}

// invalidBodyResponse mirrors the validation Graph applies to an app-only create: displayName,
// endpointConfiguration with url, authenticationConfiguration and clientConfiguration are required.
func invalidBodyResponse(body map[string]any) *http.Response {
	endpoint, _ := body["endpointConfiguration"].(map[string]any)
	if endpoint == nil || endpoint["@odata.type"] != logicAppTriggerODataType {
		return errorResponse(400, "BAD_REQUEST", "Invalid model state")
	}
	if url, _ := endpoint["url"].(string); url == "" {
		return errorResponse(400, "BAD_REQUEST", "Custom task extensions missing url")
	}
	if body["authenticationConfiguration"] == nil || body["clientConfiguration"] == nil {
		return errorResponse(400, "BAD_REQUEST", "Invalid model state")
	}
	return nil
}

// modifiedTimestamp returns a distinct lastModifiedDateTime for each write to an extension.
func modifiedTimestamp(revision int) string {
	return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(revision) * time.Minute).Format("2006-01-02T15:04:05.0000000Z")
}

// applyBody copies the writable properties of a request onto a stored extension and fills in what
// Graph derives: an empty description, replyMode from callbackConfiguration, authorizedApps as []
// and @odata.type dropped from clientConfiguration and authorized apps.
func applyBody(stored map[string]any, body map[string]any) {
	for _, key := range []string{"displayName", "description", "endpointConfiguration", "authenticationConfiguration", "clientConfiguration", "callbackConfiguration"} {
		if value, ok := body[key]; ok {
			stored[key] = value
		}
	}
	if _, ok := stored["description"]; !ok || stored["description"] == nil {
		stored["description"] = ""
	}
	if client, ok := stored["clientConfiguration"].(map[string]any); ok {
		delete(client, "@odata.type")
	}

	callback, _ := stored["callbackConfiguration"].(map[string]any)
	if callback == nil {
		stored["callbackConfiguration"] = nil
		stored["replyMode"] = "none"
		return
	}

	stored["replyMode"] = "callback"
	apps, _ := callback["authorizedApps"].([]any)
	authorizedApps := make([]any, 0, len(apps))
	for _, item := range apps {
		app, _ := item.(map[string]any)
		authorizedApps = append(authorizedApps, map[string]any{"id": app["id"]})
	}
	callback["authorizedApps"] = authorizedApps
}

func idFromPath(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

func (m *LifecycleWorkflowCustomTaskExtensionMock) RegisterMocks() {
	mockState.Lock()
	mockState.extensions = make(map[string]map[string]any)
	mockState.Unlock()

	// Create custom task extension - POST /identityGovernance/lifecycleWorkflows/customTaskExtensions
	httpmock.RegisterResponder("POST", customTaskExtensionsURL, func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return errorResponse(400, "BAD_REQUEST", "Invalid request body"), nil
		}
		if resp := invalidBodyResponse(requestBody); resp != nil {
			return resp, nil
		}

		displayName, ok := requestBody["displayName"].(string)
		if !ok {
			return errorResponse(400, "BAD_REQUEST", "Invalid model state"), nil
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
		responseObj["revision"] = 0
		responseObj["lastModifiedDateTime"] = modifiedTimestamp(0)
		applyBody(responseObj, requestBody)

		mockState.Lock()
		mockState.extensions[newId] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// Get custom task extension - GET /identityGovernance/lifecycleWorkflows/customTaskExtensions/{id}
	httpmock.RegisterResponder("GET", customTaskExtensionItemPattern, func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		extension, exists := mockState.extensions[idFromPath(req.URL.Path)]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}

		return httpmock.NewJsonResponse(200, extension)
	})

	// Update custom task extension - PATCH /identityGovernance/lifecycleWorkflows/customTaskExtensions/{id}
	httpmock.RegisterResponder("PATCH", customTaskExtensionItemPattern, func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return errorResponse(400, "BAD_REQUEST", "Invalid request body"), nil
		}

		mockState.Lock()
		defer mockState.Unlock()

		extension, exists := mockState.extensions[idFromPath(req.URL.Path)]
		if !exists {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}

		applyBody(extension, requestBody)
		revision, _ := extension["revision"].(int)
		extension["revision"] = revision + 1
		extension["lastModifiedDateTime"] = modifiedTimestamp(revision + 1)

		return httpmock.NewStringResponse(204, ""), nil
	})

	// Delete custom task extension - DELETE /identityGovernance/lifecycleWorkflows/customTaskExtensions/{id}
	httpmock.RegisterResponder("DELETE", customTaskExtensionItemPattern, func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		id := idFromPath(req.URL.Path)
		if _, exists := mockState.extensions[id]; exists {
			delete(mockState.extensions, id)
			return httpmock.NewStringResponse(204, ""), nil
		}

		return httpmock.NewStringResponse(404, notFoundBody), nil
	})
}

func (m *LifecycleWorkflowCustomTaskExtensionMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", customTaskExtensionsURL, httpmock.NewStringResponder(403, forbiddenBody))
	httpmock.RegisterResponder("GET", customTaskExtensionItemPattern, httpmock.NewStringResponder(404, notFoundBody))
	httpmock.RegisterResponder("PATCH", customTaskExtensionItemPattern, httpmock.NewStringResponder(400, badRequestBody))
	httpmock.RegisterResponder("DELETE", customTaskExtensionItemPattern, httpmock.NewStringResponder(404, notFoundBody))
}

// RegisterBadRequestMocks makes create fail with 400, as Graph does for an invalid body.
func (m *LifecycleWorkflowCustomTaskExtensionMock) RegisterBadRequestMocks() {
	httpmock.RegisterResponder("POST", customTaskExtensionsURL, httpmock.NewStringResponder(400, badRequestBody))
}

func (m *LifecycleWorkflowCustomTaskExtensionMock) CleanupMockState() {
	mockState.Lock()
	mockState.extensions = make(map[string]map[string]any)
	mockState.Unlock()
}
