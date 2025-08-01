package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	windowsUpdateRings map[string]map[string]interface{}
}

func init() {
	// Initialize mockState
	mockState.windowsUpdateRings = make(map[string]map[string]interface{})

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

// WindowsUpdateRingMock provides mock responses for Windows update ring operations
type WindowsUpdateRingMock struct{}

// RegisterMocks registers HTTP mock responses for Windows update ring operations
func (m *WindowsUpdateRingMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.windowsUpdateRings = make(map[string]map[string]interface{})
	mockState.Unlock()

	// Register GET for listing Windows update rings
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			rings := make([]map[string]interface{}, 0, len(mockState.windowsUpdateRings))
			for _, ring := range mockState.windowsUpdateRings {
				rings = append(rings, ring)
			}
			mockState.Unlock()

			response := map[string]interface{}{
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
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Windows update ring not found"}}`), nil
			}

			// Create response copy
			responseCopy := make(map[string]interface{})
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
			var requestBody map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Generate new update ring ID
			ringId := uuid.New().String()

			// Create update ring data - only include fields that were provided or have defaults
			ringData := map[string]interface{}{
				"@odata.type":                             "#microsoft.graph.windowsUpdateForBusinessConfiguration",
				"id":                                      ringId,
				"displayName":                             requestBody["displayName"],
				"microsoftUpdateServiceAllowed":           requestBody["microsoftUpdateServiceAllowed"],
				"driversExcluded":                         requestBody["driversExcluded"],
				"qualityUpdatesDeferralPeriodInDays":      requestBody["qualityUpdatesDeferralPeriodInDays"],
				"featureUpdatesDeferralPeriodInDays":      requestBody["featureUpdatesDeferralPeriodInDays"],
				"allowWindows11Upgrade":                   requestBody["allowWindows11Upgrade"],
				"skipChecksBeforeRestart":                 requestBody["skipChecksBeforeRestart"],
				"automaticUpdateMode":                     requestBody["automaticUpdateMode"],
				"featureUpdatesRollbackWindowInDays":      requestBody["featureUpdatesRollbackWindowInDays"],
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
				if schedule, ok := installationSchedule.(map[string]interface{}); ok {
					if activeHoursStart, hasStart := schedule["activeHoursStart"]; hasStart {
						ringData["installationSchedule"] = map[string]interface{}{
							"@odata.type":     "#microsoft.graph.windowsUpdateActiveHoursInstall",
							"activeHoursStart": activeHoursStart,
						}
						if activeHoursEnd, hasEnd := schedule["activeHoursEnd"]; hasEnd {
							ringData["installationSchedule"].(map[string]interface{})["activeHoursEnd"] = activeHoursEnd
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
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Windows update ring not found"}}`), nil
			}

			// Parse request body
			var requestBody map[string]interface{}
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Update ring data
			mockState.Lock()
			
			// Handle optional fields that might be removed (like going from maximal to minimal)
			// Check for specific field patterns to simulate real API behavior
			
			// For optional fields, if they're not in the request, remove them
			optionalFields := []string{"description", "businessReadyUpdatesOnly", "deliveryOptimizationMode", "prereleaseFeatures", "updateWeeks", "installationSchedule", "userPauseAccess", "userWindowsUpdateScanAccess", "updateNotificationLevel", "engagedRestartDeadlineInDays", "engagedRestartSnoozeScheduleInDays", "engagedRestartTransitionScheduleInDays", "autoRestartNotificationDismissal", "scheduleRestartWarningInHours", "scheduleImminentRestartWarningInMinutes", "featureUpdatesWillBeRolledBack", "qualityUpdatesWillBeRolledBack", "deadlineForFeatureUpdatesInDays", "deadlineForQualityUpdatesInDays", "deadlineGracePeriodInDays", "postponeRebootUntilAfterDeadline"}
			for _, field := range optionalFields {
				if _, hasField := requestBody[field]; !hasField {
					delete(ringData, field)
				}
			}
			
			for key, value := range requestBody {
				if value == nil {
					// If value is explicitly null, remove the field from the stored state
					delete(ringData, key)
				} else {
					ringData[key] = value
				}
			}
			
			// Handle individual field mappings for nested objects
			if featureUpdatesWillBeRolledBack, exists := requestBody["featureUpdatesWillBeRolledBack"]; exists {
				ringData["featureUpdatesWillBeRolledBack"] = featureUpdatesWillBeRolledBack
			}
			if qualityUpdatesWillBeRolledBack, exists := requestBody["qualityUpdatesWillBeRolledBack"]; exists {
				ringData["qualityUpdatesWillBeRolledBack"] = qualityUpdatesWillBeRolledBack
			}
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
			// Handle installation schedule for active hours  
			if installationSchedule, exists := requestBody["installationSchedule"]; exists {
				if schedule, ok := installationSchedule.(map[string]interface{}); ok {
					if activeHoursStart, hasStart := schedule["activeHoursStart"]; hasStart {
						ringData["installationSchedule"] = map[string]interface{}{
							"@odata.type":     "#microsoft.graph.windowsUpdateActiveHoursInstall",
							"activeHoursStart": activeHoursStart,
						}
						if activeHoursEnd, hasEnd := schedule["activeHoursEnd"]; hasEnd {
							ringData["installationSchedule"].(map[string]interface{})["activeHoursEnd"] = activeHoursEnd
						}
					}
				}
			}
			// Ensure the ID and @odata.type are preserved and update timestamp
			ringData["id"] = ringId
			ringData["@odata.type"] = "#microsoft.graph.windowsUpdateForBusinessConfiguration"
			ringData["lastModifiedDateTime"] = "2024-01-01T01:00:00Z"
			mockState.windowsUpdateRings[ringId] = ringData
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, ringData)
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
				return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Windows update ring not found"}}`), nil
			}

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register POST for assignments
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/[^/]+/assign$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			ringId := urlParts[len(urlParts)-2] // deviceConfigurations/{id}/assign

			// Parse request body to get assignments
			var requestBody map[string]interface{}
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
							if assignmentMap, ok := assignment.(map[string]interface{}); ok {
								if target, hasTarget := assignmentMap["target"].(map[string]interface{}); hasTarget {
									// Generate a unique assignment ID
									assignmentId := uuid.New().String()
									
									// Create assignment in the format the API returns
									// The API returns the target exactly as submitted but with additional metadata
									targetCopy := make(map[string]interface{})
									for k, v := range target {
										targetCopy[k] = v
									}
									
									graphAssignment := map[string]interface{}{
										"id": assignmentId,
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
			response := map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceConfigurations/assignments",
				"value":          []map[string]interface{}{}, // Empty assignments by default
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Dynamic mocks will handle all test cases
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *WindowsUpdateRingMock) RegisterErrorMocks() {
	// Register GET for listing Windows update rings (needed for uniqueness check)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]interface{}{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceConfigurations",
				"value":          []map[string]interface{}{}, // Empty list for error scenarios
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Register error response for creating Windows update ring with invalid data
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		factories.ErrorResponse(400, "BadRequest", "Validation error: Invalid display name"))

	// Register error response for Windows update ring not found
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/not-found-ring",
		factories.ErrorResponse(404, "ResourceNotFound", "Windows update ring not found"))
}