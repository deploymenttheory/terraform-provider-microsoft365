package mocks

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
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
	mockState.windowsUpdateRings = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_update_ring", &WindowsUpdateRingMock{})
}

// WindowsUpdateRingMock provides mock responses for Windows update ring operations
type WindowsUpdateRingMock struct{}

var _ mocks.MockRegistrar = (*WindowsUpdateRingMock)(nil)

// RegisterMocks registers HTTP mock responses for Windows update ring operations
func (m *WindowsUpdateRingMock) RegisterMocks() {
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
				content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_delete", "get_windows_update_ring_not_found.json"))
				var errorResponse map[string]any
				json.Unmarshal([]byte(content), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			// Determine which read JSON to load based on stored ring data
			scenarioFile := determineReadScenario(ringData)
			if scenarioFile == "" {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Could not determine scenario for GET request"}}`), nil
			}

			content, err := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_read", scenarioFile))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load read scenario JSON: `+err.Error()+`"}}`), nil
			}

			var response map[string]any
			if err := json.Unmarshal([]byte(content), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse read scenario JSON: `+err.Error()+`"}}`), nil
			}

			// Merge stored state into response
			response["id"] = ringId
			if displayName, hasName := ringData["displayName"]; hasName {
				response["displayName"] = displayName
			}
			if description, hasDesc := ringData["description"]; hasDesc {
				response["description"] = description
			}
			if assignments, hasAssignments := ringData["assignments"]; hasAssignments {
				response["assignments"] = assignments
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register POST for creating Windows update ring
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_create", "post_windows_update_ring_error.json"))
				var errorResponse map[string]any
				json.Unmarshal([]byte(content), &errorResponse)
				return httpmock.NewJsonResponse(400, errorResponse)
			}

			ringId := uuid.New().String()

			// Determine scenario JSON to load from validate_create/
			scenarioFile := determineCreateScenario(requestBody)
			if scenarioFile == "" {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Could not determine scenario for POST request"}}`), nil
			}

			content, err := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_create", scenarioFile))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load create scenario JSON: `+err.Error()+`"}}`), nil
			}

			var response map[string]any
			if err := json.Unmarshal([]byte(content), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse create scenario JSON: `+err.Error()+`"}}`), nil
			}

			// Override with request-specific values
			response["id"] = ringId
			if displayName, hasName := requestBody["displayName"]; hasName {
				response["displayName"] = displayName
			}
			if description, hasDesc := requestBody["description"]; hasDesc {
				response["description"] = description
			}

			mockState.Lock()
			mockState.windowsUpdateRings[ringId] = response
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, response)
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
				content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_delete", "get_windows_update_ring_not_found.json"))
				var errorResponse map[string]any
				json.Unmarshal([]byte(content), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_create", "post_windows_update_ring_error.json"))
				var errorResponse map[string]any
				json.Unmarshal([]byte(content), &errorResponse)
				return httpmock.NewJsonResponse(400, errorResponse)
			}

			// Determine which update JSON to load based on request body
			scenarioFile := determineUpdateScenario(requestBody, ringData)
			if scenarioFile == "" {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Could not determine scenario for PATCH request"}}`), nil
			}

			content, err := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_update", scenarioFile))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load update scenario JSON: `+err.Error()+`"}}`), nil
			}

			var response map[string]any
			if err := json.Unmarshal([]byte(content), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse update scenario JSON: `+err.Error()+`"}}`), nil
			}

			// Merge request updates into response
			for k, v := range requestBody {
				response[k] = v
			}
			response["id"] = ringId
			response["lastModifiedDateTime"] = "2024-01-01T00:00:00Z"

			// Preserve existing assignments if not in request
			if _, hasAssignments := requestBody["assignments"]; !hasAssignments {
				if assignments, hasExisting := ringData["assignments"]; hasExisting {
					response["assignments"] = assignments
				}
			}

			mockState.Lock()
			mockState.windowsUpdateRings[ringId] = response
			mockState.Unlock()

			return factories.SuccessResponse(200, response)(req)
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
				content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_delete", "get_windows_update_ring_not_found.json"))
				var errorResponse map[string]any
				json.Unmarshal([]byte(content), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register POST for assignments
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/[^/]+/assign$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			ringId := urlParts[len(urlParts)-2]

			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			mockState.Lock()
			if ringData, exists := mockState.windowsUpdateRings[ringId]; exists {
				if assignments, hasAssignments := requestBody["assignments"]; hasAssignments && assignments != nil {
					assignmentList := assignments.([]any)
					if len(assignmentList) > 0 {
						graphAssignments := []any{}
						for _, assignment := range assignmentList {
							if assignmentMap, ok := assignment.(map[string]any); ok {
								if target, hasTarget := assignmentMap["target"].(map[string]any); hasTarget {
									assignmentId := uuid.New().String()
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
						ringData["assignments"] = []any{}
					}
				} else {
					ringData["assignments"] = []any{}
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
				"value":          []map[string]any{},
			}
			return httpmock.NewJsonResponse(200, response)
		})
}

// determineCreateScenario determines which create JSON to load based on request body
func determineCreateScenario(requestBody map[string]any) string {
	// Check if this is an assignment test
	if assignments, hasAssignments := requestBody["assignments"]; hasAssignments {
		if assignmentList, ok := assignments.([]any); ok {
			if len(assignmentList) >= 5 {
				return "assignments_maximal.json"
			} else if len(assignmentList) > 0 {
				return "assignments_minimal.json"
			}
		}
	}

	// Determine scenario by automaticUpdateMode
	automaticUpdateMode, hasMode := requestBody["automaticUpdateMode"].(string)
	if !hasMode {
		return ""
	}

	switch automaticUpdateMode {
	case "notifyDownload":
		return "scenario_001_notify_download.json"
	case "autoInstallAtMaintenanceTime":
		return "scenario_002_auto_install_maintenance.json"
	case "autoInstallAndRebootAtMaintenanceTime":
		return "scenario_003_auto_reboot_maintenance.json"
	case "autoInstallAndRebootAtScheduledTime":
		return "scenario_004_scheduled_install.json"
	case "autoInstallAndRebootWithoutEndUserControl":
		return "scenario_005_no_end_user_control.json"
	case "windowsDefault":
		return "scenario_006_windows_default.json"
	default:
		return ""
	}
}

// determineReadScenario determines which read JSON to load based on stored ring data
func determineReadScenario(ringData map[string]any) string {
	// Check if this has assignments
	if assignments, hasAssignments := ringData["assignments"]; hasAssignments {
		if assignmentList, ok := assignments.([]any); ok && len(assignmentList) > 0 {
			if len(assignmentList) >= 5 {
				return "assignments_maximal.json"
			}
			return "assignments_minimal.json"
		}
	}

	// Determine scenario by automaticUpdateMode
	automaticUpdateMode, hasMode := ringData["automaticUpdateMode"].(string)
	if !hasMode {
		return ""
	}

	switch automaticUpdateMode {
	case "notifyDownload":
		return "scenario_001_notify_download.json"
	case "autoInstallAtMaintenanceTime":
		return "scenario_002_auto_install_maintenance.json"
	case "autoInstallAndRebootAtMaintenanceTime":
		return "scenario_003_auto_reboot_maintenance.json"
	case "autoInstallAndRebootAtScheduledTime":
		return "scenario_004_scheduled_install.json"
	case "autoInstallAndRebootWithoutEndUserControl":
		return "scenario_005_no_end_user_control.json"
	case "windowsDefault":
		return "scenario_006_windows_default.json"
	default:
		return ""
	}
}

// determineUpdateScenario determines which update JSON to load based on request and existing data
func determineUpdateScenario(requestBody map[string]any, existingData map[string]any) string {
	// Check if this is an assignment update
	if _, hasAssignments := requestBody["assignments"]; hasAssignments {
		// For updates, check the request assignments count
		if assignments, ok := requestBody["assignments"].([]any); ok {
			if len(assignments) >= 5 {
				return "assignments_maximal.json"
			} else if len(assignments) > 0 {
				return "assignments_minimal.json"
			}
		}
	}

	// Determine scenario by automaticUpdateMode (prefer request, fallback to existing)
	automaticUpdateMode := ""
	if mode, hasMode := requestBody["automaticUpdateMode"].(string); hasMode {
		automaticUpdateMode = mode
	} else if mode, hasMode := existingData["automaticUpdateMode"].(string); hasMode {
		automaticUpdateMode = mode
	}

	switch automaticUpdateMode {
	case "notifyDownload":
		return "scenario_001_notify_download.json"
	case "autoInstallAtMaintenanceTime":
		return "scenario_002_auto_install_maintenance.json"
	case "autoInstallAndRebootAtMaintenanceTime":
		return "scenario_003_auto_reboot_maintenance.json"
	case "autoInstallAndRebootAtScheduledTime":
		return "scenario_004_scheduled_install.json"
	case "autoInstallAndRebootWithoutEndUserControl":
		return "scenario_005_no_end_user_control.json"
	case "windowsDefault":
		return "scenario_006_windows_default.json"
	default:
		return ""
	}
}

// CleanupMockState clears the mock state for clean test runs
func (m *WindowsUpdateRingMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()

	for id := range mockState.windowsUpdateRings {
		delete(mockState.windowsUpdateRings, id)
	}
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *WindowsUpdateRingMock) RegisterErrorMocks() {
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		func(req *http.Request) (*http.Response, error) {
			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/deviceConfigurations",
				"value":          []map[string]any{},
			}
			return httpmock.NewJsonResponse(200, response)
		})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		func(req *http.Request) (*http.Response, error) {
			content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_create", "post_windows_update_ring_error.json"))
			var errorResponse map[string]any
			json.Unmarshal([]byte(content), &errorResponse)
			return httpmock.NewJsonResponse(400, errorResponse)
		})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/deviceConfigurations/([^/]+)$`,
		func(req *http.Request) (*http.Response, error) {
			content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_delete", "get_windows_update_ring_not_found.json"))
			var errorResponse map[string]any
			json.Unmarshal([]byte(content), &errorResponse)
			return httpmock.NewJsonResponse(404, errorResponse)
		})
}
