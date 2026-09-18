package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	softwareUpdateConfigurations map[string]map[string]any
}

func init() {
	// Initialize mockState
	mockState.softwareUpdateConfigurations = make(map[string]map[string]any)

	// Register a default 404 responder for any unmatched requests
	httpmock.RegisterNoResponder(func(_ *http.Request) (*http.Response, error) {
		return fixtureResponse(404, "validate_delete/get_configuration_not_found.json")
	})
}

// MacOSSoftwareUpdateConfigurationMock provides mock responses for macOS software update configuration operations
type MacOSSoftwareUpdateConfigurationMock struct{}

// RegisterMocks registers HTTP mock responses for macOS software update configuration operations
func (m *MacOSSoftwareUpdateConfigurationMock) RegisterMocks() {
	// Reset the state when registering mocks
	mockState.Lock()
	mockState.softwareUpdateConfigurations = make(map[string]map[string]any)
	mockState.Unlock()

	// Register GET for listing software update configurations
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		func(req *http.Request) (*http.Response, error) {
			mockState.Lock()
			configs := make([]map[string]any, 0, len(mockState.softwareUpdateConfigurations))
			for _, config := range mockState.softwareUpdateConfigurations {
				configs = append(configs, config)
			}
			mockState.Unlock()

			response, err := fixture("validate_read/get_configurations.json")
			if err != nil {
				return nil, err
			}
			response["value"] = configs

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for individual software update configuration
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			configId := urlParts[len(urlParts)-1]

			mockState.Lock()
			configData, exists := mockState.softwareUpdateConfigurations[configId]
			mockState.Unlock()

			if !exists {
				return fixtureResponse(404, "validate_delete/get_configuration_not_found.json")
			}

			scenario, err := configurationScenario(configData)
			if err != nil {
				return nil, err
			}
			responseCopy, err := fixture("validate_read/get_" + scenario + ".json")
			if err != nil {
				return nil, err
			}
			responseCopy["id"] = configId
			responseCopy["createdDateTime"] = configData["createdDateTime"]
			responseCopy["lastModifiedDateTime"] = configData["lastModifiedDateTime"]
			if strings.Contains(req.URL.Query().Get("$expand"), "assignments") {
				responseCopy["assignments"] = configData["assignments"]
			}

			return httpmock.NewJsonResponse(200, responseCopy)
		})

	// Register POST for creating software update configuration
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		func(req *http.Request) (*http.Response, error) {
			// Parse request body
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return fixtureResponse(400, "validate_create/post_invalid_request.json")
			}

			scenario, err := configurationScenario(requestBody)
			if err != nil {
				return nil, err
			}
			configData, err := fixture("validate_create/post_" + scenario + ".json")
			if err != nil {
				return nil, err
			}
			if err := validateConfigurationRequest(requestBody, configData); err != nil {
				return nil, err
			}
			configId := uuid.New().String()
			configData["id"] = configId
			mockState.Lock()
			mockState.softwareUpdateConfigurations[configId] = configData
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, configData)
		})

	// Register PATCH for updating software update configuration
	httpmock.RegisterResponder("PATCH", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			configId := urlParts[len(urlParts)-1]

			mockState.Lock()
			configData, exists := mockState.softwareUpdateConfigurations[configId]
			mockState.Unlock()

			if !exists {
				return fixtureResponse(404, "validate_delete/get_configuration_not_found.json")
			}

			// Parse request body
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return fixtureResponse(400, "validate_create/post_invalid_request.json")
			}

			scenario, err := configurationScenario(requestBody)
			if err != nil {
				return nil, err
			}
			updated, err := fixture("validate_update/patch_" + scenario + ".json")
			if err != nil {
				return nil, err
			}
			if err := validateConfigurationRequest(requestBody, updated); err != nil {
				return nil, err
			}
			updated["id"] = configId
			updated["createdDateTime"] = configData["createdDateTime"]
			updated["assignments"] = configData["assignments"]
			mockState.Lock()
			mockState.softwareUpdateConfigurations[configId] = updated
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, updated)
		})

	// Register DELETE for removing software update configuration
	httpmock.RegisterResponder("DELETE", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			configId := urlParts[len(urlParts)-1]

			mockState.Lock()
			_, exists := mockState.softwareUpdateConfigurations[configId]
			if exists {
				delete(mockState.softwareUpdateConfigurations, configId)
			}
			mockState.Unlock()

			if !exists {
				return fixtureResponse(404, "validate_delete/get_configuration_not_found.json")
			}

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register POST for assignments
	httpmock.RegisterResponder("POST", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/[^/]+/assign$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			configId := urlParts[len(urlParts)-2] // deviceConfigurations/{id}/assign

			// Parse request body to get assignments
			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return fixtureResponse(400, "validate_create/post_invalid_request.json")
			}

			mockState.Lock()
			configData, exists := mockState.softwareUpdateConfigurations[configId]
			mockState.Unlock()
			if !exists {
				return fixtureResponse(404, "validate_delete/get_configuration_not_found.json")
			}
			assignments, _ := requestBody["assignments"].([]any)
			file := "validate_read/get_assignments_empty.json"
			if len(assignments) > 0 {
				scenario, err := configurationScenario(configData)
				if err != nil {
					return nil, err
				}
				file = fmt.Sprintf("validate_read/get_%s_assignments_%d.json", scenario, len(assignments))
			}
			response, err := fixture(file)
			if err != nil {
				return nil, err
			}
			for _, assignment := range response["value"].([]any) {
				assignment.(map[string]any)["id"] = uuid.New().String()
			}
			mockState.Lock()
			configData["assignments"] = response["value"]
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register GET for assignments
	httpmock.RegisterResponder("GET", `=~^https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/[^/]+/assignments$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			configID := parts[len(parts)-2]
			response, err := fixture("validate_read/get_assignments_empty.json")
			if err != nil {
				return nil, err
			}
			mockState.Lock()
			if config, exists := mockState.softwareUpdateConfigurations[configID]; exists {
				response["value"] = config["assignments"]
			}
			mockState.Unlock()

			return httpmock.NewJsonResponse(200, response)
		})

}

// CleanupMockState resets the mock state
func (m *MacOSSoftwareUpdateConfigurationMock) CleanupMockState() {
	mockState.Lock()
	mockState.softwareUpdateConfigurations = make(map[string]map[string]any)
	mockState.Unlock()
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *MacOSSoftwareUpdateConfigurationMock) RegisterErrorMocks() {
	// Register GET for listing software update configurations (needed for uniqueness check)
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		func(req *http.Request) (*http.Response, error) {
			return fixtureResponse(200, "validate_read/get_configurations.json")
		})

	// Register error response for creating software update configuration with invalid data
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations",
		func(_ *http.Request) (*http.Response, error) {
			return fixtureResponse(400, "validate_create/post_error.json")
		})

	// Register error response for software update configuration not found
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/not-found-config",
		func(_ *http.Request) (*http.Response, error) {
			return fixtureResponse(404, "validate_delete/get_configuration_not_found.json")
		})
}

func fixture(name string) (map[string]any, error) {
	raw, err := helpers.ParseJSONFile("../tests/responses/" + name)
	if err != nil {
		return nil, err
	}
	var response map[string]any
	if err := json.Unmarshal([]byte(raw), &response); err != nil {
		return nil, fmt.Errorf("parse mock response %s: %w", name, err)
	}
	return response, nil
}

func fixtureResponse(status int, name string) (*http.Response, error) {
	response, err := fixture(name)
	if err != nil {
		return nil, err
	}
	return httpmock.NewJsonResponse(status, response)
}

func configurationScenario(config map[string]any) (string, error) {
	name, _ := config["displayName"].(string)
	switch {
	case strings.HasPrefix(name, "Test 01:"):
		return "001_minimal", nil
	case strings.HasPrefix(name, "Test 02:"):
		return "002_maximal", nil
	case strings.HasPrefix(name, "Test 03:"), strings.HasPrefix(name, "Test 04:"):
		scenario := "003_"
		if strings.HasPrefix(name, "Test 04:") {
			scenario = "004_"
		}
		description, _ := config["description"].(string)
		switch {
		case description == "":
			return scenario + "minimal", nil
		case strings.HasPrefix(description, "Intermediate"):
			return scenario + "intermediate", nil
		default:
			return scenario + "maximal", nil
		}
	case strings.HasPrefix(name, "Test 05:"):
		return "005_assignments_minimal", nil
	case strings.HasPrefix(name, "Test 06:"):
		return "006_assignments_maximal", nil
	case strings.HasPrefix(name, "Test 07:"):
		return "007_assignments_progression", nil
	case strings.HasPrefix(name, "Test 08:"):
		return "008_assignments_regression", nil
	}
	return "", fmt.Errorf("unknown software update configuration mock scenario %q", name)
}

func validateConfigurationRequest(body, response map[string]any) error {
	expected := make(map[string]any, len(response))
	for key, value := range response {
		switch key {
		case "@odata.context", "@odata.type", "id", "createdDateTime", "lastModifiedDateTime", "assignments":
			continue
		}
		expected[key] = value
	}
	actual := make(map[string]any, len(body))
	for key, value := range body {
		if key != "@odata.type" && value != nil {
			actual[key] = value
		}
	}
	if !reflect.DeepEqual(actual, expected) {
		return fmt.Errorf("software update request does not match JSON fixture for %v", response["displayName"])
	}
	return nil
}
