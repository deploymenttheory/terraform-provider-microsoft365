package mocks

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	windowsUpdateRings map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.windowsUpdateRings = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))

	// Register with global registry
	mocks.GlobalRegistry.Register("windows_update_ring", &WindowsUpdateRingMock{})
}

// WindowsUpdateRingMock provides mock responses for Windows update ring operations
type WindowsUpdateRingMock struct{}

// Ensure WindowsUpdateRingMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*WindowsUpdateRingMock)(nil)

// RegisterMocks registers HTTP mock responses for Windows update ring operations
func (m *WindowsUpdateRingMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.windowsUpdateRings = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for listing Windows update rings
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			rings := make([]map[string]any, 0, len(mockState.windowsUpdateRings))
			for _, ring := range mockState.windowsUpdateRings {
				rings = append(rings, ring)
			}
			mockState.Unlock()

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceConfigurations",
				"value":          rings,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for individual Windows update ring
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			ringId := urlParts[len(urlParts)-1]

			mockState.Lock()
			ringData, exists := mockState.windowsUpdateRings[ringId]
			mockState.Unlock()

			if !exists {
				// Check for special test IDs
				switch {
				case strings.Contains(ringId, "minimal"):
					response, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_create", "get_windows_update_ring_minimal.json"))
					if err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
					}
					response["id"] = ringId
					return factories.SuccessResponse(200, response)(req)
				case strings.Contains(ringId, "maximal"):
					response, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_create", "get_windows_update_ring_maximal.json"))
					if err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
					}
					response["id"] = ringId
					return factories.SuccessResponse(200, response)(req)
				default:
					errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_windows_update_ring_not_found.json"))
					return httpmock.NewJsonResponse(404, errorResponse)
				}
			}

			// Create response copy
			responseCopy := make(map[string]any)
			for k, v := range ringData {
				responseCopy[k] = v
			}

			// Check if expand=assignments is requested
			expandParam := req.URL.Query().Get("$expand")
			if strings.Contains(expandParam, "assignments") {
				// Include assignments if they exist in the ring data
				if assignments, hasAssignments := ringData["assignments"]; hasAssignments && assignments != nil {
					if assignmentList, ok := assignments.([]interface{}); ok && len(assignmentList) > 0 {
						responseCopy["assignments"] = assignments
					} else {
						// If assignments array is empty, return empty array (not null)
						responseCopy["assignments"] = []interface{}{}
					}
				} else {
					// If no assignments stored, return empty array (not null)
					responseCopy["assignments"] = []interface{}{}
				}
			}

			return httpmock.NewJsonResponse(200, responseCopy)
		})

	// Register POST for creating Windows update ring
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		func(req *http.Request) (*http.Response, error) {
			// Parse request body
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Generate new update ring ID
			ringId := uuid.New().String()

			// Create update ring data - only include fields that were provided or have defaults
			ringData := map[string]any{
				"@odata.type":                        "#microsoft.graph.windowsUpdateForBusinessConfiguration",
				"id":                                 ringId,
				"displayName":                        requestBody["displayName"],
				"microsoftUpdateServiceAllowed":      requestBody["microsoftUpdateServiceAllowed"],
				"driversExcluded":                    requestBody["driversExcluded"],
				"qualityUpdatesDeferralPeriodInDays": requestBody["qualityUpdatesDeferralPeriodInDays"],
				"featureUpdatesDeferralPeriodInDays": requestBody["featureUpdatesDeferralPeriodInDays"],
				"allowWindows11Upgrade":              requestBody["allowWindows11Upgrade"],
				"skipChecksBeforeRestart":            requestBody["skipChecksBeforeRestart"],
				"automaticUpdateMode":                requestBody["automaticUpdateMode"],
				"featureUpdatesRollbackWindowInDays": requestBody["featureUpdatesRollbackWindowInDays"],
			}

			// Add optional fields only if provided in request
			if description, exists := requestBody["description"]; exists {
				ringData["description"] = description
			}
			if roleScopeTagIds, exists := requestBody["roleScopeTagIds"]; exists {
				ringData["roleScopeTagIds"] = roleScopeTagIds
			} else {
				ringData["roleScopeTagIds"] = []string{"0"} // Default value
			}
			if businessReadyUpdatesOnly, exists := requestBody["businessReadyUpdatesOnly"]; exists {
				ringData["businessReadyUpdatesOnly"] = businessReadyUpdatesOnly
			}
			if deliveryOptimizationMode, exists := requestBody["deliveryOptimizationMode"]; exists {
				ringData["deliveryOptimizationMode"] = deliveryOptimizationMode
			}
			if prereleaseFeatures, exists := requestBody["prereleaseFeatures"]; exists {
				ringData["prereleaseFeatures"] = prereleaseFeatures
			}
			if updateWeeks, exists := requestBody["updateWeeks"]; exists {
				ringData["updateWeeks"] = updateWeeks
			}
			// Handle installation schedule for active hours
			if installationSchedule, exists := requestBody["installationSchedule"]; exists {
				if schedule, ok := installationSchedule.(map[string]any); ok {
					if activeHoursStart, hasStart := schedule["activeHoursStart"]; hasStart {
						ringData["installationSchedule"] = map[string]any{
							"@odata.type":      "#microsoft.graph.windowsUpdateActiveHoursInstall",
							"activeHoursStart": activeHoursStart,
						}
						if activeHoursEnd, hasEnd := schedule["activeHoursEnd"]; hasEnd {
							ringData["installationSchedule"].(map[string]any)["activeHoursEnd"] = activeHoursEnd
						}
					}
				}
			}
			if userPauseAccess, exists := requestBody["userPauseAccess"]; exists {
				ringData["userPauseAccess"] = userPauseAccess
			}
			if userWindowsUpdateScanAccess, exists := requestBody["userWindowsUpdateScanAccess"]; exists {
				ringData["userWindowsUpdateScanAccess"] = userWindowsUpdateScanAccess
			}
			if updateNotificationLevel, exists := requestBody["updateNotificationLevel"]; exists {
				ringData["updateNotificationLevel"] = updateNotificationLevel
			}
			if engagedRestartDeadlineInDays, exists := requestBody["engagedRestartDeadlineInDays"]; exists {
				ringData["engagedRestartDeadlineInDays"] = engagedRestartDeadlineInDays
			}
			if engagedRestartSnoozeScheduleInDays, exists := requestBody["engagedRestartSnoozeScheduleInDays"]; exists {
				ringData["engagedRestartSnoozeScheduleInDays"] = engagedRestartSnoozeScheduleInDays
			}
			if engagedRestartTransitionScheduleInDays, exists := requestBody["engagedRestartTransitionScheduleInDays"]; exists {
				ringData["engagedRestartTransitionScheduleInDays"] = engagedRestartTransitionScheduleInDays
			}
			if autoRestartNotificationDismissal, exists := requestBody["autoRestartNotificationDismissal"]; exists {
				ringData["autoRestartNotificationDismissal"] = autoRestartNotificationDismissal
			}
			if scheduleRestartWarningInHours, exists := requestBody["scheduleRestartWarningInHours"]; exists {
				ringData["scheduleRestartWarningInHours"] = scheduleRestartWarningInHours
			}
			if scheduleImminentRestartWarningInMinutes, exists := requestBody["scheduleImminentRestartWarningInMinutes"]; exists {
				ringData["scheduleImminentRestartWarningInMinutes"] = scheduleImminentRestartWarningInMinutes
			}
			// Handle uninstall settings - these are mapped to individual fields
			if featureUpdatesWillBeRolledBack, exists := requestBody["featureUpdatesWillBeRolledBack"]; exists {
				ringData["featureUpdatesWillBeRolledBack"] = featureUpdatesWillBeRolledBack
			}
			if qualityUpdatesWillBeRolledBack, exists := requestBody["qualityUpdatesWillBeRolledBack"]; exists {
				ringData["qualityUpdatesWillBeRolledBack"] = qualityUpdatesWillBeRolledBack
			}
			// Handle deadline settings - these are mapped to individual fields
			if deadlineForFeatureUpdatesInDays, exists := requestBody["deadlineForFeatureUpdatesInDays"]; exists {
				ringData["deadlineForFeatureUpdatesInDays"] = deadlineForFeatureUpdatesInDays
			}
			if deadlineForQualityUpdatesInDays, exists := requestBody["deadlineForQualityUpdatesInDays"]; exists {
				ringData["deadlineForQualityUpdatesInDays"] = deadlineForQualityUpdatesInDays
			}
			if deadlineGracePeriodInDays, exists := requestBody["deadlineGracePeriodInDays"]; exists {
				ringData["deadlineGracePeriodInDays"] = deadlineGracePeriodInDays
			}
			if postponeRebootUntilAfterDeadline, exists := requestBody["postponeRebootUntilAfterDeadline"]; exists {
				ringData["postponeRebootUntilAfterDeadline"] = postponeRebootUntilAfterDeadline
			}

			// Add computed fields that are always returned by the API
			ringData["createdDateTime"] = "2024-01-01T00:00:00Z"
			ringData["lastModifiedDateTime"] = "2024-01-01T00:00:00Z"

			// Initialize assignments as empty array
			ringData["assignments"] = []interface{}{}

			// Store in mock state
			mockState.Lock()
			mockState.windowsUpdateRings[ringId] = ringData
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, ringData)
		})

	// Register PATCH for updating Windows update ring
	httpmock.RegisterResponder("PATCH", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			ringId := urlParts[len(urlParts)-1]

			mockState.Lock()
			ringData, exists := mockState.windowsUpdateRings[ringId]
			mockState.Unlock()

			if !exists {
				errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_windows_update_ring_not_found.json"))
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			// Parse request body
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_create", "post_windows_update_ring_error.json"))
				return httpmock.NewJsonResponse(400, errorResponse)
			}

			// Load update template
			updatedRing, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_update", "get_windows_update_ring_updated.json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
			}

			// Start with existing data
			for k, v := range ringData {
				updatedRing[k] = v
			}

			// Apply updates from request body
			for k, v := range requestBody {
				updatedRing[k] = v
			}

			// Store updated state
			mockState.windowsUpdateRings[ringId] = updatedRing
			mockState.Unlock()

			return factories.SuccessResponse(200, updatedRing)(req)
		})

	// Register DELETE for removing Windows update ring
	httpmock.RegisterResponder("DELETE", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			ringId := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.windowsUpdateRings[ringId]
			if exists {
				delete(mockState.windowsUpdateRings, ringId)
			}
			mockState.Unlock()

			if !exists {
				errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_windows_update_ring_not_found.json"))
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register POST for assignments
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/[^/]+/assign$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			ringId := urlParts[len(urlParts)-2] // deviceConfigurations/{id}/assign

			// Parse request body to get assignments
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Store assignments in the ring
			mockState.Lock()
			if ringData, exists := mockState.windowsUpdateRings[ringId]; exists {
				if assignments, hasAssignments := requestBody["assignments"]; hasAssignments && assignments != nil {
					assignmentList := assignments.([]interface{})
					if len(assignmentList) > 0 {
						// Extract the actual assignment data from the request
						graphAssignments := []interface{}{}
						for _, assignment := range assignmentList {
							if assignmentMap, ok := assignment.(map[string]any); ok {
								if target, hasTarget := assignmentMap["target"].(map[string]any); hasTarget {
									// Generate a unique assignment ID
									assignmentId := uuid.New().String()

									// Create assignment in the format the API returns
									// The API returns the target exactly as submitted but with additional metadata
									targetCopy := make(map[string]any)
									for k, v := range target {
										targetCopy[k] = v
									}

									graphAssignment := map[string]any{
										"id":     assignmentId,
										"target": targetCopy,
									}
									graphAssignments = append(graphAssignments, graphAssignment)
								}
							}
						}
						ringData["assignments"] = graphAssignments
					} else {
						// Set empty assignments array instead of deleting
						ringData["assignments"] = []interface{}{}
					}
				} else {
					// Set empty assignments array instead of deleting
					ringData["assignments"] = []interface{}{}
				}
				mockState.windowsUpdateRings[ringId] = ringData
			}
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register GET for assignments
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/[^/]+/assignments$`,
		func(req *http.Request) (*http.Response, error) {
			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceConfigurations/assignments",
				"value":          []map[string]any{}, // Empty assignments by default
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Dynamic mocks will handle all test cases
}

// CleanupMockState clears the mock state for clean test runs
func (m *WindowsUpdateRingMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()

	// Clear all stored Windows update rings
	for id := range mockState.windowsUpdateRings {
		delete(mockState.windowsUpdateRings, id)
	}
}

// loadJSONResponse loads a JSON response from a file
func (m *WindowsUpdateRingMock) loadJSONResponse(filePath string) (map[string]any, error) {
	var response map[string]any

	content, err := os.ReadFile(filePath)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(content, &response)
	return response, err
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *WindowsUpdateRingMock) RegisterErrorMocks() {
	// Register GET for listing Windows update rings (needed for uniqueness check)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceConfigurations",
				"value":          []map[string]any{}, // Empty list for error scenarios
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Register error response for creating Windows update ring with invalid data
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		func(req *http.Request) (*http.Response, error) {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_create", "post_windows_update_ring_error.json"))
			return httpmock.NewJsonResponse(400, errorResponse)
		})

	// Register error response for Windows update ring not found
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/]+)$`,
		func(req *http.Request) (*http.Response, error) {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_windows_update_ring_not_found.json"))
			return httpmock.NewJsonResponse(404, errorResponse)
		})
}
