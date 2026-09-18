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

// DepOnboardingSettingsTestID is the fixed Apple ADE/ABM DEP token ID returned by the mocked
// /deviceManagement/depOnboardingSettings endpoint, used to exercise dep_onboarding_settings_id
// auto-resolution (resolve_id.go) when a test config omits it.
const DepOnboardingSettingsTestID = "30000000-0000-0000-0000-000000000003"

var mockState struct {
	sync.Mutex
	policies          map[string]map[string]any // policy ID -> fixture response
	membershipTargets map[string]string         // policy ID -> current device_security_group target, if set
	defaultIOSProfile map[string]string         // dep_onboarding_settings_id -> default policy ID, if set
}

func init() {
	resetMockState()
	httpmock.RegisterNoResponder(func(_ *http.Request) (*http.Response, error) {
		return fixtureResponse(404, "validate_delete/get_resource_not_found.json")
	})
}

func resetMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.policies = make(map[string]map[string]any)
	mockState.membershipTargets = make(map[string]string)
	mockState.defaultIOSProfile = make(map[string]string)
}

// IOSiPadOSDeviceEnrollmentPolicyMock provides mock responses for iOS/iPadOS ADE enrollment policy operations.
type IOSiPadOSDeviceEnrollmentPolicyMock struct{}

// RegisterMocks registers HTTP mock responders for iOS/iPadOS ADE enrollment policy operations.
//
//nolint:gocyclo // Mock registration functions naturally have high complexity
func (m *IOSiPadOSDeviceEnrollmentPolicyMock) RegisterMocks() {
	resetMockState()

	// Apple DEP token lookup - used by resolveDepOnboardingSettingsId when
	// dep_onboarding_settings_id is omitted from config.
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/depOnboardingSettings",
		func(req *http.Request) (*http.Response, error) {
			return fixtureResponse(200, "validate_read/get_dep_onboarding_settings.json")
		})

	// Read a single DEP token, expanding its default iOS enrollment profile - GET
	// /deviceManagement/depOnboardingSettings/{id}?$expand=defaultIosEnrollmentProfile. Used by
	// resolveIsDefaultPolicyAssignment on every Read and by validateRequest on demote attempts.
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/depOnboardingSettings/[0-9a-fA-F-]+`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			defaultPolicyID, hasDefault := mockState.defaultIOSProfile[id]
			mockState.Unlock()

			file := "validate_read/get_dep_onboarding_setting.json"
			if hasDefault {
				file = "validate_read/get_default_enrollment_profile.json"
			}
			response, err := fixture(file)
			if err != nil {
				return nil, err
			}
			response["id"] = id
			if hasDefault {
				response["defaultIosEnrollmentProfile"].(map[string]any)["id"] = id + "_" + defaultPolicyID
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Set default iOS enrollment profile - POST
	// /deviceManagement/depOnboardingSettings/{depId}/enrollmentProfiles/{enrollmentProfileId}/setDefaultProfile.
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/depOnboardingSettings/[0-9a-fA-F-]+/enrollmentProfiles/.+/setDefaultProfile$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			// [...  "depOnboardingSettings", "{depId}", "enrollmentProfiles", "{depId}_{policyId}", "setDefaultProfile"]
			depId := parts[len(parts)-4]
			enrollmentProfileId := parts[len(parts)-2]
			policyId := strings.TrimPrefix(enrollmentProfileId, depId+"_")

			mockState.Lock()
			mockState.defaultIOSProfile[depId] = policyId
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Security group owners - used by validateSecurityGroupOwnership. Every group is mocked as
	// owned by the Intune Provisioning Client.
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/[0-9a-fA-F-]+/owners$`,
		func(req *http.Request) (*http.Response, error) {
			return fixtureResponse(200, "validate_read/get_group_owners.json")
		})

	// Create policy - POST /deviceManagement/configurationPolicies
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		func(req *http.Request) (*http.Response, error) {
			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return fixtureResponse(400, "validate_create/post_invalid_request.json")
			}

			scenario, err := policyScenario(body)
			if err != nil {
				return nil, err
			}
			response, err := fixture("validate_create/post_" + scenario + ".json")
			if err != nil {
				return nil, err
			}
			if err := validatePolicyRequest(body, response, scenario); err != nil {
				return nil, err
			}
			id := uuid.New().String()
			response["id"] = id
			mockState.Lock()
			mockState.policies[id] = response
			mockState.Unlock()

			return httpmock.NewJsonResponse(201, response)
		})

	// Read policy base fields - GET /deviceManagement/configurationPolicies/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			policy, exists := mockState.policies[id]
			mockState.Unlock()

			if !exists {
				// Configuration policies return 400 (not 404) for a missing resource.
				return fixtureResponse(400, "validate_delete/get_policy_not_found.json")
			}

			scenario, err := policyScenario(policy)
			if err != nil {
				return nil, err
			}
			response, err := fixture("validate_read/get_" + scenario + ".json")
			if err != nil {
				return nil, err
			}
			response["id"] = id
			response["createdDateTime"] = policy["createdDateTime"]
			response["lastModifiedDateTime"] = policy["lastModifiedDateTime"]
			return httpmock.NewJsonResponse(200, response)
		})

	// Read policy settings - GET /deviceManagement/configurationPolicies/{id}/settings
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F-]+/settings$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-2]

			mockState.Lock()
			policy, exists := mockState.policies[id]
			mockState.Unlock()
			file := "validate_read/get_settings_empty.json"
			if exists {
				scenario, err := policyScenario(policy)
				if err != nil {
					return nil, err
				}
				switch scenario {
				case "002_maximal", "004_lifecycle_maximal", "005_lifecycle_maximal":
					file = "validate_read/get_" + scenario + "_settings.json"
				default:
					file = "validate_read/get_001_minimal_settings.json"
				}
			}
			response, err := fixture(file)
			if err != nil {
				return nil, err
			}
			response["@odata.context"] = "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies('" + id + "')/settings"
			return httpmock.NewJsonResponse(200, response)
		})

	// Update policy - PUT /deviceManagement/configurationPolicies('{id}') (raw request: the Graph
	// API does not allow PATCH on the 'settings' navigation property).
	httpmock.RegisterResponder("PUT", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies\('[0-9a-fA-F-]+'\)$`,
		func(req *http.Request) (*http.Response, error) {
			id := extractParenID(req.URL.Path)

			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return fixtureResponse(400, "validate_create/post_invalid_request.json")
			}

			mockState.Lock()
			existing, exists := mockState.policies[id]
			if !exists {
				mockState.Unlock()
				return fixtureResponse(400, "validate_delete/get_policy_not_found.json")
			}

			mockState.Unlock()
			scenario, err := policyScenario(body)
			if err != nil {
				return nil, err
			}
			updated, err := fixture("validate_update/put_" + scenario + ".json")
			if err != nil {
				return nil, err
			}
			if err := validatePolicyRequest(body, updated, scenario); err != nil {
				return nil, err
			}
			updated["id"] = id
			updated["createdDateTime"] = existing["createdDateTime"]
			mockState.Lock()
			mockState.policies[id] = updated
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Set enrollment time device membership target - POST
	// /deviceManagement/configurationPolicies('{id}')/setEnrollmentTimeDeviceMembershipTarget
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies\('[0-9a-fA-F-]+'\)/setEnrollmentTimeDeviceMembershipTarget$`,
		func(req *http.Request) (*http.Response, error) {
			id := extractParenID(req.URL.Path)

			var body struct {
				EnrollmentTimeDeviceMembershipTargets []struct {
					TargetId string `json:"targetId"`
				} `json:"enrollmentTimeDeviceMembershipTargets"`
			}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return fixtureResponse(400, "validate_create/post_invalid_request.json")
			}

			mockState.Lock()
			if len(body.EnrollmentTimeDeviceMembershipTargets) > 0 {
				mockState.membershipTargets[id] = body.EnrollmentTimeDeviceMembershipTargets[0].TargetId
			}
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Clear enrollment time device membership target - DELETE
	// /deviceManagement/configurationPolicies('{id}')/clearEnrollmentTimeDeviceMembershipTarget
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies\('[0-9a-fA-F-]+'\)/clearEnrollmentTimeDeviceMembershipTarget$`,
		func(req *http.Request) (*http.Response, error) {
			id := extractParenID(req.URL.Path)

			mockState.Lock()
			delete(mockState.membershipTargets, id)
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Retrieve enrollment time device membership target - GET
	// /deviceManagement/configurationPolicies('{id}')/retrieveEnrollmentTimeDeviceMembershipTarget.
	// The "beta/?" tolerates either a well-formed or malformed (missing slash) baseurl join, since
	// this raw request is built without going through the SDK's own URL templating.
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/?deviceManagement/configurationPolicies\('[0-9a-fA-F-]+'\)/retrieveEnrollmentTimeDeviceMembershipTarget$`,
		func(req *http.Request) (*http.Response, error) {
			id := extractParenID(req.URL.Path)

			mockState.Lock()
			targetId, hasTarget := mockState.membershipTargets[id]
			mockState.Unlock()

			file := "validate_read/get_membership_target_empty.json"
			if hasTarget {
				file = "validate_read/get_membership_target.json"
			}
			response, err := fixture(file)
			if err != nil {
				return nil, err
			}
			if hasTarget {
				statuses := response["enrollmentTimeDeviceMembershipTargetValidationStatuses"].([]any)
				statuses[0].(map[string]any)["targetId"] = targetId
			}
			return httpmock.NewJsonResponse(200, response)
		})

	// Delete policy - DELETE /deviceManagement/configurationPolicies/{id}
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			delete(mockState.policies, id)
			delete(mockState.membershipTargets, id)
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})
}

// extractParenID pulls the GUID out of a "...configurationPolicies('<id>')..." path segment.
func extractParenID(urlPath string) string {
	start := strings.Index(urlPath, "('")
	end := strings.Index(urlPath, "')")
	if start == -1 || end == -1 || end <= start+2 {
		return ""
	}
	return urlPath[start+2 : end]
}

// CleanupMockState resets the in-memory mock state after a test.
func (m *IOSiPadOSDeviceEnrollmentPolicyMock) CleanupMockState() {
	resetMockState()
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

func policyScenario(policy map[string]any) (string, error) {
	name, _ := policy["name"].(string)
	switch strings.TrimPrefix(name, "unit-test-ios-ade-") {
	case "minimal":
		return "001_minimal", nil
	case "maximal":
		return "002_maximal", nil
	case "update":
		return "003_lifecycle_minimal", nil
	case "update-updated":
		return "004_lifecycle_maximal", nil
	case "downgrade":
		return "005_lifecycle_maximal", nil
	case "downgrade-minimal":
		return "006_lifecycle_minimal", nil
	case "device-group":
		return "007_device_security_group", nil
	case "device-group-update":
		return "008_device_security_group_update", nil
	case "default-assignment":
		return "011_default_assignment", nil
	case "default-switch-a":
		return "012_default_switch_a", nil
	case "default-switch-b":
		return "012_default_switch_b", nil
	}
	return "", fmt.Errorf("unknown enrollment policy mock scenario %q", name)
}

// validatePolicyRequest checks the settings sent by the provider against the response fixture.
func validatePolicyRequest(body, policy map[string]any, scenario string) error {
	settingsScenario := "001_minimal"
	switch scenario {
	case "002_maximal", "004_lifecycle_maximal", "005_lifecycle_maximal":
		settingsScenario = scenario
	}
	settings, err := fixture("validate_read/get_" + settingsScenario + "_settings.json")
	if err != nil {
		return err
	}
	expected := make(map[string]any, len(policy))
	for key, value := range policy {
		switch key {
		case "@odata.context", "id", "createdDateTime", "lastModifiedDateTime", "settingCount", "isAssigned":
			continue
		}
		expected[key] = value
	}
	// GET adds setting IDs; POST and PUT submit only the setting instances.
	for _, setting := range settings["value"].([]any) {
		delete(setting.(map[string]any), "id")
	}
	expected["settings"] = settings["value"]
	if !reflect.DeepEqual(body, expected) {
		return fmt.Errorf("policy request does not match JSON fixture for %s", scenario)
	}
	return nil
}
