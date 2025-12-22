package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks/factories"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	expeditePolicies map[string]map[string]any
}

func init() {
	mockState.expeditePolicies = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("windows_quality_update_expedite_policy", &WindowsQualityUpdateExpeditePolicyMock{})
}

type WindowsQualityUpdateExpeditePolicyMock struct{}

var _ mocks.MockRegistrar = (*WindowsQualityUpdateExpeditePolicyMock)(nil)

func (m *WindowsQualityUpdateExpeditePolicyMock) RegisterMocks() {
	mockState.Lock()
	mockState.expeditePolicies = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/windowsQualityUpdateProfiles",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			list := make([]map[string]any, 0, len(mockState.expeditePolicies))
			for _, v := range mockState.expeditePolicies {
				copy := map[string]any{}
				for k, vv := range v {
					copy[k] = vv
				}
				list = append(list, copy)
			}
			mockState.Unlock()
			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/windowsQualityUpdateProfiles",
				"value":          list,
			})
		})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsQualityUpdateProfiles/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			policyData, exists := mockState.expeditePolicies[id]
			mockState.Unlock()

			if !exists {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_quality_update_expedite_policy_not_found.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			scenarioFile := determineReadScenario(policyData)
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

			for k, v := range policyData {
				response[k] = v
			}

			if _, has := policyData["expeditedUpdateSettings"]; !has {
				delete(response, "expeditedUpdateSettings")
			}

			return httpmock.NewJsonResponse(200, response)
		})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/windowsQualityUpdateProfiles",
		func(req *http.Request) (*http.Response, error) {
			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_quality_update_expedite_policy_error.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(400, errorResponse)
			}

			id := uuid.New().String()

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

			response["id"] = id
			if displayName, hasName := body["displayName"]; hasName {
				response["displayName"] = displayName
			}
			if description, hasDesc := body["description"]; hasDesc {
				response["description"] = description
			}
			if expeditedUpdateSettings, has := body["expeditedUpdateSettings"]; has {
				response["expeditedUpdateSettings"] = expeditedUpdateSettings
			} else {
				delete(response, "expeditedUpdateSettings")
			}
			if roleScopeTagIds, has := body["roleScopeTagIds"]; has {
				response["roleScopeTagIds"] = roleScopeTagIds
			}

			mockState.Lock()
			mockState.expeditePolicies[id] = response
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, response)
		})

	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsQualityUpdateProfiles/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			policyData, exists := mockState.expeditePolicies[id]
			mockState.Unlock()

			if !exists {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_quality_update_expedite_policy_not_found.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_quality_update_expedite_policy_error.json")
				var errorResponse map[string]any
				_ = json.Unmarshal([]byte(jsonStr), &errorResponse)
				return httpmock.NewJsonResponse(400, errorResponse)
			}

			for k, v := range requestBody {
				policyData[k] = v
			}

			scenarioFile := determineUpdateScenario(requestBody, policyData)
			if scenarioFile == "" {
				if _, has := policyData["expeditedUpdateSettings"]; !has {
					delete(policyData, "expeditedUpdateSettings")
				}
				mockState.Lock()
				mockState.expeditePolicies[id] = policyData
				mockState.Unlock()
				return httpmock.NewJsonResponse(200, policyData)
			}

			jsonStr, err := helpers.ParseJSONFile("../tests/responses/validate_update/" + scenarioFile)
			if err != nil {
				mockState.Lock()
				mockState.expeditePolicies[id] = policyData
				mockState.Unlock()
				return httpmock.NewJsonResponse(200, policyData)
			}

			var response map[string]any
			if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse update scenario JSON: `+err.Error()+`"}}`), nil
			}

			for k, v := range requestBody {
				response[k] = v
			}

			response["id"] = policyData["id"]

			if _, has := requestBody["expeditedUpdateSettings"]; !has {
				delete(response, "expeditedUpdateSettings")
			}

			mockState.Lock()
			mockState.expeditePolicies[id] = response
			mockState.Unlock()

			return factories.SuccessResponse(200, response)(req)
		})

	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsQualityUpdateProfiles/[^/]+/assign$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-2]
			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}
			mockState.Lock()
			if existing, ok := mockState.expeditePolicies[id]; ok {
				assignments, _ := body["assignments"].([]any)
				if assignments == nil {
					assignments = []any{}
				}
				existing["assignments"] = assignments
				mockState.expeditePolicies[id] = existing
			}
			mockState.Unlock()
			return httpmock.NewStringResponse(204, ""), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsQualityUpdateProfiles/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]
			mockState.Lock()
			delete(mockState.expeditePolicies, id)
			mockState.Unlock()
			return httpmock.NewStringResponse(204, ""), nil
		})
}

func (m *WindowsQualityUpdateExpeditePolicyMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.expeditePolicies = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/windowsQualityUpdateProfiles",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/windowsQualityUpdateProfiles",
				"value":          []any{},
			})
		})

	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/windowsQualityUpdateProfiles",
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_windows_quality_update_expedite_policy_error.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(400, errObj)
		})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/windowsQualityUpdateProfiles/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_delete/get_windows_quality_update_expedite_policy_not_found.json")
			var errObj map[string]any
			_ = json.Unmarshal([]byte(jsonStr), &errObj)
			return httpmock.NewJsonResponse(404, errObj)
		})
}

func (m *WindowsQualityUpdateExpeditePolicyMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.expeditePolicies {
		delete(mockState.expeditePolicies, id)
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

	if strings.Contains(name, "minimal") {
		return "post_001_scenario_minimal.json"
	}
	if strings.Contains(name, "maximal") {
		return "post_002_scenario_maximal.json"
	}

	return "post_001_scenario_minimal.json"
}

func determineReadScenario(policyData map[string]any) string {
	displayName, hasName := policyData["displayName"].(string)
	if !hasName {
		return "get_001_scenario_minimal.json"
	}

	name := strings.ToLower(displayName)

	if strings.Contains(name, "003") {
		if desc, hasDesc := policyData["description"]; hasDesc && desc != nil && desc != "" {
			return "get_003_lifecycle_step_2.json"
		}
		if expeditedSettings, hasSettings := policyData["expeditedUpdateSettings"]; hasSettings && expeditedSettings != nil {
			return "get_003_lifecycle_step_2.json"
		}
		return "get_003_lifecycle_step_1.json"
	}

	if strings.Contains(name, "004") {
		if desc, hasDesc := policyData["description"]; !hasDesc || desc == nil || desc == "" {
			if tags, hasTags := policyData["roleScopeTagIds"].([]any); hasTags && len(tags) == 1 {
				return "get_004_lifecycle_step_2.json"
			}
		}
		return "get_004_lifecycle_step_1.json"
	}

	if strings.Contains(name, "007") {
		if assignments, hasAssignments := policyData["assignments"].([]any); hasAssignments {
			if len(assignments) > 1 {
				return "get_007_assignments_lifecycle_step_2.json"
			}
		}
		return "get_007_assignments_lifecycle_step_1.json"
	}

	if strings.Contains(name, "008") {
		if assignments, hasAssignments := policyData["assignments"].([]any); hasAssignments {
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

	return "get_001_scenario_minimal.json"
}

func determineUpdateScenario(requestBody map[string]any, policyData map[string]any) string {
	displayName, hasName := requestBody["displayName"].(string)
	if !hasName {
		displayName, hasName = policyData["displayName"].(string)
	}

	if !hasName {
		return ""
	}

	name := strings.ToLower(displayName)

	if strings.Contains(name, "003") {
		if _, hasDesc := requestBody["description"]; hasDesc {
			return "patch_003_lifecycle_step_2.json"
		}
		if expeditedSettings, hasSettings := requestBody["expeditedUpdateSettings"]; hasSettings && expeditedSettings != nil {
			return "patch_003_lifecycle_step_2.json"
		}
	}

	if strings.Contains(name, "004") {
		if desc, hasDesc := requestBody["description"]; hasDesc && (desc == nil || desc == "") {
			return "patch_004_lifecycle_step_2.json"
		}
		if tags, hasTags := requestBody["roleScopeTagIds"].([]any); hasTags && len(tags) == 1 {
			return "patch_004_lifecycle_step_2.json"
		}
	}

	if strings.Contains(name, "007") {
		if assignments, hasAssignments := policyData["assignments"].([]any); hasAssignments && len(assignments) > 1 {
			return "patch_007_assignments_lifecycle_step_2.json"
		}
	}

	if strings.Contains(name, "008") {
		if assignments, hasAssignments := policyData["assignments"].([]any); hasAssignments && len(assignments) == 1 {
			return "patch_008_assignments_lifecycle_step_2.json"
		}
	}

	return ""
}
