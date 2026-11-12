package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
)

// mockState tracks app collections independently for each protection policy.
var mockState struct {
	sync.Mutex
	apps map[string]map[string]map[string]any
}

// AndroidManagedMobileAppMock provides mock responses for Android managed mobile app operations.
type AndroidManagedMobileAppMock struct{}

// RegisterMocks registers collection reads and full-list targetApps replacements.
func (m *AndroidManagedMobileAppMock) RegisterMocks() {
	mockState.Lock()
	mockState.apps = make(map[string]map[string]map[string]any)
	mockState.Unlock()
	registerTestApps()
	httpmock.RegisterResponder(
		"GET",
		`=~^https://graph.microsoft.com/beta/deviceAppManagement/androidManagedAppProtections/[^/]+/apps(\?.+)?$`,
		func(req *http.Request) (*http.Response, error) {
			policyID := strings.Split(req.URL.Path, "/")[4]
			mockState.Lock()
			defer mockState.Unlock()
			apps := make([]map[string]any, 0, len(mockState.apps[policyID]))
			for _, app := range mockState.apps[policyID] {
				apps = append(apps, app)
			}
			return httpmock.NewJsonResponse(http.StatusOK, map[string]any{"value": apps})
		},
	)
	httpmock.RegisterResponder(
		"POST",
		`=~^https://graph.microsoft.com/beta/deviceAppManagement/managedAppPolicies/[^/]+/targetApps$`,
		func(req *http.Request) (*http.Response, error) {
			var payload struct {
				Apps []map[string]any `json:"apps"`
			}
			if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
				return nil, fmt.Errorf("decode targetApps request: %w", err)
			}
			apps := make(map[string]map[string]any, len(payload.Apps))
			for _, app := range payload.Apps {
				id, ok := app["id"].(string)
				if !ok || id == "" {
					id = uuid.NewString()
					app["id"] = id
				}
				if app["version"] == nil {
					app["version"] = "1.0"
				}
				apps[id] = app
			}
			policyID := strings.Split(req.URL.Path, "/")[4]
			mockState.Lock()
			mockState.apps[policyID] = apps
			mockState.Unlock()
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		},
	)
}

// RegisterErrorMocks registers the targetApps failure used by error lifecycle tests.
func (m *AndroidManagedMobileAppMock) RegisterErrorMocks() {
	httpmock.RegisterResponder(
		"POST",
		`=~^https://graph.microsoft.com/beta/deviceAppManagement/managedAppPolicies/[^/]+/targetApps$`,
		factories.ErrorResponse(400, "BadRequest", "Error creating Android managed mobile app"),
	)
}

// registerTestApps registers predefined test apps
func registerTestApps() {
	// Minimal app with only required attributes
	minimalAppId := "00000000-0000-0000-0000-000000000001"
	minimalAppData := map[string]any{
		"id":      minimalAppId,
		"version": "1.0",
		"mobileAppIdentifier": map[string]any{
			"@odata.type": "#microsoft.graph.androidMobileAppIdentifier",
			"packageId":   "com.example.testapp",
		},
		"@odata.context":  "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/androidManagedAppProtections('00000000-0000-0000-0000-000000000002')/apps/$entity",
		"@odata.type":     "#microsoft.graph.managedMobileApp",
		"createdDateTime": "2023-01-01T00:00:00Z",
	}

	// Maximal app with all attributes
	maximalAppId := "00000000-0000-0000-0000-000000000002"
	maximalAppData := map[string]any{
		"id":      maximalAppId,
		"version": "1.5",
		"mobileAppIdentifier": map[string]any{
			"@odata.type": "#microsoft.graph.androidMobileAppIdentifier",
			"packageId":   "com.example.complexapp",
		},
		"@odata.context":  "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/androidManagedAppProtections('00000000-0000-0000-0000-000000000003')/apps/$entity",
		"@odata.type":     "#microsoft.graph.managedMobileApp",
		"createdDateTime": "2023-01-01T00:00:00Z",
	}

	// Store apps in mock state
	mockState.Lock()
	mockState.apps["00000000-0000-0000-0000-000000000002"] = map[string]map[string]any{
		minimalAppId: minimalAppData,
	}
	mockState.apps["00000000-0000-0000-0000-000000000003"] = map[string]map[string]any{
		maximalAppId: maximalAppData,
	}
	mockState.Unlock()
}
