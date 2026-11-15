package mocks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const settingsURL = "https://graph.microsoft.com/beta/identityGovernance/lifecycleWorkflows/settings"

// verifiedSenderDomains mirrors the domains Graph accepts as emailSettings.senderDomain in the mock tenant.
// Graph compares the value case-sensitively against the tenant domain list.
var verifiedSenderDomains = map[string]bool{
	"microsoft.com": true,
	"contoso.com":   true,
}

var mockState struct {
	sync.Mutex
	settings map[string]any
	patches  []map[string]any
	seedErr  error
}

func init() {
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("identity_and_access_lifecycle_workflow_settings", &LifecycleWorkflowSettingsMock{})
}

// LifecycleWorkflowSettingsMock provides stateful httpmock responders for the singleton lifecycle workflow settings.
type LifecycleWorkflowSettingsMock struct{}

var _ mocks.MockRegistrar = (*LifecycleWorkflowSettingsMock)(nil)

// loadSeedSettings loads the tenant settings that GET returns before any PATCH.
func loadSeedSettings() (map[string]any, error) {
	responsesPath := filepath.Join("tests", "responses", "validate_get", "get_lifecycle_workflow_settings_success.json")
	jsonData, err := os.ReadFile(responsesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load JSON response file: %w", err)
	}

	var settings map[string]any
	if err := json.Unmarshal(jsonData, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}
	return settings, nil
}

func errorResponse(status int, code, message string) (*http.Response, error) {
	resp, err := httpmock.NewJsonResponse(status, map[string]any{"error": map[string]any{"code": code, "message": message}})
	if err != nil {
		return nil, fmt.Errorf("build error response: %w", err)
	}
	return resp, nil
}

// graphError is a Graph 400 error code and message returned by validatePatch.
type graphError struct {
	code    string
	message string
}

// validatePatch mirrors the validation Graph applies to a PATCH of the settings. It returns nil when the body is valid.
func validatePatch(body map[string]any) *graphError {
	if value, ok := body["workflowScheduleIntervalInHours"]; ok {
		interval, isNumber := value.(float64)
		if !isNumber || interval < 1 || interval > 24 {
			return &graphError{"BAD_REQUEST", "Interval must be in range 1 to 24"}
		}
	}

	value, ok := body["emailSettings"]
	if !ok {
		return nil
	}
	emailSettings, isObject := value.(map[string]any)
	if !isObject {
		return &graphError{"Invalid model state", "A null value was found for the property named 'emailSettings'"}
	}
	if _, hasType := emailSettings["@odata.type"]; hasType {
		return &graphError{"Invalid model state", "A type named 'microsoft.graph.emailSettings' could not be resolved by the model."}
	}
	if domain, hasDomain := emailSettings["senderDomain"]; hasDomain {
		domainString, isString := domain.(string)
		if !isString {
			return &graphError{"Invalid model state", "A null value was found for the property named 'senderDomain'"}
		}
		if !verifiedSenderDomains[domainString] {
			return &graphError{"BAD_REQUEST", "Provided domain is not verified for tenant or is not an eligible domain"}
		}
	}
	return nil
}

func (m *LifecycleWorkflowSettingsMock) RegisterMocks() {
	settings, err := loadSeedSettings()
	mockState.Lock()
	mockState.settings = settings
	mockState.patches = nil
	mockState.seedErr = err
	mockState.Unlock()

	// Get lifecycle workflow settings - GET /identityGovernance/lifecycleWorkflows/settings
	httpmock.RegisterResponder("GET", settingsURL, func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		if mockState.seedErr != nil {
			return errorResponse(500, "InternalServerError", mockState.seedErr.Error())
		}
		return httpmock.NewJsonResponse(200, mockState.settings)
	})

	// Update lifecycle workflow settings - PATCH /identityGovernance/lifecycleWorkflows/settings
	httpmock.RegisterResponder("PATCH", settingsURL, func(req *http.Request) (*http.Response, error) {
		content, err := io.ReadAll(req.Body)
		if err != nil {
			return errorResponse(400, "BadRequest", "Invalid request body")
		}
		var requestBody map[string]any
		if err := json.Unmarshal(content, &requestBody); err != nil {
			return errorResponse(400, "BadRequest", "Invalid request body")
		}

		if invalid := validatePatch(requestBody); invalid != nil {
			return errorResponse(400, invalid.code, invalid.message)
		}

		mockState.Lock()
		defer mockState.Unlock()

		var recorded map[string]any
		_ = json.Unmarshal(content, &recorded)
		mockState.patches = append(mockState.patches, recorded)

		// Graph merges the PATCH into the current settings, including the properties inside emailSettings.
		if interval, ok := requestBody["workflowScheduleIntervalInHours"]; ok {
			mockState.settings["workflowScheduleIntervalInHours"] = interval
		}
		if emailSettings, ok := requestBody["emailSettings"].(map[string]any); ok {
			current, _ := mockState.settings["emailSettings"].(map[string]any)
			merged := map[string]any{}
			for key, value := range current {
				merged[key] = value
			}
			for key, value := range emailSettings {
				merged[key] = value
			}
			mockState.settings["emailSettings"] = merged
		}

		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *LifecycleWorkflowSettingsMock) RegisterErrorMocks() {
	// Error scenarios for testing
	forbidden := `{"error":{"code":"Authorization_RequestDenied","message":"Insufficient privileges to complete the operation."}}`
	httpmock.RegisterResponder("GET", settingsURL, httpmock.NewStringResponder(403, forbidden))
	httpmock.RegisterResponder("PATCH", settingsURL, httpmock.NewStringResponder(403, forbidden))
}

// RegisterNotFoundMocks returns 404 for every request, as Graph does when lifecycle workflows are unavailable in the tenant.
func (m *LifecycleWorkflowSettingsMock) RegisterNotFoundMocks() {
	notFound := `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`
	httpmock.RegisterResponder("GET", settingsURL, httpmock.NewStringResponder(404, notFound))
	httpmock.RegisterResponder("PATCH", settingsURL, httpmock.NewStringResponder(404, notFound))
}

func (m *LifecycleWorkflowSettingsMock) CleanupMockState() {
	mockState.Lock()
	mockState.settings = nil
	mockState.patches = nil
	mockState.seedErr = nil
	mockState.Unlock()
}

// PatchRequests returns the request bodies received by the PATCH responder, oldest first.
func (m *LifecycleWorkflowSettingsMock) PatchRequests() []map[string]any {
	mockState.Lock()
	defer mockState.Unlock()
	return append([]map[string]any{}, mockState.patches...)
}

// Settings returns a copy of the settings the mock currently stores.
func (m *LifecycleWorkflowSettingsMock) Settings() map[string]any {
	mockState.Lock()
	defer mockState.Unlock()
	settings := map[string]any{}
	for key, value := range mockState.settings {
		settings[key] = value
	}
	return settings
}
