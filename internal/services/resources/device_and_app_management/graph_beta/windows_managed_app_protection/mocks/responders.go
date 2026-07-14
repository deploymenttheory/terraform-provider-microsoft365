package windowsManagedAppProtectionMocks

import (
	"net/http"
	"sync"

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
	state map[string]string // tracks which mock ID is "active"
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
	defer m.mu.Unlock()
	m.state = make(map[string]string)

	// POST - Create minimal
	httpmock.RegisterResponder("POST", baseURL,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, minimalResponse), nil
		},
	)

	// GET - Read minimal
	httpmock.RegisterResponder("GET", baseURL+"/"+mockIDMin,
		httpmock.NewStringResponder(http.StatusOK, minimalResponse),
	)

	// PATCH - Update minimal
	httpmock.RegisterResponder("PATCH", baseURL+"/"+mockIDMin,
		httpmock.NewStringResponder(http.StatusOK, minimalResponse),
	)

	// DELETE - Delete minimal
	httpmock.RegisterResponder("DELETE", baseURL+"/"+mockIDMin,
		httpmock.NewStringResponder(http.StatusNoContent, ""),
	)

	// GET - Read maximal
	httpmock.RegisterResponder("GET", baseURL+"/"+mockIDMax,
		httpmock.NewStringResponder(http.StatusOK, maximalResponse),
	)

	// PATCH - Update maximal
	httpmock.RegisterResponder("PATCH", baseURL+"/"+mockIDMax,
		httpmock.NewStringResponder(http.StatusOK, maximalResponse),
	)

	// DELETE - Delete maximal
	httpmock.RegisterResponder("DELETE", baseURL+"/"+mockIDMax,
		httpmock.NewStringResponder(http.StatusNoContent, ""),
	)
}

// RegisterErrorMocks registers HTTP mock responders that return errors.
func (m *WindowsManagedAppProtectionMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", baseURL,
		httpmock.NewStringResponder(http.StatusBadRequest, errorResponse),
	)
}

// CleanupMockState resets the mock state between tests.
func (m *WindowsManagedAppProtectionMock) CleanupMockState() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state = make(map[string]string)
}
