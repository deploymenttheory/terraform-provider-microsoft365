package androidManagedAppProtectionMocks

import (
	"net/http"
	"sync"

	"github.com/jarcoal/httpmock"
)

const (
	baseURL   = "https://graph.microsoft.com/beta/deviceAppManagement/androidManagedAppProtections"
	mockIDMin = "00000000-0000-0000-0000-000000000001"
	mockIDMax = "00000000-0000-0000-0000-000000000002"
)

// AndroidManagedAppProtectionMock manages mock HTTP responses for unit tests.
type AndroidManagedAppProtectionMock struct {
	mu    sync.Mutex
	state map[string]string
}

const minimalResponse = `{
	"@odata.type": "#microsoft.graph.androidManagedAppProtection",
	"id": "` + mockIDMin + `",
	"displayName": "unit-test-android-managed-app-protection-minimal",
	"description": "",
	"createdDateTime": "2024-01-01T00:00:00Z",
	"lastModifiedDateTime": "2024-01-01T00:00:00Z",
	"version": "\"1\"",
	"isAssigned": false,
	"deployedAppCount": 0,
	"printBlocked": false,
	"allowedInboundDataTransferSources": "allApps",
	"allowedOutboundClipboardSharingLevel": "allApps",
	"allowedOutboundDataTransferDestinations": "allApps",
	"organizationalCredentialsRequired": false,
	"dataBackupBlocked": false,
	"deviceComplianceRequired": true,
	"managedBrowserToOpenLinksRequired": false,
	"saveAsBlocked": false,
	"pinRequired": true,
	"maximumPinRetries": 5,
	"simplePinBlocked": false,
	"minimumPinLength": 4,
	"pinCharacterSet": "numeric",
	"periodBeforePinReset": "P365D",
	"contactSyncBlocked": false,
	"fingerprintBlocked": false,
	"disableAppPinIfDevicePinIsSet": false,
	"managedBrowser": "notConfigured",
	"screenCaptureBlocked": false,
	"disableAppEncryptionIfDeviceEncryptionIsEnabled": false,
	"encryptAppData": true,
	"periodOfflineBeforeAccessCheck": "P30D",
	"periodOnlineBeforeAccessCheck": "PT30M",
	"periodOfflineBeforeWipeIsEnforced": "P90D",
	"allowedDataStorageLocations": []
}`

const maximalResponse = `{
	"@odata.type": "#microsoft.graph.androidManagedAppProtection",
	"id": "` + mockIDMax + `",
	"displayName": "unit-test-android-managed-app-protection-maximal",
	"description": "Maximal test configuration for Android managed app protection",
	"createdDateTime": "2024-01-01T00:00:00Z",
	"lastModifiedDateTime": "2024-01-01T00:00:00Z",
	"version": "\"2\"",
	"isAssigned": false,
	"deployedAppCount": 0,
	"printBlocked": true,
	"allowedInboundDataTransferSources": "none",
	"allowedOutboundClipboardSharingLevel": "blocked",
	"allowedOutboundDataTransferDestinations": "none",
	"organizationalCredentialsRequired": false,
	"dataBackupBlocked": true,
	"deviceComplianceRequired": true,
	"managedBrowserToOpenLinksRequired": false,
	"saveAsBlocked": true,
	"pinRequired": true,
	"maximumPinRetries": 10,
	"simplePinBlocked": true,
	"minimumPinLength": 6,
	"pinCharacterSet": "alphanumericAndSymbol",
	"periodBeforePinReset": "P30D",
	"contactSyncBlocked": true,
	"fingerprintBlocked": true,
	"disableAppPinIfDevicePinIsSet": false,
	"managedBrowser": "notConfigured",
	"screenCaptureBlocked": true,
	"disableAppEncryptionIfDeviceEncryptionIsEnabled": false,
	"encryptAppData": true,
	"minimumRequiredOsVersion": "9.0",
	"minimumWarningOsVersion": "8.0",
	"minimumRequiredAppVersion": "2.0.0",
	"minimumWarningAppVersion": "1.9.0",
	"minimumRequiredPatchVersion": "2024-01-01",
	"minimumWarningPatchVersion": "2023-12-01",
	"periodOfflineBeforeAccessCheck": "P30D",
	"periodOnlineBeforeAccessCheck": "PT30M",
	"periodOfflineBeforeWipeIsEnforced": "P30D",
	"allowedDataStorageLocations": ["oneDriveForBusiness", "sharePoint"]
}`

const errorResponse = `{
	"error": {
		"code": "BadRequest",
		"message": "Invalid Android Managed App Protection data"
	}
}`

// RegisterMocks registers all standard HTTP mock responders.
func (m *AndroidManagedAppProtectionMock) RegisterMocks() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state = make(map[string]string)

	httpmock.RegisterResponder("POST", baseURL,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, minimalResponse), nil
		},
	)

	httpmock.RegisterResponder("GET", baseURL+"/"+mockIDMin,
		httpmock.NewStringResponder(http.StatusOK, minimalResponse),
	)
	httpmock.RegisterResponder("PATCH", baseURL+"/"+mockIDMin,
		httpmock.NewStringResponder(http.StatusOK, minimalResponse),
	)
	httpmock.RegisterResponder("DELETE", baseURL+"/"+mockIDMin,
		httpmock.NewStringResponder(http.StatusNoContent, ""),
	)

	httpmock.RegisterResponder("GET", baseURL+"/"+mockIDMax,
		httpmock.NewStringResponder(http.StatusOK, maximalResponse),
	)
	httpmock.RegisterResponder("PATCH", baseURL+"/"+mockIDMax,
		httpmock.NewStringResponder(http.StatusOK, maximalResponse),
	)
	httpmock.RegisterResponder("DELETE", baseURL+"/"+mockIDMax,
		httpmock.NewStringResponder(http.StatusNoContent, ""),
	)
}

// RegisterErrorMocks registers HTTP mock responders that return errors.
func (m *AndroidManagedAppProtectionMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("POST", baseURL,
		httpmock.NewStringResponder(http.StatusBadRequest, errorResponse),
	)
}

// CleanupMockState resets the mock state between tests.
func (m *AndroidManagedAppProtectionMock) CleanupMockState() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state = make(map[string]string)
}
