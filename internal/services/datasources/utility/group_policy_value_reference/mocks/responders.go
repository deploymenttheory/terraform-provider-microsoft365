package mocks

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/jarcoal/httpmock"
)

func init() {
	// Register with global registry
	mocks.GlobalRegistry.Register("group_policy_value_reference", &GroupPolicyValueReferenceMock{})
}

// GroupPolicyValueReferenceMock provides mock responses for Group Policy Value Reference operations
type GroupPolicyValueReferenceMock struct{}

// Ensure GroupPolicyValueReferenceMock implements MockRegistrar interface
var _ mocks.MockRegistrar = (*GroupPolicyValueReferenceMock)(nil)

// RegisterMocks registers HTTP mock responses for Group Policy Value Reference operations
func (m *GroupPolicyValueReferenceMock) RegisterMocks() {
	m.registerCommonMocks()
}

// RegisterErrorMocks registers mock responses that simulate error conditions
func (m *GroupPolicyValueReferenceMock) RegisterErrorMocks() {
	// For now, error mocks can use the same implementation
	m.registerCommonMocks()
}

// registerCommonMocks contains the actual mock registration logic
func (m *GroupPolicyValueReferenceMock) registerCommonMocks() {
	// Register GET for group policy definitions
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyDefinitions$`,
		func(req *http.Request) (*http.Response, error) {
			filter := req.URL.Query().Get("$filter")
			filterLower := strings.ToLower(filter)

			var jsonFile string
			// Determine which response to return based on filter (case-insensitive)
			if strings.Contains(filterLower, strings.ToLower("Allow users to connect remotely by using Remote Desktop Services")) {
				jsonFile = "../tests/responses/definitions/get_scenario_01_single_rdp_policy.json"
			} else if strings.Contains(filterLower, strings.ToLower("Show Home button on toolbar")) {
				jsonFile = "../tests/responses/definitions/get_scenario_02_multiple_home_button.json"
			} else if strings.Contains(filterLower, strings.ToLower("Show Home button")) {
				// Fuzzy search - returns broader results for "contains" filter
				jsonFile = "../tests/responses/definitions/get_scenario_04_fuzzy_search_home_button.json"
			} else if strings.Contains(filterLower, "nonexistent policy") || strings.Contains(filterLower, "policy that does not exist") {
				jsonFile = "../tests/responses/definitions/get_scenario_03_no_results.json"
			} else {
				// Default empty response
				jsonFile = "../tests/responses/definitions/get_scenario_03_no_results.json"
			}

			jsonStr, _ := helpers.ParseJSONFile(jsonFile)
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		})

	// Register GET for presentations of a definition
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyDefinitions/[^/]+/presentations$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			defID := urlParts[4]

			var jsonFile string
			// Map definition IDs to presentation JSON files
			switch defID {
			case "bb67ec37-f275-484c-942c-36a07e80add8": // RDP policy
				jsonFile = "../tests/responses/presentations/get_scenario_01_rdp_policy_presentations.json"
			case "eaca0db8-9673-4487-8055-d6dc037a3ef9", // Edge machine
				"5e4239e8-de87-4994-a5ca-f744e701d6b5", // Edge machine default
				"18afaa78-2f26-4e3e-a3d3-0f6191cc868a", // Chrome machine default
				"cae87e2e-618d-4cd2-b7d4-2805f4aa83c3", // Chrome user
				"ec97fcc1-3c4b-4d70-a46e-31fdc063f22c", // Edge user
				"2b088840-f531-4090-96d1-48c399825d90", // Chrome machine
				"8b4c9402-29de-4051-ad74-9e9462bfe619", // Edge user default
				"869bcf44-7180-4a96-9ee9-8f85e89ebf4d": // Chrome user default
				jsonFile = "../tests/responses/presentations/get_scenario_02_home_button_presentations.json"
			default:
				jsonFile = "../tests/responses/presentations/get_scenario_empty_presentations.json"
			}

			jsonStr, _ := helpers.ParseJSONFile(jsonFile)
			var responseObj map[string]any
			json.Unmarshal([]byte(jsonStr), &responseObj)
			return httpmock.NewJsonResponse(200, responseObj)
		})
}
