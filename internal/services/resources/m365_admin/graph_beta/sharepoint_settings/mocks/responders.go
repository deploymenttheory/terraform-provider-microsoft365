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

const settingsURL = "https://graph.microsoft.com/beta/admin/sharepoint/settings"

// readOnlyProperties are returned by Graph but cannot be changed with PATCH.
var readOnlyProperties = []string{"availableManagedPathsForSiteCreation"}

// collectionProperties are the collection settings; Graph rejects null for them.
var collectionProperties = []string{"allowedDomainGuidsForSyncApp", "excludedFileExtensionsForSyncApp", "sharingAllowedDomainList", "sharingBlockedDomainList"}

// enumProperties are the enum settings and the values Graph accepts for them.
var enumProperties = map[string]map[string]bool{
	"imageTaggingOption":           {"disabled": true, "basic": true, "enhanced": true},
	"sharingCapability":            {"disabled": true, "externalUserSharingOnly": true, "externalUserAndGuestSharing": true, "existingExternalUserSharingOnly": true},
	"sharingDomainRestrictionMode": {"none": true, "allowList": true, "blockList": true},
}

var mockState struct {
	sync.Mutex
	settings map[string]any
	patches  []map[string]any
	seedErr  error
}

func init() {
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("m365_admin_sharepoint_settings", &SharepointSettingsMock{})
}

// SharepointSettingsMock provides stateful httpmock responders for the singleton SharePoint settings.
type SharepointSettingsMock struct{}

var _ mocks.MockRegistrar = (*SharepointSettingsMock)(nil)

// loadSeedSettings loads the tenant settings that GET returns before any PATCH.
func loadSeedSettings() (map[string]any, error) {
	responsesPath := filepath.Join("tests", "responses", "validate_get", "get_sharepoint_settings_success.json")
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
	for _, property := range readOnlyProperties {
		if _, ok := body[property]; ok {
			return &graphError{"invalidRequest", fmt.Sprintf("Property '%s' is read-only and cannot be set.", property)}
		}
	}
	for _, property := range collectionProperties {
		if value, ok := body[property]; ok && value == nil {
			return &graphError{"invalidRequest", fmt.Sprintf("A null value was found for the property named '%s'", property)}
		}
	}
	for property, allowed := range enumProperties {
		if value, ok := body[property]; ok {
			text, isString := value.(string)
			if !isString || !allowed[text] {
				return &graphError{"invalidRequest", fmt.Sprintf("The value '%v' is not valid for property '%s'.", value, property)}
			}
		}
	}
	if value, ok := body["deletedUserPersonalSiteRetentionPeriodInDays"]; ok {
		days, isNumber := value.(float64)
		if !isNumber || days < 30 || days > 3650 {
			return &graphError{"invalidRequest", "The OneDrive retention period must be between 30 and 3650 days."}
		}
	}
	if value, ok := body["idleSessionSignOut"]; ok {
		policy, isObject := value.(map[string]any)
		if !isObject {
			return &graphError{"invalidRequest", "A null value was found for the property named 'idleSessionSignOut'"}
		}
		if _, hasType := policy["@odata.type"]; hasType {
			return &graphError{"invalidRequest", "A type named 'microsoft.graph.idleSessionSignOut' could not be resolved by the model."}
		}
		warn, hasWarn := policy["warnAfterInSeconds"].(float64)
		signOut, hasSignOut := policy["signOutAfterInSeconds"].(float64)
		if hasWarn && hasSignOut && warn >= signOut {
			return &graphError{"invalidRequest", "warnAfterInSeconds must be less than signOutAfterInSeconds."}
		}
	}
	return nil
}

func (m *SharepointSettingsMock) RegisterMocks() {
	settings, err := loadSeedSettings()
	mockState.Lock()
	mockState.settings = settings
	mockState.patches = nil
	mockState.seedErr = err
	mockState.Unlock()

	// Get SharePoint settings - GET /admin/sharepoint/settings
	httpmock.RegisterResponder("GET", settingsURL, func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		if mockState.seedErr != nil {
			return errorResponse(500, "InternalServerError", mockState.seedErr.Error())
		}
		return httpmock.NewJsonResponse(200, mockState.settings)
	})

	// Update SharePoint settings - PATCH /admin/sharepoint/settings
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

		// Graph keeps every property that is not in the request body and merges the properties inside idleSessionSignOut.
		for key, value := range requestBody {
			if key != "idleSessionSignOut" {
				mockState.settings[key] = value
				continue
			}
			current, _ := mockState.settings["idleSessionSignOut"].(map[string]any)
			merged := map[string]any{}
			for property, existing := range current {
				merged[property] = existing
			}
			for property, updated := range value.(map[string]any) {
				merged[property] = updated
			}
			mockState.settings["idleSessionSignOut"] = merged
		}

		// PATCH returns 200 OK with the updated settings.
		return httpmock.NewJsonResponse(200, mockState.settings)
	})
}

func (m *SharepointSettingsMock) RegisterErrorMocks() {
	// Error scenarios for testing
	forbidden := `{"error":{"code":"accessDenied","message":"Insufficient privileges to complete the operation."}}`
	httpmock.RegisterResponder("GET", settingsURL, httpmock.NewStringResponder(403, forbidden))
	httpmock.RegisterResponder("PATCH", settingsURL, httpmock.NewStringResponder(403, forbidden))
}

// RegisterNotFoundMocks returns 404 for every request, as Graph does when SharePoint Online is unavailable in the tenant.
func (m *SharepointSettingsMock) RegisterNotFoundMocks() {
	notFound := `{"error":{"code":"itemNotFound","message":"Resource not found"}}`
	httpmock.RegisterResponder("GET", settingsURL, httpmock.NewStringResponder(404, notFound))
	httpmock.RegisterResponder("PATCH", settingsURL, httpmock.NewStringResponder(404, notFound))
}

func (m *SharepointSettingsMock) CleanupMockState() {
	mockState.Lock()
	mockState.settings = nil
	mockState.patches = nil
	mockState.seedErr = nil
	mockState.Unlock()
}

// PatchRequests returns the request bodies received by the PATCH responder, oldest first.
func (m *SharepointSettingsMock) PatchRequests() []map[string]any {
	mockState.Lock()
	defer mockState.Unlock()
	return append([]map[string]any{}, mockState.patches...)
}

// Settings returns a copy of the settings the mock currently stores.
func (m *SharepointSettingsMock) Settings() map[string]any {
	mockState.Lock()
	defer mockState.Unlock()
	settings := map[string]any{}
	for key, value := range mockState.settings {
		settings[key] = value
	}
	return settings
}
