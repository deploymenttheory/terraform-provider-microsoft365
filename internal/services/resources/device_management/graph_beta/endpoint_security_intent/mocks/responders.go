package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const (
	// ForbiddenTemplateId is a template ID the mock answers with a permission error.
	ForbiddenTemplateId = "00000000-0000-0000-0000-000000000403"
	// UnknownTemplateId is a template ID the mock rejects the way Graph rejects templates that do not exist.
	UnknownTemplateId = "00000000-0000-0000-0000-000000000400"

	guid                 = `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`
	createInstanceUrl    = `=~^https://graph\.microsoft\.com/beta/deviceManagement/templates/(` + guid + `)/createInstance$`
	intentByIdUrl        = `=~^https://graph\.microsoft\.com/beta/deviceManagement/intents/(` + guid + `)$`
	intentSettingsUrl    = `=~^https://graph\.microsoft\.com/beta/deviceManagement/intents/(` + guid + `)/settings$`
	intentAssignmentsUrl = `=~^https://graph\.microsoft\.com/beta/deviceManagement/intents/(` + guid + `)/assignments$`
	updateSettingsUrl    = `=~^https://graph\.microsoft\.com/beta/deviceManagement/intents/(` + guid + `)/updateSettings$`
	assignUrl            = `=~^https://graph\.microsoft\.com/beta/deviceManagement/intents/(` + guid + `)/assign$`
	notFoundBody         = `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`
	abstractSettingType  = "#microsoft.graph.deviceManagementSettingInstance"
)

var intentIdPattern = regexp.MustCompile(guid)

type intentState struct {
	intent      map[string]any
	settings    []map[string]any
	assignments []map[string]any
}

var mockState struct {
	sync.Mutex
	intents map[string]*intentState
}

func init() {
	mockState.intents = make(map[string]*intentState)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, notFoundBody))
	mocks.GlobalRegistry.Register("endpoint_security_intent", &EndpointSecurityIntentMock{})
}

type EndpointSecurityIntentMock struct{}

var _ mocks.MockRegistrar = (*EndpointSecurityIntentMock)(nil)

func loadJson(filename string) (map[string]any, error) {
	jsonData, err := os.ReadFile(filepath.Join("tests", "responses", "validate_create", filename))
	if err != nil {
		return nil, fmt.Errorf("failed to load JSON response file: %w", err)
	}

	var responseObj map[string]any
	if err := json.Unmarshal(jsonData, &responseObj); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}
	return responseObj, nil
}

// loadTemplateSettings loads the default settings Graph gives every intent created from the BitLocker template.
func loadTemplateSettings() ([]map[string]any, error) {
	response, err := loadJson("get_endpoint_security_intent_settings_success.json")
	if err != nil {
		return nil, err
	}

	values, _ := response["value"].([]any)
	settings := make([]map[string]any, 0, len(values))
	for _, value := range values {
		if setting, ok := value.(map[string]any); ok {
			settings = append(settings, setting)
		}
	}
	return settings, nil
}

func graphError(status int, code, message string) (*http.Response, error) {
	resp, err := httpmock.NewJsonResponse(status, map[string]any{"error": map[string]any{"code": code, "message": message}})
	if err != nil {
		return nil, fmt.Errorf("failed to build mock error response: %w", err)
	}
	return resp, nil
}

func idFromRequest(req *http.Request) string {
	return intentIdPattern.FindAllString(req.URL.Path, -1)[0]
}

// normaliseValueJson re-encodes a value the way Graph does, so unit tests see formatting differences between config and Graph.
func normaliseValueJson(valueJson string) string {
	var value any
	if err := json.Unmarshal([]byte(valueJson), &value); err != nil {
		return valueJson
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return valueJson
	}
	return string(encoded)
}

// applySettings writes setting instances onto the intent. Like Graph, it fails with a 500 when the abstract base type is sent
// and with a 400 when a definition is not part of the template.
func applySettings(state *intentState, settings []any) (int, string) {
	for _, raw := range settings {
		setting, _ := raw.(map[string]any)
		odataType, _ := setting["@odata.type"].(string)
		if odataType == "" || odataType == abstractSettingType {
			return 500, "An internal server error has occurred"
		}
		definitionId, _ := setting["definitionId"].(string)
		valueJson, _ := setting["valueJson"].(string)

		found := false
		for _, existing := range state.settings {
			if existing["definitionId"] == definitionId {
				existing["@odata.type"] = odataType
				existing["valueJson"] = normaliseValueJson(valueJson)
				found = true
				break
			}
		}
		if !found {
			return 400, fmt.Sprintf("Setting definition %s is not part of the template", definitionId)
		}
	}
	return 0, ""
}

func (m *EndpointSecurityIntentMock) RegisterMocks() {
	mockState.Lock()
	mockState.intents = make(map[string]*intentState)
	mockState.Unlock()

	// Create intent - POST /deviceManagement/templates/{templateId}/createInstance
	httpmock.RegisterResponder("POST", createInstanceUrl, func(req *http.Request) (*http.Response, error) {
		templateId := idFromRequest(req)
		switch templateId {
		case ForbiddenTemplateId:
			return graphError(403, "Forbidden", "Application is not authorized to perform this operation.")
		case UnknownTemplateId:
			return graphError(400, "BadRequest", "Template not found.")
		}

		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return graphError(400, "BadRequest", "Invalid request body")
		}
		if value, ok := requestBody["settingsDelta"]; ok && value == nil {
			return graphError(400, "BadRequest", "settingsDelta cannot be null")
		}

		intent, err := loadJson("post_endpoint_security_intent_success.json")
		if err != nil {
			return graphError(500, "InternalServerError", err.Error())
		}
		settings, err := loadTemplateSettings()
		if err != nil {
			return graphError(500, "InternalServerError", err.Error())
		}

		id := uuid.New().String()
		intent["id"] = id
		intent["templateId"] = templateId
		for _, key := range []string{"displayName", "description", "roleScopeTagIds"} {
			if value, ok := requestBody[key]; ok {
				intent[key] = value
			}
		}

		state := &intentState{intent: intent, settings: settings, assignments: []map[string]any{}}
		if delta, ok := requestBody["settingsDelta"].([]any); ok {
			if status, message := applySettings(state, delta); status != 0 {
				return graphError(status, "BadRequest", message)
			}
		}

		mockState.Lock()
		mockState.intents[id] = state
		mockState.Unlock()

		return httpmock.NewJsonResponse(200, intent)
	})

	// Get intent - GET /deviceManagement/intents/{id}
	httpmock.RegisterResponder("GET", intentByIdUrl, func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		state, exists := mockState.intents[idFromRequest(req)]
		if !exists {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}
		state.intent["isAssigned"] = len(state.assignments) > 0
		return httpmock.NewJsonResponse(200, state.intent)
	})

	// List settings - GET /deviceManagement/intents/{id}/settings
	httpmock.RegisterResponder("GET", intentSettingsUrl, func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		state, exists := mockState.intents[idFromRequest(req)]
		if !exists {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}
		return httpmock.NewJsonResponse(200, map[string]any{"value": state.settings})
	})

	// List assignments - GET /deviceManagement/intents/{id}/assignments
	httpmock.RegisterResponder("GET", intentAssignmentsUrl, func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		state, exists := mockState.intents[idFromRequest(req)]
		if !exists {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}
		return httpmock.NewJsonResponse(200, map[string]any{"value": state.assignments})
	})

	// Update intent - PATCH /deviceManagement/intents/{id}
	httpmock.RegisterResponder("PATCH", intentByIdUrl, func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return graphError(400, "BadRequest", "Invalid request body")
		}

		mockState.Lock()
		defer mockState.Unlock()

		state, exists := mockState.intents[idFromRequest(req)]
		if !exists {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}
		for _, key := range []string{"displayName", "description", "roleScopeTagIds"} {
			if value, ok := requestBody[key]; ok {
				state.intent[key] = value
			}
		}
		return httpmock.NewStringResponse(204, ""), nil
	})

	// Update settings - POST /deviceManagement/intents/{id}/updateSettings
	httpmock.RegisterResponder("POST", updateSettingsUrl, func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return graphError(400, "BadRequest", "Invalid request body")
		}

		mockState.Lock()
		defer mockState.Unlock()

		state, exists := mockState.intents[idFromRequest(req)]
		if !exists {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}
		settings, _ := requestBody["settings"].([]any)
		if status, message := applySettings(state, settings); status != 0 {
			return graphError(status, "BadRequest", message)
		}
		return httpmock.NewStringResponse(204, ""), nil
	})

	// Replace assignments - POST /deviceManagement/intents/{id}/assign
	httpmock.RegisterResponder("POST", assignUrl, func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			return graphError(400, "BadRequest", "Invalid request body")
		}

		mockState.Lock()
		defer mockState.Unlock()

		state, exists := mockState.intents[idFromRequest(req)]
		if !exists {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}
		assignments := make([]map[string]any, 0)
		values, _ := requestBody["assignments"].([]any)
		for _, value := range values {
			assignment, _ := value.(map[string]any)
			target, _ := assignment["target"].(map[string]any)
			if _, ok := target["deviceAndAppManagementAssignmentFilterType"]; !ok {
				target["deviceAndAppManagementAssignmentFilterType"] = "none"
				target["deviceAndAppManagementAssignmentFilterId"] = nil
			}
			assignments = append(assignments, map[string]any{"id": uuid.New().String(), "target": target})
		}
		state.assignments = assignments
		return httpmock.NewStringResponse(204, ""), nil
	})

	// Delete intent - DELETE /deviceManagement/intents/{id}
	httpmock.RegisterResponder("DELETE", intentByIdUrl, func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		id := idFromRequest(req)
		if _, exists := mockState.intents[id]; !exists {
			return httpmock.NewStringResponse(404, notFoundBody), nil
		}
		delete(mockState.intents, id)
		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *EndpointSecurityIntentMock) RegisterErrorMocks() {
	// Error scenarios for testing
	httpmock.RegisterResponder("POST", createInstanceUrl, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("GET", intentByIdUrl, httpmock.NewStringResponder(404, notFoundBody))
	httpmock.RegisterResponder("PATCH", intentByIdUrl, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
	httpmock.RegisterResponder("DELETE", intentByIdUrl, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

// RemoveIntent deletes an intent from the mock state to simulate an out-of-band deletion.
func (m *EndpointSecurityIntentMock) RemoveIntent(id string) {
	mockState.Lock()
	delete(mockState.intents, id)
	mockState.Unlock()
}

// SettingValue returns the value JSON the mock holds for a setting, so tests can confirm updateSettings was called.
func (m *EndpointSecurityIntentMock) SettingValue(id, definitionId string) string {
	mockState.Lock()
	defer mockState.Unlock()

	state, exists := mockState.intents[id]
	if !exists {
		return ""
	}
	for _, setting := range state.settings {
		if setting["definitionId"] == definitionId {
			value, _ := setting["valueJson"].(string)
			return value
		}
	}
	return ""
}

func (m *EndpointSecurityIntentMock) CleanupMockState() {
	mockState.Lock()
	mockState.intents = make(map[string]*intentState)
	mockState.Unlock()
}
