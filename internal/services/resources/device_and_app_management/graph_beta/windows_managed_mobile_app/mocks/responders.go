package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	apps map[string]map[string]interface{}
}

func init() {
	// Initialize mockState
	mockState.apps = make(map[string]map[string]interface{})

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// WindowsManagedMobileAppMock provides mock responses for Windows managed mobile app operations
type WindowsManagedMobileAppMock struct{}

// RegisterMocks registers HTTP mock responses for Windows managed mobile app operations
func (m *WindowsManagedMobileAppMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.apps = make(map[string]map[string]interface{})
	mockState.Unlock()

	// Register test apps
	registerTestApps()

	// Register GET for app by protection ID and app ID
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceAppManagement/windowsManagedAppProtections/[^/]+/apps/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			appId := urlParts[len(urlParts)-1]

			mockState.Lock()
			appData, exists := mockState.apps[appId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"App not found"}}`), nil
			}

			return httpmock.NewJsonResponse(200, appData)
		})

	// Register GET for listing apps
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceAppManagement/windowsManagedAppProtections/[^/]+/apps(\?.+)?$`,
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			defer mockState.Unlock()

			apps := make([]map[string]interface{}, 0, len(mockState.apps))
			for _, app := range mockState.apps {
				apps = append(apps, app)
			}

			response := map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/windowsManagedAppProtections('00000000-0000-0000-0000-000000000002')/apps",
				"value":          apps,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register POST for creating apps
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/deviceAppManagement/windowsManagedAppProtections/[^/]+/apps$`,
		func(req *http.Request) (*http.Response, error) {
			var appData map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&appData)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Extract protection ID from URL
			urlParts := strings.Split(req.URL.Path, "/")
			protectionId := urlParts[len(urlParts)-2]

			// Validate required fields
			if mobileAppId, ok := appData["mobileAppIdentifier"].(map[string]interface{}); !ok || mobileAppId["windowsAppId"] == nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"mobileAppIdentifier with windowsAppId is required"}}`), nil
			}

			// Generate ID if not provided
			if appData["id"] == nil {
				appData["id"] = uuid.New().String()
			}

			// Set computed fields
			now := time.Now().Format(time.RFC3339)
			appData["createdDateTime"] = now
			appData["@odata.context"] = "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/windowsManagedAppProtections('" + protectionId + "')/apps/$entity"
			appData["@odata.type"] = "#microsoft.graph.managedMobileApp"

			// Set default version if not provided
			if appData["version"] == nil {
				appData["version"] = "1.0"
			}

			// Ensure mobile app identifier has correct odata type
			if mobileAppId, ok := appData["mobileAppIdentifier"].(map[string]interface{}); ok {
				mobileAppId["@odata.type"] = "#microsoft.graph.windowsAppIdentifier"
			}

			// Store app in mock state
			appId := appData["id"].(string)
			mockState.Lock()
			mockState.apps[appId] = appData
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, appData)
		})

	// Register PATCH for updating apps
	httpmock.RegisterResponder("PATCH", `=~^https://graph.microsoft.com/beta/deviceAppManagement/windowsManagedAppProtections/[^/]+/apps/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			appId := urlParts[len(urlParts)-1]

			mockState.Lock()
			appData, exists := mockState.apps[appId]
			mockState.Unlock()

			if !exists {
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"App not found"}}`), nil
			}

			var updateData map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&updateData)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Update app data
			mockState.Lock()

			// Apply the updates
			for k, v := range updateData {
				appData[k] = v
			}

			mockState.apps[appId] = appData
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, appData)
		})

	// Register DELETE for removing apps
	httpmock.RegisterResponder("DELETE", `=~^https://graph.microsoft.com/beta/deviceAppManagement/windowsManagedAppProtections/[^/]+/apps/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			appId := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.apps[appId]
			if exists {
				delete(mockState.apps, appId)
			}
			mockState.Unlock()

			// Return 204 No Content for successful deletion
			return httpmock.NewStringResponse(204, ""), nil
		})
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *WindowsManagedMobileAppMock) RegisterErrorMocks() {
	// Register error response for app creation
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/deviceAppManagement/windowsManagedAppProtections/[^/]+/apps$`,
		factories.ErrorResponse(400, "BadRequest", "Error creating Windows managed mobile app"))

	// Register error response for app not found
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceAppManagement/windowsManagedAppProtections/[^/]+/apps/not-found-app$`,
		factories.ErrorResponse(404, "ResourceNotFound", "App not found"))

	// Register error response for invalid protection ID
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/deviceAppManagement/windowsManagedAppProtections/invalid-id/apps$`,
		factories.ErrorResponse(400, "BadRequest", "Invalid managed app protection ID"))
}

// registerTestApps registers predefined test apps
func registerTestApps() {
	// Minimal app with only required attributes
	minimalAppId := "00000000-0000-0000-0000-000000000001"
	minimalAppData := map[string]interface{}{
		"id":      minimalAppId,
		"version": "1.0",
		"mobileAppIdentifier": map[string]interface{}{
			"@odata.type":   "#microsoft.graph.windowsAppIdentifier",
			"windowsAppId": "com.example.testapp",
		},
		"@odata.context":  "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/windowsManagedAppProtections('00000000-0000-0000-0000-000000000002')/apps/$entity",
		"@odata.type":     "#microsoft.graph.managedMobileApp",
		"createdDateTime": "2023-01-01T00:00:00Z",
	}

	// Maximal app with all attributes
	maximalAppId := "00000000-0000-0000-0000-000000000002"
	maximalAppData := map[string]interface{}{
		"id":      maximalAppId,
		"version": "1.5",
		"mobileAppIdentifier": map[string]interface{}{
			"@odata.type":   "#microsoft.graph.windowsAppIdentifier",
			"windowsAppId": "com.example.complexapp",
		},
		"@odata.context":  "https://graph.microsoft.com/beta/$metadata#deviceAppManagement/windowsManagedAppProtections('00000000-0000-0000-0000-000000000003')/apps/$entity",
		"@odata.type":     "#microsoft.graph.managedMobileApp",
		"createdDateTime": "2023-01-01T00:00:00Z",
	}

	// Store apps in mock state
	mockState.Lock()
	mockState.apps[minimalAppId] = minimalAppData
	mockState.apps[maximalAppId] = maximalAppData
	mockState.Unlock()
}