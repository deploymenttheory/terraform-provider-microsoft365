// Package mocks provides HTTP mock responders for the visionos_device_enrollment_policy unit tests.
//
// Unlike a fixed-fixture mock (one canned JSON file per scenario), this mock echoes back whatever
// top-level fields and settings tree were sent on POST/PUT, keyed by the generated policy ID. This
// keeps the mock in lock-step with construct.go/state.go without hand-maintaining a JSON fixture
// per settings combination, since the settings catalog tree here has ~20 independent toggles.
package mocks

import (
	"encoding/json"
	"maps"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// DepOnboardingSettingsTestID is the fixed Apple ADE/ABM DEP token ID returned by the mocked
// /deviceManagement/depOnboardingSettings endpoint, used to exercise dep_onboarding_settings_id
// auto-resolution (resolve_id.go) when a test config omits it.
const DepOnboardingSettingsTestID = "30000000-0000-0000-0000-000000000003"

// intuneProvisioningClientAppID mirrors the constant in validate.go.
const intuneProvisioningClientAppID = "f1346770-5b25-470b-88bd-d5744ab7952c"

var mockState struct {
	sync.Mutex
	policies               map[string]map[string]any // policy ID -> top-level fields echoed from POST/PUT
	policySettings         map[string][]any          // policy ID -> settings array echoed from POST/PUT
	membershipTargets      map[string]string         // policy ID -> current device_security_group target, if set
	defaultVisionOSProfile map[string]string         // dep_onboarding_settings_id -> default policy ID, if set
}

func init() {
	resetMockState()
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
}

func resetMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	mockState.policies = make(map[string]map[string]any)
	mockState.policySettings = make(map[string][]any)
	mockState.membershipTargets = make(map[string]string)
	mockState.defaultVisionOSProfile = make(map[string]string)
}

// VisionOSDeviceEnrollmentPolicyMock provides mock responses for visionOS ADE enrollment policy operations.
type VisionOSDeviceEnrollmentPolicyMock struct{}

// RegisterMocks registers HTTP mock responders for visionOS ADE enrollment policy operations.
//
//nolint:gocyclo // Mock registration functions naturally have high complexity
func (m *VisionOSDeviceEnrollmentPolicyMock) RegisterMocks() {
	resetMockState()

	// Apple DEP token lookup - used by resolveDepOnboardingSettingsId when
	// dep_onboarding_settings_id is omitted from config.
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement/depOnboardingSettings",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/depOnboardingSettings",
				"value": []map[string]any{
					{
						"id":        DepOnboardingSettingsTestID,
						"tokenName": "Unit Test Apple ADE Token",
						"tokenType": "dep",
					},
				},
			})
		})

	// Read a single DEP token, expanding its default visionOS enrollment profile - GET
	// /deviceManagement/depOnboardingSettings/{id}?$expand=defaultVisionOSEnrollmentProfile. Used by
	// resolveIsDefaultPolicyAssignment on every Read and by validateRequest on demote attempts.
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/depOnboardingSettings/[0-9a-fA-F-]+`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			defaultPolicyID, hasDefault := mockState.defaultVisionOSProfile[id]
			mockState.Unlock()

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/depOnboardingSettings/$entity",
				"id":             id,
			}
			if hasDefault {
				response["defaultVisionOSEnrollmentProfile"] = map[string]any{
					"id": id + "_" + defaultPolicyID,
				}
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Set default visionOS enrollment profile - POST
	// /deviceManagement/depOnboardingSettings/{depId}/enrollmentProfiles/{enrollmentProfileId}/setDefaultProfile.
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/depOnboardingSettings/[0-9a-fA-F-]+/enrollmentProfiles/.+/setDefaultProfile$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			// [...  "depOnboardingSettings", "{depId}", "enrollmentProfiles", "{depId}_{policyId}", "setDefaultProfile"]
			depId := parts[len(parts)-4]
			enrollmentProfileId := parts[len(parts)-2]
			policyId := strings.TrimPrefix(enrollmentProfileId, depId+"_")

			mockState.Lock()
			mockState.defaultVisionOSProfile[depId] = policyId
			mockState.Unlock()

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Security group owners - used by validateSecurityGroupOwnership. Every group is mocked as
	// owned by the Intune Provisioning Client.
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/[0-9a-fA-F-]+/owners$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#directoryObjects",
				"value": []map[string]any{
					{
						"@odata.type":          "#microsoft.graph.servicePrincipal",
						"id":                   "50000000-0000-0000-0000-000000000005",
						"appId":                intuneProvisioningClientAppID,
						"displayName":          "Intune Provisioning Client",
						"servicePrincipalType": "Application",
					},
				},
			})
		})

	// Create policy - POST /deviceManagement/configurationPolicies
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies",
		func(req *http.Request) (*http.Response, error) {
			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			id := uuid.New().String()
			settings, _ := body["settings"].([]any)
			delete(body, "settings")

			response := map[string]any{
				"@odata.context":       "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies/$entity",
				"id":                   id,
				"createdDateTime":      "2024-01-01T00:00:00Z",
				"lastModifiedDateTime": "2024-01-01T00:00:00Z",
				"settingCount":         len(settings),
				"isAssigned":           false,
			}
			maps.Copy(response, body)

			mockState.Lock()
			mockState.policies[id] = response
			mockState.policySettings[id] = settings
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
				return httpmock.NewJsonResponse(400, map[string]any{
					"error": map[string]any{"code": "BadRequest", "message": "Resource not found"},
				})
			}

			return httpmock.NewJsonResponse(200, policy)
		})

	// Read policy settings - GET /deviceManagement/configurationPolicies/{id}/settings
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F-]+/settings$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-2]

			mockState.Lock()
			settings := mockState.policySettings[id]
			mockState.Unlock()

			if settings == nil {
				settings = []any{}
			}

			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/configurationPolicies('" + id + "')/settings",
				"value":          settings,
			})
		})

	// Update policy - PUT /deviceManagement/configurationPolicies('{id}') (raw request: the Graph
	// API does not allow PATCH on the 'settings' navigation property).
	httpmock.RegisterResponder("PUT", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies\('[0-9a-fA-F-]+'\)$`,
		func(req *http.Request) (*http.Response, error) {
			id := extractParenID(req.URL.Path)

			var body map[string]any
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			mockState.Lock()
			existing, exists := mockState.policies[id]
			if !exists {
				mockState.Unlock()
				return httpmock.NewJsonResponse(400, map[string]any{
					"error": map[string]any{"code": "BadRequest", "message": "Resource not found"},
				})
			}

			settings, _ := body["settings"].([]any)
			delete(body, "settings")

			// PUT is a full replacement, but id/createdDateTime are server-owned.
			updated := map[string]any{
				"id":                   id,
				"createdDateTime":      existing["createdDateTime"],
				"lastModifiedDateTime": "2024-01-02T00:00:00Z",
				"settingCount":         len(settings),
				"isAssigned":           existing["isAssigned"],
			}
			maps.Copy(updated, body)

			mockState.policies[id] = updated
			mockState.policySettings[id] = settings
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
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
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

			statuses := []map[string]any{}
			if hasTarget {
				statuses = append(statuses, map[string]any{
					"targetId":                  targetId,
					"targetValidationErrorCode": "unknown",
					"validationSucceeded":       true,
				})
			}

			return httpmock.NewJsonResponse(200, map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#microsoft.graph.enrollmentTimeDeviceMembershipTargetResult",
				"enrollmentTimeDeviceMembershipTargetValidationStatuses": statuses,
			})
		})

	// Delete policy - DELETE /deviceManagement/configurationPolicies/{id}
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/configurationPolicies/[0-9a-fA-F-]+$`,
		func(req *http.Request) (*http.Response, error) {
			parts := strings.Split(req.URL.Path, "/")
			id := parts[len(parts)-1]

			mockState.Lock()
			delete(mockState.policies, id)
			delete(mockState.policySettings, id)
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
func (m *VisionOSDeviceEnrollmentPolicyMock) CleanupMockState() {
	resetMockState()
}
