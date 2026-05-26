package mocks

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	inventoryPolicies map[string]map[string]any
}

func init() {
	mockState.inventoryPolicies = make(map[string]map[string]any)
	mocks.GlobalRegistry.Register("settings_catalog_inventory_policy", &SettingsCatalogInventoryPolicyMock{})
}

type SettingsCatalogInventoryPolicyMock struct{}

var _ mocks.MockRegistrar = (*SettingsCatalogInventoryPolicyMock)(nil)

func (m *SettingsCatalogInventoryPolicyMock) RegisterMocks() {
	mockState.Lock()
	mockState.inventoryPolicies = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/inventoryPolicies",
		func(req *http.Request) (*http.Response, error) {
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			id := uuid.New().String()

			policy := map[string]any{
				"@odata.context":       "https://graph.microsoft.com/beta/$metadata#deviceManagement/inventoryPolicies/$entity",
				"id":                   id,
				"name":                 requestBody["name"],
				"platforms":            requestBody["platforms"],
				"technologies":         requestBody["technologies"],
				"settingCount":         0,
				"createdDateTime":      "2026-05-22T14:03:42.3345811Z",
				"lastModifiedDateTime": "2026-05-22T14:03:42.3345811Z",
				"creationSource":       nil,
			}

			if description, exists := requestBody["description"]; exists {
				policy["description"] = description
			} else {
				policy["description"] = ""
			}
			if roleScopeTagIds, exists := requestBody["roleScopeTagIds"]; exists {
				policy["roleScopeTagIds"] = roleScopeTagIds
			} else {
				policy["roleScopeTagIds"] = []string{"0"}
			}
			if settings, exists := requestBody["settings"]; exists {
				if settingsList, ok := settings.([]any); ok {
					processedSettings := make([]any, 0, len(settingsList))
					for i, setting := range settingsList {
						if settingMap, ok := setting.(map[string]any); ok {
							settingCopy := make(map[string]any)
							for k, v := range settingMap {
								settingCopy[k] = v
							}
							if _, hasId := settingCopy["id"]; !hasId {
								settingCopy["id"] = strconv.Itoa(i)
							}
							processedSettings = append(processedSettings, settingCopy)
						} else {
							processedSettings = append(processedSettings, setting)
						}
					}
					policy["settings"] = processedSettings
					policy["settingCount"] = len(processedSettings)
				} else {
					policy["settings"] = settings
				}
			}

			policy["assignments"] = []any{}

			mockState.Lock()
			mockState.inventoryPolicies[id] = policy
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, policy)
		})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/inventoryPolicies\('[^']+'\)$`,
		func(req *http.Request) (*http.Response, error) {
			id := extractODataId(req.URL.Path)

			mockState.Lock()
			policy, exists := mockState.inventoryPolicies[id]
			mockState.Unlock()

			if !exists {
				switch id {
				case "00000000-0000-0000-0000-000000000001":
					response, err := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_read", "get_inventory_policy_minimal.json"))
					if err != nil {
						return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load mock response"}}`), nil
					}
					return httpmock.NewJsonResponse(200, response)
				default:
					errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_inventory_policy_not_found.json"))
					return httpmock.NewJsonResponse(404, errorResponse)
				}
			}

			policyCopy := make(map[string]any)
			for k, v := range policy {
				policyCopy[k] = v
			}

			return httpmock.NewJsonResponse(200, policyCopy)
		})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/inventoryPolicies\('[^']+'\)/settings$`,
		func(req *http.Request) (*http.Response, error) {
			id := extractODataId(req.URL.Path)

			mockState.Lock()
			policyData, exists := mockState.inventoryPolicies[id]
			mockState.Unlock()

			settings := []any{}
			if exists {
				if storedSettings, hasSettings := policyData["settings"]; hasSettings {
					if settingsArray, ok := storedSettings.([]any); ok {
						settings = settingsArray
					}
				}
			}

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/inventoryPolicies('" + id + "')/settings",
				"value":          settings,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/inventoryPolicies\('[^']+'\)/assignments$`,
		func(req *http.Request) (*http.Response, error) {
			id := extractODataId(req.URL.Path)

			mockState.Lock()
			policyData, exists := mockState.inventoryPolicies[id]
			mockState.Unlock()

			assignments := []any{}
			if exists {
				if storedAssignments, hasAssignments := policyData["assignments"]; hasAssignments {
					if assignmentArray, ok := storedAssignments.([]any); ok {
						assignments = assignmentArray
					}
				}
			}

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/inventoryPolicies('" + id + "')/assignments",
				"value":          assignments,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/inventoryPolicies\('[^']+'\)/assign$`,
		func(req *http.Request) (*http.Response, error) {
			id := extractODataId(req.URL.Path)

			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			mockState.Lock()
			if policyData, exists := mockState.inventoryPolicies[id]; exists {
				if assignments, hasAssignments := requestBody["assignments"]; hasAssignments && assignments != nil {
					assignmentList := assignments.([]any)
					if len(assignmentList) > 0 {
						graphAssignments := []any{}
						for _, assignment := range assignmentList {
							if assignmentMap, ok := assignment.(map[string]any); ok {
								assignmentId := uuid.New().String()
								graphAssignment := map[string]any{
									"id":          assignmentId,
									"@odata.type": "#microsoft.graph.deviceManagementConfigurationPolicyAssignment",
								}
								if target, hasTarget := assignmentMap["target"].(map[string]any); hasTarget {
									graphAssignment["target"] = target
								}
								graphAssignments = append(graphAssignments, graphAssignment)
							}
						}
						policyData["assignments"] = graphAssignments
					} else {
						policyData["assignments"] = []any{}
					}
				} else {
					policyData["assignments"] = []any{}
				}
				mockState.inventoryPolicies[id] = policyData
			}
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	httpmock.RegisterResponder("PUT", `=~^https://graph\.microsoft\.com/beta/deviceManagement/inventoryPolicies\('[^']+'\)$`,
		func(req *http.Request) (*http.Response, error) {
			id := extractODataId(req.URL.Path)

			var requestBody map[string]any
			if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			mockState.Lock()
			defer mockState.Unlock()

			if existing, ok := mockState.inventoryPolicies[id]; ok {
				if settings, hasSettings := requestBody["settings"]; hasSettings {
					if settingsList, okList := settings.([]any); okList {
						processedSettings := make([]any, 0, len(settingsList))
						for i, setting := range settingsList {
							if settingMap, okMap := setting.(map[string]any); okMap {
								settingCopy := make(map[string]any)
								for k, v := range settingMap {
									settingCopy[k] = v
								}
								if _, hasId := settingCopy["id"]; !hasId {
									settingCopy["id"] = strconv.Itoa(i)
								}
								processedSettings = append(processedSettings, settingCopy)
							} else {
								processedSettings = append(processedSettings, setting)
							}
						}
						requestBody["settings"] = processedSettings
					}
				}

				for k, v := range requestBody {
					existing[k] = v
				}
				existing["lastModifiedDateTime"] = "2026-05-22T15:00:00.0000000Z"
				mockState.inventoryPolicies[id] = existing
			} else {
				errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_inventory_policy_not_found.json"))
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			return httpmock.NewStringResponse(204, ""), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/inventoryPolicies\('[^']+'\)$`,
		func(req *http.Request) (*http.Response, error) {
			id := extractODataId(req.URL.Path)

			mockState.Lock()
			_, exists := mockState.inventoryPolicies[id]
			if exists {
				delete(mockState.inventoryPolicies, id)
			}
			mockState.Unlock()

			if !exists {
				errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_inventory_policy_not_found.json"))
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			return httpmock.NewStringResponse(204, ""), nil
		})
}

func (m *SettingsCatalogInventoryPolicyMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.inventoryPolicies {
		delete(mockState.inventoryPolicies, id)
	}
}

func (m *SettingsCatalogInventoryPolicyMock) loadJSONResponse(filePath string) (map[string]any, error) {
	var response map[string]any
	content, err := os.ReadFile(filePath)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(content, &response)
	return response, err
}

func (m *SettingsCatalogInventoryPolicyMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.inventoryPolicies = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/inventoryPolicies",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/inventoryPolicies\('[^']+'\)$`,
		func(req *http.Request) (*http.Response, error) {
			errorResponse, _ := m.loadJSONResponse(filepath.Join("tests", "responses", "validate_delete", "get_inventory_policy_not_found.json"))
			return httpmock.NewJsonResponse(404, errorResponse)
		})
}

func extractODataId(path string) string {
	startMarker := "inventoryPolicies('"
	startIdx := strings.Index(path, startMarker)
	if startIdx == -1 {
		return ""
	}
	startIdx += len(startMarker)
	endIdx := strings.Index(path[startIdx:], "')")
	if endIdx == -1 {
		return ""
	}
	return path[startIdx : startIdx+endIdx]
}
