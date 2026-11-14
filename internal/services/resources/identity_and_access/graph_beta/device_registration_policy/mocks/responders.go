package mocks

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const policyURL = "https://graph.microsoft.com/beta/policies/deviceRegistrationPolicy"

var mockState struct {
	sync.Mutex
	policy  map[string]any
	puts    []map[string]any
	seedErr error
}

func init() {
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("identity_and_access_device_registration_policy", &DeviceRegistrationPolicyMock{})
}

// DeviceRegistrationPolicyMock provides stateful httpmock responders for the singleton device registration policy.
type DeviceRegistrationPolicyMock struct{}

var _ mocks.MockRegistrar = (*DeviceRegistrationPolicyMock)(nil)

// loadSeedPolicy loads the tenant policy that GET returns before any PUT.
func loadSeedPolicy() (map[string]any, error) {
	responsesPath := filepath.Join("tests", "responses", "validate_get", "get_device_registration_policy_success.json")
	jsonData, err := os.ReadFile(responsesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load JSON response file: %w", err)
	}

	var policy map[string]any
	if err := json.Unmarshal(jsonData, &policy); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}
	return policy, nil
}

var errNullList = errors.New("must not be null")

func errorResponse(status int, code, message string) (*http.Response, error) {
	resp, err := httpmock.NewJsonResponse(status, map[string]any{"error": map[string]any{"code": code, "message": message}})
	if err != nil {
		return nil, fmt.Errorf("build error response: %w", err)
	}
	return resp, nil
}

// normaliseMembership mirrors Graph, which returns empty users and groups lists for an enumerated
// membership that omitted them, and rejects null lists.
func normaliseMembership(membership map[string]any) error {
	if membership == nil || membership["@odata.type"] != "#microsoft.graph.enumeratedDeviceRegistrationMembership" {
		return nil
	}
	for _, field := range []string{"users", "groups"} {
		value, ok := membership[field]
		if ok && value == nil {
			return fmt.Errorf("%s %w", field, errNullList)
		}
		if !ok {
			membership[field] = []any{}
		}
	}
	return nil
}

func (m *DeviceRegistrationPolicyMock) RegisterMocks() {
	policy, err := loadSeedPolicy()
	mockState.Lock()
	mockState.policy = policy
	mockState.puts = nil
	mockState.seedErr = err
	mockState.Unlock()

	// Get device registration policy - GET /policies/deviceRegistrationPolicy
	httpmock.RegisterResponder("GET", policyURL, func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		if mockState.seedErr != nil {
			return errorResponse(500, "InternalServerError", mockState.seedErr.Error())
		}
		return httpmock.NewJsonResponse(200, mockState.policy)
	})

	// Replace device registration policy - PUT /policies/deviceRegistrationPolicy
	httpmock.RegisterResponder("PUT", policyURL, func(req *http.Request) (*http.Response, error) {
		content, err := io.ReadAll(req.Body)
		if err != nil {
			return errorResponse(400, "BadRequest", "Invalid request body")
		}
		var requestBody, sentBody map[string]any
		if err := json.Unmarshal(content, &requestBody); err != nil {
			return errorResponse(400, "BadRequest", "Invalid request body")
		}
		_ = json.Unmarshal(content, &sentBody)

		// Graph requires every updatable property on PUT.
		if requestBody["localAdminPassword"] == nil {
			return errorResponse(400, "invalid_request", "Local admin password policy is null")
		}
		for _, field := range []string{"userDeviceQuota", "multiFactorAuthConfiguration", "azureADRegistration", "azureADJoin"} {
			if requestBody[field] == nil {
				return errorResponse(400, "invalid_request", field+" is required")
			}
		}

		join, _ := requestBody["azureADJoin"].(map[string]any)
		registration, _ := requestBody["azureADRegistration"].(map[string]any)
		localAdmins, _ := join["localAdmins"].(map[string]any)
		for _, membership := range []any{registration["allowedToRegister"], join["allowedToJoin"], localAdmins["registeringUsers"]} {
			membershipMap, _ := membership.(map[string]any)
			if err := normaliseMembership(membershipMap); err != nil {
				return errorResponse(400, "BadRequest", err.Error())
			}
		}

		mockState.Lock()
		defer mockState.Unlock()

		mockState.puts = append(mockState.puts, sentBody)

		stored := map[string]any{}
		for key, value := range requestBody {
			stored[key] = value
		}
		for _, key := range []string{"@odata.context", "id", "displayName", "description"} {
			stored[key] = mockState.policy[key]
		}

		// Graph ignores azureADRegistration changes when Intune manages registration.
		if current, ok := mockState.policy["azureADRegistration"].(map[string]any); ok && current["isAdminConfigurable"] == false {
			stored["azureADRegistration"] = current
		}

		mockState.policy = stored
		return httpmock.NewJsonResponse(200, stored)
	})
}

func (m *DeviceRegistrationPolicyMock) RegisterErrorMocks() {
	// Error scenarios for testing
	forbidden := `{"error":{"code":"Authorization_RequestDenied","message":"Insufficient privileges to complete the operation."}}`
	httpmock.RegisterResponder("GET", policyURL, httpmock.NewStringResponder(403, forbidden))
	httpmock.RegisterResponder("PUT", policyURL, httpmock.NewStringResponder(400, `{"error":{"code":"BadRequest","message":"Invalid request"}}`))
}

func (m *DeviceRegistrationPolicyMock) CleanupMockState() {
	mockState.Lock()
	mockState.policy = nil
	mockState.puts = nil
	mockState.seedErr = nil
	mockState.Unlock()
}

// PutRequests returns the request bodies received by the PUT responder, oldest first.
func (m *DeviceRegistrationPolicyMock) PutRequests() []map[string]any {
	mockState.Lock()
	defer mockState.Unlock()
	return append([]map[string]any{}, mockState.puts...)
}
