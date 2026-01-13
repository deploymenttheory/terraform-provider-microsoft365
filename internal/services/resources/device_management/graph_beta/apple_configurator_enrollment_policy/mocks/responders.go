package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	enrollmentProfiles map[string]map[string]any
	depOnboardingId    string
}

func init() {
	mockState.enrollmentProfiles = make(map[string]map[string]any)
	mockState.depOnboardingId = "54fac284-7866-43e5-860a-9c8e10fa3d7d"
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("apple_configurator_enrollment_policy", &AppleConfiguratorEnrollmentPolicyMock{})
}

type AppleConfiguratorEnrollmentPolicyMock struct{}

var _ mocks.MockRegistrar = (*AppleConfiguratorEnrollmentPolicyMock)(nil)

func (m *AppleConfiguratorEnrollmentPolicyMock) RegisterMocks() {
	mockState.Lock()
	mockState.enrollmentProfiles = make(map[string]map[string]any)
	mockState.Unlock()

	// 1. Get device management - used to resolve depOnboardingSettingsId
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/device_management_get.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)
		return httpmock.NewJsonResponse(200, responseObj)
	})

	// 2. Create enrollment profile - POST /deviceManagement/depOnboardingSettings/{depId}/enrollmentProfiles
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/depOnboardingSettings/[^/]+/enrollmentProfiles$`, func(req *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		id := mockState.depOnboardingId + "_" + strings.Replace(uuid.New().String(), "-", "", -1)[:24]
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/enrollment_profile_post.json")
		var responseObj map[string]any
		if err := json.Unmarshal([]byte(jsonStr), &responseObj); err != nil {
			responseObj = make(map[string]any)
		}

		// Copy all values from request body
		for k, v := range body {
			responseObj[k] = v
		}
		responseObj["id"] = id

		// Store in mock state
		mockState.Lock()
		if mockState.enrollmentProfiles == nil {
			mockState.enrollmentProfiles = make(map[string]map[string]any)
		}
		mockState.enrollmentProfiles[id] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// 3. Read enrollment profile - GET /deviceManagement/depOnboardingSettings/{depId}/enrollmentProfiles/{id}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/depOnboardingSettings/[^/]+/enrollmentProfiles/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]

		mockState.Lock()
		profile, exists := mockState.enrollmentProfiles[id]
		mockState.Unlock()

		if !exists {
			return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Enrollment profile not found"}}`), nil
		}

		// Load base response
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/enrollment_profile_get.json")
		var responseObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &responseObj)

		// Override with stored profile data
		for k, v := range profile {
			responseObj[k] = v
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// 4. Update enrollment profile - PATCH /deviceManagement/depOnboardingSettings/{depId}/enrollmentProfiles/{id}
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/depOnboardingSettings/[^/]+/enrollmentProfiles/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]

		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		}

		mockState.Lock()
		existing, exists := mockState.enrollmentProfiles[id]
		if !exists {
			mockState.Unlock()
			return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"Enrollment profile not found"}}`), nil
		}

		// Update existing profile
		for k, v := range body {
			existing[k] = v
		}
		mockState.enrollmentProfiles[id] = existing
		mockState.Unlock()

		return httpmock.NewJsonResponse(200, existing)
	})

	// 5. Delete enrollment profile - DELETE /deviceManagement/depOnboardingSettings/{depId}/enrollmentProfiles/{id}
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph\.microsoft\.com/beta/deviceManagement/depOnboardingSettings/.+/enrollmentProfiles/.+$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		id := parts[len(parts)-1]

		mockState.Lock()
		delete(mockState.enrollmentProfiles, id)
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})

	// 6. List enrollment profiles - GET /deviceManagement/depOnboardingSettings/{depId}/enrollmentProfiles
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/depOnboardingSettings/[^/]+/enrollmentProfiles$`, func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		responseObj := map[string]any{
			"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/depOnboardingSettings('54fac284-7866-43e5-860a-9c8e10fa3d7d')/enrollmentProfiles",
			"value":          []any{},
		}

		if mockState.enrollmentProfiles != nil && len(mockState.enrollmentProfiles) > 0 {
			list := make([]map[string]any, 0, len(mockState.enrollmentProfiles))
			for _, profile := range mockState.enrollmentProfiles {
				list = append(list, profile)
			}
			responseObj["value"] = list
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})
}

func (m *AppleConfiguratorEnrollmentPolicyMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.enrollmentProfiles = make(map[string]map[string]any)
	mockState.Unlock()

	// Make device management call fail
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/beta/deviceManagement", func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/error_500.json")
		var errObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(500, errObj)
	})

	// Make enrollment profile creation fail
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/depOnboardingSettings/[^/]+/enrollmentProfiles$`, func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/error_500.json")
		var errObj map[string]any
		_ = json.Unmarshal([]byte(jsonStr), &errObj)
		return httpmock.NewJsonResponse(500, errObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/depOnboardingSettings/[^/]+/enrollmentProfiles/[^/]+$`, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(404, `{"error":{"code":"ResourceNotFound","message":"Enrollment profile not found"}}`), nil
	})
}

func (m *AppleConfiguratorEnrollmentPolicyMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()
	for id := range mockState.enrollmentProfiles {
		delete(mockState.enrollmentProfiles, id)
	}
}
