package mocks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

const policyURL = "https://graph.microsoft.com/beta/policies/authenticationMethodsPolicy"

// methodURLPattern matches PATCH /policies/authenticationMethodsPolicy/authenticationMethodConfigurations/{id}.
var methodURLPattern = regexp.MustCompile(`^https://graph\.microsoft\.com/beta/policies/authenticationMethodsPolicy/authenticationMethodConfigurations/([A-Za-z0-9]+)$`)

// expectedODataTypes are the @odata.type values Graph uses for each method configuration ID the resource writes.
var expectedODataTypes = map[string]string{
	"Email":               "#microsoft.graph.emailAuthenticationMethodConfiguration",
	"TemporaryAccessPass": "#microsoft.graph.temporaryAccessPassAuthenticationMethodConfiguration",
	"Sms":                 "#microsoft.graph.smsAuthenticationMethodConfiguration",
	"SoftwareOath":        "#microsoft.graph.softwareOathAuthenticationMethodConfiguration",
	"HardwareOath":        "#microsoft.graph.hardwareOathAuthenticationMethodConfiguration",
}

// PatchRequest is a PATCH request body received by the mock, with the method configuration ID or "" for the policy.
type PatchRequest struct {
	MethodID string
	Body     map[string]any
}

var mockState struct {
	sync.Mutex
	policy  map[string]any
	patches []PatchRequest
	seedErr error
}

func init() {
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("identity_and_access_authentication_methods_policy", &AuthenticationMethodsPolicyMock{})
}

// AuthenticationMethodsPolicyMock provides stateful httpmock responders for the singleton authentication methods policy.
type AuthenticationMethodsPolicyMock struct{}

var _ mocks.MockRegistrar = (*AuthenticationMethodsPolicyMock)(nil)

// loadSeedPolicy loads the tenant policy that GET returns before any PATCH.
func loadSeedPolicy() (map[string]any, error) {
	responsesPath := filepath.Join("tests", "responses", "validate_get", "get_authentication_methods_policy_success.json")
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

func errorResponse(status int, code, message string) (*http.Response, error) {
	resp, err := httpmock.NewJsonResponse(status, map[string]any{"error": map[string]any{"code": code, "message": message}})
	if err != nil {
		return nil, fmt.Errorf("build error response: %w", err)
	}
	return resp, nil
}

func readBody(req *http.Request) (map[string]any, error) {
	content, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("read request body: %w", err)
	}
	var body map[string]any
	if err := json.Unmarshal(content, &body); err != nil {
		return nil, fmt.Errorf("parse request body: %w", err)
	}
	return body, nil
}

func (m *AuthenticationMethodsPolicyMock) RegisterMocks() {
	policy, err := loadSeedPolicy()
	mockState.Lock()
	mockState.policy = policy
	mockState.patches = nil
	mockState.seedErr = err
	mockState.Unlock()

	// Get authentication methods policy - GET /policies/authenticationMethodsPolicy
	httpmock.RegisterResponder("GET", policyURL, func(req *http.Request) (*http.Response, error) {
		mockState.Lock()
		defer mockState.Unlock()

		if mockState.seedErr != nil {
			return errorResponse(500, "InternalServerError", mockState.seedErr.Error())
		}
		return httpmock.NewJsonResponse(200, mockState.policy)
	})

	// Update authentication methods policy - PATCH /policies/authenticationMethodsPolicy
	httpmock.RegisterResponder("PATCH", policyURL, func(req *http.Request) (*http.Response, error) {
		body, err := readBody(req)
		if err != nil {
			return errorResponse(400, "BadRequest", "Invalid request body")
		}

		mockState.Lock()
		defer mockState.Unlock()

		mockState.patches = append(mockState.patches, PatchRequest{Body: body})
		for key, value := range body {
			if key == "@odata.type" {
				continue
			}
			mockState.policy[key] = value
		}
		return httpmock.NewStringResponse(204, ""), nil
	})

	// Update method configuration - PATCH /policies/authenticationMethodsPolicy/authenticationMethodConfigurations/{id}
	httpmock.RegisterRegexpResponder("PATCH", methodURLPattern, func(req *http.Request) (*http.Response, error) {
		methodID := methodURLPattern.FindStringSubmatch(req.URL.String())[1]
		body, err := readBody(req)
		if err != nil {
			return errorResponse(400, "BadRequest", "Invalid request body")
		}

		// Graph rejects null target collections.
		for _, field := range []string{"includeTargets", "excludeTargets"} {
			if value, ok := body[field]; ok && value == nil {
				return errorResponse(400, "badRequest", "Persistence of policy failed with error: Invalid input, IncludeTarget is null")
			}
		}

		mockState.Lock()
		defer mockState.Unlock()

		configurations, _ := mockState.policy["authenticationMethodConfigurations"].([]any)
		for _, item := range configurations {
			configuration, _ := item.(map[string]any)
			if configuration["id"] != methodID {
				continue
			}
			if expected, ok := expectedODataTypes[methodID]; !ok || body["@odata.type"] != expected {
				message := fmt.Sprintf("Persistence of policy failed with error: The configured %s authentication method configuration is invalid.", methodID)
				return errorResponse(400, "badRequest", message)
			}
			mockState.patches = append(mockState.patches, PatchRequest{MethodID: methodID, Body: body})
			for key, value := range body {
				configuration[key] = value
			}
			return httpmock.NewStringResponse(204, ""), nil
		}

		return errorResponse(404, "resourceNotFound", fmt.Sprintf("Policy configuration for authentication method %s not found", methodID))
	})
}

func (m *AuthenticationMethodsPolicyMock) RegisterErrorMocks() {
	// Error scenarios for testing
	forbidden := `{"error":{"code":"Authorization_RequestDenied","message":"Insufficient privileges to complete the operation."}}`
	badRequest := `{"error":{"code":"badRequest","message":"Persistence of policy failed with error: The configured Email authentication method configuration is invalid."}}`
	httpmock.RegisterResponder("GET", policyURL, httpmock.NewStringResponder(403, forbidden))
	httpmock.RegisterResponder("PATCH", policyURL, httpmock.NewStringResponder(400, badRequest))
	httpmock.RegisterRegexpResponder("PATCH", methodURLPattern, httpmock.NewStringResponder(400, badRequest))
}

// RegisterNotFoundMocks registers a GET responder that returns 404 for the policy.
func (m *AuthenticationMethodsPolicyMock) RegisterNotFoundMocks() {
	httpmock.RegisterResponder("GET", policyURL, httpmock.NewStringResponder(404, `{"error":{"code":"resourceNotFound","message":"Resource not found"}}`))
}

func (m *AuthenticationMethodsPolicyMock) CleanupMockState() {
	mockState.Lock()
	mockState.policy = nil
	mockState.patches = nil
	mockState.seedErr = nil
	mockState.Unlock()
}

// PatchRequests returns the PATCH request bodies received by the mock, oldest first.
func (m *AuthenticationMethodsPolicyMock) PatchRequests() []PatchRequest {
	mockState.Lock()
	defer mockState.Unlock()
	return append([]PatchRequest{}, mockState.patches...)
}
