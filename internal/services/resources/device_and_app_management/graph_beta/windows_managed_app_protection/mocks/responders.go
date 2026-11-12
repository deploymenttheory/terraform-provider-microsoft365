package windowsManagedAppProtectionMocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"sync"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

const (
	baseURL   = "https://graph.microsoft.com/beta/deviceAppManagement/windowsManagedAppProtections"
	mockIDMin = "00000000-0000-0000-0000-000000000001"
	mockIDMax = "00000000-0000-0000-0000-000000000002"
)

// WindowsManagedAppProtectionMock manages mock HTTP responses for unit tests.
type WindowsManagedAppProtectionMock struct {
	mu    sync.Mutex
	state map[string]map[string]any
}

// minimalResponse is the API response for a minimal policy configuration.
const minimalResponse = `{
	"@odata.type": "#microsoft.graph.windowsManagedAppProtection",
	"id": "` + mockIDMin + `",
	"displayName": "unit-test-windows-managed-app-protection-minimal",
	"description": "",
	"createdDateTime": "2024-01-01T00:00:00Z",
	"lastModifiedDateTime": "2024-01-01T00:00:00Z",
	"version": "\"1\"",
	"isAssigned": false,
	"deployedAppCount": 0,
	"printBlocked": false,
	"allowedInboundDataTransferSources": "allApps",
	"allowedOutboundClipboardSharingLevel": "anyDestinationAnySource",
	"allowedOutboundDataTransferDestinations": "allApps",
	"maximumAllowedDeviceThreatLevel": "notConfigured",
	"mobileThreatDefenseRemediationAction": "block",
	"periodOfflineBeforeWipeIsEnforced": "P90D",
	"periodOfflineBeforeAccessCheck": "P30D",
	"roleScopeTagIds": ["0"]
}`

// maximalResponse is the API response for a maximal policy configuration.
const maximalResponse = `{
	"@odata.type": "#microsoft.graph.windowsManagedAppProtection",
	"id": "` + mockIDMax + `",
	"displayName": "unit-test-windows-managed-app-protection-maximal",
	"description": "Maximal test configuration for Windows managed app protection",
	"createdDateTime": "2024-01-01T00:00:00Z",
	"lastModifiedDateTime": "2024-01-01T00:00:00Z",
	"version": "\"2\"",
	"isAssigned": false,
	"deployedAppCount": 0,
	"printBlocked": true,
	"allowedInboundDataTransferSources": "none",
	"allowedOutboundClipboardSharingLevel": "none",
	"allowedOutboundDataTransferDestinations": "none",
	"appActionIfUnableToAuthenticateUser": "block",
	"maximumAllowedDeviceThreatLevel": "low",
	"mobileThreatDefenseRemediationAction": "wipe",
	"minimumRequiredOsVersion": "10.0.19041",
	"minimumWarningOsVersion": "10.0.18363",
	"minimumWipeOsVersion": "10.0.17763",
	"minimumRequiredAppVersion": "1.0.0",
	"minimumWarningAppVersion": "1.1.0",
	"minimumWipeAppVersion": "0.9.0",
	"periodOfflineBeforeWipeIsEnforced": "P30D",
	"periodOfflineBeforeAccessCheck": "P7D",
	"roleScopeTagIds": ["0"]
}`

// errorResponse is returned when error mocks are registered.
const errorResponse = `{
	"error": {
		"code": "BadRequest",
		"message": "Invalid Windows Managed App Protection data"
	}
}`

// RegisterMocks registers all standard HTTP mock responders.
func (m *WindowsManagedAppProtectionMock) RegisterMocks() {
	m.mu.Lock()
	m.state = make(map[string]map[string]any)
	m.mu.Unlock()

	httpmock.RegisterResponder("POST", baseURL, func(req *http.Request) (*http.Response, error) {
		var payload map[string]any
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			return nil, fmt.Errorf("decode mock create request: %w", err)
		}
		fixture := minimalResponse
		if payload["printBlocked"] == true {
			fixture = maximalResponse
		}
		var stored map[string]any
		if err := json.Unmarshal([]byte(fixture), &stored); err != nil {
			return nil, fmt.Errorf("decode mock create fixture: %w", err)
		}
		for key, value := range payload {
			stored[key] = value
		}
		id := uuid.NewString()
		stored["id"] = id
		m.mu.Lock()
		defer m.mu.Unlock()
		m.state[id] = stored
		return newJSONResponse(http.StatusCreated, stored)
	})

	itemURL := "=~^" + regexp.QuoteMeta(baseURL) + "/[^/?]+(?:\\?.*)?$"
	httpmock.RegisterResponder("GET", itemURL, func(req *http.Request) (*http.Response, error) {
		m.mu.Lock()
		defer m.mu.Unlock()
		stored, exists := m.state[path.Base(req.URL.Path)]
		if !exists {
			return newJSONResponse(
				http.StatusNotFound,
				map[string]any{
					"error": map[string]string{
						"code":    "ResourceNotFound",
						"message": "Policy not found",
					},
				},
			)
		}
		return newJSONResponse(http.StatusOK, stored)
	})

	httpmock.RegisterResponder("PATCH", itemURL, func(req *http.Request) (*http.Response, error) {
		var payload map[string]any
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			return nil, fmt.Errorf("decode mock update request: %w", err)
		}
		m.mu.Lock()
		defer m.mu.Unlock()
		id := path.Base(req.URL.Path)
		stored, exists := m.state[id]
		if !exists {
			return newJSONResponse(
				http.StatusNotFound,
				map[string]any{
					"error": map[string]string{
						"code":    "ResourceNotFound",
						"message": "Policy not found",
					},
				},
			)
		}
		for key, value := range payload {
			stored[key] = value
		}
		stored["id"] = id
		return newJSONResponse(http.StatusOK, stored)
	})

	httpmock.RegisterResponder("DELETE", itemURL, func(req *http.Request) (*http.Response, error) {
		m.mu.Lock()
		defer m.mu.Unlock()
		delete(m.state, path.Base(req.URL.Path))
		return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
	})
}

// newJSONResponse returns a mock response with the content type required by the Graph SDK.
func newJSONResponse(status int, payload any) (*http.Response, error) {
	response, err := httpmock.NewJsonResponse(status, payload)
	if err != nil {
		return nil, fmt.Errorf("encode mock response: %w", err)
	}
	return response, nil
}

// RegisterErrorMocks registers HTTP mock responders that return errors.
func (m *WindowsManagedAppProtectionMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", baseURL,
		httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, json.RawMessage(errorResponse)),
	)
}

// CleanupMockState resets the mock state between tests.
func (m *WindowsManagedAppProtectionMock) CleanupMockState() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state = make(map[string]map[string]any)
}
