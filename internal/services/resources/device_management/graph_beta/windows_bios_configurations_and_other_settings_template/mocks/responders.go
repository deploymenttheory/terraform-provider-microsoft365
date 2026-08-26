package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const (
	collectionURL   = "https://graph.microsoft.com/beta/deviceManagement/hardwareConfigurations"
	itemURLPattern  = `=~^https://graph\.microsoft\.com/beta/deviceManagement/hardwareConfigurations/[^/]+$`
	assignURLPatter = `=~^https://graph\.microsoft\.com/beta/deviceManagement/hardwareConfigurations/[^/]+/assign$`

	notFoundFixture = "../tests/responses/validate_delete/get_windows_bios_configuration_not_found.json"
	errorFixture    = "../tests/responses/validate_create/post_009_error_scenario.json"
)

var mockState struct {
	sync.Mutex
	hardwareConfigurations map[string]map[string]any
}

func init() {
	mockState.hardwareConfigurations = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(
		httpmock.NewStringResponder(
			404,
			`{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`,
		),
	)
	mocks.GlobalRegistry.Register(
		"windows_bios_configurations_and_other_settings_template",
		&WindowsBiosConfigurationsAndOtherSettingsTemplateMock{},
	)
}

type WindowsBiosConfigurationsAndOtherSettingsTemplateMock struct{}

var _ mocks.MockRegistrar = (*WindowsBiosConfigurationsAndOtherSettingsTemplateMock)(nil)

// loadFixture reads a JSON fixture and unmarshals it into a map.
func loadFixture(path string) (map[string]any, error) {
	jsonStr, err := helpers.ParseJSONFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read fixture %s: %w", path, err)
	}

	var out map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &out); err != nil {
		return nil, fmt.Errorf("failed to parse fixture %s: %w", path, err)
	}
	return out, nil
}

func fixtureResponse(status int, path string) (*http.Response, error) {
	obj, err := loadFixture(path)
	if err != nil {
		return httpmock.NewStringResponse(
			500,
			`{"error":{"code":"InternalServerError","message":"Failed to load fixture `+path+`"}}`,
		), nil
	}
	resp, err := httpmock.NewJsonResponse(status, obj)
	if err != nil {
		return nil, fmt.Errorf("failed to build response from fixture %s: %w", path, err)
	}
	return resp, nil
}

func (m *WindowsBiosConfigurationsAndOtherSettingsTemplateMock) RegisterMocks() {
	mockState.Lock()
	mockState.hardwareConfigurations = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", collectionURL,
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			list := make([]map[string]any, 0, len(mockState.hardwareConfigurations))
			for _, v := range mockState.hardwareConfigurations {
				item := map[string]any{}
				for k, vv := range v {
					item[k] = vv
				}
				list = append(list, item)
			}
			mockState.Unlock()
			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/hardwareConfigurations",
				"value":          list,
			})
		})

	httpmock.RegisterResponder("GET", itemURLPattern,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			configData, exists := mockState.hardwareConfigurations[id]
			mockState.Unlock()

			if !exists {
				return fixtureResponse(404, notFoundFixture)
			}

			response, err := loadFixture(
				"../tests/responses/validate_read/" + determineReadScenario(configData),
			)
			if err != nil {
				return httpmock.NewStringResponse(
					500,
					`{"error":{"code":"InternalServerError","message":"Failed to load read scenario JSON: `+err.Error()+`"}}`,
				), nil
			}

			for k, v := range configData {
				response[k] = v
			}

			return httpmock.NewJsonResponse(200, response)
		})

	httpmock.RegisterResponder("POST", collectionURL,
		func(req *http.Request) (*http.Response, error) {
			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return fixtureResponse(400, errorFixture)
			}

			response, err := loadFixture(
				"../tests/responses/validate_create/" + determineCreateScenario(body),
			)
			if err != nil {
				return httpmock.NewStringResponse(
					500,
					`{"error":{"code":"InternalServerError","message":"Failed to load create scenario JSON: `+err.Error()+`"}}`,
				), nil
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
			mockState.hardwareConfigurations[id] = response
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, response)
		})

	httpmock.RegisterResponder("PATCH", itemURLPattern,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return fixtureResponse(400, errorFixture)
			}

			mockState.Lock()
			existing, ok := mockState.hardwareConfigurations[id]
			if !ok {
				mockState.Unlock()
				return fixtureResponse(404, notFoundFixture)
			}

			for k, v := range body {
				existing[k] = v
			}

			existing["lastModifiedDateTime"] = "2024-01-02T00:00:00Z"
			mockState.hardwareConfigurations[id] = existing

			scenarioFile := determineUpdateScenario(existing)
			mockState.Unlock()

			response, err := loadFixture("../tests/responses/validate_update/" + scenarioFile)
			if err != nil {
				return httpmock.NewStringResponse(
					500,
					`{"error":{"code":"InternalServerError","message":"Failed to load update scenario JSON: `+err.Error()+`"}}`,
				), nil
			}

			for k, v := range existing {
				response[k] = v
			}

			return httpmock.NewJsonResponse(200, response)
		})

	httpmock.RegisterResponder("POST", assignURLPatter,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-2]

			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return fixtureResponse(400, errorFixture)
			}

			mockState.Lock()
			if existing, ok := mockState.hardwareConfigurations[id]; ok {
				assignments, ok := body["hardwareConfigurationAssignments"].([]any)
				if !ok {
					assignments, _ = body["assignments"].([]any)
				}
				if assignments == nil {
					assignments = []any{}
				}
				existing["assignments"] = assignments
				mockState.hardwareConfigurations[id] = existing
			}
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	httpmock.RegisterResponder("DELETE", itemURLPattern,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]
			mockState.Lock()
			delete(mockState.hardwareConfigurations, id)
			mockState.Unlock()
			return httpmock.NewStringResponse(204, ""), nil
		})
}

func (m *WindowsBiosConfigurationsAndOtherSettingsTemplateMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.hardwareConfigurations = make(map[string]map[string]any)
	mockState.Unlock()

	httpmock.RegisterResponder("GET", collectionURL,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/hardwareConfigurations",
				"value":          []any{},
			})
		})

	httpmock.RegisterResponder("POST", collectionURL,
		func(req *http.Request) (*http.Response, error) {
			return fixtureResponse(400, errorFixture)
		})

	httpmock.RegisterResponder("GET", itemURLPattern,
		func(req *http.Request) (*http.Response, error) {
			return fixtureResponse(404, notFoundFixture)
		})
}

func (m *WindowsBiosConfigurationsAndOtherSettingsTemplateMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.hardwareConfigurations {
		delete(mockState.hardwareConfigurations, id)
	}
}

func determineCreateScenario(requestBody map[string]any) string {
	name := strings.ToLower(stringValue(requestBody, "displayName"))

	switch {
	case strings.Contains(name, "001"):
		return "post_001_scenario_minimal.json"
	case strings.Contains(name, "002"):
		return "post_002_scenario_maximal.json"
	case strings.Contains(name, "003"):
		return "post_003_lifecycle_step_1.json"
	case strings.Contains(name, "004"):
		return "post_004_lifecycle_step_1.json"
	case strings.Contains(name, "005"):
		return "post_005_assignments_minimal.json"
	case strings.Contains(name, "006"):
		return "post_006_assignments_maximal.json"
	case strings.Contains(name, "007"):
		return "post_007_assignments_lifecycle_step_1.json"
	case strings.Contains(name, "008"):
		return "post_008_assignments_lifecycle_step_1.json"
	case strings.Contains(name, "009"), strings.Contains(name, "error"):
		return "post_009_error_scenario.json"
	case strings.Contains(name, "maximal"):
		return "post_002_scenario_maximal.json"
	default:
		return "post_001_scenario_minimal.json"
	}
}

// determineReadScenario picks the read fixture. Lifecycle scenarios share a display name across
// both steps, so they are disambiguated by inspecting the stored state instead.
func determineReadScenario(configData map[string]any) string {
	name := strings.ToLower(stringValue(configData, "displayName"))

	switch {
	case strings.Contains(name, "003"):
		if stringValue(configData, "description") != "" {
			return "get_003_lifecycle_step_2.json"
		}
		return "get_003_lifecycle_step_1.json"
	case strings.Contains(name, "004"):
		if len(sliceValue(configData, "roleScopeTagIds")) == 1 {
			return "get_004_lifecycle_step_2.json"
		}
		return "get_004_lifecycle_step_1.json"
	case strings.Contains(name, "007"):
		if len(sliceValue(configData, "assignments")) > 1 {
			return "get_007_assignments_lifecycle_step_2.json"
		}
		return "get_007_assignments_lifecycle_step_1.json"
	case strings.Contains(name, "008"):
		if len(sliceValue(configData, "assignments")) == 1 {
			return "get_008_assignments_lifecycle_step_2.json"
		}
		return "get_008_assignments_lifecycle_step_1.json"
	case strings.Contains(name, "001"):
		return "get_001_scenario_minimal.json"
	case strings.Contains(name, "002"):
		return "get_002_scenario_maximal.json"
	case strings.Contains(name, "005"):
		return "get_005_assignments_minimal.json"
	case strings.Contains(name, "006"):
		return "get_006_assignments_maximal.json"
	case strings.Contains(name, "maximal"):
		return "get_002_scenario_maximal.json"
	default:
		return "get_001_scenario_minimal.json"
	}
}

func determineUpdateScenario(configData map[string]any) string {
	name := strings.ToLower(stringValue(configData, "displayName"))

	if strings.Contains(name, "004") {
		return "patch_004_lifecycle_step_2.json"
	}
	return "patch_003_lifecycle_step_2.json"
}

func stringValue(data map[string]any, key string) string {
	v, _ := data[key].(string)
	return v
}

// sliceValue normalises a JSON array field, which may arrive either as []any (decoded from a
// request body) or as a concrete slice written directly into the mock state.
func sliceValue(data map[string]any, key string) []any {
	switch v := data[key].(type) {
	case []any:
		return v
	case []string:
		out := make([]any, len(v))
		for i, s := range v {
			out[i] = s
		}
		return out
	default:
		return nil
	}
}
