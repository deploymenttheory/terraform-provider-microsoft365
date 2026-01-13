package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	featureProfiles map[string]map[string]any
}

func init() {
	mockState.featureProfiles = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_feature_update_profile", &WindowsFeatureUpdateProfileMock{})
}

type WindowsFeatureUpdateProfileMock struct{}

var _ mocks.MockRegistrar = (*WindowsFeatureUpdateProfileMock)(nil)

func (m *WindowsFeatureUpdateProfileMock) RegisterMocks() {
	mockState.Lock()
	mockState.featureProfiles = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/windowsFeatureUpdateProfiles",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			list := make([]map[string]any, 0, len(mockState.featureProfiles))
			for _, v := range mockState.featureProfiles {
				copy := map[string]any{}
				for k, vv := range v {
					copy[k] = vv
				}
				list = append(list, copy)
			}
			mockState.Unlock()
			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/windowsFeatureUpdateProfiles",
				"value":          list,
			})
		})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsFeatureUpdateProfiles/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			profileData, exists := mockState.featureProfiles[id]
			mockState.Unlock()

			if !exists {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_feature_update_profile_not_found.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			scenarioFile := determineReadScenario(profileData)
			if scenarioFile == "" {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Could not determine scenario for GET request"}}`), nil
			}

			jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_read/" + scenarioFile)
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load read scenario JSON: `+err.Error()+`"}}`), nil
			}

			var response map[string]any
			if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse read scenario JSON: `+err.Error()+`"}}`), nil
			}

			for k, v := range profileData {
				response[k] = v
			}

			return httpmock.NewJsonResponse(200, response)
		})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/windowsFeatureUpdateProfiles",
		func(req *http.Request) (*http.Response, error) {
			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_009_error_scenario.json")
				var errObj map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errObj)
				return httpmock.NewJsonResponse(400, errObj)
			}

			scenarioFile := determineCreateScenario(body)
			if scenarioFile == "" {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Could not determine scenario for POST request"}}`), nil
			}

			jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_create/" + scenarioFile)
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load create scenario JSON: `+err.Error()+`"}}`), nil
			}

			var response map[string]any
			if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse create scenario JSON: `+err.Error()+`"}}`), nil
			}

			id := uuid.New().String()
			response["id"] = id

			for k, v := range body {
				response[k] = v
			}

			if _, hasRoleScopeTags := body["roleScopeTagIds"]; !hasRoleScopeTags {
				response["roleScopeTagIds"] = []string{"0"}
			}

			response["assignments"] = []any{}

			mockState.Lock()
			mockState.featureProfiles[id] = response
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, response)
		})

	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsFeatureUpdateProfiles/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]
			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_009_error_scenario.json")
				var errObj map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errObj)
				return httpmock.NewJsonResponse(400, errObj)
			}

			mockState.Lock()
			existing, ok := mockState.featureProfiles[id]
			if !ok {
				mockState.Unlock()
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_feature_update_profile_not_found.json")
				var errObj map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errObj)
				return httpmock.NewJsonResponse(404, errObj)
			}

			for k, v := range body {
				existing[k] = v
			}

			existing["lastModifiedDateTime"] = "2024-01-02T00:00:00Z"

			scenarioFile := determineUpdateScenario(existing)
			if scenarioFile != "" {
				if jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_update/" + scenarioFile); err == nil {
					var templateObj map[string]any
					if err := json.Unmarshal([]byte(jsonStr), &templateObj); err == nil {
						for k, v := range existing {
							templateObj[k] = v
						}
						existing = templateObj
					}
				}
			}

			mockState.featureProfiles[id] = existing
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, existing)
		})

	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsFeatureUpdateProfiles/[^/]+/assign$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-2]
			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			mockState.Lock()
			if existing, ok := mockState.featureProfiles[id]; ok {
				assignments, ok := body["windowsFeatureUpdateProfileAssignments"].([]any)
				if !ok {
					assignments, _ = body["assignments"].([]any)
				}
				if assignments == nil {
					assignments = []any{}
				}
				existing["assignments"] = assignments
				mockState.featureProfiles[id] = existing
			}
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsFeatureUpdateProfiles/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]
			mockState.Lock()
			delete(mockState.featureProfiles, id)
			mockState.Unlock()
			return httpmock.NewStringResponse(204, ""), nil
		})
}

func (m *WindowsFeatureUpdateProfileMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.featureProfiles = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/windowsFeatureUpdateProfiles",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/windowsFeatureUpdateProfiles",
				"value":          []any{},
			})
		})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/windowsFeatureUpdateProfiles",
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_009_error_scenario.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(400, errObj)
		})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsFeatureUpdateProfiles/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_feature_update_profile_not_found.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		})
}

func (m *WindowsFeatureUpdateProfileMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.featureProfiles {
		delete(mockState.featureProfiles, id)
	}
}

func determineCreateScenario(requestBody map[string]any) string {
	displayName, hasName := requestBody["displayName"].(string)
	if !hasName {
		return "post_001_scenario_minimal.json"
	}

	name := strings.ToLower(displayName)

	if strings.Contains(name, "001") {
		return "post_001_scenario_minimal.json"
	}
	if strings.Contains(name, "002") {
		return "post_002_scenario_maximal.json"
	}
	if strings.Contains(name, "003") {
		return "post_003_lifecycle_step_1.json"
	}
	if strings.Contains(name, "004") {
		return "post_004_lifecycle_step_1.json"
	}
	if strings.Contains(name, "005") {
		return "post_005_assignments_minimal.json"
	}
	if strings.Contains(name, "006") {
		return "post_006_assignments_maximal.json"
	}
	if strings.Contains(name, "007") {
		return "post_007_assignments_lifecycle_step_1.json"
	}
	if strings.Contains(name, "008") {
		return "post_008_assignments_lifecycle_step_1.json"
	}
	if strings.Contains(name, "009") || strings.Contains(name, "error") {
		return "post_009_error_scenario.json"
	}

	if strings.Contains(name, "minimal") {
		return "post_001_scenario_minimal.json"
	}
	if strings.Contains(name, "maximal") {
		return "post_002_scenario_maximal.json"
	}

	return "post_001_scenario_minimal.json"
}

func determineReadScenario(profileData map[string]any) string {
	displayName, hasName := profileData["displayName"].(string)
	if !hasName {
		return "get_001_scenario_minimal.json"
	}

	name := strings.ToLower(displayName)

	if strings.Contains(name, "003") {
		if desc, hasDesc := profileData["description"]; hasDesc && desc != nil && desc != "" {
			return "get_003_lifecycle_step_2.json"
		}
		return "get_003_lifecycle_step_1.json"
	}

	if strings.Contains(name, "004") {
		if desc, hasDesc := profileData["description"]; !hasDesc || desc == nil || desc == "" {
			if tags, hasTags := profileData["roleScopeTagIds"].([]any); hasTags && len(tags) == 1 {
				return "get_004_lifecycle_step_2.json"
			}
		}
		return "get_004_lifecycle_step_1.json"
	}

	if strings.Contains(name, "007") {
		if assignments, hasAssignments := profileData["assignments"].([]any); hasAssignments {
			if len(assignments) > 1 {
				return "get_007_assignments_lifecycle_step_2.json"
			}
		}
		return "get_007_assignments_lifecycle_step_1.json"
	}

	if strings.Contains(name, "008") {
		if assignments, hasAssignments := profileData["assignments"].([]any); hasAssignments {
			if len(assignments) == 1 {
				return "get_008_assignments_lifecycle_step_2.json"
			}
		}
		return "get_008_assignments_lifecycle_step_1.json"
	}

	if strings.Contains(name, "001") {
		return "get_001_scenario_minimal.json"
	}
	if strings.Contains(name, "002") {
		return "get_002_scenario_maximal.json"
	}
	if strings.Contains(name, "005") {
		return "get_005_assignments_minimal.json"
	}
	if strings.Contains(name, "006") {
		return "get_006_assignments_maximal.json"
	}

	if strings.Contains(name, "minimal") {
		return "get_001_scenario_minimal.json"
	}
	if strings.Contains(name, "maximal") {
		return "get_002_scenario_maximal.json"
	}

	return "get_001_scenario_minimal.json"
}

func determineUpdateScenario(profileData map[string]any) string {
	displayName, hasName := profileData["displayName"].(string)
	if !hasName {
		return "patch_003_lifecycle_step_2.json"
	}

	name := strings.ToLower(displayName)

	if strings.Contains(name, "003") {
		return "patch_003_lifecycle_step_2.json"
	}
	if strings.Contains(name, "004") {
		return "patch_004_lifecycle_step_2.json"
	}

	return "patch_003_lifecycle_step_2.json"
}
